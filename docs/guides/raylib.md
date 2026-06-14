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

Koda can call C libraries by linking them into your binary. The wrapper has two parts:

| File | What it is |
|------|------------|
| `wrappers/raylib_shim/raylib.koda` | Koda-side declarations — binds Koda names to C symbols |
| `wrappers/raylib_min/raylib_bridge.c` | C glue — converts Koda values to C types and calls raylib |

When you run `koda build`, you point it at the C glue file with the `KODA_NATIVE_SOURCES` environment variable. The compiler includes it in the same clang link step as your program — no separate compilation needed.

The Koda side uses `// koda:extern` directives. Each line wires a Koda name to a C function and declares the argument count:

```koda
// koda:extern initwindow koda_shim_InitWindow 3
let initwindow = 0;
```

After the `#include`, `initwindow(800, 600, "My Game")` calls straight into the C bridge, which calls `InitWindow()` in raylib. No overhead beyond a normal function call.

---

## 2. How to run a raylib program

### Step 1 — Include the wrapper

At the top of your `.koda` file:

```koda
#include "../../wrappers/raylib_shim/raylib.koda"
```

Adjust the path to be relative to your file. If your game is in the project root, it's just:

```koda
#include "wrappers/raylib_shim/raylib.koda"
```

### Step 2 — Set the bridge file

Before running, tell Koda where the C glue is. Do this once in your terminal session:

**Windows (PowerShell):**
```powershell
$env:KODA_NATIVE_SOURCES = "wrappers\raylib_min\raylib_bridge.c"
```

**Windows (cmd):**
```cmd
set KODA_NATIVE_SOURCES=wrappers\raylib_min\raylib_bridge.c
```

**Linux / macOS:**
```bash
export KODA_NATIVE_SOURCES=wrappers/raylib_min/raylib_bridge.c
```

### Step 3 — Run or build

```
koda run mygame.koda
koda build mygame.koda -o mygame.exe
```

### Minimal working program

```koda
#include "wrappers/raylib_shim/raylib.koda"

func main() {
    InitWindow(800, 450, "Hello from Koda!");
    SetTargetFPS(60);

    while (!WindowShouldClose()) {
        BeginDrawing();
        ClearBackground(0x181818FF);
        DrawText("Hello, Koda!", 300, 200, 30, 0xFFFFFFFF);
        EndDrawing();
    }

    CloseWindow();
}
```

---

## 3. Available raylib functions

All names are **case-insensitive** when called from Koda.

### Window

| Function | Arguments | What it does |
|----------|-----------|--------------|
| `InitWindow(w, h, title)` | width, height, title string | Open the window |
| `CloseWindow()` | — | Close and clean up |
| `WindowShouldClose()` | — | Returns `true` when the user closes the window or presses Escape |
| `SetTargetFPS(fps)` | number | Cap the frame rate (60 is standard) |

### Drawing

Call these **between** `BeginDrawing()` and `EndDrawing()` every frame.

| Function | Arguments | What it does |
|----------|-----------|--------------|
| `BeginDrawing()` | — | Start a frame |
| `EndDrawing()` | — | End the frame and display it |
| `ClearBackground(color)` | color as hex int | Fill the screen with a color |
| `DrawRectangle(x, y, w, h, color)` | all numbers | Draw a filled rectangle |
| `DrawCircle(x, y, radius, color)` | all numbers | Draw a filled circle |
| `DrawText(text, x, y, size, color)` | string + numbers | Draw text on screen |
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

---

## 4. Colors

Colors are passed as a single **hex integer** in `0xRRGGBBAA` format (red, green, blue, alpha — all 0–255).

```koda
let black   = 0x000000FF;
let white   = 0xFFFFFFFF;
let red     = 0xFF0000FF;
let green   = 0x00FF00FF;
let blue    = 0x0000FFFF;
let yellow  = 0xFFFF00FF;
let orange  = 0xFFA500FF;
let purple  = 0x800080FF;
let gray    = 0x808080FF;
let darkgray = 0x303030FF;
```

Build a color from individual components (0–255 each):

```koda
func rgba(r, g, b, a) {
    return ((r & 255) << 24) | ((g & 255) << 16) | ((b & 255) << 8) | (a & 255);
}

let myColor = rgba(100, 200, 50, 255);
```

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

The raylib shim doesn't currently expose mouse functions. If you need mouse position, add them to `wrappers/raylib_min/raylib_bridge.c` using the same pattern as the existing functions:

