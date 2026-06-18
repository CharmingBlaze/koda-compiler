package sema

import (
	"fmt"
	"strings"

	"koda/internal/diagnostic"
	"koda/internal/lexer"
	"koda/internal/parser"
)

func (a *Analyzer) analyzeStructDecl(d *parser.StructDecl) {
	name := d.Name.Lexeme
	if existing, ok := a.currentScope.symbols[name]; ok {
		if prev, ok := existing.(*parser.StructDecl); ok && prev == d {
			return
		}
		a.record(&diagnostic.DiagnosticError{
			File:    d.Name.File,
			Line:    d.Name.Line,
			Col:     d.Name.Col,
			Message: fmt.Sprintf("duplicate declaration '%s' in the same scope", name),
		})
		return
	}
	a.currentScope.Define(name, d)
	fields := make([]string, len(d.Fields))
	defaults := make(map[string]parser.Expr)
	fieldTypes := make(map[string]string)
	for i, f := range d.Fields {
		fname := f.Name.Lexeme
		fields[i] = fname
		if f.TypeAnnot != "" {
			if !a.isValidStructFieldType(f.TypeAnnot) {
				a.record(&diagnostic.DiagnosticError{
					File:    f.Name.File,
					Line:    f.Name.Line,
					Col:     f.Name.Col,
					Message: fmt.Sprintf("unknown type '%s' on struct field '%s'", f.TypeAnnot, fname),
					Hint:    "use int, float, bool, string, or a struct type name",
				})
			} else {
				fieldTypes[fname] = normalizeTypeName(f.TypeAnnot)
			}
		}
		if f.Default != nil {
			defaults[fname] = f.Default
		}
		if f.Optional {
			if f.Default != nil {
				a.record(&diagnostic.DiagnosticError{
					File:    f.Name.File,
					Line:    f.Name.Line,
					Col:     f.Name.Col,
					Message: fmt.Sprintf("struct field '%s' cannot be both optional (?) and have a default value", f.Name.Lexeme),
					Hint:    "use either 'field?' or 'field = default', not both",
				})
			} else {
				// Optional (?) is sugar for "may omit → null", same as `field = null`.
				defaults[f.Name.Lexeme] = &parser.LiteralExpr{Token: f.Name, Value: nil}
			}
		}
	}
	a.structLayouts[name] = fields
	if len(defaults) > 0 {
		a.structFieldDefaults[name] = defaults
	}
	if len(fieldTypes) > 0 {
		a.structFieldTypes[name] = fieldTypes
	}
	methods := make(map[string]*parser.FuncDecl)
	for _, m := range d.Methods {
		mname := m.Name.Lexeme
		if _, dup := methods[mname]; dup {
			a.record(&diagnostic.DiagnosticError{
				File:    m.Name.File,
				Line:    m.Name.Line,
				Col:     m.Name.Col,
				Message: fmt.Sprintf("duplicate method '%s' on struct %s", mname, name),
			})
			continue
		}
		methods[mname] = m
		a.funcReads[m] = 0
		a.pendingStructMethods = append(a.pendingStructMethods, pendingStructMethod{
			structType: name,
			method:     m,
		})
	}
	if len(methods) > 0 {
		a.structMethods[name] = methods
	}
}

func (a *Analyzer) analyzeEnumDecl(d *parser.EnumDecl) {
	name := d.Name.Lexeme
	if _, ok := a.currentScope.symbols[name]; ok {
		a.record(&diagnostic.DiagnosticError{
			File:    d.Name.File,
			Line:    d.Name.Line,
			Col:     d.Name.Col,
			Message: fmt.Sprintf("duplicate declaration '%s' in the same scope", name),
		})
		return
	}
	a.currentScope.Define(name, d)
}

func (a *Analyzer) structFieldSlot(stName, field string) (int, bool) {
	fields, ok := a.structLayouts[stName]
	if !ok {
		return 0, false
	}
	if methods, ok := a.structMethods[stName]; ok {
		if _, isMethod := methods[field]; isMethod {
			return 0, false
		}
	}
	for i, f := range fields {
		if f == field {
			return i, true
		}
	}
	return 0, false
}

