package codegen

import (
	"fmt"
	"strconv"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"

	"koda/internal/parser"
)

func (g *Generator) emitTestDecl(t *parser.TestDecl) error {
	fd := t.SyntheticFunc()
	if err := g.emitFuncDecl(fd); err != nil {
		return err
	}
	fn := g.funcs[fd.Name.Lexeme]
	if fn == nil {
		return fmt.Errorf("codegen: missing LLVM function for test %q", t.Display.Lexeme)
	}
	g.testFuncs = append(g.testFuncs, testEntry{display: unquoteTestName(t.Display.Lexeme), fn: fn})
	return nil
}

func (g *Generator) emitTestRunner() {
	zero := constant.NewInt(types.I64, 0)
	for _, te := range g.testFuncs {
		g.block.NewCall(te.fn, zero)
		msg := g.emitStringLiteral("PASS: " + te.display)
		g.block.NewCall(g.runtimePrint, g.emitAsKodaI64(msg))
		g.block.NewCall(g.runtimePrintNewline)
	}
}

func unquoteTestName(lexeme string) string {
	if u, err := strconv.Unquote(lexeme); err == nil {
		return u
	}
	return lexeme
}
