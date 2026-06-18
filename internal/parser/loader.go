package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"koda/internal/diagnostic"
	"koda/internal/kodahome"
	"koda/internal/lexer"
)

// Parsed AST cache: absolute path → last seen mtime (nanos) + parsed program.
// Entries are invalidated when Stat mtime differs. Overlay sources are never cached.
type parseCacheEntry struct {
	modTimeNanos int64
	sizeBytes    int64
	prog         *Program
}

var parseCacheMu sync.Mutex
var parseCache = make(map[string]parseCacheEntry)
var loadModuleMu sync.Mutex

func resetParseCache() {
	parseCacheMu.Lock()
	defer parseCacheMu.Unlock()
	parseCache = make(map[string]parseCacheEntry)
}

// LoadProgram parses the entry file and all its transitive imports.
func LoadProgram(entryPath string) (*ProgramBundle, error) {
	return LoadProgramWithOverlays(entryPath, nil)
}

// LoadProgramWithOverlays is like LoadProgram but allows providing source text for specific files.
func LoadProgramWithOverlays(entryPath string, overlays map[string]string) (*ProgramBundle, error) {
	absEntry, err := filepath.Abs(entryPath)
	if err != nil {
		return nil, err
	}

	bundle := &ProgramBundle{
		Modules: make(map[string]*Program),
	}

	visited := make(map[string]bool)
	if err := loadModule(absEntry, bundle, visited, overlays); err != nil {
		return nil, err
	}

	bundle.Entry = bundle.Modules[absEntry]
	return bundle, nil
}

func parseProgramSource(absPath string, src string) (*Program, error) {
	l := lexer.NewLexer(src, absPath)
	tokens, err := l.Tokenize()
	if err != nil {
		return nil, diagnostic.WrapLexer(absPath, src, err)
	}
	pr := NewParser(tokens)
	prog, err := pr.Parse()
	if err != nil {
		return nil, diagnostic.WrapParse(absPath, src, err)
	}
	return prog, nil
}

func loadModuleImports(modulePath string, prog *Program, bundle *ProgramBundle, visited map[string]bool, overlays map[string]string) error {
	var imports []string
	findImports(prog, &imports)
	if len(imports) == 0 {
		return nil
	}
	type result struct {
		err error
	}
	results := make(chan result, len(imports))
	var wg sync.WaitGroup
	for _, rel := range imports {
		rel := rel
		wg.Add(1)
		go func() {
			defer wg.Done()
			abs, err := resolveModuleRef(modulePath, rel)
			if err != nil {
				results <- result{err: err}
				return
			}
			if err := loadModule(abs, bundle, visited, overlays); err != nil {
				results <- result{err: err}
			}
		}()
	}
	wg.Wait()
	close(results)
	for r := range results {
		if r.err != nil {
			return r.err
		}
	}
	return nil
}

func loadModule(path string, bundle *ProgramBundle, visited map[string]bool, overlays map[string]string) error {
	loadModuleMu.Lock()
	if visited[path] {
		if _, ok := bundle.Modules[path]; ok {
			loadModuleMu.Unlock()
			return nil
		}
		loadModuleMu.Unlock()
		return fmt.Errorf("import cycle detected: %s", path)
	}
	visited[path] = true
	loadModuleMu.Unlock()
	defer func() {
		loadModuleMu.Lock()
		delete(visited, path)
		loadModuleMu.Unlock()
	}()

	var src string
	hasOverlayEntry := false
	if overlays != nil {
		if s, ok := overlays[path]; ok {
			hasOverlayEntry = true
			src = s
		}
	}

	if src == "" {
		if !hasOverlayEntry {
			parseCacheMu.Lock()
			e, cached := parseCache[path]
			parseCacheMu.Unlock()

			fi, statErr := os.Stat(path)
			if statErr == nil && cached && e.modTimeNanos == fi.ModTime().UnixNano() && e.sizeBytes == fi.Size() {
				bundle.Modules[path] = e.prog
				return loadModuleImports(path, e.prog, bundle, visited, overlays)
			}
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("could not read %s: %w", path, err)
		}
		src = string(b)
	}

	prog, err := parseProgramSource(path, src)
	if err != nil {
		return err
	}
	loadModuleMu.Lock()
	bundle.Modules[path] = prog
	loadModuleMu.Unlock()

	if !hasOverlayEntry {
		if fi, err := os.Stat(path); err == nil {
			parseCacheMu.Lock()
			parseCache[path] = parseCacheEntry{
				modTimeNanos: fi.ModTime().UnixNano(),
				sizeBytes:    fi.Size(),
				prog:         prog,
			}
			parseCacheMu.Unlock()
		}
	}

	return loadModuleImports(path, prog, bundle, visited, overlays)
}

