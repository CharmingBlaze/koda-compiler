package sema

import (
	"koda/internal/parser"
)

type paramInfo struct {
	owner interface{} // *parser.FuncDecl or *parser.FuncExpr
	idx   int
	param *parser.Param
}

type envFrame struct {
	parent *envFrame
	lets   map[string]*parser.LetDecl
	params map[string]*paramInfo
}

func (e *envFrame) pushBlock() *envFrame {
	return &envFrame{parent: e, lets: make(map[string]*parser.LetDecl)}
}

func newParamEnvFuncDecl(parent *envFrame, fd *parser.FuncDecl) *envFrame {
	params := make(map[string]*paramInfo)
	for i := range fd.Params {
		p := &fd.Params[i]
		params[p.Name] = &paramInfo{owner: fd, idx: i, param: p}
	}
	return &envFrame{parent: parent, lets: make(map[string]*parser.LetDecl), params: params}
}

func newParamEnvFuncExpr(parent *envFrame, fe *parser.FuncExpr) *envFrame {
	params := make(map[string]*paramInfo)
	for i := range fe.Params {
		p := &fe.Params[i]
		params[p.Name] = &paramInfo{owner: fe, idx: i, param: p}
	}
	return &envFrame{parent: parent, lets: make(map[string]*parser.LetDecl), params: params}
}

func (e *envFrame) lookupLet(name string) *parser.LetDecl {
	for f := e; f != nil; f = f.parent {
		if f.lets != nil {
			if ld, ok := f.lets[name]; ok {
				return ld
			}
		}
	}
	return nil
}

func (e *envFrame) lookupParam(name string) *paramInfo {
	for f := e; f != nil; f = f.parent {
		if f.params != nil {
			if pi, ok := f.params[name]; ok {
				return pi
			}
		}
	}
	return nil
}

func prepareNativeAnalysis(ctx *NativeEmitContext, bundle *parser.ProgramBundle) {
	root := &envFrame{lets: make(map[string]*parser.LetDecl)}
	for _, d := range bundle.Entry.Declarations {
		walkDecl(d, root, ctx, nil)
	}
	finalizeStackDecls(ctx)
}

func finalizeStackDecls(ctx *NativeEmitContext) {
	var walk func(parser.Decl)
	var walkBlock func(*parser.BlockStmt)
	walkBlock = func(b *parser.BlockStmt) {
		if b == nil {
			return
		}
		for _, d := range b.Declarations {
			walk(d)
		}
	}
	walk = func(d parser.Decl) {
		switch x := d.(type) {
		case *parser.LetDecl:
			if ctx.letOwner[x] == nil {
				return
			}
			if ctx.EscapingDecls[x] {
				return
			}
			ctx.StackDecls[x] = true
		case *parser.FuncDecl:
			walkBlock(x.Body)
		case *parser.BlockStmt:
			walkBlock(x)
		case parser.Stmt:
			walkStmt(x, ctx, walkBlock, walk)
		}
	}
	for _, d := range ctx.Bundle.Entry.Declarations {
		walk(d)
	}
}

