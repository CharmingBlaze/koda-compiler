package sema

import (
	"testing"

	"koda/internal/parser"
)

func TestArityBoundsFromParams(t *testing.T) {
	t.Run("noDefaultNoRest", func(t *testing.T) {
		min, max, rest := arityBoundsFromParams([]parser.Param{
			{Name: "a"},
			{Name: "b"},
		})
		if min != 2 || max != 2 || rest {
			t.Fatalf("got min=%d max=%d rest=%v", min, max, rest)
		}
	})
	t.Run("defaults", func(t *testing.T) {
		min, max, rest := arityBoundsFromParams([]parser.Param{
			{Name: "a"},
			{Name: "b", Default: &parser.LiteralExpr{}},
			{Name: "c"},
		})
		if min != 2 || max != 3 || rest {
			t.Fatalf("got min=%d max=%d rest=%v", min, max, rest)
		}
	})
	t.Run("rest", func(t *testing.T) {
		min, max, rest := arityBoundsFromParams([]parser.Param{
			{Name: "a"},
			{Name: "r", IsRest: true},
		})
		if min != 1 || max != -1 || !rest {
			t.Fatalf("got min=%d max=%d rest=%v", min, max, rest)
		}
	})
}
