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
