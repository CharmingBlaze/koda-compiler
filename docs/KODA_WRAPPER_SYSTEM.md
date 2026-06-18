# Koda wrapper system

How Koda binds C/C++ libraries while keeping **raw API names** available and optional high-level helpers on top.

**Principle:**

```text
High-level Koda API = convenience layer.
Raw wrapper API     = always available.
```

> **Pipeline:** [CURRENT_COMPILER_ARCHITECTURE.md](CURRENT_COMPILER_ARCHITECTURE.md)  
> **User guide:** [wrappers.md](wrappers.md), [guides/wrapping-libraries.md](guides/wrapping-libraries.md)  
> **Roadmap:** [KODA_LANGUAGE_ROADMAP.md](KODA_LANGUAGE_ROADMAP.md) Phase 9–10

---

## Today’s stack

```text
C/C++ library (raylib, box2d, …)
  → headers parsed by kodawrap (cmd/wrapgen)
  → wrapper.c (glue) + raylib.koda (bindings)
  → META.json (metadata + header hashes)
  → koda build links via KODA_NATIVE_SOURCES + KODA_LINKFLAGS
  → Koda source calls InitWindow(…) — full Raylib names
```

### Native ABI (all wrappers)

C symbols exposed to Koda use the **argv-style** convention:

```c
Value koda_wrap_raylib_InitWindow(int arg_count, Value* args);
```

Koda declares them with:

```koda
// koda:extern InitWindow koda_wrap_raylib_InitWindow 3
let InitWindow = 0;
```

Or hand-written shims:

```koda
// koda:extern initwindow koda_shim_InitWindow 3
let initwindow = 0;
```

LLVM sees `i64 @symbol(i32, i64*)`. See `internal/codegen/runtime.go`.

---

## Import paths (today)

| Style | Example | Resolves to |
|-------|---------|-------------|
| **`use raylib;`** | Official | `wrappers/raylib/raylib.koda` (548 functions) |
| `import "@raylib"` | Legacy | Same full wrapper |
| `use koda.game;` | Helpers | `stdlib/game.koda` over full Raylib |
| `#include "path.koda"` | Legacy paste | Any `.koda` file |

**Legacy shim** (~33 fn): `koda setup raylib --shim` copies `wrappers/raylib_shim/` into a project.

**Search order** (`internal/parser/loader.go`):

1. `KODA_WRAPPERS` path list  
2. `KODA_PATH`  
3. SDK `stdlib/`, `wrappers/`  
4. Adjacent to importing file  

**Planned (Phase 2):**

```koda
use raylib;           // same as import "@raylib" / wrappers/raylib
use raylib as rl;     // future
use raylib only InitWindow, CloseWindow;  // future
```

Old `#include` remains valid.

---

## Raylib: full wrapper (default)

| Tier | Location | Functions | When to use |
|------|----------|-----------|-------------|
| **Full** | `wrappers/raylib/` from wrapgen | 548+ | **Default** — all projects, templates, `use raylib;` |
| **Shim** | `wrappers/raylib_shim/` (legacy) | ~33 | Old projects only — `koda setup raylib --shim` |

Beginners use **`koda.game`** helpers; experts call raw names in the same file:

```koda
InitWindow(1280, 720, "Game");
DrawText("Hello", 20, 20, 30, WHITE);
```

*Today:* `use raylib;`, `import "@raylib"`, or `#include` of the generated module all work.

---

## Generating wrappers

```bash
koda wrap -name raylib -headers ./raylib.h -I ./include -L ./lib -l raylib -o wrappers/raylib
```

**kodawrap** (`cmd/wrapgen`) emits:

| File | Role |
|------|------|
| `raylib.koda` | Bindings + `// koda:extern` lines |
| `wrapper.c` | C glue, type marshalling |
| `META.json` | Name, headers, hashes, link flags, counts |
| `README.md`, `api_reference.md`, `docs/index.html` | Human + offline API docs |
| `koda.json` snippet | `"native": { "sources": [...], "graphics": true }` |

**Catalog:** `internal/wrapcatalog/catalog.json` — versions, prebuilt paths, install recipes for `koda wrap install raylib`.

**Drift detection:** `internal/wrappermeta` — `koda doctor`, `koda wrap check` compare header hashes in `META.json`.

---

## Target registry layout (Phase 9)

Standard per-library folder:

```text
wrappers/
  raylib/
    raylib.koda      # or module.koda — primary import
    wrapper.c
    META.json
    README.md
    native.json      # future: link settings, resources (optional)
  box2d/
    …
```

### META.json (today)

Machine-readable record (`internal/wrappermeta/meta.go`):

- `name`, `generator`, `version`
- `headers`, `header_hashes` (drift)
- `link_flags`, `include_paths`
- `counts` (functions, structs, …)
- `import` hint (e.g. `@raylib`)

### native.json (future — resource ownership)

Design target for tooling and docs:

```json
{
  "library": "raylib",
  "version": "5.x",
  "resources": [
    {
      "type": "Texture2D",
      "create": "LoadTexture",
      "destroy": "UnloadTexture"
    },
    {
      "type": "Image",
      "create": "LoadImage",
      "destroy": "UnloadImage"
    }
  ],
  "unsafe_functions": ["MemAlloc", "MemFree"]
}
```

