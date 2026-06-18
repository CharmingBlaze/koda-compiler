package sema

import (
	"koda/internal/parser"
	"koda/internal/lexer"
	"math"
)

func (c *NativeEmitContext) foldConstants(p *parser.Program) bool {
	// 1. Find constants (variables assigned only once to a literal)
	constants := c.findConstants(p)

	// 2. Fold with constant propagation
	any := false
	for i, d := range p.Declarations {
		newD, ch := c.foldDeclWithConstants(d, constants)
		p.Declarations[i] = newD
		if ch {
			any = true
		}
	}
	return any
}

func (c *NativeEmitContext) findConstants(p *parser.Program) map[parser.Decl]any {
	constants := make(map[parser.Decl]any)
	reassigned := make(map[parser.Decl]bool)

	// Pass 1: Find all assignments and reassignments
	c.walkForAssignments(p, func(d parser.Decl, val parser.Expr, isInitial bool) {
		if !isInitial {
			reassigned[d] = true
			delete(constants, d)
			return
		}
		if lit, ok := val.(*parser.LiteralExpr); ok && !reassigned[d] {
			constants[d] = lit.Value
		} else {
			reassigned[d] = true
			delete(constants, d)
		}
	})

	return constants
}

func (c *NativeEmitContext) walkForAssignments(p *parser.Program, cb func(parser.Decl, parser.Expr, bool)) {
	for _, d := range p.Declarations {
		c.walkDeclAssignments(d, cb)
	}
}

func (c *NativeEmitContext) walkDeclAssignments(d parser.Decl, cb func(parser.Decl, parser.Expr, bool)) {
	switch decl := d.(type) {
	case *parser.LetDecl:
		if decl.Init != nil {
			cb(decl, decl.Init, true)
			c.walkExprAssignments(decl.Init, cb)
		}
	case *parser.FuncDecl:
		c.walkStmtAssignments(decl.Body, cb)
	case parser.Stmt:
		c.walkStmtAssignments(decl, cb)
	}
}

func (c *NativeEmitContext) walkStmtAssignments(s parser.Stmt, cb func(parser.Decl, parser.Expr, bool)) {
	switch stmt := s.(type) {
	case *parser.ExpressionStmt:
		c.walkExprAssignments(stmt.Expr, cb)
	case *parser.BlockStmt:
		for _, d := range stmt.Declarations {
			c.walkDeclAssignments(d, cb)
		}
	case *parser.IfStmt:
		c.walkExprAssignments(stmt.Condition, cb)
		c.walkStmtAssignments(stmt.Then, cb)
		if stmt.Else != nil {
			c.walkStmtAssignments(stmt.Else, cb)
		}
	case *parser.WhileStmt:
		c.walkExprAssignments(stmt.Condition, cb)
		c.walkStmtAssignments(stmt.Body, cb)
	case *parser.LoopStmt:
		c.walkStmtAssignments(stmt.Body, cb)
	case *parser.ForStmt:
		for _, init := range stmt.Inits {
			c.walkDeclAssignments(init, cb)
		}
		if stmt.Condition != nil {
			c.walkExprAssignments(stmt.Condition, cb)
		}
		for _, inc := range stmt.Increments {
			c.walkExprAssignments(inc, cb)
		}
		c.walkStmtAssignments(stmt.Body, cb)
	case *parser.DoWhileStmt:
		c.walkStmtAssignments(stmt.Body, cb)
		c.walkExprAssignments(stmt.Condition, cb)
	case *parser.SwitchStmt:
		c.walkExprAssignments(stmt.Subject, cb)
		for _, cas := range stmt.Cases {
			c.walkExprAssignments(cas.Value, cb)
			for _, d := range cas.Body {
				c.walkDeclAssignments(d, cb)
			}
		}
		for _, d := range stmt.Default {
			c.walkDeclAssignments(d, cb)
		}
	case *parser.ForInStmt:
		c.walkExprAssignments(stmt.Iterable, cb)
		c.walkStmtAssignments(stmt.Body, cb)
	case *parser.ReturnStmt:
		if stmt.Value != nil {
			c.walkExprAssignments(stmt.Value, cb)
		}
	case *parser.DeferStmt:
		c.walkExprAssignments(stmt.Expr, cb)
	}
}

