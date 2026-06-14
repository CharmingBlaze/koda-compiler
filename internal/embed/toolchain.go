//go:build release

package fujembed

import (
	"embed"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
)

//go:embed windows linux darwin
var assets embed.FS

var (
	once     sync.Once
	cacheDir string
	cacheErr error
)

// Extract unpacks the embedded toolchain to a temp directory the first time it is called.
func Extract() (string, error) {
	once.Do(func() {
		dir, err := os.MkdirTemp("", "koda-*")
		if err != nil {
			cacheErr = err
			return
		}
		if err := extractAll(dir); err != nil {
			_ = os.RemoveAll(dir)
			cacheErr = err
			return
		}
		cacheDir = dir
	})
	return cacheDir, cacheErr
}

func extractAll(dest string) error {
	platform := platformKey()
	files := platformFiles()

	for _, name := range files {
		src := path.Join(platform, name)
		data, err := assets.ReadFile(src)
		if err != nil {
			return fmt.Errorf("missing bundled file %s: %w", src, err)
		}

		dst := filepath.Join(dest, name)
		mode := os.FileMode(0o755)
		if filepath.Ext(name) == ".a" {
			mode = 0o644
		}
		if err := os.WriteFile(dst, data, mode); err != nil {
			return fmt.Errorf("write %s: %w", dst, err)
		}
	}
	return nil
}

func platformKey() string {
	return runtime.GOOS + "/" + runtime.GOARCH
}

func platformFiles() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{"clang.exe", "llc.exe", "lld.exe", "libkoda_runtime.a"}
	default:
		return []string{"clang", "llc", "libkoda_runtime.a"}
	}
}

// ClangPath returns the path to the embedded clang binary.
func ClangPath() (string, error) {
	dir, err := Extract()
	if err != nil {
		return "", err
	}
	name := "clang"
	if runtime.GOOS == "windows" {
		name = "clang.exe"
	}
	return filepath.Join(dir, name), nil
}

// LLCPath returns the path to the embedded llc binary.
func LLCPath() (string, error) {
	dir, err := Extract()
	if err != nil {
		return "", err
	}
	name := "llc"
	if runtime.GOOS == "windows" {
		name = "llc.exe"
	}
	return filepath.Join(dir, name), nil
}

// RuntimeLibPath returns the path to the embedded runtime library.
func RuntimeLibPath() (string, error) {
	dir, err := Extract()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "libkoda_runtime.a"), nil
}

// LLDPathWindows returns the path to embedded lld.exe (Windows release builds only).
func LLDPathWindows() (string, error) {
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("lld embedding is only used on windows")
	}
	dir, err := Extract()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "lld.exe"), nil
}
