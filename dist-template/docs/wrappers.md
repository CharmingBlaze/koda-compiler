# Distributing and using C/C++ library wrappers (Raylib, etc.)

Koda stays on a **single native pipeline**: the Go frontend lowers programs to **LLVM IR**, then **`koda build`** runs **llc** on that IR and links the resulting object with **`runtime/libkoda_runtime.a`** (C sources under `runtime/src/`) using **Clang** as the driver, plus optional **`KODA_NATIVE_SOURCES`** / **`KODA_LINKFLAGS`** for your C glue. Wrappers are **ordinary `.koda` modules** (plus your own glue) that you ship alongside—or inside—a **wrapper root directory** so `import "raylib"` (or `@raylib/core`) resolves without vendoring copies into every app.

## 1. Resolver: where imports look

For each `import "name"` / `#include "file"` string, the loader tries, in order:

1. Next to the importing file  
2. Each directory in **`KODA_PATH`** (same semantics as a normal search path)  
3. Each directory in **`KODA_WRAPPERS`** (path list, same separator as `PATH` on your OS)

Under each root it looks for `name`, `name.koda`, or `name/index.koda` (and the same for `@segment/...` style modules). Set **`KODA_WRAPPERS`** to the folder you distribute (e.g. unzip a `wrappers/` bundle next to **`koda.exe`** and point users at it once).

Example (PowerShell):

```powershell
$env:KODA_WRAPPERS = "C:\tools\koda-wrappers"
koda run .\game.koda
```

Example (Unix):

```bash
export KODA_WRAPPERS="$HOME/koda-wrappers"
koda run ./game.koda
```

## 2. Native link: actually link the C/C++ library

Generated or hand-written Koda bindings call into C symbols. **`koda build`** must pass the right flags to **clang** so the final executable links Raylib (or SQLite, etc.).

Set **`KODA_LINKFLAGS`** to extra tokens (split on whitespace) inserted **before** `-o`:

```bash
export KODA_LINKFLAGS="-L/opt/homebrew/lib -lraylib -framework CoreVideo -framework IOKit -framework Cocoa -framework GLUT -framework OpenGL"
koda build game.koda -o game
```

On Windows you might use `-L` to a MinGW lib folder and `-lraylib` (exact flags depend on how Raylib was built). Adjust per library and platform.

## 3. Producing wrappers

- **`cmd/wrapgen`** — build the binary as **`kodawrap`** (`go build -o kodawrap ./cmd/wrapgen`). Header-driven generator in the **same Go module** as `koda`. It emits readable `.koda`, `wrapper.c`, and Markdown — no Python, no extra runtime. The legacy name **`wrapgen`** is the same program.  
- **`koda wrap …`** — forwards to **`kodawrap`** (or **`wrapgen`** / **`kujiwrap`**) next to **`koda`** or on **`PATH`**.  
- Optional release tooling may embed or ship a **`wrappers/`** tree and document **`KODA_WRAPPERS`** for turnkey bundles.

Your distribution can be:

| Artifact | Role |
|----------|------|
| `koda` | Compiler + `koda build` / `koda run` |
| `wrappers/` directory (zip) | Pre-built or generated `*.koda` trees per library |
| README in `wrappers/` | Per-library **KODA_LINKFLAGS** hints |

## 4. Raylib specifically

Official SDK zips ship **`wrappers/raylib/`** (full `.koda` + `wrapper.c` + Markdown/HTML) and **raylib 5.0 prebuilds** under **`third_party/raylib_static/stage/`** (headers + `libraylib.a` + Windows `raylib.dll`). The compiler picks up that tree automatically when linking (see **`KODA_USE_VENDORED_RAYLIB`** and **`KODA_RAYLIB_STAGE`** in `koda -help`).

1. In source: `#include "raylib/raylib.koda"` (resolved against **`wrappers/`** next to **`koda`** — see [guides/raylib.md](guides/raylib.md)).  
2. Set **`KODA_NATIVE_SOURCES`** to **`wrappers/raylib/wrapper.c`**.  
3. On Windows, copy **`raylib.dll`** next to the `.exe` when using a dynamic link or bundle it with **`KODA_BUNDLE_FILES`**.  
4. If you do not use the vendored stage (e.g. Linux ARM64 without an upstream binary), set **`KODA_LINKFLAGS`** to your system `-I` / `-L` / `-lraylib` (and frameworks on macOS).

Until every symbol is lowered automatically, some declarations may still need to match whatever ABI your wrapper emits (same as any FFI layer).

## 5. Checklist for “easy for users”

- [ ] Ship **`wrappers.zip`** with one top-level folder; user sets **`KODA_WRAPPERS`** to that folder’s absolute path.  
- [ ] Document **`KODA_LINKFLAGS`** per OS for each heavy library you support.  
- [ ] Optionally ship a **small launcher script** that sets both env vars and runs `koda run` / `koda build`.

This keeps **one** language and **one** LLVM link step—wrappers are data + Koda source on the search path, not a second compiler fork.
