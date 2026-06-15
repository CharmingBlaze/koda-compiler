package project

import (
	"runtime"
	"strings"
)

// DefaultGraphicsLinkFlags returns platform Raylib link flags for beginners (no manual env vars).
func DefaultGraphicsLinkFlags() string {
	parts := []string{"-lraylib"}
	switch runtime.GOOS {
	case "windows":
		parts = append(parts, "-lopengl32", "-lgdi32", "-lwinmm")
	case "darwin":
		parts = append(parts, "-framework", "OpenGL", "-framework", "Cocoa", "-framework", "IOKit", "-framework", "CoreVideo")
	case "linux":
		parts = append(parts, "-lGL", "-lm", "-lpthread", "-ldl", "-lrt", "-lX11")
	}
	return strings.Join(parts, " ")
}

// LintMode returns normalized lint setting from koda.json ("", "beginner", "strict").
func (c *Context) LintMode() string {
	if c == nil || c.Cfg == nil {
		return ""
	}
	switch strings.ToLower(strings.TrimSpace(c.Cfg.Lint)) {
	case "beginner", "strict":
		return strings.ToLower(strings.TrimSpace(c.Cfg.Lint))
	default:
		return ""
	}
}

// PrepareOptionsForCheck maps project lint settings to sema options.
func (c *Context) PrepareOptionsForCheck(defaultWarnUnused bool) (warnUnused, beginnerLint, warnUnreachable bool) {
	warnUnused = defaultWarnUnused
	if c == nil {
		return warnUnused, beginnerLint, warnUnreachable
	}
	switch c.LintMode() {
	case "beginner":
		warnUnused = true
		beginnerLint = true
	case "strict":
		warnUnused = true
		warnUnreachable = true
	}
	return warnUnused, beginnerLint, warnUnreachable
}
