package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"koda/internal/codegen"
	"koda/internal/lexer"
	"koda/internal/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: genir <file.koda> [output.ll]")
		os.Exit(1)
	}

	sourcePath, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
		os.Exit(1)
	}
	source, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	l := lexer.NewLexer(string(source), sourcePath)
	tokens, err := l.Tokenize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Lexer error: %v\n", err)
		os.Exit(1)
	}

	p := parser.NewParser(tokens)
	program, err := p.Parse()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parser error: %v\n", err)
		os.Exit(1)
	}

	bundle := &parser.ProgramBundle{Entry: program}
	ctx, err := codegen.PrepareNativeBundle(bundle)
	if err != nil {
		fmt.Fprintf(os.Stderr, "PrepareNativeBundle error: %v\n", err)
		os.Exit(1)
	}

	mod, err := codegen.EmitLLVMIR(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "EmitLLVMIR error: %v\n", err)
		os.Exit(1)
	}

	output := "output.ll"
	if len(os.Args) >= 3 {
		output = os.Args[2]
	}

	if err := ioutil.WriteFile(output, []byte(mod.String()), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated: %s\n", output)
}
