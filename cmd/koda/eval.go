package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"koda/api"
	"koda/internal/nativebuild"
)

func runEval(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: koda eval '<koda code>'")
	}
	code := strings.Join(args, " ")
	return evalSnippet(code)
}

func runRepl() error {
	fmt.Println("Koda REPL — enter expressions or statements (; optional). Type exit or Ctrl+C to quit.")
	reader := bufio.NewReader(os.Stdin)
	opts := nativebuild.BuildOptions{NoOpt: true}
	for {
		fmt.Print("koda> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println()
			return nil
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "exit" || line == "quit" {
			return nil
		}
		if err := evalSnippetWithOpts(line, opts); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
	}
}

func evalSnippet(code string) error {
	return evalSnippetWithOpts(code, nativebuild.BuildOptions{NoOpt: true})
}

func evalSnippetWithOpts(code string, opts nativebuild.BuildOptions) error {
	wrapped := wrapEvalCode(code)
	tmp, err := os.CreateTemp("", "koda_eval_*.koda")
	if err != nil {
		return err
	}
	path := tmp.Name()
	defer func() {
		_ = os.Remove(path)
	}()
	if err := os.WriteFile(path, []byte(wrapped), 0644); err != nil {
		return err
	}
	return api.RunWithWritersOptsProgram(path, "", os.Stdout, os.Stderr, opts, nil)
}

func wrapEvalCode(line string) string {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return "print(\"\");"
	}
	// Expression-like: no semicolon and not starting with keyword statement
	if !strings.HasSuffix(trimmed, ";") && !looksLikeStatement(trimmed) {
		return fmt.Sprintf("print(%s);", trimmed)
	}
	if !strings.HasSuffix(trimmed, ";") {
		return trimmed + ";"
	}
	return trimmed
}

func looksLikeStatement(s string) bool {
	lower := strings.ToLower(s)
	prefixes := []string{"let ", "func ", "if ", "while ", "for ", "switch ", "return ", "print ", "#include", "struct ", "enum "}
	for _, p := range prefixes {
		if strings.HasPrefix(lower, p) {
			return true
		}
	}
	return strings.Contains(s, "{")
}
