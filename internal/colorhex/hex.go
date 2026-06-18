package colorhex

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Parse converts CSS-style hex digits (without '#') to a packed Raylib color (0xRRGGBBAA).
func Parse(digits string) (float64, error) {
	if digits == "" {
		return 0, fmt.Errorf("empty hex color")
	}
	for _, r := range digits {
		if !unicode.Is(unicode.ASCII_Hex_Digit, r) {
			return 0, fmt.Errorf("invalid hex color #%s: use only 0-9, A-F, or a-f", digits)
		}
	}
	d := strings.ToLower(digits)
	var expanded string
	switch len(d) {
	case 3:
		for _, c := range d {
			expanded += string(c) + string(c)
		}
		expanded += "ff"
	case 4:
		for _, c := range d[:3] {
			expanded += string(c) + string(c)
		}
		expanded += string(d[3]) + string(d[3])
	case 6:
		expanded = d + "ff"
	case 8:
		expanded = d
	default:
		return 0, fmt.Errorf(
			"invalid hex color #%s: use #RGB, #RGBA, #RRGGBB, or #RRGGBBAA (examples: #F00, #FF0000, #101018)",
			digits,
		)
	}
	u, err := strconv.ParseUint(expanded, 16, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid hex color #%s", digits)
	}
	r := (u >> 24) & 0xff
	g := (u >> 16) & 0xff
	b := (u >> 8) & 0xff
	a := u & 0xff
	packed := (r << 24) | (g << 16) | (b << 8) | a
	return float64(packed), nil
}
