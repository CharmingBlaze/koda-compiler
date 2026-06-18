# Koda compiler architecture (current)

**Audience:** contributors implementing language features, wrappers, and tooling.

**Goal:** one language with progressive depth — beginner-friendly surface, native LLVM backend, full C/C++ wrapper access underneath.

> **Product vision:** [KODA_LANGUAGE_ROADMAP.md](KODA_LANGUAGE_ROADMAP.md)  
> **Wrapper design:** [KODA_WRAPPER_SYSTEM.md](KODA_WRAPPER_SYSTEM.md)  
> **User-facing pipeline summary:** [architecture.md](architecture.md)

---

## Design principle

```text
Koda = easy entry + serious ceiling
```

One language, not two. Simple programs stay simple; expert programs use structs, methods, enums, modules, raw Raylib, and (planned) explicit `unsafe` — without a separate “beginner dialect.”

---

## End-to-end pipeline

```text
.koda source
  → Loader (imports / #include graph)
  → Lexer (internal/lexer)
  → Parser (internal/parser) → AST
  → Sema (internal/sema) → typed bundle
  → Codegen (internal/codegen) → LLVM IR (.ll)
  → Native build (internal/nativebuild) → clang + libkoda_runtime.a → executable
```

There is **one execution model**: native binary. No bytecode VM on the active path.

| Stage | Package | Key files |
|-------|---------|-----------|
| Lexer | `internal/lexer` | `lexer.go`, `token.go` |
| Parser | `internal/parser` | `parser.go`, `parser_stmt.go`, `parser_expr.go`, `ast.go`, `loader.go`, `include_flatten.go`, `prelude.go` |
| Diagnostics | `internal/diagnostic` | `diagnostic.go`, `multi_error.go` |
| Sema | `internal/sema` | `sema.go`, `typeinfer.go`, `sema_struct_enum.go`, `prepare_native.go`, `optimizer.go`, `levenshtein.go` |
| Codegen | `internal/codegen` | `codegen.go`, `statements.go`, `functions.go`, `objects.go`, `struct_methods.go`, `runtime.go` |
| Native link | `internal/nativebuild` | `build.go`, `clang.go`, `linkflags.go`, `raylib_vendored.go` |
| SDK paths | `internal/kodahome` | `home.go`, `resolve.go`, `toolchain.go` |
| C runtime | `runtime/src` | `koda_runtime.c`, `value.h`, GC, builtins |
| CLI | `cmd/koda` | `main.go`, `project.go`, `test.go`, `doctor.go`, `setup.go`, … |
| IDE | `koda-ide/` | Wails app; calls `koda/api` |
| Public API | `api/` | `studio.go`, native build helpers |

---

## Loader and modules (today)

Three mechanisms exist **now**:

| Mechanism | Syntax | Behavior |
|-----------|--------|----------|
| **Include** | `#include "path.koda"` | Paste / flatten at compile time |
| **Stdlib import** | `#include "@math"` or `import "koda.game"` | Resolve `@name` → `stdlib/name.koda` or builtin object |
| **File import** | `import "src/foo.koda"` | Load module; binding is export object |

Resolution (`internal/parser/loader.go`, `ResolveImportPath`):

1. `@module` → `KODA_WRAPPERS`, `KODA_PATH`, then bundled `stdlib/`, `wrappers/`, project-adjacent paths
2. Plain paths → relative to importer, then `KODA_PATH`, then SDK `wrappers/` / `stdlib/`

**Implemented:** `use raylib;`, `use koda.math;` — see [concepts/modules.md](concepts/modules.md). Legacy `#include` and `import "@x"` remain valid.

**Compatibility rule:** old `use raylib;` keeps working while `use` lands.

---

## Language surface — implemented vs planned

### Implemented and tested

| Feature | Status | Notes |
|---------|--------|-------|
| `let` / reassignment | ✅ | Primary mutable binding |
| `const` | ✅ | `LetDecl.IsConst`; reassignment error in sema — `tests/const_test.koda` |
| Semicolons | ✅ | Required on statements |
| Braces | ✅ | Required on control-flow bodies |
| Type annotations on `let` | ✅ | `let x: int = 3` |
| Structs | ✅ | Ordered fields, O(1) slot access |
| Struct defaults | ✅ | `field = expr` in struct decl — `tests/struct_defaults_ctor_methods.koda` |
| Struct typed fields | ✅ | `health: float = 100.0` — `tests/struct_typed_fields_test.koda` |
| Struct methods | ✅ | `func` inside `struct`; **`this`** for receiver — `tests/struct_methods.koda` |
| Struct literals | ✅ | `Rect { w: 10, h: 5 }` |
| Enums | ✅ | `enum State { A, B }`; members from 0 — parser + sema + match/switch |
| `defer` | ✅ | LIFO at function exit; codegen in `statements.go` |
| `match` / `switch` | ✅ | Pattern dispatch |
| Closures | ✅ | Upvalue capture |
| Integer types | ✅ | Opt-in `i32`, `u8`, … |
| String interpolation | ✅ | `"Score: {n}"` |
| Modules `@math`, `koda.game`, … | ✅ | Stdlib + compiler builtins |
| `// koda:extern` | ✅ | Native argv ABI bindings |
| Graphics projects | ✅ | `koda.json` `"graphics": true`, full `wrappers/raylib/wrapper.c` |

