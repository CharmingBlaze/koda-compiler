#include "object.h"
#include "gc.h"
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

static void validate_value_slot_count(int n, const char* ctx) {
    if (n < 0 || (sizeof(Value) > 0 && (size_t)n > SIZE_MAX / sizeof(Value))) {
        fprintf(stderr, "koda: %s capacity overflow: %d\n", ctx, n);
        exit(1);
    }
}

ObjString* allocate_string(int length) {
    size_t total = sizeof(ObjString) + (size_t)length + 1u;
    ObjString* string = (ObjString*)gc_alloc(total);
    string->obj.type = OBJ_STRING;
    string->obj.is_marked = false;
    string->obj.next = NULL;
    string->length = length;
    string->hash = 0;
    string->chars[length] = '\0';
    gc_register_object(&string->obj);
    return string;
}

ObjArray* allocate_array(int capacity) {
    validate_value_slot_count(capacity, "array");
    if (capacity <= KODA_ARRAY_INLINE_MAX_CAP) {
        size_t total = sizeof(ObjArray) + sizeof(Value) * (size_t)capacity;
        ObjArray* array = (ObjArray*)gc_alloc(total);
        array->obj.type = OBJ_ARRAY;
        array->obj.is_marked = false;
        array->obj.next = NULL;
        array->capacity = capacity;
        array->count = 0;
        array->inline_elements = true;
        array->elements = (Value*)(array + 1);
        gc_register_object(&array->obj);
        return array;
    }
    ObjArray* array = (ObjArray*)gc_alloc(sizeof(ObjArray));
    array->obj.type = OBJ_ARRAY;
    array->obj.is_marked = false;
    array->obj.next = NULL;
    array->capacity = capacity;
    array->count = 0;
    array->inline_elements = false;
    array->elements = NULL;
    gc_register_object(&array->obj);
    array->elements = (Value*)gc_alloc(sizeof(Value) * (size_t)capacity);
    return array;
}

ObjTable* allocate_table(int capacity) {
    validate_value_slot_count(capacity, "table");
    ObjTable* table = (ObjTable*)gc_alloc(sizeof(ObjTable));
    table->obj.type = OBJ_TABLE;
    table->obj.is_marked = false;
    table->obj.next = NULL;
    table->capacity = capacity;
    table->count = 0;
    table->keys = NULL;
    table->values = NULL;
    table->is_struct_layout = false;
    table->hashes = NULL;
    gc_register_object(&table->obj);
    table->keys = (Value*)gc_alloc(sizeof(Value) * (size_t)capacity);
    table->values = (Value*)gc_alloc(sizeof(Value) * (size_t)capacity);
    for (int i = 0; i < capacity; i++) {
        table->keys[i] = NIL_VAL;
        table->values[i] = NIL_VAL;
    }
    if (capacity >= 8) {
        table->hashes = (uint32_t*)gc_alloc(sizeof(uint32_t) * (size_t)capacity);
        memset(table->hashes, 0, sizeof(uint32_t) * (size_t)capacity);
    }
    return table;
}

ObjTable* allocate_struct_table(int field_count) {
    if (field_count < 1) {
        field_count = 1;
    }
    validate_value_slot_count(field_count, "struct fields");
    ObjTable* table = (ObjTable*)gc_alloc(sizeof(ObjTable));
    table->obj.type = OBJ_TABLE;
    table->obj.is_marked = false;
    table->obj.next = NULL;
    table->capacity = field_count;
    table->count = field_count;
    table->is_struct_layout = true;
    table->hashes = NULL;
    gc_register_object(&table->obj);
    table->keys = (Value*)gc_alloc(sizeof(Value) * (size_t)field_count);
    table->values = (Value*)gc_alloc(sizeof(Value) * (size_t)field_count);
    for (int i = 0; i < field_count; i++) {
        table->keys[i] = NIL_VAL;
        table->values[i] = NIL_VAL;
    }
    return table;
}

ObjFunction* allocate_function(int arity) {
    ObjFunction* function = (ObjFunction*)gc_alloc(sizeof(ObjFunction));
    function->obj.type = OBJ_FUNCTION;
    function->obj.is_marked = false;
    function->obj.next = NULL;
    function->arity = arity;
    function->param_count = 0;
    function->params = NULL;
    function->body = NULL;
    gc_register_object(&function->obj);
    return function;
}

