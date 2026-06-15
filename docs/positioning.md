# Koda positioning

An honest picture of what Koda is, who it is for, and what remains before broader claims hold.

> **Audience:** contributors, release planners, and technical writers. End users should start with the [Beginner's guide](beginners-guide.md).

---

## Where Koda sits

```
  C                              Koda
  Fast, manual memory       ->    Native binary - GC -
                                 approachable syntax -
                                 beginner-friendly
```

Compared to **interpreted scripting languages** (Lua, JavaScript in a VM, ...), Koda compiles to a **real native executable** -- no interpreter on the user's machine.

Koda is closest to **Lua + LLVM**: a small language that compiles to a real native binary, with automatic memory management and syntax approachable to beginners. It is **not** a drop-in replacement for C in systems programming. It **is** a legitimate answer when someone wants to **ship a native game or app** without manual memory *or* VM overhead.

---

## Is Koda a good C alternative for beginners?

**Yes -- for its target audience.**

The honest pitch: Koda is what you reach for when a beginner wants to ship a **real native game or app**, and C's manual memory is the wrong answer -- without giving up native speed to a VM.

| Audience | Fit today |
|----------|-----------|
| **"I want to make games"** | Close to ready -- LLVM backend, Raylib path, GC tuned for game loops |
| **"I want to learn systems programming instead of C"** | Not yet -- missing opt-in integers and struct methods are the main blockers |

The infrastructure is genuinely strong: the GC is more sophisticated than most hobby-language collectors, the LLVM backend is real, and Raylib integration works.

---

## Strengths already in place

These are **working today**, not roadmap aspirations.

| Area | What you get |
|------|----------------|
| **Values** | NaN-boxing -- 64-bit tagged values, low overhead |
| **GC** | Tri-generational collector: 256 KB nursery, incremental major collection, O(1) remembered set, heap lookup cache for conservative roots |
| **Per-frame allocation** | Arena builtins (`arena`, `arenaReset`, `arenaAllocArray`, `arenaAllocStruct`) |
| **Game-loop GC** | `gcFrameStep(ms)` -- spread collection work across frames |
| **Codegen** | LLVM IR → O2 when available; native binary, no VM |
| **FFI** | C wrappers via `kodawrap` and `// koda:extern`; Raylib proven in-tree |
| **Language surface** | `let`, `for-of`, closures, structs, enums, `ok`/`err`, defer |
| **Diagnostics** | Typo hints ("did you mean?") on undefined names |

See [Runtime and GC](concepts/runtime-and-gc.md) and [status.md](status.md) for engineering detail.

---

## Gaps worth addressing now

Ordered by impact for beginners and for the "C alternative" claim.

### 1. All numbers are 64-bit floats

**Severity:** highest for binary / systems-adjacent work.

You cannot represent `uint8_t` pixel data, `int32_t` network packets, or bitfield flags accurately. A game reading a PNG gets floats pretending to be bytes. Pure game logic (positions, timers, scores) is fine; anything touching **binary data, file formats, or networking** quietly breaks.

**Fix:** opt-in integer types (`int`, `i32`, `u8`, ...) lowering to real LLVM `i32` / `i64` / `i8`. They do not need to be the default -- most beginners never need them -- but they need to **exist**.

**Roadmap ref:** new surface (see [ROADMAP.md](ROADMAP.md) tier 1).

---

### 2. No methods on struct types

**Severity:** highest for beginner experience.

Structs are data-only today. You write `func area(r) { return r.w * r.h; }` outside the type. Beginners from any modern language expect `box.area()`.

**Clarification:** object-literal method shorthand and `this` on **object values** are a separate feature from **declaring methods on a named struct type**. The missing piece is syntax to attach `func` declarations to `struct` types and emit them as methods.

**Roadmap ref:** tier 1 -- see [ROADMAP.md](ROADMAP.md).

---

### 3. Unused variable / function warnings

**Severity:** high for learners.

A misspelled variable name fails silently at runtime. `--warn-unused` (MASTER_PLAN **E4/E5**) would catch a large class of beginner mistakes at compile time.

**Roadmap ref:** tier 2.

---

### 4. Enum switch exhaustiveness

**Severity:** medium -- real foot-gun for game state.

Adding a new enum case without updating every `switch` is unchecked. A **warning** (not an error) when a `switch` over an enum does not cover all cases would help beginners modelling `Idle`, `Running`, `Dead`, etc.

**Roadmap ref:** tier 2 -- MASTER_PLAN **A3** (partial).

---

### 5. ASAN / Valgrind in CI

**Severity:** medium -- safety net for runtime complexity.

The C runtime is now meaningfully complex: nursery, remembered set, arena, inline arrays, heap range cache. Manual testing and GC stress tests are the main guardrails today. One CI job running the test suite under **AddressSanitizer** would catch use-after-free on arena memory, double-free edge cases, and OOB bugs that only appear under unusual allocation patterns.

**Roadmap ref:** tier 2 -- MASTER_PLAN **H2** (bumped priority after runtime hardening).

---

### 6. String intern table retention

**Severity:** low for most programs; worth documenting.

After the hash-map fix, intern **lookup** is O(1). Strings are still retained in the intern table until GC sweep removes unmarked entries. Programs that generate many **unique** strings (e.g. thousands of formatted strings with numbers) can grow the table between collections.

**Fix:** document clearly; optional `koda_intern_clear()` for programs that need an explicit flush.

**Roadmap ref:** tier 3.

---

## Things in better shape than they look

### GC for games

Tri-generational collection, nursery, incremental major GC, arena for per-frame objects, and an O(1) write barrier path are **ahead of most languages** targeting this audience. Most beginners never need to think about any of it -- call `gcFrameStep(0.5)` in the loop and move on.

### vec2 / vec3 stdlib

Implemented in Koda rather than C, so hot paths allocate objects. Wasteful for production physics; **completely fine for learning**. Names (`add`, `dot`, `normalize`) are the right ones.

### ok / err error model

Avoids exceptions (confusing) and unchecked-null hell (error-prone) without Rust-style signatures. Missing piece: propagation sugar -- a `try`-style expression that early-returns the `err` value would cut boilerplate.

---

## Longer-term directions

Worth planning; not blockers for the game-maker audience.

| Topic | Why it matters |
|-------|----------------|
| **Optional / nullable types** (`let name: string?`) | Compile-time nudge before dereferencing `null` |
| **Interfaces / traits** | "Takes anything with `draw()` and `update()`" without a class hierarchy |
| **Package manager** | `#include` scales to one repo; community needs `koda install`-style sharing |
| **WASM target** | `wasm32-unknown-unknown` + slim runtime → browser games, embedded scripting |
| **Full DWARF debug info** | `--debug` passes `-g` today; source-line crashes need `DIFile` / `DISubprogram` in IR |

See [MASTER_PLAN.md](MASTER_PLAN.md) for detailed engineering refs (**D1**, **F1–F7**, etc.).

---

## Verdict

| Layer | Assessment |
|-------|------------|
| **Compiler / runtime internals** | Genuinely good shape -- ahead of many scripting runtimes targeting beginners |
| **Correctness (post review)** | Recent fixes closed real bugs (tombstones, loop IR, GC hot paths), not hypotheticals |
| **Language design** | Sound for the intended audience |
| **Product gap** | Better answer for **"I want to make games"** than **"I want to learn systems programming instead of C"** |

**To move the needle on the second audience:** opt-in **integer types** and **struct methods** are the two highest-impact language features.

---

## Related

- [Implementation status](status.md) -- what works in the tree today
- [Roadmap](ROADMAP.md) -- prioritized work queue
- [Master plan](MASTER_PLAN.md) -- detailed engineering matrix
- [Game development](guides/game-dev.md) -- loops, GC, arena usage
- [From C](guides/from-c.md) -- expectations for C migrants
