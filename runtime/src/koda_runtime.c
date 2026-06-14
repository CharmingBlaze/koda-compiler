#include "koda_runtime.h"
#include "value.h"
#include "object.h"
#include "gc.h"
#include "shadow_stack.h"
#include <stdbool.h>
#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <math.h>
#include <stdint.h>
#include <errno.h>
#include <ctype.h>

// Platform-specific sleep and monotonic clock
#ifdef _WIN32
#include <synchapi.h>
#include <windows.h>
#include <io.h>
#include <direct.h>
#define KODA_STAT _stat
#define KODA_STAT_FUNC _stat
#else
#include <fcntl.h>
#include <unistd.h>
#include <dirent.h>
#include <sys/stat.h>
#define KODA_STAT stat
#define KODA_STAT_FUNC stat
#endif

#ifdef _WIN32
static double koda_monotonic_seconds(void) {
    return (double)GetTickCount64() * 0.001;
}
#else
static double koda_monotonic_seconds(void) {
    struct timespec ts;
    if (clock_gettime(CLOCK_MONOTONIC, &ts) != 0) {
        return (double)time(NULL);
    }
    return (double)ts.tv_sec + (double)ts.tv_nsec * 1e-9;
}
#endif

void* gc_stack_base = NULL;

Value koda_bool(int argc, Value* args);

KodaShadowFrame* koda_shadow_stack = NULL;
int koda_shadow_depth = 0;
int koda_shadow_capacity = 0;
static int koda_shadow_depth_high_water = 0;
int koda_shadow_stack_max_capacity = KODA_SHADOW_STACK_MAX_CAPACITY;

static void koda_shadow_stack_configure_from_env(void) {
    const char* env = getenv("KODA_STACK_DEPTH");
    if (env == NULL || env[0] == '\0') {
        return;
    }
    char* end = NULL;
    long depth = strtol(env, &end, 10);
    if (end == env || depth < 256 || depth > 1048576) {
        fprintf(stderr, "koda: ignoring invalid KODA_STACK_DEPTH=%s (use 256..1048576)\n", env);
        return;
    }
    koda_shadow_stack_max_capacity = (int)depth;
}

typedef struct {
    ObjString* str;
} InternBucket;

typedef struct {
    InternBucket* buckets;
    int capacity;
    int count;
} KodaInternTable;

static KodaInternTable koda_string_intern = { NULL, 0, 0 };

static int koda_gc_debug_enabled(void) {
    const char* env = getenv("KODA_GC_DEBUG");
    return env != NULL && env[0] != '\0' && env[0] != '0';
}

static int koda_stack_base_plausible(void* stack_base) {
    if (stack_base == NULL) {
        return 0;
    }
    uintptr_t a = (uintptr_t)stack_base;
    uintptr_t b = (uintptr_t)&a;
    uintptr_t dist = (a > b) ? (a - b) : (b - a);
    return dist <= ((uintptr_t)1u << 30);
}

static void koda_shadow_stack_init(void) {
    if (koda_shadow_stack != NULL) {
        return;
    }
    koda_shadow_capacity = KODA_SHADOW_STACK_INITIAL_CAPACITY;
    koda_shadow_stack = (KodaShadowFrame*)malloc(sizeof(KodaShadowFrame) * (size_t)koda_shadow_capacity);
    if (koda_shadow_stack == NULL) {
        koda_panic_str("out of memory allocating shadow stack");
    }
}

void koda_push_frame(Value** slot_ptrs, int count) {
    koda_shadow_stack_init();
    if (koda_shadow_depth >= koda_shadow_capacity) {
        if (koda_shadow_capacity >= koda_shadow_stack_max_capacity) {
            koda_panic_str("stack overflow — maximum recursion depth reached");
        }
        int next_capacity = koda_shadow_capacity * 2;
        if (next_capacity > koda_shadow_stack_max_capacity) {
            next_capacity = koda_shadow_stack_max_capacity;
        }
        KodaShadowFrame* next = (KodaShadowFrame*)realloc(koda_shadow_stack, sizeof(KodaShadowFrame) * (size_t)next_capacity);
        if (next == NULL) {
            koda_panic_str("out of memory growing shadow stack");
        }
        koda_shadow_stack = next;
        koda_shadow_capacity = next_capacity;
    }
    koda_shadow_stack[koda_shadow_depth].slot_ptrs = slot_ptrs;
    koda_shadow_stack[koda_shadow_depth].count = count;
    koda_shadow_depth++;
    if (koda_shadow_depth > koda_shadow_depth_high_water) {
        koda_shadow_depth_high_water = koda_shadow_depth;
    }
}

void koda_pop_frame(void) {
    if (koda_shadow_depth <= 0) {
        fprintf(stderr, "koda: shadow stack underflow (push/pop mismatch)\n");
        abort();
    }
    koda_shadow_depth--;
}

int koda_get_shadow_depth(void) {
    return koda_shadow_depth;
}

int koda_shadow_stack_high_water(void) {
    return koda_shadow_depth_high_water;
}

#define KODA_CALL_STACK_MAX 256

typedef struct {
    const char* function_name;
    const char* file_name;
    int line;
} KodaCallFrame;

static KodaCallFrame koda_call_stack[KODA_CALL_STACK_MAX];
static int koda_call_stack_depth = 0;

static uint32_t fnv1a_bytes(const uint8_t* data, int len);
static uint32_t koda_string_hash(ObjString* s);
static void koda_intern_add(ObjString* str);

static int intern_probe_index(uint32_t hash, int capacity, int probe) {
    return (int)((hash + (uint32_t)probe) & (uint32_t)(capacity - 1));
}

static void koda_intern_table_grow(void) {
    int old_cap = koda_string_intern.capacity;
    InternBucket* old_buckets = koda_string_intern.buckets;
    int new_cap = old_cap == 0 ? 256 : old_cap * 2;
    InternBucket* next = (InternBucket*)calloc((size_t)new_cap, sizeof(InternBucket));
    if (next == NULL) {
        koda_panic_str("out of memory growing string intern table");
    }
    koda_string_intern.buckets = next;
    koda_string_intern.capacity = new_cap;
    koda_string_intern.count = 0;
    if (old_buckets != NULL) {
        for (int i = 0; i < old_cap; i++) {
            if (old_buckets[i].str != NULL) {
                koda_intern_add(old_buckets[i].str);
            }
        }
        free(old_buckets);
    }
}

static ObjString* koda_intern_find(const char* chars, int length) {
    if (chars == NULL || length < 0 || koda_string_intern.capacity == 0) {
        return NULL;
    }
    uint32_t hash = fnv1a_bytes((const uint8_t*)chars, length);
    int cap = koda_string_intern.capacity;
    for (int probe = 0; probe < cap; probe++) {
        int i = intern_probe_index(hash, cap, probe);
        ObjString* s = koda_string_intern.buckets[i].str;
        if (s == NULL) {
            return NULL;
        }
        if (s->length == length && memcmp(s->chars, chars, (size_t)length) == 0) {
            return s;
        }
    }
    return NULL;
}

static void koda_intern_add(ObjString* str) {
    if (str == NULL) {
        return;
    }
    if (koda_string_intern.capacity == 0) {
        koda_intern_table_grow();
    }
    if (koda_string_intern.count * 4 >= koda_string_intern.capacity * 3) {
        koda_intern_table_grow();
    }
    uint32_t hash = koda_string_hash(str);
    int cap = koda_string_intern.capacity;
    for (int probe = 0; probe < cap; probe++) {
        int i = intern_probe_index(hash, cap, probe);
        if (koda_string_intern.buckets[i].str == NULL) {
            koda_string_intern.buckets[i].str = str;
            koda_string_intern.count++;
            return;
        }
        ObjString* existing = koda_string_intern.buckets[i].str;
        if (existing->length == str->length &&
            memcmp(existing->chars, str->chars, (size_t)str->length) == 0) {
            return;
        }
    }
    koda_intern_table_grow();
    koda_intern_add(str);
}

void koda_sweep_intern_table(void) {
    if (koda_string_intern.capacity == 0) {
        return;
    }
    int cap = koda_string_intern.capacity;
    InternBucket* survivors = (InternBucket*)calloc((size_t)cap, sizeof(InternBucket));
    if (survivors == NULL) {
        koda_panic_str("out of memory sweeping string intern table");
    }
    int kept = 0;
    for (int i = 0; i < cap; i++) {
        ObjString* s = koda_string_intern.buckets[i].str;
        if (s == NULL) {
            continue;
        }
        if (!s->obj.is_marked) {
            continue;
        }
        uint32_t hash = koda_string_hash(s);
        for (int probe = 0; probe < cap; probe++) {
            int slot = intern_probe_index(hash, cap, probe);
            if (survivors[slot].str == NULL) {
                survivors[slot].str = s;
                kept++;
                break;
            }
        }
    }
    free(koda_string_intern.buckets);
    koda_string_intern.buckets = survivors;
    koda_string_intern.count = kept;
}

void koda_push_call(const char* fn_name, const char* file_name, int line) {
    if (koda_call_stack_depth < KODA_CALL_STACK_MAX) {
        koda_call_stack[koda_call_stack_depth].function_name = fn_name ? fn_name : "?";
        koda_call_stack[koda_call_stack_depth].file_name = file_name ? file_name : "?";
        koda_call_stack[koda_call_stack_depth].line = line;
        koda_call_stack_depth++;
    }
}

void koda_pop_call(void) {
    if (koda_call_stack_depth > 0) {
        koda_call_stack_depth--;
    }
}

void koda_print_stack_trace(void) {
    fprintf(stderr, "\nstack trace:\n");
    for (int i = koda_call_stack_depth - 1; i >= 0; i--) {
        fprintf(stderr, "  at %s (%s:%d)\n",
            koda_call_stack[i].function_name,
            koda_call_stack[i].file_name,
            koda_call_stack[i].line);
    }
    fflush(stderr);
}

static void koda_print_value_stderr(Value v) {
    if (IS_NIL(v)) {
        fprintf(stderr, "nil");
    } else if (IS_FALSE(v)) {
        fprintf(stderr, "false");
    } else if (IS_TRUE(v)) {
        fprintf(stderr, "true");
    } else if (IS_NUMBER(v)) {
        fprintf(stderr, "%g", AS_NUMBER(v));
    } else if (IS_OBJ(v)) {
        Obj* o = AS_OBJ(v);
        if (o->type == OBJ_STRING) {
            ObjString* s = (ObjString*)o;
            fprintf(stderr, "%.*s", s->length, s->chars);
        } else {
            fprintf(stderr, "[object]");
        }
    } else {
        fprintf(stderr, "?");
    }
}

Value koda_ok(int argc, Value* argv) {
    Value val = (argc > 0) ? argv[0] : NIL_VAL;
    Value objv = koda_allocate_object(3);
    Value key_ok = koda_copy_string("ok", 2);
    Value key_value = koda_copy_string("value", 5);
    Value key_error = koda_copy_string("error", 5);
    koda_object_set(objv, key_ok, TRUE_VAL);
    koda_object_set(objv, key_value, val);
    koda_object_set(objv, key_error, NIL_VAL);
    return objv;
}

Value koda_err(int argc, Value* argv) {
    Value msg = (argc > 0) ? argv[0] : koda_copy_string("error", 5);
    Value objv = koda_allocate_object(3);
    Value key_ok = koda_copy_string("ok", 2);
    Value key_value = koda_copy_string("value", 5);
    Value key_error = koda_copy_string("error", 5);
    koda_object_set(objv, key_ok, FALSE_VAL);
    koda_object_set(objv, key_value, NIL_VAL);
    koda_object_set(objv, key_error, msg);
    return objv;
}

Value koda_err_str(const char* msg) {
    if (msg == NULL) {
        msg = "error";
    }
    Value v = koda_copy_string(msg, (int)strlen(msg));
    Value args[1] = { v };
    return koda_err(1, args);
}

void koda_panic_str(const char* msg) {
	if (msg == NULL) {
		msg = "(null message)";
	}
	Value v = koda_copy_string(msg, (int)strlen(msg));
	Value args[1] = { v };
	koda_panic(1, args);
}

const char* koda_value_type_name(Value v) {
	if (IS_NIL(v)) {
		return "null";
	}
	if (IS_BOOL(v)) {
		return "boolean";
	}
	if (IS_NUMBER(v)) {
		return "number";
	}
	if (IS_OBJ(v)) {
		Obj* o = AS_OBJ(v);
		switch (o->type) {
		case OBJ_STRING:
			return "string";
		case OBJ_ARRAY:
			return "array";
		case OBJ_TABLE:
			return "object";
		case OBJ_CLOSURE:
		case OBJ_FUNCTION:
		case OBJ_NATIVE:
			return "function";
		case OBJ_CELL:
			return "cell";
		case OBJ_ARENA:
			return "arena";
		default:
			return "object";
		}
	}
	return "value";
}

void koda_type_error(const char* op, const char* expected, Value got) {
	char buf[512];
	snprintf(buf, sizeof(buf), "type error in '%s': expected %s, got %s",
		op != NULL ? op : "?",
		expected != NULL ? expected : "?",
		koda_value_type_name(got));
	koda_panic_str(buf);
}

void koda_null_error(const char* op) {
	char buf[256];
	snprintf(buf, sizeof(buf), "null reference in '%s'", op != NULL ? op : "?");
	koda_panic_str(buf);
}

