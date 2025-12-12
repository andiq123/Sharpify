package rules

import (
	"regexp"
)

// NameofExpression converts string literals to nameof expressions (C# 6+)
type NameofExpression struct {
	BaseVersionedRule
}

func NewNameofExpression() *NameofExpression {
	return &NameofExpression{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp6, safe: true},
	}
}

func (r *NameofExpression) Name() string {
	return "nameof-expression"
}

func (r *NameofExpression) Description() string {
	return "Use nameof() for ArgumentNullException and similar (C# 6+)"
}

func (r *NameofExpression) Apply(content string) (string, bool) {
	changed := false
	result := content

	// Pattern: throw new ArgumentNullException("paramName") -> throw new ArgumentNullException(nameof(paramName))
	pattern1 := regexp.MustCompile(`throw\s+new\s+ArgumentNullException\s*\(\s*"(\w+)"\s*\)`)
	if pattern1.MatchString(result) {
		result = pattern1.ReplaceAllString(result, "throw new ArgumentNullException(nameof(${1}))")
		changed = true
	}

	// Pattern: throw new ArgumentException("message", "paramName") -> throw new ArgumentException("message", nameof(paramName))
	pattern2 := regexp.MustCompile(`throw\s+new\s+ArgumentException\s*\(\s*("[^"]*")\s*,\s*"(\w+)"\s*\)`)
	if pattern2.MatchString(result) {
		result = pattern2.ReplaceAllString(result, "throw new ArgumentException(${1}, nameof(${2}))")
		changed = true
	}

	// Pattern: throw new ArgumentOutOfRangeException("paramName") -> throw new ArgumentOutOfRangeException(nameof(paramName))
	pattern3 := regexp.MustCompile(`throw\s+new\s+ArgumentOutOfRangeException\s*\(\s*"(\w+)"\s*\)`)
	if pattern3.MatchString(result) {
		result = pattern3.ReplaceAllString(result, "throw new ArgumentOutOfRangeException(nameof(${1}))")
		changed = true
	}

	return result, changed
}
