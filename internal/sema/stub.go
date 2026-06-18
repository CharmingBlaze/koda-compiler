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
	// StructFieldDefaults maps struct type -> field name -> default expression.
	StructFieldDefaults map[string]map[string]parser.Expr
	// StructFieldTypes maps struct type -> field name -> normalized type annotation.
	StructFieldTypes map[string]map[string]string
	// ImplicitStructField maps bare field identifiers in struct methods to slot indices.
	ImplicitStructField map[*parser.IdentifierExpr]int
	// FuncForOfVarStruct maps function name -> for-of loop variable -> element struct type.
	FuncForOfVarStruct map[string]map[string]string
	// NumericKinds maps stack LetDecl to inferred integer/float kind (A8).
	NumericKinds map[*parser.LetDecl]NumericKind
	// ConstPlainObjectLiterals maps module let name -> field -> integer literal (e.g. color components).
	ConstPlainObjectLiterals map[string]map[string]int64
	// VarIsArray maps variable name -> bound to an array literal (for method dispatch).
	VarIsArray map[string]bool
	// PlainObjectVars maps variable name -> plain object literal init (not array/struct).
	PlainObjectVars map[string]bool
	// TypedLocals maps LetDecl to explicit integer type name (P1).
	TypedLocals map[*parser.LetDecl]string
	// TypedParams maps function parameter to explicit type name.
	TypedParams map[ParamCellKey]string
	// ParamNumericKinds maps untyped parameters to inferred int/float kind (fast path without annotations).
	ParamNumericKinds map[ParamCellKey]NumericKind
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
	parser.InjectColorPrelude(bundle)
	parser.InjectRaylibPrelude(bundle)

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
		StructFieldDefaults: make(map[string]map[string]parser.Expr),
		ImplicitStructField: make(map[*parser.IdentifierExpr]int),
		NumericKinds:   make(map[*parser.LetDecl]NumericKind),
		TypedLocals:    make(map[*parser.LetDecl]string),
		TypedParams:       make(map[ParamCellKey]string),
		ParamNumericKinds: make(map[ParamCellKey]NumericKind),
	}
	if analyzer != nil {
		ctx.StructFields, ctx.StructMethods, ctx.VarStruct, ctx.VarEnum, ctx.IndexExprStructSlot, ctx.IndexExprEnumConst, ctx.StructFieldDefaults, ctx.ImplicitStructField, ctx.FuncForOfVarStruct = analyzer.ExportForCodegen()
		ctx.StructFieldTypes = analyzer.ExportStructFieldTypes()
		ctx.ConstPlainObjectLiterals = analyzer.ExportConstPlainObjectLiterals()
		ctx.VarIsArray, ctx.PlainObjectVars = analyzer.ExportVarArrayAndObject()
	}
	if bundle.Entry != nil {
		for _, d := range bundle.Entry.Declarations {
			collectTypedLocals(d, ctx.TypedLocals)
			collectTypedParams(d, ctx.TypedParams)
		}
	}
	prepareNativeAnalysis(ctx, bundle)
	ctx.NumericKinds = InferNumericKinds(bundle.Entry, ctx.EscapingDecls)
	ctx.ParamNumericKinds = InferParamNumericKinds(bundle.Entry)
	InferParamKindsFromCallSites(bundle.Entry, ctx.NumericKinds, ctx.ParamNumericKinds)
	collectTypedParamsFromProgram(bundle.Entry, ctx.TypedParams)
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

func collectTypedParams(d parser.Decl, out map[ParamCellKey]string) {
	switch x := d.(type) {
	case *parser.FuncDecl:
		recordFuncParamTypes(x, out)
		collectTypedParamsInBlock(x.Body, out)
	case *parser.StructDecl:
		for _, m := range x.Methods {
			recordFuncParamTypes(m, out)
			collectTypedParamsInBlock(m.Body, out)
		}
	case parser.Stmt:
		collectTypedParamsInStmt(x, out)
	}
}

func recordFuncParamTypes(fd *parser.FuncDecl, out map[ParamCellKey]string) {
	recordOwnerParamTypes(fd, fd.Params, out)
}

