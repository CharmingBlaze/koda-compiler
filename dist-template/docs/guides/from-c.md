# Coming from C

Koda is built for developers who want **C-level outcomes** — a single native executable, direct hardware access through C libraries, no interpreter on the player's machine — without C's ceremony around headers, build systems, and manual memory for everyday game and app logic.

This guide explains what feels familiar, what is different, and how to think about Koda as a **modern replacement for C** when you're building games and desktop applications.

---

## Why Koda instead of C?

| You want… | In C | In Koda |
|-----------|------|---------|
| Ship one binary | `gcc` + link flags + maybe CMake | `koda build` or `koda bundle` |
| Structs and enums | `struct`, `enum` | `struct`, `enum` — same idea |
| Fast game loop | `while` + manual timing | `while` + `deltatime()` |
| Call Raylib / SDL / your C lib | Headers + link + ABI glue | `kodawrap` or bundled shim + `koda.json` |
| Strings and dynamic data | `char*`, malloc, pain | Built-in strings, arrays, objects |
| Memory safety for gameplay code | You manage everything | GC for Koda values; C only at FFI boundary |
| Iteration speed | Edit, compile, link, run | `koda run`, `koda watch` |

Koda is **not** a replacement for C in kernel drivers, freestanding firmware, or code that must avoid a GC. It **is** a strong replacement for the C (or C++) you write for **game logic, tools, and desktop apps** that today link against native libraries.

---

## Philosophy

1. **Compile to native** — LLVM IR → object file → linked executable. No bytecode VM for your game.
2. **Familiar syntax** — C-style `if`, `while`, `for`, `struct`, `enum`, `func`; JavaScript-style objects and closures where they help.
3. **Zero ceremony by default** — no header files for your own modules; `#include "other.koda"` instead.
4. **C interop when you need it** — graphics, physics, audio: stay in C libraries; call them from Koda through thin wrappers.
5. **Players run your binary** — they do not install Koda, Go, or LLVM to play your game.

---

## Side-by-side: the same program

### C

```c
#include <stdio.h>
#include <math.h>

typedef struct {
    double x, y, speed, health;
} Player;

enum State { Idle, Running, Dead };

State update(Player* p, double dt) {
    p->x += p->speed * dt;
    if (p->health <= 0) return Dead;
    return Running;
}

int main(void) {
    Player p = { 0, 0, 200, 100 };
    State s = update(&p, 0.016);
    printf("%d\n", s);
    return 0;
}
```

### Koda

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

let p = Player { x: 0, y: 0, speed: 200, health: 100 };
let state = update(p, 0.016);
print(state);
```

No `malloc` for `Player` in typical gameplay code. No `printf` format strings unless you want them — `print` handles multiple values.

---

## Mapping C concepts to Koda

| C | Koda |
|---|------|
| `#include "foo.h"` | `#include "foo.koda"` |
| `int main()` | Top-level code or `func main()` |
| `struct` / `typedef` | `struct Name { fields }` |
| `enum` | `enum Name { A, B, C }` |
| `NULL` | `null` |
| `printf` | `print(...)` |
| `malloc` / `free` | Rare — GC handles Koda objects; use C only in wrappers |
| `static` globals | `let` at top level |
| Function pointers | Functions are values; closures capture locals |
| `#define` | `let` constants (no preprocessor) |
| Multi-file project | `koda.json` + `#include` + `koda new` |

---

## Types: what to expect

- **Numbers** are 64-bit floats at runtime (like JavaScript). Integer literals and game math work naturally; use `floor`, `round`, or struct fields when you need discrete grid logic.
- **Strings** are real objects — concatenation, `len`, methods — not `char*`.
- **Arrays** are growable (`push`, indexing). Out-of-bounds access panics with a clear error (safer than silent C bugs).
- **Objects** `{ x: 1, y: 2 }` for loose data; **structs** when you want named fields with compile-time checks.

---

## Memory and performance

- **Gameplay code** runs on a **tracing GC** with generational collection and write barriers — tuned for games (frame-step GC, `gcDisable` / `gcCollect` for critical sections).
- **Hot paths** can stay in **C** (physics, rendering) via wrappers; Koda orchestrates.
- **Release builds** use LLVM + native linking — overhead is small compared to hand-written C for typical script-like game logic.

Read [Game development](game-dev.md) for loop timing with `deltatime()` and [Applications](applications.md) for I/O-heavy tools.

---

## Calling C libraries (FFI)

You do not abandon C — you **stop writing all your game in C**.

1. Use **kodawrap** on a header → `.koda` bindings + `wrapper.c`.
2. List sources and link flags in **`koda.json`** or environment variables.
3. `#include` the generated module and call functions like normal Koda code.

```koda
#include "wrappers/raylib_shim/raylib.koda"

func main() {
    initwindow(800, 600, "My Game");
    settargetfps(60);
    while (!windowshouldclose()) {
        begindrawing();
        clearbackground(255);
        enddrawing();
    }
    closewindow();
}
```

Full detail: [Wrappers](../wrappers.md) and [Raylib guide](raylib.md).

---

## Project workflow (vs Makefile)

```bash
koda new mygame --template graphics   # or game / hello
cd mygame
koda run          # compile + run
koda watch        # rebuild on save
koda build -o mygame
koda bundle -o dist/mygame
```

`koda.json` holds entry point, asset folders, and native link settings — see [Commands](../commands.md#kodajson-project-manifest).

---

## What to read next

| Goal | Document |
|------|----------|
| Install and first project | [Getting started](getting-started.md) |
| Full language tutorial | [Using the language](../using-the-language.md) |
| Make a windowed game | [Game dev](game-dev.md) → [Raylib](raylib.md) |
| CLI tool or file processor | [Applications](applications.md) |
| Ship to players | [Distribution](../distribution.md) |
| Syntax lookup | [language.md](../../language.md) |
