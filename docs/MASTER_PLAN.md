# Koda — master plan to 100%

Long-term roadmap for **language completeness**, **GC/runtime hardening**, **DX**, **compiler quality**, **stdlib**, and **tooling**. This document is for **people working in the Koda repository** (not end users — see **[README.md](../README.md)**). Update the **progress matrix** when you land work.

---

## Engineering standards

- Correctness first: memory safety, GC invariants, and clear diagnostics beat feature count.
- Read existing code before changing it; match local style; do not leave known holes.
- Prefer complete implementations over TODOs and stubs.
- **`go test ./...`** green before merging substantive compiler or runtime changes.
- C runtime: keep **`runtime/src/*.c`** and **`internal/codegen/runtime.go`** declarations aligned.

---

## Do not regress (shipped baseline)

These properties are expected today; treat regressions as release blockers.

- Full **lexer → parser → sema → LLVM IR → clang → binary** pipeline.
- **`lexer.NewLexer(source, file string)`** — file path on tokens.
- **`sema.PrepareNativeBundle`** before codegen — bad programs must not reach LLVM.
- **`diagnostic.MultiError`** — aggregate sema issues with snippets where applicable.
- Call **arity** for `func` / `// koda:extern` / known argv-style methods.
- **`for`** loop head bindings scoped to the body.
- **`defer`** — LIFO; **`return`** value computed before defers run.
- **`for (let i of lo..hi)`** dynamic range — counted **`i64`** loop, unboxed bounds.
- Integer literal folding (**`+` / `-` / `*`**, unary **`-` / `+`**, nested safe cases).
- **`koda check`** — full sema, no LLVM.
- **`koda fmt`** / **`koda fmt --check`** / **`./...`** expansion.
- **`koda run --no-opt`** / **`koda build --no-opt`** / **`api.BuildOptions`**.
- **`koda build --debug`** — **`-g`** through **`llc`** and **`clang`** (symbols; not full DI source maps).
- **`koda watch`** — poll-based rebuild + restart under entry directory.
- Generational **GC**, **shadow stack**, **write barriers**, intern table sweep; string **`==`** stable.
- **`koda_runtime_init_ex` / `koda_runtime_shutdown`** from generated **`main`**.
- Shadow stack growth + hard cap panic message; **`koda_shadow_stack_high_water()`**.
- **`ObjArray` / `ObjTable`** NULL guards on partial allocation paths.
- Globals + **`koda_register_global_slot`** wiring.
- **`KODA_GC_DEBUG`** diagnostics where implemented.
- **`stdlib/math.koda`**, **`stdlib/vec2.koda`**, **`stdlib/timer.koda`**.
- **`ok` / `err`**, **`panic`**, **`assert`**, **`readFile`**, **`writeFile`**, **`parseJSON`**, **`gcFrameStep`**, **`gcStats()`**, **`gcDisable` / `gcEnable` / `gcCollect`** (GC pressure relief).
- **Runtime diagnostics** — out-of-bounds array/string index **`koda_panic_str`**; **`koda_type_error`** / **`koda_unbox_number`** on bad types for core **`[]` / unbox** paths.
- Release **`-tags release`** embed story; **CI** (`go test`, fmt, check, native smokes, GC soak).

---

## Progress matrix (maintainers: update on merge)

Legend: **Done** = meets plan intent in tree today · **Partial** = exists but missing plan details · **Open** = not started or only sketched.

