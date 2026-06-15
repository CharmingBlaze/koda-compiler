# @color — color utilities

Raylib-ready colors for games. **Beginners** use the global **`colors`** palette and `rgb()` / `rgba()` — no hex required.

```koda
clearbackground(colors.sky);
drawrectangle(0, 0, 100, 100, colors.grass);
drawcircle(50, 50, 20, rgb(255, 100, 50));
drawline(0, 0, 100, 100, rgba(255, 255, 255, 128));
```

## Built-in palette (`colors.*`)

| Name | Typical use |
|------|-------------|
| `colors.white`, `colors.black` | Text, outlines |
| `colors.red`, `colors.green`, `colors.blue` | Player, UI, accents |
| `colors.yellow`, `colors.gold`, `colors.orange` | Coins, stars |
| `colors.sky`, `colors.grass`, `colors.dirt`, `colors.forest` | Outdoor scenes |
| `colors.brown`, `colors.gray`, `colors.dark` | Enemies, castle, backgrounds |

## Custom colors

```koda
let grass = rgb(34, 139, 34);
let sky = rgba(255, 216, 168, 255);
let c = color(255, 128, 0);   // object with .r .g .b .packed
clearbackground(c.packed);
```

## Advanced

Hex literals (`0xRRGGBBAA`) still work when you need them. Import `@color` for `hsv()`, `lerp()`, and `toHex()` — see `stdlib/color.koda`.

With `@game`, `Color.white` aliases the same palette as `colors.white`.
