#include "koda.h"
#include <math.h>
#include <stdlib.h>
#include <time.h>
#ifdef _WIN32
#define WIN32_LEAN_AND_MEAN
#define NOGDI
#define NOUSER
#include <windows.h>
#else
#include <unistd.h>
#endif

#define STACK_MAX 1024
#define GRAY_STACK_MAX 1024
#define YOUNG_GEN_SIZE (1 * 1024 * 1024) // 1MB for demo
#define PROMOTION_AGE 3

typedef struct KodaTask {
    KodaValue closure;
    double target_time;
    double interval; // 0 for setTimeout, >0 for setInterval
    struct KodaTask* next;
} KodaTask;

typedef struct {
    KodaObj* objects;
    KodaValue stack[STACK_MAX];
    KodaValue* sp;
    
    KodaObj* gray_stack[GRAY_STACK_MAX];
    int gray_count;
    
    size_t bytes_allocated;
    size_t next_gc;
    
    double program_start;
    double last_frame;
    
    KodaTask* tasks;
} KodaHeap;

static double koda_get_time_seconds() {
#ifdef _WIN32
    LARGE_INTEGER freq, counter;
    QueryPerformanceFrequency(&freq);
    QueryPerformanceCounter(&counter);
    return (double)counter.QuadPart / freq.QuadPart;
#else
    struct timespec ts;
    clock_gettime(CLOCK_MONOTONIC, &ts);
    return (double)ts.tv_sec + (double)ts.tv_nsec / 1e9;
#endif
}

static KodaHeap heap;

void koda_init() {
    heap.objects = NULL;
    heap.sp = heap.stack;
    heap.bytes_allocated = 0;
    heap.next_gc = 10 * 1024 * 1024; // 10MB initial
    heap.gray_count = 0;
    
    heap.program_start = koda_get_time_seconds();
    heap.last_frame = heap.program_start;
    heap.tasks = NULL;
    
    srand((unsigned int)time(NULL));
}

void koda_runtime_init(void) {
    koda_init();
}

void koda_runtime_shutdown(void) {
    koda_shutdown();
}

void koda_push(KodaValue value) {
    if (heap.sp >= heap.stack + STACK_MAX) {
        fprintf(stderr, "Koda: stack overflow\n");
        exit(1);
    }
    *heap.sp = value;
    heap.sp++;
}

KodaValue koda_pop() {
    heap.sp--;
    return *heap.sp;
}

// Internal marking
static void gray_object(KodaObj* obj) {
    if (obj == NULL || obj->mark) return;
    obj->mark = 1;
    
    if (heap.gray_count >= GRAY_STACK_MAX) {
        // Fallback to recursive or just exit if stack is full
        // Production would reallocate gray stack
        return;
    }
    heap.gray_stack[heap.gray_count++] = obj;
}

static void blacken_object(KodaObj* obj) {
    switch (obj->type) {
        case OBJ_ARRAY:
        case OBJ_SET:
        case OBJ_TUPLE: {
            KodaArray* arr = (KodaArray*)obj;
            for (int i = 0; i < arr->count; i++) {
                if (IS_OBJ(arr->elements[i])) gray_object(AS_OBJ(arr->elements[i]));
            }
            break;
        }
        case OBJ_OBJECT:
        case OBJ_MAP: {
            KodaObject* ko = (KodaObject*)obj;
            for (int i = 0; i < ko->count; i++) {
                if (IS_OBJ(ko->entries[i].key)) gray_object(AS_OBJ(ko->entries[i].key));
                if (IS_OBJ(ko->entries[i].value)) gray_object(AS_OBJ(ko->entries[i].value));
            }
            break;
        }
        case OBJ_CLOSURE: {
            KodaClosure* closure = (KodaClosure*)obj;
            for (int i = 0; i < closure->upvalueCount; i++) {
                if (closure->upvalues[i]) gray_object((KodaObj*)closure->upvalues[i]);
            }
            break;
        }
        case OBJ_UPVALUE: {
            KodaUpvalue* upvalue = (KodaUpvalue*)obj;
            if (IS_OBJ(upvalue->closed)) gray_object(AS_OBJ(upvalue->closed));
            break;
        }
        case OBJ_BOUND_METHOD: {
            KodaBoundMethod* bound = (KodaBoundMethod*)obj;
            if (IS_OBJ(bound->receiver)) gray_object(AS_OBJ(bound->receiver));
            if (IS_OBJ(bound->method)) gray_object(AS_OBJ(bound->method));
            break;
        }
        default: break;
    }
}

void koda_collect() {
    // Mark roots
    for (KodaValue* v = heap.stack; v < heap.sp; v++) {
        if (IS_OBJ(*v)) gray_object(AS_OBJ(*v));
    }
    
    // Mark tasks
    KodaTask* task = heap.tasks;
    while (task) {
        if (IS_OBJ(task->closure)) gray_object(AS_OBJ(task->closure));
        task = task->next;
    }
    
    // Trace references
    while (heap.gray_count > 0) {
        KodaObj* obj = heap.gray_stack[--heap.gray_count];
        blacken_object(obj);
    }
    
    // Sweep
    KodaObj** prev = &heap.objects;
    while (*prev) {
        KodaObj* obj = *prev;
        if (!obj->mark) {
            *prev = obj->next;
            // Free nested memory
            if (obj->type == OBJ_ARRAY) free(((KodaArray*)obj)->elements);
            if (obj->type == OBJ_OBJECT) free(((KodaObject*)obj)->entries);
            if (obj->type == OBJ_CLOSURE) free(((KodaClosure*)obj)->upvalues);
            free(obj);
        } else {
            obj->mark = 0; // Reset for next GC
            prev = &obj->next;
        }
    }
    
    heap.next_gc = heap.bytes_allocated * 2;
}

KodaObj* koda_alloc(KodaObjType type, size_t size) {
    if (heap.bytes_allocated + size > heap.next_gc) {
        koda_collect();
    }
    
    KodaObj* obj = (KodaObj*)malloc(size);
    obj->type = type;
    obj->mark = 0;
    obj->next = heap.objects;
    heap.objects = obj;
    
    heap.bytes_allocated += size;
    return obj;
}

KodaArray* koda_new_array() {
    KodaArray* arr = (KodaArray*)koda_alloc(OBJ_ARRAY, sizeof(KodaArray));
    arr->count = 0;
    arr->capacity = 0;
    arr->elements = NULL;
    return arr;
}

KodaObject* koda_alloc_object() {
    KodaObject* obj = (KodaObject*)koda_alloc(OBJ_OBJECT, sizeof(KodaObject));
    obj->count = 0;
    obj->capacity = 0;
    obj->entries = NULL;
    return obj;
}

KodaUpvalue* koda_new_upvalue(KodaValue* slot) {
    KodaUpvalue* upvalue = (KodaUpvalue*)koda_alloc(OBJ_UPVALUE, sizeof(KodaUpvalue));
    if (slot == NULL) {
        upvalue->location = &upvalue->closed;
    } else {
        upvalue->location = slot;
    }
    upvalue->closed = NULL_VAL;
    upvalue->next = NULL;
    return upvalue;
}

KodaClosure* koda_new_closure(KodaFn fn, int upvalueCount) {
    size_t size = sizeof(KodaClosure) + sizeof(KodaUpvalue*) * upvalueCount;
    KodaClosure* closure = (KodaClosure*)koda_alloc(OBJ_CLOSURE, size);
    closure->fn = fn;
    closure->upvalueCount = upvalueCount;
    closure->upvalues = (KodaUpvalue**)((char*)closure + sizeof(KodaClosure));
    for (int i = 0; i < upvalueCount; i++) closure->upvalues[i] = NULL;
    return closure;
}

void koda_upvalue_set(KodaUpvalue* up, KodaValue val) {
    *up->location = val;
}

KodaValue koda_upvalue_get(KodaUpvalue* up) {
    return *up->location;
}

void koda_closure_set_upvalue(KodaValue closure, int index, KodaUpvalue* u) {
    AS_CLOSURE(closure)->upvalues[index] = u;
}

KodaValue koda_obj_as_value(void* obj) {
    return OBJ_VAL((KodaObj*)obj);
}

