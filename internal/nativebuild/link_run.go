package nativebuild

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"koda/internal/kodahome"
)

func objFileName() string {
	if runtime.GOOS == "windows" {
		return "main.obj"
	}
	return "main.o"
}

func runLLC(llcPath, irPath, objPath string, optFlag string, debug bool) error {
	args := []string{optFlag}
	if debug {
		args = append(args, "-g")
	}
	args = append(args, "-filetype=obj", "-o", objPath, irPath)
	cmd := exec.Command(llcPath, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func linkWithLLDGNU(tc *kodahome.Toolchain, objPath, outAbs string) error {
	args := []string{
		"-flavor", "gnu",
		objPath,
		tc.RuntimeLib,
		"-o", outAbs,
		"-lc", "-lm", "-lpthread",
	}
	cmd := exec.Command(tc.LLD, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func linkWithLLDDarwin(tc *kodahome.Toolchain, objPath, outAbs string) error {
	args := []string{
		"-flavor", "darwin",
		objPath,
		tc.RuntimeLib,
		"-lSystem",
		"-lm",
		"-o", outAbs,
	}
	cmd := exec.Command(tc.LLD, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

// runCompileAndLink compiles LLVM IR and links the Koda runtime (and system libs) in one clang invocation.
func runCompileAndLink(tc *kodahome.Toolchain, irFile, outAbs, rootDir string, opts BuildOptions, log func(string)) error {
	cc := tc.Clang
	if strings.TrimSpace(cc) == "" {
		cc = ClangDriver()
	}
	runtimeInclude := filepath.Join(rootDir, "runtime", "src")
	inputFile := irFile

	// On Windows, direct clang compilation from LLVM IR can crash on some generated programs.
	// Prefer llc -> object for stability unless explicitly disabled.
	llcPath := strings.TrimSpace(tc.LLC)
	if !llcIsUsable(llcPath) {
		if p, err := exec.LookPath("llc"); err == nil {
			llcPath = p
		}
	}
	if runtime.GOOS == "windows" && strings.HasSuffix(strings.ToLower(irFile), ".ll") && llcIsUsable(llcPath) && !disableWindowsLLCFallback() {
		objPath := filepath.Join(filepath.Dir(irFile), objFileName())
		if log != nil {
			log("  compile mode: llc -> object (windows stability fallback)\n")
		}
		if err := runLLC(llcPath, irFile, objPath, opts.llcOptFlag(), opts.Debug); err != nil {
			return fmt.Errorf("llc failed in windows fallback: %w", err)
		}
		inputFile = objPath
	}

	linkArgs := []string{opts.llcOptFlag()}
	if opts.Debug {
		linkArgs = append(linkArgs, "-g")
	}
	if tc.LLD != "" {
		if runtime.GOOS == "windows" {
			// Clang treats "-fuse-ld=C:\..." as multiple tokens (drive colon); rely on LLVM lld on PATH.
			linkArgs = append(linkArgs, "-fuse-ld=lld")
		} else {
			linkArgs = append(linkArgs, "-fuse-ld="+tc.LLD)
		}
	} else if UseLLD() {
		linkArgs = append(linkArgs, "-fuse-ld=lld")
	}
	if res, err := kodahome.BundledClangResourceFlags(); err == nil {
		linkArgs = append(linkArgs, res...)
	}
	linkArgs = append(linkArgs, inputFile, "-I", runtimeInclude)
	// LLVM IR often carries a target triple; suppress noisy override warning from clang.
	linkArgs = append(linkArgs, "-Wno-override-module")
	nativeSrc := os.Getenv("KODA_NATIVE_SOURCES")
	if strings.TrimSpace(nativeSrc) == "" {
		nativeSrc = os.Getenv("KODA_NATIVE_SOURCES")
	}
	if nativeSources := strings.Fields(nativeSrc); len(nativeSources) > 0 {
		if log != nil {
			log(fmt.Sprintf("  native sources: %s\n", strings.Join(nativeSources, " ")))
		}
		linkArgs = append(linkArgs, nativeSources...)
	}
	linkArgs = append(linkArgs, tc.RuntimeLib)
	if inc, arch, ok := vendoredRaylibStatic(rootDir); ok {
		if log != nil {
			log(fmt.Sprintf("  vendored raylib: %s\n", arch))
		}
		linkArgs = append(linkArgs, "-I", inc, arch)
	}
	linkExtra := os.Getenv("KODA_LINKFLAGS")
	if strings.TrimSpace(linkExtra) == "" {
		linkExtra = os.Getenv("KODA_LINKFLAGS")
	}
	if extra := strings.Fields(linkExtra); len(extra) > 0 {
		if log != nil {
			log(fmt.Sprintf("  link flags: %s\n\n", strings.Join(extra, " ")))
		}
		linkArgs = append(linkArgs, extra...)
	}
	linkArgs = append(linkArgs, defaultSystemLinkFlags()...)
	if runtime.GOOS == "windows" {
		linkArgs = append(linkArgs, "-lmsvcrt")
	}
	linkArgs = append(linkArgs, "-o", outAbs)

	buildArgs := func(o BuildOptions) []string {
		args := make([]string, 0, 64)
		args = append(args, o.llcOptFlag())
		if o.Debug {
			args = append(args, "-g")
		}
		if tc.LLD != "" {
			if runtime.GOOS == "windows" {
				args = append(args, "-fuse-ld=lld")
			} else {
				args = append(args, "-fuse-ld="+tc.LLD)
			}
		} else if UseLLD() {
			args = append(args, "-fuse-ld=lld")
		}
		if res, err := kodahome.BundledClangResourceFlags(); err == nil {
			args = append(args, res...)
		}
		args = append(args, inputFile, "-I", runtimeInclude)
		if nativeSources := strings.Fields(os.Getenv("KODA_NATIVE_SOURCES")); len(nativeSources) > 0 {
			args = append(args, nativeSources...)
		}
		args = append(args, tc.RuntimeLib)
		if inc, arch, ok := vendoredRaylibStatic(rootDir); ok {
			args = append(args, "-I", inc, arch)
		}
		if extra := strings.Fields(os.Getenv("KODA_LINKFLAGS")); len(extra) > 0 {
			args = append(args, extra...)
		}
		args = append(args, defaultSystemLinkFlags()...)
		if runtime.GOOS == "windows" {
			args = append(args, "-lmsvcrt")
		}
		args = append(args, "-o", outAbs)
		return args
	}

	var runOnce func(o BuildOptions) error
	runOnce = func(o BuildOptions) error {
		cmd := exec.Command(cc, buildArgs(o)...)
		var buf bytes.Buffer
		cmd.Stdout = io.MultiWriter(os.Stdout, &buf)
		cmd.Stderr = io.MultiWriter(os.Stderr, &buf)
		if runtime.GOOS == "windows" && strings.TrimSpace(tc.LLD) != "" {
			dir := filepath.Dir(tc.LLD)
			cmd.Env = append(os.Environ(), "PATH="+dir+string(os.PathListSeparator)+os.Getenv("PATH"))
		}
		if err := cmd.Run(); err != nil {
			if runtime.GOOS == "windows" && !o.NoOpt && looksLikeWindowsClangOptimizerCrash(buf.String()) {
				if log != nil {
					log("  warning: clang -O2 crashed; retrying with --no-opt\n")
				}
				o2 := o
				o2.NoOpt = true
				return runOnce(o2)
			}
			return err
		}
		return nil
	}

	return runOnce(opts)
}

func looksLikeWindowsClangOptimizerCrash(out string) bool {
	// Typical signature from LLVM/Clang on Windows when -cc1 dies during Optimizer with an access violation.
	// We keep this conservative to avoid retrying on ordinary compile/link errors.
	if !strings.Contains(out, "PLEASE submit a bug report") {
		return false
	}
	if !strings.Contains(out, "Optimizer") {
		return false
	}
	return strings.Contains(out, "Exception Code: 0xC0000005") || strings.Contains(out, "0xC0000005")
}

func disableWindowsLLCFallback() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("KODA_DISABLE_WINDOWS_LLC_FALLBACK")))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

func llcIsUsable(llc string) bool {
	llc = strings.TrimSpace(llc)
	if llc == "" {
		return false
	}
	if filepath.IsAbs(llc) {
		if fi, err := os.Stat(llc); err == nil && !fi.IsDir() {
			return true
		}
		return false
	}
	_, err := exec.LookPath(llc)
	return err == nil
}
