package nativebuild

import (
	"os"
	"path/filepath"
	"testing"

	"koda/internal/parser"
)

func TestInferNativeGlueKindShim(t *testing.T) {
	path := filepath.Join("..", "..", "wrappers", "raylib_shim", "raylib.koda")
	bundle, err := parser.LoadProgram(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := inferNativeGlueKind(bundle); got != nativeGlueShim {
		t.Fatalf("inferNativeGlueKind = %v, want shim", got)
	}
}

func TestNativeSourcesMatchKindRejectsFullWrapperForShim(t *testing.T) {
	full := filepath.Join("C:", "Koda", "wrappers", "raylib", "wrapper.c")
	if nativeSourcesMatchKind(full, nativeGlueShim) {
		t.Fatal("full raylib wrapper should not match shim kind")
	}
	shim := filepath.Join("C:", "Koda", "wrappers", "raylib_shim", "wrapper.c")
	if !nativeSourcesMatchKind(shim, nativeGlueShim) {
		t.Fatal("shim wrapper should match shim kind")
	}
}

func TestApplyNativeSourcesForBundleOverridesStaleFullWrapper(t *testing.T) {
	root := t.TempDir()
	shimDir := filepath.Join(root, "wrappers", "raylib_shim")
	if err := os.MkdirAll(shimDir, 0755); err != nil {
		t.Fatal(err)
	}
	shimC := filepath.Join(shimDir, "wrapper.c")
	if err := os.WriteFile(shimC, []byte("void koda_shim_InitWindow(void){}"), 0644); err != nil {
		t.Fatal(err)
	}

	bundle, err := parser.LoadProgram(filepath.Join("..", "..", "wrappers", "raylib_shim", "raylib.koda"))
	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("KODA_NATIVE_SOURCES", filepath.Join("C:", "Koda", "wrappers", "raylib", "wrapper.c"))
	t.Setenv("KODA_LINKFLAGS", "")

	if err := ApplyNativeSourcesForBundle(bundle, root); err != nil {
		t.Fatal(err)
	}
	got := os.Getenv("KODA_NATIVE_SOURCES")
	if got != shimC {
		t.Fatalf("KODA_NATIVE_SOURCES = %q, want %q", got, shimC)
	}
	if os.Getenv("KODA_LINKFLAGS") == "" {
		t.Fatal("expected graphics link flags to be set")
	}
}
