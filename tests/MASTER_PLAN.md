# Koda — master plan

This document is for **people working in the Koda repository**. It records what is done, what is in progress, and what comes next — with enough detail to implement each item without a separate design document.

For product positioning and the rationale behind the priority order, see **[docs/positioning.md](positioning.md)**.  
For the priority-ordered next steps in brief, see **[docs/ROADMAP.md](ROADMAP.md)**.

---

## Engineering standards

- Correctness first: memory safety, GC invariants, and clear diagnostics beat feature count.
- Read existing code before changing it; match local style; do not leave known holes.
- Prefer complete implementations over TODOs and stubs.
- **`go test ./...`** green before merging substantive compiler or runtime changes.
- C runtime: keep **`runtime/src/*.c`** and **`internal/codegen/runtime.go`** declarations aligned.

---

## Shipped baseline — do not regress

These properties are expected in the tree today. Treat regressions as release blockers.

**Pipeline**
- Full `lexer → parser → sema → LLVM IR → clang → binary` pipeline, cross-platform (Linux, macOS, Windows).
- `sema.PrepareNativeBundle` runs before codegen — bad programs never reach LLVM.
- `diagnostic.MultiError` aggregates sema issues with file/line snippets.
- `koda check` — full sema, no LLVM; `koda fmt` / `koda fmt --check` / `./...` expansion.
- `koda run`, `koda build`, `koda watch`, `koda bundle`, `koda build --debug`.

**Language surface**
- Variables, primitives (number/string/bool/null), arrays, objects, closures, upvalues.
- `struct` types with O(1) field access (`koda_struct_get` / `koda_struct_set`).
- `enum` types with ordinal constants and constant folding.
- `defer` (LIFO; return value computed before defers run).
- `for (let i of lo..hi)` dynamic range — counted `i64` loop, unboxed bounds.
- Object method shorthand + `this` on **object literals** (struct methods are a separate open item — see A7).
- Integer literal folding (`+`, `-`, `*`, unary `-`/`+`, nested safe cases).
- Call arity checks for `func` / `// koda:extern` / known argv-style methods.
- "Did you mean?" hints on undefined identifiers (Levenshtein edit distance).
- `ok` / `err`, `panic`, `assert`, `readFile`, `writeFile`, `parseJSON`.

**GC and runtime**
- Generational GC: nursery bump allocator (256 KB), young, old generations.
- Incremental major collection via `gcFrameStep(ms)`; stop-the-world `gc_collect()` for shutdown.
- Shadow stack for precise root scanning; conservative stack scan as fallback.
- Write barriers on all mutating stores into GC objects (audited by `write_barrier_test.go`).
- Remembered-set dedup O(1) — `in_remembered_set` flag on `Obj`.
- Conservative scan O(heap) per collection — heap pointer hash set + sorted range table (binary search interior pointers).
- Inline small arrays: capacity ≤ 64 allocates header + elements in one block.
- String intern table O(1) — open-addressing hash map with FNV-1a; sweep rebuilds from marked survivors.
- Tombstone sentinel `TOMBSTONE_VAL` (tag 0x05) — boolean hash-table keys are safe.
- Table `hashes` array via `gc_alloc`/`gc_free` — counted in GC threshold.
- Arena builtins: `arena(size)`, `arenaReset`, `arenaAllocArray`, `arenaAllocStruct`.
- `gcFrameStep(ms)`, `gcStats()`, `gcDisable()`, `gcEnable()`, `gcCollect()`.
- Loop terminator guards in `emitWhileStmt` / `emitDoWhileStmt` — no invalid LLVM IR from early exits.

**Numeric type inference (shipped)**
- `internal/sema/typeinfer.go` — `InferNumericKinds` pass classifies every `LetDecl` as `KindInt` or `KindFloat` at compile time. No syntax changes required.
- `internal/codegen/intfast.go` — when both operands of `+`, `-`, `*`, `%`, `&`, `|`, `^`, `<<`, `>>` are `KindInt` stack locals, emits raw LLVM `add`/`sub`/`mul`/`srem`/`and`/`or`/`xor`/`shl`/`ashr` — zero `unboxNumber`/`boxNumber` round-trips.
- Rules: integer literals → `KindInt`; `f64` with no fractional part → `KindInt`; division always `KindFloat`; bitwise ops always `KindInt`; escaped/captured variables always `KindFloat` (conservative).
- 8 unit tests in `internal/sema/typeinfer_test.go` covering literals, whole floats, fractions, int+int arithmetic, division, bitwise, and escaping decls.

