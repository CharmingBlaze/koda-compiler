package lexer

import (
	"strings"
	"testing"
)

func assertTokenTypes(t *testing.T, source string, expected []TokenType) {
	t.Helper()
	l := NewLexer(source, "")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}
	if len(tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d: %#v", len(expected), len(tokens), tokens)
	}
	for i, tok := range tokens {
		if tok.Type != expected[i] {
			t.Errorf("At index %d: expected %v, got %v (%s)", i, expected[i], tok.Type, tok.Lexeme)
		}
	}
}

func TestLexerStringUnterminatedEscapeEOF(t *testing.T) {
	l := NewLexer("\"\\", "")
	_, err := l.Tokenize()
	if err == nil {
		t.Fatal(`want error for "\"" + "\\" at EOF`)
	}
	if !strings.Contains(err.Error(), "unterminated") {
		t.Fatalf("want unterminated string error, got %v", err)
	}
}

func TestLexerLoneOpeningBacktickUnterminatedTemplate(t *testing.T) {
	l := NewLexer("`", "")
	_, err := l.Tokenize()
	if err == nil {
		t.Fatal("want error for lone opening backtick / unterminated template")
	}
	if !strings.Contains(err.Error(), "unterminated") {
		t.Fatalf("want unterminated template error, got %v", err)
	}
}

func TestLexerTemplateTokensCarryFile(t *testing.T) {
	const path = "/tmp/tpl.koda"
	l := NewLexer("`a${x}b`", path)
	toks, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	for _, tok := range toks {
		if tok.Type == TokenEOF {
			continue
		}
		if tok.File != path {
			t.Fatalf("token %s missing file: got %q", tok.Type, tok.File)
		}
	}
}

func TestLexerSetFilePropagates(t *testing.T) {
	l := NewLexer("let x = 1;", "")
	l.SetFile("/tmp/example.koda")
	toks, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	if len(toks) < 2 || toks[0].File != "/tmp/example.koda" {
		t.Fatalf("expected File on tokens, got %#v", toks[0])
	}
}

func TestLexer(t *testing.T) {
	source := `
		let x = 10;
		if (x > 5) {
			print("Hello");
		}
	`
	expected := []TokenType{
		TokenLet, TokenIdentifier, TokenEqual, TokenNumber, TokenSemicolon,
		TokenIf, TokenLParen, TokenIdentifier, TokenGreater, TokenNumber, TokenRParen, TokenLBrace,
		TokenIdentifier, TokenLParen, TokenString, TokenRParen, TokenSemicolon,
		TokenRBrace, TokenEOF,
	}
	assertTokenTypes(t, source, expected)
}

func TestLexerCaseInsensitiveKeywords(t *testing.T) {
	l := NewLexer("LET X = 1;\nFUNC F() { RETURN X; }\n", "")
	toks, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	if toks[0].Type != TokenLet || toks[0].Lexeme != "let" {
		t.Fatalf("LET: got %v %q", toks[0].Type, toks[0].Lexeme)
	}
	if toks[1].Type != TokenIdentifier || toks[1].Lexeme != "X" {
		t.Fatalf("identifier lexeme: got %q", toks[1].Lexeme)
	}
	if toks[5].Type != TokenFunc || toks[5].Lexeme != "func" {
		t.Fatalf("FUNC: got %v %q", toks[5].Type, toks[5].Lexeme)
	}
	if toks[10].Type != TokenReturn || toks[10].Lexeme != "return" {
		t.Fatalf("RETURN: got %v %q", toks[10].Type, toks[10].Lexeme)
	}
}

func TestLexerIncludeDirectiveCase(t *testing.T) {
	l := NewLexer("#INCLUDE \"a.koda\"\n", "")
	toks, err := l.Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	if toks[0].Type != TokenInclude || strings.ToLower(toks[0].Lexeme) != "#include" {
		t.Fatalf("got %#v", toks[0])
	}
}

func TestLexerComplex(t *testing.T) {
	source := `func add(a, b) { return a + b; }`
	expected := []TokenType{
		TokenFunc, TokenIdentifier, TokenLParen, TokenIdentifier, TokenComma, TokenIdentifier, TokenRParen, TokenLBrace,
		TokenReturn, TokenIdentifier, TokenPlus, TokenIdentifier, TokenSemicolon,
		TokenRBrace, TokenEOF,
	}
	assertTokenTypes(t, source, expected)
}

func TestLexerDeferKeyword(t *testing.T) {
	assertTokenTypes(t, `defer cleanup();`, []TokenType{
		TokenDefer, TokenIdentifier, TokenLParen, TokenRParen, TokenSemicolon,
		TokenEOF,
	})
}

func TestLexerNullishAssign(t *testing.T) {
	assertTokenTypes(t, `a ??= b; x ?? y;`, []TokenType{
		TokenIdentifier, TokenQuestionQuestionEqual, TokenIdentifier, TokenSemicolon,
		TokenIdentifier, TokenQuestionQuestion, TokenIdentifier, TokenSemicolon,
		TokenEOF,
	})
}

func TestLexerRejectsStrictEquality(t *testing.T) {
	l := NewLexer(`a === b`, "test.koda")
	if _, err := l.Tokenize(); err == nil {
		t.Fatal("expected error for ===")
	}
	l = NewLexer(`a !== b`, "test.koda")
	if _, err := l.Tokenize(); err == nil {
		t.Fatal("expected error for !==")
	}
}

