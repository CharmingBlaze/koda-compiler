# `@vec3` — 3D vectors

**Import:** `let v3 = import "@vec3";`

Vectors are plain objects `{ x, y, z }`.

---

## Functions

| Function | Description |
|----------|-------------|
| `create(x, y, z)` | Construct |
| `add(a, b)` | Component sum |
| `sub(a, b)` | Component difference |
| `scale(v, s)` | Multiply by scalar |
| `dot(a, b)` | Dot product |
| `cross(a, b)` | Cross product |
| `lengthsq(v)` | Squared length |
| `length(v)` | Length |
| `normalize(v)` | Unit vector |

---

## Example

```koda
let v3 = import "@vec3";

let forward = v3.create(0, 0, 1);
let right = v3.cross(forward, v3.create(0, 1, 0));
```

---

## Related

- [vec2](vec2.md)
- [Raylib guide](../guides/raylib.md)
