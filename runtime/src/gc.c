#include "gc.h"
#include "value.h"
#include "object.h"
#include "shadow_stack.h"
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#ifdef _WIN32
#include <windows.h>
#else
#include <time.h>
#endif

/* Stack base for conservative root scan; set in koda_runtime_init before any allocation. */
extern void* gc_stack_base;

extern Value* koda_globals;
extern int koda_globals_count;
extern Value** koda_global_slots;
extern int koda_global_slots_count;

void koda_mark_module_cache(void);
void koda_mark_open_upvalues(void);
void koda_sweep_intern_table(void);

#define REMEMBERED_SET_MAX 4096
#define NURSERY_SIZE (256u * 1024u)

typedef struct {
    Obj* objects;
    size_t bytes_allocated;
    size_t next_gc;
    bool gc_disabled;
    bool collecting;
    bool use_shadow_stack;
    GCStats stats;

    uint8_t nursery_buf[NURSERY_SIZE];
    uint8_t* nursery_top;
    uint8_t* nursery_end;
    size_t nursery_live_bytes;

    Obj* remembered[REMEMBERED_SET_MAX];
    int remembered_count;
    uint64_t remembered_overflow_count;
} GCState;

static GCState gc_state = {
    .objects = NULL,
    .bytes_allocated = 0,
    .next_gc = 1024u * 1024u,
    .gc_disabled = false,
    .collecting = false,
    .use_shadow_stack = false,
    .nursery_top = NULL,
    .nursery_end = NULL,
    .nursery_live_bytes = 0,
    .remembered_count = 0,
    .remembered_overflow_count = 0,
};

typedef enum {
    GC_INC_IDLE = 0,
    GC_INC_MARKING,
    GC_INC_SWEEPING,
} GCIncPhase;

static GCIncPhase gc_inc_phase = GC_INC_IDLE;
static Obj** gc_grey = NULL;
static size_t gc_grey_len;
static size_t gc_grey_cap;
static Obj** gc_inc_sweep_prev = NULL;

#ifdef _WIN32
static LARGE_INTEGER gc_win_qpc_freq;
static bool gc_win_qpc_inited;
#endif

static uint64_t gc_now_us(void) {
#ifdef _WIN32
    if (!gc_win_qpc_inited) {
        QueryPerformanceFrequency(&gc_win_qpc_freq);
        gc_win_qpc_inited = true;
    }
    LARGE_INTEGER c;
    QueryPerformanceCounter(&c);
    if (gc_win_qpc_freq.QuadPart == 0) {
        return 0;
    }
    return (uint64_t)((double)c.QuadPart * 1000000.0 / (double)gc_win_qpc_freq.QuadPart);
#else
    struct timespec ts;
#if defined(CLOCK_MONOTONIC)
    if (clock_gettime(CLOCK_MONOTONIC, &ts) != 0) {
        return 0;
    }
#else
    if (clock_gettime(CLOCK_REALTIME, &ts) != 0) {
        return 0;
    }
#endif
    return (uint64_t)ts.tv_sec * 1000000ull + (uint64_t)ts.tv_nsec / 1000ull;
#endif
}

void gc_record_pause_since(uint64_t start_us) {
    uint64_t now = gc_now_us();
    if (now <= start_us) {
        return;
    }
    uint64_t dt = now - start_us;
    gc_state.stats.total_pause_time_us += dt;
    if (dt > gc_state.stats.max_pause_time_us) {
        gc_state.stats.max_pause_time_us = dt;
    }
}

bool gc_incremental_is_idle(void) {
    return gc_inc_phase == GC_INC_IDLE;
}

static bool gc_is_heap_object_exact(Obj* candidate);
static Obj* gc_find_containing_obj(void* p);
static void gc_build_heap_lookup_cache(void);
static void gc_heap_cache_invalidate(void);
static void gc_unmark_all(void);

static void gc_incremental_abort(void) {
    gc_grey_len = 0;
    gc_inc_phase = GC_INC_IDLE;
    gc_inc_sweep_prev = NULL;
    gc_unmark_all();
}