func findImports(node Node, imports *[]string) {
	switch n := node.(type) {
	case *Program:
		for _, d := range n.Declarations {
			findImports(d, imports)
		}
	case *IncludeDecl:
		rel := strings.Trim(n.Path.Lexeme, `"'`)
		*imports = append(*imports, rel)
	case *UseDecl:
		*imports = append(*imports, useImportRef(n.ModulePath))
	case *FuncDecl:
		for _, p := range n.Params {
			if p.Default != nil {
				findImports(p.Default, imports)
			}
		}
		if n.Body != nil {
			findImports(n.Body, imports)
		}
	case *StructDecl:
		// no nested expressions
	case *EnumDecl:
		// no nested expressions
	case *LetDecl:
		if n.Init != nil {
			findImports(n.Init, imports)
		}
	case *ImportExpr:
		*imports = append(*imports, strings.Trim(n.Path.Lexeme, "\""))
	case *ExpressionStmt:
		findImports(n.Expr, imports)
	case *ReturnStmt:
		if n.Value != nil {
			findImports(n.Value, imports)
		}
	case *DeferStmt:
		if n.Expr != nil {
			findImports(n.Expr, imports)
		}
	case *BlockStmt:
		for _, s := range n.Declarations {
			findImports(s, imports)
		}
	case *IfStmt:
		if n.Condition != nil {
			findImports(n.Condition, imports)
		}
		if n.Then != nil {
			findImports(n.Then, imports)
		}
		if n.Else != nil {
			findImports(n.Else, imports)
		}
	case *WhileStmt:
		if n.Condition != nil {
			findImports(n.Condition, imports)
		}
		if n.Body != nil {
			findImports(n.Body, imports)
		}
	case *LoopStmt:
		if n.Body != nil {
			findImports(n.Body, imports)
		}
	case *PrefixExpr:
		if n.Right != nil {
			findImports(n.Right, imports)
		}
	case *InfixExpr:
		if n.Left != nil {
			findImports(n.Left, imports)
		}
		if n.Right != nil {
			findImports(n.Right, imports)
		}
	case *CallExpr:
		if n.Function != nil {
			findImports(n.Function, imports)
		}
		for _, a := range n.Arguments {
			findImports(a, imports)
		}
	case *AssignExpr:
		if n.Left != nil {
			findImports(n.Left, imports)
		}
		if n.Value != nil {
			findImports(n.Value, imports)
		}
	case *LogicalExpr:
		if n.Left != nil {
			findImports(n.Left, imports)
		}
		if n.Right != nil {
			findImports(n.Right, imports)
		}
	case *ThisExpr:
		// leaf
	case *GroupingExpr:
		if n.Expr != nil {
			findImports(n.Expr, imports)
		}
	case *UpdateExpr:
		if n.Operand != nil {
			findImports(n.Operand, imports)
		}
	case *RangeExpr:
		if n.From != nil {
			findImports(n.From, imports)
		}
		if n.To != nil {
			findImports(n.To, imports)
		}
	case *TemplateExpr:
		for _, p := range n.Parts {
			findImports(p, imports)
		}
	case *SpreadExpr:
		if n.Expr != nil {
			findImports(n.Expr, imports)
		}
	case *TupleExpr:
		for _, e := range n.Elements {
			findImports(e, imports)
		}
	case *IfExpr:
		if n.Condition != nil {
			findImports(n.Condition, imports)
		}
		if n.Then != nil {
			findImports(n.Then, imports)
		}
		if n.Else != nil {
			findImports(n.Else, imports)
		}
	case *SwitchExpr:
		if n.Subject != nil {
			findImports(n.Subject, imports)
		}
		for _, c := range n.Cases {
			if c.Value != nil {
				findImports(c.Value, imports)
			}
			if c.Body != nil {
				findImports(c.Body, imports)
			}
		}
		if n.Default != nil {
			findImports(n.Default, imports)
		}
	case *SliceExpr:
		if n.Object != nil {
			findImports(n.Object, imports)
		}
		if n.Start != nil {
			findImports(n.Start, imports)
		}
		if n.End != nil {
			findImports(n.End, imports)
		}
	case *TernaryExpr:
		if n.Condition != nil {
			findImports(n.Condition, imports)
		}
		if n.Then != nil {
			findImports(n.Then, imports)
		}
		if n.Else != nil {
			findImports(n.Else, imports)
		}
	case *ArrayExpr:
		for _, e := range n.Elements {
			findImports(e, imports)
		}
	case *ObjectExpr:
		for _, e := range n.Values {
			findImports(e, imports)
		}
		for _, e := range n.ComputedKeys {
			findImports(e, imports)
		}
	case *FuncExpr:
		for _, p := range n.Params {
			if p.Default != nil {
				findImports(p.Default, imports)
			}
		}
		if n.Body != nil {
			findImports(n.Body, imports)
		}
	case *IndexExpr:
		if n.Object != nil {
			findImports(n.Object, imports)
		}
		if n.Index != nil {
			findImports(n.Index, imports)
		}
	case *ForStmt:
		for _, init := range n.Inits {
			findImports(init, imports)
		}
		if n.Condition != nil {
			findImports(n.Condition, imports)
		}
		for _, inc := range n.Increments {
			findImports(inc, imports)
		}
		if n.Body != nil {
			findImports(n.Body, imports)
		}
	case *ForInStmt:
		if n.Iterable != nil {
			findImports(n.Iterable, imports)
		}
		if n.Body != nil {
			findImports(n.Body, imports)
		}
	case *BreakStmt:
		// leaf
	case *ContinueStmt:
		// leaf
	case *DeleteStmt:
		if n.Target != nil {
			findImports(n.Target, imports)
		}
	case *SwitchStmt:
		if n.Subject != nil {
			findImports(n.Subject, imports)
		}
		for _, c := range n.Cases {
			if c.Value != nil {
				findImports(c.Value, imports)
			}
			for _, d := range c.Body {
				findImports(d, imports)
			}
		}
		for _, d := range n.Default {
			findImports(d, imports)
		}
	case *DoWhileStmt:
		if n.Body != nil {
			findImports(n.Body, imports)
		}
		if n.Condition != nil {
			findImports(n.Condition, imports)
		}
	case *ForOfStmt:
		if n.Iterable != nil {
			findImports(n.Iterable, imports)
		}
		if n.Body != nil {
			findImports(n.Body, imports)
		}
	case *IdentifierExpr, *LiteralExpr:
		// leaf nodes
	}
}

