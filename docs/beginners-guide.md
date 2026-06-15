# Beginner's guide to Koda

Welcome. This guide assumes you have **never used Koda** and may be new to compiled languages. By the end you will install Koda, run code, understand core syntax, use modules, read files, and ship a small program.

> **Style:** This doc follows our [documentation style guide](STYLE-GUIDE.md).

**Estimated time:** 2–3 hours if you work through every section. Skim the tables and run the examples — that is enough for a first day.

---

## What is Koda?

Koda is an **easy native language for games, tools, and desktop apps** — a practical **replacement for C and C++** when you want native speed without installing Go, Python, LLVM, or a giant IDE.

| You write | You get |
|-----------|---------|
| Game logic, tools, desktop apps | One native binary (like C/C++) |
| `struct`, `func`, loops | Fast iteration with `koda run` |
| Optional Raylib / C libraries | Via wrappers and `koda.json` |

Koda is **not** a VM language (not Python, not Lua-in-a-runtime). There is no interpreter on the user's machine.

**You do not need:** Go, Python, Node, Rust, LLVM, Visual Studio, or CMake to **use** the release SDK.

---

## Who this guide is for

| You are… | Start here | Then read |
|----------|------------|-----------|
| Completely new to programming | Sections 1–8 below | [Learn path](learn/README.md) |
| Know JavaScript or Python | Sections 2–6, 9 | [Game dev](guides/game-dev.md) or [Applications](guides/applications.md) |
| Know C or C++ | [From C](guides/from-c.md) | [Language reference](../language.md) |
| Want games with graphics | Section 10 | [Raylib guide](guides/raylib.md) |

---

## 1. Install (5 minutes)

You do **not** need Go, Python, LLVM, Visual Studio, Node, or Rust to **use** Koda from releases.

