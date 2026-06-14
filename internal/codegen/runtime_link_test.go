package codegen

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/llir/llvm/ir"
)

// llvmRuntimeSymbols returns distinct LLVM symbol names the native backend links against.
func llvmRuntimeSymbols(t *testing.T) []string {
	t.Helper()
	mod := ir.NewModule()
	declareRuntimeFunctions(mod)
	seen := make(map[string]struct{})
	var names []string
	for _, fn := range mod.Funcs {
		name := fn.GlobalIdent.Name()
		if name == "llvm.dbg.value" {
			continue
		}
		if !strings.HasPrefix(name, "koda_") {
			continue
		}
		if _, ok := seen[name]; ok {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}
	return names
}

func TestRuntimeLLVMSymbolsDefinedInArchive(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	root := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", ".."))
	lib := filepath.Join(root, "runtime", "libkoda_runtime.a")
	if _, err := os.Stat(lib); err != nil {
		t.Skip("runtime/libkoda_runtime.a not present; run `make -C runtime` or scripts/build-runtime.ps1")
	}
	nm, err := exec.LookPath("llvm-nm")
	if err != nil {
		nm, err = exec.LookPath("nm")
		if err != nil {
			t.Skip("llvm-nm/nm not on PATH")
		}
	}
	out, err := exec.Command(nm, lib).CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s: %v\n%s", nm, lib, err, out)
	}
	defined := make(map[string]bool)
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasSuffix(line, ":") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		sym := fields[len(fields)-1]
		if strings.HasPrefix(sym, "koda_") {
			defined[sym] = true
		}
	}
	for _, sym := range llvmRuntimeSymbols(t) {
		if !defined[sym] {
			t.Errorf("LLVM declares %q but it is not defined in %s (rebuild runtime after changing runtime.go)", sym, lib)
		}
	}
}

func TestRuntimeLLVMSymbolsMentionedInSources(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	srcDir := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "..", "runtime", "src"))
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		t.Fatalf("read runtime src: %v", err)
	}
	var b strings.Builder
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".c") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(srcDir, e.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", e.Name(), err)
		}
		b.Write(data)
		b.WriteByte('\n')
	}
	combined := b.String()
	for _, sym := range llvmRuntimeSymbols(t) {
		re := regexp.MustCompile(`\b` + regexp.QuoteMeta(sym) + `\s*\(`)
		if !re.MatchString(combined) {
			t.Errorf("LLVM runtime symbol %q has no definition-like occurrence in runtime/src/*.c (update koda_runtime.c or runtime.go)", sym)
		}
	}
}
