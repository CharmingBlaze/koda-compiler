package codegen

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"koda/internal/parser"
)

// emitImport builds a module export object for import("@module") / import("path.koda").
func (g *Generator) emitImport(e *parser.ImportExpr) (value.Value, error) {
	rel := strings.Trim(e.Path.Lexeme, `"`)
	abs, err := parser.ResolveImportPath(g.sourcePath, rel)
	if err != nil {
		return nil, err
	}
	if slot, ok := g.importSlots[abs]; ok {
		return g.block.NewLoad(types.I64, slot), nil
	}
	slot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(slot)
	g.importSlots[abs] = slot

	obj, err := g.buildImportObject(abs)
	if err != nil {
		return nil, err
	}
	g.block.NewStore(g.emitAsKodaI64(obj), slot)
	return g.block.NewLoad(types.I64, slot), nil
}

func (g *Generator) buildImportObject(absPath string) (value.Value, error) {
	switch stdlibModuleKind(absPath) {
	case "math":
		return g.emitStdlibMathObject()
	case "json":
		return g.emitStdlibJSONObject()
	case "io":
		return g.emitStdlibIOObject()
	}
	if g.ctx == nil || g.ctx.Bundle == nil {
		return nil, fmt.Errorf("native codegen: import %s: no program bundle", absPath)
	}
	prog := g.ctx.Bundle.Modules[absPath]
	if prog == nil {
		prog = findBundleModule(g.ctx.Bundle, absPath)
		if prog == nil {
			return nil, fmt.Errorf("native codegen: import module not loaded: %s", absPath)
		}
		absPath = bundleModulePath(g.ctx.Bundle, prog)
	}
	return g.emitModuleExportObject(absPath, prog)
}

func findBundleModule(bundle *parser.ProgramBundle, resolved string) *parser.Program {
	if p := bundle.Modules[resolved]; p != nil {
		return p
	}
	base := filepath.Base(resolved)
	for path, prog := range bundle.Modules {
		if prog == bundle.Entry {
			continue
		}
		if filepath.Base(path) == base {
			return prog
		}
	}
	return nil
}

func bundleModulePath(bundle *parser.ProgramBundle, prog *parser.Program) string {
	for path, p := range bundle.Modules {
		if p == prog {
			return path
		}
	}
	return ""
}

func stdlibModuleKind(absPath string) string {
	switch strings.ToLower(filepath.Base(absPath)) {
	case "math.koda":
		return "math"
	case "json.koda":
		return "json"
	case "io.koda":
		return "io"
	default:
		return ""
	}
}

func moduleFuncLLVMName(modulePath, funcName string) string {
	sum := sha256.Sum256([]byte(modulePath + "\x00" + funcName))
	return "koda_mod_" + hex.EncodeToString(sum[:6]) + "_" + sanitizeLLVMIdent(funcName)
}

func sanitizeLLVMIdent(s string) string {
	if s == "" {
		return "fn"
	}
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteByte('_')
		}
	}
	out := b.String()
	if out == "" {
		return "fn"
	}
	if out[0] >= '0' && out[0] <= '9' {
		return "f_" + out
	}
	return out
}

func (g *Generator) emitBoxNativeFunc(fn *ir.Func) value.Value {
	if fn == nil {
		return constant.NewInt(types.I64, llvmNilTagged)
	}
	raw := g.block.NewPtrToInt(fn, types.I64)
	return g.block.NewOr(raw, constant.NewInt(types.I64, 1))
}

func (g *Generator) emitRuntimeConstValue(fn *ir.Func) value.Value {
	if fn == nil {
		return constant.NewInt(types.I64, llvmNilTagged)
	}
	return g.emitArgvRuntime(fn, nil)
}

func (g *Generator) emitObjectWithProps(props map[string]value.Value) (value.Value, error) {
	n := len(props)
	if n < 1 {
		n = 1
	}
	obj := g.block.NewCall(g.runtimeAllocObj, constant.NewInt(types.I32, int64(n)))
	objSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(objSlot)
	g.block.NewStore(obj, objSlot)
	for key, val := range props {
		keyVal := g.emitStringLiteral(key)
		objLive := g.block.NewLoad(types.I64, objSlot)
		g.block.NewCall(g.runtimeObjSet, objLive, g.emitAsKodaI64(keyVal), g.emitAsKodaI64(val))
	}
	return g.block.NewLoad(types.I64, objSlot), nil
}

