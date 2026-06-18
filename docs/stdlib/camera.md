# `@camera` / `koda.camera`

Orbit and first-person cameras for 3D games over the **full Raylib wrapper**.

## Setup

```koda
use raylib;
use koda.camera;
```

Requires `"graphics": true` and `wrappers/raylib/wrapper.c` in `koda.json` (see [game](game.md)).

## First-person camera

```koda
use raylib;
use koda.camera;

let cam = FirstPersonCamera { x: 0.0, y: 1.6, z: 14.0, yaw: 3.14159, pitch: 0.0 };

func main() {
    InitWindow(1280, 720, "FPS");
    defer CloseWindow();
    DisableCursor();

    while (!WindowShouldClose()) {
        cam.fp_look();
        BeginDrawing();
        ClearBackground(asRaylib(colors.sky));
        cam.fp_begin();
        DrawGrid(20, 1.0);
        cam.fp_end();
        EndDrawing();
    }

    EnableCursor();
}
```

| Method | Purpose |
|--------|---------|
| `fp_look()` | Mouse look (`GetMouseDelta`) |
| `fp_begin()` / `fp_end()` | `BeginMode3D` / `EndMode3D` with `Camera3D` |
| `fp_fwd_x/y/z()` | Forward vector from yaw/pitch |
| `fp_reset()` | Reset position and angles |

## Orbit camera

```koda
let cam = OrbitCamera(target_entity, 13.0, 0.42, -0.55);
cam.update();
cam.begin();
// draw 3D scene
cam.end();
```

## See also

- [FPS Arena example](../../examples/games/fps-arena/README.md)
- [Vec3](vec3.md) — math helpers
- [Input](input.md) — keyboard/mouse constants
