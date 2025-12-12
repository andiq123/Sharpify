package rules

import (
	"regexp"
	"strings"
)

// StringConcatToInterpolation converts simple string concatenation to interpolation
type StringConcatToInterpolation struct {
	BaseVersionedRule
}

func NewStringConcatToInterpolation() *StringConcatToInterpolation {
	return &StringConcatToInterpolation{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp6, safe: false}, // opt-in as it changes style
	}
}

func (r *StringConcatToInterpolation) Name() string {
	return "string-concat-interpolation"
}

func (r *StringConcatToInterpolation) Description() string {
	return "Convert string concatenation to interpolation (C# 6+)"
}

func (r *StringConcatToInterpolation) Apply(content string) (string, bool) {
	changed := false

	// Pattern: "literal" + variable or variable + "literal"
	// Match patterns like: "prefix" + variable or variable + "suffix"

	// Pattern 1: "literal" + identifier
	pattern1 := regexp.MustCompile(`"([^"]*?)"\s*\+\s*([a-zA-Z_][a-zA-Z0-9_]*(?:\.[a-zA-Z_][a-zA-Z0-9_]*)*)`)

	// Pattern 2: identifier + "literal"
	pattern2 := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*(?:\.[a-zA-Z_][a-zA-Z0-9_]*)*)\s*\+\s*"([^"]*?)"`)

	// First pass: Convert "literal" + var patterns
	result := pattern1.ReplaceAllStringFunc(content, func(match string) string {
		submatches := pattern1.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}
		literal := submatches[1]
		variable := submatches[2]

		// Skip if this looks like it's already inside an interpolated string
		// or if it's a simple case that might be intentional
		if strings.HasPrefix(literal, "{") || strings.HasSuffix(literal, "}") {
			return match
		}

		changed = true
		return `$"` + literal + `{` + variable + `}"`
	})

	// Second pass: Handle var + "literal" that aren't already converted
	// This is trickier because we need to check if we're not breaking existing interpolation
	result = pattern2.ReplaceAllStringFunc(result, func(match string) string {
		submatches := pattern2.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}
		variable := submatches[1]
		literal := submatches[2]

		// Skip common false positives
		if variable == "string" || variable == "String" {
			return match
		}

		// Skip if literal starts with $ (already interpolated somehow)
		if strings.Contains(match, `$"`) {
			return match
		}

		changed = true
		return `$"{` + variable + `}` + literal + `"`
	})

	return result, changed
}
