# Koda/Koda Memory Architecture Implementation Brief

You are implementing the memory system for the Koda/Koda compiler and runtime.
This is a JavaScript-syntax language that compiles to native binaries via LLVM IR.
The runtime is written in C with NaN-boxed values (uint64_t everywhere).

Read this entire brief before touching any code.

---

## What already exists — know this cold before writing anything

### C runtime (runtime/src/)

value.h — NaN-boxing macros:
  typedef uint64_t Value;
  #define QNAN      0x7ffc000000000000ULL
  #define SIGN_BIT  0x8000000000000000ULL
  #define NIL_VAL   (QNAN | 0x01)
  #define FALSE_VAL (QNAN | 0x02)
  #define TRUE_VAL  (QNAN | 0x03)
  IS_NUMBER(v)  — true if (v & QNAN) != QNAN
  IS_OBJ(v)     — true if (v & QNAN) == QNAN && (v & 0x7) == 0x04
  AS_OBJ(v)     — extracts Obj* from NaN-boxed pointer
  NUMBER_VAL(d) — boxes double d into Value
  AS_NUMBER(v)  — unboxes Value to double (union bitcast)

gc.h / gc.c — mark-sweep GC:
  void  gc_init(void)
  void* gc_alloc(size_t size)        — allocates, currently NO auto-trigger
  void  gc_mark_object(Obj* obj)
  void  gc_mark_value(Value v)
  void  gc_collect(void)             — currently MANUAL only
  void  gc_write_barrier(Obj* parent, Value new_val)  — NO-OP stub
  GCStats gc_get_stats(void)

object.h / object.c — heap objects:
  ObjHeader embedded in every heap object (for GC linked list)
  ObjArray, ObjObject (hash table), ObjString, ObjClosure
  koda_object_set / koda_array_set already call gc_write_barrier

koda_runtime.h / koda_runtime.c — runtime functions:
  koda_unbox_number(Value v) -> double
  koda_box_number(double d)  -> Value
  koda_get(Value obj, Value key) -> Value   (dispatches array/string/table)
  koda_set(Value obj, Value key, Value val)
  koda_alloc_cell() -> Value*               (heap cell for captured vars)
  koda_cell_read(Value* cell) -> Value
  koda_cell_write(Value* cell, Value val)

### Go codegen (internal/codegen/)

All Koda values in LLVM IR are i64.
Every function takes i64 params and returns i64.
koda_unbox_number / koda_box_number are declared in runtime.go and
wired as generator fields runtimeUnboxNumber / runtimeBoxNumber.
The 5 Phase 1 bugs (arithmetic, return boxing, short-circuit, index
dispatch, main inlining) have already been fixed.
Closures (flat capture model with heap cells) are the current work in
progress — scope stack + capture analysis + cell allocation.

---

## The three problems this brief covers

1. Automatic GC triggering — memory currently grows forever
2. Finding GC roots on the native stack — the GC cannot see C stack frames
3. Compile-time variable sizing — the compiler choosing stack vs heap
   per variable based on escape and type analysis

These are implemented in phases. Do not skip ahead.

---

## Problem 1 — Automatic GC triggering

### Current state (broken for real programs)
gc_collect() only runs when user code calls gc() explicitly.
A program that never calls gc() leaks memory forever.

### The fix — allocation-triggered collection with 2x growth

Step 1: Add GC state to gc.c

Replace the current GC globals with a single GCState struct:

  typedef struct {
      Obj*    objects;           // intrusive linked list — already exists
      size_t  bytes_allocated;   // current live heap bytes
      size_t  next_gc;           // trigger collection at this threshold
      GCStats stats;             // already exists
  } GCState;

  static GCState gc_state = {
      .objects         = NULL,
      .bytes_allocated = 0,
      .next_gc         = 1024 * 1024,   // first collection at 1MB
  };

