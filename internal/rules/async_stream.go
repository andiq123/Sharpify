package rules

import (
	"regexp"
)

type IndexRange struct {
	BaseVersionedRule
}

func NewIndexRange() *IndexRange {
	return &IndexRange{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp8, safe: true},
	}
}

func (r *IndexRange) Name() string {
	return "index-range"
}

func (r *IndexRange) Description() string {
	return "Use Index (^) operator for last element access (C# 8+)"
}

func (r *IndexRange) Apply(content string) (string, bool) {
	changed := false
	result := content

	pattern := regexp.MustCompile(`(\w+)\s*\[\s*(\w+)\.Length\s*-\s*(\d+)\s*\]`)
	matches := pattern.FindAllStringSubmatch(result, -1)
	for _, m := range matches {
		if len(m) >= 4 && m[1] == m[2] {
			old := m[0]
			replacement := m[1] + "[^" + m[3] + "]"
			result = regexp.MustCompile(regexp.QuoteMeta(old)).ReplaceAllString(result, replacement)
			changed = true
		}
	}

	return result, changed
}
