package main

import (
	"os"
	"path/filepath"
)

// toolDisplayName is the binary basename (e.g. kodawrap.exe) for CLI messages ("Run … -help").
func toolDisplayName() string {
	return filepath.Base(os.Args[0])
}

// generatedByBrand is the canonical name stamped into .koda / wrapper.c / docs output.
const generatedByBrand = "kodawrap"
