# Chapter 10 — Building and shipping

**You will learn:** `build`, `bundle`, `test`, `clean`, and `koda.json`.

**Time:** ~15 minutes.

---

## Daily commands

```bash
koda run [--debug] [-- <args>]
koda watch
koda check ./...
koda lint
koda fmt .
koda test [-v] [--failfast] [-run pattern]
koda bench <file> --count 5
koda eval 'print(1)'
koda repl
```

---

## Release build

```bash
koda build -o mygame
koda build --debug -o mygame_debug
```

Output is a **native executable** linked with `libkoda_runtime.a` (embedded in the toolchain).

---

## Bundle for players

```bash
koda bundle -o dist/MyGame
```

Copies the executable, assets from `koda.json`, and optional extra files into `dist/MyGame/`.

Example manifest:

```json
{
  "name": "mygame",
  "entry": "src/main.koda",
  "bundle": {
    "assets": ["assets"],
    "extra": ["README.txt"]
  }
}
```

Details: [Distribution guide](../guides/distribution.md).

---

## Native libraries (Raylib, etc.)

```json
{
  "native": {
    "sources": ["wrappers/raylib_shim/wrapper.c"],
    "linkflags": "-lraylib -lopengl32 -lgdi32 -lwinmm"
  }
}
```

Or set environment variables (override manifest):

```bash
export KODA_NATIVE_SOURCES="wrappers/raylib_shim/wrapper.c"
export KODA_LINKFLAGS="-lraylib -lm"
```

---

## Tests

Place `*_test.koda` or `tests/*.koda`:

```koda
assert(1 + 1 == 2, "math broke");
print("PASS my_test");
```

```bash
koda test
koda test tests/io_test.koda
```

---

## You finished the learn path

| Next step | Link |
|-----------|------|
| Games | [Game development](../guides/game-dev.md) |
| Apps & tools | [Applications](../guides/applications.md) |
| Full syntax | [Language reference](../language.md) |
| Every CLI flag | [CLI reference](../reference/cli.md) |
| FAQ | [FAQ](../faq.md) |

---

## Try it yourself

1. `koda build -o practice` in your `myapp` project.
2. Run the binary directly (not via `koda run`).
3. Add an asset file and run `koda bundle -o dist/practice`.
