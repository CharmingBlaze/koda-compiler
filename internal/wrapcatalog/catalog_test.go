package wrapcatalog

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseInstallSpec(t *testing.T) {
	spec, err := ParseInstallSpec("raylib@5.0")
	if err != nil {
		t.Fatal(err)
	}
	if spec.Name != "raylib" || spec.Version != "5.0" {
		t.Fatalf("%+v", spec)
	}
	spec, err = ParseInstallSpec("@sqlite3")
	if err != nil || spec.Name != "sqlite3" {
		t.Fatalf("%+v %v", spec, err)
	}
}

func TestCatalogFindRaylib(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	lib, ver, err := c.Find("raylib", "")
	if err != nil {
		t.Fatal(err)
	}
	if lib.Name != "raylib" || ver.Prebuilt == "" {
		t.Fatalf("%+v %+v", lib, ver)
	}
}

func TestExpandSDK(t *testing.T) {
	ctx := ExpandContext{SDK: `C:\Koda`, ProjectRoot: `D:\proj`}
	got := ctx.Expand("{sdk}/wrappers/raylib")
	want := filepath.FromSlash(`C:\Koda/wrappers/raylib`)
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestCopyDirRoundTrip(t *testing.T) {
	src := t.TempDir()
	dst := t.TempDir()
	if err := os.WriteFile(filepath.Join(src, "a.txt"), []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	sub := filepath.Join(src, "sub")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "b.txt"), []byte("y"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := CopyDir(src, dst); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dst, "sub", "b.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "y" {
		t.Fatalf("got %q", data)
	}
}

func TestResolveHeaderMissing(t *testing.T) {
	_, err := ResolveHeader([]string{filepath.Join(t.TempDir(), "nope.h")})
	if err == nil {
		t.Fatal("expected error")
	}
}