func (a *Analyzer) tryBindImplicitStructField(id *parser.IdentifierExpr, name string) bool {
	slot, ok := a.structFieldSlot(a.currentStructType, name)
	if !ok {
		return false
	}
	a.implicitStructField[id] = slot
	return true
}

func (a *Analyzer) structConstructorType(name string) string {
	decl, ok := a.currentScope.Resolve(name)
	if !ok {
		return ""
	}
	sd, ok := decl.(*parser.StructDecl)
	if !ok {
		return ""
	}
	methods, ok := a.structMethods[sd.Name.Lexeme]
	if !ok {
		return ""
	}
	if _, ok := methods["new"]; !ok {
		return ""
	}
	return sd.Name.Lexeme
}

func (a *Analyzer) checkStructConstructorCall(call *parser.CallExpr, sd *parser.StructDecl) {
	stName := sd.Name.Lexeme
	methods, ok := a.structMethods[stName]
	if !ok {
		a.record(&diagnostic.DiagnosticError{
			File:    call.Token.File,
			Line:    call.Token.Line,
			Col:     call.Token.Col,
			Message: fmt.Sprintf("struct %s has no constructor; add func new(...) inside the struct body", stName),
		})
		return
	}
	newFn, ok := methods["new"]
	if !ok {
		a.record(&diagnostic.DiagnosticError{
			File:    call.Token.File,
			Line:    call.Token.Line,
			Col:     call.Token.Col,
			Message: fmt.Sprintf("struct %s has no constructor; add func new(...) inside the struct body", stName),
		})
		return
	}
	a.checkCallArity(call, stName, 0, false, newFn.Params)
}

func (a *Analyzer) recordPlainObjectLiterals(varName string, oe *parser.ObjectExpr) {
	if oe == nil || len(oe.ComputedKeys) > 0 {
		return
	}
	fields := make(map[string]int64)
	for i, k := range oe.Keys {
		lit, ok := oe.Values[i].(*parser.LiteralExpr)
		if !ok {
			return
		}
		switch v := lit.Value.(type) {
		case int:
			fields[k.Lexeme] = int64(v)
		case float64:
			fields[k.Lexeme] = int64(v)
		default:
			return
		}
	}
	if len(fields) > 0 {
		a.constPlainObjectLiterals[varName] = fields
	}
}

// ExportConstPlainObjectLiterals returns module-level plain object lets with all-literal fields.
func (a *Analyzer) ExportConstPlainObjectLiterals() map[string]map[string]int64 {
	out := make(map[string]map[string]int64)
	for name, fields := range a.constPlainObjectLiterals {
		cp := make(map[string]int64)
		for k, v := range fields {
			cp[k] = v
		}
		out[name] = cp
	}
	return out
}

func (a *Analyzer) recordVarTypesFromInit(varName string, init parser.Expr) {
	if oe, ok := init.(*parser.ObjectExpr); ok {
		if oe.StructTag != nil {
			a.varStructType[varName] = oe.StructTag.Lexeme
		} else {
			a.varPlainObject[varName] = true
			a.recordPlainObjectLiterals(varName, oe)
		}
		return
	}
	if call, ok := init.(*parser.CallExpr); ok {
		if id, ok := call.Function.(*parser.IdentifierExpr); ok {
			if st := a.structConstructorType(id.Name.Lexeme); st != "" {
				a.varStructType[varName] = st
				return
			}
			if st := a.funcReturnStructType(id.Name.Lexeme); st != "" {
				a.varStructType[varName] = st
				return
			}
			if a.funcReturnsPlainObject(id.Name.Lexeme) {
				a.varPlainObject[varName] = true
				return
			}
		}
	}
	if ae, ok := init.(*parser.ArrayExpr); ok {
		a.varIsArray[varName] = true
		if len(ae.Elements) > 0 {
			if st := a.structTypeOfExpr(ae.Elements[0]); st != "" {
				a.varArrayElementStruct[varName] = st
			}
		}
		return
	}
	if ix, ok := init.(*parser.IndexExpr); ok {
		if id, ok2 := ix.Object.(*parser.IdentifierExpr); ok2 {
			if decl, ok3 := a.currentScope.Resolve(id.Name.Lexeme); ok3 {
				if _, ok4 := decl.(*parser.EnumDecl); ok4 {
					a.varEnumType[varName] = id.Name.Lexeme
				}
			}
		}
	}
}

