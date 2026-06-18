# Chapter 8 — Modules and imports

**You will learn:** `import()` for stdlib and local modules. `#include` is advanced.

**Time:** ~15 minutes.

---

## Import — preferred for beginners

```koda
let math = import "@math";
print(math.pi);
print(math.sqrt(25));

let io = import "@io";
let json = import "@json";
```

`@name` resolves to `stdlib/name.koda`. Local modules:

```koda
let player = import "./player";
```

| Import | Use for |
|--------|---------|
| `koda.game` | Beginner game API (`game.open`, `game.running`, …) |
| `@math` | Trig, lerp, clamp, RNG |
| `@json` | parse, stringify |
| `@io` | read, write, list, exists |
| `@array` | range, shuffle, zip |

Graphics projects use **`use raylib;`** (and optionally **`use koda.game;`**) — see [game dev guide](../guides/game-dev.md).

---

## Include — merge source files (advanced)

```koda
use raylib;
use koda.game;
```

The included file is compiled as if pasted at that line. Beginners use `import` for stdlib modules instead of `#include "stdlib/..."`.

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