static void grey_ensure_cap(size_t need) {
    if (need <= gc_grey_cap) {
        return;
    }
    size_t ncap = gc_grey_cap ? gc_grey_cap * 2u : 256u;
    while (ncap < need) {
        ncap *= 2u;
    }
    Obj** ng = (Obj**)realloc(gc_grey, ncap * sizeof(Obj*));
    if (ng == NULL) {
        fprintf(stderr, "koda: out of memory\n");
        exit(1);
    }
    gc_grey = ng;
    gc_grey_cap = ncap;
}

static void grey_push_obj(Obj* obj) {
    if (obj == NULL || obj->generation == GEN_DEAD) {
        return;
    }
    if (obj->is_marked) {
        return;
    }
    obj->is_marked = true;
    grey_ensure_cap(gc_grey_len + 1u);
    gc_grey[gc_grey_len++] = obj;
}

static void grey_push_from_value(Value v) {
    if (!IS_OBJ(v)) {
        return;
    }
    Obj* ch = AS_OBJ(v);
    if (ch == NULL) {
        return;
    }
    grey_push_obj(ch);
}

static void gc_visit_obj_edges(Obj* obj) {
    switch (obj->type) {
        case OBJ_ARRAY: {
            ObjArray* array = (ObjArray*)obj;
            if (array->elements == NULL) {
                break;
            }
            for (int i = 0; i < array->count; i++) {
                grey_push_from_value(array->elements[i]);
            }
            break;
        }
        case OBJ_TABLE: {
            ObjTable* table = (ObjTable*)obj;
            if (table->keys == NULL || table->values == NULL) {
                break;
            }
            if (table->hashes != NULL && !table->is_struct_layout) {
                for (int i = 0; i < table->capacity; i++) {
                    grey_push_from_value(table->keys[i]);
                    grey_push_from_value(table->values[i]);
                }
            } else {
                for (int i = 0; i < table->count; i++) {
                    grey_push_from_value(table->keys[i]);
                    grey_push_from_value(table->values[i]);
                }
            }
            break;
        }
        case OBJ_CLOSURE: {
            ObjClosure* closure = (ObjClosure*)obj;
            grey_push_obj((Obj*)closure->function);
            for (int i = 0; i < closure->upvalue_count; i++) {
                grey_push_from_value(closure->upvalues[i]);
            }
            break;
        }
        case OBJ_CELL: {
            ObjCell* cell = (ObjCell*)obj;
            grey_push_from_value(cell->value);
            break;
        }
        default:
            break;
    }
}

static void maybe_enqueue_stack_word(uint64_t word) {
    Value v = word;
    if (!IS_OBJ(v)) {
        return;
    }
    Obj* o = AS_OBJ(v);
    if (gc_is_heap_object_exact(o)) {
        grey_push_obj(o);
    } else {
        Obj* c = gc_find_containing_obj(o);
        if (c != NULL) {
            grey_push_obj(c);
        }
    }
}

static void maybe_enqueue_stack_addr(uintptr_t addr) {
    if (addr == 0) {
        return;
    }
    Obj* o = gc_find_containing_obj((void*)addr);
    if (o != NULL) {
        grey_push_obj(o);
    }
}

static void gc_mark_stack_conservative_grey(void) {
    if (gc_stack_base == NULL) {
        return;
    }
    gc_heap_cache_invalidate();
    void* stack_top;
    GC_GET_SP(stack_top);
    uintptr_t a = (uintptr_t)stack_top;
    uintptr_t b = (uintptr_t)gc_stack_base;
    uintptr_t lo = a < b ? a : b;
    uintptr_t hi = a < b ? b : a;
    lo = (lo + 7u) & ~(uintptr_t)7u;
    for (uintptr_t p = lo; p + sizeof(uint64_t) <= hi; p += sizeof(uint64_t)) {
        uint64_t word = *(uint64_t*)p;
        maybe_enqueue_stack_word(word);
        maybe_enqueue_stack_addr((uintptr_t)word);
    }
}

static void gc_mark_shadow_stack_grey(void) {
    for (int i = 0; i < koda_shadow_depth; i++) {
        Value** ptrs = koda_shadow_stack[i].slot_ptrs;
        int n = koda_shadow_stack[i].count;
        if (ptrs == NULL || n <= 0) {
            continue;
        }
        for (int j = 0; j < n; j++) {
            Value* p = ptrs[j];
            if (p != NULL) {
                grey_push_from_value(*p);
            }
        }
    }
}

