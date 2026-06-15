package main

import (
	"fmt"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// API represents the extracted API from C/C++ headers
type API struct {
	Name         string
	Functions    []Function
	Structs      []Struct
	Enums        []Enum
	Macros       []Macro
	Typedefs     []Typedef
	Constants    []Constant
	Headers      []string
	Dependencies []string
}

type Function struct {
	Name          string
	ReturnType    string
	Parameters    []Parameter
	Variadic      bool
	Static        bool
	Inline        bool
	Documentation string
	Header        string
}

type Parameter struct {
	Name    string
	Type    string
	Default string
}

type Struct struct {
	Name          string
	Fields        []Field
	Methods       []Function
	Documentation string
	Header        string
}

type Field struct {
	Name string
	Type string
}

type Enum struct {
	Name          string
	Values        []EnumValue
	Documentation string
	Header        string
}

type EnumValue struct {
	Name  string
	Value int64
}

type Macro struct {
	Name          string
	Value         string
	Parameters    []string
	Documentation string
	Header        string
}

type Typedef struct {
	Name          string
	TargetType    string
	Documentation string
	Header        string
}

type Constant struct {
	Name          string
	Type          string
	Value         string
	Documentation string
	Header        string
}

// wrapgenBlockedBindingNames are not all Go keywords but cannot be used as Koda `let` names
// (regex false-positives from C headers / macros, or C reserved words).
var wrapgenBlockedBindingNames = map[string]struct{}{
	"let": {}, "func": {}, "import": {}, "export": {},
	"if": {}, "else": {}, "while": {}, "for": {}, "switch": {}, "case": {}, "default": {},
	"break": {}, "continue": {}, "return": {}, "goto": {}, "do": {},
	"sizeof": {}, "typeof": {}, "asm": {}, "typeof_unqual": {},
	"struct": {}, "union": {}, "enum": {}, "typedef": {},
	"signed": {}, "unsigned": {}, "const": {}, "volatile": {},
	"static": {}, "extern": {}, "inline": {}, "restrict": {},
	"auto": {}, "register": {}, "alignas": {}, "alignof": {},
	"bool": {}, "true": {}, "false": {},
	"void": {}, "int": {}, "char": {}, "short": {}, "long": {}, "float": {}, "double": {},
}

// isValidWrapgenBindingName rejects regex false-positives (e.g. macro fragments parsed as `if(...)`) and names that cannot be emitted as Koda identifiers.
func isValidWrapgenBindingName(name string) bool {
	if name == "" {
		return false
	}
	r, w := utf8.DecodeRuneInString(name)
	if w == 0 || r == utf8.RuneError {
		return false
	}
	if r != '_' && !unicode.IsLetter(r) {
		return false
	}
	for _, ch := range name[w:] {
		if ch != '_' && !unicode.IsLetter(ch) && !unicode.IsDigit(ch) {
			return false
		}
	}
	if token.IsKeyword(name) {
		return false
	}
	if _, bad := wrapgenBlockedBindingNames[name]; bad {
		return false
	}
	return true
}

func filterAndDedupeFunctions(fs []Function) []Function {
	seen := make(map[string]struct{})
	var out []Function
	for _, f := range fs {
		if !isValidWrapgenBindingName(f.Name) {
			continue
		}
		if _, ok := seen[f.Name]; ok {
			continue
		}
		seen[f.Name] = struct{}{}
		out = append(out, f)
	}
	return out
}

type WrapperGenerator struct {
	config *WrapGenConfig
	fset   *token.FileSet
}

func NewWrapperGenerator(config *WrapGenConfig) *WrapperGenerator {
	return &WrapperGenerator{
		config: config,
		fset:   token.NewFileSet(),
	}
}

// ParseHeaders extracts API information from C/C++ header files
func (wg *WrapperGenerator) ParseHeaders() (*API, error) {
	if wg.config.UseClang {
		clang := NewClangParser(wg.config)
		api, err := clang.ParseWithClang(wg.config.InputHeaders)
		if err == nil && len(api.Functions) > 0 {
			wg.enrichAPIFromHeaders(api)
			api.Functions = filterAndDedupeFunctions(api.Functions)
			wg.analyzeDependencies(api)
			if wg.config.Verbose {
				fmt.Printf("Parsed with clang: %d functions\n", len(api.Functions))
			}
			return api, nil
		}
		if wg.config.Verbose && err != nil {
			fmt.Printf("Clang parse unavailable, using regex: %v\n", err)
		}
	}

	api := &API{
		Name: wg.config.LibraryName,
	}

	if wg.config.Verbose {
		fmt.Printf("Using regex parser for C/C++ headers...\n")
	}

	// Parse all headers with enhanced regex
	for _, header := range wg.config.InputHeaders {
		if wg.config.Verbose {
			fmt.Printf("Parsing header: %s\n", header)
		}

		headerAPI, err := wg.parseHeader(header)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %s: %v", header, err)
		}

		// Merge header API into main API
		api.Functions = append(api.Functions, headerAPI.Functions...)
		api.Structs = append(api.Structs, headerAPI.Structs...)
		api.Enums = append(api.Enums, headerAPI.Enums...)
		api.Macros = append(api.Macros, headerAPI.Macros...)
		api.Typedefs = append(api.Typedefs, headerAPI.Typedefs...)
		api.Constants = append(api.Constants, headerAPI.Constants...)
		api.Headers = append(api.Headers, header)
	}

	api.Functions = filterAndDedupeFunctions(api.Functions)

	// Analyze dependencies and relationships
	wg.analyzeDependencies(api)

	if wg.config.Verbose {
		fmt.Printf("Found %d functions, %d structs, %d enums, %d macros, %d typedefs, %d constants\n",
			len(api.Functions), len(api.Structs), len(api.Enums), len(api.Macros), len(api.Typedefs), len(api.Constants))
	}

	return api, nil
}

