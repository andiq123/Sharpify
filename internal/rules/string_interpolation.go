package rules

import (
	"regexp"
	"strconv"
	"strings"
)

// StringInterpolation converts string.Format to interpolated strings (C# 6+)
type StringInterpolation struct {
	BaseVersionedRule
}

func NewStringInterpolation() *StringInterpolation {
	return &StringInterpolation{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp6, safe: true},
	}
}

func (r *StringInterpolation) Name() string {
	return "string-interpolation"
}

func (r *StringInterpolation) Description() string {
	return "Convert string.Format to interpolated strings (C# 6+)"
}

func (r *StringInterpolation) Apply(content string) (string, bool) {
	// Match string.Format("...", arg1, arg2, ...)
	pattern := regexp.MustCompile(`string\.Format\s*\(\s*"([^"]+)"\s*,\s*([^)]+)\)`)

	if !pattern.MatchString(content) {
		return content, false
	}

	changed := false
	result := pattern.ReplaceAllStringFunc(content, func(match string) string {
		submatches := pattern.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}

		formatStr := submatches[1]
		argsStr := submatches[2]

		// Parse arguments
		args := r.parseArgs(argsStr)

		// Replace {0}, {1}, etc. with actual arguments
		newStr := formatStr
		for i, arg := range args {
			placeholder := "{" + strconv.Itoa(i) + "}"
			if strings.Contains(newStr, placeholder) {
				newStr = strings.ReplaceAll(newStr, placeholder, "{"+strings.TrimSpace(arg)+"}")
				changed = true
			}
		}

		if !changed {
			return match
		}

		return "$\"" + newStr + "\""
	})

	return result, changed
}

func (r *StringInterpolation) parseArgs(argsStr string) []string {
	var args []string
	var current strings.Builder
	depth := 0

	for _, ch := range argsStr {
		switch ch {
		case '(', '[', '{':
			depth++
			current.WriteRune(ch)
		case ')', ']', '}':
			depth--
			current.WriteRune(ch)
		case ',':
			if depth == 0 {
				args = append(args, current.String())
				current.Reset()
			} else {
				current.WriteRune(ch)
			}
		default:
			current.WriteRune(ch)
		}
	}

	if current.Len() > 0 {
		args = append(args, current.String())
	}

	return args
}
