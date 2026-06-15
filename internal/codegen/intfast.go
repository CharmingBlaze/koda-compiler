package codegen

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"koda/internal/parser"
	"koda/internal/sema"
)

func llvmIntTypeForAnnot(name string) types.Type {
	switch name {
	case "i8", "u8":
		return types.I8
	case "i16", "u16":
		return types.I16
	case "i32", "u32":
		return types.I32
	default:
		return types.I64
	}
}

func isUnsignedAnnot(name string) bool {
	switch name {
	case "u8", "u16", "u32", "u64":
		return true
	default:
		return false
	}
}

func (g *Generator) typedLocalDecl(name string) (*parser.LetDecl, string, bool) {
	if g.ctx == nil || g.ctx.TypedLocals == nil {
		return nil, "", false
	}
	for ld, annot := range g.ctx.TypedLocals {
		if ld.Name.Lexeme == name {
			return ld, annot, true
		}
	}
	return nil, "", false
}

func (g *Generator) typedIntLocalName(name string) bool {
	if g.ctx == nil || g.ctx.NumericKinds == nil || g.ctx.TypedLocals == nil {
		return false
	}
	for ld, k := range g.ctx.NumericKinds {
		if ld.Name.Lexeme == name && k == sema.KindInt {
			if _, typed := g.ctx.TypedLocals[ld]; typed {
				return true
			}
		}
	}
	return false
}

func (g *Generator) inferredIntStackName(name string) bool {
	if g.ctx == nil || g.ctx.NumericKinds == nil {
		return false
	}
	for ld, k := range g.ctx.NumericKinds {
		if ld.Name.Lexeme == name && k == sema.KindInt {
			if g.ctx.TypedLocals != nil {
				if _, typed := g.ctx.TypedLocals[ld]; typed {
					continue
				}
			}
			if g.ctx.StackDecls != nil && g.ctx.StackDecls[ld] {
				return true
			}
		}
	}
	return false
}

func (g *Generator) loadIntLocal(slot value.Value, annot string) value.Value {
	ty := llvmIntTypeForAnnot(annot)
	if ty == types.I64 {
		return g.block.NewLoad(types.I64, slot)
	}
	v := g.block.NewLoad(ty, slot)
	return g.block.NewSExt(v, types.I64)
}

func (g *Generator) storeIntLocal(slot value.Value, annot string, v value.Value) {
	ty := llvmIntTypeForAnnot(annot)
	if !v.Type().Equal(types.I64) {
		v = g.block.NewSExt(v, types.I64)
	}
	if ty == types.I64 {
		g.block.NewStore(v, slot)
		return
	}
	trunc := g.block.NewTrunc(v, ty)
	g.block.NewStore(trunc, slot)
}

func (g *Generator) tryEmitIntFastInfix(e *parser.InfixExpr) (value.Value, bool) {
	if g.ctx == nil {
		return nil, false
	}
	lid, okL := e.Left.(*parser.IdentifierExpr)
	rid, okR := e.Right.(*parser.IdentifierExpr)
	if !okL || !okR {
		return nil, false
	}
	if !g.typedIntLocalName(lid.Name.Lexeme) || !g.typedIntLocalName(rid.Name.Lexeme) {
		return nil, false
	}
	_, lAnnot, okLA := g.typedLocalDecl(lid.Name.Lexeme)
	_, rAnnot, okRA := g.typedLocalDecl(rid.Name.Lexeme)
	annot := "i64"
	if okLA {
		annot = lAnnot
	} else if okRA {
		annot = rAnnot
	}
	lslot, ok := g.locals[lid.Name.Lexeme]
	if !ok {
		return nil, false
	}
	rslot, ok := g.locals[rid.Name.Lexeme]
	if !ok {
		return nil, false
	}
	li := g.loadIntLocal(lslot, annot)
	ri := g.loadIntLocal(rslot, annot)
	ty := llvmIntTypeForAnnot(annot)
	if !li.Type().Equal(ty) {
		li = g.block.NewTrunc(li, ty)
	}
	if !ri.Type().Equal(ty) {
		ri = g.block.NewTrunc(ri, ty)
	}
	unsigned := isUnsignedAnnot(annot)
	var out value.Value
	switch e.Operator {
	case "+":
		out = g.block.NewAdd(li, ri)
	case "-":
		out = g.block.NewSub(li, ri)
	case "*":
		out = g.block.NewMul(li, ri)
	case "%":
		if unsigned {
			out = g.block.NewURem(li, ri)
		} else {
			out = g.block.NewSRem(li, ri)
		}
	case "&":
		out = g.block.NewAnd(li, ri)
	case "|":
		out = g.block.NewOr(li, ri)
	case "^":
		out = g.block.NewXor(li, ri)
	case "<<":
		out = g.block.NewShl(li, ri)
	case ">>":
		if unsigned {
			out = g.block.NewLShr(li, ri)
		} else {
			out = g.block.NewAShr(li, ri)
		}
	default:
		return nil, false
	}
	var ext value.Value
	if out.Type().Equal(types.I64) {
		ext = out
	} else {
		ext = g.block.NewSExt(out, types.I64)
	}
	return g.block.NewCall(g.runtimeBoxNumber, g.block.NewSIToFP(ext, types.Double)), true
}