func (a *Analyzer) validateStructLiteral(e *parser.ObjectExpr) {
	if e.StructTag == nil {
		return
	}
	stName := e.StructTag.Lexeme
	decl, ok := a.currentScope.Resolve(stName)
	if !ok {
		a.record(&diagnostic.DiagnosticError{
			File:    e.StructTag.File,
			Line:    e.StructTag.Line,
			Col:     e.StructTag.Col,
			Message: fmt.Sprintf("unknown struct type '%s'", stName),
		})
		return
	}
	sd, ok := decl.(*parser.StructDecl)
	if !ok {
		a.record(&diagnostic.DiagnosticError{
			File:    e.StructTag.File,
			Line:    e.StructTag.Line,
			Col:     e.StructTag.Col,
			Message: fmt.Sprintf("'%s' is not a struct type", stName),
		})
		return
	}
	want := make(map[string]bool, len(sd.Fields))
	for _, f := range sd.Fields {
		want[f.Name.Lexeme] = true
	}
	for _, k := range e.Keys {
		if !want[k.Lexeme] {
			a.record(&diagnostic.DiagnosticError{
				File:    k.File,
				Line:    k.Line,
				Col:     k.Col,
				Message: fmt.Sprintf("'%s' is not a field of struct %s", k.Lexeme, stName),
				Hint:    fmt.Sprintf("valid fields: %s", strings.Join(a.structLayouts[stName], ", ")),
			})
		}
	}
	a.checkStructLiteralFieldTypes(e, stName)
}

func (a *Analyzer) isValidStructFieldType(name string) bool {
	if isKnownTypeAnnotation(name) {
		return true
	}
	decl, ok := a.currentScope.Resolve(name)
	if !ok {
		return false
	}
	_, ok = decl.(*parser.StructDecl)
	return ok
}

func (a *Analyzer) checkStructLiteralFieldTypes(e *parser.ObjectExpr, stName string) {
	types, ok := a.structFieldTypes[stName]
	if !ok || len(types) == 0 {
		return
	}
	for i, k := range e.Keys {
		if i >= len(e.Values) {
			break
		}
		wantType, ok := types[k.Lexeme]
		if !ok || wantType == "" {
			continue
		}
		a.checkExprMatchesType(e.Values[i], wantType, k)
	}
}

func (a *Analyzer) checkExprMatchesType(expr parser.Expr, wantType string, at lexer.Token) {
	if expr == nil {
		return
	}
	wantType = normalizeTypeName(wantType)
	if a.exprMatchesType(expr, wantType) {
		return
	}
	got := a.exprTypeName(expr)
	if got == "" {
		return
	}
	if got == wantType {
		return
	}
	// float accepts integer literals; int accepts whole-number float literals (Koda number lexing).
	if wantType == "float" && got == "int" {
		return
	}
	if wantType == "int" && got == "float" {
		if lit, ok := expr.(*parser.LiteralExpr); ok {
			if f, ok := lit.Value.(float64); ok && f == float64(int64(f)) {
				return
			}
		}
	}
	a.record(&diagnostic.DiagnosticError{
		File:    at.File,
		Line:    at.Line,
		Col:     at.Col,
		Message: fmt.Sprintf("expected type '%s', got '%s'", wantType, got),
	})
}

func (a *Analyzer) exprMatchesType(expr parser.Expr, wantType string) bool {
	switch wantType {
	case "vector3":
		return a.isVector3Expr(expr)
	case "color":
		return a.isColorExpr(expr)
	default:
		return false
	}
}

