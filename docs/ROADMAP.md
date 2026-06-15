# Koda roadmap

Prioritized engineering queue for the Koda compiler and runtime. This replaces the historical week-by-week bootstrap plan with **current** ordering based on product fit, beginner impact, and runtime risk.

> **Execution model:** one pipeline — LLVM IR (`internal/codegen`) + C runtime (`runtime/src`). No bytecode VM.  
> **Positioning context:** [positioning.md](positioning.md) · **Status today:** [status.md](status.md) · **Detailed matrix:** [tests/MASTER_PLAN.md](../tests/MASTER_PLAN.md)

---

## How to read this document

| Tier | Meaning |
|------|---------|
| **Tier 1 — Now** | Production hardening and release gates |
| **Tier 2 — Next** | Strong beginner ROI or runtime safety |
| **Tier 3 — Soon** | Polish, docs, and smaller runtime gaps |
| **Long-term** | Strategic; not blockers for game-maker audience |

Each item lists **MASTER_PLAN** refs where they exist. Update this file when priorities shift.

---

## Shipped in v0.4.0 (do not regress)

These were Tier 1 items; treat breakage as a release blocker.

| Feature | Tests |
|---------|-------|
| Opt-in integer types (`i32`, `u8`, …) | `tests/integer_types.koda` |
| Struct methods (`rect.area()`) | `tests/struct_methods.koda` |
| `--warn-unused` on `koda check` | `tests/warn_unused.koda` |
| Enum switch exhaustiveness warnings | `tests/enum_exhaustive.koda` |
| Stdlib `@color`, `@easing`, `@array`, `@str`, `@pool` | `tests/stdlib_modules_test.koda` |
| `internClear()` + `assetPath()` | `tests/intern_clear_test.koda`, bundle smoke |
| `koda doctor`, `koda new`, `koda.json` manifest | `cmd/koda/doctor.go`, templates |

---

## Tier 1 — Now

### 1.1 ASAN CI as a release gate

**Problem:** runtime complexity (nursery, arena, remembered set, heap cache, inline arrays) needs automated memory bug detection before calling releases production-ready.

**Scope:**

- Linux CI job: build runtime + tests with `-fsanitize=address,undefined` (`scripts/ci-asan-smoke.sh`)
- Remove `continue-on-error: true` once green for two consecutive weeks on `main`
- Document local repro in CONTRIBUTING.md

**Success:** ASAN job blocks merges; intentional UAF test fails the job.

**MASTER_PLAN ref:** **H2**, **P5**

---

### 1.2 Release and Windows parity

**Problem:** Windows dev setup and CI coverage lagged Linux for tier-1 regressions.

**Scope (mostly landed in v0.4.0):**

- `scripts/build-runtime.ps1` matches CI clang + llvm-ar recipe
- `scripts/clang-gnu.cmd` env-configurable (`KODA_LLVM_BIN`, `KODA_MINGW_BIN`)
- `scripts/ci-release-smoke.sh` in release.yml before artifact upload
- Windows CI runs same tier-1 tests as Linux via `ci-gc-stress-timed.ps1`

**Remaining:** SDK zip smoke on publish runner; optional code signing.

---

## Tier 2 — Next

### 2.1 Typed runtime error audit (C1)

**Problem:** some argv-style / native builtins may still misbehave on bad input instead of typed panics.

**Scope:** audit all native / argv-style builtins; ensure `koda_type_error` / `koda_null_error` on bad input.

**Success:** API integration tests for bad-argument paths.

**MASTER_PLAN ref:** **C1**

---

### 2.2 Full DWARF / source-level debug info (D1)

**Problem:** `--debug` emits symbols but not complete `.koda` line mapping in gdb/lldb.

**Scope:** thread `Token.File` / `Token.Line` into `DIFile`, `DISubprogram`, `DILocation` for all statement kinds.

**Success:** panic backtrace shows `main.koda:14` not `0x00401234`.

---

### 2.3 Stress suite expansion (B6)

**Scope:** grow `tests/stress/`; keep Linux + Windows CI lists aligned.

**MASTER_PLAN ref:** **B6**

---

## Tier 3 — Soon

### 3.1 Wider constant folding (E3)

Extend beyond integer literal rules: string concat of literals, bitwise folding, power-of-two detection.

### 3.2 Parser fuzz corpus (H1)

Extend `internal/parser/FuzzParse` corpus from real programs; longer nightly budget.

### 3.3 Diagnostic consistency audit

Compound assign (`obj.x += 1`) codegen verification; align user-facing docs with compiler behavior.

---

## Long-term

| Item | Notes | MASTER_PLAN |
|------|-------|-------------|
| **Optional / nullable types** | `let name: string?`; warn on unchecked deref | new |
| **`try` / err propagation** | Reduce `ok`/`err` boilerplate | new |
| **Interfaces / traits** | Structural checks on method sets | new |
| **Package manager** | `koda install` → `packages/` | new |
| **WASM target** | `wasm32-unknown-unknown`, slim runtime | new |
| **`@app` / retained-mode UI** | Desktop apps beyond games | new |
| **Parallel module parse tuning** | Faster very large projects | **E2** (shipped baseline) |

---

## Completed foundations (do not regress)

Treat these as release blockers if they break. Detail in [tests/MASTER_PLAN.md](../tests/MASTER_PLAN.md).

**Pipeline:** lexer → parser → sema → LLVM → clang → binary; `PrepareNativeBundle` before codegen.

**Runtime / GC:**

- Tri-generational GC, shadow stack, write barriers
- O(1) remembered set; conservative root heap cache
- Arena builtins; inline small arrays; `gcFrameStep`
- Tombstone sentinel; intern hash table; table `hashes[]` GC accounting
- Loop terminator guards in while/do-while codegen

**Language:** structs, enums, closures, defer, `ok`/`err`, typo hints, `for-of`, range `lo..hi`, numeric type inference.

**Tooling:** `koda check`, `koda fmt`, `koda watch`, `koda bench`, `--no-opt`, `--debug` (`-g`), CI smokes on Ubuntu/macOS/Windows.

---

## Testing strategy

### Every merge

```bash
go test ./... -count=1
powershell -File scripts/build-runtime.ps1   # Windows
bash scripts/build-runtime.sh              # Linux / macOS
```

### Before release

```bash
bash scripts/ci-release-smoke.sh ./koda
bash scripts/ci-native-smoke.sh              # full native matrix (CI)
```

### Benchmarks (informational)

- Fibonacci (recursion / shadow stack)
- Binary trees (GC stress)
- Per-frame arena + `gcFrameStep` game loop

---

## Success criteria by audience

| Audience | Ready when |
|----------|------------|
| **Game-making beginners** | v0.4.0 SDK zip + `@game` + `koda doctor` OK |
| **General app beginners** | struct methods + warn-unused + docs hub |
| **C / systems learners** | integer types + FFI via wrapgen + ASAN CI green |
| **Contributors** | `go test` green, tests/MASTER_PLAN updated on merge |

---

## Related

- [positioning.md](positioning.md) — honest product framing
- [status.md](status.md) — what works today
- [tests/MASTER_PLAN.md](../tests/MASTER_PLAN.md) — full engineering matrix
- [handoff.md](handoff.md) — pipeline for new contributors
- [releasing.md](releasing.md) — tagging `v*` and SDK zips
