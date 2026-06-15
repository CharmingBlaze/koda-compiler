# Koda build, wrapper, and distribution guide

This guide explains the complete workflow for writing Koda programs, generating wrappers for C/C++ libraries, building native executables, and packaging applications or games for distribution.

## 1. Run a Koda program during development

`koda run` compiles with the **LLVM native** pipeline and runs the resulting binary (same as `koda build`, temp output):

```powershell
.\koda.exe run .\game.koda
```

Useful checks:

```powershell
.\koda.exe check .\game.koda
.\koda.exe disasm .\game.koda   # prints LLVM IR text
```

## 2. Build a native executable

`koda build` emits LLVM IR, lowers it with **llc** to an object file, and links that object with **`runtime/libkoda_runtime.a`** (see `runtime/src/`) using **Clang**, plus any **`KODA_NATIVE_SOURCES`** / **`KODA_LINKFLAGS`** you set for third-party libraries.

```powershell
.\koda.exe build .\game.koda -o .\game.exe
```

Optional environment variables:

| Variable | Purpose |
|----------|---------|
| `KODA_CLANG` | Full path to the Clang executable. |
| `CC` | Fallback compiler name if `KODA_CLANG` is not set. |
| `KODA_USE_LLD` | Set to `1` to request LLVM LLD linking. |
| `KODA_PATH` | Extra Koda source search paths. |
| `KODA_WRAPPERS` | Extra generated wrapper search paths. |
| `KODA_NATIVE_SOURCES` | C/C++ wrapper glue files to compile into the executable. |
| `KODA_LINKFLAGS` | Native include/library/link flags passed to Clang. |

## 3. Generate wrappers for C/C++ libraries

Build **kodawrap** from this repo (same module as `koda`; sources live under `cmd/wrapgen`):

```powershell
go build -o kodawrap.exe ./cmd/wrapgen
```

Run it on one or more headers (or use **`koda wrap …`**, which discovers **`kodawrap.exe`** next to **`koda.exe`**):

```powershell
.\kodawrap.exe -name mylib -headers .\native\mylib.h -out .\wrappers\mylib
```

Generated output:

| File | Purpose |
|------|---------|
| `mylib.koda` | Readable Koda source with `// koda:extern` lines linking to C. |
| `wrapper.c` | C ABI glue that converts Koda values to native C calls. |
| `README.md` | Human-readable usage guide. |
| `api_reference.md` | Function/type reference. |
| `examples.md` | Example usage (when docs are enabled). |
| `Makefile` / `CMakeLists.txt` | Optional (`-build`); off by default. |

Include the generated Koda wrapper from your program:

```koda
#include "wrappers/mylib/mylib.koda"

let result = my_function(1, 2);
print(result);
```

## 4. Build with a wrapper and native library

Set the generated C glue and native link flags before building:

```powershell
$env:KODA_NATIVE_SOURCES = '..\wrappers\mylib\wrapper.c ..\native\mylib.c'
$env:KODA_LINKFLAGS = '-I..\native -L..\native\build -lmylib'
.\koda.exe build .\app.koda -o .\app.exe
```

For Raylib on Windows with the current local source tree:

```powershell
$env:KODA_NATIVE_SOURCES = '..\wrappers\raylib_shim\wrapper.c'
$env:KODA_LINKFLAGS = '-I..\temp_raylib\src -L..\temp_raylib\src -lraylib -lopengl32 -lgdi32 -lwinmm'
.\koda.exe build .\raylib_brick_breaker.koda -o .\raylib_brick_breaker.exe
```

## 5. Bundle an application or game for distribution

Use `koda bundle` to create a clean folder that contains the executable, launcher, and metadata:

```powershell
.\koda.exe bundle .\game.koda -o .\dist\game
```

Bundle extra files such as assets, DLLs, licenses, or config files:

```powershell
$env:KODA_BUNDLE_FILES = '.\raylib.dll .\LICENSE.txt .\assets\logo.png'
.\koda.exe bundle .\game.koda -o .\dist\game
```

The output folder contains:

| File | Purpose |
|------|---------|
| `game.exe` | The compiled application. |
| `run.bat` | Windows launcher. |
| `README.md` | User-facing run instructions. |
| `bundle-info.txt` | Build metadata and native link settings. |
| Extra files | Any files listed in `KODA_BUNDLE_FILES`. |

## 6. Ship a wrapper package

A clean wrapper package should include:

```text
mylib-wrapper/
  mylib.koda
  wrapper.c
  README.md
  api_reference.md
  examples.md
  native/
    mylib.dll or mylib.a if redistribution is allowed
    LICENSE.txt
```

Users can either include `mylib.koda` directly or set:

```powershell
$env:KODA_WRAPPERS = 'C:\path\to\mylib-wrapper'
```

## 7. Official Releases: Windows, Linux, macOS SDK zips

GitHub **Releases** (tags `v*`) attach **offline SDK archives** built by CI for **every mainstream OS**:

| OS | Artifact |
|----|----------|
| **Windows** x64 | `koda-<tag>-sdk-windows-amd64.zip` |
| **Linux** x64 | `koda-<tag>-sdk-linux-amd64.zip` |
| **Linux** ARM64 | `koda-<tag>-sdk-linux-arm64.zip` |
| **macOS** Intel | `koda-<tag>-sdk-darwin-amd64.zip` |
| **macOS** Apple Silicon | `koda-<tag>-sdk-darwin-arm64.zip` |

Each zip unpacks to a single folder containing **`koda`** (or **`koda.exe`**), **`kodawrap`**, **`stdlib/`**, the full **`docs/`** tree, **every repo-root `*.md`**, **`wrappers/`** (including the **full raylib** binding + `wrapper.c` + reference docs), **`third_party/raylib_static/stage/`** with **raylib 5.0** headers and libraries (and **`raylib.dll`** on Windows; Linux ARM64 may ship a **README** instead of prebuilt libs — see that file), **`examples/`**, and **`SDK_README.txt`**. Keep **`stdlib/`** next to the compiler so `@` imports resolve with no extra download. Raylib workflow: [guides/raylib.md](guides/raylib.md).

Maintainers reproduce the archives by pushing a version tag; the workflow is **[`.github/workflows/release.yml`](.github/workflows/release.yml)** and **`scripts/package-release-sdk.sh`**.

### Manual “toolchain folder” layout (maintainers)

If you build from source locally, mirror the same layout next to **`koda`**:

```powershell
go build -o koda.exe .\cmd\koda
go build -o kodawrap.exe .\cmd\wrapgen
```

```text
install-root/
  koda.exe           # or `koda` on Linux/macOS
  kodawrap.exe
  stdlib/
  docs/
  *.md               # all repo-root markdown
  wrappers/
  examples/
```

## 8. Professional release checklist

Before shipping an app or game:

- Build with `koda build` or `koda bundle`.
- Run the executable from the output folder, not from the source tree.
- Include native DLLs or data files required by the app.
- Include licenses for any third-party native libraries.
- Keep wrapper docs with the wrapper package.
- Avoid source-only temp folders in the final app bundle.
- Verify on a clean machine or clean terminal session.
