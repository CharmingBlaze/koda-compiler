# CLI reference

Every `koda` subcommand. Run `koda help` or `koda <command> --help` for inline help.

---

## Command summary

| Command | Purpose |
|---------|---------|
| `koda new` / `init` | Create project from template |
| `koda run` / `native` | Compile and run (temp exe) |
| `koda watch` | Rebuild/rerun on `.koda` changes |
| `koda check` | Parse + type-check (`./...` supported) |
| `koda lint` | `check` + `fmt --check` on paths |
| `koda fmt` | Format sources (`./...`, `--check`) |
| `koda build` | Native executable |
| `koda bundle` | Package exe + assets |
| `koda test` | Run test `.koda` files |
| `koda bench` / `profile` | Time repeated runs |
| `koda debug` | Run with debug symbols |
| `koda eval` | One-line snippet |
| `koda repl` | Interactive REPL |
| `koda clean` | Remove build artifacts (`--cache`) |
| `koda doctor` | SDK health check (`--fix` refreshes stale raylib shim) |
| `koda paths` | Machine-readable paths |
| `koda env` | Print `KODA_*` (`--export`) |
| `koda completions` | Shell completion scripts |
| `koda update` | Release update hints |
| `koda doc` | Doc path helpers |
| `koda lsp` | Language server (stdio) |
| `koda disasm` | Print LLVM IR |
| `koda wrap` | Forward to `kodawrap` (generate, upgrade, install wrappers) |
| `koda setup` | Configure Raylib shim or full wrapper in a project |
| `koda version` | Version string |
| `koda help` | Help (`help <command>`) |

---

## Passing arguments to your program

Use `--` after Koda flags:

```bash
koda run game.koda -- --level 3 --fullscreen
koda watch src/main.koda -- --debug-ai
koda bench game.koda --count 10 -- game.koda --fast
```

---

## `koda run`

```bash
koda run [--no-opt] [--debug] [<file.koda>] [-- <program args...>]
```

Uses `koda.json` entry when no file is given.

---

## `koda check` and `koda lint`

```bash
koda check                    # project entry
koda check src/main.koda      # one file
koda check ./...              # all .koda under tree

koda lint                     # default ./...
koda lint src tests
```

`lint` runs semantic check plus formatting check (does not rewrite files).

---

## `koda test`

```bash
koda test [--no-opt] [-v] [--failfast] [-run <pattern>] [<files...>]
```

| Flag | Effect |
|------|--------|
| `-v` | Print each test path |
| `--failfast` | Stop on first failure |
| `-run` | Substring filter on paths |
| `--no-opt` | Unoptimised native build |

Default: `tests/*.koda` in project (if `koda.json`) or SDK repo.

---

## `koda bench` / `profile`

```bash
koda bench [--count N] [--warmup N] [--no-opt] [--debug] <file.koda> [-- <args...>]
```

Times full compile+run cycles (includes compile cost). Use for rough comparisons.

---

## `koda debug`

Same as `koda run --debug` — debug symbols, unoptimised build.

---

## `koda eval` and `koda repl`

```bash
koda eval 'print(1 + 2)'
koda eval 'let x = 3; print(x)'

koda repl
```

REPL compiles each input line (expressions are auto-wrapped in `print(...)`). Type `exit` to quit.

---

## `koda clean`

```bash
koda clean [<dir>] [--cache]
```

Removes `dist/`, `.koda_build/`, default exes. `--cache` clears temp `koda_*` dirs in the system temp folder.

---

## `koda env`

```bash
koda env
koda env --export    # export KODA_VERSION=... lines for bash
```

---

## `koda completions`

```bash
koda completions bash >> ~/.bashrc
koda completions zsh  >> ~/.zshrc
koda completions fish >> ~/.config/fish/completions/koda.fish
```

---

## `koda update`

```bash
koda update
koda update --check-only
```

Prints current version and link to GitHub Releases (manual SDK zip upgrade).

---

## `koda doc`

```bash
koda doc                      # bundled doc paths
koda doc stdlib               # list @modules
koda doc module @math         # print stdlib/math.koda
```

---

## `koda lsp`

Stdio JSON-RPC language server for editors:

```bash
koda lsp
```

Supports `initialize`, `textDocument/didOpen`, `didChange`, and `publishDiagnostics` (via `koda check` pipeline).

---

## `koda bundle`

```bash
koda bundle [<file.koda>] [-o <dir>]
```

`-o` can appear before or after the file path.

---

## `koda new` / `init`

```bash
koda new <name> [--template hello|game|graphics|raylib]
koda init <name>   # alias
```

| Template | Contents |
|----------|----------|
| **hello** | Console print demo |
| **game** | Text lunar lander |
| **graphics** | `@game` bouncing ball + Raylib shim |
| **raylib** | Full Raylib wrapper project |

---

## `koda setup raylib`

```bash
koda setup raylib [--full] [project-dir]
```

Refreshes `wrappers/raylib_shim/` (overwrites stale files), sets `"graphics": true` in `koda.json`. Use when `@game` reports undefined shim symbols.

| Flag | Result |
|------|--------|
| (default) | Beginner shim (~33 functions) + `@game` |
| `--full` | Full wrapper (548 functions) + `#include "@raylib"` |

---

## `koda wrap`

```bash
koda wrap -name mylib -headers ./include/mylib.h -out ./wrappers/mylib
koda wrap upgrade wrappers/mylib
koda wrap check wrappers/mylib
koda wrap list
koda wrap install raylib --project
```

See [Wrapping libraries](../guides/wrapping-libraries.md).

---

## `koda doctor`

```bash
koda doctor
koda doctor --fix
```

Checks SDK toolchain, raylib, wrapper header drift, and **project raylib_shim** drift (missing `@game` symbols). `--fix` copies the canonical shim from the SDK into the current project.

---

## `koda.json`

| Field | Purpose |
|-------|---------|
| `entry` | Main `.koda` file |
| `lint` | `"beginner"` or `"strict"` |
| `bundle.assets` | Dirs copied into bundle |
| `native.sources` | C/C++ glue |
| `native.graphics` | Auto Raylib link flags when `true` |
| `native.linkflags` | Extra linker flags |

Environment: `KODA_NATIVE_SOURCES`, `KODA_LINKFLAGS` override manifest when set. Graphics projects in `koda.json` clear mismatched stale env.

---

## Related

- [commands.md](../commands.md) — extended CLI document
- [Beginner's guide](../beginners-guide.md)
- [Building and shipping](../learn/10-building-and-shipping.md)
