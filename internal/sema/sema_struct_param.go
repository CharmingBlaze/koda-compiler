package sema

import "koda/internal/parser"

// recordParamStructFromCall records struct types passed to function parameters at call sites.
func (a *Analyzer) recordParamStructFromCall(call *parser.CallExpr) {
	id, ok := call.Function.(*parser.IdentifierExpr)
	if !ok {
		return
	}
	decl, ok := a.currentScope.Resolve(id.Name.Lexeme)
	if !ok {
		return
	}
	fd, ok := decl.(*parser.FuncDecl)
	if !ok {
		return
	}
	fn := fd.Name.Lexeme
	for i, arg := range call.Arguments {
		if i >= len(fd.Params) {
			break
		}
		st := a.structTypeOfExpr(arg)
		if st == "" {
			continue
		}
		if a.funcParamStruct == nil {
			a.funcParamStruct = make(map[string]map[string]string)
		}
		if a.funcParamStruct[fn] == nil {
			a.funcParamStruct[fn] = make(map[string]string)
		}
		pname := fd.Params[i].Name
		if _, exists := a.funcParamStruct[fn][pname]; !exists {
			a.funcParamStruct[fn][pname] = st
		}
	}
}

func (a *Analyzer) structTypeOfExpr(e parser.Expr) string {
	switch x := e.(type) {
	case *parser.IdentifierExpr:
		if st, ok := a.varStructType[x.Name.Lexeme]; ok {
			return st
		}
	case *parser.CallExpr:
		if id, ok := x.Function.(*parser.IdentifierExpr); ok {
			if st := a.structConstructorType(id.Name.Lexeme); st != "" {
				return st
			}
		}
	case *parser.ObjectExpr:
		if x.StructTag != nil {
			return x.StructTag.Lexeme
		}
	case *parser.IndexExpr:
		if slot, ok := a.indexExprStructSlot[x]; ok {
			if st := a.structTypeFromFieldSlot(x, slot); st != "" {
				return st
			}
		}
	}
	return ""
}

func (a *Analyzer) structTypeFromFieldSlot(ix *parser.IndexExpr, slot int) string {
	idObj, ok := ix.Object.(*parser.IdentifierExpr)
	if !ok {
		return ""
	}
	stName, ok := a.varStructType[idObj.Name.Lexeme]
	if !ok && a.activeParamStruct != nil {
		stName, ok = a.activeParamStruct[idObj.Name.Lexeme]
	}
	if !ok {
		return ""
	}
	fields, ok := a.structLayouts[stName]
	if !ok || slot < 0 || slot >= len(fields) {
		return ""
	}
	return stName
}

// refineStructFieldAccessFromCalls re-runs struct field slot binding for function parameters
// after all call sites have been analyzed.
func (a *Analyzer) refineStructFieldAccessFromCalls(prog *parser.Program) {
	if len(a.funcParamStruct) == 0 {
		return
	}
	for _, decl := range prog.Declarations {
		fd, ok := decl.(*parser.FuncDecl)
		if !ok || fd.Body == nil {
			continue
		}
		prev := a.activeParamStruct
		a.activeParamStruct = a.funcParamStruct[fd.Name.Lexeme]
		a.walkStmtStructFieldAccess(fd.Body)
		a.activeParamStruct = prev
	}
}

func (a *Analyzer) walkDeclStructFieldAccess(decl parser.Decl) {
	switch d := decl.(type) {
	case *parser.FuncDecl:
		prev := a.activeParamStruct
		a.activeParamStruct = a.funcParamStruct[d.Name.Lexeme]
		if d.Body != nil {
			a.walkStmtStructFieldAccess(d.Body)
		}
		a.activeParamStruct = prev
	case *parser.FuncExpr:
		a.walkStmtStructFieldAccess(d.Body)
	case *parser.LetDecl:
		if d.Init != nil {
			a.walkExprStructFieldAccess(d.Init)
		}
	case parser.Stmt:
		a.walkStmtStructFieldAccess(d)
	}
}

