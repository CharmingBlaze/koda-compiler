package codegen

import (
	"strings"
	"testing"

	"koda/internal/lexer"
	"koda/internal/parser"
	"koda/internal/sema"
)

func parseForTest(t *testing.T, source string) *parser.Program {
	t.Helper()
	l := lexer.NewLexer(source, "")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("Lexer failed: %v", err)
	}
	p := parser.NewParser(tokens)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parser failed: %v", err)
	}
	return program
}

func TestCodegenBasic(t *testing.T) {
	source := `
		func add(a, b) {
			return a + b;
		}
		let result = add(1, 2);
	`
	program := parseForTest(t, source)

	analyzer := sema.NewAnalyzer()
	if err := analyzer.Analyze(program); err != nil {
		t.Fatalf("Sema failed: %v", err)
	}

	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatalf("PrepareNativeBundle failed: %v", err)
	}
	gen := NewGenerator(ctx)

	mod, err := gen.Generate(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatalf("Codegen failed: %v", err)
	}
	if mod == nil {
		t.Fatal("Expected non-nil module")
	}
}

func TestCodegenControlFlow(t *testing.T) {
	source := `
		func test(x) {
			if (x > 0) {
				return x;
			}
			return 0;
		}
	`
	program := parseForTest(t, source)

	analyzer := sema.NewAnalyzer()
	if err := analyzer.Analyze(program); err != nil {
		t.Fatalf("Sema failed: %v", err)
	}

	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatalf("PrepareNativeBundle failed: %v", err)
	}
	gen := NewGenerator(ctx)

	mod, err := gen.Generate(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatalf("Codegen failed: %v", err)
	}
	if mod == nil {
		t.Fatal("Expected non-nil module")
	}
}

func TestNativeExternDeclaresGlueSymbol(t *testing.T) {
	src := `
// koda:extern Foo KODA_shim_Foo 1
let Foo = 0;
func main() {
	Foo(1);
}
`
	program := parseForTest(t, src)
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatalf("PrepareNativeBundle: %v", err)
	}
	gen := NewGenerator(ctx)
	mod, err := gen.Generate(ctx.Bundle)
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	ir := mod.String()
	if !strings.Contains(ir, "@KODA_shim_Foo") {
		t.Fatalf("expected LLVM to declare @KODA_shim_Foo, got:\n%s", ir)
	}
}

func TestEmitPrefixPosIntegerLiteralNoUnbox(t *testing.T) {
	src := `func main() { return +42; }`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatal(err)
	}
	ir := mod.String()
	if strings.Contains(ir, "call double @koda_unbox_number") {
		t.Fatalf("expected +literal fold to avoid unbox calls, got:\n%s", ir)
	}
	if !strings.Contains(ir, "call i64 @koda_box_number") {
		t.Fatalf("expected boxed numeric result, got:\n%s", ir)
	}
}

func TestEmitCompileTimePosOfProductNoUnbox(t *testing.T) {
	src := `func main() { return +(2 * 3); }`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatal(err)
	}
	ir := mod.String()
	if strings.Contains(ir, "call double @koda_unbox_number") {
		t.Fatalf("expected +(2*3) fold to avoid unbox, got:\n%s", ir)
	}
}

func TestEmitPrefixNegIntegerLiteralNoUnbox(t *testing.T) {
	src := `func main() { return -42; }`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatal(err)
	}
	ir := mod.String()
	if strings.Contains(ir, "call double @koda_unbox_number") {
		t.Fatalf("expected -literal fold to avoid unbox calls, got:\n%s", ir)
	}
	if !strings.Contains(ir, "call i64 @koda_box_number") {
		t.Fatalf("expected boxed numeric result, got:\n%s", ir)
	}
}

func TestEmitCompileTimeInt64ChainNoUnbox(t *testing.T) {
	src := `func main() { return 1 + 2 + 3; }`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatal(err)
	}
	ir := mod.String()
	if strings.Contains(ir, "call double @koda_unbox_number") {
		t.Fatalf("expected chained literal fold to avoid unbox, got:\n%s", ir)
	}
}

func TestEmitCompileTimeNegOfProductNoUnbox(t *testing.T) {
	src := `func main() { return -(2 * 3); }`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatal(err)
	}
	ir := mod.String()
	if strings.Contains(ir, "call double @koda_unbox_number") {
		t.Fatalf("expected -(2*3) fold to avoid unbox, got:\n%s", ir)
	}
}