func (c *NativeEmitContext) walkExprAssignments(e parser.Expr, cb func(parser.Decl, parser.Expr, bool)) {
	switch expr := e.(type) {
	case *parser.AssignExpr:
		if v, ok := expr.Left.(*parser.IdentifierExpr); ok {
			d := c.locals[v]
			cb(d, expr.Value, false)
		}
		c.walkExprAssignments(expr.Value, cb)
	case *parser.UpdateExpr:
		if v, ok := expr.Operand.(*parser.IdentifierExpr); ok {
			d := c.locals[v]
			cb(d, nil, false)
		}
	case *parser.InfixExpr:
		c.walkExprAssignments(expr.Left, cb)
		c.walkExprAssignments(expr.Right, cb)
	case *parser.CallExpr:
		for _, arg := range expr.Arguments {
			c.walkExprAssignments(arg, cb)
		}
	case *parser.ArrayExpr:
		for _, el := range expr.Elements {
			c.walkExprAssignments(el, cb)
		}
	case *parser.SpreadExpr:
		c.walkExprAssignments(expr.Expr, cb)
	case *parser.ObjectExpr:
		for _, v := range expr.Values {
			c.walkExprAssignments(v, cb)
		}
	}
}

func isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	if n, ok := v.(float64); ok {
		return n != 0
	}
	return true
}

func (c *NativeEmitContext) foldDeclWithConstants(d parser.Decl, constants map[parser.Decl]any) (parser.Decl, bool) {
	changed := false
	switch decl := d.(type) {
	case *parser.LetDecl:
		if decl.Init != nil {
			newInit, c1 := c.foldExprWithConstants(decl.Init, constants)
			if c1 {
				decl.Init = newInit
				changed = true
			}
		}
	case *parser.FuncDecl:
		newBody, c2 := c.foldStmtWithConstants(decl.Body, constants)
		if c2 {
			decl.Body = newBody.(*parser.BlockStmt)
			changed = true
		}
	case parser.Stmt:
		return c.foldStmtWithConstants(decl, constants)
	}
	return d, changed
}

