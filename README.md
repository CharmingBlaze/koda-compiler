# Koda

A modern language for **games and applications** — compile to a **single native binary**. No VM. No interpreter on the player's machine.

Koda is a practical alternative to **C** for game logic, tools, and desktop apps: same native output and C library access, with faster iteration and less boilerplate.

```koda
struct Player { x, y, speed, health }

func update(p, dt) {
    p.x = p.x + p.speed * dt;
    if (p.health <= 0) { return false; }
    return true;
}
```

---

## New to Koda?

| Start here | What you'll do |
|------------|----------------|
| **[Beginner's guide](docs/beginners-guide.md)** | Install, syntax, modules, ship a build |
| **[Learn path](docs/learn/README.md)** | 10 short chapters |
| **[Getting started](docs/guides/getting-started.md)** | Install + daily commands in 5 minutes |

Download **`koda`** + **`stdlib/`** from [GitHub Releases](https://github.com/CharmingBlaze/koda-compiler/releases). You do **not** need Go or LLVM to use release binaries.

---

## Documentation

**[Documentation hub →](docs/README.md)**

| Section | Links |
|---------|-------|
| Learn | [Beginner's guide](docs/beginners-guide.md) · [Learn path](docs/learn/README.md) |
| Guides | [Games](docs/guides/game-dev.md) · [Apps](docs/guides/applications.md) · [From C](docs/guides/from-c.md) · [Raylib](docs/guides/raylib.md) |
| Reference | [Language](language.md) · [CLI](docs/reference/cli.md) · [Stdlib](docs/stdlib/README.md) · [Builtins](docs/reference/builtins.md) |
| Help | [FAQ](docs/faq.md) · [Troubleshooting](docs/troubleshooting.md) · [Glossary](docs/glossary.md) |
| Concepts | [How it works](docs/concepts/README.md) |
| Style | [Documentation style guide](docs/STYLE-GUIDE.md) |

---

## Why Koda?

| | C | Koda |
|---|-----|------|
| Output | Native executable | Native executable |
| Syntax | Headers, macros | `struct`, `func`, `let`, GC |
| Modules | `.h` + link | `#include` + `import "@math"` |
| Graphics | Link Raylib yourself | Wrappers + `koda.json` |
| Iteration | Manual compile/link | `koda run`, `koda watch` |
| Ship | Binary + DLLs | `koda bundle` |

---

## Quick start

```bash
koda new mygame
cd mygame
koda run

# Or single file
koda run hello.koda
koda build hello.koda -o hello
koda bundle -o dist/MyGame
```

Templates: `hello` (default), `game` (text lander), `graphics` (Raylib).

---

## SDK downloads

Unpack an SDK zip so **`stdlib/`** and **`docs/`** sit next to **`koda`** and **`kodawrap`**.

| Platform | Artifact |
|----------|----------|
| Windows x64 | `koda-vX.Y.Z-sdk-windows-amd64.zip` |
| Linux x64 / ARM64 | `koda-vX.Y.Z-sdk-linux-*.zip` |
| macOS | `koda-vX.Y.Z-sdk-darwin-*.zip` |

---

## Examples

- [`examples/`](examples/) — demos and games
- `koda new lander --template game` — playable text game
- `koda new bounce --template graphics` — Raylib bouncing ball

---

## For contributors

Build from source requires Go 1.22+, Clang, and LLVM tools. See [CONTRIBUTING.md](CONTRIBUTING.md) and [docs/handoff.md](docs/handoff.md).

**End users should use release binaries**, not a source build.

---

[Changelog](CHANGELOG.md) · [Documentation hub](docs/README.md)