static void gc_seed_roots_grey(void) {
    if (gc_state.use_shadow_stack) {
        gc_mark_shadow_stack_grey();
    } else {
        gc_mark_stack_conservative_grey();
    }
    if (koda_globals_count > 0) {
        for (int i = 0; i < koda_globals_count; i++) {
            grey_push_from_value(koda_globals[i]);
        }
    }
    if (koda_global_slots_count > 0) {
        for (int i = 0; i < koda_global_slots_count; i++) {
            if (koda_global_slots[i] != NULL) {
                grey_push_from_value(*koda_global_slots[i]);
            }
        }
    }
    koda_mark_module_cache();
    koda_mark_open_upvalues();
}

static bool gc_incremental_mark_step(void) {
    if (gc_grey_len == 0) {
        return false;
    }
    Obj* o = gc_grey[--gc_grey_len];
    gc_visit_obj_edges(o);
    return true;
}

static bool gc_incremental_sweep_step(void) {
    if (gc_inc_sweep_prev == NULL) {
        return false;
    }
    Obj* obj = *gc_inc_sweep_prev;
    if (obj == NULL) {
        gc_inc_sweep_prev = NULL;
        return false;
    }
    if (!obj->is_marked) {
        *gc_inc_sweep_prev = obj->next;
        free_object(obj);
        return true;
    }
    obj->generation = GEN_OLD;
    obj->is_marked = false;
    gc_inc_sweep_prev = &obj->next;
    return true;
}

static size_t gc_obj_total_bytes(Obj* obj);
static void maybe_mark_stack_word(uint64_t word);
static void maybe_mark_stack_addr(uintptr_t addr);
static bool gc_debug_enabled(void);
static void gc_debug_validate_objects(void);

void gc_collect_minor(void);
static void gc_incremental_finish_sync(void);

static bool obj_header_in_nursery_range(const Obj* obj) {
    const uint8_t* p = (const uint8_t*)obj;
    return p >= gc_state.nursery_buf && p < gc_state.nursery_end;
}

void gc_register_object(Obj* obj) {
    if (obj_header_in_nursery_range(obj)) {
        obj->generation = GEN_NURSERY;
    } else {
        obj->generation = GEN_OLD;
    }
    obj->in_remembered_set = false;
    obj->in_arena = false;
    obj->next = gc_state.objects;
    gc_state.objects = obj;
}

void gc_unlink_object(Obj* victim) {
    if (victim == NULL) {
        return;
    }
    Obj** previous = &gc_state.objects;
    while (*previous != NULL) {
        if (*previous == victim) {
            *previous = victim->next;
            victim->next = NULL;
            return;
        }
        previous = &(*previous)->next;
    }
}

void gc_free(void* ptr, size_t size) {
    if (ptr == NULL) {
        return;
    }
    const uint8_t* p = (const uint8_t*)ptr;
    if (p >= gc_state.nursery_buf && p < gc_state.nursery_end) {
        (void)size;
        return;
    }
    if (size <= gc_state.bytes_allocated) {
        gc_state.bytes_allocated -= size;
    } else {
        gc_state.bytes_allocated = 0;
    }
    free(ptr);
}

static void* nursery_alloc(size_t size) {
    size = (size + 7u) & ~7u;
    if (gc_state.nursery_top == NULL || gc_state.nursery_end == NULL) {
        return NULL;
    }
    if (gc_state.nursery_top + size > gc_state.nursery_end) {
        gc_collect_minor();
        if (gc_state.nursery_top + size > gc_state.nursery_end) {
            return NULL;
        }
    }
    void* ptr = (void*)gc_state.nursery_top;
    gc_state.nursery_top += size;
    gc_state.nursery_live_bytes += size;
    return ptr;
}

void* gc_alloc(size_t size) {
    if (size == 0) {
        size = 1;
    }
    if (gc_state.gc_disabled) {
        void* ptr = malloc(size);
        if (ptr == NULL) {
            fprintf(stderr, "koda: out of memory\n");
            exit(1);
        }
        return ptr;
    }
    if (size <= 4096u && !gc_state.collecting) {
        void* n = nursery_alloc(size);
        if (n != NULL) {
            return n;
        }
    }
    gc_state.bytes_allocated += size;
    if (!gc_state.collecting && gc_state.bytes_allocated > gc_state.next_gc) {
        gc_collect();
    }
    void* ptr = malloc(size);
    if (ptr == NULL) {
        fprintf(stderr, "koda: out of memory\n");
        exit(1);
    }
    return ptr;
}

