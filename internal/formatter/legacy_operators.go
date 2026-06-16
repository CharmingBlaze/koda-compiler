package formatter

import "strings"

// normalizeLegacyOperators rewrites === / !== outside strings and comments so koda fmt
// can migrate sources removed from the lexer without touching string contents.
func normalizeLegacyOperators(src string) string {
	var b strings.Builder
	b.Grow(len(src))
	i := 0
	n := len(src)
	for i < n {
		switch {
		case src[i] == '/' && i+1 < n && src[i+1] == '/':
			j := i + 2
			for j < n && src[j] != '\n' {
				j++
			}
			b.WriteString(src[i:j])
			i = j
		case src[i] == '/' && i+1 < n && src[i+1] == '*':
			j := i + 2
			for j+1 < n && !(src[j] == '*' && src[j+1] == '/') {
				j++
			}
			if j+1 < n {
				j += 2
			}
			b.WriteString(src[i:j])
			i = j
		case src[i] == '"':
			j := i + 1
			for j < n {
				if src[j] == '\\' {
					j += 2
					continue
				}
				if src[j] == '"' {
					j++
					break
				}
				j++
			}
			b.WriteString(src[i:j])
			i = j
		case i+2 < n && src[i] == '!' && src[i+1] == '=' && src[i+2] == '=':
			b.WriteString("!=")
			i += 3
		case i+2 < n && src[i] == '=' && src[i+1] == '=' && src[i+2] == '=':
			b.WriteString("==")
			i += 3
		default:
			b.WriteByte(src[i])
			i++
		}
	}
	return b.String()
}
