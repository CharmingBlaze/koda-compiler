package codegen

import (
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"koda/internal/parser"
)

func (g *Generator) ensureFastNativeFunc(symbol string) *ir.Func {
	if fn, ok := g.fastNativeFuncs[symbol]; ok {
		return fn
	}
	fn := g.ensureNativeExternFunc(symbol)
	if g.fastNativeFuncs == nil {
		g.fastNativeFuncs = make(map[string]*ir.Func)
	}
	g.fastNativeFuncs[symbol] = fn
	return fn
}

func (g *Generator) emitFastNativeCall(symbol string, argc int, args []value.Value) value.Value {
	fn := g.ensureFastNativeFunc(symbol)
	arrTy := types.NewArray(uint64(argc), types.I64)
	zero := constant.NewInt(types.I32, 0)
	slot := g.entryAlloca(arrTy)
	for i, arg := range args {
		argI64 := g.emitAsKodaI64(arg)
		elemPtr := g.block.NewGetElementPtr(arrTy, slot, zero, constant.NewInt(types.I32, int64(i)))
		g.block.NewStore(argI64, elemPtr)
	}
	argvPtr := g.block.NewGetElementPtr(arrTy, slot, zero, zero)
	g.block.NewCall(fn, constant.NewInt(types.I32, int64(argc)), argvPtr)
	return constant.NewInt(types.I64, 0)
}

func literalInt64FromExpr(e parser.Expr) (int64, bool) {
	lit, ok := e.(*parser.LiteralExpr)
	if !ok {
		return 0, false
	}
	switch v := lit.Value.(type) {
	case int:
		return int64(v), true
	case float64:
		return int64(v), true
	default:
		return 0, false
	}
}

func (g *Generator) structTypeNameForVar(varName string) string {
	if g.ctx == nil {
		return ""
	}
	if st, ok := g.ctx.VarStruct[varName]; ok {
		return st
	}
	if g.currentEmitFuncName != "" && g.ctx.FuncForOfVarStruct != nil {
		if vars, ok := g.ctx.FuncForOfVarStruct[g.currentEmitFuncName]; ok {
			if st, ok := vars[varName]; ok {
				return st
			}
		}
	}
	return ""
}

func (g *Generator) emitVec3FieldFromIdentifier(varName, field string) (value.Value, bool, error) {
	if g.ctx != nil && g.ctx.ConstPlainObjectLiterals != nil {
		if fields, ok := g.ctx.ConstPlainObjectLiterals[varName]; ok {
			if v, ok := fields[field]; ok {
				return g.emitBoxedIntLiteral(v), true, nil
			}
		}
	}
	if stName := g.structTypeNameForVar(varName); stName != "" {
		if slot, ok := g.structFieldSlotOf(stName, field); ok {
			objI, ok := g.loadBindingI64(varName)
			if !ok {
				return nil, false, nil
			}
			return g.block.NewCall(g.runtimeStructGet, objI, constant.NewInt(types.I64, int64(slot))), true, nil
		}
	}
	objI, ok := g.loadBindingI64(varName)
	if !ok {
		return nil, false, nil
	}
	key := g.emitStringLiteral(field)
	return g.block.NewCall(g.runtimeObjGet, objI, g.emitAsKodaI64(key)), true, nil
}

func (g *Generator) emitVec3ComponentsFromExpr(e parser.Expr) (px, py, pz value.Value, ok bool, err error) {
	if vx, vy, vz, vecOk := g.parseVec3ObjectLiteral(e); vecOk {
		px, err = g.emitNumericArg(vx)
		if err != nil {
			return nil, nil, nil, false, err
		}
		py, err = g.emitNumericArg(vy)
		if err != nil {
			return nil, nil, nil, false, err
		}
		pz, err = g.emitNumericArg(vz)
		if err != nil {
			return nil, nil, nil, false, err
		}
		return px, py, pz, true, nil
	}
	if call, callOk := e.(*parser.CallExpr); callOk {
		if id, idOk := call.Function.(*parser.IdentifierExpr); idOk && strings.EqualFold(id.Name.Lexeme, "vec3") && len(call.Arguments) == 3 {
			px, err = g.emitNumericArg(call.Arguments[0])
			if err != nil {
				return nil, nil, nil, false, err
			}
			py, err = g.emitNumericArg(call.Arguments[1])
			if err != nil {
				return nil, nil, nil, false, err
			}
			pz, err = g.emitNumericArg(call.Arguments[2])
			if err != nil {
				return nil, nil, nil, false, err
			}
			return px, py, pz, true, nil
		}
	}
	id, idOk := e.(*parser.IdentifierExpr)
	if !idOk {
		return nil, nil, nil, false, nil
	}
	name := id.Name.Lexeme
	px, okX, err := g.emitVec3FieldFromIdentifier(name, "x")
	if err != nil || !okX {
		return nil, nil, nil, false, err
	}
	py, okY, err := g.emitVec3FieldFromIdentifier(name, "y")
	if err != nil || !okY {
		return nil, nil, nil, false, err
	}
	pz, okZ, err := g.emitVec3FieldFromIdentifier(name, "z")
	if err != nil || !okZ {
		return nil, nil, nil, false, err
	}
	return px, py, pz, true, nil
}

