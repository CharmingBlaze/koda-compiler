#!/usr/bin/env bash
# Build self-contained Koda SDK archives for GitHub Releases.
# Requires: unzip's companion `zip`, bash, repo checkout at stdlib/, docs/, repo-root *.md docs.
#
# Usage: package-release-sdk.sh <VERSION> <ARTIFACTS_DIR> [OUT_DIR]
# Example (from repo root): bash scripts/package-release-sdk.sh "$GITHUB_REF_NAME" artifacts ./sdk-zips

set -euo pipefail

VERSION="${1:?first arg: version/tag e.g. v1.2.3}"
ART_ROOT="${2:?second arg: artifacts root (contains windows/, linux-amd64/, …)}"
OUT_DIR="${3:-./sdk-zips}"

if ! command -v zip >/dev/null 2>&1; then
  echo "zip(1) is required" >&2
  exit 1
fi

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$REPO_ROOT"

require_file() {
  local f="$1"
  if [[ ! -f "$f" ]]; then
    echo "package-release-sdk: missing required file (check release artifacts and paths): $f" >&2
    exit 1
  fi
}

require_raylib_stage() {
  local slug="$1"
  local stage_dir="$2"
  if [[ ! -d "$stage_dir" ]]; then
    echo "package-release-sdk: missing raylib stage for ${slug}: $stage_dir" >&2
    echo "run scripts/vendor-raylib-stage.sh and set RAYLIB_VENDOR_ROOT before packaging" >&2
    exit 1
  fi
  case "$slug" in
    windows-amd64)
      require_file "$stage_dir/include/raylib.h"
      require_file "$stage_dir/lib/libraylib.a"
      require_file "$stage_dir/lib/raylib.dll"
      ;;
    linux-amd64|darwin-amd64|darwin-arm64)
      require_file "$stage_dir/include/raylib.h"
      require_file "$stage_dir/lib/libraylib.a"
      ;;
    linux-arm64)
      require_file "$stage_dir/README-Koda.txt"
      ;;
    *)
      echo "package-release-sdk: unknown slug in raylib stage check: $slug" >&2
      exit 1
      ;;
  esac
}

# Copy Koda Studio + platform launchers when CI built the IDE for this slug.
bundle_studio() {
  local slug="$1"
  local root="$2"

  case "$slug" in
    windows-amd64)
      local src="$ART_ROOT_ABS/windows/koda-studio-windows-amd64.exe"
      if [[ ! -f "$src" ]]; then
        return 0
      fi
      cp -a "$src" "$root/Koda Studio.exe"
      bash "$REPO_ROOT/scripts/write-portable-launchers.sh" "$root" windows
      ;;
    linux-amd64)
      local src="$ART_ROOT_ABS/linux-amd64/koda-studio-linux-amd64"
      if [[ ! -f "$src" ]]; then
        return 0
      fi
      cp -a "$src" "$root/Koda Studio"
      chmod +x "$root/Koda Studio"
      bash "$REPO_ROOT/scripts/write-portable-launchers.sh" "$root" unix
      ;;
    linux-arm64)
      local src="$ART_ROOT_ABS/linux-arm64/koda-studio-linux-arm64"
      if [[ ! -f "$src" ]]; then
        return 0
      fi
      cp -a "$src" "$root/Koda Studio"
      chmod +x "$root/Koda Studio"
      bash "$REPO_ROOT/scripts/write-portable-launchers.sh" "$root" unix
      ;;
    darwin-amd64|darwin-arm64)
      local src="$ART_ROOT_ABS/macos/Koda Studio.app"
      if [[ ! -d "$src" ]]; then
        return 0
      fi
      cp -a "$src" "$root/Koda Studio.app"
      bash "$REPO_ROOT/scripts/write-portable-launchers.sh" "$root" macos
      ;;
    *)
      return 0
      ;;
  esac
}

