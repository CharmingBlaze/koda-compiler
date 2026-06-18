# Bob-omb Hilltop — Mario 64 demo

A fresh 3D meadow course: red coins, orbiting bob-ombs, a windmill, bridge over a sky trench, and a star on the central hill.

## Features

- **Structs + free functions** — `Coin`, `Mario`, `Bobomb`; `tick_bobomb`, `try_collect`, `ground_height`
- **`KEY_W` / `KEY_SPACE`** from `use raylib` (preferred over `Key.*` aliases)
- **`for (let coin of coins)`** — array iteration
- **`OrbitCamera`** from `@camera`
- **`main()` sections** — setup / update / draw comments for learning the game loop

## Controls

| Input | Action |
|-------|--------|
| WASD | Move (camera-relative) |
| Space | Jump |
| Right mouse | Orbit camera |
| Mouse wheel | Zoom |
| R | Restart |

Collect **8 red coins**, climb the hill, touch the **star**.

## Run

```bash
cd examples/games/mario64-hilltop
koda run
```

## Koda Studio

From repo root:

```powershell
.\koda-ide\run-koda-studio.ps1 examples\games\mario64-hilltop
```

Press **F5** in Studio to run.
