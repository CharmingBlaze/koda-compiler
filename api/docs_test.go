package api

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListDocPagesFromSDK(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "docs", "learn"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "START_HERE.md"), []byte("# Start Here\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "faq.md"), []byte("# FAQ\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "learn", "01-welcome.md"), []byte("# Welcome\n"), 0644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("KODA_HOME", root)

	pages, err := ListDocPages()
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) < 3 {
		t.Fatalf("expected at least 3 pages, got %d", len(pages))
	}
	body, err := ReadDocPage("docs/faq.md")
	if err != nil {
		t.Fatal(err)
	}
	if body == "" {
		t.Fatal("expected faq content")
	}
}

func TestReadDocPageRejectsTraversal(t *testing.T) {
	root := t.TempDir()
	t.Setenv("KODA_HOME", root)
	if _, err := ReadDocPage("../secret.md"); err == nil {
		t.Fatal("expected error for traversal path")
	}
}
