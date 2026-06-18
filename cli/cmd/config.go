package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"KProbeCLI/db"
	"KProbeCLI/helpers"
)

const (
	INVALID_FORMAT       string = "Invalid config file format"
	DUPL_SCAN_NAME       string = "Duplicate scan name detected"
	INVALID_SCAN_NAME    string = "Invalid scan name"
	TIMEOUT_NOT_INT      string = "Timeout is not an integer"
	TIMEOUT_OUT_OF_RANGE string = "Timeout out of range"
	INVALID_CODE         string = "Invalid status code"
	CODE_TOO_LONG        string = "Status code too long"
)

var varNameRegex = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// splitConfigLine splits a config line into parts, treating quoted values as single tokens.
// Matches the web-editor's splitConfigLine function.
func splitConfigLine(line string) []string {
	var parts []string
	i := 0
	runes := []rune(line)

	for i < len(runes) {
		// Skip whitespace
		if runes[i] == ' ' || runes[i] == '\t' {
			i++
			continue
		}

		var part strings.Builder
		for i < len(runes) && runes[i] != ' ' && runes[i] != '\t' {
			if runes[i] == '"' {
				// Enter quoted value
				part.WriteRune('"')
				i++
				for i < len(runes) && runes[i] != '"' {
					if runes[i] == '\\' && i+1 < len(runes) && runes[i+1] == '"' {
						part.WriteString("\\\"")
						i += 2
					} else {
						part.WriteRune(runes[i])
						i++
					}
				}
				if i < len(runes) {
					part.WriteRune('"')
					i++
				}
			} else {
				part.WriteRune(runes[i])
				i++
			}
		}

		if part.Len() > 0 {
			parts = append(parts, part.String())
		}
	}

	return parts
}

// stripQuotes removes surrounding double quotes from a string
func stripQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

// parseConfigLine parses key=value pairs from fields[3:] and returns a Scan
func parseConfigLine(fields []string) helpers.Scan {
	scan := helpers.Scan{
		Name:    fields[0],
		Type:    fields[1],
		Address: fields[2],
		Timeout: 10,
	}

	var variables []helpers.Variable

	for p := 3; p < len(fields); p++ {
		part := fields[p]
		eqIdx := strings.Index(part, "=")
		if eqIdx == -1 {
			continue
		}

		key := part[:eqIdx]
		value := stripQuotes(part[eqIdx+1:])

		switch {
		case key == "timeout":
			t, ok := helpers.StrToInt(value)
			if ok {
				scan.Timeout = t
			}

		case key == "status_code":
			scan.StatusCode = value

		case key == "expression":
			scan.Expression = value

		case strings.HasPrefix(key, "var_bool:"):
			varName := key[len("var_bool:"):]
			variables = append(variables, helpers.Variable{
				Name:  varName,
				Mode:  "bool",
				Value: value,
			})

		case strings.HasPrefix(key, "var_json_bool:"):
			varName := key[len("var_json_bool:"):]
			variables = append(variables, helpers.Variable{
				Name:  varName,
				Mode:  "json",
				Query: value,
				As:    "bool",
			})

		case strings.HasPrefix(key, "var_json_number:"):
			varName := key[len("var_json_number:"):]
			variables = append(variables, helpers.Variable{
				Name:  varName,
				Mode:  "json",
				Query: value,
				As:    "number",
			})
		}
	}

	if len(variables) > 0 {
		scan.Variables = variables
	}

	return scan
}

func VerifyConfig(path string) {
	helpers.PrintInfo("Verifying config file")

	file, err := os.Open(path)
	if err != nil {
		helpers.PrintError(true, "Failed to open file ("+err.Error()+")")
	}
	defer file.Close()

	scanNames := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := splitConfigLine(line)
		if len(fields) == 0 {
			continue
		}
		if len(fields) < 3 {
			helpers.PrintError(true, INVALID_FORMAT)
		}

		scanName := fields[0]
		if err := validateScanName(scanName, scanNames); err != "" {
			helpers.PrintError(true, "Invalid scan name: "+err)
		}

		scanType := fields[1]
		if scanType != "http" && scanType != "ping" {
			helpers.PrintError(true, "Invalid scan type: "+scanType)
		}

		scanAddress := fields[2]
		if len(scanAddress) > 256 {
			helpers.PrintError(true, "Invalid scan address")
		}

		// Parse all key=value pairs
		parsed := parseConfigLine(fields)

		// Validate timeout
		if errMsg := validateTimeout(helpers.IntToStr(parsed.Timeout)); errMsg != "" {
			helpers.PrintError(true, "Invalid timeout: "+errMsg)
		}

		// Validate status code
		if parsed.StatusCode != "" {
			if scanType == "ping" {
				helpers.PrintError(true, "Status code is not supported for ping scans")
			}
			if errMsg := validateStatusCode(parsed.StatusCode); errMsg != "" {
				helpers.PrintError(true, "Invalid status code: "+errMsg)
			}
		}

		// Validate variables and expression (HTTP only)
		if scanType == "http" {
			if errMsg := validateVariables(parsed.Variables); errMsg != "" {
				helpers.PrintError(true, errMsg+" (scan: "+scanName+")")
			}

			if errMsg := validateExpressionConfig(parsed.Variables, parsed.Expression); errMsg != "" {
				helpers.PrintError(true, errMsg+" (scan: "+scanName+")")
			}
		} else {
			if len(parsed.Variables) > 0 || parsed.Expression != "" {
				helpers.PrintError(true, "Variables and expressions are not supported for ping scans (scan: "+scanName+")")
			}
		}

		scanNames[scanName] = true
	}

	if err := scanner.Err(); err != nil {
		helpers.PrintError(true, "Failed to read file ("+err.Error()+")")
	}

	helpers.PrintSuccess("Config file verified successfully")
}

