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

Full Raylib resolves from the **SDK** — no project-local copy required:

```json
"native": {
  "sources": ["wrappers/raylib/wrapper.c"],
  "graphics": true
}
```

Set `"graphics": true` — Koda adds platform Raylib link flags. No manual `KODA_LINKFLAGS` for standard SDK installs.

**Setup:**

```bash
koda setup raylib          # full wrapper (default)
koda setup raylib --shim   # legacy ~33-function shim only
```

Source pattern:

```koda
use raylib;
use koda.game;   // optional helpers
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
