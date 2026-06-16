package sema

import (
	"koda/internal/parser"
)

// shadowExprTempSlotCount is reserved at the end of every shadow table for
// expression temporaries (&& / || / switch-expr) that hold GC values across
// possible allocation points. See shadowStoreTemp in codegen.
const shadowExprTempSlotCount = 8

// ShadowLayout assigns a stable index in the per-function shadow slot table for
// each root slot (this, parameters, closure captures, lets, for-of bindings).
type ShadowLayout struct {
	Total        int
	TempBase     int // first index of reserved temp slots (Total = TempBase + shadowExprTempSlotCount)
	LetIndex     map[*parser.LetDecl]int
	ForOfIndex   map[*parser.ForOfStmt]int
	FreeVarIndex map[string]int
}

func newShadowLayout() *ShadowLayout {
	return &ShadowLayout{
		LetIndex:     make(map[*parser.LetDecl]int),
		ForOfIndex:   make(map[*parser.ForOfStmt]int),
		FreeVarIndex: make(map[string]int),
	}
}

func finalizeShadowTemps(L *ShadowLayout, next int) {
	L.TempBase = next
	L.Total = next + shadowExprTempSlotCount
}

func prepareShadowLayouts(ctx *NativeEmitContext, bundle *parser.ProgramBundle) {
	ctx.ShadowFuncDecl = make(map[*parser.FuncDecl]*ShadowLayout)
	ctx.ShadowFuncExpr = make(map[*parser.FuncExpr]*ShadowLayout)
	shadowVisitProgram(bundle.Entry, ctx)
	for _, mod := range bundle.Modules {
		if mod != nil {
			shadowVisitProgram(mod, ctx)
		}
	}
	ctx.ShadowEntry = shadowLayoutEntry(bundle.Entry.Declarations)
}

func shadowLayoutEntry(decls []parser.Decl) *ShadowLayout {
	L := newShadowLayout()
	next := 1 // index 0 = this; user_main has no parameters
	for _, d := range decls {
		shadowWalkDecl(d, L, &next)
	}
	finalizeShadowTemps(L, next)
	return L
}

func shadowVisitProgram(p *parser.Program, ctx *NativeEmitContext) {
	if p == nil {
		return
	}
	for _, d := range p.Declarations {
		shadowRegisterAllFuncDecls(d, ctx)
	}
}

func shadowRegisterAllFuncDecls(d parser.Decl, ctx *NativeEmitContext) {
	switch x := d.(type) {
	case *parser.FuncDecl:
		if x.Native != nil {
			return
		}
		if _, ok := ctx.ShadowFuncDecl[x]; !ok {
			ctx.ShadowFuncDecl[x] = shadowLayoutFuncDecl(x)
		}
		shadowRegisterInBlock(x.Body, ctx)
	case *parser.TestDecl:
		fd := x.SyntheticFunc()
		if _, ok := ctx.ShadowFuncDecl[fd]; !ok {
			ctx.ShadowFuncDecl[fd] = shadowLayoutFuncDecl(fd)
		}
		shadowRegisterInBlock(x.Body, ctx)
	case *parser.LetDecl:
		if x.Init != nil {
			shadowScanExprForFuncExprs(x.Init, ctx)
		}
	case *parser.BlockStmt:
		shadowRegisterInBlock(x, ctx)
	case parser.Stmt:
		shadowRegisterInStmt(x, ctx)
	}
}

func shadowRegisterInBlock(b *parser.BlockStmt, ctx *NativeEmitContext) {
	if b == nil {
		return
	}
	for _, d := range b.Declarations {
		shadowRegisterAllFuncDecls(d, ctx)
	}
}

