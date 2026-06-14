# `@vec2` — 2D vectors

**Import:** `let v2 = import "@vec2";`

Vectors are plain objects `{ x, y }`.

---

## Functions

| Function | Description |
|----------|-------------|
| `vec2(x, y)` | Construct |
| `add(a, b)` | Component sum |
| `sub(a, b)` | Component difference |
| `scale(a, s)` | Multiply by scalar |
| `dot(a, b)` | Dot product |
| `lengthsq(a)` | Squared length |
| `length(a)` | Length |
| `normalize(a)` | Unit vector |

---

## Example

```koda
let v2 = import "@vec2";

let pos = v2.vec2(10, 20);
let vel = v2.vec2(100, 0);
pos = v2.add(pos, v2.scale(vel, deltatime()));
```

---

## Related

- [math](math.md) — `lerp`, `distance`, `normalize` builtins
- [vec3](vec3.md)
