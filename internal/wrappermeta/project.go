package wrappermeta

import (
	"os"
	"path/filepath"
	"strings"
)

// FindWrapperDir resolves a wrapper package directory from a native source path
// (typically .../wrapper.c) or a directory containing META.json.
func FindWrapperDir(nativeSource string) (string, bool) {
	p := strings.TrimSpace(nativeSource)
	if p == "" {
		return "", false
	}
	if strings.HasSuffix(filepath.ToSlash(p), "/"+FileName) {
		return filepath.Dir(p), true
	}
	if fi, err := os.Stat(p); err == nil && fi.IsDir() {
		if _, err := os.Stat(filepath.Join(p, FileName)); err == nil {
			return p, true
		}
	}
	if strings.HasSuffix(filepath.ToSlash(p), "/wrapper.c") {
		dir := filepath.Dir(p)
		if _, err := os.Stat(filepath.Join(dir, FileName)); err == nil {
			return dir, true
		}
	}
	return "", false
}

// CheckProjectSources reports drift for wrapper.c paths listed in native.sources.
func CheckProjectSources(sources []string) ([]DriftReport, error) {
	seen := map[string]bool{}
	var reports []DriftReport
	for _, src := range sources {
		src = strings.TrimSpace(src)
		if src == "" {
			continue
		}
		dir, ok := FindWrapperDir(src)
		if !ok {
			continue
		}
		if seen[dir] {
			continue
		}
		seen[dir] = true
		meta, err := LoadDir(dir)
		if err != nil {
			continue
		}
		report, err := CheckDrift(meta)
		if err != nil {
			return reports, err
		}
		report.Dir = dir
		if report.Stale() {
			reports = append(reports, report)
		}
	}
	return reports, nil
}
