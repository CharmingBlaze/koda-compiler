package sema

import (
	"fmt"

	"koda/internal/parser"
)

// checkTruthyCondition warns when a condition uses array/object truthiness (differs from JavaScript).
func (a *Analyzer) checkTruthyCondition(cond parser.Expr) {
	if cond == nil {
		return
	}
	switch e := cond.(type) {
	case *parser.ArrayExpr:
		a.warnAtExpr(e, "empty arrays are truthy in Koda (unlike JavaScript); use len(arr) > 0 or arr.length > 0")
	case *parser.ObjectExpr:
		if e.StructTag == nil {
			a.warnAtExpr(e, "objects are always truthy in Koda (unlike JavaScript); check a specific field instead")
		}
	case *parser.IdentifierExpr:
		name := e.Name.Lexeme
		if a.varIsArray[name] {
			a.warnAtExpr(e, fmt.Sprintf("array '%s' is always truthy in Koda (unlike JavaScript); use len(%s) > 0", name, name))
		}
		if _, ok := a.varStructType[name]; ok {
			a.warnAtExpr(e, fmt.Sprintf("struct '%s' is always truthy; compare a field instead", name))
		}
	}
}

func (a *Analyzer) warnAtExpr(e parser.Expr, msg string) {
	var file string
	var line, col int
	switch x := e.(type) {
	case *parser.ArrayExpr:
		file, line, col = x.Token.File, x.Token.Line, x.Token.Col
	case *parser.ObjectExpr:
		file, line, col = x.Token.File, x.Token.Line, x.Token.Col
	case *parser.IdentifierExpr:
		file, line, col = x.Name.File, x.Name.Line, x.Name.Col
	default:
		return
	}
	a.warn(fmt.Sprintf("%s:%d:%d: %s", file, line, col, msg))
}
