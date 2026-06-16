package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"koda/api"
	"koda/internal/codegen"
	"koda/internal/diagnostic"
	"koda/internal/kodahome"
	"koda/internal/nativebuild"
	"koda/internal/parser"
	"koda/internal/project"
	"koda/internal/sema"
)

func parseBuildCommandArgs(args []string) (src string, out string, noOpt bool, debug bool, err error) {
	var output string
	var file string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--no-opt":
			noOpt = true
		case "--debug":
			debug = true
		case "-o":
			if i+1 >= len(args) {
				return "", "", false, false, fmt.Errorf("-o requires a path")
			}
			i++
			output = args[i]
		default:
			if strings.HasPrefix(args[i], "-") {
				return "", "", false, false, fmt.Errorf("unknown flag: %s", args[i])
			}
			if file != "" {
				return "", "", false, false, fmt.Errorf("multiple source files")
			}
			file = args[i]
		}
	}
	if file == "" {
		if output == "" {
			return "", "", noOpt, debug, nil
		}
		return "", "", false, false, fmt.Errorf("usage: koda build [--no-opt] [--debug] [<file.koda>] [-o <exe>]")
	}
	if output == "" {
		output = defaultExeName(file)
	}
	return file, output, noOpt, debug, nil
}

func parseRunCommandArgs(args []string) (src string, noOpt bool, debug bool, progArgs []string, err error) {
	before, after := splitAtDoubleDash(args)
	progArgs = after
	var file string
	for i := 0; i < len(before); i++ {
		switch before[i] {
		case "--no-opt":
			noOpt = true
		case "--debug":
			debug = true
		default:
			if strings.HasPrefix(before[i], "-") {
				return "", false, false, nil, fmt.Errorf("unknown flag: %s", before[i])
			}
			if file != "" {
				return "", false, false, nil, fmt.Errorf("multiple source files")
			}
			file = before[i]
		}
	}
	if file == "" {
		return "", noOpt, debug, progArgs, nil
	}
	return file, noOpt, debug, progArgs, nil
}

// version is set by release builds, e.g. -ldflags "-X main.version=1.0.0"
var version = "0.6.0-dev"