// parseHeader parses a single C/C++ header file
func (wg *WrapperGenerator) parseHeader(headerPath string) (*API, error) {
	content, err := ioutil.ReadFile(headerPath)
	if err != nil {
		return nil, err
	}
	raw := string(content)
	return wg.parseHeaderContent(stripCComments(raw), raw, headerPath)
}

// extractFunctions extracts function declarations from header content
func (wg *WrapperGenerator) extractFunctions(content, headerPath string) []Function {
	var functions []Function

	// Function regex pattern
	funcPattern := regexp.MustCompile(`(?m)^\s*(?:inline\s+)?(?:static\s+)?(?:\w+\s+)*?(\w+)\s*\(([^)]*)\)\s*(?:__attribute__\s*\([^)]*\))?\s*;?`)

	matches := funcPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			funcName := strings.TrimSpace(match[1])
			paramsStr := strings.TrimSpace(match[2])

			if !isValidWrapgenBindingName(funcName) {
				continue
			}

			// Skip if this looks like a struct/union definition
			if wg.isStructDefinition(funcName, paramsStr) {
				continue
			}

			function := Function{
				Name:   funcName,
				Header: headerPath,
			}

			// Extract return type (simplified)
			fullMatch := match[0]
			returnType := wg.extractReturnType(fullMatch, funcName)
			function.ReturnType = returnType

			// Parse parameters
			function.Parameters = wg.parseParameters(paramsStr)

			functions = append(functions, function)
		}
	}

	return functions
}

// extractStructs extracts struct definitions from header content
func (wg *WrapperGenerator) extractStructs(content, headerPath string) []Struct {
	var structs []Struct

	// Pattern 1: typedef struct [Name] { ... } [TypedefName];
	structPattern1 := regexp.MustCompile(`(?s)typedef\s+struct\s*(\w*)\s*\{([^}]+)\}\s*(\w+)[^;]*;`)
	// Pattern 2: struct Name { ... };
	structPattern2 := regexp.MustCompile(`(?s)struct\s+(\w+)\s*\{([^}]+)\}\s*;`)

	patterns := []struct {
		re   *regexp.Regexp
		nIdx int
		fIdx int
	}{
		{structPattern1, 3, 2},
		{structPattern2, 1, 2},
	}

	seen := make(map[string]bool)

	for _, p := range patterns {
		matches := p.re.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) > p.nIdx && len(match) > p.fIdx {
				structName := strings.TrimSpace(match[p.nIdx])
				if structName == "" && p.nIdx == 3 {
					structName = strings.TrimSpace(match[1]) // fallback to tag name
				}

				if structName == "" || seen[structName] {
					continue
				}
				seen[structName] = true

				structDef := Struct{
					Name:   structName,
					Header: headerPath,
				}

				// Parse fields
				fieldsContent := match[p.fIdx]
				structDef.Fields = wg.parseStructFields(fieldsContent)

				structs = append(structs, structDef)
			}
		}
	}

	return structs
}

// extractEnums extracts enum definitions from header content
func (wg *WrapperGenerator) extractEnums(content, headerPath string) []Enum {
	var enums []Enum

	// Enum regex patterns
	enumPattern1 := regexp.MustCompile(`(?s)typedef\s+enum\s+(\w*)\s*\{([^}]+)\}\s*(\w+)[^;]*;`)
	enumPattern2 := regexp.MustCompile(`(?s)enum\s+(\w+)\s*\{([^}]+)\}\s*;`)

	patterns := []*regexp.Regexp{enumPattern1, enumPattern2}

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				enumName := strings.TrimSpace(match[3])
				if enumName == "" {
					enumName = strings.TrimSpace(match[1])
				}

				enumDef := Enum{
					Name:   enumName,
					Header: headerPath,
				}

				// Parse enum values
				valuesContent := match[2]
				enumDef.Values = wg.parseEnumValues(valuesContent)

				enums = append(enums, enumDef)
			}
		}
	}

	return enums
}

// extractMacros extracts macro definitions from header content
func (wg *WrapperGenerator) extractMacros(content, headerPath string) []Macro {
	var macros []Macro

	// Macro regex pattern
	macroPattern := regexp.MustCompile(`(?m)^#define\s+(\w+)\s*(?:\(([^)]*)\))?\s*(.+)$`)

	matches := macroPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			macroName := strings.TrimSpace(match[1])
			paramsStr := ""
			value := ""

			if len(match) >= 4 && match[2] != "" {
				// Function-like macro
				paramsStr = strings.TrimSpace(match[2])
				value = strings.TrimSpace(match[3])
			} else if len(match) >= 3 {
				// Object-like macro
				value = strings.TrimSpace(match[2])
			}

			macro := Macro{
				Name:   macroName,
				Value:  value,
				Header: headerPath,
			}

			if paramsStr != "" {
				macro.Parameters = strings.Split(paramsStr, ",")
				for i, param := range macro.Parameters {
					macro.Parameters[i] = strings.TrimSpace(param)
				}
			}

			macros = append(macros, macro)
		}
	}

	return macros
}

