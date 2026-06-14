package nativebuild

import (
	"path/filepath"
	"runtime"
	"strings"
)

// DefaultExeName returns the default native executable name for a .koda path.
func DefaultExeName(source string) string {
	name := strings.TrimSuffix(filepath.Base(source), filepath.Ext(source))
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}