void koda_assert(KodaValue cond, KodaValue msg) {
    if (koda_is_truthy(cond)) return;
    fprintf(stderr, "assertion failed");
    if (!IS_NULL(msg)) {
        fprintf(stderr, ": ");
        koda_print_no_newline(msg);
    }
    fprintf(stderr, "\n");
    exit(1);
}
void koda_abort_arg_error(int kind, int a, int b) {
    switch (kind) {
    case 0:
        fprintf(stderr, "expected at least %d arguments but got %d\n", a, b);
        break;
    case 1:
        fprintf(stderr, "expected at most %d arguments but got %d\n", a, b);
        break;
    case 2:
        fprintf(stderr, "expected %d arguments but got %d\n", a, b);
        break;
    case 3:
        fprintf(stderr, "missing default for parameter %d\n", a);
        break;
    default:
        fprintf(stderr, "internal: bad arg error kind %d\n", kind);
        break;
    }
    exit(1);
}

KodaValue koda_argv_slice_to_array(KodaValue* argv, int start, int count) {
    KodaArray* arr = koda_new_array();
    KodaValue wrapper = koda_obj_as_value(arr);
    if (count <= 0 || argv == NULL) {
        return wrapper;
    }
    for (int i = 0; i < count; i++) {
        koda_array_push(wrapper, argv[start + i]);
    }
    return wrapper;
}

KodaValue koda_new_bound_method(KodaValue receiver, KodaValue method) {
    KodaBoundMethod* bound = (KodaBoundMethod*)koda_alloc(OBJ_BOUND_METHOD, sizeof(KodaBoundMethod));
    bound->receiver = receiver;
    bound->method = method;
    return OBJ_VAL(bound);
}

KodaValue koda_new_native(KodaNativeFn fn) {
    KodaNative* native = (KodaNative*)koda_alloc(OBJ_NATIVE, sizeof(KodaNative));
    native->fn = fn;
    return OBJ_VAL(native);
}

void koda_array_push(KodaValue array, KodaValue value) {
    if (!IS_OBJ(array)) return;
    KodaObjType t = AS_OBJ(array)->type;
    if (t != OBJ_ARRAY && t != OBJ_SET) return;
    KodaArray* arr = (KodaArray*)AS_OBJ(array);
    if (arr->count >= arr->capacity) {
        arr->capacity = arr->capacity < 8 ? 8 : arr->capacity * 2;
        arr->elements = (KodaValue*)realloc(arr->elements, sizeof(KodaValue) * arr->capacity);
    }
    arr->elements[arr->count++] = value;
    KODA_WRITE_BARRIER(arr, value);
}

// --- Native Methods ---

KodaValue koda_method_string_upper(int argCount, KodaValue* args) {
    KodaValue receiver = args[0]; // Bound method passes receiver as first arg
    if (!IS_OBJ(receiver) || AS_OBJ(receiver)->type != OBJ_STRING) return NULL_VAL;
    KodaString* s = (KodaString*)AS_OBJ(receiver);
    KodaString* res = koda_copy_string(s->chars, s->length);
    for (int i = 0; i < res->length; i++) {
        if (res->chars[i] >= 'a' && res->chars[i] <= 'z') res->chars[i] -= 32;
    }
    return OBJ_VAL(res);
}

KodaValue koda_method_array_push(int argCount, KodaValue* args) {
    KodaValue receiver = args[0];
    if (!IS_OBJ(receiver) || AS_OBJ(receiver)->type != OBJ_ARRAY) return NULL_VAL;
    KodaValue val = args[1];
    koda_array_push(receiver, val);
    return receiver;
}

KodaValue koda_method_array_pop(int argCount, KodaValue* args) {
    KodaValue receiver = args[0];
    if (!IS_OBJ(receiver) || AS_OBJ(receiver)->type != OBJ_ARRAY) return NULL_VAL;
    KodaArray* arr = (KodaArray*)AS_OBJ(receiver);
    if (arr->count <= 0) return NULL_VAL;
    arr->count--;
    return arr->elements[arr->count];
}

KodaValue koda_method_array_length(int argCount, KodaValue* args) {
    KodaValue receiver = args[0];
    if (!IS_OBJ(receiver) || AS_OBJ(receiver)->type != OBJ_ARRAY) return NULL_VAL;
    KodaArray* arr = (KodaArray*)AS_OBJ(receiver);
    return NUMBER_VAL((double)arr->count);
}

void koda_object_set(KodaValue obj, KodaValue key, KodaValue value) {
    if (!IS_OBJ(obj)) return;
    KodaObjType t = AS_OBJ(obj)->type;
    if (t != OBJ_OBJECT && t != OBJ_MAP) return;
    KodaObject* o = (KodaObject*)AS_OBJ(obj);
    
    // Simple linear search for now, can be optimized later
    for (int i = 0; i < o->count; i++) {
        if (AS_BOOL(koda_eq(o->entries[i].key, key))) {
            o->entries[i].value = value;
            return;
        }
    }
    
    if (o->count == o->capacity) {
        o->capacity = o->capacity == 0 ? 8 : o->capacity * 2;
        o->entries = (KodaEntry*)realloc(o->entries, sizeof(KodaEntry) * o->capacity);
    }
    
    o->entries[o->count].key = key;
    o->entries[o->count].value = value;
    o->count++;
}

KodaValue koda_object_delete(KodaValue obj, KodaValue key) {
    if (!IS_OBJ(obj)) return BOOL_VAL(false);
    KodaObjType t = AS_OBJ(obj)->type;
    if (t != OBJ_OBJECT && t != OBJ_MAP) return BOOL_VAL(false);
    KodaObject* o = (KodaObject*)AS_OBJ(obj);
    for (int i = 0; i < o->count; i++) {
        if (AS_BOOL(koda_eq(o->entries[i].key, key))) {
            o->entries[i] = o->entries[o->count - 1];
            o->count--;
            return BOOL_VAL(true);
        }
    }
    return BOOL_VAL(false);
}

KodaValue koda_map_method_set(int argCount, KodaValue* args) {
    if (argCount < 3) return NULL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_MAP) return NULL_VAL;
    koda_object_set(args[0], args[1], args[2]);
    return NULL_VAL;
}

KodaValue koda_map_method_get(int argCount, KodaValue* args) {
    if (argCount < 2) return NULL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_MAP) return NULL_VAL;
    KodaObject* m = (KodaObject*)AS_OBJ(args[0]);
    KodaValue key = args[1];
    for (int i = 0; i < m->count; i++) {
        if (AS_BOOL(koda_eq(m->entries[i].key, key))) return m->entries[i].value;
    }
    return NULL_VAL;
}

KodaValue koda_map_method_has(int argCount, KodaValue* args) {
    if (argCount < 2) return BOOL_VAL(false);
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_MAP) return BOOL_VAL(false);
    KodaObject* m = (KodaObject*)AS_OBJ(args[0]);
    KodaValue key = args[1];
    for (int i = 0; i < m->count; i++) {
        if (AS_BOOL(koda_eq(m->entries[i].key, key))) return BOOL_VAL(true);
    }
    return BOOL_VAL(false);
}

KodaValue koda_map_method_delete(int argCount, KodaValue* args) {
    if (argCount < 2) return BOOL_VAL(false);
    return koda_object_delete(args[0], args[1]);
}

KodaValue koda_map_method_size(int argCount, KodaValue* args) {
    if (argCount < 1) return NUMBER_VAL(0);
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_MAP) return NUMBER_VAL(0);
    KodaObject* m = (KodaObject*)AS_OBJ(args[0]);
    return NUMBER_VAL((double)m->count);
}

static int koda_set_index_of(KodaArray* s, KodaValue v) {
    for (int i = 0; i < s->count; i++) {
        if (AS_BOOL(koda_eq(s->elements[i], v))) return i;
    }
    return -1;
}

KodaValue koda_set_method_add(int argCount, KodaValue* args) {
    if (argCount < 2) return NULL_VAL;
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_SET) return NULL_VAL;
    KodaArray* s = (KodaArray*)AS_OBJ(args[0]);
    if (koda_set_index_of(s, args[1]) >= 0) return NULL_VAL;
    koda_array_push(args[0], args[1]);
    return NULL_VAL;
}