// extractTypedefs extracts typedef declarations from header content
func (wg *WrapperGenerator) extractTypedefs(content, headerPath string) []Typedef {
	var typedefs []Typedef

	// Typedef regex pattern (excluding function pointers and complex blocks for now)
	// Supports: typedef TargetType NewName;  and  typedef Type *NewName;
	typedefPattern := regexp.MustCompile(`typedef\s+([a-zA-Z_][a-zA-Z0-9_\s\*]+?)\s+([a-zA-Z_][a-zA-Z0-9_]+)\s*;`)
	typedefPtrPattern := regexp.MustCompile(`typedef\s+([a-zA-Z_][a-zA-Z0-9_\s]+)\s*\*\s*([a-zA-Z_][a-zA-Z0-9_]+)\s*;`)

	patterns := []*regexp.Regexp{typedefPtrPattern, typedefPattern}
	seen := make(map[string]struct{})

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				targetType := strings.TrimSpace(match[1])
				typeName := strings.TrimSpace(match[2])
				if pattern == typedefPtrPattern {
					targetType = strings.TrimSpace(targetType + " *")
				}
				if _, ok := seen[typeName]; ok {
					continue
				}
				seen[typeName] = struct{}{}

				typedefs = append(typedefs, Typedef{
					Name:       typeName,
					TargetType: targetType,
					Header:     headerPath,
				})
			}
		}
	}

	return typedefs
}

// extractConstants extracts constant definitions from header content
func (wg *WrapperGenerator) extractConstants(content, headerPath string) []Constant {
	var constants []Constant

	// Constant regex patterns
	constPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?m)^extern\s+const\s+(\w+)\s+(\w+)[^;]*;`),
		regexp.MustCompile(`(?m)^const\s+(\w+)\s+(\w+)\s*=\s*([^;]+);`),
	}

	for _, pattern := range constPatterns {
		matches := pattern.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) >= 3 {
				constType := strings.TrimSpace(match[1])
				constName := strings.TrimSpace(match[2])
				constValue := ""

				if len(match) >= 4 {
					constValue = strings.TrimSpace(match[3])
				}

				constant := Constant{
					Name:   constName,
					Type:   constType,
					Value:  constValue,
					Header: headerPath,
				}

				constants = append(constants, constant)
			}
		}
	}

	return constants
}

// Helper methods for parsing specific components
func (wg *WrapperGenerator) extractReturnType(fullMatch, funcName string) string {
	// Extract return type by removing function name and parameters
	beforeFunc := strings.Split(fullMatch, funcName)[0]
	returnType := strings.TrimSpace(beforeFunc)

	// Clean up return type
	returnType = strings.ReplaceAll(returnType, "inline", "")
	returnType = strings.ReplaceAll(returnType, "static", "")
	returnType = strings.ReplaceAll(returnType, "RLAPI", "")
	returnType = strings.TrimSpace(returnType)

	return returnType
}

func (wg *WrapperGenerator) parseParameters(paramsStr string) []Parameter {
	var parameters []Parameter

	if paramsStr == "" || paramsStr == "void" {
		return parameters
	}

	paramPairs := strings.Split(paramsStr, ",")
	for _, param := range paramPairs {
		param = strings.TrimSpace(param)
		if param == "" || param == "void" {
			continue
		}

		parts := strings.Fields(param)
		if len(parts) >= 2 {
			paramName := parts[len(parts)-1]
			paramTypeParts := parts[:len(parts)-1]
			for strings.HasPrefix(paramName, "*") {
				paramName = paramName[1:]
				if len(paramTypeParts) > 0 {
					paramTypeParts[len(paramTypeParts)-1] += "*"
				}
			}
			paramType := strings.Join(paramTypeParts, " ")

			// Handle default values
			if strings.Contains(paramName, "=") {
				nameValue := strings.Split(paramName, "=")
				paramName = strings.TrimSpace(nameValue[0])
				// defaultValue = strings.TrimSpace(nameValue[1])
			}

			parameters = append(parameters, Parameter{
				Name: paramName,
				Type: paramType,
			})
		}
	}

	return parameters
}

func (wg *WrapperGenerator) parseStructFields(fieldsContent string) []Field {
	var fields []Field

	// Strip comments from the entire block first
	reMulti := regexp.MustCompile(`(?s)/\*.*?\*/`)
	fieldsContent = reMulti.ReplaceAllString(fieldsContent, "")
	reSingle := regexp.MustCompile(`//.*`)
	fieldsContent = reSingle.ReplaceAllString(fieldsContent, "")

	lines := strings.Split(fieldsContent, ";")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			fieldName := parts[len(parts)-1]
			fieldTypeParts := parts[:len(parts)-1]
			for strings.HasPrefix(fieldName, "*") {
				fieldName = fieldName[1:]
				if len(fieldTypeParts) > 0 {
					fieldTypeParts[len(fieldTypeParts)-1] += "*"
				}
			}
			fieldType := strings.Join(fieldTypeParts, " ")
			if idx := strings.Index(fieldName, "["); idx >= 0 {
				fieldType += fieldName[idx:]
				fieldName = fieldName[:idx]
			}

			fields = append(fields, Field{
				Name: fieldName,
				Type: fieldType,
			})
		}
	}

	return fields
}

func (wg *WrapperGenerator) isValidIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !unicode.IsLetter(r) && r != '_' {
				return false
			}
		} else {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
				return false
			}
		}
	}
	return true
}