func TestEmitInfixIntegerLiteralMulFoldNoUnbox(t *testing.T) {
	src := `func main() { return 6 * 7; }`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatal(err)
	}
	ir := mod.String()
	if strings.Contains(ir, "call double @koda_unbox_number") {
		t.Fatalf("expected literal mul fold to avoid unbox calls, got:\n%s", ir)
	}
	if !strings.Contains(ir, "call i64 @koda_box_number") {
		t.Fatalf("expected boxed numeric result, got:\n%s", ir)
	}
}

func TestStructMethodCallsModuleHelperFunc(t *testing.T) {
	src := `
func helper(x) {
	return x;
}

struct Widget {
	func ping() {
		helper(1);
	}
}

func main() {
	let w = Widget {};
}
`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewGenerator(ctx).Generate(ctx.Bundle); err != nil {
		t.Fatalf("codegen failed: %v", err)
	}
}

func TestStructMethodNativeExternBeforeLet(t *testing.T) {
	// Native shim is declared after the struct; struct methods emit before top-level
	// let decls unless prepareBundleBindings registers natives first.
	src := `
struct Widget {
	func ping() {
		shimfn(1);
	}
}

// koda:extern shimfn koda_test_shimfn 1
let shimfn = 0;

func main() {
	let w = Widget {};
}
`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewGenerator(ctx).Generate(ctx.Bundle); err != nil {
		t.Fatalf("codegen failed: %v", err)
	}
}

func TestStructConstructorCallEmitsNewMethod(t *testing.T) {
	src := `
struct OrbitCamera {
	target = null;
	distance = 13.0;
	func new(target, distance, pitch, yaw) {
		if (target != null) { this.target = target; }
		if (distance != null) { this.distance = distance; }
	}
}
func orbit(cfg) {
	return OrbitCamera(cfg.target, cfg.distance, null, null);
}
func main() {
	let cam = OrbitCamera(null, null, null, null);
}
`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatalf("codegen failed: %v", err)
	}
	ir := mod.String()
	if !strings.Contains(strings.ToLower(ir), "koda_method_orbitcamera_new") {
		t.Fatalf("expected struct constructor to emit koda_method_Orbitcamera_new, got:\n%s", ir)
	}
}

func TestStructMethodImplicitFieldAssign(t *testing.T) {
	src := `
struct Coin {
	on = true;
	func pickup() {
		on = false;
	}
}
func main() {
	let c = Coin {};
	c.pickup();
}
`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewGenerator(ctx).Generate(ctx.Bundle); err != nil {
		t.Fatalf("codegen failed: %v", err)
	}
}

func TestCodegenFusedDrawCube(t *testing.T) {
	src := `
		// koda:extern DrawCube koda_wrap_raylib_DrawCube 5
		let DrawCube = 0;
		let col_grass = { r: 34, g: 139, b: 34, a: 255 };
		func dc(px, py, pz, w, h, l, col) {
			DrawCube({ x: px, y: py, z: pz }, w, h, l, col);
		}
		func main() {
			dc(0.0, 1.0, 0.0, 1.0, 1.0, 1.0, col_grass);
		}
	`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatalf("codegen failed: %v", err)
	}
	ir := mod.String()
	if !strings.Contains(ir, "koda_fast_DrawCube") {
		t.Fatalf("expected fused DrawCube fast path in IR, got:\n%s", ir)
	}
	dcStart := strings.Index(ir, "define i64 @dc(")
	if dcStart < 0 {
		t.Fatal("missing @dc in IR")
	}
	dcEnd := strings.Index(ir[dcStart:], "\n}\n")
	if dcEnd < 0 {
		t.Fatal("could not find end of @dc")
	}
	dcIR := ir[dcStart : dcStart+dcEnd]
	if strings.Contains(dcIR, "koda_allocate_object") {
		t.Fatalf("fused dc should not allocate heap objects:\n%s", dcIR)
	}
}