KodaValue koda_set_method_has(int argCount, KodaValue* args) {
    if (argCount < 2) return BOOL_VAL(false);
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_SET) return BOOL_VAL(false);
    return BOOL_VAL(koda_set_index_of((KodaArray*)AS_OBJ(args[0]), args[1]) >= 0);
}

KodaValue koda_set_method_remove(int argCount, KodaValue* args) {
    if (argCount < 2) return BOOL_VAL(false);
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_SET) return BOOL_VAL(false);
    KodaArray* s = (KodaArray*)AS_OBJ(args[0]);
    int i = koda_set_index_of(s, args[1]);
    if (i < 0) return BOOL_VAL(false);
    s->elements[i] = s->elements[s->count - 1];
    s->count--;
    return BOOL_VAL(true);
}

KodaValue koda_set_method_size(int argCount, KodaValue* args) {
    if (argCount < 1) return NUMBER_VAL(0);
    if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_SET) return NUMBER_VAL(0);
    return NUMBER_VAL((double)((KodaArray*)AS_OBJ(args[0]))->count);
}

KodaValue koda_map_new(int argCount, KodaValue* args) {
    (void)args;
    if (argCount != 0) return NULL_VAL;
    KodaObject* m = (KodaObject*)koda_alloc(OBJ_MAP, sizeof(KodaObject));
    m->count = 0;
    m->capacity = 0;
    m->entries = NULL;
    return OBJ_VAL((KodaObj*)m);
}

KodaValue koda_tuple_new(int argCount, KodaValue* args) {
    if (argCount < 2) return NULL_VAL;
    KodaArray* arr = (KodaArray*)koda_alloc(OBJ_TUPLE, sizeof(KodaArray));
    arr->capacity = argCount;
    arr->count = argCount;
    arr->elements = (KodaValue*)malloc(sizeof(KodaValue) * (size_t)argCount);
    if (!arr->elements) return NULL_VAL;
    memcpy(arr->elements, args, sizeof(KodaValue) * (size_t)argCount);
    return OBJ_VAL((KodaObj*)arr);
}

KodaValue koda_input(int argCount, KodaValue* args) {
    if (argCount >= 1 && IS_OBJ(args[0]) && AS_OBJ(args[0])->type == OBJ_STRING) {
        KodaString* p = (KodaString*)AS_OBJ(args[0]);
        printf("%.*s", p->length, p->chars);
        fflush(stdout);
    }
    char buf[4096];
    if (fgets(buf, sizeof buf, stdin) == NULL) {
        return OBJ_VAL(koda_copy_string("", 0));
    }
    size_t n = strlen(buf);
    while (n > 0 && (buf[n - 1] == '\n' || buf[n - 1] == '\r')) {
        buf[--n] = '\0';
    }
    return OBJ_VAL(koda_copy_string(buf, (int)n));
}

KodaValue koda_set_new(int argCount, KodaValue* args) {
    if (argCount > 1) return NULL_VAL;
    KodaArray* s = (KodaArray*)koda_alloc(OBJ_SET, sizeof(KodaArray));
    s->count = 0;
    s->capacity = 0;
    s->elements = NULL;
    KodaValue sv = OBJ_VAL((KodaObj*)s);
    if (argCount >= 1) {
        if (!IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_ARRAY) return NULL_VAL;
        KodaArray* a = (KodaArray*)AS_OBJ(args[0]);
        for (int i = 0; i < a->count; i++) {
            KodaValue v = a->elements[i];
            if (koda_set_index_of(s, v) < 0) koda_array_push(sv, v);
        }
    }
    return sv;
}

KodaValue koda_get_index(KodaValue obj, KodaValue index) {
    if (!IS_OBJ(obj)) return NULL_VAL;
    KodaObj* o = AS_OBJ(obj);
    KodaValue recv = OBJ_VAL(o);

    if (o->type == OBJ_TUPLE) {
        if (IS_NUMBER(index)) {
            KodaArray* arr = (KodaArray*)o;
            int idx = (int)AS_NUMBER(index);
            if (idx < 0) idx = arr->count + idx;
            if (idx >= 0 && idx < arr->count) return arr->elements[idx];
            return NULL_VAL;
        }
        return NULL_VAL;
    }

    if (o->type == OBJ_ARRAY) {
        if (IS_NUMBER(index)) {
            KodaArray* arr = (KodaArray*)o;
            int idx = (int)AS_NUMBER(index);
            if (idx < 0) idx = arr->count + idx;
            if (idx >= 0 && idx < arr->count) return arr->elements[idx];
            return NULL_VAL;
        }
    } else if (o->type == OBJ_STRING) {
        if (IS_NUMBER(index)) {
            KodaString* s = (KodaString*)o;
            int idx = (int)AS_NUMBER(index);
            if (idx < 0) idx = s->length + idx;
            if (idx >= 0 && idx < s->length) {
                return OBJ_VAL(koda_copy_string(s->chars + idx, 1));
            }
            return NULL_VAL;
        }
    } else if (o->type == OBJ_OBJECT || o->type == OBJ_MAP) {
        KodaObject* mobj = (KodaObject*)o;
        for (int i = 0; i < mobj->count; i++) {
            if (AS_BOOL(koda_eq(mobj->entries[i].key, index))) {
                KodaValue val = mobj->entries[i].value;
                if (IS_OBJ(val) && (AS_OBJ(val)->type == OBJ_CLOSURE || AS_OBJ(val)->type == OBJ_NATIVE)) {
                    return koda_new_bound_method(obj, val);
                }
                return val;
            }
        }
        if (o->type == OBJ_MAP && IS_OBJ(index) && AS_OBJ(index)->type == OBJ_STRING) {
            KodaString* s = (KodaString*)AS_OBJ(index);
            if (strcmp(s->chars, "set") == 0) return koda_new_bound_method(recv, koda_new_native(koda_map_method_set));
            if (strcmp(s->chars, "get") == 0) return koda_new_bound_method(recv, koda_new_native(koda_map_method_get));
            if (strcmp(s->chars, "has") == 0) return koda_new_bound_method(recv, koda_new_native(koda_map_method_has));
            if (strcmp(s->chars, "delete") == 0) return koda_new_bound_method(recv, koda_new_native(koda_map_method_delete));
            if (strcmp(s->chars, "size") == 0) return koda_new_bound_method(recv, koda_new_native(koda_map_method_size));
        }
    } else if (o->type == OBJ_SET && IS_OBJ(index) && AS_OBJ(index)->type == OBJ_STRING) {
        KodaString* s = (KodaString*)AS_OBJ(index);
        if (strcmp(s->chars, "add") == 0) return koda_new_bound_method(recv, koda_new_native(koda_set_method_add));
        if (strcmp(s->chars, "has") == 0) return koda_new_bound_method(recv, koda_new_native(koda_set_method_has));
        if (strcmp(s->chars, "remove") == 0) return koda_new_bound_method(recv, koda_new_native(koda_set_method_remove));
        if (strcmp(s->chars, "size") == 0) return koda_new_bound_method(recv, koda_new_native(koda_set_method_size));
    }

    if (o->type == OBJ_STRING) {
        if (IS_OBJ(index) && AS_OBJ(index)->type == OBJ_STRING) {
            KodaString* strk = (KodaString*)AS_OBJ(index);
            if (strcmp(strk->chars, "upper") == 0) {
                return koda_new_bound_method(recv, koda_new_native(koda_method_string_upper));
            }
        }
    } else if (o->type == OBJ_ARRAY) {
        if (IS_OBJ(index) && AS_OBJ(index)->type == OBJ_STRING) {
            KodaString* strk = (KodaString*)AS_OBJ(index);
            if (strcmp(strk->chars, "length") == 0) {
                return koda_new_bound_method(recv, koda_new_native(koda_method_array_length));
            }
            if (strcmp(strk->chars, "push") == 0) {
                return koda_new_bound_method(recv, koda_new_native(koda_method_array_push));
            }
            if (strcmp(strk->chars, "pop") == 0) {
                return koda_new_bound_method(recv, koda_new_native(koda_method_array_pop));
            }
        }
    }

    return NULL_VAL;
}

