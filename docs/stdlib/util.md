# `@util` — small helpers

**Import:** `let util = import "@util";`  
**Include:** `#include "stdlib/util.koda"`

---

## Functions

| Function | Description |
|----------|-------------|
| `clamp01(x)` | Clamp to [0, 1] |
| `sign(x)` | -1, 0, or 1 |
| `pingpong(t, length)` | Oscillate 0..length |
| `repeat_n(n, val)` | Array of `n` copies |
| `shallow_copy(obj)` | Copy object keys |
| `pick_weighted(weights)` | Random index by weight array |

---

## Example

```koda
let util = import "@util";

let t = util.pingpong(programtime(), 1.0);
let enemy = util.pick_weighted([10, 30, 5]);  // index 0, 1, or 2
```

---

## Related

- [math](math.md)
