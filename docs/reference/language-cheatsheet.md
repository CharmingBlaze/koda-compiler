# Koda language cheatsheet

Quick lookup for **Koda today** — not a wishlist. Koda uses **`func`**, not `fn`. There is no ECS `app` / `entity` / `component` syntax; use **`struct`**, **`func main()`**, and optional **`koda.game`** helpers.

**Planned evolution:** [language vision](../KODA_LANGUAGE_VISION.md) — raylib-first layers, full keyword surface, v1 identity demo at `examples/koda-3d`.

Full detail: [language.md](../../language.md) · IDE: **F1 → Language quick reference**

---

## Variables

| Syntax | Example | Notes |
|--------|---------|--------|
| `let` | `let x = 10;` | Mutable local |
| `const` | `const gravity = 9.8;` | Immutable binding |

`var` is reserved — use `let`.

---

## Functions

| Syntax | Example | Notes |
|--------|---------|--------|
| `func` | `func add(a, b) { return a + b; }` | Define a function |
| `return` | `return a + b;` | Return a value |
| default args | `func greet(name = "world") { ... }` | Optional parameters |
| variadic | `func sum(...nums) { ... }` | Rest parameter |

Type annotations are **optional**: `func tick(dt: float) { ... }` or inferred from usage.

There is no `fn` keyword and no `-> type` return syntax — use `func name(...) { ... }`.

---

## Control flow — branches

| Syntax | Example |
|--------|---------|
| `if` | `if (health <= 0) { die(); }` |
| `else` | `} else { survive(); }` |
| `else if` | `} else if (x > 5) { ... }` |
| `switch` | `switch (x) { case 1: break; default: break; }` |
| `match` | `match state { Phase.Playing { ... } default { ... } }` |

`match` arms use blocks `{ ... }`, not `=>`. Expression-style `switch` with `=>` is also supported.

---

## Control flow — loops

| Syntax | Example | Notes |
|--------|---------|--------|
| `while` | `while not WindowShouldClose() { }` | Condition loop |
| `loop` | `loop { if (done) { break; } }` | Explicit infinite loop |
| `for … in` | `for i in 0..10 { }` | Range (exclusive end) |
| `for … in … step` | `for i in 0..100 step 5 { }` | Range with step |
| `for … of` | `for (let enemy of enemies) { }` | Collection loop — **`let` required** |
| `do … while` | `do { ... } while (cond);` | At least once |
| `break` | `break;` | Exit loop / switch |
| `continue` | `continue;` | Next iteration |
| `fallthrough` | `fallthrough;` | Switch fall-through (rare) |

Ranges: `0..10` is 0–9. Use `0..=10` for inclusive end when supported in your build.

---

## Types & structures

| Syntax | Example |
|--------|---------|
| `struct` | `struct Player { x: float = 0.0; health: int = 100; }` |
| struct method | `func move(self, dt) { self.x += speed * dt; }` |
| `enum` | `enum State { Idle, Running, Dead }` |
| struct literal | `let p = Player { x: 10, health: 100 };` |
| raylib 3D camera | `Camera3D { position: vec3(4, 4, 4), target: vec3(0,0,0), up: vec3(0,1,0), fovy: 45.0, projection: CAMERA_PERSPECTIVE }` |
| field access | `p.x`, `this.x` in methods |

No `namespace`, `module`, or `package` keywords — use **`use`**, **`import`**, and files.

---

## Modules & imports

| Syntax | Example | Notes |
|--------|---------|--------|
| `use` | `use raylib;` | **Full Raylib API (548 functions)** |
| `use … as` | `use raylib as rl;` → `rl.InitWindow(...)` | Namespaced import |
| `use { … }` | `use raylib { InitWindow, DrawText };` | Selective import |
| `use` | `use koda.game;` | Optional 2D helpers |
| `import` | `let math = import "@math";` | Stdlib module |
| `#include` | `#include "helpers.koda"` | Paste file at compile time (legacy) |

`use raylib { DrawText };` selective import is **supported**.

---

## Native / C interop

| Syntax | Example |
|--------|---------|
| `// koda:extern` | `// koda:extern InitWindow koda_wrap_raylib_InitWindow 3` |
| `let` binding | `let InitWindow = 0;` |
| `koda.json` | `"native": { "sources": ["wrappers/raylib/wrapper.c"], "graphics": true }` |

