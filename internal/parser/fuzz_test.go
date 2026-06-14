package parser

import (
	"koda/internal/lexer"
	"testing"
)

func FuzzParse(f *testing.F) {
	f.Add([]byte("let x = 10;"))
	f.Add([]byte("func f() { return; }"))
	f.Add([]byte("struct s { a, b }\nlet o = s { a: 1, b: 2 };"))

	f.Fuzz(func(t *testing.T, data []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("parser panicked on input %q: %v", data, r)
			}
		}()

		l := lexer.NewLexer(string(data), "<fuzz>")
		toks, err := l.Tokenize()
		if err != nil {
			return
		}
		p := NewParser(toks)
		_, _ = p.Parse()
	})
}
