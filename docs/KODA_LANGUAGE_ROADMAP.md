# Koda language roadmap

**Identity:** a beginner-friendly but fully capable native language for applications, games, engines, tools, and full C/C++ library wrapper use.

**North-star vision:** [KODA_LANGUAGE_VISION.md](KODA_LANGUAGE_VISION.md) — raylib-first identity, layer build order, full keyword surface. Shipping syntax is documented in [language-cheatsheet.md](reference/language-cheatsheet.md).

```text
Beginner-friendly like Lua/JavaScript.
Structured like C#/Go.
Fast and native like C/C++.
Safer and cleaner than C.
Full wrapper access like a serious systems language.
Raylib-first — raw API always available.
```

**One language. Progressive depth. No crippled ceiling.**

> **Compiler internals:** [CURRENT_COMPILER_ARCHITECTURE.md](CURRENT_COMPILER_ARCHITECTURE.md)  
> **Engineering queue (ASAN, CI, runtime):** [ROADMAP.md](ROADMAP.md)  
> **Shipped today:** [status.md](status.md)

---

## Non-negotiables

These stay regardless of phase:

| Rule | Status |
|------|--------|
| Semicolons (official style) | ✅ Enforced |
| Braces on control flow | ✅ Enforced |
| Full C/C++ wrappers first-class | ✅ wrapgen + `@raylib` |
| High-level helpers optional, never exclusive | ✅ `koda.game` + raw Raylib |
| Native interop + performance path | ✅ LLVM + C runtime |
| No toy-only API surface | ✅ Policy |

---

## Progressive depth model

Users climb the same language:

```text
let / const
  → func
    → struct
      → methods
        → enum
          → modules (use)
            → raw wrappers (InitWindow, …)
              → unsafe (future)
```

Examples at each level must remain valid as users advance — no syntax rewrites required.

---

## Phase map

Work **incrementally**. Do not rewrite the compiler in one pass. Each phase: small PR, tests, doc update, existing demos still compile.

### Phase 1 — Audit ✅ (this document set)

- [x] Map lexer → parser → sema → codegen → runtime
- [x] [CURRENT_COMPILER_ARCHITECTURE.md](CURRENT_COMPILER_ARCHITECTURE.md)
- [x] [KODA_LANGUAGE_ROADMAP.md](KODA_LANGUAGE_ROADMAP.md) (this file)
- [x] [KODA_WRAPPER_SYSTEM.md](KODA_WRAPPER_SYSTEM.md)

### Phase 2 — Module imports (`use`) — **in progress**

**Goal:** official syntax replaces ugly includes without breaking them.

```koda
use raylib;              // → wrappers/raylib/raylib.koda
use koda.math;           // → stdlib/math.koda (@math)
use math;                // shorthand → @math
```

| Task | Status |
|------|--------|
| Parse `use module.path;` | ✅ |
| Resolve wrappers + `koda.*` stdlib | ✅ |
| Keep `#include` + `import "@x"` | ✅ |
| Unknown module → searched paths error | ✅ |
| Tests | ✅ `tests/use_module_test.koda`, `internal/parser/use_module_test.go` |
| `use raylib as rl` | ✅ | Phase 2 |
| `use raylib only A, B` / selective `{ }` | ✅ | Phase 2 |

**Example target (when Phase 2 lands):**

```koda
use raylib;

func main() {
    InitWindow(800, 600, "Hello Koda");
    defer CloseWindow();

    while (!WindowShouldClose()) {
        BeginDrawing();
        ClearBackground(RAYWHITE);
        DrawText("Hello", 20, 20, 30, BLACK);
        EndDrawing();
    }
}
```

*Today:* `use raylib;` is the official import. Legacy `#include` and `import "@raylib"` still work.

### Phase 3 — `const` polish — **done**

| Task | Status |
|------|--------|
| `const` bindings | ✅ |
| Reassignment error: `Cannot assign to constant 'x'.` | ✅ |
| Tests | ✅ `tests/const_test.koda`, `TestSemaConstReassignment` |

### Phase 4 — `defer` — **done** (existing)

LIFO cleanup at function exit; see `tests/` and `docs/language.md`.

### Phase 5 — Typed struct fields + defaults — **done**

```koda
struct Player {
    health: float = 100.0;
    speed: float = 8.0;
}

let player = Player {};
let hurt = Player { health: 50.0 };
```

| Task | Status |
|------|--------|
| `field: Type` syntax | ✅ |
| Defaults + `Player {}` partial literals | ✅ |
| Type check defaults + literal fields | ✅ |
| Untyped `x, y;` still works | ✅ |
| Tests | ✅ `tests/struct_typed_fields_test.koda` |

### Phase 6 — Methods (`self`) — **done**

**Status:** methods use **`this`** or **`self`** (aliases); optional explicit **`self`** first parameter.

| Task | Status |
|------|--------|
| `self` keyword as alias for `this` | ✅ |
| Optional explicit `self` first parameter | ✅ |
| Error outside struct methods | ✅ |
| Document reference semantics | [guides/game-dev.md](guides/game-dev.md) |
| Tests | ✅ `tests/struct_self_test.koda`, `tests/struct_methods.koda` |

```koda
struct Box {
    w, h;
    func area() { return self.w * self.h; }
    func scale(self, f) { self.w *= f; self.h *= f; }
}
```

Vision syntax uses `self`; `this` remains valid for compatibility.

### Phase 7 — Enums

**Status:** **done** (declaration, members, match/switch).

| Remaining | |
|-----------|---|
| Native enum import from wrappers | wrapgen emits Koda enums |
| Docs + beginner examples | [learn/07-structs-and-enums.md](learn/07-structs-and-enums.md) |

