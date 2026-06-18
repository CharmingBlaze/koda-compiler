package nativebuild

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNativeSourceCacheKeyChangesWithMtime(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "wrapper.c")
	if err := os.WriteFile(src, []byte("int x;"), 0644); err != nil {
		t.Fatal(err)
	}
	k1, err := nativeSourceCacheKey(src)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(20 * time.Millisecond)
	if err := os.WriteFile(src, []byte("int x; int y;"), 0644); err != nil {
		t.Fatal(err)
	}
	k2, err := nativeSourceCacheKey(src)
	if err != nil {
		t.Fatal(err)
	}
	if k1 == k2 {
		t.Fatalf("cache key should change when source changes: %q", k1)
	}
}

func TestNativeObjectCacheValid(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.c")
	obj := filepath.Join(dir, "a.obj")
	if err := os.WriteFile(src, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if nativeObjectCacheValid(obj, src) {
		t.Fatal("missing obj should not be valid")
	}
	if err := os.WriteFile(obj, []byte("obj"), 0644); err != nil {
		t.Fatal(err)
	}
	if !nativeObjectCacheValid(obj, src) {
		t.Fatal("obj newer than source should be valid")
	}
}