func (wg *WrapperGenerator) parseEnumValues(valuesContent string) []EnumValue {
	var values []EnumValue

	// Replace newlines with spaces to handle multi-line enums
	valuesContent = strings.ReplaceAll(valuesContent, "\n", " ")
	valuesContent = strings.ReplaceAll(valuesContent, "\r", " ")

	lines := strings.Split(valuesContent, ",")
	currentValue := int64(0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Remove comments
		if idx := strings.Index(line, "//"); idx != -1 {
			line = line[:idx]
		}
		if idx := strings.Index(line, "/*"); idx != -1 {
			if endIdx := strings.Index(line, "*/"); endIdx != -1 {
				line = line[:idx] + line[endIdx+2:]
			} else {
				line = line[:idx]
			}
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "=")
		valueName := strings.TrimSpace(parts[0])
		if valueName == "" {
			continue
		}

		// Final cleanup of valueName (remove any remaining non-identifier chars)
		partsName := strings.Fields(valueName)
		if len(partsName) > 0 {
			valueName = partsName[len(partsName)-1]
		}

		if !wg.isValidIdentifier(valueName) {
			continue
		}

		if len(parts) > 1 {
			valueStr := strings.TrimSpace(parts[1])
			if v, ok := evaluateIntLiteral(valueStr); ok {
				currentValue = v
			}
		}

		values = append(values, EnumValue{
			Name:  valueName,
			Value: currentValue,
		})
		currentValue = currentValue + 1
	}

	return values
}

func (wg *WrapperGenerator) isStructDefinition(name, params string) bool {
	// Heuristic to detect if this is actually a struct definition
	return strings.Contains(params, "{") || strings.Contains(name, "struct")
}

// GenerateWrapper generates the wrapper code for the target language
func (wg *WrapperGenerator) GenerateWrapper(api *API) error {
	switch wg.config.Language {
	case "koda", "":
		return wg.generateKodaWrapper(api)
	default:
		return fmt.Errorf("unsupported language %q - only koda (.koda) is supported", wg.config.Language)
	}
}

// generateKodaWrapper generates an elegant, readable Koda wrapper library.
func (wg *WrapperGenerator) generateKodaWrapper(api *API) error {
	wrapperFile := filepath.Join(wg.config.OutputDir, api.Name+".koda")

	if wg.config.Verbose {
		fmt.Printf("  writing %s\n", wrapperFile)
	}

	file, err := os.Create(wrapperFile)
	if err != nil {
		return err
	}
	defer file.Close()

	w := func(format string, args ...interface{}) {
		fmt.Fprintf(file, format, args...)
	}

	hdr := wg.config.PrimaryHeader
	if hdr == "" && len(api.Headers) > 0 {
		hdr = filepath.Base(api.Headers[0])
	}
	ver := wg.config.Version
	if ver == "" {
		ver = WrapgenVersion
	}
	genAt := wg.config.GeneratedAt
	if genAt == "" {
		genAt = "(unknown)"
	}

	w("// ============================================================\n")
	w("//  %s.koda\n", api.Name)
	w("//  Auto-generated by %s %s\n", generatedByBrand, ver)
	w("//  Source:  %s\n", hdr)
	w("//  Generated: %s\n", genAt)
	w("// ============================================================\n")
	w("//\n")
	w("//  Usage:\n")
	w("//    #include \"@%s\"\n", api.Name)
	w("//\n")
	w("//  Compile:\n")
	w("//    set KODA_NATIVE_SOURCES=wrapper.c\n")
	w("//    set KODA_LINKFLAGS=-I<inc> -L<lib> -l%s\n", api.Name)
	w("//    koda build game.koda -o game.exe\n")
	w("// ============================================================\n\n")

	if len(api.Constants) > 0 {
		w("// ------------------------------------------------------------\n")
		w("//  Constants\n")
		w("// ------------------------------------------------------------\n\n")
		for _, constant := range api.Constants {
			wg.writeKodaConstant(file, constant)
		}
		w("\n")
	}

	emitMacros := wg.emitableMacros(api.Macros)
	if len(emitMacros) > 0 {
		w("// ------------------------------------------------------------\n")
		w("//  Macros (object-like #define)\n")
		w("// ------------------------------------------------------------\n\n")
		for _, macro := range emitMacros {
			wg.writeKodaMacro(file, macro)
		}
		w("\n")
	}

	if len(api.Structs) > 0 {
		w("// ------------------------------------------------------------\n")
		w("//  Structs\n")
		w("// ------------------------------------------------------------\n\n")
		for _, structDef := range api.Structs {
			wg.writeKodaStructDefinition(file, structDef)
		}
		w("\n")
	}

	if len(api.Enums) > 0 {
		w("// ------------------------------------------------------------\n")
		w("//  Enums\n")
		w("// ------------------------------------------------------------\n\n")
		for _, enum := range api.Enums {
			wg.writeKodaEnumDefinition(file, enum)
		}
		w("\n")
	}

	rawBy := make(map[string]string)
	for _, h := range api.Headers {
		if b, err := ioutil.ReadFile(h); err == nil {
			rawBy[h] = string(b)
		}
	}

	if len(api.Functions) > 0 {
		cats := groupFunctionsByCategory(api.Functions)
		for _, cat := range sortCategoryKeys(cats) {
			w("// ------------------------------------------------------------\n")
			w("//  %s\n", cat)
			w("// ------------------------------------------------------------\n")
			for _, function := range cats[cat] {
				if wg.skipGlueFunction(function) {
					continue
				}
				raw := rawBy[function.Header]
				wg.writeKodaFunctionDeclaration(file, function, raw)
			}
			w("\n")
		}
	}

	return wg.generateCGlue(api)
}

