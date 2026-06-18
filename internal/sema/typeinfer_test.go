package sema

import (
	"strings"
	"testing"

	"koda/internal/lexer"
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
	if kinds[bDecl] != KindFloat {
		t.Fatalf("b want KindFloat got %v", kinds[bDecl])
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

func TestInferNumericKindsFloatTypeAnnot(t *testing.T) {
	src := `func main() {
		let x: float = 1.0;
		let y: float = 2.0;
		let n: i32 = 42;
	}`
	prog := parseForTest(t, src)
	kinds := InferNumericKinds(prog, nil)
	var xDecl, yDecl, nDecl *parser.LetDecl
	for _, d := range prog.Declarations {
		if fd, ok := d.(*parser.FuncDecl); ok {
			for _, inner := range fd.Body.Declarations {
				if ld, ok := inner.(*parser.LetDecl); ok {
					switch ld.Name.Lexeme {
					case "x":
						xDecl = ld
					case "y":
						yDecl = ld
					case "n":
						nDecl = ld
					}
				}
			}
		}
	}
	if kinds[xDecl] != KindFloat {
		t.Fatalf("x want KindFloat got %v", kinds[xDecl])
	}
	if kinds[yDecl] != KindFloat {
		t.Fatalf("y want KindFloat got %v", kinds[yDecl])
	}
	if kinds[nDecl] != KindInt {
		t.Fatalf("n want KindInt got %v", kinds[nDecl])
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

func TestSemaStructFieldDefaultTypeMismatch(t *testing.T) {
	src := `struct Bad { score: string = true; }
func main() { let x = Bad {}; }`
	prog := parseForTest(t, src)
	err := NewAnalyzer().Analyze(prog)
	if err == nil {
		t.Fatal("expected type mismatch on struct field default")
	}
	if !strings.Contains(err.Error(), "expected type 'string'") {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestSemaSelfOutsideMethod(t *testing.T) {
	src := `func main() { let x = self; }`
	prog := parseForTest(t, src)
	err := NewAnalyzer().Analyze(prog)
	if err == nil {
		t.Fatal("expected self outside method error")
	}
	if !strings.Contains(err.Error(), "'self' can only be used inside struct methods") {
		t.Fatalf("error = %q", err.Error())
	}
}

func TestSemaConstReassignment(t *testing.T) {
	src := `const gravity = 900;
func main() {
    gravity = 1;
}`
	prog := parseForTest(t, src)
	err := NewAnalyzer().Analyze(prog)
	if err == nil {
		t.Fatal("expected const reassignment error")
	}
	if !strings.Contains(err.Error(), "Cannot assign to constant 'gravity'") {
		t.Fatalf("error = %q, want Cannot assign to constant", err.Error())
	}
}

func TestSemaConstRequiresInit(t *testing.T) {
	src := `const x;`
	prog := parseForTest(t, src)
	if err := NewAnalyzer().Analyze(prog); err == nil {
		t.Fatal("expected missing initializer error")
	}
}

func TestSemaRejectsStrictEquality(t *testing.T) {
	src := `func main() { let a = 1 === 1; }`
	l := lexer.NewLexer(src, "test.koda")
	if _, err := l.Tokenize(); err == nil {
		t.Fatal("expected lexer error for ===")
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

func TestInferNumericKindsGameDelta(t *testing.T) {
	src := `
		func main() {
			while true {
				let dt = game.delta();
				let x = dt * 2.0;
			}
		}
	`
	prog := parseForTest(t, src)
	kinds := InferNumericKinds(prog, nil)
	mainFn := prog.Declarations[0].(*parser.FuncDecl)
	var dtDecl, xDecl *parser.LetDecl
	var findLets func(parser.Decl)
	findLets = func(d parser.Decl) {
		switch x := d.(type) {
		case *parser.LetDecl:
			switch x.Name.Lexeme {
			case "dt":
				dtDecl = x
			case "x":
				xDecl = x
			}
		case *parser.BlockStmt:
			for _, inner := range x.Declarations {
				findLets(inner)
			}
		case parser.Stmt:
			if bs, ok := x.(*parser.WhileStmt); ok {
				findLets(bs.Body)
			}
		}
	}
	findLets(mainFn.Body)
	if dtDecl == nil || xDecl == nil {
		t.Fatal("missing dt or x let decl")
	}
	if kinds[dtDecl] != KindFloat {
		t.Fatalf("dt want KindFloat got %v", kinds[dtDecl])
	}
	if kinds[xDecl] != KindFloat {
		t.Fatalf("x want KindFloat got %v", kinds[xDecl])
	}
}

func TestInferNumericKindsSkipsArrayInit(t *testing.T) {
	src := `
		func main() {
			let a = [10, 20, 30];
			print(a[0]);
		}
	`
	prog := parseForTest(t, src)
	kinds := InferNumericKinds(prog, nil)
	for ld := range kinds {
		if ld.Name.Lexeme == "a" {
			t.Fatalf("array let a should not get numeric kind, got %v", kinds[ld])
		}
	}
}
