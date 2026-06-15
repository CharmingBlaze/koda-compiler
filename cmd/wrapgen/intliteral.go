package main

import (
	"regexp"
	"strconv"
	"strings"
)

var intLiteralShiftRe = regexp.MustCompile(`(?i)^\(?\s*(0x[0-9a-f]+|\d+)\s*<<\s*(\d+)\s*\)?$`)

// evaluateIntLiteral parses C integer constant expressions used in enums and macros.
func evaluateIntLiteral(expr string) (int64, bool) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return 0, false
	}
	if idx := strings.Index(expr, "//"); idx >= 0 {
		expr = strings.TrimSpace(expr[:idx])
	}
	if idx := strings.Index(expr, "/*"); idx >= 0 {
		expr = strings.TrimSpace(expr[:idx])
	}
	expr = strings.TrimSpace(expr)
	expr = strings.TrimSuffix(expr, "u")
	expr = strings.TrimSuffix(expr, "U")
	expr = strings.TrimSuffix(expr, "l")
	expr = strings.TrimSuffix(expr, "L")
	expr = strings.TrimSuffix(expr, "ul")
	expr = strings.TrimSuffix(expr, "UL")
	expr = strings.TrimSuffix(expr, "ll")
	expr = strings.TrimSuffix(expr, "LL")
	expr = strings.TrimSuffix(expr, "ull")
	expr = strings.TrimSuffix(expr, "ULL")
	expr = strings.TrimSpace(expr)

	if m := intLiteralShiftRe.FindStringSubmatch(expr); len(m) == 3 {
		base, err1 := strconv.ParseInt(m[1], 0, 64)
		shift, err2 := strconv.ParseInt(m[2], 10, 64)
		if err1 == nil && err2 == nil && shift >= 0 && shift < 63 {
			return base << shift, true
		}
	}

	if v, err := strconv.ParseInt(expr, 0, 64); err == nil {
		return v, true
	}
	return 0, false
}

func formatIntLiteral(v int64) string {
	return strconv.FormatInt(v, 10)
}
