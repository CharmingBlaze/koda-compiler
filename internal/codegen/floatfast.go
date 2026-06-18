package codegen

import (
	"math"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"koda/internal/parser"
	"koda/internal/sema"
)

func isFloatTypeAnnot(name string) bool {
	switch name {
	case "float", "float32", "float64":
		return true
	default:
		return false
	}
}

func isUnboxedNumericAnnot(name string) bool {
	if isFloatTypeAnnot(name) {
		return true
	}
	switch name {
	case "int", "i8", "i16", "i32", "i64", "u8", "u16", "u32", "u64":
		return true
	default:
		return false
	}
}

func (g *Generator) typedAnnotForName(name string) (string, bool) {
	if _, annot, ok := g.typedLocalDecl(name); ok {
		return annot, true
	}
	if g.inferredFloatStackName(name) {
		return "float", true
	}
	if annot, ok := g.paramAnnotForName(name); ok {
		return annot, true
	}
	return "", false
}

func (g *Generator) paramAnnotForName(name string) (string, bool) {
	if g.ctx == nil {
		return "", false
	}
	var owner interface{}
	var params []parser.Param
	switch {
	case g.currentEmitFuncDecl != nil:
		owner = g.currentEmitFuncDecl
		params = g.currentEmitFuncDecl.Params
	case g.currentEmitFuncExpr != nil:
		owner = g.currentEmitFuncExpr
		params = g.currentEmitFuncExpr.Params
	default:
		return "", false
	}
	for i, p := range params {
		if p.Name != name {
			continue
		}
		key := sema.NewParamCellKey(owner, i)
		if annot, ok := g.ctx.TypedParams[key]; ok {
			return annot, true
		}
		if g.ctx.ParamIsCell != nil && g.ctx.ParamIsCell[key] {
			break
		}
		if g.ctx.ParamNumericKinds != nil {
			if k, ok := g.ctx.ParamNumericKinds[key]; ok && k == sema.KindFloat {
				return "float", true
			}
		}
		break
	}
	return "", false
}

func (g *Generator) typedFloatLocalName(name string) bool {
	annot, ok := g.typedAnnotForName(name)
	return ok && isFloatTypeAnnot(annot)
}

func (g *Generator) inferredFloatStackName(name string) bool {
	if g.ctx == nil || g.ctx.NumericKinds == nil || g.ctx.StackDecls == nil {
		return false
	}
	for ld, k := range g.ctx.NumericKinds {
		if ld.Name.Lexeme != name || k != sema.KindFloat {
			continue
		}
		if g.ctx.TypedLocals != nil {
			if _, typed := g.ctx.TypedLocals[ld]; typed {
				continue
			}
		}
		if g.ctx.StackDecls[ld] {
			return true
		}
	}
	return false
}

func (g *Generator) isFloatFastName(name string) bool {
	return g.typedFloatLocalName(name) || g.inferredFloatStackName(name)
}

func (g *Generator) loadFloatLocal(slot value.Value) value.Value {
	return g.block.NewLoad(types.Double, slot)
}

func (g *Generator) storeFloatLocal(slot value.Value, v value.Value) {
	if !v.Type().Equal(types.Double) {
		v = g.block.NewCall(g.runtimeUnboxNumber, g.emitAsKodaI64(v))
	}
	g.block.NewStore(v, slot)
}

func compileTimeFloat64(e parser.Expr) (float64, bool) {
	switch x := e.(type) {
	case *parser.LiteralExpr:
		switch v := x.Value.(type) {
		case float64:
			return v, true
		case int:
			return float64(v), true
		default:
			return 0, false
		}
	case *parser.GroupingExpr:
		return compileTimeFloat64(x.Expr)
	case *parser.PrefixExpr:
		if x.Operator != "+" && x.Operator != "-" {
			return 0, false
		}
		v, ok := compileTimeFloat64(x.Right)
		if !ok {
			return 0, false
		}
		if x.Operator == "-" {
			return -v, true
		}
		return v, true
	case *parser.InfixExpr:
		if x.Operator != "+" && x.Operator != "-" && x.Operator != "*" && x.Operator != "/" {
			return 0, false
		}
		a, ok1 := compileTimeFloat64(x.Left)
		b, ok2 := compileTimeFloat64(x.Right)
		if !ok1 || !ok2 {
			return 0, false
		}
		switch x.Operator {
		case "+":
			return a + b, true
		case "-":
			return a - b, true
		case "*":
			return a * b, true
		case "/":
			if b == 0 {
				return 0, false
			}
			return a / b, true
		default:
			return 0, false
		}
	default:
		return 0, false
	}
}

func (g *Generator) floatSlotForName(name string) (value.Value, bool) {
	slot, ok := g.locals[name]
	if !ok {
		return nil, false
	}
	if g.typedFloatLocalName(name) {
		return slot, true
	}
	if g.inferredFloatStackName(name) {
		return slot, true
	}
	return nil, false
}

func (g *Generator) loadFloatByName(name string) (value.Value, bool) {
	slot, ok := g.floatSlotForName(name)
	if !ok {
		return nil, false
	}
	if g.typedFloatLocalName(name) {
		return g.loadFloatLocal(slot), true
	}
	boxed := g.block.NewLoad(types.I64, slot)
	unboxed := g.block.NewCall(g.runtimeUnboxNumber, boxed)
	return unboxed, true
}

