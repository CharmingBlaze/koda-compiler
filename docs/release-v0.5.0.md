# GitHub Release v0.5.0 — announcement text

Paste into https://github.com/CharmingBlaze/koda-compiler/releases/edit/v0.5.0 (or create if the tag release has no body).

---

## Koda v0.5.0 — beginner-friendly native games, one zip

**Koda** is a small compiled language for games and desktop apps. Download one SDK zip, run `koda`, ship a native binary — no Go, no Python, no LLVM install, no CMake.

This release focuses on **beginner safety** and **stdlib gaps** called out by early users and compiler review.

### Highlights

- **Safer `switch`** — cases no longer fall through by default; use `fallthrough;` when you want C-style chaining (same mental model as `match`).
- **No more `===` / `!==`** — lexer rejects them; `koda fmt` migrates old sources safely (outside strings/comments).
- **Truthy lint** — `koda check` warns on `if (myArray)` and similar JS-isms that behave differently in Koda.
- **Optional struct fields** — `struct Node { value, next? }` (omit at construction → `null`).
- **String & array methods** — `.padStart` / `.padEnd`, `.flat` / `.flatMap`, `readDir` alias.
- **mario64-studio demo** — 3D platformer sample with orbit camera, coins, and `match` game states.
- **Linux ARM64 CI** — release-quality builds on `ubuntu-24.04-arm`.

### Try it in 2 minutes

1. Download the SDK zip for your OS below.
2. `koda doctor`
3. `koda new bounce --template graphics && cd bounce && koda run`

Or open `examples/games/mario64-studio` in **Koda Studio** and press **F5**.

### Editor support

Syntax highlighting for VS Code / Cursor: [`vscode-koda/`](https://github.com/CharmingBlaze/koda-compiler/tree/main/vscode-koda) (open folder → F5, or package a `.vsix`).

### What Koda looks like

```koda
use koda.game;

struct Mario { x, y, speed, health }

func main() {
    game.open(800, 600, "My Game");
    let player = Mario { x: 400, y: 300, speed: 220, health: 100 };
    while (game.running()) {
        let dt = game.delta();
        if (game.keyDown(Key.Right)) { player.x = player.x + player.speed * dt; }
        game.begin();
        game.clear(colors.dark);
        game.rect(player.x, player.y, 32, 32, colors.white);
        game.end();
    }
}
```

### Full changelog

<details>
<summary>Click to expand</summary>

### Added

- **`fallthrough;` keyword** — opt-in C-style switch chaining; default is no fall-through (same as `match`).
- **Optional struct fields** — `struct Node { value, next? }` (`?` = may omit at construction → `null`; equivalent to `next = null`).
- **String methods** — `.padStart(n, c?)`, `.padEnd(n, c?)` (C runtime).
- **Array methods** — `.flat(depth?)`, `.flatMap(callback)` (flatMap via map + flat in codegen).
- **`readDir(path)`** — alias of `listDir`; returns **entry names only** (not full paths).
- **Truthy lint** — `koda check` warns on `if (arr)`, `if ({})`, and known array/struct variables (differs from JavaScript).
- **Struct field hints** — `did you mean 'field'?` on typos (Levenshtein distance 1–2).
- **`scripts/resolve-raylib-stage.ps1`** — prefers `third_party/raylib_static/stage`, falls back to `raylib_lib/`.
- **mario64-studio example** — `examples/games/mario64-studio/`.
- **CI** — `ubuntu-24.04-arm` in main workflow matrix.

### Changed

- **Removed `===` / `!==`** — lexer rejects them; `koda fmt` migrates old sources (string/comment-safe rewrite).
- **Switch codegen** — cases no longer fall through by default; `break` in a case is a harmless no-op.
- **For-loop init** — `for (let i = 0, j = 10; …)` allowed (second `let` optional).
- **`koda fmt`** — rewrites legacy equality outside strings and comments.
- **`koda watch`** — project-relative paths in rebuild messages.
- **`koda doctor`** — `Fix:` hints for install-dir and temp-folder write failures.
- **`gcCollect()`** — deprecated; `gc()` is canonical (`koda check` warns).
- **Documentation** — struct defaults, `@str`, `randomChoice`, `listDir`/`readDir`, truthy/falsy in learn guide; Raylib layout in `ARCHITECTURE.md`.

### Fixed

- **API runtime hardening tests** — set `KODA_HOME` to repo root so builds link local `libkoda_runtime.a`.
- **Native extern forward refs** — `// koda:extern` names visible before their `let` (struct methods calling shims).
- **GC write barrier** — `koda_array_remove_at` uses `gc_write_barrier` on element shifts.
- **Sema** — optional `?` fields register the same null default as `field = null`.

</details>

---

**Links:** [Beginner's guide](https://github.com/CharmingBlaze/koda-compiler/blob/main/docs/beginners-guide.md) · [START_HERE](https://github.com/CharmingBlaze/koda-compiler/blob/main/START_HERE.md) · [Changelog](https://github.com/CharmingBlaze/koda-compiler/blob/main/CHANGELOG.md)

Issues and PRs welcome — especially templates, stdlib, and docs for absolute beginners.
