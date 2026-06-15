# {{name}}

A **windowed graphics demo** using Raylib through `@game` and the bundled `wrappers/raylib_shim` bridge.

## Run

```bash
koda doctor
koda run
```

The project manifest sets `"graphics": true` — Koda applies platform Raylib link flags automatically. You do **not** need to set `KODA_LINKFLAGS` for a standard SDK install.

## Source layout

```koda
#include "../wrappers/raylib_shim/raylib.koda"
#include "@game"
```

`@game` requires the shim include first. If you see undefined-variable errors for `drawline`, `getmousex`, etc., refresh the shim:

```bash
koda setup raylib
```

## Build & bundle

```bash
koda build -o {{name}}
koda bundle -o dist/{{name}}
```

See `docs/guides/raylib.md` and `docs/stdlib/game.md` in the Koda SDK.
