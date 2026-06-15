package sema

import (
	"fmt"
	"strings"

	"koda/internal/diagnostic"
	"koda/internal/parser"
)

// Analyzer performs semantic analysis on the AST.
type Analyzer struct {
	opts *AnalysisOptions
	currentScope *Scope
	scopes       []*Scope
	errors       []error

	structLayouts map[string][]string // struct type name -> ordered field names
	structMethods map[string]map[string]*parser.FuncDecl
	varStructType map[string]string   // variable name -> struct type name
	varEnumType   map[string]string   // variable name -> enum type name
	// funcParamStruct maps function name -> parameter name -> struct type (from call sites).
	funcParamStruct   map[string]map[string]string
	activeParamStruct map[string]string // param struct types while refining a function body
	warnings      []string

	indexExprStructSlot map[*parser.IndexExpr]int
	indexExprEnumConst  map[*parser.IndexExpr]int64

	structFieldDefaults   map[string]map[string]parser.Expr // struct -> field -> default expr
	implicitStructField   map[*parser.IdentifierExpr]int    // bare field refs in struct methods
	varArrayElementStruct map[string]string                 // array var -> element struct type

	letReads  map[*parser.LetDecl]int
	funcReads map[*parser.FuncDecl]int

	currentStructType string // set while analyzing a struct method body
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
	return NewAnalyzerWithOptions(nil)
}

// NewAnalyzerWithOptions creates an analyzer with optional passes enabled.
func NewAnalyzerWithOptions(opts *AnalysisOptions) *Analyzer {
	builtinRoot := NewScope(nil)
	seedGlobalBuiltins(builtinRoot)
	globalScope := NewScope(builtinRoot)
	return &Analyzer{
		opts:                opts,
		currentScope:        globalScope,
		scopes:              []*Scope{builtinRoot, globalScope},
		errors:              []error{},
		structLayouts:       make(map[string][]string),
		structMethods:       make(map[string]map[string]*parser.FuncDecl),
		varStructType:       make(map[string]string),
		varEnumType:         make(map[string]string),
		indexExprStructSlot: make(map[*parser.IndexExpr]int),
		indexExprEnumConst:  make(map[*parser.IndexExpr]int64),
		structFieldDefaults:   make(map[string]map[string]parser.Expr),
		implicitStructField:   make(map[*parser.IdentifierExpr]int),
		varArrayElementStruct: make(map[string]string),
		letReads:            make(map[*parser.LetDecl]int),
		funcReads:           make(map[*parser.FuncDecl]int),
	}
}

