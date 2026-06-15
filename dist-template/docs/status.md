# Koda implementation status

Engineering snapshot of the language surface, LLVM codegen, C runtime, and wrapper/FFI path. **Single native pipeline** -- no bytecode VM.

> **Product positioning:** see [positioning.md](positioning.md). **Prioritized work:** see [ROADMAP.md](ROADMAP.md). **Detailed matrix:** [MASTER_PLAN.md](MASTER_PLAN.md).

---

## Supported and exercised

### Language and frontend

- **Variables:** `let`, `const`, optional type annotations (`int`, `float`, `string`, …).
- **Primitives:** number (64-bit float at runtime), string, bool, null.
- **Composite values:** arrays, objects, **structs** (preferred for game data), enum members.
- **Operators:** arithmetic, comparison, **`==` / `!=`** (legacy `===` warns), logical, compound assignment.
- **Control flow:** `if` / `else`, `while`, `do-while`, `for`, `for ... of`, `switch`, **`match`**, `break`, `continue`, `defer`.
- **String interpolation:** `"Score: {score}"` in double-quoted strings; backtick `` `${expr}` `` templates.
- **Functions:** declarations, calls, recursion, return, function expressions, closures with upvalues.
- **Modules:** `import "@game"`, `import "@array"`, `#include` for low-level shims.
- **Structs / enums:** struct methods, field-checked literals, enum exhaustiveness warnings.
- **Diagnostics:** typo suggestions; `koda check --warn-unused`; unreachable warnings in `strict` lint.
- **Tooling:** `koda doctor` (OK/FAIL), `koda bench` (p50/p95/p99), `assetPath()`, parallel module parse.

### Codegen and toolchain

- **Pipeline:** lexer → parser → sema → LLVM IR → **llc** → **clang** → native executable.
- **`koda run`** uses the same pipeline as **`koda build`** (temp binary).
- **`koda build --no-opt`**, **`koda run --no-opt`**, **`koda build --debug`** (`-g` through llc/clang).
- **Shadow stack** for precise GC roots in compiled code; depth cap configurable via **`KODA_STACK_DEPTH`** (default 131072 frames).

### C runtime and GC

- **NaN-boxed values** -- single `i64` representation in IR.
- **Tri-generational GC** -- nursery (256 KB bump), young, old; incremental major collection.
- **Write barrier** -- O(1) remembered-set membership via `Obj.in_remembered_set`.
- **Conservative roots** -- heap pointer hash set + sorted range cache (per collection, not per stack word).
- **Inline small arrays** -- header + elements in one allocation when `capacity ≤ 64`.
- **Hash table `hashes[]`** -- allocated through `gc_alloc` (counted toward GC thresholds).
- **String intern table** -- open-addressing hash map (FNV-1a); sweep drops unmarked entries.
- **Table tombstones** -- dedicated `TOMBSTONE_VAL` sentinel (boolean `true` is a valid key).
- **Arena allocator** -- `arena`, `arenaReset`, `arenaAllocArray`, `arenaAllocStruct` builtins.
- **Game-loop GC** -- `gcFrameStep(ms)`, `gcStats()`, `gcDisable` / `gcEnable` / `gcCollect`.
- **Loop codegen** -- `while` / `do-while` respect existing terminators (`return` / `break` / `continue` in body).

### Standard library and builtins

- **Core:** `print`, `type`, `typeof`, `len`, `keys`, `number`, `string`, `assert`, `panic`, `ok`, `err`.
- **Time / games:** `deltaTime`, `programTime`, `time`, `sleep`, `clock`, `timestamp`.
- **Math:** trig, `lerp`, `clamp`, `distance`, `random`, `randomInt`, ... (runtime + `stdlib/math.koda`).
- **I/O:** `readFile`, `writeFile`, `parseJSON`, `toJSON`, `listDir`, ...
- **Modules:** `@math`, `@json`, `@io`, `@array`, `@game`, `@timer`, `@vec2`, `@vec3`, `@util`, `@noise`, `@str`, `@color`, `@easing`.

### Native apps, games, and FFI

