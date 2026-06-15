package sema

import (
	"fmt"
	"strings"

	"koda/internal/diagnostic"
	"koda/internal/parser"
)

func (a *Analyzer) analyzeStructDecl(d *parser.StructDecl) {
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
	fields := make([]string, len(d.Fields))
	defaults := make(map[string]parser.Expr)
	for i, f := range d.Fields {
		fields[i] = f.Name.Lexeme
		if f.Default != nil {
			defaults[f.Name.Lexeme] = f.Default
		}
	}
	a.structLayouts[name] = fields
	if len(defaults) > 0 {
		a.structFieldDefaults[name] = defaults
	}
	for _, f := range d.Fields {
		if f.Default != nil {
			a.analyzeExpr(f.Default)
		}
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
		prev := a.currentStructType
		a.currentStructType = name
		a.analyzeFuncDecl(m)
		a.currentStructType = prev
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

func (a *Analyzer) recordVarTypesFromInit(varName string, init parser.Expr) {
	if oe, ok := init.(*parser.ObjectExpr); ok && oe.StructTag != nil {
		a.varStructType[varName] = oe.StructTag.Lexeme
		return
	}
	if call, ok := init.(*parser.CallExpr); ok {
		if id, ok := call.Function.(*parser.IdentifierExpr); ok {
			if st := a.structConstructorType(id.Name.Lexeme); st != "" {
				a.varStructType[varName] = st
				return
			}
		}
	}
	if ae, ok := init.(*parser.ArrayExpr); ok && len(ae.Elements) > 0 {
		if st := a.structTypeOfExpr(ae.Elements[0]); st != "" {
			a.varArrayElementStruct[varName] = st
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
}

func (a *Analyzer) checkStructFieldAccess(e *parser.IndexExpr) {
	if te, ok := e.Object.(*parser.ThisExpr); ok && a.currentStructType != "" {
		lit, ok := e.Index.(*parser.LiteralExpr)
		if !ok {
			return
		}
		field, ok := lit.Value.(string)
		if !ok {
			return
		}
		stName := a.currentStructType
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
			Hint:    fmt.Sprintf("fields: %s", strings.Join(fields, ", ")),
		})
		_ = te
		return
	}
	idObj, ok := e.Object.(*parser.IdentifierExpr)
	if !ok {
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
	stName, ok := a.varStructType[idObj.Name.Lexeme]
	if !ok && a.activeParamStruct != nil {
		stName, ok = a.activeParamStruct[idObj.Name.Lexeme]
	}
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
		Hint:    fmt.Sprintf("fields: %s", strings.Join(fields, ", ")),
	})
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
	return structFields, structMethods, varStruct, varEnum, indexStruct, indexEnum, fieldDefaults, implicitField
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