func TestLexerOperatorsAndDirectives(t *testing.T) {
	source := `#include "x.koda"
let x = a++ + --b;
x += 1; x -= 2; x *= 3; x /= 4; x %= 5;
x &= 1; x |= 2; x ^= 3; x <<= 1; x >>= 1;
let y = a && b || !c ? d : e;
let z = a..b ... rest => value;
let shifts = x << 2 >> 1;`
	expected := []TokenType{
		TokenInclude, TokenString,
		TokenLet, TokenIdentifier, TokenEqual, TokenIdentifier, TokenPlusPlus, TokenPlus, TokenMinusMinus, TokenIdentifier, TokenSemicolon,
		TokenIdentifier, TokenPlusEqual, TokenNumber, TokenSemicolon, TokenIdentifier, TokenMinusEqual, TokenNumber, TokenSemicolon, TokenIdentifier, TokenStarEqual, TokenNumber, TokenSemicolon, TokenIdentifier, TokenSlashEqual, TokenNumber, TokenSemicolon, TokenIdentifier, TokenPercentEqual, TokenNumber, TokenSemicolon,
		TokenIdentifier, TokenAndEqual, TokenNumber, TokenSemicolon, TokenIdentifier, TokenOrEqual, TokenNumber, TokenSemicolon, TokenIdentifier, TokenCaretEqual, TokenNumber, TokenSemicolon, TokenIdentifier, TokenLessLessEqual, TokenNumber, TokenSemicolon, TokenIdentifier, TokenGreaterGreaterEqual, TokenNumber, TokenSemicolon,
		TokenLet, TokenIdentifier, TokenEqual, TokenIdentifier, TokenAndAnd, TokenIdentifier, TokenOrOr, TokenBang, TokenIdentifier, TokenQuestion, TokenIdentifier, TokenColon, TokenIdentifier, TokenSemicolon,
		TokenLet, TokenIdentifier, TokenEqual, TokenIdentifier, TokenDotDot, TokenIdentifier, TokenTripleDot, TokenIdentifier, TokenArrow, TokenIdentifier, TokenSemicolon,
		TokenLet, TokenIdentifier, TokenEqual, TokenIdentifier, TokenLessLess, TokenNumber, TokenGreaterGreater, TokenNumber, TokenSemicolon,
		TokenEOF,
	}
	assertTokenTypes(t, source, expected)
}

func TestLexerCommentsBOMAndPositions(t *testing.T) {
	source := "\ufefflet first = 1; // line comment\n/* block\ncomment */\nlet second = 2;"
	l := NewLexer(source, "")
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize failed: %v", err)
	}
	expected := []TokenType{
		TokenLet, TokenIdentifier, TokenEqual, TokenNumber, TokenSemicolon,
		TokenLet, TokenIdentifier, TokenEqual, TokenNumber, TokenSemicolon,
		TokenEOF,
	}
	if len(tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d: %#v", len(expected), len(tokens), tokens)
	}
	for i, typ := range expected {
		if tokens[i].Type != typ {
			t.Fatalf("At index %d: expected %v, got %v", i, typ, tokens[i].Type)
		}
	}
	if tokens[0].Line != 1 || tokens[0].Col != 1 {
		t.Fatalf("expected first token at 1:1, got %d:%d", tokens[0].Line, tokens[0].Col)
	}
	if tokens[5].Line != 4 || tokens[5].Col != 1 {
		t.Fatalf("expected second let at 4:1, got %d:%d", tokens[5].Line, tokens[5].Col)
	}
}

func TestVarIsReservedKeyword(t *testing.T) {
	tokens, err := NewLexer("var x = 1;", "").Tokenize()
	if err != nil {
		t.Fatal(err)
	}
	if len(tokens) < 2 || tokens[0].Type != TokenVar || tokens[0].Lexeme != "var" {
		t.Fatalf("expected leading TokenVar, got %#v", tokens)
	}
}

func TestLexerReportsErrors(t *testing.T) {
	if _, err := NewLexer("let x = @;", "").Tokenize(); err == nil {
		t.Fatal("expected invalid character error")
	}
	if _, err := NewLexer(`"unterminated`, "").Tokenize(); err == nil {
		t.Fatal("expected unterminated string error")
	}
}

func BenchmarkLexerTokenize(b *testing.B) {
	source := `
func fib(n) {
    if (n <= 1) { return n; }
    return fib(n - 1) + fib(n - 2);
}
let total = 0;
for (let i = 0; i < 100; i++) {
    total += fib(10);
}
`
	for i := 0; i < b.N; i++ {
		if _, err := NewLexer(source, "").Tokenize(); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkLexerTokenizeLarge measures throughput on a multi-KB slice (stable MB/s).
func BenchmarkLexerTokenizeLarge(b *testing.B) {
	chunk := `func f(n) { if (n <= 0) { return 0; } return n + f(n - 1); }
let x = 1 + 2 * 3; let s = "hello\"world"; // comment
/* block */ for (let i = 0; i < 10; i++) { x += i; }
`
	var sb strings.Builder
	for sb.Len() < 64*1024 {
		sb.WriteString(chunk)
	}
	source := sb.String()
	b.SetBytes(int64(len(source)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := NewLexer(source, "").Tokenize(); err != nil {
			b.Fatal(err)
		}
	}
}
