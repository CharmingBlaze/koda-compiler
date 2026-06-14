package nativebuild

import (
	"os"
	"strings"

	"koda/internal/kodahome"
)

// ClangDriver returns the C/LLVM driver used to compile IR and the embedded runtime.
// See [kodahome.Clang] for resolution order (env, bundled next to koda, PATH).
func ClangDriver() string {
	return kodahome.Clang()
}

// UseLLD returns true when the link step should pass -fuse-ld=lld to Clang.
// It is enabled if KODA_USE_LLD is truthy, or if a bundled ld.lld sits next to Clang
// (see [kodahome.HasBundledLLD]). Set KODA_USE_LLD=0 to force the platform default linker.
func UseLLD() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("KODA_USE_LLD")))
	if v == "0" || v == "false" || v == "off" || v == "no" {
		return false
	}
	if v == "1" || v == "true" || v == "yes" || v == "on" {
		return true
	}
	return kodahome.HasBundledLLD()
}
