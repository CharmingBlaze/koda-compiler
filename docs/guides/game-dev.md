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

The graphics template sets `"graphics": true` in `koda.json` — platform link flags are applied automatically. No manual `KODA_LINKFLAGS` for beginners.

### 3. Study the repo examples

| File | Description |
|------|-------------|
| `examples/games/brick_breaker.koda` | Full brick breaker |
| `examples/raylib_shim_demo.koda` | 3D camera + cube |
| `examples/games/lunar_lander_text.koda` | Text lander (standalone file) |

---

## Gold-standard windowed game

Use the `@game` wrapper — not raw Raylib names:

```koda
#include "wrappers/raylib_shim/raylib.koda"
#include "@game"

struct Player {
    x, y,
    speed,
    health
}

func main() {
    game.open(800, 600, "Koda Game");
    game.fps(60);

    let player = Player {
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

        game.begin();
        game.clear(Color.dark);
        game.rect(player.x, player.y, 32, 32, Color.white);
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

Structs are the main data model — not object literals:

```koda
struct Player {
    x, y, speed, health
}

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
| `game.rect(x, y, w, h, color)` | Draw rectangle |
| `game.circle(x, y, r, color)` | Draw circle |
| `game.keyDown(key)` | Key held this frame |
| `game.keyPressed(key)` | Key pressed this frame |
| `game.setGcBudget(ms)` | Incremental GC budget per frame |

`Key` and `Color` constants are defined in `stdlib/game.koda`. Full detail: [stdlib/game](../stdlib/game.md).

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
| Input | `game.keyDown(Key.Left)` |
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
