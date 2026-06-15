package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"koda/internal/kodahome"
)

// ClangParser uses clang to extract accurate C/C++ API information
type ClangParser struct {
	config  *WrapGenConfig
	tempDir string
}

func NewClangParser(config *WrapGenConfig) *ClangParser {
	return &ClangParser{
		config: config,
	}
}

// ParseWithClang uses clang to extract AST information from headers
func (cp *ClangParser) ParseWithClang(headers []string) (*API, error) {
	// Create temporary directory for clang output
	tempDir, err := os.MkdirTemp("", "wrapgen_clang")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cp.tempDir = tempDir

	var extractLog func(string)
	if cp.config.Verbose {
		extractLog = func(s string) { fmt.Print(s) }
	}
	if err := kodahome.EnsureEnvironment(extractLog); err != nil {
		return nil, fmt.Errorf("kodawrap needs the same bundled Clang as koda: %w", err)
	}

	api := &API{
		Name: cp.config.LibraryName,
	}

	// Generate clang AST dump for each header
	for _, header := range headers {
		headerAPI, err := cp.parseHeaderWithClang(header)
		if err != nil {
			if cp.config.Verbose {
				fmt.Printf("Warning: clang parsing failed for %s, falling back to regex: %v\n", header, err)
			}
			// Fallback to regex parsing
			content, err := os.ReadFile(header)
			if err != nil {
				return nil, fmt.Errorf("failed to read header %s: %v", header, err)
			}
			headerAPI = cp.parseWithRegex(string(content), header)
		}

		// Merge API information
		api.Functions = append(api.Functions, headerAPI.Functions...)
		api.Structs = append(api.Structs, headerAPI.Structs...)
		api.Enums = append(api.Enums, headerAPI.Enums...)
		api.Macros = append(api.Macros, headerAPI.Macros...)
		api.Typedefs = append(api.Typedefs, headerAPI.Typedefs...)
		api.Constants = append(api.Constants, headerAPI.Constants...)
		api.Headers = append(api.Headers, header)
	}

	return api, nil
}

// parseHeaderWithClang uses clang to parse a single header file
func (cp *ClangParser) parseHeaderWithClang(headerPath string) (*API, error) {
	astFile := filepath.Join(cp.tempDir, "ast.txt")

	// Run clang to generate AST dump; -isystem picks up portable headers next to the toolchain.
	var inc []string
	for _, d := range cp.config.IncludePaths {
		d = strings.TrimSpace(d)
		if d != "" {
			inc = append(inc, "-I", d)
		}
	}
	args := kodahome.ClangWrappedArgs(append(inc, cp.clangLangArgs(headerPath)...)...)
	args = append(args, "-Xclang", "-ast-dump", "-fsyntax-only", headerPath)
	cmd := exec.Command(kodahome.Clang(), args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("clang failed: %v, output: %s", err, string(output))
	}

	// Save AST dump for debugging
	if cp.config.Verbose {
		err := os.WriteFile(astFile, output, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to save AST dump: %v", err)
		}
	}

	return cp.parseASTDump(string(output), headerPath)
}

// parseASTDump parses clang AST dump output
func (cp *ClangParser) parseASTDump(astDump, headerPath string) (*API, error) {
	api := &API{}

	scanner := bufio.NewScanner(strings.NewReader(astDump))

	for scanner.Scan() {
		line := scanner.Text()

		// Parse function declarations
		if funcDecl := cp.parseFunctionDecl(line); funcDecl != nil {
			funcDecl.Header = headerPath
			api.Functions = append(api.Functions, *funcDecl)
		}

		// Parse struct declarations
		if structDecl := cp.parseStructDecl(line); structDecl != nil {
			structDecl.Header = headerPath
			api.Structs = append(api.Structs, *structDecl)
		}

		// Parse enum declarations
		if enumDecl := cp.parseEnumDecl(line); enumDecl != nil {
			enumDecl.Header = headerPath
			api.Enums = append(api.Enums, *enumDecl)
		}

		// Parse typedef declarations
		if typedef := cp.parseTypedefDecl(line); typedef != nil {
			typedef.Header = headerPath
			api.Typedefs = append(api.Typedefs, *typedef)
		}

		// Parse variable declarations (constants)
		if varDecl := cp.parseVarDecl(line); varDecl != nil {
			varDecl.Header = headerPath
			api.Constants = append(api.Constants, *varDecl)
		}
	}

	return api, nil
}