**Tooling**
- Raylib wrapper path; `koda wrap` / `kodawrap` for C library binding generation.
- CI: `go test ./...`, fmt check, native smoke on Ubuntu / macOS / Windows, GC soak.
- `koda_runtime_init_ex` / `koda_runtime_shutdown` from generated `main`.
- Shadow stack growth + hard cap panic message; `koda_shadow_stack_high_water()`.
- `ObjArray` / `ObjTable` NULL guards on partial allocation paths.
- `KODA_GC_DEBUG` diagnostics.
- `stdlib/math.koda`, `stdlib/vec2.koda`, `stdlib/vec3.koda`, `stdlib/timer.koda`.
- Release `-tags release` embed story.

---

## Progress matrix

Legend: **Done** · **Partial** (exists but gaps remain) · **Open** (not started).

| Ref | Topic | Status | Notes |
|-----|-------|--------|-------|
| **A1** | "Did you mean?" for undefined names | **Done** | `internal/sema/levenshtein.go`; `tests/typo_suggestion.koda` |
| **A2** | `struct` types + literals + field checks | **Done** | O(1) slot access; `tests/struct_typo_test.koda` |
| **A3** | `enum` + switch exhaustiveness warning | **Partial** | AST + lowering done; exhaustiveness warning not yet emitted — see P4 |
| **A4** | `math` builtins (`lerp`, `clamp`, `hypot`, full suite) | **Done** | C runtime + `stdlib/math.koda` |
| **A5** | `stdlib/timer.koda` | **Done** | `tests/timer_test.koda` |
| **A6** | Opt-in integer types (`i32`, `u8`, …) | **Open** | Type annotation syntax not yet parsed; `KindInt` inference covers pure-local arithmetic — see P1 |
| **A7** | Struct methods (`box.area()`) | **Open** | Object-literal `this` works; struct body methods not yet — see P2 |
| **A8** | Numeric type inference (`KindInt` / `KindFloat`) | **Done** | `internal/sema/typeinfer.go` + `internal/codegen/intfast.go`; 8 unit tests pass |
| **B1** | Incremental major GC | **Done** | `gc_collect_incremental`; `gcFrameStep`; `tests/incremental_gc_test.koda` in CI |
| **B2** | ObjTable open addressing | **Done** | `hashes[]` open-addressing; `gc_alloc` for hash array; `tests/table_hash_test.koda` |
| **B3** | GC pressure relief builtins | **Done** | `gcDisable`, `gcEnable`, `gcCollect`; `tests/gc_control_test.koda` |
| **B4** | `gcStats()` builtin | **Done** | Lowercase keys; `tests/gc_stats_frame_test.koda` |
| **B5** | Write-barrier static audit | **Done** | `internal/codegen/write_barrier_test.go` |
| **B6** | `tests/stress/` suite + CI job | **Partial** | `tests/stress/stress_mixed_alloc.koda` + Linux CI (`timeout 90s`); expand over time |
| **B7** | Remembered-set O(1) dedup | **Done** | `in_remembered_set` on `Obj`; `remembered_set_clear()` on all paths |
| **B8** | Conservative scan O(heap) per collection | **Done** | Heap pointer hash set + sorted range table; `qsort` + binary search |
| **B9** | Inline small arrays | **Done** | Capacity ≤ 64: one `gc_alloc`; `inline_elements` flag; grow copies without freeing inline storage |
| **B10** | Arena builtins | **Done** | `arena`, `arenaReset`, `arenaAllocArray`, `arenaAllocStruct`; `tests/arena_test.koda` |
| **B11** | String intern table O(1) | **Done** | Open-addressing hash map; sweep rebuilds from marked survivors |
| **B12** | Intern table retention docs + `koda_intern_clear()` | **Open** | Retention model undocumented; builtin not yet exposed — see P6 |
| **B13** | Tombstone sentinel | **Done** | `TOMBSTONE_VAL` tag 0x05; `IS_TOMBSTONE()`; boolean keys safe |
| **C1** | Typed runtime errors | **Partial** | `koda_type_error`, `koda_null_error`; argv/native path audit remains |
| **C2** | OOB array/string access → panic | **Done** | `tests/array_oob_*`; `api` tests |
| **C3** | Capacity overflow guards | **Done** | `validate_value_slot_count` in `object.c` |
| **C4** | Shadow push/pop balance (debug mode) | **Done** | Shutdown check under `KODA_GC_DEBUG` when depth ≠ 0 |
| **D1** | LLVM debug source locations (DI) | **Partial** | `DICompileUnit`/`DISubprogram`/`DILocation` on `--debug` builds; expand coverage |
| **D2** | `koda watch` | **Done** | Poll-based; `fsnotify` upgrade optional |
| **D3** | `koda bench` | **Done** | Warmup, avg, p50/p95/p99 reporting |
| **D4** | Runtime property "did you mean?" | **Done** | `KODA_GC_DEBUG` table key hints in `koda_object_get` |
| **E1** | Parse cache (`mtime`) | **Done** | `internal/parser/loader.go`; overlays bypass cache |
| **E2** | Parallel module parse | **Done** | Sibling `#include`/`import` paths parsed concurrently |
| **E3** | Wider constant folding | **Partial** | Extend beyond current literal + integer rules |
| **E4** | Dead-code elimination + `--warn-unused` | **Done** | `--warn-unused` + unreachable warnings in `strict` lint |
| **E5** | Unused binding warnings | **Done** | See P3 |
| **F1** | `stdlib/vec3.koda` | **Done** | `tests/vec3_test.koda` |
| **F2–F7** | `color`, `input`, `easing`, `pool`, `str`, `array` stdlib modules | **Open** | Pure Koda where possible; each needs a `tests/*.koda` file |
| **G4** | Asset embedding in `koda bundle` | **Done** | `assetPath()` builtin + `koda-assets.txt` manifest in bundles |
| **G5** | `koda doctor` depth | **Done** | OK/FAIL report, smoke build, runtime freshness, disk space (Unix) |
| **H1** | Parser fuzzing | **Partial** | `internal/parser/FuzzParse`; Linux CI smoke `-fuzztime=5s`; extend corpus |
| **H2** | ASAN / Valgrind CI | **Open** | Priority raised — see P5 |
| **H3** | Integrated smoke script | **Partial** | `ci.yml` runs Ubuntu / macOS / Windows; keep aligned with this doc |