Step 2: Trigger in gc_alloc

  void* gc_alloc(size_t size) {
      gc_state.bytes_allocated += size;

      if (gc_state.bytes_allocated > gc_state.next_gc) {
          gc_collect();
      }

      void* ptr = malloc(size);
      if (!ptr) {
          fprintf(stderr, "koda: out of memory\n");
          exit(1);
      }
      return ptr;
  }

Step 3: Grow threshold after collection

  void gc_collect(void) {
      size_t before = gc_state.bytes_allocated;

      gc_mark_roots();   // see Problem 2
      gc_sweep();

      // 2x growth: next GC at twice the current live set size
      // Same heuristic as Lua, CPython, and V8
      gc_state.next_gc = gc_state.bytes_allocated * 2;
      if (gc_state.next_gc < 1024 * 1024) {
          gc_state.next_gc = 1024 * 1024;  // floor at 1MB
      }

      gc_state.stats.bytes_freed    += before - gc_state.bytes_allocated;
      gc_state.stats.collections++;
  }

Why 2x:
  - Small programs: threshold grows fast, GC barely runs
  - Large programs: GC runs proportionally to live set, not total allocated
  - Pause time stays proportional to live object count, not heap size
  - This is the industry standard — do not invent a different heuristic

Step 4: Expose threshold tuning for games (add to koda_runtime.h)

  void koda_gc_set_threshold(size_t bytes);   // override next_gc
  void koda_gc_disable(void);                 // pause GC for pre-alloc phase
  void koda_gc_enable(void);                  // resume GC

These allow game code to do:
  koda_gc_disable()
  // ... allocate all game objects at scene load
  koda_gc_enable()
  // ... gameplay runs with no surprise collections

---

## Problem 2 — Finding GC roots on the native stack

### Why this is hard
When the LLVM-compiled code runs, Koda values (uint64_t) live in:
  - C stack frames (as local variables / function parameters)
  - CPU registers

The GC's gc_mark_roots() cannot see either of these.
If it misses a live value, the GC frees an object that is still in use.
That is a use-after-free bug — silent data corruption or crash.

### Roots the GC already knows about
  - Global variables array (exposed by codegen)
  - Module export cache
  - Open upvalue list (for closures)
  - The objects linked list (for sweep phase)

### Roots the GC cannot yet see
  - Local variables in LLVM-compiled functions on the C stack
  - Temporary values in CPU registers mid-computation

### Approach A — Conservative stack scanning (implement now, Phase 2)

Scan every 8-byte word on the C stack. If a word looks like a valid
heap object pointer (passes IS_OBJ and points inside our heap), treat
it as a live root.

This is correct because:
  - IS_OBJ checks the NaN-boxing tag bits — random integers rarely pass
  - gc_is_heap_object confirms the pointer is in our allocation range
  - False positives (keeping a dead object alive) are safe
  - False negatives (missing a live object) are impossible by construction

False positives cause minor memory leaks (dead objects kept alive until
the next collection). They do not cause crashes or corruption.
This is how the Boehm GC works. It is how MRI Ruby and CPython
handled this problem for years. It is good enough for Phase 2 and 3.

Implementation:

In koda_runtime.h, add:
  extern void* gc_stack_base;

In koda_runtime.c, add:
  void* gc_stack_base = NULL;

In koda_runtime_init(), capture the stack base as the very first thing:
  void koda_runtime_init(void) {
      gc_stack_base = __builtin_frame_address(0);
      gc_init();
      // ... rest of init unchanged
  }

In gc.c, add gc_mark_stack_conservative():
  void gc_mark_stack_conservative(void) {
      void* stack_top;

      // Capture current stack pointer
      // GCC/Clang on x86-64:
      __asm__ volatile ("mov %%rsp, %0" : "=r" (stack_top));

      uint64_t* ptr  = (uint64_t*)stack_top;
      uint64_t* base = (uint64_t*)gc_stack_base;

      while (ptr < base) {
          uint64_t word = *ptr;
          if (IS_OBJ(word)) {
              Obj* candidate = AS_OBJ(word);
              if (gc_is_heap_object(candidate)) {
                  gc_mark_object(candidate);
              }
          }
          ptr++;
      }
  }

