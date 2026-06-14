//go:build KODA_wip
// +build KODA_wip

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// kodaBin is the path to a freshly built koda executable (see TestMain).
var kodaBin string

// moduleRoot is the repository root (directory containing go.mod).
var moduleRoot string

func TestMain(m *testing.M) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	moduleRoot = filepath.Clean(filepath.Join(wd, "..", ".."))

	tmp, err := os.MkdirTemp("", "koda-cli")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer func() { _ = os.RemoveAll(tmp) }()

	name := "koda"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	kodaBin = filepath.Join(tmp, name)

	build := exec.Command("go", "build", "-o", kodaBin, ".")
	build.Dir = wd
	build.Env = os.Environ()
	if out, err := build.CombinedOutput(); err != nil {
		fmt.Fprintf(os.Stderr, "go build: %v\n%s", err, out)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func scriptPath(parts ...string) string {
	return filepath.Join(append([]string{moduleRoot}, parts...)...)
}

func TestCLI_noArgsShowsHelp(t *testing.T) {
	cmd := exec.Command(kodaBin)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("no-args: %v\n%s", err, out)
	}
	if cmd.ProcessState.ExitCode() != 0 {
		t.Fatalf("exit code = %d want 0 (print help)", cmd.ProcessState.ExitCode())
	}
	if !bytes.Contains(out, []byte("USAGE")) || !bytes.Contains(out, []byte("COMMANDS")) {
		t.Fatalf("expected help text, got:\n%s", out)
	}
}

func TestCLI_version(t *testing.T) {
	cmd := exec.Command(kodaBin, "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version: %v\n%s", err, out)
	}
	if !bytes.Contains(out, []byte("koda")) {
		t.Fatalf("expected version line, got %q", string(out))
	}
}

func TestCLI_wrap_help(t *testing.T) {
	cmd := exec.Command(kodaBin, "wrap", "--help")
	cmd.Dir = moduleRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("wrap: %v\n%s", err, out)
	}
	if !bytes.Contains(out, []byte("kodawrap")) {
		t.Fatalf("expected wrap help mentioning kodawrap, got:\n%s", out)
	}
}

func TestCLI_check_hello(t *testing.T) {
	cmd := exec.Command(kodaBin, "check", scriptPath("tests", "hello.koda"))
	cmd.Dir = moduleRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("check: %v\n%s", err, out)
	}
	if string(out) != "OK\n" {
		t.Fatalf("stdout = %q want OK + newline", string(out))
	}
}

func TestCLI_check_missingFile(t *testing.T) {
	cmd := exec.Command(kodaBin, "check", scriptPath("tests", "does_not_exist.koda"))
	cmd.Dir = moduleRoot
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for missing file")
	}
	if cmd.ProcessState.ExitCode() != 1 {
		t.Fatalf("exit code = %d want 1", cmd.ProcessState.ExitCode())
	}
	if len(out) == 0 {
		t.Fatal("expected error message on stderr/stdout")
	}
}

func TestCLI_run_hello(t *testing.T) {
	cmd := exec.Command(kodaBin, "run", scriptPath("tests", "hello.koda"))
	cmd.Dir = moduleRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	if !bytes.Contains(out, []byte("Hello, Koda!")) {
		t.Fatalf("output missing greeting:\n%s", out)
	}
}

func TestCLI_run_importModule(t *testing.T) {
	cmd := exec.Command(kodaBin, "run", scriptPath("tests", "import_test.koda"))
	cmd.Dir = moduleRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	if !bytes.Contains(out, []byte("25")) {
		t.Fatalf("output missing square result:\n%s", out)
	}
}

func TestCLI_run_phase1Surface(t *testing.T) {
	cmd := exec.Command(kodaBin, "run", scriptPath("tests", "phase1_surface.koda"))
	cmd.Dir = moduleRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, out)
	}
	// Checksum includes summing array elements via `for-of`; see tests/phase1_surface.koda.
	if !bytes.Contains(out, []byte("98")) {
		t.Fatalf("output missing expected checksum:\n%s", out)
	}
}

func TestCLI_disasm_hello(t *testing.T) {
	cmd := exec.Command(kodaBin, "disasm", scriptPath("tests", "hello.koda"))
	cmd.Dir = moduleRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("disasm: %v\n%s", err, out)
	}
	if !bytes.Contains(out, []byte("OP_")) {
		t.Fatalf("disasm output missing opcodes:\n%s", out)
	}
}

func TestCLI_build_hello(t *testing.T) {
	if _, err := exec.LookPath("clang"); err != nil {
		t.Skip("LLVM clang not on PATH:", err)
	}

	tmpDir, err := os.MkdirTemp("", "koda-build")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	outName := "hello_out"
	if runtime.GOOS == "windows" {
		outName += ".exe"
	}
	outPath := filepath.Join(tmpDir, outName)

	cmd := exec.Command(kodaBin, "build", scriptPath("tests", "hello.koda"), "-o", outPath)
	cmd.Dir = tmpDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v\n%s", err, out)
	}

	run := exec.Command(outPath)
	run.Dir = tmpDir
	runOut, err := run.CombinedOutput()
	if err != nil {
		t.Fatalf("run built binary: %v\n%s", err, runOut)
	}
	if !bytes.Contains(runOut, []byte("Hello, Koda!")) {
		t.Fatalf("built program output missing greeting:\n%s", runOut)
	}
}

