// Package nativebuild links Koda programs by running **llc** on LLVM IR text to produce an
// object file, then invoking **Clang** as the linker driver for that object plus the
// prebuilt C runtime archive (runtime/libkoda_runtime.a) and optional native glue sources.
//
// # Three levels of “Go talking to LLVM” (and where Koda sits)
//
// **1 — IR as the intermediate language (text or bitcode):** The parser builds an AST;
// the LLVM backend walks it and emits IR. Here that is a llir ir.Module
// serialized to LLVM assembly text, not hand-printed fmt.Sprintf templates.
//
// **2 — Deep in-process LLVM (Go bindings / cgo):** The official LLVM Go bindings talk
// to libLLVM for JIT, custom passes, etc. Koda deliberately does **not** link libLLVM
// into the compiler: you keep a portable `go build` with no cgo/libLLVM install burden.
// The “deep” integration is llir’s IR builder; optimization is mostly whatever Clang
// applies when invoked with -O3 (plus your usual LTO flags if you add them later).
//
// **3 — Executable and libraries (linking):** Clang invokes the platform linker by
// default. For a self-contained toolchain, ship **lld** next to clang and pass
// -fuse-ld=lld (see [UseLLD] and kodahome.HasBundledLLD). Extra objects and -l flags
// still go on the same Clang command line (KODA_NATIVE_SOURCES, KODA_LINKFLAGS).
//
// # End-to-end flow in this repo
//
// Frontend (Go): parse → bundle → (optional) sema for native.
//
// Middleware (Go + llir): AST → ir.Module string.
//
// Backend: write main.ll under .KODA_build, run llc to main.o, then Clang links main.o with
// runtime/libkoda_runtime.a, -I runtime/src, KODA_NATIVE_SOURCES, optional vendored Raylib
// (third_party/raylib_static/stage from `make raylib-lib`; disable with KODA_USE_VENDORED_RAYLIB=0),
// and KODA_LINKFLAGS.
// Set KODA_DEBUG_IR to also persist the IR on disk for inspection.
//
// Linker: the driver’s default linker, or LLD when enabled, resolves libc, your glue,
// and third-party .a/.lib/.dll as usual.
package nativebuild
