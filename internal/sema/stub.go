package sema

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"koda/internal/diagnostic"
	"koda/internal/parser"
)

// ParamCellKey identifies a function parameter for escape-driven heap-cell lowering.
type ParamCellKey struct {
	Func uintptr
	Idx  int
}

// NewParamCellKey builds a stable key for a parameter on a *parser.FuncDecl or *parser.FuncExpr.
func NewParamCellKey(owner interface{}, idx int) ParamCellKey {
	return ParamCellKey{Func: reflect.ValueOf(owner).Pointer(), Idx: idx}
}

type NativeEmitContext struct {
	Bundle       *parser.ProgramBundle
	locals       map[parser.Expr]parser.Decl
	capturedVars map[parser.Decl]bool
	paramDecls   map[*parser.FuncDecl][]parser.Decl
	currentFn    *parser.FuncDecl

	FuncCaptures map[*parser.FuncDecl][]*parser.LetDecl
	ExprCaptures map[*parser.FuncExpr][]*parser.LetDecl

	EscapingDecls map[*parser.LetDecl]bool
	StackDecls    map[*parser.LetDecl]bool
	ParamIsCell   map[ParamCellKey]bool
	letOwner      map[*parser.LetDecl]interface{}

	// FreeVarsExpr / FreeVarsDecl: outer bindings referenced by nested functions (for codegen closure seeding).
	FreeVarsExpr map[*parser.FuncExpr][]string
	FreeVarsDecl map[*parser.FuncDecl][]string

	// FuncExprEnclosing maps each function expression to its enclosing function (FuncExpr or FuncDecl)
	// while walking, so captures can be propagated to intermediate closures (nested lambdas).
	FuncExprEnclosing map[*parser.FuncExpr]interface{}

	ShadowFuncDecl map[*parser.FuncDecl]*ShadowLayout
	ShadowFuncExpr map[*parser.FuncExpr]*ShadowLayout
	ShadowEntry    *ShadowLayout

	// StructFields maps struct type name -> ordered field names (for slot access).
	StructFields map[string][]string
	// VarStruct maps variable name -> struct type name when initialized from a struct literal.
	VarStruct map[string]string
	// VarEnum maps variable name -> enum type name when initialized from Enum.Member.
	VarEnum map[string]string
	// EnumOrdinal maps "EnumName.member" -> integer ordinal for constant folding.
	EnumOrdinal map[string]int
	// IndexExprStructSlot maps struct field index expressions to slot indices.
	IndexExprStructSlot map[*parser.IndexExpr]int
	// IndexExprEnumConst maps enum member index expressions to folded integer ordinals.
	IndexExprEnumConst map[*parser.IndexExpr]int64
	// EmitDebug requests LLVM debug metadata in codegen.
	EmitDebug bool
	// StructMethods maps struct type -> method name -> FuncDecl.
	StructMethods map[string]map[string]*parser.FuncDecl
	// NumericKinds maps stack LetDecl to inferred integer/float kind (A8).
	NumericKinds map[*parser.LetDecl]NumericKind
	// TypedLocals maps LetDecl to explicit integer type name (P1).
	TypedLocals map[*parser.LetDecl]string
}

// PrepareOptions configures sema passes during native bundle preparation.
type PrepareOptions struct {
	WarnUnused      bool
	BeginnerLint    bool
	EmitDebug       bool
	WarnUnreachable bool
}

func PrepareNativeBundle(bundle *parser.ProgramBundle) (*NativeEmitContext, error) {
	return PrepareNativeBundleWithOptions(bundle, nil)
}

