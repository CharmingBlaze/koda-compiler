#!/usr/bin/env bash
# Native compile/run smoke used by GitHub Actions on Linux, macOS, and Windows (Git Bash).
# Requires: KODA_CI_BIN (path to koda or koda.exe), CI_ARTIFACT_DIR (writable output directory).
set -euo pipefail

KODA="${KODA_CI_BIN:?KODA_CI_BIN is required}"
OUT="${CI_ARTIFACT_DIR:?CI_ARTIFACT_DIR is required}"

EXE=""
case "$(uname -s)" in
  MINGW* | MSYS* | CYGWIN*) EXE=.exe ;;
esac

mkdir -p "$OUT"

run_out() {
  local lines="$1"
  local bin="$2"
  shift 2
  "$bin" "$@" | head -n "$lines"
}

echo "==> native smoke: KODA=$KODA OUT=$OUT EXE=$EXE"

"$KODA" build tests/smoke_native.koda -o "$OUT/smoke_native$EXE"
test -x "$OUT/smoke_native$EXE"
run_out 5 "$OUT/smoke_native$EXE"

"$KODA" build tests/for_of_pairs.koda -o "$OUT/for_of_pairs$EXE"
test -x "$OUT/for_of_pairs$EXE"
run_out 10 "$OUT/for_of_pairs$EXE"

"$KODA" build tests/nullish_assign.koda -o "$OUT/nullish_assign$EXE"
run_out 5 "$OUT/nullish_assign$EXE"

"$KODA" build tests/for_of_dynamic_range.koda -o "$OUT/for_of_dynamic_range$EXE"
run_out 5 "$OUT/for_of_dynamic_range$EXE"

"$KODA" build tests/defer_test.koda -o "$OUT/defer_test$EXE"
run_out 5 "$OUT/defer_test$EXE"

"$KODA" build tests/assert_string_eq.koda -o "$OUT/assert_string_eq$EXE"
run_out 5 "$OUT/assert_string_eq$EXE"

"$KODA" build tests/unary_plus_fold.koda -o "$OUT/unary_plus_fold$EXE"
run_out 5 "$OUT/unary_plus_fold$EXE"

"$KODA" build tests/gc_shadow_multi_return_pop.koda -o "$OUT/gc_shadow_multi_return_pop$EXE"
run_out 5 "$OUT/gc_shadow_multi_return_pop$EXE"

"$KODA" build tests/gc_control_test.koda -o "$OUT/gc_control_test$EXE"
run_out 5 "$OUT/gc_control_test$EXE"

"$KODA" run tests/vec3_test.koda
"$KODA" run tests/vec3_math_lerp_direct_test.koda
"$KODA" run tests/math_lerp_member_test.koda
"$KODA" run tests/array_push_growth_test.koda
"$KODA" run tests/multi_while_flow_test.koda

"$KODA" run tests/math_test.koda
"$KODA" run tests/json_test.koda
"$KODA" run tests/random_test.koda
"$KODA" run tests/io_test.koda
"$KODA" test tests/math_test.koda tests/json_test.koda tests/random_test.koda

"$KODA" build --debug tests/hello.koda -o "$OUT/hello_debug$EXE"
test -x "$OUT/hello_debug$EXE"

echo "==> native smoke OK"
