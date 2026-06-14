package kodahome

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const toolchainLockName = ".koda_toolchain_extract.lock"

// EnsureEnvironment unpacks the embedded toolchain archive when it contains a real
// Clang tree and nothing is installed yet next to the executable. Safe no-op for the
// tiny upstream placeholder archive (no clang members). Set KODA_SKIP_TOOLCHAIN_EXTRACT=1
// to disable. Publishers replace embeddata/bundled_toolchain.tar.gz before release builds.
//
// Extraction is serialized with an O_EXCL lock file (see toolchainLockName) so two
// processes do not write toolchain/bin/clang(.exe) concurrently—important on Windows
// where readers can see partially written PE files.
func EnsureEnvironment(log func(string)) error {
	if strings.TrimSpace(os.Getenv("KODA_SKIP_TOOLCHAIN_EXTRACT")) != "" {
		return nil
	}
	data := bundledToolchainTarGz
	if len(data) < 32 {
		return nil
	}
	ok, err := tarballDeclaresToolchainClang(data)
	if err != nil || !ok {
		return err
	}

	install, err := InstallDir()
	if err != nil {
		return err
	}
	if hasBundledClang(install) {
		return nil
	}

	if log != nil {
		log("First run: extracting bundled Koda toolchain (this may take a moment)...\n")
	}
	if err := withToolchainLock(install, func() error {
		if hasBundledClang(install) {
			return nil
		}
		return extractTarGz(data, install)
	}); err != nil {
		return fmt.Errorf(
			"bundled toolchain extract failed (install dir may be read-only or disk full): %w\n"+
				"Hint: move koda/kodawrap to a normal writable folder with stdlib/ beside them, "+
				"or set KODA_CLANG to an existing Clang and KODA_SKIP_TOOLCHAIN_EXTRACT=1 if you intentionally skip unpack",
			err,
		)
	}
	if log != nil {
		log("Toolchain ready.\n")
	}
	return nil
}

func hasBundledClang(install string) bool {
	for _, rel := range BundledClangRelPaths {
		p := filepath.Join(install, rel)
		fi, err := os.Stat(p)
		if err == nil && !fi.IsDir() {
			return true
		}
	}
	return false
}

// tarballDeclaresToolchainClang returns true if the gzip+tar archive appears to ship
// a clang binary (skips the tiny placeholder archive shipped in-repo).
func tarballDeclaresToolchainClang(data []byte) (bool, error) {
	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return false, nil
	}
	defer zr.Close()
	tr := tar.NewReader(zr)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return false, nil
		}
		if err != nil {
			return false, err
		}
		name := filepath.ToSlash(strings.TrimPrefix(h.Name, "./"))
		if h.Typeflag == tar.TypeDir {
			continue
		}
		if !tarEntryNameSafe(name) {
			continue
		}
		if strings.HasSuffix(name, "/bin/clang") || strings.HasSuffix(name, "/bin/clang.exe") {
			if strings.HasPrefix(name, "toolchain/") || strings.HasPrefix(name, "llvm/") {
				return true, nil
			}
		}
		mode := h.FileInfo().Mode()
		if mode.IsRegular() && h.Size > 0 {
			if _, err := io.CopyN(io.Discard, tr, h.Size); err != nil {
				return false, err
			}
		}
	}
}

func tarEntryNameSafe(name string) bool {
	name = filepath.ToSlash(strings.TrimPrefix(name, "./"))
	if name == "" || strings.Contains(name, "..") {
		return false
	}
	if filepath.IsAbs(name) {
		return false
	}
	// On Windows filepath.IsAbs("/tmp") is false; still reject Unix-rooted tar paths.
	if strings.HasPrefix(name, "/") {
		return false
	}
	return true
}

func extractTarGz(data []byte, destDir string) error {
	zr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer zr.Close()
	tr := tar.NewReader(zr)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		name := filepath.ToSlash(strings.TrimPrefix(h.Name, "./"))
		if !tarEntryNameSafe(name) {
			return fmt.Errorf("unsafe tar path: %q", h.Name)
		}
		target := filepath.Join(destDir, filepath.FromSlash(name))
		mode := h.FileInfo().Mode()
		switch {
		case mode.IsDir():
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case mode.IsRegular():
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(h.Mode&0777))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}
		default:
			if h.Size > 0 {
				if _, err := io.CopyN(io.Discard, tr, h.Size); err != nil {
					return err
				}
			}
		}
	}
}

func withToolchainLock(install string, fn func() error) error {
	lockPath := filepath.Join(install, toolchainLockName)
	for attempt := 0; attempt < 600; attempt++ {
		lf, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
		if err == nil {
			_, _ = lf.WriteString(fmt.Sprintf("%d\n", os.Getpid()))
			_ = lf.Close()
			defer func() { _ = os.Remove(lockPath) }()
			return fn()
		}
		time.Sleep(50 * time.Millisecond)
		if hasBundledClang(install) {
			return nil
		}
	}
	return fmt.Errorf(
		"timeout waiting for toolchain lock %s (another koda/kodawrap may be extracting; retry in a moment)",
		lockPath,
	)
}
