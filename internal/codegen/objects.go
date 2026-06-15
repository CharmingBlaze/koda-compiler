package codegen

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"koda/internal/lexer"
	"koda/internal/parser"
)

// storeBoxedToAssignTarget writes a NaN-boxed value to an assign target (identifier or index).
func (g *Generator) storeBoxedToAssignTarget(left parser.Expr, boxed value.Value) error {
	switch l := left.(type) {
	case *parser.IdentifierExpr:
		name := l.Name.Lexeme
		if g.ctx != nil {
			if slot, ok := g.ctx.ImplicitStructField[l]; ok {
				thisSlot, ok := g.locals["this"]
				if !ok {
					return fmt.Errorf("implicit struct field %q used outside struct method", name)
				}
				thisVal := g.block.NewLoad(types.I64, thisSlot)
				idx := constant.NewInt(types.I64, int64(slot))
				g.block.NewCall(g.runtimeStructSet, thisVal, idx, boxed)
				return nil
			}
		}
		if slot, ok := g.locals[name]; ok {
			if g.localIsCell != nil && g.localIsCell[name] {
				g.block.NewCall(g.runtimeCellWrite, slot, boxed)
			} else {
				g.block.NewStore(boxed, slot)
			}
			return nil
		}
		if slot, ok := g.moduleGlobals[name]; ok {
			if g.moduleGlobalIsCell != nil && g.moduleGlobalIsCell[name] {
				g.block.NewCall(g.runtimeCellWrite, slot, boxed)
			} else {
				g.block.NewStore(boxed, slot)
			}
			return nil
		}
		if global, ok := g.globals[name]; ok {
			g.block.NewStore(boxed, global)
			return nil
		}
		return g.undefinedVarError(name, l.Name.File, l.Name.Line, l.Name.Col)
	case *parser.IndexExpr:
		if g.ctx != nil {
			if slot, ok := g.ctx.IndexExprStructSlot[l]; ok {
				obj, err := g.emitExpr(l.Object)
				if err != nil {
					return err
				}
				g.block.NewCall(g.runtimeStructSet, g.emitAsKodaI64(obj), constant.NewInt(types.I64, int64(slot)), boxed)
				return nil
			}
		}
		obj, err := g.emitExpr(l.Object)
		if err != nil {
			return err
		}
		key, err := g.emitExpr(l.Index)
		if err != nil {
			return err
		}
		g.block.NewCall(g.runtimeSet, g.emitAsKodaI64(obj), g.emitAsKodaI64(key), boxed)
		return nil
	default:
		return fmt.Errorf("unsupported assignment target: %T", left)
	}
}

// storeNaNBoxedToAssignTarget writes boxed (i64) rhs to an identifier or index assign target.
func (g *Generator) storeNaNBoxedToAssignTarget(left parser.Expr, boxed value.Value) error {
	switch l := left.(type) {
	case *parser.IdentifierExpr:
		name := l.Name.Lexeme
		if g.ctx != nil {
			if slot, ok := g.ctx.ImplicitStructField[l]; ok {
				thisSlot, ok := g.locals["this"]
				if !ok {
					return fmt.Errorf("implicit struct field %q used outside struct method", name)
				}
				thisVal := g.block.NewLoad(types.I64, thisSlot)
				idx := constant.NewInt(types.I64, int64(slot))
				g.block.NewCall(g.runtimeStructSet, thisVal, idx, boxed)
				return nil
			}
		}
		if slot, ok := g.locals[name]; ok {
			if g.localIsCell != nil && g.localIsCell[name] {
				g.block.NewCall(g.runtimeCellWrite, slot, boxed)
			} else {
				g.block.NewStore(boxed, slot)
			}
			return nil
		}
		if slot, ok := g.moduleGlobals[name]; ok {
			if g.moduleGlobalIsCell != nil && g.moduleGlobalIsCell[name] {
				g.block.NewCall(g.runtimeCellWrite, slot, boxed)
			} else {
				g.block.NewStore(boxed, slot)
			}
			return nil
		}
		if global, ok := g.globals[name]; ok {
			g.block.NewStore(boxed, global)
			return nil
		}
		return g.undefinedVarError(name, l.Name.File, l.Name.Line, l.Name.Col)
	case *parser.IndexExpr:
		if g.ctx != nil {
			if slot, ok := g.ctx.IndexExprStructSlot[l]; ok {
				obj, err := g.emitExpr(l.Object)
				if err != nil {
					return err
				}
				g.block.NewCall(g.runtimeStructSet, g.emitAsKodaI64(obj), constant.NewInt(types.I64, int64(slot)), boxed)
				return nil
			}
		}
		obj, err := g.emitExpr(l.Object)
		if err != nil {
			return err
		}
		key, err := g.emitExpr(l.Index)
		if err != nil {
			return err
		}
		g.block.NewCall(g.runtimeSet, g.emitAsKodaI64(obj), g.emitAsKodaI64(key), boxed)
		return nil
	default:
		return fmt.Errorf("??= unsupported assignment target: %T", left)
	}
}