func normalizeStdoutNewlines(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

func TestCLI_run_and_build_print_spacing(t *testing.T) {
	want := "hello world\na b c\n\n"

	runVM := exec.Command(kodaBin, "run", scriptPath("tests", "print_spacing.koda"))
	runVM.Dir = moduleRoot
	vmOut, err := runVM.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, vmOut)
	}
	if normalizeStdoutNewlines(string(vmOut)) != want {
		t.Fatalf("VM print output mismatch\ngot:  %q\nwant: %q", string(vmOut), want)
	}

	if _, err := exec.LookPath("clang"); err != nil {
		t.Skip("LLVM clang not on PATH:", err)
	}

	tmpDir, err := os.MkdirTemp("", "koda-print-build")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	outName := "print_out"
	if runtime.GOOS == "windows" {
		outName += ".exe"
	}
	outPath := filepath.Join(tmpDir, outName)

	build := exec.Command(kodaBin, "build", scriptPath("tests", "print_spacing.koda"), "-o", outPath)
	build.Dir = tmpDir
	bout, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v\n%s", err, bout)
	}

	runNative := exec.Command(outPath)
	runNative.Dir = tmpDir
	nativeOut, err := runNative.CombinedOutput()
	if err != nil {
		t.Fatalf("run native: %v\n%s", err, nativeOut)
	}
	if normalizeStdoutNewlines(string(nativeOut)) != want {
		t.Fatalf("native print output mismatch VM\ngot:  %q\nwant: %q", string(nativeOut), want)
	}
}

func TestCLI_run_and_build_default_rest(t *testing.T) {
	want := "3\n3\n"

	runVM := exec.Command(kodaBin, "run", scriptPath("tests", "default_rest.koda"))
	runVM.Dir = moduleRoot
	vmOut, err := runVM.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, vmOut)
	}
	if normalizeStdoutNewlines(string(vmOut)) != want {
		t.Fatalf("VM output mismatch\ngot:  %q\nwant: %q", string(vmOut), want)
	}

	if _, err := exec.LookPath("clang"); err != nil {
		t.Skip("LLVM clang not on PATH:", err)
	}

	tmpDir, err := os.MkdirTemp("", "koda-default-rest")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	outName := "defrest_out"
	if runtime.GOOS == "windows" {
		outName += ".exe"
	}
	outPath := filepath.Join(tmpDir, outName)

	build := exec.Command(kodaBin, "build", scriptPath("tests", "default_rest.koda"), "-o", outPath)
	build.Dir = tmpDir
	bout, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v\n%s", err, bout)
	}

	runNative := exec.Command(outPath)
	runNative.Dir = tmpDir
	nativeOut, err := runNative.CombinedOutput()
	if err != nil {
		t.Fatalf("run native: %v\n%s", err, nativeOut)
	}
	if normalizeStdoutNewlines(string(nativeOut)) != want {
		t.Fatalf("native output mismatch VM\ngot:  %q\nwant: %q", string(nativeOut), want)
	}
}

func TestCLI_run_and_build_debug_minimal(t *testing.T) {
	runVM := exec.Command(kodaBin, "run", scriptPath("debug_minimal.koda"))
	runVM.Dir = moduleRoot
	vmOut, err := runVM.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, vmOut)
	}
	want := normalizeStdoutNewlines(string(vmOut))

	if _, err := exec.LookPath("clang"); err != nil {
		t.Skip("LLVM clang not on PATH:", err)
	}

	tmpDir, err := os.MkdirTemp("", "koda-debug-minimal")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	outName := "dbgmin_out"
	if runtime.GOOS == "windows" {
		outName += ".exe"
	}
	outPath := filepath.Join(tmpDir, outName)

	build := exec.Command(kodaBin, "build", scriptPath("debug_minimal.koda"), "-o", outPath)
	build.Dir = tmpDir
	bout, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v\n%s", err, bout)
	}

	runNative := exec.Command(outPath)
	runNative.Dir = tmpDir
	nativeOut, err := runNative.CombinedOutput()
	if err != nil {
		t.Fatalf("run native: %v\n%s", err, nativeOut)
	}
	if normalizeStdoutNewlines(string(nativeOut)) != want {
		t.Fatalf("native output mismatch VM\ngot:  %q\nwant: %q", string(nativeOut), want)
	}
}

func TestCLI_run_and_build_control_koda(t *testing.T) {
	runVM := exec.Command(kodaBin, "run", scriptPath("tests", "control.koda"))
	runVM.Dir = moduleRoot
	vmOut, err := runVM.CombinedOutput()
	if err != nil {
		t.Fatalf("run: %v\n%s", err, vmOut)
	}
	want := normalizeStdoutNewlines(string(vmOut))

	if _, err := exec.LookPath("clang"); err != nil {
		t.Skip("LLVM clang not on PATH:", err)
	}

	tmpDir, err := os.MkdirTemp("", "koda-control-build")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	outName := "control_out"
	if runtime.GOOS == "windows" {
		outName += ".exe"
	}
	outPath := filepath.Join(tmpDir, outName)

	build := exec.Command(kodaBin, "build", scriptPath("tests", "control.koda"), "-o", outPath)
	build.Dir = tmpDir
	bout, err := build.CombinedOutput()
	if err != nil {
		t.Fatalf("build: %v\n%s", err, bout)
	}

	runNative := exec.Command(outPath)
	runNative.Dir = tmpDir
	nativeOut, err := runNative.CombinedOutput()
	if err != nil {
		t.Fatalf("run native: %v\n%s", err, nativeOut)
	}
	if normalizeStdoutNewlines(string(nativeOut)) != want {
		t.Fatalf("native output mismatch VM\ngot:  %q\nwant: %q", string(nativeOut), want)
	}
}
