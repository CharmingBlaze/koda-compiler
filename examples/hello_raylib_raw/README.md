# Hello Raylib (raw API)

Minimal window using the **full Raylib wrapper** — official beginner sample.

## Run

```bash
cd examples/hello_raylib_raw
koda run
```

## Source

```koda
use raylib;
use koda.color;

func main() {
    InitWindow(800, 600, "Hello Koda");
    defer CloseWindow();
    SetTargetFPS(60);

    while (!WindowShouldClose()) {
        BeginDrawing();
        ClearBackground(asRaylib(colors.sky));
        DrawText("Hello Koda", 20, 20, 30, asRaylib(colors.black));
        EndDrawing();
    }
}
```

All 548 Raylib functions are available — see `wrappers/raylib/api_reference.md`.
