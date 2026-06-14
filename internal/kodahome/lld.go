package kodahome

import (
	"runtime"
)

func lldExeName() string {
	if runtime.GOOS == "windows" {
		return "ld.lld.exe"
	}
	return "ld.lld"
}

// HasBundledLLD reports whether LLVM's LLD driver sits next to the bundled Clang
// (toolchain/bin/ld.lld or llvm/bin/ld.lld). When true, nativebuild can pass -fuse-ld=lld
// without requiring a system linker.
func HasBundledLLD() bool {
	_, ok := BundledLLDPath()
	return ok
}
