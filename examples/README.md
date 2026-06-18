# Koda examples

Runnable samples for **Koda Studio** and `koda run`. Each folder with a `koda.json` opens as a project.

## Open in Koda Studio

1. Run `Start Koda Studio.cmd` or `.\koda-ide\run-koda-studio.ps1` from the repo root.
2. Pick any item under **Open an example** on the welcome screen.
3. Press **F5** to compile and run.

```powershell
.\koda-ide\run-koda-studio.ps1 examples\games\mario64-studio
```

## Modern patterns (all examples follow these)

**Raylib:** `use raylib;` imports the **complete** wrapper (~548 functions, original C names). `koda.game` is optional — 2D samples use it; 3D samples call `InitWindow`, `DrawCube`, `BeginMode3D`, etc. directly.

| Pattern | 2D (`koda.game`) | 3D (raw Raylib) |
|---------|------------------|-----------------|
| Entry | `func main()` + `game.open` / `defer game.close()` | `func main()` + `InitWindow` / `defer CloseWindow()` |
| Main loop | `loop { if not game.running() { break; } … }` | `loop { if WindowShouldClose() { break; } … }` |
| Clear color | `game.clear(#101018)` | `ClearBackground(#101018)` |
| Colors | `#RRGGBB` hex + `rgb()` where needed | Hoist `#hex` at module level |
| Delta time | `game.delta()` | `clampDelta(GetFrameTime())` via `use koda.util` |
| 3D camera | `game.begin3d(...)` or local `camera3d(...)` | Fusion-friendly `camera3d(px,py,pz, tx,ty,tz, fovy)` |
| GPU warmup | `game.warmup3d()` when using `@game` | One hidden 3D draw + `gcframestep(2.0)` before loop |
| GC in loop | Built into `game.begin()` / `game.end()` | `gcframestep(1.0)` at frame start **and** end |
| State | `enum` + `match` | Same where applicable |
| Syntax | `not`, `and`, `or`, `` `{score}` ``, `for i in 0..N` | Optional parens on `if` / `while` |

**Case-insensitivity:** Koda treats `ballX` and `ballx` as the same name. Brick Breaker uses parallel arrays for that reason.

**Canonical references:** `koda-3d` (identity demo), `games/pong` (2D), `games/mario64-studio` (3D platformer), `spinning-cube` (orbiting camera intro).

## Games

| Project | Description |
|---------|-------------|
| [pong](games/pong) | 2D paddle game — canonical `koda.game` |
| [brick-breaker](games/brick-breaker) | Breakout — `match`, hex colors, inferred types |
| [mario64-studio](games/mario64-studio) | Peach's Castle — canonical optimized 3D |
| [mario64](games/mario64) | Same course, simpler HUD |
| [mario64-hilltop](games/mario64-hilltop) | Bob-omb Hilltop course |
| [koda64](games/koda64) | KODA 64 branding variant |
| [fps-arena](games/fps-arena) | First-person arena shooter |
| [lunar-lander](lunar-lander) | Console text game |

## Graphics demos

| Project | Description |
|---------|-------------|
| [koda-3d](koda-3d) | **Canonical identity demo** — raw raylib, defer + loop, hex colors |
| [spinning-cube](spinning-cube) | Orbiting camera — recommended 3D intro |
| [cube3d](cube3d) | Orbiting camera + satellite cubes |
| [demo-3d](demo-3d) | WASD camera fly-through |
| [raylib-3d-demo](raylib-3d-demo) | WASD + Q/E vertical movement |
| [crystal-plaza](crystal-plaza) | Neon plaza with floating crystal ring |
| [hello_raylib_raw](hello_raylib_raw) | Smallest Raylib window |

## Language

| Project | Description |
|---------|-------------|
| [hello-use-module](hello-use-module) | `use koda.easing` |

All graphics projects use `use raylib;` and `"sources": ["wrappers/raylib/wrapper.c"]` in `koda.json`.
