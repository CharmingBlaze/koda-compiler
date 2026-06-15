# Koda roadmap

Prioritized engineering queue for the Koda compiler and runtime. This replaces the historical week-by-week bootstrap plan with **current** ordering based on product fit, beginner impact, and runtime risk.

> **Execution model:** one pipeline -- LLVM IR (`internal/codegen`) + C runtime (`runtime/src`). No bytecode VM.  
> **Positioning context:** [positioning.md](positioning.md) · **Status today:** [status.md](status.md) · **Detailed matrix:** [MASTER_PLAN.md](MASTER_PLAN.md)

---

## How to read this document

| Tier | Meaning |
|------|---------|
| **Tier 1 -- Now** | Highest impact on the "C alternative for beginners" claim; do these first |
| **Tier 2 -- Next** | Strong beginner ROI or runtime safety; follows tier 1 |
| **Tier 3 -- Soon** | Polish, docs, and smaller runtime gaps |
| **Long-term** | Strategic; not blockers for game-maker audience |

Each item lists **MASTER_PLAN** refs where they exist. Update this file when priorities shift.

---

## Tier 1 -- Now

### 1.1 Opt-in integer types

**Problem:** all runtime numbers are `f64`. Binary data, PNG pixels, network packets, and bitfields need real `i32` / `u64` / `u8` semantics.

**Scope (minimal viable):**

- Syntax: `let n: i32 = 0;` or `let byte: u8 = 255;` (exact spelling TBD in sema)
- LLVM lowering to native integer ops -- not NaN-boxed `Value` for typed locals
- Boundary conversions at calls to untyped / `Value` APIs
- Document float-default behavior for beginners who never annotate

**Success:** test that reads raw bytes or packs a 32-bit field without float rounding.

**MASTER_PLAN ref:** new (extends **A2** type story)

---

### 1.2 Struct methods

**Problem:** beginners expect `rect.area()`, not `area(rect)`. Object-literal methods exist; **named struct types** do not.

**Scope:**

- Parser: `func` declarations inside `struct { ... }`
- Sema: bind methods to struct layout; `this` / receiver lowering
- Codegen: emit as functions with implicit first argument or vtable-free slot dispatch (match existing object-method path where possible)

**Success:** `tests/struct_method_test.koda` -- `let r = Rect { w: 3, h: 4 }; assert(r.area() == 12);`

**MASTER_PLAN ref:** new (related **A2**)

---

## Tier 2 -- Next

### 2.1 Unused symbol warnings (`--warn-unused`)

**Problem:** typos in `let` bindings fail silently at runtime.

**Scope:**

- Sema pass: unused locals, params, top-level `func` declarations
- CLI: `koda build --warn-unused`, `koda check --warn-unused` (warnings, not errors by default)

**Success:** misspelled variable in test program emits warning with line number.

**MASTER_PLAN ref:** **E4**, **E5**

---

### 2.2 Enum switch exhaustiveness warnings

**Problem:** new enum variant + incomplete `switch` → no compiler signal.

**Scope:**

- When `switch` subject is enum-typed, warn on missing cases
- Allow `default:` to silence warning
- Warning only -- keep beginner-friendly

**Success:** `tests/enum_switch_warn.koda` triggers diagnostic; complete switch does not.

**MASTER_PLAN ref:** **A3** (complete partial work)

---

### 2.3 ASAN / Valgrind CI job

**Problem:** runtime complexity (nursery, arena, remembered set, heap cache, inline arrays) needs automated memory bug detection.

**Scope:**

- Linux CI job: build runtime + tests with `-fsanitize=address`
- Run `tests/arena_test.koda`, `incremental_gc_test.koda`, `gc_soak.koda`, native conformance under ASAN
- Document local repro: `CFLAGS=-fsanitize=address make -C runtime`

**Success:** CI green on main; intentional leak test fails the job.

**MASTER_PLAN ref:** **H2** -- **priority bumped** after runtime hardening merge

---

### 2.4 Incremental GC game-loop validation

**Problem:** `gcFrameStep` exists; budget guidance for real games is thin.

**Scope:**

- Expand `tests/incremental_gc_test.koda` or add `tests/stress/` game-loop scenario
- Document recommended budgets in [game-dev.md](guides/game-dev.md) (0.5–1.0 ms/frame)
- Validate pauses under `--no-opt` in CI where time-bounded

**MASTER_PLAN ref:** **B1** (partial → done)

---

## Tier 3 -- Soon

### 3.1 String intern table policy

