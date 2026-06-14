#!/usr/bin/env bash
# Build runtime/libkoda_runtime.a for local development (Linux / macOS).
# Windows maintainers: use scripts/build-runtime.ps1 (MinGW gcc + ar) or match CI (LLVM clang + llvm-ar).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if [[ "$(uname -s)" == "Darwin" ]]; then
  CC_BIN="$(xcrun --find clang)"
  AR_BIN="ar"
else
  CC_BIN="$(command -v clang-20 || command -v clang-19 || command -v clang-18 || command -v clang || command -v gcc)"
  AR_BIN="$(command -v llvm-ar-20 || command -v llvm-ar-19 || command -v llvm-ar-18 || command -v llvm-ar || command -v ar)"
fi

if [[ -z "${CC_BIN}" || ! -x "$CC_BIN" ]]; then
  echo "no C compiler found (install clang or gcc)" >&2
  exit 1
fi
if [[ -z "${AR_BIN}" || ! -x "$AR_BIN" ]]; then
  echo "no archiver found (install llvm-ar or ar)" >&2
  exit 1
fi

echo "Using CC: $CC_BIN"
echo "Using AR: $AR_BIN"

mkdir -p runtime/obj
"$CC_BIN" -c runtime/src/value.c        -O2 -std=c11 -Iruntime/src -o runtime/obj/value.o
"$CC_BIN" -c runtime/src/object.c       -O2 -std=c11 -Iruntime/src -o runtime/obj/object.o
"$CC_BIN" -c runtime/src/gc.c           -O2 -std=c11 -Iruntime/src -o runtime/obj/gc.o
"$CC_BIN" -c runtime/src/koda_runtime.c -O2 -std=c11 -Iruntime/src -o runtime/obj/koda_runtime.o
"$AR_BIN" rcs runtime/libkoda_runtime.a runtime/obj/value.o runtime/obj/object.o runtime/obj/gc.o runtime/obj/koda_runtime.o

echo "Wrote runtime/libkoda_runtime.a"