Add gc_is_heap_object() to gc.c:
  // Walk the objects list to confirm a pointer is one of ours.
  // This is O(n) but only runs during GC pauses, not hot path.
  static bool gc_is_heap_object(Obj* candidate) {
      Obj* obj = gc_state.objects;
      while (obj) {
          if (obj == candidate) return true;
          obj = obj->next;
      }
      return false;
  }

Call gc_mark_stack_conservative() from gc_mark_roots():
  void gc_mark_roots(void) {
      // Stack (conservative)
      gc_mark_stack_conservative();

      // Globals
      for (int i = 0; i < koda_globals_count; i++) {
          gc_mark_value(koda_globals[i]);
      }

      // Module export cache
      koda_mark_module_cache();

      // Open upvalues (closures)
      koda_mark_open_upvalues();
  }

Platform note on capturing stack pointer:
  x86-64 Linux/macOS:  __asm__ volatile ("mov %%rsp, %0" : "=r" (sp))
  ARM64 (Apple Silicon): __asm__ volatile ("mov %0, sp" : "=r" (sp))
  Windows x86-64:      use _AddressOfReturnAddress() or
                       __asm__ { mov stack_top, rsp }
  
  For cross-platform safety, add to gc.h:
    #if defined(__x86_64__)
    #  define GC_GET_SP(sp) __asm__ volatile ("mov %%rsp, %0" : "=r" (sp))
    #elif defined(__aarch64__)
    #  define GC_GET_SP(sp) __asm__ volatile ("mov %0, sp"    : "=r" (sp))
    #else
    #  error "Unsupported architecture — add stack pointer capture for this platform"
    #endif

### Approach B — Shadow stack (implement in Phase 3, replaces A)

A shadow stack maintains a precise list of every live Koda value across
all active stack frames. Zero false positives. Zero false negatives.

Add to koda_runtime.h:
  typedef struct {
      uint64_t* slots;   // pointer to this frame's value array
      int       count;   // number of slots in this frame
  } KodaShadowFrame;

  extern KodaShadowFrame koda_shadow_stack[];
  extern int             koda_shadow_depth;

  void koda_push_frame(uint64_t* slots, int count);
  void koda_pop_frame(void);

The LLVM codegen emits at every function entry:
  // Before:  (current — no frame tracking)
  // After:   (Phase 3)
  %frame_slots = alloca [N x i64]        ; one slot per live Koda value
  call void @koda_push_frame(%frame_slots, N)

And at every function exit (before ret):
  call void @koda_pop_frame()

gc_mark_roots() replaces gc_mark_stack_conservative() with:
  void gc_mark_shadow_stack(void) {
      for (int i = 0; i < koda_shadow_depth; i++) {
          KodaShadowFrame* frame = &koda_shadow_stack[i];
          for (int j = 0; j < frame->count; j++) {
              gc_mark_value(frame->slots[j]);
          }
      }
  }

Do not implement this until Phase 3. Conservative scanning is correct
and requires no codegen changes. Shadow stack requires the codegen to
emit push/pop at every function boundary — do that after closures work.

### Approach C — LLVM GC intrinsics (do not implement)

LLVM has gc.statepoint / gc.relocate intrinsics for precise relocating
GC. This enables a compacting collector (objects can be moved in memory,
eliminating fragmentation). It requires stack maps, safepoints at every
allocation, and significant runtime support. This is what Java HotSpot
and the .NET CLR do. It is not needed for Koda's use cases. Skip it.

---

## Problem 3 — Compile-time variable sizing

### The three storage classes