func shadowRegisterInStmt(s parser.Stmt, ctx *NativeEmitContext) {
	switch st := s.(type) {
	case *parser.BlockStmt:
		shadowRegisterInBlock(st, ctx)
	case *parser.IfStmt:
		shadowScanExprForFuncExprs(st.Condition, ctx)
		shadowRegisterInStmt(st.Then, ctx)
		if st.Else != nil {
			shadowRegisterInStmt(st.Else, ctx)
		}
	case *parser.WhileStmt:
		shadowScanExprForFuncExprs(st.Condition, ctx)
		shadowRegisterInStmt(st.Body, ctx)
	case *parser.DoWhileStmt:
		shadowRegisterInStmt(st.Body, ctx)
		shadowScanExprForFuncExprs(st.Condition, ctx)
	case *parser.ForStmt:
		for _, ini := range st.Inits {
			shadowRegisterAllFuncDecls(ini, ctx)
		}
		if st.Condition != nil {
			shadowScanExprForFuncExprs(st.Condition, ctx)
		}
		for _, inc := range st.Increments {
			shadowScanExprForFuncExprs(inc, ctx)
		}
		shadowRegisterInStmt(st.Body, ctx)
	case *parser.ForOfStmt:
		shadowScanExprForFuncExprs(st.Iterable, ctx)
		shadowRegisterInStmt(st.Body, ctx)
	case *parser.ForInStmt:
		shadowScanExprForFuncExprs(st.Iterable, ctx)
		shadowRegisterInStmt(st.Body, ctx)
	case *parser.SwitchStmt:
		shadowScanExprForFuncExprs(st.Subject, ctx)
		for _, c := range st.Cases {
			shadowScanExprForFuncExprs(c.Value, ctx)
			for _, cd := range c.Body {
				shadowRegisterAllFuncDecls(cd, ctx)
			}
		}
		for _, cd := range st.Default {
			shadowRegisterAllFuncDecls(cd, ctx)
		}
	case *parser.ExpressionStmt:
		shadowScanExprForFuncExprs(st.Expr, ctx)
	case *parser.DeferStmt:
		shadowScanExprForFuncExprs(st.Expr, ctx)
	case *parser.DeleteStmt:
		shadowScanExprForFuncExprs(st.Target, ctx)
	case *parser.ReturnStmt:
		if st.Value != nil {
			shadowScanExprForFuncExprs(st.Value, ctx)
		}
	case *parser.BreakStmt, *parser.ContinueStmt, *parser.FallthroughStmt:
	default:
	}
}

func shadowScanDeclForFuncExprs(d parser.Decl, ctx *NativeEmitContext) {
	switch x := d.(type) {
	case *parser.LetDecl:
		if x.Init != nil {
			shadowScanExprForFuncExprs(x.Init, ctx)
		}
	case *parser.FuncDecl:
		if x.Native != nil {
			return
		}
		shadowScanBlockForFuncExprs(x.Body, ctx)
	case *parser.BlockStmt:
		shadowScanBlockForFuncExprs(x, ctx)
	case parser.Stmt:
		shadowScanStmtForFuncExprs(x, ctx)
	}
}

func shadowScanBlockForFuncExprs(b *parser.BlockStmt, ctx *NativeEmitContext) {
	if b == nil {
		return
	}
	for _, d := range b.Declarations {
		shadowScanDeclForFuncExprs(d, ctx)
	}
}

func shadowScanStmtForFuncExprs(s parser.Stmt, ctx *NativeEmitContext) {
	switch st := s.(type) {
	case *parser.BlockStmt:
		shadowScanBlockForFuncExprs(st, ctx)
	case *parser.IfStmt:
		shadowScanStmtForFuncExprs(st.Then, ctx)
		if st.Else != nil {
			shadowScanStmtForFuncExprs(st.Else, ctx)
		}
	case *parser.WhileStmt:
		shadowScanStmtForFuncExprs(st.Body, ctx)
	case *parser.DoWhileStmt:
		shadowScanStmtForFuncExprs(st.Body, ctx)
	case *parser.ForStmt:
		for _, ini := range st.Inits {
			shadowScanDeclForFuncExprs(ini, ctx)
		}
		shadowScanStmtForFuncExprs(st.Body, ctx)
	case *parser.ForOfStmt:
		shadowScanStmtForFuncExprs(st.Body, ctx)
	case *parser.ForInStmt:
		shadowScanStmtForFuncExprs(st.Body, ctx)
	case *parser.SwitchStmt:
		for _, c := range st.Cases {
			for _, cd := range c.Body {
				shadowScanDeclForFuncExprs(cd, ctx)
			}
		}
		for _, cd := range st.Default {
			shadowScanDeclForFuncExprs(cd, ctx)
		}
	case *parser.ExpressionStmt:
		shadowScanExprForFuncExprs(st.Expr, ctx)
	case *parser.DeferStmt:
		shadowScanExprForFuncExprs(st.Expr, ctx)
	}
}