```c
KodaValue rl_get_mouse_x(int argCount, KodaValue* args) {
    return NUMBER_VAL(GetMouseX());
}
KodaValue rl_get_mouse_y(int argCount, KodaValue* args) {
    return NUMBER_VAL(GetMouseY());
}
KodaValue rl_is_mouse_button_down(int argCount, KodaValue* args) {
    if (argCount < 1) return BOOL_VAL(false);
    return BOOL_VAL(IsMouseButtonDown(knum(args[0])));
}
```

Then declare them in your `.koda` file:

```koda
// koda:extern getmousex rl_get_mouse_x 0
let getmousex = 0;
// koda:extern getmousey rl_get_mouse_y 0
let getmousey = 0;
// koda:extern ismousebuttondown rl_is_mouse_button_down 1
let ismousebuttondown = 0;
```

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
#include "wrappers/raylib_shim/raylib.koda"

// ── Screen ────────────────────────────────────────────────────
let SW = 800;
let SH = 600;

// ── Colors ────────────────────────────────────────────────────
let BLACK  = 0x000000FF;
let WHITE  = 0xFFFFFFFF;
let GRAY   = 0x555555FF;
let GREEN  = 0x00FF88FF;

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

Set the bridge (once per terminal session):

```powershell
$env:KODA_NATIVE_SOURCES = "wrappers\raylib_min\raylib_bridge.c"
koda run mygame\pong.koda
```

```bash
export KODA_NATIVE_SOURCES=wrappers/raylib_min/raylib_bridge.c
koda run mygame/pong.koda
```

### Building a distributable binary

```powershell
$env:KODA_NATIVE_SOURCES = "wrappers\raylib_min\raylib_bridge.c"
koda build mygame\pong.koda -o dist\pong.exe
```

Copy `raylib.dll` next to `pong.exe` if using a dynamic-link build, or use `koda bundle` to pack everything:

```powershell
$env:KODA_BUNDLE_FILES = "raylib.dll"
koda bundle mygame\pong.koda -o dist\pong
```

---

## 8. Generating a fresh wrapper

If you want access to more raylib functions than the shim provides, use `kodawrap` to generate bindings from the full `raylib.h` header:

```powershell
# Build the wrapper generator first
go build -o bin\kodawrap.exe .\cmd\wrapgen

# Generate bindings from the header
.\bin\kodawrap.exe `
  -name raylib `
  -headers .\raylib_lib\raylib-5.0_win64_mingw-w64\include\raylib.h `
  -out .\wrappers\raylib_generated
```

This creates:
```
wrappers/raylib_generated/
  raylib.koda        — all functions declared with koda:extern
  wrapper.c          — full C glue for every function
  api_reference.md   — auto-generated function reference
```

Then use the generated wrapper instead of the shim:

```koda
#include "wrappers/raylib_generated/raylib.koda"
```

And set:
```powershell
$env:KODA_NATIVE_SOURCES = "wrappers\raylib_generated\wrapper.c"
```

The full generated wrapper (`wrappers/raylib/raylib.koda`) is already included in the repo — it covers hundreds of raylib functions.

---

## 9. Environment variable reference

| Variable | What it does |
|----------|-------------|
| `KODA_NATIVE_SOURCES` | Path to the C bridge file to compile and link. Set to `wrappers\raylib_min\raylib_bridge.c` for the shim, or `wrappers\raylib\wrapper.c` for the full wrapper. |
| `KODA_LINKFLAGS` | Extra flags passed to clang when linking. Use for custom raylib installs: `-L/path/to/lib -lraylib` |
| `KODA_RAYLIB_STAGE` | Override path to a prebuilt raylib `stage/` directory (`include/` + `lib/`). |
| `KODA_USE_VENDORED_RAYLIB` | Set to `0` to skip auto-detected vendored raylib (useful on Linux ARM64 or custom builds). |
| `KODA_BUNDLE_FILES` | Extra files to copy when running `koda bundle` (e.g. `raylib.dll`). |
| `KODA_PATH` | Extra search roots for `#include` and `import()`. |
| `KODA_WRAPPERS` | Extra search roots specifically for wrapper modules. |

---

## See also

- **`wrappers/raylib_shim/raylib.koda`** — the shim declarations used by all examples
- **`wrappers/raylib_min/raylib_bridge.c`** — the C bridge source you can extend
- **`wrappers/raylib/raylib.koda`** — the full auto-generated wrapper
- **`examples/games/brick_breaker.koda`** — a complete brick breaker game using this system
- **`docs/wrappers.md`** — wrapper system internals and `KODA_LINKFLAGS` details
- **`docs/commands.md`** — all `koda` CLI commands
