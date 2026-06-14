# Chapter 2 — Install and first run

**You will learn:** how to install the SDK, verify it works, and create a project with `koda new`.

**Time:** ~10 minutes.

---

## Install the SDK

1. Go to [GitHub Releases](https://github.com/CharmingBlaze/koda-compiler/releases).
2. Download the **SDK zip** for your platform (Windows x64, Linux, macOS).
3. Unzip to a folder, for example `C:\koda` or `~/koda`.
4. Confirm layout:

```text
koda/
  koda.exe          (or `koda` on Unix)
  kodawrap.exe
  stdlib/           ← required
  docs/
```

5. Add the folder to your **PATH** (optional but convenient).

> **Note:** You do **not** install Go or LLVM to use release binaries. Those are only for building Koda from source.

Verify:

```bash
koda version
koda doctor
```

`koda doctor` checks for `stdlib/` and common configuration issues.

---

## First run — single file

Create `hello.koda` anywhere:

```koda
print("Hello, Koda!");
```

```bash
koda run hello.koda
```

Expected output:

```text
Hello, Koda!
```

---

## First project — `koda new`

```bash
koda new myapp
cd myapp
koda run
```

This creates:

| Path | Purpose |
|------|---------|
| `koda.json` | Project manifest (entry, bundle, native) |
| `src/main.koda` | Main source file |
| `assets/` | Images, sounds, data |
| `README.md` | Project notes |

Templates:

```bash
koda new lander --template game      # text lunar lander
koda new bounce --template graphics  # Raylib demo (needs link flags)
```

---

## Useful first commands

```bash
koda help              # all commands
koda check src/main.koda
koda fmt src/main.koda
koda watch             # rerun on save (in project dir)
```

---

## Try it yourself

1. Run `koda new practice` and change the message in `src/main.koda`.
2. Run `koda run` and confirm your change appears.
3. Run `koda check` with a deliberate typo and read the error message.

---

## Next chapter

[Chapter 3 — Variables and types](03-variables-and-types.md)
