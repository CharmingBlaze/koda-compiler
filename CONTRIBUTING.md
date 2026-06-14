# Contributing to Koda

## ⚠️ This is for contributors only

If you want to **use** Koda to make games or applications, **do not build from source**. On [GitHub Releases](https://github.com/CharmingBlaze/koda-compiler/releases), download the **SDK zip** for your platform (one unzip gets **`koda`**, **`kodawrap`**, **`stdlib/`**, docs, examples, and vendored Raylib where applicable). That layout works **fully offline** after download — the compiler does not fetch LLVM or libraries from the network. Loose **`koda-*`** binaries are also listed if you supply **`stdlib/`** yourself.

Building from source requires Go 1.22+, Clang + llc on PATH, and a C toolchain. The embedded-toolchain release build (`-tags release`) additionally requires having the LLVM binaries available to embed. This is a non-trivial setup that is only worth doing if you are modifying the compiler itself.

## Build

From the repo root (Go **1.22+**):

```bash
go build -o bin/koda ./cmd/koda
```

Build the C runtime archive when linking native binaries (needed for **`go test`** / **`koda build`**):

- Linux / macOS: `bash scripts/build-runtime.sh`
- Windows: `powershell -File scripts/build-runtime.ps1`

## Tests

GitHub Actions **CI** runs on **Ubuntu**, **macOS**, and **Windows** on every push and pull request (`go vet`, `go test`, **`koda fmt --check ./...`**, native smoke, GC soak, stress smoke, parser fuzz, wrapgen Raylib link). See **`.github/workflows/ci.yml`** and **`scripts/ci-native-smoke.sh`**, **`scripts/ci-wrapgen-raylib.sh`**, **`scripts/ci-gc-stress-timed.ps1`**. From the repo root locally:

```bash
go vet ./...
go test ./... -count=1
```

## Releases

See **[docs/releasing.md](docs/releasing.md)** for version bumps, changelog, and tagging **`v*`** so **`.github/workflows/release.yml`** runs.

### Offline SDK folder (Windows maintainers)

GitHub Releases attach **`koda-<tag>-sdk-windows-amd64.zip`** (compiler + **kodawrap** + stdlib + **docs** + **wrappers** + **examples** + vendored **raylib 5.0**). CI builds that zip on Linux; locally on Windows you can reproduce the same tree:

```powershell
powershell -File scripts/build-release.ps1 -PackageSdk
```

That runs **`scripts/assemble-offline-sdk.ps1`**, which downloads the official raylib Windows prebuild into **`third_party/raylib_static/stage`**. Add **`-PackageSdkZip`** to also write **`dist/koda-<version>-sdk-windows-amd64.zip`**. See **`docs/distribution.md`** §7.

## Git: commit and push to `main`

Step-by-step: **[docs/git-workflow.md](docs/git-workflow.md)**.

## Formatting

Canonical `.koda` style:

```bash
go build -o bin/koda ./cmd/koda
./bin/koda fmt ./...
./bin/koda fmt --check ./...
```

CI runs **`koda fmt --check ./...`** after building **`koda`**.

## Compiler diagnostics

Lexer and parser failures include the source line and a caret when reporting errors from **`parser.LoadProgram`** (used by **`koda run`**, **`koda build`**, etc.). **`koda check`** runs the same load step plus **`sema.PrepareNativeBundle`** (semantic errors, no LLVM).

## Native toolchain

LLVM **`clang`** / **`llc`** (and optionally **`lld`**) are required for **`koda build`** / **`koda run`** unless you use a **release** build that embeds the toolchain. See **`koda doctor`** and **`koda help`**.

## Example programs

- **`tests/`** — language and smoke tests.
- **`examples/games/`** — small standalone demos (**`examples/games/README.md`**).
- **`demos/`** — larger samples (some need Raylib + wrappers).
