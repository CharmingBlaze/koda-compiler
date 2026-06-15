# Getting started with Koda

The fastest path from zero to a **native game or app** — no Go, Python, or LLVM install.

For the full story see **[START_HERE.md](../../START_HERE.md)** or the **[Beginner's guide](../beginners-guide.md)**.

---

## Install (2 minutes)

1. Download the **SDK zip** for your OS from [GitHub Releases](https://github.com/CharmingBlaze/koda-compiler/releases).
2. Unzip so **`stdlib/`** is next to **`koda`**.
3. Run:

```bash
koda doctor
koda version
```

You do **not** need Go, Python, LLVM, Visual Studio, Node, or Rust for release binaries.

| Platform | Optional: add to PATH |
|----------|------------------------|
| Windows | `powershell -File scripts\install-koda.ps1` |
| macOS / Linux | `bash scripts/install-koda.sh` |

---

## First run

```bash
koda new myapp
cd myapp
koda run
```

Or a single file:

```koda
print("Hello, Koda!");
```

```bash
koda run hello.koda
```

---

## Daily commands

```bash
koda run [--debug] [-- <args>]    # compile + run; pass args after --
koda watch                        # rerun on save
koda build -o app                 # native executable
koda check ./...                  # check all sources
koda lint                         # check + format check
koda test -v -run io              # filtered tests
koda bench game.koda --count 5    # timing
koda eval 'print(2 + 2)'          # one-liner
koda repl                         # interactive
koda bundle -o dist/app
koda clean --cache
koda doctor
koda env --export
koda help build                   # per-command help
```

Full CLI: [reference/cli.md](reference/cli.md) · [commands.md](../commands.md)

---

## Templates

```bash
koda new lander --template game      # text lunar lander
koda new bounce --template graphics  # Raylib (needs link flags)
```

---

## Next steps

| Goal | Read |
|------|------|
| Full tutorial | [Beginner's guide](beginners-guide.md) |
| Chapter by chapter | [Learn path](learn/README.md) |
| Games | [Game development](game-dev.md) |
| Desktop / CLI apps | [Applications](applications.md) |
| From C | [From C](from-c.md) |
| Every syntax form | [Language reference](../../language.md) |
| Stuck? | [FAQ](../faq.md) · [Troubleshooting](../troubleshooting.md) |

---

[Documentation hub](../README.md)
