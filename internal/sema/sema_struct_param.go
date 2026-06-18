package sema

import (
	"strings"

	"koda/internal/parser"
)

// recordCallbackStructParamsFromMethodCall records struct element types for parameters of
// anonymous callbacks passed to array methods (.sort, .map, .filter, .each, …).
func (a *Analyzer) recordCallbackStructParamsFromMethodCall(call *parser.CallExpr) {
	ix, ok := call.Function.(*parser.IndexExpr)
	if !ok {
		return
	}
	lit, ok := ix.Index.(*parser.LiteralExpr)
	if !ok {
		return
	}
	method, ok := lit.Value.(string)
	if !ok {
		return
	}
	method = strings.ToLower(method)
	elemStruct := ""
	switch method {
	case "sort", "map", "filter", "each", "foreach", "find", "flatmap":
		elemStruct = a.arrayElementStructType(ix.Object)
	default:
		return
	}
	if elemStruct == "" || len(call.Arguments) == 0 {
		return
	}
	fe, ok := call.Arguments[0].(*parser.FuncExpr)
	if !ok {
		return
	}
	if a.funcExprParamStruct == nil {
		a.funcExprParamStruct = make(map[*parser.FuncExpr]map[string]string)
	}
	params := make(map[string]string)
	switch method {
	case "sort":
		for i := 0; i < len(fe.Params) && i < 2; i++ {
			params[fe.Params[i].Name] = elemStruct
		}
	default:
		if len(fe.Params) > 0 {
			params[fe.Params[0].Name] = elemStruct
		}
	}
	if len(params) > 0 {
		a.funcExprParamStruct[fe] = params
	}
}

func (a *Analyzer) arrayElementStructType(recv parser.Expr) string {
	if id, ok := recv.(*parser.IdentifierExpr); ok {
		name := id.Name.Lexeme
		if st, ok := a.varArrayElementStruct[name]; ok {
			return st
		}
		if a.currentFuncName != "" && a.funcParamArrayElement != nil {
			if params, ok := a.funcParamArrayElement[a.currentFuncName]; ok {
				if st, ok := params[name]; ok {
					return st
				}
			}
		}
	}
	return ""
}

// recordParamArrayElementFromCall records array element struct types for parameters from
// call sites like update_playing(dt, coins) so `for coin in coins` binds Coin fields.
func (a *Analyzer) recordParamArrayElementFromCall(call *parser.CallExpr) {
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
		pname := fd.Params[i].Name
		elemSt := ""
		if ae, ok := arg.(*parser.ArrayExpr); ok && len(ae.Elements) > 0 {
			elemSt = a.structTypeOfExpr(ae.Elements[0])
		}
		if elemSt == "" {
			if argID, ok := arg.(*parser.IdentifierExpr); ok {
				if st, ok := a.varArrayElementStruct[argID.Name.Lexeme]; ok {
					elemSt = st
				}
			}
		}
		if elemSt == "" {
			continue
		}
		if a.funcParamArrayElement == nil {
			a.funcParamArrayElement = make(map[string]map[string]string)
		}
		if a.funcParamArrayElement[fn] == nil {
			a.funcParamArrayElement[fn] = make(map[string]string)
		}
		a.funcParamArrayElement[fn][pname] = elemSt
	}
}

// refineStructFieldAccessFromCallbacks re-binds struct field slots inside anonymous callbacks
// after array element types are known from the receiver expression.
func (a *Analyzer) refineStructFieldAccessFromCallbacks() {
	if len(a.funcExprParamStruct) == 0 {
		return
	}
	for fe, params := range a.funcExprParamStruct {
		prev := a.activeParamStruct
		a.activeParamStruct = params
		a.walkStmtStructFieldAccess(fe.Body)
		a.activeParamStruct = prev
	}
}

func (a *Analyzer) isPlainObjectArg(arg parser.Expr) bool {
	switch x := arg.(type) {
	case *parser.ObjectExpr:
		return x.StructTag == nil
	case *parser.IdentifierExpr:
		return a.varPlainObject[x.Name.Lexeme]
	default:
		return false
	}
}
// recordStructMethodParamFromCall records struct types for method parameters from
// call sites like coin.pickup(player) so mario.x inside pickup() uses struct slots.
func (a *Analyzer) recordStructMethodParamFromCall(call *parser.CallExpr) {
	ix, ok := call.Function.(*parser.IndexExpr)
	if !ok {
		return
	}
	stName, ok := a.structTypeNameForObject(ix.Object)
	if !ok {
		return
	}
	lit, ok := ix.Index.(*parser.LiteralExpr)
	if !ok {
		return
	}
	mname, ok := lit.Value.(string)
	if !ok {
		return
	}
	methods, ok := a.structMethods[stName]
	if !ok {
		return
	}
	fd, ok := methods[mname]
	if !ok {
		return
	}
	fn := fd.Name.Lexeme
	params := methodParamsForCall(fd.Params)
	for i, arg := range call.Arguments {
		if i >= len(params) {
			break
		}
		pname := params[i].Name
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
		a.funcParamStruct[fn][pname] = st
	}
}

