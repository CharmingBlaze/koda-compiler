This directory holds the gzip+tar archive embedded into koda and kodawrap.

Upstream ships a tiny placeholder so the repo builds without a full LLVM tree.
For a self-contained release binary:

  1. Build a tarball whose paths start with toolchain/ or llvm/ and include
     toolchain/bin/clang (or clang.exe), toolchain/bin/ld.lld (or ld.lld.exe on Windows),
     plus lib/clang/<ver>/include for portable headers. The native linker path passes
     -fuse-ld=lld when ld.lld is present next to clang. Optional: llvm/bin/llc (or llc.exe)
     for the legacy llc+link path. You may also add stdlib/*.koda at the archive root.

  2. Replace bundled_toolchain.tar.gz with that file, then: go build -o koda ./cmd/koda

First run extracts next to the executable when no bundled clang is present yet.
