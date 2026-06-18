package sema

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"koda/internal/diagnostic"
	"koda/internal/lexer"
	"koda/internal/parser"
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

func TestSemaDuplicateBindingSameScope(t *testing.T) {
	src := `let x = 1;
let X = 2;
`
	program := parseForTest(t, src)
	err := NewAnalyzer().Analyze(program)
	if err == nil {
		t.Fatal("expected duplicate binding error")
	}
}

func TestSemaBasicScoping(t *testing.T) {
	source := `
		let x = 10;
		func add(a, b) {
			let y = a + b;
			return y;
		}
		add(x, 2);
	`
	program := parseForTest(t, source)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(program)
	if err != nil {
		t.Fatalf("Sema failed: %v", err)
	}
	if len(analyzer.Errors()) != 0 {
		t.Fatalf("Unexpected errors: %v", analyzer.Errors())
	}
}

func TestSemaUndefinedVariable(t *testing.T) {
	source := `
		let x = 10;
		y + 5;
	`
	program := parseForTest(t, source)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(program)
	if err == nil {
		t.Fatal("Expected error for undefined variable")
	}
	if len(analyzer.Errors()) == 0 {
		t.Fatal("Expected errors to be collected")
	}
}

func TestSemaBlockScoping(t *testing.T) {
	source := `
		let x = 10;
		if (x > 5) {
			let y = x;
			y = y + 1;
		}
	`
	program := parseForTest(t, source)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(program)
	if err != nil {
		t.Fatalf("Sema failed: %v", err)
	}
}

func TestSemaFunctionParams(t *testing.T) {
	source := `
		func f(a, b = 2, ...rest) {
			return a;
		}
	`
	program := parseForTest(t, source)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(program)
	if err != nil {
		t.Fatalf("Sema failed: %v", err)
	}
}

func TestSemaAssignmentTarget(t *testing.T) {
	source := `
		let x = 10;
		x = 20;
	`
	program := parseForTest(t, source)

	analyzer := NewAnalyzer()
	err := analyzer.Analyze(program)
	if err != nil {
		t.Fatalf("Sema failed: %v", err)
	}
}

func TestPrepareNativeBundleRunsSemaEarly(t *testing.T) {
	src := `let x = 1;
y + 1;
`
	program := parseForTest(t, src)
	_, err := PrepareNativeBundle(&parser.ProgramBundle{Entry: program})
	if err == nil {
		t.Fatal("expected undefined y")
	}
	var de *diagnostic.DiagnosticError
	if !errors.As(err, &de) || de.File != "<entry>" {
		t.Fatalf("want DiagnosticError with file, got %v", err)
	}
}

func TestSemaCallArityTooFew(t *testing.T) {
	src := `func draw(x, y, sprite) { return x; }
draw(1);
`
	program := parseForTest(t, src)
	a := NewAnalyzer()
	err := a.Analyze(program)
	if err == nil {
		t.Fatal("expected arity error")
	}
	if len(a.Errors()) == 0 {
		t.Fatal("expected recorded errors")
	}
}

func TestSemaCallArityTooMany(t *testing.T) {
	src := `func f(a) { return a; }
f(1, 2, 3);
`
	program := parseForTest(t, src)
	a := NewAnalyzer()
	err := a.Analyze(program)
	if err == nil {
		t.Fatal("expected arity error")
	}
}

