# Runtime and garbage collection

Koda programs link against **`libkoda_runtime.a`** — a small C runtime providing memory management, builtins, and the GC.

---

## What the runtime provides

- Tagged value representation (numbers, strings, objects, …)
- Builtin functions (`print`, `readfile`, `sin`, …)
- Incremental **mark-and-sweep** garbage collector
- FFI helpers for native wrappers

You do not ship a separate runtime DLL — it is **statically linked** into your executable.

---

## Garbage collection

Gameplay code allocates freely; the GC reclaims unreachable objects.

| Function | When to use |
|----------|-------------|
| `gcframestep()` | Once per frame in games — spreads collection work |
| `gcdisable()` / `gcenable()` | Brief critical sections (rare) |
| `gc()` / `gccollect()` | Force full collection (debugging) |
| `gcstats()` | Inspect collector state |

**Tips:**

- Reuse arrays and objects in hot loops when possible.
- Avoid building huge temporary strings every frame.
- Call `gcframestep()` in your main loop for steady frame times.

---

## Manual memory vs Koda

Unlike C, you do not `free()` Koda objects. Native C code you link via wrappers still follows C rules.

---

## Related

- [Game dev — GC](../guides/game-dev.md#gc-and-performance)
- [Builtins — GC](../reference/builtins.md)
- [handoff.md](../handoff.md) (compiler pipeline)