func walkStmt(s parser.Stmt, ctx *NativeEmitContext, walkBlock func(*parser.BlockStmt), walkDecl func(parser.Decl)) {
	switch st := s.(type) {
	case *parser.BlockStmt:
		walkBlock(st)
	case *parser.DeleteStmt:
		walkExpr(st.Target, nil, ctx, nil)
	case *parser.DeferStmt:
		walkExpr(st.Expr, nil, ctx, nil)
	case *parser.ExpressionStmt:
		walkExpr(st.Expr, nil, ctx, nil)
	case *parser.ReturnStmt:
		if st.Value != nil {
			walkExpr(st.Value, nil, ctx, nil)
		}
	case *parser.IfStmt:
		walkExpr(st.Condition, nil, ctx, nil)
		walkStmt(st.Then, ctx, walkBlock, walkDecl)
		if st.Else != nil {
			walkStmt(st.Else, ctx, walkBlock, walkDecl)
		}
	case *parser.WhileStmt:
		walkExpr(st.Condition, nil, ctx, nil)
		walkStmt(st.Body, ctx, walkBlock, walkDecl)
	case *parser.DoWhileStmt:
		walkStmt(st.Body, ctx, walkBlock, walkDecl)
		walkExpr(st.Condition, nil, ctx, nil)
	case *parser.ForStmt:
		for _, ini := range st.Inits {
			walkDecl(ini)
		}
		if st.Condition != nil {
			walkExpr(st.Condition, nil, ctx, nil)
		}
		for _, inc := range st.Increments {
			walkExpr(inc, nil, ctx, nil)
		}
		walkStmt(st.Body, ctx, walkBlock, walkDecl)
	case *parser.ForInStmt:
		walkExpr(st.Iterable, nil, ctx, nil)
		walkStmt(st.Body, ctx, walkBlock, walkDecl)
	case *parser.ForOfStmt:
		walkExpr(st.Iterable, nil, ctx, nil)
		walkStmt(st.Body, ctx, walkBlock, walkDecl)
	case *parser.SwitchStmt:
		walkExpr(st.Subject, nil, ctx, nil)
		for _, c := range st.Cases {
			walkExpr(c.Value, nil, ctx, nil)
			for _, cd := range c.Body {
				walkDecl(cd)
			}
		}
		for _, cd := range st.Default {
			walkDecl(cd)
		}
	case *parser.BreakStmt, *parser.ContinueStmt:
	default:
	}
}

func walkDecl(d parser.Decl, env *envFrame, ctx *NativeEmitContext, curFunc interface{}) {
	switch x := d.(type) {
	case *parser.StructDecl, *parser.EnumDecl:
	case *parser.LetDecl:
		if x.Init != nil {
			walkExpr(x.Init, env, ctx, curFunc)
		}
		env.lets[x.Name.Lexeme] = x
		ctx.letOwner[x] = curFunc
	case *parser.FuncDecl:
		walkFuncDecl(x, env, ctx)
	case *parser.BlockStmt:
		blockEnv := env.pushBlock()
		for _, inner := range x.Declarations {
			walkDecl(inner, blockEnv, ctx, curFunc)
		}
	case parser.Stmt:
		walkStmtFull(x, env, ctx, curFunc)
	}
}

func walkStmtFull(s parser.Stmt, env *envFrame, ctx *NativeEmitContext, curFunc interface{}) {
	switch st := s.(type) {
	case *parser.BlockStmt:
		blockEnv := env.pushBlock()
		for _, d := range st.Declarations {
			walkDecl(d, blockEnv, ctx, curFunc)
		}
	case *parser.DeleteStmt:
		walkExpr(st.Target, env, ctx, curFunc)
	case *parser.DeferStmt:
		walkExpr(st.Expr, env, ctx, curFunc)
	case *parser.ExpressionStmt:
		walkExpr(st.Expr, env, ctx, curFunc)
	case *parser.ReturnStmt:
		if st.Value != nil {
			walkExpr(st.Value, env, ctx, curFunc)
		}
		checkReturnEscape(st, env, ctx, curFunc)
	case *parser.IfStmt:
		walkExpr(st.Condition, env, ctx, curFunc)
		walkStmtFull(st.Then, env, ctx, curFunc)
		if st.Else != nil {
			walkStmtFull(st.Else, env, ctx, curFunc)
		}
	case *parser.WhileStmt:
		walkExpr(st.Condition, env, ctx, curFunc)
		walkStmtFull(st.Body, env, ctx, curFunc)
	case *parser.DoWhileStmt:
		walkStmtFull(st.Body, env, ctx, curFunc)
		walkExpr(st.Condition, env, ctx, curFunc)
	case *parser.ForStmt:
		for _, ini := range st.Inits {
			walkDecl(ini, env, ctx, curFunc)
		}
		if st.Condition != nil {
			walkExpr(st.Condition, env, ctx, curFunc)
		}
		for _, inc := range st.Increments {
			walkExpr(inc, env, ctx, curFunc)
		}
		walkStmtFull(st.Body, env, ctx, curFunc)
	case *parser.ForInStmt:
		walkExpr(st.Iterable, env, ctx, curFunc)
		walkStmtFull(st.Body, env, ctx, curFunc)
	case *parser.ForOfStmt:
		walkExpr(st.Iterable, env, ctx, curFunc)
		walkStmtFull(st.Body, env, ctx, curFunc)
	case *parser.SwitchStmt:
		walkExpr(st.Subject, env, ctx, curFunc)
		for _, c := range st.Cases {
			walkExpr(c.Value, env, ctx, curFunc)
			blockEnv := env.pushBlock()
			for _, cd := range c.Body {
				walkDecl(cd, blockEnv, ctx, curFunc)
			}
		}
		blockEnv := env.pushBlock()
		for _, cd := range st.Default {
			walkDecl(cd, blockEnv, ctx, curFunc)
		}
	case *parser.BreakStmt, *parser.ContinueStmt:
	default:
	}
}

