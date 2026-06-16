package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"syscall"
	"time"

	"koda/internal/nativebuild"
)

type watchFingerprint struct {
	size    int64
	modUnix int64
}

func parseWatchCommandArgs(args []string) (src string, noOpt bool, debug bool, progArgs []string, err error) {
	before, after := splitAtDoubleDash(args)
	progArgs = after
	var file string
	for i := 0; i < len(before); i++ {
		switch before[i] {
		case "--no-opt":
			noOpt = true
		case "--debug":
			debug = true
		default:
			if strings.HasPrefix(before[i], "-") {
				return "", false, false, nil, fmt.Errorf("unknown flag: %s", before[i])
			}
			if file != "" {
				return "", false, false, nil, fmt.Errorf("multiple source files")
			}
			file = before[i]
		}
	}
	if file == "" {
		return "", noOpt, debug, progArgs, nil
	}
	return file, noOpt, debug, progArgs, nil
}

func runWatch(path string, opts nativebuild.BuildOptions, progArgs []string) error {
	absEntry, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	watchRoot := filepath.Dir(absEntry)
	fmt.Printf("watching %s\n", watchRoot)
	fmt.Printf("entry    %s\n", absEntry)

	snap, err := collectKodaSnapshot(watchRoot)
	if err != nil {
		return err
	}

	proc, procDone, exePath, err := buildAndStartWatched(absEntry, opts, progArgs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "watch: initial build failed: %v\n", err)
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	sigc := make(chan os.Signal, 2)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigc)

	for {
		select {
		case <-ticker.C:
			next, e := collectKodaSnapshot(watchRoot)
			if e != nil {
				fmt.Fprintf(os.Stderr, "watch: scan error: %v\n", e)
				continue
			}
			if snapshotsEqual(snap, next) {
				continue
			}
			changes := snapshotChanges(snap, next, watchRoot)
			snap = next
			for _, ch := range changes {
				fmt.Printf("\n[koda watch] %s — rebuilding...\n", ch)
			}
			if len(changes) == 0 {
				fmt.Println("\nwatch: change detected, rebuilding...")
			}
			if proc != nil {
				_ = proc.Kill()
				_, _ = <-procDone
				proc = nil
				procDone = nil
			}
			if exePath != "" {
				_ = os.Remove(exePath)
				exePath = ""
			}
			proc, procDone, exePath, err = buildAndStartWatched(absEntry, opts, progArgs)
			if err != nil {
				fmt.Fprintf(os.Stderr, "watch: build failed: %v\n", err)
			}
		case err = <-procDone:
			if err != nil {
				fmt.Fprintf(os.Stderr, "watch: program exited: %v\n", err)
			} else {
				fmt.Fprintln(os.Stderr, "watch: program exited")
			}
			proc = nil
			procDone = nil
		case <-sigc:
			if proc != nil {
				_ = proc.Kill()
				if procDone != nil {
					_, _ = <-procDone
				}
			}
			if exePath != "" {
				_ = os.Remove(exePath)
			}
			return nil
		}
	}
}

func buildAndStartWatched(absEntry string, opts nativebuild.BuildOptions, progArgs []string) (*os.Process, <-chan error, string, error) {
	tmpExe, err := tempExecutablePathWatch()
	if err != nil {
		return nil, nil, "", err
	}
	log := func(s string) { _, _ = fmt.Fprint(os.Stdout, s) }
	if err := nativebuild.BuildWithOverlaysOpts(absEntry, nil, tmpExe, log, opts); err != nil {
		_ = os.Remove(tmpExe)
		return nil, nil, "", err
	}

	cmd := exec.Command(tmpExe, progArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Start(); err != nil {
		_ = os.Remove(tmpExe)
		return nil, nil, "", err
	}
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	return cmd.Process, done, tmpExe, nil
}

func tempExecutablePathWatch() (string, error) {
	pattern := "koda_watch_*"
	if runtime.GOOS == "windows" {
		pattern = "koda_watch_*.exe"
	}
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}
	name := f.Name()
	_ = f.Close()
	if err := os.Remove(name); err != nil {
		return "", err
	}
	return name, nil
}

func collectKodaSnapshot(root string) (map[string]watchFingerprint, error) {
	out := make(map[string]watchFingerprint)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == ".KODA_build" || name == "bin" || name == "dist" || name == "node_modules" || name == "vendor" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.EqualFold(filepath.Ext(path), ".koda") {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		out[path] = watchFingerprint{
			size:    info.Size(),
			modUnix: info.ModTime().UnixNano(),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

func snapshotsEqual(a, b map[string]watchFingerprint) bool {
	if len(a) != len(b) {
		return false
	}
	if len(a) == 0 {
		return true
	}
	keys := make([]string, 0, len(a))
	for k := range a {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		av, ok := a[k]
		if !ok {
			return false
		}
		bv, ok := b[k]
		if !ok {
			return false
		}
		if av != bv {
			return false
		}
	}
	return true
}

func snapshotChanges(a, b map[string]watchFingerprint, root string) []string {
	var out []string
	seen := make(map[string]bool)
	for path, av := range a {
		bv, ok := b[path]
		if !ok {
			rel, _ := filepath.Rel(root, path)
			if rel == "" {
				rel = path
			}
			out = append(out, rel+" removed")
			seen[rel] = true
			continue
		}
		if av != bv {
			rel, _ := filepath.Rel(root, path)
			if rel == "" {
				rel = path
			}
			out = append(out, rel+" changed")
			seen[rel] = true
		}
	}
	for path := range b {
		if _, ok := a[path]; !ok {
			rel, _ := filepath.Rel(root, path)
			if rel == "" {
				rel = path
			}
			if !seen[rel] {
				out = append(out, rel+" added")
			}
		}
	}
	sort.Strings(out)
	return out
}
