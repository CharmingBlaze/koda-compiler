package formatter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormatIdempotent(t *testing.T) {
	src := "func add(a, b) {\n    return a + b;\n}\n\nlet x = add(1, 2);\nif (x > 0) {\n    print(x);\n}\n"
	out, err := Format(src)
	if err != nil {
		t.Fatal(err)
	}
	out2, err := Format(out)
	if err != nil {
		t.Fatal(err)
	}
	if out != out2 {
		t.Fatalf("not idempotent:\n---first---\n%s\n---second---\n%s", out, out2)
	}
}

func TestFormatCanonicalSpacing(t *testing.T) {
	src := "func   add(a,b){\nreturn a+b;\n}\nlet x=add(1,2);\nif(x>0){print(x);}\n"
	out, err := Format(src)
	if err != nil {
		t.Fatal(err)
	}
	want := "func add(a, b) {\n    return a + b;\n}\n\nlet x = add(1, 2);\n\nif (x > 0) {\n    print(x);\n}\n"
	if out != want {
		t.Fatalf("got:\n%q\nwant:\n%q", out, want)
	}
}

func TestFormatRangeExpr(t *testing.T) {
	src := "for(let i of lo..hi){print(i);}\n"
	out, err := Format(src)
	if err != nil {
		t.Fatal(err)
	}
	want := "for (let i of lo..hi) {\n    print(i);\n}\n"
	if out != want {
		t.Fatalf("got:\n%q\nwant:\n%q", out, want)
	}
}

func TestFormatHelloKoda(t *testing.T) {
	path := filepath.Join("..", "..", "tests", "hello.koda")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Skip(err)
	}
	text := strings.ReplaceAll(string(b), "\r\n", "\n")
	out, err := Format(text)
	if err != nil {
		t.Fatal(err)
	}
	out2, err := Format(out)
	if err != nil {
		t.Fatal(err)
	}
	if out != out2 {
		t.Fatal("hello.koda format not idempotent")
	}
}

func TestFormatConsecutiveExprStmtsNoBlankLine(t *testing.T) {
	src := "print(1);\nprint(2);\nprint(3);\n"
	out, err := Format(src)
	if err != nil {
		t.Fatal(err)
	}
	want := "print(1);\nprint(2);\nprint(3);\n"
	if out != want {
		t.Fatalf("got %q want %q", out, want)
	}
}

func TestFormatExprThenLetHasBlankLine(t *testing.T) {
	src := "print(1);\nlet x = 2;\n"
	out, err := Format(src)
	if err != nil {
		t.Fatal(err)
	}
	want := "print(1);\n\nlet x = 2;\n"
	if out != want {
		t.Fatalf("got %q want %q", out, want)
	}
}

func TestFormatIfElseBraceSameLineAsElse(t *testing.T) {
	src := "if (a) {\n    print(1);\n} else {\n    print(2);\n}\n"
	out, err := Format(src)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "}\n else") {
		t.Fatalf("newline between } and else:\n%s", out)
	}
	if !strings.Contains(out, "} else {") {
		t.Fatalf("expected '} else {', got:\n%s", out)
	}
}

func TestFormatImportExpr(t *testing.T) {
	src := "let io = import \"@io\";\n"
	out, err := Format(src)
	if err != nil {
		t.Fatal(err)
	}
	want := "let io = import \"@io\";\n"
	if out != want {
		t.Fatalf("got %q want %q", out, want)
	}
}

func TestFormatUnaryPlusIdempotent(t *testing.T) {
	src := "let x=+1+(+2*3);\n"
	out, err := Format(src)
	if err != nil {
		t.Fatal(err)
	}
	out2, err := Format(out)
	if err != nil {
		t.Fatal(err)
	}
	if out != out2 {
		t.Fatalf("not idempotent:\n---first---\n%s\n---second---\n%s", out, out2)
	}
}

func TestFormatClosureTest(t *testing.T) {
	path := filepath.Join("..", "..", "tests", "closure_test.koda")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Skip(err)
	}
	text := strings.ReplaceAll(string(b), "\r\n", "\n")
	out, err := Format(text)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Format(out); err != nil {
		t.Fatal(err)
	}
}
