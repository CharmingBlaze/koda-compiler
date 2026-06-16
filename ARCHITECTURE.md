# Koda — Architecture Overview

For the full compiler pipeline and invariants, see [`docs/architecture.md`](docs/architecture.md) and [`docs/handoff.md`](docs/handoff.md).

## Repository layout

| Path | What it is |
|---|---|
| `cmd/koda/` | CLI entry point — `koda build`, `koda run`, `koda check`, etc. |
| `cmd/wrapgen/` | `kodawrap` — C header → Koda bindings generator |
| `internal/lexer/` | Tokenizer |
| `internal/parser/` | AST types + recursive-descent parser |
| `internal/sema/` | Semantic analysis, escape analysis, shadow layout |
| `internal/codegen/` | LLVM IR emission (Go, using llir/llvm) |
| `internal/nativebuild/` | Invokes llc + clang + linker |
| `internal/formatter/` | `koda fmt` |
| `internal/diagnostic/` | Error reporting with Rust-style source snippets |
| `internal/kodahome/` | Toolchain discovery and embedded binary extraction |
| `internal/embed/` | Release-build: embedded llc, lld, libkoda_runtime.a |
| `runtime/src/` | C runtime: GC, NaN-boxing, objects, shadow stack |
| `stdlib/` | Standard library as `.koda` files |
| `wrappers/` | Pre-generated Raylib and other C library bindings |
| `api/` | Go API for embedding the Koda compiler |
| `tests/` | `.koda` test programs |
| `examples/` | Sample programs and games |
| `scripts/` | Build scripts for runtime and release packages |
| `docs/` | All documentation |
| `dist-template/` | Template for what goes into the release SDK zip |
| `_legacy/` | Old artifacts kept for reference — not part of the build |

## Pipeline

```
source.koda
  → internal/lexer       (tokens)
  → internal/parser      (AST)
  → internal/sema        (analysis, escape analysis)
  → internal/codegen     (LLVM IR)
  → llc                  (object file)
  → clang + libkoda_runtime.a  (native binary)
```

## Key invariant

Every symbol in `internal/codegen/runtime.go` must have a matching implementation in `runtime/src/koda_runtime.c` with the exact same C calling convention. If these drift, the linker produces a broken binary silently.

Automated guards (run via `go test ./...`):

- `TestWriteBarrierCoverage` — GC write barriers in `koda_runtime.c`
- `TestRuntimeLLVMSymbolsDefinedInArchive` — LLVM declares match `libkoda_runtime.a`
- `TestRuntimeLLVMSymbolsMentionedInSources` — every declare appears in `runtime/src/*.c`

## Third-party layout

| Path | Role |
|---|---|
| `third_party/raylib_static/` | **Canonical** vendored Raylib — `make -C third_party/raylib_static` builds into `stage/`; used by compiler auto-detect, SDK packaging, and `koda doctor` |
| `raylib_lib/` | **Dev-only** Windows prebuild used by local `.ps1` scripts — not used by the compiler or release SDK |

## Legacy artifacts

`_legacy/` holds old wrapper-generator output for reference only. It is not part of the build, CI, or release zip. Prefer `wrappers/` and `third_party/` for current Raylib integration.

## Running tests

```bash
go test ./...
bash scripts/build-runtime.sh
./koda run tests/hello.koda
```