void koda_panic(int argc, Value* argv) {
    fprintf(stderr, "\nkoda panic: ");
    if (argc > 0) {
        koda_print_value_stderr(argv[0]);
    } else {
        fprintf(stderr, "(no message)");
    }
    fprintf(stderr, "\n");
    koda_print_stack_trace();
    fflush(stderr);
    exit(1);
}

Value koda_assert(int argc, Value* argv) {
    if (argc < 1) {
        koda_panic_str("assert() requires at least one argument");
    }
    Value b = koda_bool(1, argv);
    if (IS_TRUE(b)) {
        return NIL_VAL;
    }
    if (argc >= 2) {
        Value pargs[1] = { argv[1] };
        koda_panic(1, pargs);
    }
    koda_panic_str("assertion failed");
    return NIL_VAL;
}

void koda_gc_use_shadow_stack(bool enable) {
    gc_set_use_shadow_stack(enable);
}

Value* koda_globals = NULL;
int koda_globals_count = 0;
int koda_globals_capacity = 0;
Value** koda_global_slots = NULL;
int koda_global_slots_count = 0;
int koda_global_slots_capacity = 0;

void koda_mark_module_cache(void) {}

void koda_mark_open_upvalues(void) {}

// xoshiro128** PRNG — seeded from OS entropy on first use (not libc rand()).
static uint32_t koda_rng_state[4];
static int koda_rng_seeded = 0;

static uint64_t koda_splitmix64(uint64_t x) {
    uint64_t z = x + 0x9e3779b97f4a7c15ULL;
    z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9ULL;
    z = (z ^ (z >> 27)) * 0x94d049bb133111ebULL;
    return z ^ (z >> 31);
}

static void koda_rng_seed_u64(uint64_t seed) {
    uint64_t sm = seed;
    koda_rng_state[0] = (uint32_t)koda_splitmix64(sm);
    koda_rng_state[1] = (uint32_t)koda_splitmix64(sm);
    koda_rng_state[2] = (uint32_t)koda_splitmix64(sm);
    koda_rng_state[3] = (uint32_t)koda_splitmix64(sm);
    if (koda_rng_state[0] == 0 && koda_rng_state[1] == 0 && koda_rng_state[2] == 0 && koda_rng_state[3] == 0) {
        koda_rng_state[0] = 1;
    }
    koda_rng_seeded = 1;
}

static void koda_rng_gather_entropy(uint64_t* out) {
    uint32_t a = 0;
    uint32_t b = 0;
#ifdef _WIN32
    unsigned int s = 0;
    a = (uint32_t)GetTickCount64();
    b = (uint32_t)(GetTickCount64() >> 32) ^ (uint32_t)time(NULL);
    (void)s;
#else
    int fd = open("/dev/urandom", O_RDONLY);
    if (fd >= 0) {
        unsigned char buf[8];
        ssize_t n = read(fd, buf, 8);
        close(fd);
        if (n == 8) {
            a = (uint32_t)buf[0] | ((uint32_t)buf[1] << 8) | ((uint32_t)buf[2] << 16) | ((uint32_t)buf[3] << 24);
            b = (uint32_t)buf[4] | ((uint32_t)buf[5] << 8) | ((uint32_t)buf[6] << 16) | ((uint32_t)buf[7] << 24);
        }
    }
#endif
    if (a == 0 && b == 0) {
        a = (uint32_t)time(NULL);
        b = (uint32_t)(koda_monotonic_seconds() * 1000000.0);
    }
    *out = ((uint64_t)b << 32) | (uint64_t)a;
}

static void koda_rng_ensure_seeded(void) {
    if (koda_rng_seeded) {
        return;
    }
    uint64_t entropy = 0;
    koda_rng_gather_entropy(&entropy);
    koda_rng_seed_u64(entropy);
}

static uint32_t koda_rng_rotl32(uint32_t x, int k) {
    return (x << k) | (x >> (32 - k));
}

static uint32_t koda_rng_next_u32(void) {
    koda_rng_ensure_seeded();
    uint32_t result = koda_rng_rotl32(koda_rng_state[1] * 5, 7) * 9;
    uint32_t t = koda_rng_state[1] << 9;
    koda_rng_state[2] ^= koda_rng_state[0];
    koda_rng_state[3] ^= koda_rng_state[1];
    koda_rng_state[1] ^= koda_rng_state[2];
    koda_rng_state[0] ^= t;
    koda_rng_state[2] ^= t << 11;
    koda_rng_state[3] ^= koda_rng_state[3] >> 11;
    return result;
}

static double koda_rng_unit(void) {
    return (koda_rng_next_u32() >> 11) * (1.0 / 9007199254740992.0);
}

static int koda_frame_clock_inited = 0;
static double koda_frame_clock_start = 0;
static double koda_frame_clock_last = 0;

static void koda_frame_clock_ensure_init(void) {
    if (koda_frame_clock_inited) {
        return;
    }
    double now = koda_monotonic_seconds();
    koda_frame_clock_start = now;
    koda_frame_clock_last = now;
    koda_frame_clock_inited = 1;
}

void koda_runtime_set_stack_base(void* base) {
    gc_stack_base = base;
}

void koda_globals_init(void) {
    if (koda_globals != NULL && koda_globals_capacity > 0) {
        return;
    }
    koda_globals_capacity = 256;
    koda_globals_count = 0;
    koda_globals = (Value*)malloc(sizeof(Value) * (size_t)koda_globals_capacity);
    if (koda_globals == NULL) {
        fprintf(stderr, "koda: out of memory allocating globals array\n");
        exit(1);
    }
    koda_global_slots_capacity = 256;
    koda_global_slots_count = 0;
    koda_global_slots = (Value**)malloc(sizeof(Value*) * (size_t)koda_global_slots_capacity);
    if (koda_global_slots == NULL) {
        fprintf(stderr, "koda: out of memory allocating global slots array\n");
        exit(1);
    }
}

void koda_globals_free(void) {
    if (koda_globals != NULL) {
        free(koda_globals);
        koda_globals = NULL;
    }
    if (koda_global_slots != NULL) {
        free(koda_global_slots);
        koda_global_slots = NULL;
    }
    koda_globals_count = 0;
    koda_globals_capacity = 0;
    koda_global_slots_count = 0;
    koda_global_slots_capacity = 0;
}

void koda_register_global(Value v) {
    if (koda_globals == NULL || koda_globals_capacity <= 0) {
        koda_globals_init();
    }
    if (koda_globals_count >= koda_globals_capacity) {
        int new_cap = koda_globals_capacity * 2;
        if (new_cap < 256) {
            new_cap = 256;
        }
        Value* next = (Value*)realloc(koda_globals, sizeof(Value) * (size_t)new_cap);
        if (next == NULL) {
            fprintf(stderr, "koda: out of memory growing globals array\n");
            exit(1);
        }
        koda_globals = next;
        koda_globals_capacity = new_cap;
    }
    koda_globals[koda_globals_count++] = v;
}

void koda_register_global_slot(Value* slot) {
    if (slot == NULL) {
        return;
    }
    if (koda_global_slots == NULL || koda_global_slots_capacity <= 0) {
        koda_globals_init();
    }
    // Top-level declarations inside control-flow can execute repeatedly.
    // Keep slot registration idempotent to avoid unbounded duplicate roots.
    for (int i = 0; i < koda_global_slots_count; i++) {
        if (koda_global_slots[i] == slot) {
            return;
        }
    }
    if (koda_global_slots_count >= koda_global_slots_capacity) {
        int new_cap = koda_global_slots_capacity * 2;
        if (new_cap < 256) {
            new_cap = 256;
        }
        Value** next = (Value**)realloc(koda_global_slots, sizeof(Value*) * (size_t)new_cap);
        if (next == NULL) {
            fprintf(stderr, "koda: out of memory growing global slots array\n");
            exit(1);
        }
        koda_global_slots = next;
        koda_global_slots_capacity = new_cap;
    }
    koda_global_slots[koda_global_slots_count++] = slot;
}

void koda_runtime_init_ex(void* stack_base) {
    if (!koda_stack_base_plausible(stack_base)) {
        fprintf(stderr, "koda: invalid stack base passed to koda_runtime_init_ex\n");
        abort();
    }
    gc_stack_base = stack_base;
    koda_shadow_stack_configure_from_env();
    gc_init();
    koda_globals_init();
    koda_shadow_stack_init();
    koda_gc_use_shadow_stack(true);
}

void koda_runtime_init(void) {
#if defined(__clang__) || defined(__GNUC__)
    koda_runtime_init_ex(__builtin_frame_address(0));
#else
    uintptr_t anchor = 0;
    koda_runtime_init_ex(&anchor);
#endif
}

void koda_gc_set_threshold(size_t bytes) {
    gc_set_next_threshold(bytes);
}

void koda_gc_disable(void) {
    gc_set_disabled(true);
}

void koda_gc_enable(void) {
    gc_set_disabled(false);
}

void koda_gc_collect(void) {
    gc_collect();
}

void koda_gc_frame_step(double budget_ms) {
    if (gc_incremental_is_idle()) {
        size_t used = gc_nursery_used_bytes();
        size_t half = gc_nursery_capacity_bytes() / 2u;
        if (used >= half) {
            gc_collect_minor();
        }
    }
    if (budget_ms > 0.0) {
        double ms = budget_ms;
        if (ms < 0.0) {
            ms = 0.0;
        }
        uint64_t budget_us = (uint64_t)(ms * 1000.0 + 0.999);
        if (budget_us < 64u) {
            budget_us = 64u;
        }
        gc_frame_step_incremental(budget_us);
    }
}

Value koda_gc_stats(int argc, Value* argv) {
    (void)argc;
    (void)argv;
    GCStats st = gc_get_stats();
    Value objv = koda_allocate_object(5);
    /* Property keys are ASCII lowercase: dot access (e.g. s.bytesAllocated) compiles to s["bytesallocated"]. */
    koda_object_set(objv, koda_copy_string("collections", 11), NUMBER_VAL((double)st.collections));
    koda_object_set(objv, koda_copy_string("bytesallocated", 14), NUMBER_VAL((double)st.bytes_allocated));
    koda_object_set(objv, koda_copy_string("bytesfreed", 10), NUMBER_VAL((double)st.bytes_freed));
    koda_object_set(objv, koda_copy_string("maxpausetimeus", 14), NUMBER_VAL((double)st.max_pause_time_us));
    koda_object_set(objv, koda_copy_string("totalpausetimeus", 16), NUMBER_VAL((double)st.total_pause_time_us));
    return objv;
}

Value* koda_alloc_cell(void) {
    ObjCell* c = allocate_cell();
    return &c->value;
}

Value koda_cell_read(Value* cell) {
    return *cell;
}

void koda_cell_write(Value* cell, Value v) {
    ObjCell* cellObj = (ObjCell*)((uint8_t*)cell - offsetof(ObjCell, value));
    gc_write_barrier((Obj*)&cellObj->obj, v);
    *cell = v;
}

void koda_runtime_shutdown(void) {
    gc_collect();
    if (koda_gc_debug_enabled()) {
        fprintf(stderr,
            "koda gc debug: remembered_overflow=%llu shadow_depth_hwm=%d globals=%d/%d global_slots=%d/%d\n",
            (unsigned long long)gc_debug_remembered_overflow_count(),
            koda_shadow_depth_high_water,
            koda_globals_count,
            koda_globals_capacity,
            koda_global_slots_count,
            koda_global_slots_capacity);
    }
    koda_globals_free();
    if (koda_string_intern.buckets != NULL) {
        free(koda_string_intern.buckets);
        koda_string_intern.buckets = NULL;
    }
    koda_string_intern.count = 0;
    koda_string_intern.capacity = 0;
    if (koda_shadow_stack != NULL) {
        free(koda_shadow_stack);
        koda_shadow_stack = NULL;
    }
    koda_shadow_capacity = 0;
    koda_shadow_depth = 0;
}

Value koda_print_val(Value v) {
    print_value(v);
    printf("\n");
    fflush(stdout);
    return NIL_VAL;
}

void koda_print_newline(void) {
    printf("\n");
    fflush(stdout);
}

Value koda_print_argv(int arg_count, Value* args) {
    for (int i = 0; i < arg_count; i++) {
        print_value(args[i]);
        if (i < arg_count - 1) printf(" ");
    }
    printf("\n");
    return NIL_VAL;
}

Value koda_warn(int arg_count, Value* args) {
    for (int i = 0; i < arg_count; i++) {
        print_value(args[i]);
        if (i < arg_count - 1) {
            fprintf(stderr, " ");
        }
    }
    fprintf(stderr, "\n");
    fflush(stderr);
    return NIL_VAL;
}

