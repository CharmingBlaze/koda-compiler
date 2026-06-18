# Koda

**The beginner-friendly replacement for C and C++** — make games and desktop apps, ship native binaries, install nothing except one SDK zip.

No Go. No Python. No LLVM to install. No Visual Studio required to get started. Release builds embed the compiler; you unzip and run `koda`.

```koda
use raylib;
use koda.game;

struct Mario { x, y, speed, health }

func main() {
    game.open(800, 600, "My Game");
    defer game.close();
    game.fps(60);

    let player = Mario { x: 400, y: 300, speed: 220, health: 100 };

    loop {
        if not game.running() {
            break;
        }
        let dt = game.delta() or 0.016;
        if game.keyDown(KEY_RIGHT) {
            player.x += player.speed * dt;
        }
        game.begin();
        game.clear(#101018);
        game.rect(player.x, player.y, 32, 32, colors.white);
        game.end();
    }
}
```

**Same outcome as C/C++** (native executable, no VM on the player's PC). **None of the setup pain** (headers, CMake, manual linking, memory bugs in everyday logic).

---

## Install in 2 minutes

1. Download the **SDK zip** for your OS from [GitHub Releases](https://github.com/CharmingBlaze/koda-compiler/releases) (recommended: **v0.5.0**).
2. Unzip anywhere — keep `stdlib/` next to `koda`.
3. Run:

```bash
koda doctor
koda new bounce --template graphics
cd bounce
koda run
```

Full walkthrough: **[START_HERE.md](START_HERE.md)** · [Beginner's guide](docs/beginners-guide.md)

**VS Code / Cursor:** install syntax highlighting from [`vscode-koda/`](vscode-koda/) (F5 to try locally).

| Platform | Add to PATH (optional) |
|----------|-------------------------|
| Windows | `powershell -File scripts\install-koda.ps1` |
| macOS / Linux | `bash scripts/install-koda.sh` |

> **Contributors** who change the compiler itself need Go and LLVM — see [CONTRIBUTING.md](CONTRIBUTING.md). **Game and app makers do not.**

---

## Why Koda instead of C or C++?

| | C / C++ | Koda |
|---|---------|------|
| **Install to start** | Compiler + SDK + often CMake/vcpkg | One SDK zip |
| **Output** | Native `.exe` / binary | Native `.exe` / binary |
| **Game data** | `struct` + headers | `struct Player { x, y }` + methods |
| **Memory for gameplay** | Manual / smart pointers | GC (arena + `gcFrameStep` for games) |
| **Graphics** | `use raylib;` (548 functions) + optional `koda.game` · `"graphics": true` in `koda.json` |
| **Iteration** | compile → link → run | `koda run`, `koda watch` |
| **Beginner typos** | Silent wrong answers | `koda check --warn-unused` |
| **Ship to users** | Your binary (+ DLLs) | `koda bundle` + `assetPath()` |

Koda is **not** for kernels or embedded firmware. It **is** for the C/C++ most beginners actually want: **games, tools, and desktop apps**.

Coming from C or C++? **[From C / C++ guide](docs/guides/from-c.md)**

---

## Quick start

```bash
koda doctor
koda new bounce --template graphics   # Raylib bouncing ball
cd bounce && koda run

koda new myapp                        # console app
koda build -o myapp                   # ship standalone binary
koda bundle -o dist/MyGame            # folder + assets for players
koda check --warn-unused ./...
```

---

## Language highlights

| Feature | Notes |
|---------|-------|
| **Structs + methods** | `struct Rect { func area() { return this.w * this.h; } }` |
| **`const`** | Immutable bindings |
| **`enum` + `match`** | Game states without boolean soup |
| **String interpolation** | `"Score: {score}"` and `` `Hi ${name}` `` |
| **Core types** | `int`, `float`, `bool`, `string`, `array`, `map`, `func` (annotations optional) |
| **`use raylib;`** | Full Raylib API (548 functions, C names) |
| **`use koda.game;`** | Optional 2D helpers (`game.*`) on top of Raylib |
| **`koda doctor`** | OK/FAIL report when setup breaks |
| **`assetPath("x.png")`** | Bundle assets with your game |

Objects are for JSON and config — use **structs** for `Player`, `Enemy`, game state.

---

## Documentation

**[Documentation hub →](docs/README.md)**

| Section | Links |
|---------|-------|
| Learn | [START_HERE](START_HERE.md) · [Beginner's guide](docs/beginners-guide.md) · [Learn path](docs/learn/README.md) |
| Guides | [Games](docs/guides/game-dev.md) · [Apps](docs/guides/applications.md) · [From C/C++](docs/guides/from-c.md) |
| Reference | [Language](language.md) · [CLI](docs/reference/cli.md) · [Stdlib](docs/stdlib/README.md) |

---

## For contributors

Build from source: Go 1.22+, Clang, LLVM — **only if you hack the compiler**.

```bash
go test ./...
powershell -File scripts/build-runtime.ps1   # Windows
./scripts/build-runtime.sh                   # Linux / macOS
```

See [CONTRIBUTING.md](CONTRIBUTING.md) and [tests/MASTER_PLAN.md](tests/MASTER_PLAN.md).

---

[Changelog](CHANGELOG.md) · [Releases](https://github.com/CharmingBlaze/koda-compiler/releases)
