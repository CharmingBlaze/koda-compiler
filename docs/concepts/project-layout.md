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

Often add:

```text
wrappers/
  raylib_shim/
    raylib.koda
    wrapper.c
```

Set link flags per OS in `koda.json` or `KODA_LINKFLAGS`.

---

## Tests

```text
tests/
  io_test.koda
  math_test.koda
```

Run with `koda test` from project root.

---

## Related

- [CLI reference](../reference/cli.md)
- [Distribution](../guides/distribution.md)
