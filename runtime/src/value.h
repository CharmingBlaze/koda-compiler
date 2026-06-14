#ifndef KODA_VALUE_H
#define KODA_VALUE_H

#include <stdint.h>
#include <stdbool.h>

/*
 * Koda Value System - NaN Boxing
 * 
 * Uses NaN-boxing to represent all Koda values in a single 64-bit integer.
 * This allows efficient storage and passing of values without boxing overhead.
 * 
 * Bit layout:
 * - QNaN (quiet NaN): 0x7ffc000000000000 to 0x7ffc00000000007f
 * - Tagged values: Use the lower bits to distinguish types
 * - Numbers: Use the full 64-bit range (except NaN patterns)
 * - Objects: Store pointer in lower bits, tag in upper bits
 */

typedef uint64_t Value;

// Value tags (stored in NaN-boxed values)
#define TAG_NIL    0x01
#define TAG_FALSE  0x02
#define TAG_TRUE   0x03
#define TAG_OBJ    0x04

// NaN-boxed constants
#define QNAN      ((uint64_t)0x7ffc000000000000)
#define SIGN_BIT  ((uint64_t)0x8000000000000000)

// Tagged value constructors
#define NIL_VAL     (QNAN | TAG_NIL)
#define FALSE_VAL   (QNAN | TAG_FALSE)
#define TRUE_VAL    (QNAN | TAG_TRUE)

// Type predicates
static inline bool IS_NIL(Value v)    { return v == NIL_VAL; }
static inline bool IS_FALSE(Value v)  { return v == FALSE_VAL; }
static inline bool IS_TRUE(Value v)   { return v == TRUE_VAL; }
static inline bool IS_BOOL(Value v)   { return IS_FALSE(v) || IS_TRUE(v); }
static inline bool IS_OBJ(Value v)    { return (v & QNAN) == QNAN && (v & 0x7) == TAG_OBJ; }

// Boolean conversion
static inline bool AS_BOOL(Value v)   { return IS_TRUE(v); }

// Number handling
static inline bool IS_NUMBER(Value v) { return (v & QNAN) != QNAN; }
static inline double AS_NUMBER(Value v) { 
    union { uint64_t i; double d; } u;
    u.i = v;
    return u.d;
}
static inline Value NUMBER_VAL(double n) {
    union { uint64_t i; double d; } u;
    u.d = n;
    return u.i;
}

// Object handling
typedef struct Obj Obj;

static inline Obj* AS_OBJ(Value v) {
    return (Obj*)(uintptr_t)(v & ~(QNAN | SIGN_BIT | 0x7));
}

static inline Value OBJ_VAL(Obj* obj) {
    return (uint64_t)(uintptr_t)obj | QNAN | TAG_OBJ;
}

// Value equality
bool values_equal(Value a, Value b);
/** Same as values_equal, exposed for native codegen (returns 1 or 0). */
int64_t koda_values_equal(Value a, Value b);

// Value printing
void print_value(Value v);

#endif // KODA_VALUE_H
