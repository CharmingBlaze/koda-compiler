package api

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), ".."))
}

func chdirRepoRoot(t *testing.T) string {
	t.Helper()
	root := repoRoot(t)
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })
	return root
}

func TestRunArrayGetOutOfBoundsPanics(t *testing.T) {
	root := chdirRepoRoot(t)
	t.Setenv("KODA_HOME", root)
	p := filepath.Join(root, "tests", "array_oob_get.koda")
	var out, errBuf bytes.Buffer
	err := RunWithWritersOpts(p, "", &out, &errBuf, BuildOptions{NoOpt: true})
	if err == nil {
		t.Fatal("expected runtime exit error from out-of-bounds array read")
	}
	combined := out.String() + errBuf.String()
	if !strings.Contains(combined, "out of bounds") {
		t.Fatalf("expected panic text containing 'out of bounds', got:\n%s", combined)
	}
}

func TestRunArraySetOutOfBoundsPanics(t *testing.T) {
	root := chdirRepoRoot(t)
	t.Setenv("KODA_HOME", root)
	p := filepath.Join(root, "tests", "array_oob_set.koda")
	var out, errBuf bytes.Buffer
	err := RunWithWritersOpts(p, "", &out, &errBuf, BuildOptions{NoOpt: true})
	if err == nil {
		t.Fatal("expected runtime exit error from out-of-bounds array assign")
	}
	combined := out.String() + errBuf.String()
	if !strings.Contains(combined, "out of bounds") {
		t.Fatalf("expected panic text containing 'out of bounds', got:\n%s", combined)
	}
}
