# Koda — build from repo root with GNU Make (MSYS2/Unix/Git Bash).
# On Windows, use GNU Make (e.g. mingw32-make from MSYS2/MinGW), not a shadowed `make.bat`.
# Requires sh.exe on PATH (Git for Windows / MSYS) for mkdir -p / rm -rf in recipes.
SHELL := sh.exe
.SHELLFLAGS := -ec
#
#   make              — runtime static lib, then koda + kodawrap, then go test
#   make runtime-lib  — only runtime/libkoda_runtime.a
#   make raylib-lib   — CMake Raylib into third_party/raylib_static/stage/ (needs raylib/ + cmake)

.PHONY: all build koda kodawrap wrapgen runtime-lib raylib-lib raylib-clean test soak-test fmt clean

all: build test

build: koda kodawrap

koda:
	@mkdir -p bin
	go build -trimpath -ldflags "-s -w" -o bin/koda ./cmd/koda

kodawrap:
	@mkdir -p bin
	go build -trimpath -ldflags "-s -w" -o bin/kodawrap ./cmd/wrapgen

wrapgen: kodawrap
	@mkdir -p bin
	go build -trimpath -ldflags "-s -w" -o bin/wrapgen ./cmd/wrapgen

# Static library required for native `koda build` and for codegen tests.
runtime-lib:
	$(MAKE) -C runtime

raylib-lib:
	$(MAKE) -C third_party/raylib_static

raylib-clean:
	$(MAKE) -C third_party/raylib_static clean

test: runtime-lib
	go test ./... -count=1

soak-test: koda runtime-lib
	./bin/koda run tests/gc_pressure_expr.koda
	./bin/koda run tests/globals_perf.koda
	./bin/koda run tests/gc_soak.koda

fmt:
	gofmt -w .

clean:
	rm -rf bin .KODA_build
	go clean -cache
	$(MAKE) -C runtime clean