| Ref | Topic | Status | Notes |
|-----|--------|--------|------|
| **A1** | “Did you mean?” for undefined names | **Done** | `internal/sema/levenshtein.go`, `suggestName` hints; `tests/typo_suggestion.koda` |
| **A2** | `struct` types + literals + field checks | **Done** | `tests/struct_typo_test.koda` + CI expects **`koda check`** failure |
| **A3** | `enum` + constant folding + switch hints | **Partial** | Parser AST + lowering; verify formatter + exhaustiveness **warning** story |
| **A4** | `math.lerp` / `clamp` / `hypot` | **Done** | C + `stdlib/math.koda` + tests elsewhere |
| **A5** | `stdlib/timer.koda` | **Done** | Library + `tests/timer_test.koda` |
| **B1** | Incremental major GC | **Partial** | `gc_collect_incremental`, `koda_gc_frame_step`; **`tests/incremental_gc_test.koda`** (CI **`--no-opt`**); validate game-loop budgets vs plan |
| **B2** | ObjTable open addressing | **Partial** | `hashes[]` path exists in `object.c`; confirm threshold/probing vs plan + `tests/table_hash_test.koda` |
| **B3** | `gcDisable` / `gcEnable` / `gcCollect` builtins | **Done** | Aliases to **`koda_gc_disable` / `koda_gc_enable` / `koda_gc_collect`**; **`tests/gc_control_test.koda`** |
| **B4** | `gcStats()` object | **Done** | Lowercase keys; `tests/gc_stats_frame_test.koda` |
| **B5** | Write-barrier static audit | **Done** | **`internal/codegen/write_barrier_test.go`** |
| **B6** | `tests/stress/` suite + CI job | **Partial** | **`tests/stress/stress_mixed_alloc.koda`** + Linux CI (**`timeout 90s`**); expand directory over time |
| **C1** | Typed runtime errors (`koda_type_error`, …) | **Partial** | **`koda_value_type_name`**, **`koda_type_error`**, **`koda_null_error`**; **`koda_get` / `koda_set` / `koda_unbox_number`**; argv/runtime audit remains |
| **C2** | OOB array access → panic w/ index | **Done** | Array read/write + string index; **`tests/array_oob_*`**, **`api` tests** |
| **C3** | Capacity overflow guards | **Done** | **`validate_value_slot_count`** in **`object.c`** (array/table/struct/closure upvalues) |
| **C4** | Shadow push/pop balance (debug) | **Open** | Extra checks under `KODA_GC_DEBUG` / shutdown |
| **D1** | LLVM debug **source locations** (DI) | **Partial** | **`--debug`** emits **`-g`**; full **`DIFile` / `DISubprogram`** not implemented |
| **D2** | `koda watch` | **Done** | Polling watcher (not `fsnotify`); upgrade optional |
| **D3** | `koda bench` | **Open** | |
| **D4** | Property “did you mean?” at runtime | **Open** | Needs debug runtime path + key enumeration |
| **E1** | Parse cache (`mtime`) | **Done** | **`internal/parser/loader.go`** (abs path + **`mtime`**); overlays bypass cache; **`loader_parse_cache_test.go`** |
| **E2** | Parallel module parse | **Open** | |
| **E3** | Wider constant folding | **Partial** | Extend beyond current literal rules |
| **E4** | Dead-code elimination + `--warn-unused` | **Open** | |
| **E5** | Unused binding warnings | **Open** | |
| **F1–F7** | stdlib `vec3`, `color`, `input`, `easing`, `pool`, `str`, `array` | **Partial** | **`stdlib/vec3.koda`** + **`tests/vec3_test.koda`** (rest **Open**) |
| **G4** | Asset embedding in `koda bundle` | **Open** | |
| **G5** | `koda doctor` depth | **Partial** | Extend probes (disk, lld freshness, etc.) |
| **H1** | Parser fuzzing | **Partial** | **`internal/parser/FuzzParse`** + Linux CI smoke (**`-fuzztime=5s`**); extend corpus/time |
| **H2** | ASAN / Valgrind CI | **Open** | |
| **H3** | Final integrated smoke script | **Partial** | **`.github/workflows/ci.yml`** runs native smoke on **Ubuntu, macOS, and Windows**; unify doc list |

---

# SECTION A — Language completeness

### A1. “Did you mean?” for undefined names

**Goal:** typo’d identifiers get an edit-distance hint on diagnostics.

**Status:** implemented in **`internal/sema`** (`levenshtein`, `suggestName`). Keep **`tests/typo_suggestion.koda`** and `koda check` coverage green.

---

### A2. `struct` declarations

**Goal:** named struct types, literals, compile-time field validation, O(1) field access via **`koda_struct_get` / `koda_struct_set`**.

**Remaining:** `tests/struct_typo_test.koda` (and `koda check` expectations). Formatter polish if gaps remain.

---

### A3. `enum` declarations

**Goal:** zero-cost members; **`switch`** exhaustiveness **warnings** when subject is enum-typed.

**Remaining:** verify formatter + sema diagnostics match the plan; implement missing warning paths.

---

### A4. `math.lerp`, `math.clamp`, `math.hypot`

**Status:** present in runtime + **`stdlib/math.koda`**. Add or consolidate **`tests/math_test.koda`** if you want a single dedicated file.

---

### A5. `stdlib/timer.koda`

**Status:** shipped with tests. Keep **`update()`** + loop paths covered (shadow / multi-return).

---

# SECTION B — GC hardening for large programs

### B1. Incremental major GC

**Goal:** spread major collection work across frames; keep stop-the-world **`gc_collect()`** for shutdown/explicit full collect.

**Work:** **`tests/incremental_gc_test.koda`** exercises **`gcCollectIncremental`** + **`gcFrameStep`** in a loop (Linux CI: **`koda run --no-opt`**, time-bounded). Remaining: larger workloads and game-loop budget validation (and optional `tests/stress/` variant).

