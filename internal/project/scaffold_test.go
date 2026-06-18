package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestListTemplates(t *testing.T) {
	names := ListTemplates()
	want := []string{"hello", "game", "graphics", "pong", "raylib"}
	for _, w := range want {
		found := false
		for _, n := range names {
			if n == w {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("missing template %q in %v", w, names)
		}
	}
}

func TestScaffoldHello(t *testing.T) {
	tmp := t.TempDir()
	dir, err := Scaffold(tmp, "helloapp", "hello")
	if err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{"koda.json", "src/main.koda", "README.md"} {
		if _, err := os.Stat(filepath.Join(dir, filepath.FromSlash(rel))); err != nil {
			t.Fatalf("scaffold hello missing %s: %v", rel, err)
		}
	}
}

func TestScaffoldPong(t *testing.T) {
	tmp := t.TempDir()
	dir, err := Scaffold(tmp, "pongapp", "pong")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "src", "main.koda")); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(dir, "src", "main.koda"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "game.delta() or 0.016") {
		t.Fatal("pong template should include delta fallback")
	}
}
