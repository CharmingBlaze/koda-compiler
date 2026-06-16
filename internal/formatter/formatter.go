package formatter

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"koda/internal/lexer"
	"koda/internal/parser"
)

// Format parses Koda source and returns canonical spaced formatting.
func Format(src string) (string, error) {
	src = normalizeLegacyOperators(src)
	l := lexer.NewLexer(src, "")
	tokens, err := l.Tokenize()
	if err != nil {
		return "", err
	}
	p := parser.NewParser(tokens)
	prog, err := p.Parse()
	if err != nil {
		return "", err
	}
	var em emitter
	if err := em.formatProgram(prog); err != nil {
		return "", err
	}
	return finalize(em.String()), nil
}

func finalize(s string) string {
	lines := strings.Split(s, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " \t\r")
	}
	s = strings.Join(lines, "\n")
	s = strings.TrimRight(s, "\n")
	return s + "\n"
}

type emitter struct {
	sb     strings.Builder
	indent int
}

func (e *emitter) write(s string) { e.sb.WriteString(s) }

func (e *emitter) lineStart() {
	for i := 0; i < e.indent; i++ {
		e.write("    ")
	}
}

func (e *emitter) String() string { return e.sb.String() }

func (e *emitter) formatProgram(prog *parser.Program) error {
	for i, d := range prog.Declarations {
		if i > 0 && topLevelNeedsBlankLine(prog.Declarations[i-1], d) {
			e.write("\n")
		}
		e.lineStart()
		if err := e.emitDeclCore(d); err != nil {
			return err
		}
	}
	return nil
}

// topLevelNeedsBlankLine is true unless both neighbors are bare expression
// statements (e.g. consecutive print calls), which stay visually grouped.
func topLevelNeedsBlankLine(prev, cur parser.Decl) bool {
	_, prevExpr := prev.(*parser.ExpressionStmt)
	_, prevDefer := prev.(*parser.DeferStmt)
	_, curExpr := cur.(*parser.ExpressionStmt)
	_, curDefer := cur.(*parser.DeferStmt)
	return !((prevExpr || prevDefer) && (curExpr || curDefer))
}

func (e *emitter) emitDeclCore(d parser.Decl) error {
	switch n := d.(type) {
	case *parser.StructDecl:
		e.write("struct ")
		e.write(n.Name.Lexeme)
		e.write(" {\n")
		for _, f := range n.Fields {
			e.write("    ")
			e.write(f.Name.Lexeme)
			if f.Default != nil {
				e.write(" = ")
				if err := e.emitExpr(f.Default, precLowest); err != nil {
					return err
				}
			}
			e.write("\n")
		}
		for _, m := range n.Methods {
			if err := e.emitDeclCore(m); err != nil {
				return err
			}
		}
		e.write("}\n")
		return nil
	case *parser.EnumDecl:
		e.write("enum ")
		e.write(n.Name.Lexeme)
		e.write(" { ")
		for i, m := range n.Members {
			if i > 0 {
				e.write(", ")
			}
			e.write(m.Lexeme)
		}
		e.write(" }\n")
		return nil
	case *parser.IncludeDecl:
		e.write("#include ")
		e.write(n.Path.Lexeme)
		e.write("\n")
		return nil
	case *parser.LetDecl:
		if n.Native != nil {
			e.write("// koda:extern ")
			e.write(n.Native.BindingName)
			e.write(" ")
			e.write(n.Native.Symbol)
			e.write(" ")
			e.write(strconv.Itoa(n.Native.Arity))
			e.write("\n")
			e.lineStart()
		}
		e.write("let ")
		e.write(n.Name.Lexeme)
		if n.Init != nil {
			e.write(" = ")
			if err := e.emitExpr(n.Init, precLowest); err != nil {
				return err
			}
		}
		e.write(";\n")
		return nil
	case *parser.FuncDecl:
		if n.Native != nil {
			e.write("// koda:extern ")
			e.write(n.Native.BindingName)
			e.write(" ")
			e.write(n.Native.Symbol)
			e.write(" ")
			e.write(strconv.Itoa(n.Native.Arity))
			e.write("\n")
			e.lineStart()
		}
		e.write("func ")
		e.write(n.Name.Lexeme)
		e.write("(")
		for i, p := range n.Params {
			if i > 0 {
				e.write(", ")
			}
			if p.IsRest {
				e.write("...")
			}
			e.write(p.Name)
			if p.Default != nil {
				e.write(" = ")
				if err := e.emitExpr(p.Default, precLowest); err != nil {
					return err
				}
			}
		}
		e.write(") ")
		return e.emitBlock(n.Body, true, false, true)
	case *parser.ExpressionStmt:
		if err := e.emitExpr(n.Expr, precLowest); err != nil {
			return err
		}
		e.write(";\n")
		return nil
	case *parser.FuncExpr:
		if err := e.emitExpr(n, precLowest); err != nil {
			return err
		}
		e.write(";\n")
		return nil
	case parser.Stmt:
		return e.emitStmt(n)
	default:
		return fmt.Errorf("formatter: unsupported decl %T", d)
	}
}

