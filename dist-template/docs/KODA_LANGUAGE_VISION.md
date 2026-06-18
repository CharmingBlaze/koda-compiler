# Koda language vision

**Status:** design north star · **Shipped today:** [language-cheatsheet.md](reference/language-cheatsheet.md) · **Implementation phases:** [KODA_LANGUAGE_ROADMAP.md](KODA_LANGUAGE_ROADMAP.md)

> **Identity:** Koda is what **raylib would feel like if it had its own modern language** — C-style structure, BASIC-level ease, raylib-first, C-library-wrapper-first, game/app focused, fast native output, beginner-friendly but not toy-like.

```text
Koda is not Python-like.
Koda is not JavaScript-like.
Koda is C/raylib-like, but cleaner.
```

**Name and files:** the language is **Koda** (not a rebrand). Source files use **`.koda`**. Entry is **`func main()`** (not `fn`). The CLI is **`koda`**, not `dd`.

---

## 1. Core identity

Koda is a **native app and game language** built around C library power, raylib-style simplicity, modern syntax, and beginner-friendly tooling.

| Good for | Not trying to be |
|----------|------------------|
| 2D/3D games, desktop apps, tools, editors | Web scripting first |
| Creative coding, small engines, utilities | JVM language |
| C library wrapping, visual-scripting backends | Python replacement |
| | Full C++ clone or toy BASIC |

**Build order (non-negotiable):**

```text
Layer 1  Raw language (lexer, parser, func main, structs, modules)
Layer 2  raylib + C interop + defer
Layer 3  Color / vector / string QoL
Layer 4  Friendly wrappers (input, assets, ui)
Layer 5  Game/app sugar (app / window / update / draw)
Layer 6  Scenes, ECS, editor — later
```

Raw raylib first. C interop second. Color/vector QoL third. raygui UI fourth. Friendly game wrappers fifth. Editor and visual scripting last.

---

## 2. Syntax style

| Use | Avoid |
|-----|-------|
| C-style `{ }` braces | Python indentation |
| Semicolons after statements | BASIC line numbers |
| **`func`** functions, `let` variables | C++ `#include` |
| `use raylib;` modules | Text-paste includes |
| Struct literals with **commas** between fields | Semicolons inside struct literals |
| Raw library names (`InitWindow`, `DrawCube`) | Forced prefixes on every call |

```koda
use raylib;

func main() {
    let x = 100;
    let y = 200;
    DrawCircle(x, y, 40.0, #FF0000);
}
```

**Field vs statement rule:**

```koda
// statements → semicolons
let x = 10;
DrawText("Hi", 10, 10, 20, #FFFFFF);

// struct/object fields → commas
let camera = {
    position: { x: 4.0, y: 4.0, z: 4.0 },
    target: { x: 0.0, y: 0.0, z: 0.0 },
    fovy: 45.0,
    projection: CAMERA_PERSPECTIVE
};
```

**Shipped syntax:** typed `Camera3D { position: vec3(4, 4, 4), … }` via raylib prelude when using `use raylib;`.

**Conditions:** parentheses optional — `if health <= 0 { }` and `if (health <= 0) { }` both valid. **`not`** and **`!`** both work; **`and` / `or` / `not`** are shipped aliases for `&&` / `||` / `!`.

---

## 3. Files, entry point, project

| Item | Koda |
|------|------|
| Language name | **Koda** |
| File extension | **`.koda`** |
| Entry point | **`func main() { }`** |
| Project file | **`koda.json`** |
| CLI | `koda run`, `koda build`, `koda wrap`, `koda test`, `koda fmt` |

```text
MyGame/
  koda.json
  src/main.koda
  assets/
  wrappers/raylib/
```

**Canonical identity program (shipped syntax today):**

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
        if WindowShouldClose() {
            break;
        }
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

See **`examples/koda-3d`** — the reference implementation of this sample.

**Future app sugar (Layer 5 — compiler expands to raw loop):**

```koda
app "My Game";
window { size: 1280, 720; title: "My Game"; fps: 60; }
draw { clear(#101018); }
```

---

## 4. Modules and imports

Official style: **`use raylib;`** — public names enter scope (no prefix).

