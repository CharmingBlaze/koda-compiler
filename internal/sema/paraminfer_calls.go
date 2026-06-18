package sema

import (
	"strings"

	"koda/internal/parser"
)

// knownFloatReturningNames are call targets treated as float for param inference at call sites.
var knownFloatReturningNames = map[string]bool{
	"delta":        true,
	"clampdelta":   true,
	"getframetime": true,
	"lerp":         true,
	"random":       true,
}

// InferParamKindsFromCallSites promotes untyped parameters to KindFloat when every
// observed call passes a float-like value (literals, game.delta(), inferred float locals).
func InferParamKindsFromCallSites(prog *parser.Program, letKinds map[*parser.LetDecl]NumericKind, out map[ParamCellKey]NumericKind) {
	if prog == nil || out == nil {
		return
	}
	funcs := collectProgramFuncs(prog)
	structMethods := collectStructMethodsFromProgram(prog)
	varStructs := collectVarStructTypesFromProgram(prog)
	var scanExpr func(parser.Expr)
	scanExpr = func(e parser.Expr) {
		if e == nil {
			return
		}
		if call, ok := e.(*parser.CallExpr); ok {
			if fd, ok := resolveUserFuncCallee(call.Function, funcs); ok {
				promoteParamKindsFromCall(fd, fd.Params, call, letKinds, prog, out)
			} else if fd, ok := resolveStructMethodCallee(call, structMethods, varStructs); ok {
				promoteParamKindsFromCall(fd, fd.Params, call, letKinds, prog, out)
			}
			scanExpr(call.Function)
			for _, a := range call.Arguments {
				scanExpr(a)
			}
			return
		}
		switch x := e.(type) {
		case *parser.FuncExpr:
			recordInferredParamKinds(x, x.Params, x.Body, out)
			scanParamsInBlock(x.Body, scanExpr)
		case *parser.InfixExpr:
			scanExpr(x.Left)
			scanExpr(x.Right)
		case *parser.PrefixExpr:
			scanExpr(x.Right)
		case *parser.AssignExpr:
			scanExpr(x.Left)
			scanExpr(x.Value)
		case *parser.LogicalExpr:
			scanExpr(x.Left)
			scanExpr(x.Right)
		case *parser.IndexExpr:
			scanExpr(x.Object)
			scanExpr(x.Index)
		case *parser.GroupingExpr:
			scanExpr(x.Expr)
		case *parser.ArrayExpr:
			for _, el := range x.Elements {
				scanExpr(el)
			}
		case *parser.IfExpr:
			scanExpr(x.Condition)
			scanExpr(x.Then)
			if x.Else != nil {
				scanExpr(x.Else)
			}
		case *parser.TernaryExpr:
			scanExpr(x.Condition)
			scanExpr(x.Then)
			scanExpr(x.Else)
		}
	}
	var walkDecl func(parser.Decl)
	walkDecl = func(d parser.Decl) {
		switch x := d.(type) {
		case *parser.LetDecl:
			scanExpr(x.Init)
		case *parser.FuncDecl:
			scanParamsInBlock(x.Body, scanExpr)
		case *parser.StructDecl:
			for _, m := range x.Methods {
				scanParamsInBlock(m.Body, scanExpr)
			}
		case *parser.BlockStmt:
			for _, inner := range x.Declarations {
				walkDecl(inner)
			}
		case parser.Stmt:
			walkParamsInStmtForScan(x, walkDecl, scanExpr)
		}
	}
	for _, d := range prog.Declarations {
		walkDecl(d)
	}
}