func (wg *WrapperGenerator) writeKodaFunctionDeclaration(file *os.File, function Function, rawHeader string) {
	ret := strings.TrimSpace(function.ReturnType)
	ret = strings.TrimPrefix(ret, "RLAPI ")
	ret = strings.TrimSpace(ret)

	desc := strings.TrimSpace(function.Documentation)
	if desc == "" {
		desc = functionDescription(rawHeader, function.Name)
	}

	fmt.Fprintf(file, "\n")
	fmt.Fprintf(file, "/// %s\n", desc)
	for i, param := range function.Parameters {
		name := strings.TrimSpace(param.Name)
		if name == "" {
			name = fmt.Sprintf("arg%d", i)
		}
		ft, hint := paramDocLine(name, param.Type)
		fmt.Fprintf(file, "/// @param %s  %s  — %s\n", name, ft, hint)
	}
	if has, line := returnsClause(ret); has {
		fmt.Fprintf(file, "%s\n", line)
	}

	fmt.Fprintf(file, "// koda:extern %s %s %d\n", function.Name, wg.wrapperSymbol(function), len(function.Parameters))
	fmt.Fprintf(file, "let %s = 0;\n", function.Name)
}

func (wg *WrapperGenerator) bindingParamList(function Function) string {
	names := make([]string, 0, len(function.Parameters))
	for i, param := range function.Parameters {
		name := strings.TrimSpace(param.Name)
		if name == "" {
			name = fmt.Sprintf("arg%d", i)
		}
		names = append(names, name)
	}
	return strings.Join(names, ", ")
}

func (wg *WrapperGenerator) wrapperSymbol(function Function) string {
	return "koda_wrap_" + sanitizeIdent(wg.config.LibraryName) + "_" + sanitizeIdent(function.Name)
}

func sanitizeIdent(s string) string {
	var b strings.Builder
	for i, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_' || (i > 0 && r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteByte('_')
		}
	}
	if b.Len() == 0 {
		return "wrapper"
	}
	return b.String()
}

func (wg *WrapperGenerator) generateCGlue(api *API) error {
	path := filepath.Join(wg.config.OutputDir, "wrapper.c")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	hdr := wg.config.PrimaryHeader
	if hdr == "" && len(api.Headers) > 0 {
		hdr = filepath.Base(api.Headers[0])
	}
	ver := wg.config.Version
	if ver == "" {
		ver = WrapgenVersion
	}
	genAt := wg.config.GeneratedAt
	if genAt == "" {
		genAt = "(unknown)"
	}

	fmt.Fprintf(file, "/* ============================================================\n")
	fmt.Fprintf(file, " *  wrapper.c\n")
	fmt.Fprintf(file, " *  Auto-generated by %s %s\n", generatedByBrand, ver)
	fmt.Fprintf(file, " *  Source:  %s\n", hdr)
	fmt.Fprintf(file, " *  Generated: %s\n", genAt)
	fmt.Fprintf(file, " *\n")
	fmt.Fprintf(file, " *  Link with koda build (KODA_NATIVE_SOURCES includes this file).\n")
	fmt.Fprintf(file, " * ============================================================ */\n\n")
	fmt.Fprintf(file, "#include \"koda_wrapgen_abi.h\"\n")
	fmt.Fprintf(file, "#include <stdbool.h>\n")
	fmt.Fprintf(file, "#include <string.h>\n")
	for _, h := range api.Headers {
		fmt.Fprintf(file, "#include \"%s\"\n", filepath.Base(h))
	}
	fmt.Fprintf(file, "\n")

	seen := make(map[string]bool)
	cats := groupFunctionsByCategory(api.Functions)
	for _, cat := range sortCategoryKeys(cats) {
		fmt.Fprintf(file, "/* ------------------------------------------------------------ */\n")
		fmt.Fprintf(file, "/*  %-58s */\n", cat)
		fmt.Fprintf(file, "/* ------------------------------------------------------------ */\n\n")
		for _, function := range cats[cat] {
			if function.Variadic {
				continue
			}
			if wg.skipGlueFunction(function) {
				continue
			}
			sym := wg.wrapperSymbol(function)
			if seen[sym] {
				continue
			}
			seen[sym] = true
			wg.writeCGlueFunction(file, api, function)
		}
	}
	return nil
}

func (wg *WrapperGenerator) findStruct(api *API, typeName string) *Struct {
	typeName = strings.TrimSpace(typeName)
	// Remove 'struct ' prefix if present
	typeName = strings.TrimPrefix(typeName, "struct ")

	// Direct search
	for _, s := range api.Structs {
		if s.Name == typeName {
			return &s
		}
	}

	// Try finding via typedef
	for _, td := range api.Typedefs {
		if td.Name == typeName {
			// Resolve the underlying type
			underlying := strings.TrimSpace(td.TargetType)
			underlying = strings.TrimPrefix(underlying, "struct ")
			underlying = strings.TrimPrefix(underlying, "typedef ")
			for _, s := range api.Structs {
				if s.Name == underlying {
					return &s
				}
			}
		}
	}

	return nil
}