**Problem:** lookup is O(1) after hash-map fix; **retention** until sweep can grow memory for unique-string-heavy programs.

**Scope:**

- Document behavior in [runtime-and-gc.md](concepts/runtime-and-gc.md)
- Optional `koda_intern_clear()` builtin for explicit flush
- Consider cap + LRU only if real programs hit pain

**Success:** docs clear; optional API tested.

---

### 3.2 `status.md` / diagnostic consistency audit

- Compound assign (`obj.x += 1`) codegen verification
- Declare-before-use always blocked in sema
- Align user-facing docs with struct-method distinction

---

### 3.3 ObjTable open addressing completion

**MASTER_PLAN ref:** **B2** -- confirm probing, load factor, and `tests/table_hash_test.koda` coverage.

---

### 3.4 Stress suite expansion

**MASTER_PLAN ref:** **B6** -- grow `tests/stress/`; Linux CI timeout job.

---

## Long-term

| Item | Notes | MASTER_PLAN |
|------|-------|-------------|
| **Optional / nullable types** | `let name: string?`; warn on unchecked deref | new |
| **`try` / err propagation** | Reduce `ok`/`err` boilerplate | new |
| **Interfaces / traits** | Structural checks on method sets | new |
| **Package manager** | `koda install` → `packages/` | new |
| **WASM target** | `wasm32-unknown-unknown`, slim runtime | new |
| **Full DWARF / source DI** | `.koda` line numbers in gdb/backtraces | **D1** |
| **`koda bench`** | Compile + runtime benchmarks in CI | **D3** |
| **Parallel module parse** | Faster large projects | **E2** |
| **Runtime property hints** | "did you mean?" for missing keys | **D4** |
| **Stdlib expansion** | `color`, `input`, `easing`, `pool`, ... | **F1–F7** |
| **Bundle asset embedding** | Ship assets in `koda bundle` | **G4** |

---

## Completed foundations (do not regress)

Treat these as release blockers if they break. Detail in [MASTER_PLAN.md](MASTER_PLAN.md) § Do not regress.

**Pipeline:** lexer → parser → sema → LLVM → clang → binary; `PrepareNativeBundle` before codegen.

**Runtime / GC (2025 review merged):**

- Tri-generational GC, shadow stack, write barriers
- O(1) remembered set; conservative root heap cache
- Arena builtins; inline small arrays; `gcFrameStep`
- Tombstone sentinel; intern hash table; table `hashes[]` GC accounting
- Loop terminator guards in while/do-while codegen

**Language:** structs, enums, closures, defer, `ok`/`err`, typo hints, `for-of`, range `lo..hi`.

**Tooling:** `koda check`, `koda fmt`, `koda watch`, `--no-opt`, `--debug` (`-g`), CI smokes on Ubuntu/macOS/Windows.

---

## Testing strategy

### Every merge

```bash
go test ./... -count=1
make runtime-lib   # or equivalent compile of runtime/src
```

### Native smokes

```bash
koda run tests/native_conformance.koda
koda run tests/incremental_gc_test.koda
koda run tests/arena_test.koda
```

### Before release

- Native conformance build + run
- Raylib brick breaker (graphical gate)
- `tests/gc_soak.koda` / `tests/gc_pressure_expr.koda` under load
- ASAN job green (once **2.3** lands)

### Benchmarks (informational)

- Fibonacci (recursion / shadow stack)
- Binary trees (GC stress)
- Per-frame arena + `gcFrameStep` game loop

---

## Success criteria by audience

| Audience | Ready when |
|----------|------------|
| **Game-making beginners** | Tier 1 optional; tier 2 warnings nice-to-have; GC/arena docs clear |
| **General app beginners** | Tier 1.2 struct methods + tier 2.1 unused warnings |
| **C / systems learners** | Tier 1.1 integers + tier 1.2 methods + ASAN CI |
| **Contributors** | `go test` green, MASTER_PLAN matrix updated on merge, this file reflects priorities |

---

## Historical note

An earlier version of this file contained a **week-by-week bootstrap plan** (codegen completion → runtime → linking → v1.0.0). That bootstrap phase is largely complete. This roadmap supersedes that schedule for ongoing work. Archive context lives in git history (`ROADMAP.md` pre-2025).

---

## Related

- [positioning.md](positioning.md) -- honest product framing
- [status.md](status.md) -- what works today
- [MASTER_PLAN.md](MASTER_PLAN.md) -- full engineering matrix
- [handoff.md](handoff.md) -- pipeline for new contributors
