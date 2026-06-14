package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"koda/api"
	"koda/internal/nativebuild"
	"koda/internal/project"
)

// runTest compiles and runs .koda test files (default: tests/*.koda under repo or project).
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
		return fmt.Errorf("no test files found (pass paths or add tests/*.koda)")
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
		testsDir := filepath.Join(ctx.Root, "tests")
		if entries, err := os.ReadDir(testsDir); err == nil {
			var out []string
			for _, e := range entries {
				if e.IsDir() || !strings.HasSuffix(e.Name(), ".koda") {
					continue
				}
				if strings.HasSuffix(e.Name(), "_module.koda") {
					continue
				}
				out = append(out, filepath.Join(testsDir, e.Name()))
			}
			if len(out) > 0 {
				return out, nil
			}
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
	var out []string
	seen := make(map[string]bool)
	for _, root := range roots {
		abs, err := filepath.Abs(root)
		if err != nil {
			continue
		}
		testsDir := filepath.Join(abs, "tests")
		entries, err := os.ReadDir(testsDir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".koda") {
				continue
			}
			if strings.HasSuffix(e.Name(), "_module.koda") {
				continue
			}
			p := filepath.Join(testsDir, e.Name())
			if seen[p] {
				continue
			}
			seen[p] = true
			out = append(out, p)
		}
		if len(out) > 0 {
			break
		}
	}
	return out, nil
}
