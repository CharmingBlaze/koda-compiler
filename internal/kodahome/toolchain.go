package kodahome

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// LinkMode selects how the native backend links the object file with the runtime.
type LinkMode int

const (
	// LinkClang uses the Clang driver (same resolution as [Clang]) with the usual flags.
	LinkClang LinkMode = iota
	// LinkLLDGNU uses LLVM's LLD in GNU flavor (Linux embedded releases).
	LinkLLDGNU
	// LinkLLDDarwin uses LLVM's LLD in Darwin / ld64-compatible flavor (macOS embedded releases).
	LinkLLDDarwin
)

// Toolchain holds absolute paths to LLVM tools and the Koda static archive used by [nativebuild].
type Toolchain struct {
	LLC        string
	LLD        string // used when LinkMode == LinkLLDGNU; may be empty for LinkClang
	Clang      string // used when LinkMode == LinkClang
	RuntimeLib string
	LinkMode   LinkMode
}

// FindToolchain resolves LLVM + runtime for native builds.
// With "-tags release" and a populated internal/embed/<GOOS>/<GOARCH>/ tree, Clang + runtime are
// extracted once to a temp directory and returned.
// Otherwise the toolchain comes from the environment and PATH (see resolveDevClang, [LLCWithSource]).
func FindToolchain() (*Toolchain, error) {
	tc, err := embeddedToolchain()
	if err != nil {
		// Release-tag builds can be compiled without populated embedded assets in dev checkouts.
		// Fall back to system toolchain so local runs stay usable/noiseless.
		if errors.Is(err, ErrIncompleteEmbeddedToolchain) {
			return findSystemToolchain()
		}
		return nil, err
	}
	if tc != nil {
		return tc, nil
	}
	return findSystemToolchain()
}

func findSystemToolchain() (*Toolchain, error) {
	root, err := filepath.Abs(".")
	if err != nil {
		return nil, fmt.Errorf("project root: %w", err)
	}
	clang, err := resolveDevClang()
	if err != nil {
		return nil, err
	}
	llc, _ := LLCWithSource()
	if strings.TrimSpace(llc) == "" {
		if p := findLLCBinary(); p != "" {
			llc = p
		}
	}
	if strings.TrimSpace(llc) == "" {
		if alt := llcSiblingOfClang(clang); alt != "" {
			llc = alt
		}
	}
	lld := findLLDBinary()
	return &Toolchain{
		LLC:        llc,
		LLD:        lld,
		Clang:      clang,
		RuntimeLib: filepath.Join(root, "runtime", "libkoda_runtime.a"),
		LinkMode:   LinkClang,
	}, nil
}

func llcSiblingOfClang(clang string) string {
	clang = strings.TrimSpace(clang)
	if clang == "" {
		return ""
	}
	dir := filepath.Dir(clang)
	if dir == "." || dir == "" {
		return ""
	}
	name := "llc"
	if runtime.GOOS == "windows" {
		name = "llc.exe"
	}
	p := filepath.Join(dir, name)
	if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
		return p
	}
	return ""
}

func findLLCBinary() string {
	candidates := []string{"llc-18", "llc-14", "llc"}
	if runtime.GOOS == "windows" {
		candidates = []string{
			`C:\Program Files\LLVM\bin\llc.exe`,
			"llc.exe",
			"llc-18.exe",
			"llc-14.exe",
		}
	}
	for _, c := range candidates {
		if filepath.IsAbs(c) {
			if fi, err := os.Stat(c); err == nil && !fi.IsDir() {
				return c
			}
			continue
		}
		if p, err := exec.LookPath(c); err == nil {
			return p
		}
	}
	return ""
}

func resolveDevClang() (string, error) {
	if p := strings.TrimSpace(os.Getenv("KODA_CLANG")); p != "" {
		return p, nil
	}
	if cc := strings.TrimSpace(os.Getenv("CC")); cc != "" {
		return cc, nil
	}
	if runtime.GOOS == "windows" {
		if root, err := filepath.Abs("."); err == nil {
			gnuShim := filepath.Join(root, "scripts", "clang-gnu.cmd")
			if fi, err := os.Stat(gnuShim); err == nil && !fi.IsDir() {
				return gnuShim, nil
			}
		}
	}
	dir, err := InstallDir()
	if err == nil {
		for _, rel := range BundledClangRelPaths {
			p := filepath.Join(dir, rel)
			if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
				return p, nil
			}
		}
	}
	return findClangOnPathOrAbsolute(clangCandidates())
}

func clangCandidates() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{"clang", "clang-18", "clang-14"}
	case "darwin":
		return []string{
			"/opt/homebrew/opt/llvm/bin/clang",
			"/opt/homebrew/opt/llvm@18/bin/clang",
			"/opt/homebrew/opt/llvm@14/bin/clang",
			"/usr/local/opt/llvm/bin/clang",
			"clang",
		}
	default:
		return []string{"clang-18", "clang-14", "clang"}
	}
}

func findClangOnPathOrAbsolute(candidates []string) (string, error) {
	for _, name := range candidates {
		if filepath.IsAbs(name) {
			if fi, err := os.Stat(name); err == nil && !fi.IsDir() {
				return name, nil
			}
			continue
		}
		if path, err := exec.LookPath(name); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf(
		"clang not found for native builds.\n\n"+
			"If you use a GitHub release of koda/kodawrap: keep the executable in the SDK folder "+
			"(stdlib/ and toolchain/ beside it), run from a writable directory, and do not set "+
			"KODA_SKIP_TOOLCHAIN_EXTRACT unless you also set KODA_CLANG to a full Clang.\n\n"+
			"If you build Koda from source without -tags release, install a system Clang:\n"+
			"  Windows: choco install llvm   (or set KODA_CLANG)\n"+
			"  macOS:   brew install llvm\n"+
			"  Linux:   apt-get install clang (or clang-18)",
	)
}

func findLLDBinary() string {
	candidates := []string{"ld.lld-14", "ld.lld-15", "ld.lld"}
	if runtime.GOOS == "windows" {
		candidates = []string{"ld.lld.exe", "lld.exe"}
	}
	for _, name := range candidates {
		if p, err := exec.LookPath(name); err == nil {
			return p
		}
	}
	return ""
}

// ErrIncompleteEmbeddedToolchain means a release build was compiled with -tags release
// but internal/embed/<GOOS>/<GOARCH>/ was incomplete at compile time (see internal/embed/README.md).
var ErrIncompleteEmbeddedToolchain = errors.New("incomplete embedded toolchain (see internal/embed/README.md)")