Value koda_typeof(int arg_count, Value* args) {
    if (arg_count == 0) return NIL_VAL;
    
    Value v = args[0];
    if (IS_NIL(v)) {
        ObjString* str = allocate_string(3);
        memcpy(str->chars, "nil", 3);
        return OBJ_VAL((Obj*)str);
    } else if (IS_BOOL(v)) {
        ObjString* str = allocate_string(7);
        memcpy(str->chars, "boolean", 7);
        return OBJ_VAL((Obj*)str);
    } else if (IS_NUMBER(v)) {
        ObjString* str = allocate_string(6);
        memcpy(str->chars, "number", 6);
        return OBJ_VAL((Obj*)str);
    } else if (IS_OBJ(v)) {
        Obj* obj = AS_OBJ(v);
        const char* type_str = "object";
        switch (obj->type) {
            case OBJ_STRING: type_str = "string"; break;
            case OBJ_ARRAY: type_str = "array"; break;
            case OBJ_TABLE: type_str = "table"; break;
            case OBJ_FUNCTION: type_str = "function"; break;
            case OBJ_CLOSURE: type_str = "closure"; break;
            case OBJ_NATIVE: type_str = "native"; break;
            case OBJ_CELL: type_str = "cell"; break;
            case OBJ_ARENA: type_str = "arena"; break;
        }
        int len = strlen(type_str);
        ObjString* str = allocate_string(len);
        memcpy(str->chars, type_str, len);
        return OBJ_VAL((Obj*)str);
    }
    
    return NIL_VAL;
}

Value koda_allocate_object(int property_count) {
    ObjTable* table = allocate_table(property_count);
    return OBJ_VAL((Obj*)table);
}

Value koda_allocate_struct(int field_count) {
    ObjTable* table = allocate_struct_table(field_count);
    return OBJ_VAL((Obj*)table);
}

static ObjArena* koda_arena_from_value(Value v, const char* op) {
    if (!IS_OBJ(v)) {
        koda_type_error(op, "arena", v);
    }
    Obj* obj = AS_OBJ(v);
    if (obj->type != OBJ_ARENA) {
        koda_type_error(op, "arena", v);
    }
    return (ObjArena*)obj;
}

Value koda_arena(int argc, Value* argv) {
    if (argc < 1 || !IS_NUMBER(argv[0])) {
        koda_panic_str("arena() requires a size in bytes");
    }
    double nbytes = AS_NUMBER(argv[0]);
    if (nbytes <= 0.0 || nbytes > (double)SIZE_MAX) {
        koda_panic_str("arena size out of range");
    }
    ObjArena* arena = allocate_arena((size_t)nbytes);
    return OBJ_VAL((Obj*)arena);
}

Value koda_arena_reset(int argc, Value* argv) {
    if (argc < 1) {
        koda_panic_str("arenaReset() requires an arena");
    }
    ObjArena* arena = koda_arena_from_value(argv[0], "arenaReset");
    arena_reset(arena);
    return NIL_VAL;
}

Value koda_arena_alloc_array(int argc, Value* argv) {
    if (argc < 2) {
        koda_panic_str("arenaAllocArray() requires arena and capacity");
    }
    ObjArena* arena = koda_arena_from_value(argv[0], "arenaAllocArray");
    if (!IS_NUMBER(argv[1])) {
        koda_type_error("arenaAllocArray capacity", "number", argv[1]);
    }
    int capacity = (int)AS_NUMBER(argv[1]);
    if (capacity < 0) {
        koda_panic_str("arenaAllocArray capacity must be non-negative");
    }
    ObjArray* array = arena_allocate_array(arena, capacity);
    return OBJ_VAL((Obj*)array);
}

Value koda_arena_alloc_struct(int argc, Value* argv) {
    if (argc < 2) {
        koda_panic_str("arenaAllocStruct() requires arena and field count");
    }
    ObjArena* arena = koda_arena_from_value(argv[0], "arenaAllocStruct");
    if (!IS_NUMBER(argv[1])) {
        koda_type_error("arenaAllocStruct field count", "number", argv[1]);
    }
    int field_count = (int)AS_NUMBER(argv[1]);
    if (field_count < 1) {
        field_count = 1;
    }
    ObjTable* table = arena_allocate_struct_table(arena, field_count);
    return OBJ_VAL((Obj*)table);
}

static uint32_t fnv1a_bytes(const uint8_t* data, int len) {
    uint32_t h = 2166136261u;
    for (int i = 0; i < len; i++) {
        h ^= (uint32_t)data[i];
        h *= 16777619u;
    }
    return h == 0u ? 1u : h;
}

static uint32_t koda_string_hash(ObjString* s) {
    if (s->hash != 0u) {
        return s->hash;
    }
    s->hash = fnv1a_bytes((const uint8_t*)s->chars, s->length);
    if (s->hash == 0u) {
        s->hash = 1u;
    }
    return s->hash;
}

static uint32_t koda_hash_value(Value key) {
    if (IS_OBJ(key)) {
        Obj* ob = AS_OBJ(key);
        if (ob != NULL && ob->type == OBJ_STRING) {
            return koda_string_hash((ObjString*)ob);
        }
    }
    if (IS_NUMBER(key)) {
        union {
            Value v;
            uint64_t u;
        } u;
        u.v = key;
        uint32_t hi = (uint32_t)(u.u >> 32);
        uint32_t lo = (uint32_t)u.u;
        uint32_t h = hi ^ lo;
        return h == 0u ? 1u : h;
    }
    return 1u;
}

Value koda_struct_get(Value obj, int64_t index) {
    if (!IS_OBJ(obj)) {
        return NIL_VAL;
    }
    Obj* o = AS_OBJ(obj);
    if (o->type != OBJ_TABLE) {
        return NIL_VAL;
    }
    ObjTable* t = (ObjTable*)o;
    if (!t->is_struct_layout) {
        return NIL_VAL;
    }
    if (index < 0 || index >= t->count) {
        return NIL_VAL;
    }
    return t->values[(int)index];
}

Value koda_struct_set(Value obj, int64_t index, Value val) {
    if (!IS_OBJ(obj)) {
        return NIL_VAL;
    }
    Obj* o = AS_OBJ(obj);
    if (o->type != OBJ_TABLE) {
        return NIL_VAL;
    }
    ObjTable* t = (ObjTable*)o;
    if (!t->is_struct_layout) {
        return NIL_VAL;
    }
    if (index < 0 || index >= t->count) {
        return NIL_VAL;
    }
    gc_write_barrier(&t->obj, val);
    t->values[(int)index] = val;
    return val;
}

static Value koda_object_get_linear(ObjTable* table, Value key) {
    for (int i = 0; i < table->count; i++) {
        if (values_equal(table->keys[i], key)) {
            return table->values[i];
        }
    }
    return NIL_VAL;
}

static bool table_key_is_tombstone(Value k) {
    return IS_TOMBSTONE(k);
}

static Value koda_object_get_hash(ObjTable* table, Value key) {
    uint32_t cap = (uint32_t)table->capacity;
    if (cap == 0u) {
        return NIL_VAL;
    }
    uint32_t start = koda_hash_value(key) % cap;
    for (uint32_t p = 0u; p < cap; p++) {
        uint32_t i = (start + p) % cap;
        Value k = table->keys[i];
        if (IS_NIL(k)) {
            return NIL_VAL;
        }
        if (table_key_is_tombstone(k)) {
            continue;
        }
        if (values_equal(k, key)) {
            return table->values[i];
        }
    }
    return NIL_VAL;
}

Value koda_object_get(Value obj, Value key) {
    if (IS_NIL(obj)) {
        const char* prop = "?";
        if (IS_OBJ(key) && AS_OBJ(key)->type == OBJ_STRING) {
            prop = ((ObjString*)AS_OBJ(key))->chars;
        }
        char msg[256];
        snprintf(msg, sizeof(msg), "cannot read property '%s' of null", prop);
        koda_panic_str(msg);
    }
    if (!IS_OBJ(obj)) {
        return NIL_VAL;
    }
    Obj* o = AS_OBJ(obj);
    if (o->type != OBJ_TABLE) {
        return NIL_VAL;
    }

    ObjTable* table = (ObjTable*)o;
    if (table->is_struct_layout) {
        return koda_object_get_linear(table, key);
    }
    if (table->hashes != NULL) {
        return koda_object_get_hash(table, key);
    }
    return koda_object_get_linear(table, key);
}

static bool table_remove_pair(ObjTable* table, Value key) {
    if (table->hashes != NULL && !table->is_struct_layout) {
        uint32_t cap = (uint32_t)table->capacity;
        if (cap == 0u) {
            return false;
        }
        uint32_t start = koda_hash_value(key) % cap;
        for (uint32_t p = 0u; p < cap; p++) {
            uint32_t i = (start + p) % cap;
            Value k = table->keys[i];
            if (IS_NIL(k)) {
                return false;
            }
            if (table_key_is_tombstone(k)) {
                continue;
            }
            if (values_equal(k, key)) {
                gc_write_barrier(&table->obj, TOMBSTONE_VAL);
                table->keys[i] = TOMBSTONE_VAL;
                gc_write_barrier(&table->obj, NIL_VAL);
                table->values[i] = NIL_VAL;
                table->hashes[i] = 0u;
                table->count--;
                return true;
            }
        }
        return false;
    }
    for (int i = 0; i < table->count; i++) {
        if (values_equal(table->keys[i], key)) {
            for (int j = i; j < table->count - 1; j++) {
                gc_write_barrier(&table->obj, table->keys[j + 1]);
                table->keys[j] = table->keys[j + 1];
                gc_write_barrier(&table->obj, table->values[j + 1]);
                table->values[j] = table->values[j + 1];
            }
            table->count--;
            return true;
        }
    }
    return false;
}

/** Remove own property `key` from table object; returns boxed boolean. */
Value koda_object_remove(Value obj, Value key) {
    if (!IS_OBJ(obj)) return NIL_VAL;
    Obj* o = AS_OBJ(obj);
    if (o->type != OBJ_TABLE) return NIL_VAL;
    ObjTable* table = (ObjTable*)o;
    bool ok = table_remove_pair(table, key);
    return BOOL_VAL(ok);
}

static void table_rehash_open(Obj* parent, ObjTable* table);

static bool koda_object_set_hash(Obj* parent, ObjTable* table, Value key, Value value) {
    uint32_t cap = (uint32_t)table->capacity;
    if (cap == 0u) {
        return false;
    }
    if ((table->count + 1) * 10 > (int)cap * 7) {
        table_rehash_open(parent, table);
        cap = (uint32_t)table->capacity;
        if (cap == 0u) {
            return false;
        }
    }
    uint32_t start = koda_hash_value(key) % cap;
    int first_tomb = -1;
    for (uint32_t p = 0u; p < cap; p++) {
        uint32_t i = (start + p) % cap;
        Value k = table->keys[i];
        if (IS_NIL(k)) {
            int ins = (first_tomb >= 0) ? first_tomb : (int)i;
            gc_write_barrier(parent, key);
            gc_write_barrier(parent, value);
            table->keys[ins] = key;
            table->values[ins] = value;
            table->hashes[ins] = koda_hash_value(key);
            table->count++;
            return true;
        }
        if (table_key_is_tombstone(k)) {
            if (first_tomb < 0) {
                first_tomb = (int)i;
            }
            continue;
        }
        if (values_equal(k, key)) {
            gc_write_barrier(parent, key);
            gc_write_barrier(parent, value);
            table->keys[i] = key;
            table->values[i] = value;
            return true;
        }
    }
    table_rehash_open(parent, table);
    return koda_object_set_hash(parent, table, key, value);
}

static void table_rehash_open(Obj* parent, ObjTable* table) {
    int old_cap = table->capacity;
    Value* old_keys = table->keys;
    Value* old_vals = table->values;
    uint32_t* old_hashes = table->hashes;
    int new_cap = old_cap < 8 ? 16 : old_cap * 2;
    if (new_cap < old_cap) {
        return;
    }
    table->capacity = new_cap;
    table->count = 0;
    table->keys = (Value*)gc_alloc(sizeof(Value) * (size_t)new_cap);
    table->values = (Value*)gc_alloc(sizeof(Value) * (size_t)new_cap);
    table->hashes = (uint32_t*)gc_alloc(sizeof(uint32_t) * (size_t)new_cap);
    memset(table->hashes, 0, sizeof(uint32_t) * (size_t)new_cap);
    for (int i = 0; i < new_cap; i++) {
        table->keys[i] = NIL_VAL;
        table->values[i] = NIL_VAL;
    }
    for (int i = 0; i < old_cap; i++) {
        Value k = old_keys[i];
        if (IS_NIL(k) || table_key_is_tombstone(k)) {
            continue;
        }
        (void)koda_object_set_hash(parent, table, k, old_vals[i]);
    }
    if (old_keys != NULL) {
        gc_free(old_keys, (size_t)old_cap * sizeof(Value));
    }
    if (old_vals != NULL) {
        gc_free(old_vals, (size_t)old_cap * sizeof(Value));
    }
    if (old_hashes != NULL) {
        gc_free(old_hashes, (size_t)old_cap * sizeof(uint32_t));
    }
}

Value koda_object_set(Value obj, Value key, Value value) {
    if (!IS_OBJ(obj)) {
        return NIL_VAL;
    }
    Obj* o = AS_OBJ(obj);
    if (o->type != OBJ_TABLE) {
        return NIL_VAL;
    }

    ObjTable* table = (ObjTable*)o;
    if (table->is_struct_layout) {
        for (int i = 0; i < table->count; i++) {
            if (values_equal(table->keys[i], key)) {
                gc_write_barrier(o, key);
                gc_write_barrier(o, value);
                table->keys[i] = key;
                table->values[i] = value;
                return value;
            }
        }
        return NIL_VAL;
    }
    if (table->hashes != NULL) {
        if (koda_object_set_hash(o, table, key, value)) {
            return value;
        }
        return NIL_VAL;
    }

    for (int i = 0; i < table->count; i++) {
        if (values_equal(table->keys[i], key)) {
            gc_write_barrier(o, key);
            gc_write_barrier(o, value);
            table->keys[i] = key;
            table->values[i] = value;
            return value;
        }
    }

    if (table->count < table->capacity) {
        gc_write_barrier(o, key);
        gc_write_barrier(o, value);
        table->keys[table->count] = key;
        table->values[table->count] = value;
        table->count++;
        return value;
    }

    return NIL_VAL;
}