- **`koda build`** links **`runtime/libkoda_runtime.a`** plus optional **`KODA_NATIVE_SOURCES`** / **`KODA_LINKFLAGS`**.
- **`koda bundle`** -- distributable folder with launcher, assets manifest, and `assetPath()`.
- **`koda.json` `"graphics": true`** -- auto Raylib link flags (no manual `KODA_LINKFLAGS` for beginners).
- **`koda setup raylib`** -- refreshes project shim (overwrites stale copies).
- **`koda wrap upgrade` / `install` / `check` / `list`** -- wrapper catalog and drift detection (`koda doctor`).
- **`koda doctor --fix`** -- auto-refresh stale project `raylib_shim` when `@game` symbols are missing.
- **Raylib** via `@game` + `wrappers/raylib_shim` (~33 fn) or full `@raylib` (548 fn).
- **`args()` / `env()`** -- CLI and environment builtins.
- **`// koda:extern`** -- Koda calls lower to `Value symbol(int argCount, Value* args)`.

---

## Important distinctions

### Struct methods vs object methods

| Feature | Status |
|---------|--------|
| **Methods on object literals** (`let o = { fn draw() { ... } }`) and **`this`** in that context | Supported |
| **Methods declared on named struct types** (`struct Rect { func area() { ... } }`) | **Supported** — `tests/struct_methods.koda` |

Both paths lower through the same struct slot access codegen.

### Numbers

Runtime values remain **64-bit floats** (`Value` NaN-boxing) unless you opt in with a type annotation:

- **`let n: i32 = 0`** — native integer locals, no NaN-box round-trip in pure integer arithmetic
- **`let b: u8 = 255`** — 8-bit modular semantics at boundaries
- Untyped `let` bindings use **numeric type inference** (`KindInt` / `KindFloat`) for fast integer paths where provable

See `tests/integer_types.koda` and [ROADMAP.md](ROADMAP.md).

### Debug symbols

`--debug` emits **`-g`** and attaches LLVM **DI** metadata (`DICompileUnit`, `DISubprogram`, `DILocation`) for `.koda` line mapping. Coverage is growing (**MASTER_PLAN D1**).

---

## Known gaps (honest)

| Gap | Impact | Roadmap tier |
|-----|--------|--------------|
| Full DWARF line coverage on all statements | Complete gdb/lldb `.koda` mapping | Long-term -- **D1** |
| `@app` / retained-mode UI | Desktop apps beyond games | Long-term |
| `try` / err propagation sugar | `ok`/`err` boilerplate | Long-term |

### Smaller hardening items

- **Compound property/index assignment** (`obj.x += 1`) -- verify against current codegen.
- **Object `for-of` key/value** -- `for (let k, v of obj)` supported where codegen lowers it.
- **Declare-before-use** -- should always fail in sema before LLVM; keep tests green.

---

## Recent runtime fixes (merged)

These address findings from the full compiler/runtime review:

- Tombstone sentinel for hash tables
- `emitWhileStmt` / `emitDoWhileStmt` terminator guards
- Remembered-set O(1) dedup
- Intern table hash map
- Conservative stack scan heap cache
- Table `hashes[]` GC byte accounting
- Inline small array allocation
- `koda_string_concat` precise-GC guard + `deltaTime` clamp fix
- Arena allocator and shadow-stack depth configuration

---

## Required release gates

### Native conformance

```powershell
.\bin\koda.exe run .\tests\native_conformance.koda
.\bin\koda.exe build .\tests\native_conformance.koda -o .\tests\native_conformance.exe
.\tests\native_conformance.exe
```

### Graphical / Raylib

```powershell
$env:KODA_NATIVE_SOURCES = '..\wrappers\raylib_shim\wrapper.c'
$env:KODA_LINKFLAGS = '-I..\temp_raylib\src -L..\temp_raylib\src -lraylib -lopengl32 -lgdi32 -lwinmm'
.\koda.exe build .\raylib_brick_breaker.koda -o .\raylib_brick_breaker.exe
.\raylib_brick_breaker.exe
```

### GC and arena smokes

```bash
koda run tests/incremental_gc_test.koda
koda run tests/arena_test.koda
koda run tests/gc_soak.koda    # when exercising GC under load
```

---

## Related

- [positioning.md](positioning.md) -- who Koda is for
- [ROADMAP.md](ROADMAP.md) -- what to build next
- [handoff.md](handoff.md) -- compiler pipeline overview
- [Distribution guide](distribution.md) -- shipping binaries
