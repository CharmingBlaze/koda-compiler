package api

import (
	"os"

	"koda/internal/project"
)

// withProjectScope applies koda.json native env and chdirs to the project root for the duration of fn,
// matching the CLI behaviour of `koda run` / `koda build`.
func withProjectScope(entryPath string, fn func() error) error {
	ctx, err := project.LoadContext(entryPath)
	if err != nil {
		return err
	}
	if ctx == nil {
		return fn()
	}
	if err := ctx.ApplyNativeEnv(); err != nil {
		return err
	}
	prev, err := os.Getwd()
	if err != nil {
		return err
	}
	if err := os.Chdir(ctx.Root); err != nil {
		return err
	}
	defer func() { _ = os.Chdir(prev) }()
	return fn()
}
