package parser_test

import (
	"testing"

	"koda/internal/lexer"
	"koda/internal/parser"
)

func TestWrapModuleAsAlias(t *testing.T) {
	decls := []parser.Decl{
		&parser.LetDecl{Name: testToken("draw"), Init: &parser.LiteralExpr{Value: 0}},
		&parser.FuncDecl{
			Name: testToken("ping"),
			Body: &parser.BlockStmt{},
		},
	}
	wrapped, err := parser.WrapModuleAsAlias(decls, "rl")
	if err != nil {
		t.Fatal(err)
	}
	foundAlias := false
	foundHidden := false
	for _, d := range wrapped {
		if let, ok := d.(*parser.LetDecl); ok {
			if let.Name.Lexeme == "rl" && let.Init != nil {
				foundAlias = true
			}
			if let.Name.Lexeme == "__koda_use_rl_draw" {
				foundHidden = true
			}
		}
	}
	if !foundAlias {
		t.Fatal("expected namespace let rl")
	}
	if !foundHidden {
		t.Fatal("expected hidden let __koda_use_rl_draw")
	}
}

func TestParseUseDeclWithAlias(t *testing.T) {
	src := `use raylib as rl;
use koda.math as m;

func main() {}`
	prog := parseTestProgram(t, src)
	if len(prog.Declarations) < 2 {
		t.Fatalf("expected 2 use decls, got %d", len(prog.Declarations))
	}
	u0, ok := prog.Declarations[0].(*parser.UseDecl)
	if !ok || u0.ModulePath != "raylib" || u0.Alias != "rl" {
		t.Fatalf("use[0] = %#v", prog.Declarations[0])
	}
	u1, ok := prog.Declarations[1].(*parser.UseDecl)
	if !ok || u1.ModulePath != "koda.math" || u1.Alias != "m" {
		t.Fatalf("use[1] = %#v", prog.Declarations[1])
	}
}

func parseTestProgram(t *testing.T, src string) *parser.Program {
	t.Helper()
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
	return prog
}
