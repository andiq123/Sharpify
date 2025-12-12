package rules

import (
	"regexp"
	"strings"
)

type DefaultLiteral struct {
	BaseVersionedRule
}

func NewDefaultLiteral() *DefaultLiteral {
	return &DefaultLiteral{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp7, safe: true},
	}
}

func (r *DefaultLiteral) Name() string {
	return "default-literal"
}

func (r *DefaultLiteral) Description() string {
	return "Use default literal instead of default(T) (C# 7.1+)"
}

func (r *DefaultLiteral) Apply(content string) (string, bool) {
	changed := false
	result := content

	pattern1 := regexp.MustCompile(`(\w+(?:<[^>]+>)?)\s+(\w+)\s*=\s*default\s*\(\s*(\w+(?:<[^>]+>)?)\s*\)\s*;`)
	matches := pattern1.FindAllStringSubmatch(result, -1)
	for _, m := range matches {
		if len(m) >= 4 && strings.TrimSpace(m[1]) == strings.TrimSpace(m[3]) {
			old := m[0]
			replacement := m[1] + " " + m[2] + " = default;"
			result = strings.Replace(result, old, replacement, 1)
			changed = true
		}
	}

	pattern2 := regexp.MustCompile(`return\s+default\s*\(\s*\w+(?:<[^>]+>)?\s*\)\s*;`)
	if pattern2.MatchString(result) {
		result = pattern2.ReplaceAllString(result, "return default;")
		changed = true
	}

	return result, changed
}
