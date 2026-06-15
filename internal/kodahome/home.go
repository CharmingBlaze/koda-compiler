// package kodahome resolves install-relative paths for distributed Koda toolchains
// (bundled Clang under llvm/ or toolchain/, stdlib/, wrappers/) without importing
// the parser or codegen (avoids cycles).
package kodahome

import (
	"os"
	"path/filepath"
	"runtime"
)

// InstallDir returns the SDK root directory (stdlib/, wrappers/, toolchain siblings).
// Honors KODA_HOME when set; otherwise the executable directory, or a parent folder
// that contains stdlib/ (e.g. Koda Studio in build/bin/ with SDK at repo root).
func InstallDir() (string, error) {
	if root, ok := sdkRootFromEnv(); ok {
		return root, nil
	}
	exeDir, err := executableDir()
	if err != nil {
		return "", err
	}
	if hasStdlibRoot(exeDir) {
		return exeDir, nil
	}
	if root, ok := BootstrapSDKRoot(exeDir); ok {
		return root, nil
	}
	return exeDir, nil
}

func executableDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	if exe, err = filepath.EvalSymlinks(exe); err != nil {
		exe, _ = os.Executable()
	}
	return filepath.Abs(filepath.Dir(exe))
}

func clangFileName() string {
	if runtime.GOOS == "windows" {
		return "clang.exe"
	}
	return "clang"
}

func llcFileName() string {
	if runtime.GOOS == "windows" {
		return "llc.exe"
	}
	return "llc"
}

// BundledClangRelPaths are subpaths under InstallDir() checked for a portable LLVM.
// toolchain/ is preferred over llvm/ so a single vendor layout wins consistently.
var BundledClangRelPaths = []string{
	filepath.Join("toolchain", "bin", clangFileName()),
	filepath.Join("llvm", "bin", clangFileName()),
}

// BundledLLCRelPaths are subpaths under InstallDir() for the LLVM static compiler.
var BundledLLCRelPaths = []string{
	filepath.Join("toolchain", "bin", llcFileName()),
	filepath.Join("llvm", "bin", llcFileName()),
}

// Clang returns the C/LLVM driver. See [ClangWithSource] for the resolution order.
func Clang() string {
	p, _ := ClangWithSource()
	return p
}

// LLC returns the llc executable. See [LLCWithSource] for the resolution order.
func LLC() string {
	p, _ := LLCWithSource()
	return p
}

// StdlibDir is installDir/stdlib (may not exist).
func StdlibDir() (string, error) {
	dir, err := InstallDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "stdlib"), nil
}

// WrappersDir is installDir/wrappers (may not exist).
func WrappersDir() (string, error) {
	dir, err := InstallDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "wrappers"), nil
}