KodaValue koda_call(KodaValue callee, int argCount, KodaValue* args) {
    if (!IS_OBJ(callee)) {
        fprintf(stderr, "Koda: callee is not an object\n");
        return NULL_VAL;
    }
    KodaObj* obj = AS_OBJ(callee);
    if (obj->type == OBJ_CLOSURE) {
        KodaClosure* closure = (KodaClosure*)obj;
        return closure->fn(NULL_VAL, argCount, args, closure->upvalues);
    } else if (obj->type == OBJ_NATIVE) {
        return ((KodaNative*)obj)->fn(argCount, args);
    } else if (obj->type == OBJ_BOUND_METHOD) {
        KodaBoundMethod* bound = (KodaBoundMethod*)obj;
        if (IS_OBJ(bound->method) && AS_OBJ(bound->method)->type == OBJ_CLOSURE) {
            KodaClosure* closure = (KodaClosure*)AS_OBJ(bound->method);
            return closure->fn(bound->receiver, argCount, args, closure->upvalues);
        } else if (IS_OBJ(bound->method) && AS_OBJ(bound->method)->type == OBJ_NATIVE) {
            KodaNative* native = (KodaNative*)AS_OBJ(bound->method);
            // Native functions still expect receiver as first arg
            KodaValue* newArgs = (KodaValue*)malloc(sizeof(KodaValue) * (argCount + 1));
            newArgs[0] = bound->receiver;
            memcpy(newArgs + 1, args, sizeof(KodaValue) * argCount);
            KodaValue res = native->fn(argCount + 1, newArgs);
            free(newArgs);
            return res;
        }
    }
    
    fprintf(stderr, "Koda: callee type %d is not callable\n", obj->type);
    return NULL_VAL;
}

KodaString* koda_copy_string(const char* chars, int length) {
    KodaString* str = (KodaString*)koda_alloc(OBJ_STRING, sizeof(KodaString) + length + 1);
    str->length = length;
    memcpy(str->chars, chars, length);
    str->chars[length] = '\0';
    return str;
}

KodaValue koda_allocate_string(int length, const char* chars) {
    return OBJ_VAL(koda_copy_string(chars, length));
}

void koda_write_barrier(KodaObj* obj, KodaValue value) {
    (void)obj; (void)value;
    // No-op for Mark-Sweep GC
}

// --- Operations ---

KodaValue koda_add(KodaValue a, KodaValue b) {
    if (IS_NUMBER(a) && IS_NUMBER(b)) {
        return NUMBER_VAL(AS_NUMBER(a) + AS_NUMBER(b));
    }
    
    if (IS_OBJ(a) && AS_OBJ(a)->type == OBJ_STRING) {
        KodaString* sa = (KodaString*)AS_OBJ(a);
        char buf[64];
        const char* sb;
        int blen;
        
        if (IS_NUMBER(b)) {
            blen = sprintf(buf, "%g", AS_NUMBER(b));
            sb = buf;
        } else if (IS_OBJ(b) && AS_OBJ(b)->type == OBJ_STRING) {
            KodaString* s_obj = (KodaString*)AS_OBJ(b);
            sb = s_obj->chars;
            blen = s_obj->length;
        } else {
            return NULL_VAL;
        }
        
        int total_len = sa->length + blen;
        KodaString* res = (KodaString*)koda_alloc(OBJ_STRING, sizeof(KodaString) + total_len + 1);
        res->length = total_len;
        memcpy(res->chars, sa->chars, sa->length);
        memcpy(res->chars + sa->length, sb, blen);
        res->chars[total_len] = '\0';
        return OBJ_VAL(res);
    }
    
    return NULL_VAL;
}

KodaValue koda_to_string_val(KodaValue v) {
    char buf[256];
    int n;
    if (IS_NUMBER(v)) {
        n = snprintf(buf, sizeof(buf), "%g", AS_NUMBER(v));
        if (n < 0) n = 0;
        return OBJ_VAL(koda_copy_string(buf, n));
    }
    if (IS_BOOL(v)) {
        if (AS_BOOL(v)) {
            return OBJ_VAL(koda_copy_string("true", 4));
        }
        return OBJ_VAL(koda_copy_string("false", 5));
    }
    if (IS_NULL(v)) {
        return OBJ_VAL(koda_copy_string("null", 4));
    }
    if (IS_OBJ(v) && AS_OBJ(v)->type == OBJ_STRING) {
        return v;
    }
    n = snprintf(buf, sizeof(buf), "(value)");
    return OBJ_VAL(koda_copy_string(buf, n));
}

KodaValue koda_parse_number(KodaValue v) {
    if (IS_NUMBER(v)) {
        return v;
    }
    if (IS_OBJ(v) && AS_OBJ(v)->type == OBJ_STRING) {
        KodaString* s = (KodaString*)AS_OBJ(v);
        char* end = NULL;
        double d = strtod(s->chars, &end);
        if (end == s->chars) {
            return NULL_VAL;
        }
        while (*end == ' ' || *end == '\t' || *end == '\n' || *end == '\r') {
            end++;
        }
        if (*end != '\0') {
            return NULL_VAL;
        }
        return NUMBER_VAL(d);
    }
    return NULL_VAL;
}

KodaValue koda_range(KodaValue from, KodaValue to) {
    if (!IS_NUMBER(from) || !IS_NUMBER(to)) {
        return NULL_VAL;
    }
    long long lo = (long long)trunc(AS_NUMBER(from));
    long long hi = (long long)trunc(AS_NUMBER(to));
    KodaValue arr = OBJ_VAL(koda_new_array());
    if (lo <= hi) {
        for (long long i = lo; i <= hi; i++) {
            koda_array_push(arr, NUMBER_VAL((double)i));
        }
    } else {
        for (long long i = lo; i >= hi; i--) {
            koda_array_push(arr, NUMBER_VAL((double)i));
        }
    }
    return arr;
}

KodaValue koda_is_truthy(KodaValue v) {
    if (IS_NULL(v)) return BOOL_VAL(false);
    if (IS_BOOL(v)) return v;
    if (IS_NUMBER(v)) return BOOL_VAL(AS_NUMBER(v) != 0);
    return BOOL_VAL(true);
}

void koda_print_no_newline(KodaValue val) {
    if (IS_NUMBER(val)) {
        printf("%g", AS_NUMBER(val));
    } else if (IS_BOOL(val)) {
        printf("%s", AS_BOOL(val) ? "true" : "false");
    } else if (IS_NULL(val)) {
        printf("null");
    } else if (IS_OBJ(val)) {
        KodaObj* obj = AS_OBJ(val);
        if (obj->type == OBJ_STRING) {
            printf("%s", ((KodaString*)obj)->chars);
        } else if (obj->type == OBJ_ARRAY) {
            KodaArray* arr = (KodaArray*)obj;
            printf("[");
            for (int i = 0; i < arr->count; i++) {
                koda_print_no_newline(arr->elements[i]);
                if (i < arr->count - 1) printf(", ");
            }
            printf("]");
        } else {
            printf("<obj %p type %d>", obj, obj->type);
        }
    }
}

void koda_print_space(void) { printf(" "); }

void koda_print_newline(void) { printf("\n"); }

KodaValue koda_print(KodaValue val) {
    koda_print_no_newline(val);
    koda_print_newline();
    return NULL_VAL;
}

// Rest of arithmetic...
KodaValue koda_sub(KodaValue a, KodaValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL(AS_NUMBER(a) - AS_NUMBER(b)) : NULL_VAL; }
KodaValue koda_mul(KodaValue a, KodaValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL(AS_NUMBER(a) * AS_NUMBER(b)) : NULL_VAL; }
KodaValue koda_div(KodaValue a, KodaValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL(AS_NUMBER(a) / AS_NUMBER(b)) : NULL_VAL; }
KodaValue koda_gt(KodaValue a, KodaValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? BOOL_VAL(AS_NUMBER(a) > AS_NUMBER(b)) : BOOL_VAL(false); }
KodaValue koda_lt(KodaValue a, KodaValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? BOOL_VAL(AS_NUMBER(a) < AS_NUMBER(b)) : BOOL_VAL(false); }
KodaValue koda_eq(KodaValue a, KodaValue b) {
    if (a == b) return BOOL_VAL(true);
    if (IS_OBJ(a) && IS_OBJ(b)) {
        KodaObj* oa = AS_OBJ(a);
        KodaObj* ob = AS_OBJ(b);
        if (oa->type == OBJ_STRING && ob->type == OBJ_STRING) {
            KodaString* sa = (KodaString*)oa;
            KodaString* sb = (KodaString*)ob;
            return BOOL_VAL(sa->length == sb->length && memcmp(sa->chars, sb->chars, sa->length) == 0);
        }
    }
    return BOOL_VAL(false);
}

