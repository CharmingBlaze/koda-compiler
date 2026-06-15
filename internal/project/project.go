package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const FileName = "koda.json"

// Config is the project manifest at the repository / app root.
type Config struct {
	Name    string       `json:"name"`
	Version string       `json:"version"`
	Entry   string       `json:"entry"`
	Lint    string       `json:"lint"` // "", "beginner", "strict"
	Bundle  BundleConfig `json:"bundle"`
	Native  NativeConfig `json:"native"`
}

type BundleConfig struct {
	Assets []string `json:"assets"`
	Extra  []string `json:"extra"`
}

type NativeConfig struct {
	Sources   []string `json:"sources"`
	Linkflags string   `json:"linkflags"`
	Graphics  bool     `json:"graphics"` // auto platform link flags + raylib when linkflags empty
}

// Load reads koda.json from path (file or directory containing it).
func Load(path string) (*Config, string, error) {
	root, file, err := locateFile(path)
	if err != nil {
		return nil, "", err
	}
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, "", err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, "", fmt.Errorf("%s: %w", file, err)
	}
	if strings.TrimSpace(cfg.Entry) == "" {
		return nil, "", fmt.Errorf("%s: entry is required", file)
	}
	return &cfg, root, nil
}

// Find walks upward from dir looking for koda.json.
func Find(dir string) (*Config, string, error) {
	dir = filepath.Clean(dir)
	for {
		file := filepath.Join(dir, FileName)
		if _, err := os.Stat(file); err == nil {
			return Load(file)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return nil, "", nil
		}
		dir = parent
	}
}

func locateFile(path string) (root, file string, err error) {
	path = filepath.Clean(path)
	if strings.EqualFold(filepath.Base(path), FileName) {
		return filepath.Dir(path), path, nil
	}
	if fi, err := os.Stat(path); err == nil && fi.IsDir() {
		file = filepath.Join(path, FileName)
		if _, err := os.Stat(file); err != nil {
			return "", "", fmt.Errorf("%s not found in %s", FileName, path)
		}
		return path, file, nil
	}
	return "", "", fmt.Errorf("%s not found near %s", FileName, path)
}

// EntryPath returns the absolute entry .koda path for a project root and config.
func EntryPath(root string, cfg *Config) string {
	return filepath.Join(root, filepath.FromSlash(cfg.Entry))
}

// ResolveEntry returns explicit if set, otherwise the entry from koda.json in cwd or parents.
func ResolveEntry(cwd, explicit string) (string, error) {
	if strings.TrimSpace(explicit) != "" {
		return explicit, nil
	}
	cfg, root, err := Find(cwd)
	if err != nil {
		return "", err
	}
	if cfg == nil {
		return "", fmt.Errorf("no %s found (pass a .koda file or run from a project directory)", FileName)
	}
	return EntryPath(root, cfg), nil
}

// Context holds a loaded project and its root directory.
type Context struct {
	Root string
	Cfg  *Config
}

// LoadContext finds koda.json from the entry file path or cwd.
func LoadContext(from string) (*Context, error) {
	from = strings.TrimSpace(from)
	if from != "" {
		abs, err := filepath.Abs(from)
		if err != nil {
			return nil, err
		}
		if strings.EqualFold(filepath.Ext(abs), ".koda") {
			dir := filepath.Dir(abs)
			cfg, root, err := Find(dir)
			if err != nil {
				return nil, err
			}
			if cfg != nil {
				return &Context{Root: root, Cfg: cfg}, nil
			}
			return nil, nil
		}
	}
	cfg, root, err := Find(from)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		return nil, nil
	}
	return &Context{Root: root, Cfg: cfg}, nil
}

// ApplyNativeEnv sets KODA_NATIVE_SOURCES and KODA_LINKFLAGS from the project when unset.
func (c *Context) ApplyNativeEnv() error {
	if c == nil || c.Cfg == nil {
		return nil
	}
	if os.Getenv("KODA_NATIVE_SOURCES") == "" && len(c.Cfg.Native.Sources) > 0 {
		var abs []string
		for _, src := range c.Cfg.Native.Sources {
			abs = append(abs, filepath.Join(c.Root, filepath.FromSlash(src)))
		}
		if err := os.Setenv("KODA_NATIVE_SOURCES", strings.Join(abs, string(os.PathListSeparator))); err != nil {
			return err
		}
	}
	link := strings.TrimSpace(c.Cfg.Native.Linkflags)
	if os.Getenv("KODA_LINKFLAGS") == "" {
		if link != "" {
			if err := os.Setenv("KODA_LINKFLAGS", link); err != nil {
				return err
			}
		} else if c.Cfg.Native.Graphics {
			if err := os.Setenv("KODA_LINKFLAGS", DefaultGraphicsLinkFlags()); err != nil {
				return err
			}
		}
	}
	return nil
}

// BundleExtraPaths returns asset and extra paths relative to root for copying into bundles.
func (c *Context) BundleExtraPaths() []string {
	if c == nil || c.Cfg == nil {
		return nil
	}
	var out []string
	for _, p := range c.Cfg.Bundle.Assets {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, filepath.Join(c.Root, filepath.FromSlash(p)))
		}
	}
	for _, p := range c.Cfg.Bundle.Extra {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, filepath.Join(c.Root, filepath.FromSlash(p)))
		}
	}
	return out
}

// AppName returns the distributable name (project name or entry basename).
func (c *Context) AppName(entryPath string) string {
	if c != nil && c.Cfg != nil && strings.TrimSpace(c.Cfg.Name) != "" {
		return strings.TrimSpace(c.Cfg.Name)
	}
	base := strings.TrimSuffix(filepath.Base(entryPath), filepath.Ext(entryPath))
	if base != "" {
		return base
	}
	return "app"
}

// Save writes koda.json under root (creates the file or overwrites).
func Save(root string, cfg *Config) error {
	root = filepath.Clean(root)
	if cfg == nil {
		return fmt.Errorf("nil project config")
	}
	if strings.TrimSpace(cfg.Entry) == "" {
		return fmt.Errorf("entry is required")
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(filepath.Join(root, FileName), data, 0644)
}