func main() {
	kodahome.EnsureSDKFromExecutable()
	args := os.Args[1:]
	if len(args) == 0 {
		printHelp()
		os.Exit(0)
	}

	switch args[0] {
	case "run", "native":
		if maybeCommandHelp(args) {
			return
		}
		path, noOpt, debug, progArgs, err := parseRunCommandArgs(args[1:])
		if err != nil {
			fatal(err.Error())
		}
		if path == "" {
			path, err = resolveEntry("")
			if err != nil {
				fatal(err.Error())
			}
		}
		opts := nativebuild.BuildOptions{NoOpt: noOpt || debug, Debug: debug}
		if err := withProject(path, func() error {
			return api.RunWithBuildOptionsProgram(path, "", opts, progArgs)
		}); err != nil {
			fatalErr(err)
		}

	case "watch":
		if maybeCommandHelp(args) {
			return
		}
		path, noOpt, debug, progArgs, err := parseWatchCommandArgs(args[1:])
		if err != nil {
			fatal(err.Error())
		}
		if path == "" {
			path, err = resolveEntry("")
			if err != nil {
				fatal(err.Error())
			}
		}
		opts := nativebuild.BuildOptions{NoOpt: noOpt || debug, Debug: debug}
		if err := withProject(path, func() error {
			return runWatch(path, opts, progArgs)
		}); err != nil {
			fatalErr(err)
		}

	case "check":
		if maybeCommandHelp(args) {
			return
		}
		rest := args[1:]
		if len(rest) == 0 {
			path, err := resolveEntry("")
			if err != nil {
				fatal(err.Error())
			}
			if err := withProject(path, func() error {
				return checkFile(path)
			}); err != nil {
				fatalErr(err)
			}
			fmt.Println("OK")
			return
		}
		if len(rest) == 1 && (rest[0] == "./..." || rest[0] == "...") {
			if err := runCheckAll(rest); err != nil {
				fatalErr(err)
			}
			return
		}
		if len(rest) == 1 && !strings.HasPrefix(rest[0], "-") {
			path := rest[0]
			if err := withProject(path, func() error {
				return checkFile(path)
			}); err != nil {
				fatalErr(err)
			}
			fmt.Println("OK")
			return
		}
		if err := runCheckAll(rest); err != nil {
			fatalErr(err)
		}

	case "lint":
		if maybeCommandHelp(args) {
			return
		}
		if err := runLint(args[1:]); err != nil {
			fatalErr(err)
		}

	case "new":
		if maybeCommandHelp(args) {
			return
		}
		if err := runNew(args[1:]); err != nil {
			fatalErr(err)
		}

	case "init":
		if maybeCommandHelp(args) {
			return
		}
		if err := runNew(args[1:]); err != nil {
			fatalErr(err)
		}

	case "fmt":
		if maybeCommandHelp(args) {
			return
		}
		if err := runFmtCmd(args[1:]); err != nil {
			fatalErr(err)
		}

	case "disasm":
		if maybeCommandHelp(args) {
			return
		}
		requireArg(args, "disasm", "<file.koda>")
		if err := disasmFile(args[1]); err != nil {
			fatalErr(err)
		}

	case "build":
		if maybeCommandHelp(args) {
			return
		}
		src, output, noOpt, debug, err := parseBuildCommandArgs(args[1:])
		if err != nil {
			fatal(err.Error())
		}
		if src == "" {
			src, err = resolveEntry("")
			if err != nil {
				fatal(err.Error())
			}
		}
		ctx, err := projectContextFor(src)
		if err != nil {
			fatalErr(err)
		}
		if output == "" {
			output = defaultBuildOutput(src, ctx)
		}
		opts := nativebuild.BuildOptions{NoOpt: noOpt || debug, Debug: debug}
		if err := withProject(src, func() error {
			return buildFileOpts(src, output, opts)
		}); err != nil {
			fatalErr(err)
		}

	case "bundle":
		if maybeCommandHelp(args) {
			return
		}
		src, outputDir, err := parseBundleArgs(args[1:])
		if err != nil {
			fatal(err.Error())
		}
		if src == "" {
			src, err = resolveEntry("")
			if err != nil {
				fatal(err.Error())
			}
		}
		if err := withProject(src, func() error {
			return bundleFile(src, outputDir)
		}); err != nil {
			fatalErr(err)
		}

	case "wrap":
		if err := runWrapgen(args[1:]); err != nil {
			fatalErr(err)
		}

	case "paths":
		if err := printResolvedPaths(); err != nil {
			fatalErr(err)
		}

	case "doctor":
		if maybeCommandHelp(args) {
			return
		}
		if err := runDoctor(args[1:]); err != nil {
			fatalErr(err)
		}

	case "setup":
		if maybeCommandHelp(args) {
			return
		}
		if err := runSetup(args[1:]); err != nil {
			fatalErr(err)
		}

	case "clean":
		if maybeCommandHelp(args) {
			return
		}
		if err := runClean(args[1:]); err != nil {
			fatalErr(err)
		}

	case "test":
		if maybeCommandHelp(args) {
			return
		}
		if err := runTest(args[1:]); err != nil {
			fatalErr(err)
		}

	case "bench":
		if maybeCommandHelp(args) {
			return
		}
		if err := runBench(args[1:]); err != nil {
			fatalErr(err)
		}

	case "profile":
		if maybeCommandHelp(args) {
			return
		}
		if err := runBench(args[1:]); err != nil {
			fatalErr(err)
		}

	case "debug":
		if maybeCommandHelp(args) {
			return
		}
		path, _, _, progArgs, err := parseRunCommandArgs(args[1:])
		if err != nil {
			fatal(err.Error())
		}
		if path == "" {
			path, err = resolveEntry("")
			if err != nil {
				fatal(err.Error())
			}
		}
		opts := nativebuild.BuildOptions{NoOpt: true, Debug: true}
		if err := withProject(path, func() error {
			return api.RunWithBuildOptionsProgram(path, "", opts, progArgs)
		}); err != nil {
			fatalErr(err)
		}

	case "eval":
		if maybeCommandHelp(args) {
			return
		}
		if err := runEval(args[1:]); err != nil {
			fatalErr(err)
		}

	case "repl":
		if maybeCommandHelp(args) {
			return
		}
		if err := runRepl(); err != nil {
			fatalErr(err)
		}

	case "env":
		if maybeCommandHelp(args) {
			return
		}
		if err := runEnv(args[1:]); err != nil {
			fatalErr(err)
		}

	case "completions":
		if err := runCompletions(args[1:]); err != nil {
			fatalErr(err)
		}

	case "update":
		if maybeCommandHelp(args) {
			return
		}
		if err := runUpdate(args[1:]); err != nil {
			fatalErr(err)
		}

	case "doc":
		if maybeCommandHelp(args) {
			return
		}
		if err := runDoc(args[1:]); err != nil {
			fatalErr(err)
		}

	case "lsp":
		if maybeCommandHelp(args) {
			return
		}
		if err := runLSP(); err != nil {
			fatalErr(err)
		}

	case "version", "--version", "-v":
		fmt.Println(versionLine())

	case "help", "--help", "-h":
		if len(args) > 1 {
			printCommandHelp(args[1])
			return
		}
		printHelp()

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", args[0])
		printHelp()
		os.Exit(1)
	}
}