KodaValue koda_neq(KodaValue a, KodaValue b) { return BOOL_VAL(!AS_BOOL(koda_eq(a, b))); }
KodaValue koda_le(KodaValue a, KodaValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? BOOL_VAL(AS_NUMBER(a) <= AS_NUMBER(b)) : BOOL_VAL(false); }
KodaValue koda_ge(KodaValue a, KodaValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? BOOL_VAL(AS_NUMBER(a) >= AS_NUMBER(b)) : BOOL_VAL(false); }
KodaValue koda_mod(KodaValue a, KodaValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL(fmod(AS_NUMBER(a), AS_NUMBER(b))) : NULL_VAL; }
KodaValue koda_pow(KodaValue a, KodaValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL(pow(AS_NUMBER(a), AS_NUMBER(b))) : NULL_VAL; }
KodaValue koda_negate(KodaValue a) { return IS_NUMBER(a) ? NUMBER_VAL(-AS_NUMBER(a)) : NULL_VAL; }

KodaValue koda_bit_and(KodaValue a, KodaValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL((double)((int64_t)AS_NUMBER(a) & (int64_t)AS_NUMBER(b))) : NULL_VAL; }
KodaValue koda_bit_or(KodaValue a, KodaValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL((double)((int64_t)AS_NUMBER(a) | (int64_t)AS_NUMBER(b))) : NULL_VAL; }
KodaValue koda_bit_xor(KodaValue a, KodaValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL((double)((int64_t)AS_NUMBER(a) ^ (int64_t)AS_NUMBER(b))) : NULL_VAL; }
KodaValue koda_bit_not(KodaValue a) { return IS_NUMBER(a) ? NUMBER_VAL((double)(~(int64_t)AS_NUMBER(a))) : NULL_VAL; }
KodaValue koda_shl(KodaValue a, KodaValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL((double)((int64_t)AS_NUMBER(a) << (int64_t)AS_NUMBER(b))) : NULL_VAL; }
KodaValue koda_shr(KodaValue a, KodaValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL((double)((int64_t)AS_NUMBER(a) >> (int64_t)AS_NUMBER(b))) : NULL_VAL; }
KodaValue koda_ushr(KodaValue a, KodaValue b) { return IS_NUMBER(a) && IS_NUMBER(b) ? NUMBER_VAL((double)((uint64_t)(int64_t)AS_NUMBER(a) >> (int64_t)AS_NUMBER(b))) : NULL_VAL; }

void koda_set_index(KodaValue obj, KodaValue index, KodaValue value) {
    if (!IS_OBJ(obj)) return;
    KodaObj* o = AS_OBJ(obj);
    if (o->type == OBJ_ARRAY && IS_NUMBER(index)) {
        KodaArray* arr = (KodaArray*)o;
        int idx = (int)AS_NUMBER(index);
        if (idx < 0) idx = arr->count + idx;
        if (idx >= 0 && idx < arr->count) {
            arr->elements[idx] = value;
            KODA_WRITE_BARRIER(arr, value);
        }
    } else if (o->type == OBJ_OBJECT || o->type == OBJ_MAP) {
        koda_object_set(obj, index, value);
    }
}

int koda_for_in_len(KodaValue iterable) {
    if (!IS_OBJ(iterable)) return 0;
    KodaObj* o = AS_OBJ(iterable);
    if (o->type == OBJ_ARRAY) return ((KodaArray*)o)->count;
    if (o->type == OBJ_OBJECT) return ((KodaObject*)o)->count;
    if (o->type == OBJ_STRING) return ((KodaString*)o)->length;
    return 0;
}

KodaValue koda_for_in_get(KodaValue iterable, int index) {
    if (!IS_OBJ(iterable)) return NULL_VAL;
    KodaObj* o = AS_OBJ(iterable);
    if (o->type == OBJ_ARRAY) {
        KodaArray* arr = (KodaArray*)o;
        if (index >= 0 && index < arr->count) return NUMBER_VAL((double)index);
        return NULL_VAL;
    }
    if (o->type == OBJ_OBJECT) {
        KodaObject* obj = (KodaObject*)o;
        if (index >= 0 && index < obj->count) return obj->entries[index].key;
        return NULL_VAL;
    }
    if (o->type == OBJ_STRING) {
        KodaString* s = (KodaString*)o;
        if (index >= 0 && index < s->length) return NUMBER_VAL((double)index);
        return NULL_VAL;
    }
    return NULL_VAL;
}

KodaValue koda_slice(KodaValue obj, KodaValue start_val, KodaValue end_val) {
    if (!IS_OBJ(obj)) return NULL_VAL;
    KodaObj* o = AS_OBJ(obj);
    int length = 0;
    if (o->type == OBJ_ARRAY) length = ((KodaArray*)o)->count;
    else if (o->type == OBJ_STRING) length = ((KodaString*)o)->length;
    else return NULL_VAL;

    int start = 0;
    if (IS_NULL(start_val)) start = 0;
    else {
        start = (int)AS_NUMBER(start_val);
        if (start < 0) start = length + start;
    }

    int end = length;
    if (IS_NULL(end_val)) end = length;
    else {
        end = (int)AS_NUMBER(end_val);
        if (end < 0) end = length + end;
    }

    if (start < 0) start = 0;
    if (end > length) end = length;
    if (start > end) start = end;

    if (o->type == OBJ_ARRAY) {
        KodaArray* src = (KodaArray*)o;
        KodaArray* res = koda_new_array();
        int count = end - start;
        for (int i = 0; i < count; i++) {
            koda_array_push(OBJ_VAL(res), src->elements[start + i]);
        }
        return OBJ_VAL(res);
    } else {
        KodaString* src = (KodaString*)o;
        return OBJ_VAL(koda_copy_string(src->chars + start, end - start));
    }
}

void koda_shutdown() {
    /* Standalone Koda executables are process-lifetime runtimes.
       The OS reclaims heap pages on exit; doing a manual sweep here risks
       double-freeing objects already handled by runtime/host finalizers. */
    heap.objects = NULL;
    heap.sp = heap.stack;
}

KodaValue koda_len(KodaValue val) {
    if (!IS_OBJ(val)) return NUMBER_VAL(0);
    KodaObj* obj = AS_OBJ(val);
    switch (obj->type) {
        case OBJ_STRING: return NUMBER_VAL(((KodaString*)obj)->length);
        case OBJ_ARRAY: return NUMBER_VAL(((KodaArray*)obj)->count);
        case OBJ_OBJECT:
        case OBJ_MAP: return NUMBER_VAL(((KodaObject*)obj)->count);
        case OBJ_SET:
        case OBJ_TUPLE: return NUMBER_VAL(((KodaArray*)obj)->count);
        default: return NUMBER_VAL(0);
    }
}

KodaValue koda_type(KodaValue val) {
    const char* type_str = "unknown";
    if (IS_NUMBER(val)) type_str = "number";
    else if (IS_BOOL(val)) type_str = "bool";
    else if (IS_NULL(val)) type_str = "null";
    else if (IS_OBJ(val)) {
        switch (AS_OBJ(val)->type) {
            case OBJ_STRING: type_str = "string"; break;
            case OBJ_ARRAY: type_str = "array"; break;
            case OBJ_OBJECT: type_str = "object"; break;
            case OBJ_MAP: type_str = "map"; break;
            case OBJ_SET: type_str = "set"; break;
            case OBJ_TUPLE: type_str = "tuple"; break;
            case OBJ_CLOSURE: type_str = "function"; break;
            case OBJ_NATIVE: type_str = "function"; break;
            case OBJ_BOUND_METHOD: type_str = "function"; break;
            default: type_str = "object"; break;
        }
    }
    return OBJ_VAL(koda_copy_string(type_str, strlen(type_str)));
}

