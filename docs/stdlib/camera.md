# `@camera` — orbit camera for 3D games

**Include:** `#include "@camera"` (after the Raylib shim)

---

## Setup

```koda
#include "wrappers/raylib_shim/raylib.koda"
#include "@camera"
```

---

## Quick start

```koda
let camera = OrbitCamera {
    target: player,
    distance: 13.0,
    pitch: 0.42,
    yaw: -0.55
};

func update(dt) {
    camera.update();
}

func draw() {
    camera.begin();
    // draw world in 3D...
    camera.end();
}
```

Or use the `orbit()` helper with a config object:

```koda
let camera = orbit({ target: player, distance: 13 });
```

---

## `OrbitCamera`

Follows a **target entity** each frame. The target must have **x, y, z** as its first three struct fields (Mario `{ x, y, z, … }` works out of the box).

| Field / method | Description |
|----------------|-------------|
| `target` | Entity to orbit (struct reference) |
| `distance`, `yaw`, `pitch` | Orbit parameters |
| `look_offset` | Added to target y for look-at (default `0.5`) |
| `update()` | Mouse orbit + wheel zoom |
| `begin()` | Start 3D mode (`beginmode3d`) |
| `end()` | End 3D mode (`endmode3d`) |
| `reset()` | Restore default yaw, pitch, distance |
| `yaw` | Exposed for camera-relative movement (`sin(camera.yaw)`) |

Defaults: `distance: 13`, `yaw: -0.55`, `pitch: 0.42`, right-mouse orbit, wheel zoom.

---

## Related

- [@game](game.md) — window and 2D helpers
- [Raylib guide](../guides/raylib.md)
