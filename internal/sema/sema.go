package sema

import (
	"fmt"
	"strings"

	"koda/internal/diagnostic"
	"koda/internal/parser"
)

// Analyzer performs semantic analysis on the AST.
type Analyzer struct {
	currentScope *Scope
	scopes       []*Scope
	errors       []error

	structLayouts map[string][]string // struct type name -> ordered field names
	varStructType map[string]string   // variable name -> struct type name
	varEnumType   map[string]string   // variable name -> enum type name
	warnings      []string

	indexExprStructSlot map[*parser.IndexExpr]int
	indexExprEnumConst  map[*parser.IndexExpr]int64
}

// Scope represents a lexical scope with symbol bindings.
type Scope struct {
	parent  *Scope
	symbols map[string]parser.Decl
}

// NewScope creates a new scope.
func NewScope(parent *Scope) *Scope {
	return &Scope{
		parent:  parent,
		symbols: make(map[string]parser.Decl),
	}
}

// Define adds a symbol to the scope.
func (s *Scope) Define(name string, decl parser.Decl) {
	s.symbols[name] = decl
}

// Resolve looks up a symbol in the scope chain.
func (s *Scope) Resolve(name string) (parser.Decl, bool) {
	if decl, ok := s.symbols[name]; ok {
		return decl, true
	}
	if s.parent != nil {
		return s.parent.Resolve(name)
	}
	return nil, false
}

func (s *Scope) VisibleNames() []string {
	out := []string{}
	seen := map[string]bool{}
	for cur := s; cur != nil; cur = cur.parent {
		for name := range cur.symbols {
			if seen[name] {
				continue
			}
			seen[name] = true
			out = append(out, name)
		}
	}
	return out
}

// NewAnalyzer creates a new semantic analyzer.
func NewAnalyzer() *Analyzer {
	builtinRoot := NewScope(nil)
	seedGlobalBuiltins(builtinRoot)
	globalScope := NewScope(builtinRoot)
	return &Analyzer{
		currentScope:        globalScope,
		scopes:              []*Scope{builtinRoot, globalScope},
		errors:              []error{},
		structLayouts:       make(map[string][]string),
		varStructType:       make(map[string]string),
		varEnumType:         make(map[string]string),
		indexExprStructSlot: make(map[*parser.IndexExpr]int),
		indexExprEnumConst:  make(map[*parser.IndexExpr]int64),
	}
}

// Analyze performs semantic analysis on a program.
func (a *Analyzer) Analyze(prog *parser.Program) error {
	a.errors = nil
	a.warnings = nil
	a.structLayouts = make(map[string][]string)
	a.varStructType = make(map[string]string)
	a.varEnumType = make(map[string]string)
	a.indexExprStructSlot = make(map[*parser.IndexExpr]int)
	a.indexExprEnumConst = make(map[*parser.IndexExpr]int64)
	for _, decl := range prog.Declarations {
		a.analyzeDecl(decl)
	}
	switch len(a.errors) {
	case 0:
		return nil
	case 1:
		return a.errors[0]
	default:
		return &diagnostic.MultiError{List: append([]error(nil), a.errors...)}
	}
}

// Errors returns all errors found during analysis.
func (a *Analyzer) Errors() []error {
	return a.errors
}

// Warnings returns non-fatal diagnostics (e.g. switch exhaustiveness on enums).
func (a *Analyzer) Warnings() []string {
	return a.warnings
}

func (a *Analyzer) warn(msg string) {
	a.warnings = append(a.warnings, msg)
}

func (a *Analyzer) record(err error) {
	if err != nil {
		a.errors = append(a.errors, err)
	}
}

func (a *Analyzer) enterScope() {
	newScope := NewScope(a.currentScope)
	a.scopes = append(a.scopes, newScope)
	a.currentScope = newScope
}

func (a *Analyzer) exitScope() {
	if len(a.scopes) > 1 {
		a.scopes = a.scopes[:len(a.scopes)-1]
		a.currentScope = a.scopes[len(a.scopes)-1]
	}
}

func (a *Analyzer) analyzeDecl(decl parser.Decl) {
	switch d := decl.(type) {
	case *parser.LetDecl:
		a.analyzeLetDecl(d)
	case *parser.FuncDecl:
		a.analyzeFuncDecl(d)
	case *parser.FuncExpr:
		a.analyzeFuncExpr(d)
	case *parser.StructDecl:
		a.analyzeStructDecl(d)
	case *parser.EnumDecl:
		a.analyzeEnumDecl(d)
	case *parser.IncludeDecl:
		return
	case parser.Stmt:
		a.analyzeStmt(d)
	}
}