func (wg *WrapperGenerator) skipGlueFunction(function Function) bool {
	n := strings.TrimSpace(function.Name)
	if n == "" || strings.ContainsAny(n, "() ") {
		return true
	}
	return !isValidWrapgenBindingName(n)
}

func (wg *WrapperGenerator) writeCGlueFunction(file *os.File, api *API, function Function) {
	retType := strings.TrimSpace(function.ReturnType)
	retType = strings.TrimPrefix(retType, "RLAPI ")
	retType = strings.TrimSpace(retType)

	fmt.Fprintf(file, "/**\n")
	fmt.Fprintf(file, " * %s — Koda wrapper for %s\n", wg.wrapperSymbol(function), function.Name)
	fmt.Fprintf(file, " *\n")
	fmt.Fprintf(file, " * Koda call:  %s(%s)\n", function.Name, wg.bindingParamList(function))
	for i, param := range function.Parameters {
		nm := strings.TrimSpace(param.Name)
		if nm == "" {
			nm = fmt.Sprintf("arg%d", i)
		}
		ft, _ := paramDocLine(nm, param.Type)
		fmt.Fprintf(file, " *   %s : %s → %s\n", nm, ft, strings.TrimSpace(wg.cTypeForParam(param.Type)))
	}
	if retType == "void" {
		fmt.Fprintf(file, " *   returns: nothing\n")
	} else {
		fmt.Fprintf(file, " *   returns: %s\n", kodaDocTypeForCType(retType))
	}
	fmt.Fprintf(file, " */\n")
	fmt.Fprintf(file, "KodaValue %s(int argCount, KodaValue* args) {\n", wg.wrapperSymbol(function))
	fmt.Fprintf(file, "    if (argCount < %d) {\n", len(function.Parameters))
	fmt.Fprintf(file, "        return koda_err_str(\"%s requires %d argument(s)\");\n", function.Name, len(function.Parameters))
	fmt.Fprintf(file, "    }\n")
	for i, param := range function.Parameters {
		valExpr := fmt.Sprintf("args[%d]", i)
		fmt.Fprintf(file, "    %s arg%d = %s;\n", wg.cTypeForParam(param.Type), i, wg.cArgExprKoda(api, param.Type, valExpr))
	}
	if retType == "void" {
		fmt.Fprintf(file, "    %s(", function.Name)
		for i := range function.Parameters {
			if i > 0 {
				fmt.Fprintf(file, ", ")
			}
			fmt.Fprintf(file, "arg%d", i)
		}
		fmt.Fprintf(file, ");\n")
		fmt.Fprintf(file, "    return NULL_VAL;\n")
		fmt.Fprintf(file, "}\n\n")
		return
	}
	retVarT := wg.cVarTypeForReturn(retType)
	fmt.Fprintf(file, "    %s result = %s(", retVarT, function.Name)
	for i := range function.Parameters {
		if i > 0 {
			fmt.Fprintf(file, ", ")
		}
		fmt.Fprintf(file, "arg%d", i)
	}
	fmt.Fprintf(file, ");\n")
	fmt.Fprintf(file, "    %s\n", wg.cReturnExprKoda(retType, api))
	fmt.Fprintf(file, "}\n\n")
}

func (wg *WrapperGenerator) cParamList(function Function) string {
	parts := make([]string, 0, len(function.Parameters))
	for i, param := range function.Parameters {
		name := strings.TrimSpace(param.Name)
		if name == "" {
			name = fmt.Sprintf("arg%d", i)
		}
		parts = append(parts, strings.TrimSpace(param.Type)+" "+name)
	}
	return strings.Join(parts, ", ")
}

func (wg *WrapperGenerator) cTypeForParam(t string) string {
	t = strings.TrimSpace(t)
	if strings.Contains(t, "*") {
		if strings.Contains(t, "char") && strings.Count(t, "*") >= 2 {
			return t
		}
		if strings.Contains(t, "char") {
			if strings.Contains(t, "unsigned") {
				return "unsigned char*"
			}
			return "const char*"
		}
		return t
	}
	if strings.Contains(t, "bool") {
		return "int"
	}
	if strings.Contains(t, "double") {
		return "double"
	}
	if strings.Contains(t, "float") {
		return "float"
	}
	return t
}

func (wg *WrapperGenerator) cTypeForReturn(t string) string {
	t = strings.TrimSpace(t)
	if t == "void" {
		return "void"
	}
	if strings.Contains(t, "bool") {
		return "bool"
	}
	return wg.cTypeForParam(t)
}

func (wg *WrapperGenerator) resolveTypedef(api *API, t string) (string, bool) {
	t = strings.TrimSpace(t)
	for _, td := range api.Typedefs {
		if td.Name == t {
			return strings.TrimSpace(td.TargetType), true
		}
	}
	return t, false
}

func (wg *WrapperGenerator) isPointerType(api *API, t string) bool {
	t = strings.TrimSpace(t)
	if strings.Contains(t, "*") {
		return true
	}
	if resolved, ok := wg.resolveTypedef(api, t); ok {
		return wg.isPointerType(api, resolved)
	}
	return false
}

func isFunctionPointerParamType(t string) bool {
	t = strings.TrimSpace(t)
	if strings.Contains(t, "(*)(") {
		return true
	}
	if strings.HasSuffix(t, "Callback") {
		return true
	}
	return false
}

