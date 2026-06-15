package kodahome

import (
	"os"
	"path/filepath"
	"strings"
)

func sdkRootFromEnv() (string, bool) {
	h := strings.TrimSpace(os.Getenv("KODA_HOME"))
	if h == "" {
		return "", false
	}
	abs, err := filepath.Abs(h)
	if err != nil {
		return "", false
	}
	return abs, true
}

func hasStdlibRoot(root string) bool {
	if root == "" {
		return false
	}
	fi, err := os.Stat(filepath.Join(root, "stdlib", "math.koda"))
	return err == nil && !fi.IsDir()
}

// SDKRootCandidates returns directories that might contain stdlib/, walking up from exeDir.
func SDKRootCandidates(exeDir string) []string {
	dir, err := filepath.Abs(exeDir)
	if err != nil {
		return nil
	}
	var out []string
	for i := 0; i < 12; i++ {
		out = append(out, dir)
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return out
}

// EnsureSDKFromExecutable locates stdlib/ relative to the running binary.
func EnsureSDKFromExecutable() {
	dir, err := executableDir()
	if err != nil {
		return
	}
	BootstrapSDKRoot(dir)
}

// BootstrapSDKRoot sets KODA_HOME when unset and a stdlib/ tree is found near exeDir.
func BootstrapSDKRoot(exeDir string) (string, bool) {
	if root, ok := sdkRootFromEnv(); ok && hasStdlibRoot(root) {
		return root, true
	}
	for _, c := range SDKRootCandidates(exeDir) {
		if hasStdlibRoot(c) {
			_ = os.Setenv("KODA_HOME", c)
			return c, true
		}
	}
	return "", false
}
