package codegen

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"koda/internal/parser"
	"koda/internal/sema"
)

func (g *Generator) bindFunctionParam(key sema.ParamCellKey, paramName string, paramBoxed value.Value) {
	if annot, typed := g.ctx.TypedParams[key]; typed && isFloatTypeAnnot(annot) {
		slot := g.entryAlloca(types.Double)
		g.storeFloatLocal(slot, paramBoxed)
		g.locals[paramName] = slot
		g.localIsCell[paramName] = false
		return
	}
	if g.ctx.ParamNumericKinds != nil && g.ctx.ParamNumericKinds[key] == sema.KindFloat && !g.ctx.ParamIsCell[key] {
		slot := g.entryAlloca(types.Double)
		g.storeFloatLocal(slot, paramBoxed)
		g.locals[paramName] = slot
		g.localIsCell[paramName] = false
		return
	}
	if annot, typed := g.ctx.TypedParams[key]; typed && !isFloatTypeAnnot(annot) {
		slot := g.entryAlloca(llvmIntTypeForAnnot(annot))
		unboxed := g.block.NewCall(g.runtimeUnboxNumber, paramBoxed)
		asInt := g.block.NewFPToSI(unboxed, types.I64)
		g.storeIntLocal(slot, annot, asInt)
		g.locals[paramName] = slot
		g.localIsCell[paramName] = false
		return
	}
	if g.ctx.ParamIsCell[key] {
		cell := g.block.NewCall(g.runtimeAllocCell)
		g.block.NewCall(g.runtimeCellWrite, cell, paramBoxed)
		g.locals[paramName] = cell
		g.localIsCell[paramName] = true
		return
	}
	slot := g.entryAlloca(types.I64)
	g.block.NewStore(paramBoxed, slot)
	g.locals[paramName] = slot
	g.localIsCell[paramName] = false
}

// emitFuncDecl emits LLVM IR for function declarations.
func (g *Generator) emitFuncDecl(d *parser.FuncDecl) error {
	return g.emitFuncDeclLLVM(d, "")
}