func (c *NativeEmitContext) foldStmtWithConstants(s parser.Stmt, constants map[parser.Decl]any) (parser.Stmt, bool) {
	changed := false
	switch stmt := s.(type) {
	case *parser.ExpressionStmt:
		newExpr, c1 := c.foldExprWithConstants(stmt.Expr, constants)
		if c1 {
			stmt.Expr = newExpr
			changed = true
			if _, ok := stmt.Expr.(*parser.LiteralExpr); ok {
				return &parser.BlockStmt{}, true
			}
		}
	case *parser.BlockStmt:
		for i, d := range stmt.Declarations {
			newD, c1 := c.foldDeclWithConstants(d, constants)
			if c1 {
				stmt.Declarations[i] = newD.(parser.Stmt)
				changed = true
			}
		}
	case *parser.IfStmt:
		newCond, c1 := c.foldExprWithConstants(stmt.Condition, constants)
		if c1 {
			stmt.Condition = newCond
			changed = true
		}
		if lit, ok := stmt.Condition.(*parser.LiteralExpr); ok {
			if isTruthy(lit.Value) {
				return c.foldStmtWithConstants(stmt.Then, constants)
			} else {
				if stmt.Else != nil {
					return c.foldStmtWithConstants(stmt.Else, constants)
				}
				return &parser.BlockStmt{}, true
			}
		}
		newThen, c2 := c.foldStmtWithConstants(stmt.Then, constants)
		if c2 {
			stmt.Then = newThen
			changed = true
		}
		if stmt.Else != nil {
			newElse, c3 := c.foldStmtWithConstants(stmt.Else, constants)
			if c3 {
				stmt.Else = newElse
				changed = true
			}
		}
	case *parser.WhileStmt:
		newCond, c1 := c.foldExprWithConstants(stmt.Condition, constants)
		if c1 {
			stmt.Condition = newCond
			changed = true
		}
		if lit, ok := stmt.Condition.(*parser.LiteralExpr); ok {
			if !isTruthy(lit.Value) {
				return &parser.BlockStmt{}, true
			}
		}
		newBody, c2 := c.foldStmtWithConstants(stmt.Body, constants)
		if c2 {
			stmt.Body = newBody
			changed = true
		}
	case *parser.LoopStmt:
		newBody, c1 := c.foldStmtWithConstants(stmt.Body, constants)
		if c1 {
			stmt.Body = newBody
			changed = true
		}
	case *parser.ForStmt:
		for i, init := range stmt.Inits {
			newInit, c1 := c.foldDeclWithConstants(init, constants)
			if c1 {
				stmt.Inits[i] = newInit
				changed = true
			}
		}
		if stmt.Condition != nil {
			newCond, c2 := c.foldExprWithConstants(stmt.Condition, constants)
			if c2 {
				stmt.Condition = newCond
				changed = true
			}
		}
		for i, inc := range stmt.Increments {
			newInc, c3 := c.foldExprWithConstants(inc, constants)
			if c3 {
				stmt.Increments[i] = newInc
				changed = true
			}
		}
		newBody, c4 := c.foldStmtWithConstants(stmt.Body, constants)
		if c4 {
			stmt.Body = newBody
			changed = true
		}
	case *parser.DoWhileStmt:
		newBody, c1 := c.foldStmtWithConstants(stmt.Body, constants)
		if c1 {
			stmt.Body = newBody
			changed = true
		}
		newCond, c2 := c.foldExprWithConstants(stmt.Condition, constants)
		if c2 {
			stmt.Condition = newCond
			changed = true
		}
	case *parser.SwitchStmt:
		newSubj, c1 := c.foldExprWithConstants(stmt.Subject, constants)
		if c1 {
			stmt.Subject = newSubj
			changed = true
		}
		for i, cas := range stmt.Cases {
			newVal, c2 := c.foldExprWithConstants(cas.Value, constants)
			if c2 {
				stmt.Cases[i].Value = newVal
				changed = true
			}
			for j, d := range cas.Body {
				newD, c3 := c.foldDeclWithConstants(d, constants)
				if c3 {
					stmt.Cases[i].Body[j] = newD
					changed = true
				}
			}
		}
		for i, d := range stmt.Default {
			newD, c4 := c.foldDeclWithConstants(d, constants)
			if c4 {
				stmt.Default[i] = newD
				changed = true
			}
		}
	case *parser.ForInStmt:
		newIter, c1 := c.foldExprWithConstants(stmt.Iterable, constants)
		if c1 {
			stmt.Iterable = newIter
			changed = true
		}
		newBody, c2 := c.foldStmtWithConstants(stmt.Body, constants)
		if c2 {
			stmt.Body = newBody
			changed = true
		}
	case *parser.ReturnStmt:
		if stmt.Value != nil {
			newVal, c1 := c.foldExprWithConstants(stmt.Value, constants)
			if c1 {
				stmt.Value = newVal
				changed = true
			}
		}
	case *parser.DeferStmt:
		newExpr, c1 := c.foldExprWithConstants(stmt.Expr, constants)
		if c1 {
			stmt.Expr = newExpr
			changed = true
		}
	}
	return s, changed
}

