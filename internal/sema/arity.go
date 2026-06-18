package sema

import (
	"fmt"
	"strings"

	"koda/internal/diagnostic"
	"koda/internal/parser"
)

// arityBoundsFromParams returns the minimum required argument count, maximum
// allowed count (-1 means unbounded because of a ...rest parameter), and
// whether the parameter list ends with rest.
func arityBoundsFromParams(params []parser.Param) (min int, max int, hasRest bool) {
	for _, p := range params {
		if p.IsRest {
			hasRest = true
			break
		}
		if p.Default == nil {
			min++
		}
	}
	if hasRest {
		return min, -1, true
	}
	return min, len(params), false
}

func formatParamList(params []parser.Param) string {
	parts := make([]string, 0, len(params))
	for _, p := range params {
		if p.IsRest {
			parts = append(parts, "..."+p.Name)
			break
		}
		parts = append(parts, p.Name)
	}
	return strings.Join(parts, ", ")
}

func (a *Analyzer) checkCallArity(call *parser.CallExpr, calleeName string, nativeArity int, native bool, params []parser.Param) {
	got := len(call.Arguments)
	if native {
		if got != nativeArity {
			a.record(&diagnostic.DiagnosticError{
				File:    call.Token.File,
				Line:    call.Token.Line,
				Col:     call.Token.Col,
				Message: fmt.Sprintf("wrong number of arguments to '%s': expected %d (native), got %d", calleeName, nativeArity, got),
			})
		}
		return
	}
	min, max, rest := arityBoundsFromParams(params)
	if got < min {
		a.record(&diagnostic.DiagnosticError{
			File:    call.Token.File,
			Line:    call.Token.Line,
			Col:     call.Token.Col,
			Message: fmt.Sprintf("too few arguments to '%s': expected at least %d (%s), got %d", calleeName, min, formatParamList(params), got),
		})
		return
	}
	if !rest && max >= 0 && got > max {
		a.record(&diagnostic.DiagnosticError{
			File:    call.Token.File,
			Line:    call.Token.Line,
			Col:     call.Token.Col,
			Message: fmt.Sprintf("too many arguments to '%s': expected at most %d (%s), got %d", calleeName, max, formatParamList(params), got),
		})
	}
}

func (a *Analyzer) maybeCheckCallArity(call *parser.CallExpr) {
	switch f := call.Function.(type) {
	case *parser.IdentifierExpr:
		name := f.Name.Lexeme
		decl, ok := a.currentScope.Resolve(name)
		if !ok {
			return
		}
		switch d := decl.(type) {
		case *parser.FuncDecl:
			if !strings.EqualFold(d.Name.Lexeme, name) {
				return
			}
			if d.Native != nil {
				a.checkCallArity(call, d.Name.Lexeme, d.Native.Arity, true, nil)
				return
			}
			a.checkCallArity(call, d.Name.Lexeme, 0, false, d.Params)
		case *parser.LetDecl:
			if d.Native != nil {
				a.checkCallArity(call, d.Name.Lexeme, d.Native.Arity, true, nil)
			}
		}
	case *parser.FuncExpr:
		a.checkCallArity(call, "<function>", 0, false, f.Params)
	case *parser.IndexExpr:
		a.maybeCheckArgvMethodCallArity(call, f)
		a.maybeCheckStructMethodCallArity(call, f)
	default:
		return
	}
}
