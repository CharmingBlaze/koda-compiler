package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"koda/internal/kodahome"
)

func runEnv(args []string) error {
	export := false
	for _, a := range args {
		if a == "--export" {
			export = true
		} else if strings.HasPrefix(a, "-") {
			return fmt.Errorf("unknown flag: %s", a)
		}
	}
	emit := func(key, val string) {
		if export {
			fmt.Printf("export %s=%q\n", key, val)
		} else {
			fmt.Printf("%s=%s\n", key, val)
		}
	}
	install, _ := kodahome.InstallDir()
	stdlib, _ := kodahome.StdlibDir()
	wrap, _ := kodahome.WrappersDir()
	emit("KODA_VERSION", version)
	emit("KODA_INSTALL_DIR", install)
	emit("KODA_STDLIB_DIR", stdlib)
	emit("KODA_WRAPPERS_DIR", wrap)
	for _, k := range []string{
		"KODA_CLANG", "KODA_LLC", "CC", "KODA_PATH", "KODA_WRAPPERS",
		"KODA_NATIVE_SOURCES", "KODA_LINKFLAGS", "KODA_USE_LLD",
		"KODA_SKIP_TOOLCHAIN_EXTRACT", "KODA_DEBUG_IR", "KODA_BUNDLE_FILES",
		"KODA_RAYLIB_STAGE", "KODA_USE_VENDORED_RAYLIB",
	} {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" {
			emit(k, v)
		}
	}
	return nil
}

func runUpdate(args []string) error {
	checkOnly := false
	for _, a := range args {
		switch a {
		case "--check-only":
			checkOnly = true
		default:
			if strings.HasPrefix(a, "-") {
				return fmt.Errorf("unknown flag: %s", a)
			}
		}
	}
	fmt.Printf("current: %s\n", versionLine())
	if checkOnly {
		fmt.Println("Visit https://github.com/CharmingBlaze/koda-compiler/releases for newer SDK zips.")
		return nil
	}
	fmt.Println("Download the latest SDK from:")
	fmt.Println("  https://github.com/CharmingBlaze/koda-compiler/releases")
	fmt.Println("Replace koda, kodawrap, stdlib/, and docs/ with the new zip contents.")
	return nil
}

func runDoc(args []string) error {
	if len(args) == 0 {
		return printDocIndex()
	}
	switch args[0] {
	case "stdlib":
		return printStdlibDocIndex()
	case "module":
		if len(args) < 2 {
			return fmt.Errorf("usage: koda doc module <@name>")
		}
		name := strings.TrimPrefix(args[1], "@")
		return printModuleDoc(name)
	case "wrappers":
		return printWrappersDocIndex()
	case "wrapper":
		if len(args) < 2 {
			return fmt.Errorf("usage: koda doc wrapper <@name>")
		}
		name := strings.TrimPrefix(args[1], "@")
		return printWrapperDoc(name)
	case "--help", "-h":
		printCommandHelp("doc")
		return nil
	default:
		return fmt.Errorf("unknown doc subcommand %q (try: stdlib, module, wrappers, wrapper)", args[0])
	}
}

func printDocIndex() error {
	install, err := kodahome.InstallDir()
	if err != nil {
		return err
	}
	docs := filepath.Join(filepath.Dir(install), "docs")
	if fi, err := os.Stat(docs); err != nil || !fi.IsDir() {
		docs = filepath.Join(filepath.Dir(filepath.Dir(install)), "docs")
	}
	fmt.Println("Koda documentation paths:")
	fmt.Printf("  docs hub:     %s\n", filepath.Join(docs, "README.md"))
	fmt.Printf("  beginners:    %s\n", filepath.Join(docs, "beginners-guide.md"))
	fmt.Printf("  learn:        %s\n", filepath.Join(docs, "learn"))
	fmt.Printf("  stdlib:       %s\n", filepath.Join(docs, "stdlib"))
	fmt.Println("\nCommands:")
	fmt.Println("  koda doc stdlib")
	fmt.Println("  koda doc module @math")
	fmt.Println("  koda doc wrappers")
	fmt.Println("  koda doc wrapper @raylib")
	return nil
}

func printStdlibDocIndex() error {
	stdlib, err := kodahome.StdlibDir()
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(stdlib)
	if err != nil {
		return err
	}
	fmt.Println("stdlib modules (@import):")
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".koda") {
			continue
		}
		base := strings.TrimSuffix(e.Name(), ".koda")
		fmt.Printf("  @%s  %s\n", base, filepath.Join(stdlib, e.Name()))
	}
	return nil
}