func (g *Generator) tryEmitIntKindInfix(e *parser.InfixExpr) (value.Value, bool) {
	if v, ok := g.tryEmitIntFastInfix(e); ok {
		return v, true
	}
	if e.Operator != "+" && e.Operator != "-" && e.Operator != "*" && e.Operator != "%" &&
		e.Operator != "&" && e.Operator != "|" && e.Operator != "^" && e.Operator != "<<" && e.Operator != ">>" {
		return nil, false
	}
	lid, okL := e.Left.(*parser.IdentifierExpr)
	rid, okR := e.Right.(*parser.IdentifierExpr)
	if !okL || !okR {
		return nil, false
	}
	if !g.inferredIntStackName(lid.Name.Lexeme) || !g.inferredIntStackName(rid.Name.Lexeme) {
		return nil, false
	}
	lslot, ok := g.locals[lid.Name.Lexeme]
	if !ok {
		return nil, false
	}
	rslot, ok := g.locals[rid.Name.Lexeme]
	if !ok {
		return nil, false
	}
	li := g.block.NewLoad(types.I64, lslot)
	ri := g.block.NewLoad(types.I64, rslot)
	ld := g.block.NewCall(g.runtimeUnboxNumber, li)
	rd := g.block.NewCall(g.runtimeUnboxNumber, ri)
	liv := g.block.NewFPToSI(ld, types.I64)
	riv := g.block.NewFPToSI(rd, types.I64)
	var out value.Value
	switch e.Operator {
	case "+":
		out = g.block.NewAdd(liv, riv)
	case "-":
		out = g.block.NewSub(liv, riv)
	case "*":
		out = g.block.NewMul(liv, riv)
	case "%":
		out = g.block.NewSRem(liv, riv)
	case "&":
		out = g.block.NewAnd(liv, riv)
	case "|":
		out = g.block.NewOr(liv, riv)
	case "^":
		out = g.block.NewXor(liv, riv)
	case "<<":
		out = g.block.NewShl(liv, riv)
	case ">>":
		out = g.block.NewAShr(liv, riv)
	default:
		return nil, false
	}
	return g.block.NewCall(g.runtimeBoxNumber, g.block.NewSIToFP(out, types.Double)), true
}

// emitBoxedIntLiteral stores an integer literal into a typed local slot.
func (g *Generator) intLiteralForAnnot(v int64, annot string) value.Value {
	switch annot {
	case "i8", "u8":
		return constant.NewInt(types.I8, v)
	case "i16", "u16":
		return constant.NewInt(types.I16, v)
	case "i32", "u32":
		return constant.NewInt(types.I32, v)
	default:
		return constant.NewInt(types.I64, v)
	}
}

func (g *Generator) emitAsKodaI64FromTyped(slot value.Value, annot string) value.Value {
	v := g.loadIntLocal(slot, annot)
	return g.block.NewCall(g.runtimeBoxNumber, g.block.NewSIToFP(v, types.Double))
}

func (g *Generator) emitBoolFromIntCmp(pred enum.IPred, a, b value.Value) value.Value {
	cmp := g.block.NewICmp(pred, a, b)
	return g.emitBoxBoolNaN(cmp)
}