func (a *Analyzer) analyzeLetDecl(d *parser.LetDecl) {
	name := d.Name.Lexeme
	if _, ok := a.currentScope.symbols[name]; ok {
		a.record(&diagnostic.DiagnosticError{
			File:    d.Name.File,
			Line:    d.Name.Line,
			Col:     d.Name.Col,
			Message: fmt.Sprintf("duplicate binding '%s' in the same scope", name),
		})
		return
	}
	a.currentScope.Define(name, d)
	if d.Init != nil {
		a.analyzeExpr(d.Init)
		a.recordVarTypesFromInit(name, d.Init)
	}
}

func (a *Analyzer) analyzeFuncDecl(d *parser.FuncDecl) {
	name := d.Name.Lexeme
	if _, ok := a.currentScope.symbols[name]; ok {
		a.record(&diagnostic.DiagnosticError{
			File:    d.Name.File,
			Line:    d.Name.Line,
			Col:     d.Name.Col,
			Message: fmt.Sprintf("duplicate function '%s' in the same scope", name),
		})
		return
	}
	a.currentScope.Define(name, d)

	a.enterScope()
	defer a.exitScope()

	for _, param := range d.Params {
		a.currentScope.Define(param.Name, d)
		if param.Default != nil {
			a.analyzeExpr(param.Default)
		}
	}

	a.analyzeStmt(d.Body)
}

func (a *Analyzer) analyzeFuncExpr(e *parser.FuncExpr) {
	a.enterScope()
	defer a.exitScope()
	for _, param := range e.Params {
		a.currentScope.Define(param.Name, e)
		if param.Default != nil {
			a.analyzeExpr(param.Default)
		}
	}
	a.analyzeStmt(e.Body)
}

func (a *Analyzer) analyzeStmt(stmt parser.Stmt) {
	switch s := stmt.(type) {
	case *parser.BlockStmt:
		a.analyzeBlockStmt(s)
	case *parser.ExpressionStmt:
		a.analyzeExpr(s.Expr)
	case *parser.ReturnStmt:
		if s.Value != nil {
			a.analyzeExpr(s.Value)
		}
	case *parser.IfStmt:
		a.analyzeExpr(s.Condition)
		a.analyzeStmt(s.Then)
		if s.Else != nil {
			a.analyzeStmt(s.Else)
		}
	case *parser.WhileStmt:
		a.analyzeExpr(s.Condition)
		a.analyzeStmt(s.Body)
	case *parser.DoWhileStmt:
		a.analyzeStmt(s.Body)
		a.analyzeExpr(s.Condition)
	case *parser.ForStmt:
		for _, ini := range s.Inits {
			a.analyzeDecl(ini)
		}
		if s.Condition != nil {
			a.analyzeExpr(s.Condition)
		}
		for _, inc := range s.Increments {
			a.analyzeExpr(inc)
		}
		a.analyzeStmt(s.Body)
	case *parser.ForInStmt:
		a.analyzeExpr(s.Iterable)
		a.enterScope()
		if s.KeyVar != nil {
			a.currentScope.Define(s.KeyVar.Lexeme, s)
		}
		if s.ValueVar != nil {
			a.currentScope.Define(s.ValueVar.Lexeme, s)
		}
		a.analyzeStmt(s.Body)
		a.exitScope()
	case *parser.ForOfStmt:
		a.analyzeExpr(s.Iterable)
		a.enterScope()
		a.currentScope.Define(s.VarName.Lexeme, s)
		if s.ValueVar != nil {
			a.currentScope.Define(s.ValueVar.Lexeme, s)
		}
		a.analyzeStmt(s.Body)
		a.exitScope()
	case *parser.SwitchStmt:
		a.analyzeExpr(s.Subject)
		for _, c := range s.Cases {
			a.analyzeExpr(c.Value)
			for _, cd := range c.Body {
				a.analyzeDecl(cd)
			}
		}
		for _, cd := range s.Default {
			a.analyzeDecl(cd)
		}
		a.checkSwitchEnumExhaustive(s)
	case *parser.DeleteStmt:
		a.analyzeExpr(s.Target)
	case *parser.DeferStmt:
		a.analyzeExpr(s.Expr)
	case *parser.BreakStmt, *parser.ContinueStmt:
		return
	default:
		a.record(fmt.Errorf("unsupported statement type: %T", stmt))
	}
}

func (a *Analyzer) analyzeBlockStmt(s *parser.BlockStmt) {
	a.enterScope()
	defer a.exitScope()

	for _, decl := range s.Declarations {
		a.analyzeDecl(decl)
	}
}

// suggestName returns a hint for a typo'd identifier using case-folded edit distance (< 3).
func (a *Analyzer) suggestName(name string) string {
	ln := strings.ToLower(name)
	best, bestDist := "", 3
	for sc := a.currentScope; sc != nil; sc = sc.parent {
		for candidate := range sc.symbols {
			if candidate == name {
				continue
			}
			d := levenshtein(ln, strings.ToLower(candidate))
			if d < bestDist {
				best, bestDist = candidate, d
			}
		}
	}
	if best != "" {
		return fmt.Sprintf("did you mean '%s'?", best)
	}
	return ""
}