---

### B2. ObjTable O(1) hashing

**Goal:** open addressing over **`hashes[]`** for larger tables; small tables stay compact.

**Work:** confirm FNV / load factor / resize matches design; extend **`tests/table_hash_test.koda`** as needed.

---

### B3. GC pressure relief builtins

**Goal:** **`gcDisable()`**, **`gcEnable()`**, **`gcCollect()`** with clear semantics (distinct from shadow-stack toggles — design carefully).

**Work:** C ABI + **`builtin_globals.go`** + **`runtime.go`** + tests; document hazards (no allocations while “disabled” if that is the model).

---

### B4. `gcStats()` builtin

**Status:** shipped. Keys are **lowercase** in object form (Koda property normalization).

---

### B5. Write-barrier audit (automated)

**Goal:** CI **`go test`** scans **`koda_runtime.c`** for risky stores without **`gc_write_barrier`**.

---

### B6. Stress suite (`tests/stress/`)

**Goal:** time-bounded jobs: large live graphs, deep recursion, string pressure, mixed alloc, incremental loop stability.

---

# SECTION C — Runtime hardening

### C1–C4. Null/type errors, OOB array panic, capacity validation, shadow balance checks

Implement as described in the original backlog: prioritize **high-frequency natives** first, then tables/arrays/structs.

---

# SECTION D — Developer experience

### D1. Debug source maps / rich DI

**Today:** **`koda build --debug`** emits **`-g`**.

**Next:** thread token file/line into **`DIFile` / `DISubprogram` / `DILocation`** in **`internal/codegen`** when `Debug` is set.

### D2. `koda watch`

**Today:** directory poll of **`.koda`** files; restart child on change.

**Optional upgrade:** `fsnotify` (new dependency) + include-graph aware watching.

### D3. `koda bench`

New **`koda bench`** command + iteration protocol (`benchmarkDone()` or similar) + percentile reporting.

### D4. Runtime property suggestions

Near-match hints on missing keys (debug-gated; may require enumerating table keys safely).

---

# SECTION E — Compiler quality

### E1–E2. Parse cache + parallel parse

**E1 (done):** per-file AST cache keyed by path + **`mtime`** in **`internal/parser/loader.go`** (overlays bypass).

**E2:** parallel module parse remains **Open**.

### E3–E5. Constant folding, DCE, unused warnings

Broaden folds safely; add reachability / unused-symbol passes behind **`--warn-unused`**.

---

# SECTION F — Standard library completeness

**`vec3`** shipped (**`stdlib/vec3.koda`**, **`tests/vec3_test.koda`**); ship **`color`**, **`input`**, **`easing`**, **`pool`**, expanded **`str`** / **`array`** helpers as **pure Koda** where possible and add **`tests/*.koda`** for each module.

---

# SECTION G — Tooling

### G1–G3

**Watch**, **bench**, **debug** — see sections **D2**, **D3**, **D1**.

### G4. Asset embedding in `koda bundle`

Append/archive assets; **`assetPath()`**-style helper; document layout conventions.

### G5. `koda doctor`

Richer environment report (disk, lld, runtime archive freshness vs compiler, etc.).

---

# SECTION H — Final hardening

### H1. Parser fuzzing

**`internal/parser`** fuzz target; short CI budget.

### H2. ASAN / Valgrind

Optional CI job building runtime with sanitizers + running stress subset.

### H3. Smoke definition

Keep **`ci.yml`** native smoke list aligned with this document; extend as **`tests/stress/`** lands.

---

## Hard rules (from backlog)

- **No drive-by TODOs** for required behavior: implement or file a tracked issue with a test gap.
- **Full registration chain** for every new C builtin: **`koda_runtime.c`** / **`.h`** → **`internal/codegen/runtime.go`** → **`builtin_register.go`** / **`builtin_globals.go`** → docs/tests.
- **Write barriers** on every mutating store into GC objects (see **B5**).
- **Keep `runtime.go` and `koda_runtime.c` in sync.**
- **Every new surface area gets a `.koda` or `go test` regression.**

---

## Commands (from a repository checkout)

End users should follow **[README.md](../README.md)**. When hacking **this** repo:

```bash
go build -o koda ./cmd/koda
go test ./... -count=1
bash scripts/build-runtime.sh   # or scripts/build-runtime.ps1 on Windows
./koda check path/to/file.koda
./koda run path/to/file.koda
./koda build path/to/file.koda -o out
./koda fmt --check ./...
```

See **[CONTRIBUTING.md](../CONTRIBUTING.md)** and **[docs/handoff.md](handoff.md)** for the full contributor loop.