func (g *Generator) emitStdlibMathObject() (value.Value, error) {
	props := make(map[string]value.Value)

	// Constants (runtime argv thunks or literals matching stdlib/math.koda).
	props["pi"] = g.emitRuntimeConstValue(g.runtimePI)
	props["e"] = g.emitRuntimeConstValue(g.runtimeE)
	props["tau"] = g.block.NewCall(g.runtimeBoxNumber, constant.NewFloat(types.Double, 6.283185307179586))
	props["phi"] = g.block.NewCall(g.runtimeBoxNumber, constant.NewFloat(types.Double, 1.618033988749895))

	for name, fn := range map[string]*ir.Func{
		"sin": g.runtimeSin, "cos": g.runtimeCos, "tan": g.runtimeTan,
		"asin": g.runtimeAsin, "acos": g.runtimeAcos, "atan": g.runtimeAtan, "atan2": g.runtimeAtan2,
		"hypot": g.runtimeHypot, "sqrt": g.runtimeSqrt, "pow": g.runtimePow,
		"exp": g.runtimeExp, "log": g.runtimeLog, "log10": g.runtimeLog10, "log2": g.runtimeLog2,
		"floor": g.runtimeFloor, "ceil": g.runtimeCeil, "round": g.runtimeRound,
		"abs": g.runtimeAbs, "min": g.runtimeMin, "max": g.runtimeMax,
		"clamp": g.runtimeClamp, "sign": g.runtimeSign, "lerp": g.runtimeLerp,
		"degrees": g.runtimeDegrees, "radians": g.runtimeRadians,
		"wrap": g.runtimeWrap, "approach": g.runtimeApproach,
		"random": g.runtimeRandom, "randomint": g.runtimeRandomInt,
		"randomchoice": g.runtimeRandomChoice, "randomseed": g.runtimeRandomSeed,
		"cbrt": g.runtimeCbrt,
	} {
		if fn != nil {
			props[name] = g.emitBoxNativeFunc(fn)
		}
	}
	props["deg"] = props["degrees"]
	props["rad"] = props["radians"]

	return g.emitObjectWithProps(props)
}

func (g *Generator) emitStdlibJSONObject() (value.Value, error) {
	props := make(map[string]value.Value)
	if g.runtimeJsonParse != nil {
		props["parse"] = g.emitBoxNativeFunc(g.runtimeJsonParse)
	}
	if g.runtimeToJSON != nil {
		props["stringify"] = g.emitBoxNativeFunc(g.runtimeToJSON)
	}
	if g.runtimeJsonTryParse != nil {
		props["tryparse"] = g.emitBoxNativeFunc(g.runtimeJsonTryParse)
		props["try_parse"] = g.emitBoxNativeFunc(g.runtimeJsonTryParse)
	}
	return g.emitObjectWithProps(props)
}

func (g *Generator) emitStdlibIOObject() (value.Value, error) {
	props := make(map[string]value.Value)
	if g.runtimeReadFile != nil {
		props["read"] = g.emitBoxNativeFunc(g.runtimeReadFile)
	}
	if g.runtimeWriteFile != nil {
		props["write"] = g.emitBoxNativeFunc(g.runtimeWriteFile)
	}
	if g.runtimeAppendFile != nil {
		props["append"] = g.emitBoxNativeFunc(g.runtimeAppendFile)
	}
	if g.runtimeFileExists != nil {
		props["exists"] = g.emitBoxNativeFunc(g.runtimeFileExists)
	}
	if g.runtimeDeleteFile != nil {
		props["remove"] = g.emitBoxNativeFunc(g.runtimeDeleteFile)
	}
	if g.runtimeIsFile != nil {
		props["isfile"] = g.emitBoxNativeFunc(g.runtimeIsFile)
	}
	if g.runtimeIsDir != nil {
		props["isdir"] = g.emitBoxNativeFunc(g.runtimeIsDir)
	}
	if g.runtimeFileSize != nil {
		props["size"] = g.emitBoxNativeFunc(g.runtimeFileSize)
	}
	if g.runtimeListDir != nil {
		props["list"] = g.emitBoxNativeFunc(g.runtimeListDir)
	}
	return g.emitObjectWithProps(props)
}

