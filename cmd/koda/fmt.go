package main

import (
	"fmt"
	"os"
	"strings"

	"koda/internal/formatter"
)

func runFmtCmd(args []string) error {
	var check bool
	var paths []string
	for _, a := range args {
		switch a {
		case "--check":
			check = true
		default:
			paths = append(paths, a)
		}
	}
	if len(paths) == 0 {
		return fmt.Errorf("usage: koda fmt [--check] <file.koda> [more files...] [./...]\n       koda fmt [--check] ./...")
	}

	files, err := expandFmtTargets(paths)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("koda fmt: no .koda files matched")
	}

	var checkFailed bool
	for _, path := range files {
		raw, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		text := normalizeNL(string(raw))
		out, err := formatter.Format(text)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		if check {
			if out != text {
				fmt.Fprintf(os.Stderr, "not formatted: %s\n", path)
				checkFailed = true
			}
			continue
		}
		if out == text {
			continue
		}
		if err := os.WriteFile(path, []byte(out), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
		fmt.Println(path)
	}
	if checkFailed {
		return fmt.Errorf("koda fmt --check: one or more files need formatting")
	}
	return nil
}

func normalizeNL(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.ReplaceAll(s, "\r", "\n")
}

func expandFmtTargets(paths []string) ([]string, error) {
	return expandKodaTargets(paths)
}
