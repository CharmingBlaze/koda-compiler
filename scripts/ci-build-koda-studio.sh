#!/usr/bin/env bash
# Build Koda Studio for CI or local release packaging.
# Usage: ci-build-koda-studio.sh <output-path>
#   output-path — file (Linux/Windows exe) or .app directory (macOS)
set -euo pipefail

OUT="${1:?usage: ci-build-koda-studio.sh <output-path>}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
IDE="$ROOT/koda-ide"

cd "$IDE/frontend"
npm ci
npm run build
cd "$IDE"

if ! command -v wails >/dev/null 2>&1; then
  go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0
fi

WAILS_ARGS=(-s -m -nopackage)
if [[ "${WAILS_TAGS:-}" != "" ]]; then
  WAILS_ARGS+=(-tags "$WAILS_TAGS")
fi

wails build "${WAILS_ARGS[@]}"

BIN_DIR="$IDE/build/bin"
if [[ -d "$BIN_DIR/Koda Studio.app" ]]; then
  rm -rf "$OUT"
  cp -a "$BIN_DIR/Koda Studio.app" "$OUT"
elif [[ -f "$BIN_DIR/Koda Studio" ]]; then
  cp -a "$BIN_DIR/Koda Studio" "$OUT"
  chmod +x "$OUT"
elif [[ -f "$BIN_DIR/Koda Studio.exe" ]]; then
  cp -a "$BIN_DIR/Koda Studio.exe" "$OUT"
elif [[ -f "$BIN_DIR/koda-ide.exe" ]]; then
  cp -a "$BIN_DIR/koda-ide.exe" "$OUT"
else
  echo "Koda Studio build output not found under $BIN_DIR" >&2
  ls -la "$BIN_DIR" 2>/dev/null || true
  exit 1
fi

echo "Koda Studio -> $OUT"
