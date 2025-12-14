package rules

import (
	"regexp"
)

// StringIsNullOrEmpty converts string.IsNullOrEmpty to pattern matching
type StringIsNullOrEmpty struct {
	BaseVersionedRule
}

func NewStringIsNullOrEmpty() *StringIsNullOrEmpty {
	return &StringIsNullOrEmpty{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp11, safe: true},
	}
}

func (r *StringIsNullOrEmpty) Name() string {
	return "string-isnullorempty"
}

func (r *StringIsNullOrEmpty) Description() string {
	return "Use pattern matching for null/empty string checks (C# 11+)"
}

func (r *StringIsNullOrEmpty) Apply(content string) (string, bool) {
	changed := false
	result := content

	// Pattern 1: string.IsNullOrEmpty(x) -> x is null or ""
	pattern1 := regexp.MustCompile(`string\.IsNullOrEmpty\s*\(\s*(\w+)\s*\)`)
	if pattern1.MatchString(result) {
		result = pattern1.ReplaceAllString(result, "${1} is null or \"\"")
		changed = true
	}

	// Pattern 2: !string.IsNullOrEmpty(x) -> x is not (null or "")
	// Handle this BEFORE we convert pattern1, check for the negation
	pattern2 := regexp.MustCompile(`!\s*string\.IsNullOrEmpty\s*\(\s*(\w+)\s*\)`)
	if pattern2.MatchString(result) {
		result = pattern2.ReplaceAllString(result, "${1} is not (null or \"\")")
		changed = true
	}

	// Pattern 3: String.IsNullOrEmpty (capital S)
	pattern3 := regexp.MustCompile(`String\.IsNullOrEmpty\s*\(\s*(\w+)\s*\)`)
	if pattern3.MatchString(result) {
		result = pattern3.ReplaceAllString(result, "${1} is null or \"\"")
		changed = true
	}

	// Pattern 4: !String.IsNullOrEmpty
	pattern4 := regexp.MustCompile(`!\s*String\.IsNullOrEmpty\s*\(\s*(\w+)\s*\)`)
	if pattern4.MatchString(result) {
		result = pattern4.ReplaceAllString(result, "${1} is not (null or \"\")")
		changed = true
	}

	return result, changed
}