func (c *NativeEmitContext) foldExprWithConstants(e parser.Expr, constants map[parser.Decl]any) (parser.Expr, bool) {
	changed := false
	switch expr := e.(type) {
	case *parser.IdentifierExpr:
		d := c.locals[expr]
		if val, ok := constants[d]; ok {
			return &parser.LiteralExpr{Value: val}, true
		}
	case *parser.InfixExpr:
		left, c1 := c.foldExprWithConstants(expr.Left, constants)
		right, c2 := c.foldExprWithConstants(expr.Right, constants)
		if c1 {
			expr.Left = left
			changed = true
		}
		if c2 {
			expr.Right = right
			changed = true
		}

		lLit, lOk := left.(*parser.LiteralExpr)
		rLit, rOk := right.(*parser.LiteralExpr)

		if lOk && rOk {
			switch expr.Operator {
			case "+":
				if lNum, ok := lLit.Value.(float64); ok {
					if rNum, ok := rLit.Value.(float64); ok {
						return &parser.LiteralExpr{Value: lNum + rNum}, true
					}
				}
				if lStr, ok := lLit.Value.(string); ok {
					if rStr, ok := rLit.Value.(string); ok {
						return &parser.LiteralExpr{Value: lStr + rStr}, true
					}
				}
			case "-":
				if lNum, ok := lLit.Value.(float64); ok {
					if rNum, ok := rLit.Value.(float64); ok {
						return &parser.LiteralExpr{Value: lNum - rNum}, true
					}
				}
			case "*":
				if lNum, ok := lLit.Value.(float64); ok {
					if rNum, ok := rLit.Value.(float64); ok {
						return &parser.LiteralExpr{Value: lNum * rNum}, true
					}
				}
			case "/":
				if lNum, ok := lLit.Value.(float64); ok {
					if rNum, ok := rLit.Value.(float64); ok {
						if rNum != 0 {
							return &parser.LiteralExpr{Value: lNum / rNum}, true
						}
					}
				}
			case ">":
				if lNum, ok := lLit.Value.(float64); ok {
					if rNum, ok := rLit.Value.(float64); ok {
						return &parser.LiteralExpr{Value: lNum > rNum}, true
					}
				}
			case "<":
				if lNum, ok := lLit.Value.(float64); ok {
					if rNum, ok := rLit.Value.(float64); ok {
						return &parser.LiteralExpr{Value: lNum < rNum}, true
					}
				}
			case ">=":
				if lNum, ok := lLit.Value.(float64); ok {
					if rNum, ok := rLit.Value.(float64); ok {
						return &parser.LiteralExpr{Value: lNum >= rNum}, true
					}
				}
			case "<=":
				if lNum, ok := lLit.Value.(float64); ok {
					if rNum, ok := rLit.Value.(float64); ok {
						return &parser.LiteralExpr{Value: lNum <= rNum}, true
					}
				}
			case "==":
				return &parser.LiteralExpr{Value: lLit.Value == rLit.Value}, true
			case "!=":
				return &parser.LiteralExpr{Value: lLit.Value != rLit.Value}, true
			}
		}
	case *parser.LogicalExpr:
		left, c1 := c.foldExprWithConstants(expr.Left, constants)
		right, c2 := c.foldExprWithConstants(expr.Right, constants)
		if c1 {
			expr.Left = left
			changed = true
		}
		if c2 {
			expr.Right = right
			changed = true
		}
		lLit, lOk := left.(*parser.LiteralExpr)
		if lOk {
			switch expr.Operator.Type {
			case lexer.TokenOrOr:
				if isTruthy(lLit.Value) {
					return lLit, true
				}
				return right, true
			case lexer.TokenAndAnd:
				if !isTruthy(lLit.Value) {
					return lLit, true
				}
				return right, true
			}
		}
	case *parser.GroupingExpr:
		return c.foldExprWithConstants(expr.Expr, constants)
	case *parser.ArrayExpr:
		for i, el := range expr.Elements {
			newEl, c := c.foldExprWithConstants(el, constants)
			if c {
				expr.Elements[i] = newEl
				changed = true
			}
		}
	case *parser.SpreadExpr:
		newInner, ch := c.foldExprWithConstants(expr.Expr, constants)
		if ch {
			expr.Expr = newInner
			changed = true
		}
	case *parser.TupleExpr:
		for i, el := range expr.Elements {
			newEl, c := c.foldExprWithConstants(el, constants)
			if c {
				expr.Elements[i] = newEl
				changed = true
			}
		}
	case *parser.ObjectExpr:
		for i, v := range expr.Values {
			newV, c := c.foldExprWithConstants(v, constants)
			if c {
				expr.Values[i] = newV
				changed = true
			}
		}
	case *parser.CallExpr:
		newFunc, c1 := c.foldExprWithConstants(expr.Function, constants)
		if c1 {
			expr.Function = newFunc
			changed = true
		}
		for i, arg := range expr.Arguments {
			newArg, c2 := c.foldExprWithConstants(arg, constants)
			if c2 {
				expr.Arguments[i] = newArg
				changed = true
			}
		}
	case *parser.RangeExpr:
		newFrom, c1 := c.foldExprWithConstants(expr.From, constants)
		newTo, c2 := c.foldExprWithConstants(expr.To, constants)
		if c1 {
			expr.From = newFrom
			changed = true
		}
		if c2 {
			expr.To = newTo
			changed = true
		}
		if lf, ok := expr.From.(*parser.LiteralExpr); ok {
			if rt, ok2 := expr.To.(*parser.LiteralExpr); ok2 {
				if a, okA := lf.Value.(float64); okA {
					if b, okB := rt.Value.(float64); okB {
						lo := int(math.Trunc(a))
						hi := int(math.Trunc(b))
						var els []parser.Expr
						if lo <= hi {
							for i := lo; i <= hi; i++ {
								els = append(els, &parser.LiteralExpr{Value: float64(i)})
							}
						} else {
							for i := lo; i >= hi; i-- {
								els = append(els, &parser.LiteralExpr{Value: float64(i)})
							}
						}
						return &parser.ArrayExpr{Elements: els}, true
					}
				}
			}
		}
	case *parser.TemplateExpr:
		for i, p := range expr.Parts {
			newP, c := c.foldExprWithConstants(p, constants)
			if c {
				expr.Parts[i] = newP
				changed = true
			}
		}
	}
	return e, changed
}

