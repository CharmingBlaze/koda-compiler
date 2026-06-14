#!/usr/bin/env bash
# PR CI: Raylib header + wrapgen + koda link smoke (Linux, macOS, Windows Git Bash / MSYS).
# Vendored stdlib/sys-include/raylib.h supplies the header; a matching native Raylib is
# required to link the smoke binary (apt / brew / official MinGW prebuild zip on Windows).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if [[ ! -f "runtime/libkoda_runtime.a" ]]; then
  echo "runtime/libkoda_runtime.a missing; run scripts/build-runtime.sh first"
  exit 1
fi

uname_s="$(uname -s)"
is_windows_msys=0
case "$uname_s" in
  MINGW* | MSYS* | CYGWIN*) is_windows_msys=1 ;;
esac

HDR=""
if [[ -f "$ROOT/stdlib/sys-include/raylib.h" ]]; then
  HDR="$ROOT/stdlib/sys-include/raylib.h"
  echo "==> Using vendored raylib header: $HDR"
else
  echo "==> No stdlib/sys-include/raylib.h — installing system Raylib where supported"
  case "$uname_s" in
    Linux)
      sudo apt-get update -qq
      sudo apt-get install -y libraylib-dev pkg-config
      for c in /usr/include/raylib.h /usr/local/include/raylib.h; do
        if [[ -f "$c" ]]; then HDR="$c"; break; fi
      done
      ;;
    Darwin)
      brew list raylib &>/dev/null || brew install raylib
      _rb="$(brew --prefix raylib 2>/dev/null || true)"
      if [[ -n "$_rb" && -f "$_rb/include/raylib.h" ]]; then
        HDR="$_rb/include/raylib.h"
      fi
      ;;
    *)
      echo "raylib.h not found for OS: $uname_s (install dev package or restore stdlib/sys-include/raylib.h)" >&2
      exit 1
      ;;
  esac
  if [[ -z "$HDR" ]]; then
    echo "raylib.h not found (add stdlib/sys-include/raylib.h or install Raylib dev)"
    exit 1
  fi
  echo "    header: $HDR"
fi

RAYLIB_STAGE_WIN=""

ensure_raylib_for_link() {
  if pkg-config --exists raylib 2>/dev/null; then
    echo "==> raylib pkg-config ok"
    return 0
  fi
  case "$uname_s" in
    Linux)
      echo "==> Installing libraylib-dev for link (Linux)"
      sudo apt-get update -qq
      sudo apt-get install -y libraylib-dev pkg-config
      ;;
    Darwin)
      echo "==> Installing raylib via Homebrew for link (macOS)"
      brew list raylib &>/dev/null || brew install raylib
      _rb="$(brew --prefix raylib)"
      _pc="$_rb/lib/pkgconfig"
      if [[ -d "$_pc" ]]; then
        export PKG_CONFIG_PATH="${_pc}${PKG_CONFIG_PATH:+:}${PKG_CONFIG_PATH:-}"
      fi
      ;;
    MINGW* | MSYS* | CYGWIN*)
      echo "==> Windows MSYS: stage official raylib MinGW prebuild for link"
      RAYLIB_VER="5.0"
      RAYLIB_URL="https://github.com/raysan5/raylib/releases/download/${RAYLIB_VER}"
      st="${RUNNER_TEMP:-${TMP:-/tmp}}/koda-raylib-stage-$$"
      mkdir -p "$st"
      curl -fsSL "${RAYLIB_URL}/raylib-${RAYLIB_VER}_win64_mingw-w64.zip" -o "$st/dl.zip"
      unzip -q "$st/dl.zip" -d "$st/out"
      inner=""
      for d in "$st/out"/*; do
        [[ -d "${d}/include" && -d "${d}/lib" ]] && inner="$d" && break
      done
      if [[ -z "$inner" ]]; then
        echo "unexpected raylib zip layout" >&2
        exit 1
      fi
      RAYLIB_STAGE_WIN="$inner"
      ;;
    *)
      echo "Unsupported OS for Raylib link: $uname_s" >&2
      exit 1
      ;;
  esac
}

ensure_raylib_for_link

OUT="${WRAP_OUT:-$ROOT/wrappers/raylib_ci_generated}"
rm -rf "$OUT"
mkdir -p "$OUT"

echo "==> kodawrap / wrapgen (this may take a minute)"
go run ./cmd/wrapgen -name raylib -headers "$HDR" -out "$OUT" -docs=false -build=false -v

if [[ ! -s "$OUT/raylib.koda" ]] || [[ ! -s "$OUT/wrapper.c" ]]; then
  echo "wrapgen did not produce raylib.koda / wrapper.c"
  exit 1
fi

for artifact in koda.json META.json docs/index.html; do
  if [[ ! -s "$OUT/$artifact" ]]; then
    echo "wrapgen missing package artifact: $OUT/$artifact"
    exit 1
  fi
done
echo "==> package artifacts ok (koda.json, META.json, docs/index.html)"

echo "==> Link smoke: koda build tests/raylib_wrapgen_smoke.koda"
export KODA_WRAPPERS="$OUT"
export KODA_NATIVE_SOURCES="$OUT/wrapper.c"
if pkg-config --exists raylib 2>/dev/null; then
  export KODA_LINKFLAGS
  KODA_LINKFLAGS="$(pkg-config --libs --cflags raylib)"
else
  case "$uname_s" in
    Darwin)
      _rb="$(brew --prefix raylib)"
      export KODA_LINKFLAGS="-I${_rb}/include -L${_rb}/lib -lraylib"
      ;;
    MINGW* | MSYS* | CYGWIN*)
      if [[ -z "$RAYLIB_STAGE_WIN" ]]; then
        echo "internal error: RAYLIB_STAGE_WIN empty" >&2
        exit 1
      fi
      inc="$RAYLIB_STAGE_WIN/include"
      lib="$RAYLIB_STAGE_WIN/lib"
      export KODA_LINKFLAGS="-I${inc} -L${lib} -lraylib -lopengl32 -lgdi32 -lwinmm"
      ;;
    *)
      export KODA_LINKFLAGS="-I/usr/include -lraylib"
      ;;
  esac
fi

KODA_BIN="${KODA_BIN:-}"
if [[ -z "$KODA_BIN" || ! -x "$KODA_BIN" ]]; then
  KODA_BIN="$ROOT/.ci_koda"
  go build -trimpath -o "$KODA_BIN" ./cmd/koda
fi

SMOKE_BASE="${WRAPGEN_SMOKE_OUT:-/tmp/raylib_wrapgen_smoke}"
SMOKE_BASE="${SMOKE_BASE%.exe}"
EXE=""
if [[ "$is_windows_msys" -eq 1 ]]; then
  EXE=.exe
fi
OUT_BIN="${SMOKE_BASE}${EXE}"

"$KODA_BIN" build "$ROOT/tests/raylib_wrapgen_smoke.koda" -o "$OUT_BIN"
test -x "$OUT_BIN"
echo "==> OK: linked $OUT_BIN"
