#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"
STUDIO="$ROOT/Koda Studio"
if [[ ! -x "$STUDIO" ]]; then
  echo "Koda Studio not found." >&2
  exit 1
fi
if [[ -d "$ROOT/examples/games/pong" ]]; then
  exec "$STUDIO" "$ROOT/examples/games/pong"
fi
exec "$STUDIO"
