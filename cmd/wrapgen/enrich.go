package main

import (
	"os"
	"regexp"
	"strings"
)

func stripCComments(raw string) string {
	reMulti := regexp.MustCompile(`(?s)/\*.*?\*/`)
	out := reMulti.ReplaceAllString(raw, "")
	reSingle := regexp.MustCompile(`//.*`)
	return reSingle.ReplaceAllString(out, "")
}

func (wg *WrapperGenerator) enrichAPIFromHeaders(api *API) {
	if api == nil {
		return
	}
	enumByName := make(map[string]*Enum)
	for i := range api.Enums {
		enumByName[api.Enums[i].Name] = &api.Enums[i]
	}
	structByName := make(map[string]*Struct)
	for i := range api.Structs {
		structByName[api.Structs[i].Name] = &api.Structs[i]
	}
	macroSeen := make(map[string]struct{})
	for _, m := range api.Macros {
		macroSeen[m.Name] = struct{}{}
	}
	funcByName := make(map[string]*Function)
	for i := range api.Functions {
		funcByName[api.Functions[i].Name] = &api.Functions[i]
	}

	for _, header := range api.Headers {
		rawBytes, err := os.ReadFile(header)
		if err != nil {
			continue
		}
		raw := string(rawBytes)
		stripped := stripCComments(raw)
		parsed, err := wg.parseHeaderContent(stripped, raw, header)
		if err != nil {
			continue
		}

		for _, e := range parsed.Enums {
			if len(e.Values) == 0 {
				continue
			}
			if existing, ok := enumByName[e.Name]; ok {
				if len(existing.Values) < len(e.Values) {
					existing.Values = e.Values
				}
				continue
			}
			api.Enums = append(api.Enums, e)
			enumByName[e.Name] = &api.Enums[len(api.Enums)-1]
		}

		for _, s := range parsed.Structs {
			if len(s.Fields) == 0 {
				continue
			}
			if existing, ok := structByName[s.Name]; ok {
				if len(existing.Fields) < len(s.Fields) {
					existing.Fields = s.Fields
				}
				continue
			}
			api.Structs = append(api.Structs, s)
			structByName[s.Name] = &api.Structs[len(api.Structs)-1]
		}

		for _, m := range parsed.Macros {
			if _, ok := macroSeen[m.Name]; ok {
				continue
			}
			api.Macros = append(api.Macros, m)
			macroSeen[m.Name] = struct{}{}
		}

		for _, fn := range parsed.Functions {
			if existing, ok := funcByName[fn.Name]; ok {
				if strings.TrimSpace(existing.Documentation) == "" && strings.TrimSpace(fn.Documentation) != "" {
					existing.Documentation = fn.Documentation
				}
			}
		}
	}
}

func (wg *WrapperGenerator) parseHeaderContent(stripped, raw, headerPath string) (*API, error) {
	api := &API{}
	api.Functions = wg.extractFunctions(stripped, headerPath)
	for i := range api.Functions {
		if doc := strings.TrimSpace(docCommentBeforeFunction(raw, api.Functions[i].Name)); doc != "" {
			api.Functions[i].Documentation = doc
		}
	}
	api.Structs = wg.extractStructs(stripped, headerPath)
	api.Enums = wg.extractEnums(stripped, headerPath)
	api.Macros = wg.extractMacros(stripped, headerPath)
	api.Typedefs = wg.extractTypedefs(stripped, headerPath)
	api.Constants = wg.extractConstants(stripped, headerPath)
	return api, nil
}
