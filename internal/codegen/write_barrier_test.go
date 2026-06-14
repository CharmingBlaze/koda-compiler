package codegen

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

var gcFieldAssignPattern = regexp.MustCompile(`->\s*(?:values|elements|keys|upvalues)\s*\[[^\]]+\]\s*=`)

var writeBarrierScanExemptFuncs = map[string]string{
	// Fresh ObjTable storage zeroed before publishing; stores are initializing new buffers (no cross-gen refs yet).
	"table_rehash_open": "fresh table initialization",
}

func stripCComments(src []byte) []byte {
	out := make([]byte, 0, len(src))
	for i := 0; i < len(src); i++ {
		if i+1 < len(src) && src[i] == '/' && src[i+1] == '/' {
			i++
			for i < len(src) && src[i] != '\n' {
				i++
			}
			if i < len(src) && src[i] == '\n' {
				out = append(out, '\n')
			}
			continue
		}
		if i+1 < len(src) && src[i] == '/' && src[i+1] == '*' {
			i += 2
			for i+1 < len(src) && !(src[i] == '*' && src[i+1] == '/') {
				i++
			}
			i += 2
			continue
		}
		out = append(out, src[i])
	}
	return out
}

func matchingParenOpen(b []byte, rparen int) int {
	if rparen < 0 || rparen >= len(b) || b[rparen] != ')' {
		return -1
	}
	depth := 1
	for i := rparen - 1; i >= 0; i-- {
		switch b[i] {
		case ')':
			depth++
		case '(':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func matchingBraceClose(b []byte, lbrace int) int {
	if lbrace < 0 || lbrace >= len(b) || b[lbrace] != '{' {
		return -1
	}
	depth := 1
	for i := lbrace + 1; i < len(b); i++ {
		switch b[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func nameBeforeParen(b []byte, openParen int) string {
	i := openParen - 1
	for i >= 0 && (b[i] == ' ' || b[i] == '\t' || b[i] == '\n' || b[i] == '\r' || b[i] == '*') {
		i--
	}
	end := i + 1
	for i >= 0 && (isIdentByte(b[i])) {
		i--
	}
	return string(b[i+1 : end])
}

func isIdentByte(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_'
}

var cKeywords = map[string]bool{
	"if": true, "for": true, "while": true, "switch": true, "catch": true,
	"sizeof": true, "return": true,
}

// TestWriteBarrierCoverage ensures Value stores into GC object arrays in koda_runtime.c
// are paired with gc_write_barrier in the same function (see GC mutator invariants).
func TestWriteBarrierCoverage(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	path := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "..", "runtime", "src", "koda_runtime.c"))
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	src := stripCComments(raw)

	reBody := regexp.MustCompile(`\)\s*\{`)
	indices := reBody.FindAllIndex(src, -1)
	for _, loc := range indices {
		rparen := loc[0]
		lbrace := bytes.LastIndexByte(src[loc[0]:loc[1]], '{')
		if lbrace < 0 {
			continue
		}
		lbrace += loc[0]
		openParen := matchingParenOpen(src, rparen)
		if openParen < 0 {
			continue
		}
		name := nameBeforeParen(src, openParen)
		if name == "" || cKeywords[strings.ToLower(name)] {
			continue
		}
		bodyClose := matchingBraceClose(src, lbrace)
		if bodyClose < 0 {
			continue
		}
		body := src[lbrace : bodyClose+1]
		if !gcFieldAssignPattern.Match(body) {
			continue
		}
		if _, exempt := writeBarrierScanExemptFuncs[name]; exempt {
			continue
		}
		if !bytes.Contains(body, []byte("gc_write_barrier")) {
			t.Errorf("function %q in koda_runtime.c assigns through GC object fields but has no gc_write_barrier in the same function body", name)
		}
	}
}
