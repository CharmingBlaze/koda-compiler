package nativebuild

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"koda/internal/kodahome"
	"koda/internal/parser"
	"koda/internal/project"
)

type nativeGlueKind int

const (
	nativeGlueNone nativeGlueKind = iota
	nativeGlueShim
	nativeGlueFullRaylib
)

// ApplyNativeSourcesForBundle sets KODA_NATIVE_SOURCES (and link flags when needed) from
// koda:extern symbols in the program when the project manifest did not already configure them.
func ApplyNativeSourcesForBundle(bundle *parser.ProgramBundle, projectRoot string) error {
	kind := inferNativeGlueKind(bundle)
	if kind == nativeGlueNone {
		return nil
	}
	want, ok := resolveNativeWrapperPath(kind, projectRoot)
	if !ok {
		return fmt.Errorf("program needs raylib native glue but wrapper.c was not found; run: koda setup raylib (shim) or koda setup raylib --full")
	}
	cur := strings.TrimSpace(os.Getenv("KODA_NATIVE_SOURCES"))
	if cur != "" && nativeSourcesMatchKind(cur, kind) {
		return nil
	}
	if err := os.Setenv("KODA_NATIVE_SOURCES", want); err != nil {
		return err
	}
	if strings.TrimSpace(os.Getenv("KODA_LINKFLAGS")) == "" {
		if err := os.Setenv("KODA_LINKFLAGS", project.DefaultGraphicsLinkFlags()); err != nil {
			return err
		}
	}
	return nil
}

func inferNativeGlueKind(bundle *parser.ProgramBundle) nativeGlueKind {
	if bundle == nil {
		return nativeGlueNone
	}
	hasShim, hasWrap := false, false
	for path, prog := range bundle.Modules {
		if strings.Contains(filepath.ToSlash(path), "raylib_shim") {
			hasShim = true
		}
		for _, decl := range prog.Declarations {
			native := nativeDirectiveFromDecl(decl)
			if native == nil {
				continue
			}
			sym := native.Symbol
			if strings.HasPrefix(sym, "koda_shim_") {
				hasShim = true
			}
			if strings.HasPrefix(sym, "koda_wrap_raylib_") {
				hasWrap = true
			}
		}
	}
	if hasShim {
		return nativeGlueShim
	}
	if hasWrap {
		return nativeGlueFullRaylib
	}
	return nativeGlueNone
}

func nativeDirectiveFromDecl(decl parser.Decl) *parser.NativeDirective {
	switch d := decl.(type) {
	case *parser.LetDecl:
		return d.Native
	case *parser.FuncDecl:
		return d.Native
	default:
		return nil
	}
}

func resolveNativeWrapperPath(kind nativeGlueKind, projectRoot string) (string, bool) {
	var rel string
	switch kind {
	case nativeGlueShim:
		rel = filepath.Join("wrappers", "raylib_shim", "wrapper.c")
	case nativeGlueFullRaylib:
		rel = filepath.Join("wrappers", "raylib", "wrapper.c")
	default:
		return "", false
	}
	for _, root := range nativeWrapperRoots(projectRoot) {
		p := filepath.Join(root, rel)
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p, true
		}
	}
	return "", false
}

func nativeWrapperRoots(projectRoot string) []string {
	seen := map[string]bool{}
	var roots []string
	add := func(p string) {
		p = strings.TrimSpace(p)
		if p == "" || seen[p] {
			return
		}
		seen[p] = true
		roots = append(roots, p)
	}
	add(projectRoot)
	if inst, err := kodahome.InstallDir(); err == nil {
		add(inst)
	}
	return roots
}

func nativeSourcesMatchKind(sources string, kind nativeGlueKind) bool {
	for _, p := range splitNativeSources(sources) {
		if nativePathMatchesKind(p, kind) {
			return true
		}
	}
	return false
}

func nativePathMatchesKind(path string, kind nativeGlueKind) bool {
	slash := filepath.ToSlash(path)
	switch kind {
	case nativeGlueShim:
		return strings.Contains(slash, "raylib_shim") && strings.HasSuffix(slash, "wrapper.c")
	case nativeGlueFullRaylib:
		return strings.Contains(slash, "wrappers/raylib/") && strings.HasSuffix(slash, "wrapper.c")
	default:
		return false
	}
}

func splitNativeSources(sources string) []string {
	if sources == "" {
		return nil
	}
	return strings.FieldsFunc(sources, func(r rune) bool {
		return r == filepath.ListSeparator || r == ';' || r == ' '
	})
}
