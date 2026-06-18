# Game development with Koda

Koda compiles game logic to a **native binary** — same end result as a C + Raylib project, with faster iteration (`koda run`, `koda watch`) and less boilerplate.

**Complete Raylib reference:** [raylib.md](raylib.md) — all 548 functions, original C names  
**Coming from C:** [from-c.md](from-c.md)

> **`use raylib;` gives you the entire Raylib API.** Names like `InitWindow`, `DrawTexture`, and `CheckCollisionBoxes` work exactly as in C. `koda.game` is optional sugar for 2D loops — it does not hide or limit the underlying wrapper.

---

## Three ways to start

### 1. Text game (no graphics libraries)

Works immediately after install:

```bash
koda new lander --template game
cd lander
koda run
```

### 2. Graphics project (`koda.game` API — optional)

The graphics template links the **full Raylib wrapper** (548 functions). `koda.game` is a thin helper module; you can call any Raylib function directly with `use raylib;`.

```bash
koda new bounce --template graphics
cd bounce
koda doctor
koda run
```

The graphics template sets `"graphics": true` in `koda.json` and links the **full Raylib wrapper** (`wrappers/raylib/wrapper.c`) from the SDK. No manual `KODA_LINKFLAGS` or `KODA_NATIVE_SOURCES` for beginners.

```koda
use raylib;
use koda.game;
```

Koda infers the correct C glue from your imports (`koda_wrap_raylib_*` → full wrapper). Stale shell `KODA_NATIVE_SOURCES` is overridden when it does not match `koda.json`.

### 3. Study the repo examples

| Project | Description |
|---------|-------------|
| `examples/koda-3d/` | **Identity demo** — `Camera3D`, `vec3()`, `RAYWHITE`, defer + loop |
| `examples/games/pong/` | Canonical 2D `koda.game` loop |
| `examples/games/brick-breaker/` | Breakout — enums, hex colors, `match` |
| `examples/games/mario64-studio/` | Optimized 3D — fusion camera, packed colors |
| `examples/spinning-cube/` | Minimal 3D intro |
| `examples/lunar-lander/` | Text lander (no graphics) |

Run **`koda setup raylib`** in older projects missing `koda.json` native settings. Use **`koda setup raylib --shim`** only for legacy shim-based code.

---

## Gold-standard windowed game

Use **`koda.game`** helpers (full Raylib underneath):

```koda
use raylib;
use koda.game;

struct Mario {
    x, y,
    speed,
    health
}

func main() {
    game.open(800, 600, "Koda Game");
    defer game.close();
    game.fps(60);

    let player = Mario {
        x: 400,
        y: 300,
        speed: 220,
        health: 100
    };

    while (game.running()) {
        let dt = game.delta();

        if (game.keyDown(KEY_LEFT)) {
            player.x = player.x - player.speed * dt;
        }
        if (game.keyDown(KEY_RIGHT)) {
            player.x = player.x + player.speed * dt;
        }
        if (game.mouseDown(MOUSE_BUTTON_LEFT)) {
            player.x = game.mouseX();
            player.y = game.mouseY();
        }

        game.begin();
        game.clear(colors.dark);
        game.rect(player.x, player.y, 32, 32, colors.white);
        game.end();
    }

}
```

### `koda.json` for graphics

```json
{
  "name": "mygame",
  "entry": "src/main.koda",
  "lint": "beginner",
  "native": {
    "sources": ["wrappers/raylib/wrapper.c"],
    "graphics": true
  }
}
```

---

## Structs for game state

Structs are the main data model — use **dot notation** like JavaScript:

```koda
struct Mario {
    x, y, speed, health
}

let player = Mario {
    x: 400,
    y: 300,
    speed: 220,
    health: 100
};

player.x = player.x + player.speed * dt;
```

**Naming:** the struct type name is a binding in the same scope — do not reuse it as a variable (`struct Player` + `let player` causes a duplicate-binding error). Use `struct Mario` / `let player` or `struct GoombaNpc` / `let goomba_npc`.

Object literals with methods also work (good for cameras and subsystems):

```koda
let cam = {
    yaw: 0.0,
    dist: 12.0,
    update: func() {
        this.yaw = this.yaw + getmousex() * 0.001;
    }
};
cam.update();
```

Pass structs to functions by reference (fields mutate in place):

```koda
func update(player, dt) {
    player.x = player.x + player.speed * dt;
    if (player.health <= 0) {
        return false;
    }
    return true;
}
```

Use `const` for tuning values:

```koda
const gravity = 900;
const screenWidth = 800;
```

Use **`enum`** and **`match`** for game phases instead of many booleans:

```koda
enum GameState { Playing, Won, GameOver }

let state = GameState.Playing;

match state {
    GameState.Playing {
        update_game(dt);
    }
    GameState.Won {
        draw.text("STAR GET!", 380, 340, 40, colors.yellow);
    }
    GameState.GameOver {
        draw.text("GAME OVER - press R", 400, 340, 36, colors.red);
    }
}
```

Draw HUD text with string interpolation — no manual concatenation:

```koda
draw.text("Score: {score}   Lives: {lives}", 20, 20, 24, colors.white);
```

