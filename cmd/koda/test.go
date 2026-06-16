package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"koda/api"
	"koda/internal/nativebuild"
	"koda/internal/project"
)

// runTest compiles and runs .koda test files via the native pipeline (same as koda run).
func runTest(args []string) error {
	var paths []string
	noOpt := false
	verbose := false
	failFast := false
	runPattern := ""
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--no-opt":
			noOpt = true
		case "-v", "--verbose":
			verbose = true
		case "--failfast":
			failFast = true
		case "-run":
			if i+1 >= len(args) {
				return fmt.Errorf("-run requires a pattern")
			}
			i++
			runPattern = args[i]
		default:
			if strings.HasPrefix(args[i], "-") {
				return fmt.Errorf("unknown flag: %s", args[i])
			}
			paths = append(paths, args[i])
		}
	}
	if len(paths) == 0 {
		found, err := defaultTestFiles()
		if err != nil {
			return err
		}
		paths = found
	}
	if len(paths) == 0 {
		return fmt.Errorf("no test files found (add *_test.koda files or pass paths explicitly)")
	}
	if runPattern != "" {
		filtered := make([]string, 0, len(paths))
		for _, p := range paths {
			if strings.Contains(filepath.Base(p), runPattern) || strings.Contains(p, runPattern) {
				filtered = append(filtered, p)
			}
		}
		paths = filtered
		if len(paths) == 0 {
			return fmt.Errorf("no tests matched -run %q", runPattern)
		}
	}

	opts := nativebuild.BuildOptions{NoOpt: noOpt}
	fail := 0
	for _, p := range paths {
		abs, err := filepath.Abs(p)
		if err != nil {
			return err
		}
		if _, err := os.Stat(abs); err != nil {
			fmt.Fprintf(os.Stderr, "FAIL %s: %v\n", p, err)
			fail++
			if failFast {
				break
			}
			continue
		}
		if verbose {
			fmt.Printf("==> %s\n", abs)
		}
		if err := withProject(abs, func() error {
			return api.RunWithBuildOptions(abs, "", opts)
		}); err != nil {
			fmt.Fprintf(os.Stderr, "FAIL %s: %v\n", p, err)
			fail++
			if failFast {
				break
			}
			continue
		}
		if verbose {
			fmt.Printf("ok %s\n", p)
		}
	}
	if fail > 0 {
		return fmt.Errorf("%d test(s) failed", fail)
	}
	fmt.Printf("%d passed\n", len(paths))
	return nil
}

func defaultTestFiles() ([]string, error) {
	if ctx, err := project.LoadContext(cwd()); err == nil && ctx != nil {
		if out, err := discoverProjectTestFiles(ctx.Root); err != nil {
			return nil, err
		} else if len(out) > 0 {
			return out, nil
		}
	}
	exe, err := os.Executable()
	if err != nil {
		return nil, err
	}
	roots := []string{
		filepath.Join(filepath.Dir(exe), ".."),
		".",
	}
	for _, root := range roots {
		abs, err := filepath.Abs(root)
		if err != nil {
			continue
		}
		if out, err := discoverCompilerTestFiles(abs); err != nil {
			return nil, err
		} else if len(out) > 0 {
			return out, nil
		}
	}
	return nil, nil
}

// discoverProjectTestFiles finds Go-style *_test.koda files under a user project.
func discoverProjectTestFiles(root string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			base := d.Name()
			if base == ".git" || base == "node_modules" || base == "dist" || base == ".koda_build" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(d.Name(), "_test.koda") {
			return nil
		}
		out = append(out, path)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sortStrings(out)
	return out, nil
}

// discoverCompilerTestFiles prefers *_test.koda in tests/; falls back to all tests/*.koda
// (compiler harness) when no *_test.koda files exist.
func discoverCompilerTestFiles(root string) ([]string, error) {
	testsDir := filepath.Join(root, "tests")
	entries, err := os.ReadDir(testsDir)
	if err != nil {
		return nil, nil
	}
	var named []string
	var legacy []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".koda") {
			continue
		}
		if strings.HasSuffix(e.Name(), "_module.koda") {
			continue
		}
		p := filepath.Join(testsDir, e.Name())
		if strings.HasSuffix(e.Name(), "_test.koda") {
			named = append(named, p)
			continue
		}
		legacy = append(legacy, p)
	}
	if len(named) > 0 {
		sortStrings(named)
		return named, nil
	}
	sortStrings(legacy)
	return legacy, nil
}

func sortStrings(ss []string) {
	sort.Strings(ss)
}
