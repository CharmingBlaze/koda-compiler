# Wrapping any C library for Koda

Use **kodawrap** (`koda wrap …`) to turn C/C++ headers into an organized Koda package: bindings, C glue, and documentation.

For **graphics (Raylib)** in a game project, prefer the beginner path first:

```bash
koda new mygame --template graphics
# or in an existing project:
koda setup raylib
koda run
```

`koda setup raylib` writes or **refreshes** `wrappers/raylib_shim/` (overwrites stale copies) and sets `koda.json`:

```json
{
  "native": {
    "sources": ["wrappers/raylib_shim/wrapper.c"],
    "graphics": true
  }
}
```

When `"graphics": true` and `linkflags` is empty, Koda applies platform Raylib link flags automatically. You do **not** need `KODA_LINKFLAGS` for a standard graphics template.

---

## Workflow

```bash
# 1. Generate package (clang parses headers when available)
koda wrap -name sqlite3 -headers /usr/include/sqlite3.h \
  -I /usr/include -L /usr/lib -l sqlite3 \
  -o wrappers/sqlite3

# 2. Point Koda at wrappers
export KODA_WRAPPERS="$PWD/wrappers"

# 3. Use in your game or app
#include "@sqlite3"
# or: let sqlite = import "@sqlite3";

# 4. Build — merge generated koda.json or set native fields
```

**Beginner (recommended):** merge the generated `koda.json` fragment into your project manifest:

```json
{
  "native": {
    "sources": ["wrappers/sqlite3/wrapper.c"],
    "linkflags": "-I/usr/include -L/usr/lib -lsqlite3"
  }
}
```

**Advanced:** environment variables still work when you need ad-hoc builds:

```bash
export KODA_NATIVE_SOURCES="wrappers/sqlite3/wrapper.c"
export KODA_LINKFLAGS="-I/usr/include -L/usr/lib -lsqlite3"
koda build src/main.koda -o app
```

---

## Generated package layout

Every successful `koda wrap` run produces:

| File | Purpose |
|------|---------|
| `<name>.koda` | Bindings grouped by category with comments |
| `wrapper.c` | C glue — add to `native.sources` in `koda.json` |
| `README.md` | Quick start, link flags, troubleshooting |
| `api_reference.md` | Full API in Markdown |
| `examples.md` | Copy-paste sample calls |
| `koda.json` | Fragment to merge into your project manifest |
| `META.json` | Stats, import path, link flags (for tooling) |
| `docs/index.html` | Searchable offline documentation |

Open `docs/index.html` in a browser for the best reading experience.

List installed wrappers from the CLI:

```bash
koda doc wrappers
koda doc wrapper @raylib
```

---

## CLI reference

```bash
koda wrap -name <lib> -headers <a.h>[,b.h] -out <dir>   # legacy

koda wrap [options] <header.h> [more.h ...]              # modern

Options:
  -name <lib>     library name
  -o <dir>        output directory
  -I <dir>        include path for clang (repeatable)
  -L <libdir>     linker search path (stored in koda.json)
  -l <lib>        link library (stored in koda.json)
  --linkflags     raw linker flags string
  --no-clang      regex-only parsing (fallback)
  --no-docs       skip Markdown docs
  --no-html       skip docs/index.html
  -v              verbose
```

Run `koda wrap --help` for the full list.

---

## Organizing many libraries

Ship a **`wrappers/`** tree next to `koda` and set **`KODA_WRAPPERS`** once:

```text
wrappers/
  README.md
  raylib/
    raylib.koda  wrapper.c  README.md  docs/
  sqlite3/
    sqlite3.koda  wrapper.c  ...
```

Users import with:

```koda
#include "@raylib"
let db = import "@sqlite3";
```

The loader searches `KODA_WRAPPERS` for `@name` → `name/name.koda`.

---

## Tips for complex headers

| Issue | Approach |
|-------|----------|
| Missing types | Add `-I` paths for dependencies |
| Clang fails | Use `--no-clang` for regex mode |
| Huge APIs | Wrap one public header; exclude internals |
| C++ only | Wrap a thin C API header if available |
| Windows DLLs | Copy `.dll` next to exe or bundle with `koda bundle` |

---

## Upgrading a wrapper after the native library changes

Every `koda wrap` run writes **`META.json`** with header paths, include/link flags, and content hashes. Regenerate without remembering the original CLI:

```bash
# After updating raylib, sqlite, etc.
koda wrap upgrade wrappers/raylib
koda wrap upgrade @raylib          # resolves via KODA_WRAPPERS / SDK wrappers/

# Check whether headers changed since last wrap
koda wrap check wrappers/raylib
```

`koda doctor` warns when your project's `native.sources` point at a stale wrapper.

**Workflow:**

1. Update the native library (rebuild `third_party/`, vcpkg upgrade, etc.)
2. `koda wrap upgrade wrappers/mylib`
3. `koda build` and fix any renamed/removed API calls in your Koda code

Record the upstream version when generating:

```bash
koda wrap -name raylib -I ./include -o wrappers/raylib \
  --library-version 5.5 ./include/raylib.h
```

Do **not** hand-edit generated `wrapper.c` or `*.koda` — changes are lost on upgrade.

### Install from catalog

```bash
koda wrap list
koda wrap install raylib --project    # copies SDK wrapper + updates koda.json
koda wrap install sqlite3@3          # generates from system headers when found
```

### C++ headers

Use `--cpp` or pass a `.hpp` file — wrapgen runs clang in C++17 mode and picks up free functions and class methods where visible in the AST. Full class APIs still work best through a thin C shim.

---

## See also

- [wrappers.md](../wrappers.md) — resolver and link flags
- [Raylib guide](raylib.md) — graphics workflow
- [Distribution](distribution.md) — shipping binaries with native libs
