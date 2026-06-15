# Example games

## Text-only (no Raylib)

Runs anywhere **`koda`** runs:

```bash
koda run examples/games/lunar_lander_text.koda
```

## Koda64 / Mario 64 (Mario 64-style 3D)

| Project | Path |
|---------|------|
| **Mario 64** (new) | `examples/games/mario64/` |
| Koda64 (original) | `examples/games/koda64/` |

```bash
cd examples/games/mario64
koda run
```

Open in **Koda Studio**: run `koda-ide/run-koda-studio.ps1` with the project path, or pass the folder as the first argument after build.

## Graphics with `@game` (beginner)

```bash
koda new bounce --template graphics
cd bounce
koda run
```

Uses `wrappers/raylib_shim` + `@game`. Refresh stale shims with **`koda setup raylib`**.

## Raylib (full wrapper / advanced)

For hundreds of Raylib functions, use **`koda setup raylib --full`** or **`koda wrap install raylib --project`**, then `#include "@raylib"`.

See **`docs/guides/raylib.md`** and **`docs/guides/wrapping-libraries.md`**.