# Raylib + Koda — Complete Guide

Raylib is Koda's built-in graphics library for making games and apps with windows, drawing, input, audio, and more. Koda ships with a ready-to-use wrapper — you don't have to write any C to make a game.

This guide covers:
1. [How the wrapper system works](#1-how-the-wrapper-system-works)
2. [How to run a raylib program](#2-how-to-run-a-raylib-program)
3. [All available raylib functions](#3-available-raylib-functions)
4. [Colors](#4-colors)
5. [Key codes](#5-key-codes)
6. [Mouse input](#6-mouse-input)
7. [Step-by-step: build a complete game](#7-step-by-step-build-a-complete-game)
8. [Generating a fresh wrapper from any C header](#8-generating-a-fresh-wrapper)
9. [Environment variable reference](#9-environment-variable-reference)

---

## 1. How the wrapper system works

Koda links Raylib through the **full generated wrapper** (548 functions, original C names like `InitWindow`, `DrawText`).

| File | What it is |
|------|------------|
| `wrappers/raylib/raylib.koda` | Koda bindings — `// koda:extern` to `koda_wrap_raylib_*` |
| `wrappers/raylib/wrapper.c` | C glue — marshals Koda values to Raylib structs |

When you run `koda build`, the compiler links `wrappers/raylib/wrapper.c` from your project's `koda.json` (`native.sources`) or infers it from `koda_wrap_raylib_*` symbols in your program.

```koda
use raylib;

func main() {
    InitWindow(800, 600, "My Game");
    defer CloseWindow();
    ...
}
```

Optional **`koda.game`** helpers sit on top of the same full wrapper — they never hide the raw API.

Legacy **shim** (~33 functions, lowercase names): `koda setup raylib --shim` — not recommended for new projects.

---

## 2. How to run a raylib program

### Step 1 — Include the wrapper

At the top of your `.koda` file:

```koda
use raylib;
```

Adjust the path to be relative to your file. If your game is in the project root, it's just:

```koda
use raylib;
```

### Step 2 — Configure native sources (recommended: `koda.json`)

Graphics projects should set `"graphics": true` and point at the full wrapper in `koda.json`:

```json
{
  "native": {
    "sources": ["wrappers/raylib/wrapper.c"],
    "graphics": true
  }
}
```

Koda applies platform link flags automatically. You do **not** need to set `KODA_LINKFLAGS` for a standard graphics project.

**Configure graphics in an existing project:**

```bash
koda setup raylib
```

**Beginner helpers — `koda.game` on top of full Raylib:**

```koda
use raylib;
use koda.game;

func main() {
    game.open(800, 450, "Hello from Koda!");
    defer game.close();
    game.fps(60);
    while (game.running()) {
        game.begin();
        game.clear(colors.dark);
        game.text("Hello, Koda!", 300, 200, 30, colors.white);
        game.end();
    }
}
```

See [stdlib/game.md](../stdlib/game.md) for the full `koda.game` API.

**Manual override (optional):**

```powershell
$env:KODA_NATIVE_SOURCES = "wrappers\raylib\wrapper.c"
```

**Linux / macOS:**

```bash
export KODA_NATIVE_SOURCES=wrappers/raylib/wrapper.c
```

### Step 3 — Run or build

```
koda run mygame.koda
koda build mygame.koda -o mygame.exe
```

### Minimal working program

```koda
use raylib;
use koda.color;

func main() {
    InitWindow(800, 450, "Hello from Koda!");
    defer CloseWindow();
    SetTargetFPS(60);

    while (!WindowShouldClose()) {
        BeginDrawing();
        ClearBackground(asRaylib(colors.dark));
        DrawText("Hello, Koda!", 300, 200, 30, asRaylib(colors.white));
        EndDrawing();
    }
}
```

---

## 3. Available raylib functions

With `use raylib;` you can call **every** function in the generated wrapper (**548** bindings, same names as the C API). The table below is a **beginner cheat sheet** — not the full list.

**Complete reference:** `wrappers/raylib/api_reference.md` in the SDK (or open `wrappers/raylib/docs/index.html` in a browser).

All names are **case-insensitive** when called from Koda.

### Window

| Function | Arguments | What it does |
|----------|-----------|--------------|
| `InitWindow(w, h, title)` | width, height, title string | Open the window |
| `CloseWindow()` | — | Close and clean up |
| `WindowShouldClose()` | — | Returns `true` when the user closes the window or presses Escape |
| `SetTargetFPS(fps)` | number | Cap the frame rate (60 is standard) |
| `GetScreenWidth()` | — | Window width in pixels |
| `GetScreenHeight()` | — | Window height in pixels |
| `SetWindowTitle(title)` | string | Change the window title |
| `GetFPS()` | — | Current frames per second |

### Drawing

Call these **between** `BeginDrawing()` and `EndDrawing()` every frame.

| Function | Arguments | What it does |
|----------|-----------|--------------|
| `BeginDrawing()` | — | Start a frame |
| `EndDrawing()` | — | End the frame and display it |
| `ClearBackground(color)` | color as hex int | Fill the screen with a color |
| `DrawRectangle(x, y, w, h, color)` | all numbers | Draw a filled rectangle |
| `DrawCircle(x, y, radius, color)` | all numbers | Draw a filled circle |
| `DrawCircleLines(x, y, radius, color)` | all numbers | Draw a circle outline |
| `DrawLine(x1, y1, x2, y2, color)` | 5 numbers | Draw a 2D line |
| `DrawRectangleLines(x, y, w, h, color)` | all numbers | Draw a rectangle outline |
| `DrawText(text, x, y, size, color)` | string + numbers | Draw text on screen |
| `LoadTexture(path)` | string | Load an image; returns a texture handle |
| `DrawTexture(tex, x, y, color)` | handle + numbers | Draw a loaded texture |
| `UnloadTexture(tex)` | handle | Free a texture |
| `DrawLine3D(x1,y1,z1, x2,y2,z2, color)` | 7 numbers | Draw a line in 3D space |

### 3D

| Function | Arguments | What it does |
|----------|-----------|--------------|
| `BeginMode3D(eyeX,eyeY,eyeZ, targetX,targetY,targetZ, upX,upY,upZ, fov)` | 10 numbers | Start 3D mode with a camera |
| `EndMode3D()` | — | End 3D mode |
| `DrawGrid(slices, spacing)` | 2 numbers | Draw a ground grid |
| `DrawCube(x,y,z, w,h,d, color)` | 7 numbers | Draw a filled cube |
| `DrawCubeWires(x,y,z, w,h,d, color)` | 7 numbers | Draw a wireframe cube |

### Keyboard input

| Function | Arguments | What it does |
|----------|-----------|--------------|
| `IsKeyDown(key)` | key code number | `true` while the key is held |
| `IsKeyPressed(key)` | key code number | `true` for one frame when the key is first pressed |

### Mouse input

| Function | Arguments | What it does |
|----------|-----------|--------------|
| `GetMouseX()` | — | Mouse X in window coordinates |
| `GetMouseY()` | — | Mouse Y in window coordinates |
| `IsMouseButtonDown(btn)` | button index | `true` while the button is held (0=left, 1=right, 2=middle) |
| `IsMouseButtonPressed(btn)` | button index | `true` for one frame when first pressed |
| `GetMouseWheelMove()` | — | Scroll wheel delta this frame |

---

## 4. Colors

Use the built-in **`color`** palette and **`rgb()` / `rgba()`** helpers — no hex required for beginners.

```koda
let sky = colors.sky;
let grass = colors.grass;
let white = colors.white;

let custom = rgb(34, 139, 34);
let fade = rgba(255, 255, 255, 128);
```

| Name | Typical use |
|------|-------------|
| `colors.white`, `colors.black` | Text, outlines |
| `colors.red`, `colors.green`, `colors.blue` | Player, UI, accents |
| `colors.yellow`, `colors.gold`, `colors.orange` | Coins, stars |
| `colors.sky`, `colors.grass`, `colors.dirt` | Outdoor scenes |
| `colors.gray`, `colors.dark` | Backgrounds, panels |

With `koda.game`, `colors.white` aliases the same palette as `colors.white`.

**Advanced:** colors can still be passed as hex integers in `0xRRGGBBAA` format (red, green, blue, alpha — all 0–255). Import `@color` for `hsv()`, `lerp()`, and `toHex()`.

---

## 5. Key codes

The easiest way to use key codes is to define constants at the top of your file, the same way the examples do:

```koda
let KEY_LEFT  = 263;
let KEY_RIGHT = 262;
let KEY_UP    = 265;
let KEY_DOWN  = 264;
let KEY_SPACE = 32;
let KEY_ENTER = 257;
let KEY_ESCAPE = 256;

// Letters — just the ASCII code
let KEY_A = 65;
let KEY_W = 87;
let KEY_S = 83;
let KEY_D = 68;
let KEY_R = 82;
let KEY_P = 80;
```

**Common key codes table:**

| Key | Code | Key | Code |
|-----|------|-----|------|
| Space | 32 | Enter | 257 |
| Escape | 256 | Backspace | 259 |
| Left | 263 | Right | 262 |
| Up | 265 | Down | 264 |
| A–Z | 65–90 | 0–9 | 48–57 |
| F1–F12 | 290–301 | Left Shift | 340 |
| Left Ctrl | 341 | Left Alt | 342 |

---

## 6. Mouse input

The full wrapper exposes `GetMouseX`, `GetMouseY`, `IsMouseButtonDown`, and the rest of Raylib's input API (see [section 3](#3-available-raylib-functions)). For 2D games you can use optional `koda.game` helpers instead:

```koda
use raylib;
use koda.game;

func main() {
    game.open(800, 600, "Click demo");
    while (game.running()) {
        if (game.mouseDown(Mouse.Left)) {
            game.setTitle("Clicked at " + string(game.mouseX()) + "," + string(game.mouseY()));
        }
        game.begin();
        game.clear(colors.dark);
        game.circle(game.mouseX(), game.mouseY(), 8, colors.white);
        game.end();
    }
}
```

Or `#include "@input"` for standalone helpers (`mousePos()`, `mouseButton()`, …).

---

## 7. Step-by-step: build a complete game

We'll build a **Pong** game from scratch. Two paddles, a ball, score tracking, and game states.

### File layout

```
mygame/
  pong.koda
```

Run from the project root (where `wrappers/` lives):
```
koda run mygame/pong.koda
```

### The complete game

```koda
use raylib;

// ── Screen ────────────────────────────────────────────────────
let SW = 800;
let SH = 600;

// ── Colors ────────────────────────────────────────────────────
let BLACK  = colors.black;
let WHITE  = colors.white;
let GRAY   = colors.gray;
let GREEN  = colors.green;

// ── Key codes ─────────────────────────────────────────────────
let KEY_W      = 87;
let KEY_S      = 83;
let KEY_UP     = 265;
let KEY_DOWN   = 264;
let KEY_SPACE  = 32;
let KEY_R      = 82;

// ── Game state enum ───────────────────────────────────────────
enum State {
    Start,
    Playing,
    GameOver
}

// ── Paddle struct ─────────────────────────────────────────────
struct Paddle {
    x,
    y,
    w,
    h,
    speed,
    score
}

// ── Ball struct ───────────────────────────────────────────────
struct Ball {
    x,
    y,
    r,
    vx,
    vy
}

// ── Init helpers ──────────────────────────────────────────────
func makePaddle(x) {
    return Paddle { x: x, y: SH / 2 - 40, w: 12, h: 80, speed: 5, score: 0 };
}

func makeBall() {
    return Ball { x: SW / 2, y: SH / 2, r: 8, vx: 4, vy: 3 };
}

// ── Collision helper ──────────────────────────────────────────
func ballHitsPaddle(b, p) {
    return b.x - b.r < p.x + p.w &&
           b.x + b.r > p.x &&
           b.y - b.r < p.y + p.h &&
           b.y + b.r > p.y;
}

// ── Main ──────────────────────────────────────────────────────
func main() {
    InitWindow(SW, SH, "Pong — Koda");
    SetTargetFPS(60);

    let p1    = makePaddle(20);
    let p2    = makePaddle(SW - 32);
    let ball  = makeBall();
    let state = State.Start;

    while (!WindowShouldClose()) {

        // ── Update ────────────────────────────────────────────
        if (state == State.Start) {
            if (IsKeyPressed(KEY_SPACE)) {
                state = State.Playing;
            }

        } else if (state == State.Playing) {

            // Player 1 — W/S
            if (IsKeyDown(KEY_W) && p1.y > 0) {
                p1.y -= p1.speed;
            }
            if (IsKeyDown(KEY_S) && p1.y + p1.h < SH) {
                p1.y += p1.speed;
            }

            // Player 2 — Up/Down arrows
            if (IsKeyDown(KEY_UP) && p2.y > 0) {
                p2.y -= p2.speed;
            }
            if (IsKeyDown(KEY_DOWN) && p2.y + p2.h < SH) {
                p2.y += p2.speed;
            }

            // Move ball
            ball.x += ball.vx;
            ball.y += ball.vy;

            // Bounce off top/bottom
            if (ball.y - ball.r <= 0 || ball.y + ball.r >= SH) {
                ball.vy = -ball.vy;
            }

            // Bounce off paddles
            if (ballHitsPaddle(ball, p1)) {
                ball.vx = math.abs(ball.vx);   // always move right
            }
            if (ballHitsPaddle(ball, p2)) {
                ball.vx = -math.abs(ball.vx);  // always move left
            }

            // Score
            if (ball.x + ball.r < 0) {
                p2.score += 1;
                ball = makeBall();
                ball.vx = -ball.vx;
            }
            if (ball.x - ball.r > SW) {
                p1.score += 1;
                ball = makeBall();
            }

            // Win condition
            if (p1.score >= 7 || p2.score >= 7) {
                state = State.GameOver;
            }

        } else if (state == State.GameOver) {
            if (IsKeyPressed(KEY_R)) {
                p1 = makePaddle(20);
                p2 = makePaddle(SW - 32);
                ball = makeBall();
                state = State.Start;
            }
        }

        // ── Draw ──────────────────────────────────────────────
        BeginDrawing();
        ClearBackground(BLACK);

        // Center divider
        let i = 0;
        while (i < SH) {
            DrawRectangle(SW / 2 - 2, i, 4, 10, GRAY);
            i += 20;
        }

        // Paddles
        DrawRectangle(p1.x, p1.y, p1.w, p1.h, WHITE);
        DrawRectangle(p2.x, p2.y, p2.w, p2.h, WHITE);

        // Ball
        DrawCircle(ball.x, ball.y, ball.r, WHITE);

        // Scores
        DrawText(string(p1.score), SW / 2 - 60, 30, 40, WHITE);
        DrawText(string(p2.score), SW / 2 + 30, 30, 40, WHITE);

        // Overlays
        if (state == State.Start) {
            DrawText("PONG", SW / 2 - 60, SH / 2 - 60, 50, GREEN);
            DrawText("W/S and UP/DOWN to move", SW / 2 - 150, SH / 2, 20, WHITE);
            DrawText("Press SPACE to start", SW / 2 - 130, SH / 2 + 30, 20, WHITE);
        }

        if (state == State.GameOver) {
            let winner = if (p1.score >= 7) { "Player 1 Wins!" } else { "Player 2 Wins!" };
            DrawText(winner, SW / 2 - 120, SH / 2 - 30, 35, GREEN);
            DrawText("Press R to restart", SW / 2 - 110, SH / 2 + 20, 22, WHITE);
        }

        EndDrawing();
    }

    CloseWindow();
}
```

### Running it

Set the bridge (once per terminal session) if not using `koda.json`:

```powershell
$env:KODA_NATIVE_SOURCES = "wrappers\raylib\wrapper.c"
koda run mygame\pong.koda
```

```bash
export KODA_NATIVE_SOURCES=wrappers/raylib/wrapper.c
koda run mygame/pong.koda
```

### Building a distributable binary

```powershell
$env:KODA_NATIVE_SOURCES = "wrappers\raylib\wrapper.c"
koda build mygame\pong.koda -o dist\pong.exe
```

Copy `raylib.dll` next to `pong.exe` if using a dynamic-link build, or use `koda bundle` to pack everything:

```powershell
$env:KODA_BUNDLE_FILES = "raylib.dll"
koda bundle mygame\pong.koda -o dist\pong
```

---

## 8. Regenerating the wrapper (contributors)

The SDK **already includes** the full wrapper at `wrappers/raylib/` — you do **not** need to generate anything to use all Raylib commands. Run `use raylib;` and call any function from the [API reference](../../wrappers/raylib/api_reference.md).

To refresh bindings from a **newer** `raylib.h` (maintainers / contributors):

```powershell
# Build the wrapper generator first
go build -o bin\kodawrap.exe .\cmd\wrapgen

# Generate bindings from the header
.\bin\kodawrap.exe `
  -name raylib `
  -headers .\third_party\raylib_static\stage\include\raylib.h `
  -out .\wrappers\raylib_generated
```

This creates:
```
wrappers/raylib_generated/
  raylib.koda        — all functions declared with koda:extern
  wrapper.c          — full C glue for every function
  api_reference.md   — auto-generated function reference
```

Point your project at the new output if you are testing a regenerated tree:

```koda
#include "wrappers/raylib_generated/raylib.koda"
```

```powershell
$env:KODA_NATIVE_SOURCES = "wrappers\raylib_generated\wrapper.c"
```

The repo's default **`wrappers/raylib/`** tree is what `use raylib;` loads — **548 functions**, ready to use.

---

## 9. Environment variable reference

| Variable | What it does |
|----------|-------------|
| `KODA_NATIVE_SOURCES` | Path to the C bridge file to compile and link. Defaults from `koda.json` or inferred from `koda:extern` symbols. Default: `wrappers/raylib/wrapper.c` (full 548-function wrapper). Legacy shim: `koda setup raylib --shim`. |
| `KODA_LINKFLAGS` | Extra flags passed to clang when linking. Use for custom raylib installs: `-L/path/to/lib -lraylib` |
| `KODA_RAYLIB_STAGE` | Override path to a prebuilt raylib `stage/` directory (`include/` + `lib/`). |
| `KODA_USE_VENDORED_RAYLIB` | Set to `0` to skip auto-detected vendored raylib (useful on Linux ARM64 or custom builds). |
| `KODA_BUNDLE_FILES` | Extra files to copy when running `koda bundle` (e.g. `raylib.dll`). |
| `KODA_PATH` | Extra search roots for `#include` and `import()`. |
| `KODA_WRAPPERS` | Extra search roots specifically for wrapper modules. |

---

## See also

- **`wrappers/raylib/raylib.koda`** — full generated bindings used by `use raylib;`
- **`wrappers/raylib/api_reference.md`** — complete function list (548 entries)
- **`wrappers/raylib/docs/index.html`** — searchable offline API browser
- **`wrappers/raylib/wrapper.c`** — C bridge linked by `koda.json` / `KODA_NATIVE_SOURCES`
- **`examples/cube3d/`**, **`examples/spinning-cube/`** — raw Raylib 3D (no `koda.game` required)
- **`examples/games/brick-breaker/`** — canonical 2D sample using optional `koda.game`
- **`docs/wrappers.md`** — wrapper system internals and `KODA_LINKFLAGS` details
- **`docs/commands.md`** — all `koda` CLI commands
