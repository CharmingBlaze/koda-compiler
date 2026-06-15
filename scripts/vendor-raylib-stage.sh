#!/usr/bin/env bash
# Download official raylib 6.0 prebuilt SDKs and normalize to third_party/raylib_static/stage layout
# (include/, lib/) for Koda [vendoredRaylibStatic] + offline zips. Run from repo root in CI before SDK packaging.
#
# Usage: vendor-raylib-stage.sh <OUT_BASE>
#   OUT_BASE/windows-amd64/third_party/raylib_static/stage/...
#   OUT_BASE/linux-amd64/...
#   ...
#
# Linux ARM64: upstream does not publish raylib-6.0 prebuilds; we write NOTES-linux-arm64.txt only.
#
set -euo pipefail

OUT_BASE="${1:?OUT_BASE directory (e.g. .raylib-vendor)}"
RAYLIB_VER="6.0"
RAYLIB_URL="https://github.com/raysan5/raylib/releases/download/${RAYLIB_VER}"

mkdir -p "$OUT_BASE"

tmpdir_root() {
  mktemp -d "${TMPDIR:-/tmp}/koda-raylib-XXXXXX"
}

fetch_zip() {
  local name="$1"
  local dest="$2"
  local t
  t="$(tmpdir_root)"
  curl -fsSL "${RAYLIB_URL}/${name}" -o "$t/dl.zip" || {
    echo "vendor-raylib-stage: download failed: ${RAYLIB_URL}/${name} (network or GitHub rate limit)" >&2
    rm -rf "$t"
    exit 1
  }
  unzip -q "$t/dl.zip" -d "$t/out"
  rm -f "$t/dl.zip"
  local inner=""
  local d
  for d in "$t/out"/*; do
    [[ -d "${d}/include" && -d "${d}/lib" ]] && inner="$d" && break
  done
  if [[ -z "$inner" ]]; then
    echo "unexpected zip layout for $name (no include+lib dir)" >&2
    exit 1
  fi
  mkdir -p "$dest"
  cp -a "$inner/include" "$inner/lib" "$dest/"
  # Optional: ship LICENSE / CHANGELOG if present
  for x in LICENSE CHANGELOG README.md; do
    if [[ -f "$inner/$x" ]]; then
      cp -a "$inner/$x" "$dest/"
    fi
  done
  rm -rf "$t"
}

fetch_tgz() {
  local name="$1"
  local dest="$2"
  local t
  t="$(tmpdir_root)"
  curl -fsSL "${RAYLIB_URL}/${name}" -o "$t/dl.tar.gz" || {
    echo "vendor-raylib-stage: download failed: ${RAYLIB_URL}/${name} (network or GitHub rate limit)" >&2
    rm -rf "$t"
    exit 1
  }
  tar -xzf "$t/dl.tar.gz" -C "$t"
  rm -f "$t/dl.tar.gz"
  local inner=""
  local d
  for d in "$t"/*; do
    [[ -d "$d" ]] || continue
    [[ -d "${d}/include" && -d "${d}/lib" ]] && inner="$d" && break
  done
  if [[ -z "$inner" ]]; then
    echo "unexpected tar layout for $name (no include+lib dir)" >&2
    exit 1
  fi
  mkdir -p "$dest"
  cp -a "$inner/include" "$inner/lib" "$dest/"
  for x in LICENSE CHANGELOG README.md; do
    if [[ -f "$inner/$x" ]]; then
      cp -a "$inner/$x" "$dest/"
    fi
  done
  rm -rf "$t"
}

write_stage() {
  local slug="$1"
  local dest="$OUT_BASE/${slug}/third_party/raylib_static/stage"
  mkdir -p "$dest"
  case "$slug" in
    windows-amd64)
      fetch_zip "raylib-${RAYLIB_VER}_win64_mingw-w64.zip" "$dest"
      ;;
    linux-amd64)
      fetch_tgz "raylib-${RAYLIB_VER}_linux_amd64.tar.gz" "$dest"
      ;;
    linux-arm64)
      mkdir -p "$dest"
      cat >"$dest/README-Koda.txt" <<'EOF'
Raylib 6.0 prebuilt binaries are not published by upstream for Linux ARM64.

Options:
  • Install system raylib (e.g. Debian/Ubuntu: libraylib-dev) and set KODA_LINKFLAGS to your
    -I and -l flags, plus KODA_NATIVE_SOURCES to wrappers/raylib/wrapper.c
  • Build from source and point KODA_RAYLIB_STAGE at a local stage/ tree with include/ + lib/

The full Koda wrapper and headers remain in this SDK under wrappers/raylib/ and docs/guides/raylib.md.
EOF
      ;;
    darwin-amd64|darwin-arm64)
      fetch_tgz "raylib-${RAYLIB_VER}_macos.tar.gz" "$dest"
      ;;
    *)
      echo "unknown slug: $slug" >&2
      exit 1
      ;;
  esac
}

for slug in windows-amd64 linux-amd64 linux-arm64 darwin-amd64 darwin-arm64; do
  write_stage "$slug"
  stage="$OUT_BASE/${slug}/third_party/raylib_static/stage"
  case "$slug" in
    windows-amd64)
      [[ -f "$stage/include/raylib.h" ]] || { echo "missing raylib.h for $slug" >&2; exit 1; }
      [[ -f "$stage/lib/libraylib.a" ]] || { echo "missing libraylib.a for $slug" >&2; exit 1; }
      [[ -f "$stage/lib/raylib.dll" ]] || { echo "missing raylib.dll for $slug" >&2; exit 1; }
      ;;
    linux-amd64|darwin-amd64|darwin-arm64)
      [[ -f "$stage/include/raylib.h" ]] || { echo "missing raylib.h for $slug" >&2; exit 1; }
      [[ -f "$stage/lib/libraylib.a" ]] || { echo "missing libraylib.a for $slug" >&2; exit 1; }
      ;;
    linux-arm64)
      [[ -f "$stage/README-Koda.txt" ]] || { echo "missing README-Koda.txt for $slug" >&2; exit 1; }
      ;;
  esac
  echo "ok: $slug -> $OUT_BASE/$slug/third_party/raylib_static/stage"
done

echo "Raylib vendor complete under $OUT_BASE"
