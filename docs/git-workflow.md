# Git workflow (contributors)

This repository uses **`main`** as the default branch. The remote is conventionally named **`origin`** (for example `https://github.com/CharmingBlaze/koda-compiler.git`).

## Committing changes

From the repository root:

```bash
git status
git add -A
git commit -m "Short imperative summary (optional body)"
```

Use complete sentences in the commit body when the change needs context. Follow **[CONTRIBUTING.md](../CONTRIBUTING.md)** (`go vet`, `go test`, `koda fmt --check` where relevant).

## Pushing to `main`

```bash
git pull --rebase origin main
git push origin main
```

- Prefer **`git pull --rebase`** before **`git push`** so local commits replay cleanly on top of **`origin/main`**.
- Avoid **`git push --force`** to **`main`** unless project policy explicitly allows it.

## Credentials

Pushing requires Git authentication for **`origin`** (SSH key, HTTPS credential helper, or **`gh auth`**). Environments without configured credentials cannot push until a maintainer signs in or installs credentials.

## Large or generated paths

See **`.gitignore`** (e.g. **`.KODA_build/`**, **`dist/`**). After **`scripts/build-runtime.sh`** / **`build-runtime.ps1`**, **`runtime/libkoda_runtime.a`** and **`runtime/obj/*.o`** may change; include them in commits when **C/runtime sources** changed so codegen link tests and local builds stay consistent.

## Releases

Tagging **`v*`** triggers **[`.github/workflows/release.yml`](../.github/workflows/release.yml)**. See **[releasing.md](releasing.md)**.
