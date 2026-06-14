# Release builds

## `koda` CLI (Go)

Build a standalone binary (no cgo):

```bash
CGO_ENABLED=0 go build -ldflags="-s -w" -o koda ./cmd/koda
```

On Windows (PowerShell):

```powershell
$env:CGO_ENABLED=0; go build -ldflags="-s -w" -o koda.exe ./cmd/koda
```

Cross-compile from Linux, for example:

```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o koda.exe ./cmd/koda
```

`koda run`, `check`, and `disasm` need no C toolchain. **`koda build`** emits **LLVM IR** (`.ll`), runs **llc** to produce an object file, then invokes **Clang** (on `PATH`, or **`KODA_CLANG`** / **`CC`**) to link **`runtime/libkoda_runtime.a`** and headers from **`runtime/src/`**. Optional **`KODA_USE_LLD=1`** adds **`-fuse-ld=lld`** when you want the LLVM linker.

## Checksums

When publishing releases, attach `SHA256SUMS` for each artifact (generate with `sha256sum` / `Get-FileHash`).
