package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"koda/internal/wrappermeta"
)

func runUpgradeCLI(args []string) error {
	dir, rest, err := parseWrapperDirArg(args)
	if err != nil {
		return err
	}
	cfg, meta, err := configFromMetaDir(dir)
	if err != nil {
		return err
	}
	cfg.GeneratedAt = time.Now().UTC().Format(time.RFC3339)
	if cfg.LibraryVersion == "" {
		cfg.LibraryVersion = meta.LibraryVersion
	}
	applyUpgradeFlags(cfg, rest)
	if cfg.Verbose {
		fmt.Printf("%s: upgrading @%s from %s\n", toolDisplayName(), meta.Name, dir)
		fmt.Printf("  headers: %s\n\n", strings.Join(cfg.InputHeaders, ", "))
	}
	return runGenerate(cfg)
}

func runCheckCLI(args []string) error {
	dir, _, err := parseWrapperDirArg(args)
	if err != nil {
		return err
	}
	meta, err := wrappermeta.LoadDir(dir)
	if err != nil {
		return err
	}
	report, err := wrappermeta.CheckDrift(meta)
	if err != nil {
		return err
	}
	report.Dir = dir
	printDriftReport(report)
	if report.Stale() {
		fmt.Println()
		fmt.Printf("Run: %s upgrade %s\n", toolDisplayName(), dir)
		return fmt.Errorf("wrapper @%s is stale", meta.Name)
	}
	return nil
}

func parseWrapperDirArg(args []string) (string, []string, error) {
	if len(args) == 0 {
		return "", nil, fmt.Errorf("usage: %s upgrade <wrapper-dir>\n       %s check <wrapper-dir>", toolDisplayName(), toolDisplayName())
	}
	dir := strings.TrimSpace(args[0])
	rest := args[1:]
	if strings.HasPrefix(dir, "@") {
		resolved, err := resolveWrapperName(dir[1:])
		if err != nil {
			return "", nil, err
		}
		dir = resolved
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return "", nil, err
	}
	if _, err := os.Stat(filepath.Join(abs, wrappermeta.FileName)); err != nil {
		return "", nil, fmt.Errorf("%s: no %s (not a wrapper package)", abs, wrappermeta.FileName)
	}
	return abs, rest, nil
}

func resolveWrapperName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", fmt.Errorf("empty wrapper name")
	}
	seen := map[string]bool{}
	var roots []string
	add := func(p string) {
		p = strings.TrimSpace(p)
		if p == "" || seen[p] {
			return
		}
		if fi, err := os.Stat(p); err != nil || !fi.IsDir() {
			return
		}
		seen[p] = true
		roots = append(roots, p)
	}
	if v := os.Getenv("KODA_WRAPPERS"); v != "" {
		for _, part := range strings.Split(v, string(os.PathListSeparator)) {
			add(part)
		}
	}
	if exe, err := os.Executable(); err == nil {
		add(filepath.Join(filepath.Dir(exe), "wrappers"))
		add(filepath.Join(filepath.Dir(exe), "..", "wrappers"))
	}
	for _, root := range roots {
		dir := filepath.Join(root, name)
		if fi, err := os.Stat(filepath.Join(dir, wrappermeta.FileName)); err == nil && !fi.IsDir() {
			return dir, nil
		}
	}
	return "", fmt.Errorf("wrapper @%s not found (set KODA_WRAPPERS)", name)
}

func configFromMetaDir(dir string) (*WrapGenConfig, wrappermeta.PackageMeta, error) {
	meta, err := wrappermeta.LoadDir(dir)
	if err != nil {
		return nil, wrappermeta.PackageMeta{}, err
	}
	useClang := meta.UseClangEnabled()
	cfg := &WrapGenConfig{
		LibraryName:   meta.Name,
		InputHeaders:  append([]string{}, meta.Headers...),
		OutputDir:     dir,
		PrimaryHeader: meta.PrimaryHeader,
		IncludePaths:  meta.ResolveIncludePaths(),
		LinkFlags:     meta.ResolveLinkFlags(),
		UseClang:      useClang,
		Language:      "koda",
		Version:       WrapgenVersion,
	}
	if cfg.PrimaryHeader == "" && len(cfg.InputHeaders) > 0 {
		cfg.PrimaryHeader = filepath.Base(cfg.InputHeaders[0])
	}
	for _, h := range cfg.InputHeaders {
		if _, err := os.Stat(h); err != nil {
			return nil, meta, fmt.Errorf("header not found: %s (update headers in %s or reinstall the native library)", h, wrappermeta.FileName)
		}
	}
	return cfg, meta, nil
}

func applyUpgradeFlags(cfg *WrapGenConfig, args []string) {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-v", "--verbose":
			cfg.Verbose = true
		case "--no-clang":
			cfg.UseClang = false
		case "--no-docs":
			cfg.NoDocsMarkdown = true
			cfg.NoHTML = true
		case "--no-html":
			cfg.NoHTML = true
		case "--library-version":
			if i+1 < len(args) {
				i++
				cfg.LibraryVersion = args[i]
			}
		}
	}
}

func printDriftReport(report wrappermeta.DriftReport) {
	meta := report.Meta
	fmt.Printf("Wrapper @%s (%s)\n", meta.Name, report.Dir)
	if meta.LibraryVersion != "" {
		fmt.Printf("  library version: %s\n", meta.LibraryVersion)
	}
	fmt.Printf("  generated: %s\n", meta.GeneratedAt)
	fmt.Printf("  generator: %s %s\n", meta.Generator, meta.Version)
	if len(report.Unchanged) > 0 {
		fmt.Printf("  unchanged headers: %d\n", len(report.Unchanged))
	}
	for _, p := range report.Changed {
		fmt.Printf("  changed: %s\n", p)
	}
	for _, p := range report.Missing {
		fmt.Printf("  missing: %s\n", p)
	}
	if !report.Stale() {
		fmt.Println("  status: up to date")
	}
}

func runGenerate(cfg *WrapGenConfig) error {
	if err := validateConfig(cfg); err != nil {
		return err
	}
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return fmt.Errorf("could not create output directory: %w", err)
	}
	generator := NewWrapperGenerator(cfg)
	api, err := generator.ParseHeaders()
	if err != nil {
		return fmt.Errorf("failed to parse headers: %w", err)
	}
	if err := generator.GenerateWrapper(api); err != nil {
		return fmt.Errorf("failed to generate wrapper: %w", err)
	}
	if err := generator.emitPackageArtifacts(api); err != nil {
		return fmt.Errorf("failed to write package files: %w", err)
	}
	if err := generator.emitProfessionalDocs(api); err != nil {
		return fmt.Errorf("failed to write documentation: %w", err)
	}
	printSuccess(cfg, api)
	return nil
}