KodaValue koda_clock() {
    return NUMBER_VAL((double)clock() / CLOCKS_PER_SEC);
}

KodaValue koda_wall_time(void) {
    return NUMBER_VAL((double)time(NULL));
}

void koda_sleep_ms(KodaValue ms_val) {
    if (!IS_NUMBER(ms_val)) return;
    double ms = AS_NUMBER(ms_val);
    if (ms <= 0) return;
#ifdef _WIN32
    Sleep((DWORD)(ms + 0.5));
#else
    struct timespec ts;
    ts.tv_sec = (time_t)(ms / 1000.0);
    ts.tv_nsec = (long)fmod(ms, 1000.0) * 1000000L;
    if (ts.tv_nsec < 0) ts.tv_nsec = 0;
    nanosleep(&ts, NULL);
#endif
}

static int koda_rand_seeded;

KodaValue koda_random_unit(void) {
    if (!koda_rand_seeded) {
        srand((unsigned int)time(NULL));
        koda_rand_seeded = 1;
    }
    return NUMBER_VAL((double)rand() / (double)RAND_MAX);
}

KodaValue koda_io_read_file(int argCount, KodaValue* args) {
    if (argCount < 1 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NULL_VAL;
    const char* path = ((KodaString*)AS_OBJ(args[0]))->chars;
    
    FILE* file = fopen(path, "rb");
    if (file == NULL) return NULL_VAL;
    
    fseek(file, 0L, SEEK_END);
    size_t fileSize = ftell(file);
    rewind(file);
    
    char* buffer = (char*)malloc(fileSize + 1);
    if (buffer == NULL) {
        fclose(file);
        return NULL_VAL;
    }
    
    size_t bytesRead = fread(buffer, sizeof(char), fileSize, file);
    buffer[bytesRead] = '\0';
    fclose(file);
    
    KodaValue res = OBJ_VAL(koda_copy_string(buffer, bytesRead));
    free(buffer);
    return res;
}

KodaValue koda_io_write_file(int argCount, KodaValue* args) {
    if (argCount < 2 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NULL_VAL;
    if (!IS_OBJ(args[1]) || AS_OBJ(args[1])->type != OBJ_STRING) return NULL_VAL;
    
    const char* path = ((KodaString*)AS_OBJ(args[0]))->chars;
    KodaString* content = (KodaString*)AS_OBJ(args[1]);
    
    FILE* file = fopen(path, "wb");
    if (file == NULL) return BOOL_VAL(false);
    
    fwrite(content->chars, sizeof(char), content->length, file);
    fclose(file);
    return BOOL_VAL(true);
}

KodaValue koda_module_init_io() {
    static KodaValue module_exports = 0;
    if (module_exports != 0) return module_exports;
    module_exports = OBJ_VAL(koda_alloc_object());
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("read", 4)), koda_new_native(koda_io_read_file));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("write", 5)), koda_new_native(koda_io_write_file));
    return module_exports;
}

// --- JSON Module ---

typedef struct {
    char* buffer;
    int capacity;
    int length;
} KodaStringBuilder;

static void sb_init(KodaStringBuilder* sb) {
    sb->capacity = 128;
    sb->length = 0;
    sb->buffer = (char*)malloc(sb->capacity);
}

static void sb_append(KodaStringBuilder* sb, const char* str, int len) {
    if (sb->length + len >= sb->capacity) {
        while (sb->length + len >= sb->capacity) sb->capacity *= 2;
        sb->buffer = (char*)realloc(sb->buffer, sb->capacity);
    }
    memcpy(sb->buffer + sb->length, str, len);
    sb->length += len;
    sb->buffer[sb->length] = '\0';
}

static void koda_json_stringify_recursive(KodaStringBuilder* sb, KodaValue val) {
    if (IS_NUMBER(val)) {
        char buf[32];
        int len = sprintf(buf, "%g", AS_NUMBER(val));
        sb_append(sb, buf, len);
    } else if (IS_BOOL(val)) {
        if (AS_BOOL(val)) sb_append(sb, "true", 4);
        else sb_append(sb, "false", 5);
    } else if (IS_NULL(val)) {
        sb_append(sb, "null", 4);
    } else if (IS_OBJ(val)) {
        KodaObj* obj = AS_OBJ(val);
        if (obj->type == OBJ_STRING) {
            KodaString* s = (KodaString*)obj;
            sb_append(sb, "\"", 1);
            sb_append(sb, s->chars, s->length);
            sb_append(sb, "\"", 1);
        } else if (obj->type == OBJ_ARRAY) {
            KodaArray* arr = (KodaArray*)obj;
            sb_append(sb, "[", 1);
            for (int i = 0; i < arr->count; i++) {
                koda_json_stringify_recursive(sb, arr->elements[i]);
                if (i < arr->count - 1) sb_append(sb, ",", 1);
            }
            sb_append(sb, "]", 1);
        } else if (obj->type == OBJ_OBJECT) {
            KodaObject* o = (KodaObject*)obj;
            sb_append(sb, "{", 1);
            for (int i = 0; i < o->count; i++) {
                KodaString* key = (KodaString*)AS_OBJ(o->entries[i].key);
                sb_append(sb, "\"", 1);
                sb_append(sb, key->chars, key->length);
                sb_append(sb, "\":", 2);
                koda_json_stringify_recursive(sb, o->entries[i].value);
                if (i < o->count - 1) sb_append(sb, ",", 1);
            }
            sb_append(sb, "}", 1);
        }
    }
}

KodaValue koda_json_stringify(int argCount, KodaValue* args) {
    if (argCount < 1) return NULL_VAL;
    KodaStringBuilder sb;
    sb_init(&sb);
    koda_json_stringify_recursive(&sb, args[0]);
    KodaValue res = OBJ_VAL(koda_copy_string(sb.buffer, sb.length));
    free(sb.buffer);
    return res;
}

typedef struct {
    const char* current;
} KodaJsonParser;

static void skip_ws(KodaJsonParser* p) {
    while (*p->current == ' ' || *p->current == '\t' || *p->current == '\n' || *p->current == '\r') p->current++;
}

static KodaValue koda_json_parse_recursive(KodaJsonParser* p);

static KodaValue koda_json_parse_string(KodaJsonParser* p) {
    p->current++; // skip "
    const char* start = p->current;
    while (*p->current != '"' && *p->current != '\0') p->current++;
    int len = p->current - start;
    KodaString* s = koda_copy_string(start, len);
    if (*p->current == '"') p->current++;
    return OBJ_VAL(s);
}

static KodaValue koda_json_parse_array(KodaJsonParser* p) {
    p->current++; // skip [
    KodaArray* arr = koda_new_array();
    skip_ws(p);
    while (*p->current != ']' && *p->current != '\0') {
        koda_array_push(OBJ_VAL(arr), koda_json_parse_recursive(p));
        skip_ws(p);
        if (*p->current == ',') { p->current++; skip_ws(p); }
    }
    if (*p->current == ']') p->current++;
    return OBJ_VAL(arr);
}

static KodaValue koda_json_parse_object(KodaJsonParser* p) {
    p->current++; // skip {
    KodaObject* obj = koda_alloc_object();
    skip_ws(p);
    while (*p->current != '}' && *p->current != '\0') {
        KodaValue key = koda_json_parse_string(p);
        skip_ws(p);
        if (*p->current == ':') p->current++;
        skip_ws(p);
        koda_object_set(OBJ_VAL(obj), key, koda_json_parse_recursive(p));
        skip_ws(p);
        if (*p->current == ',') { p->current++; skip_ws(p); }
    }
    if (*p->current == '}') p->current++;
    return OBJ_VAL(obj);
}

static KodaValue koda_json_parse_recursive(KodaJsonParser* p) {
    skip_ws(p);
    char c = *p->current;
    if (c == '"') return koda_json_parse_string(p);
    if (c == '[') return koda_json_parse_array(p);
    if (c == '{') return koda_json_parse_object(p);
    if (c == 't') { p->current += 4; return BOOL_VAL(true); }
    if (c == 'f') { p->current += 5; return BOOL_VAL(false); }
    if (c == 'n') { p->current += 4; return NULL_VAL; }
    if ((c >= '0' && c <= '9') || c == '-') {
        char* end;
        double d = strtod(p->current, &end);
        p->current = end;
        return NUMBER_VAL(d);
    }
    return NULL_VAL;
}

