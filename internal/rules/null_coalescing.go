package rules

import (
	"regexp"
)

// NullCoalescing converts null checks to null-coalescing operators
type NullCoalescing struct {
	BaseVersionedRule
}

func NewNullCoalescing() *NullCoalescing {
	return &NullCoalescing{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp8, safe: true},
	}
}

func (r *NullCoalescing) Name() string {
	return "null-coalescing-assignment"
}

func (r *NullCoalescing) Description() string {
	return "Use null-coalescing assignment operator (??=) (C# 8+)"
}

func (r *NullCoalescing) Apply(content string) (string, bool) {
	changed := false
	result := content

	// Pattern: if (x == null) x = value; -> x ??= value;
	pattern1 := regexp.MustCompile(`if\s*\(\s*(\w+)\s*==\s*null\s*\)\s*(\w+)\s*=\s*([^;]+);`)
	matches := pattern1.FindAllStringSubmatch(result, -1)
	for _, m := range matches {
		if len(m) >= 4 && m[1] == m[2] {
			old := m[0]
			replacement := m[1] + " ??= " + m[3] + ";"
			result = regexp.MustCompile(regexp.QuoteMeta(old)).ReplaceAllString(result, replacement)
			changed = true
		}
	}

	return result, changed
}