func SetConfig(path string) {
	helpers.PrintInfo("Replacing config file")

	helpers.PrintQuestion("Do you want to replace the config file? (y/n)")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input != "y" && input != "Y" {
		helpers.PrintWarning("Config file replacement aborted")
		return
	}

	helpers.PrintInfo("Deleting old config file")
	db.DeleteScans()
	helpers.PrintSuccess("Old config file deleted successfully")
	helpers.PrintInfo("Adding new scans")

	file, err := os.Open(path)
	if err != nil {
		helpers.PrintError(true, "Failed to open file ("+err.Error()+")")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		fields := splitConfigLine(line)
		scan := parseConfigLine(fields)

		helpers.PrintInfo("Adding scan: " + scan.Name)
		db.AddScan(scan)

		db.InsertValue("config_set", "true")
		helpers.PrintSuccess("Scan added successfully")
	}

	if err := scanner.Err(); err != nil {
		helpers.PrintError(true, "Failed to read file ("+err.Error()+")")
	}

	helpers.PrintSuccess("Config file replaced successfully")
}

func ViewConfig() {
	scans := db.GetScans()
	if len(scans) == 0 {
		helpers.PrintWarning("No scans found")
		return
	}

	for _, scan := range scans {
		fmt.Println("\033[1m" + scan.Name + "\033[0m")
		fmt.Println(" -> Type: " + scan.Type)
		fmt.Println(" -> Address: " + scan.Address)
		fmt.Println(" -> Timeout: " + helpers.IntToStr(scan.Timeout) + "ms")
		if scan.StatusCode != "" {
			fmt.Println(" -> Status code(s): " + scan.StatusCode)
		}
		if len(scan.Variables) > 0 {
			fmt.Println(" -> Variables:")
			for _, v := range scan.Variables {
				switch v.Mode {
				case "bool":
					fmt.Println("    - " + v.Name + " (bool): search=\"" + v.Value + "\"")
				case "json":
					fmt.Println("    - " + v.Name + " (json, as " + v.As + "): query=\"" + v.Query + "\"")
				}
			}
		}
		if scan.Expression != "" {
			fmt.Println(" -> Expression: \"" + scan.Expression + "\"")
		}
		fmt.Println()
	}
}

func validateScanName(scanName string, scanNames map[string]bool) string {
	if _, exists := scanNames[scanName]; exists {
		return DUPL_SCAN_NAME
	}
	if scanName == "" || len(scanName) > 32 {
		return INVALID_SCAN_NAME
	}
	for _, c := range scanName {
		if !('a' <= c && c <= 'z') && !('0' <= c && c <= '9') && c != '_' {
			return INVALID_SCAN_NAME
		}
	}
	return ""
}

func validateTimeout(scanTimeout string) string {
	num, correct := helpers.StrToInt(scanTimeout)
	if !correct {
		return TIMEOUT_NOT_INT
	}

	if num < 0 || num > 30000 {
		return TIMEOUT_OUT_OF_RANGE
	}

	return ""
}

func validateStatusCode(statusCode string) string {
	if len(statusCode) > 256 {
		return CODE_TOO_LONG
	}

	codes := strings.Split(statusCode, ",")
	if len(codes) == 0 {
		return INVALID_CODE
	}

	for _, code := range codes {
		if num, ok := helpers.StrToInt(code); !ok || num < 100 || num > 599 {
			return INVALID_CODE
		}
	}

	return ""
}

