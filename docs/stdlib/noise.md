# `@noise` — 1D value noise

**Import:** `let noise = import "@noise";`

Simple deterministic noise for terrain, wind, procedural animation. Seeded via `randomseed`.

---

## Functions

| Function | Description |
|----------|-------------|
| `seed(s)` | Set noise seed (calls `randomseed`) |
| `value1d(x)` | Smooth noise at coordinate `x` |

---

## Example

```koda
let noise = import "@noise";

noise.seed(12345);
let i = 0;
while (i < 10) {
    print(noise.value1d(i * 0.5));
    i = i + 1;
}
```

---

## Notes

- 1D only — for 2D/3D, combine multiple `value1d` calls or use a C library via wrappers.
- Deterministic when you call `seed` before sampling.

---

## Related

- [math — random](math.md)
