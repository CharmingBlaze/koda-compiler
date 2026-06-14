package main

import (
	"fmt"
	"os"

	"koda/internal/formatter"
)

func runLint(args []string) error {
	paths := args
	if len(paths) == 0 {
		paths = []string{"./..."}
	}
	files, err := expandKodaTargets(paths)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("koda lint: no .koda files matched")
	}
	var checkErrs int
	var fmtErrs int
	for _, path := range files {
		if err := checkFile(path); err != nil {
			fmt.Fprintf(os.Stderr, "lint check: %s: %v\n", path, err)
			checkErrs++
			continue
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		text := normalizeNL(string(raw))
		out, err := formatter.Format(text)
		if err != nil {
			fmt.Fprintf(os.Stderr, "lint fmt: %s: %v\n", path, err)
			fmtErrs++
			continue
		}
		if out != text {
			fmt.Fprintf(os.Stderr, "lint fmt: not formatted: %s\n", path)
			fmtErrs++
		}
	}
	if checkErrs > 0 || fmtErrs > 0 {
		return fmt.Errorf("lint failed: %d check error(s), %d format issue(s)", checkErrs, fmtErrs)
	}
	fmt.Printf("lint OK (%d files)\n", len(files))
	return nil
}

func runCheckAll(args []string) error {
	paths := args
	if len(paths) == 0 {
		paths = []string{"./..."}
	}
	files, err := expandKodaTargets(paths)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("koda check: no .koda files matched")
	}
	for _, path := range files {
		if err := checkFile(path); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
	}
	fmt.Printf("OK (%d files)\n", len(files))
	return nil
}