// validateVariables checks variable definitions for errors
func validateVariables(variables []helpers.Variable) string {
	varNames := make(map[string]bool)

	for _, v := range variables {
		if v.Name == "" {
			return "Variable name cannot be empty"
		}

		if !varNameRegex.MatchString(v.Name) {
			return "Variable name \"" + v.Name + "\" is invalid (must start with a letter, contain only lowercase letters, digits, and underscores)"
		}

		if varNames[v.Name] {
			return "Duplicate variable name: \"" + v.Name + "\""
		}
		varNames[v.Name] = true

		if v.Mode == "bool" && strings.TrimSpace(v.Value) == "" {
			return "Variable \"" + v.Name + "\" (bool) has an empty search value"
		}

		if v.Mode == "json" && strings.TrimSpace(v.Query) == "" {
			return "Variable \"" + v.Name + "\" (json) has an empty JSON query"
		}
	}

	return ""
}

// validateExpressionConfig validates the expression against defined variables
func validateExpressionConfig(variables []helpers.Variable, expression string) string {
	expression = strings.TrimSpace(expression)

	if expression == "" && len(variables) > 0 {
		return "Variables defined but no expression set"
	}

	if expression == "" {
		return ""
	}

	// Build defined variable names set and type map
	definedVars := make(map[string]bool)
	varTypes := make(map[string]string) // "bool" or "number"
	for _, v := range variables {
		definedVars[v.Name] = true
		switch v.Mode {
		case "bool":
			varTypes[v.Name] = "bool"
		case "json":
			if v.As == "number" {
				varTypes[v.Name] = "number"
			} else {
				varTypes[v.Name] = "bool"
			}
		}
	}

	// Tokenize
	tokens, err := tokenizeExpression(expression)
	if err != "" {
		return "Expression error: " + err
	}

	// Check balanced parentheses
	depth := 0
	for _, t := range tokens {
		if t == "(" {
			depth++
		}
		if t == ")" {
			depth--
		}
		if depth < 0 {
			return "Unbalanced parentheses: extra closing ')'"
		}
	}
	if depth > 0 {
		return "Unbalanced parentheses: missing closing ')'"
	}

	// Check identifiers reference defined variables and type correctness
	comparisonOps := map[string]bool{">": true, "<": true, ">=": true, "<=": true, "==": true, "!=": true}

	for i, t := range tokens {
		if varNameRegex.MatchString(t) {
			if !definedVars[t] {
				return "Undefined variable: \"" + t + "\""
			}

			prevToken := ""
			nextToken := ""
			if i > 0 {
				prevToken = tokens[i-1]
			}
			if i+1 < len(tokens) {
				nextToken = tokens[i+1]
			}

			switch varTypes[t] {
			case "bool":
				if comparisonOps[prevToken] || comparisonOps[nextToken] {
					return "Variable \"" + t + "\" is bool and cannot be used with comparison operators"
				}
			case "number":
				if !comparisonOps[prevToken] && !comparisonOps[nextToken] {
					return "Variable \"" + t + "\" is a number and must be used with a comparison operator (e.g. " + t + " > 0)"
				}
			}
		}
	}

	return ""
}

// tokenizeExpression tokenizes an expression string into tokens for validation.
// Returns tokens and an error message (empty string on success).
func tokenizeExpression(expr string) ([]string, string) {
	var tokens []string
	i := 0
	s := []rune(strings.TrimSpace(expr))

	for i < len(s) {
		// Whitespace
		if s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r' {
			i++
			continue
		}

		// Parentheses
		if s[i] == '(' || s[i] == ')' {
			tokens = append(tokens, string(s[i]))
			i++
			continue
		}

		// Two-character operators
		if i+1 < len(s) {
			two := string(s[i : i+2])
			if two == "&&" || two == "||" || two == ">=" || two == "<=" || two == "==" || two == "!=" {
				tokens = append(tokens, two)
				i += 2
				continue
			}
		}

		// Single-character operators
		if s[i] == '!' || s[i] == '>' || s[i] == '<' {
			tokens = append(tokens, string(s[i]))
			i++
			continue
		}

		// Numeric literals (including negative numbers and decimals)
		if (s[i] >= '0' && s[i] <= '9') || (s[i] == '-' && i+1 < len(s) && s[i+1] >= '0' && s[i+1] <= '9') {
			start := i
			if s[i] == '-' {
				i++
			}
			for i < len(s) && ((s[i] >= '0' && s[i] <= '9') || s[i] == '.') {
				i++
			}
			tokens = append(tokens, string(s[start:i]))
			continue
		}

		// Identifiers (variable names)
		if (s[i] >= 'a' && s[i] <= 'z') || s[i] == '_' {
			start := i
			for i < len(s) && ((s[i] >= 'a' && s[i] <= 'z') || (s[i] >= '0' && s[i] <= '9') || s[i] == '_') {
				i++
			}
			tokens = append(tokens, string(s[start:i]))
			continue
		}

		return nil, "Unexpected character '" + string(s[i]) + "' at position " + helpers.IntToStr(i+1)
	}

	return tokens, ""
}
