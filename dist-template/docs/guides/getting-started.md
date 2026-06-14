# Getting started with Koda (beginner path)

Koda is a **modern language for games and applications** that compiles to a **native binary** — a practical alternative to C for game logic, tools, and desktop apps.

**Documentation hub:** [README.md](../README.md) · **Coming from C:** [from-c.md](from-c.md)

**You do not install Go or LLVM** to use Koda. Download **`koda`** and **`kodawrap`** (and unpack an SDK zip so **`stdlib/`** sits next to them). Release binaries embed the compiler toolchain; that is separate from the Go + LLVM setup **only maintainers** use to build Koda from source.

## 1) Install Koda the easy way

Download the latest `koda` (and `kodawrap`) from [GitHub Releases](https://github.com/CharmingBlaze/koda-compiler/releases), use an SDK zip so **`stdlib/`** is beside the executables, then run:

```bash
koda version
```

Do not compile Koda from source for normal usage. The source tree is maintainer-focused and not the intended beginner install route.

---

## 2) Your first project

Create a new project:

```bash
koda new mygame              # hello template (default)
koda new lander --template game       # text game, runs immediately
koda new bounce --template graphics   # Raylib window (set link flags first)
cd mygame
koda run
```

Or create `hello.koda` by hand:

```koda
print("Hello, Koda!");
```

Run it:

```bash
koda run hello.koda
```

---

## 3) Commands you will use every day

```bash
# Run (temporary executable)
koda run game.koda

# Build a native executable
koda build game.koda -o game.exe

# Debug-friendly build
koda build --debug game.koda -o game_debug.exe

# Check parse + semantic errors
koda check game.koda

# Format source
koda fmt game.koda
koda fmt --check .

# Rebuild/rerun when files change
koda watch game.koda

# Package for sharing
koda bundle game.koda -o dist/MyGame
```

For **every** command, flags, and copy-paste examples, see **[docs/commands.md](../commands.md)** (or run **`koda help`**).

---

## 4) Using the wrapper tool (`kodawrap`)

`kodawrap` generates `.koda` bindings + `wrapper.c` from C/C++ headers.

```bash
koda wrap --help
```

Typical flow:

1. Generate bindings from a header.
2. Import generated `.koda` module in your game.
3. Build/run with native glue via `KODA_NATIVE_SOURCES` and linker flags via `KODA_LINKFLAGS`.

Full details: `docs/wrappers.md`.

---

## 5) Learn the whole language

- `docs/using-the-language.md` — **start here**: how to use the language end-to-end (syntax, types, control flow, modules, builtins, stdlib)
- `language.md` — complete language catalog (operators, statements, builtins, methods)
- `docs/user_guide.md` — practical beginner/intermediate walkthrough
- `docs/reference.md` — builtins and runtime-facing APIs
- `docs/guides/game-dev.md` — game-focused usage patterns
- `docs/distribution.md` — shipping and bundles

If docs and behavior ever differ, run a tiny `.koda` example and trust the CLI result.