Enables:

- `defer UnloadTexture(tex);` guidance  
- Future `own LoadTexture(path) with UnloadTexture`  
- Doc generators listing cleanup pairs  

Not required for manual `defer` today.

---

## Linking and projects

**Project manifest** (`koda.json`):

```json
{
  "name": "arena",
  "entry": "src/main.koda",
  "native": {
    "sources": ["wrappers/raylib/wrapper.c"],
    "graphics": true
  }
}
```

`"graphics": true` → vendored Raylib static lib + platform flags (`internal/nativebuild/raylib_vendored.go`).

**Environment overrides:**

| Variable | Purpose |
|----------|---------|
| `KODA_NATIVE_SOURCES` | Extra `.c` glue |
| `KODA_LINKFLAGS` | `-lraylib`, frameworks, `-L` paths |
| `KODA_WRAPPERS` | Wrapper search root |
| `KODA_HOME` | SDK root (stdlib, wrappers) |

**CLI:**

```bash
koda run          # compile + run entry from koda.json
koda build        # native exe
koda setup raylib # refresh project shim from SDK
koda doctor       # SDK + shim drift
koda wrap install raylib --project
```

Target project file (future): TOML with `[wrappers] raylib = true` — optional; `koda.json` remains supported.

---

## High-level modules (must not replace raw API)

| Module | Layer | Raw still available |
|--------|-------|---------------------|
| `koda.game` / `stdlib/game.koda` | Window loop, draw helpers | Yes — `@raylib` / shim |
| `@input` | Key constants + helpers | Yes — `IsKeyDown` in full wrapper |
| `@camera` | Orbit / FPS camera structs | Yes — `BeginMode3D` |
| `@color` | Named colors + `rgb()` | Yes — Raylib `Color` in full wrapper |

**Not allowed:** shipping only five toy functions with no path to full library.

---

## Native declaration format (Phase 10 — design)

Target syntax for hand-written or generated blocks:

```koda
native library raylib {
    link "raylib";

    struct Vector3 {
        x: float;
        y: float;
        z: float;
    }

    const RAYWHITE: Color;

    func InitWindow(width: int, height: int, title: c_string);
    func CloseWindow();
}
```

**Interop types (target):** `c_int`, `c_string`, `pointer<T>`, `f32`, `u8`, …

**Migration path:** wrapgen continues emitting `// koda:extern`; parser optionally accepts `native library` as sugar that lowers to the same LLVM declarations.

---

## Unsafe wrappers (future)

Rules:

- Normal Koda code: safe defaults, GC, no raw pointer deref  
- `unsafe { … }` required for: pointer deref, manual alloc/free, calling marked-unsafe natives  
- High-level wrappers may hide unsafe inside safe functions  

Example (future):

```koda
unsafe native func MemFree(ptr: void_ptr);

unsafe {
    let ptr = MemAlloc(1024);
    MemFree(ptr);
}
```

Mark unsafe functions in `native.json` / META for codegen + sema.

---

## Wrapper-first examples

### Raw (target style after `use` — today use `@raylib` or include)

```koda
// FUTURE: use raylib;
// TODAY:
import "@raylib";

func main() {
    InitWindow(1280, 720, "Raw Raylib");
    defer CloseWindow();
    SetTargetFPS(60);
    while (!WindowShouldClose()) {
        BeginDrawing();
        ClearBackground(RAYWHITE);
        DrawText("Full Raylib access", 20, 20, 30, WHITE);
        EndDrawing();
    }
}
```

### Beginner + high-level (today)

```koda
use koda.game;

func main() {
    game.open(800, 600, "Hello");
    while (game.running()) {
        game.begin();
        game.clear(colors.dark);
        game.text("Hello Koda", 20, 20, 24, colors.white);
        game.end();
    }
    game.close();
}
```

Both are valid Koda. Same compiler. Same binary pipeline.

---

## Testing requirements (wrappers)

| Test | Purpose |
|------|---------|
| `tests/stdlib_modules_test.koda` | `koda.game`, `@math`, … load |
| Native bind / codegen tests | `internal/codegen/native_bind_debug_test.go` |
| Project shim drift | `internal/project/shim_drift_test.go` |
| `koda doctor` in CI | SDK layout |

Add when Phase 2 lands:

- `tests/use_raylib_test.koda` — `use raylib` resolves and compiles  
- `tests/use_unknown_module_test.koda` — error lists search paths  

---

## Checklist for new library support

1. Run `koda wrap` against headers → `wrappers/<name>/`  
2. Commit `META.json` + document `KODA_LINKFLAGS` per OS  
3. Add to `wrapcatalog/catalog.json` if installable via `koda wrap install`  
4. Optional: beginner shim subset (like `raylib_shim`) — **never** the only export  
5. Document in `docs/guides/` with raw + helper examples  
6. Register `use <name>` mapping when Phase 2 ships  

---

## Related

- [guides/raylib.md](guides/raylib.md) — Raylib walkthrough  
- [guides/wrapping-libraries.md](guides/wrapping-libraries.md) — wrap any C library  
- [concepts/project-layout.md](concepts/project-layout.md) — `koda.json`, shims  
