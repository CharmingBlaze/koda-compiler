# {{name}}

Windowed demo using the **full Raylib wrapper** (548 functions from `raylib.h`).

## Quick start

```bash
koda run
koda doctor   # if linking fails
```

## Import

```koda
use raylib;
use koda.color;   // asRaylib() for draw colors
```

Bindings live in the Koda SDK at `wrappers/raylib/`. `koda.json` points at `wrappers/raylib/wrapper.c`; the compiler resolves it from the SDK.

## API docs

- `wrappers/raylib/api_reference.md`
- `wrappers/raylib/docs/index.html`

## Beginner helpers

For a smaller API surface (window loop, 2D draw), use `koda new mygame --template graphics` with `use koda.game` on top of the same full wrapper.
