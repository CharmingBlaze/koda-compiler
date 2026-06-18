# `koda.game`

Beginner-friendly game loop API over the **full Raylib wrapper** (548 functions). **`use raylib;` gives you every Raylib command** — `InitWindow`, `LoadModel`, `PlaySound`, etc. This module only adds optional shortcuts; it does not limit the underlying API.

## Setup

```koda
use raylib;
use koda.game;
```

```json
{
  "native": {
    "sources": ["wrappers/raylib/wrapper.c"],
    "graphics": true
  }
}
```

Run **`koda setup raylib`** in the project folder if `koda.json` is missing native settings. Use **`koda setup raylib --shim`** only for legacy ~33-function projects.

## Minimal loop

```koda
use raylib;
use koda.game;

func main() {
    game.open(800, 600, "My Game");
    defer game.close();
    game.fps(60);

    while (game.running()) {
        let dt = game.delta() or 0.016;
        game.begin();
        game.clear(colors.dark);
        game.text("Hello", 20, 20, 24, colors.white);
        game.end();
    }
}
```

`game.clear`, `game.text`, and other draw helpers accept prelude `colors.*` (packed `rgb()` values) or Raylib `{r,g,b,a}` objects. Packed colors pass through with zero allocation; the compiler fuses draw calls to native fast paths.

## API

| Member | Raylib call | Notes |
|--------|-------------|-------|
| `game.open(w, h, title)` | `InitWindow` | Open window |
| `game.close()` | `CloseWindow` | Close window |
| `game.running()` | `!WindowShouldClose()` | Main-loop condition |
| `game.delta()` | `GetFrameTime()` | Frame delta (seconds); clamps 0/spikes to ~16 ms |
| `game.fps(n)` | `SetTargetFPS` | Target FPS |
| `game.begin()` | `BeginDrawing` + `gcframestep` | Start 2D frame; spreads GC (~1 ms) |
| `game.end()` | `EndDrawing` + `gcframestep` | End 2D frame; spreads GC (~1 ms) |
| `game.begin3d(px…fovy)` | `BeginMode3D(camera3d(…))` | Fusion-friendly 3D camera |
| `game.end3d()` | `EndMode3D` | End 3D pass |
| `game.warmup3d()` | hidden 3D frame | Shader compile before main loop |
| `game.clear(color)` | `ClearBackground` | Background color |
| `game.text(msg, x, y, size, color)` | `DrawText` | Text |
| `game.rect(x, y, w, h, color)` | `DrawRectangle` | Filled rect |
| `game.circle(x, y, r, color)` | `DrawCircle` | Filled circle |
| `game.circleLines(x, y, r, color)` | `DrawCircleLines` | Circle outline |
| `game.rectLines(x, y, w, h, color)` | `DrawRectangleLines` | Rect outline |
| `game.line(x1, y1, x2, y2, color)` | `DrawLine` | Line |
| `game.keyDown(key)` | `IsKeyDown` | Held key |
| `game.keyPressed(key)` | `IsKeyPressed` | Key just pressed |
| `game.mouseX()` / `game.mouseY()` | `GetMousePosition` | Cursor |
| `game.setGcBudget(ms)` | `gcframestep` | Extra GC budget (usually unnecessary with `begin`/`end`) |

## Performance

`game.begin()` and `game.end()` each call `gcframestep(1.0)` so incremental GC work is spread across the frame instead of stalling in one big pause. Use them every frame in your main loop.

**Draw calls:** `game.rect`, `game.circle`, `game.circleLines`, `game.rectLines`, `game.line`, `game.clear`, and `game.text` pass packed `colors.*` straight through — no per-call color object allocation. The compiler fuses these to native fast paths (`koda_fast_Draw*`) when linked with the full Raylib wrapper.

**Rules of thumb:**

| Scale | Pattern |
|-------|---------|
| Any game | `game.begin()` / `game.end()` every frame |
| Many draws | Use `colors.*` or hoisted packed colors (not `{r,g,b,a}` literals in loops) |
| Physics / movement | `let vx: float = 0.0` — unboxed stack floats in hot loops |
| HUD text | Cache strings; update only when values change |
| Heavy temp data | `let scratch = arena(65536)` + `arenaReset(scratch)` each frame |
| Ship build | `koda build` (defaults to `-O3`) or `koda run --release` |

For **raw Raylib loops** (no `@game`), call `gcframestep(1.0)` at the start and end of each frame yourself.

3D projects should call `game.warmup3d()` once after `game.open()`. Use `BeginMode3D(camera3d(px, py, pz, tx, ty, tz, fovy))` with seven numeric args for a zero-allocation camera path.

## Input and colors

Use Raylib **`KEY_*`** and **`MOUSE_BUTTON_*`** constants from `use raylib` (reliable in all builds):

```koda
if (game.keyDown(KEY_W)) { ... }
if (game.keyPressed(KEY_SPACE)) { ... }
if (game.mouseDown(MOUSE_BUTTON_LEFT)) { ... }
game.clear(colors.sky);
```

`Key.W` / `Mouse.Left` aliases exist in `koda.game` for readability but **`KEY_*` is preferred** in examples and shipped games.

## Raw API

Nothing blocks calling Raylib directly in the same file:

```koda
use raylib;
use koda.game;

func main() {
    game.open(640, 480, "Both APIs");
    defer game.close();
    while (game.running()) {
        game.begin();
        DrawFPS(10, 10);
        game.end();
    }
}
```

## See also

- [Raylib guide](../guides/raylib.md) — full API reference
- [Game development](../guides/game-dev.md) — loops, input, shipping
- [Modules](../concepts/modules.md) — `use raylib`, `use koda.game`
