# Koda vs `compiler.md` — implementation checklist

This file tracks the **[compiler.md](compiler.md)** master specification against **this repository** (`cmd/koda`, `internal/parser`, `internal/codegen`, `cmd/wrapgen` / **wrapgen**, `runtime/src`).

**Legend**

- [x] **Implemented** — usable via `koda run` / `koda build` as documented in [KODA_PROGRAMMER_REFERENCE.md](KODA_PROGRAMMER_REFERENCE.md) or tests.
- [~] **Partial / different** — behavior exists but differs from the spec snippet, is LLVM-only, or needs the noted caveat.
- [ ] **Not implemented** — no lexer/AST/runtime support yet, or tooling absent.

**Re-verify after changes**

```powershell
go test ./... -count=1
.\bin\koda.exe run .\tests\native_conformance.koda
# optional native gate (needs clang):
.\bin\koda.exe build .\tests\native_conformance.koda -o .\tests\native_conformance.exe
```

---

## Language overview & goals

- [x] JavaScript-like syntax core (`let`, `func`, `if`, loops, objects, arrays).
- [x] Single native backend: LLVM IR + C runtime (`koda run` / `koda build` / `koda bundle`).
- [x] Access to C libraries via generated glue (`// koda:extern` + `wrapper.c` from **wrapgen**, not `#native` / `#ffi` syntax).
- [ ] `#native` / `#ffi` directives as in spec grammar — use **wrapgen** + extern lines instead.
- [ ] Optional static type hints (`x: number`) — spec “future”; not in parser.

---

## Lexical structure

- [x] UTF-8 source, identifiers with letters/digits/underscore (Unicode-aware start per lexer).
- [x] `//` single-line comments.
- [x] `/* … */` multi-line comments.
- [~] `///` doc comments — treated as ordinary `//` (no **kodadoc** tool).
- [x] Double-quoted and single-quoted strings with escapes (`\n`, `\t`, `\r`, `\\`, `\"`, `\'`, `\xNN`, `\uNNNN`, `\u{...}`).
- [x] `` `template ${expr}` `` template literals.
- [x] `"""` … `"""` multi-line strings.
- [x] `r"…"` raw strings (lexer path via `r` + `"`).
- [x] Numbers: decimal, float, scientific, `0x` hex, `0b` binary, `0o` octal, `_` separators.
- [x] `true`, `false`, `null`.
- [x] Reserved **future** keywords from spec (`try`, `class`, …) — not special tokens; many parse as identifiers until rejected in context.

---

## Keywords & operators (core language)

- [x] Control: `if`, `else`, `for`, `while`, `do`/`while`, `switch`, `case`, `default`, `break`, `continue`, `return`.
- [x] Declarations: `let`, `func`.
- [x] `import` as **expression** `import "path"` (see modules).
- [x] `in` and `of` in `for (let x in/of iterable)`.
- [ ] `this` as a dedicated keyword / `TokenThis` — **not** reserved; object methods use identifiers per [KODA_PROGRAMMER_REFERENCE.md](KODA_PROGRAMMER_REFERENCE.md) §6.
- [x] Arithmetic `+ - * / % **`, comparisons, `==` `!=` `===` `!==`, logical `&& || !`, bitwise `& | ^ ~ << >> >>>`, compound assigns, ternary `?:`.
- [x] `>>=` and `>>>=` compound assignment tokens ([lexer.go](internal/lexer.go)).
- [x] Unary `+ - ! ~`, prefix/postfix `++ --`.
- [x] **Arrow functions** (`=>`) — parsed as sugar for `func` expressions (`id => …`, `(a, b) => …`, block or expression body).

---

## Type system & builtins (runtime)

- [x] Dynamic types: number, string, bool, null, array, object, function.
- [x] `type(x)` builtin.
- [x] `is_number`, `is_string`, `is_bool`, `is_null`, `is_array`, `is_object`, `is_function`.
- [x] `number()`, `string()`, `bool()` conversions (see evaluator/native for edge cases).
- [~] Truthiness — largely JS-like; empty string `""` falsy; confirm `[]` / `{}` vs spec in tests if you rely on it.

