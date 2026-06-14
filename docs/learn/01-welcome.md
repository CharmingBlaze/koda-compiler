# Chapter 1 — Welcome to Koda

**You will learn:** what Koda is, how it compares to C, and the basic compile-run workflow.

**Time:** ~5 minutes.

---

## What you are learning

Koda is a language for **games and applications** that **compiles to native code**. You write `.koda` files and use the `koda` command-line tool to run or build them.

```text
  hello.koda  ──►  koda run  ──►  native program  ──►  output
```

Unlike interpreted languages, your users receive a **single executable** (plus assets you choose to bundle). There is no Koda runtime they must install.

---

## Why people choose Koda

| Benefit | Detail |
|---------|--------|
| Native output | Same end result as C — one binary |
| Fast iteration | `koda run` and `koda watch` without manual link steps |
| Familiar syntax | `struct`, `func`, `if`, `while` like C; objects like JS |
| Library access | Call C/C++ via wrappers (`kodawrap`) |
| Batteries included | `stdlib/` — math, json, io, timers, vectors |

---

## What Koda is not

- Not a browser language (no DOM).
- Not a managed VM language (no JVM, no bytecode interpreter for users).
- Not a replacement for your renderer — use Raylib/SDL via wrappers for graphics.

---

## The three commands you will use first

```bash
koda run hello.koda    # compile + run (temporary exe)
koda build -o hello    # compile to permanent exe
koda check hello.koda  # errors only, no binary
```

---

## Try it yourself

If you already have Koda installed, create `hello.koda`:

```koda
print("I am learning Koda");
```

Run `koda run hello.koda`. If you see the message, you are ready for [Chapter 2 — Install](02-install-and-first-run.md).

If not installed yet, start with Chapter 2.

---

## Next chapter

[Chapter 2 — Install and first run](02-install-and-first-run.md)
