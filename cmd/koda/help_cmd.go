package main

import (
	"fmt"
	"os"
)

var commandHelp = map[string]string{
	"new": `Create a project directory with koda.json and src/main.koda.

  koda new <name> [--template hello|game|graphics|pong|raylib]

Templates:
  hello     Minimal print program (default)
  game      Text lunar lander
  graphics  Full Raylib wrapper + optional koda.game (bouncing ball demo)
  raylib    Full Raylib API sample (raw InitWindow / DrawText)
  pong      Two-player Pong with koda.game
`,
	"init": `Alias for koda new.

  koda init <name> [--template hello|game|graphics|pong|raylib]
`,
	"run": `Compile to a native binary and run it (temporary executable).

  koda run [--no-opt] [--release] [--debug] [<file.koda>] [-- <program args...>]

Uses koda.json entry when no file is given. Arguments after -- are passed to your program.
`,
	"native": `Same as koda run (backward-compatible alias).
`,
	"watch": `Rebuild and rerun when .koda files change under the entry directory.

  koda watch [--no-opt] [--release] [--debug] [<file.koda>] [-- <program args...>]
`,
	"check": `Parse and type-check without codegen.

  koda check [<file.koda>]
  koda check ./...
`,
	"lint": `Run check + fmt --check on sources (CI-friendly).

  koda lint [<paths...>] [./...]
`,
	"fmt": `Canonical formatting (4-space indent).

  koda fmt [--check] <file.koda> [more...] [./...]
`,
	"build": `Build a native executable.

  koda build [--release] [--no-opt] [--debug] [<file.koda>] [-o <exe>]

Defaults to -O3 optimization unless --no-opt or --debug is set.
`,
	"bundle": `Build and package executable + assets for distribution.

  koda bundle [<file.koda>] [-o <dir>]
`,
	"test": `Run .koda test files (assert/expect or test "name" { } blocks).

  koda test [--no-opt] [-v] [--failfast] [-run <pattern>] [<files...>]

Default: discovers *_test.koda under the project (Go-style), or tests/*.koda in the compiler repo.
Uses the native compile-and-run pipeline (same as koda run).
`,
	"bench": `Time repeated runs of a program.

  koda bench [--count N] [--warmup N] [--no-opt] <file.koda> [-- <args...>]
`,
	"profile": `Alias for koda bench (execution timing).

  koda profile [--count N] <file.koda> [-- <args...>]
`,
	"debug": `Run with debug symbols (koda run --debug).

  koda debug [<file.koda>] [-- <args...>]
`,
	"eval": `Evaluate a one-line Koda expression or snippet.

  koda eval 'print(1 + 2)'
  koda eval 'let x = 3; print(x)'
`,
	"repl": `Interactive read-eval-print loop (compile per line).

  koda repl
`,
	"clean": `Remove build artifacts.

  koda clean [<dir>] [--cache]

Removes dist/, .koda_build/, and default executables. --cache also clears temp toolchain dirs when possible.
`,
	"doctor": `Human-readable SDK health check (stdlib, clang, writable install, Raylib wrapper).

  koda doctor [--fix]

  --fix   Refresh an outdated project raylib_shim from the SDK (legacy --shim projects only).
`,
	"setup": `Configure optional project integrations.

  koda setup raylib [--shim] [project-dir]

Writes koda.json native.graphics + native.sources.
  Default: full SDK Raylib wrapper (548 functions, use raylib).
  --shim:  legacy ~33-function shim copied into the project.
`,
	"paths": `Machine-readable toolchain paths for scripts/CI.
`,
	"env": `Print Koda environment variables and resolved paths.

  koda env
  koda env --export    # shell export lines
`,
	"completions": `Generate shell completion scripts.

  koda completions bash
  koda completions zsh
  koda completions fish
`,
	"update": `Check GitHub Releases for a newer Koda version.

  koda update [--check-only]
`,
	"doc": `Documentation helpers.

  koda doc                    List bundled doc paths
  koda doc stdlib             Stdlib module index
  koda doc module <@name>     Show stdlib source
  koda doc wrappers           List C library wrapper packages
  koda doc wrapper <@name>    Show wrapper META and doc paths
`,
	"lsp": `Language Server Protocol (stdio JSON-RPC) for editors.

  koda lsp

Supports initialize, textDocument/didOpen, didChange, and publishDiagnostics.
`,
	"disasm": `Print LLVM IR after codegen.

  koda disasm <file.koda>
`,
	"wrap": `Forward to kodawrap (C header → organized package).

  koda wrap -name mylib -headers ./include/mylib.h -I ./include -L ./lib -l mylib -out ./wrappers/mylib
  koda wrap -name mylib -I ./include -l mylib -o wrappers/mylib ./include/mylib.h

Regenerate after upgrading a native library:
  koda wrap upgrade wrappers/mylib
  koda wrap upgrade @raylib
  koda wrap check wrappers/mylib

Install from the built-in catalog (known link flags + headers):
  koda wrap list
  koda wrap install raylib --project
  koda wrap install sqlite3@3 -o wrappers/sqlite3

C++ headers: add --cpp or use .hpp (wrapgen uses -std=c++17).

Output: .koda, wrapper.c, README, api_reference, examples, koda.json, META.json, docs/index.html
See docs/guides/wrapping-libraries.md
`,
	"version": `Print version and platform.
`,
	"help": `Show help. Use koda <command> --help for command details.
`,
}

func printCommandHelp(cmd string) {
	text, ok := commandHelp[cmd]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", cmd)
		printHelp()
		os.Exit(1)
	}
	fmt.Printf("koda %s\n\n%s", cmd, text)
}

func maybeCommandHelp(args []string) bool {
	if len(args) >= 2 && (args[1] == "--help" || args[1] == "-h") {
		printCommandHelp(args[0])
		return true
	}
	return false
}