func (g *Generator) parseVec3ObjectLiteral(e parser.Expr) (x, y, z parser.Expr, ok bool) {
	oe, ok := e.(*parser.ObjectExpr)
	if !ok || len(oe.ComputedKeys) != 0 {
		return nil, nil, nil, false
	}
	if oe.StructTag != nil && !strings.EqualFold(oe.StructTag.Lexeme, "vector3") {
		return nil, nil, nil, false
	}
	var fx, fy, fz parser.Expr
	for i, k := range oe.Keys {
		switch strings.ToLower(k.Lexeme) {
		case "x":
			fx = oe.Values[i]
		case "y":
			fy = oe.Values[i]
		case "z":
			fz = oe.Values[i]
		}
	}
	if fx == nil || fy == nil || fz == nil {
		return nil, nil, nil, false
	}
	return fx, fy, fz, true
}

func (g *Generator) colorComponentsFromObjectExpr(e parser.Expr) (r, gch, b, a int64, ok bool) {
	oe, ok := e.(*parser.ObjectExpr)
	if !ok || oe.StructTag != nil || len(oe.ComputedKeys) != 0 {
		return 0, 0, 0, 0, false
	}
	fields := make(map[string]int64)
	for i, k := range oe.Keys {
		v, ok := literalInt64FromExpr(oe.Values[i])
		if !ok {
			return 0, 0, 0, 0, false
		}
		fields[strings.ToLower(k.Lexeme)] = v
	}
	r, okR := fields["r"]
	gch, okG := fields["g"]
	b, okB := fields["b"]
	a = int64(255)
	if av, okA := fields["a"]; okA {
		a = av
	}
	if !okR || !okG || !okB {
		return 0, 0, 0, 0, false
	}
	return r, gch, b, a, true
}

func (g *Generator) colorComponentsFromIdentifier(name string) (r, gch, b, a int64, ok bool) {
	if g.ctx == nil || g.ctx.ConstPlainObjectLiterals == nil {
		return 0, 0, 0, 0, false
	}
	fields, ok := g.ctx.ConstPlainObjectLiterals[name]
	if !ok {
		return 0, 0, 0, 0, false
	}
	r, okR := fields["r"]
	gch, okG := fields["g"]
	b, okB := fields["b"]
	if !okR || !okG || !okB {
		return 0, 0, 0, 0, false
	}
	a = int64(255)
	if av, okA := fields["a"]; okA {
		a = av
	}
	return r, gch, b, a, true
}

func (g *Generator) emitBoxedIntLiteral(v int64) value.Value {
	return g.block.NewCall(g.runtimeBoxNumber, constant.NewFloat(types.Double, float64(v)))
}

func (g *Generator) emitNumericArg(e parser.Expr) (value.Value, error) {
	return g.emitExpr(e)
}

func (g *Generator) tryEmitFusedDrawCube(name string, args []parser.Expr) (value.Value, bool, error) {
	if len(args) != 5 {
		return nil, false, nil
	}
	fastSym := "koda_fast_DrawCube"
	fast8Sym := "koda_fast_DrawCube8"
	if strings.EqualFold(name, "DrawCubeWires") {
		fastSym = "koda_fast_DrawCubeWires"
		fast8Sym = "koda_fast_DrawCubeWires8"
	} else if !strings.EqualFold(name, "DrawCube") {
		return nil, false, nil
	}

	px, py, pz, ok, err := g.emitVec3ComponentsFromExpr(args[0])
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return nil, false, nil
	}
	w, err := g.emitNumericArg(args[1])
	if err != nil {
		return nil, false, err
	}
	h, err := g.emitNumericArg(args[2])
	if err != nil {
		return nil, false, err
	}
	l, err := g.emitNumericArg(args[3])
	if err != nil {
		return nil, false, err
	}

	var cr, cg, cb, ca int64
	var colorPacked bool
	switch c := args[4].(type) {
	case *parser.ObjectExpr:
		cr, cg, cb, ca, colorPacked = g.colorComponentsFromObjectExpr(c)
	case *parser.IdentifierExpr:
		cr, cg, cb, ca, colorPacked = g.colorComponentsFromIdentifier(c.Name.Lexeme)
	}

	if colorPacked {
		nargs := []value.Value{
			px, py, pz, w, h, l,
			g.emitBoxedIntLiteral(cr),
			g.emitBoxedIntLiteral(cg),
			g.emitBoxedIntLiteral(cb),
			g.emitBoxedIntLiteral(ca),
		}
		return g.emitFastNativeCall(fastSym, 10, nargs), true, nil
	}

	col, err := g.emitExpr(args[4])
	if err != nil {
		return nil, false, err
	}
	nargs := []value.Value{px, py, pz, w, h, l, col}
	return g.emitFastNativeCall(fast8Sym, 7, nargs), true, nil
}