func (g *Generator) emitModuleExportObject(absPath string, prog *parser.Program) (value.Value, error) {
	savedLocals := g.locals
	savedLocalIsCell := g.localIsCell
	savedModuleGlobals := g.moduleGlobals
	savedModuleGlobalIsCell := g.moduleGlobalIsCell
	savedEmitPath := g.moduleEmitPath

	g.locals = make(map[string]value.Value)
	g.localIsCell = make(map[string]bool)
	g.moduleGlobals = make(map[string]value.Value)
	g.moduleGlobalIsCell = make(map[string]bool)
	g.moduleEmitPath = absPath

	moduleFuncs := make(map[string]*ir.Func)
	exportNames := make([]string, 0, len(prog.Declarations))

	for _, decl := range prog.Declarations {
		switch d := decl.(type) {
		case *parser.FuncDecl:
			exportNames = append(exportNames, d.Name.Lexeme)
			if err := g.emitFuncDecl(d); err != nil {
				g.restoreModuleEmit(savedLocals, savedLocalIsCell, savedModuleGlobals, savedModuleGlobalIsCell, savedEmitPath)
				return nil, err
			}
			if fn, ok := g.funcs[d.Name.Lexeme]; ok {
				moduleFuncs[d.Name.Lexeme] = fn
			}
		case *parser.LetDecl:
			exportNames = append(exportNames, d.Name.Lexeme)
			if err := g.emitLetDecl(d); err != nil {
				g.restoreModuleEmit(savedLocals, savedLocalIsCell, savedModuleGlobals, savedModuleGlobalIsCell, savedEmitPath)
				return nil, err
			}
		}
	}

	props := make(map[string]value.Value, len(exportNames))
	for _, name := range exportNames {
		val, err := g.loadModuleExportValue(name, moduleFuncs)
		if err != nil {
			g.restoreModuleEmit(savedLocals, savedLocalIsCell, savedModuleGlobals, savedModuleGlobalIsCell, savedEmitPath)
			return nil, err
		}
		props[name] = val
	}

	for name := range moduleFuncs {
		delete(g.funcs, name)
	}

	g.restoreModuleEmit(savedLocals, savedLocalIsCell, savedModuleGlobals, savedModuleGlobalIsCell, savedEmitPath)

	return g.emitObjectWithProps(props)
}

func (g *Generator) restoreModuleEmit(
	locals map[string]value.Value,
	localIsCell map[string]bool,
	moduleGlobals map[string]value.Value,
	moduleGlobalIsCell map[string]bool,
	emitPath string,
) {
	g.locals = locals
	g.localIsCell = localIsCell
	g.moduleGlobals = moduleGlobals
	g.moduleGlobalIsCell = moduleGlobalIsCell
	g.moduleEmitPath = emitPath
}

func (g *Generator) loadModuleExportValue(name string, moduleFuncs map[string]*ir.Func) (value.Value, error) {
	if fn, ok := moduleFuncs[name]; ok {
		return g.emitBoxNativeFunc(fn), nil
	}
	if slot, ok := g.locals[name]; ok {
		if g.localIsCell != nil && g.localIsCell[name] {
			return g.block.NewCall(g.runtimeCellRead, slot), nil
		}
		return g.block.NewLoad(types.I64, slot), nil
	}
	if slot, ok := g.moduleGlobals[name]; ok {
		if g.moduleGlobalIsCell != nil && g.moduleGlobalIsCell[name] {
			return g.block.NewCall(g.runtimeCellRead, slot), nil
		}
		return g.block.NewLoad(types.I64, slot), nil
	}
	return nil, fmt.Errorf("native codegen: module export %q not found after emission", name)
}