func (a *Analyzer) walkStmtStructFieldAccess(stmt parser.Stmt) {
	switch s := stmt.(type) {
	case *parser.BlockStmt:
		for _, decl := range s.Declarations {
			a.walkDeclStructFieldAccess(decl)
		}
	case *parser.ExpressionStmt:
		a.walkExprStructFieldAccess(s.Expr)
	case *parser.ReturnStmt:
		if s.Value != nil {
			a.walkExprStructFieldAccess(s.Value)
		}
	case *parser.IfStmt:
		a.walkExprStructFieldAccess(s.Condition)
		a.walkStmtStructFieldAccess(s.Then)
		if s.Else != nil {
			a.walkStmtStructFieldAccess(s.Else)
		}
	case *parser.WhileStmt:
		a.walkExprStructFieldAccess(s.Condition)
		a.walkStmtStructFieldAccess(s.Body)
	case *parser.DoWhileStmt:
		a.walkStmtStructFieldAccess(s.Body)
		a.walkExprStructFieldAccess(s.Condition)
	case *parser.ForStmt:
		for _, ini := range s.Inits {
			a.walkDeclStructFieldAccess(ini)
		}
		if s.Condition != nil {
			a.walkExprStructFieldAccess(s.Condition)
		}
		for _, inc := range s.Increments {
			a.walkExprStructFieldAccess(inc)
		}
		a.walkStmtStructFieldAccess(s.Body)
	case *parser.ForInStmt:
		a.walkExprStructFieldAccess(s.Iterable)
		a.walkStmtStructFieldAccess(s.Body)
	case *parser.SwitchStmt:
		a.walkExprStructFieldAccess(s.Subject)
		for _, c := range s.Cases {
			a.walkExprStructFieldAccess(c.Value)
			for _, decl := range c.Body {
				a.walkDeclStructFieldAccess(decl)
			}
		}
		for _, decl := range s.Default {
			a.walkDeclStructFieldAccess(decl)
		}
	case *parser.DeferStmt:
		a.walkExprStructFieldAccess(s.Expr)
	}
}

func (a *Analyzer) walkExprStructFieldAccess(e parser.Expr) {
	if e == nil {
		return
	}
	switch x := e.(type) {
	case *parser.IndexExpr:
		a.walkExprStructFieldAccess(x.Object)
		a.walkExprStructFieldAccess(x.Index)
		a.checkStructFieldAccess(x)
	case *parser.AssignExpr:
		a.walkExprStructFieldAccess(x.Value)
		if ix, ok := x.Left.(*parser.IndexExpr); ok {
			a.walkExprStructFieldAccess(ix.Object)
			a.walkExprStructFieldAccess(ix.Index)
			a.checkStructFieldAccess(ix)
		}
	case *parser.InfixExpr:
		a.walkExprStructFieldAccess(x.Left)
		a.walkExprStructFieldAccess(x.Right)
	case *parser.LogicalExpr:
		a.walkExprStructFieldAccess(x.Left)
		a.walkExprStructFieldAccess(x.Right)
	case *parser.CallExpr:
		a.walkExprStructFieldAccess(x.Function)
		for _, arg := range x.Arguments {
			a.walkExprStructFieldAccess(arg)
		}
	case *parser.PrefixExpr:
		a.walkExprStructFieldAccess(x.Right)
	case *parser.GroupingExpr:
		a.walkExprStructFieldAccess(x.Expr)
	case *parser.ObjectExpr:
		for _, v := range x.Values {
			a.walkExprStructFieldAccess(v)
		}
	case *parser.ArrayExpr:
		for _, el := range x.Elements {
			a.walkExprStructFieldAccess(el)
		}
	case *parser.FuncExpr:
		a.walkStmtStructFieldAccess(x.Body)
	case *parser.RangeExpr:
		a.walkExprStructFieldAccess(x.From)
		a.walkExprStructFieldAccess(x.To)
	case *parser.UpdateExpr:
		a.walkExprStructFieldAccess(x.Operand)
	case *parser.TemplateExpr:
		for _, p := range x.Parts {
			a.walkExprStructFieldAccess(p)
		}
	case *parser.TupleExpr:
		for _, el := range x.Elements {
			a.walkExprStructFieldAccess(el)
		}
	case *parser.SpreadExpr:
		a.walkExprStructFieldAccess(x.Expr)
	}
}
