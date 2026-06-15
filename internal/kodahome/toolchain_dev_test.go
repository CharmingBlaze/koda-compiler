package kodahome

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindSystemToolchainRuntimeUsesInstallDirNotCWD(t *testing.T) {
	root := t.TempDir()
	sdk := filepath.Join(root, "sdk")
	if err := os.MkdirAll(filepath.Join(sdk, "runtime"), 0755); err != nil {
		t.Fatal(err)
	}
	lib := filepath.Join(sdk, "runtime", "libkoda_runtime.a")
	if err := os.WriteFile(lib, []byte("fake"), 0644); err != nil {
		t.Fatal(err)
	}
	proj := filepath.Join(root, "my-project")
	if err := os.MkdirAll(proj, 0755); err != nil {
		t.Fatal(err)
	}

	t.Setenv("KODA_HOME", sdk)
	t.Setenv("KODA_CLANG", filepath.Join(sdk, "fake-clang"))
	if err := os.WriteFile(filepath.Join(sdk, "fake-clang"), []byte("@echo off\n"), 0755); err != nil {
		t.Fatal(err)
	}

	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(proj); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })

	tc, err := findSystemToolchain()
	if err != nil {
		t.Fatal(err)
	}
	if tc.RuntimeLib != lib {
		t.Fatalf("RuntimeLib=%q want %q", tc.RuntimeLib, lib)
	}
}