// emitNullishAssign lowers `lhs ??= rhs` to: if lhs is nil then lhs = rhs; result is lhs after.
func (g *Generator) emitNullishAssign(e *parser.AssignExpr) (value.Value, error) {
	rhs, err := g.emitExpr(e.Value)
	if err != nil {
		return nil, err
	}
	rhsI := g.emitAsKodaI64(rhs)

	cur, err := g.emitExpr(e.Left)
	if err != nil {
		return nil, err
	}
	curI := g.emitAsKodaI64(cur)

	nilTag := constant.NewInt(types.I64, llvmNilTagged)
	isNil := g.block.NewICmp(enum.IPredEQ, curI, nilTag)

	g.tempN++
	suf := fmt.Sprintf("qna%d", g.tempN)
	ifB := g.currentFn.NewBlock("qna.if." + suf)
	elseB := g.currentFn.NewBlock("qna.el." + suf)
	mergeB := g.currentFn.NewBlock("qna.mg." + suf)
	g.block.NewCondBr(isNil, ifB, elseB)

	g.block = ifB
	if err := g.storeNaNBoxedToAssignTarget(e.Left, rhsI); err != nil {
		return nil, err
	}
	ifB.NewBr(mergeB)

	g.block = elseB
	elseB.NewBr(mergeB)

	g.block = mergeB
	out := mergeB.NewPhi(
		ir.NewIncoming(rhsI, ifB),
		ir.NewIncoming(curI, elseB),
	)
	return out, nil
}

func (g *Generator) emitStructObject(e *parser.ObjectExpr) (value.Value, error) {
	stName := e.StructTag.Lexeme
	layout := g.ctx.StructFields[stName]
	n := len(layout)
	if n < 1 {
		n = 1
	}
	count := constant.NewInt(types.I32, int64(n))
	obj := g.block.NewCall(g.runtimeAllocStruct, count)
	objSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(objSlot)
	g.block.NewStore(obj, objSlot)
	for i, fname := range layout {
		var val parser.Expr
		for j, k := range e.Keys {
			if k.Lexeme == fname {
				val = e.Values[j]
				break
			}
		}
		if val == nil {
			if g.ctx != nil && g.ctx.StructFieldDefaults != nil {
				if defs, ok := g.ctx.StructFieldDefaults[stName]; ok {
					if def, ok := defs[fname]; ok {
						val = def
					}
				}
			}
		}
		if val == nil {
			val = &parser.LiteralExpr{Token: e.Token, Value: nil}
		}
		vv, err := g.emitExpr(val)
		if err != nil {
			return nil, err
		}
		objLive := g.block.NewLoad(types.I64, objSlot)
		idx := constant.NewInt(types.I64, int64(i))
		g.block.NewCall(g.runtimeStructSet, objLive, idx, g.emitAsKodaI64(vv))
	}
	return g.block.NewLoad(types.I64, objSlot), nil
}

func (g *Generator) emitStructWithDefaults(stName string) (value.Value, error) {
	tok := lexer.Token{Lexeme: stName}
	tag := tok
	return g.emitStructObject(&parser.ObjectExpr{Token: tok, StructTag: &tag})
}