| Form | Example | Status |
|------|---------|--------|
| Full import | `use raylib;` → `InitWindow(...)` | **shipped** |
| Alias | `use raylib as rl;` → `rl.InitWindow(...)` | **shipped** |
| Selective | `use raylib { InitWindow, DrawText };` | **shipped** |
| Selective | `use raylib { InitWindow, DrawText };` | planned |
| Module-only | `import raylib;` → `raylib.InitWindow(...)` | future |
| Legacy include | `#include "file.koda"` | legacy only |

No `#include <raylib.h>`. Modules are real compilation units, not text pasting.

**Wrapper generation (today):**

```bash
koda wrap raylib.h --module raylib --out wrappers/raylib
```

**Target (Phase 10):**

```koda
extern c {
    link "raylib";
    func InitWindow(width: int, height: int, title: cstring);
}
```

---

## 5. Comments

| Form | Example |
|------|---------|
| Line | `// single-line` |
| Block | `/* block comment */` |
| Doc (future) | `/// Loads a texture from disk.` |

---

## 6. Variables and types

```koda
let score = 0;              // mutable local (default for games)
const gravity = 9.8;        // immutable binding
let speed: float = 4.5;     // explicit type (optional)
```

**Core types:** `bool`, `int`, `float`, `string`, `number` (inferred default).

**Exact-width (opt-in):** `i8`…`i64`, `u8`…`u64`, `f32`, `f64` · aliases: `int` = `i32`, `float` = `f32`.

**Strings:** `"Hello"`, `"Score: {score}"` interpolation. **`cstring`** for C interop (future). Literals auto-convert at C call sites.

**Functions:**

```koda
func add(a: int, b: int) {
    return a + b;
}
```

Koda uses **`func`**, not `fn`. Return type `-> T` syntax is **future**; infer from usage today.

---

## 7. Control flow

| Construct | Example | Status |
|-----------|---------|--------|
| `if` / `else` | `if (health <= 0) { die(); }` | shipped |
| `while` | `while not WindowShouldClose() { }` | shipped |
| `loop` | `loop { if done { break; } }` | shipped |
| Range | `for i in 0..10 { }` → 0–9 | shipped |
| Inclusive range | `for i in 0..=10 { }` → 0–10 | shipped |
| Collection | `for (let enemy of enemies) { }` | shipped |
| `match` | `match state { Phase.Menu { } default { } }` | shipped |
| `break` / `continue` | standard | shipped |

**`loop`** is preferred over `while true { }` for game main loops — explicit infinite loop, exit with `break`.

---

## 8. Operators

| Category | Operators |
|----------|-----------|
| Arithmetic | `+` `-` `*` `/` `%` |
| Assignment | `=` `+=` `-=` `*=` `/=` `%=` |
| Comparison | `==` `!=` `<` `<=` `>` `>=` |
| Logical (C-style) | `!` `&&` `\|\|` |
| Logical (aliases) | `not` `and` `or` |
| Bitwise | `&` `\|` `^` `~` `<<` `>>` |
| Range | `..` (exclusive), `..=` (inclusive) |
| Access | `.` |
| Pointer (future) | `&` `*` |
| Cast (future) | `as` |
| Optional (future) | `?` `?.` `??` |
| Call / group | `()` `[]` `{ }` |
| Punctuation | `;` `,` `:` |

---

## 9. Structs, enums, methods

```koda
struct Player {
    position: float = 0.0;
    health: int = 100;
    speed: float = 5.0;

    func damage(amount) {
        health -= amount;
    }
}

let player = Player {};
player.damage(10);
```

**Enums:**

```koda
enum GameState { Menu, Playing, GameOver }

match state {
    GameState.Playing { update(); }
    default { drawMenu(); }
}
```

Raylib constants (`CAMERA_PERSPECTIVE`, `KEY_W`) work directly from `use raylib;`.

**Namespaces (future):**

```koda
namespace draw {
    func text(value, options) { … }
}
// draw.text("Hello", { position: vec2(10, 10), color: colors.white });
```

---

## 10. Arrays, maps, optionals (future)