---

## `koda.game` API reference

| Function | Purpose |
|----------|---------|
| `game.open(w, h, title)` | Open window |
| `game.close()` | Close window |
| `game.running()` | `true` while window is open |
| `game.delta()` | Seconds since last frame |
| `game.fps(n)` | Target frame rate |
| `game.begin()` / `game.end()` | Start/end drawing |
| `game.clear(color)` | Clear background |
| `game.text(msg, x, y, size, color)` | Draw text |
| `draw.text(msg, x, y, size, color)` | Same (alias on `draw` object) |
| `game.rect(x, y, w, h, color)` | Draw rectangle |
| `game.circle(x, y, r, color)` | Draw circle |
| `game.line(x1, y1, x2, y2, color)` | Draw line |
| `game.circleLines(x, y, r, color)` | Circle outline |
| `game.rectLines(x, y, w, h, color)` | Rectangle outline |
| `game.loadImage(path)` / `game.drawImage(tex, x, y, color)` / `game.unloadImage(tex)` | Textures |
| `game.keyDown(key)` | Key held this frame |
| `game.keyPressed(key)` | Key pressed this frame |
| `game.mouseX()` / `game.mouseY()` | Cursor position |
| `game.mouseDown(btn)` / `game.mousePressed(btn)` | Mouse buttons (`Mouse.Left`, …) |
| `game.mouseWheel()` | Scroll delta |
| `game.width()` / `game.height()` | Window size |
| `game.setTitle(title)` | Window title |
| `game.fpsCounter()` | Current FPS |
| `game.setGcBudget(ms)` | Extra incremental GC budget (optional) |
| `game.begin3d(…)` / `game.end3d()` | Fusion-friendly 3D camera pass |
| `game.warmup3d()` | One hidden 3D frame before the main loop |

**Full Raylib API:** `use raylib;` exposes all **548** C functions (`InitWindow`, `DrawCube`, `LoadModel`, …). The table below lists optional `koda.game` shortcuts only.

`Key`, `Mouse`, and `Color` constants are defined in `stdlib/game.koda`. Full detail: [stdlib/game](../stdlib/game.md).

For raw Raylib (3D, textures, audio, anything in the API reference), skip `koda.game` and call Raylib by name — see [raylib.md](raylib.md) and `wrappers/raylib/api_reference.md`.

---

## Build, run, ship

```bash
koda run src/main.koda
koda watch
koda build -o mygame
koda bundle -o dist/mygame
```

Run `koda doctor` before your first graphics build.

---

## Structs in game code

Struct values passed to functions are **shared by reference** (like objects in JavaScript): mutating `bob.x` inside `reset_bobomb(bob, …)` updates the caller’s `bob1` / `bob2`. That is intentional for enemies, coins, and the player — pass the struct, modify fields in place, no `return` needed.

Use `for (let coin of coins)` to iterate arrays; `for coin in coins` (no `let`) is equivalent sugar — see [language.md](../language.md#forin-values-and-ranges).

---

## Quick reference (cheatsheet)

| Need | API |
|------|-----|
| Window + loop | `use raylib;` + optional `use koda.game;` |
| Frame delta | `game.delta()` |
| Input | `game.keyDown(KEY_LEFT)`, `game.mouseX()`, `game.mouseDown(MOUSE_BUTTON_LEFT)` |
| Draw | `game.begin()`, `game.clear()`, `game.rect()`, `game.end()` |
| RNG | `random()`, `randomint(a,b)`, `import "@math"` |
| Timers | `import "@timer"` |
| JSON config | `import "@json"` |
| GC per frame | `game.begin()` / `game.end()` (built-in); raw loops: `gcframestep(1.0)` ×2 |
| Diagnostics | `koda doctor`, `koda check` |

---

## GC and performance

`game.begin()` and `game.end()` each spread ~1 ms of incremental GC work — use them every frame instead of calling `game.setGcBudget()` manually.

**2D draws:** `game.rect` / `game.circle` / `game.clear` / `game.text` with `colors.*` are zero-allocation (packed colors + compiler fusion to `koda_fast_Draw*`). Cache HUD strings; avoid `` `{score}` `` every frame unless the value changed.

**Build:** `koda build` defaults to `-O3`; use `koda run --release` when profiling gameplay.

For **raw Raylib loops**, call `gcframestep(1.0)` at the start and end of each frame. **3D games** should warm up GPU shaders once (`game.warmup3d()`). Use `BeginMode3D(camera3d(...))` for zero-allocation cameras.

Advanced: `arena()` + `arenaReset()` per frame for scratch data — see [runtime and GC](../concepts/runtime-and-gc.md).

---

## Where to go next

| Document | Contents |
|----------|----------|
| [Beginner's guide](../beginners-guide.md) | Full onboarding |
| [stdlib/game](../stdlib/game.md) | `koda.game` module |
| [raylib.md](raylib.md) | Full Raylib API (548 functions) + cheat sheet |
| [wrappers.md](../wrappers.md) | Extending bindings with kodawrap |
| [distribution.md](distribution.md) | Shipping builds |