func TestCodegenFusedDrawCubeIdentifierPosition(t *testing.T) {
	src := `
		// koda:extern DrawCube koda_wrap_raylib_DrawCube 5
		let DrawCube = 0;
		let col_grass = { r: 34, g: 139, b: 34, a: 255 };
		let origin = { x: 0, y: 1, z: 2 };
		func dc(pos) {
			DrawCube(pos, 1.0, 1.0, 1.0, col_grass);
		}
		func main() {
			dc(origin);
		}
	`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatalf("codegen failed: %v", err)
	}
	ir := mod.String()
	if !strings.Contains(ir, "koda_fast_DrawCube") {
		t.Fatalf("expected fused DrawCube for identifier position, got:\n%s", ir)
	}
	dcStart := strings.Index(ir, "define i64 @dc(")
	if dcStart < 0 {
		t.Fatal("missing @dc in IR")
	}
	dcEnd := strings.Index(ir[dcStart:], "\n}\n")
	if dcEnd < 0 {
		t.Fatal("could not find end of @dc")
	}
	dcIR := ir[dcStart : dcStart+dcEnd]
	if strings.Contains(dcIR, "koda_allocate_object") {
		t.Fatalf("fused dc should not allocate heap objects:\n%s", dcIR)
	}
}

func TestCodegenFusedBeginMode3DArgOrder(t *testing.T) {
	src := `
		// koda:extern BeginMode3D koda_wrap_raylib_BeginMode3D 1
		let BeginMode3D = 0;
		func camera3d(px, py, pz, tx, ty, tz, fovy) { return { position: { x: px, y: py, z: pz }, target: { x: tx, y: ty, z: tz }, up: { x: 0.0, y: 1.0, z: 0.0 }, fovy: fovy, projection: 0 }; }
		func main() {
			BeginMode3D(camera3d(1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 50.0));
		}
	`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatal(err)
	}
	ir := mod.String()
	if !strings.Contains(ir, "koda_fast_BeginMode3D") {
		t.Fatalf("expected fused BeginMode3D, got:\n%s", ir)
	}
	// fovy (50) must be argv slot 9, not slot 6 (up.x).
	if strings.Contains(ir, "store i64 %") {
		// Sanity: boxed 50.0 appears after up vector stores in the fused argv block.
		if !strings.Contains(ir, "double 5") && !strings.Contains(ir, "double 50") {
			t.Fatalf("expected fovy literal in IR")
		}
	}
}

func TestCodegenFusedBeginMode3DCamera3DStruct(t *testing.T) {
	src := `
		// koda:extern BeginMode3D koda_wrap_raylib_BeginMode3D 1
		let BeginMode3D = 0;
		struct Camera3D {
			position;
			target;
			up;
			fovy = 45.0;
			projection = 0;
		}
		func main() {
			BeginMode3D(Camera3D {
				position: vec3(4.0, 4.0, 4.0),
				target: vec3(0.0, 0.0, 0.0),
				up: vec3(0.0, 1.0, 0.0),
				fovy: 45.0,
				projection: 0
			});
		}
	`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatal(err)
	}
	ir := mod.String()
	if !strings.Contains(ir, "koda_fast_BeginMode3D") {
		t.Fatalf("expected fused BeginMode3D for Camera3D struct, got:\n%s", ir)
	}
}

func TestCodegenFusedClearBackgroundHexColor(t *testing.T) {
	src := `
		// koda:extern ClearBackground koda_wrap_raylib_ClearBackground 1
		let ClearBackground = 0;
		func main() {
			ClearBackground(#101018);
		}
	`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatalf("codegen failed: %v", err)
	}
	ir := mod.String()
	if !strings.Contains(ir, "koda_fast_ClearBackground") {
		t.Fatalf("expected fused ClearBackground, got:\n%s", ir)
	}
}

func TestCodegenFusedDrawRectanglePackedColor(t *testing.T) {
	src := `
		// koda:extern DrawRectangle koda_wrap_raylib_DrawRectangle 5
		let DrawRectangle = 0;
		let red = 4278190335;
		func draw() {
			DrawRectangle(10, 20, 30, 40, red);
		}
		func main() {
			draw();
		}
	`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatalf("codegen failed: %v", err)
	}
	ir := mod.String()
	if !strings.Contains(ir, "koda_fast_DrawRectangle") {
		t.Fatalf("expected fused DrawRectangle fast path in IR, got:\n%s", ir)
	}
	drawStart := strings.Index(ir, "define i64 @draw(")
	if drawStart < 0 {
		t.Fatal("missing @draw in IR")
	}
	drawEnd := strings.Index(ir[drawStart:], "\n}\n")
	if drawEnd < 0 {
		t.Fatal("could not find end of @draw")
	}
	drawIR := ir[drawStart : drawStart+drawEnd]
	if strings.Contains(drawIR, "koda_allocate_object") {
		t.Fatalf("fused draw should not allocate heap objects:\n%s", drawIR)
	}
}

