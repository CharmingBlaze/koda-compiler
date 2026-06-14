package kodahome

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
)

// BundledClangResourceFlags returns extra clang flags (-isystem ...) for portable
// lib/clang/*/include trees shipped next to the executable (extracted or laid out by hand).
func BundledClangResourceFlags() ([]string, error) {
	install, err := InstallDir()
	if err != nil {
		return nil, err
	}
	var flags []string
	for _, inc := range discoverClangBuiltinIncludes(install) {
		flags = append(flags, "-isystem", normalizeClangFilesystemArg(inc))
	}
	return flags, nil
}

// normalizeClangFilesystemArg makes paths safe for Clang argv on Windows by using
// forward slashes, avoiding backslash escape issues in some driver / response-file paths.
func normalizeClangFilesystemArg(p string) string {
	if runtime.GOOS == "windows" {
		return filepath.ToSlash(p)
	}
	return p
}

func discoverClangBuiltinIncludes(install string) []string {
	var bases []string
	for _, root := range []string{"toolchain", "llvm"} {
		bases = append(bases, filepath.Join(install, root, "lib", "clang"))
	}
	// Optional: headers vendored under stdlib/sys-include
	bases = append(bases, filepath.Join(install, "stdlib", "sys-include"))

	var out []string
	seen := make(map[string]bool)
	for _, base := range bases {
		matches, _ := filepath.Glob(filepath.Join(base, "*", "include"))
		sort.Strings(matches)
		for _, m := range matches {
			fi, err := os.Stat(m)
			if err != nil || !fi.IsDir() {
				continue
			}
			abs, err := filepath.Abs(m)
			if err != nil {
				continue
			}
			if seen[abs] {
				continue
			}
			seen[abs] = true
			out = append(out, abs)
		}
	}
	// Flat portable headers (stddef/stdint/stdbool shims) for wrapgen / CI without a host SDK.
	sysFlat := filepath.Join(install, "stdlib", "sys-include")
	if fi, err := os.Stat(sysFlat); err == nil && fi.IsDir() {
		if abs, err := filepath.Abs(sysFlat); err == nil && !seen[abs] {
			seen[abs] = true
			out = append(out, abs)
		}
	}
	return out
}

// ClangWrappedArgs prepends resource flags before other clang arguments.
func ClangWrappedArgs(rest ...string) []string {
	flags, _ := BundledClangResourceFlags()
	return append(flags, rest...)
}
