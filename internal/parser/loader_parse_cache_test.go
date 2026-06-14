package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseCacheUsesMtime(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "one.koda")
	if err := os.WriteFile(p, []byte("let a = 1;\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	resetParseCache()
	b1, err := LoadProgram(p)
	if err != nil {
		t.Fatalf("first load: %v", err)
	}
	decl0 := firstTopLevelLetName(t, b1, p)
	if decl0 != "a" {
		t.Fatalf("want first binding a, got %q", decl0)
	}

	b2, err := LoadProgram(p)
	if err != nil {
		t.Fatalf("second load: %v", err)
	}
	if firstTopLevelLetName(t, b2, p) != "a" {
		t.Fatal("second load should still see a")
	}

	if err := os.WriteFile(p, []byte("let b = 99;\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	b3, err := LoadProgram(p)
	if err != nil {
		t.Fatalf("after edit: %v", err)
	}
	if got := firstTopLevelLetName(t, b3, p); got != "b" {
		t.Fatalf("after mtime change want binding b, got %q", got)
	}
}

func firstTopLevelLetName(t *testing.T, bundle *ProgramBundle, modulePath string) string {
	t.Helper()
	modPath, err := filepath.Abs(modulePath)
	if err != nil {
		t.Fatal(err)
	}
	prog := bundle.Modules[modPath]
	if prog == nil {
		t.Fatalf("missing module %q (have %d modules)", modPath, len(bundle.Modules))
	}
	if len(prog.Declarations) == 0 {
		t.Fatal("empty program")
	}
	ld, ok := prog.Declarations[0].(*LetDecl)
	if !ok {
		t.Fatalf("want let, got %T", prog.Declarations[0])
	}
	return ld.Name.Lexeme
}