Every variable the compiler emits falls into exactly one class:

  Class 1 — Stack scalar
    Condition: provably not captured, provably not escaping function scope
    LLVM IR:   %slot = alloca i64
               store i64 %val, %slot
               %x = load i64, %slot        (→ register after mem2reg)
    GC cost:   zero — not a heap object, not a root

  Class 2 — Heap cell
    Condition: captured by a closure OR escapes function scope
    LLVM IR:   %cell = call i64* @koda_alloc_cell()
               call void @koda_cell_write(%cell, %val)
               %x = call i64 @koda_cell_read(%cell)
    GC cost:   one allocation, write barrier on writes

  Class 3 — Heap object
    Condition: array literal, object literal, string, closure struct
    LLVM IR:   %arr = call i64 @koda_allocate_array(i32 %n)
    GC cost:   allocation proportional to size, full GC management

The compiler currently emits Class 1 (alloca) for all locals and Class 3
for arrays/objects. Class 2 (heap cells for captured vars) is the closure
work in progress. The goal here is making the Class 1 / Class 2 decision
automatic based on escape analysis.

### Escape analysis — implement in sema, Phase 2

A variable escapes if it outlives the stack frame that declared it.

Rules in order of implementation difficulty:

  Rule 1 (easy): A variable captured by a closure escapes.
    func f(x) { return func() { return x } }
    x escapes — it is referenced after f() returns.
    This is already identified by closure capture analysis.
    Use this result directly: EscapingDecls[d] = true for all captured vars.

  Rule 2 (easy): A variable directly returned escapes.
    func f() { let x = makeObj(); return x }
    x escapes — the caller holds it after f() returns.
    Detection: ReturnStmt whose value resolves to a local LetDecl.

  Rule 3 (medium): A variable stored into an escaping container escapes.
    func f() { let arr = []; arr.push(x); return arr }
    x escapes if arr escapes (Rule 2 applies to arr).
    Detection: assignment to field/index of an escaping variable.

  Rule 4 (hard — skip for now): interprocedural escape.
    func store(container, val) { container.data = val }
    func f() { let x = 5; store(global_obj, x); }
    x escapes through store() into global_obj.
    This requires tracking what called functions do with their arguments.
    Skip for Phase 2. Be conservative: if passed to an unknown function,
    assume the variable escapes.

Add to internal/sema/stub.go NativeEmitContext:

  type NativeEmitContext struct {
      Bundle       *parser.ProgramBundle
      locals       map[parser.Expr]parser.Decl
      capturedVars map[parser.Decl]bool
      paramDecls   map[*parser.FuncDecl][]parser.Decl
      FuncCaptures map[*parser.FuncDecl][]*parser.LetDecl
      ExprCaptures map[*parser.FuncExpr][]*parser.LetDecl

      // Add these for escape analysis:
      EscapingDecls map[*parser.LetDecl]bool  // must be heap cell
      StackDecls    map[*parser.LetDecl]bool  // safe as stack alloca
  }

Populate them in PrepareNativeBundle() after capture analysis:

  func analyzeEscape(ctx *NativeEmitContext, bundle *parser.ProgramBundle) {
      // Start with all captured vars as escaping (from closure analysis)
      for decl := range ctx.capturedVars {
          if ld, ok := decl.(*parser.LetDecl); ok {
              ctx.EscapingDecls[ld] = true
          }
      }

      // Add Rule 2: returned locals
      walkReturns(bundle, func(ret *parser.ReturnStmt) {
          if ident, ok := ret.Value.(*parser.IdentifierExpr); ok {
              if decl, ok := ctx.locals[ident]; ok {
                  if ld, ok := decl.(*parser.LetDecl); ok {
                      ctx.EscapingDecls[ld] = true
                  }
              }
          }
      })

      // Everything not in EscapingDecls is a stack scalar
      walkLetDecls(bundle, func(ld *parser.LetDecl) {
          if !ctx.EscapingDecls[ld] {
              ctx.StackDecls[ld] = true
          }
      })
  }

