#ifndef KODA_GC_H
#define KODA_GC_H

#include "object.h"
#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>

#if defined(__clang__) || defined(__GNUC__)
#if defined(__x86_64__) || defined(_M_X64)
#define GC_GET_SP(sp) __asm__ volatile("mov %%rsp, %0" : "=r"(sp))
#elif defined(__aarch64__) || defined(_M_ARM64)
#define GC_GET_SP(sp) __asm__ volatile("mov %0, sp" : "=r"(sp))
#endif
#endif

#ifndef GC_GET_SP
/* Address of the local `sp` slot: stack anchor when no arch-specific SP read exists. */
#define GC_GET_SP(sp) ((sp) = (void*)(uintptr_t)&(sp))
#endif

typedef struct {
    uint64_t collections;
    uint64_t bytes_allocated;
    uint64_t bytes_freed;
    uint64_t max_pause_time_us;
    uint64_t total_pause_time_us;
} GCStats;

void gc_init(void);
void* gc_alloc(size_t size);
void gc_register_object(Obj* obj);
void gc_unlink_object(Obj* obj);
void gc_free(void* ptr, size_t size);

void gc_mark_object(Obj* obj);
void gc_mark_value(Value v);
void gc_write_barrier(Obj* parent, Value new_val);

void gc_mark_roots(void);
void gc_mark_stack_conservative(void);
void gc_mark_shadow_stack(void);
void gc_sweep(void);
void gc_collect(void);
void gc_collect_minor(void);
void gc_collect_incremental(uint64_t budget_us);
void gc_frame_step_incremental(uint64_t budget_us);
bool gc_incremental_is_idle(void);
void gc_record_pause_since(uint64_t start_us);
void gc_set_use_shadow_stack(bool enabled);
bool gc_uses_shadow_stack(void);

size_t gc_nursery_used_bytes(void);
size_t gc_nursery_capacity_bytes(void);

GCStats gc_get_stats(void);
void gc_set_disabled(bool disabled);
void gc_set_next_threshold(size_t bytes);
uint64_t gc_debug_remembered_overflow_count(void);

#endif
