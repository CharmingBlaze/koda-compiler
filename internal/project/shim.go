package project

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"koda/internal/kodahome"
)

var raylibShimFiles = []string{"raylib.koda", "wrapper.c"}

// SyncRaylibShim copies the canonical raylib_shim into dest, overwriting existing files.
// Prefers SDK wrappers/raylib_shim; falls back to the embedded graphics template.
func SyncRaylibShim(dest string) error {
	dest = filepath.Clean(dest)
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}
	if src, ok := canonicalRaylibShimDir(); ok {
		return copyRaylibShimFiles(src, dest)
	}
	return copyTemplateSubtreeOverwrite("graphics", "wrappers/raylib_shim", dest)
}

func canonicalRaylibShimDir() (string, bool) {
	wdir, err := kodahome.WrappersDir()
	if err != nil {
		return "", false
	}
	dir := filepath.Join(wdir, "raylib_shim")
	if raylibShimComplete(dir) {
		return dir, true
	}
	return "", false
}

func raylibShimComplete(dir string) bool {
	for _, name := range raylibShimFiles {
		fi, err := os.Stat(filepath.Join(dir, name))
		if err != nil || fi.IsDir() {
			return false
		}
	}
	return true
}

func copyRaylibShimFiles(src, dest string) error {
	for _, name := range raylibShimFiles {
		data, err := os.ReadFile(filepath.Join(src, name))
		if err != nil {
			return fmt.Errorf("raylib shim %s: %w", name, err)
		}
		target := filepath.Join(dest, name)
		if err := os.WriteFile(target, data, 0644); err != nil {
			return err
		}
	}
	return nil
}

func copyTemplateSubtreeOverwrite(template, rel, dest string) error {
	tplRoot := filepath.ToSlash(filepath.Join("templates", template, rel))
	if _, err := fs.Stat(templateFS, tplRoot); err != nil {
		return fmt.Errorf("template path %q: %w", tplRoot, err)
	}
	return fs.WalkDir(templateFS, tplRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		sub, err := filepath.Rel(filepath.FromSlash(tplRoot), filepath.FromSlash(path))
		if err != nil || sub == "." {
			return err
		}
		target := filepath.Join(dest, sub)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, err := templateFS.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		return os.WriteFile(target, data, 0644)
	})
}
