package codegen

import (
	"fmt"
	"strings"

	"koda/internal/parser"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// tryEmitMethodCall lowers argv runtime calls and LLVM-backed `.map()` / `.filter()` / `.find()` / `.reduce()`.
func (g *Generator) tryEmitMethodCall(member *parser.IndexExpr, recvVal value.Value, call *parser.CallExpr) (value.Value, bool, error) {
	lit, ok := member.Index.(*parser.LiteralExpr)
	if !ok {
		return nil, false, nil
	}
	name, ok := lit.Value.(string)
	if !ok {
		return nil, false, nil
	}
	name = strings.ToLower(name)

	if v, handled, err := g.tryEmitNativeNamespaceCall(member, name, call.Arguments); handled {
		return v, true, err
	}

	switch name {
	case "concat":
		args := []value.Value{recvVal}
		for _, arg := range call.Arguments {
			v, err := g.emitExpr(arg)
			if err != nil {
				return nil, true, err
			}
			args = append(args, v)
		}
		return g.emitArgvRuntime(g.runtimeArrayConcat, args), true, nil

	case "split", "trim", "toupper", "tolower", "replace", "replaceall", "startswith", "endswith":
		v, err := g.emitStringMethod(name, recvVal, call.Arguments)
		return v, true, err

	case "join", "sort", "reverse", "push", "pop":
		v, err := g.emitArrayOnlyMethod(name, recvVal, call.Arguments)
		return v, true, err

	case "slice":
		v, err := g.emitSliceAmbiguous(recvVal, call.Arguments)
		return v, true, err
	case "indexof", "includes":
		v, err := g.emitIndexOrIncludesAmbiguous(name, recvVal, call.Arguments)
		return v, true, err
	case "length":
		if len(call.Arguments) != 0 {
			return nil, true, fmt.Errorf("length expects 0 arguments")
		}
		return g.block.NewCall(g.runtimeLen, g.emitAsKodaI64(recvVal)), true, nil

	case "map":
		return g.emitArrayMethodMap(recvVal, call.Arguments)
	case "filter":
		return g.emitArrayMethodFilter(recvVal, call.Arguments)
	case "find":
		return g.emitArrayMethodFind(recvVal, call.Arguments)
	case "reduce":
		return g.emitArrayMethodReduce(recvVal, call.Arguments)
	}

	return nil, false, nil
}

// tryEmitNativeNamespaceCall handles calls like math.lerp(...), where method dispatch
// should route to a known argv-native symbol instead of generic indirect-call ABI.
func (g *Generator) tryEmitNativeNamespaceCall(member *parser.IndexExpr, name string, args []parser.Expr) (value.Value, bool, error) {
	id, ok := member.Object.(*parser.IdentifierExpr)
	if !ok {
		return nil, false, nil
	}
	idName := strings.ToLower(id.Name.Lexeme)
	var fn *ir.Func
	switch idName {
	case "math":
		fn = g.mathNamespaceNative(name)
	case "json":
		fn = g.jsonNamespaceNative(name)
	case "io":
		fn = g.ioNamespaceNative(name)
	default:
		return nil, false, nil
	}
	if fn == nil || !isNativeArgvCallee(fn) {
		if fn == nil {
			return nil, false, nil
		}
		if len(fn.Params) == 1 && fn.Params[0].Typ.Equal(types.I64) && len(args) == 1 {
			a0, err := g.emitExpr(args[0])
			if err != nil {
				return nil, true, err
			}
			return g.block.NewCall(fn, g.emitAsKodaI64(a0)), true, nil
		}
		return nil, false, nil
	}
	argv := make([]value.Value, 0, len(args))
	for _, a := range args {
		v, err := g.emitExpr(a)
		if err != nil {
			return nil, true, err
		}
		argv = append(argv, v)
	}
	return g.emitArgvRuntime(fn, argv), true, nil
}

func (g *Generator) mathNamespaceNative(name string) *ir.Func {
	switch name {
	case "lerp":
		return g.runtimeLerp
	case "clamp":
		return g.runtimeClamp
	case "distance":
		return g.runtimeDistance
	case "anglebetween":
		return g.runtimeAngleBetween
	case "map":
		return g.runtimeMap
	case "sin":
		return g.runtimeSin
	case "cos":
		return g.runtimeCos
	case "tan":
		return g.runtimeTan
	case "asin":
		return g.runtimeAsin
	case "acos":
		return g.runtimeAcos
	case "atan":
		return g.runtimeAtan
	case "atan2":
		return g.runtimeAtan2
	case "pow":
		return g.runtimePow
	case "exp":
		return g.runtimeExp
	case "log":
		return g.runtimeLog
	case "log10":
		return g.runtimeLog10
	case "log2":
		return g.runtimeLog2
	case "floor":
		return g.runtimeFloor
	case "ceil":
		return g.runtimeCeil
	case "round":
		return g.runtimeRound
	case "trunc":
		return g.runtimeTrunc
	case "sign":
		return g.runtimeSign
	case "min":
		return g.runtimeMin
	case "max":
		return g.runtimeMax
	case "smoothstep":
		return g.runtimeSmoothstep
	case "distancesq":
		return g.runtimeDistanceSq
	case "normalize":
		return g.runtimeNormalize
	case "hypot":
		return g.runtimeHypot
	case "sqrt":
		return g.runtimeSqrt
	case "cbrt":
		return g.runtimeCbrt
	case "abs":
		return g.runtimeAbs
	case "fmod":
		return g.runtimeFmod
	case "degrees":
		return g.runtimeDegrees
	case "radians":
		return g.runtimeRadians
	case "wrap":
		return g.runtimeWrap
	case "approach":
		return g.runtimeApproach
	case "smoothdamp":
		return g.runtimeSmoothdamp
	case "random":
		return g.runtimeRandom
	case "randomint":
		return g.runtimeRandomInt
	case "randomchoice":
		return g.runtimeRandomChoice
	case "randomseed":
		return g.runtimeRandomSeed
	case "randomrange":
		return g.runtimeRandom
	default:
		return nil
	}
}

func (g *Generator) jsonNamespaceNative(name string) *ir.Func {
	switch name {
	case "parse":
		return g.runtimeJsonParse
	case "stringify":
		return g.runtimeToJSON
	case "try_parse", "tryparse":
		return g.runtimeJsonTryParse
	default:
		return nil
	}
}

func (g *Generator) ioNamespaceNative(name string) *ir.Func {
	switch name {
	case "read":
		return g.runtimeReadFile
	case "write":
		return g.runtimeWriteFile
	case "append":
		return g.runtimeAppendFile
	case "exists":
		return g.runtimeFileExists
	case "remove":
		return g.runtimeDeleteFile
	case "isfile":
		return g.runtimeIsFile
	case "isdir":
		return g.runtimeIsDir
	case "size":
		return g.runtimeFileSize
	case "list":
		return g.runtimeListDir
	default:
		return nil
	}
}

func (g *Generator) emitStringMethod(name string, recv value.Value, args []parser.Expr) (value.Value, error) {
	recvSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(recvSlot)
	g.block.NewStore(g.emitAsKodaI64(recv), recvSlot)
	loadRecv := func() value.Value { return g.block.NewLoad(types.I64, recvSlot) }

	switch name {
	case "split":
		if len(args) != 1 {
			return nil, fmt.Errorf("split expects 1 argument (delimiter)")
		}
		a0, err := g.emitExpr(args[0])
		if err != nil {
			return nil, err
		}
		return g.emitArgvRuntime(g.runtimeStringSplit, []value.Value{loadRecv(), a0}), nil
	case "trim":
		return g.emitArgvRuntime(g.runtimeStringTrim, []value.Value{loadRecv()}), nil
	case "toupper":
		return g.emitArgvRuntime(g.runtimeStringUpper, []value.Value{loadRecv()}), nil
	case "tolower":
		return g.emitArgvRuntime(g.runtimeStringLower, []value.Value{loadRecv()}), nil
	case "replace":
		if len(args) != 2 {
			return nil, fmt.Errorf("replace expects 2 arguments")
		}
		a0, err := g.emitExpr(args[0])
		if err != nil {
			return nil, err
		}
		a1, err := g.emitExpr(args[1])
		if err != nil {
			return nil, err
		}
		return g.emitArgvRuntime(g.runtimeStringReplace, []value.Value{loadRecv(), a0, a1}), nil
	case "replaceall":
		if len(args) != 2 {
			return nil, fmt.Errorf("replaceAll expects 2 arguments")
		}
		a0, err := g.emitExpr(args[0])
		if err != nil {
			return nil, err
		}
		a1, err := g.emitExpr(args[1])
		if err != nil {
			return nil, err
		}
		return g.emitArgvRuntime(g.runtimeStringReplaceAll, []value.Value{loadRecv(), a0, a1}), nil
	case "startswith":
		if len(args) != 1 {
			return nil, fmt.Errorf("startsWith expects 1 argument")
		}
		a0, err := g.emitExpr(args[0])
		if err != nil {
			return nil, err
		}
		return g.emitArgvRuntime(g.runtimeStringStartsWith, []value.Value{loadRecv(), a0}), nil
	case "endswith":
		if len(args) != 1 {
			return nil, fmt.Errorf("endsWith expects 1 argument")
		}
		a0, err := g.emitExpr(args[0])
		if err != nil {
			return nil, err
		}
		return g.emitArgvRuntime(g.runtimeStringEndsWith, []value.Value{loadRecv(), a0}), nil
	default:
		return nil, fmt.Errorf("unknown string method %q", name)
	}
}

func (g *Generator) emitArrayOnlyMethod(name string, recv value.Value, args []parser.Expr) (value.Value, error) {
	recvSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(recvSlot)
	g.block.NewStore(g.emitAsKodaI64(recv), recvSlot)
	loadRecv := func() value.Value { return g.block.NewLoad(types.I64, recvSlot) }

	switch name {
	case "join":
		if len(args) != 1 {
			return nil, fmt.Errorf("join expects 1 argument (separator)")
		}
		a0, err := g.emitExpr(args[0])
		if err != nil {
			return nil, err
		}
		return g.emitArgvRuntime(g.runtimeArrayJoin, []value.Value{loadRecv(), a0}), nil
	case "sort":
		return g.emitArgvRuntime(g.runtimeArraySort, []value.Value{loadRecv()}), nil
	case "reverse":
		return g.emitArgvRuntime(g.runtimeArrayReverse, []value.Value{loadRecv()}), nil
	case "push":
		if len(args) != 1 {
			return nil, fmt.Errorf("push expects 1 argument")
		}
		a0, err := g.emitExpr(args[0])
		if err != nil {
			return nil, err
		}
		g.block.NewCall(g.runtimeArrayPush, loadRecv(), g.emitAsKodaI64(a0))
		return loadRecv(), nil
	case "pop":
		if len(args) != 0 {
			return nil, fmt.Errorf("pop expects 0 arguments")
		}
		return g.block.NewCall(g.runtimeArrayPop, loadRecv()), nil
	default:
		return nil, fmt.Errorf("unknown array method %q", name)
	}
}

func (g *Generator) emitSliceAmbiguous(recv value.Value, args []parser.Expr) (value.Value, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("slice expects 2 arguments (start, end)")
	}
	a0, err := g.emitExpr(args[0])
	if err != nil {
		return nil, err
	}
	a1, err := g.emitExpr(args[1])
	if err != nil {
		return nil, err
	}
	recvSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(recvSlot)
	g.block.NewStore(g.emitAsKodaI64(recv), recvSlot)
	loadRecv := func() value.Value { return g.block.NewLoad(types.I64, recvSlot) }
	return g.emitArgvNilTaggedFallback(g.emitArgvRuntime(g.runtimeArraySlice, []value.Value{loadRecv(), a0, a1}), func() value.Value {
		return g.emitArgvRuntime(g.runtimeStringSlice, []value.Value{loadRecv(), a0, a1})
	}), nil
}

