# `@array` — array utilities

**Import:** `let array = import "@array";`  
**Include:** `#include "stdlib/array.koda"`

---

## Re-exports from runtime

`push`, `pop`, `shift`, `unshift`, `concat`, `reverse`, `includes`, `indexof`, `sort`, `slice`, `flatten`

---

## Helper functions

| Function | Description |
|----------|-------------|
| `range(from, to)` | Integers `[from, to)` |
| `fill(n, val)` | `n` copies of `val` |
| `sum(arr)` | Sum of numeric elements |
| `zip(a, b)` | `[[a0,b0], [a1,b1], …]` |
| `shuffle(arr)` | Fisher–Yates shuffle (in place) |
| `sample(arr, count)` | `count` random elements without replacement |

---

## Example

```koda
let array = import "@array";

let levels = array.range(1, 11);
array.shuffle(levels);
let hand = array.sample(["A", "K", "Q", "J"], 2);
```

---

## Related

- [Language reference — Arrays](../../language.md#10-arrays)
