# Cutting a Koda release

**Audience:** maintainers tagging **`v*`** on the default branch so **[`.github/workflows/release.yml`](../.github/workflows/release.yml)** publishes artifacts.

## Checklist

1. **`CHANGELOG.md`** — Move notes from **`[Unreleased]`** into a dated **`[X.Y.Z]`** section (see **`[0.2.0]`** as a template).
2. **Version strings** — Align **`cmd/koda/main.go`** `version` and **`cmd/wrapgen/wrapgen_version.go`** `WrapgenVersion` with the tag (release builds may also pass `-ldflags "-X main.version=..."` if you use that flow).
3. **CI green** — **`go vet ./...`**, **`go test ./...`**, and **`koda fmt --check`** (see **[CONTRIBUTING.md](../CONTRIBUTING.md)**).
4. **Tag** — From a **git clone** of `main` with the version bump merged:
   ```powershell
   powershell -ExecutionPolicy Bypass -File scripts/push-release-tag.ps1 -Tag v0.5.0
   git push origin main
   git push origin v0.5.0
   ```
   Or manually: `git tag -a v0.5.0 -m "Release v0.5.0"` then `git push origin v0.5.0`. The workflow builds **`-tags release`** `koda` binaries with embedded Clang + **`libkoda_runtime.a`** (and **lld** on Windows). Release jobs run **`scripts/ci-release-smoke.sh`** before uploading artifacts.
5. **Post-release** — Open **`[Unreleased]`** again; bump dev versions (e.g. **`0.6.0-dev`**) on `main` after **`v0.5.0`**.

## Windows: test the SDK layout before tagging

From repo root, after LLVM + MinGW paths match **`scripts/build-release.ps1`** defaults:

```powershell
powershell -File scripts/build-release.ps1 -PackageSdk
```

This mirrors the **`koda-*-sdk-windows-amd64`** folder inside the release zip (see **`scripts/assemble-offline-sdk.ps1`**).

## Notes

- **Linux CI** (**`ci.yml`**) does not embed the toolchain; **`release.yml`** repopulates **`internal/embed/...`** per job from the runner’s Clang/LLVM before **`go build -tags release`**.
- Smoke tests in release jobs run **`scripts/ci-release-smoke.sh`** (hello, GC shadow, struct methods, integer types, stdlib modules, warn-unused) after each platform binary is produced.
