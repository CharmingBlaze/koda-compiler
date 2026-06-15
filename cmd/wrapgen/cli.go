package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func parseAllCLI(args []string) (*WrapGenConfig, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("usage: %s [options] <header.h> [...]\n       %s upgrade <wrapper-dir>\n       %s check <wrapper-dir>\n       %s -name <lib> -headers <a.h>[,b.h] -out <dir>  (legacy)\nRun %s --help", toolDisplayName(), toolDisplayName(), toolDisplayName(), toolDisplayName(), toolDisplayName())
	}
	switch args[0] {
	case "--version":
		fmt.Println(WrapgenVersion)
		os.Exit(0)
	case "-h", "--help":
		printModernUsage()
		os.Exit(0)
	case "upgrade":
		return nil, fmt.Errorf("use: %s upgrade <wrapper-dir> (not combined with header args)", toolDisplayName())
	case "check":
		return nil, fmt.Errorf("use: %s check <wrapper-dir>", toolDisplayName())
	case "install":
		return nil, fmt.Errorf("use: %s install <name[@version]>", toolDisplayName())
	case "list":
		return nil, fmt.Errorf("use: %s list", toolDisplayName())
	}
	for _, a := range args {
		if a == "-headers" {
			return parseLegacyCLI(args)
		}
	}
	return parseModernCLI(args)
}

func parseModernCLI(args []string) (*WrapGenConfig, error) {
	cfg := &WrapGenConfig{
		Version:     WrapgenVersion,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		OutputDir:   ".",
		Language:    "koda",
		UseClang:    true,
	}
	var headers []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case (a == "-o" || a == "-out") && i+1 < len(args):
			cfg.OutputDir = args[i+1]
			i++
		case (a == "-name" || a == "--name") && i+1 < len(args):
			cfg.LibraryName = args[i+1]
			i++
		case a == "-I" && i+1 < len(args):
			cfg.IncludePaths = append(cfg.IncludePaths, args[i+1])
			i++
		case a == "-L" && i+1 < len(args):
			cfg.LinkFlags = append(cfg.LinkFlags, "-L"+args[i+1])
			i++
		case a == "-l" && i+1 < len(args):
			cfg.LinkFlags = append(cfg.LinkFlags, "-l"+args[i+1])
			i++
		case a == "--linkflags" && i+1 < len(args):
			cfg.LinkFlags = append(cfg.LinkFlags, strings.Fields(args[i+1])...)
			i++
		case a == "--no-docs":
			cfg.NoDocsMarkdown = true
			cfg.NoHTML = true
		case a == "--no-html":
			cfg.NoHTML = true
		case a == "--cpp":
			cfg.UseCPP = true
		case a == "--no-clang":
			cfg.UseClang = false
		case a == "--library-version" && i+1 < len(args):
			cfg.LibraryVersion = args[i+1]
			i++
		case a == "-v" || a == "--verbose":
			cfg.Verbose = true
		case strings.HasPrefix(a, "-"):
			return nil, fmt.Errorf("unknown flag: %s", a)
		default:
			headers = append(headers, a)
		}
	}
	if len(headers) == 0 {
		return nil, fmt.Errorf("expected at least one header file (.h)")
	}
	cfg.InputHeaders = headers
	base := filepath.Base(headers[0])
	cfg.PrimaryHeader = base
	cfg.LibraryName = strings.TrimSuffix(base, filepath.Ext(base))
	if cfg.UseCPP || isCPPHeader(headers[0]) {
		cfg.UseCPP = true
	}
	if cfg.LibraryName == "" {
		cfg.LibraryName = "bindings"
	}
	return cfg, nil
}

