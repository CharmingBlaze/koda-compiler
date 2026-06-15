package parser

import (
	"fmt"
	"koda/internal/lexer"
	"strings"
)

// Node is the base interface for all AST nodes.
type Node interface {
	fmt.Stringer
}

// Expr is the interface for all expression nodes.
type Expr interface {
	Node
	exprNode()
}

// Stmt is the interface for all statement nodes.
type Stmt interface {
	Decl
	stmtNode()
}

// Decl is the interface for all declaration nodes.
type Decl interface {
	Node
	declNode()
}

// Program is the root of the AST for a single source file.
type Program struct {
	Declarations []Decl
}

// ProgramBundle represents a complete multi-module program.
type ProgramBundle struct {
	Entry   *Program
	Modules map[string]*Program // Absolute path -> AST
}

func (p *Program) String() string {
	var b strings.Builder
	for _, d := range p.Declarations {
		b.WriteString(d.String())
		b.WriteByte('\n')
	}
	return b.String()
}

type NativeDirective struct {
	BindingName string // Koda/Koda binding name in source (first token after // koda:extern)
	Symbol      string // C symbol to link
	Arity       int
}

type LetDecl struct {
	Token      lexer.Token
	Name       lexer.Token
	TypeAnnot  string           // optional explicit type, e.g. "i32", "u8"
	IsConst    bool             // true for `const` bindings (immutable)
	Init       Expr             // optional
	Native     *NativeDirective // from // koda:extern (legacy extern directive still accepted)
}

func (d *LetDecl) declNode() {}
func (d *LetDecl) stmtNode() {}
func (d *LetDecl) String() string {
	prefix := ""
	if d.Native != nil {
		prefix = fmt.Sprintf("// koda:extern %s %s %d\n", d.Native.BindingName, d.Native.Symbol, d.Native.Arity)
	}
	if d.Init == nil {
		return prefix + fmt.Sprintf("let %s;", d.Name.Lexeme)
	}
	return prefix + fmt.Sprintf("let %s = %s;", d.Name.Lexeme, d.Init.String())
}

type FuncDecl struct {
	Token  lexer.Token
	Name   lexer.Token
	Params []Param
	Body   *BlockStmt
	Native *NativeDirective
}

// StructDecl declares a named struct type with ordered fields (O(1) slot access at compile time).
type StructDecl struct {
	Token   lexer.Token
	Name    lexer.Token
	Fields  []lexer.Token
	Methods []*FuncDecl
}

func (s *StructDecl) declNode() {}
func (s *StructDecl) String() string { return "struct " + s.Name.Lexeme }

// EnumDecl declares an enum namespace; members are integer constants from 0 upward.
type EnumDecl struct {
	Token   lexer.Token
	Name    lexer.Token
	Members []lexer.Token
}

func (e *EnumDecl) declNode() {}
func (e *EnumDecl) String() string { return "enum " + e.Name.Lexeme }

type Param struct {
	Name    string
	Default Expr // optional
	IsRest  bool
}

func (d *FuncDecl) declNode() {}
func (d *FuncDecl) stmtNode() {}
func (d *FuncDecl) String() string {
	params := []string{}
	for _, p := range d.Params {
		params = append(params, p.Name)
	}
	prefix := ""
	if d.Native != nil {
		prefix = fmt.Sprintf("// koda:extern %s %s %d\n", d.Native.BindingName, d.Native.Symbol, d.Native.Arity)
	}
	return prefix + fmt.Sprintf("func %s(%s) %s", d.Name.Lexeme, strings.Join(params, ", "), d.Body.String())
}

// TestDecl is a named unit test block: test "name" { ... }.
type TestDecl struct {
	Token     lexer.Token
	Display   lexer.Token // string literal token
	SynthName lexer.Token // internal function name (__koda_test_N_slug)
	Body      *BlockStmt
	synth     *FuncDecl
}