func (wg *WrapperGenerator) cArgExprKoda(api *API, t string, valExpr string) string {
	t = strings.TrimSpace(t)
	if isFunctionPointerParamType(t) {
		return "(" + t + ")0"
	}
	if strings.Contains(t, "unsigned") && strings.Contains(t, "char") && strings.Contains(t, "*") {
		return fmt.Sprintf(`((unsigned char*)(void*)(IS_OBJ(%s) && AS_OBJ(%s)->type == OBJ_STRING ? ((ObjString*)AS_OBJ(%s))->chars : (const char*)""))`, valExpr, valExpr, valExpr)
	}
	if strings.Contains(t, "char") && strings.Contains(t, "*") {
		if strings.Count(t, "*") >= 2 {
			return "(" + t + ")NULL"
		}
		return fmt.Sprintf(`(IS_OBJ(%s) && AS_OBJ(%s)->type == OBJ_STRING ? ((ObjString*)AS_OBJ(%s))->chars : "")`, valExpr, valExpr, valExpr)
	}
	if strings.Contains(t, "bool") {
		return fmt.Sprintf("(IS_BOOL(%s) ? (AS_BOOL(%s) ? 1 : 0) : 0)", valExpr, valExpr)
	}

	st := wg.findStruct(api, t)
	if st != nil {
		// Generate recursive field extraction
		fields := ""
		for i, field := range st.Fields {
			if i > 0 {
				fields += ", "
			}
			if strings.Contains(field.Type, "[") {
				fields += "{0}"
				continue
			}
			fieldValExpr := fmt.Sprintf("koda_get_index(%s, koda_copy_string(\"%s\", %d))", valExpr, field.Name, len(field.Name))
			fields += wg.cArgExprKoda(api, field.Type, fieldValExpr)
		}
		return fmt.Sprintf("(%s){ %s }", st.Name, fields)
	}

	if wg.isPointerType(api, t) {
		pt := t
		if resolved, ok := wg.resolveTypedef(api, t); ok && strings.Contains(resolved, "*") {
			pt = resolved
		} else if !strings.Contains(pt, "*") {
			pt = pt + "*"
		}
		return "(" + pt + ")NULL"
	}
	if resolved, ok := wg.resolveTypedef(api, t); ok && wg.isPointerType(api, resolved) {
		return "(" + resolved + ")NULL"
	}
	if strings.Contains(t, "float") || strings.Contains(t, "double") {
		return fmt.Sprintf("(IS_NUMBER(%s) ? (double)AS_NUMBER(%s) : 0.0)", valExpr, valExpr)
	}
	if strings.Contains(t, "int") || strings.Contains(t, "long") || strings.Contains(t, "short") || strings.Contains(t, "size_t") {
		return fmt.Sprintf("(IS_NUMBER(%s) ? (int)AS_NUMBER(%s) : 0)", valExpr, valExpr)
	}
	return fmt.Sprintf("(IS_NUMBER(%s) ? (int)AS_NUMBER(%s) : 0)", valExpr, valExpr)
}

func (wg *WrapperGenerator) cVarTypeForReturn(t string) string {
	t = strings.TrimSpace(t)
	if strings.Contains(t, "bool") {
		return "bool"
	}
	if strings.Contains(t, "char") && strings.Contains(t, "*") {
		return "const char*"
	}
	if strings.Contains(t, "float") && !strings.Contains(t, "double") {
		return "float"
	}
	if strings.Contains(t, "double") {
		return "double"
	}
	if strings.Contains(t, "*") {
		return "void*"
	}
	return wg.cTypeForParam(t)
}

func (wg *WrapperGenerator) cReturnExprKoda(t string, api *API) string {
	t = strings.TrimSpace(t)
	if isFunctionPointerParamType(t) {
		return "return NULL_VAL;"
	}
	if strings.Contains(t, "char") && strings.Contains(t, "*") {
		return "return result ? koda_copy_string(result, (int)strlen(result)) : NULL_VAL;"
	}
	if strings.Contains(t, "bool") {
		return "return BOOL_VAL(result);"
	}

	st := wg.findStruct(api, t)
	if st != nil {
		return fmt.Sprintf("return %s;", wg.cReturnStructExpr(api, st, "result"))
	}

	if strings.Contains(t, "*") {
		return "return NULL_VAL;"
	}
	return "return NUMBER_VAL((double)result);"
}

func (wg *WrapperGenerator) cReturnStructExpr(api *API, st *Struct, valExpr string) string {
	// Create a new table object and fill its fields (koda_allocate_object returns Value).
	res := fmt.Sprintf("({ KodaValue _obj = koda_allocate_object(%d); ", len(st.Fields))
	for _, field := range st.Fields {
		fType := strings.TrimSpace(field.Type)
		fValExpr := fmt.Sprintf("%s.%s", valExpr, field.Name)

		// Handle nested structs recursively
		nestedSt := wg.findStruct(api, fType)

		var fieldValCode string
		if nestedSt != nil {
			fieldValCode = wg.cReturnStructExpr(api, nestedSt, fValExpr)
		} else if strings.Contains(fType, "[") {
			fieldValCode = "NULL_VAL"
		} else if strings.Contains(fType, "**") || strings.Count(fType, "*") >= 2 {
			fieldValCode = "NULL_VAL"
		} else if strings.Contains(fType, "char") && strings.Contains(fType, "*") {
			fieldValCode = fmt.Sprintf("%s ? koda_copy_string(%s, (int)strlen(%s)) : NULL_VAL", fValExpr, fValExpr, fValExpr)
		} else if strings.Contains(fType, "*") {
			fieldValCode = "NULL_VAL"
		} else if strings.Contains(fType, "bool") {
			fieldValCode = fmt.Sprintf("BOOL_VAL(%s)", fValExpr)
		} else if strings.Contains(fType, "float") || strings.Contains(fType, "double") || strings.Contains(fType, "int") || strings.Contains(fType, "long") || strings.Contains(fType, "short") || strings.Contains(fType, "size_t") || strings.Contains(fType, "unsigned") {
			fieldValCode = fmt.Sprintf("NUMBER_VAL((double)%s)", fValExpr)
		} else {
			fieldValCode = "NULL_VAL"
		}

		res += fmt.Sprintf("koda_object_set(_obj, koda_copy_string(\"%s\", %d), %s); ", field.Name, len(field.Name), fieldValCode)
	}
	res += "_obj; })"
	return res
}

