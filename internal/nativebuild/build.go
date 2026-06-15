package nativebuild

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"koda/internal/codegen"
	"koda/internal/kodahome"
	"koda/internal/parser"
	"koda/internal/sema"
)

// Build compiles a loaded program bundle to a native executable using the same
// pipeline as `koda build`: LLVM IR -> object file, linked with the Koda runtime library.
// sourceDisplay is a short label for logs (e.g. base name of the entry file); may be empty.
func Build(bundle *parser.ProgramBundle, output, sourceDisplay string, log func(string)) error {
	return BuildWithOptions(bundle, output, sourceDisplay, log, BuildOptions{})
}

// BuildWithOptions is like [Build] but accepts optimisation flags (e.g. from `koda build --no-opt`).
func BuildWithOptions(bundle *parser.ProgramBundle, output, sourceDisplay string, log func(string), opts BuildOptions) error {
	if log == nil {
		log = func(string) {}
	}
	if err := kodahome.EnsureEnvironment(log); err != nil {
		return err
	}
	ctx, err := sema.PrepareNativeBundleWithOptions(bundle, &sema.PrepareOptions{
		EmitDebug: opts.Debug,
	})
	if err != nil {
		return fmt.Errorf("prepare native bundle: %w", err)
	}
	mod, err := codegen.EmitLLVMIR(ctx)
	if err != nil {
		return fmt.Errorf("llvm ir emit: %w", err)
	}

	irText := mod.String()
	if !opts.NoOpt {
		optimised, err := codegen.OptimiseIR(irText)
		if err != nil {
			if log != nil {
				log(fmt.Sprintf("  warning: LLVM IR optimisation skipped: %v\n", err))
			}
		} else {
			irText = optimised
		}
	}

	tmpDir := ".KODA_build"
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return err
	}
	irPath := filepath.Join(tmpDir, "main.ll")
	if err := os.WriteFile(irPath, []byte(irText), 0644); err != nil {
		return err
	}

	outAbs, err := filepath.Abs(output)
	if err != nil {
		return fmt.Errorf("output path: %w", err)
	}

	rootDir, err := filepath.Abs(".")
	if err != nil {
		return fmt.Errorf("project root: %w", err)
	}

	tc, err := kodahome.FindToolchain()
	if err != nil {
		return fmt.Errorf("toolchain: %w", err)
	}

	objPath := filepath.Join(tmpDir, objFileName())
	label := sourceDisplay
	if label == "" {
		label = "program"
	}
	log(fmt.Sprintf("  compiling  %s\n", label))
	log(fmt.Sprintf("  output     %s\n", outAbs))
	driverLabel := tc.LLC
	if tc.LinkMode == kodahome.LinkClang {
		if strings.TrimSpace(tc.Clang) != "" {
			driverLabel = tc.Clang
		} else {
			driverLabel = ClangDriver()
		}
	} else {
		driverLabel = tc.LLD
	}
	log(fmt.Sprintf("  driver     %s\n\n", driverLabel))

	if tc.LinkMode == kodahome.LinkLLDGNU || tc.LinkMode == kodahome.LinkLLDDarwin {
		if err := runLLC(tc.LLC, irPath, objPath, opts.llcOptFlag(), opts.Debug); err != nil {
			return fmt.Errorf("llc failed: %w", err)
		}
	}

	switch tc.LinkMode {
	case kodahome.LinkLLDGNU:
		if err := linkWithLLDGNU(tc, objPath, outAbs); err != nil {
			return fmt.Errorf("lld gnu link: %w", err)
		}
	case kodahome.LinkLLDDarwin:
		if err := linkWithLLDDarwin(tc, objPath, outAbs); err != nil {
			return fmt.Errorf("lld darwin link: %w", err)
		}
	default:
		if err := runCompileAndLink(tc, irPath, outAbs, rootDir, opts, log); err != nil {
			return fmt.Errorf("clang failed: %w", err)
		}
	}

	log(fmt.Sprintf("\n  built: %s\n", outAbs))
	return nil
}

// BuildWithOverlays loads a program (with optional per-file source overlays) and builds it.
func BuildWithOverlays(entryPath string, overlays map[string]string, output string, log func(string)) error {
	return BuildWithOverlaysOpts(entryPath, overlays, output, log, BuildOptions{})
}

// BuildWithOverlaysOpts is like [BuildWithOverlays] with explicit [BuildOptions].
func BuildWithOverlaysOpts(entryPath string, overlays map[string]string, output string, log func(string), opts BuildOptions) error {
	bundle, err := parser.LoadProgramWithOverlays(entryPath, overlays)
	if err != nil {
		return err
	}
	return BuildWithOptions(bundle, output, filepath.Base(entryPath), log, opts)
}
