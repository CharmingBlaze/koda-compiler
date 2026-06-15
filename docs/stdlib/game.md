# @game — beginner game API

Thin wrapper over the Raylib shim. Ships in `stdlib/game.koda`.

---

## Setup

Graphics projects include the shim and set `"graphics": true` in `koda.json`:

```koda
#include "wrappers/raylib_shim/raylib.koda"
#include "@game"
```

---

## API

### Window and loop

| Function | Description |
|----------|-------------|
| `game.open(w, h, title)` | Open window |
| `game.close()` | Close window |
| `game.running()` | `true` while window should stay open |
| `game.fps(n)` | Set target FPS |
| `game.delta()` | Frame delta time in seconds |

### Drawing

| Function | Description |
|----------|-------------|
| `game.begin()` | Start drawing |
| `game.end()` | End drawing |
| `game.clear(color)` | Clear background |
| `game.text(msg, x, y, size, color)` | Draw text |
| `game.rect(x, y, w, h, color)` | Filled rectangle |
| `game.circle(x, y, r, color)` | Filled circle |

### Input

| Function | Description |
|----------|-------------|
| `game.keyDown(key)` | Key held |
| `game.keyPressed(key)` | Key pressed this frame |

### Runtime

| Function | Description |
|----------|-------------|
| `game.setGcBudget(ms)` | Incremental GC budget per frame |

---

## Constants

### `Key`

`Key.Left`, `Key.Right`, `Key.Up`, `Key.Down`, `Key.Space`, `Key.Escape`

### `Color`

`Color.dark`, `Color.white`, `Color.red`, `Color.green`, `Color.blue`

---

## Example

```koda
#include "wrappers/raylib_shim/raylib.koda"
#include "@game"

struct Ball { x, y, velx, vely, radius }

func main() {
    game.open(640, 480, "Bounce");
    game.fps(60);

    let ball = Ball { x: 320, y: 240, velx: 140, vely: 100, radius: 20 };

    while (game.running()) {
        let dt = game.delta();
        ball.x = ball.x + ball.velx * dt;
        ball.y = ball.y + ball.vely * dt;

        game.begin();
        game.clear(Color.dark);
        game.circle(ball.x, ball.y, ball.radius, Color.white);
        game.end();
        game.setGcBudget(0.5);
    }
}
```

---

## Related

- [Game development guide](../guides/game-dev.md)
- [Raylib shim (advanced)](../guides/raylib.md)
- Source: `stdlib/game.koda`