---

# SECTION A — Language completeness

### A1. "Did you mean?" for undefined names

**Status: Done.**  
Levenshtein edit-distance hints via `internal/sema/levenshtein.go` and `suggestName`. Covered by `tests/typo_suggestion.koda` and `koda check` CI.

---

### A2. `struct` declarations

**Status: Done.**  
Named struct types, struct literals, compile-time field validation, O(1) slot access via `koda_struct_get` / `koda_struct_set`. Keep `tests/struct_typo_test.koda` green.

---

### A3. `enum` declarations

**Status: Partial.**  
Parser AST, lowering, and ordinal constant folding are done. The remaining work is the **exhaustiveness warning** when a `switch` subject is a known enum type and not all members are covered. See **P4** for implementation details.

---

### A4. Math builtins

**Status: Done.**  
Full suite in C runtime and `stdlib/math.koda`: `lerp`, `clamp`, `hypot`, trig, `smoothstep`, `smoothdamp`, `approach`, `angleBetween`, `normalize`, `wrap`, `fmod`, `degrees`, `radians`, etc.

---

### A5. `stdlib/timer.koda`

**Status: Done.**  
Cooldown, interval, and countdown helpers. `tests/timer_test.koda` covers the main paths.

---

### A6. Opt-in integer types (`i32`, `u8`, …)

**Status: Open.** See **P1** for full scope.

The numeric type **inference** system (A8) handles pure local arithmetic automatically. Explicit integer types are the separate tool for binary data, C interop, and network work — things inference cannot prove safe.

---

### A7. Struct methods (`box.area()`)

**Status: Open.** See **P2** for full scope.

Object-literal `this` already works. The gap is functions declared inside a `struct` body.

---

### A8. Numeric type inference (`KindInt` / `KindFloat`)

**Status: Done.**