func (c *NativeEmitContext) inlineFunctions(p *parser.Program) bool {
	inlinable := make(map[*parser.FuncDecl]parser.Expr)

	// 1. Identify inlinable functions (single return expr)
	for _, d := range p.Declarations {
		if fd, ok := d.(*parser.FuncDecl); ok {
			skip := false
			for _, p := range fd.Params {
				if p.Default != nil || p.IsRest {
					skip = true
					break
				}
			}
			if !skip && len(fd.Body.Declarations) == 1 {
				if rs, ok := fd.Body.Declarations[0].(*parser.ReturnStmt); ok {
					if rs.Value != nil {
						// Check for recursion (simple check: name not in expr)
						if !c.containsCallTo(rs.Value, fd.Name.Lexeme) {
							inlinable[fd] = rs.Value
						}
					}
				}
			}
		}
	}

	// 2. Perform inlining
	return c.performInlining(p, inlinable)
}

func (c *NativeEmitContext) containsCallTo(e parser.Expr, name string) bool {
	switch expr := e.(type) {
	case *parser.CallExpr:
		if v, ok := expr.Function.(*parser.IdentifierExpr); ok && v.Name.Lexeme == name {
			return true
		}
		for _, arg := range expr.Arguments {
			if c.containsCallTo(arg, name) {
				return true
			}
		}
	case *parser.InfixExpr:
		return c.containsCallTo(expr.Left, name) || c.containsCallTo(expr.Right, name)
	case *parser.GroupingExpr:
		return c.containsCallTo(expr.Expr, name)
	case *parser.RangeExpr:
		return c.containsCallTo(expr.From, name) || c.containsCallTo(expr.To, name)
	case *parser.TupleExpr:
		for _, el := range expr.Elements {
			if c.containsCallTo(el, name) {
				return true
			}
		}
	case *parser.SpreadExpr:
		return c.containsCallTo(expr.Expr, name)
	case *parser.TemplateExpr:
		for _, p := range expr.Parts {
			if c.containsCallTo(p, name) {
				return true
			}
		}
	}
	return false
}

func (c *NativeEmitContext) performInlining(p *parser.Program, inlinable map[*parser.FuncDecl]parser.Expr) bool {
	changed := false
	for i, d := range p.Declarations {
		newD, ok := c.inlineDecl(d, inlinable)
		if ok {
			p.Declarations[i] = newD
			changed = true
		}
	}
	return changed
}

func (c *NativeEmitContext) inlineDecl(d parser.Decl, inlinable map[*parser.FuncDecl]parser.Expr) (parser.Decl, bool) {
	changed := false
	switch decl := d.(type) {
	case *parser.LetDecl:
		if decl.Init != nil {
			newInit, ok := c.inlineExpr(decl.Init, inlinable)
			if ok {
				decl.Init = newInit
				changed = true
			}
		}
	case *parser.FuncDecl:
		newBody, ok := c.inlineStmt(decl.Body, inlinable)
		if ok {
			decl.Body = newBody.(*parser.BlockStmt)
			changed = true
		}
	case parser.Stmt:
		return c.inlineStmt(decl, inlinable)
	}
	return d, changed
}