func (a *Analyzer) isVector3Expr(expr parser.Expr) bool {
	if call, ok := expr.(*parser.CallExpr); ok {
		if id, ok := call.Function.(*parser.IdentifierExpr); ok && strings.EqualFold(id.Name.Lexeme, "vec3") {
			return len(call.Arguments) == 3
		}
	}
	oe, ok := expr.(*parser.ObjectExpr)
	if !ok || len(oe.ComputedKeys) > 0 {
		return false
	}
	if oe.StructTag != nil && !strings.EqualFold(oe.StructTag.Lexeme, "vector3") {
		return false
	}
	hasX, hasY, hasZ := false, false, false
	for _, k := range oe.Keys {
		switch strings.ToLower(k.Lexeme) {
		case "x":
			hasX = true
		case "y":
			hasY = true
		case "z":
			hasZ = true
		}
	}
	return hasX && hasY && hasZ
}

func (a *Analyzer) isColorExpr(expr parser.Expr) bool {
	if call, ok := expr.(*parser.CallExpr); ok {
		if id, ok := call.Function.(*parser.IdentifierExpr); ok {
			switch strings.ToLower(id.Name.Lexeme) {
			case "rgb", "rgba", "color":
				return true
			}
		}
	}
	if _, ok := expr.(*parser.LiteralExpr); ok {
		// hex colors lex as numbers at runtime; sema may not see them here.
		return false
	}
	oe, ok := expr.(*parser.ObjectExpr)
	if !ok || len(oe.ComputedKeys) > 0 {
		return false
	}
	if oe.StructTag != nil && !strings.EqualFold(oe.StructTag.Lexeme, "color") {
		return false
	}
	hasR, hasG, hasB := false, false, false
	for _, k := range oe.Keys {
		switch strings.ToLower(k.Lexeme) {
		case "r":
			hasR = true
		case "g":
			hasG = true
		case "b":
			hasB = true
		}
	}
	return hasR && hasG && hasB
}

func (a *Analyzer) exprTypeName(expr parser.Expr) string {
	switch x := expr.(type) {
	case *parser.LiteralExpr:
		switch x.Value.(type) {
		case int:
			return "int"
		case float64:
			return "float"
		case bool:
			return "bool"
		case string:
			return "string"
		default:
			if x.Value == nil {
				return "null"
			}
			return ""
		}
	case *parser.ObjectExpr:
		if x.StructTag != nil {
			return normalizeTypeName(x.StructTag.Lexeme)
		}
		return "object"
	case *parser.IdentifierExpr:
		if decl, ok := a.currentScope.Resolve(x.Name.Lexeme); ok {
			if ld, ok := decl.(*parser.LetDecl); ok && ld.TypeAnnot != "" {
				return normalizeTypeName(ld.TypeAnnot)
			}
			if st, ok := a.varStructType[x.Name.Lexeme]; ok {
				return normalizeTypeName(st)
			}
		}
		return ""
	default:
		return ""
	}
}

func (a *Analyzer) checkStructFieldAccess(e *parser.IndexExpr) {
	if te, ok := e.Object.(*parser.ThisExpr); ok && a.currentStructType != "" {
		_ = te
		a.bindReceiverFieldAccess(e, a.currentStructType)
		return
	}
	if id, ok := e.Object.(*parser.IdentifierExpr); ok && a.currentStructType != "" && isSelfReceiverParam(id.Name.Lexeme) {
		a.bindReceiverFieldAccess(e, a.currentStructType)
		return
	}
	lit, ok := e.Index.(*parser.LiteralExpr)
	if !ok {
		return
	}
	field, ok := lit.Value.(string)
	if !ok {
		return
	}
	stName, ok := a.structTypeNameForObject(e.Object)
	if !ok {
		return
	}
	if methods, ok := a.structMethods[stName]; ok {
		if _, isMethod := methods[field]; isMethod {
			return
		}
	}
	fields, ok := a.structLayouts[stName]
	if !ok {
		return
	}
	for i, f := range fields {
		if f == field {
			a.indexExprStructSlot[e] = i
			return
		}
	}
	a.record(&diagnostic.DiagnosticError{
		File:    lit.Token.File,
		Line:    lit.Token.Line,
		Col:     lit.Token.Col,
		Message: fmt.Sprintf("'%s' is not a field of struct %s", field, stName),
		Hint:    a.structFieldHint(field, fields),
	})
}

