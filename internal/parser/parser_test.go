package parser

import (
	"path/filepath"
	"strings"
	"testing"

	"koda/internal/lexer"
)

func TestParseClassicForLoop(t *testing.T) {
	src := `func main() {
		let total = 0;
		for (let i = 0; i < 5; i += 1) {
			if (i == 3) { continue; }
			total += i;
		}
	}`
	l := lexer.NewLexer(src, "")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	prog, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	fn := prog.Declarations[0].(*FuncDecl)
	forStmt := fn.Body.Declarations[1].(*ForStmt)
	if len(forStmt.Inits) != 1 {
		t.Fatalf("inits: want 1, got %d", len(forStmt.Inits))
	}
	if forStmt.Condition == nil {
		t.Fatal("expected condition")
	}
	if len(forStmt.Increments) != 1 {
		t.Fatalf("increments: want 1, got %d", len(forStmt.Increments))
	}
}

func TestParseForOfDestructuring(t *testing.T) {
	src := `func main() {
		let o = { a: 1 };
		for (let [k, v] of o) { print(k); }
	}`
	l := lexer.NewLexer(src, "")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	prog, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	fn := prog.Declarations[0].(*FuncDecl)
	forOf := fn.Body.Declarations[1].(*ForOfStmt)
	if forOf.ValueVar == nil {
		t.Fatal("expected ValueVar for [k, v] binding")
	}
	if forOf.VarName.Lexeme != "k" || forOf.ValueVar.Lexeme != "v" {
		t.Fatalf("bindings: got %q, %q", forOf.VarName.Lexeme, forOf.ValueVar.Lexeme)
	}
}

func TestParseClassicForInfinite(t *testing.T) {
	src := `func main() { for (;;) { break; } }`
	l := lexer.NewLexer(src, "")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	prog, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	fn := prog.Declarations[0].(*FuncDecl)
	forStmt := fn.Body.Declarations[0].(*ForStmt)
	if len(forStmt.Inits) != 0 || forStmt.Condition != nil || len(forStmt.Increments) != 0 {
		t.Fatalf("for (;;): got inits=%d cond=%v inc=%d", len(forStmt.Inits), forStmt.Condition != nil, len(forStmt.Increments))
	}
}

func TestVarReservedUseLet(t *testing.T) {
	l := lexer.NewLexer(`var x = 1;`, "")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	p := NewParser(tokens)
	_, err = p.Parse()
	if err == nil {
		t.Fatal("expected error for var declaration")
	}
	if !strings.Contains(err.Error(), "let") {
		t.Fatalf("expected hint to use let, got: %v", err)
	}
}

func TestVarReservedInExpression(t *testing.T) {
	l := lexer.NewLexer(`func main() { let a = var; }`, "")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	p := NewParser(tokens)
	_, err = p.Parse()
	if err == nil {
		t.Fatal("expected error for var in expression")
	}
	if !strings.Contains(err.Error(), "let") {
		t.Fatalf("expected hint to use let, got: %v", err)
	}
}

func TestParseRangeInForOf(t *testing.T) {
	src := `func main() {
		let lo = 1;
		let hi = 4;
		for (let i of lo..hi) { print(i); }
	}`
	prog := parseForTest(t, src)
	fn := prog.Declarations[0].(*FuncDecl)
	forOf := fn.Body.Declarations[2].(*ForOfStmt)
	r, ok := forOf.Iterable.(*RangeExpr)
	if !ok {
		t.Fatalf("expected RangeExpr iterable, got %T", forOf.Iterable)
	}
	if _, ok := r.From.(*IdentifierExpr); !ok {
		t.Fatalf("range From: want ident, got %T", r.From)
	}
	if _, ok := r.To.(*IdentifierExpr); !ok {
		t.Fatalf("range To: want ident, got %T", r.To)
	}
}

func TestParseDeferStmt(t *testing.T) {
	src := `func main() { defer print(1); }`
	l := lexer.NewLexer(src, "")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	prog, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	fn := prog.Declarations[0].(*FuncDecl)
	ds := fn.Body.Declarations[0].(*DeferStmt)
	if _, ok := ds.Expr.(*CallExpr); !ok {
		t.Fatalf("expected call in defer, got %T", ds.Expr)
	}
}