### Storage class decision in codegen — Phase 2

Update emitLetDecl in internal/codegen/codegen.go:

  func (g *Generator) emitLetDecl(d *parser.LetDecl) error {
      name := d.Name.Lexeme

      if d.Native != nil {
          return g.emitNativeExternLet(d)
      }

      var storageSlot value.Value

      switch {

      case g.ctx.StackDecls[d]:
          // Class 1: stack alloca — zero GC cost
          // LLVM's mem2reg will promote this to a register
          slot := g.block.NewAlloca(types.I64)
          storageSlot = slot
          g.defineVar(name, slot)

      case g.ctx.EscapingDecls[d]:
          // Class 2: heap cell — survives stack frame
          cell := g.block.NewCall(g.runtimeAllocCell)
          storageSlot = cell
          g.defineVar(name, cell)
          g.capturedCells[d] = cell

      default:
          // Conservative fallback: treat as escaping
          // Better to allocate on heap than to corrupt memory
          cell := g.block.NewCall(g.runtimeAllocCell)
          storageSlot = cell
          g.defineVar(name, cell)
      }

      if d.Init != nil {
          initVal, err := g.emitExpr(d.Init)
          if err != nil {
              return err
          }
          boxed := g.emitAsKodaI64(initVal)

          if _, isCell := g.capturedCells[d]; isCell {
              g.block.NewCall(g.runtimeCellWrite, storageSlot, boxed)
          } else {
              g.block.NewStore(boxed, storageSlot)
          }
      }

      return nil
  }

### Type inference feeding unboxed storage — Phase 3

Once type inference is implemented (Phase 3), add a fourth storage class:

  Class 1b — Unboxed scalar
    Condition: stack scalar AND compile-time proven to always be float64
    LLVM IR:   %slot = alloca double    (not i64 — raw double)
               store double %val, %slot
    GC cost:   zero
    Arithmetic cost: zero — no koda_unbox_number / koda_box_number calls
                     LLVM emits fadd/fmul directly on the double

  Example after type inference:
    func distance(x1, y1, x2, y2) {
        let dx = x2 - x1   // TypeInfo: Float64, StackDecl → alloca double
        let dy = y2 - y1   // same
        return sqrt(dx*dx + dy*dy)
    }
    After mem2reg: dx and dy are SSA double values in registers.
    The entire function is ~6 floating-point instructions. No memory access.
    Performance: identical to hand-written C.

Add to NativeEmitContext in Phase 3:
  TypeInfo map[*parser.LetDecl]InferredType

  type InferredType int
  const (
      TypeUnknown InferredType = iota
      TypeFloat64
      TypeBool
      TypeString
      TypeArray
      TypeObject
      TypeFunction
  )

Codegen checks in Phase 3:
  if g.ctx.TypeInfo[d] == TypeFloat64 && g.ctx.StackDecls[d] {
      slot := g.block.NewAlloca(types.Double)
      // emit fadd/fmul directly, skip koda_unbox/box
  }

---

## Generational GC — Phase 3, implement after shadow stack

Do not implement this until the shadow stack is working. Generational GC
requires precise root information — conservative scanning produces false
positives that can incorrectly classify nursery objects as reachable from
the old generation.

### Why generational GC

The current mark-sweep pauses the entire program for every collection.
At 60fps, any GC pause longer than ~2ms drops a frame. A full mark-sweep
over 50,000 live objects easily takes 5-20ms.

The generational hypothesis: most objects die young. A temporary array
in a loop iteration is dead by the next iteration. Scanning it every
frame wastes time. Only scan it once — in the nursery — and promote
survivors to an older generation that is scanned less frequently.

### Memory layout

  Nursery   (Gen 0): 256KB fixed buffer, bump allocator
                     Collected every ~0.5-1ms
                     All new allocations land here
  
  Young gen (Gen 1): 2MB, free-list allocator
                     Survivors from nursery promoted here
                     Collected every ~10ms
  
  Old gen   (Gen 2): Unlimited, mark-sweep
                     Long-lived objects (player, level data, etc.)
                     Collected every ~100ms or when explicitly requested

