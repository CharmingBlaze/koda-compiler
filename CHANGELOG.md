# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- **Documentation hub** — `docs/README.md` indexes all user guides; README positions Koda as a modern C alternative for games and apps.
- **[Coming from C](docs/guides/from-c.md)** — side-by-side C vs Koda, migration table, FFI and project workflow.
- **[Applications guide](docs/guides/applications.md)** — CLI tools, file I/O, JSON, shipping desktop apps.
- **Updated [game-dev.md](docs/guides/game-dev.md)** — project templates, correct Raylib shim API, `koda.json` workflow.
- **`koda new`** — templates: `hello`, `game` (text lunar lander), `graphics` (Raylib bouncing-ball demo with bundled shim).
- **Project manifest (`koda.json`)** — `entry`, `bundle.assets`, and `native.sources` / `native.linkflags` are applied automatically for `run`, `build`, `watch`, `check`, and `bundle` when you omit the entry file.

### Fixed

## [0.3.2] - 2026-05-13

Fixes **`.github/workflows/release.yml`** so tagged releases build on current GitHub runners: LLVM **clang/llc/llvm-ar** version fallbacks on Linux, resilient **`llc*.exe`** discovery on Windows (Chocolatey LLVM), and **`chmod +x`** before macOS smoke tests.

## [0.3.1] - 2026-05-13

Cross-platform CI (Windows + Linux + macOS), LF-safe workflows, offline SDK packaging helpers, and contributor docs for GitHub downloads.

### Added

- **CI — Ubuntu, macOS, Windows** — **`.github/workflows/ci.yml`** matrix (**`ubuntu-latest`**, **`macos-latest`**, **`windows-latest`**) runs the same checks on every push/PR: C runtime build, **`go vet` / `go test`**, **`koda fmt --check`**, **`koda check`**, native compile/run smoke (**`scripts/ci-native-smoke.sh`**), time-bounded GC + stress (**GNU `timeout`** on Unix; **`scripts/ci-gc-stress-timed.ps1`** on Windows), parser fuzz, and Raylib wrapgen + link (**`scripts/ci-wrapgen-raylib.sh`**, including MinGW Raylib fetch on Windows Git Bash).
- **Local C runtime (Unix)** — **`scripts/build-runtime.sh`** builds **`runtime/libkoda_runtime.a`** with the same compiler discovery as CI (clang-20 … gcc / llvm-ar … ar; Xcode **`clang`** on macOS).
- **Docs / GitHub** — README and getting-started emphasize the **SDK zip** and that **`koda` does not download LLVM or Raylib** at compile time (embedded toolchain unpacks locally only).
- **`scripts/assemble-offline-sdk.ps1`** — Windows maintainers can assemble the same **offline SDK folder** as GitHub **`koda-*-sdk-windows-amd64.zip`**. **`scripts/build-release.ps1`** supports **`-PackageSdk`** / **`-PackageSdkZip`** alongside **`-NoZip`**.

### Fixed

- **CI (Linux / macOS / Windows)** — workflow YAML and repository **`*.sh`** scripts are normalized to **LF** line endings and **`.gitattributes`** forces **`eol=lf`** for them. CRLF in **`run: |`** blocks previously made bash fail immediately on GitHub-hosted runners (e.g. **`$'\r': command not found`**).
- **`koda fmt ./...`** — **`_legacy/`** is skipped (like **`_wg_*`** trees) so **`koda fmt --check ./...`** passes when a legacy mirror exists in the tree.

## [0.3.0] - 2026-05-10

Third release: **`stdlib/vec3`**, parser per-file parse cache, GC control builtins, runtime and lexer hardening, codegen loop and method-lowering fixes, and expanded CI (stress, fuzz).

### Added

