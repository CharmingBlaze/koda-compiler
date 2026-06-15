package nativebuild

import "runtime"

func defaultSystemLinkFlags() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{"-lopengl32", "-lgdi32", "-lwinmm"}
	case "darwin":
		return []string{"-framework", "OpenGL", "-framework", "Cocoa", "-framework", "IOKit", "-framework", "CoreVideo"}
	case "linux":
		return []string{"-lGL", "-lm", "-lpthread", "-ldl", "-lrt", "-lX11"}
	default:
		return nil
	}
}

func omitLinkFlag(flags []string, flag string) []string {
	out := make([]string, 0, len(flags))
	for _, f := range flags {
		if f != flag {
			out = append(out, f)
		}
	}
	return out
}