ObjClosure* allocate_closure(ObjFunction* function, int upvalue_count) {
    validate_value_slot_count(upvalue_count, "closure upvalues");
    ObjClosure* closure = (ObjClosure*)gc_alloc(sizeof(ObjClosure));
    closure->obj.type = OBJ_CLOSURE;
    closure->obj.is_marked = false;
    closure->obj.next = NULL;
    closure->function = function;
    closure->upvalues = NULL;
    closure->upvalue_count = 0;
    gc_register_object(&closure->obj);
    closure->upvalues = (Value*)gc_alloc(sizeof(Value) * (size_t)upvalue_count);
    closure->upvalue_count = upvalue_count;
    for (int i = 0; i < upvalue_count; i++) {
        closure->upvalues[i] = NIL_VAL;
    }
    return closure;
}

ObjNative* allocate_native(Value (*native)(int, Value*)) {
    ObjNative* native_obj = (ObjNative*)gc_alloc(sizeof(ObjNative));
    native_obj->obj.type = OBJ_NATIVE;
    native_obj->obj.is_marked = false;
    native_obj->obj.next = NULL;
    native_obj->native = native;
    gc_register_object(&native_obj->obj);
    return native_obj;
}

ObjCell* allocate_cell(void) {
    ObjCell* cell = (ObjCell*)gc_alloc(sizeof(ObjCell));
    cell->obj.type = OBJ_CELL;
    cell->obj.is_marked = false;
    cell->obj.next = NULL;
    cell->value = NIL_VAL;
    gc_register_object(&cell->obj);
    return cell;
}

static void arena_track_allocation(ObjArena* arena, Obj* obj) {
    if (arena->alloc_count >= arena->alloc_capacity) {
        int next_cap = arena->alloc_capacity == 0 ? 64 : arena->alloc_capacity * 2;
        Obj** next = (Obj**)realloc(arena->allocations, sizeof(Obj*) * (size_t)next_cap);
        if (next == NULL) {
            fprintf(stderr, "koda: out of memory tracking arena allocations\n");
            exit(1);
        }
        arena->allocations = next;
        arena->alloc_capacity = next_cap;
    }
    arena->allocations[arena->alloc_count++] = obj;
}

static void* arena_bump(ObjArena* arena, size_t nbytes) {
    size_t top = (arena->top + 7u) & ~7u;
    if (top + nbytes > arena->capacity) {
        fprintf(stderr, "koda: arena out of memory\n");
        exit(1);
    }
    void* ptr = arena->buffer + top;
    arena->top = top + nbytes;
    return ptr;
}

ObjArena* allocate_arena(size_t capacity) {
    if (capacity == 0) {
        capacity = 1;
    }
    ObjArena* arena = (ObjArena*)gc_alloc(sizeof(ObjArena));
    arena->obj.type = OBJ_ARENA;
    arena->obj.is_marked = false;
    arena->obj.next = NULL;
    arena->buffer = (uint8_t*)malloc(capacity);
    if (arena->buffer == NULL) {
        fprintf(stderr, "koda: out of memory allocating arena buffer\n");
        exit(1);
    }
    arena->capacity = capacity;
    arena->top = 0;
    arena->allocations = NULL;
    arena->alloc_count = 0;
    arena->alloc_capacity = 0;
    gc_register_object(&arena->obj);
    return arena;
}

void arena_reset(ObjArena* arena) {
    if (arena == NULL) {
        return;
    }
    for (int i = 0; i < arena->alloc_count; i++) {
        Obj* child = arena->allocations[i];
        if (child != NULL) {
            gc_unlink_object(child);
            child->in_arena = false;
        }
    }
    arena->alloc_count = 0;
    arena->top = 0;
}

ObjArray* arena_allocate_array(ObjArena* arena, int capacity) {
    validate_value_slot_count(capacity, "arena array");
    if (capacity <= KODA_ARRAY_INLINE_MAX_CAP) {
        size_t total = sizeof(ObjArray) + sizeof(Value) * (size_t)capacity;
        ObjArray* array = (ObjArray*)arena_bump(arena, total);
        array->obj.type = OBJ_ARRAY;
        array->obj.is_marked = false;
        array->obj.next = NULL;
        array->capacity = capacity;
        array->count = 0;
        array->inline_elements = true;
        array->elements = (Value*)(array + 1);
        gc_register_object(&array->obj);
        array->obj.in_arena = true;
        arena_track_allocation(arena, &array->obj);
        return array;
    }
    size_t header = sizeof(ObjArray);
    ObjArray* array = (ObjArray*)arena_bump(arena, header);
    array->obj.type = OBJ_ARRAY;
    array->obj.is_marked = false;
    array->obj.next = NULL;
    array->capacity = capacity;
    array->count = 0;
    array->inline_elements = false;
    array->elements = NULL;
    gc_register_object(&array->obj);
    array->obj.in_arena = true;
    arena_track_allocation(arena, &array->obj);
    array->elements = (Value*)arena_bump(arena, sizeof(Value) * (size_t)capacity);
    return array;
}

