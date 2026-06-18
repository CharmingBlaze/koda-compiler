package codegen

import (
	"strings"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"koda/internal/parser"
)

func (g *Generator) structTypesForName(stName string) map[string]string {
	if g.ctx == nil || g.ctx.StructFieldTypes == nil || stName == "" {
		return nil
	}
	if fields, ok := g.ctx.StructFieldTypes[stName]; ok {
		return fields
	}
	lower := strings.ToLower(stName)
	for k, v := range g.ctx.StructFieldTypes {
		if strings.ToLower(k) == lower {
			return v
		}
	}
	return nil
}

func (g *Generator) isFloatStructField(stName, field string) bool {
	fields := g.structTypesForName(stName)
	if fields == nil {
		return false
	}
	annot, ok := fields[strings.ToLower(field)]
	return ok && isFloatTypeAnnot(annot)
}

func (g *Generator) structLayoutForName(stName string) []string {
	if g.ctx == nil || g.ctx.StructFields == nil || stName == "" {
		return nil
	}
	if layout, ok := g.ctx.StructFields[stName]; ok {
		return layout
	}
	lower := strings.ToLower(stName)
	for k, v := range g.ctx.StructFields {
		if strings.ToLower(k) == lower {
			return v
		}
	}
	return nil
}

func (g *Generator) structFieldSlotOfName(stName, field string) (int, bool) {
	layout := g.structLayoutForName(stName)
	if layout == nil {
		return 0, false
	}
	lower := strings.ToLower(field)
	for i, f := range layout {
		if strings.EqualFold(f, lower) {
			return i, true
		}
	}
	return 0, false
}

func structFieldNameFromIndex(ix *parser.IndexExpr) string {
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

func (g *Generator) structTypeForFieldIndex(ix *parser.IndexExpr) (string, bool) {
	if ix == nil {
		return "", false
	}
	switch obj := ix.Object.(type) {
	case *parser.ThisExpr:
		if g.currentStructMethodType != "" {
			return g.currentStructMethodType, true
		}
	case *parser.IdentifierExpr:
		if g.currentStructMethodType != "" && strings.EqualFold(obj.Name.Lexeme, "self") {
			return g.currentStructMethodType, true
		}
		if st := g.structTypeNameForExpr(obj); st != "" {
			return st, true
		}
	}
	return "", false
}

func (g *Generator) isStructFieldIndex(ix *parser.IndexExpr) bool {
	stName, ok := g.structTypeForFieldIndex(ix)
	if !ok {
		return false
	}
	field := structFieldNameFromIndex(ix)
	if field == "" {
		return false
	}
	_, ok = g.structFieldSlotOfName(stName, field)
	return ok
}

func (g *Generator) floatStructFieldSlot(ix *parser.IndexExpr) (int, bool) {
	stName, ok := g.structTypeForFieldIndex(ix)
	if !ok {
		return 0, false
	}
	if g.ctx != nil {
		if slot, ok := g.ctx.IndexExprStructSlot[ix]; ok {
			layout := g.structLayoutForName(stName)
			if layout != nil && slot >= 0 && slot < len(layout) {
				if g.isFloatStructField(stName, layout[slot]) {
					return slot, true
				}
			}
		}
	}
	field := structFieldNameFromIndex(ix)
	if field == "" {
		return 0, false
	}
	slot, ok := g.structFieldSlotOfName(stName, field)
	if !ok || !g.isFloatStructField(stName, field) {
		return 0, false
	}
	return slot, true
}

func (g *Generator) loadStructObjectValue(obj parser.Expr) (value.Value, bool) {
	switch x := obj.(type) {
	case *parser.ThisExpr:
		thisSlot, ok := g.locals["this"]
		if !ok {
			return nil, false
		}
		return g.block.NewLoad(types.I64, thisSlot), true
	case *parser.IdentifierExpr:
		if objI, ok := g.loadBindingI64(x.Name.Lexeme); ok {
			return objI, true
		}
	}
	return nil, false
}

func (g *Generator) loadFloatFromStructField(ix *parser.IndexExpr, slot int) (value.Value, bool) {
	objVal, ok := g.loadStructObjectValue(ix.Object)
	if !ok {
		return nil, false
	}
	idx := constant.NewInt(types.I64, int64(slot))
	boxed := g.block.NewCall(g.runtimeStructGet, objVal, idx)
	return g.block.NewCall(g.runtimeUnboxNumber, boxed), true
}

func (g *Generator) isFloatFastIdentifier(id *parser.IdentifierExpr) bool {
	if g.isFloatFastName(id.Name.Lexeme) {
		return true
	}
	if g.ctx != nil && g.currentStructMethodType != "" {
		if slot, ok := g.ctx.ImplicitStructField[id]; ok {
			layout := g.structLayoutForName(g.currentStructMethodType)
			if layout != nil && slot >= 0 && slot < len(layout) {
				return g.isFloatStructField(g.currentStructMethodType, layout[slot])
			}
		}
	}
	return false
}

func (g *Generator) loadFloatByIdentifier(id *parser.IdentifierExpr) (value.Value, bool) {
	if lf, ok := g.loadFloatByName(id.Name.Lexeme); ok {
		return lf, true
	}
	if g.ctx != nil && g.currentStructMethodType != "" {
		if slot, ok := g.ctx.ImplicitStructField[id]; ok {
			thisSlot, ok := g.locals["this"]
			if !ok {
				return nil, false
			}
			thisVal := g.block.NewLoad(types.I64, thisSlot)
			idx := constant.NewInt(types.I64, int64(slot))
			boxed := g.block.NewCall(g.runtimeStructGet, thisVal, idx)
			return g.block.NewCall(g.runtimeUnboxNumber, boxed), true
		}
	}
	return nil, false
}

func (g *Generator) isFloatFastExpr(e parser.Expr) bool {
	switch x := e.(type) {
	case *parser.IdentifierExpr:
		return g.isFloatFastIdentifier(x)
	case *parser.IndexExpr:
		if !g.isStructFieldIndex(x) {
			return false
		}
		_, ok := g.floatStructFieldSlot(x)
		return ok
	default:
		return false
	}
}

func (g *Generator) loadFloatFromExpr(e parser.Expr) (value.Value, bool) {
	switch x := e.(type) {
	case *parser.IdentifierExpr:
		return g.loadFloatByIdentifier(x)
	case *parser.IndexExpr:
		if !g.isStructFieldIndex(x) {
			return nil, false
		}
		if slot, ok := g.floatStructFieldSlot(x); ok {
			return g.loadFloatFromStructField(x, slot)
		}
	}
	return nil, false
}