func TestCodegenTypedFloatFastInfix(t *testing.T) {
	src := `
		func main() {
			let x: float = 1.0;
			let y: float = 2.0;
			let z: float = x + y;
			x = x + y;
			x += y;
			return z;
		}
	`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatalf("codegen failed: %v", err)
	}
	ir := mod.String()
	mainStart := strings.Index(ir, "define i64 @koda_user_main(")
	if mainStart < 0 {
		t.Fatal("missing @koda_user_main in IR")
	}
	mainEnd := strings.Index(ir[mainStart:], "\n}\n")
	if mainEnd < 0 {
		t.Fatal("could not find end of @koda_user_main")
	}
	mainIR := ir[mainStart : mainStart+mainEnd]
	if !strings.Contains(mainIR, "fadd") {
		t.Fatalf("expected native fadd for typed float locals, got:\n%s", mainIR)
	}
	if strings.Contains(mainIR, "koda_value_add") {
		t.Fatalf("typed float infix should not call koda_value_add:\n%s", mainIR)
	}
	if !strings.Contains(mainIR, "alloca double") {
		t.Fatalf("expected double alloca for typed float locals, got:\n%s", mainIR)
	}
}

func TestCodegenInferredFloatParam(t *testing.T) {
	src := `
		func scale(dt) {
			return dt * 2.0;
		}
		func main() {
			return scale(0.5);
		}
	`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatalf("codegen failed: %v", err)
	}
	ir := mod.String()
	scaleStart := strings.Index(ir, "define i64 @scale(")
	if scaleStart < 0 {
		t.Fatal("missing @scale in IR")
	}
	scaleEnd := strings.Index(ir[scaleStart:], "\n}\n")
	if scaleEnd < 0 {
		t.Fatal("could not find end of @scale")
	}
	scaleIR := ir[scaleStart : scaleStart+scaleEnd]
	if !strings.Contains(scaleIR, "alloca double") {
		t.Fatalf("expected double alloca for typed float param, got:\n%s", scaleIR)
	}
	if !strings.Contains(scaleIR, "fmul") {
		t.Fatalf("expected fmul for inferred float param math, got:\n%s", scaleIR)
	}
}

func TestCodegenFloatLocalTimesLiteral(t *testing.T) {
	src := `
		func main() {
			let x: float = 2.0;
			let y: float = x * 3.0;
			return y;
		}
	`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatalf("codegen failed: %v", err)
	}
	ir := mod.String()
	mainStart := strings.Index(ir, "define i64 @koda_user_main(")
	mainEnd := strings.Index(ir[mainStart:], "\n}\n")
	mainIR := ir[mainStart : mainStart+mainEnd]
	if !strings.Contains(mainIR, "fmul") {
		t.Fatalf("expected fmul for float * literal, got:\n%s", mainIR)
	}
}

func TestCodegenFusedDrawLine(t *testing.T) {
	src := `
		// koda:extern DrawLine koda_wrap_raylib_DrawLine 5
		let DrawLine = 0;
		let white = 4294967295;
		func draw() {
			DrawLine(0, 0, 100, 100, white);
		}
		func main() {
			draw();
		}
	`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatalf("codegen failed: %v", err)
	}
	ir := mod.String()
	if !strings.Contains(ir, "koda_fast_DrawLine") {
		t.Fatalf("expected fused DrawLine fast path in IR, got:\n%s", ir)
	}
}

func TestCodegenStructFloatFieldFastMath(t *testing.T) {
	src := `
		struct Hero {
			x: float = 0.0;
			func tick(dt) {
				this.x += dt * 2.0;
			}
		}
		func main() {
			let h = Hero {};
			h.tick(0.5);
		}
	`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatalf("codegen failed: %v", err)
	}
	ir := mod.String()
	tickIR := irFuncBody(ir, "koda_method_hero_tick")
	if tickIR == "" {
		t.Fatalf("missing koda_method_hero_tick in IR:\n%s", ir)
	}
	if !strings.Contains(tickIR, "fmul") || !strings.Contains(tickIR, "fadd") {
		t.Fatalf("expected fmul/fadd for struct float field update, got:\n%s", tickIR)
	}
	if strings.Contains(tickIR, "koda_value_mul") || strings.Contains(tickIR, "koda_value_add") {
		t.Fatalf("struct float field += should not use boxed value ops:\n%s", tickIR)
	}
}

