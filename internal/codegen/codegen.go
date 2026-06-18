package codegen

import (
	"fmt"
	"strings"

	"koda/internal/diagnostic"
	"koda/internal/parser"
	"koda/internal/sema"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// PrepareNativeBundle prepares a bundle for native emission by performing
// capture analysis and local variable mapping.
func PrepareNativeBundle(bundle *parser.ProgramBundle) (*sema.NativeEmitContext, error) {
	return sema.PrepareNativeBundle(bundle)
}

func (g *Generator) undefinedVarError(name string, file string, line int, col int) error {
	candidates := make([]string, 0, len(g.locals)+len(g.moduleGlobals)+len(g.globals)+len(g.funcs))
	for k := range g.locals {
		candidates = append(candidates, k)
	}
	for k := range g.moduleGlobals {
		candidates = append(candidates, k)
	}
	for k := range g.globals {
		candidates = append(candidates, k)
	}
	for k := range g.funcs {
		candidates = append(candidates, k)
	}
	hint := ""
	if s, ok := diagnostic.BestSuggestion(name, candidates, 2); ok {
		hint = fmt.Sprintf("did you mean '%s'?", s)
	}
	srcPath := file
	if srcPath == "" {
		srcPath = g.sourcePath
	}
	return &diagnostic.DiagnosticError{
		File:    srcPath,
		Line:    line,
		Col:     col,
		Message: fmt.Sprintf("undefined variable '%s'", name),
		Hint:    hint,
	}
}

// EmitLLVMIR generates LLVM IR for the bundle carried by ctx.
func EmitLLVMIR(ctx *sema.NativeEmitContext) (*ir.Module, error) {
	gen := NewGenerator(ctx)
	return gen.Generate(ctx.Bundle)
}

type Generator struct {
	mod                     *ir.Module
	block                   *ir.Block
	ctx                     *sema.NativeEmitContext
	currentStructMethodType string // struct type while emitting a struct method body
	currentEmitFuncName     string // Koda function name while emitting a function body
	currentEmitFuncDecl     *parser.FuncDecl
	currentEmitFuncExpr     *parser.FuncExpr
	funcs                   map[string]*ir.Func    // Function name → LLVM function
	funcStubs               map[string]bool        // Names pre-declared for struct-method forward refs
	fastNativeFuncs         map[string]*ir.Func    // Lazily declared fusion fast-path natives
	locals                  map[string]value.Value // Variable name → stack slot
	globals                 map[string]*ir.Global  // Global variables
	currentFn               *ir.Func
	runtimeInit             *ir.Func
	runtimeInitEx           *ir.Func
	runtimeSetArgv          *ir.Func
	runtimeShutdown         *ir.Func
	runtimeRegisterGlobal   *ir.Func
	runtimeSetStackBase     *ir.Func
	runtimeDeltaTime        *ir.Func
	runtimeProgramTime      *ir.Func
	runtimeTimestamp        *ir.Func
	runtimeTime             *ir.Func
	runtimeClock            *ir.Func
	runtimeSleep            *ir.Func
	runtimePrint            *ir.Func
	runtimePrintArgv        *ir.Func
	runtimePrintNewline     *ir.Func
	runtimeWarn             *ir.Func
	runtimeRandom           *ir.Func
	runtimeRandomInt        *ir.Func
	runtimeRandomChoice     *ir.Func
	runtimeRandomSeed       *ir.Func
	runtimeLerp             *ir.Func
	runtimeClamp            *ir.Func
	runtimeDistance         *ir.Func
	runtimeAngleBetween     *ir.Func
	runtimeMap              *ir.Func
	runtimePI               *ir.Func
	runtimeE                *ir.Func
	runtimeSin              *ir.Func
	runtimeCos              *ir.Func
	runtimeTan              *ir.Func
	runtimeAsin             *ir.Func
	runtimeAcos             *ir.Func
	runtimeAtan             *ir.Func
	runtimeAtan2            *ir.Func
	runtimePow              *ir.Func
	runtimeExp              *ir.Func
	runtimeLog              *ir.Func
	runtimeLog10            *ir.Func
	runtimeLog2             *ir.Func
	runtimeFloor            *ir.Func
	runtimeCeil             *ir.Func
	runtimeRound            *ir.Func
	runtimeTrunc            *ir.Func
	runtimeSign             *ir.Func
	runtimeMin              *ir.Func
	runtimeMax              *ir.Func
	runtimeSmoothstep       *ir.Func
	runtimeDistanceSq       *ir.Func
	runtimeNormalize        *ir.Func
	runtimeHypot            *ir.Func
	runtimeFmod             *ir.Func
	runtimeDegrees          *ir.Func
	runtimeRadians          *ir.Func
	runtimeWrap             *ir.Func
	runtimeApproach         *ir.Func
	runtimeSmoothdamp       *ir.Func
	runtimeIsNumber         *ir.Func
	runtimeIsString         *ir.Func
	runtimeIsBool           *ir.Func
	runtimeIsNull           *ir.Func
	runtimeIsArray          *ir.Func
	runtimeIsObject         *ir.Func
	runtimeIsFunction       *ir.Func
	runtimeBool             *ir.Func
	runtimeFormat           *ir.Func
	runtimeArrayMap         *ir.Func
	runtimeArrayFilter      *ir.Func
	runtimeArrayForEach     *ir.Func
	runtimeArrayFind        *ir.Func
	runtimeArrayFindIndex   *ir.Func
	runtimeArraySome        *ir.Func
	runtimeArrayEvery       *ir.Func
	runtimeArrayReduce      *ir.Func
	runtimeArraySort        *ir.Func
	runtimeArrayReverse     *ir.Func
	runtimeArrayIndexOf     *ir.Func
	runtimeArrayIncludes    *ir.Func
	runtimeArraySlice       *ir.Func
	runtimeArrayConcat      *ir.Func
	runtimeArrayJoin        *ir.Func
	runtimeArrayFlat        *ir.Func
	runtimeStringSplit      *ir.Func
	runtimeStringTrim       *ir.Func
	runtimeStringUpper      *ir.Func
	runtimeStringLower      *ir.Func
	runtimeStringStartsWith *ir.Func
	runtimeStringEndsWith   *ir.Func
	runtimeStringIndexOf    *ir.Func
	runtimeStringSlice      *ir.Func
	runtimeStringReplace    *ir.Func
	runtimeStringReplaceAll *ir.Func
	runtimeStringPadStart   *ir.Func
	runtimeStringPadEnd     *ir.Func
	runtimeReadFile         *ir.Func
	runtimeAssetPath        *ir.Func
	runtimeArgs             *ir.Func
	runtimeEnv              *ir.Func
	runtimeRgb              *ir.Func
	runtimeRgba             *ir.Func
	runtimeVec2             *ir.Func
	runtimeVec3             *ir.Func
	runtimeRect             *ir.Func
	runtimeBox              *ir.Func
	runtimeColor            *ir.Func
	runtimeValueSub         *ir.Func
	runtimeValueMul         *ir.Func

	runtimeWriteFile        *ir.Func
	runtimeAppendFile       *ir.Func
	runtimeFileExists       *ir.Func
	runtimeDeleteFile       *ir.Func
	runtimeIsFile           *ir.Func
	runtimeIsDir            *ir.Func
	runtimeFileSize         *ir.Func
	runtimeListDir          *ir.Func
	runtimeKeys             *ir.Func
	runtimeAssert           *ir.Func
	runtimeTrace            *ir.Func
	runtimeParseJSON        *ir.Func
	runtimeJsonParse        *ir.Func
	runtimeJsonTryParse     *ir.Func
	runtimeToJSON           *ir.Func
	runtimeAllocObj         *ir.Func
	runtimeAllocStruct      *ir.Func
	runtimeStructGet        *ir.Func
	runtimeStructField      *ir.Func
	runtimeStructSet        *ir.Func
	runtimeObjGet           *ir.Func
	runtimeObjSet           *ir.Func
	runtimeObjRemove        *ir.Func
	runtimeUnboxNumber      *ir.Func
	runtimeBoxNumber        *ir.Func
	runtimeSet              *ir.Func
	runtimeAllocStr         *ir.Func
	runtimeAllocArray       *ir.Func
	runtimeArrayGet         *ir.Func
	runtimeArraySet         *ir.Func
	runtimeArrayPush        *ir.Func
	runtimeArrayPushArgv    *ir.Func
	runtimeArrayPop         *ir.Func
	runtimeArrayPopArgv     *ir.Func
	runtimeArrayRemoveAt    *ir.Func
	runtimeArrayClear       *ir.Func
	runtimeArrayLen         *ir.Func
	runtimeForOfLength      *ir.Func
	runtimeForOfKeyAt       *ir.Func
	runtimeForOfValueAt     *ir.Func
	runtimeType             *ir.Func
	runtimeLen              *ir.Func
	runtimeAbs              *ir.Func
	runtimeSqrt             *ir.Func
	runtimeCbrt             *ir.Func
	runtimeNumber           *ir.Func
	runtimeString           *ir.Func
	runtimeStringConcat     *ir.Func
	runtimeValueAdd         *ir.Func
	runtimeRange            *ir.Func
	runtimeAllocCell        *ir.Func
	runtimeCellRead         *ir.Func
	runtimeCellWrite        *ir.Func
	runtimeGcCollect        *ir.Func
	runtimeGcDisable        *ir.Func
	runtimeGcEnable         *ir.Func
	runtimeGcFrameStep      *ir.Func
	runtimeGcStats          *ir.Func
	runtimeInternClear      *ir.Func
	runtimeArena            *ir.Func
	runtimeArenaReset       *ir.Func
	runtimeArenaAllocArray  *ir.Func
	runtimeArenaAllocStruct *ir.Func
	runtimePushFrame        *ir.Func
	runtimePopFrame         *ir.Func
	runtimeOk               *ir.Func
	runtimeErr              *ir.Func
	runtimeValuesEqual      *ir.Func
	runtimePanic            *ir.Func
	runtimeMatches          *ir.Func
	runtimePushCall         *ir.Func
	runtimePopCall          *ir.Func
	runtimeMalloc           *ir.Func
	sourcePath              string
	dbg                     *debugInfo
	localIsCell             map[string]bool
	moduleGlobals           map[string]value.Value
	moduleGlobalIsCell      map[string]bool
	loopStack               []loopContext
	tempN                   int

	shadowLayout     *sema.ShadowLayout
	shadowFramePtr   value.Value
	shadowFrameArrTy *types.ArrayType
	shadowPushed     bool
	shadowTempNext   int

	// deferLayers: per active Koda function (and user_main), deferred expressions in source order (LIFO at exit).
	deferLayers [][]parser.Expr

	moduleEmitPath string
	importSlots    map[string]value.Value

	testFuncs []testEntry

	// Compile-time string literals cached in module globals (initialized once at startup).
	stringLiteralSlots map[string]*ir.Global
	stringLiteralData  []stringLiteralEntry
	stringLiteralInit  *ir.Func
}

type stringLiteralEntry struct {
	text string
	slot *ir.Global
}

type testEntry struct {
	display string
	fn      *ir.Func
}

// loopContext tracks information about the current loop for break/continue.
type loopContext struct {
	condBlock         *ir.Block
	incBlock          *ir.Block
	afterBlock        *ir.Block
	fallthroughTarget *ir.Block
}

// NewGenerator creates a new LLVM IR generator.
func NewGenerator(ctx *sema.NativeEmitContext) *Generator {
	mod := ir.NewModule()

	// Declare runtime functions
	runtimeFuncs := declareRuntimeFunctions(mod)

	gen := &Generator{
		mod:                     mod,
		ctx:                     ctx,
		funcs:                   make(map[string]*ir.Func),
		funcStubs:               make(map[string]bool),
		locals:                  make(map[string]value.Value),
		moduleGlobals:           make(map[string]value.Value),
		moduleGlobalIsCell:      make(map[string]bool),
		importSlots:             make(map[string]value.Value),
		globals:                 make(map[string]*ir.Global),
		runtimeInit:             runtimeFuncs["KODA_runtime_init"],
		runtimeInitEx:           runtimeFuncs["KODA_runtime_init_ex"],
		runtimeSetArgv:            runtimeFuncs["KODA_runtime_set_argv"],
		runtimeShutdown:         runtimeFuncs["KODA_runtime_shutdown"],
		runtimeRegisterGlobal:   runtimeFuncs["KODA_register_global_slot"],
		runtimeSetStackBase:     runtimeFuncs["KODA_runtime_set_stack_base"],
		runtimeDeltaTime:        runtimeFuncs["KODA_delta_time"],
		runtimeProgramTime:      runtimeFuncs["KODA_program_time"],
		runtimeTimestamp:        runtimeFuncs["KODA_timestamp"],
		runtimeTime:             runtimeFuncs["KODA_time"],
		runtimeClock:            runtimeFuncs["KODA_clock"],
		runtimeSleep:            runtimeFuncs["KODA_sleep"],
		runtimePrint:            runtimeFuncs["KODA_print"],
		runtimePrintArgv:        runtimeFuncs["KODA_print_argv"],
		runtimePrintNewline:     runtimeFuncs["KODA_print_newline"],
		runtimeWarn:             runtimeFuncs["KODA_warn"],
		runtimeRandom:           runtimeFuncs["KODA_random"],
		runtimeRandomInt:        runtimeFuncs["KODA_randomInt"],
		runtimeRandomChoice:     runtimeFuncs["KODA_randomChoice"],
		runtimeRandomSeed:       runtimeFuncs["KODA_randomSeed"],
		runtimeLerp:             runtimeFuncs["KODA_lerp"],
		runtimeClamp:            runtimeFuncs["KODA_clamp"],
		runtimeDistance:         runtimeFuncs["KODA_distance"],
		runtimeAngleBetween:     runtimeFuncs["KODA_angleBetween"],
		runtimeMap:              runtimeFuncs["KODA_map"],
		runtimePI:               runtimeFuncs["KODA_pi"],
		runtimeE:                runtimeFuncs["KODA_e"],
		runtimeSin:              runtimeFuncs["KODA_sin"],
		runtimeCos:              runtimeFuncs["KODA_cos"],
		runtimeTan:              runtimeFuncs["KODA_tan"],
		runtimeAsin:             runtimeFuncs["KODA_asin"],
		runtimeAcos:             runtimeFuncs["KODA_acos"],
		runtimeAtan:             runtimeFuncs["KODA_atan"],
		runtimeAtan2:            runtimeFuncs["KODA_atan2"],
		runtimePow:              runtimeFuncs["KODA_pow"],
		runtimeExp:              runtimeFuncs["KODA_exp"],
		runtimeLog:              runtimeFuncs["KODA_log"],
		runtimeLog10:            runtimeFuncs["KODA_log10"],
		runtimeLog2:             runtimeFuncs["KODA_log2"],
		runtimeFloor:            runtimeFuncs["KODA_floor"],
		runtimeCeil:             runtimeFuncs["KODA_ceil"],
		runtimeRound:            runtimeFuncs["KODA_round"],
		runtimeTrunc:            runtimeFuncs["KODA_trunc"],
		runtimeSign:             runtimeFuncs["KODA_sign"],
		runtimeMin:              runtimeFuncs["KODA_min"],
		runtimeMax:              runtimeFuncs["KODA_max"],
		runtimeSmoothstep:       runtimeFuncs["KODA_smoothstep"],
		runtimeDistanceSq:       runtimeFuncs["KODA_distanceSq"],
		runtimeNormalize:        runtimeFuncs["KODA_normalize"],
		runtimeHypot:            runtimeFuncs["KODA_hypot"],
		runtimeFmod:             runtimeFuncs["KODA_fmod"],
		runtimeDegrees:          runtimeFuncs["KODA_degrees"],
		runtimeRadians:          runtimeFuncs["KODA_radians"],
		runtimeWrap:             runtimeFuncs["KODA_wrap"],
		runtimeApproach:         runtimeFuncs["KODA_approach"],
		runtimeSmoothdamp:       runtimeFuncs["KODA_smoothdamp"],
		runtimeIsNumber:         runtimeFuncs["KODA_isNumber"],
		runtimeIsString:         runtimeFuncs["KODA_isString"],
		runtimeIsBool:           runtimeFuncs["KODA_isBool"],
		runtimeIsNull:           runtimeFuncs["KODA_isNull"],
		runtimeIsArray:          runtimeFuncs["KODA_isArray"],
		runtimeIsObject:         runtimeFuncs["KODA_isObject"],
		runtimeIsFunction:       runtimeFuncs["KODA_isFunction"],
		runtimeBool:             runtimeFuncs["KODA_bool"],
		runtimeFormat:           runtimeFuncs["KODA_format"],
		runtimeArrayMap:         runtimeFuncs["KODA_array_map"],
		runtimeArrayFilter:      runtimeFuncs["KODA_array_filter"],
		runtimeArrayForEach:     runtimeFuncs["KODA_array_forEach"],
		runtimeArrayFind:        runtimeFuncs["KODA_array_find"],
		runtimeArrayFindIndex:   runtimeFuncs["KODA_array_findIndex"],
		runtimeArraySome:        runtimeFuncs["KODA_array_some"],
		runtimeArrayEvery:       runtimeFuncs["KODA_array_every"],
		runtimeArrayReduce:      runtimeFuncs["KODA_array_reduce"],
		runtimeArraySort:        runtimeFuncs["KODA_array_sort"],
		runtimeArrayReverse:     runtimeFuncs["KODA_array_reverse"],
		runtimeArrayIndexOf:     runtimeFuncs["KODA_array_indexOf"],
		runtimeArrayIncludes:    runtimeFuncs["KODA_array_includes"],
		runtimeArraySlice:       runtimeFuncs["KODA_array_slice"],
		runtimeArrayConcat:      runtimeFuncs["KODA_array_concat"],
		runtimeArrayJoin:        runtimeFuncs["KODA_array_join"],
		runtimeArrayFlat:        runtimeFuncs["KODA_array_flat"],
		runtimeStringSplit:      runtimeFuncs["KODA_string_split"],
		runtimeStringTrim:       runtimeFuncs["KODA_string_trim"],
		runtimeStringUpper:      runtimeFuncs["KODA_string_upper"],
		runtimeStringLower:      runtimeFuncs["KODA_string_lower"],
		runtimeStringStartsWith: runtimeFuncs["KODA_string_startsWith"],
		runtimeStringEndsWith:   runtimeFuncs["KODA_string_endsWith"],
		runtimeStringIndexOf:    runtimeFuncs["KODA_string_indexOf"],
		runtimeStringSlice:      runtimeFuncs["KODA_string_slice"],
		runtimeStringReplace:    runtimeFuncs["KODA_string_replace"],
		runtimeStringReplaceAll: runtimeFuncs["KODA_string_replaceAll"],
		runtimeStringPadStart:   runtimeFuncs["KODA_string_padStart"],
		runtimeStringPadEnd:     runtimeFuncs["KODA_string_padEnd"],
		runtimeReadFile:         runtimeFuncs["KODA_readFile"],
		runtimeAssetPath:        runtimeFuncs["KODA_asset_path"],
		runtimeArgs:             runtimeFuncs["KODA_args"],
		runtimeEnv:              runtimeFuncs["KODA_env"],
		runtimeRgb:              runtimeFuncs["KODA_rgb"],
		runtimeRgba:             runtimeFuncs["KODA_rgba"],
		runtimeVec2:             runtimeFuncs["KODA_vec2"],
		runtimeVec3:             runtimeFuncs["KODA_vec3"],
		runtimeRect:             runtimeFuncs["KODA_rect"],
		runtimeBox:              runtimeFuncs["KODA_box"],
		runtimeColor:            runtimeFuncs["KODA_color"],
		runtimeValueSub:         runtimeFuncs["KODA_value_sub"],
		runtimeValueMul:         runtimeFuncs["KODA_value_mul"],
		runtimeWriteFile:        runtimeFuncs["KODA_writeFile"],
		runtimeAppendFile:       runtimeFuncs["KODA_appendFile"],
		runtimeFileExists:       runtimeFuncs["KODA_fileExists"],
		runtimeDeleteFile:       runtimeFuncs["KODA_deleteFile"],
		runtimeIsFile:           runtimeFuncs["KODA_isFile"],
		runtimeIsDir:            runtimeFuncs["KODA_isDir"],
		runtimeFileSize:         runtimeFuncs["KODA_fileSize"],
		runtimeListDir:          runtimeFuncs["KODA_listDir"],
		runtimeKeys:             runtimeFuncs["KODA_keys"],
		runtimeAssert:           runtimeFuncs["KODA_assert"],
		runtimeTrace:            runtimeFuncs["KODA_trace"],
		runtimeParseJSON:        runtimeFuncs["KODA_parseJSON"],
		runtimeJsonParse:        runtimeFuncs["KODA_jsonParse"],
		runtimeJsonTryParse:     runtimeFuncs["KODA_jsonTryParse"],
		runtimeToJSON:           runtimeFuncs["KODA_toJSON"],
		runtimeAllocObj:         runtimeFuncs["KODA_allocate_object"],
		runtimeAllocStruct:      runtimeFuncs["KODA_allocate_struct"],
		runtimeStructGet:        runtimeFuncs["KODA_struct_get"],
		runtimeStructField:      runtimeFuncs["KODA_struct_field"],
		runtimeStructSet:        runtimeFuncs["KODA_struct_set"],
		runtimeObjGet:           runtimeFuncs["KODA_object_get"],
		runtimeObjSet:           runtimeFuncs["KODA_object_set"],
		runtimeObjRemove:        runtimeFuncs["KODA_object_remove"],
		runtimeUnboxNumber:      runtimeFuncs["KODA_unbox_number"],
		runtimeBoxNumber:        runtimeFuncs["KODA_box_number"],
		runtimeSet:              runtimeFuncs["KODA_set"],
		runtimeAllocStr:         runtimeFuncs["KODA_allocate_string"],
		runtimeAllocArray:       runtimeFuncs["KODA_allocate_array"],
		runtimeArrayGet:         runtimeFuncs["KODA_array_get"],
		runtimeArraySet:         runtimeFuncs["KODA_array_set"],
		runtimeArrayPush:        runtimeFuncs["KODA_array_push"],
		runtimeArrayPushArgv:    runtimeFuncs["KODA_array_push_argv"],
		runtimeArrayPop:         runtimeFuncs["KODA_array_pop"],
		runtimeArrayPopArgv:     runtimeFuncs["KODA_array_pop_argv"],
		runtimeArrayRemoveAt:    runtimeFuncs["KODA_array_remove_at"],
		runtimeArrayClear:       runtimeFuncs["KODA_array_clear"],
		runtimeArrayLen:         runtimeFuncs["KODA_array_length"], // distinct symbol from KODA_len
		runtimeForOfLength:      runtimeFuncs["KODA_forof_length"],
		runtimeForOfKeyAt:       runtimeFuncs["KODA_forof_key_at"],
		runtimeForOfValueAt:     runtimeFuncs["KODA_forof_value_at"],
		runtimeType:             runtimeFuncs["KODA_type"],
		runtimeLen:              runtimeFuncs["KODA_len"],
		runtimeAbs:              runtimeFuncs["KODA_abs"],
		runtimeSqrt:             runtimeFuncs["KODA_sqrt"],
		runtimeCbrt:             runtimeFuncs["KODA_cbrt"],
		runtimeNumber:           runtimeFuncs["KODA_number"],
		runtimeString:           runtimeFuncs["KODA_string"],
		runtimeStringConcat:     runtimeFuncs["KODA_string_concat"],
		runtimeValueAdd:         runtimeFuncs["KODA_value_add"],
		runtimeRange:            runtimeFuncs["KODA_range"],
		runtimeAllocCell:        runtimeFuncs["KODA_alloc_cell"],
		runtimeCellRead:         runtimeFuncs["KODA_cell_read"],
		runtimeCellWrite:        runtimeFuncs["KODA_cell_write"],
		runtimeGcCollect:        runtimeFuncs["KODA_gc_collect"],
		runtimeGcDisable:        runtimeFuncs["KODA_gc_disable"],
		runtimeGcEnable:         runtimeFuncs["KODA_gc_enable"],
		runtimeGcFrameStep:      runtimeFuncs["KODA_gc_frame_step"],
		runtimeGcStats:          runtimeFuncs["KODA_gc_stats"],
		runtimeInternClear:      runtimeFuncs["KODA_intern_clear"],
		runtimeArena:            runtimeFuncs["KODA_arena"],
		runtimeArenaReset:       runtimeFuncs["KODA_arena_reset"],
		runtimeArenaAllocArray:  runtimeFuncs["KODA_arena_alloc_array"],
		runtimeArenaAllocStruct: runtimeFuncs["KODA_arena_alloc_struct"],
		runtimePushFrame:        runtimeFuncs["KODA_push_frame"],
		runtimePopFrame:         runtimeFuncs["KODA_pop_frame"],
		runtimeOk:               runtimeFuncs["KODA_ok"],
		runtimeErr:              runtimeFuncs["KODA_err"],
		runtimeValuesEqual:      runtimeFuncs["KODA_values_equal"],
		runtimePanic:            runtimeFuncs["KODA_panic"],
		runtimeMatches:          runtimeFuncs["KODA_matches"],
		runtimePushCall:         runtimeFuncs["KODA_push_call"],
		runtimePopCall:          runtimeFuncs["KODA_pop_call"],
		runtimeMalloc:           runtimeFuncs["malloc"],
		tempN:                   0,
	}

	return gen
}

// Generate emits LLVM IR for a program bundle.
func (g *Generator) Generate(bundle *parser.ProgramBundle) (*ir.Module, error) {
	entry := bundle.Entry
	if entry == nil {
		return nil, fmt.Errorf("no entry program")
	}

	g.sourcePath = "<entry>"
	if ep, err := parser.BundleEntryPath(bundle); err == nil {
		g.sourcePath = ep
	}
	if g.ctx != nil && g.ctx.EmitDebug {
		g.dbg = newDebugInfo(g.mod, g.sourcePath)
	}
	g.registerBuiltinFuncs()

	// Create user_main function (implicit entry point for all top-level code)
	userMain := g.mod.NewFunc("user_main", types.I64, ir.NewParam("this", types.I64))
	g.funcs["user_main"] = userMain
	if g.dbg != nil {
		g.dbg.subprogram(userMain, "user_main", g.sourcePath, 1)
	}
	g.currentFn = userMain

	// Start emitting into user_main's entry block
	g.block = userMain.NewBlock("entry")
	g.locals = make(map[string]value.Value)
	g.localIsCell = make(map[string]bool)

	// Allocate 'this' slot
	thisSlot := g.entryAlloca(types.I64)
	g.locals["this"] = thisSlot
	g.block.NewStore(constant.NewInt(types.I64, 0), thisSlot)

	prevShadowLayout := g.shadowLayout
	prevShadowFramePtr := g.shadowFramePtr
	prevShadowFrameArrTy := g.shadowFrameArrTy
	prevShadowPushed := g.shadowPushed
	prevShadowTempNext := g.shadowTempNext

	g.shadowLayout = g.ctx.ShadowEntry
	var emptyParams []string
	g.beginShadowFrame(g.ctx.ShadowEntry, thisSlot, emptyParams)
	g.emitCallTracePush("user_main", g.sourcePath, 0)
	g.pushDeferLayer()
	defer g.popDeferLayer()

	if err := g.prepareBundleBindings(bundle); err != nil {
		return nil, err
	}

	if err := g.emitStructMethods(); err != nil {
		return nil, err
	}

	// Emit all top-level declarations and statements
	for _, decl := range entry.Declarations {
		if err := g.emitDecl(decl); err != nil {
			return nil, err
		}
	}

	// User-defined `function main()` is emitted as LLVM symbol `koda_user_main`; invoke after top-level runs.
	if mainFn := g.funcs["main"]; mainFn != nil {
		g.block.NewCall(mainFn, constant.NewInt(types.I64, 0))
	} else if len(g.testFuncs) > 0 {
		g.emitTestRunner()
	}

	// Terminate user_main
	if g.block.Term == nil {
		if err := g.emitDefersForCurrentLayer(); err != nil {
			return nil, err
		}
		g.emitCallTracePop()
		g.emitShadowPop()
		g.block.NewRet(constant.NewInt(types.I64, 0))
	}
	g.shadowLayout = prevShadowLayout
	g.shadowFramePtr = prevShadowFramePtr
	g.shadowFrameArrTy = prevShadowFrameArrTy
	g.shadowPushed = prevShadowPushed
	g.shadowTempNext = prevShadowTempNext
	g.currentFn = nil

	g.emitStringLiteralsInit()

	// Create main function that calls the entry point
	mainFn := g.mod.NewFunc("main", types.I32,
		ir.NewParam("argc", types.I32),
		ir.NewParam("argv", types.NewPointer(types.NewPointer(types.I8))))
	entryBlock := mainFn.NewBlock("entry")
	g.block = entryBlock
	g.currentFn = mainFn

	// Call runtime init
	stackBaseSlot := g.block.NewAlloca(types.I8)
	stackBase := g.block.NewBitCast(stackBaseSlot, types.NewPointer(types.I8))
	g.block.NewCall(g.runtimeInitEx, stackBase)
	g.block.NewCall(g.runtimeSetArgv, mainFn.Params[0], mainFn.Params[1])
	if g.stringLiteralInit != nil {
		g.block.NewCall(g.stringLiteralInit, constant.NewInt(types.I64, 0))
	}

	// Call the user's main function if it exists
	if um := g.funcs["user_main"]; um != nil {
		g.block.NewCall(um, constant.NewInt(types.I64, 0))
	}

	g.block.NewCall(g.runtimeShutdown)
	g.block.NewRet(constant.NewInt(types.I32, 0))
	g.currentFn = nil

	return g.mod, nil
}

func (g *Generator) emitDecl(decl parser.Decl) error {
	switch d := decl.(type) {
	case *parser.FuncDecl:
		return g.emitFuncDecl(d)
	case *parser.TestDecl:
		return g.emitTestDecl(d)
	case *parser.LetDecl:
		return g.emitLetDecl(d)
	case *parser.IncludeDecl:
		return g.emitIncludeDecl(d)
	case *parser.UseDecl:
		return nil
	case parser.Stmt:
		return g.emitStmt(d)
	}
	return nil
}

func (g *Generator) emitIncludeDecl(_ *parser.IncludeDecl) error {
	// Module includes are handled at the parser level by loading and parsing the included file
	// The codegen just needs to skip the include declaration since the included
	// declarations are already merged into the AST before code generation
	return nil
}

// prepareBundleBindings registers module-level names before struct methods are emitted.
// Struct method bodies may call shim natives (drawcube, getmousex, …) declared in #includes
// that appear after the struct in the flattened entry, or only in included modules.
func (g *Generator) prepareBundleBindings(bundle *parser.ProgramBundle) error {
	if bundle == nil {
		return nil
	}
	userMain := g.funcs["user_main"]
	seenNative := make(map[string]bool)
	seenGlobal := make(map[string]bool)

	registerDecls := func(decls []parser.Decl) error {
		for _, decl := range decls {
			ld, ok := decl.(*parser.LetDecl)
			if !ok {
				continue
			}
			name := ld.Name.Lexeme
			if ld.Native != nil {
				key := strings.ToLower(name)
				if seenNative[key] {
					continue
				}
				if _, exists := g.funcs[name]; exists {
					seenNative[key] = true
					continue
				}
				if err := g.emitNativeExternLet(ld); err != nil {
					return err
				}
				seenNative[key] = true
				continue
			}
			if g.currentFn != userMain {
				continue
			}
			if seenGlobal[name] {
				continue
			}
			if _, exists := g.globals[name]; exists {
				seenGlobal[name] = true
				continue
			}
			global := g.mod.NewGlobalDef("koda_global_"+name, constant.NewInt(types.I64, llvmNilTagged))
			g.globals[name] = global
			g.moduleGlobals[name] = global
			g.moduleGlobalIsCell[name] = false
			seenGlobal[name] = true
		}
		return nil
	}

	if bundle.Entry != nil {
		if err := registerDecls(bundle.Entry.Declarations); err != nil {
			return err
		}
	}
	structMethodDecls := g.structMethodDeclSet()
	registerFuncDecls := func(decls []parser.Decl, modulePath string) {
		for _, decl := range decls {
			fd, ok := decl.(*parser.FuncDecl)
			if !ok || fd.Native != nil || structMethodDecls[fd] {
				continue
			}
			g.declareFuncStub(fd, modulePath)
		}
	}
	if bundle.Entry != nil {
		registerFuncDecls(bundle.Entry.Declarations, "")
	}
	for modulePath, prog := range bundle.Modules {
		if prog == nil {
			continue
		}
		registerFuncDecls(prog.Declarations, modulePath)
	}
	for _, prog := range bundle.Modules {
		if prog == nil {
			continue
		}
		if err := registerDecls(prog.Declarations); err != nil {
			return err
		}
	}
	return nil
}

func (g *Generator) structMethodDeclSet() map[*parser.FuncDecl]bool {
	set := make(map[*parser.FuncDecl]bool)
	if g.ctx == nil {
		return set
	}
	for _, methods := range g.ctx.StructMethods {
		for _, fd := range methods {
			set[fd] = true
		}
	}
	return set
}

func (g *Generator) declareFuncStub(d *parser.FuncDecl, modulePath string) {
	name := d.Name.Lexeme
	if _, exists := g.funcs[name]; exists {
		return
	}
	llvmName := name
	if modulePath != "" && !strings.EqualFold(name, "main") {
		llvmName = moduleFuncLLVMName(modulePath, name)
	} else if strings.EqualFold(name, "main") {
		llvmName = "koda_user_main"
	}
	skipSelf := structMethodSkipsSelfParam(g, d)
	params := []*ir.Param{ir.NewParam("this", types.I64)}
	for i, param := range d.Params {
		if skipSelf && i == 0 {
			continue
		}
		params = append(params, ir.NewParam(param.Name, types.I64))
	}
	fn := g.mod.NewFunc(llvmName, types.I64, params...)
	g.funcs[llvmName] = fn
	g.funcs[name] = fn
	g.funcs[strings.ToLower(name)] = fn
	g.funcStubs[llvmName] = true
	g.funcStubs[name] = true
}

func (g *Generator) emitLetDecl(d *parser.LetDecl) error {
	name := d.Name.Lexeme

	if d.Native != nil {
		if _, exists := g.funcs[name]; exists {
			return nil
		}
		return g.emitNativeExternLet(d)
	}

	if g.currentFn == g.funcs["user_main"] {
		if global, exists := g.globals[name]; exists {
			g.locals[name] = global
			g.localIsCell[name] = false
			if g.runtimeRegisterGlobal != nil {
				g.block.NewCall(g.runtimeRegisterGlobal, global)
			}
			if d.Init != nil {
				initVal, err := g.emitExpr(d.Init)
				if err != nil {
					return err
				}
				g.block.NewStore(g.emitAsKodaI64(initVal), global)
			}
			return nil
		}
		global := g.mod.NewGlobalDef("koda_global_"+name, constant.NewInt(types.I64, llvmNilTagged))
		g.globals[name] = global
		g.locals[name] = global
		g.localIsCell[name] = false
		if g.runtimeRegisterGlobal != nil {
			g.block.NewCall(g.runtimeRegisterGlobal, global)
		}
		if d.Init != nil {
			initVal, err := g.emitExpr(d.Init)
			if err != nil {
				return err
			}
			g.block.NewStore(g.emitAsKodaI64(initVal), global)
		}
		return nil
	}

	// If we're not in a function (g.block is nil), create a global variable
	if g.block == nil {
		// Create a global variable with pointer type
		global := g.mod.NewGlobalDef(name, constant.NewNull(types.NewPointer(types.I64)))
		g.globals[name] = global
		g.locals[name] = global // Store in locals for easy access

		// Initialize if there's an initial value
		if d.Init != nil {
			// For global variables, we need to initialize them in a constructor
			// For now, just store the initial value in the global
			initVal, err := g.emitExpr(d.Init)
			if err != nil {
				return err
			}
			// Store the initial value into the global
			// This needs to be done in the main function, not here
			// For now, we'll skip this and handle it in main
			_ = initVal
		}
		return nil
	}

	var storageSlot value.Value
	var useStack bool
	var typedAnnot string
	if g.ctx != nil && g.ctx.TypedLocals != nil {
		if annot, ok := g.ctx.TypedLocals[d]; ok {
			typedAnnot = annot
			useStack = true
			if isFloatTypeAnnot(annot) {
				storageSlot = g.entryAlloca(types.Double)
			} else {
				storageSlot = g.entryAlloca(llvmIntTypeForAnnot(annot))
			}
			g.localIsCell[name] = false
		}
	}
	if typedAnnot == "" {
		if k, ok := g.ctx.NumericKinds[d]; ok && g.ctx.StackDecls[d] {
			switch k {
			case sema.KindFloat:
				typedAnnot = "float"
				useStack = true
				storageSlot = g.entryAlloca(types.Double)
				g.localIsCell[name] = false
			case sema.KindInt:
				useStack = true
				storageSlot = g.entryAlloca(types.I64)
				g.localIsCell[name] = false
			}
		}
		if !useStack {
			switch {
			case g.ctx.StackDecls[d]:
				useStack = true
				storageSlot = g.entryAlloca(types.I64)
				g.localIsCell[name] = false
			case g.ctx.EscapingDecls[d]:
				useStack = false
				storageSlot = g.block.NewCall(g.runtimeAllocCell)
				g.localIsCell[name] = true
			default:
				/* Conservative: heap cell if escape status unknown (must not keep stack slot past return). */
				useStack = false
				storageSlot = g.block.NewCall(g.runtimeAllocCell)
				g.localIsCell[name] = true
			}
		}
	}
	g.locals[name] = storageSlot
	if typedAnnot == "" {
		g.shadowStoreLet(d, storageSlot)
	}
	if g.currentFn == g.funcs["user_main"] {
		g.moduleGlobals[name] = storageSlot
		g.moduleGlobalIsCell[name] = !useStack
	}
	if g.currentFn == g.funcs["user_main"] && g.runtimeRegisterGlobal != nil {
		g.block.NewCall(g.runtimeRegisterGlobal, storageSlot)
	}

	if useStack {
		if typedAnnot != "" {
			if isFloatTypeAnnot(typedAnnot) {
				g.block.NewStore(constant.NewFloat(types.Double, 0.0), storageSlot)
			} else {
				g.storeIntLocal(storageSlot, typedAnnot, constant.NewInt(types.I64, 0))
			}
		} else {
			g.block.NewStore(constant.NewInt(types.I64, llvmNilTagged), storageSlot)
		}
	} else {
		g.block.NewCall(g.runtimeCellWrite, storageSlot, constant.NewInt(types.I64, llvmNilTagged))
	}

	if d.Init != nil {
		initVal, err := g.emitExpr(d.Init)
		if err != nil {
			return err
		}
		if typedAnnot != "" {
			if isFloatTypeAnnot(typedAnnot) {
				if lit, ok := d.Init.(*parser.LiteralExpr); ok {
					if fv, ok2 := lit.Value.(float64); ok2 {
						g.block.NewStore(g.floatLiteralForInit(fv), storageSlot)
						return nil
					}
					if iv, ok2 := lit.Value.(int); ok2 {
						g.block.NewStore(g.floatLiteralForInit(float64(iv)), storageSlot)
						return nil
					}
				}
				g.storeFloatLocal(storageSlot, initVal)
				return nil
			}
			if lit, ok := d.Init.(*parser.LiteralExpr); ok {
				if iv, ok2 := literalInt64(lit); ok2 {
					g.storeIntLocal(storageSlot, typedAnnot, g.intLiteralForAnnot(iv, typedAnnot))
					return nil
				}
			}
			boxed := g.emitAsKodaI64(initVal)
			unboxed := g.block.NewCall(g.runtimeUnboxNumber, boxed)
			asInt := g.block.NewFPToSI(unboxed, types.I64)
			g.storeIntLocal(storageSlot, typedAnnot, asInt)
			return nil
		}
		boxed := g.emitAsKodaI64(initVal)
		if useStack {
			g.block.NewStore(boxed, storageSlot)
		} else {
			g.block.NewCall(g.runtimeCellWrite, storageSlot, boxed)
		}
	}

	return nil
}

// emitNativeExternLet binds a Koda name to an LLVM declaration for a C symbol
// implementing KodaValue (*)(int argCount, KodaValue* args) (i64, i32, i64* in IR).
func (g *Generator) emitNativeExternLet(d *parser.LetDecl) error {
	sym := strings.TrimSpace(d.Native.Symbol)
	if sym == "" {
		return fmt.Errorf("native extern missing C symbol for %q", d.Name.Lexeme)
	}
	fn := g.ensureNativeExternFunc(sym)
	name := d.Name.Lexeme
	g.funcs[name] = fn
	g.funcs[strings.ToLower(name)] = fn
	if kn := strings.TrimSpace(d.Native.BindingName); kn != "" && kn != name {
		g.funcs[kn] = fn
		g.funcs[strings.ToLower(kn)] = fn
	}
	// Native bindings are declarations only; ignore initializer (typically `0`).
	return nil
}

func (g *Generator) ensureNativeExternFunc(symbol string) *ir.Func {
	for _, f := range g.mod.Funcs {
		if f.Name() == symbol {
			return f
		}
	}
	return g.mod.NewFunc(symbol, types.I64,
		ir.NewParam("arg_count", types.I32),
		ir.NewParam("args", types.NewPointer(types.I64)))
}

func (g *Generator) emitExpr(expr parser.Expr) (value.Value, error) {
	switch e := expr.(type) {
	case *parser.LiteralExpr:
		return g.emitLiteral(e)
	case *parser.IdentifierExpr:
		return g.emitIdentifier(e)
	case *parser.PrefixExpr:
		return g.emitPrefix(e)
	case *parser.InfixExpr:
		return g.emitInfix(e)
	case *parser.CallExpr:
		return g.emitCall(e)
	case *parser.ObjectExpr:
		return g.emitObject(e)
	case *parser.ArrayExpr:
		return g.emitArray(e)
	case *parser.IndexExpr:
		return g.emitIndex(e)
	case *parser.SliceExpr:
		return g.emitSlice(e)
	case *parser.AssignExpr:
		return g.emitAssign(e)
	case *parser.FuncExpr:
		return g.emitFuncExpr(e)
	case *parser.ThisExpr:
		return g.emitThis(e)
	case *parser.TemplateExpr:
		return g.emitTemplate(e)
	case *parser.RangeExpr:
		return g.emitRange(e)
	case *parser.TupleExpr:
		return g.emitTuple(e)
	case *parser.ImportExpr:
		return g.emitImport(e)
	case *parser.SwitchExpr:
		return g.emitSwitchExpr(e)
	case *parser.UpdateExpr:
		return g.emitUpdate(e)
	default:
		return nil, fmt.Errorf("unsupported expression type: %T", expr)
	}
}

func (g *Generator) emitPrefix(e *parser.PrefixExpr) (value.Value, error) {
	switch e.Operator {
	case "typeof":
		rv, err := g.emitExpr(e.Right)
		if err != nil {
			return nil, err
		}
		return g.block.NewCall(g.runtimeType, g.emitAsKodaI64(rv)), nil
	case "!", "not":
		right, err := g.emitExpr(e.Right)
		if err != nil {
			return nil, err
		}
		truthy := g.emitTruthy(right)
		notTruthy := g.block.NewXor(truthy, constant.NewBool(true))
		return g.emitBoxBoolNaN(notTruthy), nil
	case "+":
		if res, ok := compileTimeInt64(e); ok {
			return g.block.NewCall(g.runtimeBoxNumber, constant.NewFloat(types.Double, float64(res))), nil
		}
		return g.emitExpr(e.Right)
	case "-":
		if res, ok := compileTimeInt64(e); ok {
			return g.block.NewCall(g.runtimeBoxNumber, constant.NewFloat(types.Double, float64(res))), nil
		}
		right, err := g.emitExpr(e.Right)
		if err != nil {
			return nil, err
		}
		ri := g.emitAsKodaI64(right)
		rd := g.block.NewCall(g.runtimeUnboxNumber, ri)
		neg := g.block.NewFNeg(rd)
		return g.block.NewCall(g.runtimeBoxNumber, neg), nil
	default:
		return nil, fmt.Errorf("unsupported prefix operator: %s", e.Operator)
	}
}

// emitBoxBoolNaN maps i1 (predicate result) to Koda true/false NaN-boxed values.
func (g *Generator) emitBoxBoolNaN(cmp value.Value) value.Value {
	falseVal := constant.NewInt(types.I64, 0x7ffc000000000002)
	trueVal := constant.NewInt(types.I64, 0x7ffc000000000003)
	return g.block.NewSelect(cmp, trueVal, falseVal)
}

func (g *Generator) emitTruthy(v value.Value) value.Value {
	if v.Type().Equal(types.I1) {
		return v
	}
	if !v.Type().Equal(types.I64) {
		v = g.emitAsKodaI64(v)
	}
	isZero := g.block.NewICmp(enum.IPredEQ, v, constant.NewInt(types.I64, 0))
	isNil := g.block.NewICmp(enum.IPredEQ, v, constant.NewInt(types.I64, 0x7ffc000000000001))
	isFalse := g.block.NewICmp(enum.IPredEQ, v, constant.NewInt(types.I64, 0x7ffc000000000002))
	zeroOrNil := g.block.NewOr(isZero, isNil)
	falsey := g.block.NewOr(zeroOrNil, isFalse)
	return g.block.NewXor(falsey, constant.NewBool(true))
}

// isNativeArgvCallee reports whether fn uses the embedded runtime convention
// KodaValue fn(int argCount, KodaValue* args) (KodaValue as i64).
func isNativeArgvCallee(fn *ir.Func) bool {
	if fn == nil || len(fn.Params) != 2 {
		return false
	}
	if fn.Params[0].Typ != types.I32 {
		return false
	}
	return fn.Params[1].Typ.Equal(types.NewPointer(types.I64))
}

func (g *Generator) emitAsKodaI64(v value.Value) value.Value {
	if fn, ok := v.(*ir.Func); ok {
		return g.block.NewPtrToInt(fn, types.I64)
	}
	switch v.Type() {
	case types.I64:
		return v
	case types.I32:
		return g.block.NewSExt(v, types.I64)
	case types.I1:
		return g.block.NewZExt(v, types.I64)
	case types.Double:
		return g.block.NewBitCast(v, types.I64)
	default:
		return g.block.NewBitCast(v, types.I64)
	}
}

// emitArgvRuntime calls a runtime function declared as i64 name(i32 arg_count, i64* argv).
func (g *Generator) emitArgvRuntime(fn *ir.Func, args []value.Value) value.Value {
	n := len(args)
	zero := constant.NewInt(types.I32, 0)
	if n == 0 {
		return g.block.NewCall(fn, constant.NewInt(types.I32, 0), constant.NewNull(types.NewPointer(types.I64)))
	}
	arrTy := types.NewArray(uint64(n), types.I64)
	slot := g.entryAlloca(arrTy)
	for i, arg := range args {
		argI64 := g.emitAsKodaI64(arg)
		elemPtr := g.block.NewGetElementPtr(arrTy, slot, zero, constant.NewInt(types.I32, int64(i)))
		g.block.NewStore(argI64, elemPtr)
	}
	argvPtr := g.block.NewGetElementPtr(arrTy, slot, zero, zero)
	return g.block.NewCall(fn, constant.NewInt(types.I32, int64(n)), argvPtr)
}

// indirectKodaFuncPtrType is i64 (i64 this, i64 × argCount)* — matches every user/closure function in this backend.
func indirectKodaFuncPtrType(argCount int) *types.PointerType {
	paramTys := make([]types.Type, 0, 1+argCount)
	for i := 0; i < 1+argCount; i++ {
		paramTys = append(paramTys, types.I64)
	}
	return types.NewPointer(types.NewFunc(types.I64, paramTys...))
}

// maxClosureFreeVars is the maximum number of captured cells we can indirect-call through LLVM
// with a fixed function-pointer type per arity (see emitIndirectI64Callee).
const maxClosureFreeVars = 16

// indirectClosureFuncPtrType is i64 (i64 this, i64* × nCells, i64 × nArgs)* for closure LLVM functions.
func indirectClosureFuncPtrType(nCells, nArgs int) *types.PointerType {
	paramTys := make([]types.Type, 0, 1+nCells+nArgs)
	paramTys = append(paramTys, types.I64)
	cellPtr := types.NewPointer(types.I64)
	for i := 0; i < nCells; i++ {
		paramTys = append(paramTys, cellPtr)
	}
	for i := 0; i < nArgs; i++ {
		paramTys = append(paramTys, types.I64)
	}
	return types.NewPointer(types.NewFunc(types.I64, paramTys...))
}

// emitIndirectI64Callee calls a Koda "function value" stored as i64: either a tagged raw function
// pointer (LSB set) with ABI (this, args...), or an untagged heap pointer to { fn, n, cell0..n-1 }.
func (g *Generator) emitIndirectI64Callee(fnVal, thisVal value.Value, args []value.Value) (value.Value, error) {
	g.tempN++
	suf := fmt.Sprintf(".ic%d", g.tempN)

	v := g.emitAsKodaI64(fnVal)
	one := constant.NewInt(types.I64, 1)
	isRaw := g.block.NewICmp(enum.IPredEQ, g.block.NewAnd(v, one), one)

	rawB := g.currentFn.NewBlock("indcall.raw" + suf)
	heapB := g.currentFn.NewBlock("indcall.heap" + suf)
	g.block.NewCondBr(isRaw, rawB, heapB)

	mergeB := g.currentFn.NewBlock("indcall.merge" + suf)

	// --- Raw (zero-capture) function pointer: (ptr | 1) ---
	g.block = rawB
	maskedFn := g.block.NewAnd(v, constant.NewInt(types.I64, -2))
	rawFnPtr := g.block.NewIntToPtr(maskedFn, indirectKodaFuncPtrType(len(args)))
	rawArgs := append([]value.Value{thisVal}, args...)
	rawOut := g.block.NewCall(rawFnPtr, rawArgs...)
	g.block.NewBr(mergeB)

	// --- Heap closure: i64* env at ptrtoint ---
	g.block = heapB
	i64PtrTy := types.NewPointer(types.I64)
	envPtr := g.block.NewIntToPtr(v, i64PtrTy)
	fnRaw := g.block.NewLoad(types.I64, envPtr)
	nPtr := g.block.NewGetElementPtr(types.I64, envPtr, constant.NewInt(types.I32, 1))
	nCells := g.block.NewLoad(types.I64, nPtr)

	defaultB := g.currentFn.NewBlock("indcall.bad" + suf)
	caseBlocks := make([]*ir.Block, maxClosureFreeVars)
	cases := make([]*ir.Case, maxClosureFreeVars)
	outs := make([]value.Value, maxClosureFreeVars)
	for k := 1; k <= maxClosureFreeVars; k++ {
		bk := g.currentFn.NewBlock(fmt.Sprintf("indcall.n%d%s", k, suf))
		caseBlocks[k-1] = bk
		cases[k-1] = ir.NewCase(constant.NewInt(types.I64, int64(k)), bk)
	}
	g.block.NewSwitch(nCells, defaultB, cases...)

	g.block = defaultB
	g.block.NewUnreachable()

	cellPtrTy := types.NewPointer(types.I64)
	for k := 1; k <= maxClosureFreeVars; k++ {
		g.block = caseBlocks[k-1]
		finalArgs := []value.Value{thisVal}
		for i := 0; i < k; i++ {
			elemPtr := g.block.NewGetElementPtr(types.I64, envPtr, constant.NewInt(types.I32, int64(2+i)))
			ci := g.block.NewLoad(types.I64, elemPtr)
			finalArgs = append(finalArgs, g.block.NewIntToPtr(ci, cellPtrTy))
		}
		finalArgs = append(finalArgs, args...)
		clTy := indirectClosureFuncPtrType(k, len(args))
		clFnPtr := g.block.NewIntToPtr(fnRaw, clTy)
		outs[k-1] = g.block.NewCall(clFnPtr, finalArgs...)
		g.block.NewBr(mergeB)
	}

	g.block = mergeB
	incomings := make([]*ir.Incoming, 0, 1+maxClosureFreeVars)
	incomings = append(incomings, ir.NewIncoming(rawOut, rawB))
	for k := 0; k < maxClosureFreeVars; k++ {
		incomings = append(incomings, ir.NewIncoming(outs[k], caseBlocks[k]))
	}
	return g.block.NewPhi(incomings...), nil
}

func (g *Generator) loadBindingI64(name string) (value.Value, bool) {
	if slot, ok := g.locals[name]; ok {
		if g.localIsCell != nil && g.localIsCell[name] {
			return g.block.NewCall(g.runtimeCellRead, slot), true
		}
		return g.block.NewLoad(types.I64, slot), true
	}
	if slot, ok := g.moduleGlobals[name]; ok {
		if g.moduleGlobalIsCell != nil && g.moduleGlobalIsCell[name] {
			return g.block.NewCall(g.runtimeCellRead, slot), true
		}
		return g.block.NewLoad(types.I64, slot), true
	}
	if global, ok := g.globals[name]; ok {
		return g.block.NewLoad(types.I64, global), true
	}
	return nil, false
}

// tryEmitStaticCall emits a direct call when the callee is a known *ir.Func (native extern or user fn).
func (g *Generator) tryEmitStaticCall(fn *ir.Func, thisVal value.Value, args []value.Value, e *parser.CallExpr) (value.Value, bool, error) {
	if fn == nil {
		return nil, false, nil
	}
	zero := constant.NewInt(types.I64, 0)

	if fn == g.runtimePrint {
		if len(args) == 0 {
			g.block.NewCall(g.runtimePrintNewline)
			return zero, true, nil
		}
		if len(args) == 1 {
			return g.block.NewCall(g.runtimePrint, g.emitAsKodaI64(args[0])), true, nil
		}
		return g.emitArgvRuntime(g.runtimePrintArgv, args), true, nil
	}

	if fn == g.runtimeWarn {
		if len(args) == 0 {
			return zero, true, nil
		}
		return g.emitArgvRuntime(g.runtimeWarn, args), true, nil
	}

	if fn == g.runtimeGcFrameStep {
		var budget value.Value
		if len(args) == 0 {
			budget = constant.NewFloat(types.Double, 0)
		} else {
			budget = g.block.NewCall(g.runtimeUnboxNumber, g.emitAsKodaI64(args[0]))
		}
		g.block.NewCall(g.runtimeGcFrameStep, budget)
		return zero, true, nil
	}

	if isNativeArgvCallee(fn) {
		argCount := len(args)
		zero32 := constant.NewInt(types.I32, 0)
		var argvPtr value.Value
		if argCount == 0 {
			argvPtr = constant.NewNull(types.NewPointer(types.I64))
		} else {
			arrTy := types.NewArray(uint64(argCount), types.I64)
			slot := g.entryAlloca(arrTy)
			for i, arg := range args {
				argI64 := g.emitAsKodaI64(arg)
				elemPtr := g.block.NewGetElementPtr(arrTy, slot, zero32, constant.NewInt(types.I32, int64(i)))
				g.block.NewStore(argI64, elemPtr)
			}
			argvPtr = g.block.NewGetElementPtr(arrTy, slot, zero32, zero32)
		}
		call := g.block.NewCall(fn, constant.NewInt(types.I32, int64(argCount)), argvPtr)
		if g.dbg != nil && e != nil {
			attachDbg(call, g.dbg.loc(e.Token.File, e.Token.Line))
		}
		if types.Equal(fn.Sig.RetType, types.Void) {
			_ = call
			return zero, true, nil
		}
		return call, true, nil
	}

	if len(fn.Params) == len(args)+1 {
		finalArgs := []value.Value{thisVal}
		finalArgs = append(finalArgs, args...)
		call := g.block.NewCall(fn, finalArgs...)
		if types.Equal(fn.Sig.RetType, types.Void) {
			_ = call
			return zero, true, nil
		}
		return call, true, nil
	}

	return nil, false, nil
}

func (g *Generator) emitCall(e *parser.CallExpr) (value.Value, error) {
	if v, handled, err := g.tryEmitStructConstructor(e); handled {
		return v, err
	}
	// Check if this is a method call (e.g., obj.method())
	// Set this value if the function is called on an object
	var objForThis value.Value
	var memberExpr *parser.IndexExpr
	if ix, ok := e.Function.(*parser.IndexExpr); ok {
		memberExpr = ix
		obj, err := g.emitExpr(ix.Object)
		if err != nil {
			return nil, err
		}
		objForThis = obj
	}

	if memberExpr != nil && objForThis != nil {
		if v, handled, err := g.tryEmitStructMethodCall(memberExpr, objForThis, e); handled {
			return v, err
		}
		if v, handled, err := g.tryEmitMethodCall(memberExpr, objForThis, e); handled {
			return v, err
		}
	}

	var staticFn *ir.Func
	var calleeName string
	if id, ok := e.Function.(*parser.IdentifierExpr); ok {
		calleeName = id.Name.Lexeme
		staticFn = g.funcs[id.Name.Lexeme]
	}

	if calleeName != "" {
		if v, handled, err := g.tryEmitFusedNativeCall(calleeName, e.Arguments); handled {
			return v, err
		}
	}

	var args []value.Value
	for _, arg := range e.Arguments {
		val, err := g.emitExpr(arg)
		if err != nil {
			return nil, err
		}
		args = append(args, val)
	}

	if staticFn != nil && memberExpr == nil {
		thisVal := constant.NewInt(types.I64, 0)
		if v, handled, err := g.tryEmitStaticCall(staticFn, thisVal, args, e); handled {
			return v, err
		}
	}

	// Emit the function (dynamic callee or method value)
	fnVal, err := g.emitExpr(e.Function)
	if err != nil {
		return nil, err
	}

	if fn, ok := fnVal.(*ir.Func); ok && memberExpr == nil {
		thisVal := constant.NewInt(types.I64, 0)
		if v, handled, err := g.tryEmitStaticCall(fn, thisVal, args, e); handled {
			return v, err
		}
	}

	fnSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(fnSlot)
	g.block.NewStore(g.emitAsKodaI64(fnVal), fnSlot)

	// Set this for this call (JavaScript-like behavior)
	var thisVal value.Value
	if objForThis != nil {
		thisVal = objForThis
	} else {
		// In standalone function calls, this is undefined/null
		thisVal = constant.NewInt(types.I64, 0)
	}
	thisSlot := g.entryAlloca(types.I64)
	g.shadowStoreTemp(thisSlot)
	g.block.NewStore(g.emitAsKodaI64(thisVal), thisSlot)

	fnLive := g.block.NewLoad(types.I64, fnSlot)
	thisLive := g.block.NewLoad(types.I64, thisSlot)

	if fn, ok := fnVal.(*ir.Func); ok && len(fn.Params) == len(args)+1 {
		finalArgs := []value.Value{thisLive}
		finalArgs = append(finalArgs, args...)
		call := g.block.NewCall(fn, finalArgs...)
		if types.Equal(fn.Sig.RetType, types.Void) {
			_ = call
			return constant.NewInt(types.I64, 0), nil
		}
		return call, nil
	}

	// Function value as i64: tagged raw pointer (this + args) or heap closure (this + cells + args).
	if fnVal.Type().Equal(types.I64) {
		return g.emitIndirectI64Callee(fnLive, thisLive, args)
	}

	call := g.block.NewCall(fnVal, args...)
	if fn, ok := fnVal.(*ir.Func); ok && types.Equal(fn.Sig.RetType, types.Void) {
		_ = call
		return constant.NewInt(types.I64, 0), nil
	}
	return call, nil
}