func shadowScanExprForFuncExprs(e parser.Expr, ctx *NativeEmitContext) {
	switch x := e.(type) {
	case *parser.FuncExpr:
		if _, ok := ctx.ShadowFuncExpr[x]; !ok {
			ctx.ShadowFuncExpr[x] = shadowLayoutFuncExpr(x, ctx)
		}
		shadowScanBlockForFuncExprs(x.Body, ctx)
	case *parser.CallExpr:
		shadowScanExprForFuncExprs(x.Function, ctx)
		for _, a := range x.Arguments {
			shadowScanExprForFuncExprs(a, ctx)
		}
	case *parser.InfixExpr:
		shadowScanExprForFuncExprs(x.Left, ctx)
		shadowScanExprForFuncExprs(x.Right, ctx)
	case *parser.PrefixExpr:
		shadowScanExprForFuncExprs(x.Right, ctx)
	case *parser.LogicalExpr:
		shadowScanExprForFuncExprs(x.Left, ctx)
		shadowScanExprForFuncExprs(x.Right, ctx)
	case *parser.AssignExpr:
		shadowScanExprForFuncExprs(x.Left, ctx)
		shadowScanExprForFuncExprs(x.Value, ctx)
	case *parser.GroupingExpr:
		shadowScanExprForFuncExprs(x.Expr, ctx)
	case *parser.IndexExpr:
		shadowScanExprForFuncExprs(x.Object, ctx)
		shadowScanExprForFuncExprs(x.Index, ctx)
	case *parser.ArrayExpr:
		for _, el := range x.Elements {
			shadowScanExprForFuncExprs(el, ctx)
		}
	case *parser.SpreadExpr:
		shadowScanExprForFuncExprs(x.Expr, ctx)
	case *parser.ObjectExpr:
		for _, v := range x.Values {
			shadowScanExprForFuncExprs(v, ctx)
		}
		for _, ck := range x.ComputedKeys {
			shadowScanExprForFuncExprs(ck, ctx)
		}
	case *parser.TemplateExpr:
		for _, p := range x.Parts {
			shadowScanExprForFuncExprs(p, ctx)
		}
	case *parser.RangeExpr:
		shadowScanExprForFuncExprs(x.From, ctx)
		shadowScanExprForFuncExprs(x.To, ctx)
	case *parser.SwitchExpr:
		shadowScanExprForFuncExprs(x.Subject, ctx)
		for _, c := range x.Cases {
			shadowScanExprForFuncExprs(c.Value, ctx)
			shadowScanExprForFuncExprs(c.Body, ctx)
		}
		if x.Default != nil {
			shadowScanExprForFuncExprs(x.Default, ctx)
		}
	case *parser.TupleExpr:
		for _, el := range x.Elements {
			shadowScanExprForFuncExprs(el, ctx)
		}
	case *parser.IfExpr:
		shadowScanExprForFuncExprs(x.Condition, ctx)
		shadowScanExprForFuncExprs(x.Then, ctx)
		if x.Else != nil {
			shadowScanExprForFuncExprs(x.Else, ctx)
		}
	case *parser.SliceExpr:
		shadowScanExprForFuncExprs(x.Object, ctx)
		if x.Start != nil {
			shadowScanExprForFuncExprs(x.Start, ctx)
		}
		if x.End != nil {
			shadowScanExprForFuncExprs(x.End, ctx)
		}
	case *parser.TernaryExpr:
		shadowScanExprForFuncExprs(x.Condition, ctx)
		shadowScanExprForFuncExprs(x.Then, ctx)
		shadowScanExprForFuncExprs(x.Else, ctx)
	case *parser.UpdateExpr:
		shadowScanExprForFuncExprs(x.Operand, ctx)
	}
}

func shadowLayoutFuncDecl(fd *parser.FuncDecl) *ShadowLayout {
	L := newShadowLayout()
	next := 1 + len(fd.Params)
	shadowWalkBlock(fd.Body, L, &next)
	finalizeShadowTemps(L, next)
	return L
}

func shadowLayoutFuncExpr(fe *parser.FuncExpr, ctx *NativeEmitContext) *ShadowLayout {
	L := newShadowLayout()
	next := 1 + len(fe.Params)
	for _, name := range ctx.FreeVarsExpr[fe] {
		L.FreeVarIndex[name] = next
		next++
	}
	shadowWalkBlock(fe.Body, L, &next)
	finalizeShadowTemps(L, next)
	return L
}

func shadowWalkDecl(d parser.Decl, L *ShadowLayout, next *int) {
	switch x := d.(type) {
	case *parser.LetDecl:
		if L != nil && next != nil {
			L.LetIndex[x] = *next
			*next++
		}
		if x.Init != nil {
			shadowWalkExprLets(x.Init, L, next)
		}
	case *parser.FuncDecl:
		return
	case *parser.TestDecl:
		return
	case *parser.IncludeDecl:
		return
	case *parser.BlockStmt:
		shadowWalkBlock(x, L, next)
	case parser.Stmt:
		shadowWalkStmt(x, L, next)
	}
}

func shadowWalkBlock(b *parser.BlockStmt, L *ShadowLayout, next *int) {
	if b == nil {
		return
	}
	for _, d := range b.Declarations {
		shadowWalkDecl(d, L, next)
	}
}

