#include "value.h"
#include "object.h"
#include <stdio.h>
#include <string.h>

bool values_equal(Value a, Value b) {
    if (IS_NUMBER(a) && IS_NUMBER(b)) {
        return AS_NUMBER(a) == AS_NUMBER(b);
    }
    if (IS_OBJ(a) && IS_OBJ(b)) {
        Obj* oa = AS_OBJ(a);
        Obj* ob = AS_OBJ(b);
        if (oa->type == OBJ_STRING && ob->type == OBJ_STRING) {
            if (oa == ob) {
                return true;
            }
            ObjString* sa = (ObjString*)oa;
            ObjString* sb = (ObjString*)ob;
            if (sa->length != sb->length) {
                return false;
            }
            return memcmp(sa->chars, sb->chars, (size_t)sa->length) == 0;
        }
    }
    return a == b;
}

int64_t koda_values_equal(Value a, Value b) {
    return values_equal(a, b) ? 1 : 0;
}

void print_value(Value v) {
    if (IS_NIL(v)) {
        printf("nil");
    } else if (IS_FALSE(v)) {
        printf("false");
    } else if (IS_TRUE(v)) {
        printf("true");
    } else if (IS_NUMBER(v)) {
        printf("%g", AS_NUMBER(v));
    } else if (IS_OBJ(v)) {
        Obj* o = AS_OBJ(v);
        if (o->type == OBJ_STRING) {
            ObjString* s = (ObjString*)o;
            printf("%.*s", s->length, s->chars);
        } else {
            printf("[object]");
        }
    } else {
        printf("?");
    }
}