---

## Variables & scope

- [x] `let x = expr;`, multiple `let` inits where parser allows.
- [x] `let x;` → null.
- [x] Block scope, function scope, closures / upvalues.
- [~] `let a = 1, b = 2` in **for-loop** init — supported; top-level comma `let` chains depend on parser rules.
- [x] Shadowing in nested blocks (typical patterns).

---

## Functions

- [x] `func name(args) { … }`, `return`, recursion.
- [x] First-class `func (…) { … }`, IIFE patterns.
- [~] Default parameters — supported with limitations (literal defaults; see HANDOFF if present).
- [x] `…rest` parameters (VM + native per prior docs).
- [x] **Arrow functions** — see lexical (`=>`).
- [ ] **Variadic** spec examples using `.length()` on rest bundle — verify rest is a real array with `length()` in all backends.

---

## Control flow

- [x] `if` / `else`, `while`, `do`/`while`, C-style `for`, `break`, `continue`.
- [x] `switch` / `case` / `default` — **C-style fall-through**; use `break` to exit the switch or to avoid running the next `case` / `default`.
- [ ] `switch` as **expression** (`let m = switch …`).
- [ ] `if` as **expression** (`let max = if (a > b) a else b`).
- [x] `for (let x in arr)` (**keys**, numeric indices for arrays) and `for (let x of arr)` (**values**); `for (let [k, v] of …)` destructuring for keys+values (`tests/for_of_pairs.koda`).
- [ ] `for (let i, item in items)` — **second** binding not parsed (only single `let ident`).
- [ ] `for (let key, value in obj)` — not parsed.
- [ ] Labeled `break` / `continue`.
- [x] Ternary `cond ? a : b` as expression.

---

## Data structures

- [x] Array literals, indexing, nesting.
- [x] Object literals, dot / bracket access, dynamic assignment of fields.
- [ ] **Computed keys** `{ [expr]: value }` — object parser requires `ident :` only today.
- [x] **Shorthand properties** `{ name, health }` — same as `{ name: name, … }` for identifier keys.
- [~] **Method shorthand** `heal(amount) { … }` on objects — see programmer reference §7 / IMPLEMENTATION_STATUS.
- [ ] **Negative indices** and **slicing** `[a:b]` on arrays/strings — spec §618–644 / §768–769.
- [~] Array methods: subset on natives (`push`, etc. per `native.go`); not full spec `map`/`filter`/`reduce` unless implemented on array type.
- [~] String methods: `upper`, etc. per `StringMethods` — not full spec list (`trim`, `split`, …) unless present in `native.go`.
- [x] `delete` on **object** properties via bracket form `delete obj["key"]` (returns whether the key existed); not for array elements.
- [ ] Object spread `{ …obj }`.
- [ ] Object destructuring `let {a,b} = obj`.
- [x] `Map()` / `Set()` builtins (constructors + `[]` / methods per [native.go](internal/sema.go) / [mapset.go](internal/sema.go); VM + C runtime).
- [x] **Tuples** `(a, b, …)` — at least two elements; immutable; indexed with numbers only.

---

## Module system

- [x] `#include "path"` and `#include <name>` with `KODA_PATH` / `KODA_WRAPPERS` resolution ([loader.go](internal/parser/loader.go)).
- [x] `import "path"` expression form (and `@` modules per loader docs).
- [ ] `#include "file" as alias` — not implemented.
- [ ] `#include "file" { sym1, sym2 }` selective import — not implemented.
- [ ] `#include "file" *` — not implemented.
- [x] Circular dependency detection for includes/imports.
- [ ] Namespaced `file.symbol` access pattern — depends on flat merge vs module object (verify intended style).

---

## FFI & C interop

