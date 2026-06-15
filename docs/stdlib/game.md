# @game — beginner game API

Thin wrapper over the Raylib shim. Ships in `stdlib/game.koda`.

---

## Setup

Graphics projects include the shim and set `"graphics": true` in `koda.json`:

```koda
#include "wrappers/raylib_shim/raylib.koda"
#include "@game"
```

The shim **must** be included before `@game`. If compile errors mention undefined shim names (`drawline`, `getmousex`, …), run **`koda setup raylib`** in the project to refresh the shim.

```json
{
  "native": {
    "sources": ["wrappers/raylib_shim/wrapper.c"],
    "graphics": true
  }
}
```

---

## Dot notation

Koda uses JavaScript-style field and method access:

```koda
game.open(640, 480, "Demo");
game.circle(ball.x, ball.y, ball.radius, colors.white);

let cam = {
    yaw: 0.0,
    update: func() { this.yaw = this.yaw + 0.01; }
};
cam.update();
```

Structs work the same way: `player.x`, `player.speed`. Object methods use **`this`**. Optional chaining: `obj?.field`.

**Naming tip:** avoid struct type names that match your variable names (e.g. use `struct Mario` with `let player = Mario { ... }`, not `struct Player` with `let player`).

---

### `draw` alias

`@game` also exports a `draw` object with the same drawing helpers:

```koda
draw.text("Score: {score}", 20, 20, 24, colors.white);
draw.rect(x, y, w, h, colors.red);
draw.line(x1, y1, x2, y2, colors.white);
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
| `game.width()` | Window width in pixels |
| `game.height()` | Window height in pixels |
| `game.setTitle(title)` | Change window title |
| `game.fpsCounter()` | Current FPS |

### Drawing

| Function | Description |
|----------|-------------|
| `game.begin()` | Start drawing |
| `game.end()` | End drawing |
| `game.clear(color)` | Clear background |
| `game.text(msg, x, y, size, color)` | Draw text |
| `draw.text(msg, x, y, size, color)` | Same (alias) |
| `game.rect(x, y, w, h, color)` | Filled rectangle |
| `game.circle(x, y, r, color)` | Filled circle |
| `game.line(x1, y1, x2, y2, color)` | Line segment |
| `game.circleLines(x, y, r, color)` | Circle outline |
| `game.rectLines(x, y, w, h, color)` | Rectangle outline |
| `game.loadImage(path)` | Load image; returns texture handle |
| `game.drawImage(tex, x, y, color)` | Draw loaded texture |
| `game.unloadImage(tex)` | Free texture |

### Keyboard input

| Function | Description |
|----------|-------------|
| `game.keyDown(key)` | Key held |
| `game.keyPressed(key)` | Key pressed this frame |

### Mouse input

| Function | Description |
|----------|-------------|
| `game.mouseX()` / `game.mouseY()` | Cursor position |
| `game.mouseDown(btn)` | Button held |
| `game.mousePressed(btn)` | Button pressed this frame |
| `game.mouseWheel()` | Scroll wheel delta |

### Runtime

| Function | Description |
|----------|-------------|
| `game.setGcBudget(ms)` | Incremental GC budget per frame |

---

## Constants

### `Key`

`Key.Left`, `Key.Right`, `Key.Up`, `Key.Down`, `Key.Space`, `Key.Escape`

### `Mouse`

`Mouse.Left` (0), `Mouse.Right` (1), `Mouse.Middle` (2)

### `Color`

`colors.dark`, `colors.white`, `colors.red`, `colors.green`, `colors.blue`

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

        if (ball.x < ball.radius || ball.x > game.width() - ball.radius) {
            ball.velx = -ball.velx;
        }
        if (ball.y < ball.radius || ball.y > game.height() - ball.radius) {
            ball.vely = -ball.vely;
        }

        game.begin();
        game.clear(colors.dark);
        game.circle(ball.x, ball.y, ball.radius, colors.white);
        game.text("FPS: {game.fpsCounter()}", 8, 8, 18, colors.white);
        game.end();
        game.setGcBudget(0.5);
    }
}
```

---

## Related

- [Game development guide](../guides/game-dev.md)
- [@input](input.md) — input helpers without the `game` object
- [Raylib shim (advanced)](../guides/raylib.md)
- Source: `stdlib/game.koda`
