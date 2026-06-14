#ifndef KODA_WRAPGEN_ABI_H
#define KODA_WRAPGEN_ABI_H

/*
 * ABI shim for cmd/wrapgen-generated wrapper.c when linked with libkoda_runtime.a.
 * Wrapgen historically emitted Koda* names from the embedded alternate header; the
 * shipped runtime uses Value / ObjString / NIL_VAL instead.
 */
#include "value.h"
#include "object.h"
#include "koda_runtime.h"

typedef Value KodaValue;
typedef ObjString KodaString;

#ifndef NULL_VAL
#define NULL_VAL NIL_VAL
#endif

#endif /* KODA_WRAPGEN_ABI_H */