func recordOwnerParamTypes(owner interface{}, params []parser.Param, out map[ParamCellKey]string) {
	if owner == nil {
		return
	}
	for i, p := range params {
		if p.TypeAnnot != "" && isKnownTypeAnnotation(p.TypeAnnot) {
			out[NewParamCellKey(owner, i)] = normalizeTypeName(p.TypeAnnot)
		}
	}
}

func collectTypedParamsFromProgram(prog *parser.Program, out map[ParamCellKey]string) {
	if prog == nil {
		return
	}
	var scanExpr func(parser.Expr)
	scanExpr = func(e parser.Expr) {
		if e == nil {
			return
		}
		switch x := e.(type) {
		case *parser.FuncExpr:
			recordOwnerParamTypes(x, x.Params, out)
			scanParamsInBlock(x.Body, scanExpr)
		case *parser.CallExpr:
			scanExpr(x.Function)
			for _, a := range x.Arguments {
				scanExpr(a)
			}
		case *parser.InfixExpr:
			scanExpr(x.Left)
			scanExpr(x.Right)
		case *parser.PrefixExpr:
			scanExpr(x.Right)
		case *parser.AssignExpr:
			scanExpr(x.Left)
			scanExpr(x.Value)
		case *parser.LogicalExpr:
			scanExpr(x.Left)
			scanExpr(x.Right)
		case *parser.IndexExpr:
			scanExpr(x.Object)
			scanExpr(x.Index)
		case *parser.GroupingExpr:
			scanExpr(x.Expr)
		case *parser.ArrayExpr:
			for _, el := range x.Elements {
				scanExpr(el)
			}
		case *parser.IfExpr:
			scanExpr(x.Condition)
			scanExpr(x.Then)
			if x.Else != nil {
				scanExpr(x.Else)
			}
		}
	}
	var walkDecl func(parser.Decl)
	walkDecl = func(d parser.Decl) {
		switch x := d.(type) {
		case *parser.LetDecl:
			scanExpr(x.Init)
		case *parser.FuncDecl:
			scanParamsInBlock(x.Body, scanExpr)
		case *parser.StructDecl:
			for _, m := range x.Methods {
				scanParamsInBlock(m.Body, scanExpr)
			}
		case *parser.BlockStmt:
			for _, inner := range x.Declarations {
				walkDecl(inner)
			}
		case parser.Stmt:
			walkParamsInStmtForScan(x, walkDecl, scanExpr)
		}
	}
	for _, d := range prog.Declarations {
		walkDecl(d)
	}
}

func collectTypedParamsInBlock(b *parser.BlockStmt, out map[ParamCellKey]string) {
	if b == nil {
		return
	}
	for _, inner := range b.Declarations {
		collectTypedParams(inner, out)
	}
}

func collectTypedParamsInStmt(s parser.Stmt, out map[ParamCellKey]string) {
	if s == nil {
		return
	}
	switch x := s.(type) {
	case *parser.BlockStmt:
		collectTypedParamsInBlock(x, out)
	case *parser.IfStmt:
		collectTypedParamsInStmt(x.Then, out)
		collectTypedParamsInStmt(x.Else, out)
	case *parser.WhileStmt:
		collectTypedParamsInStmt(x.Body, out)
	case *parser.LoopStmt:
		collectTypedParamsInStmt(x.Body, out)
	case *parser.DoWhileStmt:
		collectTypedParamsInStmt(x.Body, out)
	case *parser.ForStmt:
		collectTypedParamsInStmt(x.Body, out)
	case *parser.ForInStmt:
		collectTypedParamsInStmt(x.Body, out)
	case *parser.ForOfStmt:
		collectTypedParamsInStmt(x.Body, out)
	case *parser.SwitchStmt:
		for _, c := range x.Cases {
			for _, cd := range c.Body {
				collectTypedParams(cd, out)
			}
		}
		for _, cd := range x.Default {
			collectTypedParams(cd, out)
		}
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
	case *parser.LoopStmt:
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
