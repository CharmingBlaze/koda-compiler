#!/usr/bin/env bash
# Semantic check for all canonical example projects (no native link — fast CI gate).
set -euo pipefail

KODA="${KODA_CI_BIN:?KODA_CI_BIN is required}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"

PROJECTS=(
  examples/koda-3d
  examples/cube3d
  examples/demo-3d
  examples/spinning-cube
  examples/raylib-3d-demo
  examples/crystal-plaza
  examples/hello_raylib_raw
  examples/hello-use-module
  examples/lunar-lander
  examples/games/pong
  examples/games/brick-breaker
  examples/games/fps-arena
  examples/games/koda64
  examples/games/mario64
  examples/games/mario64-hilltop
  examples/games/mario64-studio
)

echo "==> koda check all example projects (${#PROJECTS[@]})"
for proj in "${PROJECTS[@]}"; do
  echo "    $proj"
  (cd "$ROOT/$proj" && "$KODA" check src/main.koda)
done

echo "==> koda check dist-template samples"
"$KODA" check "$ROOT/dist-template/examples/hello.koda"
"$KODA" check "$ROOT/dist-template/examples/games/brick_breaker.koda"

echo "==> examples check OK"
