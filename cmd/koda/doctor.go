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
)

// runDoctor prints a diagnostic report for zero-install / native builds.
func runDoctor() error {
	fmt.Println("=== koda doctor ===")
	fmt.Println()

	releaseToolchain := false
	if embedDir, err := fujembed.Extract(); err == nil {
		releaseToolchain = true
		fmt.Printf("Release build — embedded toolchain at %s\n", embedDir)
		fmt.Println("You only need koda + kodawrap (and stdlib next to them). No Go install and no LLVM install on your machine.")
		clang, _ := fujembed.ClangPath()
		lib, _ := fujembed.RuntimeLibPath()
		fmt.Printf("  clang:   %s\n", clang)
		fmt.Printf("  runtime: %s\n", lib)
		if runtime.GOOS == "windows" {
			if lld, err := fujembed.LLDPathWindows(); err == nil {
				fmt.Printf("  lld:     %s\n", lld)
			}
		}
	} else if errors.Is(err, fujembed.ErrDevelopmentBuild) {
		fmt.Println("Development build — system toolchain")

		tc, err := kodahome.FindToolchain()
		if err != nil {
			fmt.Printf("x toolchain: %v\n", err)
		} else {
			fmt.Printf("ok clang:   %s\n", tc.Clang)
			fmt.Printf("ok runtime: %s\n", tc.RuntimeLib)
			fmt.Printf("  llc:     %s\n", tc.LLC)
			if tc.LLD != "" {
				fmt.Printf("  lld:     %s\n", tc.LLD)
			}
		}

		llcPath, llcSrc := kodahome.LLCWithSource()
		fmt.Printf("llc resolution: %s (from %s)\n", llcPath, llcSrc)
		printToolProbe("llc", llcPath)
	} else {
		fmt.Printf("Embedded toolchain failed: %v\n", err)
	}

	fmt.Println()

	install, err := kodahome.InstallDir()
	if err != nil {
		return fmt.Errorf("install dir: %w", err)
	}
	fmt.Printf("install_dir: %s\n", install)

	if ok, why := kodahome.InstallDirWritable(install); ok {
		fmt.Println("install_writable: ok")
	} else {
		fmt.Printf("install_writable: NO (%s)\n", why)
		fmt.Println("  hint: move the binary to a normal writable folder; Gatekeeper translocation can block writes.")
	}
	for _, w := range kodahome.InstallDirWarnings(install) {
		fmt.Printf("install_warning: %s\n", w)
	}
	fmt.Println()

	clangPath, clangSrc := kodahome.ClangWithSource()
	fmt.Printf("clang (resolution): %s (from %s)\n", clangPath, clangSrc)
	printToolProbe("clang", clangPath)

	if p, ok := kodahome.BundledLLDPath(); ok {
		fmt.Printf("lld next to binary: %s\n", p)
		if fi, err := os.Stat(p); err != nil || fi.IsDir() {
			fmt.Println("lld_status: missing_or_invalid")
		} else {
			fmt.Println("lld_status: ok")
		}
	} else {
		fmt.Println("lld next to binary: (none)")
	}
	fmt.Println()

	stdlib, err := kodahome.StdlibDir()
	if err != nil {
		return err
	}
	fmt.Printf("stdlib_dir: %s\n", stdlib)
	printDirStatus("stdlib", stdlib)
	if fi, err := os.Stat(filepath.Join(stdlib, "math.koda")); err != nil || fi.IsDir() {
		fmt.Println("stdlib_hint: @math / @json imports need stdlib/ next to koda (or set KODA_PATH)")
	}

	wrap, err := kodahome.WrappersDir()
	if err != nil {
		return err
	}
	fmt.Printf("wrappers_dir: %s\n", wrap)
	printDirStatus("wrappers", wrap)

	flags, err := kodahome.BundledClangResourceFlags()
	if err != nil {
		fmt.Printf("bundled_clang_isystem: error: %v\n", err)
	} else {
		fmt.Printf("bundled_clang_isystem_entries: %d\n", len(flags))
		if len(flags) > 0 && os.Getenv("KODA_DOCTOR_VERBOSE") != "" {
			fmt.Println(strings.Join(flags, "\n"))
		}
	}

	fmt.Println()
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	if releaseToolchain {
		fmt.Println("koda is a normal native program; you do not install Go separately to use it.")
	} else {
		fmt.Printf("Go version (this koda binary was built with): %s\n", runtime.Version())
	}
	fmt.Println()
	fmt.Println("Doctor finished.")
	return nil
}

func printDirStatus(label, path string) {
	fi, err := os.Stat(path)
	if err != nil {
		fmt.Printf("%s_exists: 0 (%v)\n", label, err)
		return
	}
	if !fi.IsDir() {
		fmt.Printf("%s_exists: 0 (not a directory)\n", label)
		return
	}
	fmt.Printf("%s_exists: 1\n", label)
}

func printToolProbe(name, path string) {
	if path == "" {
		fmt.Printf("%s_probe: skip (empty path)\n", name)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path, "--version")
	cmd.Stderr = nil
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("%s_probe: failed (%v)\n", name, err)
		return
	}
	line := strings.TrimSpace(string(out))
	if idx := strings.IndexByte(line, '\n'); idx >= 0 {
		line = line[:idx]
	}
	if len(line) > 120 {
		line = line[:120] + "..."
	}
	fmt.Printf("%s_probe: %s\n", name, line)
}