// When allowPlain is false (initial analysis), unknown/dynamic args are skipped so inner
// calls like dot(v, v) do not mark v as plain before refine propagates param types.
func (a *Analyzer) recordParamStructFromCall(call *parser.CallExpr, allowPlain bool) {
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
		pname := fd.Params[i].Name
		st := a.structTypeOfExpr(arg)
		if st == "" {
			if !allowPlain || !a.isPlainObjectArg(arg) {
				continue
			}
			if a.funcParamPlain == nil {
				a.funcParamPlain = make(map[string]map[string]bool)
			}
			if a.funcParamPlain[fn] == nil {
				a.funcParamPlain[fn] = make(map[string]bool)
			}
			a.funcParamPlain[fn][pname] = true
			if a.funcParamStruct != nil && a.funcParamStruct[fn] != nil {
				delete(a.funcParamStruct[fn], pname)
			}
			continue
		}
		if a.funcParamPlain != nil && a.funcParamPlain[fn] != nil && a.funcParamPlain[fn][pname] {
			continue
		}
		if a.funcParamStruct == nil {
			a.funcParamStruct = make(map[string]map[string]string)
		}
		if a.funcParamStruct[fn] == nil {
			a.funcParamStruct[fn] = make(map[string]string)
		}
		if a.funcParamPlain != nil && a.funcParamPlain[fn] != nil {
			delete(a.funcParamPlain[fn], pname)
		}
		if existing, exists := a.funcParamStruct[fn][pname]; exists {
			if !strings.EqualFold(existing, st) {
				delete(a.funcParamStruct[fn], pname)
			}
		} else {
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
		if a.activeParamStruct != nil {
			if st, ok := a.activeParamStruct[x.Name.Lexeme]; ok {
				return st
			}
		}
	case *parser.CallExpr:
		if id, ok := x.Function.(*parser.IdentifierExpr); ok {
			if st := a.structConstructorType(id.Name.Lexeme); st != "" {
				return st
			}
			if st := a.funcReturnStructType(id.Name.Lexeme); st != "" {
				return st
			}
			if a.funcReturnsPlainObject(id.Name.Lexeme) {
				return ""
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
		if id, ok := x.Object.(*parser.IdentifierExpr); ok {
			if _, ok := x.Index.(*parser.LiteralExpr); ok {
				if st, ok := a.varArrayElementStruct[id.Name.Lexeme]; ok {
					return st
				}
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
// after all call sites have been analyzed. Multiple passes propagate types through nested calls
// (e.g. length(v) calling dot(v, v)).
func (a *Analyzer) refineStructFieldAccessFromCalls(prog *parser.Program) {
	if len(a.funcParamStruct) == 0 {
		return
	}
	for pass := 0; pass < 16; pass++ {
		before := a.funcParamStructEntryCount()
		for _, decl := range prog.Declarations {
			switch d := decl.(type) {
			case *parser.FuncDecl:
				if d.Body == nil {
					continue
				}
				prevFn := a.currentFuncName
				a.currentFuncName = d.Name.Lexeme
				prev := a.activeParamStruct
				a.activeParamStruct = a.structParamsForFunc(d.Name.Lexeme)
				a.walkStmtStructFieldAccess(d.Body)
				a.activeParamStruct = prev
				a.currentFuncName = prevFn
			case *parser.StructDecl:
				for _, m := range d.Methods {
					if m.Body == nil {
						continue
					}
					prevFn := a.currentFuncName
					a.currentFuncName = m.Name.Lexeme
					prev := a.activeParamStruct
					a.activeParamStruct = a.structParamsForFunc(m.Name.Lexeme)
					a.walkStmtStructFieldAccess(m.Body)
					a.activeParamStruct = prev
					a.currentFuncName = prevFn
				}
			}
		}
		if a.funcParamStructEntryCount() == before {
			break
		}
	}
}

func (a *Analyzer) funcParamStructEntryCount() int {
	n := 0
	for _, params := range a.funcParamStruct {
		n += len(params)
	}
	return n
}

func (a *Analyzer) structParamsForFunc(fn string) map[string]string {
	structMap := a.funcParamStruct[fn]
	if structMap == nil {
		return nil
	}
	plainMap := a.funcParamPlain[fn]
	if plainMap == nil {
		return structMap
	}
	out := make(map[string]string, len(structMap))
	for name, st := range structMap {
		if plainMap[name] {
			continue
		}
		out[name] = st
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func (a *Analyzer) walkDeclStructFieldAccess(decl parser.Decl) {
	switch d := decl.(type) {
	case *parser.FuncDecl:
		prev := a.activeParamStruct
		a.activeParamStruct = a.structParamsForFunc(d.Name.Lexeme)
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
	case *parser.LoopStmt:
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
	case *parser.ForOfStmt:
		a.walkExprStructFieldAccess(s.Iterable)
		loopVar := s.VarName.Lexeme
		prev, hadPrev := a.varStructType[loopVar]
		if st := a.arrayElementStructType(s.Iterable); st != "" {
			a.varStructType[loopVar] = st
			if a.currentFuncName != "" {
				if a.forOfVarStruct[a.currentFuncName] == nil {
					a.forOfVarStruct[a.currentFuncName] = make(map[string]string)
				}
				a.forOfVarStruct[a.currentFuncName][loopVar] = st
			}
		}
		a.walkStmtStructFieldAccess(s.Body)
		if st := a.arrayElementStructType(s.Iterable); st != "" {
			if hadPrev {
				a.varStructType[loopVar] = prev
			} else {
				delete(a.varStructType, loopVar)
			}
		}
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
		a.recordParamStructFromCall(x, true)
		a.recordParamArrayElementFromCall(x)
		a.recordStructMethodParamFromCall(x)
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
