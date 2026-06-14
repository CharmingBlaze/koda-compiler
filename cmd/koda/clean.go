package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// runClean removes common native build artifacts from the project tree.
func runClean(args []string) error {
	root := "."
	cache := false
	for _, a := range args {
		switch a {
		case "--cache":
			cache = true
		default:
			if strings.HasPrefix(a, "-") {
				return fmt.Errorf("unknown flag: %s", a)
			}
			root = a
		}
	}
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	for _, name := range []string{".koda_build", ".KODA_build", "dist"} {
		p := filepath.Join(absRoot, name)
		if err := removeTreeIfExists(p); err != nil {
			return err
		}
	}
	for _, rel := range []string{"main.koda", "src/main.koda", "tests/smoke_native.koda"} {
		kodaPath := filepath.Join(absRoot, rel)
		if _, err := os.Stat(kodaPath); err != nil {
			continue
		}
		for _, exe := range []string{defaultExeName(kodaPath), strings.TrimSuffix(kodaPath, ".koda") + ".exe"} {
			if err := removeFileIfExists(exe); err != nil {
				return err
			}
		}
	}
	if cache {
		if err := cleanTempCaches(); err != nil {
			fmt.Fprintf(os.Stderr, "clean cache: %v\n", err)
		}
	}
	fmt.Println("clean OK")
	return nil
}

func cleanTempCaches() error {
	tmp := os.TempDir()
	entries, err := os.ReadDir(tmp)
	if err != nil {
		return err
	}
	for _, e := range entries {
		name := e.Name()
		if strings.HasPrefix(name, "koda-") || strings.HasPrefix(name, "koda_run_") || strings.HasPrefix(name, "koda_watch_") || strings.HasPrefix(name, "koda_eval_") {
			p := filepath.Join(tmp, name)
			_ = os.RemoveAll(p)
		}
	}
	return nil
}

func removeTreeIfExists(path string) error {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if !fi.IsDir() {
		return os.Remove(path)
	}
	return os.RemoveAll(path)
}

func removeFileIfExists(path string) error {
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
