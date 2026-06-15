package wrappermeta

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const FileName = "META.json"

// PackageMeta is the machine-readable record for a generated wrapper package.
// It stores enough information to regenerate bindings and detect header drift.
type PackageMeta struct {
	Name           string            `json:"name"`
	Generator      string            `json:"generator"`
	Version        string            `json:"version"`
	GeneratedAt    string            `json:"generated_at"`
	PrimaryHeader  string            `json:"primary_header"`
	Headers        []string          `json:"headers"`
	IncludePaths   []string          `json:"include_paths,omitempty"`
	LinkFlags      []string          `json:"link_flags,omitempty"`
	Linkflags      string            `json:"linkflags,omitempty"`
	UseClang       *bool             `json:"use_clang,omitempty"`
	LibraryVersion string            `json:"library_version,omitempty"`
	HeaderHashes   map[string]string `json:"header_hashes,omitempty"`
	Counts         map[string]int    `json:"counts"`
	Import         string            `json:"import"`
	Docs           map[string]string `json:"docs"`
}

// DriftReport compares recorded header hashes with files on disk.
type DriftReport struct {
	Dir       string
	Meta      PackageMeta
	Changed   []string
	Missing   []string
	Unchanged []string
}

func (d DriftReport) Stale() bool {
	return len(d.Changed) > 0 || len(d.Missing) > 0
}

func LoadDir(dir string) (PackageMeta, error) {
	return Load(filepath.Join(dir, FileName))
}

func Load(path string) (PackageMeta, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return PackageMeta{}, err
	}
	var meta PackageMeta
	if err := json.Unmarshal(data, &meta); err != nil {
		return PackageMeta{}, fmt.Errorf("invalid %s: %w", FileName, err)
	}
	if meta.Name == "" {
		return PackageMeta{}, fmt.Errorf("%s: missing name", FileName)
	}
	if len(meta.Headers) == 0 {
		return PackageMeta{}, fmt.Errorf("%s: missing headers (cannot upgrade)", FileName)
	}
	return meta, nil
}

func Write(path string, meta PackageMeta) error {
	b, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0644)
}

// HashHeaders returns sha256 hex prefixes for each header path (key = path as recorded).
func HashHeaders(paths []string) (map[string]string, error) {
	out := make(map[string]string, len(paths))
	for _, p := range paths {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		h, err := hashFile(p)
		if err != nil {
			return nil, fmt.Errorf("hash %s: %w", p, err)
		}
		out[p] = h
	}
	return out, nil
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	sum := sha256.New()
	if _, err := io.Copy(sum, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(sum.Sum(nil))[:16], nil
}

// CheckDrift compares META header_hashes with current files.
func CheckDrift(meta PackageMeta) (DriftReport, error) {
	report := DriftReport{Meta: meta}
	if len(meta.HeaderHashes) == 0 {
		return report, nil
	}
	for _, header := range meta.Headers {
		header = strings.TrimSpace(header)
		if header == "" {
			continue
		}
		cur, err := hashFile(header)
		if err != nil {
			if os.IsNotExist(err) {
				report.Missing = append(report.Missing, header)
				continue
			}
			return report, err
		}
		prev, ok := meta.HeaderHashes[header]
		if !ok {
			report.Changed = append(report.Changed, header)
			continue
		}
		if prev != cur {
			report.Changed = append(report.Changed, header)
		} else {
			report.Unchanged = append(report.Unchanged, header)
		}
	}
	return report, nil
}

// LinkFlagsString joins link flags and -I include paths for koda.json / README.
func LinkFlagsString(includePaths, linkFlags []string) string {
	parts := append([]string{}, linkFlags...)
	for _, d := range includePaths {
		d = strings.TrimSpace(d)
		if d != "" {
			parts = append(parts, "-I"+d)
		}
	}
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, " ")
}

// ParseLinkflags splits a legacy linkflags string into include paths and linker tokens.
func ParseLinkflags(s string) (includePaths, linkFlags []string) {
	for _, tok := range strings.Fields(strings.TrimSpace(s)) {
		if strings.HasPrefix(tok, "-I") && len(tok) > 2 {
			includePaths = append(includePaths, tok[2:])
			continue
		}
		linkFlags = append(linkFlags, tok)
	}
	return includePaths, linkFlags
}

// ResolveIncludePaths returns include paths from meta, inferring from linkflags when needed.
func (m PackageMeta) ResolveIncludePaths() []string {
	if len(m.IncludePaths) > 0 {
		return append([]string{}, m.IncludePaths...)
	}
	inc, _ := ParseLinkflags(m.Linkflags)
	return inc
}

// ResolveLinkFlags returns linker tokens without -I paths.
func (m PackageMeta) ResolveLinkFlags() []string {
	if len(m.LinkFlags) > 0 {
		return append([]string{}, m.LinkFlags...)
	}
	_, link := ParseLinkflags(m.Linkflags)
	return link
}

// UseClangEnabled defaults to true when unset (matches wrapgen default).
func (m PackageMeta) UseClangEnabled() bool {
	if m.UseClang == nil {
		return true
	}
	return *m.UseClang
}
