package sema

import "koda/internal/parser"

// nativeBuiltinSentinel marks names reserved by the native runtime (see codegen.registerBuiltinFuncs).
var nativeBuiltinSentinel parser.Decl = (*parser.IncludeDecl)(nil)

// nativeBuiltinNames must stay aligned with internal/codegen/builtin_register.go pairs.
var nativeBuiltinNames = []string{
	"this",
	"ok", "err", "panic", "assert", "expect", "warn",
	"deltatime", "print", "len", "keys", "type", "typeof", "matches",
	"time", "clock", "timestamp", "programtime", "sleep",
	"random", "randomint", "randomchoice", "randomseed",
	"lerp", "clamp", "distance", "anglebetween", "map",
	"pi", "e",
	"sin", "cos", "tan", "asin", "acos", "atan", "atan2",
	"pow", "exp", "log", "log10", "log2",
	"floor", "ceil", "round", "trunc", "sign", "min", "max",
	"smoothstep", "distancesq", "normalize",
	"hypot", "fmod", "degrees", "radians", "wrap", "approach", "smoothdamp",
	"isnumber", "isstring", "isbool", "isnull", "isarray", "isobject", "isfunction",
	"bool", "format",
	"readfile", "writefile", "appendfile", "fileexists", "deletefile",
	"isfile", "isdir", "filesize", "listdir",
	"trace", "parsejson", "tojson",
	"abs", "sqrt", "cbrt", "number", "string",
	"gc", "gccollect", "gcdisable", "gcenable", "gcframestep", "gcstats",
	"internclear",
	"arena", "arenareset", "arenaallocarray", "arenaallocstruct",
	"arraypush", "arraypop", "arrayslice", "arraysort", "arrayreverse",
	"arrayincludes", "arrayindexof", "arrayconcat", "assetpath", "args", "env", "rgb", "rgba",
	"vec2", "vec3", "rect", "box", "color", "structfield",
	"replace", "replaceall",
}

func seedGlobalBuiltins(s *Scope) {
	for _, name := range nativeBuiltinNames {
		s.Define(name, nativeBuiltinSentinel)
	}
}
