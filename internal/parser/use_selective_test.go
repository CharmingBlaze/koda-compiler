package parser_test

import (
	"strings"
	"testing"

	"koda/internal/parser"
)

func TestFilterModuleDecls(t *testing.T) {
	decls := []parser.Decl{
		&parser.LetDecl{Name: testToken("draw"), Init: &parser.LiteralExpr{Value: 0}},
		&parser.FuncDecl{Name: testToken("easeIn"), Body: &parser.BlockStmt{}},
		&parser.FuncDecl{Name: testToken("easeOut"), Body: &parser.BlockStmt{}},
	}
	filtered, err := parser.FilterModuleDecls(decls, []string{"easeIn"}, "koda.easing")
	if err != nil {
		t.Fatal(err)
	}
	if len(filtered) != 1 {
		t.Fatalf("expected 1 decl, got %d", len(filtered))
	}
	if fn, ok := filtered[0].(*parser.FuncDecl); !ok || fn.Name.Lexeme != "easeIn" {
		t.Fatalf("filtered = %#v", filtered[0])
	}
	_, err = parser.FilterModuleDecls(decls, []string{"missing"}, "koda.easing")
	if err == nil {
		t.Fatal("expected error for unknown selective import")
	}
}

func TestParseUseDeclSelective(t *testing.T) {
	src := `use raylib { InitWindow, DrawText };
use koda.easing { easeIn, easeOut };

func main() {}`
	prog := parseTestProgram(t, src)
	u0, ok := prog.Declarations[0].(*parser.UseDecl)
	if !ok || u0.ModulePath != "raylib" || len(u0.Selective) != 2 {
		t.Fatalf("use[0] = %#v", prog.Declarations[0])
	}
	u1, ok := prog.Declarations[1].(*parser.UseDecl)
	if !ok || u1.ModulePath != "koda.easing" || len(u1.Selective) != 2 {
		t.Fatalf("use[1] = %#v", prog.Declarations[1])
	}
}

func TestParseUseDeclSelectiveWithAlias(t *testing.T) {
	src := `use raylib { InitWindow } as rl;

func main() {}`
	prog := parseTestProgram(t, src)
	u, ok := prog.Declarations[0].(*parser.UseDecl)
	if !ok || u.Alias != "rl" || len(u.Selective) != 1 || !strings.EqualFold(u.Selective[0], "InitWindow") {
		t.Fatalf("use = %#v", prog.Declarations[0])
	}
}
