package sema

import (
	"fmt"
	"strings"

	"koda/internal/diagnostic"
	"koda/internal/parser"
)

// argvMethodArityBounds returns inclusive [min,max] for known argv / array methods
// lowered in codegen (internal/codegen/methods.go). known is false for names we do not validate here.
func argvMethodArityBounds(lowerName string) (min, max int, known bool) {
	switch lowerName {
	case "trim", "toupper", "tolower", "sort", "reverse", "pop":
		return 0, 0, true
	case "split", "startswith", "endswith", "join", "indexof", "includes":
		return 1, 1, true
	case "replace", "replaceall", "slice":
		return 2, 2, true
	case "push":
		return 1, 1, true
	case "map", "filter", "find":
		return 1, 1, true
	case "reduce":
		return 1, 2, true
	case "concat":
		// Variadic: receiver plus zero or more arrays — codegen accepts any count.
		return 0, 0, false
	default:
		return 0, 0, false
	}
}

func formatArgvMethodArityExpect(min, max int) string {
	if min == max {
		return fmt.Sprintf("%d", min)
	}
	return fmt.Sprintf("%d or %d", min, max)
}

// maybeCheckArgvMethodCallArity validates argument counts for obj.prop(...) when prop is a
// string literal naming a builtin argv method. User-defined methods and dynamic property names are skipped.
func (a *Analyzer) maybeCheckArgvMethodCallArity(call *parser.CallExpr, idx *parser.IndexExpr) {
	lit, ok := idx.Index.(*parser.LiteralExpr)
	if !ok {
		return
	}
	s, ok := lit.Value.(string)
	if !ok {
		return
	}
	name := strings.ToLower(s)
	min, max, known := argvMethodArityBounds(name)
	if !known {
		return
	}
	got := len(call.Arguments)
	if got < min || got > max {
		a.record(&diagnostic.DiagnosticError{
			File:    call.Token.File,
			Line:    call.Token.Line,
			Col:     call.Token.Col,
			Message: fmt.Sprintf("wrong number of arguments to '.%s()': expected %s, got %d", name, formatArgvMethodArityExpect(min, max), got),
		})
	}
}