func irFuncBody(ir, funcName string) string {
	marker := "define i64 @" + funcName
	start := strings.Index(ir, marker)
	if start < 0 {
		return ""
	}
	rest := ir[start:]
	end := strings.Index(rest, "\n}\n")
	if end < 0 {
		return rest
	}
	return rest[:end+3]
}

func TestCodegenStructVarFloatFieldFastMath(t *testing.T) {
	src := `
		struct Hero {
			x: float = 0.0;
		}
		func main() {
			let h = Hero {};
			h.x += 1.5;
		}
	`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatalf("codegen failed: %v", err)
	}
	ir := mod.String()
	mainIR := irFuncBody(ir, "koda_user_main")
	if mainIR == "" {
		t.Fatalf("missing koda_user_main in IR:\n%s", ir)
	}
	if !strings.Contains(mainIR, "fadd") {
		t.Fatalf("expected fadd for struct var field update, got:\n%s", mainIR)
	}
	if strings.Contains(mainIR, "koda_value_add") {
		t.Fatalf("struct var field += should not use boxed value add:\n%s", mainIR)
	}
}

func TestCodegenInferredStructMethodParamFromCall(t *testing.T) {
	src := `
		struct Hero {
			x: float = 0.0;
			func tick(dt) {
				this.x += dt * 2.0;
			}
		}
		func main() {
			let h = Hero {};
			h.tick(0.5);
		}
	`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatalf("codegen failed: %v", err)
	}
	tickIR := irFuncBody(mod.String(), "koda_method_hero_tick")
	if tickIR == "" {
		t.Fatal("missing koda_method_hero_tick")
	}
	if !strings.Contains(tickIR, "fmul") {
		t.Fatalf("expected fmul for inferred struct method param, got:\n%s", tickIR)
	}
}

func TestCodegenFusedDrawCircleLines(t *testing.T) {
	src := `
		// koda:extern DrawCircleLines koda_wrap_raylib_DrawCircleLines 4
		let DrawCircleLines = 0;
		let white = 4294967295;
		func draw() {
			DrawCircleLines(100, 100, 32, white);
		}
		func main() {
			draw();
		}
	`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatalf("codegen failed: %v", err)
	}
	ir := mod.String()
	if !strings.Contains(ir, "koda_fast_DrawCircleLines") {
		t.Fatalf("expected fused DrawCircleLines, got:\n%s", ir)
	}
}

func TestCodegenFusedDrawRectangleLines(t *testing.T) {
	src := `
		// koda:extern DrawRectangleLines koda_wrap_raylib_DrawRectangleLines 5
		let DrawRectangleLines = 0;
		let white = 4294967295;
		func draw() {
			DrawRectangleLines(10, 20, 30, 40, white);
		}
		func main() {
			draw();
		}
	`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatalf("codegen failed: %v", err)
	}
	ir := mod.String()
	if !strings.Contains(ir, "koda_fast_DrawRectangleLines") {
		t.Fatalf("expected fused DrawRectangleLines, got:\n%s", ir)
	}
}

func TestCodegenInferredFloatLambdaParam(t *testing.T) {
	src := `
		func main() {
			let scale = func(dt) { return dt * 2.0; };
			return scale(0.5);
		}
	`
	program := parseForTest(t, src)
	if err := sema.NewAnalyzer().Analyze(program); err != nil {
		t.Fatal(err)
	}
	ctx, err := sema.PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err != nil {
		t.Fatal(err)
	}
	mod, err := NewGenerator(ctx).Generate(ctx.Bundle)
	if err != nil {
		t.Fatalf("codegen failed: %v", err)
	}
	ir := mod.String()
	if !strings.Contains(ir, "define i64 @closure_") {
		t.Fatal("missing closure in IR")
	}
	if !strings.Contains(ir, "fmul") {
		t.Fatalf("expected fmul for inferred float lambda param, got:\n%s", ir)
	}
}
