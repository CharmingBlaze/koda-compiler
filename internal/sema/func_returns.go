package sema

import "koda/internal/parser"

func (a *Analyzer) recordFuncReturnStruct(fd *parser.FuncDecl) {
	if fd.Body == nil {
		return
	}
	for _, decl := range fd.Body.Declarations {
		rs, ok := decl.(*parser.ReturnStmt)
		if !ok || rs.Value == nil {
			continue
		}
		if st := a.structTypeOfExpr(rs.Value); st != "" {
			if a.funcReturnStruct == nil {
				a.funcReturnStruct = make(map[*parser.FuncDecl]string)
			}
			a.funcReturnStruct[fd] = st
		} else if oe, ok := rs.Value.(*parser.ObjectExpr); ok && oe.StructTag == nil {
			if a.funcReturnPlain == nil {
				a.funcReturnPlain = make(map[*parser.FuncDecl]bool)
			}
			a.funcReturnPlain[fd] = true
		}
		return
	}
}

func (a *Analyzer) funcReturnStructType(name string) string {
	decl, ok := a.currentScope.Resolve(name)
	if !ok {
		return ""
	}
	fd, ok := decl.(*parser.FuncDecl)
	if !ok {
		return ""
	}
	return a.funcReturnStruct[fd]
}

func (a *Analyzer) funcReturnsPlainObject(name string) bool {
	decl, ok := a.currentScope.Resolve(name)
	if !ok {
		return false
	}
	fd, ok := decl.(*parser.FuncDecl)
	if !ok {
		return false
	}
	return a.funcReturnPlain[fd]
}