func parseLegacyCLI(args []string) (*WrapGenConfig, error) {
	fs := flag.NewFlagSet("kodawrap", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	cfg := &WrapGenConfig{
		Version:        WrapgenVersion,
		GeneratedAt:    time.Now().UTC().Format(time.RFC3339),
		OutputDir:      "./wrapper",
		Language:       "koda",
		ComplexCPP:     true,
		Documentation:  true,
		BuildSystem:    false,
		IncludeTests:   false,
		NoDocsMarkdown: false,
		NoHTML:         false,
	}
	var headers string
	fs.StringVar(&cfg.LibraryName, "name", "", "library name (required)")
	fs.StringVar(&cfg.OutputDir, "out", "./wrapper", "output directory (-out or -o)")
	var outShort string
	fs.StringVar(&outShort, "o", "", "alias for -out (same value)")
	fs.StringVar(&cfg.Language, "lang", "koda", "target language (koda only)")
	fs.BoolVar(&cfg.Documentation, "docs", true, "generate README + api_reference (+ HTML unless --no-html in modern mode)")
	fs.BoolVar(&cfg.BuildSystem, "build", false, "deprecated; ignored")
	fs.BoolVar(&cfg.IncludeTests, "tests", false, "deprecated; ignored")
	fs.BoolVar(&cfg.ComplexCPP, "cpp", true, "reserved")
	fs.BoolVar(&cfg.Verbose, "v", false, "verbose")
	fs.StringVar(&headers, "headers", "", "comma-separated headers (required)")

	fs.Usage = func() { printModernUsage() }

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if strings.TrimSpace(outShort) != "" {
		cfg.OutputDir = outShort
	}
	if headers != "" {
		for _, h := range strings.Split(headers, ",") {
			if h = strings.TrimSpace(h); h != "" {
				cfg.InputHeaders = append(cfg.InputHeaders, h)
			}
		}
	}
	if cfg.LibraryName == "" {
		return nil, fmt.Errorf("-name is required in legacy mode")
	}
	if len(cfg.InputHeaders) == 0 {
		return nil, fmt.Errorf("-headers is required in legacy mode")
	}
	if len(cfg.InputHeaders) > 0 {
		cfg.PrimaryHeader = filepath.Base(cfg.InputHeaders[0])
	}
	cfg.NoDocsMarkdown = !cfg.Documentation
	cfg.NoHTML = !cfg.Documentation
	return cfg, nil
}

func printModernUsage() {
	me := toolDisplayName()
	fmt.Fprintf(os.Stderr, `%s — Koda C/C++ header → .koda, wrapper.c, and docs

USAGE (preferred)
  %s [options] <header.h> [<more.h> ...]

UPGRADE (regenerate from META.json)
  %s upgrade <wrapper-dir|@name>
  %s check <wrapper-dir|@name>

CATALOG (known libraries)
  %s list
  %s install <name[@version]> [-o dir] [--project]

OPTIONS
  -name <lib>   library name (default: header basename)
  -o <dir>      output directory (default: .)
  -I <dir>      extra include path for clang (repeatable)
  -L <libdir>   linker search path (recorded in koda.json)
  -l <lib>      link library name (recorded in koda.json)
  --linkflags   extra linker flags as one string
  --cpp         parse headers as C++ (-std=c++17)
  --no-clang    regex-only parsing (no clang AST)
  --no-docs     skip README.md, api_reference.md, examples.md
  --no-html     skip docs/index.html only
  -v            verbose
  --version     print version and exit
  -h, --help    this help

LEGACY (still supported)
  %s -name <lib> -headers <a.h>[,b.h] -out <dir> [-docs=false] [-v]

OUTPUT (organized package)
  <name>.koda       wrapper.c       README.md
  api_reference.md  examples.md     koda.json       META.json
  docs/index.html

`, me, me, me, me, me, me, me)
}

func validateConfig(config *WrapGenConfig) error {
	if config.LibraryName == "" {
		return fmt.Errorf("library name missing")
	}
	if len(config.InputHeaders) == 0 {
		return fmt.Errorf("no header files")
	}
	for _, header := range config.InputHeaders {
		if _, err := os.Stat(header); err != nil {
			return fmt.Errorf("header file not found: %s: %w", header, err)
		}
	}
	return nil
}
