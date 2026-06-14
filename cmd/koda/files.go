package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var kodaSkipDirNames = map[string]bool{
	".git":                true,
	".KODA_build":         true,
	".koda_build":         true,
	"bin":                 true,
	"node_modules":        true,
	"vendor":              true,
	"dist":                true,
	"_legacy":             true,
	"_wg_legacy":          true,
	"_wg_out":             true,
	"raylib_full_wrapper": true,
	"test_output":         true,
}

func splitAtDoubleDash(args []string) (before, after []string) {
	for i, a := range args {
		if a == "--" {
			return args[:i], args[i+1:]
		}
	}
	return args, nil
}

func expandKodaTargets(paths []string) ([]string, error) {
	seen := map[string]struct{}{}
	for _, p := range paths {
		switch p {
		case "./...", "...":
			list, err := collectKodaFiles(".")
			if err != nil {
				return nil, err
			}
			for _, f := range list {
				seen[f] = struct{}{}
			}
			continue
		}
		st, err := os.Stat(p)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", p, err)
		}
		if st.IsDir() {
			list, err := collectKodaFiles(p)
			if err != nil {
				return nil, err
			}
			for _, f := range list {
				seen[f] = struct{}{}
			}
			continue
		}
		if !strings.HasSuffix(strings.ToLower(p), ".koda") {
			return nil, fmt.Errorf("%s: expected a .koda file or directory", p)
		}
		seen[p] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for f := range seen {
		out = append(out, f)
	}
	sort.Strings(out)
	return out, nil
}

func collectKodaFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != root && kodaSkipDirNames[d.Name()] {
				return fs.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(strings.ToLower(path), ".koda") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}
