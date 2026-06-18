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

// logWriter adapts a log func to io.Writer (for forwarding clang stderr to IDEs).
type logWriter struct {
	log func(string)
}

func (w logWriter) Write(p []byte) (int, error) {
	if w.log != nil && len(p) > 0 {
		w.log(string(p))
	}
	return len(p), nil
}

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
func runCompileAndLink(tc *kodahome.Toolchain, irFile, outAbs, sdkRoot, projectRoot string, opts BuildOptions, log func(string)) error {
	cc := tc.Clang
	if strings.TrimSpace(cc) == "" {
		cc = ClangDriver()
	}
	runtimeInclude := filepath.Join(sdkRoot, "runtime", "src")
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

	nativeSrc := os.Getenv("KODA_NATIVE_SOURCES")
	if strings.TrimSpace(nativeSrc) == "" {
		nativeSrc = os.Getenv("KODA_NATIVE_SOURCES")
	}
	var nativeLinkObjects []string
	if nativeSources := splitNativeSources(nativeSrc); len(nativeSources) > 0 {
		if log != nil {
			log(fmt.Sprintf("  native sources: %s\n", strings.Join(nativeSources, " ")))
		}
		var err error
		nativeLinkObjects, err = materializeNativeObjects(cc, nativeSources, projectRoot, sdkRoot, opts, log)
		if err != nil {
			return err
		}
	}
	if inc, arch, ok := vendoredRaylibStatic(projectRoot); ok {
		if log != nil {
			log(fmt.Sprintf("  vendored raylib: %s\n", arch))
		}
		_ = inc
	}
	linkExtra := os.Getenv("KODA_LINKFLAGS")
	if strings.TrimSpace(linkExtra) == "" {
		linkExtra = os.Getenv("KODA_LINKFLAGS")
	}
	if extra := strings.Fields(linkExtra); len(extra) > 0 {
		if _, _, vendored := vendoredRaylibStatic(projectRoot); vendored {
			extra = omitLinkFlag(extra, "-lraylib")
		}
		if log != nil {
			log(fmt.Sprintf("  link flags: %s\n\n", strings.Join(extra, " ")))
		}
	}

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
		args = append(args, inputFile, "-I", runtimeInclude, "-Wno-override-module")
		if len(nativeLinkObjects) > 0 {
			args = append(args, nativeLinkObjects...)
		}
		args = append(args, tc.RuntimeLib)
		if inc, arch, ok := vendoredRaylibStatic(projectRoot); ok {
			args = append(args, "-I", inc, arch)
		}
		if extra := strings.Fields(linkExtra); len(extra) > 0 {
			if _, _, vendored := vendoredRaylibStatic(projectRoot); vendored {
				extra = omitLinkFlag(extra, "-lraylib")
			}
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
		var stderr io.Writer = &buf
		if log != nil {
			stderr = io.MultiWriter(&buf, logWriter{log})
		}
		cmd.Stdout = stderr
		cmd.Stderr = stderr
		if runtime.GOOS == "windows" && strings.TrimSpace(tc.LLD) != "" {
			dir := filepath.Dir(tc.LLD)
			cmd.Env = append(os.Environ(), "PATH="+dir+string(os.PathListSeparator)+os.Getenv("PATH"))
		}
		if err := cmd.Run(); err != nil {
			detail := strings.TrimSpace(buf.String())
			if runtime.GOOS == "windows" && !o.NoOpt && (looksLikeWindowsClangOptimizerCrash(detail) || looksLikeWindowsClangExit74(detail, err)) {
				if log != nil {
					log("  warning: clang failed; retrying with --no-opt\n")
				}
				o2 := o
				o2.NoOpt = true
				return runOnce(o2)
			}
			if detail != "" {
				return fmt.Errorf("%s: %w", detail, err)
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

func looksLikeWindowsClangExit74(detail string, err error) bool {
	if err == nil {
		return false
	}
	if !strings.Contains(err.Error(), "exit status 74") {
		return false
	}
	// Empty or unhelpful clang output — often a flaky -O2 / llc path on Windows.
	return detail == "" || strings.Contains(detail, "PLEASE submit a bug report")
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
