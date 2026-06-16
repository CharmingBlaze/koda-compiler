# Mario 64 — Koda Studio demo

Peach's Castle grounds in 3D — a beginner-friendly Mario 64-style platformer built with modern Koda.

## What this demo shows

- **`const`** for screen size, keys, and physics tuning
- **`enum GameState`** + **`match`** for win/lose overlays
- **String interpolation** — `"Score: {score}   Lives: {lives}"`
- **Struct methods** — `coin.show()`, `coin.pickup(player)` on `Coin`
- **`OrbitCamera`** from `@camera` — right-mouse orbit, wheel zoom

## Controls

| Input | Action |
|-------|--------|
| WASD | Move (camera-relative) |
| Space | Jump |
| Right mouse drag | Orbit camera |
| Mouse wheel | Zoom |
| R | Restart |

Collect **5 coins**, climb the hill, touch the **star**.

## Run

**Koda Studio:** open this folder as workspace, press **F5**.

**Terminal:**

```bash
koda run
```

## Open in Koda Studio (from repo root)

```powershell
.\koda-ide\run-koda-studio.ps1 examples\games\mario64-studio
```
