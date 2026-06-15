package project

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// GameShimSymbols are raylib_shim bindings required by stdlib/game.koda (@game).
var GameShimSymbols = []string{
	"initwindow", "closewindow", "windowshouldclose", "settargetfps",
	"begindrawing", "enddrawing", "clearbackground",
	"drawtext", "drawrectangle", "drawcircle",
	"drawline", "drawcirclelines", "drawrectanglelines",
	"loadtexture", "drawtexture", "unloadtexture",
	"iskeydown", "iskeypressed",
	"getmousex", "getmousey", "ismousebuttondown", "ismousebuttonpressed", "getmousewheelmove",
	"getscreenwidth", "getscreenheight", "setwindowtitle", "getfps",
}

// RaylibShimReport describes whether a project's raylib_shim matches the SDK copy.
type RaylibShimReport struct {
	ProjectRoot    string
	ShimDir        string
	Stale          bool
	MissingSymbols []string
	OutdatedFiles  []string
	MissingFiles   []string
}

func (r RaylibShimReport) NeedsRefresh() bool {
	return r.Stale
}

func (r RaylibShimReport) Detail() string {
	if !r.Stale {
		return "up to date"
	}
	var parts []string
	if len(r.MissingFiles) > 0 {
		parts = append(parts, "missing "+strings.Join(r.MissingFiles, ", "))
	}
	if len(r.OutdatedFiles) > 0 {
		parts = append(parts, "outdated "+strings.Join(r.OutdatedFiles, ", "))
	}
	if len(r.MissingSymbols) > 0 {
		n := r.MissingSymbols
		if len(n) > 4 {
			n = append(append([]string{}, n[:4]...), "...")
		}
		parts = append(parts, "missing @game symbols: "+strings.Join(n, ", "))
	}
	if len(parts) == 0 {
		return "out of date"
	}
	return strings.Join(parts, "; ")
}

// CheckRaylibShim compares the project's raylib_shim to the SDK canonical copy.
// Returns a non-stale report when the project has no graphics shim configured.
func CheckRaylibShim(projectRoot string) (RaylibShimReport, error) {
	report := RaylibShimReport{ProjectRoot: projectRoot}
	cfg, root, err := Find(projectRoot)
	if err != nil {
		return report, err
	}
	if cfg == nil {
		return report, nil
	}
	report.ProjectRoot = root

	shimDir, ok := projectRaylibShimDir(root, cfg)
	if !ok {
		return report, nil
	}
	report.ShimDir = shimDir

	canonical, hasCanonical := canonicalRaylibShimDir()
	for _, name := range raylibShimFiles {
		projectPath := filepath.Join(shimDir, name)
		if _, err := os.Stat(projectPath); err != nil {
			report.Stale = true
			report.MissingFiles = append(report.MissingFiles, name)
			continue
		}
		if hasCanonical {
			same, err := filesEqual(projectPath, filepath.Join(canonical, name))
			if err != nil {
				return report, err
			}
			if !same {
				report.Stale = true
				report.OutdatedFiles = append(report.OutdatedFiles, name)
			}
		}
	}

	kodaPath := filepath.Join(shimDir, "raylib.koda")
	if missing := missingShimSymbols(kodaPath); len(missing) > 0 {
		report.Stale = true
		report.MissingSymbols = missing
	}

	return report, nil
}

func projectRaylibShimDir(root string, cfg *Config) (string, bool) {
	if cfg == nil {
		return "", false
	}
	for _, src := range cfg.Native.Sources {
		slash := filepath.ToSlash(strings.TrimSpace(src))
		if !strings.Contains(slash, "raylib_shim") {
			continue
		}
		dir := filepath.Dir(filepath.Join(root, filepath.FromSlash(src)))
		if raylibShimComplete(dir) || fileExists(filepath.Join(dir, "raylib.koda")) {
			return dir, true
		}
	}
	if cfg.Native.Graphics {
		dir := filepath.Join(root, "wrappers", "raylib_shim")
		if fileExists(filepath.Join(dir, "raylib.koda")) || fileExists(filepath.Join(dir, "wrapper.c")) {
			return dir, true
		}
	}
	return "", false
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}

func missingShimSymbols(kodaPath string) []string {
	data, err := os.ReadFile(kodaPath)
	if err != nil {
		return append([]string{}, GameShimSymbols...)
	}
	text := strings.ToLower(string(data))
	var missing []string
	for _, sym := range GameShimSymbols {
		if !shimDeclaresSymbol(text, sym) {
			missing = append(missing, sym)
		}
	}
	return missing
}

func shimDeclaresSymbol(lowerFile, sym string) bool {
	sym = strings.ToLower(sym)
	return strings.Contains(lowerFile, "koda:extern "+sym+" ") ||
		strings.Contains(lowerFile, "let "+sym+" ") ||
		strings.Contains(lowerFile, "func "+sym+"(")
}

func filesEqual(a, b string) (bool, error) {
	ha, err := fileHashPrefix(a)
	if err != nil {
		return false, err
	}
	hb, err := fileHashPrefix(b)
	if err != nil {
		return false, err
	}
	return ha == hb, nil
}

func fileHashPrefix(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	sum := sha256.New()
	if _, err := io.Copy(sum, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(sum.Sum(nil))[:16], nil
}

// RefreshRaylibShimIfStale updates the project shim when CheckRaylibShim reports drift.
func RefreshRaylibShimIfStale(projectRoot string) (RaylibShimReport, bool, error) {
	report, err := CheckRaylibShim(projectRoot)
	if err != nil {
		return report, false, err
	}
	if !report.Stale {
		return report, false, nil
	}
	dest := report.ShimDir
	if dest == "" {
		dest = filepath.Join(report.ProjectRoot, "wrappers", "raylib_shim")
	}
	if err := SyncRaylibShim(dest); err != nil {
		return report, false, fmt.Errorf("refresh raylib shim: %w", err)
	}
	after, err := CheckRaylibShim(projectRoot)
	if err != nil {
		return report, true, err
	}
	return after, true, nil
}
