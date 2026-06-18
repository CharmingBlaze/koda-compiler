package parser

import (
	"fmt"
	"path/filepath"
	"strings"

	"koda/internal/kodahome"
)

const useImportPrefix = "use:"

// ModuleSearchRoots lists directories consulted when resolving `use` paths (for error messages).
type ModuleSearchRoots struct {
	Wrappers []string
	Stdlib   []string
	Packages []string
}

func (r ModuleSearchRoots) displayPaths(modulePath string) []string {
	safe := filepath.FromSlash(strings.ToLower(modulePath))
	base := filepath.Base(safe)
	var out []string
	appendCandidates := func(root string) {
		root = strings.TrimSpace(root)
		if root == "" {
			return
		}
		out = append(out, filepath.ToSlash(filepath.Join(root, safe+".koda")))
		out = append(out, filepath.ToSlash(filepath.Join(root, safe, "index.koda")))
		out = append(out, filepath.ToSlash(filepath.Join(root, safe, base+".koda")))
	}
	for _, root := range r.Wrappers {
		appendCandidates(root)
	}
	for _, root := range r.Stdlib {
		appendCandidates(root)
	}
	for _, root := range r.Packages {
		appendCandidates(root)
	}
	if len(out) == 0 {
		out = append(out,
			filepath.ToSlash(filepath.Join("wrappers", safe)),
			filepath.ToSlash(filepath.Join("stdlib", safe)),
			filepath.ToSlash(filepath.Join("packages", safe)),
		)
	}
	seen := make(map[string]bool)
	uniq := make([]string, 0, len(out))
	for _, p := range out {
		if seen[p] {
			continue
		}
		seen[p] = true
		uniq = append(uniq, p)
	}
	return uniq
}

func moduleSearchRoots(importerPath string) ModuleSearchRoots {
	var r ModuleSearchRoots
	appendUnique := func(list *[]string, paths ...string) {
		for _, p := range paths {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			found := false
			for _, existing := range *list {
				if strings.EqualFold(existing, p) {
					found = true
					break
				}
			}
			if !found {
				*list = append(*list, p)
			}
		}
	}

	appendUnique(&r.Wrappers, pathListEnv("KODA_WRAPPERS")...)
	appendUnique(&r.Stdlib, pathListEnv("KODA_PATH")...)

	if inst, err := kodahome.InstallDir(); err == nil {
		appendUnique(&r.Wrappers,
			filepath.Join(inst, "wrappers"),
		)
		appendUnique(&r.Stdlib,
			filepath.Join(inst, "stdlib"),
			filepath.Join(inst, "lib"),
		)
		appendUnique(&r.Packages,
			filepath.Join(inst, "packages"),
		)
		parent := filepath.Dir(inst)
		appendUnique(&r.Wrappers, filepath.Join(parent, "wrappers"))
		appendUnique(&r.Stdlib, filepath.Join(parent, "stdlib"))
		appendUnique(&r.Packages, filepath.Join(parent, "packages"))
	}

	if importerPath != "" {
		dir := filepath.Dir(importerPath)
		appendUnique(&r.Wrappers, dir)
		appendUnique(&r.Stdlib, dir)
	}

	return r
}

// ResolveUseModulePath resolves `use module.path` to an absolute .koda file path.
func ResolveUseModulePath(importerPath, modulePath string) (string, error) {
	modulePath = strings.TrimSpace(modulePath)
	if modulePath == "" {
		return "", fmt.Errorf("empty use module path")
	}

	roots := moduleSearchRoots(importerPath)

	// Official stdlib namespace: use koda.math → @math
	if strings.HasPrefix(strings.ToLower(modulePath), "koda.") {
		sub := modulePath[len("koda."):]
		if sub == "" {
			return "", unknownModuleError(modulePath, roots)
		}
		if p, err := ResolveImportPath(importerPath, "@"+strings.ToLower(sub)); err == nil {
			return p, nil
		}
		return "", unknownModuleError(modulePath, roots)
	}

	mod := strings.ToLower(modulePath)

	// Wrappers first (raylib, box2d, …)
	if p, ok := tryResolveModuleUnderRoots(mod, roots.Wrappers); ok {
		return p, nil
	}

	// Stdlib shorthand: use math → @math
	if p, err := ResolveImportPath(importerPath, "@"+mod); err == nil {
		return p, nil
	}

	// Plain module under stdlib/ (e.g. use timer)
	if p, ok := tryResolveModuleUnderRoots(mod, roots.Stdlib); ok {
		return p, nil
	}

	return "", unknownModuleError(modulePath, roots)
}

func unknownModuleError(modulePath string, roots ModuleSearchRoots) error {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Unknown module %q.", modulePath))
	b.WriteString("\n\nSearched:")
	for _, p := range roots.displayPaths(modulePath) {
		b.WriteString("\n  ")
		b.WriteString(p)
	}
	return fmt.Errorf("%s", b.String())
}

func isUseImportRef(rel string) bool {
	return strings.HasPrefix(rel, useImportPrefix)
}

func useModuleFromImportRef(rel string) string {
	return strings.TrimPrefix(rel, useImportPrefix)
}

func resolveModuleRef(importerPath, rel string) (string, error) {
	if isUseImportRef(rel) {
		return ResolveUseModulePath(importerPath, useModuleFromImportRef(rel))
	}
	return ResolveImportPath(importerPath, rel)
}

// useImportRef encodes a use declaration for the transitive loader.
func useImportRef(modulePath string) string {
	return useImportPrefix + modulePath
}
