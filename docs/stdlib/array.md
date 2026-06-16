# `@array` — array utilities

**Import:** `let array = import "@array";`  
**Include:** `#include "stdlib/array.koda"`

---

## Re-exports from runtime

`push`, `pop`, `shift`, `unshift`, `concat`, `reverse`, `includes`, `indexof`, `sort`, `slice`, `flatten`

---

## Instance methods (on array values)

| Method / property | Description |
|-------------------|-------------|
| `arr.add(x)` | Append an element (alias for `push`) |
| `arr.remove_at(i)` | Remove element at index `i`, return it |
| `arr.clear()` | Remove all elements |
| `arr.count` | Number of elements (property; also `arr.length()`) |
| `for x in arr { … }` | Loop over elements (preferred for beginners) |
| `arr.each(fn)` | Call `fn(element)` for each element |

Also available: `push`, `pop`, `map`, `filter`, `find`, `reduce`, `slice`, `sort`, `sort(comparator)`, `reverse`, `join`, `includes`, `indexof`.

`sort()` sorts numbers/strings in place (default order). `sort(func(a, b) { return a.y - b.y; })` accepts a comparator returning a number (negative if `a` sorts before `b`). Works with **struct arrays** (`enemies[0].health`) and object arrays.

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
