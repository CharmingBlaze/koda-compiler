package parser_test

import (
	"strings"
	"testing"

	"koda/internal/lexer"
	"koda/internal/parser"
)

func testToken(lexeme string) lexer.Token {
	return lexer.Token{Type: lexer.TokenIdentifier, Lexeme: lexeme, File: "<test>"}
}

func TestInjectRaylibPrelude(t *testing.T) {
	bundle := &parser.ProgramBundle{
		Entry: &parser.Program{
			Declarations: []parser.Decl{
				&parser.LetDecl{Name: testToken("initwindow")},
			},
		},
	}
	parser.InjectRaylibPrelude(bundle)
	foundCamera := false
	foundRaywhite := false
	for _, d := range bundle.Entry.Declarations {
		switch x := d.(type) {
		case *parser.StructDecl:
			if strings.EqualFold(x.Name.Lexeme, "camera3d") {
				foundCamera = true
			}
		case *parser.LetDecl:
			if strings.EqualFold(x.Name.Lexeme, "raywhite") {
				foundRaywhite = true
			}
		}
	}
	if !foundCamera {
		t.Fatal("expected Camera3D struct from raylib prelude")
	}
	if !foundRaywhite {
		t.Fatal("expected RAYWHITE constant from raylib prelude")
	}
}

func TestInjectRaylibPreludeSkipsWithoutRaylib(t *testing.T) {
	bundle := &parser.ProgramBundle{
		Entry: &parser.Program{
			Declarations: []parser.Decl{
				&parser.FuncDecl{Name: testToken("main")},
			},
		},
	}
	before := len(bundle.Entry.Declarations)
	parser.InjectRaylibPrelude(bundle)
	if len(bundle.Entry.Declarations) != before {
		t.Fatalf("expected no prelude without raylib symbols, got %d decls", len(bundle.Entry.Declarations))
	}
}