KodaValue koda_json_parse(int argCount, KodaValue* args) {
    if (argCount < 1 || !IS_OBJ(args[0]) || AS_OBJ(args[0])->type != OBJ_STRING) return NULL_VAL;
    KodaJsonParser p;
    p.current = ((KodaString*)AS_OBJ(args[0]))->chars;
    return koda_json_parse_recursive(&p);
}

KodaValue koda_module_init_json() {
    static KodaValue module_exports = 0;
    if (module_exports != 0) return module_exports;
    module_exports = OBJ_VAL(koda_alloc_object());
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("parse", 5)), koda_new_native(koda_json_parse));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("stringify", 9)), koda_new_native(koda_json_stringify));
    return module_exports;
}

KodaValue koda_abs_num(KodaValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(fabs(AS_NUMBER(v))); }
KodaValue koda_sqrt_num(KodaValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(sqrt(AS_NUMBER(v))); }
KodaValue koda_cbrt_num(KodaValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(cbrt(AS_NUMBER(v))); }
KodaValue koda_sin_num(KodaValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(sin(AS_NUMBER(v))); }
KodaValue koda_cos_num(KodaValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(cos(AS_NUMBER(v))); }
KodaValue koda_tan_num(KodaValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(tan(AS_NUMBER(v))); }
KodaValue koda_asin_num(KodaValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(asin(AS_NUMBER(v))); }
KodaValue koda_acos_num(KodaValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(acos(AS_NUMBER(v))); }
KodaValue koda_atan_num(KodaValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(atan(AS_NUMBER(v))); }
KodaValue koda_atan2_num(KodaValue y, KodaValue x) { if (!IS_NUMBER(y) || !IS_NUMBER(x)) return NULL_VAL; return NUMBER_VAL(atan2(AS_NUMBER(y), AS_NUMBER(x))); }
KodaValue koda_log_num(KodaValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(log(AS_NUMBER(v))); }
KodaValue koda_log2_num(KodaValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(log2(AS_NUMBER(v))); }
KodaValue koda_log10_num(KodaValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(log10(AS_NUMBER(v))); }
KodaValue koda_exp_num(KodaValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(exp(AS_NUMBER(v))); }
KodaValue koda_floor_num(KodaValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(floor(AS_NUMBER(v))); }
KodaValue koda_ceil_num(KodaValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(ceil(AS_NUMBER(v))); }
KodaValue koda_round_num(KodaValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(round(AS_NUMBER(v))); }
KodaValue koda_trunc_num(KodaValue v) { if (!IS_NUMBER(v)) return NULL_VAL; return NUMBER_VAL(trunc(AS_NUMBER(v))); }

KodaValue koda_min_num(int argCount, KodaValue* args) {
    if (argCount == 0) return NULL_VAL;
    double res = AS_NUMBER(args[0]);
    for (int i = 1; i < argCount; i++) {
        double v = AS_NUMBER(args[i]);
        if (v < res) res = v;
    }
    return NUMBER_VAL(res);
}

KodaValue koda_max_num(int argCount, KodaValue* args) {
    if (argCount == 0) return NULL_VAL;
    double res = AS_NUMBER(args[0]);
    for (int i = 1; i < argCount; i++) {
        double v = AS_NUMBER(args[i]);
        if (v > res) res = v;
    }
    return NUMBER_VAL(res);
}

KodaValue koda_clamp_num(KodaValue v, KodaValue min, KodaValue max) {
    if (!IS_NUMBER(v) || !IS_NUMBER(min) || !IS_NUMBER(max)) return NULL_VAL;
    double dv = AS_NUMBER(v);
    double dmin = AS_NUMBER(min);
    double dmax = AS_NUMBER(max);
    if (dv < dmin) return min;
    if (dv > dmax) return max;
    return v;
}

KodaValue koda_lerp_num(KodaValue a, KodaValue b, KodaValue t) {
    if (!IS_NUMBER(a) || !IS_NUMBER(b) || !IS_NUMBER(t)) return NULL_VAL;
    double da = AS_NUMBER(a);
    double db = AS_NUMBER(b);
    double dt = AS_NUMBER(t);
    return NUMBER_VAL(da + (db - da) * dt);
}

KodaValue koda_smoothstep_num(KodaValue a, KodaValue b, KodaValue t) {
    if (!IS_NUMBER(a) || !IS_NUMBER(b) || !IS_NUMBER(t)) return NULL_VAL;
    double da = AS_NUMBER(a);
    double db = AS_NUMBER(b);
    double tv = AS_NUMBER(t);
    double dt = (tv - da) / (db - da);
    if (dt < 0) dt = 0;
    if (dt > 1) dt = 1;
    return NUMBER_VAL(dt * dt * (3 - 2 * dt));
}

KodaValue koda_map_num(KodaValue v, KodaValue inMin, KodaValue inMax, KodaValue outMin, KodaValue outMax) {
    if (!IS_NUMBER(v) || !IS_NUMBER(inMin) || !IS_NUMBER(inMax) || !IS_NUMBER(outMin) || !IS_NUMBER(outMax)) return NULL_VAL;
    double dv = AS_NUMBER(v);
    double dim = AS_NUMBER(inMin);
    double dix = AS_NUMBER(inMax);
    double dom = AS_NUMBER(outMin);
    double dox = AS_NUMBER(outMax);
    return NUMBER_VAL(dom + (dox - dom) * ((dv - dim) / (dix - dim)));
}

KodaValue koda_sign_num(KodaValue v) {
    if (!IS_NUMBER(v)) return NULL_VAL;
    double dv = AS_NUMBER(v);
    if (dv > 0) return NUMBER_VAL(1);
    if (dv < 0) return NUMBER_VAL(-1);
    return NUMBER_VAL(0);
}

KodaValue koda_hypot_num(KodaValue a, KodaValue b) {
    if (!IS_NUMBER(a) || !IS_NUMBER(b)) return NULL_VAL;
    return NUMBER_VAL(hypot(AS_NUMBER(a), AS_NUMBER(b)));
}

KodaValue koda_distance_num(KodaValue x1, KodaValue y1, KodaValue x2, KodaValue y2) {
    if (!IS_NUMBER(x1) || !IS_NUMBER(y1) || !IS_NUMBER(x2) || !IS_NUMBER(y2)) return NULL_VAL;
    double dx = AS_NUMBER(x2) - AS_NUMBER(x1);
    double dy = AS_NUMBER(y2) - AS_NUMBER(y1);
    return NUMBER_VAL(sqrt(dx * dx + dy * dy));
}

KodaValue koda_distance_sq_num(KodaValue x1, KodaValue y1, KodaValue x2, KodaValue y2) {
    if (!IS_NUMBER(x1) || !IS_NUMBER(y1) || !IS_NUMBER(x2) || !IS_NUMBER(y2)) return NULL_VAL;
    double dx = AS_NUMBER(x2) - AS_NUMBER(x1);
    double dy = AS_NUMBER(y2) - AS_NUMBER(y1);
    return NUMBER_VAL(dx * dx + dy * dy);
}

KodaValue koda_angle_between_num(KodaValue x1, KodaValue y1, KodaValue x2, KodaValue y2) {
    if (!IS_NUMBER(x1) || !IS_NUMBER(y1) || !IS_NUMBER(x2) || !IS_NUMBER(y2)) return NULL_VAL;
    return NUMBER_VAL(atan2(AS_NUMBER(y2) - AS_NUMBER(y1), AS_NUMBER(x2) - AS_NUMBER(x1)));
}

KodaValue koda_normalize_num(KodaValue x, KodaValue y) {
    if (!IS_NUMBER(x) || !IS_NUMBER(y)) return NULL_VAL;
    double dx = AS_NUMBER(x);
    double dy = AS_NUMBER(y);
    double len = sqrt(dx * dx + dy * dy);
    double nx = 0, ny = 0;
    if (len > 0) {
        nx = dx / len;
        ny = dy / len;
    }
    KodaValue objVal = OBJ_VAL(koda_alloc_object());
    koda_object_set(objVal, OBJ_VAL(koda_copy_string("x", 1)), NUMBER_VAL(nx));
    koda_object_set(objVal, OBJ_VAL(koda_copy_string("y", 1)), NUMBER_VAL(ny));
    return objVal;
}

