# FAQ — frequently asked questions

Short answers to common questions about Koda.

---

## General

### What is Koda?

A language that **compiles to native executables** for games and applications. Syntax blends C-style control flow with JavaScript-style objects.

### Is Koda interpreted?

No. `koda run` compiles to a native binary (often a temp file). Users run your built executable without installing Koda.

### Do I need Go, LLVM, Python, or C++ to use Koda?

**No** for release SDK binaries. Download the **SDK zip**, unzip, run `koda`. The compiler embeds Clang and the runtime.

You need Go/LLVM **only** if you modify the Koda compiler itself ([CONTRIBUTING.md](../CONTRIBUTING.md)).

Koda is **not** Python — it compiles to a native executable. Your players do not install Python or Koda to run your game.

### Do I need Visual Studio or CMake?

**No** to get started with the SDK. Graphics templates use bundled Raylib and `"graphics": true` in `koda.json`.

### Is Koda case-sensitive?

Keywords and builtins are **case-insensitive**. `print`, `Print`, and `PRINT` are equivalent.

### What platforms are supported?

Windows x64, Linux x64/ARM64, macOS Intel and Apple Silicon — see [GitHub Releases](https://github.com/CharmingBlaze/koda-compiler/releases).

---

## Language

### C or JavaScript?

Closer to **C** for structs, enums, and control flow. Objects and closures feel like **JavaScript**. Not a C superset — no headers, no manual `malloc` for everyday code (GC handles gameplay allocations).

### Does Koda have classes?

Use **structs** + functions, or **objects** with methods. No `class` keyword.

### How do I import code?

`#include "path.koda"` merges files. `import "@math"` or `import "src/lib.koda"` loads modules.

### Why does `delete` not work as a function name?

`delete` is a **reserved keyword** (property deletion). Use `deletefile` or `io.remove` for files.

### String escapes?

`\"`, `\\`, `\n`, `\t`, `\r`, `\'` in double-quoted strings.

---

## Tools

### Where must `stdlib/` live?

Next to the `koda` executable (SDK zip layout). `koda doctor` warns if missing.

### `koda run` vs `koda build`?

`run` — quick compile + execute (temp output). `build` — permanent executable you can ship.

### How do I run tests?

`koda test` with optional `-v`, `--failfast`, `-run <pattern>`.

### REPL or one-liner?

`koda repl` for interactive use; `koda eval 'print(1)'` for scripts.

### Pass arguments to my game?

`koda run game.koda -- --level 3`

### Graphics / Raylib link errors?

Run `koda doctor` and fix any FAIL lines. Graphics templates set `"graphics": true` in `koda.json` so link flags are applied automatically. See [Game dev guide](guides/game-dev.md).

---

## Games

### Best way to learn game dev?

1. Text game: `koda new lander --template game`
2. [Game dev guide](guides/game-dev.md)
3. Graphics: `koda new bounce --template graphics`

### Frame timing?

`deltatime()` for per-frame delta. `programtime()` for elapsed since start. `import "@timer"` for cooldowns.

### Random numbers reproducible?

`randomseed(n)` before `random` / `randomint` calls.

### GC stutter in games?

Call `gcframestep()` once per frame; reduce per-frame allocations.

---

## Interop

### Call C libraries?

Use `kodawrap` / `koda wrap` to generate bindings. See [wrappers.md](wrappers.md).

### Ship to players?

`koda bundle -o dist/GameName` — see [distribution](guides/distribution.md).

---

## More help

- [Troubleshooting](troubleshooting.md)
- [Glossary](glossary.md)
- [Beginner's guide](beginners-guide.md)