static void remembered_set_clear(void) {
    for (int i = 0; i < gc_state.remembered_count; i++) {
        if (gc_state.remembered[i] != NULL) {
            gc_state.remembered[i]->in_remembered_set = false;
        }
    }
    gc_state.remembered_count = 0;
}

static void remembered_set_add(Obj* obj) {
    if (obj == NULL || obj->in_remembered_set) {
        return;
    }
    if (gc_state.remembered_count < REMEMBERED_SET_MAX) {
        obj->in_remembered_set = true;
        gc_state.remembered[gc_state.remembered_count++] = obj;
        return;
    }
    gc_state.remembered_overflow_count++;
    gc_collect();
    remembered_set_clear();
    obj->in_remembered_set = true;
    gc_state.remembered[gc_state.remembered_count++] = obj;
}

static void gc_reset_nursery(void) {
    gc_state.nursery_top = gc_state.nursery_buf;
    gc_state.nursery_live_bytes = 0;
}

static void gc_unlink_gen_dead(void) {
    Obj** previous = &gc_state.objects;
    Obj* obj = gc_state.objects;
    while (obj != NULL) {
        if (obj->generation == GEN_DEAD) {
            *previous = obj->next;
            free_object(obj);
            obj = *previous;
        } else {
            previous = &obj->next;
            obj = obj->next;
        }
    }
}

static bool nursery_has_live_object(void) {
    for (Obj* o = gc_state.objects; o != NULL; o = o->next) {
        const uint8_t* p = (const uint8_t*)o;
        if (p >= gc_state.nursery_buf && p < gc_state.nursery_end && o->generation != GEN_DEAD) {
            return true;
        }
    }
    return false;
}

void gc_collect_minor(void) {
    if (gc_state.collecting) {
        return;
    }
    gc_incremental_finish_sync();
    uint64_t pause_start = gc_now_us();
    gc_state.collecting = true;

    gc_unmark_all();
    if (gc_state.use_shadow_stack) {
        gc_mark_shadow_stack();
    } else {
        gc_mark_stack_conservative();
    }
    if (koda_globals_count > 0) {
        for (int i = 0; i < koda_globals_count; i++) {
            gc_mark_value(koda_globals[i]);
        }
    }
    if (koda_global_slots_count > 0) {
        for (int i = 0; i < koda_global_slots_count; i++) {
            if (koda_global_slots[i] != NULL) {
                gc_mark_value(*koda_global_slots[i]);
            }
        }
    }
    koda_mark_module_cache();
    koda_mark_open_upvalues();

    for (int i = 0; i < gc_state.remembered_count; i++) {
        if (gc_state.remembered[i] != NULL) {
            gc_mark_object(gc_state.remembered[i]);
        }
    }

    /* Drop intern-table slots for nursery strings that did not survive this mark phase
     * (while is_marked is still valid), before we repurpose flags in the nursery walk. */
    koda_sweep_intern_table();

    for (Obj* obj = gc_state.objects; obj != NULL; obj = obj->next) {
        if (obj->generation == GEN_NURSERY) {
            if (obj->is_marked) {
                obj->generation = GEN_YOUNG;
            } else {
                obj->generation = GEN_DEAD;
            }
            obj->is_marked = false;
        }
    }

    gc_unlink_gen_dead();
    remembered_set_clear();

    if (!nursery_has_live_object()) {
        gc_reset_nursery();
    }

    gc_unmark_all();
    if (gc_debug_enabled()) {
        gc_debug_validate_objects();
    }
    gc_state.stats.collections++;
    gc_state.collecting = false;
    gc_record_pause_since(pause_start);
}

size_t gc_nursery_used_bytes(void) {
    if (gc_state.nursery_top == NULL) {
        return 0;
    }
    return (size_t)(gc_state.nursery_top - gc_state.nursery_buf);
}

size_t gc_nursery_capacity_bytes(void) {
    return (size_t)sizeof(gc_state.nursery_buf);
}

