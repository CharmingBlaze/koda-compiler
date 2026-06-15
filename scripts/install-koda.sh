#!/usr/bin/env bash
# Install Koda SDK binaries into ~/.local/bin (user PATH).
# Run from the root of an extracted SDK zip (where koda lives).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
if [[ ! -x "$ROOT/koda" && ! -f "$ROOT/koda.exe" ]]; then
  if [[ -x "$ROOT/koda" || -f "$ROOT/koda.exe" ]]; then
    :
  else
    ROOT="$(pwd)"
  fi
fi
if [[ ! -x "$ROOT/koda" && ! -f "$ROOT/koda.exe" ]]; then
  echo "Run from your Koda SDK folder (must contain koda)." >&2
  exit 1
fi

DEST="${HOME}/.local/bin"
mkdir -p "$DEST"

install -m 755 "$ROOT/koda" "$DEST/koda" 2>/dev/null || cp "$ROOT/koda" "$DEST/koda" && chmod +x "$DEST/koda"
if [[ -x "$ROOT/kodawrap" ]]; then
  install -m 755 "$ROOT/kodawrap" "$DEST/kodawrap" 2>/dev/null || cp "$ROOT/kodawrap" "$DEST/kodawrap" && chmod +x "$DEST/kodawrap"
fi

echo ""
echo "Installed koda to: $DEST/koda"
echo ""
echo "Ensure ~/.local/bin is on your PATH (add to ~/.bashrc or ~/.zshrc):"
echo '  export PATH="$HOME/.local/bin:$PATH"'
echo ""
echo "Keep stdlib/ next to the SDK or set KODA_HOME:"
echo "  export KODA_HOME=\"$ROOT\""
echo ""
echo "Then run:"
echo "  koda doctor"
echo "  koda new mygame --template graphics"
echo ""
