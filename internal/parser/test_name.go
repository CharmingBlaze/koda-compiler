package parser

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// synthTestName builds a stable internal symbol for a test block.
func synthTestName(display string, seq int) string {
	if u, err := strconv.Unquote(display); err == nil {
		display = u
	}
	var slug strings.Builder
	for _, r := range display {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			slug.WriteRune(unicode.ToLower(r))
		} else if r == ' ' || r == '-' || r == '_' {
			slug.WriteByte('_')
		}
	}
	s := strings.Trim(slug.String(), "_")
	if s == "" {
		s = "case"
	}
	return fmt.Sprintf("__koda_test_%d_%s", seq, s)
}