### Phase 8 — `koda.math` / vector types

**Status:** **partial** — `@vec3` has typed **`Vec3`** struct + helpers; operators deferred.

| Task | Status |
|------|--------|
| `Vec3` struct type + `zero()` | ✅ `@vec3` |
| `Vector2` struct + `zero()` / `create()` | ✅ `@vec2` (`vec2()` kept as alias) |
| `ColorBytes` struct + `to_raylib()` | ✅ `@color` |
| Nested struct param propagation | ✅ multi-pass refine (`length` → `dot`) |
| Plain + typed args same helper in one scope | ⚠️ use separate tests/calls until call-site specialization |
| `Vec4` / `Mat4` types | Planned |
| Operator overloading (`+=`, `*`) | Later |
| Tests | ✅ `tests/vec3_struct_test.koda`, `tests/vector2_struct_test.koda`, `tests/color_struct_test.koda` |

### Phase 9 — Wrapper registry

**Status:** **partial** — `META.json`, `wrapcatalog`, `koda wrap install`.

| Task | Notes |
|------|-------|
| Standard layout: `wrappers/raylib/{wrapper.koda,META.json,README.md}` | Document in [KODA_WRAPPER_SYSTEM.md](KODA_WRAPPER_SYSTEM.md) |
| Resource metadata (create/destroy pairs) | Extend META schema |
| `use raylib` official entry | Phase 2 + 9 together |
| `koda doctor` drift | Already checks shim |

### Phase 10 — Native declaration format

**Goal (design first):**

```koda
native library raylib {
    link "raylib";
    func InitWindow(width: int, height: int, title: c_string);
    …
}
```

| Task | Notes |
|------|-------|
| Spec document | Map to today’s `// koda:extern` |
| Optional parser stub | Behind flag |
| wrapgen emits `native library` blocks | Long-term |

### Phase 11 — Better errors

| Task | Notes |
|------|-------|
| Module not found → searched paths | Phase 2 |
| Field typo → did you mean | Extend `levenshtein.go` |
| `Cannot assign to constant` | Phase 3 |
| `Cannot call unsafe X outside unsafe block` | Phase 12 precursor |
| Arity errors for wrappers | Partially done |

### Phase 12 — High-level libraries (optional layers)

**Rule:** convenience only; raw API always available.

| Module | Purpose |
|--------|---------|
| `koda.window` | Window + main loop helpers |
| `koda.input` | Already `@input` |
| `koda.assets` | Textures, load/unload + defer |
| `koda.scene` | Entity lists (future) |
| `koda.render` | Batch helpers (future) |

**Future application shape (not implemented):**

```koda
use koda.app;

func main() {
    let app = Application { title: "Arena", width: 1280, height: 720 };
    app.run(Game {});
}
```

Mark as **Future** in docs until code exists.

---

## Future features (explicitly not promised yet)

Document only — do not use in shipping examples:

| Feature | Notes |
|---------|-------|
| `unsafe { … }` | Pointer deref, manual alloc |
| `Result<T,E>` + `?` | Typed errors |
| `own T with destroy` | RAII sugar |
| Optional semicolons | Low priority |
| Implicit `self` | Low priority |
| `use only` import filter | Scaffold when `use` lands |

---

## Example files (to add as phases land)

| File | Phase |
|------|-------|
| `examples/koda-3d/` | **V1 identity target** ✅ |
| `examples/hello_raylib_raw/` | 2 ✅ |
| `examples/spinning-cube/` | 8 (camera helpers) ✅ |
| `examples/math_vec3_demo.koda` | 8 |
| `examples/wrapper_resource_defer.koda` | 4 |

Existing demos stay on **shipped syntax** until phases merge; north-star forms (`Camera3D { }`, `vec3()`) documented as **Future** until Phase 8.

---

## Style rules for implementers

- **No Python** in toolchain
- Small, reviewable diffs
- Tests for every feature
- Update [status.md](status.md) on user-visible changes
- **Do not invent syntax in working examples** — use **Future** sections in docs
- Keep semicolons and braces as official style

---

## Success criteria

**V1 identity target** — this must compile first:

```koda
use raylib;

func main() {
    InitWindow(800, 600, "Koda 3D");
    defer CloseWindow();
    SetTargetFPS(60);

    let camera = {
        position: { x: 4.0, y: 4.0, z: 4.0 },
        target: { x: 0.0, y: 0.0, z: 0.0 },
        up: { x: 0.0, y: 1.0, z: 0.0 },
        fovy: 45.0,
        projection: CAMERA_PERSPECTIVE
    };

    loop {
        if WindowShouldClose() { break; }
        BeginDrawing();
        ClearBackground(#101018);
        BeginMode3D(camera);
        DrawGrid(10, 1.0);
        DrawCube({ x: 0.0, y: 1.0, z: 0.0 }, 2.0, 2.0, 2.0, #663399);
        EndMode3D();
        DrawText("Koda 3D — raylib", 10, 10, 20, #FFFFFF);
        EndDrawing();
    }
}
```

See **`examples/koda-3d/`**. Full north-star syntax (`Camera3D { position: vec3(…) }`, `colors.rebeccaPurple`, `RAYWHITE`) lands in Phases 8–9.

A beginner can write:

```koda
use raylib;

func main() {
    InitWindow(800, 600, "Hello");
    defer CloseWindow();
    …
}
```

An expert can ship a wrapper-heavy native game or tool with full Raylib, custom engine code, and (later) explicit unsafe — **same language**, no ceiling removed.