Value koda_allocate_string(int length, char* chars) {
    ObjString* interned = koda_intern_find(chars, length);
    if (interned != NULL) {
        return OBJ_VAL((Obj*)interned);
    }
    ObjString* str = allocate_string(length);
    memcpy(str->chars, chars, length);
    str->chars[length] = '\0';
    koda_intern_add(str);
    return OBJ_VAL((Obj*)str);
}

Value koda_copy_string(const char* chars, int length) {
    if (chars == NULL || length < 0) {
        return NIL_VAL;
    }
    ObjString* interned = koda_intern_find(chars, length);
    if (interned != NULL) {
        return OBJ_VAL((Obj*)interned);
    }
    ObjString* str = allocate_string(length);
    memcpy(str->chars, chars, (size_t)length);
    str->chars[length] = '\0';
    koda_intern_add(str);
    return OBJ_VAL((Obj*)str);
}

Value koda_allocate_array(int length) {
    ObjArray* arr = allocate_array(length);
    return OBJ_VAL((Obj*)arr);
}

Value koda_array_get(Value arr, int index) {
	if (!IS_OBJ(arr)) {
		koda_type_error("[]", "array", arr);
	}
	Obj* obj = AS_OBJ(arr);
	if (obj->type != OBJ_ARRAY) {
		koda_type_error("[]", "array", arr);
	}

	ObjArray* array = (ObjArray*)obj;
	if (index < 0 || index >= array->count) {
		char msg[256];
		snprintf(msg, sizeof(msg), "index %d out of bounds (array length %d)", index, array->count);
		koda_panic_str(msg);
	}

	return array->elements[index];
}

double koda_unbox_number(Value v) {
	if (!IS_NUMBER(v)) {
		koda_type_error("unbox", "number", v);
	}
	return AS_NUMBER(v);
}

Value koda_box_number(double d) {
    return NUMBER_VAL(d);
}

Value koda_get(Value obj, Value key) {
	if (!IS_OBJ(obj)) {
		koda_type_error("[]", "object", obj);
	}
	Obj* o = AS_OBJ(obj);
	if (o->type == OBJ_ARRAY) {
		if (!IS_NUMBER(key)) {
			koda_type_error("array index", "number", key);
		}
		int idx = (int)AS_NUMBER(key);
		return koda_array_get(obj, idx);
	}
	if (o->type == OBJ_STRING) {
		if (!IS_NUMBER(key)) {
			koda_type_error("string index", "number", key);
		}
		int idx = (int)AS_NUMBER(key);
		ObjString* s = (ObjString*)o;
		if (idx < 0 || idx >= s->length) {
			char msg[256];
			snprintf(msg, sizeof(msg), "index %d out of bounds (string length %d)", idx, s->length);
			koda_panic_str(msg);
		}
		return koda_copy_string(&s->chars[idx], 1);
	}
	if (o->type == OBJ_TABLE) {
		return koda_object_get(obj, key);
	}
	koda_type_error("[]", "indexable object", obj);
}

Value koda_get_index(Value obj, Value key) {
    return koda_get(obj, key);
}

void koda_array_set(Value arr, int64_t index, Value value) {
	if (!IS_OBJ(arr)) {
		koda_type_error("[]= (array)", "array", arr);
	}
	Obj* obj = AS_OBJ(arr);
	if (obj->type != OBJ_ARRAY) {
		koda_type_error("[]= (array)", "array", arr);
	}

	ObjArray* array = (ObjArray*)obj;
	if (index < 0 || index >= (int64_t)array->capacity) {
		char msg[320];
		snprintf(msg, sizeof(msg),
			"index %lld out of bounds for assignment (logical length %d, capacity %d)",
			(long long)index, array->count, array->capacity);
		koda_panic_str(msg);
	}

	gc_write_barrier(obj, value);
	array->elements[(int)index] = value;
	if ((int64_t)index + 1 > (int64_t)array->count) {
		array->count = (int)index + 1;
	}
}

Value koda_set(Value obj, Value key, Value val) {
	if (!IS_OBJ(obj)) {
		koda_type_error("[]=", "object", obj);
	}
	Obj* o = AS_OBJ(obj);
	if (o->type == OBJ_ARRAY) {
		if (!IS_NUMBER(key)) {
			koda_type_error("array index (assign)", "number", key);
		}
		int64_t i = (int64_t)AS_NUMBER(key);
		koda_array_set(obj, i, val);
		return val;
	}
	if (o->type == OBJ_TABLE) {
		return koda_object_set(obj, key, val);
	}
	koda_type_error("[]=", "array or object", obj);
}

void koda_array_push(Value arr, Value value) {
    if (!IS_OBJ(arr)) return;
    Obj* obj = AS_OBJ(arr);
    if (obj->type != OBJ_ARRAY) return;
    
    ObjArray* array = (ObjArray*)obj;
    if (array->count >= array->capacity) {
        int old_cap = array->capacity;
        int new_cap = old_cap < 1 ? 1 : old_cap * 2;
        if (new_cap < old_cap) {
            koda_panic_str("array push capacity overflow");
        }
        if ((size_t)new_cap > SIZE_MAX / sizeof(Value)) {
            koda_panic_str("array push capacity overflow");
        }
        Value* next = (Value*)gc_alloc(sizeof(Value) * (size_t)new_cap);
        if (array->elements != NULL && old_cap > 0) {
            memcpy(next, array->elements, sizeof(Value) * (size_t)old_cap);
            if (!array->inline_elements) {
                gc_free(array->elements, sizeof(Value) * (size_t)old_cap);
            }
        }
        array->elements = next;
        array->inline_elements = false;
        array->capacity = new_cap;
    }
    gc_write_barrier(obj, value);
    array->elements[array->count] = value;
    array->count++;
}

Value koda_array_pop(Value arr) {
    if (!IS_OBJ(arr)) return NIL_VAL;
    Obj* obj = AS_OBJ(arr);
    if (obj->type != OBJ_ARRAY) return NIL_VAL;
    
    ObjArray* array = (ObjArray*)obj;
    if (array->count > 0) {
        array->count--;
        return array->elements[array->count];
    }
    return NIL_VAL;
}

Value koda_array_length(Value arr) {
    if (!IS_OBJ(arr)) return NIL_VAL;
    Obj* obj = AS_OBJ(arr);
    if (obj->type != OBJ_ARRAY) return NIL_VAL;
    
    ObjArray* array = (ObjArray*)obj;
    return NUMBER_VAL(array->count);
}

/** Slot iteration length for native `for-of` / `for-in` (arrays + tables). */
Value koda_forof_length(Value v) {
    if (!IS_OBJ(v)) return NUMBER_VAL(0);
    Obj* o = AS_OBJ(v);
    if (o->type == OBJ_ARRAY) {
        return NUMBER_VAL((double)((ObjArray*)o)->count);
    }
    if (o->type == OBJ_TABLE) {
        return NUMBER_VAL((double)((ObjTable*)o)->count);
    }
    return NUMBER_VAL(0);
}

/** Key at linear slot `idx` (number): array → numeric index; table → insertion-order key. */
Value koda_forof_key_at(Value v, Value idx_val) {
    if (!IS_OBJ(v) || !IS_NUMBER(idx_val)) return NIL_VAL;
    int i = (int)AS_NUMBER(idx_val);
    Obj* o = AS_OBJ(v);
    if (o->type == OBJ_ARRAY) {
        ObjArray* a = (ObjArray*)o;
        if (i < 0 || i >= a->count) return NIL_VAL;
        return NUMBER_VAL((double)i);
    }
    if (o->type == OBJ_TABLE) {
        ObjTable* t = (ObjTable*)o;
        if (i < 0 || i >= t->count) return NIL_VAL;
        return t->keys[i];
    }
    return NIL_VAL;
}

/** Value at linear slot `idx` (insertion order for tables). */
Value koda_forof_value_at(Value v, Value idx_val) {
    if (!IS_OBJ(v) || !IS_NUMBER(idx_val)) return NIL_VAL;
    int i = (int)AS_NUMBER(idx_val);
    Obj* o = AS_OBJ(v);
    if (o->type == OBJ_ARRAY) {
        ObjArray* a = (ObjArray*)o;
        if (i < 0 || i >= a->count) return NIL_VAL;
        return a->elements[i];
    }
    if (o->type == OBJ_TABLE) {
        ObjTable* t = (ObjTable*)o;
        if (i < 0 || i >= t->count) return NIL_VAL;
        return t->values[i];
    }
    return NIL_VAL;
}

// Standard library functions
Value koda_type(Value value) {
    if (IS_NIL(value)) {
        ObjString* str = allocate_string(3);
        memcpy(str->chars, "nil", 3);
        return OBJ_VAL((Obj*)str);
    } else if (IS_BOOL(value)) {
        ObjString* str = allocate_string(7);
        memcpy(str->chars, "boolean", 7);
        return OBJ_VAL((Obj*)str);
    } else if (IS_NUMBER(value)) {
        ObjString* str = allocate_string(6);
        memcpy(str->chars, "number", 6);
        return OBJ_VAL((Obj*)str);
    } else if (IS_OBJ(value)) {
        Obj* obj = AS_OBJ(value);
        const char* type_str = "object";
        switch (obj->type) {
            case OBJ_STRING: type_str = "string"; break;
            case OBJ_ARRAY: type_str = "array"; break;
            case OBJ_TABLE: type_str = "table"; break;
            case OBJ_FUNCTION: type_str = "function"; break;
            case OBJ_CLOSURE: type_str = "closure"; break;
            case OBJ_NATIVE: type_str = "native"; break;
            case OBJ_CELL: type_str = "cell"; break;
            case OBJ_ARENA: type_str = "arena"; break;
        }
        int len = strlen(type_str);
        ObjString* str = allocate_string(len);
        memcpy(str->chars, type_str, len);
        return OBJ_VAL((Obj*)str);
    }
    return NIL_VAL;
}

Value koda_len(Value value) {
    if (IS_OBJ(value)) {
        Obj* obj = AS_OBJ(value);
        switch (obj->type) {
            case OBJ_STRING: {
                ObjString* str = (ObjString*)obj;
                return NUMBER_VAL(str->length);
            }
            case OBJ_ARRAY: {
                ObjArray* arr = (ObjArray*)obj;
                return NUMBER_VAL(arr->count);
            }
            case OBJ_TABLE: {
                ObjTable* table = (ObjTable*)obj;
                return NUMBER_VAL((double)table->count);
            }
            default:
                break;
        }
    }
    return NUMBER_VAL(0);
}

Value koda_abs(Value value) {
    if (!IS_NUMBER(value)) return NIL_VAL;
    double num = AS_NUMBER(value);
    if (num < 0) num = -num;
    return NUMBER_VAL(num);
}

Value koda_sqrt(Value value) {
    if (!IS_NUMBER(value)) return NIL_VAL;
    double num = AS_NUMBER(value);
    if (num < 0) return NIL_VAL;
    return NUMBER_VAL(sqrt(num));
}

Value koda_cbrt(Value value) {
    if (!IS_NUMBER(value)) return NIL_VAL;
    return NUMBER_VAL(cbrt(AS_NUMBER(value)));
}

Value koda_time() {
    return NUMBER_VAL((double)time(NULL));
}

Value koda_clock(void) {
    return NUMBER_VAL((double)clock() / (double)CLOCKS_PER_SEC);
}

Value koda_wall_time(void) {
    return NUMBER_VAL((double)time(NULL));
}

void koda_sleep(int64_t ms) {
    // Sleep implementation (platform-specific)
    #ifdef _WIN32
        Sleep((DWORD)ms);
    #else
        usleep(ms * 1000);
    #endif
}

Value koda_number(Value value) {
    if (IS_NUMBER(value)) return value;
    if (IS_FALSE(value)) return NUMBER_VAL(0);
    if (IS_TRUE(value)) return NUMBER_VAL(1);
    if (IS_NIL(value)) return NUMBER_VAL(0);
    return NIL_VAL;
}

Value koda_string(Value value) {
    if (IS_OBJ(value)) {
        Obj* obj = AS_OBJ(value);
        if (obj->type == OBJ_STRING) return value;
    }
    if (IS_NUMBER(value)) {
        char buffer[32];
        snprintf(buffer, 32, "%g", AS_NUMBER(value));
        int len = strlen(buffer);
        ObjString* str = allocate_string(len);
        memcpy(str->chars, buffer, len);
        return OBJ_VAL((Obj*)str);
    }
    if (IS_TRUE(value)) {
        ObjString* str = allocate_string(4);
        memcpy(str->chars, "true", 4);
        return OBJ_VAL((Obj*)str);
    }
    if (IS_FALSE(value)) {
        ObjString* str = allocate_string(5);
        memcpy(str->chars, "false", 5);
        return OBJ_VAL((Obj*)str);
    }
    if (IS_NIL(value)) {
        ObjString* str = allocate_string(3);
        memcpy(str->chars, "nil", 3);
        return OBJ_VAL((Obj*)str);
    }
    return NIL_VAL;
}

