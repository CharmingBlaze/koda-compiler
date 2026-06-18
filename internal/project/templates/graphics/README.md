# Graphics template (bouncing ball)

A **windowed graphics demo** using the **full Raylib wrapper** (548 functions) with optional **`koda.game`** helpers. Call **any** Raylib function via `use raylib;` — the template does not restrict the API.

## Run

```bash
koda run
```

## Source

```koda
use raylib;
use koda.game;
```

`koda.json` links `wrappers/raylib/wrapper.c` from the SDK — no project-local shim copy required.

## Troubleshooting

- **`koda setup raylib`** — writes `koda.json` native section for graphics
- **`koda doctor`** — SDK health check
- Undefined `InitWindow` / `game` → ensure `use raylib` and `use koda.game` appear before `main`

## See also

- [Game stdlib](../../../docs/stdlib/game.md)
- [Raylib guide](../../../docs/guides/raylib.md)
