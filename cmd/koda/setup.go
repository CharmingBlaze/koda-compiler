package main

import (
	"fmt"
	"io/fs"
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
		return fmt.Errorf("usage: koda setup <target>\n\nTargets:\n  raylib   Configure graphics linking (koda.json + raylib shim)")
	}
	switch strings.ToLower(args[0]) {
	case "raylib":
		return setupRaylib(args[1:])
	default:
		return fmt.Errorf("unknown setup target %q (available: raylib)", args[0])
	}
}

func setupRaylib(args []string) error {
	root := cwd()
	if len(args) > 0 {
		root = args[0]
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

	if err := ensureRaylibShim(projectRoot); err != nil {
		return err
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
	cfg.Native.Sources = mergeUnique(cfg.Native.Sources, "wrappers/raylib_shim/wrapper.c")
	cfg.Native.Graphics = true

	if err := project.Save(cfgRoot, cfg); err != nil {
		return err
	}

	fmt.Println("Raylib setup complete.")
	fmt.Printf("  project: %s\n", filepath.Join(cfgRoot, project.FileName))
	fmt.Printf("  shim:    %s\n", filepath.Join(projectRoot, "wrappers", "raylib_shim"))
	fmt.Println()
	fmt.Println("Graphics linking is enabled via koda.json (\"graphics\": true).")
	fmt.Println("Include the shim in your .koda file:")
	fmt.Println(`  #include "wrappers/raylib_shim/raylib.koda"`)
	fmt.Println("Or use @game from stdlib after importing graphics helpers.")
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

func ensureRaylibShim(root string) error {
	dest := filepath.Join(root, "wrappers", "raylib_shim")
	tplRoot := "templates/graphics/wrappers/raylib_shim"
	return fs.WalkDir(templateFS, tplRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(filepath.FromSlash(tplRoot), filepath.FromSlash(path))
		if err != nil || rel == "." {
			return err
		}
		out := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(out, 0755)
		}
		if _, err := os.Stat(out); err == nil {
			return nil
		}
		data, err := templateFS.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(out), 0755); err != nil {
			return err
		}
		return os.WriteFile(out, data, 0644)
	})
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