- **`stdlib/vec3.koda`** — 3D **`{x, y, z}`** helpers (**`create`**, **`add`**, **`sub`**, **`scale`**, **`dot`**, **`cross`**, **`length`**, **`normalize`**, …); **`lerp`** now uses **`math.lerp`** per component after namespace-call lowering was fixed.
- **`tests/vec3_test.koda`**, **`tests/incremental_gc_test.koda`** — stdlib **`vec3`** coverage; incremental GC (**`gcCollectIncremental`** + **`gcFrameStep`**) loop stability (Linux CI runs incremental test with **`koda run --no-opt`** for time-bounded safety).
- **Parser — per-file parse cache** — **`internal/parser/loader.go`** caches parsed programs by **absolute path** + **`mtime`** nanoseconds; overlays skip the cache; **`internal/parser/loader_parse_cache_test.go`** (**`TestParseCacheUsesMtime`**).
- **`docs/MASTER_PLAN.md`** — maintainer-facing roadmap to “100%” (language, GC, runtime, DX, compiler, stdlib, tooling) with a **progress matrix** aligned to the current tree.
- **GC builtins** — **`gcDisable()`**, **`gcEnable()`**, **`gcCollect()`** (alias of **`gc()`**) wired to **`koda_gc_disable` / `koda_gc_enable` / `koda_gc_collect`**; **`tests/gc_control_test.koda`**.
- **`tests/struct_typo_test.koda`** — **`koda check`** regression for invalid struct field access (CI expects failure).
- **`internal/codegen/TestWriteBarrierCoverage`** — static scan of **`koda_runtime.c`** so mutators that assign through **`values` / `elements` / `keys` / `upvalues`** keep **`gc_write_barrier`** in the same function (with a narrow exempt list for **`table_rehash_open`**).
- **Runtime hardening (C)** — **`koda_value_type_name`**, **`koda_type_error`**, **`koda_null_error`**; **`koda_get` / `koda_set`** use typed panics instead of silent/fprintf failures where applicable; **`koda_unbox_number`** panics on non-number; array **`[]=`** out-of-bounds panics (**`logical length` + `capacity`**); string character index out-of-range panics. **`object.c`** — shared **`validate_value_slot_count`** for array/table/struct allocations and closure upvalue buffers.
- **Tests / CI** — **`tests/array_oob_get.koda`**, **`tests/array_oob_set.koda`** (API integration in **`api/runtime_hardening_test.go`**; tests **`Chdir`** to repo root so **`nativebuild`** finds **`runtime/libkoda_runtime.a`**); **`tests/stress/stress_mixed_alloc.koda`** with bounded Linux CI step; **`tests/vec3_test.koda`**, **`tests/vec3_math_lerp_direct_test.koda`**, and **`tests/math_lerp_member_test.koda`** in native smoke; **`tests/incremental_gc_test.koda`** in GC soak (**`timeout 120s`**, **`--no-opt`**); **`internal/parser`** **`FuzzParse`** (Lexer + Parser must not panic) with **`-fuzztime=5s`** on CI.

### Fixed

