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

func TestSemaArgvMethodArityReduceZeroArgs(t *testing.T) {
	src := `func main() {
		let a = [1];
		a.reduce();
	}`
	program := parseForTest(t, src)
	a := NewAnalyzer()
	if err := a.Analyze(program); err == nil {
		t.Fatal("expected arity error for reduce with 0 arguments")
	}
}
