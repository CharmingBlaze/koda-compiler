package sema

import "strings"

// IntegerTypeNames are supported explicit integer type annotations (P1 / advanced).
var IntegerTypeNames = map[string]bool{
	"i8": true, "i16": true, "i32": true, "i64": true,
	"u8": true, "u16": true, "u32": true, "u64": true,
	"byte": true,
}

// BeginnerTypeNames are optional type annotations for learners and FFI prep.
var BeginnerTypeNames = map[string]bool{
	"int": true, "float": true, "bool": true, "string": true, "byte": true,
}

func normalizeTypeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func isKnownTypeAnnotation(name string) bool {
	n := normalizeTypeName(name)
	return IntegerTypeNames[n] || BeginnerTypeNames[n]
}

func isIntegerTypeName(name string) bool {
	n := normalizeTypeName(name)
	switch n {
	case "int", "byte":
		return true
	default:
		return IntegerTypeNames[n]
	}
}
