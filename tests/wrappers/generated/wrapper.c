#include "koda.h"
#include <string.h>
#include "mylib.h"

KodaValue koda_wrap_mylib_add_ints(int argCount, KodaValue* args) {
    if (argCount < 2) return NULL_VAL;
    int arg0 = (IS_NUMBER(args[0]) ? AS_NUMBER(args[0]) : 0);
    int arg1 = (IS_NUMBER(args[1]) ? AS_NUMBER(args[1]) : 0);
    int result = add_ints(arg0, arg1);
    return NUMBER_VAL((double)result);
}

KodaValue koda_wrap_mylib_triple_int(int argCount, KodaValue* args) {
    if (argCount < 1) return NULL_VAL;
    int arg0 = (IS_NUMBER(args[0]) ? AS_NUMBER(args[0]) : 0);
    int result = triple_int(arg0);
    return NUMBER_VAL((double)result);
}

