package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	fujembed "koda/internal/embed"
	"koda/internal/kodahome"
	"koda/internal/project"
	"koda/internal/wrappermeta"
)

type doctorLine struct {
	ok   bool
	label string
	detail string
	fix   string
}

// runDoctor prints a beginner-friendly SDK health report.
func runDoctor(args []string) error {
	fix := false
	for _, a := range args {
		if a == "--fix" {
			fix = true
		}
	}
	var lines []doctorLine
	var failures int

	lines = append(lines, doctorLine{ok: true, label: "koda", detail: strings.TrimSpace(version)})

	install, err := kodahome.InstallDir()
	if err != nil {
		lines = append(lines, doctorLine{
			ok: false, label: "install dir", detail: err.Error(),
			fix: "Set KODA_HOME to your SDK folder or run koda from an unpacked release zip.",
		})
		failures++
	} else {
		lines = append(lines, doctorLine{ok: true, label: "SDK path", detail: install})
		if ok, why := kodahome.InstallDirWritable(install); ok {
			lines = append(lines, doctorLine{ok: true, label: "writable install", detail: "ok"})
		} else {
			lines = append(lines, doctorLine{
				ok: false, label: "writable install", detail: why,
				fix: "Move koda.exe to a normal folder (not Downloads with Gatekeeper translocation).",
			})
			failures++
		}
	}

	stdlib, err := kodahome.StdlibDir()
	if err != nil {
		lines = append(lines, doctorLine{ok: false, label: "stdlib", detail: err.Error(), fix: "Unpack the SDK zip so stdlib/ sits next to koda.exe."})
		failures++
	} else if fi, err := os.Stat(filepath.Join(stdlib, "math.koda")); err != nil || fi.IsDir() {
		lines = append(lines, doctorLine{ok: false, label: "stdlib", detail: stdlib, fix: "Missing math.koda — reinstall the SDK bundle."})
		failures++
	} else {
		lines = append(lines, doctorLine{ok: true, label: "stdlib", detail: stdlib})
	}

	tc, tcErr := kodahome.FindToolchain()
	if tcErr != nil {
		lines = append(lines, doctorLine{ok: false, label: "runtime archive", detail: tcErr.Error(), fix: "Run scripts/build-runtime.ps1 or make -C runtime from a source checkout."})
		failures++
	} else if fi, err := os.Stat(tc.RuntimeLib); err != nil || fi.IsDir() {
		lines = append(lines, doctorLine{ok: false, label: "runtime archive", detail: tc.RuntimeLib, fix: "Rebuild runtime/libkoda_runtime.a"})
		failures++
	} else {
		lines = append(lines, doctorLine{ok: true, label: "runtime archive", detail: tc.RuntimeLib})
		if exe, err := os.Executable(); err == nil {
			if exeFi, err := os.Stat(exe); err == nil {
				if exeFi.ModTime().After(fi.ModTime()) {
					lines = append(lines, doctorLine{
						ok: false, label: "runtime freshness", detail: "koda.exe is newer than libkoda_runtime.a",
						fix: "Run scripts/build-runtime.ps1 or scripts/build-runtime.sh",
					})
					failures++
				} else {
					lines = append(lines, doctorLine{ok: true, label: "runtime freshness", detail: "archive up to date"})
				}
			}
		}
	}

	clangPath, _ := kodahome.ClangWithSource()
	if probeOK(clangPath) {
		lines = append(lines, doctorLine{ok: true, label: "clang", detail: shortVersion(clangPath)})
	} else {
		lines = append(lines, doctorLine{ok: false, label: "clang", detail: clangPath, fix: "Install LLVM/clang or use a release build with an embedded toolchain."})
		failures++
	}

	llcPath, _ := kodahome.LLCWithSource()
	if probeOK(llcPath) {
		lines = append(lines, doctorLine{ok: true, label: "llc", detail: shortVersion(llcPath)})
	} else {
		lines = append(lines, doctorLine{ok: false, label: "llc", detail: llcPath, fix: "Install LLVM (llc) or use a release SDK zip."})
		failures++
	}

	if p, ok := kodahome.BundledLLDPath(); ok && fileOK(p) {
		lines = append(lines, doctorLine{ok: true, label: "lld", detail: p})
	} else if runtime.GOOS == "windows" {
		lines = append(lines, doctorLine{ok: false, label: "lld", detail: "not bundled", fix: "Use a release build with lld.exe next to koda, or install LLVM."})
		failures++
	}

	raylibOK, raylibDetail := detectRaylib()
	if raylibOK {
		lines = append(lines, doctorLine{ok: true, label: "raylib", detail: raylibDetail})
	} else {
		lines = append(lines, doctorLine{
			ok: false, label: "raylib", detail: raylibDetail,
			fix: "For graphics: build third_party/raylib_static (make raylib-lib) or set KODA_LINKFLAGS manually. Console games work without raylib.",
		})
	}

	if stale, err := doctorWrapperDrift(); err == nil {
		for _, r := range stale {
			lines = append(lines, doctorLine{
				ok: false, label: "wrapper " + r.Meta.Import, detail: "native headers changed since last koda wrap",
				fix: fmt.Sprintf("koda wrap upgrade %s", r.Dir),
			})
			failures++
		}
	}

	if shimReport, err := doctorRaylibShim(fix); err == nil && shimReport.ShimDir != "" {
		if shimReport.Stale {
			lines = append(lines, doctorLine{
				ok: false, label: "raylib shim", detail: shimReport.Detail(),
				fix: "koda setup raylib   (or: koda doctor --fix)",
			})
			failures++
		} else {
			detail := shimReport.Detail()
			if fix {
				detail = "up to date (refreshed if needed)"
			}
			lines = append(lines, doctorLine{ok: true, label: "raylib shim", detail: detail})
		}
	}

	tmpDir := os.TempDir()
	if f, err := os.CreateTemp(tmpDir, "koda_doctor_*"); err != nil {
		lines = append(lines, doctorLine{
			ok: false, label: "temp folder write", detail: err.Error(),
			fix: "Check permissions on " + tmpDir + " or set TMP/TEMP to a writable folder.",
		})
		failures++
	} else {
		_ = f.Close()
		_ = os.Remove(f.Name())
		lines = append(lines, doctorLine{ok: true, label: "temp folder write", detail: tmpDir})
	}

	if free, ok := diskFreeBytes(tmpDir); ok {
		if free < 64*1024*1024 {
			lines = append(lines, doctorLine{
				ok: false, label: "disk space", detail: fmt.Sprintf("%.0f MB free on temp drive", float64(free)/(1024*1024)),
				fix: "Free at least 64 MB on your temp drive for builds.",
			})
			failures++
		} else {
			lines = append(lines, doctorLine{ok: true, label: "disk space", detail: fmt.Sprintf("%.0f MB free", float64(free)/(1024*1024))})
		}
	}

	if exe, err := os.Executable(); err == nil {
		if err := doctorSmokeBuild(exe, "hello"); err != nil {
			lines = append(lines, doctorLine{ok: false, label: "sample hello build", detail: err.Error(), fix: "Fix clang/runtime issues above, then retry."})
			failures++
		} else {
			lines = append(lines, doctorLine{ok: true, label: "sample hello build", detail: "ok"})
		}
	}

	fmt.Println("Koda Doctor")
	fmt.Println()
	for _, ln := range lines {
		tag := "OK  "
		if !ln.ok {
			tag = "FAIL"
		}
		fmt.Printf("%s  %s", tag, ln.label)
		if ln.detail != "" {
			fmt.Printf(" — %s", ln.detail)
		}
		fmt.Println()
		if !ln.ok && ln.fix != "" {
			fmt.Printf("      Fix: %s\n", ln.fix)
		}
	}
	fmt.Println()
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	if _, err := fujembed.Extract(); err == nil {
		fmt.Println("Toolchain: release (embedded clang/runtime)")
	} else if errors.Is(err, fujembed.ErrDevelopmentBuild) {
		fmt.Println("Toolchain: development (system clang)")
	}
	fmt.Printf("Graphics link hint: %s\n", project.DefaultGraphicsLinkFlags())
	fmt.Println()
	if failures > 0 {
		fmt.Printf("Result: %d issue(s) found. Fix FAIL lines before building games.\n", failures)
		return fmt.Errorf("%d doctor check(s) failed", failures)
	}
	fmt.Println("Result: all checks passed.")
	return nil
}