The compiler automatically classifies every `let` variable as `KindInt` (whole-number arithmetic only) or `KindFloat` (may be fractional or unproven). No syntax changes — the programmer writes exactly what they always wrote.

**How it works:**

1. `internal/sema/typeinfer.go` — `InferNumericKinds` pass, called at the end of `PrepareNativeBundle`. Two-pass analysis: seed from initialisers, then narrow from all subsequent assignments.

2. `internal/codegen/intfast.go` — fast path in `emitInfix`. When both operands of `+`, `-`, `*`, `%`, `&`, `|`, `^`, `<<`, `>>` are `KindInt` stack locals, codegen emits direct LLVM integer instructions (`add`, `sub`, `mul`, `srem`, `and`, `or`, `xor`, `shl`, `ashr`) and skips the `unboxNumber` / `boxNumber` round-trips entirely.

**Classification rules:**

| Expression | Result | Reason |
|---|---|---|
| Integer literal (`42`, `int(0)`) | `KindInt` | Provably whole |
| `float64` with zero fraction (`3.0`) | `KindInt` | Same bits as integer |
| `float64` with fraction (`3.14`) | `KindFloat` | Has fractional part |
| `a + b` where both `KindInt` | `KindInt` | Sum of integers is integer |
| `a / b` | `KindFloat` | May produce fraction (`5/2 == 2.5`) |
| `a & b`, `a | b`, `a ^ b`, shifts | `KindInt` | Bitwise ops always produce integers |
| Function call result | `KindFloat` | Return type unknown |
| Escaped / captured variable | `KindFloat` | Escapes to boxed `Value` context |
| Reassigned a `KindFloat` value | `KindFloat` | Narrowed by second pass |

**What this does not fix:**

Inference is limited to provably-integer *local* arithmetic. It does not make `arr[i]` return a raw integer, cannot represent `uint8_t` pixel buffers, and cannot make C FFI calls that take typed pointers safe. That requires explicit type annotations (P1).

**Tests:** `internal/sema/typeinfer_test.go` — 8 unit tests covering all classification cases.

---

# SECTION B — GC hardening

### B1. Incremental major GC

**Status: Done.**  
`gc_collect_incremental` spreads major collection across frames. `gcFrameStep(ms)` is the game-loop API. `tests/incremental_gc_test.koda` in CI. `gc_collect()` remains for shutdown / explicit full collect.

---

### B2. ObjTable open addressing

**Status: Done.**  
Open-addressing via `hashes[]` for tables with capacity ≥ 8; linear scan for smaller tables. FNV-1a hash; 0.75 load factor triggers rehash. `hashes` array allocated with `gc_alloc` / freed with `gc_free` (counted in GC threshold). `tests/table_hash_test.koda`.

---

### B3–B4. GC builtins

**Status: Done.**  
`gcDisable()`, `gcEnable()`, `gcCollect()`, `gcFrameStep(ms)`, `gcStats()`. Tests in `tests/gc_control_test.koda` and `tests/gc_stats_frame_test.koda`.

---

### B5. Write-barrier audit

**Status: Done.**  
`internal/codegen/write_barrier_test.go` scans `koda_runtime.c` for mutation sites and verifies `gc_write_barrier` is called. Run in CI.

---

### B6. Stress suite

**Status: Partial.**  
`tests/stress/stress_mixed_alloc.koda` runs in Linux CI with `timeout 90s`. Expand with: deep recursion, large live object graphs, string pressure, incremental GC under allocation load. Each new stress test should have a comment explaining what failure it would catch.

---

### B7–B13. Runtime memory system

All shipped. See baseline section above for details. The only open item is **B12** (intern table retention docs + `koda_intern_clear()` — see P6).

---

# SECTION C — Runtime hardening

### C1. Typed runtime errors

**Status: Partial.**  
`koda_value_type_name`, `koda_type_error`, `koda_null_error` exist and are used in core paths (`koda_get`, `koda_set`, `koda_unbox_number`). Remaining: audit all native / argv-style builtins and ensure they produce typed errors rather than silent null returns on bad input.

### C2. OOB array / string access

**Status: Done.**  
`koda_panic_str` with index and length on all array read/write/string index paths. `tests/array_oob_get.koda`, `tests/array_oob_set.koda`, `api` diagnostic tests.

### C3. Capacity overflow guards

