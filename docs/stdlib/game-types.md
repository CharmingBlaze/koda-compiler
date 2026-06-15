# Built-in game types

Koda ships native **vec2**, **vec3**, **color**, **rect**, and **box** constructors for game code. They support component access (`.x`, `.y`, `.z`) and vector math with `+`, `-`, `*`, and `+=`.

## vec2 / vec3

```koda
let pos = vec3(0, 1, 8);
let velocity = vec3(1, 0, 0);
let dt = 0.016;
let speed = 5.0;

pos = pos + velocity * speed * dt;
pos += vec3(0, 0.5, 0);

print(pos.x, pos.y, pos.z);
```

Use **vec2** for 2D positions and UVs:

```koda
let p = vec2(320, 240);
let scaled = p * 2;
```

## Struct fields

Store vectors on structs and use dot notation:

```koda
struct Mario { pos, vel_y, h, r, grounded }

let player = Mario {
    pos: vec3(0, 1, 8),
    vel_y: 0.0,
    h: 1.2,
    r: 0.45,
    grounded: true
};

player.pos = player.pos + movement * speed * dt;
drawcube(player.pos.x, player.pos.y, player.pos.z, 1, 1, 1, colors.red);
```

## color(r, g, b)

Builds an RGBA color object with `.r`, `.g`, `.b`, `.a`, and `.packed` (Raylib-ready integer):

```koda
let c = color(255, 128, 0);
clearbackground(c.packed);
```

Named palette colors live on **`colors`** (auto-injected):

```koda
clearbackground(colors.sky);
drawrectangle(0, 0, 100, 100, colors.grass);
```

Use **`rgb()` / `rgba()`** when you need a packed integer directly.

## rect / box

```koda
let r = rect(10, 20, 100, 50);   // x, y, w, h
let b = box(vec3(0, 1, 0), vec3(2, 3, 4));  // center, size
```

## @vec2 / @vec3 modules

`#include "@vec3"` still provides helpers like `dot`, `cross`, `normalize`, and `lerp` for advanced math.
