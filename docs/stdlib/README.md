# Koda standard library

Modules live in `stdlib/` beside the `koda` binary. Load them with **`import "@name"`** or **`#include "stdlib/name.koda"`**.

Built-in `@math`, `@json`, and `@io` objects are optimized in the compiler; file-based modules re-export helpers.

---

## Module index

| Import | File | Purpose |
|--------|------|---------|
| `@math` | [math.koda](math.md) | Trigonometry, lerp, clamp, RNG |
| `@json` | [json.koda](json.md) | parse, stringify, try_parse |
| `@io` | [io.koda](io.md) | Files and directories |
| `@array` | [array.koda](array.md) | range, shuffle, zip, sum |
| `@timer` | [timer.koda](timer.md) | Cooldowns, intervals, countdowns |
| `@vec2` | [vec2.koda](vec2.md) | 2D vectors |
| `@vec3` | [vec3.koda](vec3.md) | 3D vectors |
| `@util` | [util.koda](util.md) | clamp01, pick_weighted, pingpong |
| `@noise` | [noise.koda](noise.md) | 1D value noise |
| `@str` | [str.koda](str.md) | String helper aliases |
| `@game` | [game.koda](game.md) | Beginner game loop API over Raylib (`game.*` and `draw.*`) |
| `@camera` | [camera.koda](camera.md) | Orbit camera for 3D games |
| `@input` | [input.koda](input.md) | Keyboard/mouse helpers (requires Raylib shim) |
| `@color` | [color.koda](color.md) | RGBA/HSV color utilities (pure Koda) |
| *(built-in)* | [game-types.md](game-types.md) | `vec2`, `vec3`, `color`, `rect`, `box` |
| `@easing` | [easing.koda](easing.md) | Animation easing curves |
| `@pool` | [pool.koda](pool.md) | Object pool for hot loops |

**Graphics note:** For windowed games use `@game` (not `@input` alone). `@input` is for lower-level input when you call shim functions directly.

---

## Usage patterns

**Import object (recommended):**

```koda
let math = import "@math";
let x = math.lerp(0, 100, 0.5);
```

**Include (merges into current file):**

```koda
#include "stdlib/timer.koda"
let cd = cooldown(0.5);
```

**Global builtins** — many stdlib functions exist as globals too (`sqrt`, `readfile`, `randomint`). Imports group them under one namespace.

---

## Dot notation (JavaScript-style)

Koda supports field and method access like JavaScript:

```koda
player.x = player.x + speed * dt;
game.circle(ball.x, ball.y, r, colors.white);

let m = import "@math";
m.sin(m.pi / 2);

let cam = {
    yaw: 0.0,
    update: func() { this.yaw = this.yaw + 0.01; }
};
cam.update();
```

Optional chaining: `obj?.field`. Struct type names are case-insensitive with variable names — `struct Player` conflicts with `let player`.

The compiler routes `math.*`, `json.*`, and `io.*` on identifier bindings to native argv functions.

---

## Related

- [Builtins and globals](reference/builtins.md)
- [Language reference](../../language.md)
- [Beginner's guide](../beginners-guide.md)
