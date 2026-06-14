# Distribution guide

See [DISTRIBUTION_GUIDE.md](../../DISTRIBUTION_GUIDE.md) for the full reference.

This page summarizes the key commands.

---

## Quick reference

| Goal | Command |
|------|---------|
| Run during development | `koda run game.koda` |
| Build native executable | `koda build game.koda -o game.exe` |
| Create distributable folder | `koda bundle game.koda -o dist/game` |
| Include extra files in bundle | Set `KODA_BUNDLE_FILES` |
| Use a C library wrapper | Set `KODA_NATIVE_SOURCES` and `KODA_LINKFLAGS` |

---

## Bundle output

```
dist/mygame/
  mygame.exe          <- compiled application
  run.bat             <- Windows launcher
  README.md           <- user-facing instructions
  bundle-info.txt     <- build metadata
  (extra files)       <- DLLs, assets, licenses from KODA_BUNDLE_FILES
```

---

## Environment variables

| Variable | Purpose |
|----------|---------|
| `KODA_CLANG` | Path to Clang executable |
| `CC` | Fallback compiler |
| `KODA_USE_LLD` | Set `1` to use LLD linker |
| `KODA_PATH` | Extra Koda source search paths |
| `KODA_WRAPPERS` | Extra wrapper search paths |
| `KODA_NATIVE_SOURCES` | C/C++ wrapper glue files |
| `KODA_LINKFLAGS` | Flags passed to Clang linker |
| `KODA_BUNDLE_FILES` | Extra files copied into bundle |

---

## Cross-compiling the CLI

```powershell
$env:GOOS = "linux"; $env:GOARCH = "amd64"
go build -o koda-linux-amd64 ./cmd/koda
```

Supported targets: `windows/amd64`, `linux/amd64`, `linux/arm64`, `darwin/amd64`, `darwin/arm64`.