func (g *Generator) emitIndexOrIncludesAmbiguous(name string, recv value.Value, args []parser.Expr) (value.Value, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("%s expects 1 argument", name)
	}
	a0, err := g.emitExpr(args[0])
	if err != nil {
		return nil, err
	}
	recvSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(recvSlot)
	g.block.NewStore(g.emitAsKodaI64(recv), recvSlot)
	loadRecv := func() value.Value { return g.block.NewLoad(types.I64, recvSlot) }
	argv := []value.Value{loadRecv(), a0}
	switch name {
	case "indexof":
		va := g.emitArgvRuntime(g.runtimeArrayIndexOf, argv)
		return g.emitArgvNilTaggedFallback(va, func() value.Value {
			return g.emitArgvRuntime(g.runtimeStringIndexOf, argv)
		}), nil
	case "includes":
		va := g.emitArgvRuntime(g.runtimeArrayIncludes, argv)
		return g.emitArgvNilTaggedFallback(va, func() value.Value {
			// Strings use contains-like helpers via argv wrappers only for booleans — approximate via indexOf >= 0.
			io := g.emitArgvRuntime(g.runtimeStringIndexOf, argv)
			negOne := constant.NewFloat(types.Double, -1)
			boxNeg := g.block.NewCall(g.runtimeBoxNumber, negOne)
			cmp := g.block.NewFCmp(enum.FPredOGT, g.block.NewCall(g.runtimeUnboxNumber, g.emitAsKodaI64(io)), g.block.NewCall(g.runtimeUnboxNumber, boxNeg))
			return g.emitBoxBoolNaN(cmp)
		}), nil
	default:
		return nil, fmt.Errorf("unknown ambiguous method %q", name)
	}
}