static size_t gc_obj_total_bytes(Obj* obj) {
    switch (obj->type) {
        case OBJ_STRING:
            return sizeof(ObjString) + (size_t)((ObjString*)obj)->length + 1u;
        case OBJ_ARRAY: {
            ObjArray* a = (ObjArray*)obj;
            return sizeof(ObjArray) + (size_t)a->capacity * sizeof(Value);
        }
        case OBJ_TABLE: {
            ObjTable* t = (ObjTable*)obj;
            size_t extra = 0;
            if (t->hashes != NULL) {
                extra = (size_t)t->capacity * sizeof(uint32_t);
            }
            return sizeof(ObjTable) + (size_t)t->capacity * 2u * sizeof(Value) + extra;
        }
        case OBJ_FUNCTION:
            return sizeof(ObjFunction);
        case OBJ_CLOSURE: {
            ObjClosure* c = (ObjClosure*)obj;
            return sizeof(ObjClosure) + (size_t)c->upvalue_count * sizeof(Value);
        }
        case OBJ_NATIVE:
            return sizeof(ObjNative);
        case OBJ_CELL:
            return sizeof(ObjCell);
        case OBJ_ARENA:
            return sizeof(ObjArena);
        default:
            return sizeof(Obj);
    }
}

void gc_mark_object(Obj* obj) {
    if (obj == NULL || obj->is_marked) {
        return;
    }
    if (obj->generation == GEN_DEAD) {
        return;
    }
    obj->is_marked = true;

    switch (obj->type) {
        case OBJ_ARRAY: {
            ObjArray* array = (ObjArray*)obj;
            /* elements buffer is allocated after the ObjArray header; GC may run in between. */
            if (array->elements == NULL) {
                break;
            }
            for (int i = 0; i < array->count; i++) {
                gc_mark_value(array->elements[i]);
            }
            break;
        }
        case OBJ_TABLE: {
            ObjTable* table = (ObjTable*)obj;
            if (table->keys == NULL || table->values == NULL) {
                break; /* same partial-allocation window as ObjArray.elements */
            }
            if (table->hashes != NULL && !table->is_struct_layout) {
                for (int i = 0; i < table->capacity; i++) {
                    gc_mark_value(table->keys[i]);
                    gc_mark_value(table->values[i]);
                }
            } else {
                for (int i = 0; i < table->count; i++) {
                    gc_mark_value(table->keys[i]);
                    gc_mark_value(table->values[i]);
                }
            }
            break;
        }
        case OBJ_CLOSURE: {
            ObjClosure* closure = (ObjClosure*)obj;
            gc_mark_object((Obj*)closure->function);
            for (int i = 0; i < closure->upvalue_count; i++) {
                gc_mark_value(closure->upvalues[i]);
            }
            break;
        }
        case OBJ_CELL: {
            ObjCell* cell = (ObjCell*)obj;
            gc_mark_value(cell->value);
            break;
        }
        default:
            break;
    }
}

void gc_mark_value(Value v) {
    if (IS_OBJ(v)) {
        gc_mark_object(AS_OBJ(v));
    }
}

void gc_write_barrier(Obj* parent, Value new_val) {
    if (gc_state.collecting || parent == NULL) {
        return;
    }
    if (!IS_OBJ(new_val)) {
        return;
    }
    Obj* child = AS_OBJ(new_val);
    if (child == NULL) {
        return;
    }
    if (parent->generation == GEN_OLD && child->generation != GEN_OLD) {
        remembered_set_add(parent);
    }
}

static int gc_ptr_cache_probe(Obj* key, int cap) {
    uintptr_t h = (uintptr_t)key;
    h ^= h >> 33;
    h *= 0xff51afd7ed558ccdULL;
    h ^= h >> 33;
    return (int)(h & (uint32_t)(cap - 1));
}

typedef struct {
    uintptr_t base;
    uintptr_t end;
    Obj* obj;
} GcObjRange;

static Obj** gc_ptr_cache = NULL;
static int gc_ptr_cache_cap = 0;
static GcObjRange* gc_range_cache = NULL;
static int gc_range_cache_len = 0;
static int gc_range_cache_cap = 0;
static bool gc_heap_cache_built = false;

static void gc_heap_cache_invalidate(void) {
    gc_heap_cache_built = false;
}

static int gc_range_cmp(const void* a, const void* b) {
    uintptr_t ba = ((const GcObjRange*)a)->base;
    uintptr_t bb = ((const GcObjRange*)b)->base;
    if (ba < bb) {
        return -1;
    }
    if (ba > bb) {
        return 1;
    }
    return 0;
}

