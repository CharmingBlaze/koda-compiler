package sema

import (
	"koda/internal/parser"
)

// InferParamNumericKinds classifies untyped function parameters as KindInt or KindFloat
// from how they are used in the function body. Explicit : type annotations are handled
// separately via TypedParams — beginners can omit them and still get fast native math.
func InferParamNumericKinds(prog *parser.Program) map[ParamCellKey]NumericKind {
	out := make(map[ParamCellKey]NumericKind)
	if prog == nil {
		return out
	}
	var walkDecl func(parser.Decl)
	walkDecl = func(d parser.Decl) {
		switch x := d.(type) {
		case *parser.FuncDecl:
			recordInferredParamKinds(x, x.Params, x.Body, out)
			walkParamsInBlock(x.Body, walkDecl)
		case *parser.StructDecl:
			for _, m := range x.Methods {
				recordInferredParamKinds(m, m.Params, m.Body, out)
				walkParamsInBlock(m.Body, walkDecl)
			}
		case *parser.BlockStmt:
			for _, inner := range x.Declarations {
				walkDecl(inner)
			}
		case parser.Stmt:
			walkParamsInStmt(x, walkDecl)
		}
	}
	for _, d := range prog.Declarations {
		walkDecl(d)
	}
	inferFuncExprParamsInProgram(prog, out)
	return out
}

func walkParamsInBlock(b *parser.BlockStmt, walkDecl func(parser.Decl)) {
	if b == nil {
		return
	}
	for _, d := range b.Declarations {
		walkDecl(d)
	}
}

func walkParamsInStmt(s parser.Stmt, walkDecl func(parser.Decl)) {
	if s == nil {
		return
	}
	switch x := s.(type) {
	case *parser.BlockStmt:
		walkParamsInBlock(x, walkDecl)
	case *parser.IfStmt:
		walkParamsInStmt(x.Then, walkDecl)
		walkParamsInStmt(x.Else, walkDecl)
	case *parser.WhileStmt:
		walkParamsInStmt(x.Body, walkDecl)
	case *parser.LoopStmt:
		walkParamsInStmt(x.Body, walkDecl)
	case *parser.DoWhileStmt:
		walkParamsInStmt(x.Body, walkDecl)
	case *parser.ForStmt:
		walkParamsInStmt(x.Body, walkDecl)
	case *parser.ForInStmt:
		walkParamsInStmt(x.Body, walkDecl)
	case *parser.ForOfStmt:
		walkParamsInStmt(x.Body, walkDecl)
	case *parser.SwitchStmt:
		for _, c := range x.Cases {
			for _, cd := range c.Body {
				walkDecl(cd)
			}
		}
		for _, cd := range x.Default {
			walkDecl(cd)
		}
	}
}

func recordInferredParamKinds(owner interface{}, params []parser.Param, body *parser.BlockStmt, out map[ParamCellKey]NumericKind) {
	if owner == nil || len(params) == 0 || body == nil {
		return
	}
	hasUntyped := false
	for _, p := range params {
		if p.TypeAnnot == "" {
			hasUntyped = true
			break
		}
	}
	if !hasUntyped {
		return
	}
	kinds := inferParamKinds(owner, params, body)
	for idx, k := range kinds {
		if params[idx].TypeAnnot != "" {
			continue
		}
		out[NewParamCellKey(owner, idx)] = k
	}
}