1. Download the **SDK zip** for your platform from [GitHub Releases](https://github.com/CharmingBlaze/koda-compiler/releases).
2. Unzip so **`stdlib/`** sits next to **`koda`** (or `koda.exe`).
3. Open a terminal in that folder:

```bash
koda doctor
koda version
```

If `doctor` shows **FAIL**, fix those lines before continuing. See [START_HERE.md](../START_HERE.md) for platform-specific steps.

Optional — add Koda to PATH:

| Platform | Command (from SDK folder) |
|----------|---------------------------|
| Windows | `powershell -File scripts\install-koda.ps1` |
| macOS / Linux | `bash scripts/install-koda.sh` |

> **Note:** Building Koda from source is for **contributors** only. See [CONTRIBUTING.md](../CONTRIBUTING.md).

Full install details: [Learn — Install](learn/02-install-and-first-run.md).

---

## 2. First program (10 minutes)

Create `hello.koda`:

```koda
print("Hello, Koda!");
```

Run it:

```bash
koda run hello.koda
```

You should see:

```text
Hello, Koda!
```

**What happened:** Koda compiled your file to a temporary native executable and ran it. No separate compile step is required for quick tests.

### Project template

```bash
koda new myapp
cd myapp
koda run
```

This creates `koda.json`, `src/main.koda`, and `assets/`.

| Template | Command | What you get |
|----------|---------|--------------|
| Hello | `koda new myapp` | Minimal print program |
| Text game | `koda new lander --template game` | Lunar lander in the terminal |
| Graphics | `koda new bounce --template graphics` | Bouncing ball with `@game` API |

---

## 3. Variables and types (15 minutes)

```koda
let name = "Ada";
let score = 42;
let alive = true;
let empty = null;

const gravity = 900;
const screenWidth = 800;

print(type(name));   // string
print(type(score));   // number
```

| Keyword | Meaning |
|---------|---------|
| `let` | Mutable binding |
| `const` | Immutable — cannot reassign |

Optional type annotations (add when you want clarity):

```koda
let lives: int = 3;
let dt: float = 0.016;
let label: string = "Player";
```

| Type | Example | Notes |
|------|---------|-------|
| Number | `42`, `3.14` | Inferred as 64-bit float; use `int` when needed |
| String | `"hello"` | Escapes: `\n`, `\t`, `\"`, `\\` |
| Bool | `true`, `false` | |
| Null | `null` | “No value” |
| Array | `[1, 2, 3]` | Ordered list |
| Struct | `Player { x: 1 }` | Named game/app data (preferred) |
| Object | `{ x: 1 }` | JSON/config maps (not for core game state) |
| Function | `func(x) { return x; }` | Callable value |

**Equality:** use `==` and `!=` only. One obvious comparison operator.

**Case-insensitive:** `print`, `Print`, and `PRINT` are the same for keywords and builtins.

Reassign `let` with `=`:

```koda
let x = 1;
x = x + 1;
print(x);  // 2
```

More: [Learn — Variables](learn/03-variables-and-types.md).

---

## 4. Control flow (15 minutes)

```koda
let health = 75;

if (health <= 0) {
    print("dead");
} else if (health < 30) {
    print("low");
} else {
    print("ok");
}

let i = 0;
while (i < 3) {
    print(i);
    i = i + 1;
}

for (let n of [10, 20, 30]) {
    print(n);
}
```

`switch` works like C — use `break` to avoid fall-through:

```koda
let state = 2;
switch (state) {
    case 1:
        print("one");
        break;
    case 2:
        print("two");
        break;
    default:
        print("other");
}
```

More: [Learn — Control flow](learn/04-control-flow.md).

---

## 5. Functions (15 minutes)

```koda
func greet(name) {
    return "Hello, " + name;
}

print(greet("Koda"));

func add(a, b) {
    return a + b;
}

let double = func(x) {
    return x * 2;
};

print(double(21));
```

Closures capture variables:

```koda
func makeCounter() {
    let n = 0;
    return func() {
        n = n + 1;
        return n;
    };
}

let counter = makeCounter();
print(counter());  // 1
print(counter());  // 2
```

More: [Learn — Functions](learn/05-functions.md).

---

## 6. Structs and game data (20 minutes)

**Structs are the main data model** for games and apps. Objects are for JSON and config — not your first choice for a player or enemy.

```koda
struct Player {
    x, y,
    speed,
    health
}

let player = Player {
    x: 400,
    y: 300,
    speed: 220,
    health: 100
};

func updatePlayer(player, dt) {
    if (isKeyDown(KEY_RIGHT)) {
        player.x = player.x + player.speed * dt;
    }
}
```

Enums for phases and states:

```koda
enum State {
    Idle, Running, Dead
}

let state = State.Running;
if (state == State.Dead) {
    print("game over");
}
```

More: [Learn — Structs and enums](learn/07-structs-and-enums.md).

---

## 7. Arrays and JSON objects (15 minutes)

Arrays hold lists. Object literals are for config and parsed JSON — not core game entities.

```koda
let items = ["sword", "shield"];
items.push("potion");
print(len(items));

let config = { volume: 80, fullscreen: true };
```

**Array helpers** — `import "@array"`:

```koda
let nums = range(0, 5);   // [0,1,2,3,4]
let total = sum(nums);
shuffle(nums);
```

More: [Learn — Objects and arrays](learn/06-objects-and-arrays.md).

---

## 8. Modules and imports (15 minutes)

**Import** loads a standard or local module (preferred for beginners):

```koda
let math = import "@math";
print(math.sqrt(16));

let io = import "@io";
let json = import "@json";
```

Graphics projects use the `@game` wrapper:

```koda
#include "wrappers/raylib_shim/raylib.koda"
#include "@game"

func main() {
    game.open(800, 600, "My Game");
    game.fps(60);
    while (game.running()) {
        let dt = game.delta();
        game.begin();
        game.clear(Color.dark);
        game.text("Hello, Koda!", 300, 280, 30, Color.white);
        game.end();
    }
}
```

| Import | Provides |
|--------|----------|
| `@game` | `open`, `running`, `delta`, `begin`, `end`, `clear`, `text`, `rect`, … |
| `@math` | `sin`, `lerp`, `random`, `pi`, … |
| `@json` | `parse`, `stringify`, `try_parse` |
| `@io` | `read`, `write`, `exists`, `list`, … |
| `@timer` | Countdowns, cooldowns, intervals |
| `@array` | `range`, `shuffle`, `zip`, … |

`#include` merges source at compile time — advanced/low-level; beginners use `import` instead.

More: [Learn — Modules](learn/08-modules-and-imports.md) · [Stdlib overview](stdlib/README.md).

---

## 9. Files and JSON (15 minutes)

### Global builtins

```koda
let ok = writefile("save.txt", "hello");
let text = readfile("save.txt");
if (fileexists("save.txt")) {
    print(text);
}
```

`readfile` / `writefile` return **result objects** with `.ok`, `.value`, `.error`.

### `@io` module

```koda
let io = import "@io";

io.write("config.json", "{\"level\":1}");
let entries = io.list("assets");
if (io.isfile("save.txt")) {
    print(io.size("save.txt"));
}
```

### JSON

```koda
let json = import "@json";

let cfg = json.parse("{\"level\":1,\"name\":\"Ada\"}");
let text = json.stringify(cfg, 2);   // pretty-print with indent 2

let result = json.try_parse("{bad}");
if (!result.ok) {
    print(result.error);
}
```

More: [Learn — Files and JSON](learn/09-files-and-json.md) · [Applications guide](guides/applications.md).

---

## 10. Build, test, and ship (15 minutes)

```bash
# Daily development
koda run src/main.koda
koda watch                    # rebuild on save
koda check                    # parse + type check only
koda test                     # run tests/*.koda

# Release
koda build -o myapp           # native executable
koda build --debug -o myapp   # debug symbols
koda bundle -o dist/MyApp     # exe + assets for players
koda clean                    # remove build artifacts
```

`koda.json` at the project root defines entry point, native C sources, and bundle assets:

```json
{
  "name": "mygame",
  "entry": "src/main.koda",
  "lint": "beginner",
  "bundle": { "assets": ["assets"] },
  "native": {
    "sources": ["wrappers/raylib_shim/wrapper.c"],
    "graphics": true
  }
}
```

Set `"graphics": true` for Raylib projects — the compiler adds platform link flags automatically. Advanced users can still set `KODA_LINKFLAGS`.

More: [Learn — Building and shipping](learn/10-building-and-shipping.md) · [CLI reference](reference/cli.md) · [Distribution](guides/distribution.md).

---

## 11. Games (optional)

**Text game** — no extra libraries:

```bash
koda new lander --template game
cd lander
koda run
```

**Graphics** — `@game` wrapper over Raylib:

```bash
koda new bounce --template graphics
cd bounce
koda doctor    # check toolchain + raylib
koda run
```

Gold-standard game loop:

```koda
#include "wrappers/raylib_shim/raylib.koda"
#include "@game"

struct Player {
    x, y, speed, health
}

func main() {
    game.open(800, 600, "Koda Game");
    game.fps(60);

    let player = Player {
        x: 400, y: 300, speed: 220, health: 100
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
}
```

Cheatsheet: [Game development](guides/game-dev.md) · [Raylib guide](guides/raylib.md).

---

## 12. When something goes wrong

| Problem | Try |
|---------|-----|
| `koda` not found | Add SDK folder to PATH |
| `stdlib` missing | Keep `stdlib/` next to `koda.exe` |
| Link errors on graphics | Run `koda doctor` — follow FAIL lines |
| Parse error | `koda check file.koda` |
| Weird runtime behavior | Small repro + [FAQ](faq.md) |

Run diagnostics first:

```bash
koda doctor
koda help
```

Full list: [Troubleshooting](troubleshooting.md).

---

## Learning paths

### Path A — Language first

1. [Learn path](learn/README.md) (10 chapters)
2. [Language reference](../language.md) (lookup)
3. [Game dev](guides/game-dev.md) or [Applications](guides/applications.md)

### Path B — Game first

1. Sections 1–2 of this guide
2. [Game development](guides/game-dev.md)
3. [Raylib guide](guides/raylib.md)
4. Backfill syntax from [Learn path](learn/README.md) as needed

### Path C — From C

1. [From C](guides/from-c.md)
2. [Language reference](../language.md)
3. [Wrappers](wrappers.md) for library interop

---

## Quick reference card

| Need | Command / API |
|------|----------------|
| Run | `koda run file.koda -- --arg` |
| Check all | `koda check ./...` |
| Lint CI | `koda lint` |
| REPL | `koda repl` |
| Bench | `koda bench file.koda --count 10` |
| Build | `koda build -o app` |
| Math | `import "@math"` → `math.lerp`, `math.randomint` |
| Files | `import "@io"` → `io.read`, `io.write` |
| JSON | `import "@json"` → `json.parse`, `json.stringify` |
| Timers | `import "@timer"` → `cooldown`, `interval` |
| RNG | `random()`, `randomint(a,b)`, `randomseed(n)` |
| Frame time | `deltatime()`, `programtime()` |
| Print debug | `print()`, `warn()` |
| Tests | `koda test` |

---

## Next steps

| Document | Why read it |
|----------|-------------|
| [Documentation hub](README.md) | Full index |
| [Learn path](learn/README.md) | Chapter-by-chapter tutorial |
| [Language reference](../language.md) | Every syntax form |
| [Stdlib](stdlib/README.md) | Module APIs |
| [FAQ](faq.md) | Common questions |

Happy building.
