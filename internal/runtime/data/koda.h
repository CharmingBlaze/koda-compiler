#ifndef KODA_H
#define KODA_H

#include <stdint.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>

/**
 * KODA RUNTIME HEADER
 * -------------------
 * Optimized for performance using NaN-boxing and a generational GC.
 */

typedef uint64_t KodaValue;

#define QNAN     ((uint64_t)0x7ffc000000000000)
#define SIGN_BIT ((uint64_t)0x8000000000000000)

#define TAG_NULL   1
#define TAG_FALSE  2
#define TAG_TRUE   3

// --- Value Encoding/Decoding ---

static inline double value_to_num(KodaValue v) {
    union { uint64_t u; double d; } cast;
    cast.u = v;
    return cast.d;
}

static inline KodaValue num_to_value(double n) {
    union { uint64_t u; double d; } cast;
    cast.d = n;
    return cast.u;
}

#define IS_NUMBER(v) (((v) & QNAN) != QNAN)
#define IS_NULL(v)   ((v) == (QNAN | TAG_NULL))
#define IS_BOOL(v)   (((v) | 1) == (QNAN | TAG_TRUE))
#define IS_OBJ(v)    (((v) & (QNAN | SIGN_BIT)) == (QNAN | SIGN_BIT))

#define AS_NUMBER(v) value_to_num(v)
#define AS_BOOL(v)   ((v) == (QNAN | TAG_TRUE))
#define AS_OBJ(v)    ((struct KodaObj*)(uintptr_t)((v) & ~(SIGN_BIT | QNAN)))
#define AS_CLOSURE(v) ((KodaClosure*)AS_OBJ(v))

#define NUMBER_VAL(n) num_to_value(n)
#define BOOL_VAL(b)   ((b) ? (QNAN | TAG_TRUE) : (QNAN | TAG_FALSE))
#define NULL_VAL      (QNAN | TAG_NULL)
#define OBJ_VAL(o)    (KodaValue)(SIGN_BIT | QNAN | (uintptr_t)(o))

// --- Object System ---

typedef enum {
    OBJ_STRING,
    OBJ_ARRAY,
    OBJ_OBJECT,
    OBJ_MAP,
    OBJ_SET,
    OBJ_TUPLE,
    OBJ_CLOSURE,
    OBJ_UPVALUE,
    OBJ_NATIVE,
    OBJ_BOUND_METHOD,
    OBJ_FFI,
} KodaObjType;

typedef struct KodaObj {
    KodaObjType type;
    uint8_t generation; // 0=young, 1=old
    uint8_t mark;       // 0=white, 1=gray/black
    uint8_t age;        // for promotion
    struct KodaObj* next;
} KodaObj;

typedef struct {
    KodaObj header;
    int length;
    uint32_t hash;
    char chars[];
} KodaString;

typedef struct {
    KodaObj header;
    int count;
    int capacity;
    KodaValue* elements;
} KodaArray;

typedef struct KodaUpvalue {
    KodaObj header;
    KodaValue* location;
    KodaValue closed;
    struct KodaUpvalue* next;
} KodaUpvalue;

typedef KodaValue (*KodaFn)(KodaValue thisVal, int argCount, KodaValue* args, KodaUpvalue** upvalues);

typedef struct {
    KodaObj header;
    KodaFn fn;
    KodaUpvalue** upvalues;
    int upvalueCount;
} KodaClosure;

// FFI wrapper
typedef void (*Finalizer)(void*);
typedef struct {
    KodaObj header;
    void* ptr;
    Finalizer finalizer;
    int ref_count;
} KodaFFI;
typedef KodaValue (*KodaNativeFn)(int argCount, KodaValue* args);

typedef struct {
    KodaObj header;
    KodaNativeFn fn;
} KodaNative;

typedef struct {
    KodaValue key;
    KodaValue value;
} KodaEntry;

typedef struct {
    KodaObj header;
    int count;
    int capacity;
    KodaEntry* entries;
} KodaObject;