func walkParamsInStmtForScan(s parser.Stmt, walkDecl func(parser.Decl), scanExpr func(parser.Expr)) {
	if s == nil {
		return
	}
	switch x := s.(type) {
	case *parser.BlockStmt:
		for _, d := range x.Declarations {
			walkDecl(d)
		}
	case *parser.ExpressionStmt:
		scanExpr(x.Expr)
	case *parser.ReturnStmt:
		scanExpr(x.Value)
	case *parser.IfStmt:
		scanExpr(x.Condition)
		walkParamsInStmtForScan(x.Then, walkDecl, scanExpr)
		walkParamsInStmtForScan(x.Else, walkDecl, scanExpr)
	case *parser.WhileStmt:
		scanExpr(x.Condition)
		walkParamsInStmtForScan(x.Body, walkDecl, scanExpr)
	case *parser.LoopStmt:
		walkParamsInStmtForScan(x.Body, walkDecl, scanExpr)
	case *parser.DoWhileStmt:
		walkParamsInStmtForScan(x.Body, walkDecl, scanExpr)
		scanExpr(x.Condition)
	case *parser.ForStmt:
		for _, ini := range x.Inits {
			walkDecl(ini)
		}
		scanExpr(x.Condition)
		for _, inc := range x.Increments {
			scanExpr(inc)
		}
		walkParamsInStmtForScan(x.Body, walkDecl, scanExpr)
	case *parser.ForInStmt:
		scanExpr(x.Iterable)
		walkParamsInStmtForScan(x.Body, walkDecl, scanExpr)
	case *parser.ForOfStmt:
		scanExpr(x.Iterable)
		walkParamsInStmtForScan(x.Body, walkDecl, scanExpr)
	case *parser.SwitchStmt:
		scanExpr(x.Subject)
		for _, c := range x.Cases {
			scanExpr(c.Value)
			for _, cd := range c.Body {
				walkDecl(cd)
			}
		}
		for _, cd := range x.Default {
			walkDecl(cd)
		}
	case *parser.DeferStmt:
		scanExpr(x.Expr)
	}
}

func scanParamsInBlock(b *parser.BlockStmt, scanExpr func(parser.Expr)) {
	if b == nil {
		return
	}
	var walkDecl func(parser.Decl)
	walkDecl = func(d parser.Decl) {
		switch x := d.(type) {
		case *parser.LetDecl:
			scanExpr(x.Init)
		case *parser.FuncDecl:
			scanParamsInBlock(x.Body, scanExpr)
		case *parser.BlockStmt:
			for _, inner := range x.Declarations {
				walkDecl(inner)
			}
		case parser.Stmt:
			walkParamsInStmtForScan(x, walkDecl, scanExpr)
		}
	}
	for _, d := range b.Declarations {
		walkDecl(d)
	}
}

func collectProgramFuncs(prog *parser.Program) map[string]*parser.FuncDecl {
	out := make(map[string]*parser.FuncDecl)
	if prog == nil {
		return out
	}
	var walk func(parser.Decl)
	walk = func(d parser.Decl) {
		switch x := d.(type) {
		case *parser.FuncDecl:
			out[x.Name.Lexeme] = x
			if x.Body != nil {
				for _, inner := range x.Body.Declarations {
					walk(inner)
				}
			}
		case *parser.BlockStmt:
			for _, inner := range x.Declarations {
				walk(inner)
			}
		case parser.Stmt:
			walkParamsInStmt(x, walk)
		}
	}
	for _, d := range prog.Declarations {
		walk(d)
	}
	return out
}

func resolveUserFuncCallee(fn parser.Expr, funcs map[string]*parser.FuncDecl) (*parser.FuncDecl, bool) {
	if id, ok := fn.(*parser.IdentifierExpr); ok {
		fd, ok := funcs[id.Name.Lexeme]
		return fd, ok
	}
	return nil, false
}

func promoteParamKindsFromCall(fd *parser.FuncDecl, params []parser.Param, call *parser.CallExpr, letKinds map[*parser.LetDecl]NumericKind, prog *parser.Program, out map[ParamCellKey]NumericKind) {
	if fd == nil || call == nil {
		return
	}
	for i, arg := range call.Arguments {
		if i >= len(params) {
			break
		}
		if params[i].TypeAnnot != "" {
			continue
		}
		key := NewParamCellKey(fd, i)
		if out[key] == KindFloat {
			continue
		}
		if inferArgNumericKind(arg, letKinds, prog) == KindFloat {
			out[key] = KindFloat
		}
	}
}

func collectStructMethodsFromProgram(prog *parser.Program) map[string]map[string]*parser.FuncDecl {
	out := make(map[string]map[string]*parser.FuncDecl)
	if prog == nil {
		return out
	}
	for _, d := range prog.Declarations {
		sd, ok := d.(*parser.StructDecl)
		if !ok {
			continue
		}
		stName := sd.Name.Lexeme
		methods := make(map[string]*parser.FuncDecl)
		for _, m := range sd.Methods {
			methods[m.Name.Lexeme] = m
		}
		out[stName] = methods
	}
	return out
}

