package kodahome

import (
	"testing"
)

func TestTarEntryNameSafe(t *testing.T) {
	if tarEntryNameSafe("toolchain/bin/clang.exe") != true {
		t.Fatal("expected safe path")
	}
	if tarEntryNameSafe("../evil") {
		t.Fatal("expected unsafe")
	}
	if tarEntryNameSafe("/abs") {
		t.Fatal("expected unsafe")
	}
}

func TestPlaceholderBundleDoesNotDeclareClang(t *testing.T) {
	ok, err := tarballDeclaresToolchainClang(bundledToolchainTarGz)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("placeholder archive must not declare a clang binary")
	}
}

func TestBundledIncludeFlagsNoPanic(t *testing.T) {
	_, err := BundledClangResourceFlags()
	if err != nil {
		t.Fatal(err)
	}
}

func TestLLCResolves(t *testing.T) {
	s := LLC()
	if s == "" {
		t.Fatal("empty llc path")
	}
}