// emitBlock formats a block. If braceSameLine, "{" follows ")" or "else " on the same line.
// If openLineIndented is true, the opening "{" begins the current line (caller already called lineStart).
// If newlineAfterClose is false, the closing "}" is not followed by a newline (caller adds ";\n", etc.).
func (e *emitter) emitBlock(b *parser.BlockStmt, braceSameLine, openLineIndented, newlineAfterClose bool) error {
	if braceSameLine {
		e.write("{\n")
	} else if openLineIndented {
		e.write("{\n")
	} else {
		e.lineStart()
		e.write("{\n")
	}
	e.indent++
	for _, d := range b.Declarations {
		e.lineStart()
		if err := e.emitDeclCore(d); err != nil {
			return err
		}
	}
	e.indent--
	e.lineStart()
	e.write("}")
	if newlineAfterClose {
		e.write("\n")
	}
	return nil
}

func (e *emitter) emitStmt(s parser.Stmt) error {
	switch n := s.(type) {
	case *parser.BlockStmt:
		return e.emitBlock(n, false, true, true)
	case *parser.ExpressionStmt:
		if err := e.emitExpr(n.Expr, precLowest); err != nil {
			return err
		}
		e.write(";\n")
		return nil
	case *parser.ReturnStmt:
		e.write("return")
		if n.Value != nil {
			e.write(" ")
			if err := e.emitExpr(n.Value, precLowest); err != nil {
				return err
			}
		}
		e.write(";\n")
		return nil
	case *parser.DeferStmt:
		e.write("defer ")
		if err := e.emitExpr(n.Expr, precLowest); err != nil {
			return err
		}
		e.write(";\n")
		return nil
	case *parser.IfStmt:
		e.write("if (")
		if err := e.emitExpr(n.Condition, precLowest); err != nil {
			return err
		}
		e.write(") ")
		if n.Else != nil {
			if tb, ok := n.Then.(*parser.BlockStmt); ok {
				// Keep "} else {" on one line (no newline after then-block's "}").
				if err := e.emitBlock(tb, true, false, false); err != nil {
					return err
				}
				e.write(" else ")
				return e.emitStmtAsChild(n.Else)
			}
		}
		return e.emitStmtAsChild(n.Then)
	case *parser.WhileStmt:
		e.write("while (")
		if err := e.emitExpr(n.Condition, precLowest); err != nil {
			return err
		}
		e.write(") ")
		return e.emitStmtAsChild(n.Body)
	case *parser.DoWhileStmt:
		e.write("do ")
		if err := e.emitStmtAsChild(n.Body); err != nil {
			return err
		}
		e.write(" while (")
		if err := e.emitExpr(n.Condition, precLowest); err != nil {
			return err
		}
		e.write(");\n")
		return nil
	case *parser.ForStmt:
		e.write("for (")
		for i, ini := range n.Inits {
			if i > 0 {
				e.write(", ")
			}
			let, ok := ini.(*parser.LetDecl)
			if !ok {
				return fmt.Errorf("formatter: unsupported for-loop initializer")
			}
			e.write("let ")
			e.write(let.Name.Lexeme)
			if let.Init != nil {
				e.write(" = ")
				if err := e.emitExpr(let.Init, precLowest); err != nil {
					return err
				}
			}
		}
		e.write("; ")
		if n.Condition != nil {
			if err := e.emitExpr(n.Condition, precLowest); err != nil {
				return err
			}
		}
		e.write("; ")
		for i, inc := range n.Increments {
			if i > 0 {
				e.write(", ")
			}
			if err := e.emitExpr(inc, precLowest); err != nil {
				return err
			}
		}
		e.write(") ")
		return e.emitStmtAsChild(n.Body)
	case *parser.ForInStmt:
		e.write("for (let ")
		e.write(n.KeyVar.Lexeme)
		e.write(" in ")
		if err := e.emitExpr(n.Iterable, precLowest); err != nil {
			return err
		}
		e.write(") ")
		return e.emitStmtAsChild(n.Body)
	case *parser.ForOfStmt:
		if n.ValueVar != nil {
			e.write("for (let [")
			e.write(n.VarName.Lexeme)
			e.write(", ")
			e.write(n.ValueVar.Lexeme)
			e.write("] of ")
			if err := e.emitExpr(n.Iterable, precLowest); err != nil {
				return err
			}
			e.write(") ")
			return e.emitStmtAsChild(n.Body)
		}
		e.write("for ")
		e.write(n.VarName.Lexeme)
		e.write(" in ")
		if err := e.emitExpr(n.Iterable, precLowest); err != nil {
			return err
		}
		e.write(" ")
		return e.emitStmtAsChild(n.Body)
	case *parser.BreakStmt:
		e.write("break;\n")
		return nil
	case *parser.ContinueStmt:
		e.write("continue;\n")
		return nil
	case *parser.FallthroughStmt:
		e.write("fallthrough;\n")
		return nil
	case *parser.SwitchStmt:
		if n.Token.Type == lexer.TokenMatch {
			e.write("match ")
			if err := e.emitExpr(n.Subject, precLowest); err != nil {
				return err
			}
			e.write(" {\n")
			e.indent++
			for _, c := range n.Cases {
				e.lineStart()
				if err := e.emitExpr(c.Value, precLowest); err != nil {
					return err
				}
				e.write(" {\n")
				e.indent++
				for _, bd := range c.Body {
					e.lineStart()
					if err := e.emitDeclCore(bd); err != nil {
						return err
					}
				}
				e.indent--
				e.lineStart()
				e.write("}\n")
			}
			if len(n.Default) > 0 {
				e.lineStart()
				e.write("default {\n")
				e.indent++
				for _, bd := range n.Default {
					e.lineStart()
					if err := e.emitDeclCore(bd); err != nil {
						return err
					}
				}
				e.indent--
				e.lineStart()
				e.write("}\n")
			}
			e.indent--
			e.lineStart()
			e.write("}\n")
			return nil
		}
		e.write("switch (")
		if err := e.emitExpr(n.Subject, precLowest); err != nil {
			return err
		}
		e.write(") {\n")
		e.indent++
		for _, c := range n.Cases {
			e.lineStart()
			e.write("case ")
			if err := e.emitExpr(c.Value, precLowest); err != nil {
				return err
			}
			e.write(":\n")
			e.indent++
			for _, bd := range c.Body {
				e.lineStart()
				if err := e.emitDeclCore(bd); err != nil {
					return err
				}
			}
			e.indent--
		}
		if len(n.Default) > 0 {
			e.lineStart()
			e.write("default:\n")
			e.indent++
			for _, bd := range n.Default {
				e.lineStart()
				if err := e.emitDeclCore(bd); err != nil {
					return err
				}
			}
			e.indent--
		}
		e.indent--
		e.lineStart()
		e.write("}\n")
		return nil
	default:
		return fmt.Errorf("formatter: unsupported stmt %T", s)
	}
}

