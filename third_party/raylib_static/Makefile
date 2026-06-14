# Build Raylib from the repo's `raylib/` tree (CMake), then install into ./stage/ for Koda native links.
#
# From repo root:
#   make -C third_party/raylib_static          # static lib (default) -> stage/lib/libraylib.a + headers
#   make -C third_party/raylib_static shared   # shared -> stage-shared/ with raylib.dll (Windows) + import lib
#
# After `make` (static), `koda build` picks up stage/ automatically (see KODA_USE_VENDORED_RAYLIB in nativebuild).
# For shared builds, use the DLL/import lib with your own KODA_LINKFLAGS; the compiler still prefers static stage/ when present.

SHELL := sh.exe
.SHELLFLAGS := -ec

CMAKE ?= cmake
RAYLIB_ROOT ?= $(abspath ../../raylib)

SHARED ?=
ifeq ($(SHARED),1)
  BUILD_DIR := $(abspath cmake-build-shared)
  STAGE := $(abspath stage-shared)
  SHARED_FLAG := ON
else
  BUILD_DIR := $(abspath cmake-build)
  STAGE := $(abspath stage)
  SHARED_FLAG := OFF
endif

.PHONY: all static shared configure build install clean

all: install

static:
	@$(MAKE) SHARED=0 install

shared:
	@$(MAKE) SHARED=1 install

configure:
	@test -f "$(RAYLIB_ROOT)/CMakeLists.txt" || (echo "Expected Raylib sources at $(RAYLIB_ROOT) (e.g. git submodule)."; exit 1)
	"$(CMAKE)" -S "$(RAYLIB_ROOT)" -B "$(BUILD_DIR)" \
		-DCMAKE_BUILD_TYPE=Release \
		-DBUILD_SHARED_LIBS=$(SHARED_FLAG) \
		-DBUILD_EXAMPLES=OFF \
		-DCMAKE_INSTALL_PREFIX="$(STAGE)"

build: configure
	"$(CMAKE)" --build "$(BUILD_DIR)" --parallel

install: build
	rm -rf "$(STAGE)"
	"$(CMAKE)" --install "$(BUILD_DIR)" --prefix "$(STAGE)"
	@echo "Raylib installed under $(STAGE)/ (lib/, include/; shared builds also have bin/ with .dll on Windows)"

clean:
	rm -rf "$(abspath cmake-build)" "$(abspath cmake-build-shared)" "$(abspath stage)" "$(abspath stage-shared)"
