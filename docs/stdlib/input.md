# @input — keyboard and mouse helpers

Standalone input functions over the Raylib shim. Ships in `stdlib/input.koda`.

For windowed games, prefer [`@game`](game.md). Use `@input` when you want input helpers without the `game` object namespace.

---

## Setup

Include the shim first:

```koda
#include "wrappers/raylib_shim/raylib.koda"
#include "@input"
```

---

## API

| Function | Description |
|----------|-------------|
| `keyDown(k)` | Key held |
| `keyPressed(k)` | Key pressed this frame |
| `mousePos()` | `{ x, y }` cursor position |
| `mouseButton(b)` | Mouse button held (0=left, 1=right, 2=middle) |
| `mousePressed(b)` | Button pressed this frame |
| `mouseWheel()` | Scroll wheel delta |

Use `Key` / `Mouse` constants from `@game` if you include both modules, or define your own key codes.

---

## Example

```koda
#include "wrappers/raylib_shim/raylib.koda"
#include "@input"

func main() {
    initwindow(400, 300, "Input demo");
    while (!windowshouldclose()) {
        let pos = mousePos();
        if (mousePressed(0)) {
            print("click at", pos.x, pos.y);
        }
        begindrawing();
        clearbackground(colors.dark);
        drawcircle(pos.x, pos.y, 6, colors.white);
        enddrawing();
    }
    closewindow();
}
```

---

## Related

- [@game](game.md)
- Source: `stdlib/input.koda`
