# Troubleshooting

Symptoms, causes, and fixes for common Koda issues.

---

## Installation

### `koda` is not recognized

| Cause | Fix |
|-------|-----|
| Not on PATH | Add SDK folder to PATH or run with full path |
| Wrong download | Use SDK zip, not source repo alone |

### `stdlib/math.koda` missing

| Cause | Fix |
|-------|-----|
| Only copied `koda.exe` | Unpack full SDK; keep `stdlib/` beside binary |
| Running from wrong dir | `cd` to SDK or set PATH; run `koda doctor` |

---

## Compile and link errors

### Undefined variable / parse error

```bash
koda check file.koda
```

Read line:column in the message. Common issues: missing `;` after `let`, wrong keyword, unclosed `}`.

### Link error: `raylib`, `opengl32`, etc.

Graphics projects need native libraries:

```powershell
# Windows example
$env:KODA_LINKFLAGS = "-lraylib -lopengl32 -lgdi32 -lwinmm"
```

Or set `native.linkflags` in `koda.json`. Install Raylib for your OS.

### `KODA_NATIVE_SOURCES` not found

Add C glue files to `koda.json`:

```json
"native": { "sources": ["wrappers/raylib_shim/wrapper.c"] }
```

---

## Runtime

### Panic: assert failed

Your `assert(condition, message)` failed. Run under `koda run` and read the message. Use `print` before the assert to inspect values.

### Panic: index out of bounds

Array index ≥ `len(arr)`. Check loop bounds and empty arrays.

### `readfile` / `writefile` — `.ok` is false

| Cause | Fix |
|-------|-----|
| Wrong path | Paths relative to cwd — run from project root |
| Permission | Check file locks and permissions |
| Missing file | Use `fileexists` or `try_parse` first |

### `json.try_parse` always fails

Check string content — invalid JSON or wrong escapes in source. Use `"{\"x\":1}"` in `.koda` (escaped quotes).

### `io.exists` returns false but file is there

Same cwd issue — `koda run` cwd should be project root. Compare with `fileexists` on same path string.

### Type error in unbox / expected number

Often mixing strings and numbers in arithmetic. Use `number(x)` or fix operand types.

---

## Performance

### Stutter every few seconds

GC collection — call `gcframestep()` each frame; avoid allocating large arrays every frame.

### Slow `koda run` every time

Expected — full compile+link each run. Use `koda watch` during dev; `koda build` for release.

---

## Diagnostics checklist

```bash
koda version
koda doctor
koda paths
koda check src/main.koda
koda help
```

---

## Still stuck?

1. Reduce to a **minimal** `.koda` file (10 lines).
2. Note OS, `koda version`, exact command, full error text.
3. Search [FAQ](faq.md) and [GitHub Issues](https://github.com/CharmingBlaze/koda-compiler/issues).

---

## Related

- [FAQ](faq.md)
- [CLI reference](reference/cli.md)
- [Beginner's guide](beginners-guide.md)
