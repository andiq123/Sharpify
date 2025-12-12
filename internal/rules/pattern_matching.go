package rules

import (
	"regexp"
)

// PatternMatching converts type checks to pattern matching (C# 7+)
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

	// Pattern: (x as Type) != null -> x is Type
	pattern := regexp.MustCompile(`\(\s*(\w+)\s+as\s+(\w+)\s*\)\s*!=\s*null`)
	if pattern.MatchString(result) {
		result = pattern.ReplaceAllString(result, "${1} is ${2}")
		changed = true
	}

	return result, changed
}

// PatternMatchingNull converts null checks to is null / is not null (C# 9+)
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

	// Pattern: x == null -> x is null
	pattern1 := regexp.MustCompile(`(\w+)\s*==\s*null`)
	if pattern1.MatchString(result) {
		result = pattern1.ReplaceAllString(result, "${1} is null")
		changed = true
	}

	// Pattern: x != null -> x is not null
	pattern2 := regexp.MustCompile(`(\w+)\s*!=\s*null`)
	if pattern2.MatchString(result) {
		result = pattern2.ReplaceAllString(result, "${1} is not null")
		changed = true
	}

	return result, changed
}