func (g *Generator) canonicalStructName(name string) string {
	if g.ctx == nil || name == "" {
		return name
	}
	for k := range g.ctx.StructFields {
		if strings.EqualFold(k, name) {
			return k
		}
	}
	return name
}

func (g *Generator) isCamera3DStructType(name string) bool {
	return strings.EqualFold(name, "camera3d")
}

func structLiteralFieldExpr(oe *parser.ObjectExpr, fname string) (parser.Expr, bool) {
	for i, k := range oe.Keys {
		if strings.EqualFold(k.Lexeme, fname) {
			return oe.Values[i], true
		}
	}
	return nil, false
}

func (g *Generator) emitVec3ComponentsFromStructSlot(structVal value.Value, stName, field string) (px, py, pz value.Value, ok bool, err error) {
	stName = g.canonicalStructName(stName)
	slot, ok := g.structFieldSlotOf(stName, field)
	if !ok {
		return nil, nil, nil, false, nil
	}
	vecI := g.block.NewCall(g.runtimeStructGet, structVal, constant.NewInt(types.I64, int64(slot)))
	// Field may be a vec3 object ({x,y,z}) or a nested Vector3 struct.
	if g.ctx != nil && g.ctx.StructFieldTypes != nil {
		if fields, ok := g.ctx.StructFieldTypes[stName]; ok {
			if ft, ok := fields[field]; ok && strings.EqualFold(ft, "vector3") {
				vec3Name := g.canonicalStructName("Vector3")
				for i, axis := range []string{"x", "y", "z"} {
					axSlot, ok := g.structFieldSlotOf(vec3Name, axis)
					if !ok {
						return nil, nil, nil, false, nil
					}
					v := g.block.NewCall(g.runtimeStructGet, vecI, constant.NewInt(types.I64, int64(axSlot)))
					switch i {
					case 0:
						px = v
					case 1:
						py = v
					case 2:
						pz = v
					}
				}
				return px, py, pz, true, nil
			}
		}
	}
	for i, axis := range []string{"x", "y", "z"} {
		key := g.emitStringLiteral(axis)
		v := g.block.NewCall(g.runtimeObjGet, vecI, g.emitAsKodaI64(key))
		switch i {
		case 0:
			px = v
		case 1:
			py = v
		case 2:
			pz = v
		}
	}
	return px, py, pz, true, nil
}

func (g *Generator) tryEmitFusedBeginMode3DFromCamera3D(arg parser.Expr) (value.Value, bool, error) {
	vecFields := []string{"position", "target", "up"}
	vals := make([]value.Value, 0, 10)

	switch e := arg.(type) {
	case *parser.ObjectExpr:
		if e.StructTag == nil || !g.isCamera3DStructType(e.StructTag.Lexeme) {
			return nil, false, nil
		}
		for _, fname := range vecFields {
			fe, ok := structLiteralFieldExpr(e, fname)
			if !ok {
				return nil, false, nil
			}
			px, py, pz, ok, err := g.emitVec3ComponentsFromExpr(fe)
			if err != nil {
				return nil, false, err
			}
			if !ok {
				return nil, false, nil
			}
			vals = append(vals, px, py, pz)
		}
		fovyExpr, ok := structLiteralFieldExpr(e, "fovy")
		if !ok {
			fovyExpr = &parser.LiteralExpr{Value: float64(45)}
		}
		fovy, err := g.emitNumericArg(fovyExpr)
		if err != nil {
			return nil, false, err
		}
		vals = append(vals, fovy)
	case *parser.IdentifierExpr:
		stName := g.canonicalStructName(g.structTypeNameForVar(e.Name.Lexeme))
		if !g.isCamera3DStructType(stName) {
			return nil, false, nil
		}
		camI, ok := g.loadBindingI64(e.Name.Lexeme)
		if !ok {
			return nil, false, nil
		}
		for _, fname := range vecFields {
			px, py, pz, ok, err := g.emitVec3ComponentsFromStructSlot(camI, stName, fname)
			if err != nil {
				return nil, false, err
			}
			if !ok {
				return nil, false, nil
			}
			vals = append(vals, px, py, pz)
		}
		slot, ok := g.structFieldSlotOf(stName, "fovy")
		if !ok {
			return nil, false, nil
		}
		vals = append(vals, g.block.NewCall(g.runtimeStructGet, camI, constant.NewInt(types.I64, int64(slot))))
	default:
		return nil, false, nil
	}
	return g.emitFastNativeCall("koda_fast_BeginMode3D", 10, vals), true, nil
}

