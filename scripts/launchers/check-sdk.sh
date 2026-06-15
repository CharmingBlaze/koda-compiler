#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"
if [[ ! -x "$ROOT/koda" ]]; then
  echo "koda not found in this folder." >&2
  exit 1
fi
exec "$ROOT/koda" doctor