func (t *TestDecl) declNode() {}
func (t *TestDecl) String() string {
	return fmt.Sprintf("test %q %s", t.Display.Lexeme, t.Body.String())
}

// SyntheticFunc returns a FuncDecl used for sema/codegen/shadow layout.
func (t *TestDecl) SyntheticFunc() *FuncDecl {
	if t.synth == nil {
		t.synth = &FuncDecl{Token: t.Token, Name: t.SynthName, Body: t.Body}
	}
	return t.synth
}

// Statements

type BlockStmt struct {
	Token        lexer.Token
	Declarations []Decl
}

func (s *BlockStmt) declNode() {}
func (s *BlockStmt) stmtNode() {}
func (s *BlockStmt) String() string {
	var b strings.Builder
	b.WriteString("{\n")
	for _, stmt := range s.Declarations {
		b.WriteString("  ")
		b.WriteString(stmt.String())
		b.WriteByte('\n')
	}
	b.WriteString("}")
	return b.String()
}

type ExpressionStmt struct {
	Token lexer.Token
	Expr  Expr
}

func (s *ExpressionStmt) declNode() {}
func (s *ExpressionStmt) stmtNode() {}
func (s *ExpressionStmt) String() string {
	return s.Expr.String() + ";"
}

type ReturnStmt struct {
	Token lexer.Token
	Value Expr // optional
}

func (s *ReturnStmt) declNode() {}
func (s *ReturnStmt) stmtNode() {}
func (s *ReturnStmt) String() string {
	if s.Value == nil {
		return "return;"
	}
	return fmt.Sprintf("return %s;", s.Value.String())
}

// DeferStmt schedules Expr to run (for side effects) when the enclosing function returns, LIFO vs other defers.
type DeferStmt struct {
	Token lexer.Token
	Expr  Expr
}

func (s *DeferStmt) declNode() {}
func (s *DeferStmt) stmtNode() {}
func (s *DeferStmt) String() string {
	return fmt.Sprintf("defer %s;", s.Expr.String())
}

type IfStmt struct {
	Token     lexer.Token
	Condition Expr
	Then      Stmt
	Else      Stmt // optional
}

func (s *IfStmt) declNode() {}
func (s *IfStmt) stmtNode() {}
func (s *IfStmt) String() string {
	if s.Else == nil {
		return fmt.Sprintf("if (%s) %s", s.Condition.String(), s.Then.String())
	}
	return fmt.Sprintf("if (%s) %s else %s", s.Condition.String(), s.Then.String(), s.Else.String())
}

type WhileStmt struct {
	Token     lexer.Token
	Condition Expr
	Body      Stmt
}

func (s *WhileStmt) declNode() {}
func (s *WhileStmt) stmtNode() {}
func (s *WhileStmt) String() string {
	return fmt.Sprintf("while (%s) %s", s.Condition.String(), s.Body.String())
}

// Expressions

type IdentifierExpr struct {
	Token lexer.Token
	Name  lexer.Token
}

func (e *IdentifierExpr) exprNode()      {}
func (e *IdentifierExpr) String() string { return e.Name.Lexeme }

type LiteralExpr struct {
	Token lexer.Token
	Value any
}

func (e *LiteralExpr) exprNode() {}
func (e *LiteralExpr) String() string {
	return fmt.Sprintf("%v", e.Value)
}

type PrefixExpr struct {
	Token    lexer.Token
	Operator string
	Right    Expr
}

func (e *PrefixExpr) exprNode() {}
func (e *PrefixExpr) String() string {
	return fmt.Sprintf("(%s%s)", e.Operator, e.Right.String())
}

type InfixExpr struct {
	Token    lexer.Token
	Left     Expr
	Operator string
	Right    Expr
}

func (e *InfixExpr) exprNode() {}
func (e *InfixExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", e.Left.String(), e.Operator, e.Right.String())
}

type CallExpr struct {
	Token     lexer.Token
	Function  Expr
	Arguments []Expr
}