func probeOK(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path, "--version")
	return cmd.Run() == nil
}

func shortVersion(path string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, path, "--version").Output()
	if err != nil {
		return path
	}
	line := strings.TrimSpace(string(out))
	if idx := strings.IndexByte(line, '\n'); idx >= 0 {
		line = line[:idx]
	}
	if len(line) > 80 {
		line = line[:80] + "..."
	}
	return line
}

func fileOK(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}

func detectRaylib() (bool, string) {
	if stage, _, ok := nativeRaylibStage(); ok {
		return true, stage
	}
	if os.Getenv("KODA_LINKFLAGS") != "" {
		return true, "KODA_LINKFLAGS set in environment"
	}
	return false, "not found (optional for console programs)"
}

func nativeRaylibStage() (stage string, archive string, ok bool) {
	inst, err := kodahome.InstallDir()
	if err != nil {
		return "", "", false
	}
	candidates := []string{
		filepath.Join(inst, "third_party", "raylib_static", "stage"),
		filepath.Join(filepath.Dir(inst), "third_party", "raylib_static", "stage"),
	}
	for _, c := range candidates {
		lib := filepath.Join(c, "lib", "libraylib.a")
		if runtime.GOOS == "windows" {
			lib = filepath.Join(c, "lib", "raylib.lib")
		}
		if fileOK(lib) {
			return c, lib, true
		}
	}
	return "", "", false
}