// parseFunctionDecl extracts function declarations from AST line
func (cp *ClangParser) parseFunctionDecl(line string) *Function {
	patterns := []string{
		`FunctionDecl\s+0x[0-9a-f]+\s+<[^>]+>\s+(?:inline\s+)?(?:static\s+)?(\w+)\s+'([^']+)'\s*`,
		`CXXMethodDecl\s+0x[0-9a-f]+\s+<[^>]+>\s+(?:[\w:]+\s+)?(\w+)\s+'([^']+)'\s*`,
	}
	for _, pat := range patterns {
		re := regexp.MustCompile(pat)
		matches := re.FindStringSubmatch(line)
		if len(matches) >= 3 {
			funcName := matches[1]
			signature := matches[2]
			function := &Function{Name: funcName}
			cp.parseFunctionSignature(signature, function)
			return function
		}
	}
	return nil
}

func (cp *ClangParser) clangLangArgs(headerPath string) []string {
	if cp.config.UseCPP || isCPPHeader(headerPath) {
		return []string{"-x", "c++", "-std=c++17"}
	}
	return []string{"-x", "c"}
}

// parseStructDecl extracts struct declarations from AST line
func (cp *ClangParser) parseStructDecl(line string) *Struct {
	// Match struct declaration pattern
	re := regexp.MustCompile(`RecordDecl\s+0x[0-9a-f]+\s+<[^>]+>\s+(?:struct|class)\s+(\w+)`)
	matches := re.FindStringSubmatch(line)

	if len(matches) >= 2 {
		structName := matches[1]

		return &Struct{
			Name: structName,
		}
	}

	return nil
}

// parseEnumDecl extracts enum declarations from AST line
func (cp *ClangParser) parseEnumDecl(line string) *Enum {
	// Match enum declaration pattern
	re := regexp.MustCompile(`EnumDecl\s+0x[0-9a-f]+\s+<[^>]+>\s+(?:typedef\s+)?(\w+)`)
	matches := re.FindStringSubmatch(line)

	if len(matches) >= 2 {
		enumName := matches[1]

		return &Enum{
			Name: enumName,
		}
	}

	return nil
}

// parseTypedefDecl extracts typedef declarations from AST line
func (cp *ClangParser) parseTypedefDecl(line string) *Typedef {
	// Match typedef declaration pattern
	re := regexp.MustCompile(`TypedefDecl\s+0x[0-9a-f]+\s+<[^>]+>\s+(\w+)\s+'([^']+)'\s*(?:\w+\s+)?`)
	matches := re.FindStringSubmatch(line)

	if len(matches) >= 3 {
		typeName := matches[1]
		targetType := matches[2]

		return &Typedef{
			Name:       typeName,
			TargetType: targetType,
		}
	}

	return nil
}

// parseVarDecl extracts variable declarations from AST line
func (cp *ClangParser) parseVarDecl(line string) *Constant {
	// Match variable declaration pattern
	re := regexp.MustCompile(`VarDecl\s+0x[0-9a-f]+\s+<[^>]+>\s+(?:const\s+)?(\w+)\s+'([^']+)'\s+(?:extern|static)?\s*(?:\w+\s+)?([^'\s]+)`)
	matches := re.FindStringSubmatch(line)

	if len(matches) >= 4 {
		varName := matches[1]
		varType := matches[2]
		varValue := matches[3]

		return &Constant{
			Name:  varName,
			Type:  varType,
			Value: varValue,
		}
	}

	return nil
}

// parseFunctionSignature parses function signature string
func (cp *ClangParser) parseFunctionSignature(signature string, function *Function) {
	// Simple parsing - in practice you'd want more sophisticated parsing
	parts := strings.Split(signature, "(")
	if len(parts) >= 2 {
		returnType := strings.TrimSpace(parts[0])
		function.ReturnType = returnType

		paramsStr := strings.Join(parts[1:], "(")
		paramsStr = strings.TrimSuffix(paramsStr, ")")

		function.Parameters = cp.parseParametersFromSignature(paramsStr)
	}
}

