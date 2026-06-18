package parser

import (
	"fmt"
	"strings"
)

// FlattenEntryIncludes expands top-level #include directives in the entry program
// into one merged Program. Requires bundle.Modules to contain every included file
// (see findImports + loadModule in loader.go).
func FlattenEntryIncludes(bundle *ProgramBundle) error {
	if bundle == nil || bundle.Entry == nil {
		return fmt.Errorf("invalid program bundle")
	}
	if len(bundle.Modules) == 0 {
		// In-memory bundles (e.g. tests) have no file paths; nothing to flatten.
		return nil
	}
	entryPath, err := BundleEntryPath(bundle)
	if err != nil {
		return err
	}
	decls, err := expandIncludes(entryPath, bundle, make(map[string]bool))
	if err != nil {
		return err
	}
	bundle.Entry = &Program{Declarations: decls}
	return nil
}

func expandIncludes(modulePath string, bundle *ProgramBundle, stack map[string]bool) ([]Decl, error) {
	if stack[modulePath] {
		return nil, fmt.Errorf("include cycle involving %q", modulePath)
	}
	stack[modulePath] = true
	defer delete(stack, modulePath)

	prog, modulePath, err := lookupBundleModule(bundle, modulePath)
	if err != nil {
		return nil, err
	}

	var out []Decl
	for _, d := range prog.Declarations {
		var rel string
		switch decl := d.(type) {
		case *IncludeDecl:
			rel = strings.Trim(decl.Path.Lexeme, `"'`)
		case *UseDecl:
			expanded, err := expandUseDecl(modulePath, decl, bundle, stack)
			if err != nil {
				return nil, err
			}
			out = append(out, expanded...)
			continue
		default:
			out = append(out, d)
			continue
		}
		abs, err := resolveModuleRef(modulePath, rel)
		if err != nil {
			return nil, fmt.Errorf("%s: include %q: %w", modulePath, rel, err)
		}
		inner, err := expandIncludes(abs, bundle, stack)
		if err != nil {
			return nil, err
		}
		out = append(out, inner...)
	}
	return out, nil
}

func lookupBundleModule(bundle *ProgramBundle, modulePath string) (*Program, string, error) {
	if bundle == nil {
		return nil, "", fmt.Errorf("missing module %q (not loaded)", modulePath)
	}
	if prog := bundle.Modules[modulePath]; prog != nil {
		return prog, modulePath, nil
	}
	for path, prog := range bundle.Modules {
		if strings.EqualFold(path, modulePath) {
			return prog, path, nil
		}
	}
	return nil, "", fmt.Errorf("missing module %q (not loaded)", modulePath)
}