// pathListEnv reads a PATH-style environment variable (colon- or semicolon-separated).
func pathListEnv(key string) []string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return nil
	}
	return strings.Split(v, string(os.PathListSeparator))
}

// tryResolveModuleUnderRoots looks for moduleID (e.g. "array" or "raylib/core") under each root.
func tryResolveModuleUnderRoots(moduleID string, roots []string) (string, bool) {
	safe := filepath.FromSlash(moduleID)
	base := filepath.Base(safe)
	for _, root := range roots {
		root = strings.TrimSpace(root)
		if root == "" {
			continue
		}
		candidates := []string{
			filepath.Join(root, safe+".koda"),
			filepath.Join(root, safe, "index.koda"),
			filepath.Join(root, safe, base+".koda"),
		}
		for _, p := range candidates {
			fi, err := os.Stat(p)
			if err != nil || fi.IsDir() {
				continue
			}
			abs, err := filepath.Abs(p)
			if err != nil {
				continue
			}
			return abs, true
		}
	}
	return "", false
}

// ResolveImportPath resolves a relative or @module path relative to the importer.
func ResolveImportPath(importerPath, relPath string) (string, error) {
	if strings.HasPrefix(relPath, "@") {
		name := strings.ToLower(relPath[1:])

		// 1) Explicit search paths (override)
		if p, ok := tryResolveModuleUnderRoots(name, pathListEnv("KODA_WRAPPERS")); ok {
			return p, nil
		}
		if p, ok := tryResolveModuleUnderRoots(name, pathListEnv("KODA_PATH")); ok {
			return p, nil
		}

		// 2) Shipped layout next to koda / stdlib/, wrappers/
		if inst, err := kodahome.InstallDir(); err == nil {
			bundled := []string{
				filepath.Join(inst, "stdlib"),
				filepath.Join(inst, "wrappers"),
				filepath.Join(inst, "lib"),
			}
			if p, ok := tryResolveModuleUnderRoots(name, bundled); ok {
				return p, nil
			}
			// Developer repo layout: <repo>/bin/koda with stdlib at <repo>/stdlib
			parent := filepath.Dir(inst)
			devRoots := []string{
				filepath.Join(parent, "stdlib"),
				filepath.Join(parent, "wrappers"),
				parent,
			}
			if p, ok := tryResolveModuleUnderRoots(name, devRoots); ok {
				return p, nil
			}
		}

		// 3) Next to the importing file
		dir := filepath.Dir(importerPath)
		if p, ok := tryResolveModuleUnderRoots(name, []string{dir}); ok {
			return p, nil
		}

		return "", fmt.Errorf("could not resolve @ module %s", relPath)
	}

	return resolvePlainIncludePath(importerPath, relPath)
}

func resolvePlainIncludePath(importerPath, relPath string) (string, error) {
	if filepath.IsAbs(relPath) {
		return filepath.Clean(relPath), nil
	}

	tryFile := func(p string) (string, bool) {
		fi, err := os.Stat(p)
		if err != nil || fi.IsDir() {
			return "", false
		}
		abs, err := filepath.Abs(p)
		if err != nil {
			return "", false
		}
		return abs, true
	}

	if p, ok := tryFile(filepath.Join(filepath.Dir(importerPath), relPath)); ok {
		return p, nil
	}

	for _, r := range pathListEnv("KODA_PATH") {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}
		if p, ok := tryFile(filepath.Join(r, relPath)); ok {
			return p, nil
		}
	}

	if inst, err := kodahome.InstallDir(); err == nil {
		for _, base := range []string{
			filepath.Join(inst, "wrappers"),
			filepath.Join(inst, "stdlib"),
			inst,
		} {
			if p, ok := tryFile(filepath.Join(base, relPath)); ok {
				return p, nil
			}
		}
	}

	return filepath.Abs(filepath.Join(filepath.Dir(importerPath), relPath))
}

func BundleEntryPath(bundle *ProgramBundle) (string, error) {
	for path, prog := range bundle.Modules {
		if prog == bundle.Entry {
			return path, nil
		}
	}
	return "", fmt.Errorf("entry path not found in bundle")
}
