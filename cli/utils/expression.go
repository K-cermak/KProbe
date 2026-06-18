package utils

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"KProbeCLI/helpers"
)

// EvaluateVariables resolves all variables against the HTTP response body.
// Returns a map of variable name -> resolved value (bool or float64).
func EvaluateVariables(body string, variables []helpers.Variable) (map[string]interface{}, error) {
	values := make(map[string]interface{})

	for _, v := range variables {
		switch v.Mode {
		case "bool":
			values[v.Name] = strings.Contains(body, v.Value)

		case "json":
			var jsonData interface{}
			if err := json.Unmarshal([]byte(body), &jsonData); err != nil {
				return nil, fmt.Errorf("failed to parse response as JSON: %w", err)
			}

			result, found := resolveJSONPath(jsonData, v.Query)

			switch v.As {
			case "number":
				if !found {
					return nil, fmt.Errorf("JSON path %q not found for variable %q", v.Query, v.Name)
				}
				num, err := toFloat64(result)
				if err != nil {
					return nil, fmt.Errorf("variable %q: value at %q is not a number: %w", v.Name, v.Query, err)
				}
				values[v.Name] = num

			default: // "bool"
				values[v.Name] = found && isTruthy(result)
			}
		}
	}

	return values, nil
}

// resolveJSONPath walks a parsed JSON structure using a dot-notation path.
// Supports: $.key.nested, $.key[0], $.key[0].nested
func resolveJSONPath(data interface{}, path string) (interface{}, bool) {
	if !strings.HasPrefix(path, "$.") {
		return nil, false
	}

	path = path[2:] // strip "$."
	if path == "" {
		return data, true
	}

	parts := splitJSONPath(path)
	current := data

	for _, part := range parts {
		if part == "" {
			continue
		}

		// Check for array index: key[0]
		if idx := strings.Index(part, "["); idx != -1 {
			key := part[:idx]
			indexStr := part[idx+1 : len(part)-1]
			index, err := strconv.Atoi(indexStr)
			if err != nil {
				return nil, false
			}

			if key != "" {
				obj, ok := current.(map[string]interface{})
				if !ok {
					return nil, false
				}
				current, ok = obj[key]
				if !ok {
					return nil, false
				}
			}

			arr, ok := current.([]interface{})
			if !ok || index < 0 || index >= len(arr) {
				return nil, false
			}
			current = arr[index]
		} else {
			obj, ok := current.(map[string]interface{})
			if !ok {
				return nil, false
			}
			current, ok = obj[part]
			if !ok {
				return nil, false
			}
		}
	}

	return current, true
}

