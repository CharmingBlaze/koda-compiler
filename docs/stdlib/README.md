# Stdlib modules

Koda ships `.koda` modules under `stdlib/`. Import with **`use koda.name;`** (recommended) or legacy **`import "@name"`** / **`#include "@name"`**.

## Graphics (full Raylib)

**`use raylib;`** loads the complete generated wrapper (**548 functions**, same names as C). Nothing is hidden behind `koda.game` — helpers are optional.

Load **`use raylib;`** first, then optional modules:

| Module | Doc | Purpose |
|--------|-----|---------|
| `koda.game` | [game.md](game.md) | Window loop, draw helpers |
| `koda.input` | [input.md](input.md) | Keyboard/mouse |
| `koda.camera` | [camera.md](camera.md) | 3D orbit / FPS camera |
| `koda.color` | [color.md](color.md) | `asRaylib()` for draw colors |
| `koda.ui` | [ui.md](ui.md) | Score/life pip HUD widgets |

```koda
use raylib;
use koda.game;
```

## General

| Module | Doc |
|--------|-----|
| `@math` / `koda.math` | [math.md](math.md) |
| `@json` / `koda.json` | [json.md](json.md) |
| `@io` / `koda.io` | [io.md](io.md) |
| `@array` / `koda.array` | [array.md](array.md) |
| `@vec2` / `koda.vec2` | [vec2.md](vec2.md) |
| `@vec3` / `koda.vec3` | [vec3.md](vec3.md) |
| `@timer` / `koda.timer` | [timer.md](timer.md) |
| `@easing` / `koda.easing` | [easing.md](easing.md) |
| `@noise` / `koda.noise` | [noise.md](noise.md) |

**Graphics note:** `koda.game` is a thin 2D helper layer. **All** Raylib commands remain available in the same file via `use raylib;` — see `wrappers/raylib/api_reference.md`.