Value koda_string_concat(Value a, Value b) {
    /* Intermediate strings live in C stack slots; conservative stack scan keeps them
     * rooted across chained allocations. Precise shadow-stack mode does not scan C
     * locals — runtime helpers that chain allocations must not run with shadow stack on. */
    if (gc_uses_shadow_stack()) {
        koda_panic_str("internal error: koda_string_concat called with precise GC active");
    }
    Value sa = koda_string(a);
    Value sb = koda_string(b);
    if (!IS_OBJ(sa) || AS_OBJ(sa)->type != OBJ_STRING) return sb;
    if (!IS_OBJ(sb) || AS_OBJ(sb)->type != OBJ_STRING) return sa;
    ObjString* A = (ObjString*)AS_OBJ(sa);
    ObjString* B = (ObjString*)AS_OBJ(sb);
    int len = A->length + B->length;
    ObjString* out = allocate_string(len);
    memcpy(out->chars, A->chars, A->length);
    memcpy(out->chars + A->length, B->chars, B->length);
    return OBJ_VAL((Obj*)out);
}

// Time functions for game development
Value koda_delta_time(int arg_count, Value* args) {
    (void)arg_count;
    (void)args;
    koda_frame_clock_ensure_init();
    double now = koda_monotonic_seconds();
    double dt = now - koda_frame_clock_last;
    koda_frame_clock_last = now;
    if (dt <= 0.0 || dt > 0.25) {
        dt = 0.25;
    }
    return NUMBER_VAL(dt);
}

Value koda_program_time(int arg_count, Value* args) {
    (void)arg_count;
    (void)args;
    koda_frame_clock_ensure_init();
    double now = koda_monotonic_seconds();
    return NUMBER_VAL(now - koda_frame_clock_start);
}

Value koda_timestamp(int arg_count, Value* args) {
    // Unix timestamp
    return NUMBER_VAL((double)time(NULL));
}

// Random functions for game development (xoshiro128** PRNG)
Value koda_random(int arg_count, Value* args) {
    switch (arg_count) {
    case 0:
        return NUMBER_VAL(koda_rng_unit());
    case 1:
        if (!IS_NUMBER(args[0])) return NIL_VAL;
        double max = AS_NUMBER(args[0]);
        return NUMBER_VAL(koda_rng_unit() * max);
    case 2:
        if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1])) return NIL_VAL;
        double min = AS_NUMBER(args[0]);
        double max2 = AS_NUMBER(args[1]);
        if (max2 < min) {
            double t = min;
            min = max2;
            max2 = t;
        }
        return NUMBER_VAL(min + koda_rng_unit() * (max2 - min));
    default:
        return NIL_VAL;
    }
}

Value koda_randomInt(int arg_count, Value* args) {
    switch (arg_count) {
    case 1:
        if (!IS_NUMBER(args[0])) return NIL_VAL;
        int max = (int)AS_NUMBER(args[0]);
        if (max <= 0) return NUMBER_VAL(0);
        return NUMBER_VAL((double)(koda_rng_next_u32() % (uint32_t)max));
    case 2:
        if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1])) return NIL_VAL;
        int min = (int)AS_NUMBER(args[0]);
        int max2 = (int)AS_NUMBER(args[1]);
        if (max2 <= min) return NUMBER_VAL((double)min);
        int span = max2 - min;
        return NUMBER_VAL((double)(min + (int)(koda_rng_next_u32() % (uint32_t)span)));
    default:
        return NIL_VAL;
    }
}

Value koda_randomChoice(int arg_count, Value* args) {
    if (arg_count < 1 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_ARRAY) {
        return NIL_VAL;
    }
    ObjArray* arr = (ObjArray*)AS_OBJ(args[0]);
    if (arr->count <= 0) {
        return NIL_VAL;
    }
    int idx = (int)(koda_rng_next_u32() % (uint32_t)arr->count);
    return arr->elements[idx];
}

Value koda_randomSeed(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    double d = AS_NUMBER(args[0]);
    uint64_t seed = (uint64_t)d;
    if (seed == 0) {
        seed = (uint64_t)(d * 1000000.0);
    }
    seed ^= (uint64_t)(uintptr_t)&koda_rng_state;
    koda_rng_seed_u64(seed);
    return NIL_VAL;
}

// Math functions for game development
Value koda_lerp(int arg_count, Value* args) {
    if (arg_count != 3) return NIL_VAL;
    if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1]) || !IS_NUMBER(args[2])) return NIL_VAL;
    double a = AS_NUMBER(args[0]);
    double b = AS_NUMBER(args[1]);
    double t = AS_NUMBER(args[2]);
    return NUMBER_VAL(a + (b - a) * t);
}

Value koda_clamp(int arg_count, Value* args) {
    if (arg_count != 3) return NIL_VAL;
    if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1]) || !IS_NUMBER(args[2])) return NIL_VAL;
    double val = AS_NUMBER(args[0]);
    double min = AS_NUMBER(args[1]);
    double max2 = AS_NUMBER(args[2]);
    
    if (val < min) return NUMBER_VAL(min);
    if (val > max2) return NUMBER_VAL(max2);
    return NUMBER_VAL(val);
}

Value koda_distance(int arg_count, Value* args) {
    if (arg_count != 4) return NIL_VAL;
    if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1]) || !IS_NUMBER(args[2]) || !IS_NUMBER(args[3])) return NIL_VAL;
    double x1 = AS_NUMBER(args[0]);
    double y1 = AS_NUMBER(args[1]);
    double x2 = AS_NUMBER(args[2]);
    double y2 = AS_NUMBER(args[3]);
    
    double dx = x2 - x1;
    double dy = y2 - y1;
    return NUMBER_VAL(sqrt(dx * dx + dy * dy));
}

Value koda_angleBetween(int arg_count, Value* args) {
    if (arg_count != 4) return NIL_VAL;
    if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1]) || !IS_NUMBER(args[2]) || !IS_NUMBER(args[3])) return NIL_VAL;
    double x1 = AS_NUMBER(args[0]);
    double y1 = AS_NUMBER(args[1]);
    double x2 = AS_NUMBER(args[2]);
    double y2 = AS_NUMBER(args[3]);
    
    return NUMBER_VAL(atan2(y2 - y1, x2 - x1));
}

Value koda_map(int arg_count, Value* args) {
    if (arg_count != 5) return NIL_VAL;
    if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1]) || !IS_NUMBER(args[2]) || !IS_NUMBER(args[3]) || !IS_NUMBER(args[4])) return NIL_VAL;
    double val = AS_NUMBER(args[0]);
    double inMin = AS_NUMBER(args[1]);
    double inMax = AS_NUMBER(args[2]);
    double outMin = AS_NUMBER(args[3]);
    double outMax = AS_NUMBER(args[4]);
    
    double normalized = (val - inMin) / (inMax - inMin);
    return NUMBER_VAL(outMin + normalized * (outMax - outMin));
}

Value koda_pi(int arg_count, Value* args) {
    return NUMBER_VAL(3.141592653589793);
}

Value koda_e(int arg_count, Value* args) {
    return NUMBER_VAL(2.718281828459045);
}

Value koda_sin(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(sin(AS_NUMBER(args[0])));
}

Value koda_cos(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(cos(AS_NUMBER(args[0])));
}

Value koda_tan(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(tan(AS_NUMBER(args[0])));
}

Value koda_asin(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(asin(AS_NUMBER(args[0])));
}

Value koda_acos(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(acos(AS_NUMBER(args[0])));
}

Value koda_atan(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(atan(AS_NUMBER(args[0])));
}

Value koda_atan2(int arg_count, Value* args) {
    if (arg_count != 2 || !IS_NUMBER(args[0]) || !IS_NUMBER(args[1])) return NIL_VAL;
    return NUMBER_VAL(atan2(AS_NUMBER(args[0]), AS_NUMBER(args[1])));
}

Value koda_pow(int arg_count, Value* args) {
    if (arg_count != 2 || !IS_NUMBER(args[0]) || !IS_NUMBER(args[1])) return NIL_VAL;
    return NUMBER_VAL(pow(AS_NUMBER(args[0]), AS_NUMBER(args[1])));
}

Value koda_exp(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(exp(AS_NUMBER(args[0])));
}

Value koda_log(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(log(AS_NUMBER(args[0])));
}

Value koda_log10(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(log10(AS_NUMBER(args[0])));
}

Value koda_log2(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(log2(AS_NUMBER(args[0])));
}

Value koda_floor(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(floor(AS_NUMBER(args[0])));
}

Value koda_ceil(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(ceil(AS_NUMBER(args[0])));
}

Value koda_round(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(round(AS_NUMBER(args[0])));
}

Value koda_trunc(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(trunc(AS_NUMBER(args[0])));
}

Value koda_sign(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    double val = AS_NUMBER(args[0]);
    if (val > 0) return NUMBER_VAL(1);
    if (val < 0) return NUMBER_VAL(-1);
    return NUMBER_VAL(0);
}

Value koda_min(int arg_count, Value* args) {
    if (arg_count == 0) return NIL_VAL;
    double result = AS_NUMBER(args[0]);
    for (int i = 1; i < arg_count; i++) {
        if (!IS_NUMBER(args[i])) return NIL_VAL;
        double val = AS_NUMBER(args[i]);
        if (val < result) result = val;
    }
    return NUMBER_VAL(result);
}

Value koda_max(int arg_count, Value* args) {
    if (arg_count == 0) return NIL_VAL;
    double result = AS_NUMBER(args[0]);
    for (int i = 1; i < arg_count; i++) {
        if (!IS_NUMBER(args[i])) return NIL_VAL;
        double val = AS_NUMBER(args[i]);
        if (val > result) result = val;
    }
    return NUMBER_VAL(result);
}

Value koda_smoothstep(int arg_count, Value* args) {
    if (arg_count != 3) return NIL_VAL;
    if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1]) || !IS_NUMBER(args[2])) return NIL_VAL;
    double a = AS_NUMBER(args[0]);
    double b = AS_NUMBER(args[1]);
    double t = AS_NUMBER(args[2]);
    
    // Clamp t to [0, 1]
    if (t < 0) t = 0;
    if (t > 1) t = 1;
    
    // Smoothstep interpolation
    double result = t * t * (3 - 2 * t);
    return NUMBER_VAL(a + result * (b - a));
}

Value koda_distanceSq(int arg_count, Value* args) {
    if (arg_count != 4) return NIL_VAL;
    if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1]) || !IS_NUMBER(args[2]) || !IS_NUMBER(args[3])) return NIL_VAL;
    double x1 = AS_NUMBER(args[0]);
    double y1 = AS_NUMBER(args[1]);
    double x2 = AS_NUMBER(args[2]);
    double y2 = AS_NUMBER(args[3]);
    
    double dx = x2 - x1;
    double dy = y2 - y1;
    return NUMBER_VAL(dx * dx + dy * dy);
}

Value koda_normalize(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1])) return NIL_VAL;

    double x = AS_NUMBER(args[0]);
    double y = AS_NUMBER(args[1]);

    double len = sqrt(x * x + y * y);
    if (len == 0) {
        return NIL_VAL;
    }

    Value objv = koda_allocate_object(2);
    koda_object_set(objv, koda_copy_string("x", 1), NUMBER_VAL(x / len));
    koda_object_set(objv, koda_copy_string("y", 1), NUMBER_VAL(y / len));
    return objv;
}

Value koda_hypot(int arg_count, Value* args) {
    if (arg_count != 2 || !IS_NUMBER(args[0]) || !IS_NUMBER(args[1])) return NIL_VAL;
    return NUMBER_VAL(hypot(AS_NUMBER(args[0]), AS_NUMBER(args[1])));
}

Value koda_fmod(int arg_count, Value* args) {
    if (arg_count != 2 || !IS_NUMBER(args[0]) || !IS_NUMBER(args[1])) return NIL_VAL;
    return NUMBER_VAL(fmod(AS_NUMBER(args[0]), AS_NUMBER(args[1])));
}

Value koda_degrees(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(AS_NUMBER(args[0]) * (180.0 / 3.141592653589793));
}

Value koda_radians(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    return NUMBER_VAL(AS_NUMBER(args[0]) * (3.141592653589793 / 180.0));
}

Value koda_wrap(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_NUMBER(args[0])) return NIL_VAL;
    double x = AS_NUMBER(args[0]);
    const double pi = 3.141592653589793;
    const double two_pi = 2.0 * pi;
    x = fmod(x + pi, two_pi);
    if (x < 0.0) x += two_pi;
    return NUMBER_VAL(x - pi);
}

Value koda_approach(int arg_count, Value* args) {
    if (arg_count != 3 || !IS_NUMBER(args[0]) || !IS_NUMBER(args[1]) || !IS_NUMBER(args[2])) return NIL_VAL;
    double value = AS_NUMBER(args[0]);
    double target = AS_NUMBER(args[1]);
    double delta = AS_NUMBER(args[2]);
    if (delta < 0.0) delta = -delta;
    if (value < target) {
        value += delta;
        if (value > target) value = target;
    } else if (value > target) {
        value -= delta;
        if (value < target) value = target;
    }
    return NUMBER_VAL(value);
}