func collectVarStructTypesFromProgram(prog *parser.Program) map[string]string {
	out := make(map[string]string)
	if prog == nil {
		return out
	}
	structNames := make(map[string]bool)
	for st, methods := range collectStructMethodsFromProgram(prog) {
		structNames[strings.ToLower(st)] = true
		_ = methods
	}
	var walk func(parser.Decl)
	walk = func(d parser.Decl) {
		switch x := d.(type) {
		case *parser.LetDecl:
			if st := structTypeFromInitExpr(x.Init, structNames); st != "" {
				out[x.Name.Lexeme] = st
			}
		case *parser.FuncDecl:
			if x.Body != nil {
				for _, inner := range x.Body.Declarations {
					walk(inner)
				}
			}
		case *parser.StructDecl:
			for _, m := range x.Methods {
				if m.Body != nil {
					for _, inner := range m.Body.Declarations {
						walk(inner)
					}
				}
			}
		case *parser.BlockStmt:
			for _, inner := range x.Declarations {
				walk(inner)
			}
		case parser.Stmt:
			walkParamsInStmt(x, walk)
		}
	}
	for _, d := range prog.Declarations {
		walk(d)
	}
	return out
}

func structTypeFromInitExpr(e parser.Expr, structNames map[string]bool) string {
	if e == nil {
		return ""
	}
	switch x := e.(type) {
	case *parser.ObjectExpr:
		if x.StructTag != nil {
			return x.StructTag.Lexeme
		}
	case *parser.CallExpr:
		if id, ok := x.Function.(*parser.IdentifierExpr); ok {
			name := id.Name.Lexeme
			if structNames[strings.ToLower(name)] {
				return name
			}
		}
	}
	return ""
}

func resolveStructMethodCallee(call *parser.CallExpr, structMethods map[string]map[string]*parser.FuncDecl, varStructs map[string]string) (*parser.FuncDecl, bool) {
	if call == nil {
		return nil, false
	}
	ix, ok := call.Function.(*parser.IndexExpr)
	if !ok {
		return nil, false
	}
	methodName := methodNameFromIndexExpr(ix)
	if methodName == "" {
		return nil, false
	}
	id, ok := ix.Object.(*parser.IdentifierExpr)
	if !ok {
		return nil, false
	}
	stName, ok := varStructs[id.Name.Lexeme]
	if !ok {
		for k, v := range varStructs {
			if strings.EqualFold(k, id.Name.Lexeme) {
				stName = v
				ok = true
				break
			}
		}
	}
	if !ok {
		return nil, false
	}
	methods := structMethodsForName(structMethods, stName)
	if methods == nil {
		return nil, false
	}
	if fd, ok := methods[methodName]; ok {
		return fd, true
	}
	for k, fd := range methods {
		if strings.EqualFold(k, methodName) {
			return fd, true
		}
	}
	return nil, false
}

func structMethodsForName(structMethods map[string]map[string]*parser.FuncDecl, stName string) map[string]*parser.FuncDecl {
	if structMethods == nil || stName == "" {
		return nil
	}
	if methods, ok := structMethods[stName]; ok {
		return methods
	}
	lower := strings.ToLower(stName)
	for k, v := range structMethods {
		if strings.ToLower(k) == lower {
			return v
		}
	}
	return nil
}

func methodNameFromIndexExpr(ix *parser.IndexExpr) string {
	if ix == nil {
		return ""
	}
	switch idx := ix.Index.(type) {
	case *parser.IdentifierExpr:
		return idx.Name.Lexeme
	case *parser.LiteralExpr:
		if s, ok := idx.Value.(string); ok {
			return s
		}
	}
	return ""
}

