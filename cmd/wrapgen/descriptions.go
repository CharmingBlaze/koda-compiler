package main

import (
	"strings"
	"unicode"
)

func functionDescription(rawHeader, funcName string) string {
	if d := strings.TrimSpace(docCommentBeforeFunction(rawHeader, funcName)); d != "" {
		return d
	}
	return describeFromCamelName(funcName)
}

func describeFromCamelName(name string) string {
	parts := splitCamelCase(name)
	if len(parts) == 0 {
		return "No description available."
	}
	verb := parts[0]
	rest := strings.ToLower(strings.Join(parts[1:], " "))
	verbMap := map[string]string{
		"Init":    "Initialises",
		"Close":   "Closes",
		"Get":     "Returns",
		"Set":     "Sets",
		"Draw":    "Draws",
		"Load":    "Loads",
		"Unload":  "Unloads",
		"Play":    "Plays",
		"Stop":    "Stops",
		"Pause":   "Pauses",
		"Update":  "Updates",
		"Enable":  "Enables",
		"Disable": "Disables",
		"Check":   "Checks",
		"Is":      "Returns true if",
		"Begin":   "Begins",
		"End":     "Ends",
		"Show":    "Shows",
		"Hide":    "Hides",
	}
	vw, ok := verbMap[verb]
	if !ok {
		vw = "Calls"
		if rest == "" {
			return vw + " " + strings.ToLower(name) + "."
		}
		return vw + " " + rest + "."
	}
	if rest == "" {
		if verb == "Is" {
			return vw + " the condition holds for " + strings.ToLower(name[len(verb):]) + "."
		}
		return vw + " " + strings.ToLower(strings.Join(parts[1:], " ")) + "."
	}
	if verb == "Is" {
		return vw + " " + rest + "."
	}
	return vw + " " + rest + "."
}

func splitCamelCase(s string) []string {
	var parts []string
	var cur strings.Builder
	for _, r := range s {
		if unicode.IsUpper(r) && cur.Len() > 0 {
			parts = append(parts, cur.String())
			cur.Reset()
		}
		cur.WriteRune(r)
	}
	if cur.Len() > 0 {
		parts = append(parts, cur.String())
	}
	return parts
}

func kodaDocTypeForCType(t string) string {
	t = strings.TrimSpace(t)
	t = strings.ReplaceAll(t, "const ", "")
	switch {
	case strings.Contains(t, "void") && strings.Contains(t, "*"):
		return "object"
	case strings.Contains(t, "char") && strings.Contains(t, "*"):
		return "string"
	case strings.HasPrefix(t, "bool") || t == "_Bool":
		return "bool"
	case strings.Contains(t, "float") || strings.Contains(t, "double") || strings.Contains(t, "int") ||
		strings.Contains(t, "short") || strings.Contains(t, "long") || strings.Contains(t, "size_t"):
		return "number"
	case t == "void":
		return "nothing"
	default:
		if strings.ContainsAny(t, "*") {
			return "object"
		}
		return "object"
	}
}

func paramDocLine(paramName, cType string) (kodaType, hint string) {
	ft := kodaDocTypeForCType(cType)
	hint = "argument value"
	switch ft {
	case "string":
		hint = "null-terminated C string passed as Koda string"
	case "number":
		hint = "numeric value"
	case "bool":
		hint = "truthy Koda value"
	case "object":
		hint = "boxed object (struct layout depends on the C API)"
	}
	return ft, hint
}

func returnsClause(ret string) (hasReturns bool, line string) {
	ret = strings.TrimSpace(ret)
	ret = strings.TrimPrefix(ret, "RLAPI ")
	ret = strings.TrimSpace(ret)
	if ret == "" || ret == "void" {
		return false, ""
	}
	ft := kodaDocTypeForCType(ret)
	return true, "/// @returns " + ft + " — result from the native call"
}
