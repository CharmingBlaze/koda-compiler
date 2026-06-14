# Example games

## Text-only (no Raylib)

Runs anywhere **`koda`** runs:

```bash
koda run examples/games/lunar_lander_text.koda
```

## Raylib (graphics)

Raylib samples need generated bindings plus linker flags. See **`tests/raylib_wrapgen_smoke.koda`** and **`scripts/ci-wrapgen-raylib.sh`** for the usual **`KODA_WRAPPERS`**, **`KODA_NATIVE_SOURCES`**, and **`KODA_LINKFLAGS`** setup.

After wrappers are built:

```bash
export KODA_WRAPPERS=/path/to/wrappers/raylib
export KODA_NATIVE_SOURCES=/path/to/wrappers/raylib/wrapper.c
export KODA_LINKFLAGS="$(pkg-config --libs --cflags raylib)"   # platform-specific
koda run demos/demo_3d.koda    # or another raylib-backed demo under demos/
```
