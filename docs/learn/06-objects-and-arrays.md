# Chapter 6 — Objects and arrays

**You will learn:** arrays for lists, and object literals for JSON-like data.

**Time:** ~15 minutes.

> **Modeling a player or enemy?** Read [chapter 7 — Structs](07-structs-and-enums.md) first. Use structs for game data, not `{ x: 10, y: 20 }`.

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

## Objects — for config and JSON

Object literals are flexible maps. Use them when the shape comes from a file or API:

```koda
let config = {
    volume: 80,
    fullscreen: true
};

print(config.volume);
```

Shorthand when keys match variable names:

```koda
let volume = 80;
let cfg = { volume };
```

`delete` and advanced object patterns are in the [language reference](../../language.md) — not needed on day one.

**Methods** — functions inside objects use `this`:

```koda
let cam = {
    yaw: 0.0,
    update: func() { this.yaw = this.yaw + 0.01; }
};
cam.update();
```

For game entities with fixed fields, prefer [structs in chapter 7](07-structs-and-enums.md).

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
```

| Function | Purpose |
|----------|---------|
| `range(from, to)` | `[from, …, to-1]` |
| `fill(n, val)` | Array of `n` copies of `val` |
| `sum(arr)` | Numeric sum |
| `shuffle(arr)` | Random reorder (in place) |

Full API: [stdlib/array](../stdlib/array.md).

---

## Try it yourself

Keep an inventory **array**, load settings from a JSON **object** with `import "@json"`, and store game entities as **structs** (chapter 7).

---

## Next chapter

[Chapter 7 — Structs and enums](07-structs-and-enums.md)
