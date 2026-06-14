//go:build release

package kodahome

import (
	"errors"
	"runtime"
	"sync"

	fujembed "koda/internal/embed"
)

var (
	embeddedOnce sync.Once
	embeddedTC   *Toolchain
	embeddedErr  error
)

func embeddedToolchain() (*Toolchain, error) {
	embeddedOnce.Do(func() {
		embeddedTC, embeddedErr = setupEmbeddedClangToolchain()
	})
	return embeddedTC, embeddedErr
}

func setupEmbeddedClangToolchain() (*Toolchain, error) {
	if _, err := fujembed.Extract(); err != nil {
		return nil, errors.Join(ErrIncompleteEmbeddedToolchain, err)
	}
	clang, err := fujembed.ClangPath()
	if err != nil {
		return nil, errors.Join(ErrIncompleteEmbeddedToolchain, err)
	}
	llc, err := fujembed.LLCPath()
	if err != nil {
		return nil, errors.Join(ErrIncompleteEmbeddedToolchain, err)
	}
	lib, err := fujembed.RuntimeLibPath()
	if err != nil {
		return nil, errors.Join(ErrIncompleteEmbeddedToolchain, err)
	}
	tc := &Toolchain{
		LLC:        llc,
		LLD:        "",
		Clang:      clang,
		RuntimeLib: lib,
		LinkMode:   LinkClang,
	}
	if runtime.GOOS == "windows" {
		lld, err := fujembed.LLDPathWindows()
		if err != nil {
			return nil, errors.Join(ErrIncompleteEmbeddedToolchain, err)
		}
		tc.LLD = lld
	}
	return tc, nil
}
