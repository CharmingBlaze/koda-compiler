package sema

import (
	"fmt"
	"strings"

	"koda/internal/parser"
)

func (a *Analyzer) noteBindingRead(decl parser.Decl) {
	switch d := decl.(type) {
	case *parser.LetDecl:
		if d.Native != nil {
			return
		}
		a.letReads[d]++
	case *parser.FuncDecl:
		if d.Native != nil {
			return
		}
		a.funcReads[d]++
	}
}

func (a *Analyzer) checkUnusedBindings() {
	if a.opts == nil || !a.opts.WarnUnused {
		return
	}
	for ld, n := range a.letReads {
		if n > 0 {
			continue
		}
		a.warn(fmt.Sprintf("%s:%d:%d: unused variable '%s'", ld.Name.File, ld.Name.Line, ld.Name.Col, ld.Name.Lexeme))
	}
	for fd, n := range a.funcReads {
		if n > 0 {
			continue
		}
		if strings.EqualFold(fd.Name.Lexeme, "main") {
			continue
		}
		a.warn(fmt.Sprintf("%s:%d:%d: unused function '%s'", fd.Name.File, fd.Name.Line, fd.Name.Col, fd.Name.Lexeme))
	}
}
