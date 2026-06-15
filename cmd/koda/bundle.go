package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"koda/internal/nativebuild"
	"koda/internal/project"
)

func bundleFile(path string, outputDir string) error {
	ctx, err := projectContextFor(path)
	if err != nil {
		return err
	}
	appName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	if ctx != nil {
		appName = ctx.AppName(path)
	}
	exeName := appName
	if runtime.GOOS == "windows" {
		exeName += ".exe"
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	fmt.Printf("Bundling %s -> %s/\n\n", filepath.Base(path), outputDir)

	exePath := filepath.Join(outputDir, exeName)
	if err := buildFileOpts(path, exePath, nativebuild.BuildOptions{}); err != nil {
		return err
	}
	if err := writeBundleLauncher(outputDir, exeName); err != nil {
		return err
	}
	if err := writeBundleReadme(outputDir, appName, exeName); err != nil {
		return err
	}
	if err := writeBundleInfo(outputDir, path); err != nil {
		return err
	}
	if err := copyBundleExtraFiles(outputDir); err != nil {
		return err
	}
	if ctx != nil {
		if err := copyBundleProjectAssets(outputDir, ctx); err != nil {
			return err
		}
		if err := writeBundleAssetsManifest(outputDir, ctx); err != nil {
			return err
		}
	}

	fmt.Println()
	fmt.Printf("  Bundle ready: %s\n", outputDir)
	fmt.Printf("  Share the whole folder - friends only run %s (or the launcher script).\n", exeName)
	fmt.Printf("  No Go, Python, or C++ needed on their machine.\n")
	return nil
}

func copyBundleProjectAssets(outputDir string, ctx *project.Context) error {
	for _, item := range ctx.BundleExtraPaths() {
		info, err := os.Stat(item)
		if err != nil {
			return fmt.Errorf("bundle asset %s: %w", item, err)
		}
		base := filepath.Base(item)
		dst := filepath.Join(outputDir, base)
		if info.IsDir() {
			if err := copyDir(item, dst); err != nil {
				return fmt.Errorf("copy bundle dir %s: %w", item, err)
			}
		} else {
			if err := copyFile(item, dst); err != nil {
				return fmt.Errorf("copy bundle file %s: %w", item, err)
			}
		}
		fmt.Printf("  copied: %s\n", dst)
	}
	return nil
}

func writeBundleLauncher(outputDir, exeName string) error {
	if runtime.GOOS == "windows" {
		content := fmt.Sprintf("@echo off\r\n\"%%~dp0%s\" %%*\r\n", exeName)
		return os.WriteFile(filepath.Join(outputDir, "run.bat"), []byte(content), 0644)
	}
	content := fmt.Sprintf("#!/bin/sh\nDIR=\"$(CDPATH= cd -- \"$(dirname -- \"$0\")\" && pwd)\"\n\"$DIR/%s\" \"$@\"\n", exeName)
	return os.WriteFile(filepath.Join(outputDir, "run.sh"), []byte(content), 0755)
}

func writeBundleReadme(outputDir, appName, exeName string) error {
	launcher := "run.bat"
	if runtime.GOOS != "windows" {
		launcher = "run.sh"
	}
	lines := []string{
		"# " + appName,
		"",
		"Thanks for trying this app. It was built with **Koda** - you do not need to install Go, Python, or a C++ compiler.",
		"",
		"## Run",
		"",
		"- **Windows:** double-click `" + launcher + "` or run `" + exeName + "`.",
		"- **macOS / Linux:** in a terminal: `./" + launcher + "` or `./" + exeName + "`",
		"",
		"Keep any DLLs, `.dylib`s, or `assets` folder **in the same folder** as the executable when you move or zip this directory.",
		"",
		"## Files",
		"",
		"| File | What it is |",
		"|------|------------|",
		"| `" + exeName + "` | Your application. |",
		"| `" + launcher + "` | Optional double-click launcher. |",
		"| `README.md` | This file. |",
		"| `bundle-info.txt` | Build notes for developers. |",
		"",
	}
	return os.WriteFile(filepath.Join(outputDir, "README.md"), []byte(strings.Join(lines, "\n")), 0644)
}

func writeBundleInfo(outputDir, sourcePath string) error {
	lines := []string{
		"# Koda bundle metadata",
		"source=" + sourcePath,
		"native_sources=" + os.Getenv("KODA_NATIVE_SOURCES"),
		"linkflags=" + os.Getenv("KODA_LINKFLAGS"),
		"KODA_version=" + version,
	}
	return os.WriteFile(filepath.Join(outputDir, "bundle-info.txt"), []byte(strings.Join(lines, "\n")+"\n"), 0644)
}

func copyBundleExtraFiles(outputDir string) error {
	items := parseBundleFilesEnv(os.Getenv("KODA_BUNDLE_FILES"))
	for _, item := range items {
		info, err := os.Stat(item)
		if err != nil {
			return fmt.Errorf("bundle extra %s: %w", item, err)
		}
		dst := filepath.Join(outputDir, filepath.Base(item))
		if info.IsDir() {
			if err := copyDir(item, dst); err != nil {
				return fmt.Errorf("copy bundle dir %s: %w", item, err)
			}
		} else {
			if err := copyFile(item, dst); err != nil {
				return fmt.Errorf("copy bundle file %s: %w", item, err)
			}
		}
		fmt.Printf("  copied: %s\n", dst)
	}
	return nil
}

func parseBundleFilesEnv(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := splitRespectingQuotes(raw)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		for _, q := range strings.Split(p, string(os.PathListSeparator)) {
			q = strings.TrimSpace(q)
			if q != "" {
				out = append(out, q)
			}
		}
	}
	return out
}

