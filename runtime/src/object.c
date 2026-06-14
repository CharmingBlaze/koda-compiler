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
    ObjArray* array = (ObjArray*)gc_alloc(sizeof(ObjArray));
    array->obj.type = OBJ_ARRAY;
    array->obj.is_marked = false;
    array->obj.next = NULL;
    array->capacity = capacity;
    array->count = 0;
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
        table->hashes = (uint32_t*)calloc((size_t)capacity, sizeof(uint32_t));
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

void free_object(Obj* obj) {
    switch (obj->type) {
        case OBJ_STRING: {
            ObjString* string = (ObjString*)obj;
            size_t total = sizeof(ObjString) + (size_t)string->length + 1u;
            gc_free(string, total);
            break;
        }
        case OBJ_ARRAY: {
            ObjArray* array = (ObjArray*)obj;
            if (array->elements != NULL) {
                gc_free(array->elements, (size_t)array->capacity * sizeof(Value));
            }
            gc_free(array, sizeof(ObjArray));
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
                // hashes is allocated with calloc (not gc_alloc), so release it with raw free
                // to avoid corrupting GC bytes_allocated accounting.
                free(table->hashes);
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
    }
}