func (e *emitter) emitStmtAsChild(s parser.Stmt) error {
	if _, ok := s.(*parser.BlockStmt); ok {
		return e.emitBlock(s.(*parser.BlockStmt), true, false, true)
	}
	e.write("\n")
	e.indent++
	e.lineStart()
	if err := e.emitStmt(s); err != nil {
		return err
	}
	e.indent--
	return nil
}

func (e *emitter) emitExpr(ex parser.Expr, minPrec int) error {
	switch v := ex.(type) {
	case *parser.LiteralExpr:
		e.write(literalText(v))
		return nil
	case *parser.IdentifierExpr:
		e.write(v.Name.Lexeme)
		return nil
	case *parser.ThisExpr:
		e.write("this")
		return nil
	case *parser.GroupingExpr:
		e.write("(")
		if err := e.emitExpr(v.Expr, precLowest); err != nil {
			return err
		}
		e.write(")")
		return nil
	case *parser.PrefixExpr:
		p := precPrefix
		if p < minPrec {
			e.write("(")
			if err := e.emitExpr(ex, precLowest); err != nil {
				return err
			}
			e.write(")")
			return nil
		}
		e.write(v.Operator)
		return e.emitExpr(v.Right, precPrefix)
	case *parser.InfixExpr:
		opPrec := precOf(v.Token.Type)
		if opPrec < minPrec {
			e.write("(")
			if err := e.emitExpr(ex, precLowest); err != nil {
				return err
			}
			e.write(")")
			return nil
		}
		if err := e.emitExpr(v.Left, opPrec); err != nil {
			return err
		}
		e.write(" ")
		op := strings.TrimSpace(v.Operator)
		switch op {
		case "===":
			op = "=="
		case "!==":
			op = "!="
		}
		e.write(op)
		e.write(" ")
		return e.emitExpr(v.Right, opPrec+1)
	case *parser.LogicalExpr:
		opPrec := precOf(v.Operator.Type)
		if opPrec < minPrec {
			e.write("(")
			if err := e.emitExpr(ex, precLowest); err != nil {
				return err
			}
			e.write(")")
			return nil
		}
		if err := e.emitExpr(v.Left, opPrec); err != nil {
			return err
		}
		e.write(" ")
		e.write(v.Operator.Lexeme)
		e.write(" ")
		return e.emitExpr(v.Right, opPrec+1)
	case *parser.AssignExpr:
		if precAssign < minPrec {
			e.write("(")
			if err := e.emitExpr(ex, precLowest); err != nil {
				return err
			}
			e.write(")")
			return nil
		}
		if err := e.emitExpr(v.Left, precAssign); err != nil {
			return err
		}
		e.write(" ")
		e.write(strings.TrimSpace(v.Token.Lexeme))
		e.write(" ")
		return e.emitExpr(v.Value, precLowest)
	case *parser.CallExpr:
		if precCall < minPrec {
			e.write("(")
			if err := e.emitExpr(ex, precLowest); err != nil {
				return err
			}
			e.write(")")
			return nil
		}
		if err := e.emitExpr(v.Function, precCall); err != nil {
			return err
		}
		e.write("(")
		for i, a := range v.Arguments {
			if i > 0 {
				e.write(", ")
			}
			if err := e.emitExpr(a, precLowest); err != nil {
				return err
			}
		}
		e.write(")")
		return nil
	case *parser.IndexExpr:
		if precIndex < minPrec {
			e.write("(")
			if err := e.emitExpr(ex, precLowest); err != nil {
				return err
			}
			e.write(")")
			return nil
		}
		if err := e.emitExpr(v.Object, precIndex); err != nil {
			return err
		}
		if lit, ok := v.Index.(*parser.LiteralExpr); ok {
			if s, ok := lit.Value.(string); ok && isIdent(s) && lit.Token.Type == lexer.TokenIdentifier {
				e.write(".")
				e.write(s)
				return nil
			}
		}
		e.write("[")
		if err := e.emitExpr(v.Index, precLowest); err != nil {
			return err
		}
		e.write("]")
		return nil
	case *parser.ArrayExpr:
		e.write("[")
		for i, el := range v.Elements {
			if i > 0 {
				e.write(", ")
			}
			if err := e.emitExpr(el, precLowest); err != nil {
				return err
			}
		}
		e.write("]")
		return nil
	case *parser.ObjectExpr:
		if v.StructTag != nil {
			e.write(v.StructTag.Lexeme)
			e.write(" ")
		}
		e.write("{")
		if len(v.Keys) > 0 {
			e.write(" ")
		}
		for i := range v.Keys {
			if i > 0 {
				e.write(", ")
			}
			e.write(v.Keys[i].Lexeme)
			e.write(": ")
			if err := e.emitExpr(v.Values[i], precLowest); err != nil {
				return err
			}
		}
		if len(v.Keys) > 0 {
			e.write(" ")
		}
		e.write("}")
		return nil
	case *parser.FuncExpr:
		e.write("func(")
		for i, p := range v.Params {
			if i > 0 {
				e.write(", ")
			}
			if p.IsRest {
				e.write("...")
			}
			e.write(p.Name)
			if p.Default != nil {
				e.write(" = ")
				if err := e.emitExpr(p.Default, precLowest); err != nil {
					return err
				}
			}
		}
		e.write(") ")
		return e.emitBlock(v.Body, true, false, false)
	case *parser.UpdateExpr:
		if v.IsPrefix {
			e.write(v.Operator.Lexeme)
			return e.emitExpr(v.Operand, precCall)
		}
		if err := e.emitExpr(v.Operand, precCall); err != nil {
			return err
		}
		e.write(v.Operator.Lexeme)
		return nil
	case *parser.RangeExpr:
		opPrec := precEquals
		if opPrec < minPrec {
			e.write("(")
			if err := e.emitExpr(ex, precLowest); err != nil {
				return err
			}
			e.write(")")
			return nil
		}
		if err := e.emitExpr(v.From, opPrec); err != nil {
			return err
		}
		e.write("..")
		return e.emitExpr(v.To, opPrec+1)
	case *parser.ImportExpr:
		e.write("import ")
		e.write(v.Path.Lexeme)
		return nil
	default:
		return fmt.Errorf("formatter: unsupported expr %T", ex)
	}
}

func isIdent(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 && !unicode.IsLetter(r) && r != '_' {
			return false
		}
		if i > 0 && !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}

func literalText(e *parser.LiteralExpr) string {
	switch e.Token.Type {
	case lexer.TokenString:
		return e.Token.Lexeme
	case lexer.TokenNumber:
		if e.Token.Lexeme != "" {
			return e.Token.Lexeme
		}
	case lexer.TokenTrue:
		return "true"
	case lexer.TokenFalse:
		return "false"
	case lexer.TokenNull:
		return "null"
	}
	switch v := e.Value.(type) {
	case string:
		return strconv.Quote(v)
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case nil:
		return "null"
	}
	return e.Token.Lexeme
}
