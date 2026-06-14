package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiagnoseValidProgram(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "ok.koda")
	src := "let x = 1;\nprint(x);\n"
	if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	d := Diagnose(p, "")
	if len(d) != 0 {
		t.Fatalf("want no diagnostics, got %#v", d)
	}
}

func TestDiagnoseTypoSuggestionHint(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "typo.koda")
	src := "let drawRectangle = 1;\ndrawRectange;\n"
	if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	d := Diagnose(p, "")
	if len(d) == 0 {
		t.Fatal("want diagnostics")
	}
	if !strings.Contains(strings.ToLower(d[0].Message), "undefined") {
		t.Fatalf("message: %q", d[0].Message)
	}
	if !strings.Contains(strings.ToLower(d[0].Hint), "did you mean") {
		t.Fatalf("want hint with 'did you mean', got hint=%q msg=%q", d[0].Hint, d[0].Message)
	}
}

func TestDiagnoseSemaUndefined(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "bad.koda")
	if err := os.WriteFile(p, []byte("notARealBinding;\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	d := Diagnose(p, "")
	if len(d) == 0 {
		t.Fatal("want diagnostics")
	}
	msg := strings.ToLower(d[0].Message)
	if !strings.Contains(msg, "undefined") {
		t.Fatalf("message: %q", d[0].Message)
	}
}

func TestDiagnoseOverlayReplacesOnDiskSource(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "entry.koda")
	if err := os.WriteFile(p, []byte("brokenOnDisk;\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	overlay := "let ok = 1;\nprint(ok);\n"
	d := Diagnose(p, overlay)
	if len(d) != 0 {
		t.Fatalf("overlay should supply valid sema input, got %#v", d)
	}
}

func TestDiagnoseCallArity(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "arity.koda")
	src := "func f(a) { return a; }\nf(1, 2);\n"
	if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	d := Diagnose(p, "")
	if len(d) == 0 {
		t.Fatal("want arity diagnostic")
	}
	if !strings.Contains(strings.ToLower(d[0].Message), "too many") {
		t.Fatalf("message: %q", d[0].Message)
	}
}

func TestDiagnoseArgvMethodArity(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "mArity.koda")
	src := "func main() {\n\t\"x\".split();\n}\n"
	if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	d := Diagnose(p, "")
	if len(d) == 0 {
		t.Fatal("want argv method arity diagnostic")
	}
	msg := strings.ToLower(d[0].Message)
	if !strings.Contains(msg, "split") || !strings.Contains(msg, "wrong number") {
		t.Fatalf("message: %q", d[0].Message)
	}
}

func TestDiagnoseMultipleSemaErrorsAggregated(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "multi.koda")
	src := "aaa;\nbbb;\n"
	if err := os.WriteFile(p, []byte(src), 0o644); err != nil {
		t.Fatal(err)
	}
	d := Diagnose(p, "")
	if len(d) < 2 {
		t.Fatalf("want two diagnostics, got %d: %#v", len(d), d)
	}
	m0 := strings.ToLower(d[0].Message)
	m1 := strings.ToLower(d[1].Message)
	if !strings.Contains(m0, "aaa") && !strings.Contains(m1, "aaa") {
		t.Fatalf("want aaa in diagnostics: %#v", d)
	}
	if !strings.Contains(m0, "bbb") && !strings.Contains(m1, "bbb") {
		t.Fatalf("want bbb in diagnostics: %#v", d)
	}
}