func (e *CallExpr) exprNode() {}
func (e *CallExpr) String() string {
	args := []string{}
	for _, a := range e.Arguments {
		args = append(args, a.String())
	}
	return fmt.Sprintf("%s(%s)", e.Function.String(), strings.Join(args, ", "))
}

type AssignExpr struct {
	Token lexer.Token
	Left  Expr
	Value Expr
}

func (e *AssignExpr) exprNode() {}
func (e *AssignExpr) String() string {
	return fmt.Sprintf("%s = %s", e.Left.String(), e.Value.String())
}

type LogicalExpr struct {
	Token    lexer.Token
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

func (e *LogicalExpr) exprNode() {}
func (e *LogicalExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", e.Left.String(), e.Operator.Lexeme, e.Right.String())
}

type ThisExpr struct {
	Token lexer.Token
}

func (e *ThisExpr) exprNode()      {}
func (e *ThisExpr) String() string { return "this" }

type GroupingExpr struct {
	Token lexer.Token
	Expr  Expr
}

func (e *GroupingExpr) exprNode() {}
func (e *GroupingExpr) String() string {
	return fmt.Sprintf("(%s)", e.Expr.String())
}

type UpdateExpr struct {
	Token    lexer.Token
	Operator lexer.Token
	Operand  Expr
	IsPrefix bool
}

func (e *UpdateExpr) exprNode() {}
func (e *UpdateExpr) String() string {
	if e.IsPrefix {
		return fmt.Sprintf("%s%s", e.Operator.Lexeme, e.Operand.String())
	}
	return fmt.Sprintf("%s%s", e.Operand.String(), e.Operator.Lexeme)
}

type RangeExpr struct {
	Token lexer.Token
	From  Expr
	To    Expr
}

func (e *RangeExpr) exprNode() {}
func (e *RangeExpr) String() string {
	return fmt.Sprintf("%s..%s", e.From.String(), e.To.String())
}

type ImportExpr struct {
	Token lexer.Token
	Path  lexer.Token
}

func (e *ImportExpr) exprNode()      {}
func (e *ImportExpr) String() string { return fmt.Sprintf("import %s", e.Path.Lexeme) }

type TemplateExpr struct {
	Token lexer.Token
	Parts []Expr
}

func (e *TemplateExpr) exprNode()      {}
func (e *TemplateExpr) String() string { return "template" }

// SpreadExpr appears inside array literals (`[...xs]`).
type SpreadExpr struct {
	Token lexer.Token
	Expr  Expr
}

func (e *SpreadExpr) exprNode()      {}
func (e *SpreadExpr) String() string { return "…" }

type TupleExpr struct {
	Token    lexer.Token
	Elements []Expr
}

func (e *TupleExpr) exprNode()      {}
func (e *TupleExpr) String() string { return "tuple" }

type IfExpr struct {
	Token     lexer.Token
	Condition Expr
	Then      Expr
	Else      Expr
}

func (e *IfExpr) exprNode()      {}
func (e *IfExpr) String() string { return "if-expr" }

type SwitchExpr struct {
	Token   lexer.Token
	Subject Expr
	Cases   []SwitchCaseExpr
	Default Expr
}

type SwitchCaseExpr struct {
	Value Expr
	Body  Expr
}

func (e *SwitchExpr) exprNode()      {}
func (e *SwitchExpr) String() string { return "switch-expr" }

type SliceExpr struct {
	Token  lexer.Token
	Object Expr
	Start  Expr // optional
	End    Expr // optional
}

func (e *SliceExpr) exprNode()      {}
func (e *SliceExpr) String() string { return "slice" }

type TernaryExpr struct {
	Token     lexer.Token
	Condition Expr
	Then      Expr
	Else      Expr
}

func (e *TernaryExpr) exprNode()      {}
func (e *TernaryExpr) String() string { return "ternary" }

type ArrayExpr struct {
	Token    lexer.Token
	Elements []Expr
}

func (e *ArrayExpr) exprNode()      {}
func (e *ArrayExpr) String() string { return "array" }

type ObjectExpr struct {
	Token        lexer.Token
	StructTag    *lexer.Token // non-nil for `StructName { ... }` struct literals
	Keys         []lexer.Token
	Values       []Expr
	ComputedKeys []Expr
}

func (e *ObjectExpr) exprNode()      {}
func (e *ObjectExpr) String() string { return "object" }

type FuncExpr struct {
	Token  lexer.Token
	Params []Param
	Body   *BlockStmt
}

func (e *FuncExpr) exprNode()      {}
func (e *FuncExpr) declNode()      {} // Params are bound as Decl in semantic analysis
func (e *FuncExpr) String() string { return "func-expr" }

type IndexExpr struct {
	Token    lexer.Token
	Object   Expr
	Index    Expr
	Optional bool // true for `?.` member or `?.[` access
}

func (e *IndexExpr) exprNode()      {}
func (e *IndexExpr) String() string { return "index" }

// More Statements

type ForStmt struct {
	Token      lexer.Token
	Inits      []Decl
	Condition  Expr
	Increments []Expr
	Body       Stmt
}

func (s *ForStmt) declNode()      {}
func (s *ForStmt) stmtNode()      {}
func (s *ForStmt) String() string { return "for" }

type ForInStmt struct {
	Token    lexer.Token
	KeyVar   *lexer.Token
	ValueVar *lexer.Token
	Iterable Expr
	Body     Stmt
}

func (s *ForInStmt) declNode()      {}
func (s *ForInStmt) stmtNode()      {}
func (s *ForInStmt) String() string { return "for-in" }

type BreakStmt struct {
	Token lexer.Token
}

func (s *BreakStmt) declNode()      {}
func (s *BreakStmt) stmtNode()      {}
func (s *BreakStmt) String() string { return "break;" }

type ContinueStmt struct {
	Token lexer.Token
}

func (s *ContinueStmt) declNode()      {}
func (s *ContinueStmt) stmtNode()      {}
func (s *ContinueStmt) String() string { return "continue;" }

// DeleteStmt deletes an own property from a table object (`delete obj["k"]`).
type DeleteStmt struct {
	Token  lexer.Token
	Target Expr
}

func (s *DeleteStmt) declNode()      {}
func (s *DeleteStmt) stmtNode()      {}
func (s *DeleteStmt) String() string { return "delete …;" }

type SwitchStmt struct {
	Token   lexer.Token
	Subject Expr
	Cases   []SwitchCase
	Default []Decl
}

type SwitchCase struct {
	Value Expr
	Body  []Decl
}

func (s *SwitchStmt) declNode()      {}
func (s *SwitchStmt) stmtNode()      {}
func (s *SwitchStmt) String() string { return "switch" }

type DoWhileStmt struct {
	Token     lexer.Token
	Body      Stmt
	Condition Expr
}

func (s *DoWhileStmt) declNode()      {}
func (s *DoWhileStmt) stmtNode()      {}
func (s *DoWhileStmt) String() string { return "do-while" }

type ForOfStmt struct {
	Token   lexer.Token
	VarName lexer.Token // single-var `for-of` value binding, or key in `for (let [k, v] of …)`
	// ValueVar non-nil selects destructuring: bind VarName=key, *ValueVar=value per insertion slot.
	ValueVar *lexer.Token
	Iterable Expr
	Body     Stmt
}

func (s *ForOfStmt) declNode()      {}
func (s *ForOfStmt) stmtNode()      {}
func (s *ForOfStmt) String() string { return "for-of" }

type IncludeDecl struct {
	Token lexer.Token
	Path  lexer.Token
}

func (d *IncludeDecl) declNode()      {}
func (d *IncludeDecl) String() string { return fmt.Sprintf("include %s;", d.Path.Lexeme) }
