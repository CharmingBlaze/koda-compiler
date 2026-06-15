package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"koda/internal/project"
)

//go:embed templates/*
var templateFS embed.FS

func runNew(args []string) error {
	name := ""
	template := "hello"
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--template", "-t":
			if i+1 >= len(args) {
				return fmt.Errorf("--template requires a name")
			}
			i++
			template = args[i]
		default:
			if strings.HasPrefix(args[i], "-") {
				return fmt.Errorf("unknown flag: %s", args[i])
			}
			if name != "" {
				return fmt.Errorf("multiple project names")
			}
			name = args[i]
		}
	}
	if name == "" {
		return fmt.Errorf("usage: koda new <name> [--template hello|game|graphics]")
	}
	if err := validateProjectName(name); err != nil {
		return err
	}

	dest := filepath.Join(".", name)
	if _, err := os.Stat(dest); err == nil {
		return fmt.Errorf("directory already exists: %s", dest)
	} else if !os.IsNotExist(err) {
		return err
	}

	tplRoot := "templates/" + template
	if _, err := fs.Stat(templateFS, tplRoot); err != nil {
		avail := strings.Join(listTemplates(), ", ")
		return fmt.Errorf("unknown template %q (available: %s)", template, avail)
	}

	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}
	if err := copyTemplate(tplRoot, dest, name); err != nil {
		_ = os.RemoveAll(dest)
		return err
	}
	if err := os.MkdirAll(filepath.Join(dest, "assets"), 0755); err != nil {
		_ = os.RemoveAll(dest)
		return err
	}

	fmt.Printf("Created Koda project %s (template: %s)\n", name, template)
	fmt.Printf("  %s\n", filepath.Join(name, project.FileName))
	fmt.Printf("  %s\n", filepath.Join(name, "src", "main.koda"))
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", name)
	switch template {
	case "graphics":
		fmt.Println("  koda setup raylib   # configure koda.json graphics linking")
		fmt.Println("  koda run")
	case "game":
		fmt.Println("  koda run    # text lunar lander — no extra libraries")
	default:
		fmt.Println("  koda run")
	}
	return nil
}

func listTemplates() []string {
	entries, err := fs.ReadDir(templateFS, "templates")
	if err != nil {
		return []string{"hello"}
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return names
}

func validateProjectName(name string) error {
	if name == "" || name == "." || name == ".." {
		return fmt.Errorf("invalid project name: %s", name)
	}
	if strings.ContainsAny(name, `/\`) {
		return fmt.Errorf("project name must not contain path separators")
	}
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			continue
		}
		return fmt.Errorf("project name %q: use letters, digits, hyphen, underscore", name)
	}
	return nil
}

func copyTemplate(tplRoot, dest, name string) error {
	return fs.WalkDir(templateFS, tplRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(filepath.FromSlash(tplRoot), filepath.FromSlash(path))
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		target := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, err := templateFS.ReadFile(path)
		if err != nil {
			return err
		}
		text := strings.ReplaceAll(string(data), "{{name}}", name)
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		return os.WriteFile(target, []byte(text), 0644)
	})
}