// parseParametersFromSignature parses parameters from function signature
func (cp *ClangParser) parseParametersFromSignature(paramsStr string) []Parameter {
	var parameters []Parameter

	if paramsStr == "" || paramsStr == "void" {
		return parameters
	}

	paramList := strings.Split(paramsStr, ",")
	for _, param := range paramList {
		param = strings.TrimSpace(param)
		if param == "" {
			continue
		}

		// Split type and name
		parts := strings.Fields(param)
		if len(parts) >= 2 {
			paramName := parts[len(parts)-1]
			paramType := strings.Join(parts[:len(parts)-1], " ")

			parameters = append(parameters, Parameter{
				Name: paramName,
				Type: paramType,
			})
		}
	}

	return parameters
}

// parseWithRegex fallback parsing using regex patterns
func (cp *ClangParser) parseWithRegex(content, headerPath string) *API {
	api := &API{}

	// Extract functions
	api.Functions = cp.extractFunctionsRegex(content, headerPath)

	// Extract structs
	api.Structs = cp.extractStructsRegex(content, headerPath)

	// Extract enums
	api.Enums = cp.extractEnumsRegex(content, headerPath)

	// Extract macros
	api.Macros = cp.extractMacrosRegex(content, headerPath)

	// Extract typedefs
	api.Typedefs = cp.extractTypedefsRegex(content, headerPath)

	// Extract constants
	api.Constants = cp.extractConstantsRegex(content, headerPath)

	return api
}

// Regex-based extraction methods (enhanced versions)
func (cp *ClangParser) extractFunctionsRegex(content, headerPath string) []Function {
	var functions []Function

	// Enhanced function regex
	funcPattern := regexp.MustCompile(`(?m)^\s*(?:extern\s+)?(?:inline\s+)?(?:static\s+)?(?:__attribute__\s*\([^)]*\)\s+)*(?:\w+\s+)*?(\w+)\s*\(([^)]*)\)\s*(?:__attribute__\s*\([^)]*\))?\s*;?`)

	matches := funcPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			funcName := strings.TrimSpace(match[1])
			paramsStr := strings.TrimSpace(match[2])

			function := Function{
				Name:   funcName,
				Header: headerPath,
			}

			// Extract return type
			fullMatch := match[0]
			beforeFunc := strings.Split(fullMatch, funcName)[0]
			returnType := strings.TrimSpace(beforeFunc)
			returnType = regexp.MustCompile(`\b(?:extern|inline|static|__attribute__\s*\([^)]*\))\b`).ReplaceAllString(returnType, "")
			function.ReturnType = strings.TrimSpace(returnType)

			// Parse parameters
			function.Parameters = cp.parseParametersFromSignature(paramsStr)

			// Check for variadic
			if strings.Contains(paramsStr, "...") {
				function.Variadic = true
			}

			functions = append(functions, function)
		}
	}

	return functions
}

func (cp *ClangParser) extractStructsRegex(content, headerPath string) []Struct {
	var structs []Struct

	// Enhanced struct regex
	structPattern := regexp.MustCompile(`(?s)typedef\s+struct\s+(\w*)\s*\{([^}]+)\}\s*(\w+)[^;]*;|(?s)struct\s+(\w+)\s*\{([^}]+)\}\s*;`)

	matches := structPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		var structName string
		var fieldsContent string

		// Handle both typedef struct and struct patterns
		if match[1] != "" && match[3] != "" {
			// typedef struct pattern
			structName = strings.TrimSpace(match[3])
			fieldsContent = match[2]
		} else if match[4] != "" && match[5] != "" {
			// struct pattern
			structName = strings.TrimSpace(match[4])
			fieldsContent = match[5]
		}

		if structName != "" {
			structDef := Struct{
				Name:   structName,
				Header: headerPath,
			}

			// Parse fields
			structDef.Fields = cp.parseStructFields(fieldsContent)

			structs = append(structs, structDef)
		}
	}

	return structs
}

