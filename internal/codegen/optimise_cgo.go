//go:build cgo && llvm14

package codegen

import (
	"fmt"
	"os"

	tgllvm "tinygo.org/x/go-llvm"
)

// OptimiseIR parses LLVM IR text, runs the LLVM O2 default pipeline (includes mem2reg,
// instcombine, GVN, CFG simplification, and related passes), and returns updated IR text.
// On failure it returns the original IR and a non-nil error so callers can log and continue.
func OptimiseIR(irText string) (string, error) {
	if err := tgllvm.InitializeNativeTarget(); err != nil {
		return irText, fmt.Errorf("llvm init target: %w", err)
	}
	if err := tgllvm.InitializeNativeAsmPrinter(); err != nil {
		return irText, fmt.Errorf("llvm init asm printer: %w", err)
	}

	ctx := tgllvm.NewContext()
	defer ctx.Dispose()

	f, err := os.CreateTemp("", "koda-opt-*.ll")
	if err != nil {
		return irText, fmt.Errorf("temp ir file: %w", err)
	}
	tmpName := f.Name()
	_, _ = f.WriteString(irText)
	if err := f.Close(); err != nil {
		_ = os.Remove(tmpName)
		return irText, fmt.Errorf("temp ir close: %w", err)
	}
	defer func() { _ = os.Remove(tmpName) }()

	buf, err := tgllvm.NewMemoryBufferFromFile(tmpName)
	if err != nil {
		return irText, fmt.Errorf("llvm memory buffer: %w", err)
	}
	defer buf.Dispose()

	mod, err := ctx.ParseIR(buf)
	if err != nil {
		return irText, fmt.Errorf("parse ir: %w", err)
	}
	defer mod.Dispose()

	if err := tgllvm.VerifyModule(mod, tgllvm.ReturnStatusAction); err != nil {
		return irText, fmt.Errorf("verify: %w", err)
	}

	triple := tgllvm.DefaultTargetTriple()
	targ, err := tgllvm.GetTargetFromTriple(triple)
	if err != nil {
		return irText, fmt.Errorf("get target: %w", err)
	}
	mt := targ.CreateTargetMachine(triple, "", "", tgllvm.CodeGenLevelDefault, tgllvm.RelocDefault, tgllvm.CodeModelDefault)
	defer mt.Dispose()

	pbo := tgllvm.NewPassBuilderOptions()
	defer pbo.Dispose()

	// default<O2> runs the standard O2 pipeline (mem2reg, instcombine, GVN, etc.).
	if err := mod.RunPasses("default<O2>", mt, pbo); err != nil {
		return irText, fmt.Errorf("run passes: %w", err)
	}

	return mod.String(), nil
}
