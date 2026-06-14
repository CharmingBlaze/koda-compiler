package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"koda/internal/project"
)

func cwd() string {
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	return dir
}

func withProject(entryPath string, fn func() error) error {
	ctx, err := project.LoadContext(entryPath)
	if err != nil {
		return err
	}
	if ctx != nil {
		if err := ctx.ApplyNativeEnv(); err != nil {
			return err
		}
		if err := os.Chdir(ctx.Root); err != nil {
			return err
		}
	}
	return fn()
}

func projectContextFor(entryPath string) (*project.Context, error) {
	ctx, err := project.LoadContext(entryPath)
	if err != nil {
		return nil, err
	}
	if ctx != nil {
		if err := ctx.ApplyNativeEnv(); err != nil {
			return nil, err
		}
	}
	return ctx, nil
}

func resolveEntry(explicit string) (string, error) {
	return project.ResolveEntry(cwd(), explicit)
}

func defaultBuildOutput(entryPath string, ctx *project.Context) string {
	name := "app"
	if ctx != nil {
		name = ctx.AppName(entryPath)
	} else {
		name = strings.TrimSuffix(filepath.Base(entryPath), filepath.Ext(entryPath))
	}
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}
