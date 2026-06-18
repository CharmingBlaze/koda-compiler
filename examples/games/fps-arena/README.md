# FPS Arena

First-person shooter demo: mouse look, WASD movement, click to shoot red targets in a 3D arena.

## Controls

| Input | Action |
|-------|--------|
| Mouse | Look around |
| WASD | Move |
| Left click | Shoot |
| R | Restart |

Clear all targets before your health runs out.

## Run

```bash
cd examples/games/fps-arena
koda run
```

## Koda Studio

```powershell
.\koda-ide\run-koda-studio.ps1 examples\games\fps-arena
```

Press **F5** in Studio to run.

## Features

- **Full Raylib API** via `use raylib;` (548 functions)
- **`FirstPersonCamera`** from `koda.camera`
- **`KEY_*`** constants from `use raylib` (WASD, R, mouse via Raylib)
- Struct arrays, aim-cone hit detection, cover pillars