**Status: Done.**  
`validate_value_slot_count` in `object.c` guards array, table, struct, and closure upvalue allocations against size_t overflow.

### C4. Shadow push/pop balance

**Status: Open.**  
Under `KODA_GC_DEBUG`, add checks at runtime shutdown that `koda_shadow_depth == 0`. A mismatched push/pop is a latent GC bug that only surfaces on unusual code paths. Low priority but catches a whole class of codegen errors.

---

# SECTION D — Developer experience

### D1. Full DWARF debug info

**Today:** `koda build --debug` passes `-g` through to `llc` and `clang`. Symbols exist; source-level mapping does not.

**Next:** thread `Token.File` / `Token.Line` into `DIFile`, `DISubprogram`, and `DILocation` in `internal/codegen` when `EmitDebug` is set. This turns "panic at 0x00401234" into "panic at main.koda:14" — disproportionately valuable for beginners. Dependent on stable line-number tracking through sema (currently partial).

### D2. `koda watch`

**Status: Done.** Poll-based file watching under entry directory; restarts child on change. Optional upgrade: `fsnotify` for event-driven watching + include-graph-aware invalidation (only re-parse files that transitively changed).

### D3. `koda bench`

**Status: Open.** New CLI subcommand. Needs: an iteration protocol the user calls in their Koda program (`benchmarkDone()` or `benchmarkIter(n)`), timing loop in the runner, percentile reporting (p50 / p95 / p99). The fibonacci test (`tests/bench_fib.koda`) can be the first target.

### D4. Runtime property suggestions

**Status: Open.** When accessing a missing property on a table or struct at runtime, emit a "did you mean?" hint. Requires: enumerating table keys safely under `KODA_GC_DEBUG` and running Levenshtein against the accessed name. Debug-gated to avoid production overhead.

---

# SECTION E — Compiler quality

### E1. Parse cache

**Status: Done.** Per-file AST cache keyed by absolute path + `mtime` in `internal/parser/loader.go`. Overlays bypass cache. `loader_parse_cache_test.go`.

### E2. Parallel module parse

**Status: Open.** Meaningful for projects with 10+ included files. Requires: goroutine-safe access to the parser and loader; dependency graph to avoid reprocessing in wrong order. Build on E1.

### E3. Wider constant folding

**Status: Partial.** Current: integer literal arithmetic (`+`, `-`, `*`), unary, and simple `if` / `while` dead-branch elimination. Extend to: string concatenation of literals, bitwise folding, power-of-two detection for `%` and `/`.

### E4–E5. `--warn-unused` and unused binding warnings

**Status: Open.** See **P3** for full scope. These are sema-only; no codegen changes needed.

---

# SECTION F — Standard library completeness

### F1. `stdlib/vec3.koda`

**Status: Done.** `tests/vec3_test.koda`.

### F2–F7. Remaining stdlib modules

**Status: Open.** Implement as pure Koda where possible (fall back to C only for performance-critical paths). Each needs a `tests/*.koda` regression.

| Module | Key exports |
|--------|-------------|
| `@color` | `rgba(r,g,b,a)`, `hsv(h,s,v)`, `lerp(a,b,t)`, `toHex()` |
| `@input` | `keyDown(k)`, `keyPressed(k)`, `mousePos()`, `mouseButton(b)` — thin wrappers over raylib shim |
| `@easing` | `easeIn(t)`, `easeOut(t)`, `easeInOut(t)`, `elastic(t)`, `bounce(t)` |
| `@pool` | `pool(size, construct)`, `pool.get()`, `pool.release(obj)`, `pool.reset()` |
| `@str` | `padStart(s,n,c)`, `padEnd(s,n,c)`, `repeat(s,n)`, `format(tpl,...args)` (beyond current `format`) |
| `@array` | `range(lo,hi)`, `zip(a,b)`, `flatten(arr)`, `unique(arr)`, `sum(arr)`, `max(arr)`, `min(arr)`, `shuffle(arr)` |

---

# SECTION G — Tooling

### G4. Asset embedding in `koda bundle`

**Status: Open.** Append asset directory to the binary or produce a self-extracting archive. Expose `assetPath("file.png")` builtin that resolves relative to the executable at runtime. Document expected folder layout conventions.

### G5. `koda doctor`

