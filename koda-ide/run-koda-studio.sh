#!/usr/bin/env bash
# Build and run Koda Studio from source (macOS / Linux). Windows: run-koda-studio.ps1
set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
REPO="$(cd "$ROOT/.." && pwd)"

export GOCACHE="${GOCACHE:-$ROOT/.gocache}"
export GOMODCACHE="${GOMODCACHE:-$ROOT/.gomodcache}"
export GOPATH="${GOPATH:-$ROOT/.gopath}"

if ! command -v wails >/dev/null 2>&1; then
  echo "Installing Wails CLI…"
  go install github.com/wailsapp/wails/v2/cmd/wails@v2.12.0
fi

if [[ ! -d "$ROOT/frontend/node_modules" ]]; then
  (cd "$ROOT/frontend" && npm install)
fi

(cd "$ROOT/frontend" && npm run build)

(cd "$ROOT" && wails build -s -m -nopackage)

if [[ -d "$ROOT/build/bin/Koda Studio.app" ]]; then
  STUDIO="$ROOT/build/bin/Koda Studio.app"
elif [[ -x "$ROOT/build/bin/Koda Studio" ]]; then
  STUDIO="$ROOT/build/bin/Koda Studio"
else
  echo "Studio binary not found under build/bin" >&2
  exit 1
fi

if [[ -f "$REPO/stdlib/math.koda" ]]; then
  export KODA_HOME="$REPO"
  echo "Using SDK at: $REPO"
fi

echo "Starting Koda Studio…"
if [[ -d "$STUDIO" ]]; then
  if [[ $# -gt 0 ]]; then
    open "$STUDIO" --args "$@"
  else
    open "$STUDIO"
  fi
else
  exec "$STUDIO" "$@"
fi