func printModuleDoc(name string) error {
	stdlib, err := kodahome.StdlibDir()
	if err != nil {
		return err
	}
	path := filepath.Join(stdlib, name+".koda")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("module @%s: %w", name, err)
	}
	fmt.Printf("# @%s (%s)\n\n", name, path)
	fmt.Print(string(data))
	return nil
}

type wrapperMETA struct {
	Name         string            `json:"name"`
	Generator    string            `json:"generator"`
	Version      string            `json:"version"`
	Import       string            `json:"import"`
	Linkflags    string            `json:"linkflags"`
	PrimaryHeader string           `json:"primary_header"`
	Counts       map[string]int    `json:"counts"`
	Docs         map[string]string `json:"docs"`
}

func wrapperSearchRoots() []string {
	seen := make(map[string]bool)
	var roots []string
	add := func(p string) {
		p = strings.TrimSpace(p)
		if p == "" || seen[p] {
			return
		}
		if fi, err := os.Stat(p); err != nil || !fi.IsDir() {
			return
		}
		seen[p] = true
		roots = append(roots, p)
	}
	if v := strings.TrimSpace(os.Getenv("KODA_WRAPPERS")); v != "" {
		for _, part := range strings.Split(v, string(os.PathListSeparator)) {
			add(part)
		}
	}
	if d, err := kodahome.WrappersDir(); err == nil {
		add(d)
	}
	return roots
}

func findWrapperDir(name string) (string, error) {
	for _, root := range wrapperSearchRoots() {
		dir := filepath.Join(root, name)
		if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
			return dir, nil
		}
	}
	return "", fmt.Errorf("wrapper @%s not found (set KODA_WRAPPERS or install SDK wrappers/)", name)
}

func printWrappersDocIndex() error {
	roots := wrapperSearchRoots()
	if len(roots) == 0 {
		return fmt.Errorf("no wrapper roots (set KODA_WRAPPERS or install wrappers/ next to koda)")
	}
	fmt.Println("C library wrappers (@import / #include):")
	for _, root := range roots {
		entries, err := os.ReadDir(root)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			name := e.Name()
			dir := filepath.Join(root, name)
			line := fmt.Sprintf("  @%s  %s", name, dir)
			metaPath := filepath.Join(dir, "META.json")
			if data, err := os.ReadFile(metaPath); err == nil {
				var meta wrapperMETA
				if json.Unmarshal(data, &meta) == nil && len(meta.Counts) > 0 {
					line += fmt.Sprintf("  (%d functions)", meta.Counts["functions"])
				}
			}
			fmt.Println(line)
		}
	}
	return nil
}

func printWrapperDoc(name string) error {
	dir, err := findWrapperDir(name)
	if err != nil {
		return err
	}
	metaPath := filepath.Join(dir, "META.json")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return fmt.Errorf("wrapper @%s: no META.json in %s", name, dir)
	}
	var meta wrapperMETA
	if err := json.Unmarshal(data, &meta); err != nil {
		return fmt.Errorf("wrapper @%s: invalid META.json: %w", name, err)
	}
	display := meta.Name
	if display == "" {
		display = name
	}
	fmt.Printf("# @%s\n\n", display)
	if meta.Generator != "" || meta.Version != "" {
		fmt.Printf("Generated by **%s** %s\n\n", meta.Generator, meta.Version)
	}
	if meta.Import != "" {
		fmt.Printf("Import: `%s` or `#include \"%s\"`\n\n", meta.Import, meta.Import)
	}
	if meta.PrimaryHeader != "" {
		fmt.Printf("Header: %s\n\n", meta.PrimaryHeader)
	}
	if len(meta.Counts) > 0 {
		fmt.Printf("API: %d functions, %d structs, %d enums\n\n",
			meta.Counts["functions"], meta.Counts["structs"], meta.Counts["enums"])
	}
	if meta.Linkflags != "" {
		fmt.Printf("Link flags: %s\n\n", meta.Linkflags)
	}
	fmt.Printf("Package: %s\n", dir)
	fmt.Printf("Glue:    %s\n", filepath.Join(dir, "wrapper.c"))
	if len(meta.Docs) > 0 {
		fmt.Println("\nDocumentation:")
		keys := []string{"html", "readme", "api_reference", "examples"}
		for _, k := range keys {
			rel, ok := meta.Docs[k]
			if !ok || rel == "" {
				continue
			}
			fmt.Printf("  %s: %s\n", k, filepath.Join(dir, rel))
		}
	}
	kodaJSON := filepath.Join(dir, "koda.json")
	if _, err := os.Stat(kodaJSON); err == nil {
		fmt.Printf("\nManifest fragment: %s\n", kodaJSON)
	}
	return nil
}
