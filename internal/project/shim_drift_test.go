package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckRaylibShimDetectsMissingSymbols(t *testing.T) {
	root := t.TempDir()
	shimDir := filepath.Join(root, "wrappers", "raylib_shim")
	if err := os.MkdirAll(shimDir, 0755); err != nil {
		t.Fatal(err)
	}
	oldShim := `// koda:extern initwindow koda_shim_InitWindow 3
let initwindow = 0;
// koda:extern drawtext koda_shim_DrawText 5
let drawtext = 0;
`
	if err := os.WriteFile(filepath.Join(shimDir, "raylib.koda"), []byte(oldShim), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(shimDir, "wrapper.c"), []byte("void koda_shim_InitWindow(void){}"), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := `{
  "name": "test",
  "entry": "src/main.koda",
  "native": {
    "sources": ["wrappers/raylib_shim/wrapper.c"],
    "graphics": true
  }
}`
	if err := os.WriteFile(filepath.Join(root, FileName), []byte(cfg), 0644); err != nil {
		t.Fatal(err)
	}

	report, err := CheckRaylibShim(root)
	if err != nil {
		t.Fatal(err)
	}
	if !report.Stale {
		t.Fatal("expected stale report for outdated shim")
	}
	if len(report.MissingSymbols) == 0 {
		t.Fatal("expected missing @game symbols")
	}
}

func TestRefreshRaylibShimIfStale(t *testing.T) {
	root := t.TempDir()
	shimDir := filepath.Join(root, "wrappers", "raylib_shim")
	if err := os.MkdirAll(shimDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(shimDir, "raylib.koda"), []byte("let initwindow = 0;\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(shimDir, "wrapper.c"), []byte("void x(void){}"), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := `{
  "name": "test",
  "entry": "src/main.koda",
  "native": {
    "sources": ["wrappers/raylib_shim/wrapper.c"],
    "graphics": true
  }
}`
	if err := os.WriteFile(filepath.Join(root, FileName), []byte(cfg), 0644); err != nil {
		t.Fatal(err)
	}

	after, refreshed, err := RefreshRaylibShimIfStale(root)
	if err != nil {
		t.Fatal(err)
	}
	if !refreshed {
		t.Fatal("expected refresh")
	}
	if after.Stale {
		t.Fatalf("expected fresh shim after refresh: %s", after.Detail())
	}
	missing := missingShimSymbols(filepath.Join(shimDir, "raylib.koda"))
	if len(missing) > 0 {
		t.Fatalf("still missing symbols after refresh: %v", missing)
	}
}
