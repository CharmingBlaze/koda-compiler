package codegen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"koda/internal/kodahome"

	"github.com/llir/llvm/ir"
)

// Linker handles linking LLVM IR text to executables via llc + a system linker.
// Prefer [koda/internal/nativebuild.Build], which lowers IR with llc then links with
// runtime/libkoda_runtime.a via Clang.
// and matches `koda build`; Linker remains for toolchains that compile IR to .o first.
type Linker struct {
	runtimeLib string
}

// NewLinker creates a new linker.
func NewLinker(runtimeLib string) *Linker {
	return &Linker{
		runtimeLib: runtimeLib,
	}
}

// LinkModule writes mod to irFile, runs llc to produce an object file, then links.
func (l *Linker) LinkModule(mod *ir.Module, irFile, output string) error {
	if err := os.WriteFile(irFile, []byte(mod.String()), 0644); err != nil {
		return fmt.Errorf("write ir: %w", err)
	}
	return l.LinkExecutable(irFile, output)
}

// LinkExecutable compiles LLVM IR file to an object with llc, then links to an executable.
func (l *Linker) LinkExecutable(irFile, output string) error {
	objFile := output + ".o"

	llcPath := kodahome.LLC()
	llcCmd := exec.Command(llcPath, "-filetype=obj", irFile, "-o", objFile)
	llcCmd.Stderr = os.Stderr
	llcCmd.Stdout = os.Stdout
	if err := llcCmd.Run(); err != nil {
		return fmt.Errorf("%s failed: %w", llcPath, err)
	}

	var linkCmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		linkCmd = exec.Command("gcc",
			"-static",
			"-O3",
			"-s",
			objFile,
			l.runtimeLib,
			"-lm",
			"-o", output,
		)

	case "windows":
		// Default matches MinGW-w64 cross prefix; set KODA_MINGW_GCC=gcc on MSYS2/LLVM-only installs.
		cc := os.Getenv("KODA_MINGW_GCC")
		if cc == "" {
			cc = "x86_64-w64-mingw32-gcc"
		}
		linkCmd = exec.Command(cc,
			"-static",
			"-static-libgcc",
			"-O3",
			objFile,
			l.runtimeLib,
			"-o", output,
		)

	case "darwin":
		linkCmd = exec.Command("clang",
			"-O3",
			"-Wl,-dead_strip",
			objFile,
			l.runtimeLib,
			"-o", output,
		)

	default:
		_ = os.Remove(objFile)
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	linkCmd.Stderr = os.Stderr
	linkCmd.Stdout = os.Stdout

	if err := linkCmd.Run(); err != nil {
		_ = os.Remove(objFile)
		return fmt.Errorf("linker failed: %w", err)
	}

	if err := os.Remove(objFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove temp object %s: %w", objFile, err)
	}
	return nil
}

// LinkExecutableInDir is like LinkExecutable but places the object file in dir (avoids cwd clutter).
func (l *Linker) LinkExecutableInDir(irFile, output, workDir string) error {
	base := filepath.Base(output)
	objFile := filepath.Join(workDir, base+".o")

	llcPath := kodahome.LLC()
	llcCmd := exec.Command(llcPath, "-filetype=obj", irFile, "-o", objFile)
	llcCmd.Dir = workDir
	llcCmd.Stderr = os.Stderr
	llcCmd.Stdout = os.Stdout
	if err := llcCmd.Run(); err != nil {
		return fmt.Errorf("%s failed: %w", llcPath, err)
	}

	var linkCmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		linkCmd = exec.Command("gcc", "-static", "-O3", "-s", objFile, l.runtimeLib, "-lm", "-o", output)
	case "windows":
		cc := os.Getenv("KODA_MINGW_GCC")
		if cc == "" {
			cc = "x86_64-w64-mingw32-gcc"
		}
		linkCmd = exec.Command(cc, "-static", "-static-libgcc", "-O3", objFile, l.runtimeLib, "-o", output)
	case "darwin":
		linkCmd = exec.Command("clang", "-O3", "-Wl,-dead_strip", objFile, l.runtimeLib, "-o", output)
	default:
		_ = os.Remove(objFile)
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	linkCmd.Dir = workDir
	linkCmd.Stderr = os.Stderr
	linkCmd.Stdout = os.Stdout

	if err := linkCmd.Run(); err != nil {
		_ = os.Remove(objFile)
		return fmt.Errorf("linker failed: %w", err)
	}
	_ = os.Remove(objFile)
	return nil
}
