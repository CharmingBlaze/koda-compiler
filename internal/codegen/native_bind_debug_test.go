package codegen

import (
	"testing"

	"koda/internal/parser"
)

func TestNativeLetParsedForStructMethodTest(t *testing.T) {
	src := `
struct Widget {
	func ping() {
		shimfn(1);
	}
}

// koda:extern shimfn koda_test_shimfn 1
let shimfn = 0;
`
	program := parseForTest(t, src)
	for _, d := range program.Declarations {
		ld, ok := d.(*parser.LetDecl)
		if !ok || ld.Name.Lexeme != "shimfn" {
			continue
		}
		if ld.Native == nil {
			t.Fatal("shimfn let missing Native directive")
		}
		t.Logf("native symbol=%s", ld.Native.Symbol)
		return
	}
	t.Fatal("shimfn let not found")
}