Value koda_smoothdamp(int arg_count, Value* args) {
    if (arg_count != 4 || !IS_NUMBER(args[0]) || !IS_NUMBER(args[1]) || !IS_NUMBER(args[2]) || !IS_NUMBER(args[3])) return NIL_VAL;
    double current = AS_NUMBER(args[0]);
    double target = AS_NUMBER(args[1]);
    double velocity = AS_NUMBER(args[2]);
    double smooth_time = AS_NUMBER(args[3]);
    if (smooth_time <= 0.0001) smooth_time = 0.0001;
    double omega = 2.0 / smooth_time;
    double x = omega * (1.0 / 60.0);
    double exp_term = 1.0 / (1.0 + x + 0.48 * x * x + 0.235 * x * x * x);
    double change = current - target;
    double temp = (velocity + omega * change) * (1.0 / 60.0);
    double next_velocity = (velocity - omega * temp) * exp_term;
    double output = target + (change + temp) * exp_term;
    (void)next_velocity;
    return NUMBER_VAL(output);
}

// Type checking functions
Value koda_isNumber(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    return BOOL_VAL(IS_NUMBER(args[0]));
}

Value koda_isString(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    return BOOL_VAL(IS_OBJ(args[0]) && AS_OBJ(args[0])->type == OBJ_STRING);
}

Value koda_isBool(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    return BOOL_VAL(IS_BOOL(args[0]));
}

Value koda_isNull(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    return BOOL_VAL(IS_NIL(args[0]));
}

Value koda_isArray(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    return BOOL_VAL(IS_OBJ(args[0]) && AS_OBJ(args[0])->type == OBJ_ARRAY);
}

Value koda_isObject(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    return BOOL_VAL(IS_OBJ(args[0]) && AS_OBJ(args[0])->type == OBJ_TABLE);
}

Value koda_isFunction(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    return BOOL_VAL(IS_OBJ(args[0]) && AS_OBJ(args[0])->type == OBJ_CLOSURE);
}

// Conversion functions
Value koda_bool(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    
    Value val = args[0];
    if (IS_BOOL(val)) return val;
    if (IS_NUMBER(val)) return BOOL_VAL(AS_NUMBER(val) != 0);
    if (IS_NIL(val)) return BOOL_VAL(false);
    if (IS_OBJ(val)) return BOOL_VAL(true);
    
    return BOOL_VAL(false);
}

// Format function: args[0] is the format string; each "{}" is replaced by koda_string(args[1]), args[2], ...
// If there are more placeholders than arguments, remaining "{}" are copied literally.
static int koda_format_piece_len(Value v) {
    Value s = koda_string(v);
    if (!IS_OBJ(s) || AS_OBJ(s)->type != OBJ_STRING) return 0;
    return ((ObjString*)AS_OBJ(s))->length;
}

Value koda_format(int arg_count, Value* args) {
    if (arg_count < 1) return NIL_VAL;

    Value fmtv = koda_string(args[0]);
    if (!IS_OBJ(fmtv) || AS_OBJ(fmtv)->type != OBJ_STRING) return NIL_VAL;
    ObjString* fmt = (ObjString*)AS_OBJ(fmtv);

    int total = 0;
    int i = 0;
    int ap = 1;
    while (i < fmt->length) {
        if (i + 1 < fmt->length && fmt->chars[i] == '{' && fmt->chars[i + 1] == '}') {
            if (ap < arg_count) {
                total += koda_format_piece_len(args[ap]);
                ap++;
            } else {
                total += 2;
            }
            i += 2;
        } else {
            total++;
            i++;
        }
    }

    ObjString* out = allocate_string(total);
    int pos = 0;
    i = 0;
    ap = 1;
    while (i < fmt->length) {
        if (i + 1 < fmt->length && fmt->chars[i] == '{' && fmt->chars[i + 1] == '}') {
            if (ap < arg_count) {
                Value piece = koda_string(args[ap]);
                ap++;
                if (IS_OBJ(piece) && AS_OBJ(piece)->type == OBJ_STRING) {
                    ObjString* p = (ObjString*)AS_OBJ(piece);
                    memcpy(out->chars + pos, p->chars, p->length);
                    pos += p->length;
                }
            } else {
                out->chars[pos++] = '{';
                out->chars[pos++] = '}';
            }
            i += 2;
        } else {
            out->chars[pos++] = fmt->chars[i];
            i++;
        }
    }
    return OBJ_VAL((Obj*)out);
}

// Integer range as half-open [from, to): koda_range(from, to) -> array of numbers.
Value koda_range(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_NUMBER(args[0]) || !IS_NUMBER(args[1])) return NIL_VAL;
    int start = (int)AS_NUMBER(args[0]);
    int end = (int)AS_NUMBER(args[1]);
    if (end < start) {
        ObjArray* empty = allocate_array(1);
        return OBJ_VAL((Obj*)empty);
    }
    int n = end - start;
    ObjArray* out = allocate_array(n < 1 ? 1 : n);
    Value ov = OBJ_VAL((Obj*)out);
    for (int i = 0; i < n; i++) {
        koda_array_push(ov, NUMBER_VAL((double)(start + i)));
    }
    return ov;
}

// Array methods for game development
//
// Higher-order helpers taking Koda callbacks are lowered by the LLVM backend on `.map()` /
// `.filter()` etc.; argv variants remain unavailable from native Tier-1 glue.

Value koda_array_map(int arg_count, Value* args) {
    (void)args;
    if (arg_count != 2) return NIL_VAL;
    return NIL_VAL;
}

Value koda_array_filter(int arg_count, Value* args) {
    (void)args;
    if (arg_count != 2) return NIL_VAL;
    return NIL_VAL;
}

Value koda_array_forEach(int arg_count, Value* args) {
    (void)args;
    if (arg_count != 2) return NIL_VAL;
    return NIL_VAL;
}

Value koda_array_find(int arg_count, Value* args) {
    (void)args;
    if (arg_count != 2) return NIL_VAL;
    return NIL_VAL;
}

Value koda_array_findIndex(int arg_count, Value* args) {
    (void)args;
    if (arg_count != 2) return NIL_VAL;
    return NUMBER_VAL(-1);
}

Value koda_array_some(int arg_count, Value* args) {
    (void)args;
    if (arg_count != 2) return NIL_VAL;
    return BOOL_VAL(false);
}

Value koda_array_every(int arg_count, Value* args) {
    (void)args;
    if (arg_count != 2) return NIL_VAL;
    return BOOL_VAL(false);
}

Value koda_array_reduce(int arg_count, Value* args) {
    (void)args;
    if (arg_count != 2) return NIL_VAL;
    return NIL_VAL;
}

static int cmp_sort_values(Value a, Value b) {
    if (IS_NUMBER(a) && IS_NUMBER(b)) {
        double da = AS_NUMBER(a);
        double db = AS_NUMBER(b);
        return (da > db) - (da < db);
    }
    Value sa = koda_string(a);
    Value sb = koda_string(b);
    if (!IS_OBJ(sa) || AS_OBJ(sa)->type != OBJ_STRING) {
        sa = koda_copy_string("", 0);
    }
    if (!IS_OBJ(sb) || AS_OBJ(sb)->type != OBJ_STRING) {
        sb = koda_copy_string("", 0);
    }
    ObjString* A = (ObjString*)AS_OBJ(sa);
    ObjString* B = (ObjString*)AS_OBJ(sb);
    int minlen = A->length < B->length ? A->length : B->length;
    int c = memcmp(A->chars, B->chars, (size_t)minlen);
    if (c != 0) {
        return c > 0 ? 1 : -1;
    }
    return (A->length > B->length) - (A->length < B->length);
}

Value koda_array_sort(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_ARRAY) return NIL_VAL;
    ObjArray* arr = (ObjArray*)AS_OBJ(args[0]);
    int n = arr->count;
    for (int i = 0; i < n - 1; i++) {
        for (int j = 0; j < n - 1 - i; j++) {
            if (cmp_sort_values(arr->elements[j], arr->elements[j + 1]) > 0) {
                Value tmp = arr->elements[j];
                gc_write_barrier((Obj*)arr, arr->elements[j + 1]);
                arr->elements[j] = arr->elements[j + 1];
                gc_write_barrier((Obj*)arr, tmp);
                arr->elements[j + 1] = tmp;
            }
        }
    }
    return args[0];
}

Value koda_array_reverse(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_ARRAY) return NIL_VAL;
    ObjArray* src = (ObjArray*)AS_OBJ(args[0]);
    ObjArray* out = allocate_array(src->count < 1 ? 1 : src->count);
    Value ov = OBJ_VAL((Obj*)out);
    for (int i = src->count - 1; i >= 0; i--) {
        koda_array_push(ov, src->elements[i]);
    }
    return ov;
}

Value koda_array_indexOf(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_ARRAY) return NIL_VAL;
    ObjArray* arr = (ObjArray*)AS_OBJ(args[0]);
    Value needle = args[1];
    for (int i = 0; i < arr->count; i++) {
        if (values_equal(arr->elements[i], needle)) {
            return NUMBER_VAL((double)i);
        }
    }
    return NUMBER_VAL(-1);
}

Value koda_array_includes(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_ARRAY) return NIL_VAL;
    ObjArray* arr = (ObjArray*)AS_OBJ(args[0]);
    Value needle = args[1];
    for (int i = 0; i < arr->count; i++) {
        if (values_equal(arr->elements[i], needle)) {
            return BOOL_VAL(true);
        }
    }
    return BOOL_VAL(false);
}

Value koda_array_slice(int arg_count, Value* args) {
    if (arg_count != 3) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_ARRAY) return NIL_VAL;
    if (!IS_NUMBER(args[1]) || !IS_NUMBER(args[2])) return NIL_VAL;

    ObjArray* src = (ObjArray*)AS_OBJ(args[0]);
    int count = src->count;
    int start = (int)AS_NUMBER(args[1]);
    int end = (int)AS_NUMBER(args[2]);

    if (start < 0) start = count + start;
    if (end < 0) end = count + end;
    if (start < 0) start = 0;
    if (end > count) end = count;
    if (end < start) end = start;

    int len = end - start;
    ObjArray* out = allocate_array(len < 1 ? 1 : len);
    Value outv = OBJ_VAL((Obj*)out);
    for (int i = 0; i < len; i++) {
        koda_array_push(outv, src->elements[start + i]);
    }
    return outv;
}

Value koda_array_concat(int arg_count, Value* args) {
    if (arg_count < 1) return NIL_VAL;

    int total = 0;
    for (int a = 0; a < arg_count; a++) {
        if (!IS_OBJ(args[a]) || AS_OBJ(args[a])->type != OBJ_ARRAY) return NIL_VAL;
        ObjArray* src = (ObjArray*)AS_OBJ(args[a]);
        total += src->count;
    }

    ObjArray* out = allocate_array(total < 1 ? 1 : total);
    Value ov = OBJ_VAL((Obj*)out);
    for (int a = 0; a < arg_count; a++) {
        ObjArray* src = (ObjArray*)AS_OBJ(args[a]);
        for (int j = 0; j < src->count; j++) {
            koda_array_push(ov, src->elements[j]);
        }
    }
    return ov;
}

Value koda_array_join(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_ARRAY) return NIL_VAL;
    Value sepv = koda_string(args[1]);
    if (!IS_OBJ(sepv) || AS_OBJ(sepv)->type != OBJ_STRING) return NIL_VAL;
    ObjArray* arr = (ObjArray*)AS_OBJ(args[0]);
    ObjString* sep = (ObjString*)AS_OBJ(sepv);
    Value acc = koda_copy_string("", 0);
    for (int i = 0; i < arr->count; i++) {
        if (i > 0) {
            acc = koda_string_concat(acc, sepv);
        }
        Value part = koda_string(arr->elements[i]);
        acc = koda_string_concat(acc, part);
    }
    return acc;
}

static void push_substring(Value arrv, const char* start, int len) {
    Value chunk = koda_copy_string(start, len);
    koda_array_push(arrv, chunk);
}

// String methods for game development
Value koda_string_split(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(args[1]) || AS_OBJ(args[1])->type != OBJ_STRING) return NIL_VAL;

    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    ObjString* delim = (ObjString*)AS_OBJ(args[1]);

    if (delim->length == 0) {
        Value ov = OBJ_VAL((Obj*)allocate_array(s->length < 1 ? 1 : s->length));
        for (int i = 0; i < s->length; i++) {
            push_substring(ov, &s->chars[i], 1);
        }
        return ov;
    }

    Value ov = OBJ_VAL((Obj*)allocate_array(8));
    int pos = 0;
    const char* hay = s->chars;
    int sl = s->length;
    int dl = delim->length;

    while (pos <= sl) {
        int found = -1;
        for (int i = pos; i <= sl - dl; i++) {
            if (memcmp(hay + i, delim->chars, (size_t)dl) == 0) {
                found = i;
                break;
            }
        }
        if (found < 0) {
            push_substring(ov, hay + pos, sl - pos);
            break;
        }
        push_substring(ov, hay + pos, found - pos);
        pos = found + dl;
    }
    return ov;
}