| Feature | Target syntax | Status |
|---------|---------------|--------|
| Arrays | `[1, 2, 3]`, `for (let x of arr)` | partial |
| Maps | `map { "wood": 10, "stone": 5 }` | planned |
| Optionals | `let target: Enemy? = none;` `target?.damage(10);` | future |
| Result | `loadSave(path) else { return; }` | future |

**Rule:** `none` = Koda optional empty · `null` = C pointer null.

---

## 11. Memory model

| Layer | Model |
|-------|-------|
| Default | GC for Koda objects |
| C resources | explicit load/unload + **`defer`** |
| Advanced | `arena()`, `Pool<T>`, `gcframestep()` in raw loops |

```koda
let tex = LoadTexture("player.png");
defer UnloadTexture(tex);
```

**Friendly assets (future):** `assets.texture("player.png")` — auto-unload at shutdown.

---

## 12. C interop

**Today:** `// koda:extern` + `koda wrap` + `wrappers/*/wrapper.c`.

**Target:**

```koda
extern c {
    link "raylib";
    func InitWindow(width: int, height: int, title: cstring);
    struct Color { r: u8; g: u8; b: u8; a: u8; }
}
```

**Future types:** `ptr<T>`, `mut ptr<T>`, `voidptr`, `cstring`, `@cstruct` layout.

---

## 13. Colors

Accept colors anywhere Raylib expects `Color`:

```koda
ClearBackground(#101018);
DrawText("Hello", 10, 10, 20, #FFFFFF);
DrawCube({ x: 0, y: 1, z: 0 }, 2, 2, 2, rgb(255, 0, 0));
```

| Form | Status |
|------|--------|
| `#RGB`, `#RGBA`, `#RRGGBB`, `#RRGGBBAA` | **shipped** |
| `rgb()`, `rgba()` | **shipped** |
| `colors.*` namespace | partial (`use koda.color`) |
| `css("rebeccapurple")`, `hsl()` | planned |
| Raylib names `RAYWHITE`, `BLACK`, `RED` | **shipped** (raylib prelude) |

---

## 14. Vectors and math

**Target constructors:**

```koda
vec2(10.0, 20.0)
vec3(0.0, 1.0, 0.0)
```

**Shipped:** global `vec3()` / `vec2()` builtins; `Camera3D` / `Vector3` struct types from raylib prelude.

**Also available:** `use koda.vec3` → `Vec3` struct helpers. Operator overloading (`a + b`) is **later**.

**Functions:** `dot`, `cross`, `length`, `normalize`, `lerp`, `clamp` — via `@vec3` / `@math`.

---

## 15. Game lifecycle layers

**Layer 1 — Raw (canonical, always available):**

```koda
func main() {
    InitWindow(800, 600, "Koda");
    defer CloseWindow();
    loop { … }
}
```

**Layer 2 — Optional helpers:** `koda.game` (`game.open`, `game.delta`, `game.clear(#…)`).

**Layer 5 — Sugar (future):** `app`, `window`, `start`, `update`, `draw` — compiler expands to raw raylib loop. Raw API never removed.

**Layer 6 — Engine (v4+):** `scene`, `entity`, `component`, `system` — only after layers 1–5.

---

## 16. UI, input, audio, assets

| Area | Raw (raylib) | Friendly (future) |
|------|--------------|-------------------|
| Input | `IsKeyDown(KEY_W)` | `input.down(keys.w)` |
| UI | raygui via wrapper | `use koda.ui` |
| Audio | `LoadSound` / `PlaySound` + defer | `assets.sound("jump.wav")` |
| Assets | `LoadTexture` + defer | `assets.texture("player.png")` |

Raw C-first names must always work.

---

## 17. Testing and formatting

```koda
test "hex colors parse" {
    expect(#FF0000 == rgb(255, 0, 0));
}
```

```bash
koda test
koda fmt
```

---

## 18. Compiler errors

Errors must be **human** — file, line, caret, fix hint:

```text
Expected `;` after this statement.

File: src/main.koda
Line: 8

    SetTargetFPS(60)
                    ^
```

---

## 19. Standard library map