func PrepareNativeBundleWithOptions(bundle *parser.ProgramBundle, opts *PrepareOptions) (*NativeEmitContext, error) {
	if err := parser.FlattenEntryIncludes(bundle); err != nil {
		return nil, err
	}
	parser.InjectNativeMathPrelude(bundle)

	entryPath := "<entry>"
	if p, err := parser.BundleEntryPath(bundle); err == nil {
		entryPath = p
	}
	var analyzer *Analyzer
	if bundle.Entry != nil {
		var aopts *AnalysisOptions
		if opts != nil && (opts.WarnUnused || opts.BeginnerLint || opts.WarnUnreachable) {
			aopts = &AnalysisOptions{
				WarnUnused:      opts.WarnUnused,
				BeginnerLint:    opts.BeginnerLint,
				WarnUnreachable: opts.WarnUnreachable,
			}
		}
		analyzer = NewAnalyzerWithOptions(aopts)
		if err := analyzer.Analyze(bundle.Entry); err != nil {
			var me *diagnostic.MultiError
			if errors.As(err, &me) && me != nil {
				me.Label = entryPath
			}
			patchEmptyDiagnosticFiles(err, entryPath)
			return nil, err
		}
		for _, w := range analyzer.Warnings() {
			_, _ = fmt.Fprintf(os.Stderr, "warning: %s\n", w)
		}
	}

	ctx := &NativeEmitContext{
		Bundle:         bundle,
		locals:         make(map[parser.Expr]parser.Decl),
		capturedVars:   make(map[parser.Decl]bool),
		paramDecls:     make(map[*parser.FuncDecl][]parser.Decl),
		FuncCaptures:   make(map[*parser.FuncDecl][]*parser.LetDecl),
		ExprCaptures:   make(map[*parser.FuncExpr][]*parser.LetDecl),
		EscapingDecls:  make(map[*parser.LetDecl]bool),
		StackDecls:     make(map[*parser.LetDecl]bool),
		ParamIsCell:    make(map[ParamCellKey]bool),
		letOwner:       make(map[*parser.LetDecl]interface{}),
		FreeVarsExpr:      make(map[*parser.FuncExpr][]string),
		FreeVarsDecl:      make(map[*parser.FuncDecl][]string),
		FuncExprEnclosing: make(map[*parser.FuncExpr]interface{}),
		ShadowFuncDecl:    make(map[*parser.FuncDecl]*ShadowLayout),
		ShadowFuncExpr: make(map[*parser.FuncExpr]*ShadowLayout),
		EnumOrdinal:    enumOrdinalMap(bundle.Entry),
		StructMethods:  make(map[string]map[string]*parser.FuncDecl),
		NumericKinds:   make(map[*parser.LetDecl]NumericKind),
		TypedLocals:    make(map[*parser.LetDecl]string),
	}
	if analyzer != nil {
		ctx.StructFields, ctx.StructMethods, ctx.VarStruct, ctx.VarEnum, ctx.IndexExprStructSlot, ctx.IndexExprEnumConst = analyzer.ExportForCodegen()
	}
	if bundle.Entry != nil {
		for _, d := range bundle.Entry.Declarations {
			collectTypedLocals(d, ctx.TypedLocals)
		}
	}
	prepareNativeAnalysis(ctx, bundle)
	ctx.NumericKinds = InferNumericKinds(bundle.Entry, ctx.EscapingDecls)
	if opts != nil {
		ctx.EmitDebug = opts.EmitDebug
	}
	prepareShadowLayouts(ctx, bundle)
	return ctx, nil
}

func collectTypedLocals(d parser.Decl, out map[*parser.LetDecl]string) {
	switch x := d.(type) {
	case *parser.LetDecl:
		if x.TypeAnnot != "" && isKnownTypeAnnotation(x.TypeAnnot) {
			out[x] = normalizeTypeName(x.TypeAnnot)
		}
	case *parser.FuncDecl:
		collectTypedLocalsInBlock(x.Body, out)
	case *parser.StructDecl:
		for _, m := range x.Methods {
			collectTypedLocalsInBlock(m.Body, out)
		}
	case parser.Stmt:
		collectTypedLocalsInStmt(x, out)
	}
}

func collectTypedLocalsInBlock(b *parser.BlockStmt, out map[*parser.LetDecl]string) {
	if b == nil {
		return
	}
	for _, inner := range b.Declarations {
		collectTypedLocals(inner, out)
	}
}

func collectTypedLocalsInStmt(s parser.Stmt, out map[*parser.LetDecl]string) {
	if s == nil {
		return
	}
	switch x := s.(type) {
	case *parser.BlockStmt:
		collectTypedLocalsInBlock(x, out)
	case *parser.IfStmt:
		collectTypedLocalsInStmt(x.Then, out)
		collectTypedLocalsInStmt(x.Else, out)
	case *parser.WhileStmt:
		collectTypedLocalsInStmt(x.Body, out)
	case *parser.DoWhileStmt:
		collectTypedLocalsInStmt(x.Body, out)
	case *parser.ForStmt:
		for _, ini := range x.Inits {
			collectTypedLocals(ini, out)
		}
		collectTypedLocalsInStmt(x.Body, out)
	case *parser.ForInStmt:
		collectTypedLocalsInStmt(x.Body, out)
	case *parser.ForOfStmt:
		collectTypedLocalsInStmt(x.Body, out)
	case *parser.SwitchStmt:
		for _, c := range x.Cases {
			for _, cd := range c.Body {
				collectTypedLocals(cd, out)
			}
		}
		for _, cd := range x.Default {
			collectTypedLocals(cd, out)
		}
	}
}

func ValidateNativeEmitSupport(ctx *NativeEmitContext) error {
	_ = ctx
	return nil
}

func patchEmptyDiagnosticFiles(err error, entryPath string) {
	var me *diagnostic.MultiError
	if errors.As(err, &me) && me != nil {
		for _, e := range me.List {
			patchOneDiagnosticFile(e, entryPath)
		}
		return
	}
	patchOneDiagnosticFile(err, entryPath)
}

func patchOneDiagnosticFile(err error, entryPath string) {
	var de *diagnostic.DiagnosticError
	if errors.As(err, &de) && de != nil && de.File == "" {
		de.File = entryPath
	}
}