Value koda_string_trim(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    int lo = 0;
    int hi = s->length;
    while (lo < hi && isspace((unsigned char)s->chars[lo])) {
        lo++;
    }
    while (hi > lo && isspace((unsigned char)s->chars[hi - 1])) {
        hi--;
    }
    return koda_copy_string(s->chars + lo, hi - lo);
}

Value koda_string_upper(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    ObjString* out = allocate_string(s->length);
    for (int i = 0; i < s->length; i++) {
        unsigned char c = (unsigned char)s->chars[i];
        out->chars[i] = (char)toupper((int)c);
    }
    return OBJ_VAL((Obj*)out);
}

Value koda_string_lower(int arg_count, Value* args) {
    if (arg_count != 1) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    ObjString* out = allocate_string(s->length);
    for (int i = 0; i < s->length; i++) {
        unsigned char c = (unsigned char)s->chars[i];
        out->chars[i] = (char)tolower((int)c);
    }
    return OBJ_VAL((Obj*)out);
}

Value koda_string_startsWith(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(args[1]) || AS_OBJ(args[1])->type != OBJ_STRING) return NIL_VAL;
    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    ObjString* pre = (ObjString*)AS_OBJ(args[1]);
    if (pre->length > s->length) return BOOL_VAL(false);
    return BOOL_VAL(memcmp(s->chars, pre->chars, (size_t)pre->length) == 0);
}

Value koda_string_endsWith(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(args[1]) || AS_OBJ(args[1])->type != OBJ_STRING) return NIL_VAL;
    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    ObjString* suf = (ObjString*)AS_OBJ(args[1]);
    if (suf->length > s->length) return BOOL_VAL(false);
    int off = s->length - suf->length;
    return BOOL_VAL(memcmp(s->chars + off, suf->chars, (size_t)suf->length) == 0);
}

Value koda_string_indexOf(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(args[1]) || AS_OBJ(args[1])->type != OBJ_STRING) return NIL_VAL;
    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    ObjString* needle = (ObjString*)AS_OBJ(args[1]);
    if (needle->length == 0) return NUMBER_VAL(0);
    if (needle->length > s->length) return NUMBER_VAL(-1);
    for (int i = 0; i <= s->length - needle->length; i++) {
        if (memcmp(s->chars + i, needle->chars, (size_t)needle->length) == 0) {
            return NUMBER_VAL((double)i);
        }
    }
    return NUMBER_VAL(-1);
}

Value koda_string_slice(int arg_count, Value* args) {
    if (arg_count != 3) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    if (!IS_NUMBER(args[1]) || !IS_NUMBER(args[2])) return NIL_VAL;

    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    int len = s->length;
    int start = (int)AS_NUMBER(args[1]);
    int end = (int)AS_NUMBER(args[2]);

    if (start < 0) start = len + start;
    if (end < 0) end = len + end;
    if (start < 0) start = 0;
    if (end > len) end = len;
    if (end < start) end = start;

    int span = end - start;
    return koda_copy_string(s->chars + start, span);
}

Value koda_string_replace(int arg_count, Value* args) {
    if (arg_count != 3) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(args[1]) || AS_OBJ(args[1])->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(args[2]) || AS_OBJ(args[2])->type != OBJ_STRING) return NIL_VAL;

    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    ObjString* from = (ObjString*)AS_OBJ(args[1]);
    ObjString* to = (ObjString*)AS_OBJ(args[2]);

    if (from->length == 0) return args[0];

    int idx = -1;
    for (int i = 0; i <= s->length - from->length; i++) {
        if (memcmp(s->chars + i, from->chars, (size_t)from->length) == 0) {
            idx = i;
            break;
        }
    }
    if (idx < 0) return args[0];

    int newLen = s->length - from->length + to->length;
    ObjString* out = allocate_string(newLen);
    memcpy(out->chars, s->chars, (size_t)idx);
    memcpy(out->chars + idx, to->chars, (size_t)to->length);
    memcpy(out->chars + idx + to->length, s->chars + idx + from->length,
           (size_t)(s->length - idx - from->length));
    return OBJ_VAL((Obj*)out);
}

Value koda_string_replaceAll(int arg_count, Value* args) {
    if (arg_count != 3) return NIL_VAL;
    Value cur = args[0];
    if (!IS_OBJ(cur) || AS_OBJ(cur)->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(args[1]) || AS_OBJ(args[1])->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(args[2]) || AS_OBJ(args[2])->type != OBJ_STRING) return NIL_VAL;

    ObjString* from = (ObjString*)AS_OBJ(args[1]);
    if (from->length == 0) return cur;

    while (true) {
        Value args3[3] = { cur, args[1], args[2] };
        Value next = koda_string_replace(3, args3);
        if (values_equal(next, cur)) {
            break;
        }
        cur = next;
    }
    return cur;
}

/** Literal substring match after string coercion (`matches(haystack, pattern)`). */
Value koda_matches(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    Value a = koda_string(args[0]);
    Value b = koda_string(args[1]);
    if (!IS_OBJ(a) || AS_OBJ(a)->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(b) || AS_OBJ(b)->type != OBJ_STRING) return NIL_VAL;
    ObjString* hay = (ObjString*)AS_OBJ(a);
    ObjString* needle = (ObjString*)AS_OBJ(b);
    if (needle->length == 0) return TRUE_VAL;
    char* hit = strstr(hay->chars, needle->chars);
    return BOOL_VAL(hit != NULL);
}

Value koda_keys(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_OBJ(args[0])) {
        return NIL_VAL;
    }
    Obj* o = AS_OBJ(args[0]);
    if (o->type != OBJ_TABLE) {
        return NIL_VAL;
    }
    ObjTable* t = (ObjTable*)o;
    Value arr = koda_allocate_array(t->count);
    for (int i = 0; i < t->count; i++) {
        koda_array_push(arr, t->keys[i]);
    }
    return arr;
}

// File I/O functions
#ifdef _WIN32
#include <sys/stat.h>
#endif

static bool koda_stat_path(const char* path, struct KODA_STAT* st) {
    return KODA_STAT_FUNC(path, st) == 0;
}

#ifdef _WIN32
#define KODA_ISREG(m) ((m) & _S_IFREG)
#define KODA_ISDIR(m) ((m) & _S_IFDIR)
#else
#define KODA_ISREG(m) S_ISREG(m)
#define KODA_ISDIR(m) S_ISDIR(m)
#endif

Value koda_readFile(int arg_count, Value* args) {
    if (arg_count < 1 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) {
        return koda_err_str("readFile requires a string path");
    }
    ObjString* path = (ObjString*)AS_OBJ(args[0]);

    FILE* f = fopen(path->chars, "rb");
    if (!f) {
        char msg[512];
        snprintf(msg, sizeof(msg), "could not open '%s': %s", path->chars, strerror(errno));
        return koda_err_str(msg);
    }
    if (fseek(f, 0, SEEK_END) != 0) {
        fclose(f);
        return koda_err_str("could not seek file");
    }
    long size = ftell(f);
    if (size < 0) {
        fclose(f);
        return koda_err_str("could not determine file size");
    }
    rewind(f);
    if (size == 0) {
        fclose(f);
        Value content = koda_copy_string("", 0);
        Value inner[1] = { content };
        return koda_ok(1, inner);
    }
    char* buf = (char*)malloc((size_t)size + 1u);
    if (!buf) {
        fclose(f);
        return koda_err_str("out of memory reading file");
    }
    size_t got = fread(buf, 1, (size_t)size, f);
    fclose(f);
    if (got != (size_t)size) {
        free(buf);
        return koda_err_str("could not read full file");
    }
    buf[size] = '\0';
    Value content = koda_copy_string(buf, (int)size);
    free(buf);
    Value inner[1] = { content };
    return koda_ok(1, inner);
}

Value koda_writeFile(int arg_count, Value* args) {
    if (arg_count < 2 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) {
        return koda_err_str("writeFile requires path string and content string");
    }
    if (!IS_OBJ(args[1]) || AS_OBJ(args[1])->type != OBJ_STRING) {
        return koda_err_str("writeFile requires path string and content string");
    }
    ObjString* path = (ObjString*)AS_OBJ(args[0]);
    ObjString* body = (ObjString*)AS_OBJ(args[1]);

    FILE* f = fopen(path->chars, "wb");
    if (!f) {
        char msg[512];
        snprintf(msg, sizeof(msg), "could not open '%s' for write: %s", path->chars, strerror(errno));
        return koda_err_str(msg);
    }
    if (body->length > 0) {
        size_t w = fwrite(body->chars, 1, (size_t)body->length, f);
        if (w != (size_t)body->length) {
            fclose(f);
            return koda_err_str("writeFile: short write");
        }
    }
    if (fclose(f) != 0) {
        return koda_err_str("writeFile: fclose failed");
    }
    Value inner[1] = { NIL_VAL };
    return koda_ok(1, inner);
}

Value koda_appendFile(int arg_count, Value* args) {
    if (arg_count != 2) return NIL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    if (!IS_OBJ(args[1]) || AS_OBJ(args[1])->type != OBJ_STRING) return NIL_VAL;

    ObjString* path = (ObjString*)AS_OBJ(args[0]);
    ObjString* body = (ObjString*)AS_OBJ(args[1]);
    FILE* f = fopen(path->chars, "ab");
    if (!f) {
        return BOOL_VAL(false);
    }
    if (body->length > 0) {
        size_t w = fwrite(body->chars, 1, (size_t)body->length, f);
        if (w != (size_t)body->length) {
            fclose(f);
            return BOOL_VAL(false);
        }
    }
    fclose(f);
    return BOOL_VAL(true);
}

Value koda_fileExists(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    ObjString* path = (ObjString*)AS_OBJ(args[0]);
    struct KODA_STAT st;
    return BOOL_VAL(koda_stat_path(path->chars, &st));
}

Value koda_deleteFile(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    ObjString* path = (ObjString*)AS_OBJ(args[0]);
    return BOOL_VAL(remove(path->chars) == 0);
}

Value koda_is_file(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    ObjString* path = (ObjString*)AS_OBJ(args[0]);
    struct KODA_STAT st;
    if (!koda_stat_path(path->chars, &st)) {
        return BOOL_VAL(false);
    }
    return BOOL_VAL(KODA_ISREG(st.st_mode));
}

Value koda_is_dir(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    ObjString* path = (ObjString*)AS_OBJ(args[0]);
    struct KODA_STAT st;
    if (!koda_stat_path(path->chars, &st)) {
        return BOOL_VAL(false);
    }
    return BOOL_VAL(KODA_ISDIR(st.st_mode));
}

Value koda_file_size(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    ObjString* path = (ObjString*)AS_OBJ(args[0]);
    struct KODA_STAT st;
    if (!koda_stat_path(path->chars, &st)) {
        return NUMBER_VAL(0);
    }
    return NUMBER_VAL((double)st.st_size);
}