func (cp *ClangParser) extractEnumsRegex(content, headerPath string) []Enum {
	var enums []Enum

	// Enhanced enum regex
	enumPattern := regexp.MustCompile(`(?s)typedef\s+enum\s+(\w*)\s*\{([^}]+)\}\s*(\w+)[^;]*;|(?s)enum\s+(\w+)\s*\{([^}]+)\}\s*;`)

	matches := enumPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		var enumName string
		var valuesContent string

		// Handle both typedef enum and enum patterns
		if match[1] != "" && match[3] != "" {
			// typedef enum pattern
			enumName = strings.TrimSpace(match[3])
			valuesContent = match[2]
		} else if match[4] != "" && match[5] != "" {
			// enum pattern
			enumName = strings.TrimSpace(match[4])
			valuesContent = match[5]
		}

		if enumName != "" {
			enumDef := Enum{
				Name:   enumName,
				Header: headerPath,
			}

			// Parse enum values
			enumDef.Values = cp.parseEnumValues(valuesContent)

			enums = append(enums, enumDef)
		}
	}

	return enums
}

func (cp *ClangParser) extractMacrosRegex(content, headerPath string) []Macro {
	var macros []Macro

	// Enhanced macro regex
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

func (cp *ClangParser) extractTypedefsRegex(content, headerPath string) []Typedef {
	var typedefs []Typedef

	// Enhanced typedef regex
	typedefPattern := regexp.MustCompile(`(?m)typedef\s+(?!enum|struct|union)([^;]+);`)

	matches := typedefPattern.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) >= 2 {
			typedefStr := strings.TrimSpace(match[1])
			parts := strings.Fields(typedefStr)

			if len(parts) >= 2 {
				typeName := parts[len(parts)-1]
				targetType := strings.Join(parts[:len(parts)-1], " ")

				typedef := Typedef{
					Name:       typeName,
					TargetType: targetType,
					Header:     headerPath,
				}

				typedefs = append(typedefs, typedef)
			}
		}
	}

	return typedefs
}

func (cp *ClangParser) extractConstantsRegex(content, headerPath string) []Constant {
	var constants []Constant

	// Enhanced constant regex
	constPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?m)^extern\s+const\s+(\w+)\s+(\w+)[^;]*;`),
		regexp.MustCompile(`(?m)^const\s+(\w+)\s+(\w+)\s*=\s*([^;]+);`),
		regexp.MustCompile(`(?m)^#define\s+(\w+)\s+(.+)$`),
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

				// Handle #define constants
				if pattern.String()[0] == '#' {
					constType = "macro"
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

func (cp *ClangParser) parseStructFields(fieldsContent string) []Field {
	var fields []Field

	lines := strings.Split(fieldsContent, ";")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// Remove comments
		if commentIndex := strings.Index(line, "//"); commentIndex != -1 {
			line = line[:commentIndex]
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			fieldName := parts[len(parts)-1]
			fieldType := strings.Join(parts[:len(parts)-1], " ")
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

func (cp *ClangParser) parseEnumValues(valuesContent string) []EnumValue {
	var values []EnumValue

	lines := strings.Split(valuesContent, ",")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "//") {
			continue
		}

		// Remove comments
		if commentIndex := strings.Index(line, "//"); commentIndex != -1 {
			line = line[:commentIndex]
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, "=")
		valueName := strings.TrimSpace(parts[0])

		var value int64 = int64(i)
		if len(parts) > 1 {
			valueStr := strings.TrimSpace(parts[1])
			if v, ok := evaluateIntLiteral(valueStr); ok {
				value = v
			}
		}

		values = append(values, EnumValue{
			Name:  valueName,
			Value: value,
		})
	}

	return values
}

// evaluateExpression evaluates simple constant expressions
func (cp *ClangParser) evaluateExpression(expr string) (int64, error) {
	// Remove common prefixes
	expr = strings.TrimSpace(expr)
	expr = strings.ReplaceAll(expr, "U", "")
	expr = strings.ReplaceAll(expr, "L", "")
	expr = strings.ReplaceAll(expr, "u", "")
	expr = strings.ReplaceAll(expr, "l", "")

	// Handle hex, octal, and decimal
	if strings.HasPrefix(expr, "0x") || strings.HasPrefix(expr, "0X") {
		var value int64
		_, err := fmt.Sscanf(expr, "%x", &value)
		return value, err
	} else if strings.HasPrefix(expr, "0") && len(expr) > 1 {
		var value int64
		_, err := fmt.Sscanf(expr, "%o", &value)
		return value, err
	} else {
		var value int64
		_, err := fmt.Sscanf(expr, "%d", &value)
		return value, err
	}
}
