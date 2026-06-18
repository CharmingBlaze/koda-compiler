# Koda handoff — language, examples, and continuation guide

**Date:** June 2025 · **Audience:** Next engineer, contributor, or AI session picking up this repo.

This document is the **single handoff** for the current Koda product state: language identity, shipped syntax, example catalog, compiler map, verification, and what to build next.

**Related docs:**

| Doc | Purpose |
|-----|---------|
| [KODA_LANGUAGE_VISION.md](KODA_LANGUAGE_VISION.md) | North-star design (raylib-first layers) |
| [KODA_LANGUAGE_ROADMAP.md](KODA_LANGUAGE_ROADMAP.md) | Implementation phases |
| [reference/language-cheatsheet.md](reference/language-cheatsheet.md) | Accurate syntax today |
| [handoff.md](handoff.md) | Compiler/runtime engineering checklist |
| [status.md](status.md) | Feature matrix |

---

## 1. Product identity (do not rebrand)

**Koda** is a native app and game language — C-style structure, BASIC-level ease, **raylib-first**, C-library-wrapper-first, fast native output, beginner-friendly but not toy-like.

```text
Koda is not Python-like.
Koda is not JavaScript-like.
Koda is C/raylib-like, but cleaner.
```

- Language name: **Koda** (not dataDream)
- File extension: **`.koda`**
- Entry point: **`func main()`** (not `fn`)
- CLI: **`koda run`**, **`koda build`**, **`koda check`**, **`koda fmt`**
- Project file: **`koda.json`**

---

## 2. Canonical code patterns (all examples follow these)

### Raw 3D (identity demo)

See **`examples/koda-3d/`** — the reference program:

```koda
use raylib;

func main() {
    InitWindow(800, 600, "Koda 3D");
    defer CloseWindow();
    SetTargetFPS(60);

    let camera = Camera3D {
        position: vec3(4.0, 4.0, 4.0),
        target: vec3(0.0, 0.0, 0.0),
        up: vec3(0.0, 1.0, 0.0),
        fovy: 45.0,
        projection: CAMERA_PERSPECTIVE
    };

    loop {
        if WindowShouldClose() {
            break;
        }
        BeginDrawing();
        ClearBackground(#101018);
        BeginMode3D(camera);
        DrawGrid(10, 1.0);
        DrawCube(vec3(0.0, 1.0, 0.0), 2.0, 2.0, 2.0, colors.rebeccaPurple);
        EndMode3D();
        DrawText("Koda 3D - raylib", 10, 10, 20, RAYWHITE);
        EndDrawing();
    }
}
```

### 2D games (`koda.game` optional layer)

See **`examples/games/pong/`**:

```koda
use raylib;
use koda.game;

func main() {
    game.open(800, 600, "Pong");
    defer game.close();
    game.fps(60);

    loop {
        if not game.running() {
            break;
        }
        game.begin();
        game.clear(#101018);
        // …
        game.end();
    }
}
```

### Pattern table

| Concern | 2D (`koda.game`) | 3D (raw Raylib) |
|---------|------------------|-----------------|
| Entry | `func main()` + `game.open` / `defer game.close()` | `func main()` + `InitWindow` / `defer CloseWindow()` |
| Main loop | `loop { if not game.running() { break; } }` | `loop { if WindowShouldClose() { break; } }` |
| Clear | `game.clear(#101018)` | `ClearBackground(#101018)` |
| Colors | `#RRGGBB`, `colors.*` | Same + `RAYWHITE`, prelude palette |
| Collection loop | `for (let x of arr)` | Same |
| Range loop | `for i in 0..10` | Same |
| GC in raw loop | Built into `game.begin/end` | `gcframestep(1.0)` start + end |
| Camera | `koda.camera` helpers | `camera3d()` fusion or `Camera3D { }` |

### Import styles (all shipped)

```koda
use raylib;                              // full API in scope
use raylib as rl;                        // rl.InitWindow(...)
use raylib { InitWindow, DrawText };     // selective
use koda.math;                           // stdlib namespace
```

---