func versionLine() string {
	return fmt.Sprintf("koda %s (%s/%s)", version, runtime.GOOS, runtime.GOARCH)
}

func printHelp() {
	fmt.Printf(`%s

Koda compiles .koda programs to native apps. Your players only need the bundle folder
(executable + assets). They do not need Go, Python, or C++ toolchains.

Release builds from GitHub embed the native compiler (Clang, llc, linker bits, runtime lib).
End users install only **koda** and **kodawrap** (and ship **stdlib/** with them). They do **not** install Go, LLVM, or any other compiler for normal **build / run / wrap**.

USAGE
  koda <command> [arguments]

COMMANDS
  new, init   <name> [--template hello|game|graphics|pong|raylib]   Create a project
  run, native [--no-opt] [--debug] [<file.koda>] [-- <program args>]
  watch       [--no-opt] [--debug] [<file.koda>] [-- <program args>]
  check       [<file.koda> | ./...]                     Parse + sema
  lint        [<paths...>] [./...]                      check + fmt --check
  fmt         [--check] <files...> [./...]
  build       [--no-opt] [--debug] [<file.koda>] [-o <exe>]
  bundle      [<file.koda>] [-o <dir>]
  test        [--no-opt] [-v] [--failfast] [-run <pattern>] [<files...>]
  bench       [--count N] [--warmup N] <file.koda> [-- <args>]
  profile     Same as bench
  debug       Run with debug symbols (like run --debug)
  eval        '<koda code>'                             One-shot snippet
  repl        Interactive compile-per-line REPL
  clean       [<dir>] [--cache]
  doctor      SDK health check
  setup       raylib | …                            Project setup helpers
  paths       Machine-readable toolchain paths
  env         [--export]                              KODA_* environment
  completions bash|zsh|fish
  update      [--check-only]                          Release check hints
  doc         [stdlib | module <@name>]               Doc helpers
  lsp         Language server (stdio JSON-RPC)
  disasm      <file.koda>                             LLVM IR
  wrap        ...                                     Forward to kodawrap

  help        [command]                               Command help
  version     Version and platform

OPTIONS
  -o <path>   Output executable (build) or output directory (bundle)

MAINTAINER BUILD (not for end users)
  Building Koda from source is for contributors only. See CONTRIBUTING.md in the repository.

C / C++ LIBRARIES (readable .koda wrappers)
  Use the kodawrap binary from Releases (or next to koda). Generate bindings from headers:

    koda wrap -name mylib -headers ./include/mylib.h -out ./wrappers/mylib

  That writes readable mylib.koda plus wrapper.c and docs. To compile your game:

    set KODA_NATIVE_SOURCES=wrappers\mylib\wrapper.c
    set KODA_LINKFLAGS=-I.\include -L.\lib -lmylib
    koda bundle game.koda -o dist\mygame

  See docs/wrappers.md and docs/distribution.md.

ENVIRONMENT
  KODA_CLANG / CC     Override Clang path (only when not using the embedded release toolchain; contributors / custom SDKs)
  KODA_LLC            Override llc path (same; release builds embed llc)
  KODA_USE_LLD        1 to force -fuse-ld=lld; 0 to disable even if bundled ld.lld exists
  KODA_PATH           Extra @module search dirs (path list, same separator as PATH)
  KODA_WRAPPERS       Pre-built .koda libraries (path list; overrides KODA_PATH)
  KODA_NATIVE_SOURCES C/C++ sources linked into your app (e.g. wrapper.c)
  KODA_LINKFLAGS      Extra linker flags (-lraylib, -L..., frameworks, etc.)
  KODA_RAYLIB_STAGE   Override path to third_party/.../stage (include/ + lib/) for vendored raylib
  KODA_USE_VENDORED_RAYLIB  If third_party/raylib_static/stage exists (cwd or next to koda), prepends -I and links libraylib.a; set 0/false to skip
  KODA_BUNDLE_FILES   Extra files/dirs copied into the bundle (path-list or quoted paths with spaces)
  KODA_SKIP_TOOLCHAIN_EXTRACT  If set, skip unpacking the embedded compiler (then set KODA_CLANG or use a dev build)
  KODA_DEBUG_IR       If set, writes .KODA_build/main.ll in addition to piping IR to clang

SINGLE-EXE / EMBEDDED TOOLCHAIN
  Release builds embed Clang, llc, the linker, and the Koda runtime archive inside the executable.
  On first koda build or kodawrap use, they are unpacked to a writable cache — no LLVM install on the machine.

ZERO-SETUP DISTRIBUTION (same folder as koda.exe / koda)
  Ship stdlib (and optional wrappers/) next to the release binary so @ imports resolve with no extra install.
  You do not need a separate llvm/ folder when using a proper release build (embedded toolchain).
  KODA_CLANG / KODA_PATH are only for overriding defaults (e.g. contributor builds without embed).

EXAMPLES
  koda new mygame
  cd mygame && koda run
  koda run tests\hello.koda
  koda run --no-opt tests\loops.koda
  koda build game.koda -o game%s
  koda bundle game.koda -o dist\mygame
`, versionLine(), exeExt())
}