**Status: Partial.** Extend environment report: disk space, lld freshness check (compare lld version against what was used to build the runtime archive), runtime archive mtime vs compiler mtime, KODA_PATH resolution.

---

# SECTION H — Final hardening

### H1. Parser fuzzing

**Status: Partial.** `internal/parser/FuzzParse` fuzz target; Linux CI smoke with `-fuzztime=5s`. Extend: larger corpus from real Koda programs, coverage-guided fuzzing with longer budget in nightly CI.

### H2. ASAN / Valgrind CI

**Status: Open.** Priority raised after the 28533c9 round added substantial C complexity (nursery, inline arrays, arena, heap range cache). A sanitizer run catches use-after-free and double-free that the normal test suite cannot detect.

**Scope:** Linux-only CI job (non-blocking initially). Build `runtime/src/*.c` with `-fsanitize=address,undefined`. Link against the sanitized runtime. Run `tests/stress/` and a targeted subset of `tests/*.koda`. Mark H2 Done once the job is green in CI for two consecutive weeks.

### H3. Smoke definition

Keep `ci.yml` native smoke list aligned with this document. Extend the list as `tests/stress/` grows.

---

# Priority stack — what to do next

These are in impact order. The numeric type inference system (A8) is already shipped and does not appear here.

## P1 — Opt-in integer types (`i32`, `u8`, …)

**Why first:** inference (A8) handles pure local arithmetic automatically. The remaining gap is binary data: reading PNG files, network packets, C APIs that take `uint8_t*`. `f64` silently loses bits on 32-bit modular arithmetic. This is the one thing beginners trying to touch real data will hit immediately.

**Relationship to A8:** these are complementary. A8 makes existing code faster for free. P1 adds a new surface area that makes impossible things possible. You write `let pixel: u8 = buf[i]` and get an actual 8-bit integer with correct modular arithmetic and safe C FFI passing.

**Scope:**
- Lexer / parser: type annotation syntax — `let n: i32 = 0`, `let buf: u8`.
- Sema: `LetDecl` carries an optional `TypeAnnotation`; sema validates that assignments are compatible; error on narrowing without explicit cast.
- Codegen: `KindInt` locals with a declared type skip NaN-boxing entirely — stored as raw `i32` / `i64` / `i8` in their stack slot. At Value-boundary sites (array store, function call, return from untyped context), box once.
- Stdlib: document which builtins accept typed integers and what they return.
- Start with `i32` and `u8`. Add `i64`, `i16`, `u16`, `u32`, `u64` once the plumbing is proven.
- Tests: `tests/integer_types.koda` — arithmetic, wrapping overflow, cast round-trips, passing `u8` to a C FFI that takes `uint8_t*`.

---

## P2 — Struct methods (`box.area()`)

**Why second:** the gap between "this language has struct" and "this language feels like any other language I've used". Every modern language beginners have seen lets you put functions on data types. The workaround (`func area(r) { ... }` at module level) reads as broken and causes naming collisions as programs grow.

**Scope:**
- Parser: `struct Rect { w, h; func area() { return this.w * this.h; } }` — `FuncDecl` nodes inside `StructDecl.Body`.
- Sema: bind method names to the struct type; resolve `this` to a typed struct instance inside method bodies; validate that `this.field` only accesses declared fields.
- Codegen: lower method calls on struct-typed variables to `koda_struct_get` / `koda_struct_set` for field access — the same path already used for object-literal `this`.
- `docs/status.md` / `language.md`: update to distinguish object-literal methods (existing, working) from struct body methods (this item).
- Tests: `tests/struct_methods.koda` — method declaration, `this` field access, `this` field mutation, method calling another method.

---

## P3 — `--warn-unused`

**Why third:** the single highest beginner-help-per-effort item. A misspelled variable (`scroe` instead of `score`) compiles today with no warning. The programmer gets a wrong answer and no idea why.

**Scope:**
- Sema pass: after existing analysis completes, walk all `LetDecl` and `FuncDecl` nodes; record a read-count per declaration; warn on any with read-count zero.
- CLI: `--warn-unused` flag, off by default (avoids breaking existing code). Enable by default in `koda check`.
- MASTER_PLAN: mark E4 and E5 Done when landed.
- Tests: `tests/warn_unused.koda` — declare unused let and unused func, expect specific warning text from `koda check --warn-unused`.