// emitFuncDeclLLVM emits a function declaration. When llvmName is non-empty it is used as the
// LLVM symbol instead of d.Name (struct methods need stable AST pointers for sema keys).
func (g *Generator) emitFuncDeclLLVM(d *parser.FuncDecl, llvmName string) error {
	name := d.Name.Lexeme

	// Handle native function declarations from // koda:extern (legacy extern directive)
	if d.Native != nil {
		// Native functions use the VM ABI: Value fn(int argCount, Value* args)
		// which maps to: i64 fn(i32, i64*)
		fn := g.mod.NewFunc(d.Native.Symbol, types.I64, ir.NewParam("argCount", types.I32), ir.NewParam("args", types.NewPointer(types.I64)))
		g.funcs[name] = fn
		return nil
	}

	if llvmName == "" {
		llvmName = name
		if g.moduleEmitPath != "" && !strings.EqualFold(name, "main") {
			llvmName = moduleFuncLLVMName(g.moduleEmitPath, name)
		} else if strings.EqualFold(name, "main") {
			llvmName = "koda_user_main"
		}
	}

	// Create function parameters, always starting with 'this'
	skipSelf := structMethodSkipsSelfParam(g, d)
	params := []*ir.Param{ir.NewParam("this", types.I64)}
	for i, param := range d.Params {
		if skipSelf && i == 0 {
			continue
		}
		params = append(params, ir.NewParam(param.Name, types.I64))
	}

	var fn *ir.Func
	if existing, ok := g.funcs[llvmName]; ok && g.funcStubs[llvmName] {
		if len(existing.Blocks) > 0 {
			return nil
		}
		fn = existing
	} else if existing, ok := g.funcs[name]; ok && g.funcStubs[name] {
		if len(existing.Blocks) > 0 {
			return nil
		}
		fn = existing
	} else {
		fn = g.mod.NewFunc(llvmName, types.I64, params...)
	}
	g.funcs[llvmName] = fn
	if strings.EqualFold(name, "main") {
		g.funcs[name] = fn
	} else if llvmName == name || llvmName == moduleFuncLLVMName(g.moduleEmitPath, name) {
		g.funcs[name] = fn
	}

	prevFn := g.currentFn
	prevBlock := g.block
	prevLocals := g.locals
	prevLocalIsCell := g.localIsCell
	prevShadowLayout := g.shadowLayout
	prevShadowFramePtr := g.shadowFramePtr
	prevShadowFrameArrTy := g.shadowFrameArrTy
	prevShadowPushed := g.shadowPushed
	prevShadowTempNext := g.shadowTempNext

	prevEmitFunc := g.currentEmitFuncName
	prevEmitFuncDecl := g.currentEmitFuncDecl
	prevEmitFuncExpr := g.currentEmitFuncExpr
	g.currentEmitFuncName = name
	g.currentEmitFuncDecl = d
	g.currentEmitFuncExpr = nil

	g.currentFn = fn
	entry := fn.NewBlock("entry")
	g.block = entry

	g.locals = make(map[string]value.Value)
	g.localIsCell = make(map[string]bool)
	if prevFn != nil && prevFn != g.funcs["user_main"] {
		for _, name := range g.ctx.FreeVarsDecl[d] {
			if slot, ok := prevLocals[name]; ok {
				g.locals[name] = slot
				if prevLocalIsCell != nil {
					g.localIsCell[name] = prevLocalIsCell[name]
				}
			}
		}
	}

	thisSlot := g.entryAlloca(types.I64)
	g.locals["this"] = thisSlot
	g.localIsCell["this"] = false
	g.block.NewStore(fn.Params[0], thisSlot)
	if skipSelf {
		g.locals["self"] = thisSlot
		g.localIsCell["self"] = false
	}

	llvmParam := 1
	for i, param := range d.Params {
		if skipSelf && i == 0 {
			continue
		}
		key := sema.NewParamCellKey(d, i)
		g.bindFunctionParam(key, param.Name, g.emitAsKodaI64(fn.Params[llvmParam]))
		llvmParam++
	}

	g.shadowLayout = g.ctx.ShadowFuncDecl[d]
	paramNames := make([]string, 0, len(d.Params))
	for i, p := range d.Params {
		if skipSelf && i == 0 {
			continue
		}
		paramNames = append(paramNames, p.Name)
	}
	g.beginShadowFrame(g.shadowLayout, thisSlot, paramNames)
	g.emitCallTracePush(llvmName, g.sourcePath, d.Name.Line)
	g.pushDeferLayer()
	defer g.popDeferLayer()

	for _, decl := range d.Body.Declarations {
		if err := g.emitDecl(decl); err != nil {
			return err
		}
	}

	if g.block.Term == nil {
		if err := g.emitDefersForCurrentLayer(); err != nil {
			return err
		}
		g.emitCallTracePop()
		g.emitShadowPop()
		g.block.NewRet(constant.NewInt(types.I64, 0))
	}
	g.shadowLayout = prevShadowLayout
	g.shadowFramePtr = prevShadowFramePtr
	g.shadowFrameArrTy = prevShadowFrameArrTy
	g.shadowPushed = prevShadowPushed
	g.shadowTempNext = prevShadowTempNext
	g.currentEmitFuncName = prevEmitFunc
	g.currentEmitFuncDecl = prevEmitFuncDecl
	g.currentEmitFuncExpr = prevEmitFuncExpr
	g.currentFn = prevFn
	g.block = prevBlock
	g.locals = prevLocals
	g.localIsCell = prevLocalIsCell
	return nil
}