func exeExt() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

func runWrapgen(args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Println("koda wrap — runs kodawrap (header → .koda + wrapper.c). Example:")
		fmt.Println("  koda wrap -name mylib -headers ./include/mylib.h -out ./wrappers/mylib")
		fmt.Println("\nInstall kodawrap from the same GitHub Releases page as koda, or place kodawrap next to koda.")
		fmt.Println("Release kodawrap embeds the same Clang as koda — no separate LLVM install for header parsing.")
		fmt.Println("(Legacy names wrapgen / kujiwrap are still discovered if present.)")
		return nil
	}
	names := []string{
		"kodawrap", "kodawrap.exe",
		"wrapgen", "wrapgen.exe",
		"kujiwrap", "kujiwrap.exe",
	}
	if self, err := os.Executable(); err == nil {
		dir := filepath.Dir(self)
		for _, name := range names {
			p := filepath.Join(dir, name)
			if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
				return runPassthrough(p, args)
			}
		}
	}
	for _, name := range names {
		if p, err := exec.LookPath(name); err == nil {
			return runPassthrough(p, args)
		}
	}
	return fmt.Errorf("kodawrap not found (looked next to koda and on PATH).\nDownload kodawrap from GitHub Releases (same page as koda), or add it next to this executable.\n(legacy binary name wrapgen also works.)")
}

func runPassthrough(path string, args []string) error {
	cmd := exec.Command(path, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if x, ok := err.(*exec.ExitError); ok && x.ExitCode() != 0 {
			os.Exit(x.ExitCode())
		}
		return err
	}
	return nil
}

func requireArg(args []string, cmd, arg string) {
	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: koda %s %s\n", cmd, arg)
		os.Exit(1)
	}
}

