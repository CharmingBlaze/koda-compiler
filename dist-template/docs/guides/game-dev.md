# Game development with Koda

Koda compiles game logic to a **native binary** — same end result as a C + Raylib project, with faster iteration (`koda run`, `koda watch`) and less boilerplate.

**Complete Raylib reference:** [raylib.md](raylib.md)  
**Coming from C:** [from-c.md](from-c.md)

---

## Three ways to start

### 1. Text game (no graphics libraries)

Works immediately after install:

```bash
koda new lander --template game
cd lander
koda run
```

### 2. Graphics project (`@game` API)

```bash
koda new bounce --template graphics
cd bounce
koda doctor
koda run
```

The graphics template sets `"graphics": true` in `koda.json` — platform link flags and the Raylib shim wrapper are applied automatically. No manual `KODA_LINKFLAGS` or `KODA_NATIVE_SOURCES` for beginners.

Koda also infers the correct C glue from your `#include` (`koda_shim_*` → shim, `koda_wrap_raylib_*` → full wrapper), so stale shell environment variables are overridden when they do not match.

### 3. Study the repo examples

| File | Description |
|------|-------------|
| `examples/games/koda64/` | Mario 64-style 3D platformer (orbit camera, structs, dot notation) |
| `examples/games/brick_breaker.koda` | Full brick breaker |
| `examples/raylib_shim_demo.koda` | 3D camera + cube |
| `examples/games/lunar_lander_text.koda` | Text lander (standalone file) |

Refresh an older project's shim with **`koda setup raylib`** before using `@game`.

---

## Gold-standard windowed game

Use the `@game` wrapper — not raw Raylib names:

```koda
#include "wrappers/raylib_shim/raylib.koda"
#include "@game"

struct Mario {
    x, y,
    speed,
    health
}

func main() {
    game.open(800, 600, "Koda Game");
    game.fps(60);

    let player = Mario {
        x: 400,
        y: 300,
        speed: 220,
        health: 100
    };

    while (game.running()) {
        let dt = game.delta();

        if (game.keyDown(Key.Left)) {
            player.x = player.x - player.speed * dt;
        }
        if (game.keyDown(Key.Right)) {
            player.x = player.x + player.speed * dt;
        }
        if (game.mouseDown(Mouse.Left)) {
            player.x = game.mouseX();
            player.y = game.mouseY();
        }

        game.begin();
        game.clear(colors.dark);
        game.rect(player.x, player.y, 32, 32, colors.white);
        game.end();
        game.setGcBudget(0.5);
    }

    game.close();
}
```

### `koda.json` for graphics

```json
{
  "name": "mygame",
  "entry": "src/main.koda",
  "lint": "beginner",
  "native": {
    "sources": ["wrappers/raylib_shim/wrapper.c"],
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

## `@game` API reference

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
| `game.setGcBudget(ms)` | Incremental GC budget per frame |

`Key`, `Mouse`, and `Color` constants are defined in `stdlib/game.koda`. Full detail: [stdlib/game](../stdlib/game.md).

Raw Raylib shim names (`initwindow`, `drawtext`, …) remain available for advanced wrapping — see [raylib.md](raylib.md).

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

## Quick reference (cheatsheet)

| Need | API |
|------|-----|
| Window + loop | `import` / `#include "@game"` |
| Frame delta | `game.delta()` |
| Input | `game.keyDown(Key.Left)`, `game.mouseX()`, `game.mouseDown(Mouse.Left)` |
| Draw | `game.begin()`, `game.clear()`, `game.rect()`, `game.end()` |
| RNG | `random()`, `randomint(a,b)`, `import "@math"` |
| Timers | `import "@timer"` |
| JSON config | `import "@json"` |
| GC per frame | `game.setGcBudget(0.5)` |
| Diagnostics | `koda doctor`, `koda check` |

---

## GC and performance

Beginners should call `game.setGcBudget(0.5)` once per frame — the runtime spreads GC work automatically. Advanced: `arena()`, `gcDisable()`, `gcFrameStep()` — see [runtime and GC](../concepts/runtime-and-gc.md).

---

## Where to go next

| Document | Contents |
|----------|----------|
| [Beginner's guide](../beginners-guide.md) | Full onboarding |
| [stdlib/game](../stdlib/game.md) | `@game` module |
| [raylib.md](raylib.md) | Low-level shim (advanced) |
| [wrappers.md](../wrappers.md) | Extending bindings with kodawrap |
| [distribution.md](distribution.md) | Shipping builds |
