# `@input` / `koda.input`

Keyboard and mouse helpers over the **full Raylib wrapper**.

## Setup

```koda
use raylib;
use koda.input;
```

Prefer **`koda.game`** for windowed games with a built-in loop. Use **`koda.input`** when you manage the window yourself but want `Key.*` / `Mouse.*` constants and thin wrappers.

## API

```koda
if (keyDown(Key.W)) { ... }
if (keyPressed(Key.Space)) { ... }
if (mousePressed(Mouse.Left)) { ... }
let pos = mousePos();  // { x, y } from GetMousePosition()
```

| Function | Raylib |
|----------|--------|
| `keyDown(k)` | `IsKeyDown` |
| `keyPressed(k)` | `IsKeyPressed` |
| `mouseButton(b)` | `IsMouseButtonDown` |
| `mousePressed(b)` | `IsMouseButtonPressed` |
| `mouseWheel()` | `GetMouseWheelMove` |

## Constants

`Key` and `Mouse` map to Raylib enums (`KEY_W`, `MOUSE_BUTTON_LEFT`, …). Load **`use raylib`** before **`use koda.input`**.

## Example

```koda
use raylib;
use koda.input;

func main() {
    InitWindow(800, 600, "Input");
    defer CloseWindow();

    while (!WindowShouldClose()) {
        if (keyPressed(Key.Escape)) { break; }
        BeginDrawing();
        ClearBackground(asRaylib(colors.dark));
        if (keyDown(Key.W)) {
            DrawText("W held", 20, 20, 20, asRaylib(colors.white));
        }
        EndDrawing();
    }
}
```

Use **`use koda.color`** for `asRaylib()` when passing colors to full-wrapper draw calls.

## See also

- [Game module](game.md)
- [Game dev guide](../guides/game-dev.md)