// emitArgvNilTaggedFallback returns primary if it is non-nil-tagged; otherwise evaluates alt() on an alternate block.
func (g *Generator) emitArgvNilTaggedFallback(primary value.Value, alt func() value.Value) value.Value {
	entry := g.block
	nilTag := constant.NewInt(types.I64, llvmNilTagged)
	isNil := entry.NewICmp(enum.IPredEQ, g.emitAsKodaI64(primary), nilTag)

	g.tempN++
	suf := fmt.Sprintf(".nf%d", g.tempN)
	okB := g.currentFn.NewBlock("nilfb.ok" + suf)
	altB := g.currentFn.NewBlock("nilfb.alt" + suf)
	mergeB := g.currentFn.NewBlock("nilfb.mg" + suf)
	entry.NewCondBr(isNil, altB, okB)

	g.block = okB
	okB.NewBr(mergeB)

	g.block = altB
	vAlt := alt()
	altTerm := g.block
	altTerm.NewBr(mergeB)

	g.block = mergeB
	return mergeB.NewPhi(
		ir.NewIncoming(primary, okB),
		ir.NewIncoming(vAlt, altTerm),
	)
}

func (g *Generator) emitArrayLenAsDouble(recv value.Value) value.Value {
	lv := g.block.NewCall(g.runtimeLen, g.emitAsKodaI64(recv))
	return g.block.NewCall(g.runtimeUnboxNumber, g.emitAsKodaI64(lv))
}

