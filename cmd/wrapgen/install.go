package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"koda/internal/project"
	"koda/internal/wrapcatalog"
)

type installOptions struct {
	outDir    string
	project   bool
	verbose   bool
	regenerate bool
}

func runListCLI(args []string) error {
	_ = args
	cat, err := wrapcatalog.Load()
	if err != nil {
		return err
	}
	fmt.Println("Known libraries (koda wrap install <name>[@version]):")
	for _, lib := range cat.Libraries {
		vers := make([]string, 0, len(lib.Versions))
		for _, v := range lib.Versions {
			vers = append(vers, v.Version)
		}
		line := fmt.Sprintf("  %-10s %s", lib.Name, lib.Description)
		if lib.DefaultVersion != "" {
			line += fmt.Sprintf("  (default @%s)", lib.DefaultVersion)
		}
		if len(vers) > 0 {
			line += fmt.Sprintf("  [%s]", strings.Join(vers, ", "))
		}
		fmt.Println(line)
	}
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  koda wrap install raylib")
	fmt.Println("  koda wrap install sqlite3@3 --project")
	return nil
}

func runInstallCLI(args []string) error {
	opts := installOptions{outDir: "", project: false}
	var specArg string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-o", "-out":
			if i+1 >= len(args) {
				return fmt.Errorf("-o requires a directory")
			}
			i++
			opts.outDir = args[i]
		case "--project", "-p":
			opts.project = true
		case "-v", "--verbose":
			opts.verbose = true
		case "--regenerate":
			opts.regenerate = true
		default:
			if strings.HasPrefix(args[i], "-") {
				return fmt.Errorf("unknown flag: %s", args[i])
			}
			if specArg != "" {
				return fmt.Errorf("unexpected argument: %s", args[i])
			}
			specArg = args[i]
		}
	}
	if specArg == "" {
		return fmt.Errorf("usage: %s install <name[@version]> [-o dir] [--project]\n       %s list", toolDisplayName(), toolDisplayName())
	}
	spec, err := wrapcatalog.ParseInstallSpec(specArg)
	if err != nil {
		return err
	}
	cat, err := wrapcatalog.Load()
	if err != nil {
		return err
	}
	lib, ver, err := cat.Find(spec.Name, spec.Version)
	if err != nil {
		return err
	}

	cwd, _ := os.Getwd()
	ctx, err := wrapcatalog.DefaultExpandContext(cwd)
	if err != nil {
		return err
	}

	outDir := opts.outDir
	if outDir == "" {
		outDir = filepath.Join("wrappers", lib.Name)
	}
	absOut, err := filepath.Abs(outDir)
	if err != nil {
		return err
	}

	if ver.Prebuilt != "" && !opts.regenerate {
		prebuilt := ctx.Expand(ver.Prebuilt)
		if st, err := os.Stat(prebuilt); err == nil && st.IsDir() {
			if opts.verbose {
				fmt.Printf("copying prebuilt wrapper from %s\n", prebuilt)
			}
			if err := os.RemoveAll(absOut); err != nil {
				return err
			}
			if err := wrapcatalog.CopyDir(prebuilt, absOut); err != nil {
				return fmt.Errorf("copy prebuilt wrapper: %w", err)
			}
			fmt.Printf("Installed @%s from SDK prebuilt wrapper\n", lib.Name)
			fmt.Printf("  %s\n", absOut)
			if opts.project {
				if err := mergeProjectNative(absOut, lib.Name, ver); err != nil {
					return err
				}
			}
			printInstallNextSteps(lib.Name, absOut)
			return nil
		}
		if opts.verbose {
			fmt.Printf("prebuilt not found at %s — generating from headers\n", prebuilt)
		}
	}

	headers := ctx.ExpandAll(ver.Headers)
	header, err := wrapcatalog.ResolveHeader(headers)
	if err != nil {
		return fmt.Errorf("@%s: %w\nInstall the native library or pass headers manually with koda wrap", lib.Name, err)
	}
	includePaths := ver.ResolvedIncludePaths(ctx, header)
	cfg := &WrapGenConfig{
		LibraryName:    lib.Name,
		InputHeaders:   []string{header},
		OutputDir:      absOut,
		PrimaryHeader:  filepath.Base(header),
		IncludePaths:   includePaths,
		LinkFlags:      append([]string{}, ver.LinkFlags...),
		UseClang:       true,
		UseCPP:         isCPPHeader(header),
		Language:       "koda",
		Version:        WrapgenVersion,
		GeneratedAt:    time.Now().UTC().Format(time.RFC3339),
		LibraryVersion: ver.LibraryVersion,
		Verbose:        opts.verbose,
	}
	if err := runGenerate(cfg); err != nil {
		return err
	}
	if opts.project {
		if err := mergeProjectNative(absOut, lib.Name, ver); err != nil {
			return err
		}
	}
	return nil
}

func mergeProjectNative(outDir, name string, ver wrapcatalog.VersionRecipe) error {
	cfg, root, err := project.Find(cwd())
	if err != nil {
		return err
	}
	if cfg == nil {
		return fmt.Errorf("--project requires koda.json in the current directory (run koda new first)")
	}
	rel, err := filepath.Rel(root, filepath.Join(outDir, "wrapper.c"))
	if err != nil {
		rel = filepath.Join("wrappers", name, "wrapper.c")
	} else {
		rel = filepath.ToSlash(rel)
	}
	cfg.Native.Sources = mergeUniqueSources(cfg.Native.Sources, rel)
	if ver.Graphics {
		cfg.Native.Graphics = true
	}
	if strings.TrimSpace(cfg.Native.Linkflags) == "" && len(ver.LinkFlags) > 0 {
		cfg.Native.Linkflags = strings.Join(ver.LinkFlags, " ")
	}
	if err := project.Save(root, cfg); err != nil {
		return err
	}
	fmt.Printf("Updated %s (native.sources + graphics/link flags)\n", filepath.Join(root, project.FileName))
	return nil
}

func mergeUniqueSources(base []string, add ...string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range append(base, add...) {
		s = strings.TrimSpace(s)
		if s == "" || seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}

func printInstallNextSteps(name, outDir string) {
	fmt.Println()
	fmt.Printf("  #include \"@%s\"\n", name)
	fmt.Println()
	fmt.Println("Next: koda wrap check", outDir)
}

func isCPPHeader(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".hpp" || ext == ".hxx" || ext == ".hh"
}

func cwd() string {
	d, err := os.Getwd()
	if err != nil {
		return "."
	}
	return d
}
