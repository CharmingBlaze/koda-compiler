package wrapcatalog

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"koda/internal/kodahome"
)

//go:embed catalog.json
var catalogJSON []byte

type VersionRecipe struct {
	Version        string   `json:"version"`
	LibraryVersion string   `json:"library_version"`
	Graphics       bool     `json:"graphics"`
	Prebuilt       string   `json:"prebuilt"`
	Headers        []string `json:"headers"`
	IncludePaths   []string `json:"include_paths"`
	LinkFlags      []string `json:"link_flags"`
}

type LibraryRecipe struct {
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	DefaultVersion string          `json:"default_version"`
	Versions       []VersionRecipe `json:"versions"`
}

type Catalog struct {
	Libraries []LibraryRecipe `json:"libraries"`
}

type InstallSpec struct {
	Name    string
	Version string
}

func Load() (Catalog, error) {
	var c Catalog
	if err := json.Unmarshal(catalogJSON, &c); err != nil {
		return c, err
	}
	return c, nil
}

func ParseInstallSpec(s string) (InstallSpec, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return InstallSpec{}, fmt.Errorf("empty library name")
	}
	if strings.HasPrefix(s, "@") {
		s = s[1:]
	}
	spec := InstallSpec{}
	if i := strings.Index(s, "@"); i >= 0 {
		spec.Name = strings.TrimSpace(s[:i])
		spec.Version = strings.TrimSpace(s[i+1:])
	} else {
		spec.Name = s
	}
	if spec.Name == "" {
		return InstallSpec{}, fmt.Errorf("invalid install spec %q", s)
	}
	return spec, nil
}

func (c Catalog) Find(name, version string) (LibraryRecipe, VersionRecipe, error) {
	name = strings.ToLower(strings.TrimSpace(name))
	for _, lib := range c.Libraries {
		if strings.ToLower(lib.Name) != name {
			continue
		}
		if version == "" {
			version = lib.DefaultVersion
		}
		for _, v := range lib.Versions {
			if v.Version == version {
				return lib, v, nil
			}
		}
		avail := make([]string, 0, len(lib.Versions))
		for _, v := range lib.Versions {
			avail = append(avail, v.Version)
		}
		return lib, VersionRecipe{}, fmt.Errorf("unknown version %q for %s (available: %s)", version, lib.Name, strings.Join(avail, ", "))
	}
	names := make([]string, 0, len(c.Libraries))
	for _, lib := range c.Libraries {
		names = append(names, lib.Name)
	}
	return LibraryRecipe{}, VersionRecipe{}, fmt.Errorf("unknown library %q (catalog: %s)", name, strings.Join(names, ", "))
}

type ExpandContext struct {
	SDK         string
	ProjectRoot string
}

func (ctx ExpandContext) Expand(path string) string {
	path = strings.TrimSpace(path)
	path = strings.ReplaceAll(path, "{sdk}", ctx.SDK)
	path = strings.ReplaceAll(path, "{root}", ctx.ProjectRoot)
	if home := os.Getenv("USERPROFILE"); home != "" {
		path = strings.ReplaceAll(path, "{home}", home)
	} else if home := os.Getenv("HOME"); home != "" {
		path = strings.ReplaceAll(path, "{home}", home)
	}
	return filepath.FromSlash(path)
}

func (ctx ExpandContext) ExpandAll(paths []string) []string {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		p = ctx.Expand(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func DefaultExpandContext(projectRoot string) (ExpandContext, error) {
	sdk, err := kodahome.InstallDir()
	if err != nil {
		sdk = ""
	}
	if projectRoot == "" {
		projectRoot, _ = os.Getwd()
	}
	return ExpandContext{SDK: sdk, ProjectRoot: projectRoot}, nil
}

// ResolveHeader returns the first existing header from candidates.
func ResolveHeader(candidates []string) (string, error) {
	for _, p := range candidates {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p, nil
		}
	}
	return "", fmt.Errorf("no header found; tried:\n  %s", strings.Join(candidates, "\n  "))
}

func (v VersionRecipe) ResolvedIncludePaths(ctx ExpandContext, headerPath string) []string {
	paths := ctx.ExpandAll(v.IncludePaths)
	if len(paths) > 0 {
		return paths
	}
	if headerPath != "" {
		return []string{filepath.Dir(headerPath)}
	}
	return nil
}
