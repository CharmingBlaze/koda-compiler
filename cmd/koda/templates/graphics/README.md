# {{name}}

A **windowed graphics demo** using Raylib through the bundled `wrappers/raylib_shim` bridge.

## Prerequisites

You need **Raylib** available to the linker (headers + library). Common options:

- Use a **Koda SDK zip** that includes a vendored Raylib stage, and set `KODA_RAYLIB_STAGE` if needed.
- Install Raylib on your system and set linker flags (examples below).

## Windows (PowerShell)

```powershell
$env:KODA_LINKFLAGS = "-lraylib -lopengl32 -lgdi32 -lwinmm"
koda run
```

## Linux / macOS

```bash
export KODA_LINKFLAGS="$(pkg-config --libs --cflags raylib 2>/dev/null || echo '-lraylib')"
koda run
```

Or add `linkflags` to `koda.json` once you know the flags for your machine.

## Build & bundle

```bash
koda build -o {{name}}
koda bundle -o dist/{{name}}
```

See `docs/wrappers.md` and `docs/guides/raylib.md` in the Koda SDK for full Raylib setup.
