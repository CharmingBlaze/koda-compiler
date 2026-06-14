# Suggested GitHub issues after bytecode VM removal

Create these (or merge with existing issues) so no capability disappears without tracking.

---

## LLVM path missing: guaranteed tail-call optimization

**Body (draft):** The legacy bytecode VM documented `OpTailCall` for tail-position calls. The LLVM pipeline should emit tail calls where possible (`musttail` / sibling call optimization) so deep recursion does not overflow the stack. Scope: identify tail positions in codegen and emit appropriate LLVM IR; add tests (e.g. recursive fib-style).

---

## LLVM path missing: user-callable `gc()` builtin

**Body (draft):** If the language exposes `gc()` to scripts, wire it to the C runtime’s collector (or document that manual GC is not exposed). Verify against standard tests such as `tests/gc_test.koda` if present.

---

## LLVM path missing: timers (`setTimeout` / `setInterval`)

**Body (draft):** If these were VM-era APIs, either implement in C runtime + codegen or explicitly document as out of scope for the native host.

---

## LLVM path missing: REPL / evaluate-expression mode

**Body (draft):** If tooling relied on a bytecode REPL, define whether the replacement is: compile-and-run snippets via native pipeline, or a parse-only / check-only mode.

---

## LLVM path missing: profiling flags parity (`--profile`, `--memprofile`)

**Body (draft):** If the old CLI exposed profiling tied to the VM, specify LLVM/native equivalents (e.g. `perf`, `/Fd`, sampling) or document workflow.

---

## Documentation cleanup: dual-path references

**Body (draft):** Search docs for “bytecode VM”, “koda run (VM)”, and `internal/parser` paths; align with single pipeline (lexer → parser → sema → codegen → `.ll` → clang).