KodaValue koda_delta_time(void) {
    double now = koda_get_time_seconds();
    double dt = now - heap.last_frame;
    heap.last_frame = now;
    if (dt > 0.1) dt = 0.1; 
    return NUMBER_VAL(dt);
}

KodaValue koda_time(void) {
    return NUMBER_VAL(koda_get_time_seconds() - heap.program_start);
}

KodaValue koda_timestamp(void) {
    return NUMBER_VAL((double)time(NULL));
}

KodaValue koda_random(int argCount, KodaValue* args) {
    double r = (double)rand() / (double)RAND_MAX;
    if (argCount == 0) return NUMBER_VAL(r);
    if (argCount == 1) return NUMBER_VAL(r * AS_NUMBER(args[0]));
    double min = AS_NUMBER(args[0]);
    double max = AS_NUMBER(args[1]);
    return NUMBER_VAL(min + r * (max - min));
}

KodaValue koda_random_int(int argCount, KodaValue* args) {
    if (argCount == 1) return NUMBER_VAL(rand() % (int)AS_NUMBER(args[0]));
    int min = (int)AS_NUMBER(args[0]);
    int max = (int)AS_NUMBER(args[1]);
    if (max <= min) return NUMBER_VAL(min);
    return NUMBER_VAL(min + (rand() % (max - min)));
}

KodaValue koda_random_choice(KodaValue array) {
    if (!IS_OBJ(array) || AS_OBJ(array)->type != OBJ_ARRAY) return NULL_VAL;
    KodaArray* arr = (KodaArray*)AS_OBJ(array);
    if (arr->count == 0) return NULL_VAL;
    return arr->elements[rand() % arr->count];
}

void koda_random_seed(KodaValue seed) {
    if (IS_NUMBER(seed)) srand((unsigned int)AS_NUMBER(seed));
}

KodaValue koda_set_timeout(KodaValue closure, KodaValue ms) {
    if (!IS_NUMBER(ms)) return NULL_VAL;
    KodaTask* task = (KodaTask*)malloc(sizeof(KodaTask));
    task->closure = closure;
    task->target_time = koda_get_time_seconds() + (AS_NUMBER(ms) / 1000.0);
    task->interval = 0;
    task->next = heap.tasks;
    heap.tasks = task;
    return NULL_VAL;
}

KodaValue koda_set_interval(KodaValue closure, KodaValue ms) {
    if (!IS_NUMBER(ms)) return NULL_VAL;
    double interval = AS_NUMBER(ms) / 1000.0;
    KodaTask* task = (KodaTask*)malloc(sizeof(KodaTask));
    task->closure = closure;
    task->target_time = koda_get_time_seconds() + interval;
    task->interval = interval;
    task->next = heap.tasks;
    heap.tasks = task;
    return NULL_VAL;
}

void koda_poll_tasks() {
    double now = koda_get_time_seconds();
    KodaTask** prev = &heap.tasks;
    while (*prev) {
        KodaTask* task = *prev;
        if (now >= task->target_time) {
            // Call closure
            koda_call(task->closure, 0, NULL);
            
            if (task->interval > 0) {
                task->target_time = now + task->interval;
                prev = &task->next;
            } else {
                *prev = task->next;
                free(task);
            }
        } else {
            prev = &task->next;
        }
    }
}

KodaValue koda_math_sin(int argCount, KodaValue* args) { if (argCount < 1 || !IS_NUMBER(args[0])) return NULL_VAL; return NUMBER_VAL(sin(AS_NUMBER(args[0]))); }
KodaValue koda_math_cos(int argCount, KodaValue* args) { if (argCount < 1 || !IS_NUMBER(args[0])) return NULL_VAL; return NUMBER_VAL(cos(AS_NUMBER(args[0]))); }
KodaValue koda_math_tan(int argCount, KodaValue* args) { if (argCount < 1 || !IS_NUMBER(args[0])) return NULL_VAL; return NUMBER_VAL(tan(AS_NUMBER(args[0]))); }
KodaValue koda_math_sqrt(int argCount, KodaValue* args) { if (argCount < 1 || !IS_NUMBER(args[0])) return NULL_VAL; return NUMBER_VAL(sqrt(AS_NUMBER(args[0]))); }
KodaValue koda_math_abs(int argCount, KodaValue* args) { if (argCount < 1 || !IS_NUMBER(args[0])) return NULL_VAL; return NUMBER_VAL(fabs(AS_NUMBER(args[0]))); }
KodaValue koda_math_floor(int argCount, KodaValue* args) { if (argCount < 1 || !IS_NUMBER(args[0])) return NULL_VAL; return NUMBER_VAL(floor(AS_NUMBER(args[0]))); }
KodaValue koda_math_ceil(int argCount, KodaValue* args) { if (argCount < 1 || !IS_NUMBER(args[0])) return NULL_VAL; return NUMBER_VAL(ceil(AS_NUMBER(args[0]))); }
KodaValue koda_math_clamp(int argCount, KodaValue* args) { if (argCount < 3) return NULL_VAL; return koda_clamp_num(args[0], args[1], args[2]); }
KodaValue koda_math_lerp(int argCount, KodaValue* args) { if (argCount < 3) return NULL_VAL; return koda_lerp_num(args[0], args[1], args[2]); }
KodaValue koda_math_sign(int argCount, KodaValue* args) { if (argCount < 1) return NULL_VAL; return koda_sign_num(args[0]); }
KodaValue koda_math_random(int argCount, KodaValue* args) { return koda_random(argCount, args); }
KodaValue koda_math_randomInt(int argCount, KodaValue* args) { return koda_random_int(argCount, args); }
KodaValue koda_math_randomChoice(int argCount, KodaValue* args) { if (argCount < 1) return NULL_VAL; return koda_random_choice(args[0]); }
KodaValue koda_math_randomSeed(int argCount, KodaValue* args) { if (argCount >= 1) koda_random_seed(args[0]); return NULL_VAL; }
KodaValue koda_math_deltaTime(int argCount, KodaValue* args) { (void)argCount; (void)args; return koda_delta_time(); }
KodaValue koda_math_time(int argCount, KodaValue* args) { (void)argCount; (void)args; return koda_time(); }
KodaValue koda_math_timestamp(int argCount, KodaValue* args) { (void)argCount; (void)args; return koda_timestamp(); }
KodaValue koda_math_sleep(int argCount, KodaValue* args) { if (argCount >= 1) koda_sleep_ms(args[0]); return NULL_VAL; }

KodaValue koda_module_init_math() {
    static KodaValue module_exports = 0;
    if (module_exports != 0) return module_exports;
    module_exports = OBJ_VAL(koda_alloc_object());
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("sin", 3)), koda_new_native(koda_math_sin));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("cos", 3)), koda_new_native(koda_math_cos));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("tan", 3)), koda_new_native(koda_math_tan));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("sqrt", 4)), koda_new_native(koda_math_sqrt));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("abs", 3)), koda_new_native(koda_math_abs));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("floor", 5)), koda_new_native(koda_math_floor));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("ceil", 4)), koda_new_native(koda_math_ceil));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("clamp", 5)), koda_new_native(koda_math_clamp));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("lerp", 4)), koda_new_native(koda_math_lerp));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("sign", 4)), koda_new_native(koda_math_sign));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("random", 6)), koda_new_native(koda_math_random));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("randomInt", 9)), koda_new_native(koda_math_randomInt));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("randomChoice", 12)), koda_new_native(koda_math_randomChoice));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("randomSeed", 10)), koda_new_native(koda_math_randomSeed));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("deltaTime", 9)), koda_new_native(koda_math_deltaTime));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("time", 4)), koda_new_native(koda_math_time));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("timestamp", 9)), koda_new_native(koda_math_timestamp));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("sleep", 5)), koda_new_native(koda_math_sleep));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("pi", 2)), NUMBER_VAL(3.14159265358979323846));
    koda_object_set(module_exports, OBJ_VAL(koda_copy_string("e", 1)), NUMBER_VAL(2.71828182845904523536));
    return module_exports;
}

