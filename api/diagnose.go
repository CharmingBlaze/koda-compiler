package api

import (
	"errors"
	"path/filepath"
	"regexp"
	"strconv"

	"koda/internal/diagnostic"
	"koda/internal/parser"
	"koda/internal/sema"
)

var reLineCol = regexp.MustCompile(`\[line (\d+):(\d+)\]`)
var reLine = regexp.MustCompile(`\[line (\d+)\]`)

func diagnosticsFromError(path string, err error) []Diagnostic {
	if err == nil {
		return nil
	}
	abs, _ := filepath.Abs(path)
	msg := err.Error()
	d := Diagnostic{Path: abs, Severity: "error", Message: msg, Line: 1, Col: 1}
	if m := reLineCol.FindStringSubmatch(msg); len(m) == 3 {
		d.Line, _ = strconv.Atoi(m[1])
		d.Col, _ = strconv.Atoi(m[2])
		return []Diagnostic{d}
	}
	if m := reLine.FindStringSubmatch(msg); len(m) == 2 {
		d.Line, _ = strconv.Atoi(m[1])
		d.Col = 1
		return []Diagnostic{d}
	}
	return []Diagnostic{d}
}

func diagnosticFromSemaError(path string, err error) []Diagnostic {
	var me *diagnostic.MultiError
	if errors.As(err, &me) && me != nil && len(me.List) > 0 {
		out := make([]Diagnostic, 0, len(me.List))
		for _, e := range me.List {
			out = append(out, diagnosticFromSemaError(path, e)...)
		}
		return out
	}
	var d *diagnostic.DiagnosticError
	if errors.As(err, &d) && d != nil {
		p := path
		if d.File != "" {
			if abs, err := filepath.Abs(d.File); err == nil {
				p = abs
			} else {
				p = d.File
			}
		} else {
			if abs, err := filepath.Abs(path); err == nil {
				p = abs
			}
		}
		line, col := d.Line, d.Col
		if line <= 0 {
			line = 1
		}
		if col <= 0 {
			col = 1
		}
		return []Diagnostic{{Path: p, Line: line, Col: col, Message: d.Message, Hint: d.Hint, Severity: "error"}}
	}
	return diagnosticsFromError(path, err)
}

// Diagnose runs load + semantic preparation for path; optional overlay replaces on-disk source for that file only.
func Diagnose(path, overlay string) []Diagnostic {
	abs, err := filepath.Abs(path)
	if err != nil {
		return []Diagnostic{{Path: path, Line: 1, Col: 1, Message: err.Error(), Severity: "error"}}
	}
	overlays := map[string]string{}
	if overlay != "" {
		overlays[abs] = overlay
	}
	bundle, err := parser.LoadProgramWithOverlays(path, overlays)
	if err != nil {
		return diagnosticsFromError(path, err)
	}
	_, err = sema.PrepareNativeBundle(bundle)
	if err != nil {
		return diagnosticFromSemaError(path, err)
	}
	return []Diagnostic{}
}
