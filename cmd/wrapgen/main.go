package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// WrapGenConfig drives wrapper generation (shared with cmd/wrapgen generator).
type WrapGenConfig struct {
	LibraryName    string
	InputHeaders   []string
	OutputDir      string
	Language       string
	Version        string
	GeneratedAt    string
	PrimaryHeader  string
	IncludePaths   []string
	LinkFlags      []string // -L, -l, --linkflags tokens for koda.json / README
	NoDocsMarkdown bool
	NoHTML         bool
	Verbose        bool
	UseClang       bool // try clang AST first (default true)
	// Legacy flags (still accepted when -name is used):
	Documentation bool
	BuildSystem     bool
	IncludeTests    bool
	ComplexCPP      bool
}

func main() {
	config, err := parseAllCLI(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		printModernUsage()
		os.Exit(2)
	}

	if config.Verbose {
		fmt.Printf("%s: generating bindings for %q\n", toolDisplayName(), config.LibraryName)
		fmt.Printf("  headers : %s\n", strings.Join(config.InputHeaders, ", "))
		fmt.Printf("  output  : %s\n\n", config.OutputDir)
	}

	if err := validateConfig(config); err != nil {
		log.Fatalf("error: %v\n\nRun %s --help for usage.", err, toolDisplayName())
	}

	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		log.Fatalf("error: could not create output directory: %v", err)
	}

	generator := NewWrapperGenerator(config)

	api, err := generator.ParseHeaders()
	if err != nil {
		log.Fatalf("error: failed to parse headers: %v", err)
	}

	if err := generator.GenerateWrapper(api); err != nil {
		log.Fatalf("error: failed to generate wrapper: %v", err)
	}

	if err := generator.emitPackageArtifacts(api); err != nil {
		log.Fatalf("error: failed to write package files: %v", err)
	}

	if err := generator.emitProfessionalDocs(api); err != nil {
		log.Fatalf("error: failed to write documentation: %v", err)
	}

	printSuccess(config, api)
}

func printSuccess(config *WrapGenConfig, api *API) {
	out := config.OutputDir
	lib := config.LibraryName

	fmt.Printf("Generated bindings for %s\n\n", lib)
	fmt.Printf("  Output folder : %s\n", out)
	fmt.Printf("  %s/%s.koda\n", out, lib)
	fmt.Printf("  %s/wrapper.c\n", out)
	if !config.NoDocsMarkdown {
		fmt.Printf("  %s/README.md\n", out)
		fmt.Printf("  %s/api_reference.md\n", out)
		fmt.Printf("  %s/examples.md\n", out)
		fmt.Printf("  %s/koda.json\n", out)
		fmt.Printf("  %s/META.json\n", out)
	}
	if !config.NoHTML {
		fmt.Printf("  %s/docs/index.html\n", out)
	}
	fmt.Printf("\n")
	fmt.Printf("  Functions : %d\n", len(api.Functions))
	fmt.Printf("  Structs   : %d\n", len(api.Structs))
	fmt.Printf("  Enums     : %d\n", len(api.Enums))
	fmt.Printf("  Constants : %d\n", len(api.Constants))

	fmt.Fprintf(os.Stdout, "\nTo use in your Koda program:\n\n  #include \"@%s\"\n\n", lib)
	fmt.Fprintf(os.Stdout, "Or set KODA_WRAPPERS to this folder and: import \"@%s\"\n\n", lib)
	fmt.Fprintf(os.Stdout, "Build (merge koda.json native section or set env):\n\n")
	fmt.Fprintf(os.Stdout, "  set KODA_NATIVE_SOURCES=%s\\\\wrapper.c\n", out)
	fmt.Fprintf(os.Stdout, "  set KODA_LINKFLAGS=-I<include-dir> -L<lib-dir> -l%s\n", lib)
	fmt.Fprintf(os.Stdout, "  koda build mygame.koda -o mygame.exe\n")
	if !config.NoDocsMarkdown {
		fmt.Fprintf(os.Stdout, "\nSee %s/README.md, %s/examples.md, and %s/docs/index.html\n", out, out, out)
	} else if !config.NoHTML {
		fmt.Fprintf(os.Stdout, "\nSee %s/docs/index.html for browsable API.\n", out)
	}
}
