#!/usr/bin/env bash
# ASAN/UBSAN smoke — Linux only. Non-blocking in CI until stable for two weeks.
set -euo pipefail

KODA="${KODA_CI_BIN:?KODA_CI_BIN is required}"
CC_BIN="${CC_BIN:-clang}"

echo "==> ASAN smoke: CC=$CC_BIN KODA=$KODA"

SAN_FLAGS="-fsanitize=address,undefined -fno-omit-frame-pointer -g"
mkdir -p runtime/obj_asan
"$CC_BIN" -c runtime/src/value.c        -O1 -std=c11 -Iruntime/src $SAN_FLAGS -o runtime/obj_asan/value.o
"$CC_BIN" -c runtime/src/object.c       -O1 -std=c11 -Iruntime/src $SAN_FLAGS -o runtime/obj_asan/object.o
"$CC_BIN" -c runtime/src/gc.c           -O1 -std=c11 -Iruntime/src $SAN_FLAGS -o runtime/obj_asan/gc.o
"$CC_BIN" -c runtime/src/koda_runtime.c -O1 -std=c11 -Iruntime/src $SAN_FLAGS -o runtime/obj_asan/koda_runtime.o
llvm-ar rcs runtime/libkoda_runtime_asan.a runtime/obj_asan/value.o runtime/obj_asan/object.o runtime/obj_asan/gc.o runtime/obj_asan/koda_runtime.o

cp runtime/libkoda_runtime.a runtime/libkoda_runtime_normal.a
cp runtime/libkoda_runtime_asan.a runtime/libkoda_runtime.a
trap 'cp runtime/libkoda_runtime_normal.a runtime/libkoda_runtime.a' EXIT

export ASAN_OPTIONS=detect_leaks=1:halt_on_error=1
export UBSAN_OPTIONS=halt_on_error=1

run_one() {
  local f="$1"
  echo "==> $f"
  timeout 120s "$KODA" run "$f"
}

run_one tests/arena_test.koda
run_one tests/incremental_gc_test.koda
run_one tests/gc_control_test.koda
run_one tests/intern_clear_test.koda
run_one tests/stress/stress_mixed_alloc.koda
run_one tests/stress/stress_deep_recursion.koda
run_one tests/stress/stress_string_pressure.koda

echo "ASAN smoke OK"