func (g *Generator) emitArrayMethodMap(recv value.Value, args []parser.Expr) (value.Value, bool, error) {
	if len(args) != 1 {
		return nil, true, fmt.Errorf("map expects 1 callback")
	}
	recvSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(recvSlot)
	g.block.NewStore(g.emitAsKodaI64(recv), recvSlot)

	fnVal, err := g.emitExpr(args[0])
	if err != nil {
		return nil, true, err
	}
	nThis := constant.NewInt(types.I64, 0)

	recvLive := g.block.NewLoad(types.I64, recvSlot)
	lenF := g.emitArrayLenAsDouble(recvLive)
	capI := g.block.NewFPToSI(lenF, types.I32)
	out := g.block.NewCall(g.runtimeAllocArray, capI)
	outSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(outSlot)
	g.block.NewStore(g.emitAsKodaI64(out), outSlot)

	g.tempN++
	suf := fmt.Sprintf(".map%d", g.tempN)
	hdr := g.currentFn.NewBlock("map.hdr" + suf)
	body := g.currentFn.NewBlock("map.body" + suf)
	step := g.currentFn.NewBlock("map.step" + suf)
	done := g.currentFn.NewBlock("map.done" + suf)

	idx := g.entryAlloca(types.Double)
	g.block.NewStore(constant.NewFloat(types.Double, 0), idx)
	g.block.NewBr(hdr)

	g.block = hdr
	idxNow := g.block.NewLoad(types.Double, idx)
	ok := g.block.NewFCmp(enum.FPredOLT, idxNow, lenF)
	g.block.NewCondBr(ok, body, done)

	g.block = body
	idxBox := g.block.NewCall(g.runtimeBoxNumber, idxNow)
	el := g.block.NewCall(g.runtimeArrayGet, recvLive, g.emitAsKodaI64(idxBox))
	elSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(elSlot)
	g.block.NewStore(g.emitAsKodaI64(el), elSlot)
	elLive := g.block.NewLoad(types.I64, elSlot)
	mv, err := g.emitIndirectI64Callee(fnVal, nThis, []value.Value{elLive})
	if err != nil {
		return nil, true, err
	}
	idxInt := g.block.NewFPToSI(idxNow, types.I64)
	outLive := g.block.NewLoad(types.I64, outSlot)
	g.block.NewCall(g.runtimeArraySet, outLive, idxInt, g.emitAsKodaI64(mv))
	g.block.NewBr(step)

	g.block = step
	one := constant.NewFloat(types.Double, 1)
	next := g.block.NewFAdd(g.block.NewLoad(types.Double, idx), one)
	g.block.NewStore(next, idx)
	g.block.NewBr(hdr)

	g.block = done
	return g.block.NewLoad(types.I64, outSlot), true, nil
}