func checkReturnEscape(ret *parser.ReturnStmt, env *envFrame, ctx *NativeEmitContext, curFunc interface{}) {
	if ret.Value == nil {
		return
	}
	id, ok := ret.Value.(*parser.IdentifierExpr)
	if !ok {
		return
	}
	if ld := env.lookupLet(id.Name.Lexeme); ld != nil {
		if o := ctx.letOwner[ld]; o == curFunc {
			ctx.EscapingDecls[ld] = true
		}
	}
	if pi := env.lookupParam(id.Name.Lexeme); pi != nil && pi.owner == curFunc {
		ctx.ParamIsCell[NewParamCellKey(pi.owner, pi.idx)] = true
	}
}

func walkFuncDecl(fd *parser.FuncDecl, env *envFrame, ctx *NativeEmitContext) {
	paramEnv := newParamEnvFuncDecl(env, fd)
	walkBlockStmt(fd.Body, paramEnv, ctx, fd)
}

func walkBlockStmt(b *parser.BlockStmt, env *envFrame, ctx *NativeEmitContext, curFunc interface{}) {
	if b == nil {
		return
	}
	blockEnv := env.pushBlock()
	for _, d := range b.Declarations {
		walkDecl(d, blockEnv, ctx, curFunc)
	}
}

func visitFuncExpr(fe *parser.FuncExpr, env *envFrame, ctx *NativeEmitContext, enclosing interface{}) {
	if ctx.FuncExprEnclosing != nil {
		ctx.FuncExprEnclosing[fe] = enclosing
	}
	paramEnv := newParamEnvFuncExpr(env, fe)
	walkBlockStmt(fe.Body, paramEnv, ctx, fe)
}

func nextEnclosingFunc(ctx *NativeEmitContext, cur interface{}) interface{} {
	fe, ok := cur.(*parser.FuncExpr)
	if !ok || ctx.FuncExprEnclosing == nil {
		return nil
	}
	enc, ok := ctx.FuncExprEnclosing[fe]
	if !ok {
		return nil
	}
	return enc
}

func appendUniqueName(names []string, name string) []string {
	for _, y := range names {
		if y == name {
			return names
		}
	}
	return append(names, name)
}

func recordFreeVar(ctx *NativeEmitContext, curFunc interface{}, name string) {
	if curFunc == nil {
		return
	}
	switch f := curFunc.(type) {
	case *parser.FuncExpr:
		ctx.FreeVarsExpr[f] = appendUniqueName(ctx.FreeVarsExpr[f], name)
	case *parser.FuncDecl:
		ctx.FreeVarsDecl[f] = appendUniqueName(ctx.FreeVarsDecl[f], name)
	}
}

func markIdentifierCapture(id *parser.IdentifierExpr, env *envFrame, ctx *NativeEmitContext, curFunc interface{}) {
	if ld := env.lookupLet(id.Name.Lexeme); ld != nil {
		owner := ctx.letOwner[ld]
		if owner != nil && owner != curFunc {
			ctx.EscapingDecls[ld] = true
			recordFreeVar(ctx, curFunc, id.Name.Lexeme)
			for x := nextEnclosingFunc(ctx, curFunc); x != nil && x != owner; x = nextEnclosingFunc(ctx, x) {
				recordFreeVar(ctx, x, id.Name.Lexeme)
			}
		}
		return
	}
	if pi := env.lookupParam(id.Name.Lexeme); pi != nil {
		if pi.owner != curFunc {
			ctx.ParamIsCell[NewParamCellKey(pi.owner, pi.idx)] = true
			recordFreeVar(ctx, curFunc, id.Name.Lexeme)
			for x := nextEnclosingFunc(ctx, curFunc); x != nil && x != pi.owner; x = nextEnclosingFunc(ctx, x) {
				recordFreeVar(ctx, x, id.Name.Lexeme)
			}
		}
	}
}

