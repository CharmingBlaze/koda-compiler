# Glossary

Terms used in Koda documentation.

| Term | Definition |
|------|------------|
| **Koda** | The programming language |
| **`koda`** | CLI tool that compiles and runs `.koda` files |
| **`kodawrap`** | Tool to generate Koda bindings from C/C++ headers |
| **`.koda`** | Source file extension |
| **Native binary** | Executable produced by `koda build` — no VM required at runtime |
| **SDK** | Release zip containing `koda`, `kodawrap`, `stdlib/`, `docs/` |
| **`stdlib/`** | Standard library modules (`math.koda`, `json.koda`, …) |
| **Builtin** | Global function implemented in the C runtime (`print`, `len`, …) |
| **Import** | `import "@math"` — load module exports |
| **Include** | `#include "file.koda"` — merge source at compile time |
| **`@module`** | Stdlib import path resolving to `stdlib/module.koda` |
| **`koda.json`** | Project manifest (entry, bundle, native link settings) |
| **Argv native** | Runtime function called with argument array (most builtins) |
| **Struct** | Fixed-field record type (`struct Player { x, y }`) |
| **Enum** | Named constant group (`enum State { Idle, Run }`) |
| **Closure** | Function capturing outer variables |
| **GC** | Garbage collector in `libkoda_runtime` |
| **`deltatime()`** | Seconds since previous frame (games) |
| **`programtime()`** | Seconds since program start |
| **Shim** | Thin C wrapper layer (e.g. Raylib shim) |
| **`KODA_LINKFLAGS`** | Environment variable for linker flags |
| **`KODA_NATIVE_SOURCES`** | C/C++ files compiled with your project |
| **Bundle** | Output folder from `koda bundle` (exe + assets) |
| **Result object** | `{ ok, value, error }` from I/O and parse helpers |
| **Raylib** | C graphics library; used via wrappers for windowed games |

---

## Related

- [Beginner's guide](beginners-guide.md)
- [Documentation style guide](STYLE-GUIDE.md)
