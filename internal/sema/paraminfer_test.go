package sema

import (
	"testing"

	"koda/internal/lexer"
	"koda/internal/parser"
)

func parseParamInfer(t *testing.T, src string) *parser.Program {
	t.Helper()
	l := lexer.NewLexer(src, "test.koda")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	prog, err := parser.NewParser(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	return prog
}

func TestInferParamNumericKindsDelta(t *testing.T) {
	src := `
		func updatePlaying(dt) {
			paddlex += 480.0 * dt;
		}
	`
	prog := parseParamInfer(t, src)
	kinds := InferParamNumericKinds(prog)
	fn := prog.Declarations[0].(*parser.FuncDecl)
	key := NewParamCellKey(fn, 0)
	if kinds[key] != KindFloat {
		t.Fatalf("dt want KindFloat got %v", kinds[key])
	}
}

func TestInferParamNumericKindsExplicitAnnotSkipped(t *testing.T) {
	src := `
		func scale(dt: float) {
			return dt;
		}
	`
	prog := parseParamInfer(t, src)
	kinds := InferParamNumericKinds(prog)
	fn := prog.Declarations[0].(*parser.FuncDecl)
	key := NewParamCellKey(fn, 0)
	if _, ok := kinds[key]; ok {
		t.Fatalf("explicitly typed param should not appear in ParamNumericKinds")
	}
}

func TestInferParamNumericKindsIntOnly(t *testing.T) {
	src := `
		func repeat(count) {
			let i = 0;
			while i < count {
				i = i + 1;
			}
		}
	`
	prog := parseParamInfer(t, src)
	kinds := InferParamNumericKinds(prog)
	fn := prog.Declarations[0].(*parser.FuncDecl)
	key := NewParamCellKey(fn, 0)
	if kinds[key] != KindInt {
		t.Fatalf("count want KindInt got %v", kinds[key])
	}
}

func TestInferParamKindsFromCallSite(t *testing.T) {
	src := `
		func drawFrame(dt) {
			print("frame");
		}
		func main() {
			drawFrame(game.delta());
		}
	`
	prog := parseParamInfer(t, src)
	letKinds := InferNumericKinds(prog, nil)
	kinds := InferParamNumericKinds(prog)
	InferParamKindsFromCallSites(prog, letKinds, kinds)
	fn := prog.Declarations[0].(*parser.FuncDecl)
	key := NewParamCellKey(fn, 0)
	if kinds[key] != KindFloat {
		t.Fatalf("dt from game.delta() call want KindFloat got %v", kinds[key])
	}
}

func TestInferParamKindsFromStructMethodCallSite(t *testing.T) {
	src := `
		struct Hero {
			x: float = 0.0;
			func tick(dt) {
				print("ok");
			}
		}
		func main() {
			let h = Hero {};
			h.tick(0.5);
		}
	`
	prog := parseParamInfer(t, src)
	letKinds := InferNumericKinds(prog, nil)
	kinds := InferParamNumericKinds(prog)
	InferParamKindsFromCallSites(prog, letKinds, kinds)
	hero := prog.Declarations[0].(*parser.StructDecl)
	tick := hero.Methods[0]
	key := NewParamCellKey(tick, 0)
	if kinds[key] != KindFloat {
		t.Fatalf("struct method dt from float literal call want KindFloat got %v", kinds[key])
	}
}

func TestInferFuncExprParamNumericKinds(t *testing.T) {
	src := `
		func main() {
			let tickFn = func(t) { return t * 2.0; };
			print(tickFn(1.0));
		}
	`
	prog := parseParamInfer(t, src)
	kinds := InferParamNumericKinds(prog)
	mainFn := prog.Declarations[0].(*parser.FuncDecl)
	letStep := mainFn.Body.Declarations[0].(*parser.LetDecl)
	fe := letStep.Init.(*parser.FuncExpr)
	key := NewParamCellKey(fe, 0)
	if kinds[key] != KindFloat {
		t.Fatalf("lambda param t want KindFloat got %v", kinds[key])
	}
}
