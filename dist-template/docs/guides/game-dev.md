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

Lunar lander simulation using `print`, loops, and `randomint` — good for learning the language.

### 2. Graphics project (Raylib)

```bash
koda new bounce --template graphics
cd bounce
```

The template includes `wrappers/raylib_shim/` and `koda.json` native settings. Set link flags for your platform, then:

```powershell
# Windows example
$env:KODA_LINKFLAGS = "-lraylib -lopengl32 -lgdi32 -lwinmm"
koda run
```

### 3. Study the repo examples

| File | Description |
|------|-------------|
| `examples/games/brick_breaker.koda` | Full brick breaker |
| `examples/raylib_shim_demo.koda` | 3D camera + cube |
| `examples/games/lunar_lander_text.koda` | Text lander (standalone file) |

---

## Minimal windowed game

The Raylib **shim** uses lowercase names matching the wrapper file:

```koda
#include "wrappers/raylib_shim/raylib.koda"

func main() {
    let width = 800;
    let height = 600;
    let dark = 255;
    let white = 4294967295;

    initwindow(width, height, "My Game");
    settargetfps(60);

    while (!windowshouldclose()) {
        let dt = deltatime();

        begindrawing();
        clearbackground(dark);
        drawtext("Hello, Koda!", 300, 280, 30, white);
        enddrawing();
    }

    closewindow();
}
```

**Project workflow** — put native glue in `koda.json`:

```json
{
  "name": "mygame",
  "entry": "src/main.koda",
  "native": {
    "sources": ["wrappers/raylib_shim/wrapper.c"],
    "linkflags": ""
  }
}
```

Set `KODA_LINKFLAGS` in the shell or fill `linkflags` once you know your platform flags.

---

## The game loop

Every Koda game follows the same pattern as C:

```koda
#include "wrappers/raylib_shim/raylib.koda"

func main() {
    initwindow(800, 600, "Game");
    settargetfps(60);

    let playerx = 400.0;
    let playery = 300.0;

    while (!windowshouldclose()) {
        let dt = deltatime();

        // Update — input, physics, AI
        if (iskeydown(87)) {  // W
            playery = playery - 200.0 * dt;
        }

        // Draw
        begindrawing();
        clearbackground(255);
        drawrectangle(playerx - 20, playery - 20, 40, 40, 4294967295);
        enddrawing();
    }

    closewindow();
}
```

| Phase | Koda helpers |
|-------|----------------|
| Timing | `deltatime()`, `sleep()`, `clock()` |
| Input | Shim: `iskeydown`, `iskeypressed` — see [raylib.md](raylib.md) |
| Random | `random`, `randomint`, `randomseed` |
| Math | `lerp`, `clamp`, `distance`, `stdlib/math.koda` |
| 2D vectors | `stdlib/vec2.koda` |
| 3D vectors | `stdlib/vec3.koda` |
| Frame pacing | `stdlib/timer.koda` |

---

## Structs and enums for game state

```koda
struct Player {
    x, y, speed, health
}

enum State {
    Idle, Running, Dead
}

func update(player, dt) {
    player.x = player.x + player.speed * dt;
    if (player.health <= 0) {
        return State.Dead;
    }
    return State.Running;
}
```

Same mental model as C structs — field access is checked at compile time.

---

## Build, run, ship

```bash
koda run src/main.koda       # fast iteration
koda watch                   # rebuild on every .koda save
koda build -o mygame         # release binary
koda bundle -o dist/mygame   # exe + assets + README for players
```

Add sprites and sounds to `assets/` and list them in `koda.json`:

```json
"bundle": { "assets": ["assets"] }
```

---

## GC and performance

- Call **`gcFrameStep()`** once per frame in heavy games (see `tests/incremental_gc_test.koda`).
- Use **`gcDisable()`** / **`gcEnable()`** around critical sections if needed.
- Keep per-frame allocation low; reuse structs and arrays where possible.
- Push hot math to C libraries (physics, rendering) via wrappers.

---

## Where to go next

| Document | Contents |
|----------|----------|
| [raylib.md](raylib.md) | Functions, colors, keys, full Pong-style walkthrough |
| [wrappers.md](../wrappers.md) | Extending bindings with kodawrap |
| [distribution.md](../distribution.md) | Shipping builds |
| [language.md](../../language.md) | Full syntax reference |
| [commands.md](../commands.md) | CLI and `koda.json` |