## 3. Example catalog (16 projects)

All under **`examples/`** with **`koda.json`**. Verified with **`koda check`** from project directory (`KODA_HOME` = repo root).

### Games (`examples/games/`)

| Project | Style | Notes |
|---------|-------|-------|
| **pong** | `koda.game` | Canonical 2D — `loop`, `match`, hex colors |
| **brick-breaker** | `koda.game` | Breakout — parallel arrays, `#hex` |
| **mario64-studio** | Raw Raylib + helpers | Canonical 3D platformer |
| **mario64** | Raw Raylib | Peach's Castle, simpler HUD |
| **mario64-hilltop** | Raw Raylib | Bob-omb Hilltop |
| **koda64** | Raw Raylib | Branding variant |
| **fps-arena** | Raw Raylib + FP cam | FPS gallery |

### Graphics demos

| Project | Notes |
|---------|-------|
| **koda-3d** | **Identity demo** — north-star syntax |
| **spinning-cube** | Recommended 3D intro |
| **cube3d** | Orbiting camera + satellites |
| **demo-3d** | WASD fly-through |
| **raylib-3d-demo** | WASD + Q/E vertical |
| **crystal-plaza** | Neon plaza |
| **hello_raylib_raw** | Smallest window |

### Other

| Project | Notes |
|---------|-------|
| **lunar-lander** | Console text game |
| **hello-use-module** | `use koda.easing` |

### Project templates (`koda new`)

| Template | Path | Output |
|----------|------|--------|
| `hello` | `internal/project/templates/hello/` | Console hello |
| `game` | `internal/project/templates/game/` | Text lunar lander |
| `graphics` | `internal/project/templates/graphics/` | Bouncing ball (`koda.game`) |
| `pong` | `internal/project/templates/pong/` | Full pong |
| `raylib` | `internal/project/templates/raylib/` | Raw Raylib window |

All templates use **`loop`** + **`defer`** where applicable.

---

## 4. Language surface shipped (quick reference)

| Feature | Status |
|---------|--------|
| `func main()`, `let`, `const`, `defer` | ✅ |
| `loop`, `while`, `for`, `match`, `enum`, `struct` + methods | ✅ |
| `use raylib`, `use as`, `use { selective }` | ✅ |
| `#RRGGBB` hex colors | ✅ |
| `Camera3D { }`, `vec3()`, `RAYWHITE` (raylib prelude) | ✅ |
| `colors.rebeccaPurple` (color prelude) | ✅ |
| `and` / `or` / `not` aliases | ✅ |
| `koda.game`, `koda.camera`, `koda.ui`, `@math`, … | ✅ |
| Full Raylib wrapper (548 fn) | ✅ |
| `app` / `window` / `draw` sugar | ❌ planned Phase 12 |
| `extern c { }` blocks | ❌ planned Phase 10 |
| `css()` / `hsl()` colors | ❌ planned V2 |

Full detail: [reference/language-cheatsheet.md](reference/language-cheatsheet.md).

---

## 5. Compiler architecture (where code lives)

```text
cmd/koda/           CLI entry
internal/lexer/     Tokens (#hex, loop keyword, …)
internal/parser/    AST, use/import flatten, preludes (math, colors, raylib types)
internal/sema/      PrepareNativeBundle, types, struct layout
internal/codegen/   LLVM IR, native fusion (DrawCube, BeginMode3D, …)
internal/nativebuild/ llc + clang link
runtime/src/        koda_runtime.c — GC, builtins, vec3 game types
wrappers/raylib/    Full Raylib binding + fast_paths.c
stdlib/             koda.game, koda.camera, @math, …
```

**Preludes** (auto-injected in `internal/parser/prelude.go` + `raylib_types.koda`):

- `let math = { … }` — always
- `let colors = { … }` — always
- `struct Camera3D`, `RAYWHITE`, … — when `use raylib` in any form

**Import flattening** (`internal/parser/include_flatten.go`, `use_selective.go`, `use_alias.go`):

- `use raylib;` → inline all module decls
- `use raylib as rl;` → hidden bindings + `let rl = { … }`
- `use raylib { A, B };` → filter to named exports