| Layer | Modules |
|-------|---------|
| Core std | `@io`, `@math`, `@str`, `@array`, `@json`, `@util` |
| Game | `koda.game`, `@input`, `@color`, `@camera`, `@ui` |
| Bindings | `use raylib;`, raygui (future) |
| Future | `koda.assets`, `koda.window`, `koda.scene` |

---

## 20. Naming conventions

| Kind | Style | Example |
|------|-------|---------|
| Koda files | `snake_case.koda` | `main.koda`, `player.koda` |
| C imports | Original C names | `InitWindow`, `CAMERA_PERSPECTIVE` |
| Friendly API | module helpers | `game.open`, `colors.dark` |
| Structs | PascalCase | `Player`, `Camera3D` |
| Constants | C/Raylib style | `KEY_W`, `MOUSE_BUTTON_LEFT` |

---

## 21. Full keyword index

### Layer 1 — Core (v1)

```text
use  as                    func  return
let  const                 struct  enum
if  else                   while  loop  for  in  break  continue
true  false  null          defer
```

### Layer 2 — C interop (v1 + future)

```text
extern  c  link             ptr  mut  voidptr  cstring
```

### Layer 3 — Safety & tooling (v1.5+)

```text
match  case                 test  assert  expect
none                        map (literal keyword)
```

### Layer 4 — Organization (v2+)

```text
namespace  module  package
```

### Layer 5 — Game/app sugar (v3+)

```text
app  window  start  update  draw  ui
scene  asset  state
```

### Layer 6 — Engine (v4+)

```text
entity  component  system  arena  pool
```

### Future

```text
async  await  try  catch  throw  public  private  export  import
```

---

## 22. Recommended compiler rollout

### V1 (minimum useful language)

```text
use  func main  let  const  defer  loop
int/float/bool/string  if/while/for  struct/enum/match
use raylib  hex colors  // koda:extern + koda wrap
koda.json + koda run/build
```

**V1 target must compile:** `examples/koda-3d`.

### V2

```text
use raylib as rl  selective use { … }
css()/hsl()  namespace  map { }
vec2/vec3 operator overloading
```

### V3

```text
app/window/update/draw sugar
koda.assets  koda.ui (raygui)
extern c { } blocks
```

### V4

```text
entity/component/system  arena/pool  ptr<T>  Result<T>
```

---

## 23. Shipped vs planned (quick reference)

| Vision target | Koda today | Phase |
|---------------|------------|-------|
| `use raylib;` (548 fn) | ✅ | shipped |
| `func main()` | ✅ | shipped |
| `.koda` + `koda.json` | ✅ | shipped |
| `defer`, `const`, `loop` | ✅ | shipped |
| struct + typed fields + methods | ✅ | shipped |
| `enum` + `match` | ✅ | shipped |
| `and` / `or` / `not` | ✅ | shipped |
| `#RRGGBB` inline | ✅ | shipped |
| `while not …` / `!` | ✅ | shipped |
| `extern c { }` | `koda wrap` + `// koda:extern` | Phase 10 |
| `Camera3D { … }` typed literal | ✅ raylib prelude | shipped |
| `vec3()` constructors | ✅ builtin | shipped |
| `RAYWHITE`, `BLACK`, `RED` | ✅ raylib prelude | shipped |
| `colors.rebeccaPurple` | ✅ color prelude | shipped |
| `use raylib as rl` | ✅ | shipped |
| `use raylib { InitWindow, DrawText }` | ✅ | shipped |
| `app` / `window` / `draw` | `koda.game` only | Phase 12 |
| `ptr<T>`, `Result<T>` | GC + `ok`/`err` | v4 |

**Do not** use unimplemented syntax in shipping examples. Mark as **Future** in docs until the compiler accepts it.

---

## 24. Success criteria

A beginner compiles a raw raylib 3D demo with **`use raylib;`** and C API names — no forced wrapper subset. Start from **`examples/koda-3d`**.

An expert ships wrapper-heavy native games with **defer**, **structs**, **match**, and full library access — **same language**, no artificial ceiling.

Friendly helpers (`koda.game`, future `app`/`draw` sugar) are **optional layers**, never the only path to the hardware.