// emitFuncExpr emits LLVM IR for function expressions (closures).
func (g *Generator) emitFuncExpr(e *parser.FuncExpr) (value.Value, error) {
	// Generate a unique name for the nested function
	name := fmt.Sprintf("closure_%d", g.tempN)
	g.tempN++

	freeVars := g.ctx.FreeVarsExpr[e]
	cellPtrTy := types.NewPointer(types.I64)

	params := []*ir.Param{ir.NewParam("this", types.I64)}
	for _, fv := range freeVars {
		params = append(params, ir.NewParam("__cap_"+fv, cellPtrTy))
	}
	for _, param := range e.Params {
		params = append(params, ir.NewParam(param.Name, types.I64))
	}

	fn := g.mod.NewFunc(name, types.I64, params...)
	g.funcs[name] = fn

	prevFn := g.currentFn
	prevBlock := g.block
	prevLocals := g.locals
	prevLocalIsCell := g.localIsCell
	prevShadowLayout := g.shadowLayout
	prevShadowFramePtr := g.shadowFramePtr
	prevShadowFrameArrTy := g.shadowFrameArrTy
	prevShadowPushed := g.shadowPushed
	prevShadowTempNext := g.shadowTempNext
	prevEmitFunc := g.currentEmitFuncName
	prevEmitFuncDecl := g.currentEmitFuncDecl
	prevEmitFuncExpr := g.currentEmitFuncExpr

	g.currentFn = fn
	g.currentEmitFuncName = name
	g.currentEmitFuncDecl = nil
	g.currentEmitFuncExpr = e
	entry := fn.NewBlock("entry")
	g.block = entry
	g.locals = make(map[string]value.Value)
	g.localIsCell = make(map[string]bool)

	for i, fv := range freeVars {
		cellParam := fn.Params[1+i]
		g.locals[fv] = cellParam
		g.localIsCell[fv] = true
	}

	paramOffset := 1 + len(freeVars)

	thisSlot := g.entryAlloca(types.I64)
	g.locals["this"] = thisSlot
	g.localIsCell["this"] = false
	g.block.NewStore(fn.Params[0], thisSlot)

	for i, param := range e.Params {
		key := sema.NewParamCellKey(e, i)
		g.bindFunctionParam(key, param.Name, g.emitAsKodaI64(fn.Params[paramOffset+i]))
	}

	layout := g.ctx.ShadowFuncExpr[e]
	g.shadowLayout = layout
	paramNames := make([]string, len(e.Params))
	for i := range e.Params {
		paramNames[i] = e.Params[i].Name
	}
	g.beginShadowFrame(layout, thisSlot, paramNames)
	g.emitCallTracePush(name, g.sourcePath, e.Token.Line)
	g.pushDeferLayer()
	defer g.popDeferLayer()
	if layout != nil {
		for _, name := range g.ctx.FreeVarsExpr[e] {
			if idx, ok := layout.FreeVarIndex[name]; ok {
				if slot, ok2 := g.locals[name]; ok2 {
					g.shadowStoreIndex(idx, slot)
				}
			}
		}
	}

	// Emit function body
	for _, decl := range e.Body.Declarations {
		if err := g.emitDecl(decl); err != nil {
			return nil, err
		}
	}

	if g.block.Term == nil {
		if err := g.emitDefersForCurrentLayer(); err != nil {
			return nil, err
		}
		g.emitCallTracePop()
		g.emitShadowPop()
		g.block.NewRet(constant.NewInt(types.I64, 0))
	}

	// Restore previous function state
	g.shadowLayout = prevShadowLayout
	g.shadowFramePtr = prevShadowFramePtr
	g.shadowFrameArrTy = prevShadowFrameArrTy
	g.shadowPushed = prevShadowPushed
	g.shadowTempNext = prevShadowTempNext
	g.currentFn = prevFn
	g.block = prevBlock
	g.locals = prevLocals
	g.localIsCell = prevLocalIsCell
	g.currentEmitFuncName = prevEmitFunc
	g.currentEmitFuncDecl = prevEmitFuncDecl
	g.currentEmitFuncExpr = prevEmitFuncExpr

	nFV := len(freeVars)
	if nFV == 0 {
		// Tag raw function pointers so indirect calls use (this, args…) without cell slots.
		fnI64 := g.block.NewPtrToInt(fn, types.I64)
		return g.block.NewOr(fnI64, constant.NewInt(types.I64, 1)), nil
	}
	if nFV > maxClosureFreeVars {
		return nil, fmt.Errorf("closure %s: too many captured variables (%d > %d)", name, nFV, maxClosureFreeVars)
	}

	sz := int64(8 * (2 + nFV))
	raw := g.block.NewCall(g.runtimeMalloc, constant.NewInt(types.I64, sz))
	pI64 := g.block.NewBitCast(raw, types.NewPointer(types.I64))
	g.block.NewStore(g.block.NewPtrToInt(fn, types.I64), pI64)
	pN := g.block.NewGetElementPtr(types.I64, pI64, constant.NewInt(types.I32, 1))
	g.block.NewStore(constant.NewInt(types.I64, int64(nFV)), pN)
	for i, fv := range freeVars {
		slot, ok := prevLocals[fv]
		if !ok {
			return nil, fmt.Errorf("closure %s: missing parent slot for captured %q", name, fv)
		}
		pSlot := g.block.NewGetElementPtr(types.I64, pI64, constant.NewInt(types.I32, int64(2+i)))
		g.block.NewStore(g.block.NewPtrToInt(slot, types.I64), pSlot)
	}
	return g.block.NewPtrToInt(raw, types.I64), nil
}