zip_sdk() {
  local slug="$1"          # windows-amd64
  local fj_src="$2"        # source path to compiler binary
  local fw_src="$3"        # source path to kodawrap binary
  local fj_out="$4"        # basename inside archive (koda or koda.exe)
  local fw_out="$5"        # kodawrap or kodawrap.exe

  require_file "$fj_src"
  require_file "$fw_src"
  require_dir "$REPO_ROOT/stdlib"
  require_dir "$REPO_ROOT/docs"

  local root_name="koda-${VERSION}-${slug}"
  local stage
  stage="$(mktemp -d)"

  mkdir -p "$stage/$root_name"
  cp -a "$fj_src" "$stage/$root_name/$fj_out"
  cp -a "$fw_src" "$stage/$root_name/$fw_out"
  cp -a "$REPO_ROOT/stdlib" "$stage/$root_name/stdlib"
  cp -a "$REPO_ROOT/docs" "$stage/$root_name/docs"

  mkdir -p "$stage/$root_name/scripts"
  cp -a "$REPO_ROOT/scripts/install-koda.ps1" "$REPO_ROOT/scripts/install-koda.sh" "$stage/$root_name/scripts/"

  # All top-level documentation markdown (CHANGELOG, CONTRIBUTING, language.md, START_HERE.md, …).
  shopt -s nullglob
  local md
  for md in "$REPO_ROOT"/*.md; do
    cp -a "$md" "$stage/$root_name/"
  done
  shopt -u nullglob

  # Wrapper corpus + markdown (e.g. raylib README, api_reference) not duplicated under docs/.
  if [[ -d "$REPO_ROOT/wrappers" ]]; then
    cp -a "$REPO_ROOT/wrappers" "$stage/$root_name/wrappers"
  fi

  # Official raylib 5.0 prebuilts (see scripts/vendor-raylib-stage.sh). RAYLIB_VENDOR_ROOT defaults
  # to REPO_ROOT/.raylib-vendor when that directory exists; CI sets it explicitly.
  local rv_root="${RAYLIB_VENDOR_ROOT:-}"
  if [[ -z "$rv_root" ]] && [[ -d "$REPO_ROOT/.raylib-vendor" ]]; then
    rv_root="$REPO_ROOT/.raylib-vendor"
  fi
  local raylib_stage="$rv_root/$slug/third_party/raylib_static/stage"
  if [[ -z "$rv_root" ]]; then
    echo "package-release-sdk: RAYLIB_VENDOR_ROOT is not set and .raylib-vendor is missing; refusing to build incomplete SDKs" >&2
    exit 1
  fi
  require_raylib_stage "$slug" "$raylib_stage"
  mkdir -p "$stage/$root_name/third_party/raylib_static"
  cp -a "$raylib_stage" "$stage/$root_name/third_party/raylib_static/stage"

  # Small runnable corpus (optional; helps offline users).
  if [[ -d "$REPO_ROOT/examples" ]]; then
    cp -a "$REPO_ROOT/examples" "$stage/$root_name/examples"
  fi

  bundle_studio "$slug" "$stage/$root_name"

  cat >"$stage/$root_name/SDK_README.txt" <<EOF
Koda SDK ${VERSION} (${slug})
==============================

The beginner-friendly replacement for C/C++ — native games and apps from one zip.
No Go. No Python. No LLVM install required.

Platforms: Windows, Linux, and macOS each have dedicated release zips (this archive is ${slug}).
Everything offline: compiler, kodawrap, stdlib, docs, examples — unzip and run koda doctor.

Each SDK zip includes Koda Studio (IDE) and one-click launchers:
  • Windows — Start Koda Studio.bat
  • Linux — ./start-koda-studio.sh (or koda-studio.desktop)
  • macOS — Start Koda Studio.command (double-click in Finder)

The compiler never fetches LLVM, Raylib, or other dependencies from the network. Release builds embed Clang + llc + the Koda runtime; they unpack to a local temp directory on first use only. Raylib (where included in this zip) lives under third_party/raylib_static/stage/.

Quick start
-----------
  Read START_HERE.md, then:

  • Windows:
      Double-click Start Koda Studio.bat
      — or — .\koda.exe doctor && .\koda.exe new bounce --template graphics

  • Linux:
      chmod +x koda kodawrap start-koda-studio.sh
      ./start-koda-studio.sh
      — or — ./koda doctor && ./koda new bounce --template graphics

  • macOS:
      Double-click Start Koda Studio.command
      — or — chmod +x koda kodawrap && ./koda doctor

  Optional PATH: scripts/install-koda.ps1 (Windows) or scripts/install-koda.sh (Unix)

Layout
------
  $(printf '%s' "$fj_out")     — Koda compiler (${slug})
  $(printf '%s' "$fw_out")     — kodawrap (C header → .koda + wrapper.c)
  stdlib/         — shipped .koda modules (@ imports, #includes)
  docs/           — full documentation tree (guides, commands, language notes, …)
  *.md (root)     — all repo-root markdown (language ref, README, CHANGELOG, …)
  wrappers/       — full raylib wrapper (raylib.koda + wrapper.c + docs/HTML)
  third_party/raylib_static/stage/ — raylib 5.0 headers + libs (+ raylib.dll on Windows) when vendored

Use
---
  Keep this folder together so stdlib/ sits next to koda (or run from this directory):

  • Windows:
      .\koda.exe version
      .\koda.exe run examples\hello.koda

  • Linux / macOS:
      chmod +x koda kodawrap
      ./koda version
      ./koda run examples/hello.koda

Wrapper tool
------------
  ./koda wrap --help          (same as ./kodawrap … when both sit here)
  kodawrap uses the same bundled C toolchain as koda after unpack — no system LLVM/Clang required.

Embedded toolchain (release builds): Clang + llc + runtime are bundled inside $(printf '%s' "$fj_out") — no separate LLVM download for Koda itself.
Raylib games: see docs/guides/raylib.md — vendored static lib is used automatically when
third_party/raylib_static/stage exists next to this SDK (or set KODA_RAYLIB_STAGE).
Copy raylib.dll next to your .exe on Windows when using the dynamic build from lib/.

Other native libs: see docs/wrappers.md.

EOF

  mkdir -p "$OUT_DIR"
  local out_zip
  out_zip="$(cd "$OUT_DIR" && pwd)/koda-${VERSION}-sdk-${slug}.zip"
  rm -f "$out_zip"
  if [[ "$fj_out" != *.exe ]]; then
    chmod +x "$stage/$root_name/$fj_out" "$stage/$root_name/$fw_out"
  fi
  (cd "$stage" && zip -rq "$out_zip" "$root_name")
  rm -rf "$stage"
  echo "wrote $out_zip"
}

require_dir() {
  local d="$1"
  if [[ ! -d "$d" ]]; then
    echo "package-release-sdk: missing required directory (repo checkout incomplete?): $d" >&2
    exit 1
  fi
}

ART_ROOT_ABS="$(cd "$ART_ROOT" && pwd)"

zip_sdk "windows-amd64" \
  "$ART_ROOT_ABS/windows/koda-windows-amd64.exe" \
  "$ART_ROOT_ABS/windows/kodawrap-windows-amd64.exe" \
  "koda.exe" \
  "kodawrap.exe"

zip_sdk "linux-amd64" \
  "$ART_ROOT_ABS/linux-amd64/koda-linux-amd64" \
  "$ART_ROOT_ABS/linux-amd64/kodawrap-linux-amd64" \
  "koda" \
  "kodawrap"

zip_sdk "linux-arm64" \
  "$ART_ROOT_ABS/linux-arm64/koda-linux-arm64" \
  "$ART_ROOT_ABS/linux-arm64/kodawrap-linux-arm64" \
  "koda" \
  "kodawrap"

zip_sdk "darwin-amd64" \
  "$ART_ROOT_ABS/macos/koda-darwin-amd64" \
  "$ART_ROOT_ABS/macos/kodawrap-darwin-amd64" \
  "koda" \
  "kodawrap"

zip_sdk "darwin-arm64" \
  "$ART_ROOT_ABS/macos/koda-darwin-arm64" \
  "$ART_ROOT_ABS/macos/kodawrap-darwin-arm64" \
  "koda" \
  "kodawrap"

echo "All SDK zips OK under $OUT_DIR"
