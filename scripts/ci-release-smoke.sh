#!/usr/bin/env bash
# Pre-publish smoke for release.yml — tier-1 language + runtime regressions.
set -euo pipefail

KODA="${1:?usage: ci-release-smoke.sh <path-to-koda>}"
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

EXE=""
case "$(uname -s)" in
  MINGW* | MSYS* | CYGWIN*) EXE=.exe ;;
esac

OUT="${CI_ARTIFACT_DIR:-/tmp/koda-release-smoke}"
mkdir -p "$OUT"

echo "==> release smoke: KODA=$KODA"

"$KODA" version

"$KODA" build tests/hello.koda -o "$OUT/hello$EXE"
"$OUT/hello$EXE"

"$KODA" build tests/gc_shadow_multi_return_pop.koda -o "$OUT/gc_shadow$EXE"
"$OUT/gc_shadow$EXE"

"$KODA" run tests/struct_methods.koda
"$KODA" run tests/integer_types.koda
"$KODA" run tests/intern_clear_test.koda
"$KODA" run tests/stdlib_modules_test.koda
"$KODA" run tests/enum_exhaustive.koda

if "$KODA" check --warn-unused tests/warn_unused.koda 2>&1 | grep -q "unused variable"; then
  echo "warn-unused OK"
else
  echo "expected warn-unused warnings" >&2
  exit 1
fi

echo "==> release smoke OK"
