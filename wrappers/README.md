# Koda wrapper packages

Pre-built and generated **C library bindings** for Koda. Each subdirectory is a self-contained package.

## Using a wrapper

1. Set **`KODA_WRAPPERS`** to this folder (or keep it next to `koda` in the SDK).
2. `#include "@library"` or `import "@library"` in your `.koda` file.
3. Set **`KODA_NATIVE_SOURCES`** to `library/wrapper.c`.
4. Set **`KODA_LINKFLAGS`** per the library `README.md` or `koda.json`.

## Generate a new wrapper

```bash
koda wrap -name mylib -headers /path/to/mylib.h \
  -I /path/to/include -L /path/to/lib -l mylib \
  -o wrappers/mylib
```

See [docs/guides/wrapping-libraries.md](../docs/guides/wrapping-libraries.md).

## Packages in this tree

| Directory | Library |
|-----------|---------|
| `raylib/` | **Default** — Raylib 5.x full bindings (548 functions) + HTML docs |
| `raylib_shim/` | **Legacy only** — ~33-function subset for old `--shim` projects |
| `raylib_min/` | Reduced Raylib bridge |

Each package includes `docs/index.html` for offline browsing.
