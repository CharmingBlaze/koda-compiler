# `@str` — string helpers

**Import:** `let str = import "@str";`

Aliases and helpers around string **methods** (`trim`, `toupper`, etc.).

---

## Functions

| Function | Description |
|----------|-------------|
| `upper(s)` | `s.toupper()` |
| `lower(s)` | `s.tolower()` |
| `trim(s)` | `s.trim()` |
| `split(s, sep)` | `s.split(sep)` |
| `join(arr, sep)` | `arr.join(sep)` |
| `format(tmpl, obj)` | Replace `{key}` in template |
| `padStart(s, n, c?)` | Left-pad to length `n` |
| `padEnd(s, n, c?)` | Right-pad to length `n` |
| `repeat(s, n)` | Repeat string `n` times |

---

## String methods (on any string)

| Method | Example |
|--------|---------|
| `trim()` | `"  hi  ".trim()` |
| `toupper()`, `tolower()` | Case change |
| `split(sep)` | Split to array |
| `replace(a, b)` | Replace first |
| `replaceall(a, b)` | Replace all |
| `startswith(s)`, `endswith(s)` | Prefix/suffix test |
| `padStart(n, c?)`, `padEnd(n, c?)` | Pad to length (also on string values directly) |

Array method: `arr.join(sep)`.

---

## Example

```koda
let str = import "@str";

let msg = str.format("Hello, {name}!", { name: "Ada" });
let parts = str.split("one,two,three", ",");
```

---

## Related

- [Language reference — Strings](../../language.md#13-string-methods)