typedef struct {
    KodaObj header;
    KodaValue receiver;
    KodaValue method; // Can be a Closure or NativeFn
} KodaBoundMethod;

// --- GC & Memory ---

void koda_init(void);
void koda_shutdown(void);
void koda_runtime_init(void);
void koda_runtime_shutdown(void);
KodaObj* koda_alloc(KodaObjType type, size_t size);
KodaString* koda_copy_string(const char* chars, int length);
/** Boxed string for compiler string literals (length, UTF-8 bytes). */
KodaValue koda_allocate_string(int length, const char* chars);
KodaArray* koda_new_array();
KodaObject* koda_alloc_object();
KodaUpvalue* koda_new_upvalue(KodaValue* slot);
KodaClosure* koda_new_closure(KodaFn fn, int upvalueCount);

void koda_upvalue_set(KodaUpvalue* up, KodaValue val);
KodaValue koda_upvalue_get(KodaUpvalue* up);
void koda_closure_set_upvalue(KodaValue closure, int index, KodaUpvalue* u);
KodaValue koda_obj_as_value(void* obj);

/** Fatal arity mismatch (stderr + exit). kind: 0=min, 1=max, 2=exact, 3=missing default index (b unused). */
void koda_assert(KodaValue cond, KodaValue msg);
void koda_abort_arg_error(int kind, int a, int b);

/** Build a Koda array from argv[start : start+count). count may be 0 (argv may be NULL). */
KodaValue koda_argv_slice_to_array(KodaValue* argv, int start, int count);

void koda_array_push(KodaValue array, KodaValue value);
void koda_object_set(KodaValue obj, KodaValue key, KodaValue value);
/** Remove entry whose key compares equal; returns true if a key was removed. */
KodaValue koda_object_delete(KodaValue obj, KodaValue key);
KodaValue koda_map_new(int argCount, KodaValue* args);
KodaValue koda_set_new(int argCount, KodaValue* args);
KodaValue koda_tuple_new(int argCount, KodaValue* args);
KodaValue koda_input(int argCount, KodaValue* args);
KodaValue koda_get_index(KodaValue obj, KodaValue index);
KodaValue koda_call(KodaValue callee, int argCount, KodaValue* args);

KodaValue koda_len(KodaValue val);
KodaValue koda_type(KodaValue val);
KodaValue koda_clock();
/** Wall-clock seconds since Unix epoch (maps to `time()` in the language reference). */
KodaValue koda_wall_time(void);
/** Sleep for a wall duration in milliseconds (NaN-boxed number). */
void koda_sleep_ms(KodaValue ms);
KodaValue koda_abs_num(KodaValue v);
KodaValue koda_sqrt_num(KodaValue v);
KodaValue koda_cbrt_num(KodaValue v);
KodaValue koda_sin_num(KodaValue v);
KodaValue koda_cos_num(KodaValue v);
KodaValue koda_tan_num(KodaValue v);
KodaValue koda_asin_num(KodaValue v);
KodaValue koda_acos_num(KodaValue v);
KodaValue koda_atan_num(KodaValue v);
KodaValue koda_atan2_num(KodaValue y, KodaValue x);
KodaValue koda_log_num(KodaValue v);
KodaValue koda_log2_num(KodaValue v);
KodaValue koda_log10_num(KodaValue v);
KodaValue koda_exp_num(KodaValue v);
KodaValue koda_floor_num(KodaValue v);
KodaValue koda_ceil_num(KodaValue v);
KodaValue koda_round_num(KodaValue v);
KodaValue koda_trunc_num(KodaValue v);
KodaValue koda_min_num(int argCount, KodaValue* args);
KodaValue koda_max_num(int argCount, KodaValue* args);
KodaValue koda_clamp_num(KodaValue v, KodaValue min, KodaValue max);
KodaValue koda_lerp_num(KodaValue a, KodaValue b, KodaValue t);
KodaValue koda_smoothstep_num(KodaValue a, KodaValue b, KodaValue t);
KodaValue koda_map_num(KodaValue v, KodaValue inMin, KodaValue inMax, KodaValue outMin, KodaValue outMax);
KodaValue koda_sign_num(KodaValue v);
KodaValue koda_hypot_num(KodaValue a, KodaValue b);
KodaValue koda_distance_num(KodaValue x1, KodaValue y1, KodaValue x2, KodaValue y2);
KodaValue koda_distance_sq_num(KodaValue x1, KodaValue y1, KodaValue x2, KodaValue y2);
KodaValue koda_angle_between_num(KodaValue x1, KodaValue y1, KodaValue x2, KodaValue y2);
KodaValue koda_normalize_num(KodaValue x, KodaValue y);

