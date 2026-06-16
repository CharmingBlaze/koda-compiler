package sema

import (
	"fmt"

	"koda/internal/parser"
)

func (a *Analyzer) checkUnreachableCode(body *parser.BlockStmt) {
	if body == nil {
		return
	}
	a.walkBlockUnreachable(body, false)
}

func (a *Analyzer) walkBlockUnreachable(block *parser.BlockStmt, afterTerminal bool) {
	if block == nil {
		return
	}
	terminal := afterTerminal
	for _, d := range block.Declarations {
		if terminal {
			a.warnUnreachableDecl(d)
			continue
		}
		terminal = declEndsWithTerminal(d) || terminal
	}
}

func (a *Analyzer) warnUnreachableDecl(d parser.Decl) {
	switch x := d.(type) {
	case *parser.LetDecl:
		a.warn(fmt.Sprintf("%s:%d:%d: unreachable code after return/break/continue",
			x.Name.File, x.Name.Line, x.Name.Col))
	case *parser.FuncDecl:
		a.warn(fmt.Sprintf("%s:%d:%d: unreachable function '%s' after return/break/continue",
			x.Name.File, x.Name.Line, x.Name.Col, x.Name.Lexeme))
	case *parser.TestDecl:
		a.warn(fmt.Sprintf("%s:%d:%d: unreachable test %q after return/break/continue",
			x.Token.File, x.Token.Line, x.Token.Col, x.Display.Lexeme))
	case parser.Stmt:
		a.warnUnreachableStmt(x)
	}
}

func (a *Analyzer) warnUnreachableStmt(s parser.Stmt) {
	switch x := s.(type) {
	case *parser.ExpressionStmt:
		if x.Expr != nil {
			if ce, ok := x.Expr.(*parser.CallExpr); ok {
				a.warn(fmt.Sprintf("%s:%d:%d: unreachable code after return/break/continue",
					ce.Token.File, ce.Token.Line, ce.Token.Col))
				return
			}
		}
	case *parser.ReturnStmt:
		a.warn(fmt.Sprintf("%s:%d:%d: unreachable code after return/break/continue",
			x.Token.File, x.Token.Line, x.Token.Col))
	}
}

func declEndsWithTerminal(d parser.Decl) bool {
	if s, ok := d.(parser.Stmt); ok {
		return stmtEndsWithTerminal(s)
	}
	return false
}

func stmtEndsWithTerminal(s parser.Stmt) bool {
	switch x := s.(type) {
	case *parser.ReturnStmt:
		return true
	case *parser.BreakStmt, *parser.ContinueStmt, *parser.FallthroughStmt:
		return true
	case *parser.BlockStmt:
		term := false
		for _, inner := range x.Declarations {
			if term {
				return true
			}
			term = declEndsWithTerminal(inner)
		}
		return term
	case *parser.IfStmt:
		thenTerm := stmtEndsWithTerminal(x.Then)
		elseTerm := x.Else != nil && stmtEndsWithTerminal(x.Else)
		return thenTerm && elseTerm
	default:
		return false
	}
}