func inferParamKinds(owner interface{}, params []parser.Param, body *parser.BlockStmt) []NumericKind {
	kinds := make([]NumericKind, len(params))
	paramIdx := make(map[string]int, len(params))
	for i, p := range params {
		kinds[i] = KindInt
		paramIdx[p.Name] = i
	}

	localKinds := make(map[*parser.LetDecl]NumericKind)

	narrowParam := func(idx int, k NumericKind) {
		if k == KindFloat {
			kinds[idx] = KindFloat
		}
	}

	narrowLocal := func(ld *parser.LetDecl, k NumericKind) {
		if ld == nil {
			return
		}
		if cur, ok := localKinds[ld]; ok && cur == KindFloat {
			return
		}
		localKinds[ld] = k
	}

	var exprKind func(parser.Expr) NumericKind
	var propagateKind func(parser.Expr, NumericKind)
	var walkExpr func(parser.Expr)
	var walkStmt func(parser.Stmt)
	var walkDecl func(parser.Decl)

	exprKind = func(e parser.Expr) NumericKind {
		if e == nil {
			return KindFloat
		}
		switch x := e.(type) {
		case *parser.LiteralExpr:
			switch x.Value.(type) {
			case int:
				return KindInt
			case float64:
				return KindFloat
			default:
				return KindFloat
			}
		case *parser.IdentifierExpr:
			name := x.Name.Lexeme
			if idx, ok := paramIdx[name]; ok {
				return kinds[idx]
			}
			if ld := findLetInBlock(body, name); ld != nil {
				if k, ok := localKinds[ld]; ok {
					return k
				}
			}
			return KindFloat
		case *parser.InfixExpr:
			lk := exprKind(x.Left)
			rk := exprKind(x.Right)
			switch x.Operator {
			case "+", "-", "*", "%", "&", "|", "^", "<<", ">>", ">>>":
				if lk == KindFloat || rk == KindFloat {
					return KindFloat
				}
				if lk == KindInt && rk == KindInt && x.Operator != "/" {
					return KindInt
				}
			case "/":
				return KindFloat
			case "<", "<=", ">", ">=", "==", "!=":
				if lk == KindInt && rk == KindInt {
					return KindInt
				}
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

	propagateKind = func(e parser.Expr, k NumericKind) {
		if e == nil {
			return
		}
		switch x := e.(type) {
		case *parser.IdentifierExpr:
			if idx, ok := paramIdx[x.Name.Lexeme]; ok {
				narrowParam(idx, k)
			}
		case *parser.InfixExpr:
			propagateKind(x.Left, k)
			propagateKind(x.Right, k)
		case *parser.PrefixExpr:
			propagateKind(x.Right, k)
		case *parser.GroupingExpr:
			propagateKind(x.Expr, k)
		case *parser.AssignExpr:
			propagateKind(x.Left, k)
			propagateKind(x.Value, k)
		}
	}

	walkExpr = func(e parser.Expr) {
		if e == nil {
			return
		}
		switch x := e.(type) {
		case *parser.InfixExpr:
			lk := exprKind(x.Left)
			rk := exprKind(x.Right)
			switch x.Operator {
			case "<", "<=", ">", ">=", "==", "!=":
				if lk == KindInt {
					propagateKind(x.Right, KindInt)
				}
				if rk == KindInt {
					propagateKind(x.Left, KindInt)
				}
			default:
				propagateKind(e, exprKind(e))
			}
			walkExpr(x.Left)
			walkExpr(x.Right)
		case *parser.PrefixExpr:
			propagateKind(e, exprKind(e))
			walkExpr(x.Right)
		case *parser.AssignExpr:
			walkExpr(x.Left)
			walkExpr(x.Value)
			if id, ok := x.Left.(*parser.IdentifierExpr); ok {
				if ld := findLetInBlock(body, id.Name.Lexeme); ld != nil {
					narrowLocal(ld, exprKind(x.Value))
				}
			}
			propagateKind(x.Value, exprKind(x.Value))
		case *parser.CallExpr:
			walkExpr(x.Function)
			for _, a := range x.Arguments {
				walkExpr(a)
			}
		case *parser.IndexExpr:
			walkExpr(x.Object)
			walkExpr(x.Index)
		case *parser.GroupingExpr:
			walkExpr(x.Expr)
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
		case *parser.LoopStmt:
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
				if isFloatTypeName(x.TypeAnnot) {
					narrowLocal(x, KindFloat)
				} else if isIntegerTypeName(x.TypeAnnot) {
					narrowLocal(x, KindInt)
				} else {
					narrowLocal(x, KindFloat)
				}
			} else if x.Init != nil {
				narrowLocal(x, exprKind(x.Init))
			}
			walkExpr(x.Init)
		case *parser.FuncDecl:
			walkStmt(x.Body)
		case parser.Stmt:
			walkStmt(x)
		}
	}

	for pass := 0; pass < 2; pass++ {
		walkStmt(body)
	}
	return kinds
}

func findLetInBlock(b *parser.BlockStmt, name string) *parser.LetDecl {
	if b == nil {
		return nil
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
		}
		return nil
	}
	for _, d := range b.Declarations {
		if ld := findInDecl(d); ld != nil {
			return ld
		}
	}
	return nil
}
