package rules

import (
	"regexp"
)

type PatternMatching struct {
	BaseVersionedRule
}

func NewPatternMatching() *PatternMatching {
	return &PatternMatching{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp7, safe: true},
	}
}

func (r *PatternMatching) Name() string {
	return "pattern-matching"
}

func (r *PatternMatching) Description() string {
	return "Use pattern matching for type checks (C# 7+)"
}

func (r *PatternMatching) Apply(content string) (string, bool) {
	changed := false
	result := content

	pattern := regexp.MustCompile(`\(\s*(\w+)\s+as\s+(\w+)\s*\)\s*!=\s*null`)
	if pattern.MatchString(result) {
		result = pattern.ReplaceAllString(result, "${1} is ${2}")
		changed = true
	}

	return result, changed
}

type PatternMatchingNull struct {
	BaseVersionedRule
}

func NewPatternMatchingNull() *PatternMatchingNull {
	return &PatternMatchingNull{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp9, safe: true},
	}
}

func (r *PatternMatchingNull) Name() string {
	return "pattern-matching-null"
}

func (r *PatternMatchingNull) Description() string {
	return "Use 'is null' and 'is not null' patterns (C# 9+)"
}

func (r *PatternMatchingNull) Apply(content string) (string, bool) {
	changed := false
	result := content

	// Only convert null checks in simple statement contexts (if, while, etc.)
	// Avoid converting inside lambdas (=>) as it may break Expression trees in LINQ

	// Pattern for: if (x == null) or while (x != null) - statement level only
	// Use word boundary and statement context to be safe
	pattern1 := regexp.MustCompile(`(\bif\s*\(\s*)(\w+)\s*==\s*null(\s*\))`)
	if pattern1.MatchString(result) {
		result = pattern1.ReplaceAllString(result, "${1}${2} is null${3}")
		changed = true
	}

	pattern2 := regexp.MustCompile(`(\bif\s*\(\s*)(\w+)\s*!=\s*null(\s*\))`)
	if pattern2.MatchString(result) {
		result = pattern2.ReplaceAllString(result, "${1}${2} is not null${3}")
		changed = true
	}

	// Also handle while statements
	pattern3 := regexp.MustCompile(`(\bwhile\s*\(\s*)(\w+)\s*==\s*null(\s*\))`)
	if pattern3.MatchString(result) {
		result = pattern3.ReplaceAllString(result, "${1}${2} is null${3}")
		changed = true
	}

	pattern4 := regexp.MustCompile(`(\bwhile\s*\(\s*)(\w+)\s*!=\s*null(\s*\))`)
	if pattern4.MatchString(result) {
		result = pattern4.ReplaceAllString(result, "${1}${2} is not null${3}")
		changed = true
	}

	// Handle ternary conditionals at assignment level (safe context)
	// x == null ? a : b - only when not preceded by =>
	pattern5 := regexp.MustCompile(`([=,\(]\s*)(\w+)\s*==\s*null\s*\?`)
	if pattern5.MatchString(result) {
		// Check we're not in a lambda context
		result = pattern5.ReplaceAllStringFunc(result, func(match string) string {
			// Skip if this looks like it's in a lambda context
			if regexp.MustCompile(`=>\s*$`).MatchString(match) {
				return match
			}
			submatches := pattern5.FindStringSubmatch(match)
			if submatches == nil {
				return match
			}
			changed = true
			return submatches[1] + submatches[2] + " is null ?"
		})
	}

	return result, changed
}