func (g *Generator) emitArrayMethodFilter(recv value.Value, args []parser.Expr) (value.Value, bool, error) {
	if len(args) != 1 {
		return nil, true, fmt.Errorf("filter expects 1 callback")
	}
	recvSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(recvSlot)
	g.block.NewStore(g.emitAsKodaI64(recv), recvSlot)

	fnVal, err := g.emitExpr(args[0])
	if err != nil {
		return nil, true, err
	}
	out := g.block.NewCall(g.runtimeAllocArray, constant.NewInt(types.I32, 1))
	outSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(outSlot)
	g.block.NewStore(g.emitAsKodaI64(out), outSlot)
	nThis := constant.NewInt(types.I64, 0)

	recvLive := g.block.NewLoad(types.I64, recvSlot)
	lenF := g.emitArrayLenAsDouble(recvLive)

	g.tempN++
	suf := fmt.Sprintf(".filt%d", g.tempN)
	hdr := g.currentFn.NewBlock("filt.hdr" + suf)
	body := g.currentFn.NewBlock("filt.body" + suf)
	pushB := g.currentFn.NewBlock("filt.push" + suf)
	skipB := g.currentFn.NewBlock("filt.skip" + suf)
	step := g.currentFn.NewBlock("filt.step" + suf)
	done := g.currentFn.NewBlock("filt.done" + suf)

	idx := g.entryAlloca(types.Double)
	g.block.NewStore(constant.NewFloat(types.Double, 0), idx)
	g.block.NewBr(hdr)

	g.block = hdr
	idxNow := g.block.NewLoad(types.Double, idx)
	ok := g.block.NewFCmp(enum.FPredOLT, idxNow, lenF)
	g.block.NewCondBr(ok, body, done)

	g.block = body
	idxBox := g.block.NewCall(g.runtimeBoxNumber, idxNow)
	el := g.block.NewCall(g.runtimeArrayGet, recvLive, g.emitAsKodaI64(idxBox))
	elSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(elSlot)
	g.block.NewStore(g.emitAsKodaI64(el), elSlot)
	elLive := g.block.NewLoad(types.I64, elSlot)
	predV, err := g.emitIndirectI64Callee(fnVal, nThis, []value.Value{elLive})
	if err != nil {
		return nil, true, err
	}
	brOk := g.emitTruthy(predV)
	g.block.NewCondBr(brOk, pushB, skipB)

	g.block = pushB
	outLive := g.block.NewLoad(types.I64, outSlot)
	pushVal := g.block.NewLoad(types.I64, elSlot)
	g.block.NewCall(g.runtimeArrayPush, outLive, pushVal)
	g.block.NewBr(step)

	g.block = skipB
	skipB.NewBr(step)

	g.block = step
	one := constant.NewFloat(types.Double, 1)
	next := g.block.NewFAdd(g.block.NewLoad(types.Double, idx), one)
	g.block.NewStore(next, idx)
	g.block.NewBr(hdr)

	g.block = done
	return g.block.NewLoad(types.I64, outSlot), true, nil
}

