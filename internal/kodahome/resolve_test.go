package kodahome

import (
	"strings"
	"testing"
)

func TestInstallDirWarnings_TranslocationHeuristic(t *testing.T) {
	p := "/private/var/folders/xy/AppTranslocation/ABC/koda.app/Contents/MacOS"
	w := InstallDirWarnings(p)
	if len(w) != 1 || !strings.Contains(w[0], "translocation") {
		t.Fatalf("expected translocation hint, got %#v", w)
	}
	if len(InstallDirWarnings("/Applications/Koda.app/Contents/MacOS")) != 0 {
		t.Fatal("expected no warning for normal Applications path")
	}
}
