package parser

import (
	"fmt"
	"strings"

	"koda/internal/lexer"
)

// WrapModuleAsAlias exports aliased module wrapping for tests and tooling.
func WrapModuleAsAlias(decls []Decl, alias string) ([]Decl, error) {
	return wrapModuleAsAlias(decls, alias)
}

func useAliasPrefix(alias string) string {
	return "__koda_use_" + strings.ToLower(alias) + "_"
}

// wrapModuleAsAlias rewrites flattened module declarations under a single namespace object.
// Exported lets and funcs become `alias.name`; other decls (structs, enums) stay global.
func wrapModuleAsAlias(decls []Decl, alias string) ([]Decl, error) {
	if strings.TrimSpace(alias) == "" {
		return nil, fmt.Errorf("empty use alias")
	}
	prefix := useAliasPrefix(alias)
	var out []Decl
	var keys []lexer.Token
	var values []Expr

	for _, d := range decls {
		switch x := d.(type) {
		case *LetDecl:
			hidden := cloneLetWithName(x, prefix+x.Name.Lexeme)
			out = append(out, hidden)
			keys = append(keys, x.Name)
			values = append(values, &IdentifierExpr{Name: hidden.Name})
		case *FuncDecl:
			if strings.EqualFold(x.Name.Lexeme, "main") {
				out = append(out, d)
				continue
			}
			hidden := cloneFuncWithName(x, prefix+x.Name.Lexeme)
			out = append(out, hidden)
			keys = append(keys, x.Name)
			values = append(values, &IdentifierExpr{Name: hidden.Name})
		default:
			out = append(out, d)
		}
	}

	aliasTok := lexer.Token{
		Type:   lexer.TokenIdentifier,
		Lexeme: alias,
		File:   "<use-alias>",
	}
	out = append(out, &LetDecl{
		Token: aliasTok,
		Name:  aliasTok,
		Init: &ObjectExpr{
			Token:  aliasTok,
			Keys:   keys,
			Values: values,
		},
	})
	return out, nil
}

func cloneLetWithName(src *LetDecl, name string) *LetDecl {
	if src == nil {
		return nil
	}
	cp := *src
	cp.Name = lexer.Token{
		Type:   src.Name.Type,
		Lexeme: name,
		Line:   src.Name.Line,
		Col:    src.Name.Col,
		File:   src.Name.File,
	}
	return &cp
}

func cloneFuncWithName(src *FuncDecl, name string) *FuncDecl {
	if src == nil {
		return nil
	}
	cp := *src
	cp.Name = lexer.Token{
		Type:   src.Name.Type,
		Lexeme: name,
		Line:   src.Name.Line,
		Col:    src.Name.Col,
		File:   src.Name.File,
	}
	return &cp
}
