# {{name}}

Classic **Pong** — two-player paddle game using Raylib and the `@game` helper.

## Controls

| Player | Keys |
|--------|------|
| Left paddle | **W** / **S** |
| Right paddle | **↑** / **↓** |
| Serve / restart | **Space** |

First to **11** points wins.

## Run

```bash
cd {{name}}
koda run
```

If linking fails, set Raylib flags (see the [graphics template README](../graphics/README.md) or `docs/guides/raylib.md`).

## Build

```bash
koda build -o {{name}}
```
