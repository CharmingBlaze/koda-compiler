package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"koda/internal/kodahome"
	"koda/internal/project"
)

// SDKLine is one row in the IDE SDK health panel.
type SDKLine struct {
	OK     bool   `json:"ok"`
	Label  string `json:"label"`
	Detail string `json:"detail"`
	Fix    string `json:"fix"`
}

// SDKStatus summarizes whether Koda Studio can compile and run programs.
type SDKStatus struct {
	OK         bool      `json:"ok"`
	Version    string    `json:"version"`
	InstallDir string    `json:"installDir"`
	StdlibDir  string    `json:"stdlibDir"`
	Lines      []SDKLine `json:"lines"`
}

// BootstrapSDKRoot sets KODA_HOME when unset and stdlib/ is found near exeDir.
func BootstrapSDKRoot(exeDir string) {
	kodahome.BootstrapSDKRoot(exeDir)
}

// EnsureSDKFromExecutable locates stdlib/ relative to the running binary (for IDE hosts).
func EnsureSDKFromExecutable() {
	kodahome.EnsureSDKFromExecutable()
}

// CheckSDK returns a lightweight health report for the welcome screen.
func CheckSDK(studioVersion string) SDKStatus {
	st := SDKStatus{Version: studioVersion, Lines: []SDKLine{}}
	failures := 0

	add := func(ok bool, label, detail, fix string) {
		st.Lines = append(st.Lines, SDKLine{OK: ok, Label: label, Detail: detail, Fix: fix})
		if !ok {
			failures++
		}
	}

	install, err := kodahome.InstallDir()
	if err != nil {
		add(false, "SDK folder", err.Error(), "Place Koda Studio next to koda.exe and stdlib/ (SDK zip layout).")
	} else {
		st.InstallDir = install
		add(true, "SDK folder", install, "")
	}

	stdlib, err := kodahome.StdlibDir()
	if err != nil {
		add(false, "stdlib", err.Error(), "Unzip the full SDK so stdlib/ sits next to Koda Studio.")
	} else {
		st.StdlibDir = stdlib
		if fi, err := os.Stat(filepath.Join(stdlib, "math.koda")); err != nil || fi.IsDir() {
			add(false, "stdlib", stdlib, "Missing math.koda — download the SDK zip from GitHub Releases.")
		} else {
			add(true, "stdlib", stdlib, "")
		}
	}

	if _, err := kodahome.FindToolchain(); err != nil {
		add(false, "compiler runtime", err.Error(), "Use a release SDK or build runtime/libkoda_runtime.a from source.")
	} else {
		add(true, "compiler runtime", "ready", "")
	}

	st.OK = failures == 0
	return st
}

// ListProjectTemplates returns hello, game, graphics, … for the new-project wizard.
func ListProjectTemplates() []string {
	return project.ListTemplates()
}

// ScaffoldProject creates parentDir/name from an embedded template.
func ScaffoldProject(parentDir, name, template string) (string, error) {
	return project.Scaffold(parentDir, name, template)
}

// ExampleEntry is one runnable SDK sample for Koda Studio.
type ExampleEntry struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Blurb    string `json:"blurb"`
	Category string `json:"category"`
}

type exampleManifest struct {
	Name   string `json:"name"`
	Studio *struct {
		Title    string `json:"title"`
		Blurb    string `json:"blurb"`
		Category string `json:"category"`
	} `json:"studio"`
}

// ExampleGamePath returns the absolute path to examples/<id> under KODA_HOME.
// id uses forward slashes, e.g. "games/fps-arena" or "hello_raylib_raw".
func ExampleGamePath(name string) (string, error) {
	install, err := kodahome.InstallDir()
	if err != nil {
		return "", err
	}
	clean := strings.Trim(filepath.ToSlash(strings.TrimSpace(name)), "/")
	if clean == "" || strings.Contains(clean, "..") {
		return "", os.ErrNotExist
	}
	p := filepath.Join(install, "examples", filepath.FromSlash(clean))
	fi, err := os.Stat(p)
	if err != nil {
		return "", err
	}
	if !fi.IsDir() {
		return "", os.ErrNotExist
	}
	return filepath.Abs(p)
}

// ListExamples discovers runnable example projects under SDK examples/.
func ListExamples() ([]ExampleEntry, error) {
	install, err := kodahome.InstallDir()
	if err != nil {
		return nil, err
	}
	root := filepath.Join(install, "examples")
	var out []ExampleEntry
	err = filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() != "koda.json" {
			return nil
		}
		dir := filepath.Dir(path)
		rel, err := filepath.Rel(root, dir)
		if err != nil {
			return nil
		}
		id := filepath.ToSlash(rel)
		b, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		var man exampleManifest
		_ = json.Unmarshal(b, &man)
		title := humanizeExampleID(id)
		blurb := "Press F5 to run."
		category := "game"
		if man.Studio != nil {
			if t := strings.TrimSpace(man.Studio.Title); t != "" {
				title = t
			}
			if b := strings.TrimSpace(man.Studio.Blurb); b != "" {
				blurb = b
			}
			if c := strings.TrimSpace(man.Studio.Category); c != "" {
				category = c
			}
		} else if strings.HasPrefix(id, "hello") || id == "hello-use-module" {
			category = "language"
		} else if !strings.HasPrefix(id, "games/") {
			category = "graphics"
		}
		out = append(out, ExampleEntry{ID: id, Title: title, Blurb: blurb, Category: category})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(out, func(i, j int) bool {
		ci, cj := exampleCategoryOrder(out[i].Category), exampleCategoryOrder(out[j].Category)
		if ci != cj {
			return ci < cj
		}
		return strings.ToLower(out[i].Title) < strings.ToLower(out[j].Title)
	})
	return out, nil
}

func exampleCategoryOrder(cat string) int {
	switch strings.ToLower(cat) {
	case "game":
		return 0
	case "graphics":
		return 1
	case "language", "console":
		return 2
	default:
		return 3
	}
}

func humanizeExampleID(id string) string {
	base := filepath.Base(id)
	parts := strings.Fields(strings.NewReplacer("-", " ", "_", " ").Replace(base))
	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, " ")
}