func defaultExeName(source string) string {
	name := strings.TrimSuffix(filepath.Base(source), filepath.Ext(source))
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func checkFile(path string) error {
	return checkFileWithOptions(path, true)
}

func checkFileWithOptions(path string, warnUnused bool) error {
	bundle, err := parser.LoadProgram(path)
	if err != nil {
		return err
	}
	beginnerLint := false
	warnUnreachable := false
	if ctx, err := project.LoadContext(path); err == nil && ctx != nil {
		warnUnused, beginnerLint, warnUnreachable = ctx.PrepareOptionsForCheck(warnUnused)
	}
	opts := &sema.PrepareOptions{WarnUnused: warnUnused, BeginnerLint: beginnerLint, WarnUnreachable: warnUnreachable}
	_, err = sema.PrepareNativeBundleWithOptions(bundle, opts)
	return err
}
func disasmFile(path string) error {
	bundle, err := parser.LoadProgram(path)
	if err != nil {
		return err
	}
	ctx, err := codegen.PrepareNativeBundle(bundle)
	if err != nil {
		return err
	}
	mod, err := codegen.EmitLLVMIR(ctx)
	if err != nil {
		return err
	}
	fmt.Print(mod.String())
	return nil
}
func buildFileOpts(path string, output string, opts nativebuild.BuildOptions) error {
	bundle, err := parser.LoadProgram(path)
	if err != nil {
		return err
	}
	return nativebuild.BuildWithOptions(bundle, output, filepath.Base(path), func(s string) { fmt.Print(s) }, opts)
}

// printResolvedPaths prints install-relative toolchain and stdlib resolution for
// debugging portable releases (no KODA_HOME — layout is always relative to the binary).
func printResolvedPaths() error {
	install, err := kodahome.InstallDir()
	if err != nil {
		return err
	}
	fmt.Printf("KODA_INSTALL_DIR=%s\n", install)

	clang := kodahome.Clang()
	fmt.Printf("KODA_CLANG_EFFECTIVE=%s\n", clang)
	fmt.Printf("KODA_CLANG_STATUS=%s\n", toolResolveStatus(clang))

	llc := kodahome.LLC()
	fmt.Printf("KODA_LLC_EFFECTIVE=%s\n", llc)
	fmt.Printf("KODA_LLC_STATUS=%s\n", toolResolveStatus(llc))

	if kodahome.HasBundledLLD() {
		fmt.Println("KODA_BUNDLED_LLD=1")
	} else {
		fmt.Println("KODA_BUNDLED_LLD=0")
	}

	stdlib, err := kodahome.StdlibDir()
	if err != nil {
		return err
	}
	fmt.Printf("stdlib_dir=%s\n", stdlib)
	if fi, err := os.Stat(stdlib); err == nil && fi.IsDir() {
		fmt.Println("stdlib_exists=1")
	} else {
		fmt.Println("stdlib_exists=0")
	}

	wrap, err := kodahome.WrappersDir()
	if err != nil {
		return err
	}
	fmt.Printf("wrappers_dir=%s\n", wrap)
	if fi, err := os.Stat(wrap); err == nil && fi.IsDir() {
		fmt.Println("wrappers_exists=1")
	} else {
		fmt.Println("wrappers_exists=0")
	}

	flags, err := kodahome.BundledClangResourceFlags()
	if err != nil {
		fmt.Printf("bundled_clang_include_flags_error=%v\n", err)
	} else {
		fmt.Printf("bundled_clang_include_flags_count=%d\n", len(flags))
	}

	for _, k := range []string{
		"KODA_CLANG", "KODA_LLC", "CC", "KODA_PATH", "KODA_WRAPPERS",
		"KODA_SKIP_TOOLCHAIN_EXTRACT", "KODA_USE_LLD",
	} {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" {
			fmt.Printf("%s=%s\n", k, v)
		}
	}
	fmt.Println("hint: for a human-readable report (writable install dir, tool probes), run: koda doctor")
	return nil
}

func toolResolveStatus(tool string) string {
	if tool == "" {
		return "empty"
	}
	if filepath.IsAbs(tool) || strings.ContainsRune(tool, filepath.Separator) ||
		(runtime.GOOS == "windows" && strings.ContainsRune(tool, '/')) {
		fi, err := os.Stat(tool)
		if err == nil && !fi.IsDir() {
			return "ok"
		}
		return "missing"
	}
	if _, err := exec.LookPath(tool); err == nil {
		return "on_path"
	}
	return "missing"
}

func fatal(msg string) {
	fmt.Fprintf(os.Stderr, "\nerror: %s\n\nRun 'koda help' for usage.\n", msg)
	os.Exit(1)
}

func fatalErr(err error) {
	fatal(diagnostic.FormatError(err))
}
