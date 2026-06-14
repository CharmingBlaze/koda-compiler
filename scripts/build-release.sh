#!/bin/bash
# scripts/build-release.sh — local release binary for testing (requires clang + llvm-ar).

set -e

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then ARCH="amd64"; fi
if [ "$ARCH" = "aarch64" ]; then ARCH="arm64"; fi

echo "Building Koda release binary for $OS/$ARCH..."

echo "Building C runtime..."
mkdir -p runtime/obj_release

if CLANG=$(command -v clang-18 2>/dev/null); then
  :
elif CLANG=$(command -v clang 2>/dev/null); then
  :
else
  echo "clang not found" >&2
  exit 1
fi

"$CLANG" -c runtime/src/value.c        -O2 -std=c11 -Iruntime/src -o runtime/obj_release/value.o
"$CLANG" -c runtime/src/object.c       -O2 -std=c11 -Iruntime/src -o runtime/obj_release/object.o
"$CLANG" -c runtime/src/gc.c           -O2 -std=c11 -Iruntime/src -o runtime/obj_release/gc.o
"$CLANG" -c runtime/src/koda_runtime.c -O2 -std=c11 -Iruntime/src -o runtime/obj_release/koda_runtime.o

if AR=$(command -v llvm-ar-18 2>/dev/null); then
  :
elif AR=$(command -v llvm-ar 2>/dev/null); then
  :
elif AR=$(command -v ar 2>/dev/null); then
  :
else
  echo "llvm-ar/ar not found" >&2
  exit 1
fi

"$AR" rcs runtime/libkoda_runtime.a runtime/obj_release/*.o

echo "Populating embed directory..."
mkdir -p "internal/embed/$OS/$ARCH"
cp "$CLANG"                    "internal/embed/$OS/$ARCH/clang"
cp runtime/libkoda_runtime.a "internal/embed/$OS/$ARCH/"

echo "Building koda..."
go build -trimpath -tags release -ldflags="-s -w" -o koda-release ./cmd/koda

echo ""
echo "Done. Test with:"
echo "  ./koda-release build tests/hello.koda -o hello"
echo "  ./hello"