ObjTable* arena_allocate_struct_table(ObjArena* arena, int field_count) {
    if (field_count < 1) {
        field_count = 1;
    }
    validate_value_slot_count(field_count, "arena struct fields");
    size_t header = sizeof(ObjTable);
    size_t slots = sizeof(Value) * (size_t)field_count;
    ObjTable* table = (ObjTable*)arena_bump(arena, header);
    table->obj.type = OBJ_TABLE;
    table->obj.is_marked = false;
    table->obj.next = NULL;
    table->capacity = field_count;
    table->count = field_count;
    table->is_struct_layout = true;
    table->hashes = NULL;
    gc_register_object(&table->obj);
    table->obj.in_arena = true;
    arena_track_allocation(arena, &table->obj);
    table->keys = (Value*)arena_bump(arena, slots);
    table->values = (Value*)arena_bump(arena, slots);
    for (int i = 0; i < field_count; i++) {
        table->keys[i] = NIL_VAL;
        table->values[i] = NIL_VAL;
    }
    return table;
}

void free_object(Obj* obj) {
    if (obj->in_arena) {
        return;
    }
    switch (obj->type) {
        case OBJ_STRING: {
            ObjString* string = (ObjString*)obj;
            size_t total = sizeof(ObjString) + (size_t)string->length + 1u;
            gc_free(string, total);
            break;
        }
        case OBJ_ARRAY: {
            ObjArray* array = (ObjArray*)obj;
            if (array->inline_elements) {
                size_t total = sizeof(ObjArray) + (size_t)array->capacity * sizeof(Value);
                gc_free(array, total);
            } else {
                if (array->elements != NULL) {
                    gc_free(array->elements, (size_t)array->capacity * sizeof(Value));
                }
                gc_free(array, sizeof(ObjArray));
            }
            break;
        }
        case OBJ_TABLE: {
            ObjTable* table = (ObjTable*)obj;
            if (table->keys != NULL) {
                gc_free(table->keys, (size_t)table->capacity * sizeof(Value));
            }
            if (table->values != NULL) {
                gc_free(table->values, (size_t)table->capacity * sizeof(Value));
            }
            if (table->hashes != NULL) {
                gc_free(table->hashes, (size_t)table->capacity * sizeof(uint32_t));
            }
            gc_free(table, sizeof(ObjTable));
            break;
        }
        case OBJ_FUNCTION: {
            ObjFunction* function = (ObjFunction*)obj;
            gc_free(function, sizeof(ObjFunction));
            break;
        }
        case OBJ_CLOSURE: {
            ObjClosure* closure = (ObjClosure*)obj;
            if (closure->upvalues != NULL) {
                gc_free(closure->upvalues, (size_t)closure->upvalue_count * sizeof(Value));
            }
            gc_free(closure, sizeof(ObjClosure));
            break;
        }
        case OBJ_NATIVE: {
            ObjNative* native_obj = (ObjNative*)obj;
            gc_free(native_obj, sizeof(ObjNative));
            break;
        }
        case OBJ_CELL: {
            ObjCell* cell = (ObjCell*)obj;
            gc_free(cell, sizeof(ObjCell));
            break;
        }
        case OBJ_ARENA: {
            ObjArena* arena = (ObjArena*)obj;
            for (int i = 0; i < arena->alloc_count; i++) {
                if (arena->allocations[i] != NULL) {
                    gc_unlink_object(arena->allocations[i]);
                }
            }
            if (arena->allocations != NULL) {
                free(arena->allocations);
            }
            if (arena->buffer != NULL) {
                free(arena->buffer);
            }
            gc_free(arena, sizeof(ObjArena));
            break;
        }
    }
}

void print_object(Value v) {
    if (!IS_OBJ(v)) {
        print_value(v);
        return;
    }

    Obj* obj = AS_OBJ(v);
    switch (obj->type) {
        case OBJ_STRING: {
            ObjString* string = (ObjString*)obj;
            printf("%s", string->chars);
            break;
        }
        case OBJ_ARRAY: {
            ObjArray* array = (ObjArray*)obj;
            printf("[");
            for (int i = 0; i < array->count; i++) {
                if (i > 0) printf(", ");
                print_value(array->elements[i]);
            }
            printf("]");
            break;
        }
        case OBJ_TABLE:
            printf("[table]");
            break;
        case OBJ_FUNCTION:
            printf("[function]");
            break;
        case OBJ_CLOSURE:
            printf("[closure]");
            break;
        case OBJ_NATIVE:
            printf("[native]");
            break;
        case OBJ_CELL:
            printf("[cell]");
            break;
        case OBJ_ARENA:
            printf("[arena]");
            break;
    }
}