---

## P4 — Enum switch exhaustiveness warning

**Why fourth:** game state machines built on enums are the most common beginner enum use. Adding `Transitioning` to a `State` enum and forgetting to handle it in the main switch loop is a silent hole.

**Scope:**
- Sema: when a `switch` subject resolves to a known enum type (via the `varEnum` map already exported by the analyzer), collect all `case` values covered, compare against the full enum ordinal set, emit a warning for each uncovered member.
- Must be a **warning**, not an error — programs using `default:` or intentionally covering only a subset must not break.
- MASTER_PLAN: mark A3 Done when landed.
- Tests: `tests/enum_exhaustive.koda` — partial coverage → warning; `default:` present → no warning; full coverage → no warning.

---

## P5 — ASAN CI job

**Why fifth:** the 28533c9 round added substantial C complexity (nursery, inline array header+elements in one allocation, arena with `gc_unlink_object`, heap range cache with `qsort`). The GC stress suite runs under the normal allocator and will not catch use-after-free on arena-reset objects, double-free on inline array grow paths, or out-of-bounds in the heap cache under unusual allocation sequences.

**Scope:**
- Linux CI job, non-blocking initially. Build `runtime/src/*.c` with `-fsanitize=address,undefined`. Link against the sanitized runtime. Run `tests/stress/` + targeted `tests/*.koda` subset.
- MASTER_PLAN: mark H2 Done once the job is green in CI.

---

## P6 — Intern table retention docs + `koda_intern_clear()`

**Why sixth:** small, bounded, completes the intern table story. After the hash map fix (B11), lookup is O(1). The remaining issue is eviction: strings stay in the intern table until GC sweep. Programs that generate many unique strings (e.g. `format("item_%d", i)` in a tight loop) accumulate them silently until the next full collection.

**Scope:**
- Add one paragraph to `docs/concepts/runtime-and-gc.md` explaining the retention model.
- Add `koda_intern_clear()` to `koda_runtime.c` / `.h`, `internal/codegen/runtime.go`, `builtin_register.go`. No GC invariants change — it just resets the hash table to empty; strings still alive will be re-interned on next allocation.

---

## Longer-term

These become relevant once P1–P6 are closed.

| Item | Depends on | Notes |
|------|-----------|-------|
| Optional / nullable types (`string?`) | P1 type system | Compile-time null safety; requires type annotation layer |
| Interfaces / structural traits | P2 struct methods | `interface Drawable { draw() }` without classes |
| `try` propagation sugar | — | `try expr` early-returns `err(...)` without `.ok` boilerplate |
| Package manager | — | `koda install github.com/user/pkg` |
| Full DWARF debug info | — | `DIFile`/`DISubprogram`/`DILocation`; source-level gdb/lldb |
| WASM target | — | `wasm32-unknown-unknown` via LLVM; drop file-IO builtins |
| `koda bench` | — | D3; iteration protocol + percentile reporting |
| Parallel module parse | E1 (done) | E2; significant for large projects |
| Asset embedding in `koda bundle` | — | G4; `assetPath()` helper |

---

## Hard rules

These apply to every PR regardless of priority.

- **`go test ./...`** green before merging compiler or runtime changes.
- **Full registration chain** for every new C builtin: `koda_runtime.c` / `.h` → `internal/codegen/runtime.go` → `builtin_register.go` / `builtin_globals.go` → docs + tests.
- **Write barriers** on every mutating store into GC objects — enforced by `internal/codegen/write_barrier_test.go`.
- **Keep `runtime.go` and `koda_runtime.c` in sync.**
- **Every new surface area gets a `.koda` or `go test` regression.**
- **No drive-by TODOs** for required behaviour: implement or file a tracked issue with a test gap.

---

## Commands (from a repository checkout)

End users follow **README.md**. For contributors:

```bash
go build -o koda ./cmd/koda
go test ./... -count=1
bash scripts/build-runtime.sh   # or scripts/build-runtime.ps1 on Windows
./koda check path/to/file.koda
./koda run path/to/file.koda
./koda build path/to/file.koda -o out
./koda fmt --check ./...
```

See **CONTRIBUTING.md** and **docs/handoff.md** for the full contributor loop.
