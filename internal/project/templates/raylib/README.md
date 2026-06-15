# {{name}}

Windowed demo using the **full Raylib wrapper** (`548` functions from `raylib.h`).

## Quick start

```bash
koda run
koda doctor   # if linking fails
```

## Include

```koda
#include "@raylib"
```

The bindings live in the Koda SDK at `wrappers/raylib/` — you do not copy them into this project. `koda.json` points at `wrappers/raylib/wrapper.c`; the compiler resolves it from the SDK when missing locally.

## API docs

- `koda doc wrapper @raylib`
- SDK: `wrappers/raylib/api_reference.md` and `wrappers/raylib/docs/index.html`

## Beginner API

For a smaller surface (window, 2D draw, input), use `koda new mygame --template graphics` with `@game` instead.
