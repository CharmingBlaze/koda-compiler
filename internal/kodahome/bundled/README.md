# Legacy: bundled LLVM layout

Release pipeline used to populate **`internal/kodahome/bundled/`** with `llc`, `lld`, and `libkoda_runtime.a`.

That layout is **replaced** by **`internal/embed/<GOOS>/<GOARCH>/`** (Clang + runtime, and `lld.exe` on Windows). See **`internal/embed/README.md`**.

This directory remains for older notes and tooling references only.