func splitRespectingQuotes(s string) []string {
	var out []string
	var b strings.Builder
	var quote rune
	flush := func() {
		if b.Len() == 0 {
			return
		}
		out = append(out, b.String())
		b.Reset()
	}
	for _, r := range s {
		switch {
		case quote != 0:
			if r == quote {
				quote = 0
			} else {
				b.WriteRune(r)
			}
		case r == '\'' || r == '"':
			quote = r
		case r == ' ' || r == '\t' || r == '\n' || r == '\r':
			flush()
		default:
			b.WriteRune(r)
		}
	}
	flush()
	return out
}

func copyDir(srcDir, dstDir string) error {
	return filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dstDir, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func writeBundleAssetsManifest(outputDir string, ctx *project.Context) error {
	paths := ctx.BundleExtraPaths()
	if len(paths) == 0 {
		return nil
	}
	var lines []string
	lines = append(lines, "# Koda bundled assets — use assetPath(\"file.png\") at runtime")
	for _, p := range paths {
		base := filepath.Base(p)
		if fi, err := os.Stat(p); err == nil && fi.IsDir() {
			_ = filepath.WalkDir(p, func(path string, d os.DirEntry, err error) error {
				if err != nil || d.IsDir() {
					return nil
				}
				rel, _ := filepath.Rel(p, path)
				lines = append(lines, filepath.ToSlash(filepath.Join(base, rel)))
				return nil
			})
		} else {
			lines = append(lines, base)
		}
	}
	manifest := filepath.Join(outputDir, "koda-assets.txt")
	if err := os.WriteFile(manifest, []byte(strings.Join(lines, "\n")+"\n"), 0644); err != nil {
		return err
	}
	fmt.Printf("  assets manifest: %s\n", manifest)
	return nil
}

func parseBundleArgs(args []string) (src string, outputDir string, err error) {
	outputDir = "dist"
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-o":
			if i+1 >= len(args) {
				return "", "", fmt.Errorf("-o requires a path")
			}
			i++
			outputDir = args[i]
		default:
			if strings.HasPrefix(args[i], "-") {
				return "", "", fmt.Errorf("unknown flag: %s", args[i])
			}
			if src != "" {
				return "", "", fmt.Errorf("multiple source files")
			}
			src = args[i]
		}
	}
	return src, outputDir, nil
}
