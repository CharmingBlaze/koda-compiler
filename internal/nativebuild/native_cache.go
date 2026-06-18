package nativebuild

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"koda/internal/kodahome"
)

func nativeCacheDisabled() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("KODA_NATIVE_CACHE")))
	return v == "0" || v == "false" || v == "no" || v == "off"
}

func verboseNativeCompile() bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv("KODA_VERBOSE")))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

func nativeCacheDir(projectRoot string) string {
	return filepath.Join(projectRoot, ".KODA_build", "native")
}

func nativeSourceCacheKey(srcAbs string) (string, error) {
	fi, err := os.Stat(srcAbs)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	_, _ = h.Write([]byte(filepath.Clean(srcAbs)))
	_, _ = h.Write([]byte{0})
	_, _ = h.Write([]byte(fmt.Sprintf("%d:%d", fi.Size(), fi.ModTime().UnixNano())))
	sum := hex.EncodeToString(h.Sum(nil))
	base := strings.TrimSuffix(filepath.Base(srcAbs), filepath.Ext(srcAbs))
	if base == "" {
		base = "native"
	}
	return base + "-" + sum[:16], nil
}

func nativeObjectCacheValid(objPath, srcAbs string) bool {
	objInfo, err := os.Stat(objPath)
	if err != nil {
		return false
	}
	srcInfo, err := os.Stat(srcAbs)
	if err != nil {
		return false
	}
	return !objInfo.IsDir() && objInfo.ModTime().After(srcInfo.ModTime()) || objInfo.ModTime().Equal(srcInfo.ModTime())
}

// materializeNativeObjects compiles KODA_NATIVE_SOURCES .c files to cached object files.
// Non-.c entries are passed through unchanged. Returns paths suitable for the link step.
func materializeNativeObjects(cc string, sources []string, projectRoot, sdkRoot string, opts BuildOptions, log func(string)) ([]string, error) {
	if len(sources) == 0 || nativeCacheDisabled() {
		return sources, nil
	}

	cacheDir := nativeCacheDir(projectRoot)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}

	runtimeInclude := filepath.Join(sdkRoot, "runtime", "src")
	includeDirs := []string{runtimeInclude}
	if inc, _, ok := vendoredRaylibStatic(projectRoot); ok {
		includeDirs = append(includeDirs, inc)
	}

	out := make([]string, 0, len(sources))
	for _, src := range sources {
		src = strings.TrimSpace(src)
		if src == "" {
			continue
		}
		srcAbs, err := filepath.Abs(src)
		if err != nil {
			return nil, err
		}
		if !strings.EqualFold(filepath.Ext(srcAbs), ".c") {
			out = append(out, srcAbs)
			continue
		}
		objPath, err := compileNativeSourceCached(cc, srcAbs, cacheDir, includeDirs, opts, log)
		if err != nil {
			return nil, err
		}
		out = append(out, objPath)
	}
	return out, nil
}

func compileNativeSourceCached(cc, srcAbs, cacheDir string, includeDirs []string, opts BuildOptions, log func(string)) (string, error) {
	key, err := nativeSourceCacheKey(srcAbs)
	if err != nil {
		return "", err
	}
	ext := ".o"
	if runtime.GOOS == "windows" {
		ext = ".obj"
	}
	objPath := filepath.Join(cacheDir, key+ext)
	if nativeObjectCacheValid(objPath, srcAbs) {
		if log != nil {
			log(fmt.Sprintf("  native cache hit: %s\n", filepath.Base(srcAbs)))
		}
		return objPath, nil
	}

	args := []string{"-c", srcAbs, "-o", objPath, opts.llcOptFlag()}
	if opts.Debug {
		args = append(args, "-g")
	}
	if res, err := kodahome.BundledClangResourceFlags(); err == nil {
		args = append(args, res...)
	}
	for _, inc := range includeDirs {
		args = append(args, "-I", inc)
	}
	args = append(args, "-Wno-pointer-sign", "-Wno-incompatible-pointer-types-discards-qualifiers")

	if log != nil {
		log(fmt.Sprintf("  compiling native: %s\n", filepath.Base(srcAbs)))
	}
	cmd := exec.Command(cc, args...)
	if verboseNativeCompile() {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = io.Discard
		cmd.Stderr = io.Discard
	}
	if err := cmd.Run(); err != nil {
		_ = os.Remove(objPath)
		return "", fmt.Errorf("compile %s: %w", srcAbs, err)
	}
	return objPath, nil
}
