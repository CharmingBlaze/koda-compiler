# Koda implementation status

This file tracks the language surface against the Go frontend, LLVM codegen, C runtime, and wrapper/FFI path (**single native pipeline** — no bytecode VM).

## Supported and exercised

- **Variables:** `let name = expr;`, `let name;`, mutable assignment, global and local scope.
- **Primitives:** number, string, bool, null.
- **Composite values:** arrays and objects in the C / LLVM runtime.
- **Operators:** arithmetic, comparison, equality, logical operators, compound assignment, prefix/postfix update.
- **Control flow:** `if`, `else`, `while`, `for`, `for ... of`, `switch`, `break`, `continue`.
- **Functions:** declarations, calls, recursion, return, function expressions.
- **Closures:** LLVM codegen path with upvalue / capture work (see `internal/codegen`, `internal/sema`).
- **Modules:** `#include` and import resolution through local path, `KODA_PATH`, and `KODA_WRAPPERS`.
- **Standard library:** `print`, `type`, `number`, `string`, `len`, `time`, `sleep`, `abs`, `sqrt`, `random`.
- **Native apps/games:** `koda build` emits LLVM IR, lowers it with **llc**, then links with **`runtime/libkoda_runtime.a`** and optional native wrapper sources through `KODA_NATIVE_SOURCES` and `KODA_LINKFLAGS`.
- **Distributable apps/games:** `koda bundle` builds an executable and writes a clean distribution folder with a launcher, README, and build metadata.
- **Library wrappers:** line directives of the form `// koda:extern` compile Koda calls to C ABI wrapper functions `Value symbol(int argCount, Value* args)` (NaN-boxed `uint64_t` / `i64` in IR).

## Important implementation notes

- **`koda run`** uses the **same** LLVM native pipeline as **`koda build`** (temp executable).
- Raylib is not built into Koda. It is used through a normal wrapper bridge under `wrappers/raylib_min`.
- `raylib_brick_breaker.koda` proves Koda can compile and run a real graphical application.

## Known gaps / work to harden

- **Object method shorthand and `this`:** supported in the LLVM path; verify with tests you care about.
- **Object for-of key/value form:** supported via `for (let k, v of obj)` where codegen lowers it.
- **Compound property/index assignment:** `obj.x += 1` and `arr[i] -= 1` — verify against current codegen.
- **Diagnostics:** declare-before-use and invalid assignment errors should be consistently reported before code generation.
- **Documentation:** keep `KODA_PROGRAMMER_REFERENCE.md` user-facing and this file engineering-facing.
- **Distribution guide:** see `DISTRIBUTION_GUIDE.md` for compile, run, wrapper, library-linking, and app bundle workflows.

## Required release gate

Before claiming a feature is complete, both of these should pass when applicable:

```powershell
.\bin\koda.exe run .\tests\native_conformance.koda
.\bin\koda.exe build .\tests\native_conformance.koda -o .\tests\native_conformance.exe
.\tests\native_conformance.exe
```

Graphical/game wrapper release gate:

```powershell
$env:KODA_NATIVE_SOURCES = '..\wrappers\raylib_min\raylib_bridge.c'
$env:KODA_LINKFLAGS = '-I..\temp_raylib\src -L..\temp_raylib\src -lraylib -lopengl32 -lgdi32 -lwinmm'
.\koda.exe build .\raylib_brick_breaker.koda -o .\raylib_brick_breaker.exe
.\raylib_brick_breaker.exe
```
