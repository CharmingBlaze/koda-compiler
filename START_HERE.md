# Start here — Koda in 5 minutes

**Koda replaces C and C++ for beginners** who want to make **games and desktop apps** without installing Go, Python, LLVM, Visual Studio, or a package manager.

Download one SDK zip. Unzip. Run `koda`. Ship a native `.exe` (or binary on Mac/Linux).

---

## What you get (and what you do not need)

| You get | You do **not** need |
|---------|---------------------|
| Native `.exe` / binary (like C/C++) | Go |
| Structs, enums, `func`, game loops | Python |
| Built-in graphics path (`@game` + Raylib in the zip) | Node.js or npm |
| Automatic memory (GC) for game logic | A separate C/C++ compiler to **use** Koda |
| `koda run` — edit and play instantly | Rust, Cargo, or CMake to get started |

Release builds **embed Clang, LLVM, and the runtime inside `koda`**. The first compile unpacks tools to a temp folder on your machine. Nothing is downloaded from the internet at compile time.

---

## Install (Windows)

1. Download **`koda-*-sdk-windows-amd64.zip`** from [GitHub Releases](https://github.com/CharmingBlaze/koda-compiler/releases).
2. Unzip to a folder, e.g. `C:\koda`.
3. Open **PowerShell** in that folder:

```powershell
.\koda.exe doctor
.\koda.exe new bounce --template graphics
cd bounce
..\koda.exe run
```

Optional — add Koda to your PATH (run once from the SDK folder):

```powershell
powershell -ExecutionPolicy Bypass -File scripts\install-koda.ps1
```

---

## Install (macOS / Linux)

1. Download the SDK zip for your platform from [GitHub Releases](https://github.com/CharmingBlaze/koda-compiler/releases).
2. Unzip and open a terminal in that folder:

```bash
chmod +x koda kodawrap
./koda doctor
./koda new bounce --template graphics
cd bounce
../koda run
```

Optional — install to `~/.local/bin`:

```bash
bash scripts/install-koda.sh
```

---

## Your first console app

```bash
koda new mytool
cd mytool
koda run
koda build -o mytool
```

You now have a standalone program you can copy to friends. They **do not** install Koda to run it.

---

## Why Koda instead of C or C++?

| | C / C++ | Koda |
|---|---------|------|
| First game | Install compiler, SDK, learn headers & linking | Unzip SDK, `koda new bounce --template graphics` |
| Memory | You `malloc` / `new` / smart pointers | GC for gameplay; C only at library edges |
| Strings & JSON | Libraries or pain | Built in |
| Graphics (Raylib) | Download, link, configure CMake | `@game` + `"graphics": true` in `koda.json` |
| Typo in a variable | Silent wrong behavior | `koda check --warn-unused` |
| Ship to players | Your `.exe` + maybe DLLs | `koda bundle` + `assetPath()` |

Koda is **not** for operating-system kernels or firmware. It **is** for the C/C++ you write for games, tools, and apps today.

---

## Learn more

| Doc | Purpose |
|-----|---------|
| [Beginner's guide](docs/beginners-guide.md) | Full tutorial |
| [Game development](docs/guides/game-dev.md) | Loops, input, shipping |
| [From C / C++](docs/guides/from-c.md) | Side-by-side migration |
| [FAQ](docs/faq.md) | Common questions |

When something breaks: **`koda doctor`** — fix every line marked **FAIL**.