func (g *Generator) emitArrayMethodFind(recv value.Value, args []parser.Expr) (value.Value, bool, error) {
	if len(args) != 1 {
		return nil, true, fmt.Errorf("find expects 1 callback")
	}
	recvSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(recvSlot)
	g.block.NewStore(g.emitAsKodaI64(recv), recvSlot)

	fnVal, err := g.emitExpr(args[0])
	if err != nil {
		return nil, true, err
	}
	nThis := constant.NewInt(types.I64, 0)
	recvLive := g.block.NewLoad(types.I64, recvSlot)
	lenF := g.emitArrayLenAsDouble(recvLive)

	g.tempN++
	suf := fmt.Sprintf(".find%d", g.tempN)
	merge := g.currentFn.NewBlock("find.merge" + suf)
	hdr := g.currentFn.NewBlock("find.hdr" + suf)
	body := g.currentFn.NewBlock("find.body" + suf)
	foundHit := g.currentFn.NewBlock("find.hit" + suf)
	step := g.currentFn.NewBlock("find.step" + suf)
	miss := g.currentFn.NewBlock("find.miss" + suf)

	idx := g.entryAlloca(types.Double)
	g.block.NewStore(constant.NewFloat(types.Double, 0), idx)
	g.block.NewBr(hdr)

	g.block = hdr
	idxNow := g.block.NewLoad(types.Double, idx)
	ok := g.block.NewFCmp(enum.FPredOLT, idxNow, lenF)
	g.block.NewCondBr(ok, body, miss)

	g.block = body
	idxBox := g.block.NewCall(g.runtimeBoxNumber, idxNow)
	el := g.block.NewCall(g.runtimeArrayGet, recvLive, g.emitAsKodaI64(idxBox))
	elSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(elSlot)
	g.block.NewStore(g.emitAsKodaI64(el), elSlot)
	elLive := g.block.NewLoad(types.I64, elSlot)
	predV, err := g.emitIndirectI64Callee(fnVal, nThis, []value.Value{elLive})
	if err != nil {
		return nil, true, err
	}
	brOk := g.emitTruthy(predV)
	g.block.NewCondBr(brOk, foundHit, step)

	g.block = foundHit
	foundHit.NewBr(merge)

	g.block = step
	one := constant.NewFloat(types.Double, 1)
	next := g.block.NewFAdd(g.block.NewLoad(types.Double, idx), one)
	g.block.NewStore(next, idx)
	g.block.NewBr(hdr)

	g.block = miss
	nilV := constant.NewInt(types.I64, llvmNilTagged)
	miss.NewBr(merge)

	g.block = merge
	return merge.NewPhi(
		ir.NewIncoming(elLive, foundHit),
		ir.NewIncoming(nilV, miss),
	), true, nil
}