static void gc_build_heap_lookup_cache(void) {
    if (gc_heap_cache_built) {
        return;
    }
    int live = 0;
    for (Obj* obj = gc_state.objects; obj != NULL; obj = obj->next) {
        if (obj->generation != GEN_DEAD) {
            live++;
        }
    }
    int cap = 256;
    while (cap < live * 2) {
        cap *= 2;
    }
    if (gc_ptr_cache_cap < cap) {
        Obj** next = (Obj**)realloc(gc_ptr_cache, sizeof(Obj*) * (size_t)cap);
        if (next == NULL) {
            fprintf(stderr, "koda: out of memory building GC pointer cache\n");
            exit(1);
        }
        gc_ptr_cache = next;
        gc_ptr_cache_cap = cap;
    }
    memset(gc_ptr_cache, 0, sizeof(Obj*) * (size_t)gc_ptr_cache_cap);

    if (gc_range_cache_cap < live) {
        int next_cap = gc_range_cache_cap == 0 ? 256 : gc_range_cache_cap;
        while (next_cap < live) {
            next_cap *= 2;
        }
        GcObjRange* next = (GcObjRange*)realloc(gc_range_cache, sizeof(GcObjRange) * (size_t)next_cap);
        if (next == NULL) {
            fprintf(stderr, "koda: out of memory building GC range cache\n");
            exit(1);
        }
        gc_range_cache = next;
        gc_range_cache_cap = next_cap;
    }

    int range_n = 0;
    for (Obj* obj = gc_state.objects; obj != NULL; obj = obj->next) {
        if (obj->generation == GEN_DEAD) {
            continue;
        }
        int slot = gc_ptr_cache_probe(obj, gc_ptr_cache_cap);
        for (int p = 0; p < gc_ptr_cache_cap; p++) {
            int i = (slot + p) % gc_ptr_cache_cap;
            if (gc_ptr_cache[i] == NULL) {
                gc_ptr_cache[i] = obj;
                break;
            }
        }
        if (range_n < gc_range_cache_cap) {
            gc_range_cache[range_n].base = (uintptr_t)obj;
            gc_range_cache[range_n].end = (uintptr_t)obj + gc_obj_total_bytes(obj);
            gc_range_cache[range_n].obj = obj;
            range_n++;
        }
    }
    gc_range_cache_len = range_n;
    if (range_n > 1) {
        qsort(gc_range_cache, (size_t)range_n, sizeof(GcObjRange), gc_range_cmp);
    }
    gc_heap_cache_built = true;
}

static bool gc_is_heap_object_exact(Obj* candidate) {
    if (candidate == NULL) {
        return false;
    }
    gc_build_heap_lookup_cache();
    int start = gc_ptr_cache_probe(candidate, gc_ptr_cache_cap);
    for (int p = 0; p < gc_ptr_cache_cap; p++) {
        int i = (start + p) % gc_ptr_cache_cap;
        Obj* slot = gc_ptr_cache[i];
        if (slot == NULL) {
            return false;
        }
        if (slot == candidate) {
            return true;
        }
    }
    return false;
}

static Obj* gc_find_containing_obj(void* p) {
    if (p == NULL) {
        return NULL;
    }
    uintptr_t c = (uintptr_t)p;
    gc_build_heap_lookup_cache();
    int lo = 0;
    int hi = gc_range_cache_len - 1;
    while (lo <= hi) {
        int mid = lo + (hi - lo) / 2;
        GcObjRange* r = &gc_range_cache[mid];
        if (c < r->base) {
            hi = mid - 1;
        } else if (c >= r->end) {
            lo = mid + 1;
        } else {
            return r->obj;
        }
    }
    return NULL;
}

void gc_mark_stack_conservative(void) {
    if (gc_stack_base == NULL) {
        return;
    }
    gc_heap_cache_invalidate();
    void* stack_top;
    GC_GET_SP(stack_top);
    uintptr_t a = (uintptr_t)stack_top;
    uintptr_t b = (uintptr_t)gc_stack_base;
    uintptr_t lo = a < b ? a : b;
    uintptr_t hi = a < b ? b : a;
    lo = (lo + 7u) & ~(uintptr_t)7u;
    for (uintptr_t p = lo; p + sizeof(uint64_t) <= hi; p += sizeof(uint64_t)) {
        uint64_t word = *(uint64_t*)p;
        maybe_mark_stack_word(word);
        maybe_mark_stack_addr((uintptr_t)word);
    }
}

