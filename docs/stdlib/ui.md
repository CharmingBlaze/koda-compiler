# `koda.ui` — HUD widgets

Small 2D helpers for score counters, lives, ammo slots, and similar row-of-squares UI.

```koda
use raylib;
use koda.color;
use koda.ui;
```

## Functions

| Function | Purpose |
|----------|---------|
| `draw_pips(x, y, filled, count, size, gap, fill, empty)` | Generic pip row |
| `score_pips(x, y, score, goal)` | Gold squares for collected coins (22×22, gap 28) |
| `life_pips(x, y, lives, max)` | Red squares for lives (18×18, gap 26) |

`filled` is how many slots are lit; `count` is total slots. Ranges are half-open: `for i in 0..count` → `0 .. count-1`.

## Example

```koda
score_pips(20, 134, score, goal);
life_pips(50, 172, lives, 3);
```

Custom colors and spacing:

```koda
draw_pips(90, 98, score, 8, 18, 22, colors.red, colors.black);
```

Raw Raylib (`DrawRectangle`, …) remains available in the same file.