### Partial / different from target syntax

| Target (vision doc) | Today | Gap |
|---------------------|-------|-----|
| `use raylib;` / `use koda.math;` | `#include` / `import "@raylib"` | ✅ `use` landed; `as` / `only` filters future |
| `struct Player { pos: Vec3; }` | `struct Player { x, y, z; }` (untyped fields) | Field type annotations |
| `func hurt(self, …)` | `func hurt() { this.health -= … }` | Explicit `self` param (optional future) |
| `Vec3` type + operators | `@vec3` object helpers | Named type, `+=`, `.to_raylib()` |
| `native library { … }` | `// koda:extern` + wrapgen output | Declarative FFI block |
| `unsafe { … }` | Not present | Pointer deref, manual alloc |
| `Result<T,E>` / `?` | `ok` / `err` builtins only | Typed error propagation |
| `use koda.math` | `@math` / `#include "@math"` | Dotted module paths |
| `koda.toml` | `koda.json` | TOML project file (optional) |

### Not started (documented as future)

- `use module as alias`
- `use module only A, B`
- Owned resources: `own LoadTexture(...) with UnloadTexture`
- High-level `koda.app` / `Application.run(Game {})` (design only)
- Math operator overloading (`pos += vel * dt`)

---

## Type system (today)

- **Runtime default:** most values are **NaN-boxed `i64`** (64-bit float semantics for untyped numbers).
- **Inference:** sema infers types for `let` where possible; explicit annotations for integers.
- **Structs:** compile-time field layout; no separate runtime struct type object for user structs.
- **Enums:** integer constants in a namespace; exhaustiveness **warnings** in strict lint.

Target direction: keep inference for beginners; encourage explicit struct field types in libraries (Phase 5).

---

## Native interop (today)

**ABI:** generated and hand-written wrappers use:

```c
Value symbol(int arg_count, Value* args);
```

Declared from Koda via:

```koda
// koda:extern InitWindow koda_shim_InitWindow 3
let initwindow = 0;
```

Or wrapgen output under `wrappers/raylib/`.

**Linking:** `KODA_NATIVE_SOURCES`, `KODA_LINKFLAGS`, vendored Raylib in `third_party/raylib_static/stage/`, project `koda.json` `"native"` block.

**Two Raylib tiers:**

| Tier | Import | Symbols | Audience |
|------|--------|---------|----------|
| **Full** | SDK `wrappers/raylib/` + `use raylib;` | 548+ functions | **Default** |
| **Shim** | Project `wrappers/raylib_shim/` (legacy) | ~33 functions | `koda setup raylib --shim` only |
| Full | `import "@raylib"` / `wrappers/raylib/` | 548+ via wrapgen | Serious / raw API |

High-level `koda.game` is a **convenience layer**; raw names remain available.

---

## Runtime and GC

- Generational GC, write barrier, shadow stack for compiled roots
- Game-loop hooks: `gcFrameStep`, `gcStats`, arena allocators
- Builtins registered in `internal/codegen/builtin_register.go` — must match `runtime/src`

---

## Tooling (CLI)

| Command | Role |
|---------|------|
| `koda run` | Compile + execute (temp binary) |
| `koda build` | Native executable |
| `koda check` | Sema only; `--warn-unused` |
| `koda test` | `*_test.koda` discovery + run |
| `koda fmt` | Formatter |
| `koda doctor` | SDK / shim drift |
| `koda new` | Project templates (hello, game, graphics, pong) |
| `koda wrap` / `koda setup raylib` | Wrapper install & refresh |
| `koda bundle` | Ship folder + assets |
| `koda disasm` | Print LLVM IR |

---

## Test layout

| Location | Purpose |
|----------|---------|
| `tests/*.koda` | Language + runtime integration |
| `internal/*_test.go` | Unit tests (parser, sema, codegen, formatter) |
| `examples/` | Demos (games, raylib, 3D) |
| `examples/games/` | Full projects with `koda.json` |
| `cmd/koda/main_test.go` | CLI smoke |

Add a new `tests/<feature>_test.koda` for every language change (project rule).

---

## Koda Studio (IDE)

- `koda-ide/` — Wails + Svelte
- Uses `api.CheckSDK`, `api.RunWithWriters`, `api.BuildNativeHost`
- Welcome screen: templates + demo projects (`ExampleGamePath` → `examples/games/…`)

---

## Contributor checklist (before merging language work)

1. Parser + AST + sema + codegen (if needed)
2. `tests/<name>_test.koda` + `go test ./...`
3. Update `docs/language.md` or `docs/status.md` if user-visible
4. Keep existing demos compiling (`examples/games/*`)
5. Do not document syntax that is not implemented — mark **Future** in roadmap docs

---

## Related documents

- [status.md](status.md) — feature matrix for releases
- [ROADMAP.md](ROADMAP.md) — engineering tiers (ASAN, Windows, runtime audit)
- [compiler/architecture.md](compiler/architecture.md) — short LLVM pipeline reference
- [concepts/modules.md](concepts/modules.md) — user-facing import guide (today’s syntax)