static void maybe_mark_stack_word(uint64_t word) {
    Value v = word;
    if (IS_OBJ(v)) {
        Obj* o = AS_OBJ(v);
        if (gc_is_heap_object_exact(o)) {
            gc_mark_object(o);
        } else {
            Obj* c = gc_find_containing_obj(o);
            if (c != NULL) {
                gc_mark_object(c);
            }
        }
    }
}

static void maybe_mark_stack_addr(uintptr_t addr) {
    if (addr == 0) {
        return;
    }
    Obj* o = gc_find_containing_obj((void*)addr);
    if (o != NULL) {
        gc_mark_object(o);
    }
}

void gc_mark_shadow_stack(void) {
    for (int i = 0; i < koda_shadow_depth; i++) {
        Value** ptrs = koda_shadow_stack[i].slot_ptrs;
        int n = koda_shadow_stack[i].count;
        if (ptrs == NULL || n <= 0) {
            continue;
        }
        for (int j = 0; j < n; j++) {
            Value* p = ptrs[j];
            if (p != NULL) {
                gc_mark_value(*p);
            }
        }
    }
}

void gc_set_use_shadow_stack(bool enabled) {
    gc_state.use_shadow_stack = enabled;
}

bool gc_uses_shadow_stack(void) {
    return gc_state.use_shadow_stack;
}

void gc_mark_roots(void) {
    if (gc_state.use_shadow_stack) {
        gc_mark_shadow_stack();
    } else {
        gc_mark_stack_conservative();
    }
    if (koda_globals_count > 0) {
        for (int i = 0; i < koda_globals_count; i++) {
            gc_mark_value(koda_globals[i]);
        }
    }
    if (koda_global_slots_count > 0) {
        for (int i = 0; i < koda_global_slots_count; i++) {
            if (koda_global_slots[i] != NULL) {
                gc_mark_value(*koda_global_slots[i]);
            }
        }
    }
    koda_mark_module_cache();
    koda_mark_open_upvalues();
}

static void gc_unmark_all(void) {
    for (Obj* obj = gc_state.objects; obj != NULL; obj = obj->next) {
        obj->is_marked = false;
    }
}

void gc_sweep(void) {
    Obj** previous = &gc_state.objects;
    Obj* obj = gc_state.objects;
    while (obj != NULL) {
        if (!obj->is_marked) {
            *previous = obj->next;
            free_object(obj);
            obj = *previous;
        } else {
            previous = &obj->next;
            obj = obj->next;
        }
    }
}

static size_t gc_inc_bytes_before_sweep = 0;

void gc_collect_incremental(uint64_t budget_us) {
    if (gc_state.gc_disabled) {
        return;
    }
    uint64_t deadline;
    if (budget_us == UINT64_MAX) {
        deadline = UINT64_MAX;
    } else {
        uint64_t now = gc_now_us();
        deadline = now + budget_us;
    }
    while (gc_now_us() < deadline) {
        if (gc_inc_phase == GC_INC_IDLE) {
            size_t soft = (gc_state.next_gc / 4u) * 3u;
            if (soft == 0u) {
                soft = gc_state.next_gc;
            }
            if (gc_state.bytes_allocated < soft) {
                return;
            }
            gc_unmark_all();
            gc_inc_bytes_before_sweep = gc_state.bytes_allocated;
            gc_seed_roots_grey();
            for (int i = 0; i < gc_state.remembered_count; i++) {
                if (gc_state.remembered[i] != NULL) {
                    grey_push_obj(gc_state.remembered[i]);
                }
            }
            gc_inc_phase = GC_INC_MARKING;
            continue;
        }
        if (gc_inc_phase == GC_INC_MARKING) {
            if (gc_incremental_mark_step()) {
                continue;
            }
            koda_sweep_intern_table();
            gc_inc_phase = GC_INC_SWEEPING;
            gc_inc_sweep_prev = &gc_state.objects;
            continue;
        }
        if (gc_inc_phase == GC_INC_SWEEPING) {
            if (gc_incremental_sweep_step()) {
                continue;
            }
            remembered_set_clear();
            gc_state.next_gc = gc_state.bytes_allocated * 2;
            if (gc_state.next_gc < 1024u * 1024u) {
                gc_state.next_gc = 1024u * 1024u;
            }
            size_t freed = 0;
            if (gc_inc_bytes_before_sweep > gc_state.bytes_allocated) {
                freed = gc_inc_bytes_before_sweep - gc_state.bytes_allocated;
            }
            gc_state.stats.bytes_freed += freed;
            if (!nursery_has_live_object()) {
                gc_reset_nursery();
            }
            gc_inc_phase = GC_INC_IDLE;
            gc_state.stats.collections++;
            gc_state.stats.bytes_allocated = gc_state.bytes_allocated;
            return;
        }
    }
}