// emitObject emits LLVM IR for object expressions.
func (g *Generator) emitObject(e *parser.ObjectExpr) (value.Value, error) {
	if e.StructTag != nil {
		return g.emitStructObject(e)
	}

	nKeys := len(e.Keys) + len(e.ComputedKeys)
	if nKeys < 1 {
		nKeys = 1 // runtime table expects non-zero capacity for key/value slots
	}
	count := constant.NewInt(types.I32, int64(nKeys))
	obj := g.block.NewCall(g.runtimeAllocObj, count)
	objSlot := g.entryAlloca(types.I64)
	g.block.NewStore(obj, objSlot)
	g.shadowStoreTemp(objSlot)

	// Set properties
	for i, key := range e.Keys {
		keyVal := g.emitStringLiteral(key.Lexeme)
		keySlot := g.entryAlloca(types.I64)
		g.shadowStoreTemp(keySlot)
		g.block.NewStore(g.emitAsKodaI64(keyVal), keySlot)
		val, err := g.emitExpr(e.Values[i])
		if err != nil {
			return nil, err
		}
		objLive := g.block.NewLoad(types.I64, objSlot)
		keyLive := g.block.NewLoad(types.I64, keySlot)
		g.block.NewCall(g.runtimeObjSet, objLive, keyLive, g.emitAsKodaI64(val))
	}

	// Set computed keys (if any)
	for i, keyExpr := range e.ComputedKeys {
		keyVal, err := g.emitExpr(keyExpr)
		if err != nil {
			return nil, err
		}
		keySlot := g.entryAlloca(types.I64)
		g.shadowStoreTemp(keySlot)
		g.block.NewStore(g.emitAsKodaI64(keyVal), keySlot)
		valIdx := len(e.Keys) + i
		val, err := g.emitExpr(e.Values[valIdx])
		if err != nil {
			return nil, err
		}
		objLive := g.block.NewLoad(types.I64, objSlot)
		keyLive := g.block.NewLoad(types.I64, keySlot)
		g.block.NewCall(g.runtimeObjSet, objLive, keyLive, g.emitAsKodaI64(val))
	}

	return g.block.NewLoad(types.I64, objSlot), nil
}

// emitIndex emits LLVM IR for index expressions (property or array access).
func (g *Generator) emitIndex(e *parser.IndexExpr) (value.Value, error) {
	if g.ctx != nil {
		if ord, ok := g.ctx.IndexExprEnumConst[e]; ok {
			return g.block.NewCall(g.runtimeBoxNumber, constant.NewFloat(types.Double, float64(ord))), nil
		}
		if slot, ok := g.ctx.IndexExprStructSlot[e]; ok {
			obj, err := g.emitExpr(e.Object)
			if err != nil {
				return nil, err
			}
			objSlot := g.entryAlloca(types.I64)
			g.shadowStoreTemp(objSlot)
			g.block.NewStore(g.emitAsKodaI64(obj), objSlot)
			idx := constant.NewInt(types.I64, int64(slot))
			if !e.Optional {
				objLive := g.block.NewLoad(types.I64, objSlot)
				return g.block.NewCall(g.runtimeStructGet, objLive, idx), nil
			}
			nilTag := constant.NewInt(types.I64, llvmNilTagged)
			objI := g.block.NewLoad(types.I64, objSlot)
			isNil := g.block.NewICmp(enum.IPredEQ, objI, nilTag)
			g.tempN++
			suf := fmt.Sprintf(".stopt%d", g.tempN)
			skip := g.currentFn.NewBlock("stnil" + suf)
			cont := g.currentFn.NewBlock("stget" + suf)
			merge := g.currentFn.NewBlock("stmg" + suf)
			g.block.NewCondBr(isNil, skip, cont)
			g.block = skip
			skip.NewBr(merge)
			g.block = cont
			got := g.block.NewCall(g.runtimeStructGet, objI, idx)
			cont.NewBr(merge)
			g.block = merge
			return merge.NewPhi(
				ir.NewIncoming(nilTag, skip),
				ir.NewIncoming(got, cont),
			), nil
		}
	}

	obj, err := g.emitExpr(e.Object)
	if err != nil {
		return nil, err
	}
	objSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(objSlot)
	g.block.NewStore(g.emitAsKodaI64(obj), objSlot)

	if !e.Optional {
		key, err := g.emitExpr(e.Index)
		if err != nil {
			return nil, err
		}
		objLive := g.block.NewLoad(types.I64, objSlot)
		return g.block.NewCall(g.runtimeArrayGet, objLive, g.emitAsKodaI64(key)), nil
	}

	nilTag := constant.NewInt(types.I64, llvmNilTagged)
	objI := g.block.NewLoad(types.I64, objSlot)
	isNil := g.block.NewICmp(enum.IPredEQ, objI, nilTag)

	g.tempN++
	suf := fmt.Sprintf(".opt%d", g.tempN)
	skip := g.currentFn.NewBlock("optnil" + suf)
	cont := g.currentFn.NewBlock("optget" + suf)
	merge := g.currentFn.NewBlock("optmg" + suf)
	g.block.NewCondBr(isNil, skip, cont)

	g.block = skip
	skip.NewBr(merge)

	g.block = cont
	key, err := g.emitExpr(e.Index)
	if err != nil {
		return nil, err
	}
	got := g.block.NewCall(g.runtimeArrayGet, objI, g.emitAsKodaI64(key))
	cont.NewBr(merge)

	g.block = merge
	return merge.NewPhi(
		ir.NewIncoming(nilTag, skip),
		ir.NewIncoming(got, cont),
	), nil
}