func shadowWalkStmt(s parser.Stmt, L *ShadowLayout, next *int) {
	switch st := s.(type) {
	case *parser.BlockStmt:
		shadowWalkBlock(st, L, next)
	case *parser.IfStmt:
		shadowWalkStmt(st.Then, L, next)
		if st.Else != nil {
			shadowWalkStmt(st.Else, L, next)
		}
	case *parser.WhileStmt:
		shadowWalkStmt(st.Body, L, next)
	case *parser.DoWhileStmt:
		shadowWalkStmt(st.Body, L, next)
	case *parser.ForStmt:
		for _, ini := range st.Inits {
			shadowWalkDecl(ini, L, next)
		}
		shadowWalkStmt(st.Body, L, next)
	case *parser.ForOfStmt:
		if L != nil && next != nil {
			L.ForOfIndex[st] = *next
			*next++
			if st.ValueVar != nil {
				*next++
			}
		}
		shadowWalkStmt(st.Body, L, next)
	case *parser.ForInStmt:
		shadowWalkStmt(st.Body, L, next)
	case *parser.SwitchStmt:
		for _, c := range st.Cases {
			for _, cd := range c.Body {
				shadowWalkDecl(cd, L, next)
			}
		}
		for _, cd := range st.Default {
			shadowWalkDecl(cd, L, next)
		}
	case *parser.ExpressionStmt:
		shadowWalkExprLets(st.Expr, L, next)
	case *parser.DeferStmt:
		shadowWalkExprLets(st.Expr, L, next)
	}
}

func shadowWalkExprLets(e parser.Expr, L *ShadowLayout, next *int) {
	switch x := e.(type) {
	case *parser.FuncExpr:
		shadowWalkBlock(x.Body, L, next)
	case *parser.CallExpr:
		for _, a := range x.Arguments {
			shadowWalkExprLets(a, L, next)
		}
		shadowWalkExprLets(x.Function, L, next)
	case *parser.InfixExpr:
		shadowWalkExprLets(x.Left, L, next)
		shadowWalkExprLets(x.Right, L, next)
	case *parser.PrefixExpr:
		shadowWalkExprLets(x.Right, L, next)
	case *parser.LogicalExpr:
		shadowWalkExprLets(x.Left, L, next)
		shadowWalkExprLets(x.Right, L, next)
	case *parser.AssignExpr:
		shadowWalkExprLets(x.Left, L, next)
		shadowWalkExprLets(x.Value, L, next)
	case *parser.IndexExpr:
		shadowWalkExprLets(x.Object, L, next)
		shadowWalkExprLets(x.Index, L, next)
	case *parser.ArrayExpr:
		for _, el := range x.Elements {
			shadowWalkExprLets(el, L, next)
		}
	case *parser.SpreadExpr:
		shadowWalkExprLets(x.Expr, L, next)
	case *parser.ObjectExpr:
		for _, v := range x.Values {
			shadowWalkExprLets(v, L, next)
		}
		for _, ck := range x.ComputedKeys {
			shadowWalkExprLets(ck, L, next)
		}
	case *parser.TemplateExpr:
		for _, p := range x.Parts {
			shadowWalkExprLets(p, L, next)
		}
	case *parser.RangeExpr:
		shadowWalkExprLets(x.From, L, next)
		shadowWalkExprLets(x.To, L, next)
	case *parser.SwitchExpr:
		shadowWalkExprLets(x.Subject, L, next)
		for _, c := range x.Cases {
			shadowWalkExprLets(c.Value, L, next)
			shadowWalkExprLets(c.Body, L, next)
		}
		if x.Default != nil {
			shadowWalkExprLets(x.Default, L, next)
		}
	case *parser.GroupingExpr:
		shadowWalkExprLets(x.Expr, L, next)
	case *parser.IfExpr:
		shadowWalkExprLets(x.Condition, L, next)
		shadowWalkExprLets(x.Then, L, next)
		if x.Else != nil {
			shadowWalkExprLets(x.Else, L, next)
		}
	case *parser.SliceExpr:
		shadowWalkExprLets(x.Object, L, next)
		if x.Start != nil {
			shadowWalkExprLets(x.Start, L, next)
		}
		if x.End != nil {
			shadowWalkExprLets(x.End, L, next)
		}
	case *parser.TernaryExpr:
		shadowWalkExprLets(x.Condition, L, next)
		shadowWalkExprLets(x.Then, L, next)
		shadowWalkExprLets(x.Else, L, next)
	case *parser.UpdateExpr:
		shadowWalkExprLets(x.Operand, L, next)
	}
}
