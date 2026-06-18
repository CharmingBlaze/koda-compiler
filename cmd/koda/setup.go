package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"koda/internal/kodahome"
	"koda/internal/nativebuild"
	"koda/internal/project"
)

func runSetup(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: koda setup <target>\n\nTargets:\n  raylib [--shim]   Configure graphics linking (full Raylib wrapper by default)")
	}
	switch strings.ToLower(args[0]) {
	case "raylib":
		return setupRaylib(args[1:])
	default:
		return fmt.Errorf("unknown setup target %q (available: raylib)", args[0])
	}
}

func setupRaylib(args []string) error {
	shim := false
	var positional []string
	for _, a := range args {
		switch strings.ToLower(a) {
		case "--shim":
			shim = true
		case "--full", "-full":
			// Legacy alias; full wrapper is the default.
		default:
			positional = append(positional, a)
		}
	}

	root := cwd()
	if len(positional) > 0 {
		root = positional[0]
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	root = abs

	cfg, cfgRoot, err := project.Find(root)
	if err != nil {
		return err
	}
	projectRoot := root
	if cfg != nil {
		projectRoot = cfgRoot
	}

	if shim {
		if err := ensureRaylibShim(projectRoot); err != nil {
			return err
		}
	}

	if cfg == nil {
		cfg = &project.Config{
			Name:    filepath.Base(root),
			Version: "0.1.0",
			Entry:   "src/main.koda",
			Lint:    "beginner",
		}
		cfgRoot = root
	}
	if cfgRoot == "" {
		cfgRoot = projectRoot
	}
	if strings.TrimSpace(cfg.Lint) == "" {
		cfg.Lint = "beginner"
	}
	if shim {
		cfg.Native.Sources = mergeUnique(stripRaylibFullSources(cfg.Native.Sources), "wrappers/raylib_shim/wrapper.c")
	} else {
		cfg.Native.Sources = []string{"wrappers/raylib/wrapper.c"}
	}
	cfg.Native.Graphics = true

	if err := project.Save(cfgRoot, cfg); err != nil {
		return err
	}

	if shim {
		fmt.Println("Raylib shim setup complete.")
		fmt.Printf("  project: %s\n", filepath.Join(cfgRoot, project.FileName))
		fmt.Printf("  shim:    %s\n", filepath.Join(projectRoot, "wrappers", "raylib_shim"))
		fmt.Println()
		fmt.Println("Graphics linking is enabled via koda.json (\"graphics\": true).")
		fmt.Println("Include the shim in your .koda file:")
		fmt.Println(`  #include "wrappers/raylib_shim/raylib.koda"`)
		fmt.Println()
		fmt.Println("For the full Raylib API (default): koda setup raylib")
	} else {
		fmt.Println("Full Raylib setup complete.")
		fmt.Printf("  project: %s\n", filepath.Join(cfgRoot, project.FileName))
		fmt.Println()
		fmt.Println("Graphics linking is enabled via koda.json (\"graphics\": true).")
		fmt.Println("Import the full wrapper in your .koda file:")
		fmt.Println(`  use raylib;`)
		fmt.Println()
		fmt.Println("548 functions — see: koda doc wrapper @raylib")
	}
	fmt.Println()

	raylibOK, detail := detectRaylibForRoot(cfgRoot)
	if raylibOK {
		fmt.Printf("Raylib library: %s\n", detail)
		fmt.Println("Run: koda run")
		return nil
	}
	if !raylibOK {
		fmt.Println("Raylib library: not found yet (console programs still work).")
		fmt.Println()
		fmt.Println("Next steps for windowed graphics:")
		if runtime.GOOS == "windows" {
			fmt.Println("  - Install raylib for MinGW, or build vendored sources:")
		}
		fmt.Println("  - From repo root: make -C third_party/raylib_static")
		fmt.Println("  - Or set KODA_LINKFLAGS with your platform raylib flags")
		fmt.Println("  - Run: koda doctor")
	}
	return nil
}

func stripRaylibFullSources(sources []string) []string {
	out := make([]string, 0, len(sources))
	for _, s := range sources {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		slash := filepath.ToSlash(s)
		if strings.Contains(slash, "wrappers/raylib/wrapper.c") {
			continue
		}
		out = append(out, s)
	}
	return out
}

func ensureRaylibShim(root string) error {
	dest := filepath.Join(root, "wrappers", "raylib_shim")
	return project.SyncRaylibShim(dest)
}

func detectRaylibForRoot(root string) (bool, string) {
	if stage, arch, ok := nativebuild.VendoredRaylibStatic(root); ok {
		return true, fmt.Sprintf("vendored stage %s (%s)", stage, filepath.Base(arch))
	}
	if inst, err := kodahome.InstallDir(); err == nil {
		if stage, arch, ok := nativebuild.VendoredRaylibStatic(inst); ok {
			return true, fmt.Sprintf("SDK stage %s (%s)", stage, filepath.Base(arch))
		}
	}
	if os.Getenv("KODA_LINKFLAGS") != "" {
		return true, "KODA_LINKFLAGS set in environment"
	}
	return false, "not found"
}

func mergeUnique(base []string, add ...string) []string {
	seen := make(map[string]bool, len(base)+len(add))
	out := make([]string, 0, len(base)+len(add))
	for _, s := range base {
		s = strings.TrimSpace(s)
		if s == "" || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	for _, s := range add {
		s = strings.TrimSpace(s)
		if s == "" || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}
