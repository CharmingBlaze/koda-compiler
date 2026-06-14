# test

Koda bindings for the **test** library.

---

## Files in this folder

| File | Description |
|------|-------------|
| `test.koda` | Include this in your Koda program. |
| `wrapper.c` | Compiled automatically by `koda build` and `koda bundle`. You do not need to touch this file. |
| `api_reference.md` | Full reference for every function, struct, and constant. |
| `examples.md` | Ready-to-run code examples. |

---

## Library summary

- **3** functions

---

## Quick start

**Step 1.** Include the bindings at the top of your Koda program:

```koda
#include "test.koda"
```

**Step 2.** Call functions directly by name:

```koda
let result = add(a, b);
print(result);
```

**Step 3.** Build or bundle:

```powershell
set KODA_NATIVE_SOURCES=test\wrapper.c
set KODA_LINKFLAGS=-I<include-dir> -L<lib-dir> -ltest

koda build  mygame.koda -o mygame.exe
koda bundle mygame.koda -o dist\mygame
```

---

## Troubleshooting

**Undefined symbol**  
Make sure `KODA_NATIVE_SOURCES` points to `wrapper.c`.

**Missing header or library**  
Add `-I<dir>` for headers and `-L<dir> -ltest` for the library in `KODA_LINKFLAGS`.

**Unexpected return values**  
Check the type conversions in `wrapper.c`. Pointer and struct types may need manual adjustment for complex cases.

---

## See also

- [API Reference](api_reference.md)
- [Examples](examples.md)
