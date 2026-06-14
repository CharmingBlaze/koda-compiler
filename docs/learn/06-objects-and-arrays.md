# Chapter 6 — Objects and arrays

**You will learn:** object literals, arrays, indexing, methods, and the `@array` module.

**Time:** ~15 minutes.

---

## Objects

```koda
let player = {
    name: "Hero",
    x: 100,
    y: 200
};

player.x = player.x + 10;
print(player["name"]);

let copy = { name, x, y };  // shorthand when keys match variable names
```

Delete a property:

```koda
delete player["temp"];
```

---

## Arrays

```koda
let items = ["sword", "shield"];
items.push("bow");
print(len(items));
print(items[0]);

items[1] = "armor";
```

Slice and sort:

```koda
let part = items.slice(1, 3);
items.sort();
```

---

## String methods

```koda
let s = "  hello world  ";
print(s.trim());
print(s.toupper());
print(s.split(" "));
print(s.startswith("hello"));
```

Or `import "@str"` for helper aliases.

---

## Array module

```koda
let array = import "@array";

let nums = array.range(0, 10);
let total = array.sum(nums);
array.shuffle(nums);

let pairs = array.zip(["a", "b"], [1, 2]);
```

| Function | Purpose |
|----------|---------|
| `range(from, to)` | `[from, …, to-1]` |
| `fill(n, val)` | Array of `n` copies of `val` |
| `sum(arr)` | Numeric sum |
| `zip(a, b)` | Pairs of elements |
| `shuffle(arr)` | Random reorder (in place) |
| `sample(arr, n)` | `n` random elements |

Full API: [stdlib/array](../stdlib/array.md).

---

## Try it yourself

Build an inventory array, push three items, shuffle it, and print each item in a `for` loop.

---

## Next chapter

[Chapter 7 — Structs and enums](07-structs-and-enums.md)
