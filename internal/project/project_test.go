package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindAndResolveEntry(t *testing.T) {
	tmp := t.TempDir()
	root := filepath.Join(tmp, "myapp")
	if err := os.MkdirAll(filepath.Join(root, "src"), 0755); err != nil {
		t.Fatal(err)
	}
	cfg := `{"name":"myapp","version":"0.1.0","entry":"src/main.koda"}`
	if err := os.WriteFile(filepath.Join(root, FileName), []byte(cfg), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "src", "main.koda"), []byte("print(1);"), 0644); err != nil {
		t.Fatal(err)
	}

	got, err := ResolveEntry(filepath.Join(root, "src"), "")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(root, "src", "main.koda")
	if got != want {
		t.Fatalf("entry = %q want %q", got, want)
	}

	ctx, err := LoadContext(got)
	if err != nil {
		t.Fatal(err)
	}
	if ctx == nil || ctx.Cfg.Name != "myapp" {
		t.Fatalf("context: %+v", ctx)
	}
	if ctx.AppName(got) != "myapp" {
		t.Fatalf("app name = %q", ctx.AppName(got))
	}
}

func TestApplyNativeEnv(t *testing.T) {
	tmp := t.TempDir()
	cfg := `{
		"name":"g","entry":"main.koda",
		"native":{"sources":["glue/wrapper.c"],"linkflags":"-lfoo"}
	}`
	if err := os.WriteFile(filepath.Join(tmp, FileName), []byte(cfg), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tmp, "glue"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "glue", "wrapper.c"), []byte("int x;"), 0644); err != nil {
		t.Fatal(err)
	}

	ctx, err := LoadContext(tmp)
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("KODA_NATIVE_SOURCES", "")
	os.Setenv("KODA_LINKFLAGS", "")
	if err := ctx.ApplyNativeEnv(); err != nil {
		t.Fatal(err)
	}
	if os.Getenv("KODA_NATIVE_SOURCES") == "" {
		t.Fatal("expected KODA_NATIVE_SOURCES")
	}
	if os.Getenv("KODA_LINKFLAGS") != "-lfoo" {
		t.Fatalf("linkflags = %q", os.Getenv("KODA_LINKFLAGS"))
	}
}
