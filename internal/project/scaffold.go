package project

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

//go:embed templates/*
var templateFS embed.FS

// ListTemplates returns sorted template names (hello, game, graphics, …).
func ListTemplates() []string {
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

// ValidateProjectName checks a single path segment project folder name.
func ValidateProjectName(name string) error {
	if name == "" || name == "." || name == ".." {
		return fmt.Errorf("invalid project name: %s", name)
	}
	if strings.ContainsAny(name, `/\`) {
		return fmt.Errorf("project name must not contain path separators")
	}
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == ' ' {
			continue
		}
		return fmt.Errorf("project name %q: use letters, digits, spaces, hyphen, or underscore", name)
	}
	return nil
}

// Scaffold creates parentDir/name from an embedded template and returns the new project root.
func Scaffold(parentDir, name, template string) (string, error) {
	name = strings.TrimSpace(name)
	if err := ValidateProjectName(name); err != nil {
		return "", err
	}
	template = strings.TrimSpace(template)
	if template == "" {
		template = "hello"
	}
	parent, err := filepath.Abs(parentDir)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(parent)
	if err != nil {
		return "", err
	}
	if !fi.IsDir() {
		return "", fmt.Errorf("parent is not a directory")
	}
	dest := filepath.Join(parent, name)
	if _, err := os.Stat(dest); err == nil {
		return "", fmt.Errorf("folder already exists: %s", name)
	} else if !os.IsNotExist(err) {
		return "", err
	}
	tplRoot := "templates/" + template
	if _, err := fs.Stat(templateFS, tplRoot); err != nil {
		avail := strings.Join(ListTemplates(), ", ")
		return "", fmt.Errorf("unknown template %q (available: %s)", template, avail)
	}
	if err := copyTemplate(tplRoot, dest, name); err != nil {
		_ = os.RemoveAll(dest)
		return "", err
	}
	if template == "graphics" || template == "pong" {
		// Full Raylib wrapper resolves from the SDK via koda.json native.sources.
	}
	if err := os.MkdirAll(filepath.Join(dest, "assets"), 0755); err != nil {
		_ = os.RemoveAll(dest)
		return "", err
	}
	return filepath.Abs(dest)
}

// CopyTemplateSubtree copies templates/<template>/<rel> into dest (existing files are skipped).
func CopyTemplateSubtree(template, rel, dest string) error {
	template = strings.TrimSpace(template)
	rel = strings.TrimSpace(rel)
	if template == "" {
		return fmt.Errorf("template name required")
	}
	tplRoot := filepath.ToSlash(filepath.Join("templates", template, rel))
	if _, err := fs.Stat(templateFS, tplRoot); err != nil {
		return fmt.Errorf("template path %q: %w", tplRoot, err)
	}
	return fs.WalkDir(templateFS, tplRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		sub, err := filepath.Rel(filepath.FromSlash(tplRoot), filepath.FromSlash(path))
		if err != nil || sub == "." {
			return err
		}
		target := filepath.Join(dest, sub)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		if _, err := os.Stat(target); err == nil {
			return nil
		}
		data, err := templateFS.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		return os.WriteFile(target, data, 0644)
	})
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
