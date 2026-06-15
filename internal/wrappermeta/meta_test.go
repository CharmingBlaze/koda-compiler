package wrappermeta

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHashAndDrift(t *testing.T) {
	dir := t.TempDir()
	header := filepath.Join(dir, "lib.h")
	if err := os.WriteFile(header, []byte("void foo();\n"), 0644); err != nil {
		t.Fatal(err)
	}
	hashes, err := HashHeaders([]string{header})
	if err != nil {
		t.Fatal(err)
	}
	meta := PackageMeta{
		Name:         "lib",
		Headers:      []string{header},
		HeaderHashes: hashes,
	}
	report, err := CheckDrift(meta)
	if err != nil {
		t.Fatal(err)
	}
	if report.Stale() {
		t.Fatalf("expected fresh, got %+v", report)
	}
	if err := os.WriteFile(header, []byte("void foo(int x);\n"), 0644); err != nil {
		t.Fatal(err)
	}
	report, err = CheckDrift(meta)
	if err != nil {
		t.Fatal(err)
	}
	if !report.Stale() || len(report.Changed) != 1 {
		t.Fatalf("expected changed header, got %+v", report)
	}
}

func TestParseLinkflags(t *testing.T) {
	inc, link := ParseLinkflags("-I/usr/include -L/usr/lib -lsqlite3")
	if len(inc) != 1 || inc[0] != "/usr/include" {
		t.Fatalf("inc: %v", inc)
	}
	if len(link) != 2 {
		t.Fatalf("link: %v", link)
	}
}