---

## 6. Verification commands

From repo root (contributors need Go; end users use release binary):

```powershell
$env:KODA_HOME = "C:\path\to\Koda"

# Unit tests
go test ./internal/...

# Check every example project
Get-ChildItem -Recurse examples -Filter koda.json | ForEach-Object {
  Push-Location $_.DirectoryName
  koda check src/main.koda
  Pop-Location
}

# Run identity demo
koda run examples/koda-3d/src/main.koda

# Run canonical 2D
koda run examples/games/pong/src/main.koda
```

Language tests: `tests/*.koda` (run via `go test` harness or `koda test` when runtime is built).

---

## 7. Known gotchas (examples + docs)

1. **Case-insensitivity** — `ballX` and `ballx` are the same; brick-breaker uses parallel arrays.
2. **For-of** — must write `for (let item of arr)`, not `for item of arr`.
3. **String interpolation** — avoid `-` touching identifiers in strings (e.g. use `"Right mouse"` not `"Right-mouse"` which can parse as `{mouse}`).
4. **Mouse constants** — use `MOUSE_BUTTON_LEFT`, not `Mouse.Left`.
5. **Match arms** — one case per arm; no comma-combined cases.
6. **Raylib prelude** — `Camera3D` / `RAYWHITE` are global even with `use raylib as rl;` (module funcs are namespaced, types/colors are prelude).
7. **Custom `func camera3d()`** — skip raylib types prelude if a top-level `camera3d` function exists (name clash).
8. **DrawText / window titles** — Raylib’s default font is ASCII-ish; use `-` not `—` in `DrawText`, `InitWindow`, and `game.text` strings (otherwise you get `?` glyphs).

---

## 8. What to build next (priority order)

### Phase 12 — App sugar (high user value)

Compiler expands to raw raylib loop:

```koda
app "My Game";
window { size: 800, 600; title: "My Game"; fps: 60; }
draw { clear(#101018); }
```

Never replace raw `func main()` path.

### V2 QoL

- `css("rebeccapurple")`, `hsl()` color helpers
- `vec3` operator overloading (`a + b`)
- Merge prelude colors into `rl.*` for aliased imports

### Phase 10 — `extern c { }`

Replace `// koda:extern` + wrapgen with first-class blocks; wrapgen emits them.

### Phase 11 — Better errors

Module-not-found paths, struct field typos, color parse hints.

---

## 9. Files changed in the examples/docs pass

| Area | Action |
|------|--------|
| `examples/**/src/main.koda` (16 projects) | Modernized: `loop`, `#hex`, `use raylib`, struct methods |
| `examples/koda-3d/` | Added as identity demo |
| `internal/project/templates/*` | Updated to `loop` + modern patterns |
| `wrappers/raylib/types.koda` | Raylib struct types + palette (mirrors prelude) |
| `internal/parser/prelude.go` | math, colors, raylib preludes |
| `internal/parser/use_alias.go` | `use as` |
| `internal/parser/use_selective.go` | `use { }` |
| `docs/KODA_LANGUAGE_VISION.md` | Full language map |
| `docs/status.md`, `examples/README.md` | Updated pattern tables |
| `README.md` | Hero sample uses `loop` + hex |

---

## 10. Session start checklist

1. Read this file + [language-cheatsheet.md](reference/language-cheatsheet.md).
2. Open **`examples/koda-3d`** and **`examples/games/pong`** as reference.
3. Run `go test ./internal/...` and example `koda check` batch.
4. Pick next item from §8; update [KODA_LANGUAGE_ROADMAP.md](KODA_LANGUAGE_ROADMAP.md) when landing.
5. For compiler/runtime work, also read [handoff.md](handoff.md) invariants.

---

## 11. One-line summary

**Koda today:** a native raylib-first language with `func main()`, `use raylib`, `loop`/`defer`/`match`, hex colors, full C API access, optional `koda.game` — examples and templates all demonstrate the same patterns; next big win is app/window/draw sugar or color/vector QoL.
