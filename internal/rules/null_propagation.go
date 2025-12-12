package rules

import (
	"regexp"
)

// NullPropagation converts null checks to null-conditional operators (C# 6+)
type NullPropagation struct {
	BaseVersionedRule
}

func NewNullPropagation() *NullPropagation {
	return &NullPropagation{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp6, safe: true},
	}
}

func (r *NullPropagation) Name() string {
	return "null-propagation"
}

func (r *NullPropagation) Description() string {
	return "Use null-conditional operators (?.) (C# 6+)"
}

func (r *NullPropagation) Apply(content string) (string, bool) {
	changed := false
	result := content

	// Pattern: x != null ? x.Property : null -> x?.Property
	pattern := regexp.MustCompile(`(\w+)\s*!=\s*null\s*\?\s*(\w+)\.(\w+)\s*:\s*null`)
	matches := pattern.FindAllStringSubmatch(result, -1)
	for _, m := range matches {
		if len(m) >= 4 && m[1] == m[2] {
			old := m[0]
			replacement := m[1] + "?." + m[3]
			result = regexp.MustCompile(regexp.QuoteMeta(old)).ReplaceAllString(result, replacement)
			changed = true
		}
	}

	return result, changed
}