- **Lexer** — Template literal (**`` ` ``**) handling no longer calls **`advance()`** twice after the opening backtick, which could skip content and **panic** at EOF on input **[`` ` ``]** alone (**`fuzz` / `testdata`** regression).
- **Lexer** — **`string()`** treats a **`\`** immediately before EOF as an unterminated escape (no second **`advance()`** past the buffer).
- **Codegen (loops)** — while/for/do-while blocks now use unique LLVM block labels per statement; this prevents control-flow aliasing when multiple top-level loops appear in one function.
- **Method lowering** — `.length()` now lowers directly to **`koda_len`** with arity checking (`0` args), avoiding fallback paths that could treat property values as callables.
- **Method lowering (`math.*`)** — namespace math calls now dispatch to dedicated runtime argv natives even when user code defines overlapping global names (e.g. **`lerp`**), fixing member-call crashes from ABI mismatches.
- **Runtime arrays** — **`koda_array_push`** now grows array capacity dynamically (with overflow checks) instead of silently dropping writes once `count == capacity`; this restores expected `push` semantics and fixes growth-sensitive stress tests.
- **Regression tests** — **`tests/array_push_growth_test.koda`** validates long `push` growth + **`length()`**; **`tests/multi_while_flow_test.koda`** guards multi-loop control flow. Linux CI now runs both in native smoke and runs **`tests/nursery_test.koda`** in GC soak (`--no-opt`, time-bounded).
- **Game-scale stress** — **`tests/stress/large_game_sim.koda`** simulates a larger entity/projectile workload over hundreds of frames with incremental GC stepping; Linux CI runs it in stress smoke (`--no-opt`, time-bounded).

## [0.2.0] - 2026-05-06

Second release: richer language surface, **`koda fmt`**, shadow-stack / GC fixes, **`koda watch`**, **`build --debug`**, and **`run --no-opt`**.

### Added

- **`koda run` / `koda native`** — optional **`--no-opt`** (same meaning as **`koda build`**): skips the LLVM IR optimisation pass and uses **`-O0`** for **`llc`/Clang**, which avoids some **Clang-on-IR** optimiser crashes on large programs (e.g. **`tests/loops.koda`** on certain Windows LLVM builds). **`api.RunWithBuildOptions`**, **`api.RunWithWritersOpts`**, and type alias **`api.BuildOptions`** (`[nativebuild.BuildOptions]`).
- **`koda build --debug`** — emits native debug symbols (`-g`) and implies **`--no-opt`** for better stepping/stack traces in external debuggers. This flows through **`nativebuild.BuildOptions.Debug`** to both **`llc`** and **`clang`** paths.
- **`koda watch`** — watches `.koda` files under the entry file directory, rebuilds on change, and reruns the compiled temp executable (kills prior run before restart). Supports **`--no-opt`**.
- **`defer` statement** — defers a call (or any expression) until the enclosing **`func`** / closure / **`user_main`** exits: **LIFO** vs other defers, after the **`return`** value is computed, immediately before call-trace pop and shadow-stack pop. Lexer **`TokenDefer`**, AST **`DeferStmt`**, sema, shadow walk, formatter, and LLVM codegen. **`tests/defer_test.koda`** covers LIFO order and defers before an early **`return`**.
- **`??=` (nullish assignment)** — lexer/parser/formatter support; LLVM lowers to nil-check + conditional store (identifier and index targets). **`a..b` range expressions** are now parsed as **`RangeExpr`** (the lexer had **`..`** but the Pratt table did not). **`for (let i of lo..hi)`** with non-literal bounds uses a counted **`i64`** loop and **no** range object allocation.
- **`koda_shadow_stack_high_water()`** — C API returning max shadow-stack depth since **`koda_runtime_init`** (for diagnostics; high-water is tracked unconditionally).
- **`stdlib/vec2.koda`** — small **`{x, y}`** helpers (**`add`**, **`sub`**, **`scale`**, **`dot`**, **`length`**, **`normalize`**, etc.) on top of existing **`sqrt`**.
- **`stdlib/math.koda`** — re-exports **`degrees`**, **`radians`**, **`wrap`**, and **`approach`** from the injected **`math`** prelude; **`deg`** / **`rad`** remain as aliases.
- **Codegen — integer literal `+` / `-` / `*`** — when **both** operands are integer **literals**, the result is **constant-folded** at compile time (**`koda_box_number` only**, no **`koda_unbox_number`** on that subexpression). Applies per infix node (e.g. **`(1+2)+3`** folds in steps).
- **Codegen — unary `-` on integer literals** — folded to a single **`koda_box_number`** when the value round-trips through **`float64`** (same safety rule as literal infix fold).
- **Parser / codegen — unary `+`** — **`+expr`** is accepted; on integer-only literal trees it folds like unary **`-`**; otherwise it is a no-op on the emitted operand.
- **Codegen — nested integer literal folds** — **`compileTimeInt64`** folds trees of **`+` / `-` / `*`**, unary **`-`**, and parentheses over integer literals only (overflow / precision loss → runtime path).
- **`tests/assert_string_eq.koda`** — regression for **`assert`** with **`==`** on string literals and variables (interning / **`values_equal`**).
- **`tests/unary_plus_fold.koda`** — native smoke for unary **`+`** (literal fold and identity on a variable).
- **`tests/gc_shadow_multi_return_pop.koda`** — regression for shadow-stack **`koda_pop_frame`** on every **`ret`** path (multiple returns after **`if`**).
- **Tests / CI** — **`api/diagnose_test.go`** covers load + **`sema.PrepareNativeBundle`** (same pipeline as **`koda check`**): valid program, undefined name, overlay override, and aggregated multi-error text. Linux CI expects **`koda check tests/undefined_var.koda`** and **`tests/sema_errors.koda`** to fail; native-smoke **`koda build`** covers several **`tests/*.koda`** programs. **`internal/sema/arity_test.go`** and **`TestSemaCallArity*`** cover **call arity** rules.
- **Language surface (native / LLVM)** — template literals with **`${}`**; unary **`typeof`**; **`delete obj.key`**; array literal spread **`[...a]`**; **`??`** and **`?.`**; **`let { x, y } = obj`** destructuring; **`matches(haystack, pattern)`** (substring / **`strstr`**, not full regex); string methods (**`split`**, **`join`** on arrays, **`trim`**, **`replace`** / **`replaceAll`**, **`indexOf`**, **`toUpper`** / **`tolower`**, **`slice`**, **`startsWith`** / **`endsWith`**); array methods (**`map`**, **`filter`**, **`reduce`**, **`find`**, **`slice`**, **`sort`**, **`reverse`**, **`includes`**, **`concat`**, **`join`**); prelude **`math`** object (**`math.floor`**, …) alongside existing math globals; runtime symbols **`koda_array_join`**, **`koda_object_remove`**, **`koda_matches`** (rebuild **`runtime/libkoda_runtime.a`** after pulling).
- **`koda fmt`** — canonical AST-based formatting (`internal/formatter`): 4-space indent, spacing around operators and after commas, `if (` / `while (` / `for (` style, `} else {` kept on one line when both branches are blocks, `import "…"` expressions, **`koda fmt --check`**, directory and **`./...`** expansion (skips `.git`, `.KODA_build`, `bin`, `node_modules`, `vendor`). Top-level **consecutive** expression statements stay compact where appropriate.
- **Classic `for (init; cond; step)`** — parsed and lowered alongside `while` / `for-in` / `for-of` (`internal/parser`, `internal/codegen`).
- **`for (let [k, v] of iterable)`** — destructuring **`for-of`** for **arrays** (numeric index + element) and **objects/tables** (insertion-order key + value); runtime helpers **`koda_forof_length`**, **`koda_forof_key_at`**, **`koda_forof_value_at`** (`runtime/src/koda_runtime.c`).
- **Release binaries (`-tags release`)** — embedded Clang + runtime under **`internal/embed/`** (Windows also bundles **`lld.exe`** for linking); **`scripts/build-release.ps1`** / **`scripts/build-release.sh`**; **`release.yml`** publishes **`kodawrap`** where applicable.
- **Runtime hardening (GC/debug):** dynamic global-slot root registration (`koda_register_global_slot` + generated top-level slot wiring), opt-in **`KODA_GC_DEBUG`** stats/checks (remembered-set overflow counter, shadow stack depth high-water mark, global slot stats, object-header sanity checks), and stress scenarios (`tests/gc_pressure_expr.koda`, `tests/globals_perf.koda`, `tests/gc_soak.koda`).

### Changed

- **`language.md`** / **`docs/language.md`** — precedence summary and unary operator list now spell out **`+`**, **`-`**, **`!`**, and **`typeof`** (binary **`+` / `-`** called out where both appear in the same line).
- **`README.md`** — links **`docs/handoff.md`** for the compiler/runtime layout and invariants.
- **`docs/handoff.md`** — rewritten for the **LLVM-only** pipeline (no VM / **`koda`**); points at **`docs/architecture.md`** and invariants for new contributors.
- **Lexer API** — `lexer.NewLexer` now takes `(source, file string)` so the source path is always threaded at construction time; `loader` passes each module’s absolute path. **Template literal** tokens (start, string segments, `${}` interp markers, close) now set `Token.File` like other tokens. **Embedded template expressions** are re-lexed with the surrounding file path so sema diagnostics inside `` `${...}` `` point at the real `.koda` file. Injected **math prelude** tokens use the synthetic path `<builtin:math-prelude>`.
- **Sema diagnostics** — analysis records **all** issues in a function body or expression tree (undefined names, invalid targets, etc.), not only the first. Two or more problems return a **`diagnostic.MultiError`** (with `Unwrap() []error` so `errors.As` still finds the first `*DiagnosticError`). **`koda build` / `FormatError`** print a short summary plus each error with Rust-style snippets. A single error is still returned as a plain `*DiagnosticError`.
- **Sema — call and argv method arity** — for callees resolved to a **`func`** declaration (or an **`// koda:extern`** **`func` / `let`**), argument count is checked against parameters: required slots (no default), optional defaults, **`...rest`** (unbounded tail), and exact count for native **`Arity`**. For **`recv.method(args)`** / **`recv["method"](args)`** with a **compile-time string** method name, known argv string/array helpers (**`split`**, **`trim`**, **`reduce`**, **`map`**, …) use the same counts as **`internal/codegen/methods.go`**; **`concat`** stays variadic (unchecked). **`print`** and other builtins are unchanged (not modeled as `FuncDecl`). Arbitrary user-defined **`obj.m()`** is still not modeled.
- **GC intern table** — interned strings are no longer treated as unconditional GC roots. After each mark phase, **`koda_sweep_intern_table`** compacts the table to entries whose `ObjString` was marked reachable from real roots, so dead strings can be collected and intern slots never dangle across **`gc_sweep`**. (Full collect: sweep runs after mark; minor collect: sweep runs after remembered set, before nursery demotion.) **`koda_allocate_string` / `koda_copy_string`** still intern so equal text shares one object when live.
- **String `==`** — `values_equal` short-circuits identical string object pointers before `memcmp`.
- **`koda check`** — now runs **`sema.PrepareNativeBundle`** after **`parser.LoadProgram`** (flatten, math prelude, full semantic analysis) so undefined names and other sema errors are reported **without** invoking LLVM or the native linker. Still prints `OK` on success.
- **`len(table)`** — returns **entry count** for object/table values (`koda_len`).
- **`for-of`** / **`for-in`** lowering — uses slot-ordered runtime iteration for arrays and tables (**`for-in`** binds **keys**; **`for-of`** binds **values**; **`tests/phase1_surface.koda`** uses **`for-of`** where element sums are intended).
- **Documentation** — loop-style guide (`docs/user_guide.md`), **`for-of`** / **`syntax.md`** examples, single-line vs multi-line braced bodies.
- **Language (case):** reserved words and ASCII identifiers are **case-insensitive** (identifiers and object property keys are normalized to lowercase in the AST). **`@` module** specifiers and **`#include`** are matched case-insensitively; **`// koda:…`** lines treat the `koda:` prefix and the **`extern`** keyword case-insensitively, while the **C symbol** on extern lines stays spelled exactly for the linker. **`main`** / **`Main`** / etc. still map to the native entry (`koda_user_main`) case-insensitively. **Semicolons between statements are still required.**
- **Branding and paths:** sources use the **`.koda`** extension; C runtime and LLVM symbols use the **`koda_`** prefix; static archive **`runtime/libkoda_runtime.a`**. **Wrapper generator:** canonical binary name **`kodawrap`** (`go build -o kodawrap ./cmd/wrapgen`); **`koda wrap`** prefers **`kodawrap`** next to **`koda`**, then **`wrapgen`**, then legacy **`kujiwrap`**.
- **CI hardening:** Linux CI runs a serialized, time-bounded GC soak pass (`gc_pressure_expr`, `globals_perf`, `gc_soak`) to catch long-run/rooting regressions.

### Fixed

- **`defer` + `return value`** — **`return`** now evaluates and spills the result **before** running deferred calls, matching Go-style ordering for side effects.
- **Shadow stack hard cap** — when the growable shadow stack hits **`KODA_SHADOW_STACK_MAX_CAPACITY`**, **`koda_panic_str`** now reports **`stack overflow — maximum recursion depth reached`** instead of a bare **`stack overflow`** (Tier 6 polish).
- **Sema — `for-in` / `for-of` loop heads** — loop bindings (`let i`, `let [k, v]`, etc.) are now defined in an inner scope for the body so uses like **`sum + i`** are not reported as undefined.
- **Codegen — dynamic `for-of` range** — lower and upper bounds are **unboxed to integer `i64`** for the counted loop; comparing raw NaN-boxed tags with **`icmp`** could hang or mis-compare.
- **GC mark (`OBJ_ARRAY` / `OBJ_TABLE`)** — documented the existing **NULL guard** on `elements` / `keys` / `values` during partial allocation (GC may run between header and buffer allocation).
- **Codegen — shadow-stack pop on multiple returns** — **`emitShadowPop`** no longer clears **`shadowPushed`** after the first emitted **`return`**, so every **`ret`** basic block still emits **`koda_pop_frame`** (previously, tail/merge returns skipped the pop, leaking shadow depth and faulting during teardown on Windows). **`tests/gc_shadow_multi_return_pop.koda`** and **`tests/gc_shadow_branch_test.koda`** exercise this path.
- **`kodahome`**: removed duplicate **`KODA_CLANG`** / **`KODA_LLC`** env branches in **`ClangWithSource`** / **`LLCWithSource`**.

## [0.1.0] - 2026-05-05

First self-contained **Koda** release: native LLVM pipeline, embedded toolchain option, GC stack, and Result-based errors.

### Added

- **`ok(value)`** / **`err(message)`** — Result objects `{ ok, value, error }` for recoverable errors (plain tables; no exceptions).
- **`panic(message)`** — unrecoverable error: message to stderr, optional call-stack trace, **`exit(1)`** (not catchable).
- **`assert(condition, message?)`** — debug assertions using truthy rules; failures use the panic path.
- **`readFile`** / **`writeFile`** — return **`ok(...)`** or **`err(...)`** (real I/O; errno messages on failure).
- **`parseJSON`** — returns **`ok(parsed)`** or **`err(message)`**; JSON subset (objects, arrays, strings, numbers, booleans, null).
- **Call stack tracking** — codegen emits call-frame push/pop into the C runtime so panics can print a Koda stack trace.
- **Self-contained release binary** (Windows, Linux, macOS): **`go build -tags release`** embeds LLVM **`llc`** / **`lld`** and **`libkoda_runtime.a`**; **`kodahome.FindToolchain`** extracts them on first use. GitHub Actions **`.github/workflows/release.yml`** publishes **`koda-windows-amd64.exe`**, **`koda-linux-amd64`**, **`koda-linux-arm64`**, **`koda-darwin-amd64`**, **`koda-darwin-arm64`**.
- **`koda build --no-opt`** — disables IR optimisation and forces **`llc -O0`**; **`nativebuild.BuildOptions`** and related API surface the same.
- **`codegen.OptimiseIR`** — stub when **`CGO_ENABLED=0`** or without **`-tags llvm14`**; optional **tinygo.org/x/go-llvm** **`default<O2>`** pipeline when CGo + LLVM dev libs are available.
- Builtin **`gcFrameStep(ms)`** — maps to **`koda_gc_frame_step`** for game-loop GC integration (optional budget parameter).
- **Nursery GC** and **shadow stack** — bump nursery, precise roots via shadow frames, write barrier (see runtime and sema packages).
- **Builtin registration** in codegen — global stdlib names (**`readFile`**, **`sin`**, **`ok`**, **`assert`**, etc.) resolve consistently for native emission.
- Lexer/parser: **`var`** is reserved; use **`let`** for declarations.
- GitHub Actions **`.github/workflows/ci.yml`** — **`go vet`**, **`go test`**, C runtime build, native smoke.
- **`scripts/build-runtime.sh`** and **`scripts/build-runtime.ps1`** — produce **`runtime/libkoda_runtime.a`**.
- **`tests/smoke_native.koda`**, **`tests/error_handling_test.koda`** — CI and local checks.
- **`internal/codegen`** tests that LLVM runtime symbols exist in the static archive.
- **`cmd/wrapgen`** parse regression test (**`testdata/minimal.h`**).

### Fixed

- **`null`** literal codegen now uses the correct NaN-boxed **`NIL_VAL`** pattern (e.g. **`type(null)`** in **`tests/hello.koda`**).
- Array literal **`[1, 2, 3]`** and index expressions **`arr[i]`** (including chained indexing) parse correctly (Pratt handlers for **`[`** and indexing).

### Changed

- Docs and scripts consistently reference **`runtime/libkoda_runtime.a`** and **`runtime/src/koda_runtime.c`**.
- **`runtime/Makefile`** writes **`runtime/libkoda_runtime.a`** at the path **`koda build`** expects.
- Root **`make test`** depends on **`make -C runtime`** so the archive exists before link tests.
- Documentation describes the native pipeline: LLVM IR → **llc** → object → **Clang** + **`runtime/libkoda_runtime.a`** (or embedded toolchain when built with **`release`** tag).
- LLVM **`declare`** names aligned with **`koda_runtime.c`** (including **`koda_print_val`**, **`koda_get`**, argv-style builtins).

### Removed

- Bytecode VM package (**`internal/vm/`**) — execution is the **LLVM native** pipeline only (**`koda run`**, **`koda build`**, **`koda bundle`**, **`api.Run`**).
- Standalone interpreter / **`Value = any`** Go runtime — values are NaN-boxed in C and **i64** in IR.
- Separate **`koda native`** as a distinct pipeline — **`koda run`** is the primary path (legacy alias may remain).

### Deprecated

- The large embed under **`internal/runtime/data/`** is not used by default **`koda build`** (see **`internal/runtime/embed.go`**).
