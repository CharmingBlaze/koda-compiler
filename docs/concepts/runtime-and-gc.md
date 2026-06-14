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
| `gcFrameStep(ms)` | Once per frame in games — spreads collection work (try `0.5`–`1.0` ms) |
| `gcDisable()` / `gcEnable()` | Brief critical sections (rare) |
| `gc()` / `gcCollect()` | Force full collection (debugging) |
| `gcStats()` | Inspect collector state |
| `arena(bytes)` | Per-frame bump allocator for short-lived objects |
| `arenaReset(arena)` | O(1) reset at end of frame (do not keep arena object refs after reset) |
| `arenaAllocArray(arena, cap)` | Array allocated inside an arena |
| `arenaAllocStruct(arena, fields)` | Struct table allocated inside an arena |

**Tips:**

- Reuse arrays and objects in hot loops when possible.
- Avoid building huge temporary strings every frame.
- Call **`gcFrameStep()`** in your main loop for steady frame times.
- Set **`KODA_STACK_DEPTH`** (256–1048576) if deep recursion hits the shadow-stack cap.

---

## Manual memory vs Koda

Unlike C, you do not `free()` Koda objects. Native C code you link via wrappers still follows C rules.

---

## Related

- [Game dev — GC](../guides/game-dev.md#gc-and-performance)
- [Builtins — GC](../reference/builtins.md)
- [handoff.md](../handoff.md) (compiler pipeline)
