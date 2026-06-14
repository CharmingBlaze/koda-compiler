package kodahome

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	fujembed "koda/internal/embed"
)

// Resolution order for Clang (and similarly LLC):
//  1. KODA_CLANG / KODA_LLC — explicit override (CI, custom SDKs)
//  2. CC (Clang only) — common third-party convention; treated like an explicit override
//  3. Bundled next to the executable: toolchain/bin then llvm/bin (after first-run extract)
//  4. Embedded release binaries (temp dir; internal/embed) when the install tree is not present yet
//  5. PATH lookup
//  6. Bare name ("clang" / "llc") so exec errors are obvious if nothing resolves
//
// Preferring (3) before (4) keeps kodawrap’s Clang AST step on the full unpacked toolchain next to
// koda / kodawrap so release users do not need a host LLVM installation.

// ClangWithSource returns the same executable as [Clang] and a short label for diagnostics.
func ClangWithSource() (path, source string) {
	if s := strings.TrimSpace(os.Getenv("KODA_CLANG")); s != "" {
		return s, "KODA_CLANG"
	}
	if cc := strings.TrimSpace(os.Getenv("CC")); cc != "" {
		return cc, "CC"
	}
	dir, err := InstallDir()
	if err == nil {
		for _, rel := range BundledClangRelPaths {
			p := filepath.Join(dir, rel)
			if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
				return p, "bundled"
			}
		}
	}
	if p, err := fujembed.ClangPath(); err == nil {
		return p, "embedded"
	} else if err != nil && !errors.Is(err, fujembed.ErrDevelopmentBuild) {
		// Keep falling back so local/dev toolchains still work if embedded extraction is unavailable.
	}
	if p, err := exec.LookPath("clang"); err == nil {
		return p, "PATH"
	}
	return "clang", "default"
}

// LLCWithSource returns the same executable as [LLC] and a short label for diagnostics.
func LLCWithSource() (path, source string) {
	if s := strings.TrimSpace(os.Getenv("KODA_LLC")); s != "" {
		return s, "KODA_LLC"
	}
	dir, err := InstallDir()
	if err == nil {
		for _, rel := range BundledLLCRelPaths {
			p := filepath.Join(dir, rel)
			if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
				return p, "bundled"
			}
		}
	}
	if p, err := fujembed.LLCPath(); err == nil {
		return p, "embedded"
	} else if err != nil && !errors.Is(err, fujembed.ErrDevelopmentBuild) {
		// Keep falling back so local/dev toolchains still work if embedded extraction is unavailable.
	}
	if p, err := exec.LookPath("llc"); err == nil {
		return p, "PATH"
	}
	return "llc", "default"
}

// BundledLLDPath returns the absolute path to a bundled ld.lld next to the install,
// preferring toolchain/bin over llvm/bin (same roots as [HasBundledLLD]).
func BundledLLDPath() (path string, ok bool) {
	dir, err := InstallDir()
	if err != nil {
		return "", false
	}
	for _, root := range []string{"toolchain", "llvm"} {
		p := filepath.Join(dir, root, "bin", lldExeName())
		fi, err := os.Stat(p)
		if err == nil && !fi.IsDir() {
			return p, true
		}
	}
	return "", false
}

// InstallDirWarnings returns non-fatal hints about the install location (e.g. macOS Gatekeeper).
func InstallDirWarnings(install string) []string {
	if install == "" {
		return nil
	}
	var out []string
	// Heuristic paths are macOS-specific, but we evaluate on every OS so tests and
	// diagnostics behave the same in CI (Linux/Windows) as on darwin.
	if strings.Contains(install, "AppTranslocation") ||
		strings.Contains(install, "/private/var/folders/") {
		out = append(out, "Executable appears to run from a Gatekeeper translocation path (often read-only). "+
			"Move koda into /Applications or a normal folder so stdlib/ and toolchain/ beside the binary are writable and discoverable.")
	}
	return out
}