/** Time functions */
KodaValue koda_delta_time(void);
KodaValue koda_time(void);
KodaValue koda_timestamp(void);
void koda_sleep_ms(KodaValue ms);

/** Random functions */
KodaValue koda_random(int argCount, KodaValue* args);
KodaValue koda_random_int(int argCount, KodaValue* args);
KodaValue koda_random_choice(KodaValue array);
KodaValue koda_set_timeout(KodaValue closure, KodaValue ms);
KodaValue koda_set_interval(KodaValue closure, KodaValue ms);
void koda_poll_tasks(void);

KodaValue koda_io_read_file(int argCount, KodaValue* args);
KodaValue koda_io_write_file(int argCount, KodaValue* args);
KodaValue koda_module_init_io();
KodaValue koda_json_parse(int argCount, KodaValue* args);
KodaValue koda_json_stringify(int argCount, KodaValue* args);
KodaValue koda_module_init_json();

void koda_push(KodaValue value);
KodaValue koda_pop();
void koda_collect();

KodaValue koda_add(KodaValue a, KodaValue b);
KodaValue koda_to_string_val(KodaValue v);
/** Parse a number from a string, or pass through numeric values; else null. */
KodaValue koda_parse_number(KodaValue v);
KodaValue koda_range(KodaValue from, KodaValue to);
KodaValue koda_sub(KodaValue a, KodaValue b);
KodaValue koda_mul(KodaValue a, KodaValue b);
KodaValue koda_div(KodaValue a, KodaValue b);
KodaValue koda_print(KodaValue val);
void koda_print_no_newline(KodaValue val);
void koda_print_space(void);
void koda_print_newline(void);

KodaValue koda_is_truthy(KodaValue v);
KodaValue koda_gt(KodaValue a, KodaValue b);
KodaValue koda_lt(KodaValue a, KodaValue b);
KodaValue koda_eq(KodaValue a, KodaValue b);
KodaValue koda_neq(KodaValue a, KodaValue b);
KodaValue koda_le(KodaValue a, KodaValue b);
KodaValue koda_ge(KodaValue a, KodaValue b);
KodaValue koda_mod(KodaValue a, KodaValue b);
KodaValue koda_pow(KodaValue a, KodaValue b);
KodaValue koda_bit_and(KodaValue a, KodaValue b);
KodaValue koda_bit_or(KodaValue a, KodaValue b);
KodaValue koda_bit_xor(KodaValue a, KodaValue b);
KodaValue koda_bit_not(KodaValue a);
KodaValue koda_shl(KodaValue a, KodaValue b);
KodaValue koda_shr(KodaValue a, KodaValue b);
KodaValue koda_ushr(KodaValue a, KodaValue b);
void koda_set_index(KodaValue obj, KodaValue index, KodaValue value);
int koda_for_in_len(KodaValue iterable);
KodaValue koda_for_in_get(KodaValue iterable, int index);
KodaValue koda_slice(KodaValue obj, KodaValue start, KodaValue end);

KodaValue koda_negate(KodaValue a);

// Write Barrier
void koda_write_barrier(KodaObj* obj, KodaValue value);

#define KODA_WRITE_BARRIER(obj, val) \
    do { if (IS_OBJ(val)) koda_write_barrier((KodaObj*)(obj), (val)); } while(0)

#endif // KODA_H