func (wg *WrapperGenerator) exampleArgs(function Function) string {
	args := make([]string, 0, len(function.Parameters))
	for i, param := range function.Parameters {
		name := strings.TrimSpace(param.Name)
		if name == "" {
			name = fmt.Sprintf("arg%d", i)
		}
		switch {
		case strings.Contains(param.Type, "char") && strings.Contains(param.Type, "*"):
			args = append(args, "\"text\"")
		case strings.Contains(param.Type, "*"):
			args = append(args, "0")
		default:
			args = append(args, name)
		}
	}
	return strings.Join(args, ", ")
}

func (wg *WrapperGenerator) writeKodaStructDefinition(file *os.File, structDef Struct) {
	fmt.Fprintf(file, "// Struct: %s\n", structDef.Name)
	for _, field := range structDef.Fields {
		fmt.Fprintf(file, "//   %s: %s\n", field.Name, field.Type)
	}
	fmt.Fprintf(file, "\n")
}

func (wg *WrapperGenerator) writeKodaEnumDefinition(file *os.File, enum Enum) {
	fmt.Fprintf(file, "// Enum: %s\n", enum.Name)
	seen := make(map[string]struct{})
	for _, value := range enum.Values {
		if !isValidWrapgenBindingName(value.Name) {
			continue
		}
		if _, dup := seen[value.Name]; dup {
			continue
		}
		seen[value.Name] = struct{}{}
		fmt.Fprintf(file, "let %s = %s; // %s\n", value.Name, formatIntLiteral(value.Value), enum.Name)
	}
	fmt.Fprintf(file, "\n")
}

var skipMacroNames = map[string]struct{}{
	"bool": {}, "true": {}, "false": {}, "NULL": {}, "RLAPI": {}, "CLITERAL": {},
	"RAYLIB_H": {}, "RLGL_H": {}, "RAYMATH_H": {},
}

func (wg *WrapperGenerator) emitableMacros(macros []Macro) []Macro {
	var out []Macro
	seen := make(map[string]struct{})
	for _, m := range macros {
		if len(m.Parameters) > 0 {
			continue
		}
		if _, skip := skipMacroNames[m.Name]; skip {
			continue
		}
		if !isValidWrapgenBindingName(m.Name) {
			continue
		}
		if _, dup := seen[m.Name]; dup {
			continue
		}
		val := strings.TrimSpace(m.Value)
		if val == "" || strings.ContainsAny(val, "{(\\") {
			continue
		}
		if strings.HasPrefix(val, "CLITERAL") {
			continue
		}
		if _, ok := evaluateIntLiteral(val); !ok {
			continue
		}
		seen[m.Name] = struct{}{}
		out = append(out, m)
	}
	return out
}

func (wg *WrapperGenerator) writeKodaMacro(file *os.File, macro Macro) {
	val := strings.TrimSpace(macro.Value)
	if v, ok := evaluateIntLiteral(val); ok {
		fmt.Fprintf(file, "let %s = %s; // #define\n", macro.Name, formatIntLiteral(v))
	}
}

func (wg *WrapperGenerator) writeKodaConstant(file *os.File, constant Constant) {
	fmt.Fprintf(file, "let %s = %s; // %s\n", constant.Name, constant.Value, constant.Type)
}

// analyzeDependencies analyzes dependencies and relationships in the API
func (wg *WrapperGenerator) analyzeDependencies(api *API) {
	// Build type dependency graph
	typeMap := make(map[string]bool)

	// Collect all defined types
	for _, structDef := range api.Structs {
		typeMap[structDef.Name] = true
	}
	for _, enum := range api.Enums {
		typeMap[enum.Name] = true
	}
	for _, typedef := range api.Typedefs {
		typeMap[typedef.Name] = true
	}

	// Analyze function parameter and return types for dependencies
	for _, function := range api.Functions {
		// Check return type
		if typeMap[function.ReturnType] {
			api.Dependencies = append(api.Dependencies, function.ReturnType)
		}

		// Check parameter types
		for _, param := range function.Parameters {
			if typeMap[param.Type] {
				api.Dependencies = append(api.Dependencies, param.Type)
			}
		}
	}

	// Remove duplicates
	uniqueDeps := make(map[string]bool)
	var cleanDeps []string
	for _, dep := range api.Dependencies {
		if !uniqueDeps[dep] {
			uniqueDeps[dep] = true
			cleanDeps = append(cleanDeps, dep)
		}
	}
	api.Dependencies = cleanDeps
}