// Analyze performs semantic analysis on a program.
func (a *Analyzer) Analyze(prog *parser.Program) error {
	a.errors = nil
	a.warnings = nil
	a.structLayouts = make(map[string][]string)
	a.structMethods = make(map[string]map[string]*parser.FuncDecl)
	a.varStructType = make(map[string]string)
	a.varEnumType = make(map[string]string)
	a.funcParamStruct = make(map[string]map[string]string)
	a.activeParamStruct = nil
	a.indexExprStructSlot = make(map[*parser.IndexExpr]int)
	a.indexExprEnumConst = make(map[*parser.IndexExpr]int64)
	a.structFieldDefaults = make(map[string]map[string]parser.Expr)
	a.implicitStructField = make(map[*parser.IdentifierExpr]int)
	a.varArrayElementStruct = make(map[string]string)
	a.letReads = make(map[*parser.LetDecl]int)
	a.funcReads = make(map[*parser.FuncDecl]int)
	a.currentStructType = ""
	for _, decl := range prog.Declarations {
		a.analyzeDecl(decl)
	}
	a.refineStructFieldAccessFromCalls(prog)
	a.checkUnusedBindings()
	if a.opts != nil && a.opts.WarnUnreachable {
		for _, decl := range prog.Declarations {
			if fd, ok := decl.(*parser.FuncDecl); ok && fd.Body != nil {
				a.checkUnreachableCode(fd.Body)
			}
			if td, ok := decl.(*parser.TestDecl); ok && td.Body != nil {
				a.checkUnreachableCode(td.Body)
			}
		}
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
	case *parser.TestDecl:
		a.analyzeTestDecl(d)
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
	if d.TypeAnnot != "" {
		if !isKnownTypeAnnotation(d.TypeAnnot) {
			a.record(&diagnostic.DiagnosticError{
				File:    d.Name.File,
				Line:    d.Name.Line,
				Col:     d.Name.Col,
				Message: fmt.Sprintf("unknown type '%s' (supported: int, float, bool, string, byte, i32, u8, …)", d.TypeAnnot),
			})
		}
	}
	if d.IsConst && d.Init == nil {
		a.record(&diagnostic.DiagnosticError{
			File:    d.Name.File,
			Line:    d.Name.Line,
			Col:     d.Name.Col,
			Message: fmt.Sprintf("const '%s' requires an initializer", name),
		})
	}
	// Analyze initializer before binding the name so enum/type names are not shadowed (e.g. let phase = Phase.Menu).
	if d.Init != nil {
		a.analyzeExpr(d.Init)
		a.recordVarTypesFromInit(name, d.Init)
	}
	a.currentScope.Define(name, d)
	a.letReads[d] = 0
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
	a.funcReads[d] = 0

	a.enterScope()
	defer a.exitScope()

	for _, param := range d.Params {
		paramTok := d.Name
		paramTok.Lexeme = param.Name
		a.currentScope.Define(param.Name, &parser.LetDecl{Name: paramTok})
		if param.Default != nil {
			a.analyzeExpr(param.Default)
		}
	}

	a.analyzeStmt(d.Body)
}

func (a *Analyzer) analyzeTestDecl(d *parser.TestDecl) {
	a.enterScope()
	defer a.exitScope()
	a.analyzeStmt(d.Body)
}

func (a *Analyzer) analyzeFuncExpr(e *parser.FuncExpr) {
	a.enterScope()
	defer a.exitScope()
	for _, param := range e.Params {
		paramTok := e.Token
		paramTok.Lexeme = param.Name
		a.currentScope.Define(param.Name, &parser.LetDecl{Name: paramTok})
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
		if id, ok := s.Iterable.(*parser.IdentifierExpr); ok {
			if st, ok := a.varArrayElementStruct[id.Name.Lexeme]; ok {
				a.varStructType[s.VarName.Lexeme] = st
			}
		}
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
		if a.opts != nil && a.opts.BeginnerLint {
			a.warn(fmt.Sprintf("%s:%d:%d: delete is an advanced feature; prefer struct fields or new objects for game data",
				s.Token.File, s.Token.Line, s.Token.Col))
		}
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

func (a *Analyzer) noteStructMethodRead(ix *parser.IndexExpr) {
	idObj, ok := ix.Object.(*parser.IdentifierExpr)
	if !ok {
		return
	}
	lit, ok := ix.Index.(*parser.LiteralExpr)
	if !ok {
		return
	}
	mname, ok := lit.Value.(string)
	if !ok {
		return
	}
	stName, ok := a.varStructType[idObj.Name.Lexeme]
	if !ok {
		return
	}
	methods, ok := a.structMethods[stName]
	if !ok {
		return
	}
	if fd, ok := methods[mname]; ok {
		a.funcReads[fd]++
	}
}

func (a *Analyzer) analyzeExpr(expr parser.Expr) {
	switch e := expr.(type) {
	case *parser.IdentifierExpr:
		name := e.Name.Lexeme
		if decl, ok := a.currentScope.Resolve(name); ok {
			a.noteBindingRead(decl)
		} else if a.currentStructType != "" {
			if a.tryBindImplicitStructField(e, name) {
				return
			}
			hint := a.suggestName(name)
			a.record(&diagnostic.DiagnosticError{
				File:    e.Name.File,
				Line:    e.Name.Line,
				Col:     e.Name.Col,
				Message: fmt.Sprintf("undefined variable '%s'", name),
				Hint:    hint,
			})
		} else {
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
		if e.Operator == "===" || e.Operator == "!==" {
			alt := "=="
			if e.Operator == "!==" {
				alt = "!="
			}
			a.warn(fmt.Sprintf("%s:%d:%d: '%s' is deprecated; use '%s' instead (Koda has no loose equality)",
				e.Token.File, e.Token.Line, e.Token.Col, e.Operator, alt))
		}
		a.analyzeExpr(e.Left)
		a.analyzeExpr(e.Right)
	case *parser.LogicalExpr:
		a.analyzeExpr(e.Left)
		a.analyzeExpr(e.Right)
	case *parser.CallExpr:
		a.analyzeExpr(e.Function)
		if id, ok := e.Function.(*parser.IdentifierExpr); ok {
			if decl, ok2 := a.currentScope.Resolve(id.Name.Lexeme); ok2 {
				a.noteBindingRead(decl)
				if sd, ok3 := decl.(*parser.StructDecl); ok3 {
					a.checkStructConstructorCall(e, sd)
				}
			}
		}
		if ix, ok := e.Function.(*parser.IndexExpr); ok {
			a.noteStructMethodRead(ix)
		}
		for _, arg := range e.Arguments {
			a.analyzeExpr(arg)
		}
		a.recordParamStructFromCall(e)
		a.maybeCheckCallArity(e)
	case *parser.AssignExpr:
		a.analyzeExpr(e.Value)
		if ident, ok := e.Left.(*parser.IdentifierExpr); ok {
			name := ident.Name.Lexeme
			if decl, ok := a.currentScope.Resolve(name); ok {
				if ld, ok := decl.(*parser.LetDecl); ok && ld.IsConst {
					a.record(&diagnostic.DiagnosticError{
						File:    ident.Name.File,
						Line:    ident.Name.Line,
						Col:     ident.Name.Col,
						Message: fmt.Sprintf("cannot assign to const '%s'", name),
						Hint:    "use let for mutable bindings; const is immutable",
					})
					return
				}
			}
			if _, ok := a.currentScope.Resolve(name); !ok {
				if a.currentStructType != "" && a.tryBindImplicitStructField(ident, name) {
					return
				}
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
		if id, ok := e.Operand.(*parser.IdentifierExpr); ok {
			if _, ok := a.currentScope.Resolve(id.Name.Lexeme); !ok && a.currentStructType != "" {
				a.tryBindImplicitStructField(id, id.Name.Lexeme)
			}
		}
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
