package nativebuild

import (
	"os"
	"path/filepath"
	"strings"

	"koda/internal/kodahome"
)

// vendoredRaylibStatic returns include dir and a static (or import) library path when a
// third_party/raylib_static/stage tree exists. Resolution order:
//  1) KODA_RAYLIB_STAGE — explicit path to a stage directory with include/ + lib/
//  2) <cwd>/third_party/raylib_static/stage (current project root)
//  3) <dir of koda.exe>/third_party/raylib_static/stage (offline SDK layout)
//
// Set KODA_USE_VENDORED_RAYLIB=0 to disable even when a stage exists.
func vendoredRaylibStatic(rootDir string) (includeDir, archive string, ok bool) {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("KODA_USE_VENDORED_RAYLIB")))
	if v == "0" || v == "false" || v == "no" {
		return "", "", false
	}
	if s := strings.TrimSpace(os.Getenv("KODA_RAYLIB_STAGE")); s != "" {
		return raylibStageAt(filepath.Clean(s))
	}
	tries := []string{filepath.Join(rootDir, "third_party", "raylib_static", "stage")}
	if inst, err := kodahome.InstallDir(); err == nil {
		tries = append(tries, filepath.Join(inst, "third_party", "raylib_static", "stage"))
	}
	for _, stage := range tries {
		if inc, arch, ok := raylibStageAt(stage); ok {
			return inc, arch, true
		}
	}
	return "", "", false
}

func raylibStageAt(stage string) (includeDir, archive string, ok bool) {
	inc := filepath.Join(stage, "include")
	archives := []string{
		filepath.Join(stage, "lib", "libraylib.a"),
		filepath.Join(stage, "lib", "libraylibdll.a"),
		filepath.Join(stage, "lib", "libraylib.dylib"),
		filepath.Join(stage, "lib", "raylib.lib"),
	}
	for _, a := range archives {
		if st, err := os.Stat(a); err == nil && !st.IsDir() {
			if fi, err := os.Stat(inc); err == nil && fi.IsDir() {
				return inc, a, true
			}
		}
	}
	return "", "", false
}
