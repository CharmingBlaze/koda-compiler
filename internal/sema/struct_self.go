package sema

import (
	"fmt"
	"strings"

	"koda/internal/diagnostic"
	"koda/internal/parser"
)

func isSelfReceiverParam(name string) bool {
	return strings.EqualFold(name, "self")
}

func methodParamsForCall(params []parser.Param) []parser.Param {
	if len(params) > 0 && isSelfReceiverParam(params[0].Name) {
		return params[1:]
	}
	return params
}

func (a *Analyzer) checkReceiverOutsideMethod(te *parser.ThisExpr) {
	if a.currentStructType != "" {
		return
	}
	kw := te.Token.Lexeme
	if kw == "" {
		kw = "this"
	}
	a.record(&diagnostic.DiagnosticError{
		File:    te.Token.File,
		Line:    te.Token.Line,
		Col:     te.Token.Col,
		Message: fmt.Sprintf("'%s' can only be used inside struct methods", kw),
		Hint:    "move this expression into a func inside a struct body, or use a variable name",
	})
}

func (a *Analyzer) bindReceiverFieldAccess(e *parser.IndexExpr, stName string) bool {
	lit, ok := e.Index.(*parser.LiteralExpr)
	if !ok {
		return false
	}
	field, ok := lit.Value.(string)
	if !ok {
		return false
	}
	fields, ok := a.structLayouts[stName]
	if !ok {
		return false
	}
	for i, f := range fields {
		if f == field {
			a.indexExprStructSlot[e] = i
			return true
		}
	}
	a.record(&diagnostic.DiagnosticError{
		File:    lit.Token.File,
		Line:    lit.Token.Line,
		Col:     lit.Token.Col,
		Message: fmt.Sprintf("'%s' is not a field of struct %s", field, stName),
		Hint:    a.structFieldHint(field, fields),
	})
	return true
}
