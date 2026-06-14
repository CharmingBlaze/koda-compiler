package diagnostic

import (
	"errors"
	"testing"
)

func TestErrorsAsFindsFirstDiagnosticInsideMultiError(t *testing.T) {
	e1 := &DiagnosticError{File: "a.koda", Line: 1, Col: 1, Message: "m1"}
	e2 := &DiagnosticError{File: "a.koda", Line: 2, Col: 1, Message: "m2"}
	m := &MultiError{List: []error{e1, e2}}
	var got *DiagnosticError
	if !errors.As(m, &got) {
		t.Fatal("errors.As expected to find first diagnostic")
	}
	if got.Message != "m1" {
		t.Fatalf("got %q", got.Message)
	}
}

func TestMultiErrorErrorString(t *testing.T) {
	m := &MultiError{
		Label: "game.koda",
		List: []error{
			&DiagnosticError{File: "game.koda", Line: 2, Col: 3, Message: "bad a"},
			&DiagnosticError{File: "game.koda", Line: 5, Col: 1, Message: "bad b"},
		},
	}
	s := m.Error()
	if s == "" {
		t.Fatal("empty Error()")
	}
}