func doctorSmokeBuild(kodaExe, kind string) error {
	tmp, err := os.MkdirTemp("", "koda_doctor_build_*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	src := filepath.Join(tmp, "hello.koda")
	if err := os.WriteFile(src, []byte(`func main() { print("doctor"); }`), 0o644); err != nil {
		return err
	}
	out := filepath.Join(tmp, "hello")
	if runtime.GOOS == "windows" {
		out += ".exe"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, kodaExe, "build", src, "-o", out)
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	_ = kind
	return nil
}

func doctorWrapperDrift() ([]wrappermeta.DriftReport, error) {
	ctx, err := project.LoadContext(cwd())
	if err != nil || ctx == nil || ctx.Cfg == nil {
		return nil, err
	}
	var abs []string
	for _, src := range ctx.Cfg.Native.Sources {
		src = strings.TrimSpace(src)
		if src == "" {
			continue
		}
		p := filepath.Join(ctx.Root, filepath.FromSlash(src))
		abs = append(abs, p)
	}
	if inst, err := kodahome.InstallDir(); err == nil {
		for _, src := range ctx.Cfg.Native.Sources {
			src = strings.TrimSpace(src)
			if src == "" {
				continue
			}
			p := filepath.Join(inst, filepath.FromSlash(src))
			if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
				abs = append(abs, p)
			}
		}
	}
	return wrappermeta.CheckProjectSources(abs)
}

func doctorRaylibShim(fix bool) (project.RaylibShimReport, error) {
	root := cwd()
	if fix {
		after, refreshed, err := project.RefreshRaylibShimIfStale(root)
		if err != nil {
			return after, err
		}
		if refreshed {
			return after, nil
		}
	}
	return project.CheckRaylibShim(root)
}