static void gc_incremental_finish_sync(void) {
    while (gc_inc_phase != GC_INC_IDLE) {
        gc_collect_incremental(UINT64_MAX);
    }
}

void gc_frame_step_incremental(uint64_t budget_us) {
    uint64_t t0 = gc_now_us();
    gc_collect_incremental(budget_us);
    gc_record_pause_since(t0);
}

void gc_collect(void) {
    if (gc_state.collecting) {
        return;
    }
    gc_incremental_abort();
    uint64_t pause_start = gc_now_us();
    gc_state.collecting = true;
    size_t before = gc_state.bytes_allocated;

    gc_unmark_all();
    gc_mark_roots();
    if (gc_debug_enabled()) {
        gc_debug_validate_objects();
    }
    koda_sweep_intern_table();
    gc_sweep();

    for (Obj* obj = gc_state.objects; obj != NULL; obj = obj->next) {
        obj->generation = GEN_OLD;
    }
    remembered_set_clear();

    gc_state.next_gc = gc_state.bytes_allocated * 2;
    if (gc_state.next_gc < 1024u * 1024u) {
        gc_state.next_gc = 1024u * 1024u;
    }

    size_t freed = 0;
    if (before > gc_state.bytes_allocated) {
        freed = before - gc_state.bytes_allocated;
    }
    gc_state.stats.bytes_freed += freed;
    gc_state.stats.collections++;
    gc_state.stats.bytes_allocated = gc_state.bytes_allocated;

    if (!nursery_has_live_object()) {
        gc_reset_nursery();
    }

    gc_state.collecting = false;
    gc_record_pause_since(pause_start);
}

void gc_init(void) {
    gc_state.objects = NULL;
    gc_state.bytes_allocated = 0;
    gc_state.next_gc = 1024u * 1024u;
    gc_state.gc_disabled = false;
    gc_state.collecting = false;
    gc_state.use_shadow_stack = false;
    gc_state.nursery_top = gc_state.nursery_buf;
    gc_state.nursery_end = gc_state.nursery_buf + sizeof(gc_state.nursery_buf);
    gc_state.nursery_live_bytes = 0;
    gc_state.remembered_count = 0;
    gc_state.remembered_overflow_count = 0;
    memset(&gc_state.stats, 0, sizeof(gc_state.stats));
    gc_incremental_abort();
}

GCStats gc_get_stats(void) {
    gc_state.stats.bytes_allocated = gc_state.bytes_allocated;
    return gc_state.stats;
}

void gc_set_disabled(bool disabled) {
    gc_state.gc_disabled = disabled;
}

bool gc_is_disabled(void) {
    return gc_state.gc_disabled;
}

void gc_set_next_threshold(size_t bytes) {
    gc_state.next_gc = bytes;
}

uint64_t gc_debug_remembered_overflow_count(void) {
    return gc_state.remembered_overflow_count;
}

static bool gc_debug_enabled(void) {
    static int cached = -1;
    if (cached >= 0) {
        return cached != 0;
    }
    const char* env = getenv("KODA_GC_DEBUG");
    if (env == NULL || env[0] == '\0' || env[0] == '0') {
        cached = 0;
    } else {
        cached = 1;
    }
    return cached != 0;
}

static void gc_debug_validate_objects(void) {
    for (Obj* obj = gc_state.objects; obj != NULL; obj = obj->next) {
        if (obj->type < OBJ_STRING || obj->type > OBJ_CELL) {
            fprintf(stderr, "koda gc debug: invalid object header type=%d\n", (int)obj->type);
            abort();
        }
        if (obj->generation > GEN_DEAD) {
            fprintf(stderr, "koda gc debug: invalid generation=%u\n", (unsigned)obj->generation);
            abort();
        }
    }
}
