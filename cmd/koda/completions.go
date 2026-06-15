package main

import (
	"fmt"
	"strings"
)

func runCompletions(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: koda completions bash|zsh|fish")
	}
	switch args[0] {
	case "bash":
		fmt.Print(bashCompletionScript())
	case "zsh":
		fmt.Print(zshCompletionScript())
	case "fish":
		fmt.Print(fishCompletionScript())
	default:
		return fmt.Errorf("unknown shell %q (use bash, zsh, or fish)", args[0])
	}
	return nil
}

func kodaCommandsList() string {
	cmds := []string{
		"new", "init", "run", "native", "watch", "check", "lint", "fmt",
		"build", "bundle", "test", "bench", "profile", "debug", "eval", "repl",
		"clean", "doctor", "setup", "paths", "env", "completions", "update", "doc", "lsp",
		"disasm", "wrap", "version", "help",
	}
	return strings.Join(cmds, " ")
}

func bashCompletionScript() string {
	return `# koda bash completion
_koda_completion() {
  local cur prev
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"
  local commands="` + kodaCommandsList() + `"
  if [[ ${COMP_CWORD} -eq 1 ]]; then
    COMPREPLY=($(compgen -W "${commands}" -- "${cur}"))
    return
  fi
  case "${COMP_WORDS[1]}" in
    new|init)
      COMPREPLY=($(compgen -W "--template hello game graphics pong raylib" -- "${cur}"))
      ;;
    run|native|watch|build|check|bench|profile|debug|disasm|bundle|test|eval)
      COMPREPLY=($(compgen -f -X '!*.koda' -- "${cur}"))
      ;;
    fmt|lint|check)
      COMPREPLY=($(compgen -f -X '!*.koda' -- "${cur}"))
      COMPREPLY+=($(compgen -W "./..." -- "${cur}"))
      ;;
    completions)
      COMPREPLY=($(compgen -W "bash zsh fish" -- "${cur}"))
      ;;
    doc)
      COMPREPLY=($(compgen -W "stdlib module" -- "${cur}"))
      ;;
  esac
}
complete -F _koda_completion koda
`
}

func zshCompletionScript() string {
	return `#compdef koda
_koda() {
  local -a commands
  commands=(` + strings.ReplaceAll(kodaCommandsList(), " ", "\n") + `)
  _arguments -C \
    '1:command:->command' \
    '*:arg:->args'
  case $state in
    command) _describe 'command' commands ;;
    args)
      case $words[1] in
        completions) _values 'shell' bash zsh fish ;;
        doc) _values 'subcommand' stdlib module ;;
      esac
  esac
}
_koda
`
}

func fishCompletionScript() string {
	return `complete -c koda -n "__fish_use_subcommand" -a "` + kodaCommandsList() + `"
complete -c koda -n "__fish_seen_subcommand_from completions" -a "bash zsh fish"
complete -c koda -n "__fish_seen_subcommand_from doc" -a "stdlib module"
`
}
