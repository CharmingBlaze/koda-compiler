# Bytecode VM retirement — pre-removal audit

Generated as part of retiring Path A (bytecode VM) in favor of Path B (LLVM → native only).

## Section A — Files owned entirely by the bytecode VM

| Path | Description |
|------|-------------|
| `internal/vm/stub.go` | Only file under `internal/vm/`: `package vm` with `type Value = any`; **no other package imported this module** (verified: `grep -r "\"koda/internal/vm\"" --include="*.go"` → empty). |

**Not present in this tree (verify on your checkout):**

- `internal/codegen.go`, `internal/codegen.go`, `internal/runtime.go`, `internal/sema.go` — **do not exist** in this workspace; older docs still reference them.

**Import audit**

- `grep -r "\"koda/internal/vm\"" --include="*.go" .` → **no matches**.
- `grep -r "\"koda/internal/parser\"" --include="*.go" .` → **no matches**.

No rerouting required before deleting `internal/vm/`.

## Section B — Files shared between paths (LLVM-only going forward)

| Area | Notes |
|------|--------|
| `internal/parser/`, `internal/lexer/`, `internal/sema/` | Shared frontend; **keep**. |
| `internal/codegen/`, `internal/nativebuild/` | Path B only; **keep**. |
| `cmd/koda/main.go` | Previously returned `ErrRuntimeUnavailable` for `run` / stub APIs; **updated** to call the native pipeline (`koda/api`). |
| `api/run.go` | Previously stubbed VM; **updated** to build + exec native binary. |
| `api/build_host.go` | Already native-only; **unchanged**. |
| `api/diagnose.go` | Parse + `sema.PrepareNativeBundle` only; **unchanged** (no VM). |
| `koda-ide/app.go` | Called `RunVM` / `SetVMPrintHook`; **updated** to `RunWithWriters`. |

## Section C — Capabilities the historical bytecode VM had that LLVM path may still lack

The repository **does not ship** a bytecode opcode table or `internal/vm` interpreter in this tree—only the stub above. Older documentation (e.g. `docs/architecture.md` before rewrite) described a full VM with `OpTailCall`, `gc()`, etc.

For each item below, **open a tracking issue** if not already implemented on the LLVM + C runtime path:

| Capability | LLVM / runtime check |
|------------|-------------------------|
| **Tail calls (`OpTailCall`)** | LLVM path does not guaranteed emit `musttail` / tail-call IR for recursive patterns; file issue if TCO is required for parity. |
| **`gc()` native** | C runtime has internal GC; exposing `gc()` to user `.koda` may be missing in codegen; verify against `koda_runtime.c` / builtins wiring. |
| **`setTimeout` / `setInterval`** | Not part of C runtime in tree; file issue if required. |
| **REPL** | `cmd/koda` has no REPL in current `main.go`; if a REPL existed on VM, parity is “LLVM compile snippet” or static-only—file issue. |
| **`koda disasm` (bytecode)** | Replaced with **LLVM IR** text from codegen (`PrepareNativeBundle` + `EmitLLVMIR`). File issue if bytecode-style disasm is still needed. |
| **`--debug` / `--profile` / `--memprofile`** | Not present on current CLI; if they were VM-only, file issue for LLVM equivalents (e.g. IR dump, `KODA_DEBUG_IR`). |

See **`docs/GITHUB_ISSUES_VM_RETIREMENT.md`** for copy-paste issue titles and bodies.
