package kodahome

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSDKRootCandidatesWalksUp(t *testing.T) {
	root := t.TempDir()
	leaf := filepath.Join(root, "a", "b", "c")
	if err := os.MkdirAll(leaf, 0755); err != nil {
		t.Fatal(err)
	}
	cands := SDKRootCandidates(leaf)
	if len(cands) < 3 {
		t.Fatalf("expected at least 3 candidates, got %d: %v", len(cands), cands)
	}
	if cands[0] != leaf {
		t.Fatalf("first candidate should be leaf dir, got %q", cands[0])
	}
	foundRoot := false
	for _, c := range cands {
		if c == root {
			foundRoot = true
			break
		}
	}
	if !foundRoot {
		t.Fatalf("expected temp root in candidates: %v", cands)
	}
}

func TestBootstrapSDKRootFindsParentStdlib(t *testing.T) {
	t.Setenv("KODA_HOME", "")
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "stdlib"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "stdlib", "math.koda"), []byte("// test\n"), 0644); err != nil {
		t.Fatal(err)
	}
	exeDir := filepath.Join(root, "koda-ide", "build", "bin")
	if err := os.MkdirAll(exeDir, 0755); err != nil {
		t.Fatal(err)
	}
	got, ok := BootstrapSDKRoot(exeDir)
	if !ok {
		t.Fatal("BootstrapSDKRoot returned false")
	}
	if got != root {
		t.Fatalf("got SDK root %q, want %q", got, root)
	}
	if os.Getenv("KODA_HOME") != root {
		t.Fatalf("KODA_HOME=%q, want %q", os.Getenv("KODA_HOME"), root)
	}
}
