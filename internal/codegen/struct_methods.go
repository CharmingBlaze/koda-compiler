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

func cloneFuncDeclName(fd *parser.FuncDecl, llvmName string) *parser.FuncDecl {
	cp := *fd
	tok := cp.Name
	tok.Lexeme = llvmName
	cp.Name = tok
	return &cp
}

func (g *Generator) emitStructMethods() error {
	if g.ctx == nil || g.ctx.StructMethods == nil {
		return nil
	}
	for stName, methods := range g.ctx.StructMethods {
		for mname, fd := range methods {
			llvmName := structMethodLLVMName(stName, mname)
			clone := cloneFuncDeclName(fd, llvmName)
			prev := g.currentStructMethodType
			g.currentStructMethodType = stName
			if err := g.emitFuncDecl(clone); err != nil {
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

func (g *Generator) structTypeForMethodReceiver(member *parser.IndexExpr) (string, bool) {
	switch obj := member.Object.(type) {
	case *parser.IdentifierExpr:
		if g.ctx == nil {
			return "", false
		}
		stName, ok := g.ctx.VarStruct[obj.Name.Lexeme]
		return stName, ok
	case *parser.ThisExpr:
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
