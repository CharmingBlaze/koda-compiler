package codegen

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

// argvI64 declares Value name(int arg_count, Value* args) as i64(i32, i64*).
func argvI64(mod *ir.Module, cName string) *ir.Func {
	return mod.NewFunc(cName, types.I64,
		ir.NewParam("arg_count", types.I32),
		ir.NewParam("args", types.NewPointer(types.I64)))
}

// declareRuntimeFunctions declares all runtime functions used by the generated code.
// LLVM symbol names and calling conventions must match runtime/src/koda_runtime.c.
func declareRuntimeFunctions(mod *ir.Module) map[string]*ir.Func {
	functions := make(map[string]*ir.Func)

	functions["KODA_runtime_init"] = mod.NewFunc("koda_runtime_init", types.Void)
	functions["KODA_runtime_init_ex"] = mod.NewFunc("koda_runtime_init_ex", types.Void,
		ir.NewParam("stack_base", types.NewPointer(types.I8)))
	functions["KODA_runtime_shutdown"] = mod.NewFunc("koda_runtime_shutdown", types.Void)
	functions["KODA_register_global_slot"] = mod.NewFunc("koda_register_global_slot", types.Void,
		ir.NewParam("slot", types.NewPointer(types.I64)))
	functions["KODA_runtime_set_stack_base"] = mod.NewFunc("koda_runtime_set_stack_base", types.Void,
		ir.NewParam("base", types.NewPointer(types.I8)))

	// print: koda_print_val for one value; koda_print_argv for multiple (space-separated, one newline).
	functions["KODA_print"] = mod.NewFunc("koda_print_val", types.I64,
		ir.NewParam("val", types.I64))
	functions["KODA_print_argv"] = argvI64(mod, "koda_print_argv")
	functions["KODA_print_newline"] = mod.NewFunc("koda_print_newline", types.Void)
	functions["KODA_warn"] = argvI64(mod, "koda_warn")
	typeFn := mod.NewFunc("koda_type", types.I64, ir.NewParam("val", types.I64))
	functions["KODA_type"] = typeFn
	functions["KODA_typeof"] = typeFn

	// Time (C uses argv for delta/timestamp/program; scalars for clock/wall/time/sleep).
	functions["KODA_delta_time"] = argvI64(mod, "koda_delta_time")
	functions["KODA_clock"] = mod.NewFunc("koda_clock", types.I64)
	functions["KODA_timestamp"] = argvI64(mod, "koda_timestamp")
	kodaTimeFn := mod.NewFunc("koda_time", types.I64)
	functions["KODA_time"] = kodaTimeFn
	functions["KODA_program_time"] = argvI64(mod, "koda_program_time")
	functions["KODA_wall_time"] = mod.NewFunc("koda_wall_time", types.I64)
	functions["KODA_sleep"] = mod.NewFunc("koda_sleep", types.Void,
		ir.NewParam("ms", types.I64))

	// Random
	functions["KODA_random"] = argvI64(mod, "koda_random")
	functions["KODA_randomInt"] = argvI64(mod, "koda_randomInt")
	functions["KODA_randomChoice"] = argvI64(mod, "koda_randomChoice")
	functions["KODA_randomSeed"] = argvI64(mod, "koda_randomSeed")

	// Math (C implementations are argv-style)
	functions["KODA_lerp"] = argvI64(mod, "koda_lerp")
	functions["KODA_clamp"] = argvI64(mod, "koda_clamp")
	functions["KODA_distance"] = argvI64(mod, "koda_distance")
	functions["KODA_angleBetween"] = argvI64(mod, "koda_angleBetween")
	functions["KODA_map"] = argvI64(mod, "koda_map")

	functions["KODA_pi"] = argvI64(mod, "koda_pi")
	functions["KODA_e"] = argvI64(mod, "koda_e")

	functions["KODA_sin"] = argvI64(mod, "koda_sin")
	functions["KODA_cos"] = argvI64(mod, "koda_cos")
	functions["KODA_tan"] = argvI64(mod, "koda_tan")
	functions["KODA_asin"] = argvI64(mod, "koda_asin")
	functions["KODA_acos"] = argvI64(mod, "koda_acos")
	functions["KODA_atan"] = argvI64(mod, "koda_atan")
	functions["KODA_atan2"] = argvI64(mod, "koda_atan2")

	functions["KODA_pow"] = argvI64(mod, "koda_pow")
	functions["KODA_exp"] = argvI64(mod, "koda_exp")
	functions["KODA_log"] = argvI64(mod, "koda_log")
	functions["KODA_log10"] = argvI64(mod, "koda_log10")
	functions["KODA_log2"] = argvI64(mod, "koda_log2")

	functions["KODA_floor"] = argvI64(mod, "koda_floor")
	functions["KODA_ceil"] = argvI64(mod, "koda_ceil")
	functions["KODA_round"] = argvI64(mod, "koda_round")
	functions["KODA_trunc"] = argvI64(mod, "koda_trunc")

	functions["KODA_sign"] = argvI64(mod, "koda_sign")
	functions["KODA_min"] = argvI64(mod, "koda_min")
	functions["KODA_max"] = argvI64(mod, "koda_max")
	functions["KODA_smoothstep"] = argvI64(mod, "koda_smoothstep")
	functions["KODA_distanceSq"] = argvI64(mod, "koda_distanceSq")
	functions["KODA_normalize"] = argvI64(mod, "koda_normalize")
	functions["KODA_hypot"] = argvI64(mod, "koda_hypot")
	functions["KODA_fmod"] = argvI64(mod, "koda_fmod")
	functions["KODA_degrees"] = argvI64(mod, "koda_degrees")
	functions["KODA_radians"] = argvI64(mod, "koda_radians")
	functions["KODA_wrap"] = argvI64(mod, "koda_wrap")
	functions["KODA_approach"] = argvI64(mod, "koda_approach")
	functions["KODA_smoothdamp"] = argvI64(mod, "koda_smoothdamp")

	// Type checks (argv in C)
	functions["KODA_isNumber"] = argvI64(mod, "koda_isNumber")
	functions["KODA_isString"] = argvI64(mod, "koda_isString")
	functions["KODA_isBool"] = argvI64(mod, "koda_isBool")
	functions["KODA_isNull"] = argvI64(mod, "koda_isNull")
	functions["KODA_isArray"] = argvI64(mod, "koda_isArray")
	functions["KODA_isObject"] = argvI64(mod, "koda_isObject")
	functions["KODA_isFunction"] = argvI64(mod, "koda_isFunction")

	functions["KODA_bool"] = argvI64(mod, "koda_bool")

	functions["KODA_format"] = argvI64(mod, "koda_format")

	functions["KODA_array_map"] = argvI64(mod, "koda_array_map")
	functions["KODA_array_filter"] = argvI64(mod, "koda_array_filter")
	functions["KODA_array_forEach"] = argvI64(mod, "koda_array_forEach")
	functions["KODA_array_find"] = argvI64(mod, "koda_array_find")
	functions["KODA_array_findIndex"] = argvI64(mod, "koda_array_findIndex")
	functions["KODA_array_some"] = argvI64(mod, "koda_array_some")
	functions["KODA_array_every"] = argvI64(mod, "koda_array_every")
	functions["KODA_array_reduce"] = argvI64(mod, "koda_array_reduce")
	functions["KODA_array_sort"] = argvI64(mod, "koda_array_sort")
	functions["KODA_array_reverse"] = argvI64(mod, "koda_array_reverse")
	functions["KODA_array_indexOf"] = argvI64(mod, "koda_array_indexOf")
	functions["KODA_array_includes"] = argvI64(mod, "koda_array_includes")
	functions["KODA_array_slice"] = argvI64(mod, "koda_array_slice")
	functions["KODA_array_concat"] = argvI64(mod, "koda_array_concat")
	functions["KODA_array_join"] = argvI64(mod, "koda_array_join")

	functions["KODA_string_split"] = argvI64(mod, "koda_string_split")
	functions["KODA_string_trim"] = argvI64(mod, "koda_string_trim")
	functions["KODA_string_upper"] = argvI64(mod, "koda_string_upper")
	functions["KODA_string_lower"] = argvI64(mod, "koda_string_lower")
	functions["KODA_string_startsWith"] = argvI64(mod, "koda_string_startsWith")
	functions["KODA_string_endsWith"] = argvI64(mod, "koda_string_endsWith")
	functions["KODA_string_indexOf"] = argvI64(mod, "koda_string_indexOf")
	functions["KODA_string_slice"] = argvI64(mod, "koda_string_slice")
	functions["KODA_string_replace"] = argvI64(mod, "koda_string_replace")
	functions["KODA_string_replaceAll"] = argvI64(mod, "koda_string_replaceAll")

	functions["KODA_readFile"] = argvI64(mod, "koda_readFile")
	functions["KODA_writeFile"] = argvI64(mod, "koda_writeFile")
	functions["KODA_appendFile"] = argvI64(mod, "koda_appendFile")
	functions["KODA_fileExists"] = argvI64(mod, "koda_fileExists")
	functions["KODA_deleteFile"] = argvI64(mod, "koda_deleteFile")
	functions["KODA_isFile"] = argvI64(mod, "koda_is_file")
	functions["KODA_isDir"] = argvI64(mod, "koda_is_dir")
	functions["KODA_fileSize"] = argvI64(mod, "koda_file_size")
	functions["KODA_listDir"] = argvI64(mod, "koda_list_dir")
	functions["KODA_keys"] = argvI64(mod, "koda_keys")

	functions["KODA_ok"] = argvI64(mod, "koda_ok")
	functions["KODA_err"] = argvI64(mod, "koda_err")
	functions["KODA_values_equal"] = mod.NewFunc("koda_values_equal", types.I64,
		ir.NewParam("a", types.I64), ir.NewParam("b", types.I64))
	functions["KODA_panic"] = mod.NewFunc("koda_panic", types.Void,
		ir.NewParam("arg_count", types.I32),
		ir.NewParam("args", types.NewPointer(types.I64)))
	functions["KODA_assert"] = argvI64(mod, "koda_assert")
	functions["KODA_assert_llvm"] = mod.NewFunc("koda_assert_llvm", types.Void,
		ir.NewParam("cond", types.I64),
		ir.NewParam("msg", types.I64))
	functions["KODA_trace"] = argvI64(mod, "koda_trace")

	functions["KODA_parseJSON"] = argvI64(mod, "koda_parseJSON")
	functions["KODA_jsonParse"] = argvI64(mod, "koda_json_parse")
	functions["KODA_jsonTryParse"] = argvI64(mod, "koda_json_try_parse")
	functions["KODA_toJSON"] = argvI64(mod, "koda_toJSON")

	functions["KODA_allocate_object"] = mod.NewFunc("koda_allocate_object", types.I64,
		ir.NewParam("property_count", types.I32))
	functions["KODA_allocate_struct"] = mod.NewFunc("koda_allocate_struct", types.I64,
		ir.NewParam("field_count", types.I32))
	functions["KODA_struct_get"] = mod.NewFunc("koda_struct_get", types.I64,
		ir.NewParam("obj", types.I64),
		ir.NewParam("index", types.I64))
	functions["KODA_struct_set"] = mod.NewFunc("koda_struct_set", types.I64,
		ir.NewParam("obj", types.I64),
		ir.NewParam("index", types.I64),
		ir.NewParam("value", types.I64))

	functions["KODA_unbox_number"] = mod.NewFunc("koda_unbox_number", types.Double,
		ir.NewParam("v", types.I64))
	functions["KODA_box_number"] = mod.NewFunc("koda_box_number", types.I64,
		ir.NewParam("d", types.Double))

	getFn := mod.NewFunc("koda_get", types.I64,
		ir.NewParam("obj", types.I64),
		ir.NewParam("key", types.I64))
	functions["KODA_get"] = getFn
	functions["KODA_object_get"] = getFn
	functions["KODA_array_get"] = getFn

	functions["KODA_set"] = mod.NewFunc("koda_set", types.I64,
		ir.NewParam("obj", types.I64),
		ir.NewParam("key", types.I64),
		ir.NewParam("val", types.I64))

	functions["KODA_object_set"] = mod.NewFunc("koda_object_set", types.I64,
		ir.NewParam("obj", types.I64),
		ir.NewParam("key", types.I64),
		ir.NewParam("value", types.I64))
	functions["KODA_object_remove"] = mod.NewFunc("koda_object_remove", types.I64,
		ir.NewParam("obj", types.I64),
		ir.NewParam("key", types.I64))

	functions["KODA_allocate_string"] = mod.NewFunc("koda_allocate_string", types.I64,
		ir.NewParam("length", types.I32),
		ir.NewParam("chars", types.NewPointer(types.I8)))

	functions["KODA_allocate_array"] = mod.NewFunc("koda_allocate_array", types.I64,
		ir.NewParam("length", types.I32))
	functions["KODA_array_set"] = mod.NewFunc("koda_array_set", types.Void,
		ir.NewParam("arr", types.I64),
		ir.NewParam("index", types.I64),
		ir.NewParam("value", types.I64))
	functions["KODA_array_push"] = mod.NewFunc("koda_array_push", types.Void,
		ir.NewParam("arr", types.I64),
		ir.NewParam("value", types.I64))
	functions["KODA_array_pop"] = mod.NewFunc("koda_array_pop", types.I64,
		ir.NewParam("arr", types.I64))
	lenFn := mod.NewFunc("koda_len", types.I64,
		ir.NewParam("value", types.I64))
	functions["KODA_len"] = lenFn
	arrayLenFn := mod.NewFunc("koda_array_length", types.I64,
		ir.NewParam("value", types.I64))
	functions["KODA_array_length"] = arrayLenFn

	functions["KODA_forof_length"] = mod.NewFunc("koda_forof_length", types.I64,
		ir.NewParam("v", types.I64))
	functions["KODA_forof_key_at"] = mod.NewFunc("koda_forof_key_at", types.I64,
		ir.NewParam("v", types.I64),
		ir.NewParam("idx", types.I64))
	functions["KODA_forof_value_at"] = mod.NewFunc("koda_forof_value_at", types.I64,
		ir.NewParam("v", types.I64),
		ir.NewParam("idx", types.I64))

	functions["KODA_abs"] = mod.NewFunc("koda_abs", types.I64,
		ir.NewParam("value", types.I64))
	functions["KODA_sqrt"] = mod.NewFunc("koda_sqrt", types.I64,
		ir.NewParam("value", types.I64))
	functions["KODA_cbrt"] = mod.NewFunc("koda_cbrt", types.I64,
		ir.NewParam("value", types.I64))
	functions["KODA_number"] = mod.NewFunc("koda_number", types.I64,
		ir.NewParam("value", types.I64))
	functions["KODA_string"] = mod.NewFunc("koda_string", types.I64,
		ir.NewParam("value", types.I64))
	functions["KODA_string_concat"] = mod.NewFunc("koda_string_concat", types.I64,
		ir.NewParam("a", types.I64),
		ir.NewParam("b", types.I64))

	functions["KODA_range"] = argvI64(mod, "koda_range")
	functions["KODA_matches"] = argvI64(mod, "koda_matches")

	cellPtr := types.NewPointer(types.I64)
	functions["KODA_alloc_cell"] = mod.NewFunc("koda_alloc_cell", cellPtr)
	functions["KODA_cell_read"] = mod.NewFunc("koda_cell_read", types.I64,
		ir.NewParam("cell", cellPtr))
	functions["KODA_cell_write"] = mod.NewFunc("koda_cell_write", types.Void,
		ir.NewParam("cell", cellPtr),
		ir.NewParam("val", types.I64))
	functions["KODA_gc_collect"] = mod.NewFunc("koda_gc_collect", types.Void)
	functions["KODA_gc_disable"] = mod.NewFunc("koda_gc_disable", types.Void)
	functions["KODA_gc_enable"] = mod.NewFunc("koda_gc_enable", types.Void)
	functions["KODA_gc_frame_step"] = mod.NewFunc("koda_gc_frame_step", types.Void,
		ir.NewParam("budget_ms", types.Double))
	functions["KODA_gc_stats"] = argvI64(mod, "koda_gc_stats")
	functions["KODA_arena"] = argvI64(mod, "koda_arena")
	functions["KODA_arena_reset"] = argvI64(mod, "koda_arena_reset")
	functions["KODA_arena_alloc_array"] = argvI64(mod, "koda_arena_alloc_array")
	functions["KODA_arena_alloc_struct"] = argvI64(mod, "koda_arena_alloc_struct")

	ptrPtr := types.NewPointer(types.NewPointer(types.I64))
	functions["KODA_push_frame"] = mod.NewFunc("koda_push_frame", types.Void,
		ir.NewParam("slots", ptrPtr),
		ir.NewParam("count", types.I32))
	functions["KODA_pop_frame"] = mod.NewFunc("koda_pop_frame", types.Void)

	functions["KODA_push_call"] = mod.NewFunc("koda_push_call", types.Void,
		ir.NewParam("fn_name", types.NewPointer(types.I8)),
		ir.NewParam("file_name", types.NewPointer(types.I8)),
		ir.NewParam("line", types.I32))
	functions["KODA_pop_call"] = mod.NewFunc("koda_pop_call", types.Void)

	// C malloc for packed closure environments (linker resolves from CRT).
	functions["malloc"] = mod.NewFunc("malloc", types.NewPointer(types.I8), ir.NewParam("size", types.I64))

	return functions
}
