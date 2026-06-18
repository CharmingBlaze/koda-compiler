package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"koda/internal/project"
)

func runNew(args []string) error {
	name := ""
	template := "hello"
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--template", "-t":
			if i+1 >= len(args) {
				return fmt.Errorf("--template requires a name")
			}
			i++
			template = args[i]
		default:
			if strings.HasPrefix(args[i], "-") {
				return fmt.Errorf("unknown flag: %s", args[i])
			}
			if name != "" {
				return fmt.Errorf("multiple project names")
			}
			name = args[i]
		}
	}
	if name == "" {
		return fmt.Errorf("usage: koda new <name> [--template hello|game|graphics|pong|raylib]")
	}
	if err := project.ValidateProjectName(name); err != nil {
		return err
	}

	dest := filepath.Join(".", name)
	if _, err := os.Stat(dest); err == nil {
		return fmt.Errorf("directory already exists: %s", dest)
	} else if !os.IsNotExist(err) {
		return err
	}

	absDest, err := project.Scaffold(".", name, template)
	if err != nil {
		return err
	}

	fmt.Printf("Created Koda project %s (template: %s)\n", name, template)
	fmt.Printf("  %s\n", filepath.Join(name, project.FileName))
	fmt.Printf("  %s\n", filepath.Join(name, "src", "main.koda"))
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", name)
	switch template {
	case "graphics", "pong":
		fmt.Println("  koda run    # use raylib + koda.game")
	case "raylib":
		fmt.Println("  koda run    # full Raylib API (#include \"@raylib\")")
	case "game":
		fmt.Println("  koda run    # text lunar lander — no extra libraries")
	default:
		fmt.Println("  koda run")
	}
	_ = absDest
	return nil
}

func listTemplates() []string {
	return project.ListTemplates()
}
