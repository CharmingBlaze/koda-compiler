#ifndef KODA_RUNTIME_H
#define KODA_RUNTIME_H

#include "value.h"
#include "object.h"
#include "gc.h"
#include <stddef.h>
#include <stdbool.h>

extern void* gc_stack_base;
extern Value* koda_globals;
extern int koda_globals_count;
extern int koda_globals_capacity;
extern Value** koda_global_slots;
extern int koda_global_slots_count;
extern int koda_global_slots_capacity;

// Runtime initialization
void koda_runtime_set_stack_base(void* base);
void koda_runtime_init_ex(void* stack_base);
void koda_runtime_init(void);
void koda_runtime_set_argv(int argc, char** argv);

// Runtime cleanup
void koda_runtime_shutdown(void);
void koda_globals_init(void);
void koda_globals_free(void);
void koda_register_global(Value v);
void koda_register_global_slot(Value* slot);

void koda_gc_set_threshold(size_t bytes);
void koda_gc_disable(void);
void koda_gc_enable(void);
void koda_gc_collect(void);
void koda_gc_use_shadow_stack(bool enable);
void koda_gc_frame_step(double budget_ms);
Value koda_gc_stats(int argc, Value* argv);
void koda_intern_clear(void);

#include "shadow_stack.h"

Value* koda_alloc_cell(void);
Value koda_cell_read(Value* cell);
void koda_cell_write(Value* cell, Value v);

// Native function declarations
Value koda_print_argv(int arg_count, Value* args);
Value koda_print_val(Value v);
void koda_print_newline(void);
Value koda_typeof(int arg_count, Value* args);
Value koda_get_index(Value obj, Value index);
void koda_assert_llvm(Value cond, Value msg);

/** Result helpers and panic (see language docs). */
Value koda_ok(int argc, Value* argv);
Value koda_err(int argc, Value* argv);
void koda_panic(int argc, Value* argv);
Value koda_assert(int argc, Value* argv);
Value koda_err_str(const char* msg);
void koda_panic_str(const char* msg);

/** Short type name for diagnostics (static buffer — do not free). */
const char* koda_value_type_name(Value v);
/** Format a type error and panic (via [koda_panic_str]). */
void koda_type_error(const char* op, const char* expected, Value got);
void koda_null_error(const char* op);
/** After a full mark pass, drop intern slots for unmarked strings (weak refs); call before gc_sweep. */
void koda_sweep_intern_table(void);

/** Optional call stack for panic stack traces (native codegen). */
void koda_push_call(const char* fn_name, const char* file_name, int line);
void koda_pop_call(void);
void koda_print_stack_trace(void);
Value koda_clock(void);
Value koda_wall_time(void);
Value koda_allocate_string(int length, char* chars);
/** Copy UTF-8 bytes into a string object (wrapgen / FFI helpers). */
Value koda_copy_string(const char* chars, int length);
Value koda_allocate_object(int property_count);
Value koda_allocate_struct(int field_count);
Value koda_struct_get(Value obj, int64_t index);
Value koda_struct_set(Value obj, int64_t index, Value val);
Value koda_object_get(Value obj, Value key);
Value koda_object_set(Value obj, Value key, Value value);
Value koda_object_remove(Value obj, Value key);
Value koda_array_push_argv(int argc, Value* argv);
Value koda_array_pop_argv(int argc, Value* argv);
Value koda_asset_path(int argc, Value* argv);
Value koda_args(int argc, Value* argv);
Value koda_env(int argc, Value* argv);
Value koda_rgb(int argc, Value* argv);
Value koda_rgba(int argc, Value* argv);
Value koda_vec2(int argc, Value* argv);
Value koda_vec3(int argc, Value* argv);
Value koda_rect(int argc, Value* argv);
Value koda_box(int argc, Value* argv);
Value koda_color(int argc, Value* argv);
Value koda_value_add(Value a, Value b);
Value koda_value_sub(Value a, Value b);
Value koda_value_mul(Value a, Value b);

/** Native lowering for `for-in` / `for-of` over arrays and tables (slot order). */
Value koda_forof_length(Value v);
Value koda_forof_key_at(Value v, Value idx_val);
Value koda_forof_value_at(Value v, Value idx_val);

/** NaN-box helpers for LLVM codegen — never bitcast i64 to double without these. */
double koda_unbox_number(Value v);
Value koda_box_number(double d);

/** Unified index read / write (arrays, strings, tables). */
Value koda_get(Value obj, Value key);
Value koda_set(Value obj, Value key, Value value);
int koda_get_shadow_depth(void);
/** High-water shadow stack depth since runtime init (for diagnostics). */
int koda_shadow_stack_high_water(void);

// Bool value helper (not in value.h, needed for wrapper)
static inline Value BOOL_VAL(bool b) {
    return b ? TRUE_VAL : FALSE_VAL;
}

#endif // KODA_RUNTIME_H
