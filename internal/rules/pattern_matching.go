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

	pattern1 := regexp.MustCompile(`(\w+)\s*==\s*null`)
	if pattern1.MatchString(result) {
		result = pattern1.ReplaceAllString(result, "${1} is null")
		changed = true
	}

	pattern2 := regexp.MustCompile(`(\w+)\s*!=\s*null`)
	if pattern2.MatchString(result) {
		result = pattern2.ReplaceAllString(result, "${1} is not null")
		changed = true
	}

	return result, changed
}
