# Launch post drafts — Koda v0.5.0

Use one of these for [dev.to](https://dev.to) and a shorter variant for [Hacker News](https://news.ycombinator.com/submit).  
**Repo:** https://github.com/CharmingBlaze/koda-compiler  
**Release:** https://github.com/CharmingBlaze/koda-compiler/releases/tag/v0.5.0

**Release page copy:** see [`docs/release-v0.5.0.md`](release-v0.5.0.md) for GitHub Release announcement text.

**Post timing:** dev.to first, wait a few hours for engagement, then HN.

**Suggested title:** I built a beginner-friendly compiled language for games — one zip, native binaries, no CMake

**Suggested tags:** `go`, `gamedev`, `opensource`, `programming`, `showdev`

---

I wanted to teach a kid to make games without handing them C++, CMake, and a week of toolchain setup. So I built **Koda** — a small language that compiles to **native executables**, ships as **one SDK zip**, and feels closer to JavaScript than to `malloc`.

No Go runtime on the player's machine. No Python. No "install LLVM first." Unzip, run `koda`, ship a `.exe` (or binary on macOS/Linux).

### What you get in the zip

- **`koda`** CLI — `new`, `run`, `build`, `check`, `fmt`, `test`, `doctor`
- **Embedded Clang + runtime** — release builds bundle the compiler; beginners don't configure a toolchain
- **`@game` stdlib** — Raylib-backed API so your first program can open a window and draw rectangles
- **Koda Studio** — desktop IDE (Wails) with docs, templates, and F5 to run
- **Examples** — bouncing ball, lunar lander, and a Mario 64-style 3D demo (`mario64-studio`)

### Hello, native game

```koda
#include "wrappers/raylib_shim/raylib.koda"
#include "@game"

struct Mario { x, y, speed, health }

func main() {
    game.open(800, 600, "My Game");
    game.fps(60);

    let player = Mario { x: 400, y: 300, speed: 220, health: 100 };

    while (game.running()) {
        let dt = game.delta();
        if (game.keyDown(Key.Right)) {
            player.x = player.x + player.speed * dt;
        }
        game.begin();
        game.clear(colors.dark);
        game.rect(player.x, player.y, 32, 32, colors.white);
        game.end();
    }
}
```

`koda run` compiles this to a native binary with GC, structs, and a game loop — not an interpreted script.

### Why not just use C++ / C# / Godot?

- **vs C++:** Same end result (native binary), but no headers, linking puzzles, or UB in everyday game logic
- **vs C# / Java:** No VM requirement for players; one folder to zip and share
- **vs scripting in Godot:** You're learning a *language* that stands alone — CLI tools, small apps, and games share one toolchain

Koda is intentionally small. It's not trying to replace Rust or Zig. It's the **on-ramp**: compiled, typed enough to catch mistakes, friendly enough that `struct Player { x, y, health }` is valid syntax.

### What's in v0.5.0

Recent work focused on **beginner footguns** and **stdlib gaps**:

- **Switch** no longer falls through by default (`fallthrough;` when you really want C behavior)
- **Removed `===` / `!==`** — `koda fmt` migrates old code safely
- **Truthy lint** — `koda check` warns if you write `if (myArray)` expecting JavaScript truthiness
- **Optional struct fields** — `struct Node { value, next? }`
- **`padStart` / `padEnd`**, **`flat` / `flatMap`**, **`readDir`**
- **mario64-studio** — coins, orbit camera, `match` game states, struct methods

### Try it (2 minutes)

1. Download the SDK for your OS from [GitHub Releases v0.5.0](https://github.com/CharmingBlaze/koda-compiler/releases/tag/v0.5.0)
2. Unzip, run `koda doctor`
3. `koda new bounce --template graphics && cd bounce && koda run`

Or open `examples/games/mario64-studio` in Koda Studio and press F5.

### Under the hood (for the curious)

- **Compiler:** Go frontend (lexer → parser → sema → LLVM IR → Clang)
- **Runtime:** C11 GC'd value model (`libkoda_runtime.a`)
- **Graphics:** Raylib via generated wrappers (`koda wrap`, `koda setup raylib`)

The repo is MIT. Issues and PRs welcome — especially templates, stdlib, and docs for absolute beginners.

**Links**

- GitHub: https://github.com/CharmingBlaze/koda-compiler
- Changelog: https://github.com/CharmingBlaze/koda-compiler/blob/main/CHANGELOG.md
- Beginner guide: https://github.com/CharmingBlaze/koda-compiler/blob/main/docs/beginners-guide.md

If you try it, I'd love to hear what broke first — that's the roadmap.

---

## Hacker News — submission

**Title (pick one):**

- `Show HN: Koda – beginner-friendly compiled language for games, one SDK zip, native binaries`
- `Show HN: I built a JS-like language that compiles to native games (no LLVM install for users)`

**URL:** https://github.com/CharmingBlaze/koda-compiler/releases/tag/v0.5.0

**First comment (post immediately after submit):**

Author here. Koda is a small compiled language aimed at beginners who want **native games/apps** without C++ toolchain pain.

**What works today:** one SDK zip (compiler + stdlib + Raylib + IDE), `koda run` → native binary, `@game` API, structs/enums/match, GC, `koda new --template graphics`.

**v0.5.0** just shipped: safer switch semantics, truthy lint, optional struct fields, mario64-studio demo.

Quick start: download the zip from the release page, `koda doctor`, `koda new bounce --template graphics`, `koda run`.

Happy to answer questions about the compiler pipeline (Go → LLVM → Clang), GC design, or why I'm not targeting "replace Rust."

---

## One-liner bios (Twitter / Mastodon / Bluesky)

> Shipped Koda v0.5.0 — a beginner-friendly compiled language for games. One SDK zip, native binaries, Raylib `@game` API, no CMake/LLVM for users. https://github.com/CharmingBlaze/koda-compiler/releases/tag/v0.5.0
