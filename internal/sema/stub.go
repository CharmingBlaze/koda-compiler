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
}

func PrepareNativeBundle(bundle *parser.ProgramBundle) (*NativeEmitContext, error) {
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
		analyzer = NewAnalyzer()
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
	}
	if analyzer != nil {
		ctx.StructFields, ctx.VarStruct, ctx.VarEnum, ctx.IndexExprStructSlot, ctx.IndexExprEnumConst = analyzer.ExportForCodegen()
	}
	prepareNativeAnalysis(ctx, bundle)
	prepareShadowLayouts(ctx, bundle)
	return ctx, nil
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
