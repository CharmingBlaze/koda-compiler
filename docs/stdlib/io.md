# `@io` — files and directories

**Import:** `let io = import "@io";`

Wraps file builtins with a consistent namespace.

---

## API

| Method | Global equivalent | Returns |
|--------|-------------------|---------|
| `read(path)` | `readfile` | `{ ok, value, error }` |
| `write(path, text)` | `writefile` | `{ ok, value, error }` |
| `append(path, text)` | `appendfile` | `bool` |
| `exists(path)` | `fileexists` | `bool` |
| `remove(path)` | `deletefile` | `bool` |
| `isfile(path)` | `isfile` | `bool` |
| `isdir(path)` | `isdir` | `bool` |
| `size(path)` | `filesize` | number (bytes) |
| `list(path)` | `listdir` | array of names |

> **Note:** `remove` not `delete` — `delete` is a reserved keyword in Koda.

---

## Example

```koda
let io = import "@io";

if (!io.exists("save.dat")) {
    io.write("save.dat", "{\"level\":1}");
}

let data = io.read("save.dat");
if (data.ok) {
    print(data.value);
}

for (let name of io.list("assets")) {
    print(name);
}
```

---

## Paths

Use forward slashes or platform paths relative to the **process working directory** (usually project root when using `koda run`).

---

## Related

- [Applications guide](../guides/applications.md)
- [Learn — Files and JSON](../learn/09-files-and-json.md)
