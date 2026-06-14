# Koda implementation handoff

**Audience:** Engineers working on the **native LLVM** compiler (lexer → parser → sema → LLVM IR → clang → static binary).  
**Architecture overview:** [architecture.md](architecture.md).  
**Language surface:** [language.md](../language.md) (root) or [docs/language.md](language.md), [README.md](../README.md).  
**Shipped changes:** [CHANGELOG.md](../CHANGELOG.md).  
**Long-term roadmap (maintainers):** [MASTER_PLAN.md](MASTER_PLAN.md).

**Stack:** Go **1.22+**, module **`koda`** ([go.mod](../go.mod)), LLVM IR via [llir/llvm](https://github.com/llir/llvm) v0.3.6, C11 runtime under **`runtime/src/`** (linked as **`runtime/libkoda_runtime.a`**).

There is **no bytecode VM** in the supported path: one pipeline end-to-end.

---

## 1. Product (CLI)

| Command | What it does | Toolchain |
|--------|----------------|-----------|
| **`koda run`** | Load bundle → sema → LLVM → temp native exe → run | **llc** + **clang** + runtime archive (or embedded release binary) |
| **`koda watch`** | Watch `.koda` files under entry dir; rebuild + restart temp exe on change | Same as `run` |
| **`koda build`** | Same lowering → linked **`*.exe`** / binary | Same |
| **`koda check`** | Lexer + parser + **`sema.PrepareNativeBundle`** only (no LLVM) | Go only |
| **`koda fmt`** | AST-based canonical format | Go only |
| **`koda bundle`** | Package entry + assets for distribution | See **`cmd/koda`** |

Release builds (`**-tags release**`) embed **llc** / **clang** (and **lld** on Windows) plus **`libkoda_runtime.a`** so end users need no LLVM install.

---

## 2. Repository map

| Path | Role |
|------|------|
| **`internal/lexer`** | Tokens; **`NewLexer(src, file)`** threads diagnostics paths |
| **`internal/parser`** | AST, **`LoadProgram`**, `#include` / import flattening, math prelude injection |
| **`internal/sema`** | **`PrepareNativeBundle`**, escape/shadow layout, arity, builtins prelude, **`Analyze`** |
| **`internal/codegen`** | LLVM **`ir.Module`** emission; **`internal/codegen/runtime.go`** declares **`koda_*`** — must match **`runtime/src/koda_runtime.c`** |
| **`internal/nativebuild`** | **llc** + **clang** invocation |
| **`internal/diagnostic`** | **`DiagnosticError`**, **`MultiError`**, snippet formatting |
| **`internal/formatter`** | **`koda fmt`** |
| **`runtime/src`** | NaN-boxed **`Value`**, GC, objects, **`koda_runtime_init`**, builtins |
| **`stdlib/*.koda`** | Optional **`@math`**, **`@vec2`**, etc. |
| **`tests/*.koda`**, **`examples/`** | Regression and demos |

---

## 3. Invariants (do not break)

1. **Sema blocks codegen** — **`PrepareNativeBundle`** errors must prevent emitting broken IR.
2. **Runtime / LLVM declare drift** — every **`koda_*`** used from generated IR exists in C with the same ABI (**`internal/codegen/runtime.go`** ↔ **`koda_runtime.c`**).
3. **Builtin names** — **`internal/sema/builtin_globals.go`** (sema prelude) and **`internal/codegen/builtin_register.go`** (codegen prelude) must agree for globals **`koda build`** can link.
4. **Shadow stack** — every **`koda_push_frame`** matched by **`koda_pop_frame`** on **every** exit edge: each LLVM **`ret`** path must emit a pop (functions with multiple **`return`** statements have multiple pop sites). **`defer`** runs before pop on **`return`**.
5. **GC write barriers** — any new mutator that stores an **old → young** reference must use the same barrier pattern as **`koda_object_set`** / **`koda_array_set`**.

---

## 4. Session checklist

1. **`go test ./...`** (and **`go vet ./...`**) — same coverage as **`.github/workflows/ci.yml`** (Ubuntu, macOS, Windows).
2. **`go build -o koda ./cmd/koda`** then **`koda run tests/hello.koda`**
3. After runtime C changes: rebuild **`runtime/libkoda_runtime.a`** (see **`scripts/build-runtime.ps1`** / **`.sh`**)
4. **`CHANGELOG.md`** — update **`[Unreleased]`** for user-visible compiler or runtime changes; see **[releasing.md](releasing.md)** to cut **`v*`** tags.

---

## 5. Historical note

Older docs referred to a **bytecode VM** and **`koda`** binaries; that dual path was removed so the tree matches **one** execution model. See **`docs/architecture.md`** § *What happened to the bytecode VM?* for context.
