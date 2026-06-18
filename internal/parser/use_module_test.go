package parser_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"koda/internal/kodahome"
	"koda/internal/lexer"
	"koda/internal/parser"
)

func TestParseUseDecl(t *testing.T) {
	src := `use raylib;
use koda.math;

func main() {
    print("ok");
}`
	l := lexer.NewLexer(src, "test.koda")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	p := parser.NewParser(tokens)
	prog, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	if len(prog.Declarations) < 2 {
		t.Fatalf("expected at least 2 decls, got %d", len(prog.Declarations))
	}
	u0, ok := prog.Declarations[0].(*parser.UseDecl)
	if !ok {
		t.Fatalf("decl[0]: want *UseDecl, got %T", prog.Declarations[0])
	}
	if u0.ModulePath != "raylib" {
		t.Fatalf("module path = %q, want raylib", u0.ModulePath)
	}
	u1, ok := prog.Declarations[1].(*parser.UseDecl)
	if !ok {
		t.Fatalf("decl[1]: want *UseDecl, got %T", prog.Declarations[1])
	}
	if u1.ModulePath != "koda.math" {
		t.Fatalf("module path = %q, want koda.math", u1.ModulePath)
	}
}

func TestResolveUseKodaMath(t *testing.T) {
	root := repoRoot(t)
	t.Setenv("KODA_HOME", root)
	kodahome.BootstrapSDKRoot(root)

	path, err := parser.ResolveUseModulePath(filepath.Join(root, "tests", "dummy.koda"), "koda.math")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(strings.ToLower(path), "math.koda") {
		t.Fatalf("resolved %q, want math.koda", path)
	}
}

func TestResolveUseRaylib(t *testing.T) {
	root := repoRoot(t)
	t.Setenv("KODA_HOME", root)
	kodahome.BootstrapSDKRoot(root)

	path, err := parser.ResolveUseModulePath(filepath.Join(root, "tests", "dummy.koda"), "raylib")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.ToLower(path), "raylib") {
		t.Fatalf("resolved %q, expected raylib wrapper path", path)
	}
}

func TestResolveUseUnknownModule(t *testing.T) {
	root := repoRoot(t)
	t.Setenv("KODA_HOME", root)
	kodahome.BootstrapSDKRoot(root)

	_, err := parser.ResolveUseModulePath(filepath.Join(root, "tests", "dummy.koda"), "raylib3d")
	if err == nil {
		t.Fatal("expected error for unknown module")
	}
	msg := err.Error()
	if !strings.Contains(msg, "Unknown module") {
		t.Fatalf("error = %q, want Unknown module", msg)
	}
	if !strings.Contains(msg, "Searched:") {
		t.Fatalf("error = %q, want Searched: list", msg)
	}
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "stdlib", "math.koda")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find repo root (stdlib/math.koda)")
		}
		dir = parent
	}
}
