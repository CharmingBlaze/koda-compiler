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

func TestApplyNativeEnvClearsStaleNativeSources(t *testing.T) {
	tmp := t.TempDir()
	cfg := `{"name":"x","version":"0.1.0","entry":"src/main.koda"}`
	if err := os.WriteFile(filepath.Join(tmp, FileName), []byte(cfg), 0644); err != nil {
		t.Fatal(err)
	}
	cfgObj, root, err := Load(tmp)
	if err != nil {
		t.Fatal(err)
	}
	ctx := &Context{Root: root, Cfg: cfgObj}
	t.Setenv("KODA_NATIVE_SOURCES", filepath.Join(tmp, "stale", "wrapper.c"))
	if err := ctx.ApplyNativeEnv(); err != nil {
		t.Fatal(err)
	}
	if v := os.Getenv("KODA_NATIVE_SOURCES"); v != "" {
		t.Fatalf("expected unset KODA_NATIVE_SOURCES, got %q", v)
	}
}

func TestApplyNativeEnvGraphicsDefaultsShim(t *testing.T) {
	tmp := t.TempDir()
	shim := filepath.Join(tmp, "wrappers", "raylib_shim")
	if err := os.MkdirAll(shim, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(shim, "wrapper.c"), []byte("int x;"), 0644); err != nil {
		t.Fatal(err)
	}
	cfg := `{"name":"x","version":"0.1.0","entry":"src/main.koda","native":{"graphics":true}}`
	if err := os.WriteFile(filepath.Join(tmp, FileName), []byte(cfg), 0644); err != nil {
		t.Fatal(err)
	}
	cfgObj, root, err := Load(tmp)
	if err != nil {
		t.Fatal(err)
	}
	ctx := &Context{Root: root, Cfg: cfgObj}
	t.Setenv("KODA_NATIVE_SOURCES", filepath.Join(tmp, "wrong", "wrapper.c"))
	if err := ctx.ApplyNativeEnv(); err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(tmp, "wrappers", "raylib_shim", "wrapper.c")
	if os.Getenv("KODA_NATIVE_SOURCES") != want {
		t.Fatalf("sources = %q, want %q", os.Getenv("KODA_NATIVE_SOURCES"), want)
	}
}
