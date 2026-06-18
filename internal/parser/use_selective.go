package parser

import (
	"fmt"
	"strings"
)

// FilterModuleDecls exports selective import filtering for tests.
func FilterModuleDecls(decls []Decl, names []string, modulePath string) ([]Decl, error) {
	return filterModuleDecls(decls, names, modulePath)
}

// filterModuleDecls keeps only module lets/funcs whose names appear in names (case-insensitive).
func filterModuleDecls(decls []Decl, names []string, modulePath string) ([]Decl, error) {
	if len(names) == 0 {
		return decls, nil
	}
	want := make(map[string]bool, len(names))
	for _, n := range names {
		want[strings.ToLower(strings.TrimSpace(n))] = true
	}
	found := make(map[string]bool, len(names))
	var out []Decl
	for _, d := range decls {
		switch x := d.(type) {
		case *LetDecl:
			key := strings.ToLower(x.Name.Lexeme)
			if want[key] {
				out = append(out, d)
				found[key] = true
			}
		case *FuncDecl:
			if strings.EqualFold(x.Name.Lexeme, "main") {
				continue
			}
			key := strings.ToLower(x.Name.Lexeme)
			if want[key] {
				out = append(out, d)
				found[key] = true
			}
		}
	}
	for _, n := range names {
		key := strings.ToLower(strings.TrimSpace(n))
		if key == "" || found[key] {
			continue
		}
		return nil, fmt.Errorf("use %q: unknown import %q (not exported from module)", modulePath, n)
	}
	return out, nil
}

func expandUseDecl(modulePath string, decl *UseDecl, bundle *ProgramBundle, stack map[string]bool) ([]Decl, error) {
	abs, err := resolveModuleRef(modulePath, useImportRef(decl.ModulePath))
	if err != nil {
		return nil, fmt.Errorf("%s: use %q: %w", modulePath, decl.ModulePath, err)
	}
	inner, err := expandIncludes(abs, bundle, stack)
	if err != nil {
		return nil, err
	}
	if bundle != nil && strings.EqualFold(decl.ModulePath, "raylib") {
		bundle.RaylibImported = true
	}
	if len(decl.Selective) > 0 {
		inner, err = filterModuleDecls(inner, decl.Selective, decl.ModulePath)
		if err != nil {
			return nil, err
		}
	}
	if decl.Alias != "" {
		wrapped, err := wrapModuleAsAlias(inner, decl.Alias)
		if err != nil {
			return nil, err
		}
		if bundle != nil && strings.EqualFold(decl.ModulePath, "raylib") {
			bundle.RaylibAlias = decl.Alias
		}
		return wrapped, nil
	}
	return inner, nil
}