func (g *Generator) emitArrayMethodReduce(recv value.Value, args []parser.Expr) (value.Value, bool, error) {
	if len(args) != 1 && len(args) != 2 {
		return nil, true, fmt.Errorf("reduce expects 1 or 2 arguments")
	}
	recvSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(recvSlot)
	g.block.NewStore(g.emitAsKodaI64(recv), recvSlot)

	fnVal, err := g.emitExpr(args[0])
	if err != nil {
		return nil, true, err
	}
	nThis := constant.NewInt(types.I64, 0)

	recvLive := g.block.NewLoad(types.I64, recvSlot)
	lenF := g.emitArrayLenAsDouble(recvLive)
	zero := constant.NewFloat(types.Double, 0)
	isEmpty := g.block.NewFCmp(enum.FPredOEQ, lenF, zero)

	g.tempN++
	suf := fmt.Sprintf(".red%d", g.tempN)
	emptyB := g.currentFn.NewBlock("red.empty" + suf)
	startB := g.currentFn.NewBlock("red.setup" + suf)
	mergeFinal := g.currentFn.NewBlock("red.merge" + suf)

	g.block.NewCondBr(isEmpty, emptyB, startB)

	var emptyOut value.Value
	g.block = emptyB
	if len(args) == 1 {
		emptyOut = constant.NewInt(types.I64, llvmNilTagged)
	} else {
		initV, err := g.emitExpr(args[1])
		if err != nil {
			return nil, true, err
		}
		emptyOut = g.emitAsKodaI64(initV)
	}
	emptyB.NewBr(mergeFinal)

	g.block = startB
	var accSlot value.Value
	var startIdx value.Value
	if len(args) == 1 {
		zBox := g.block.NewCall(g.runtimeBoxNumber, zero)
		firstEl := g.block.NewCall(g.runtimeArrayGet, recvLive, g.emitAsKodaI64(zBox))
		accSlot = g.entryAlloca(types.I64)
		g.shadowStoreTemp(accSlot)
		g.block.NewStore(g.emitAsKodaI64(firstEl), accSlot)
		startIdx = constant.NewFloat(types.Double, 1)
	} else {
		initV, err := g.emitExpr(args[1])
		if err != nil {
			return nil, true, err
		}
		accSlot = g.entryAlloca(types.I64)
		g.shadowStoreTemp(accSlot)
		g.block.NewStore(g.emitAsKodaI64(initV), accSlot)
		startIdx = zero
	}

	hdr := g.currentFn.NewBlock("red.hdr" + suf)
	body := g.currentFn.NewBlock("red.body" + suf)
	step := g.currentFn.NewBlock("red.step" + suf)
	done := g.currentFn.NewBlock("red.done" + suf)

	idx := g.entryAlloca(types.Double)
	g.block.NewStore(startIdx, idx)
	g.block.NewBr(hdr)

	g.block = hdr
	idxNow := g.block.NewLoad(types.Double, idx)
	ok := g.block.NewFCmp(enum.FPredOLT, idxNow, lenF)
	g.block.NewCondBr(ok, body, done)

	g.block = body
	idxBox := g.block.NewCall(g.runtimeBoxNumber, idxNow)
	el := g.block.NewCall(g.runtimeArrayGet, recvLive, g.emitAsKodaI64(idxBox))
	elSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(elSlot)
	g.block.NewStore(g.emitAsKodaI64(el), elSlot)
	elLive := g.block.NewLoad(types.I64, elSlot)
	accLoad := g.block.NewLoad(types.I64, accSlot)
	rv, err := g.emitIndirectI64Callee(fnVal, nThis, []value.Value{accLoad, elLive})
	if err != nil {
		return nil, true, err
	}
	g.block.NewStore(g.emitAsKodaI64(rv), accSlot)
	g.block.NewBr(step)

	g.block = step
	one := constant.NewFloat(types.Double, 1)
	next := g.block.NewFAdd(g.block.NewLoad(types.Double, idx), one)
	g.block.NewStore(next, idx)
	g.block.NewBr(hdr)

	g.block = done
	loopOut := g.block.NewLoad(types.I64, accSlot)
	done.NewBr(mergeFinal)

	g.block = mergeFinal
	return mergeFinal.NewPhi(
		ir.NewIncoming(emptyOut, emptyB),
		ir.NewIncoming(loopOut, done),
	), true, nil
}