### Nursery bump allocator

  static uint8_t  g_nursery[256 * 1024];
  static uint8_t* g_nursery_top = g_nursery;
  static uint8_t* g_nursery_end = g_nursery + sizeof(g_nursery);

  void* nursery_alloc(size_t size) {
      size = (size + 7) & ~7;               // align to 8 bytes

      if (g_nursery_top + size > g_nursery_end) {
          gc_collect_minor();               // nursery full — minor GC
          // After minor GC: survivors promoted, bump pointer reset
          if (g_nursery_top + size > g_nursery_end) {
              // Survived minor GC but still no room — promote to old gen
              return old_gen_alloc(size);
          }
      }

      void* ptr = g_nursery_top;
      g_nursery_top += size;
      return ptr;
  }

  // Reset nursery after minor collection
  void nursery_reset(void) {
      g_nursery_top = g_nursery;
  }

This is 2-3 CPU instructions per allocation. malloc() is 50-200
instructions. The nursery bump allocator is the single biggest
allocation performance win available.

### Write barrier implementation (fill in the stub)

The gc_write_barrier stub is already called from koda_object_set and
koda_array_set. In Phase 3, fill it in:

First, add generation tracking to ObjHeader in object.h:
  typedef struct Obj {
      ObjType  type;
      bool     is_marked;
      uint8_t  generation;    // ADD: 0=nursery, 1=young, 2=old
      struct Obj* next;
  } Obj;

  #define GEN_NURSERY 0
  #define GEN_YOUNG   1
  #define GEN_OLD     2

Add remembered set to gc.c:
  #define REMEMBERED_SET_MAX 4096
  static Obj* remembered_set[REMEMBERED_SET_MAX];
  static int  remembered_set_count = 0;

  static void remembered_set_add(Obj* obj) {
      // Avoid duplicates
      for (int i = 0; i < remembered_set_count; i++) {
          if (remembered_set[i] == obj) return;
      }
      if (remembered_set_count < REMEMBERED_SET_MAX) {
          remembered_set[remembered_set_count++] = obj;
      }
      // If full: trigger a full GC to clear it
  }

Fill in the write barrier:
  void gc_write_barrier(Obj* parent, Value new_val) {
      if (!IS_OBJ(new_val)) return;

      Obj* child = AS_OBJ(new_val);

      // Only matters when old object points to young/nursery object
      if (parent->generation == GEN_OLD &&
          child->generation  != GEN_OLD) {
          remembered_set_add(parent);
      }
  }

During minor GC, the remembered set provides additional roots — old-gen
objects that reference nursery objects keep those nursery objects alive.

### Minor collection (nursery GC)

  void gc_collect_minor(void) {
      // Roots: shadow stack + globals + remembered set
      gc_mark_shadow_stack();
      gc_mark_globals();
      gc_mark_remembered_set();   // old→young cross-gen pointers

      // Promote survivors from nursery to young gen
      // Any nursery object still marked after this pass is alive
      // Move it to young gen, update all pointers to it
      gc_promote_nursery_survivors();

      // Reset nursery — all unmarked nursery objects are dead
      nursery_reset();

      // Clear remembered set (rebuild on next write barrier hits)
      remembered_set_count = 0;
  }

### Game loop integration

For games, the GC should be called explicitly each frame with a time budget:

  // Add to koda_runtime.h:
  void koda_gc_frame_step(double budget_ms);

  // Implementation in gc.c:
  void koda_gc_frame_step(double budget_ms) {
      uint64_t start_us = koda_timestamp_us();

      // Always do a minor collection (fast — <0.5ms)
      if (nursery_is_full_enough()) {
          gc_collect_minor();
      }

      // Do incremental old-gen work within time budget
      double elapsed_ms = (koda_timestamp_us() - start_us) / 1000.0;
      if (elapsed_ms < budget_ms) {
          gc_incremental_step(budget_ms - elapsed_ms);
      }
  }

  // Game calls this each frame:
  // koda_gc_frame_step(0.5)  // give GC 0.5ms per frame = 30ms/sec GC budget

