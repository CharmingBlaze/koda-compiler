#ifndef KODA_SHADOW_STACK_H
#define KODA_SHADOW_STACK_H

#include "value.h"
#include <stdbool.h>

#define KODA_SHADOW_STACK_INITIAL_CAPACITY 4096
#define KODA_SHADOW_STACK_MAX_CAPACITY 65536

typedef struct {
    Value** slot_ptrs;
    int count;
} KodaShadowFrame;

extern KodaShadowFrame* koda_shadow_stack;
extern int koda_shadow_depth;
extern int koda_shadow_capacity;

void koda_push_frame(Value** slot_ptrs, int count);
void koda_pop_frame(void);
int koda_get_shadow_depth(void);

#endif
