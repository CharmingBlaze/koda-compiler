# Koda CLI — every command

Use the **`koda`** binary from [GitHub Releases](https://github.com/CharmingBlaze/koda-compiler/releases). For full inline help:

```bash
koda help
```

---

## `koda new`

Create a new project directory with `koda.json`, `src/main.koda`, `assets/`, and a template README.

```bash
koda new <name> [--template hello|game|graphics|raylib]
```

| Template | What you get |
|----------|----------------|
| **hello** (default) | Minimal print program — runs immediately |
| **game** | Text-based lunar lander — no native libraries |
| **graphics** | Bouncing-ball window via `@game` + bundled Raylib shim |
| **raylib** | Full Raylib wrapper (548 functions) for advanced graphics |

**Examples**

```bash
koda new myapp
koda new lander --template game
koda new bounce --template graphics
cd bounce
koda doctor
koda run
```

The **graphics** template sets `"graphics": true` in `koda.json` — platform Raylib link flags are applied automatically. Run `koda doctor` if the first build fails.

---

## `koda.json` (project manifest)

When `koda.json` sits in the project root (or a parent directory), you can omit the entry file on `run`, `build`, `watch`, `check`, and `bundle`.

```json
{
  "name": "mygame",
  "version": "0.1.0",
  "entry": "src/main.koda",
  "lint": "beginner",
  "bundle": {
    "assets": ["assets"],
    "extra": ["LICENSE"]
  },
  "native": {
    "sources": ["wrappers/raylib_shim/wrapper.c"],
    "graphics": true
  }
}
```

| Field | Purpose |
|-------|---------|
| `name` | Project name (default executable / bundle name) |
| `entry` | Main `.koda` file (required) |
| `lint` | `"beginner"` or `"strict"` diagnostics |
| `bundle.assets` | Directories copied into `koda bundle` output |
| `bundle.extra` | Additional files or folders to copy |
| `native.sources` | C/C++ files linked on every build |
| `native.graphics` | When `true`, apply platform Raylib link flags automatically |
| `native.linkflags` | Extra linker flags (overrides auto graphics flags when set) |

Environment variables override manifest values when set. For graphics projects, `koda.json` also wins over stale `KODA_NATIVE_SOURCES` left in the shell.

**Project-aware commands**

```bash
cd mygame
koda run
koda build -o mygame
koda check
koda bundle -o dist/mygame
```

---

## `koda run` / `koda native`

Compile the entry **`.koda`** file to a native executable (temporary), run it, then delete the temp binary.

```bash
koda run [--no-opt] [--debug] [<file.koda>] [-- <program args...>]
koda native [--no-opt] [<file.koda>]   # same as run
```

- **`--no-opt`** — skips LLVM IR optimisation and uses a less aggressive native compile.
- **`--debug`** — debug symbols, unoptimised build.
- Omit the file when `koda.json` defines `entry`.
- Arguments after **`--`** are passed to your program (`args()` in Koda).

**Example**

```bash
koda run src/main.koda
koda run              # uses koda.json entry
koda run -- --level 3
```

---

## `koda watch`

Rebuild and rerun whenever **`.koda`** files under the entry file’s directory change.

```bash
koda watch [--no-opt] [<file.koda>] [-- <program args...>]
```

---

## `koda test`

```bash
koda test [--no-opt] [-v] [--failfast] [-run <pattern>] [<files...>]
```

---

## `koda lint`

Run `check` plus `fmt --check` on paths (default `./...`).

```bash
koda lint [paths...]
```

---

## `koda bench` / `profile`

Time repeated compile+run cycles.

```bash
koda bench [--count N] [--warmup N] [--no-opt] [--debug] <file.koda> [-- <args...>]
```

---

## `koda debug`

Run with debug symbols (`koda run --debug`).

---

## `koda eval` / `repl`

```bash
koda eval 'print(1 + 2)'
koda repl
```

---

## `koda init`

Alias for `koda new`.

---

## `koda env` / `completions` / `update` / `doc` / `lsp`

```bash
koda env [--export]
koda completions bash|zsh|fish
koda update [--check-only]
koda doc [stdlib | module @name | wrappers | wrapper @name]
koda lsp
```

See **[reference/cli.md](reference/cli.md)** for full details.

---

## `koda clean`

```bash
koda clean [<dir>] [--cache]
```

---

## `koda check`

Parse, resolve imports, and run semantic analysis only — **no** native compile.

```bash
koda check [<file.koda>]
koda check ./...
koda check --warn-unused
```

Prints **`OK`** if the program is valid.

---

## `koda fmt`

Format **`.koda`** sources with the canonical formatter.

```bash
koda fmt [--check] <file.koda> [more files...]
koda fmt [--check] ./...
```

- **`--check`** — do not write files; exit with an error if any file would change (for CI).

---

## `koda disasm`

Print **LLVM IR** for the program (after parse, sema, and codegen).

```bash
koda disasm <file.koda>
```

---

## `koda build`

Produce a **native executable** you keep.

```bash
koda build [--no-opt] [--debug] <file.koda> [-o <exe>]
```

---

## `koda bundle`

Build the program and write a **folder** you can zip or ship.

```bash
koda bundle <file.koda> [-o <dir>]
```

Copy extra files with **`KODA_BUNDLE_FILES`**. Use `assetPath("file.png")` in code for bundled assets.

---

## `koda setup`

Configure optional project integrations.

```bash
koda setup raylib [--full] [project-dir]
```

| Mode | What it does |
|------|----------------|
| Default | Copies/refreshes **`wrappers/raylib_shim/`** from the SDK (overwrites stale files), sets `"graphics": true` and shim `native.sources` |
| **`--full`** | Points at the full **`wrappers/raylib/`** wrapper (548 functions); use `#include "@raylib"` |

Run this when `@game` errors mention undefined shim symbols like `drawline` or `getmousex` — your project shim is out of date.

```bash
cd my-koda-project
koda setup raylib
koda run
```

---

## `koda wrap`

Forward to **`kodawrap`** (C/C++ header → organized package). **`koda`** looks for **`kodawrap`** next to itself, then **`PATH`**.

```bash
koda wrap --help
koda wrap -name mylib -headers ./include/mylib.h -out ./wrappers/mylib

# Regenerate after upgrading native headers
koda wrap upgrade wrappers/mylib
koda wrap upgrade @raylib
koda wrap check wrappers/mylib

# Install from built-in catalog
koda wrap list
koda wrap install raylib --project
koda wrap install sqlite3@3 -o wrappers/sqlite3
```

Full workflow: **[wrappers.md](wrappers.md)** · **[Wrapping libraries](guides/wrapping-libraries.md)**.

---

## `koda paths`

Print **machine-readable** toolchain paths (for scripts and CI).

```bash
koda paths
```

---

## `koda doctor`

Human-readable **health check**: clang, runtime library, stdlib, raylib, wrapper drift, smoke build.

```bash
koda doctor
```

Fix every **FAIL** line before graphics projects. Wrapper drift suggests `koda wrap upgrade <dir>`.

---

## `koda version`

```bash
koda version
koda --version
```

---

## `koda help`

```bash
koda help
koda help setup
koda --help
```

---

## Environment variables (quick reference)

| Variable | Purpose |
|----------|---------|
| **`KODA_HOME`** | SDK root when `koda.exe` is not beside `stdlib/` |
| **`KODA_CLANG`** / **`CC`** | Clang driver for native builds |
| **`KODA_PATH`** | Extra **`@module`** search directories |
| **`KODA_WRAPPERS`** | Pre-built wrapper **`.koda`** trees |
| **`KODA_NATIVE_SOURCES`** | C/C++ sources linked into your app (e.g. **`wrapper.c`**) |
| **`KODA_LINKFLAGS`** | Extra linker flags (**`-l`**, **`-L`**, frameworks, …) |
| **`KODA_BUNDLE_FILES`** | Extra files copied into **`koda bundle`** output |
| **`KODA_RAYLIB_STAGE`** | Override vendored raylib `stage/` directory |

See **`koda help`** for the complete list.

---

## Related

- [CLI reference](reference/cli.md)
- [Raylib guide](guides/raylib.md)
- [Game development](guides/game-dev.md)
- [Troubleshooting](troubleshooting.md)
