package sema

import (
	"strings"
	"testing"

	"koda/internal/parser"
)

func TestInferNumericKindsLiterals(t *testing.T) {
	src := `func main() {
		let a = 1;
		let b = 3.0;
		let c = 3.14;
	}`
	prog := parseForTest(t, src)
	kinds := InferNumericKinds(prog, nil)
	var aDecl, bDecl, cDecl *parser.LetDecl
	for _, d := range prog.Declarations {
		if fd, ok := d.(*parser.FuncDecl); ok {
			for _, inner := range fd.Body.Declarations {
				if ld, ok := inner.(*parser.LetDecl); ok {
					switch ld.Name.Lexeme {
					case "a":
						aDecl = ld
					case "b":
						bDecl = ld
					case "c":
						cDecl = ld
					}
				}
			}
		}
	}
	if kinds[aDecl] != KindInt {
		t.Fatalf("a want KindInt got %v", kinds[aDecl])
	}
	if kinds[bDecl] != KindInt {
		t.Fatalf("b want KindInt got %v", kinds[bDecl])
	}
	if kinds[cDecl] != KindFloat {
		t.Fatalf("c want KindFloat got %v", kinds[cDecl])
	}
}

func TestInferNumericKindsDivision(t *testing.T) {
	src := `func main() {
		let x = 10;
		let y = 2;
		let z = x / y;
	}`
	prog := parseForTest(t, src)
	kinds := InferNumericKinds(prog, nil)
	var zDecl *parser.LetDecl
	for _, d := range prog.Declarations {
		if fd, ok := d.(*parser.FuncDecl); ok {
			for _, inner := range fd.Body.Declarations {
				if ld, ok := inner.(*parser.LetDecl); ok && ld.Name.Lexeme == "z" {
					zDecl = ld
				}
			}
		}
	}
	if kinds[zDecl] != KindFloat {
		t.Fatalf("z want KindFloat got %v", kinds[zDecl])
	}
}

func TestSemaWarnUnusedLetAndFunc(t *testing.T) {
	src := `let unusedVar = 1;
func unusedFn() { return 0; }
func main() { print("ok"); }
`
	prog := parseForTest(t, src)
	a := NewAnalyzerWithOptions(&AnalysisOptions{WarnUnused: true})
	if err := a.Analyze(prog); err != nil {
		t.Fatalf("analyze: %v", err)
	}
	w := a.Warnings()
	if len(w) < 2 {
		t.Fatalf("want >=2 warnings, got %v", w)
	}
	joined := strings.Join(w, "\n")
	if !strings.Contains(joined, "unused variable 'unusedvar'") {
		t.Fatalf("missing unused let warning: %v", w)
	}
	if !strings.Contains(joined, "unused function 'unusedfn'") {
		t.Fatalf("missing unused func warning: %v", w)
	}
}

func TestSemaEnumSwitchExhaustiveWarning(t *testing.T) {
	src := `enum Phase { Menu, Playing, GameOver }
func main() {
	let p = Phase.Menu;
	switch (p) {
		case Phase.Menu: print("menu");
	}
}
`
	prog := parseForTest(t, src)
	a := NewAnalyzer()
	if err := a.Analyze(prog); err != nil {
		t.Fatalf("analyze: %v", err)
	}
	w := a.Warnings()
	if len(w) != 1 {
		t.Fatalf("want 1 warning, got %v", w)
	}
	if !strings.Contains(w[0], "not exhaustive") || !strings.Contains(w[0], "playing") {
		t.Fatalf("unexpected warning: %q", w[0])
	}
}

func TestSemaEnumSwitchDefaultNoWarning(t *testing.T) {
	src := `enum Phase { Menu, Playing }
func main() {
	let phase = Phase.Menu;
	switch (phase) {
		case Phase.Menu: print("menu");
		default: print("other");
	}
}
`
	prog := parseForTest(t, src)
	a := NewAnalyzer()
	if err := a.Analyze(prog); err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if len(a.Warnings()) != 0 {
		t.Fatalf("want no warnings, got %v", a.Warnings())
	}
}

func TestSemaIntegerTypeAnnot(t *testing.T) {
	src := `func main() {
		let n: i32 = 42;
		let b: u8 = 255;
	}`
	prog := parseForTest(t, src)
	if err := NewAnalyzer().Analyze(prog); err != nil {
		t.Fatalf("analyze: %v", err)
	}
}

func TestSemaConstReassignment(t *testing.T) {
	src := `const gravity = 900;
func main() {
    gravity = 1;
}`
	prog := parseForTest(t, src)
	if err := NewAnalyzer().Analyze(prog); err == nil {
		t.Fatal("expected const reassignment error")
	}
}

func TestSemaConstRequiresInit(t *testing.T) {
	src := `const x;`
	prog := parseForTest(t, src)
	if err := NewAnalyzer().Analyze(prog); err == nil {
		t.Fatal("expected missing initializer error")
	}
}

func TestSemaStrictEqualityDeprecated(t *testing.T) {
	src := `func main() { let a = 1 === 1; }`
	prog := parseForTest(t, src)
	a := NewAnalyzer()
	if err := a.Analyze(prog); err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if len(a.Warnings()) == 0 {
		t.Fatal("expected === deprecation warning")
	}
}

func TestSemaBeginnerTypeAnnot(t *testing.T) {
	src := `func main() {
		let lives: int = 3;
		let name: string = "Jesse";
	}`
	prog := parseForTest(t, src)
	if err := NewAnalyzer().Analyze(prog); err != nil {
		t.Fatalf("analyze: %v", err)
	}
}

func TestSemaUnknownTypeAnnot(t *testing.T) {
	src := `func main() { let n: float32 = 1; }`
	prog := parseForTest(t, src)
	if err := NewAnalyzer().Analyze(prog); err == nil {
		t.Fatal("expected unknown type error")
	}
}