func TestSemaCallArityRestAllowsExtra(t *testing.T) {
	src := `func sum(...rest) { return 0; }
sum(1, 2, 3, 4);
`
	program := parseForTest(t, src)
	a := NewAnalyzer()
	err := a.Analyze(program)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSemaCollectsMultipleErrorsInFunctionBody(t *testing.T) {
	src := `func main() {
		aaa;
		bbb;
		ccc;
	}`
	program := parseForTest(t, src)
	a := NewAnalyzer()
	err := a.Analyze(program)
	if err == nil {
		t.Fatal("expected sema errors")
	}
	if len(a.Errors()) < 3 {
		t.Fatalf("want at least 3 recorded errors, got %d: %v", len(a.Errors()), a.Errors())
	}
	var me *diagnostic.MultiError
	if !errors.As(err, &me) {
		t.Fatalf("want MultiError for multiple issues, got %T: %v", err, err)
	}
	if len(me.List) < 3 {
		t.Fatalf("MultiError list: want len >= 3, got %d", len(me.List))
	}
}

func TestPrepareNativeBundleSemaFileFromIncludedModule(t *testing.T) {
	tmp := t.TempDir()
	mainPath := filepath.Join(tmp, "main.koda")
	libPath := filepath.Join(tmp, "lib.koda")
	mainAbs, err := filepath.Abs(mainPath)
	if err != nil {
		t.Fatal(err)
	}
	libAbs, err := filepath.Abs(libPath)
	if err != nil {
		t.Fatal(err)
	}
	overlays := map[string]string{
		mainAbs: "#include \"lib.koda\"\n",
		libAbs:  "typoName + 1;\n",
	}
	bundle, err := parser.LoadProgramWithOverlays(mainPath, overlays)
	if err != nil {
		t.Fatal(err)
	}
	_, err = PrepareNativeBundle(bundle)
	if err == nil {
		t.Fatal("expected undefined typoName")
	}
	var de *diagnostic.DiagnosticError
	if !errors.As(err, &de) {
		t.Fatalf("want DiagnosticError, got %v", err)
	}
	got := filepath.Clean(de.File)
	want := filepath.Clean(libAbs)
	if !strings.EqualFold(got, want) {
		t.Fatalf("error file: got %q want %q (lib)", got, want)
	}
}

func TestSemaForOfBindsLoopVarInBody(t *testing.T) {
	src := `
		let lo = 0;
		let hi = 2;
		let sum = 0;
		for (let i of lo..hi) {
			sum = sum + i;
		}
	`
	program := parseForTest(t, src)
	if err := NewAnalyzer().Analyze(program); err != nil {
		t.Fatalf("sema: %v", err)
	}
}

func TestSemaDeferStmtAnalyzesCallee(t *testing.T) {
	src := `func main() {
		defer print(1);
	}`
	program := parseForTest(t, src)
	if err := NewAnalyzer().Analyze(program); err != nil {
		t.Fatalf("sema: %v", err)
	}
}

func TestSemaArgvMethodAritySplitTooFew(t *testing.T) {
	src := `func main() {
		"a".split();
	}`
	program := parseForTest(t, src)
	a := NewAnalyzer()
	if err := a.Analyze(program); err == nil {
		t.Fatal("expected arity error for split with 0 arguments")
	}
}

func TestSemaArgvMethodArityTrimTooMany(t *testing.T) {
	src := `func main() {
		"a".trim(1);
	}`
	program := parseForTest(t, src)
	a := NewAnalyzer()
	if err := a.Analyze(program); err == nil {
		t.Fatal("expected arity error for trim with 1 argument")
	}
}

func TestSemaArgvMethodArityTrimOk(t *testing.T) {
	src := `func main() {
		"a".trim();
	}`
	program := parseForTest(t, src)
	if err := NewAnalyzer().Analyze(program); err != nil {
		t.Fatalf("sema: %v", err)
	}
}

func TestSemaArgvMethodArityClearWithArgSkipped(t *testing.T) {
	src := `func main() {
		let game = { clear: func(c) { } };
		game.clear(1);
	}`
	program := parseForTest(t, src)
	if err := NewAnalyzer().Analyze(program); err != nil {
		t.Fatalf("game.clear(color) must not use argv .clear() arity: %v", err)
	}
}

func TestSemaArgvMethodArityComputedIndexSkipped(t *testing.T) {
	src := `func main() {
		let m = "split";
		"x"[m]();
	}`
	program := parseForTest(t, src)
	if err := NewAnalyzer().Analyze(program); err != nil {
		t.Fatalf("dynamic method name must not trigger argv arity: %v", err)
	}
}

func TestSemaParamCallableNotCheckedAsEnclosingFunc(t *testing.T) {
	src := `func pool(size, makeItem) {
		makeItem();
	}
	func main() {
		pool(1, func() { return 0; });
	}`
	program := parseForTest(t, src)
	if err := NewAnalyzer().Analyze(program); err != nil {
		t.Fatalf("calling a func param must not use enclosing func arity: %v", err)
	}
}

func TestSemaStructParamFieldAccess(t *testing.T) {
	src := `struct Point {
		x, y
	}
	func move(p) {
		p.x = p.x + 1;
		p.y = p.y + 2;
	}
	func main() {
		let p = Point { x: 0, y: 0 };
		move(p);
	}`
	program := parseForTest(t, src)
	a := NewAnalyzer()
	if err := a.Analyze(program); err != nil {
		t.Fatalf("sema: %v", err)
	}
	_, _, _, _, indexStruct, _, _, _, _ := a.ExportForCodegen()
	if len(indexStruct) < 2 {
		t.Fatalf("expected struct field slots for param p.x/p.y, got %d", len(indexStruct))
	}
}

func TestSemaAmbiguousStructParamTypes(t *testing.T) {
	src := `struct Vector2 { x, y }
	func dot(a, b) {
		return a.x * b.x + a.y * b.y;
	}
	func main() {
		let v = Vector2 { x: 1, y: 2 };
		dot(v, v);
		let plain = { x: 3, y: 4 };
		dot(plain, plain);
	}`
	program := parseForTest(t, src)
	a := NewAnalyzer()
	if err := a.Analyze(program); err != nil {
		t.Fatalf("sema: %v", err)
	}
	params := a.structParamsForFunc("dot")
	if params != nil {
		if _, ok := params["a"]; ok {
			t.Fatal("expected ambiguous dot param 'a' to fall back to dynamic field access")
		}
	}
	if a.funcParamPlain["dot"] == nil || !a.funcParamPlain["dot"]["a"] {
		t.Fatal("expected plain call site to mark dot param 'a'")
	}
}

func TestSemaNestedStructParamPropagation(t *testing.T) {
	src := `struct Vec3 { x, y, z }
	func dot(a, b) {
		return a.x * b.x + a.y * b.y + a.z * b.z;
	}
	func lengthsq(v) {
		return dot(v, v);
	}
	func length(v) {
		return lengthsq(v);
	}
	func main() {
		let v = Vec3 { x: 3, y: 0, z: 4 };
		length(v);
	}`
	program := parseForTest(t, src)
	a := NewAnalyzer()
	if err := a.Analyze(program); err != nil {
		t.Fatalf("sema: %v", err)
	}
	_, _, _, _, indexStruct, _, _, _, _ := a.ExportForCodegen()
	if len(indexStruct) < 3 {
		t.Fatalf("expected struct field slots inside dot() from nested call, got %d", len(indexStruct))
	}
}

func TestSemaStructParamMutateInHelper(t *testing.T) {
	src := `struct Ball { x, vx }
	func serve(b) {
		b.x = 400;
		b.vx = 360;
	}
	func main() {
		let ball = Ball { x: 0, vx: 0 };
		serve(ball);
	}`
	program := parseForTest(t, src)
	a := NewAnalyzer()
	if err := a.Analyze(program); err != nil {
		t.Fatalf("sema: %v", err)
	}
	_, _, _, _, indexStruct, _, _, _, _ := a.ExportForCodegen()
	if len(indexStruct) < 2 {
		t.Fatalf("expected struct field slots inside serve(), got %d", len(indexStruct))
	}
}

func TestSemaForwardReferences(t *testing.T) {
	src := `
		func early() {
			late_helper();
			let x = player.x;
		}
		struct Mario { x = 0.0; }
		struct Coin {
			on = true;
			func show() {
				if (on) { dc(0.0, 0.0, 0.0, 1.0, 1.0, 1.0, col_gold); }
			}
		}
		let player = Mario { x: 0.0 };
		let col_gold = { r: 255, g: 215, b: 0, a: 255 };
		func late_helper() { return; }
		func dc(a, b, c, d, e, f, g) { return; }
	`
	program := parseForTest(t, src)
	if err := NewAnalyzer().Analyze(program); err != nil {
		t.Fatalf("expected forward refs to resolve: %v", err)
	}
}