Value koda_list_dir(int arg_count, Value* args) {
    if (arg_count != 1 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NIL_VAL;
    ObjString* dir = (ObjString*)AS_OBJ(args[0]);
    Value arr = koda_allocate_array(0);
#ifdef _WIN32
    char pattern[1024];
    snprintf(pattern, sizeof(pattern), "%s\\*", dir->chars);
    WIN32_FIND_DATAA fd;
    HANDLE h = FindFirstFileA(pattern, &fd);
    if (h == INVALID_HANDLE_VALUE) {
        return NIL_VAL;
    }
    do {
        if (strcmp(fd.cFileName, ".") == 0 || strcmp(fd.cFileName, "..") == 0) {
            continue;
        }
        int n = (int)strlen(fd.cFileName);
        koda_array_push(arr, koda_copy_string(fd.cFileName, n));
    } while (FindNextFileA(h, &fd));
    FindClose(h);
#else
    DIR* d = opendir(dir->chars);
    if (!d) {
        return NIL_VAL;
    }
    struct dirent* ent;
    while ((ent = readdir(d)) != NULL) {
        if (strcmp(ent->d_name, ".") == 0 || strcmp(ent->d_name, "..") == 0) {
            continue;
        }
        int n = (int)strlen(ent->d_name);
        koda_array_push(arr, koda_copy_string(ent->d_name, n));
    }
    closedir(d);
#endif
    return arr;
}

// Debug utilities (legacy LLVM helper; uses truthy + panic)
void koda_assert_llvm(Value cond, Value msg) {
    Value a[1] = { cond };
    Value b = koda_bool(1, a);
    if (IS_TRUE(b)) {
        return;
    }
    if (IS_OBJ(msg) && AS_OBJ(msg)->type == OBJ_STRING) {
        Value pargs[1] = { msg };
        koda_panic(1, pargs);
    }
    koda_panic_str("assertion failed");
}

Value koda_trace(int arg_count, Value* args) {
    for (int i = 0; i < arg_count; i++) {
        print_value(args[i]);
        if (i < arg_count - 1) printf(" ");
    }
    printf("\n");
    return NIL_VAL;
}

// JSON parsing (subset) — returns ok(value) or err(message)
typedef struct {
    const char* current;
    char err[256];
} KodaJsonParser;

static void koda_json_skip_ws(KodaJsonParser* p) {
    while (*p->current == ' ' || *p->current == '\t' || *p->current == '\n' || *p->current == '\r') {
        p->current++;
    }
}

static void koda_json_set_err(KodaJsonParser* p, const char* msg) {
    if (p->err[0] != '\0') {
        return;
    }
    snprintf(p->err, sizeof(p->err), "%s", msg);
}

static Value koda_json_parse_value(KodaJsonParser* p);

static Value koda_json_parse_string(KodaJsonParser* p) {
    if (*p->current != '"') {
        koda_json_set_err(p, "expected '\"' for string");
        return NIL_VAL;
    }
    p->current++;
    const char* start = p->current;
    while (*p->current != '"' && *p->current != '\0') {
        p->current++;
    }
    if (*p->current != '"') {
        koda_json_set_err(p, "unterminated string");
        return NIL_VAL;
    }
    int len = (int)(p->current - start);
    Value sv = koda_copy_string(start, len);
    p->current++;
    return sv;
}

static Value koda_json_parse_array(KodaJsonParser* p) {
    if (*p->current != '[') {
        koda_json_set_err(p, "expected '['");
        return NIL_VAL;
    }
    p->current++;
    Value arrv = koda_allocate_array(0);
    koda_json_skip_ws(p);
    while (*p->current != ']' && *p->current != '\0') {
        Value elem = koda_json_parse_value(p);
        if (p->err[0]) {
            return NIL_VAL;
        }
        koda_array_push(arrv, elem);
        koda_json_skip_ws(p);
        if (*p->current == ',') {
            p->current++;
            koda_json_skip_ws(p);
        }
    }
    if (*p->current != ']') {
        koda_json_set_err(p, "expected ']'");
        return NIL_VAL;
    }
    p->current++;
    return arrv;
}

static Value koda_json_parse_object(KodaJsonParser* p) {
    if (*p->current != '{') {
        koda_json_set_err(p, "expected '{'");
        return NIL_VAL;
    }
    p->current++;
    Value objv = koda_allocate_object(64);
    koda_json_skip_ws(p);
    while (*p->current != '}' && *p->current != '\0') {
        Value key = koda_json_parse_string(p);
        if (p->err[0]) {
            return NIL_VAL;
        }
        koda_json_skip_ws(p);
        if (*p->current != ':') {
            koda_json_set_err(p, "expected ':' in object");
            return NIL_VAL;
        }
        p->current++;
        koda_json_skip_ws(p);
        Value val = koda_json_parse_value(p);
        if (p->err[0]) {
            return NIL_VAL;
        }
        koda_object_set(objv, key, val);
        koda_json_skip_ws(p);
        if (*p->current == ',') {
            p->current++;
            koda_json_skip_ws(p);
        }
    }
    if (*p->current != '}') {
        koda_json_set_err(p, "expected '}'");
        return NIL_VAL;
    }
    p->current++;
    return objv;
}

static Value koda_json_parse_value(KodaJsonParser* p) {
    koda_json_skip_ws(p);
    char c = *p->current;
    if (c == '\0') {
        koda_json_set_err(p, "unexpected end of input");
        return NIL_VAL;
    }
    if (c == '"') {
        return koda_json_parse_string(p);
    }
    if (c == '[') {
        return koda_json_parse_array(p);
    }
    if (c == '{') {
        return koda_json_parse_object(p);
    }
    if (c == 't') {
        if (strncmp(p->current, "true", 4) == 0) {
            p->current += 4;
            return TRUE_VAL;
        }
        koda_json_set_err(p, "invalid literal");
        return NIL_VAL;
    }
    if (c == 'f') {
        if (strncmp(p->current, "false", 5) == 0) {
            p->current += 5;
            return FALSE_VAL;
        }
        koda_json_set_err(p, "invalid literal");
        return NIL_VAL;
    }
    if (c == 'n') {
        if (strncmp(p->current, "null", 4) == 0) {
            p->current += 4;
            return NIL_VAL;
        }
        koda_json_set_err(p, "invalid literal");
        return NIL_VAL;
    }
    if ((c >= '0' && c <= '9') || c == '-') {
        char* end = NULL;
        double d = strtod(p->current, &end);
        if (end == p->current) {
            koda_json_set_err(p, "invalid number");
            return NIL_VAL;
        }
        p->current = end;
        return NUMBER_VAL(d);
    }
    koda_json_set_err(p, "unexpected character");
    return NIL_VAL;
}

Value koda_json_parse(int arg_count, Value* args) {
    if (arg_count < 1 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) {
        return NIL_VAL;
    }
    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    KodaJsonParser p;
    p.current = s->chars;
    p.err[0] = '\0';
    Value val = koda_json_parse_value(&p);
    if (p.err[0]) {
        return NIL_VAL;
    }
    koda_json_skip_ws(&p);
    if (*p.current != '\0') {
        return NIL_VAL;
    }
    return val;
}

Value koda_parseJSON(int arg_count, Value* args) {
    if (arg_count < 1 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) {
        return koda_err_str("parseJSON requires a string");
    }
    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    KodaJsonParser p;
    p.current = s->chars;
    p.err[0] = '\0';
    Value val = koda_json_parse_value(&p);
    if (p.err[0]) {
        return koda_err_str(p.err);
    }
    koda_json_skip_ws(&p);
    if (*p.current != '\0') {
        return koda_err_str("parseJSON: trailing data after value");
    }
    Value inner[1] = { val };
    return koda_ok(1, inner);
}

Value koda_json_try_parse(int arg_count, Value* args) {
    if (arg_count < 1 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) {
        return koda_err_str("try_parse requires a string");
    }
    ObjString* s = (ObjString*)AS_OBJ(args[0]);
    KodaJsonParser p;
    p.current = s->chars;
    p.err[0] = '\0';
    Value val = koda_json_parse_value(&p);
    if (p.err[0]) {
        return koda_err_str(p.err);
    }
    koda_json_skip_ws(&p);
    if (*p.current != '\0') {
        return koda_err_str("try_parse: trailing data after value");
    }
    Value inner[1] = { val };
    return koda_ok(1, inner);
}

typedef struct {
    char* data;
    int length;
    int capacity;
} KodaJsonBuf;

static void koda_json_buf_reserve(KodaJsonBuf* b, int need) {
    if (need <= b->capacity) {
        return;
    }
    int cap = b->capacity > 0 ? b->capacity : 64;
    while (cap < need) {
        cap *= 2;
    }
    char* next = (char*)realloc(b->data, (size_t)cap);
    if (!next) {
        return;
    }
    b->data = next;
    b->capacity = cap;
}

static void koda_json_buf_append(KodaJsonBuf* b, const char* s, int len) {
    if (!s || len < 0) {
        return;
    }
    koda_json_buf_reserve(b, b->length + len + 1);
    if (!b->data) {
        return;
    }
    memcpy(b->data + b->length, s, (size_t)len);
    b->length += len;
    b->data[b->length] = '\0';
}

static void koda_json_buf_append_cstr(KodaJsonBuf* b, const char* s) {
    if (!s) {
        return;
    }
    koda_json_buf_append(b, s, (int)strlen(s));
}

static void koda_json_buf_append_char(KodaJsonBuf* b, char c) {
    koda_json_buf_append(b, &c, 1);
}

static void koda_json_escape_string(KodaJsonBuf* out, ObjString* s) {
    koda_json_buf_append_char(out, '"');
    for (int i = 0; i < s->length; i++) {
        char c = s->chars[i];
        switch (c) {
            case '"': koda_json_buf_append_cstr(out, "\\\""); break;
            case '\\': koda_json_buf_append_cstr(out, "\\\\"); break;
            case '\n': koda_json_buf_append_cstr(out, "\\n"); break;
            case '\r': koda_json_buf_append_cstr(out, "\\r"); break;
            case '\t': koda_json_buf_append_cstr(out, "\\t"); break;
            default: koda_json_buf_append_char(out, c); break;
        }
    }
    koda_json_buf_append_char(out, '"');
}

typedef struct {
    KodaJsonBuf* out;
    int indent_step;
    int depth;
} KodaJsonFmt;

static void koda_json_indent(KodaJsonFmt* fmt) {
    if (fmt->indent_step <= 0) {
        return;
    }
    koda_json_buf_append_char(fmt->out, '\n');
    int spaces = fmt->depth * fmt->indent_step;
    for (int i = 0; i < spaces; i++) {
        koda_json_buf_append_char(fmt->out, ' ');
    }
}

static void koda_json_stringify_value(KodaJsonFmt* fmt, Value v);

static void koda_json_stringify_array(KodaJsonFmt* fmt, ObjArray* a) {
    koda_json_buf_append_char(fmt->out, '[');
    if (fmt->indent_step > 0 && a->count > 0) {
        fmt->depth++;
        koda_json_indent(fmt);
    }
    for (int i = 0; i < a->count; i++) {
        if (i > 0) {
            koda_json_buf_append_char(fmt->out, ',');
            if (fmt->indent_step > 0) {
                koda_json_indent(fmt);
            }
        }
        koda_json_stringify_value(fmt, a->elements[i]);
    }
    if (fmt->indent_step > 0 && a->count > 0) {
        fmt->depth--;
        koda_json_indent(fmt);
    }
    koda_json_buf_append_char(fmt->out, ']');
}

static void koda_json_stringify_object(KodaJsonFmt* fmt, ObjTable* t) {
    koda_json_buf_append_char(fmt->out, '{');
    if (fmt->indent_step > 0 && t->count > 0) {
        fmt->depth++;
        koda_json_indent(fmt);
    }
    for (int i = 0; i < t->count; i++) {
        if (i > 0) {
            koda_json_buf_append_char(fmt->out, ',');
            if (fmt->indent_step > 0) {
                koda_json_indent(fmt);
            }
        }
        Value key = t->keys[i];
        if (IS_OBJ(key) && AS_OBJ(key)->type == OBJ_STRING) {
            koda_json_escape_string(fmt->out, (ObjString*)AS_OBJ(key));
        } else {
            koda_json_stringify_value(fmt, key);
        }
        if (fmt->indent_step > 0) {
            koda_json_buf_append_cstr(fmt->out, ": ");
        } else {
            koda_json_buf_append_char(fmt->out, ':');
        }
        koda_json_stringify_value(fmt, t->values[i]);
    }
    if (fmt->indent_step > 0 && t->count > 0) {
        fmt->depth--;
        koda_json_indent(fmt);
    }
    koda_json_buf_append_char(fmt->out, '}');
}

static void koda_json_stringify_value(KodaJsonFmt* fmt, Value v) {
    KodaJsonBuf* out = fmt->out;
    if (IS_NIL(v)) {
        koda_json_buf_append_cstr(out, "null");
        return;
    }
    if (IS_BOOL(v)) {
        koda_json_buf_append_cstr(out, AS_BOOL(v) ? "true" : "false");
        return;
    }
    if (IS_NUMBER(v)) {
        char numbuf[64];
        snprintf(numbuf, sizeof(numbuf), "%.17g", AS_NUMBER(v));
        koda_json_buf_append_cstr(out, numbuf);
        return;
    }
    if (!IS_OBJ(v)) {
        koda_json_buf_append_cstr(out, "null");
        return;
    }
    Obj* o = AS_OBJ(v);
    switch (o->type) {
        case OBJ_STRING:
            koda_json_escape_string(out, (ObjString*)o);
            break;
        case OBJ_ARRAY:
            koda_json_stringify_array(fmt, (ObjArray*)o);
            break;
        case OBJ_TABLE:
            koda_json_stringify_object(fmt, (ObjTable*)o);
            break;
        default:
            koda_json_buf_append_cstr(out, "null");
            break;
    }
}

Value koda_toJSON(int arg_count, Value* args) {
    if (arg_count < 1 || arg_count > 2) {
        return NIL_VAL;
    }
    int indent = 0;
    if (arg_count == 2 && IS_NUMBER(args[1])) {
        indent = (int)AS_NUMBER(args[1]);
        if (indent < 0) {
            indent = 0;
        }
    }
    KodaJsonBuf buf = { NULL, 0, 0 };
    KodaJsonFmt fmt = { &buf, indent, 0 };
    koda_json_stringify_value(&fmt, args[0]);
    Value result = buf.data ? koda_copy_string(buf.data, buf.length) : koda_copy_string("", 0);
    if (buf.data) {
        free(buf.data);
    }
    return result;
}

// Graphics: not linked here (no Raylib). Use koda build + KODA_NATIVE_SOURCES / wrapgen for native games.
Value koda_gfx_init_window(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value koda_gfx_set_target_fps(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value koda_gfx_window_should_close(int arg_count, Value* args) { (void)arg_count; (void)args; return NUMBER_VAL(1); }
Value koda_gfx_close_window(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value koda_gfx_begin_drawing(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value koda_gfx_end_drawing(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value koda_gfx_clear_background(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value koda_gfx_draw_text(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value koda_gfx_draw_rectangle(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value koda_gfx_draw_circle(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value koda_gfx_is_key_down(int arg_count, Value* args) { (void)arg_count; (void)args; return NUMBER_VAL(0); }
Value koda_gfx_is_key_pressed(int arg_count, Value* args) { (void)arg_count; (void)args; return NUMBER_VAL(0); }
Value koda_gfx_begin_mode3d(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value koda_gfx_end_mode3d(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value koda_gfx_draw_grid(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value koda_gfx_draw_cube(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
Value koda_gfx_draw_cube_wires(int arg_count, Value* args) { (void)arg_count; (void)args; return NIL_VAL; }