// emitAssign emits LLVM IR for assignment expressions.
func (g *Generator) emitAssign(e *parser.AssignExpr) (value.Value, error) {
	if e.Token.Type == lexer.TokenQuestionQuestionEqual {
		return g.emitNullishAssign(e)
	}
	if e.Token.Type != lexer.TokenEqual {
		return g.emitCompoundAssign(e, e.Token.Lexeme)
	}

	val, err := g.emitExpr(e.Value)
	if err != nil {
		return nil, err
	}

	switch left := e.Left.(type) {
	case *parser.IdentifierExpr:
		name := left.Name.Lexeme
		if slot, ok := g.locals[name]; ok {
			boxed := g.emitAsKodaI64(val)
			if g.localIsCell != nil && g.localIsCell[name] {
				g.block.NewCall(g.runtimeCellWrite, slot, boxed)
			} else {
				g.block.NewStore(boxed, slot)
			}
			return val, nil
		}
		if slot, ok := g.moduleGlobals[name]; ok {
			boxed := g.emitAsKodaI64(val)
			if g.moduleGlobalIsCell != nil && g.moduleGlobalIsCell[name] {
				g.block.NewCall(g.runtimeCellWrite, slot, boxed)
			} else {
				g.block.NewStore(boxed, slot)
			}
			return val, nil
		}
		if global, ok := g.globals[name]; ok {
			g.block.NewStore(g.emitAsKodaI64(val), global)
			return val, nil
		}
		return nil, g.undefinedVarError(name, left.Name.File, left.Name.Line, left.Name.Col)
	case *parser.IndexExpr:
		if g.ctx != nil {
			if slot, ok := g.ctx.IndexExprStructSlot[left]; ok {
				obj, err := g.emitExpr(left.Object)
				if err != nil {
					return nil, err
				}
				idx := constant.NewInt(types.I64, int64(slot))
				boxed := g.emitAsKodaI64(val)
				return g.block.NewCall(g.runtimeStructSet, g.emitAsKodaI64(obj), idx, boxed), nil
			}
		}
		// Property or array assignment
		obj, err := g.emitExpr(left.Object)
		if err != nil {
			return nil, err
		}
		key, err := g.emitExpr(left.Index)
		if err != nil {
			return nil, err
		}
		// If the value is a double, convert it to i64 (NaN-boxed)
		if val.Type().String() == "double" {
			val = g.block.NewBitCast(val, types.I64)
		}
		return g.block.NewCall(g.runtimeSet, g.emitAsKodaI64(obj), g.emitAsKodaI64(key), g.emitAsKodaI64(val)), nil
	default:
		return nil, fmt.Errorf("unsupported assignment target: %T", left)
	}
}

// emitCompoundAssign emits LLVM IR for compound assignment expressions (+=, -=, *=, /=).
func (g *Generator) emitCompoundAssign(e *parser.AssignExpr, op string) (value.Value, error) {
	// Get current value
	current, err := g.emitExpr(e.Left)
	if err != nil {
		return nil, err
	}

	// Get new value
	newVal, err := g.emitExpr(e.Value)
	if err != nil {
		return nil, err
	}

	// Perform operation on unboxed numbers, then re-box (same semantics as emitInfix for + - * /).
	leftI := g.emitAsKodaI64(current)
	rightI := g.emitAsKodaI64(newVal)
	ld := g.block.NewCall(g.runtimeUnboxNumber, leftI)
	rd := g.block.NewCall(g.runtimeUnboxNumber, rightI)
	var result value.Value
	switch op {
	case "+=":
		result = g.block.NewCall(g.runtimeValueAdd, leftI, rightI)
	case "-=":
		result = g.block.NewCall(g.runtimeValueSub, leftI, rightI)
	case "*=":
		result = g.block.NewCall(g.runtimeValueMul, leftI, rightI)
	case "/=":
		result = g.block.NewCall(g.runtimeBoxNumber, g.block.NewFDiv(ld, rd))
	case "%=":
		result = g.block.NewCall(g.runtimeBoxNumber, g.block.NewFRem(ld, rd))
	default:
		return nil, fmt.Errorf("unsupported compound operator: %s", op)
	}

	// Assign result
	boxed := g.emitAsKodaI64(result)
	if err := g.storeBoxedToAssignTarget(e.Left, boxed); err != nil {
		return nil, err
	}
	return result, nil
}
