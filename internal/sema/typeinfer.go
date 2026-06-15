package sema

import (
	"math"

	"koda/internal/parser"
)

// NumericKind classifies stack-local numeric bindings for fast integer codegen (A8).
type NumericKind int

const (
	KindFloat NumericKind = iota
	KindInt
)

// InferNumericKinds classifies every LetDecl as KindInt or KindFloat.
func InferNumericKinds(prog *parser.Program, escaping map[*parser.LetDecl]bool) map[*parser.LetDecl]NumericKind {
	out := make(map[*parser.LetDecl]NumericKind)
	if prog == nil {
		return out
	}

	narrow := func(ld *parser.LetDecl, k NumericKind) {
		if ld == nil {
			return
		}
		if escaping != nil && escaping[ld] {
			out[ld] = KindFloat
			return
		}
		if cur, ok := out[ld]; ok && cur == KindFloat {
			return
		}
		out[ld] = k
	}

	var exprKind func(parser.Expr) NumericKind
	var walkExpr func(parser.Expr)
	var walkStmt func(parser.Stmt)
	var walkDecl func(parser.Decl)

	exprKind = func(e parser.Expr) NumericKind {
		if e == nil {
			return KindFloat
		}
		switch x := e.(type) {
		case *parser.LiteralExpr:
			switch v := x.Value.(type) {
			case int:
				return KindInt
			case float64:
				if math.IsNaN(v) || math.IsInf(v, 0) {
					return KindFloat
				}
				if v == math.Trunc(v) {
					return KindInt
				}
				return KindFloat
			default:
				return KindFloat
			}
		case *parser.IdentifierExpr:
			if decl, ok := resolveLetInProgram(prog, x.Name.Lexeme); ok {
				if k, ok2 := out[decl]; ok2 {
					return k
				}
			}
			return KindFloat
		case *parser.InfixExpr:
			lk := exprKind(x.Left)
			rk := exprKind(x.Right)
			switch x.Operator {
			case "+", "-", "*", "%", "&", "|", "^", "<<", ">>", ">>>":
				if lk == KindInt && rk == KindInt && x.Operator != "/" {
					return KindInt
				}
			case "/":
				return KindFloat
			}
			return KindFloat
		case *parser.PrefixExpr:
			if x.Operator == "+" || x.Operator == "-" {
				return exprKind(x.Right)
			}
			return KindFloat
		case *parser.GroupingExpr:
			return exprKind(x.Expr)
		default:
			return KindFloat
		}
	}

	walkExpr = func(e parser.Expr) {
		if e == nil {
			return
		}
		switch x := e.(type) {
		case *parser.InfixExpr:
			walkExpr(x.Left)
			walkExpr(x.Right)
		case *parser.PrefixExpr:
			walkExpr(x.Right)
		case *parser.AssignExpr:
			walkExpr(x.Left)
			walkExpr(x.Value)
			if id, ok := x.Left.(*parser.IdentifierExpr); ok {
				if decl, ok2 := resolveLetInProgram(prog, id.Name.Lexeme); ok2 {
					narrow(decl, exprKind(x.Value))
				}
			}
		case *parser.CallExpr:
			walkExpr(x.Function)
			for _, a := range x.Arguments {
				walkExpr(a)
			}
		case *parser.IndexExpr:
			walkExpr(x.Object)
			walkExpr(x.Index)
		}
	}

	walkStmt = func(s parser.Stmt) {
		if s == nil {
			return
		}
		switch x := s.(type) {
		case *parser.BlockStmt:
			for _, d := range x.Declarations {
				walkDecl(d)
			}
		case *parser.ExpressionStmt:
			walkExpr(x.Expr)
		case *parser.ReturnStmt:
			walkExpr(x.Value)
		case *parser.IfStmt:
			walkExpr(x.Condition)
			walkStmt(x.Then)
			walkStmt(x.Else)
		case *parser.WhileStmt:
			walkExpr(x.Condition)
			walkStmt(x.Body)
		case *parser.DoWhileStmt:
			walkStmt(x.Body)
			walkExpr(x.Condition)
		case *parser.ForStmt:
			for _, ini := range x.Inits {
				walkDecl(ini)
			}
			walkExpr(x.Condition)
			for _, inc := range x.Increments {
				walkExpr(inc)
			}
			walkStmt(x.Body)
		case *parser.ForInStmt:
			walkExpr(x.Iterable)
			walkStmt(x.Body)
		case *parser.ForOfStmt:
			walkExpr(x.Iterable)
			walkStmt(x.Body)
		case *parser.SwitchStmt:
			walkExpr(x.Subject)
			for _, c := range x.Cases {
				walkExpr(c.Value)
				for _, d := range c.Body {
					walkDecl(d)
				}
			}
			for _, d := range x.Default {
				walkDecl(d)
			}
		case *parser.DeferStmt:
			walkExpr(x.Expr)
		}
	}

	walkDecl = func(d parser.Decl) {
		switch x := d.(type) {
		case *parser.LetDecl:
			if x.TypeAnnot != "" {
				narrow(x, KindInt)
			} else if x.Init != nil {
				narrow(x, exprKind(x.Init))
			} else {
				narrow(x, KindFloat)
			}
			walkExpr(x.Init)
		case *parser.FuncDecl:
			walkStmt(x.Body)
		case *parser.TestDecl:
			walkStmt(x.Body)
		case parser.Stmt:
			walkStmt(x)
		}
	}

	for pass := 0; pass < 2; pass++ {
		for _, d := range prog.Declarations {
			walkDecl(d)
		}
	}
	return out
}

func resolveLetInProgram(prog *parser.Program, name string) (*parser.LetDecl, bool) {
	if prog == nil {
		return nil, false
	}
	var findInDecl func(parser.Decl) *parser.LetDecl
	findInDecl = func(d parser.Decl) *parser.LetDecl {
		switch x := d.(type) {
		case *parser.LetDecl:
			if x.Name.Lexeme == name {
				return x
			}
		case *parser.FuncDecl:
			if x.Body != nil {
				for _, inner := range x.Body.Declarations {
					if ld := findInDecl(inner); ld != nil {
						return ld
					}
				}
			}
		case *parser.BlockStmt:
			for _, inner := range x.Declarations {
				if ld := findInDecl(inner); ld != nil {
					return ld
				}
			}
		case parser.Stmt:
			// handled via BlockStmt
		}
		return nil
	}
	for _, d := range prog.Declarations {
		if ld := findInDecl(d); ld != nil {
			return ld, true
		}
	}
	return nil, false
}