func inferArgNumericKind(e parser.Expr, letKinds map[*parser.LetDecl]NumericKind, prog *parser.Program) NumericKind {
	if e == nil {
		return KindInt
	}
	switch x := e.(type) {
	case *parser.LiteralExpr:
		switch x.Value.(type) {
		case int:
			return KindInt
		case float64:
			return KindFloat
		default:
			return KindInt
		}
	case *parser.IdentifierExpr:
		if ld, ok := resolveLetInProgram(prog, x.Name.Lexeme); ok && letKinds != nil {
			if k, ok2 := letKinds[ld]; ok2 {
				return k
			}
		}
		return KindInt
	case *parser.CallExpr:
		if calleeReturnsFloat(x.Function) {
			return KindFloat
		}
		return KindInt
	case *parser.InfixExpr:
		lk := inferArgNumericKind(x.Left, letKinds, prog)
		rk := inferArgNumericKind(x.Right, letKinds, prog)
		if x.Operator == "/" {
			return KindFloat
		}
		if lk == KindFloat || rk == KindFloat {
			return KindFloat
		}
		return KindInt
	case *parser.LogicalExpr:
		lk := inferArgNumericKind(x.Left, letKinds, prog)
		rk := inferArgNumericKind(x.Right, letKinds, prog)
		if lk == KindFloat || rk == KindFloat {
			return KindFloat
		}
		return KindInt
	case *parser.PrefixExpr:
		if x.Operator == "+" || x.Operator == "-" {
			return inferArgNumericKind(x.Right, letKinds, prog)
		}
		return KindInt
	case *parser.GroupingExpr:
		return inferArgNumericKind(x.Expr, letKinds, prog)
	default:
		return KindInt
	}
}

func calleeReturnsFloat(fn parser.Expr) bool {
	name := calleeBaseName(fn)
	if name == "" {
		return false
	}
	return knownFloatReturningNames[strings.ToLower(name)]
}

func calleeBaseName(fn parser.Expr) string {
	switch x := fn.(type) {
	case *parser.IdentifierExpr:
		return x.Name.Lexeme
	case *parser.IndexExpr:
		if id, ok := x.Index.(*parser.IdentifierExpr); ok {
			return id.Name.Lexeme
		}
		if lit, ok := x.Index.(*parser.LiteralExpr); ok {
			if s, ok := lit.Value.(string); ok {
				return s
			}
		}
	}
	return ""
}

// inferFuncExprParamsInProgram walks nested lambdas for body-based param inference.
func inferFuncExprParamsInProgram(prog *parser.Program, out map[ParamCellKey]NumericKind) {
	if prog == nil {
		return
	}
	var scanExpr func(parser.Expr)
	scanExpr = func(e parser.Expr) {
		if e == nil {
			return
		}
		switch x := e.(type) {
		case *parser.FuncExpr:
			recordInferredParamKinds(x, x.Params, x.Body, out)
			scanParamsInBlock(x.Body, scanExpr)
		case *parser.CallExpr:
			scanExpr(x.Function)
			for _, a := range x.Arguments {
				scanExpr(a)
			}
		case *parser.InfixExpr:
			scanExpr(x.Left)
			scanExpr(x.Right)
		case *parser.PrefixExpr:
			scanExpr(x.Right)
		case *parser.AssignExpr:
			scanExpr(x.Left)
			scanExpr(x.Value)
		case *parser.LogicalExpr:
			scanExpr(x.Left)
			scanExpr(x.Right)
		case *parser.IndexExpr:
			scanExpr(x.Object)
			scanExpr(x.Index)
		case *parser.GroupingExpr:
			scanExpr(x.Expr)
		case *parser.ArrayExpr:
			for _, el := range x.Elements {
				scanExpr(el)
			}
		case *parser.IfExpr:
			scanExpr(x.Condition)
			scanExpr(x.Then)
			if x.Else != nil {
				scanExpr(x.Else)
			}
		}
	}
	var walkDecl func(parser.Decl)
	walkDecl = func(d parser.Decl) {
		switch x := d.(type) {
		case *parser.LetDecl:
			scanExpr(x.Init)
		case *parser.FuncDecl:
			scanParamsInBlock(x.Body, scanExpr)
		case *parser.StructDecl:
			for _, m := range x.Methods {
				scanParamsInBlock(m.Body, scanExpr)
			}
		case *parser.BlockStmt:
			for _, inner := range x.Declarations {
				walkDecl(inner)
			}
		case parser.Stmt:
			walkParamsInStmtForScan(x, walkDecl, scanExpr)
		}
	}
	for _, d := range prog.Declarations {
		walkDecl(d)
	}
}
