#!/usr/bin/env bash
# Assemble a portable "zero-install" directory layout next to a built koda binary.
# Does not build or embed the LLVM toolchain tarball (maintainers handle that separately).
#
# Usage:
#   ./scripts/package-zero-install.sh /path/to/koda.exe [output-dir]
# Default output-dir: ./dist/zero-install
set -euo pipefail

KODA_EXE="${1:?usage: $0 <path-to-koda-binary> [output-dir]}"
OUT="${2:-dist/zero-install}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"

mkdir -p "$OUT"
cp -f "$KODA_EXE" "$OUT/"

for d in stdlib wrappers runtime; do
  if [[ -d "$ROOT/$d" ]]; then
    rm -rf "$OUT/$d"
    cp -R "$ROOT/$d" "$OUT/"
  fi
done

# Ship static runtime for host packagers who symlink it into the bundle.
if [[ -f "$ROOT/runtime/libkoda_runtime.a" ]]; then
  mkdir -p "$OUT/runtime"
  cp -f "$ROOT/runtime/libkoda_runtime.a" "$OUT/runtime/"
fi

cat >"$OUT/README_ZERO_INSTALL.txt" <<'EOF'
Zero-install layout
-------------------
- koda (or koda.exe): compiler CLI
- stdlib/:           @module resolution (e.g. @array)
- wrappers/:        optional pre-generated bindings (e.g. @raylib)
- runtime/:          libkoda_runtime.a for native links

Run: koda doctor
Build: koda build your.koda -o your
EOF

echo "Wrote layout to $OUT"
ls -la "$OUT"
