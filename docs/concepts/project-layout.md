# Project layout

Standard structure for a Koda application or game.

---

## `koda new` template

```text
myapp/
  koda.json          # manifest
  src/
    main.koda        # entry point
  assets/            # images, sounds, data
  README.md
```

---

## `koda.json` essentials

```json
{
  "name": "myapp",
  "version": "0.1.0",
  "entry": "src/main.koda",
  "bundle": {
    "assets": ["assets"],
    "extra": ["README.md"]
  },
  "native": {
    "sources": [],
    "linkflags": ""
  }
}
```

| Field | Purpose |
|-------|---------|
| `entry` | Main source file — required |
| `bundle.assets` | Copied into `koda bundle` output |
| `native.sources` | C/C++ glue compiled with every build |
| `native.linkflags` | Passed to the linker (`-lraylib`, …) |

---

## Graphics projects

```text
wrappers/
  raylib_shim/
    raylib.koda    # koda:extern declarations
    wrapper.c      # C glue
```

```json
"native": {
  "sources": ["wrappers/raylib_shim/wrapper.c"],
  "graphics": true
}
```

Set `"graphics": true` — Koda adds platform Raylib link flags. No manual `KODA_LINKFLAGS` for standard SDK installs.

**Setup / refresh:**

```bash
koda setup raylib          # beginner shim + @game
koda setup raylib --full   # 548-function wrapper + @raylib
```

Source pattern:

```koda
#include "wrappers/raylib_shim/raylib.koda"
#include "@game"
```

---

## Tests

```text
tests/
  io_test.koda
  math_test.koda
```

Run with `koda test` from project root. Convention: **`*_test.koda`** files (Go-style), e.g. `tests/io_test.koda` or `src/player_test.koda`.

---

## Related

- [CLI reference](../reference/cli.md)
- [Distribution](../guides/distribution.md)
