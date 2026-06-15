#!/usr/bin/env bash
# Portable launchers for an extracted Koda SDK folder.
# Usage: write-portable-launchers.sh <sdk-dir> <platform>
#   platform: windows | unix | macos
set -euo pipefail

SDK_DIR="${1:?sdk directory}"
PLATFORM="${2:?platform: windows, unix, or macos}"
SDK_DIR="$(cd "$SDK_DIR" && pwd)"
LAUNCHERS="$(cd "$(dirname "$0")/launchers" && pwd)"

case "$PLATFORM" in
  windows)
    cp "$LAUNCHERS/start-koda-studio.bat" "$SDK_DIR/Start Koda Studio.bat"
    cp "$LAUNCHERS/open-pong-example.bat" "$SDK_DIR/Open Pong Example.bat"
    cp "$LAUNCHERS/check-sdk.bat" "$SDK_DIR/Check SDK (doctor).bat"
    ;;
  unix)
    cp "$LAUNCHERS/start-koda-studio.sh" "$SDK_DIR/start-koda-studio.sh"
    cp "$LAUNCHERS/open-pong-example.sh" "$SDK_DIR/open-pong-example.sh"
    cp "$LAUNCHERS/check-sdk.sh" "$SDK_DIR/check-sdk.sh"
    cp "$LAUNCHERS/koda-studio.desktop" "$SDK_DIR/koda-studio.desktop"
    chmod +x "$SDK_DIR/start-koda-studio.sh" "$SDK_DIR/open-pong-example.sh" "$SDK_DIR/check-sdk.sh"
    ;;
  macos)
    cp "$LAUNCHERS/start-koda-studio.command" "$SDK_DIR/Start Koda Studio.command"
    cp "$LAUNCHERS/open-pong-example.command" "$SDK_DIR/Open Pong Example.command"
    cp "$LAUNCHERS/check-sdk.command" "$SDK_DIR/Check SDK (doctor).command"
    chmod +x "$SDK_DIR/Start Koda Studio.command" "$SDK_DIR/Open Pong Example.command" "$SDK_DIR/Check SDK (doctor).command"
    ;;
  *)
    echo "unknown platform: $PLATFORM" >&2
    exit 1
    ;;
esac

echo "Wrote $PLATFORM launchers in $SDK_DIR"