func (a *Analyzer) structTypeNameForObject(obj parser.Expr) (string, bool) {
	switch x := obj.(type) {
	case *parser.IdentifierExpr:
		name := x.Name.Lexeme
		if st, ok := a.varStructType[name]; ok {
			return st, true
		}
		if a.activeParamStruct != nil {
			if st, ok := a.activeParamStruct[name]; ok {
				return st, true
			}
		}
	case *parser.IndexExpr:
		if id, ok := x.Object.(*parser.IdentifierExpr); ok {
			if st, ok := a.varArrayElementStruct[id.Name.Lexeme]; ok {
				return st, true
			}
		}
	}
	return "", false
}

func (a *Analyzer) structFieldHint(field string, fields []string) string {
	if hint := suggestFromList(field, fields); hint != "" {
		return hint
	}
	return fmt.Sprintf("fields: %s", strings.Join(fields, ", "))
}

func suggestFromList(name string, candidates []string) string {
	const maxDist = 2 // distance 1–2 only; 3+ produces misleading hints
	ln := strings.ToLower(name)
	best, bestDist := "", maxDist+1
	for _, candidate := range candidates {
		if candidate == name {
			continue
		}
		d := levenshtein(ln, strings.ToLower(candidate))
		if d >= 1 && d <= maxDist && d < bestDist {
			best, bestDist = candidate, d
		}
	}
	if best != "" {
		return fmt.Sprintf("did you mean '%s'?", best)
	}
	return ""
}

func (a *Analyzer) checkEnumMemberAccess(e *parser.IndexExpr) {
	idObj, ok := e.Object.(*parser.IdentifierExpr)
	if !ok {
		return
	}
	lit, ok := e.Index.(*parser.LiteralExpr)
	if !ok {
		return
	}
	mem, ok := lit.Value.(string)
	if !ok {
		return
	}
	decl, ok := a.currentScope.Resolve(idObj.Name.Lexeme)
	if !ok {
		return
	}
	ed, ok := decl.(*parser.EnumDecl)
	if !ok {
		return
	}
	var memList []string
	for i, m := range ed.Members {
		memList = append(memList, m.Lexeme)
		if m.Lexeme == mem {
			a.indexExprEnumConst[e] = int64(i)
			return
		}
	}
	a.record(&diagnostic.DiagnosticError{
		File:    lit.Token.File,
		Line:    lit.Token.Line,
		Col:     lit.Token.Col,
		Message: fmt.Sprintf("'%s' is not a member of %s (members: %s)", mem, ed.Name.Lexeme, strings.Join(memList, ", ")),
	})
}

func (a *Analyzer) enumCaseOrdinal(val parser.Expr, enumName string) (int, bool) {
	ix, ok := val.(*parser.IndexExpr)
	if !ok {
		return 0, false
	}
	id, ok := ix.Object.(*parser.IdentifierExpr)
	if !ok || id.Name.Lexeme != enumName {
		return 0, false
	}
	lit, ok := ix.Index.(*parser.LiteralExpr)
	if !ok {
		return 0, false
	}
	mem, ok := lit.Value.(string)
	if !ok {
		return 0, false
	}
	decl, ok := a.currentScope.Resolve(enumName)
	if !ok {
		return 0, false
	}
	ed, ok := decl.(*parser.EnumDecl)
	if !ok {
		return 0, false
	}
	for i, m := range ed.Members {
		if m.Lexeme == mem {
			return i, true
		}
	}
	return 0, false
}

// ExportVarArrayAndObject returns array vs plain-object variable bindings for codegen method dispatch.
func (a *Analyzer) ExportVarArrayAndObject() (map[string]bool, map[string]bool) {
	arr := make(map[string]bool)
	for k, v := range a.varIsArray {
		if v {
			arr[k] = true
		}
	}
	obj := make(map[string]bool)
	for k, v := range a.varPlainObject {
		if v {
			obj[k] = true
		}
	}
	return arr, obj
}