func walkExpr(ex parser.Expr, env *envFrame, ctx *NativeEmitContext, curFunc interface{}) {
	switch e := ex.(type) {
	case *parser.IdentifierExpr:
		if env != nil && ctx != nil {
			markIdentifierCapture(e, env, ctx, curFunc)
		}
	case *parser.FuncExpr:
		visitFuncExpr(e, env, ctx, curFunc)
	case *parser.CallExpr:
		walkExpr(e.Function, env, ctx, curFunc)
		for _, a := range e.Arguments {
			walkExpr(a, env, ctx, curFunc)
		}
	case *parser.AssignExpr:
		walkExpr(e.Left, env, ctx, curFunc)
		walkExpr(e.Value, env, ctx, curFunc)
	case *parser.InfixExpr:
		walkExpr(e.Left, env, ctx, curFunc)
		walkExpr(e.Right, env, ctx, curFunc)
	case *parser.PrefixExpr:
		walkExpr(e.Right, env, ctx, curFunc)
	case *parser.LogicalExpr:
		walkExpr(e.Left, env, ctx, curFunc)
		walkExpr(e.Right, env, ctx, curFunc)
	case *parser.IndexExpr:
		walkExpr(e.Object, env, ctx, curFunc)
		walkExpr(e.Index, env, ctx, curFunc)
	case *parser.SliceExpr:
		walkExpr(e.Object, env, ctx, curFunc)
		if e.Start != nil {
			walkExpr(e.Start, env, ctx, curFunc)
		}
		if e.End != nil {
			walkExpr(e.End, env, ctx, curFunc)
		}
	case *parser.ArrayExpr:
		for _, el := range e.Elements {
			walkExpr(el, env, ctx, curFunc)
		}
	case *parser.SpreadExpr:
		walkExpr(e.Expr, env, ctx, curFunc)
	case *parser.ObjectExpr:
		for _, v := range e.Values {
			walkExpr(v, env, ctx, curFunc)
		}
		for _, k := range e.ComputedKeys {
			walkExpr(k, env, ctx, curFunc)
		}
	case *parser.GroupingExpr:
		walkExpr(e.Expr, env, ctx, curFunc)
	case *parser.LiteralExpr, *parser.ThisExpr:
	case *parser.TemplateExpr:
		for _, p := range e.Parts {
			walkExpr(p, env, ctx, curFunc)
		}
	case *parser.RangeExpr:
		walkExpr(e.From, env, ctx, curFunc)
		walkExpr(e.To, env, ctx, curFunc)
	case *parser.TupleExpr:
		for _, el := range e.Elements {
			walkExpr(el, env, ctx, curFunc)
		}
	case *parser.ImportExpr:
	case *parser.SwitchExpr:
		walkExpr(e.Subject, env, ctx, curFunc)
		for _, c := range e.Cases {
			walkExpr(c.Value, env, ctx, curFunc)
			walkExpr(c.Body, env, ctx, curFunc)
		}
		if e.Default != nil {
			walkExpr(e.Default, env, ctx, curFunc)
		}
	case *parser.UpdateExpr:
		walkExpr(e.Operand, env, ctx, curFunc)
	case *parser.IfExpr:
		walkExpr(e.Condition, env, ctx, curFunc)
		walkExpr(e.Then, env, ctx, curFunc)
		if e.Else != nil {
			walkExpr(e.Else, env, ctx, curFunc)
		}
	case *parser.TernaryExpr:
		walkExpr(e.Condition, env, ctx, curFunc)
		walkExpr(e.Then, env, ctx, curFunc)
		walkExpr(e.Else, env, ctx, curFunc)
	default:
	}
}