---

## Object shapes for inline property access — Phase 4

Once objects have a predictable field layout, property access can bypass
the hash table entirely.

Add shape system to object.h:
  typedef struct {
      uint32_t  shape_id;
      uint32_t  field_count;
      char**    field_names;   // ordered list of field names
  } ObjShape;

  // Modified ObjObject: shape ID + dense slot array instead of hash table
  typedef struct {
      ObjHeader header;
      uint32_t  shape_id;
      uint32_t  slot_count;
      Value     slots[];       // field values indexed by position, not name
  } ObjShapedObject;

Shape registry in koda_runtime.c:
  static ObjShape* shape_registry[4096];
  static uint32_t  shape_count = 0;

  uint32_t koda_get_or_create_shape(char** field_names, int count) {
      // Look up existing shape with same fields in same order
      for (uint32_t i = 0; i < shape_count; i++) {
          if (shapes_equal(shape_registry[i], field_names, count)) {
              return i;
          }
      }
      // Create new shape
      ObjShape* s = malloc(sizeof(ObjShape));
      s->shape_id    = shape_count;
      s->field_count = count;
      s->field_names = copy_field_names(field_names, count);
      shape_registry[shape_count] = s;
      return shape_count++;
  }

When the compiler sees an object literal with known fields, it emits
a shape lookup at compile time and a direct slot load at runtime:

  ; property access: player.x
  ; compiler determined shape_id=7 means slot 0 = "x"
  %shape_ptr = getelementptr inbounds %ObjShapedObject, %player, 0, 1
  %shape_id  = load i32, i32* %shape_ptr
  %is_right_shape = icmp eq i32 %shape_id, 7
  br i1 %is_right_shape, label %fast, label %slow

  fast:
  %slots_ptr = getelementptr inbounds %ObjShapedObject, %player, 0, 3
  %x_ptr     = getelementptr inbounds i64, i64* %slots_ptr, i64 0
  %x_val     = load i64, i64* %x_ptr
  br label %merge

  slow:
  %x_val_slow = call i64 @koda_object_get(%player_val, %x_key)
  br label %merge

  merge:
  %x = phi i64 [ %x_val, %fast ], [ %x_val_slow, %slow ]

This is 10-50x faster than the hash table path for hot object access
(game entities updated every frame). The slow path ensures correctness
when shape assumptions are wrong.

---

## Implementation sequence — strict order, no skipping

### Now — Phase 2 (alongside closure work)

  [ ] 1. Add GCState struct to gc.c, migrate existing fields into it
  [ ] 2. Add bytes_allocated tracking and next_gc threshold to gc_alloc
  [ ] 3. Add 2x growth heuristic to gc_collect
  [ ] 4. Add koda_gc_disable / koda_gc_enable / koda_gc_set_threshold to
         koda_runtime.h and koda_runtime.c
  [ ] 5. Save gc_stack_base in koda_runtime_init using
         __builtin_frame_address(0)
  [ ] 6. Add GC_GET_SP macro to gc.h for cross-platform stack pointer capture
  [ ] 7. Implement gc_mark_stack_conservative in gc.c
  [ ] 8. Add gc_is_heap_object helper to gc.c
  [ ] 9. Wire gc_mark_stack_conservative into gc_mark_roots
  [ ] 10. Add escape analysis (Rules 1 and 2) to sema/stub.go
  [ ] 11. Add EscapingDecls / StackDecls to NativeEmitContext
  [ ] 12. Update emitLetDecl in codegen to use storage class decision

  Verify:
    koda build tests/gc_test.koda && ./gc_test
    koda build tests/closure_test.koda && ./closure_test
    koda build tests/phase1_surface.koda && ./phase1_surface
    // Memory should not grow unboundedly on a long-running test
    // Closures should capture variables correctly
    // phase1_surface should pass completely