No `extern c { }` blocks. Generate wrappers with **`koda wrap`**.

---

## Memory & cleanup

| Syntax | Example | Notes |
|--------|---------|--------|
| `defer` | `defer CloseWindow();` | Run on scope exit |
| `arena()` | `let scratch = arena();` | Scratch allocator |
| `arenaReset` | `arenaReset(scratch);` | Reset arena each frame |
| `gcframestep` | `gcframestep(1.0);` | Spread GC in raw Raylib loops |
| `delete` | `delete obj.field;` | Remove object property |

No `ptr<T>`, `*`, `&`, or `Result<T>` in the language surface today.

---

## Type annotations (optional)

| Type | Example |
|------|---------|
| `int` | `let lives: int = 3;` |
| `float` | `let dt: float = 0.016;` |
| `string` | `let label: string = "Hi";` |
| `bool` | `let alive: bool = true;` |
| `i8` … `i64`, `u8` … `u64` | Exact-width integers |
| `float` / `float32` / `float64` | Floating point |

Types are inferred when omitted — beginners can skip annotations.

---

## Literals & values

| Form | Example |
|------|---------|
| numbers | `42`, `3.14`, `1.0` |
| booleans | `true`, `false` |
| null | `null` |
| strings | `"Score: {score}"` |
| templates | `` `Hello ${name}` `` |
| hex colors | `#101018`, `#FF0000`, `#F00` | CSS-style Raylib color |
| colors (stdlib) | `colors.red`, `colors.rebeccaPurple`, `rgb(255, 0, 0)` | Named / channel colors |
| raylib colors | `RAYWHITE`, `BLACK`, `RED`, `PURPLE` | With `use raylib;` (auto prelude) |
| vectors | `vec3(0, 1, 0)`, `vec2(10, 20)` | Global builtins |
| range | `0..10`, `0..=10` |

---

## Operators

| Category | Operators |
|----------|-----------|
| Arithmetic | `+ - * / % **` |
| Assignment | `= += -= *= /=` |
| Comparison | `== != < <= > >=` |
| Logical | `&& \|\| !` or **`and` `or` `not`** |
| Bitwise | `& \| ^ ~ << >> >>>` |
| Optional coalesce | `??`, `?.` (nullish / optional index only) |
| Grouping | `( )` |

No `?.` method optional chaining on arbitrary calls.

---

## Testing

| Syntax | Example |
|--------|---------|
| `test` | `test "adds" { expect(add(1, 2) == 3); }` |
| `expect` | `expect(x > 0);` | Preferred in tests |
| `assert` | `assert(cond, "msg");` | Runtime assert |

CLI: **`koda test`**, **`koda bench file.koda --count 10`**. No `assert_eq` keyword — use `expect(a == b)`.

---

## Game development (libraries, not keywords)

Koda does **not** have `app`, `window`, `update`, `draw`, `scene`, `entity`, or ECS syntax. Use:

```koda
use raylib;        // all 548 Raylib functions
use koda.game;     // optional shortcuts

func main() {
    game.open(800, 600, "My Game");
    defer game.close();
    while (game.running()) {
        let dt = game.delta();
        game.begin();
        game.clear(colors.dark);
        game.end();
    }
}
```

Or raw Raylib: `InitWindow`, `BeginDrawing`, `DrawCube`, … — see [guides/raylib.md](../guides/raylib.md).

---

## Comments

| Form | Use |
|------|-----|
| `// …` | Line comment |
| `/* … */` | Block comment |

---

## Not in Koda (avoid in new code)

| You may see elsewhere | Use in Koda instead |
|-----------------------|---------------------|
| `fn`, `-> int` | `func` |
| `match x { 1 => ... }` (Rust-style) | `match x { Case.One { ... } }` |
| `extern c { }` | `// koda:extern` + `koda wrap` |
| `ptr<T>`, `Result<T>`, `?.` | GC values, `ok`/`err`, plain calls |
| `app` / `update` / `draw` / ECS | `func main`, `while`, `struct`, `koda.game` |
| `namespace` / `package` | files + `use` / `import` |
