# Game examples

| Project | API style | Notes |
|---------|-----------|-------|
| **pong** | `koda.game` | Canonical 2D — `loop`, `match`, hex colors |
| **brick-breaker** | `koda.game` | Breakout — parallel arrays, `#hex`, `match` |
| **mario64-studio** | Raw Raylib + `koda.camera` / `koda.ui` | Canonical 3D — `loop`, fusion camera, `#hex` |
| **mario64** | Raw Raylib + helpers | Same Peach's Castle course |
| **mario64-hilltop** | Raw Raylib + orbit cam | Rolling meadow + Bob-ombs |
| **koda64** | Raw Raylib + helpers | KODA 64 variant |
| **fps-arena** | Raw Raylib + `FirstPersonCamera` | FPS shooting gallery |

Run any project:

```powershell
cd examples/games/pong
..\..\koda.exe run
```

Or open in Koda Studio and press **F5**.