func (c *NativeEmitContext) inlineStmt(s parser.Stmt, inlinable map[*parser.FuncDecl]parser.Expr) (parser.Stmt, bool) {
	changed := false
	switch stmt := s.(type) {
	case *parser.ExpressionStmt:
		newExpr, ok := c.inlineExpr(stmt.Expr, inlinable)
		if ok {
			stmt.Expr = newExpr
			changed = true
		}
	case *parser.BlockStmt:
		for i, d := range stmt.Declarations {
			newD, ok := c.inlineDecl(d, inlinable)
			if ok {
				stmt.Declarations[i] = newD
				changed = true
			}
		}
	case *parser.IfStmt:
		newCond, ok := c.inlineExpr(stmt.Condition, inlinable)
		if ok {
			stmt.Condition = newCond
			changed = true
		}
		newThen, c1 := c.inlineStmt(stmt.Then, inlinable)
		if c1 {
			stmt.Then = newThen
			changed = true
		}
		if stmt.Else != nil {
			newElse, c2 := c.inlineStmt(stmt.Else, inlinable)
			if c2 {
				stmt.Else = newElse
				changed = true
			}
		}
	case *parser.WhileStmt:
		newCond, ok := c.inlineExpr(stmt.Condition, inlinable)
		if ok {
			stmt.Condition = newCond
			changed = true
		}
		newBody, c1 := c.inlineStmt(stmt.Body, inlinable)
		if c1 {
			stmt.Body = newBody
			changed = true
		}
	case *parser.LoopStmt:
		newBody, ok := c.inlineStmt(stmt.Body, inlinable)
		if ok {
			stmt.Body = newBody
			changed = true
		}
	case *parser.ForStmt:
		newBody, ok := c.inlineStmt(stmt.Body, inlinable)
		if ok {
			stmt.Body = newBody
			changed = true
		}
	case *parser.DoWhileStmt:
		newBody, c1 := c.inlineStmt(stmt.Body, inlinable)
		if c1 {
			stmt.Body = newBody
			changed = true
		}
		newCond, c2 := c.inlineExpr(stmt.Condition, inlinable)
		if c2 {
			stmt.Condition = newCond
			changed = true
		}
	case *parser.SwitchStmt:
		newSubj, ok := c.inlineExpr(stmt.Subject, inlinable)
		if ok {
			stmt.Subject = newSubj
			changed = true
		}
	case *parser.ForInStmt:
		newIter, c1 := c.inlineExpr(stmt.Iterable, inlinable)
		if c1 {
			stmt.Iterable = newIter
			changed = true
		}
		newBody, c2 := c.inlineStmt(stmt.Body, inlinable)
		if c2 {
			stmt.Body = newBody
			changed = true
		}
	case *parser.ReturnStmt:
		if stmt.Value != nil {
			newVal, ok := c.inlineExpr(stmt.Value, inlinable)
			if ok {
				stmt.Value = newVal
				changed = true
			}
		}
	case *parser.DeferStmt:
		newExpr, ok := c.inlineExpr(stmt.Expr, inlinable)
		if ok {
			stmt.Expr = newExpr
			changed = true
		}
	}
	return s, changed
}

func (c *NativeEmitContext) inlineExpr(e parser.Expr, inlinable map[*parser.FuncDecl]parser.Expr) (parser.Expr, bool) {
	changed := false
	switch expr := e.(type) {
	case *parser.CallExpr:
		// Check if we can inline
		if v, ok := expr.Function.(*parser.IdentifierExpr); ok {
			d := c.locals[v]
			if fd, ok := d.(*parser.FuncDecl); ok {
				if body, ok := inlinable[fd]; ok {
					// Inline!
					// Substitute params with args
					return c.substitute(body, fd.Params, expr.Arguments), true
				}
			}
		}
		// Recurse on args
		for i, arg := range expr.Arguments {
			newArg, ok := c.inlineExpr(arg, inlinable)
			if ok {
				expr.Arguments[i] = newArg
				changed = true
			}
		}
	case *parser.InfixExpr:
		newLeft, c1 := c.inlineExpr(expr.Left, inlinable)
		newRight, c2 := c.inlineExpr(expr.Right, inlinable)
		if c1 {
			expr.Left = newLeft
			changed = true
		}
		if c2 {
			expr.Right = newRight
			changed = true
		}
	case *parser.GroupingExpr:
		return c.inlineExpr(expr.Expr, inlinable)
	}
	return e, changed
}

func (c *NativeEmitContext) substitute(e parser.Expr, params []parser.Param, args []parser.Expr) parser.Expr {
	// Deep copy and substitute
	switch expr := e.(type) {
	case *parser.IdentifierExpr:
		for i, p := range params {
			if p.Name == expr.Name.Lexeme {
				if i < len(args) {
					return args[i]
				}
				if p.Default != nil {
					return p.Default
				}
				return &parser.LiteralExpr{Value: nil}
			}
		}
	case *parser.InfixExpr:
		return &parser.InfixExpr{
			Left:     c.substitute(expr.Left, params, args),
			Operator: expr.Operator,
			Right:    c.substitute(expr.Right, params, args),
		}
	case *parser.GroupingExpr:
		return &parser.GroupingExpr{Expr: c.substitute(expr.Expr, params, args)}
	case *parser.LiteralExpr:
		return expr
	}
	return e
}