- [ ] Spec `#native` / `#ffi` pipeline.
- [x] **`// koda:extern`** lines + **`KodaValue` wrapper symbols** in C ([native_emit.go](internal/codegen.go), [WRAPPERS.md](WRAPPERS.md)).
- [x] **wrapgen** (`go build -o wrapgen ./cmd/wrapgen`) generates `.koda` + `wrapper.c` + docs; merged into root `go.mod`.
- [x] **`koda wrap`** forwards to `wrapgen` (or legacy `kujiwrap`) when installed beside `koda`.
- [ ] Callback / finalizer attributes from spec (`[finalizer: …]`).

---

## Standard library (spec modules)

- [x] Core: `print`, `type`, conversions, `len`, `time`, `sleep`, math helpers (`abs`, `sqrt`, `random`, … per [native.go](internal/sema.go)).
- [~] File I/O natives (`readFile`, `writeFile`, …) — present as builtins; not necessarily `file.read` namespaced API from spec.
- [ ] Namespaced `#include <math>` / `math.sin` style stdlib modules from spec.
- [~] `json` — global **`json`** object with `parse`, `stringify`, and `try_parse` (VM + tree + C runtime); not the spec namespaced `#include <json>` module shape.
- [ ] `http`, `os`, `regex` modules as described in spec.
- [x] `input` optional prompt — one line from **stdin** (`input()` / `input("> ")`); VM + native.
- [ ] `warn` as in spec (use `print` today).

---

## Error handling & assertions

- [ ] `try` / `catch` / `finally`.
- [ ] `assert()` builtin with release stripping.

---

## Memory management

- [x] GC for heap objects in C runtime (mark/sweep; see `runtime/src/gc.c`).
- [x] `gc()` no-op or stub per natives.
- [ ] `alloc` / `free` / `with` blocks from spec.
- [ ] Reference-counted C handles with finalizers as in spec examples.

---

## Compiler architecture (this repo)

- [x] Hand-written lexer ([internal/lexer](../internal/lexer/)).
- [x] Recursive descent parser → AST ([internal/parser](../internal/parser/)).
- [x] **Sema** and native emit prep ([internal/sema](../internal/sema/)).
- [x] LLVM IR via **llir** + **llc** + **clang** linking **`runtime/libkoda_runtime.a`** ([internal/codegen](../internal/codegen), [runtime/src](../runtime/src)).
- [~] Optimizer — sema/constant folding; not full spec optimizer list.

---

## Runtime (values)

- [x] NaN-boxing in C ([runtime/src/value.h](../runtime/src/value.h)); **i64** in LLVM IR.
- [x] Closures, upvalues, globals, locals.

---

## Tooling (spec “Tooling Ecosystem”)

- [x] `koda run`, `check`, `disasm`, `build`, `bundle`, `wrap`, `paths`, `doctor`, `version`, `help`.
- [x] **wrapgen** + `koda wrap` (wrapper generation).
- [ ] `koda watch`, `koda fmt`, REPL, `koda test`, `koda profile`, `koda debug` (if added).
- [ ] Package manager, LSP, etc. (roadmap).

---

## Grammar notes (EBNF in spec)

Anything listed under **“Not implemented yet”** in [compiler.md § Implementation Status](compiler.md) plus:

- [x] **Expression** range `a..b` (`TokenDotDot`) — inclusive integer sequence as an **array** after truncating bounds toward zero; descending ranges supported. **`switch` / `case` with `90..99`** is not implemented.
- [ ] `templateString` fully aligned with all edge cases in spec.
- [x] `doWhileStmt`, `forStmt` with `in`/`of`, `switchStmt` core cases.

---

## Summary counts (manual audit)

| Category | Approx. |
|----------|-----------|
| Checked **implemented** | Core language, lexer, LLVM path, `#include`, `import"`, wrapgen/extern, CLI |
| Partial / differs | for-in arity vs spec, stdlib naming vs spec modules, `json` as global object not `#include` module |
| Not implemented | `#native`/`#ffi`, rich modules, most tooling, many spec sugar features (`switch` case ranges, object spread, …) |

When you close a gap, update the matching line to `[x]` and add a short note with the PR or commit.
