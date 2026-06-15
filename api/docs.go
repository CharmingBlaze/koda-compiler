package api

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"koda/internal/kodahome"
)

// DocPage is one markdown help page shipped with the SDK.
type DocPage struct {
	Rel      string `json:"rel"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Beginner bool   `json:"beginner"`
}

var docSkipNames = map[string]bool{
	"MASTER_PLAN.md":              true,
	"VM_RETIREMENT_CHECKLIST.md":  true,
	"VM_RETIREMENT_AUDIT.md":      true,
	"GITHUB_ISSUES_VM_RETIREMENT.md": true,
	"REORGANIZATION.md":           true,
	"handoff.md":                  true,
	"list.md":                     true,
	"status.md":                   true,
	"ROADMAP.md":                  true,
	"positioning.md":              true,
	"STYLE-GUIDE.md":              true,
	"windows-native-toolchain.md": true,
	"release.md":                  true,
	"releasing.md":                true,
	"git-workflow.md":             true,
	"architecture.md":             true,
	"compiler.md":                   true,
	"koda-language-legacy.md":       true,
}

var docRootFiles = []string{
	"START_HERE.md",
	"README.md",
	"language.md",
	"CHANGELOG.md",
}

func docCategory(rel string) string {
	r := strings.ReplaceAll(rel, "\\", "/")
	switch {
	case r == "START_HERE.md" || r == "docs/beginners-guide.md":
		return "Start here"
	case strings.HasPrefix(r, "docs/learn/"):
		return "Learn step-by-step"
	case strings.HasPrefix(r, "docs/guides/"):
		return "Guides"
	case strings.HasPrefix(r, "docs/stdlib/"):
		return "Stdlib modules"
	case strings.HasPrefix(r, "docs/reference/"):
		return "Reference"
	case strings.HasPrefix(r, "docs/concepts/"):
		return "Concepts"
	case r == "language.md" || r == "docs/language.md" || strings.HasPrefix(r, "docs/language/"):
		return "Reference"
	case r == "docs/faq.md" || r == "docs/troubleshooting.md" || r == "docs/glossary.md":
		return "Help & FAQ"
	case strings.HasPrefix(r, "docs/compiler/"):
		return "Advanced"
	case strings.HasPrefix(r, "docs/"):
		return "Documentation"
	default:
		return "Documentation"
	}
}

func docBeginner(rel string) bool {
	r := strings.ReplaceAll(rel, "\\", "/")
	if r == "START_HERE.md" || r == "docs/beginners-guide.md" || r == "docs/faq.md" ||
		r == "docs/troubleshooting.md" || r == "docs/glossary.md" ||
		r == "docs/guides/getting-started.md" || r == "docs/guides/game-dev.md" ||
		r == "docs/guides/raylib.md" || r == "docs/guides/applications.md" {
		return true
	}
	if strings.HasPrefix(r, "docs/learn/") && !strings.HasSuffix(r, "README.md") {
		return true
	}
	if strings.HasPrefix(r, "docs/concepts/") {
		return true
	}
	if strings.HasPrefix(r, "docs/stdlib/") {
		return true
	}
	if r == "docs/reference/builtins.md" || r == "docs/reference/cli.md" {
		return true
	}
	return false
}

func docTitleFromFile(rel, content string) string {
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	base := filepath.Base(rel)
	base = strings.TrimSuffix(base, filepath.Ext(base))
	base = strings.ReplaceAll(base, "-", " ")
	base = strings.ReplaceAll(base, "_", " ")
	if base == "" {
		return rel
	}
	return strings.ToUpper(base[:1]) + base[1:]
}

func resolveDocPath(install, rel string) (string, error) {
	if install == "" {
		return "", fmt.Errorf("SDK not found")
	}
	clean := filepath.Clean(filepath.FromSlash(rel))
	if clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid doc path")
	}
	if filepath.Ext(clean) != ".md" {
		return "", fmt.Errorf("not a markdown file")
	}
	full := filepath.Join(install, clean)
	absInstall, err := filepath.Abs(install)
	if err != nil {
		return "", err
	}
	absFull, err := filepath.Abs(full)
	if err != nil {
		return "", err
	}
	prefix := absInstall + string(os.PathSeparator)
	if !strings.EqualFold(absFull, absInstall) &&
		!strings.HasPrefix(strings.ToLower(absFull), strings.ToLower(prefix)) {
		return "", fmt.Errorf("doc path outside SDK")
	}
	return absFull, nil
}

// ListDocPages returns markdown help shipped with the SDK (docs/ + key root files).
func ListDocPages() ([]DocPage, error) {
	install, err := kodahome.InstallDir()
	if err != nil {
		return nil, err
	}
	var pages []DocPage
	seen := map[string]bool{}

	add := func(rel string) {
		if seen[rel] {
			return
		}
		if docSkipNames[filepath.Base(rel)] {
			return
		}
		full, err := resolveDocPath(install, rel)
		if err != nil {
			return
		}
		if fi, err := os.Stat(full); err != nil || fi.IsDir() {
			return
		}
		b, err := os.ReadFile(full)
		if err != nil {
			return
		}
		seen[rel] = true
		content := string(b)
		pages = append(pages, DocPage{
			Rel:      strings.ReplaceAll(rel, "\\", "/"),
			Title:    docTitleFromFile(rel, content),
			Category: docCategory(rel),
			Beginner: docBeginner(rel),
		})
	}

	for _, rel := range docRootFiles {
		add(rel)
	}

	docsDir := filepath.Join(install, "docs")
	_ = filepath.Walk(docsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}
		rel, err := filepath.Rel(install, path)
		if err != nil {
			return nil
		}
		add(filepath.ToSlash(rel))
		return nil
	})

	sort.Slice(pages, func(i, j int) bool {
		ci, cj := pages[i].Category, pages[j].Category
		if ci != cj {
			return categoryOrder(ci) < categoryOrder(cj)
		}
		if pages[i].Beginner != pages[j].Beginner {
			return pages[i].Beginner
		}
		return strings.ToLower(pages[i].Title) < strings.ToLower(pages[j].Title)
	})
	return pages, nil
}

func categoryOrder(cat string) int {
	order := []string{
		"Start here",
		"Learn step-by-step",
		"Guides",
		"Help & FAQ",
		"Concepts",
		"Stdlib modules",
		"Reference",
		"Documentation",
		"Advanced",
	}
	for i, c := range order {
		if c == cat {
			return i
		}
	}
	return len(order)
}

// ReadDocPage reads one SDK markdown file by relative path (e.g. docs/faq.md).
func ReadDocPage(rel string) (string, error) {
	install, err := kodahome.InstallDir()
	if err != nil {
		return "", err
	}
	full, err := resolveDocPath(install, rel)
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(full)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
