# Koda

**Easy native C for games and apps.** Compile to a single executable — no VM, no interpreter on the player's machine.

Koda is a practical alternative to C for game logic, tools, and desktop apps: native output, structs for game data, beginner-friendly `@game` API, and `koda doctor` when something breaks.

```koda
#include "wrappers/raylib_shim/raylib.koda"
#include "@game"

struct Player { x, y, speed, health }

func main() {
    game.open(800, 600, "Koda Game");
    game.fps(60);

    let player = Player { x: 400, y: 300, speed: 220, health: 100 };

    while (game.running()) {
        let dt = game.delta();
        if (game.keyDown(Key.Right)) {
            player.x = player.x + player.speed * dt;
        }
        game.begin();
        game.clear(Color.dark);
        game.rect(player.x, player.y, 32, 32, Color.white);
        game.end();
    }
}
```

---

## New to Koda?

| Start here | What you'll do |
|------------|----------------|
| **[Beginner's guide](docs/beginners-guide.md)** | Install, structs-first syntax, `@game`, ship a build |
| **[Learn path](docs/learn/README.md)** | Short chapters (structs before objects) |
| **[Game development](docs/guides/game-dev.md)** | Graphics without manual link flags |

Download **`koda`** + **`stdlib/`** from [GitHub Releases](https://github.com/CharmingBlaze/koda-compiler/releases). Run **`koda doctor`** after install.

---

## Quick start

```bash
koda doctor
koda new bounce --template graphics
cd bounce
koda run

# Console
koda new myapp
koda run

koda build -o myapp
koda bundle -o dist/MyGame
koda check --warn-unused
koda bench tests/hello.koda --count 5
```

Graphics templates set `"graphics": true` in `koda.json` — Raylib link flags are applied automatically.

---

## Language highlights

| Feature | Notes |
|---------|-------|
| **Structs** | Main data model for game/app state |
| **`const`** | Immutable bindings |
| **`==` / `!=`** | One equality operator (no `===` confusion) |
| **`import "@game"`** | Beginner game API over Raylib |
| **`import "@array"`** | `arraypush`, `range`, `shuffle`, … registered as builtins |
| **`assetPath("x.png")`** | Resolve bundled assets at runtime |
| **`koda check --warn-unused`** | Catch typos in unused variables |
| **Optional types** | `let lives: int = 3` when you need them |

Objects are for JSON and config — not your first choice for a `Player`.

---

## Documentation

**[Documentation hub →](docs/README.md)**

| Section | Links |
|---------|-------|
| Learn | [Beginner's guide](docs/beginners-guide.md) · [Learn path](docs/learn/README.md) |
| Guides | [Games](docs/guides/game-dev.md) · [Apps](docs/guides/applications.md) · [From C](docs/guides/from-c.md) |
| Reference | [Language](language.md) · [CLI](docs/reference/cli.md) · [Stdlib](docs/stdlib/README.md) |
| Tooling | [koda doctor](docs/reference/cli.md) · [MASTER_PLAN](tests/MASTER_PLAN.md) |

---

## Why Koda?

| | C | Koda |
|---|-----|------|
| Output | Native executable | Native executable |
| Game data | Structs + headers | `struct Player { x, y }` |
| Graphics | Link Raylib yourself | `@game` + `koda.json` `"graphics": true` |
| Iteration | Manual compile/link | `koda run`, `koda watch` |
| Diagnostics | Cryptic linker errors | `koda doctor` OK/FAIL report |
| Ship | Binary + DLLs | `koda bundle` + `assetPath()` |

---

## For contributors

Build from source: Go 1.22+, Clang, LLVM. After pulling:

```bash
go test ./...
# Windows
powershell -File scripts/build-runtime.ps1
# Linux/macOS
./scripts/build-runtime.sh
```

See [CONTRIBUTING.md](CONTRIBUTING.md) and [tests/MASTER_PLAN.md](tests/MASTER_PLAN.md).

---

[Changelog](CHANGELOG.md) · [Documentation hub](docs/README.md)
