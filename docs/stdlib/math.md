# `@math` — mathematics and RNG

**Import:** `let math = import "@math";`  
**Include:** `#include "stdlib/math.koda"`

---

## Constants

| Name | Value |
|------|-------|
| `pi` | π |
| `e` | Euler's number |
| `tau` | 2π |
| `phi` | Golden ratio |
| `inf`, `nan` | IEEE special values |

---

## Trigonometry

`sin`, `cos`, `tan`, `asin`, `acos`, `atan`, `atan2`, `hypot`

Radians unless you use `degrees` / `radians` (aliases `deg`, `rad`).

---

## Powers and logs

| Function | Notes |
|----------|-------|
| `sqrt`, `cbrt` | Roots |
| `pow`, `exp` | Power / e^x |
| `log`, `log2`, `log10` | Natural, base-2, base-10 |

---

## Rounding and sign

`floor`, `ceil`, `round`, `trunc`, `abs`, `sign`, `min`, `max`

---

## Game math

| Function | Purpose |
|----------|---------|
| `lerp(a, b, t)` | Linear interpolation |
| `clamp(v, lo, hi)` | Bound value |
| `wrap(v, lo, hi)` | Repeat value in range |
| `approach(cur, target, maxDelta)` | Move toward target |
| `smoothstep` | Smooth 0–1 curve |
| `distance`, `distancesq` | 2D distance |
| `normalize` | Unit vector from x,y |
| `map(v, inLo, inHi, outLo, outHi)` | Range remap |

---

## Random numbers

Uses **xoshiro128\*\*** seeded from OS entropy (not libc `rand()`).

| Function | Returns |
|----------|---------|
| `random()` | Float in [0, 1) |
| `randomrange(lo, hi)` | Float in [lo, hi) |
| `randomint(lo, hi)` | Integer in [lo, hi) |
| `randomchoice(arr)` | One random element |
| `randomseed(n)` | Seeded sequence (reproducible runs) |

```koda
let math = import "@math";
math.randomseed(42);
let roll = math.randomint(1, 7);
```

---

## Examples

```koda
let math = import "@math";

let angle = math.atan2(dy, dx);
let t = math.clamp(progress, 0.0, 1.0);
let x = math.lerp(startX, endX, t);
```

See also: [vec2](vec2.md), [Game dev guide](../guides/game-dev.md).