// splitJSONPath splits a JSON path like "data.items[0].value" into parts.
func splitJSONPath(path string) []string {
	var parts []string
	var current strings.Builder
	bracketDepth := 0

	for _, ch := range path {
		if ch == '[' {
			bracketDepth++
			current.WriteRune(ch)
		} else if ch == ']' {
			bracketDepth--
			current.WriteRune(ch)
		} else if ch == '.' && bracketDepth == 0 {
			parts = append(parts, current.String())
			current.Reset()
		} else {
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

func toFloat64(v interface{}) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	case json.Number:
		return val.Float64()
	case string:
		return strconv.ParseFloat(val, 64)
	case bool:
		if val {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %T to number", v)
	}
}

func isTruthy(v interface{}) bool {
	if v == nil {
		return false
	}
	switch val := v.(type) {
	case bool:
		return val
	case float64:
		return val != 0
	case string:
		return val != ""
	case []interface{}:
		return len(val) > 0
	case map[string]interface{}:
		return len(val) > 0
	default:
		return true
	}
}

// ===== EXPRESSION EVALUATION =====

// Token types for the expression parser
const (
	tokIdent  = iota // variable name
	tokNumber        // numeric literal
	tokAnd           // &&
	tokOr            // ||
	tokNot           // !
	tokLParen        // (
	tokRParen        // )
	tokGt            // >
	tokLt            // <
	tokGe            // >=
	tokLe            // <=
	tokEq            // ==
	tokNe            // !=
)

type token struct {
	typ int
	val string
}

// EvaluateExpression evaluates a logical expression against resolved variable values.
func EvaluateExpression(expr string, values map[string]interface{}) (bool, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return true, nil
	}

	tokens, err := tokenize(expr)
	if err != nil {
		return false, err
	}

	p := &parser{tokens: tokens, pos: 0, values: values}
	result, err := p.parseOr()
	if err != nil {
		return false, err
	}

	if p.pos < len(p.tokens) {
		return false, fmt.Errorf("unexpected token %q at position %d", p.tokens[p.pos].val, p.pos)
	}

	return result, nil
}

func tokenize(expr string) ([]token, error) {
	var tokens []token
	i := 0
	s := []rune(expr)

	for i < len(s) {
		// Skip whitespace
		if unicode.IsSpace(s[i]) {
			i++
			continue
		}

		// Parentheses
		if s[i] == '(' {
			tokens = append(tokens, token{tokLParen, "("})
			i++
			continue
		}
		if s[i] == ')' {
			tokens = append(tokens, token{tokRParen, ")"})
			i++
			continue
		}

		// Two-character operators
		if i+1 < len(s) {
			two := string(s[i : i+2])
			switch two {
			case "&&":
				tokens = append(tokens, token{tokAnd, "&&"})
				i += 2
				continue
			case "||":
				tokens = append(tokens, token{tokOr, "||"})
				i += 2
				continue
			case ">=":
				tokens = append(tokens, token{tokGe, ">="})
				i += 2
				continue
			case "<=":
				tokens = append(tokens, token{tokLe, "<="})
				i += 2
				continue
			case "==":
				tokens = append(tokens, token{tokEq, "=="})
				i += 2
				continue
			case "!=":
				tokens = append(tokens, token{tokNe, "!="})
				i += 2
				continue
			}
		}

		// Single-character operators
		if s[i] == '!' {
			tokens = append(tokens, token{tokNot, "!"})
			i++
			continue
		}
		if s[i] == '>' {
			tokens = append(tokens, token{tokGt, ">"})
			i++
			continue
		}
		if s[i] == '<' {
			tokens = append(tokens, token{tokLt, "<"})
			i++
			continue
		}

		// Numeric literals (including negative numbers and decimals)
		if unicode.IsDigit(s[i]) || (s[i] == '-' && i+1 < len(s) && unicode.IsDigit(s[i+1])) {
			start := i
			if s[i] == '-' {
				i++
			}
			for i < len(s) && (unicode.IsDigit(s[i]) || s[i] == '.') {
				i++
			}
			tokens = append(tokens, token{tokNumber, string(s[start:i])})
			continue
		}

		// Identifiers (variable names)
		if s[i] == '_' || unicode.IsLetter(s[i]) {
			start := i
			for i < len(s) && (unicode.IsLetter(s[i]) || unicode.IsDigit(s[i]) || s[i] == '_') {
				i++
			}
			tokens = append(tokens, token{tokIdent, string(s[start:i])})
			continue
		}

		return nil, fmt.Errorf("unexpected character %q at position %d", string(s[i]), i+1)
	}

	return tokens, nil
}

// Recursive descent parser
// Precedence (lowest to highest): OR, AND, comparison, unary NOT, primary
type parser struct {
	tokens []token
	pos    int
	values map[string]interface{}
}

func (p *parser) peek() *token {
	if p.pos < len(p.tokens) {
		return &p.tokens[p.pos]
	}
	return nil
}

func (p *parser) next() *token {
	t := p.peek()
	if t != nil {
		p.pos++
	}
	return t
}

// parseOr: expr || expr
func (p *parser) parseOr() (bool, error) {
	left, err := p.parseAnd()
	if err != nil {
		return false, err
	}

	for p.peek() != nil && p.peek().typ == tokOr {
		p.next()
		right, err := p.parseAnd()
		if err != nil {
			return false, err
		}
		left = left || right
	}

	return left, nil
}

// parseAnd: expr && expr
func (p *parser) parseAnd() (bool, error) {
	left, err := p.parseComparison()
	if err != nil {
		return false, err
	}

	for p.peek() != nil && p.peek().typ == tokAnd {
		p.next()
		right, err := p.parseComparison()
		if err != nil {
			return false, err
		}
		left = left && right
	}

	return left, nil
}

// parseComparison: value (> | < | >= | <= | == | !=) value
func (p *parser) parseComparison() (bool, error) {
	leftVal, leftIsBool, err := p.parseUnary()
	if err != nil {
		return false, err
	}

	t := p.peek()
	if t == nil {
		if leftIsBool {
			return leftVal.(bool), nil
		}
		return false, fmt.Errorf("numeric value used without comparison operator")
	}

	switch t.typ {
	case tokGt, tokLt, tokGe, tokLe, tokEq, tokNe:
		op := p.next()
		rightVal, _, err := p.parseUnary()
		if err != nil {
			return false, err
		}

		leftNum, errL := toFloat64(leftVal)
		rightNum, errR := toFloat64(rightVal)
		if errL != nil || errR != nil {
			return false, fmt.Errorf("comparison requires numeric values")
		}

		switch op.typ {
		case tokGt:
			return leftNum > rightNum, nil
		case tokLt:
			return leftNum < rightNum, nil
		case tokGe:
			return leftNum >= rightNum, nil
		case tokLe:
			return leftNum <= rightNum, nil
		case tokEq:
			return leftNum == rightNum, nil
		case tokNe:
			return leftNum != rightNum, nil
		}
	}

	if leftIsBool {
		return leftVal.(bool), nil
	}

	// A bare number that isn't compared — treat as truthy (non-zero)
	num, _ := toFloat64(leftVal)
	return num != 0, nil
}

// parseUnary: !expr or primary
func (p *parser) parseUnary() (interface{}, bool, error) {
	if p.peek() != nil && p.peek().typ == tokNot {
		p.next()
		val, isBool, err := p.parseUnary()
		if err != nil {
			return nil, false, err
		}
		if isBool {
			return !val.(bool), true, nil
		}
		num, _ := toFloat64(val)
		return num == 0, true, nil
	}
	return p.parsePrimary()
}

// parsePrimary: (expr) | identifier | number
func (p *parser) parsePrimary() (interface{}, bool, error) {
	t := p.peek()
	if t == nil {
		return nil, false, fmt.Errorf("unexpected end of expression")
	}

	switch t.typ {
	case tokLParen:
		p.next()
		result, err := p.parseOr()
		if err != nil {
			return nil, false, err
		}
		if p.peek() == nil || p.peek().typ != tokRParen {
			return nil, false, fmt.Errorf("expected closing parenthesis")
		}
		p.next()
		return result, true, nil

	case tokIdent:
		p.next()
		val, ok := p.values[t.val]
		if !ok {
			return nil, false, fmt.Errorf("undefined variable %q", t.val)
		}
		switch v := val.(type) {
		case bool:
			return v, true, nil
		case float64:
			return v, false, nil
		default:
			return val, false, nil
		}

	case tokNumber:
		p.next()
		num, err := strconv.ParseFloat(t.val, 64)
		if err != nil {
			return nil, false, fmt.Errorf("invalid number %q", t.val)
		}
		return num, false, nil

	default:
		return nil, false, fmt.Errorf("unexpected token %q", t.val)
	}
}
