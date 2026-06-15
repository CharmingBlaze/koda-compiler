#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")" && pwd)"
cd "$ROOT"
STUDIO="$ROOT/Koda Studio"
if [[ ! -x "$STUDIO" ]]; then
  echo "Koda Studio not found. Unzip the full Koda SDK (Linux) so 'Koda Studio' sits next to koda and stdlib/." >&2
  exit 1
fi
exec "$STUDIO" "$@"