func (a *Analyzer) analyzeExpr(expr parser.Expr) {
	switch e := expr.(type) {
	case *parser.IdentifierExpr:
		name := e.Name.Lexeme
		if _, ok := a.currentScope.Resolve(name); !ok {
			hint := a.suggestName(name)
			a.record(&diagnostic.DiagnosticError{
				File:    e.Name.File,
				Line:    e.Name.Line,
				Col:     e.Name.Col,
				Message: fmt.Sprintf("undefined variable '%s'", name),
				Hint:    hint,
			})
		}
	case *parser.LiteralExpr:
		return
	case *parser.PrefixExpr:
		a.analyzeExpr(e.Right)
	case *parser.InfixExpr:
		a.analyzeExpr(e.Left)
		a.analyzeExpr(e.Right)
	case *parser.LogicalExpr:
		a.analyzeExpr(e.Left)
		a.analyzeExpr(e.Right)
	case *parser.CallExpr:
		a.analyzeExpr(e.Function)
		for _, arg := range e.Arguments {
			a.analyzeExpr(arg)
		}
		a.maybeCheckCallArity(e)
	case *parser.AssignExpr:
		a.analyzeExpr(e.Value)
		if ident, ok := e.Left.(*parser.IdentifierExpr); ok {
			name := ident.Name.Lexeme
			if _, ok := a.currentScope.Resolve(name); !ok {
				hint := a.suggestName(name)
				a.record(&diagnostic.DiagnosticError{
					File:    ident.Name.File,
					Line:    ident.Name.Line,
					Col:     ident.Name.Col,
					Message: fmt.Sprintf("undefined variable '%s'", name),
					Hint:    hint,
				})
			}
			return
		}
		if ix, ok := e.Left.(*parser.IndexExpr); ok {
			a.analyzeExpr(ix.Object)
			a.analyzeExpr(ix.Index)
			a.checkStructFieldAccess(ix)
			return
		}
		a.record(&diagnostic.DiagnosticError{
			File:    e.Token.File,
			Line:    e.Token.Line,
			Col:     e.Token.Col,
			Message: "invalid assignment target",
			Hint:    "left side of '=' must be a variable or index expression",
		})
	case *parser.GroupingExpr:
		a.analyzeExpr(e.Expr)
	case *parser.ImportExpr:
		return
	case *parser.IndexExpr:
		a.analyzeExpr(e.Object)
		a.analyzeExpr(e.Index)
		a.checkEnumMemberAccess(e)
		a.checkStructFieldAccess(e)
	case *parser.SpreadExpr:
		a.analyzeExpr(e.Expr)
	case *parser.TemplateExpr:
		for _, p := range e.Parts {
			a.analyzeExpr(p)
		}
	case *parser.ThisExpr:
		return
	case *parser.ArrayExpr:
		for _, el := range e.Elements {
			a.analyzeExpr(el)
		}
	case *parser.ObjectExpr:
		if e.StructTag != nil {
			a.validateStructLiteral(e)
		}
		for _, v := range e.Values {
			a.analyzeExpr(v)
		}
		for _, ck := range e.ComputedKeys {
			if ck == nil {
				continue
			}
			a.analyzeExpr(ck)
		}
	case *parser.FuncExpr:
		a.analyzeFuncExpr(e)
	case *parser.RangeExpr:
		a.analyzeExpr(e.From)
		a.analyzeExpr(e.To)
	case *parser.UpdateExpr:
		a.analyzeExpr(e.Operand)
	case *parser.TupleExpr:
		for _, el := range e.Elements {
			a.analyzeExpr(el)
		}
	case *parser.IfExpr:
		a.analyzeExpr(e.Condition)
		a.analyzeExpr(e.Then)
		if e.Else != nil {
			a.analyzeExpr(e.Else)
		}
	case *parser.SwitchExpr:
		a.analyzeExpr(e.Subject)
		for _, c := range e.Cases {
			a.analyzeExpr(c.Value)
			a.analyzeExpr(c.Body)
		}
		if e.Default != nil {
			a.analyzeExpr(e.Default)
		}
	case *parser.SliceExpr:
		a.analyzeExpr(e.Object)
		if e.Start != nil {
			a.analyzeExpr(e.Start)
		}
		if e.End != nil {
			a.analyzeExpr(e.End)
		}
	case *parser.TernaryExpr:
		a.analyzeExpr(e.Condition)
		a.analyzeExpr(e.Then)
		a.analyzeExpr(e.Else)
	default:
		a.record(fmt.Errorf("unsupported expression type: %T", expr))
	}
}
