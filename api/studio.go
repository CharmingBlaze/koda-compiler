package api

import (
	"os"
	"path/filepath"

	"koda/internal/kodahome"
	"koda/internal/project"
)

// SDKLine is one row in the IDE SDK health panel.
type SDKLine struct {
	OK     bool   `json:"ok"`
	Label  string `json:"label"`
	Detail string `json:"detail"`
	Fix    string `json:"fix"`
}

// SDKStatus summarizes whether Koda Studio can compile and run programs.
type SDKStatus struct {
	OK         bool      `json:"ok"`
	Version    string    `json:"version"`
	InstallDir string    `json:"installDir"`
	StdlibDir  string    `json:"stdlibDir"`
	Lines      []SDKLine `json:"lines"`
}

// BootstrapSDKRoot sets KODA_HOME when unset and stdlib/ is found near exeDir.
func BootstrapSDKRoot(exeDir string) {
	kodahome.BootstrapSDKRoot(exeDir)
}

// EnsureSDKFromExecutable locates stdlib/ relative to the running binary (for IDE hosts).
func EnsureSDKFromExecutable() {
	kodahome.EnsureSDKFromExecutable()
}

// CheckSDK returns a lightweight health report for the welcome screen.
func CheckSDK(studioVersion string) SDKStatus {
	st := SDKStatus{Version: studioVersion, Lines: []SDKLine{}}
	failures := 0

	add := func(ok bool, label, detail, fix string) {
		st.Lines = append(st.Lines, SDKLine{OK: ok, Label: label, Detail: detail, Fix: fix})
		if !ok {
			failures++
		}
	}

	install, err := kodahome.InstallDir()
	if err != nil {
		add(false, "SDK folder", err.Error(), "Place Koda Studio next to koda.exe and stdlib/ (SDK zip layout).")
	} else {
		st.InstallDir = install
		add(true, "SDK folder", install, "")
	}

	stdlib, err := kodahome.StdlibDir()
	if err != nil {
		add(false, "stdlib", err.Error(), "Unzip the full SDK so stdlib/ sits next to Koda Studio.")
	} else {
		st.StdlibDir = stdlib
		if fi, err := os.Stat(filepath.Join(stdlib, "math.koda")); err != nil || fi.IsDir() {
			add(false, "stdlib", stdlib, "Missing math.koda — download the SDK zip from GitHub Releases.")
		} else {
			add(true, "stdlib", stdlib, "")
		}
	}

	if _, err := kodahome.FindToolchain(); err != nil {
		add(false, "compiler runtime", err.Error(), "Use a release SDK or build runtime/libkoda_runtime.a from source.")
	} else {
		add(true, "compiler runtime", "ready", "")
	}

	st.OK = failures == 0
	return st
}

// ListProjectTemplates returns hello, game, graphics, … for the new-project wizard.
func ListProjectTemplates() []string {
	return project.ListTemplates()
}

// ScaffoldProject creates parentDir/name from an embedded template.
func ScaffoldProject(parentDir, name, template string) (string, error) {
	return project.Scaffold(parentDir, name, template)
}
