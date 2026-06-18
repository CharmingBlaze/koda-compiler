# Koda positioning

An honest picture of what Koda is, who it is for, and how it replaces C/C++ for beginners making games and apps.

> **Audience:** contributors, release planners, and technical writers. End users should start with **[START_HERE.md](../START_HERE.md)** or the [Beginner's guide](beginners-guide.md).

---

## The pitch

**Koda is the best beginner language for native games and applications** when the alternative is learning C or C++ first — with a **serious ceiling**: full Raylib, generated C/C++ wrappers, structs, methods, and native performance. See [KODA_LANGUAGE_ROADMAP.md](KODA_LANGUAGE_ROADMAP.md).

| Promise | Delivery |
|---------|----------|
| **Easy install** | One SDK zip — no Go, Python, LLVM, or IDE required to compile |
| **Native output** | Same as C/C++: a real executable, no VM on the player's machine |
| **Games + apps** | `koda.game`, Raylib in the zip, `koda new --template graphics` |
| **Beginner-safe** | GC, typo hints, `--warn-unused`, `koda doctor` |

Koda is **not** Python (no interpreter workflow). Koda is **not** a toy — LLVM backend, generational GC, cross-platform CI.

---

## Where Koda sits

```
  C / C++                         Koda
  Native, manual setup       →    Native, one zip install
  malloc / new / leaks       →    GC for gameplay code
  CMake + link flags         →    koda.json + koda build
  Steep first game           →    koda new bounce --template graphics
```

Compared to **Python or JavaScript** for games: Koda compiles to a **real native binary** — players do not install a runtime.

Compared to **C/C++**: Koda trades kernel-grade control for **approachable syntax, built-in stdlib, and zero toolchain setup** for end users.

---

## Is Koda a good C/C++ alternative for beginners?

**Yes — for games, tools, and desktop apps.**

| Audience | Fit today |
|----------|-----------|
| **"I want to make games"** | Ready — full Raylib (`use raylib;`, 548 fn), optional `koda.game`, templates, `koda bundle` |
| **"I want desktop/CLI apps"** | Ready — I/O, JSON, native binary |
| **"I want to learn C/C++ syntax without the pain"** | Ready — structs, methods, enums, native output |
| **"I want OS kernels / firmware"** | Not Koda — use C/Rust |

---

## Strengths already in place

| Area | What you get |
|------|----------------|
| **Install** | SDK zip with embedded Clang/llc/runtime — no network fetch at compile time |
| **Language** | Structs **with methods**, opt-in `i32`/`u8`, enums, closures, defer |
| **Diagnostics** | Typo hints, `--warn-unused`, enum exhaustiveness warnings |
| **GC** | Tri-generational, arena, `gcFrameStep` for game loops |
| **Graphics** | `koda.game` over Raylib; `"graphics": true` in `koda.json` |
| **Ship** | `koda build`, `koda bundle`, `assetPath()` |

See [status.md](status.md) for engineering detail.

---

## Remaining gaps (honest)

| Gap | Impact | Status |
|-----|--------|--------|
| ASAN CI blocking | Contributor safety net | Non-blocking CI job exists |
| Full gdb line mapping | Debug experience | Partial (`--debug`) |
| Package manager | Community libraries | Long-term |
| WASM target | Browser games | Long-term |

These do **not** block the beginner install-and-ship story.

---

## Verdict

| Layer | Assessment |
|-------|------------|
| **End-user install** | One zip, `koda doctor`, optional PATH scripts |
| **vs C/C++ for beginners** | Strong — native output without toolchain ceremony |
| **vs Python for games** | Native binary, no runtime dependency for players |
| **Compiler internals** | Production-grade LLVM + GC |

---

## Related

- [START_HERE.md](../START_HERE.md) — 5-minute install
- [From C / C++](guides/from-c.md) — migration guide
- [status.md](status.md) — what works today
- [ROADMAP.md](ROADMAP.md) — contributor priorities