// ExportForCodegen copies struct/enum binding maps for LLVM emission.
func (a *Analyzer) ExportForCodegen() (
	structFields map[string][]string,
	structMethods map[string]map[string]*parser.FuncDecl,
	varStruct map[string]string,
	varEnum map[string]string,
	indexStruct map[*parser.IndexExpr]int,
	indexEnum map[*parser.IndexExpr]int64,
	fieldDefaults map[string]map[string]parser.Expr,
	implicitField map[*parser.IdentifierExpr]int,
	forOfVarStruct map[string]map[string]string,
) {
	structFields = make(map[string][]string)
	for k, v := range a.structLayouts {
		cp := make([]string, len(v))
		copy(cp, v)
		structFields[k] = cp
	}
	structMethods = make(map[string]map[string]*parser.FuncDecl)
	for st, ms := range a.structMethods {
		cp := make(map[string]*parser.FuncDecl)
		for k, v := range ms {
			cp[k] = v
		}
		structMethods[st] = cp
	}
	fieldDefaults = make(map[string]map[string]parser.Expr)
	for st, defs := range a.structFieldDefaults {
		cp := make(map[string]parser.Expr)
		for k, v := range defs {
			cp[k] = v
		}
		fieldDefaults[st] = cp
	}
	implicitField = make(map[*parser.IdentifierExpr]int)
	for k, v := range a.implicitStructField {
		implicitField[k] = v
	}
	varStruct = make(map[string]string)
	for k, v := range a.varStructType {
		varStruct[k] = v
	}
	varEnum = make(map[string]string)
	for k, v := range a.varEnumType {
		varEnum[k] = v
	}
	indexStruct = make(map[*parser.IndexExpr]int)
	for k, v := range a.indexExprStructSlot {
		indexStruct[k] = v
	}
	indexEnum = make(map[*parser.IndexExpr]int64)
	for k, v := range a.indexExprEnumConst {
		indexEnum[k] = v
	}
	forOfVarStruct = make(map[string]map[string]string)
	for fn, vars := range a.forOfVarStruct {
		cp := make(map[string]string, len(vars))
		for k, v := range vars {
			cp[k] = v
		}
		forOfVarStruct[fn] = cp
	}
	return structFields, structMethods, varStruct, varEnum, indexStruct, indexEnum, fieldDefaults, implicitField, forOfVarStruct
}

// ExportStructFieldTypes copies struct field type annotations for codegen fast paths.
func (a *Analyzer) ExportStructFieldTypes() map[string]map[string]string {
	out := make(map[string]map[string]string)
	for st, fields := range a.structFieldTypes {
		cp := make(map[string]string, len(fields))
		for k, v := range fields {
			cp[k] = v
		}
		out[st] = cp
	}
	return out
}

func enumOrdinalMap(entry *parser.Program) map[string]int {
	out := make(map[string]int)
	if entry == nil {
		return out
	}
	for _, d := range entry.Declarations {
		ed, ok := d.(*parser.EnumDecl)
		if !ok {
			continue
		}
		for i, m := range ed.Members {
			key := ed.Name.Lexeme + "." + m.Lexeme
			out[key] = i
		}
	}
	return out
}

func (a *Analyzer) checkSwitchEnumExhaustive(sw *parser.SwitchStmt) {
	id, ok := sw.Subject.(*parser.IdentifierExpr)
	if !ok {
		return
	}
	enumName, ok := a.varEnumType[id.Name.Lexeme]
	if !ok || len(sw.Default) > 0 {
		return
	}
	decl, ok := a.currentScope.Resolve(enumName)
	if !ok {
		return
	}
	ed, ok := decl.(*parser.EnumDecl)
	if !ok {
		return
	}
	covered := make([]bool, len(ed.Members))
	for _, c := range sw.Cases {
		ord, ok := a.enumCaseOrdinal(c.Value, enumName)
		if !ok {
			return
		}
		if ord >= 0 && ord < len(covered) {
			covered[ord] = true
		}
	}
	var missing []string
	for i, m := range ed.Members {
		if i < len(covered) && !covered[i] {
			missing = append(missing, m.Lexeme)
		}
	}
	if len(missing) > 0 {
		a.warn(fmt.Sprintf("%s:%d:%d: switch on %s is not exhaustive; missing cases: %s (add cases or default)", sw.Token.File, sw.Token.Line, sw.Token.Col, enumName, strings.Join(missing, ", ")))
	}
}