### Phase 3 — Performance

  [ ] 1. Implement shadow stack in koda_runtime.c
  [ ] 2. Emit koda_push_frame / koda_pop_frame in codegen at fn boundaries
  [ ] 3. Replace gc_mark_stack_conservative with gc_mark_shadow_stack
  [ ] 4. Add generation field to ObjHeader
  [ ] 5. Implement nursery bump allocator (nursery_alloc)
  [ ] 6. Implement gc_collect_minor with remembered set
  [ ] 7. Fill in gc_write_barrier with generation check
  [ ] 8. Implement gc_incremental_step for old gen
  [ ] 9. Add koda_gc_frame_step for game loop integration
  [ ] 10. Add type inference pass to sema
  [ ] 11. Add TypeInfo to NativeEmitContext
  [ ] 12. Update emitLetDecl to use alloca double for proven Float64 scalars
  [ ] 13. Switch from llir/llvm to tinygo-org/go-llvm
  [ ] 14. Run mem2reg + instcombine + GVN passes after IR generation

  Verify:
    go test ./...
    koda build tests/gc_test.koda && valgrind --leak-check=full ./gc_test
    koda build demos/demo_3d.koda && ./demo_3d   // no frame hitches at 60fps
    Benchmark: fibonacci(35) should run within 2x of equivalent C

### Phase 4 — Optimisation

  [ ] 1. Implement ObjShape and shape registry
  [ ] 2. Assign shapes to object literals at compile time
  [ ] 3. Emit shape-guarded direct slot access for known-shape objects
  [ ] 4. Implement scalar replacement for small non-escaping objects
  [ ] 5. Add DWARF debug info via DIBuilder (go-llvm)

---

## Invariants — never violate these

1. gc_alloc is the only allocation entry point. Never call malloc() directly
   for GC-managed objects. Everything goes through gc_alloc so bytes_allocated
   stays accurate.

2. Every GC-managed object has ObjHeader as its first field. The GC walks
   the objects linked list by casting Obj* — this requires the header to
   be at offset 0.

3. gc_write_barrier is called on every field write to a heap object. It is
   already wired into koda_object_set and koda_array_set. When you add new
   heap object types, add the write barrier call to their setter functions.

4. The GC must never run while an allocation is partially initialised.
   Allocate first, then call gc_alloc, then set fields. If gc_alloc
   triggers a collection before fields are set, the half-initialised
   object must be safe to scan (nil/zero fields are safe — the GC
   treats zero as a non-pointer).

5. koda_gc_disable / koda_gc_enable must be balanced. Disable for
   pre-allocation phases only. Never leave GC disabled during normal
   program execution.

6. The stack base must be captured before any Koda values are allocated.
   koda_runtime_init captures it as its first action. Do not move this.

---

## Verification commands

After every implementation step, run all of these before committing:

  make -C runtime                              # C runtime builds clean
  go build ./cmd/koda/...                      # compiler builds
  go test ./...                                # all Go tests pass
  koda build tests/hello.koda    && ./hello    # basic execution works
  koda build tests/gc_test.koda  && ./gc_test  # GC collects correctly
  koda build tests/closure_test.koda && ./closure_test  # closures work
  koda build tests/phase1_surface.koda && ./phase1_surface  # full surface
  koda build demos/demo_3d.koda  && ./demo_3d  # 3D demo runs

For memory correctness (install valgrind on Linux or use AddressSanitizer):
  clang -fsanitize=address -o runtime/src/koda_runtime.o ...
  koda build tests/gc_test.koda && ./gc_test
  // AddressSanitizer will catch use-after-free from GC bugs immediately