func (g *Generator) tryEmitFusedBeginMode3D(args []parser.Expr) (value.Value, bool, error) {
	if len(args) != 1 {
		return nil, false, nil
	}
	if v, ok, err := g.tryEmitFusedBeginMode3DFromCamera3D(args[0]); ok || err != nil {
		return v, ok, err
	}
	call, ok := args[0].(*parser.CallExpr)
	if !ok {
		return nil, false, nil
	}
	id, ok := call.Function.(*parser.IdentifierExpr)
	if !ok || !strings.EqualFold(id.Name.Lexeme, "camera3d") {
		return nil, false, nil
	}
	if len(call.Arguments) != 7 {
		return nil, false, nil
	}
	vals := make([]value.Value, 0, 10)
	for i := 0; i < 6; i++ {
		v, err := g.emitNumericArg(call.Arguments[i])
		if err != nil {
			return nil, false, err
		}
		vals = append(vals, v)
	}
	fovy, err := g.emitNumericArg(call.Arguments[6])
	if err != nil {
		return nil, false, err
	}
	vals = append(vals,
		g.emitBoxedIntLiteral(0),
		g.emitBoxedIntLiteral(1),
		g.emitBoxedIntLiteral(0),
		fovy,
	)
	return g.emitFastNativeCall("koda_fast_BeginMode3D", 10, vals), true, nil
}

func (g *Generator) tryEmitFusedColorDraw(name string, argc int, args []parser.Expr) (value.Value, bool, error) {
	if len(args) != argc {
		return nil, false, nil
	}
	vals := make([]value.Value, argc)
	for i := 0; i < argc-1; i++ {
		v, err := g.emitNumericArg(args[i])
		if err != nil {
			return nil, false, err
		}
		vals[i] = v
	}
	col, err := g.emitExpr(args[argc-1])
	if err != nil {
		return nil, false, err
	}
	vals[argc-1] = col
	return g.emitFastNativeCall(name, argc, vals), true, nil
}

func (g *Generator) tryEmitFusedDrawText(args []parser.Expr) (value.Value, bool, error) {
	if len(args) != 5 {
		return nil, false, nil
	}
	text, err := g.emitExpr(args[0])
	if err != nil {
		return nil, false, err
	}
	vals := []value.Value{text}
	for i := 1; i < 4; i++ {
		v, err := g.emitNumericArg(args[i])
		if err != nil {
			return nil, false, err
		}
		vals = append(vals, v)
	}
	col, err := g.emitExpr(args[4])
	if err != nil {
		return nil, false, err
	}
	vals = append(vals, col)
	return g.emitFastNativeCall("koda_fast_DrawText", 5, vals), true, nil
}

// tryEmitFusedNativeCall recognizes hot Raylib patterns and emits numeric fast paths.
func (g *Generator) tryEmitFusedNativeCall(name string, args []parser.Expr) (value.Value, bool, error) {
	if strings.EqualFold(name, "BeginMode3D") {
		return g.tryEmitFusedBeginMode3D(args)
	}
	if strings.EqualFold(name, "DrawCube") || strings.EqualFold(name, "DrawCubeWires") {
		return g.tryEmitFusedDrawCube(name, args)
	}
	if strings.EqualFold(name, "ClearBackground") {
		if len(args) != 1 {
			return nil, false, nil
		}
		col, err := g.emitExpr(args[0])
		if err != nil {
			return nil, false, err
		}
		return g.emitFastNativeCall("koda_fast_ClearBackground", 1, []value.Value{col}), true, nil
	}
	if strings.EqualFold(name, "DrawRectangle") {
		return g.tryEmitFusedColorDraw("koda_fast_DrawRectangle", 5, args)
	}
	if strings.EqualFold(name, "DrawCircle") {
		return g.tryEmitFusedColorDraw("koda_fast_DrawCircle", 4, args)
	}
	if strings.EqualFold(name, "DrawCircleLines") {
		return g.tryEmitFusedColorDraw("koda_fast_DrawCircleLines", 4, args)
	}
	if strings.EqualFold(name, "DrawLine") {
		return g.tryEmitFusedColorDraw("koda_fast_DrawLine", 5, args)
	}
	if strings.EqualFold(name, "DrawRectangleLines") {
		return g.tryEmitFusedColorDraw("koda_fast_DrawRectangleLines", 5, args)
	}
	if strings.EqualFold(name, "DrawLine3D") {
		return g.tryEmitFusedColorDraw("koda_fast_DrawLine3D", 7, args)
	}
	if strings.EqualFold(name, "DrawText") {
		return g.tryEmitFusedDrawText(args)
	}
	return nil, false, nil
}
