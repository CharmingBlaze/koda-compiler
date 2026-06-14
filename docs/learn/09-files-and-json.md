# Chapter 9 — Files and JSON

**You will learn:** reading and writing files, result objects, and JSON config.

**Time:** ~15 minutes.

---

## Result objects

Many I/O builtins return `{ ok, value, error }`:

```koda
let result = readfile("config.json");
if (result.ok) {
    print(result.value);
} else {
    print(result.error);
}
```

`writefile` and `readfile` use this pattern. Booleans like `fileexists` return `true`/`false` directly.

---

## Global file builtins

| Builtin | Returns |
|---------|---------|
| `readfile(path)` | ok/value/error |
| `writefile(path, text)` | ok/value/error |
| `appendfile(path, text)` | bool |
| `fileexists(path)` | bool |
| `deletefile(path)` | bool |
| `isfile(path)` | bool |
| `isdir(path)` | bool |
| `filesize(path)` | number |
| `listdir(path)` | array of names |

---

## `@io` module

```koda
let io = import "@io";

io.write("save.txt", "level=3");
let text = io.read("save.txt");
if (text.ok) {
    print(text.value);
}

let files = io.list("assets");
if (io.isfile("save.txt")) {
    io.remove("save.txt");
}
```

API mirrors globals: `read`, `write`, `append`, `exists`, `remove`, `isfile`, `isdir`, `size`, `list`.

---

## JSON

```koda
let json = import "@json";

let cfg = { width: 800, height: 600, title: "My Game" };
let text = json.stringify(cfg, 2);   // optional indent
writefile("settings.json", text);

let loaded = json.parse(readfile("settings.json").value);
print(loaded.width);

let safe = json.try_parse("{not json}");
if (!safe.ok) {
    warn(safe.error);
}
```

| API | Behavior |
|-----|----------|
| `json.parse(s)` | Value or error at runtime |
| `json.stringify(v, indent?)` | JSON string |
| `json.try_parse(s)` | Always `{ ok, value, error }` |

---

## Try it yourself

Save `{ "score": 0, "name": "Player" }` to `save.json`, reload it, increment score, save again.

---

## Next chapter

[Chapter 10 — Building and shipping](10-building-and-shipping.md)
