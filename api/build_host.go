package api

import (
	"path/filepath"

	"koda/internal/nativebuild"
)

// BuildNativeHost compiles entryPath to a native executable (same pipeline as koda build):
// Go codegen emits LLVM IR, llc lowers it to an object file, then Clang links that object with
// runtime/libkoda_runtime.a and headers under runtime/src (plus optional KODA_NATIVE_SOURCES / KODA_LINKFLAGS).
func BuildNativeHost(entryPath, overlay, output string, log func(string)) error {
	EnsureSDKFromExecutable()
	return withProjectScope(entryPath, func() error {
		absEntry, err := filepath.Abs(entryPath)
		if err != nil {
			return err
		}
		overlays := map[string]string{}
		if overlay != "" {
			overlays[absEntry] = overlay
		}
		return nativebuild.BuildWithOverlays(entryPath, overlays, output, log)
	})
}

// DefaultExeName returns the default native executable name for a .koda path.
func DefaultExeName(source string) string {
	return nativebuild.DefaultExeName(source)
}
