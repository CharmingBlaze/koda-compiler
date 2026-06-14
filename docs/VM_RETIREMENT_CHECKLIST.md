# Bytecode VM retirement — completion checklist

Cross-check against the retirement task. *(Separate git commits are recommended for history clarity; this doc tracks substance.)*

| Step | Requirement | Status |
|------|----------------|--------|
| **1** | Written audit (Sections A–B–C) | [`VM_RETIREMENT_AUDIT.md`](VM_RETIREMENT_AUDIT.md); suggested issues in [`GITHUB_ISSUES_VM_RETIREMENT.md`](GITHUB_ISSUES_VM_RETIREMENT.md) |
| **2** | `cmd/koda`: no VM branch; `koda run` → LLVM via `api.Run` | Done |
| **3** | `api/run.go`: native only (`Run`, `RunWithWriters`); `diagnose.go` static | Done |
| **4** | `internal/vm/` removed; grep `"koda/internal/vm"` → empty | Done |
| **5** | No standalone Go `type Value` for VM; AST `interface{}` literals unchanged | Verified (`grep ^type Value`) |
| **6** | `internal/parser/` tree-walking files | **N/A** — directory not present in this tree |
| **7** | `go mod tidy` after cleanup | Run locally (`go mod tidy`) |
| **8** | README, `docs/architecture.md`, CHANGELOG, ROADMAP | Updated (architecture pipeline + “what happened”; CHANGELOG **Removed**; ROADMAP note) |
| **9** | `go build ./...`, `go test ./...`, `go vet ./...`; smoke `koda build` | Run locally (Windows may block `go test` exes — see Defender exclusions) |

**CLI notes**

- `koda native <file>` is kept as an **alias** of `koda run` (same LLVM path), not a second interpreter.
- `koda disasm` prints **LLVM IR** (not historical bytecode).