func TestParseNullishAssign(t *testing.T) {
	src := `func main() { let x = null; x ??= 1; }`
	l := lexer.NewLexer(src, "")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	prog, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	fn := prog.Declarations[0].(*FuncDecl)
	es := fn.Body.Declarations[1].(*ExpressionStmt)
	as, ok := es.Expr.(*AssignExpr)
	if !ok || as.Token.Type != lexer.TokenQuestionQuestionEqual {
		t.Fatalf("expected ??= assign, got %#v", es.Expr)
	}
}

func TestParserCaseFoldsBoundNames(t *testing.T) {
	src := `LET X = 1;
func G() { RETURN X; }`
	l := lexer.NewLexer(src, "")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	prog, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	let := prog.Declarations[0].(*LetDecl)
	if let.Name.Lexeme != "x" {
		t.Fatalf("let name: want x, got %q", let.Name.Lexeme)
	}
	fn := prog.Declarations[1].(*FuncDecl)
	if fn.Name.Lexeme != "g" {
		t.Fatalf("func name: want g, got %q", fn.Name.Lexeme)
	}
	ret := fn.Body.Declarations[0].(*ReturnStmt)
	id := ret.Value.(*IdentifierExpr)
	if id.Name.Lexeme != "x" {
		t.Fatalf("return ref: want x, got %q", id.Name.Lexeme)
	}
}

func TestParserRejectDuplicateParams(t *testing.T) {
	l := lexer.NewLexer(`func f(A, a) { return 1; }`, "")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	_, err = NewParser(tokens).Parse()
	if err == nil {
		t.Fatal("expected duplicate parameter error")
	}
	if !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("got: %v", err)
	}
}

func TestParser(t *testing.T) {
	source := `
		let x = 10;
		func add(a, b) {
			return a + b;
		}
		if (x > 5) {
			add(x, 2);
		}
	`
	l := lexer.NewLexer(source, "")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("Lexer failed: %v", err)
	}

	p := NewParser(tokens)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parser failed: %v", err)
	}

	if len(program.Declarations) != 3 {
		t.Errorf("Expected 3 declarations, got %d", len(program.Declarations))
	}

	// Verify one node type
	if _, ok := program.Declarations[0].(*LetDecl); !ok {
		t.Errorf("Expected first decl to be LetDecl, got %T", program.Declarations[0])
	}
}

func parseForTest(t *testing.T, source string) *Program {
	t.Helper()
	l := lexer.NewLexer(source, "")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("Lexer failed: %v", err)
	}
	p := NewParser(tokens)
	program, err := p.Parse()
	if err != nil {
		t.Fatalf("Parser failed: %v", err)
	}
	return program
}

func TestParserImportAssignmentAndBlockDeclarations(t *testing.T) {
	program := parseForTest(t, `
		let math = import "@math";
		let x = 1;
		x = x + 2;
		if (x > 1) {
			let y = x;
			y = y + 1;
		}
	`)
	if len(program.Declarations) != 4 {
		t.Fatalf("expected 4 declarations, got %d", len(program.Declarations))
	}
	first := program.Declarations[0].(*LetDecl)
	if _, ok := first.Init.(*ImportExpr); !ok {
		t.Fatalf("expected import initializer, got %T", first.Init)
	}
	assignStmt := program.Declarations[2].(*ExpressionStmt)
	if _, ok := assignStmt.Expr.(*AssignExpr); !ok {
		t.Fatalf("expected assignment expression, got %T", assignStmt.Expr)
	}
	ifStmt := program.Declarations[3].(*IfStmt)
	block := ifStmt.Then.(*BlockStmt)
	if len(block.Declarations) != 2 {
		t.Fatalf("expected 2 block statements, got %d", len(block.Declarations))
	}
	if _, ok := block.Declarations[0].(*LetDecl); !ok {
		t.Fatalf("expected let declaration inside block, got %T", block.Declarations[0])
	}
}

func TestParserDefaultAndRestParams(t *testing.T) {
	program := parseForTest(t, `
		func collect(a, b = 2, ...rest) {
			return a;
		}
	`)
	fn := program.Declarations[0].(*FuncDecl)
	if len(fn.Params) != 3 {
		t.Fatalf("expected 3 params, got %d", len(fn.Params))
	}
	if fn.Params[1].Default == nil {
		t.Fatal("expected default value for second param")
	}
	if !fn.Params[2].IsRest || fn.Params[2].Name != "rest" {
		t.Fatalf("expected rest param, got %#v", fn.Params[2])
	}
}

