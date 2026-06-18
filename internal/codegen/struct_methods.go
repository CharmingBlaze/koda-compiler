package codegen

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir/value"

	"koda/internal/parser"
)

func structMethodLLVMName(stName, methodName string) string {
	return fmt.Sprintf("koda_method_%s_%s", stName, methodName)
}

func structMethodSkipsSelfParam(g *Generator, d *parser.FuncDecl) bool {
	return g.currentStructMethodType != "" &&
		len(d.Params) > 0 &&
		strings.EqualFold(d.Params[0].Name, "self")
}

func (g *Generator) emitStructMethods() error {
	if g.ctx == nil || g.ctx.StructMethods == nil {
		return nil
	}
	for stName, methods := range g.ctx.StructMethods {
		for mname, fd := range methods {
			llvmName := structMethodLLVMName(stName, mname)
			prev := g.currentStructMethodType
			g.currentStructMethodType = stName
			if err := g.emitFuncDeclLLVM(fd, llvmName); err != nil {
				g.currentStructMethodType = prev
				return err
			}
			g.currentStructMethodType = prev
			key := stName + "." + mname
			if fn, ok := g.funcs[llvmName]; ok {
				g.funcs[key] = fn
				g.funcs[strings.ToLower(key)] = fn
			}
		}
	}
	return nil
}

func (g *Generator) structTypeNameForExpr(expr parser.Expr) string {
	if g.ctx == nil {
		return ""
	}
	id, ok := expr.(*parser.IdentifierExpr)
	if !ok {
		return ""
	}
	name := id.Name.Lexeme
	if st, ok := g.ctx.VarStruct[name]; ok {
		return st
	}
	if g.currentEmitFuncName != "" && g.ctx.FuncForOfVarStruct != nil {
		if vars, ok := g.ctx.FuncForOfVarStruct[g.currentEmitFuncName]; ok {
			if st, ok := vars[name]; ok {
				return st
			}
		}
	}
	return ""
}

func (g *Generator) structFieldSlotOf(stName, field string) (int, bool) {
	if g.ctx == nil {
		return 0, false
	}
	layout, ok := g.ctx.StructFields[stName]
	if !ok {
		return 0, false
	}
	for i, f := range layout {
		if f == field {
			return i, true
		}
	}
	return 0, false
}

// structFieldSlotForIndex resolves a string field index on a struct-typed receiver,
// including for-of loop variables (e.g. coin.on, p.x).
func (g *Generator) structFieldSlotForIndex(ix *parser.IndexExpr) (int, bool) {
	lit, ok := ix.Index.(*parser.LiteralExpr)
	if !ok {
		return 0, false
	}
	field, ok := lit.Value.(string)
	if !ok {
		return 0, false
	}
	stName := g.structTypeNameForExpr(ix.Object)
	if stName == "" {
		return 0, false
	}
	return g.structFieldSlotOf(stName, field)
}

func (g *Generator) structTypeForMethodReceiver(member *parser.IndexExpr) (string, bool) {
	stName := g.structTypeNameForExpr(member.Object)
	if stName != "" {
		return stName, true
	}
	if _, ok := member.Object.(*parser.ThisExpr); ok {
		if g.currentStructMethodType != "" {
			return g.currentStructMethodType, true
		}
	}
	return "", false
}

// tryEmitStructConstructor lowers Coin(7, -5) to alloc-with-defaults + new(...).
func (g *Generator) tryEmitStructConstructor(call *parser.CallExpr) (value.Value, bool, error) {
	id, ok := call.Function.(*parser.IdentifierExpr)
	if !ok || g.ctx == nil || g.ctx.StructMethods == nil {
		return nil, false, nil
	}
	stName := id.Name.Lexeme
	methods, ok := g.ctx.StructMethods[stName]
	if !ok {
		return nil, false, nil
	}
	if _, ok := methods["new"]; !ok {
		return nil, false, nil
	}
	obj, err := g.emitStructWithDefaults(stName)
	if err != nil {
		return nil, true, err
	}
	key := stName + ".new"
	fn, ok := g.funcs[key]
	if !ok {
		fn, ok = g.funcs[strings.ToLower(key)]
	}
	if !ok {
		return nil, true, fmt.Errorf("constructor %s.new not emitted", stName)
	}
	thisVal := g.emitAsKodaI64(obj)
	var args []value.Value
	for _, arg := range call.Arguments {
		v, err := g.emitExpr(arg)
		if err != nil {
			return nil, true, err
		}
		args = append(args, v)
	}
	finalArgs := append([]value.Value{thisVal}, args...)
	g.block.NewCall(fn, finalArgs...)
	return obj, true, nil
}

// tryEmitStructMethodCall dispatches rect.area() on struct-typed receivers.
func (g *Generator) tryEmitStructMethodCall(member *parser.IndexExpr, recvVal value.Value, call *parser.CallExpr) (value.Value, bool, error) {
	if g.ctx == nil || g.ctx.StructMethods == nil {
		return nil, false, nil
	}
	stName, ok := g.structTypeForMethodReceiver(member)
	if !ok {
		return nil, false, nil
	}
	lit, ok := member.Index.(*parser.LiteralExpr)
	if !ok {
		return nil, false, nil
	}
	mname, ok := lit.Value.(string)
	if !ok {
		return nil, false, nil
	}
	methods, ok := g.ctx.StructMethods[stName]
	if !ok {
		return nil, false, nil
	}
	if _, ok := methods[mname]; !ok {
		return nil, false, nil
	}
	key := stName + "." + mname
	fn, ok := g.funcs[key]
	if !ok {
		fn, ok = g.funcs[strings.ToLower(key)]
	}
	if !ok {
		return nil, false, nil
	}
	thisVal := g.emitAsKodaI64(recvVal)
	var args []value.Value
	for _, arg := range call.Arguments {
		v, err := g.emitExpr(arg)
		if err != nil {
			return nil, true, err
		}
		args = append(args, v)
	}
	finalArgs := append([]value.Value{thisVal}, args...)
	return g.block.NewCall(fn, finalArgs...), true, nil
}
