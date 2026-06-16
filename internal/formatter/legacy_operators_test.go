package formatter

import (
	"strings"
	"testing"
)

func TestNormalizeLegacyOperatorsSkipsStrings(t *testing.T) {
	src := `let s = "a === b"; /* !== */ let x = 1 === 1;`
	got := normalizeLegacyOperators(src)
	if !strings.Contains(got, `"a === b"`) {
		t.Fatalf("string not preserved: %q", got)
	}
	if !strings.Contains(got, `/* !== */`) {
		t.Fatalf("block comment not preserved: %q", got)
	}
	if strings.Contains(got, "1 === 1") {
		t.Fatalf("code === not rewritten: %q", got)
	}
	if !strings.Contains(got, "1 == 1") {
		t.Fatalf("expected == rewrite: %q", got)
	}
}