func (g *Generator) emitFloatBinop(op string, lf, rf value.Value) (value.Value, bool) {
	var out value.Value
	switch op {
	case "+":
		out = g.block.NewFAdd(lf, rf)
	case "-":
		out = g.block.NewFSub(lf, rf)
	case "*":
		out = g.block.NewFMul(lf, rf)
	case "/":
		out = g.block.NewFDiv(lf, rf)
	default:
		return nil, false
	}
	return g.block.NewCall(g.runtimeBoxNumber, out), true
}

func (g *Generator) tryEmitFloatFastInfix(e *parser.InfixExpr) (value.Value, bool) {
	if e.Operator != "+" && e.Operator != "-" && e.Operator != "*" && e.Operator != "/" {
		return nil, false
	}
	if v, ok := g.tryEmitFloatIdIdInfix(e); ok {
		return v, true
	}
	if v, ok := g.tryEmitFloatIdLitInfix(e); ok {
		return v, true
	}
	if e.Operator == "+" || e.Operator == "*" {
		if v, ok := g.tryEmitFloatLitIdInfix(e); ok {
			return v, true
		}
	}
	return nil, false
}

func (g *Generator) tryEmitFloatIdIdInfix(e *parser.InfixExpr) (value.Value, bool) {
	if !g.isFloatFastExpr(e.Left) || !g.isFloatFastExpr(e.Right) {
		return nil, false
	}
	lf, ok := g.loadFloatFromExpr(e.Left)
	if !ok {
		return nil, false
	}
	rf, ok := g.loadFloatFromExpr(e.Right)
	if !ok {
		return nil, false
	}
	return g.emitFloatBinop(e.Operator, lf, rf)
}

func (g *Generator) tryEmitFloatIdLitInfix(e *parser.InfixExpr) (value.Value, bool) {
	if !g.isFloatFastExpr(e.Left) {
		return nil, false
	}
	lit, ok := compileTimeFloat64(e.Right)
	if !ok {
		return nil, false
	}
	lf, ok := g.loadFloatFromExpr(e.Left)
	if !ok {
		return nil, false
	}
	rf := constant.NewFloat(types.Double, lit)
	return g.emitFloatBinop(e.Operator, lf, rf)
}

func (g *Generator) tryEmitFloatLitIdInfix(e *parser.InfixExpr) (value.Value, bool) {
	if !g.isFloatFastExpr(e.Right) {
		return nil, false
	}
	lit, ok := compileTimeFloat64(e.Left)
	if !ok {
		return nil, false
	}
	rf, ok := g.loadFloatFromExpr(e.Right)
	if !ok {
		return nil, false
	}
	lf := constant.NewFloat(types.Double, lit)
	return g.emitFloatBinop(e.Operator, lf, rf)
}

func (g *Generator) emitAsKodaI64FromFloatSlot(slot value.Value) value.Value {
	return g.block.NewCall(g.runtimeBoxNumber, g.loadFloatLocal(slot))
}

func (g *Generator) floatLiteralForInit(v float64) value.Value {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return constant.NewFloat(types.Double, 0.0)
	}
	return constant.NewFloat(types.Double, v)
}

func (g *Generator) canFloatFastCompoundRHS(e parser.Expr) bool {
	if g.isFloatFastExpr(e) {
		return true
	}
	if infix, ok := e.(*parser.InfixExpr); ok {
		if infix.Operator != "+" && infix.Operator != "-" && infix.Operator != "*" && infix.Operator != "/" {
			return false
		}
		if _, ok := compileTimeFloat64(infix.Right); ok {
			return g.isFloatFastExpr(infix.Left)
		}
		if _, ok := compileTimeFloat64(infix.Left); ok && (infix.Operator == "+" || infix.Operator == "*") {
			return g.isFloatFastExpr(infix.Right)
		}
		return g.isFloatFastExpr(infix.Left) && g.isFloatFastExpr(infix.Right)
	}
	_, ok := compileTimeFloat64(e)
	return ok
}

func (g *Generator) emitFloatFastCompoundRHS(e parser.Expr) (value.Value, bool) {
	if g.isFloatFastExpr(e) {
		return g.loadFloatFromExpr(e)
	}
	if infix, ok := e.(*parser.InfixExpr); ok {
		if boxed, ok2 := g.tryEmitFloatFastInfix(infix); ok2 {
			return g.block.NewCall(g.runtimeUnboxNumber, boxed), true
		}
	}
	if lit, ok := compileTimeFloat64(e); ok {
		return constant.NewFloat(types.Double, lit), true
	}
	return nil, false
}

func (g *Generator) tryEmitFloatFastCompoundAssign(e *parser.AssignExpr, op string) (value.Value, bool) {
	if !g.isFloatFastExpr(e.Left) || !g.canFloatFastCompoundRHS(e.Value) {
		return nil, false
	}
	lf, ok := g.loadFloatFromExpr(e.Left)
	if !ok {
		return nil, false
	}
	rf, ok := g.emitFloatFastCompoundRHS(e.Value)
	if !ok {
		return nil, false
	}
	return g.emitFloatBinop(compoundToBinop(op), lf, rf)
}

func compoundToBinop(op string) string {
	switch op {
	case "+=":
		return "+"
	case "-=":
		return "-"
	case "*=":
		return "*"
	case "/=":
		return "/"
	default:
		return ""
	}
}
