# Chapter 8 — Modules and imports

**You will learn:** `#include`, `import()`, `@` stdlib modules, and dot notation.

**Time:** ~15 minutes.

---

## Include — merge source files

```koda
#include "stdlib/timer.koda"
#include "src/utils.koda"
```

The included file is compiled as if pasted at that line. Use for small helpers and stdlib `.koda` files.

---

## Import — module exports

```koda
let math = import "@math";
print(math.pi);
print(math.sqrt(25));
```

`@name` resolves to `stdlib/name.koda` or a built-in export object for `@math`, `@json`, `@io`.

---

## Dot notation

After import, call through the binding:

```koda
let io = import "@io";
let json = import "@json";

if (io.exists("save.dat")) {
    let data = json.parse(io.read("save.dat").value);
}
```

Works the same as global builtins where names overlap: `fileexists` ≡ `io.exists` when routed correctly.

---

## Stdlib quick map

| Import | Use for |
|--------|---------|
| `@math` | Trig, lerp, clamp, RNG |
| `@json` | parse, stringify, try_parse |
| `@io` | read, write, list, exists |
| `@timer` | Cooldowns, intervals, countdowns |
| `@array` | range, shuffle, zip |
| `@vec2`, `@vec3` | Vector math |
| `@util` | clamp01, pick_weighted |
| `@noise` | 1D value noise |
| `@str` | String helper aliases |

Full index: [Stdlib overview](../stdlib/README.md).

---

## Project modules

```koda
let utils = import "src/utils.koda";
```

Relative paths load from the project tree (see loader rules in [Language reference](../language.md)).

---

## Try it yourself

1. `import "@math"` and print `sin(math.pi / 2)`.
2. `import "@json"` and round-trip an object with `stringify` then `parse`.

---

## Next chapter

[Chapter 9 — Files and JSON](09-files-and-json.md)