func TestParserBreakContinueAndInclude(t *testing.T) {
	program := parseForTest(t, `
		#include "helpers.koda"
		while (true) {
			continue;
			break;
		}
	`)
	if _, ok := program.Declarations[0].(*IncludeDecl); !ok {
		t.Fatalf("expected include declaration, got %T", program.Declarations[0])
	}
	loop := program.Declarations[1].(*WhileStmt)
	body := loop.Body.(*BlockStmt)
	if _, ok := body.Declarations[0].(*ContinueStmt); !ok {
		t.Fatalf("expected continue statement, got %T", body.Declarations[0])
	}
	if _, ok := body.Declarations[1].(*BreakStmt); !ok {
		t.Fatalf("expected break statement, got %T", body.Declarations[1])
	}
}

func TestParserRejectsInvalidRestParams(t *testing.T) {
	l := lexer.NewLexer(`func bad(...rest, x) { return x; }`, "")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("Lexer failed: %v", err)
	}
	_, err = NewParser(tokens).Parse()
	if err == nil {
		t.Fatal("expected parser error for non-final rest parameter")
	}
}

func TestProgramIncludeLoadsShim(t *testing.T) {
	repoExamples := filepath.Join("..", "..", "examples", "raylib_shim_demo.koda")
	bundle, err := LoadProgram(repoExamples)
	if err != nil {
		t.Fatal(err)
	}
	if len(bundle.Modules) < 2 {
		t.Fatalf("expected included module loaded, have %d modules", len(bundle.Modules))
	}
	if err := FlattenEntryIncludes(bundle); err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, d := range bundle.Entry.Declarations {
		if let, ok := d.(*LetDecl); ok {
			names = append(names, let.Name.Lexeme)
		}
	}
	found := false
	for _, n := range names {
		if n == "initwindow" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("initwindow not in entry after flatten; let names: %v", names)
	}
}

func TestParseNotPropertyAccess(t *testing.T) {
	src := `assert(!bad.ok, "m");`
	l := lexer.NewLexer(src, "")
	toks, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	p := NewParser(toks)
	prog, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	es := prog.Declarations[0].(*ExpressionStmt)
	call := es.Expr.(*CallExpr)
	pfx := call.Arguments[0].(*PrefixExpr)
	if pfx.Operator != "!" {
		t.Fatalf("want ! prefix, got %q", pfx.Operator)
	}
	if _, ok := pfx.Right.(*IndexExpr); !ok {
		t.Fatalf("want IndexExpr under !, got %T", pfx.Right)
	}
}

func TestParseCompoundAssignStatement(t *testing.T) {
	src := `let x = 10;
x += 3;
`
	l := lexer.NewLexer(src, "")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	p := NewParser(tokens)
	prog, err := p.Parse()
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if len(prog.Declarations) != 2 {
		t.Fatalf("expected 2 declarations, got %d", len(prog.Declarations))
	}
	es, ok := prog.Declarations[1].(*ExpressionStmt)
	if !ok {
		t.Fatalf("second decl: got %T", prog.Declarations[1])
	}
	ae, ok := es.Expr.(*AssignExpr)
	if !ok {
		t.Fatalf("expr: got %T", es.Expr)
	}
	if ae.Token.Type != lexer.TokenPlusEqual {
		t.Fatalf("assign token: got %v", ae.Token.Type)
	}
}

func TestParseDeferStatement(t *testing.T) {
	src := `func main() {
		defer print(1);
		print(0);
	}`
	l := lexer.NewLexer(src, "")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	prog, err := NewParser(tokens).Parse()
	if err != nil {
		t.Fatal(err)
	}
	fn := prog.Declarations[0].(*FuncDecl)
	if len(fn.Body.Declarations) != 2 {
		t.Fatalf("body decls: want 2, got %d", len(fn.Body.Declarations))
	}
	ds, ok := fn.Body.Declarations[0].(*DeferStmt)
	if !ok {
		t.Fatalf("first stmt: want *DeferStmt, got %T", fn.Body.Declarations[0])
	}
	if _, ok := ds.Expr.(*CallExpr); !ok {
		t.Fatalf("defer expr: want *CallExpr, got %T", ds.Expr)
	}
}
