#ifndef KODA_OBJECT_H
#define KODA_OBJECT_H

#include "value.h"
#include <stddef.h>

// Object types
typedef enum {
    OBJ_STRING,
    OBJ_ARRAY,
    OBJ_TABLE,
    OBJ_CLOSURE,
    OBJ_FUNCTION,
    OBJ_NATIVE,
    OBJ_CELL,
    OBJ_ARENA
} ObjType;

#define GEN_NURSERY 0
#define GEN_YOUNG   1
#define GEN_OLD     2
#define GEN_DEAD    3

// Base object structure
struct Obj {
    ObjType type;
    bool is_marked;
    uint8_t generation;
    bool in_remembered_set;
    bool in_arena;
    struct Obj* next;
};

// String object
typedef struct {
    Obj obj;
    int length;
    uint32_t hash; /* FNV-1a; 0 = not computed yet */
    char chars[];
} ObjString;

#define KODA_ARRAY_INLINE_MAX_CAP 64

// Array object
typedef struct {
    Obj obj;
    int capacity;
    int count;
    Value* elements;
    bool inline_elements;
} ObjArray;

// Table object (hash table)
typedef struct {
    Obj obj;
    int capacity;
    int count;
    Value* keys;
    Value* values;
    /* When true, values[i] is field slot i (struct instance); property hash is disabled. */
    bool is_struct_layout;
    uint32_t* hashes; /* parallel to keys/values when using open addressing; NULL for linear / struct */
} ObjTable;

// Function object
typedef struct {
    Obj obj;
    int arity;
    int param_count;
    Value* params;
    Value* body;
} ObjFunction;

// Closure object
typedef struct {
    Obj obj;
    ObjFunction* function;
    Value* upvalues;
    int upvalue_count;
} ObjClosure;

// Native function object
typedef struct {
    Obj obj;
    Value (*native)(int arg_count, Value* args);
} ObjNative;

/* Mutable binding cell (closure capture / escaping local). Header is gc-managed. */
typedef struct {
    Obj obj;
    Value value;
} ObjCell;

/* Bump allocator for per-frame / scoped short-lived objects (see arena builtins). */
typedef struct {
    Obj obj;
    uint8_t* buffer;
    size_t capacity;
    size_t top;
    Obj** allocations;
    int alloc_count;
    int alloc_capacity;
} ObjArena;

// Object allocation
ObjString* allocate_string(int length);
ObjArray* allocate_array(int capacity);
ObjTable* allocate_table(int capacity);
ObjTable* allocate_struct_table(int field_count);
ObjFunction* allocate_function(int arity);
ObjClosure* allocate_closure(ObjFunction* function, int upvalue_count);
ObjNative* allocate_native(Value (*native)(int, Value*));
ObjCell* allocate_cell(void);
ObjArena* allocate_arena(size_t capacity);
void arena_reset(ObjArena* arena);
ObjArray* arena_allocate_array(ObjArena* arena, int capacity);
ObjTable* arena_allocate_struct_table(ObjArena* arena, int field_count);

// Object operations
void free_object(Obj* obj);
void print_object(Value v);

#endif // KODA_OBJECT_H
