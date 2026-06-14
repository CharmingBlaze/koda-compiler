# `@json` — JSON parse and stringify

**Import:** `let json = import "@json";`

---

## Functions

| Function | Args | Returns |
|----------|------|---------|
| `parse(text)` | JSON string | Parsed value (object/array/…) |
| `stringify(value, indent?)` | Value, optional indent spaces | JSON string |
| `try_parse(text)` | JSON string | `{ ok, value, error }` |

Aliases: `tryparse` ≡ `try_parse`.

---

## Examples

```koda
let json = import "@json";

let cfg = json.parse("{\"width\":800}");
let pretty = json.stringify({ name: "Koda", ok: true }, 2);

let result = json.try_parse("{bad}");
if (result.ok) {
    print(result.value);
} else {
    warn(result.error);
}
```

---

## Notes

- `stringify` with indent `2` produces multi-line pretty JSON.
- `parse` throws/panics on invalid input at runtime; prefer `try_parse` for user files.
- Object keys in parsed JSON are strings; use `cfg.width` when keys are valid identifiers.

---

## Related

- [Files and JSON chapter](../learn/09-files-and-json.md)
- [Applications guide](../guides/applications.md)
