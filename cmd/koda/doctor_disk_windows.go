//go:build windows

package main

func diskFreeBytes(path string) (uint64, bool) {
	return 0, false
}
