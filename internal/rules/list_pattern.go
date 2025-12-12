package rules

import (
	"regexp"
)

type ListPattern struct {
	BaseVersionedRule
}

func NewListPattern() *ListPattern {
	return &ListPattern{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp11, safe: false},
	}
}

func (r *ListPattern) Name() string {
	return "list-pattern"
}

func (r *ListPattern) Description() string {
	return "Convert empty collection checks to list patterns (C# 11+)"
}

func (r *ListPattern) Apply(content string) (string, bool) {
	changed := false
	result := content

	emptyPattern := regexp.MustCompile(`(\w+)\.(?:Count|Length)\s*==\s*0`)
	if emptyPattern.MatchString(result) {
		result = emptyPattern.ReplaceAllStringFunc(result, func(match string) string {
			submatches := emptyPattern.FindStringSubmatch(match)
			if submatches == nil {
				return match
			}
			changed = true
			return submatches[1] + " is []"
		})
	}

	notAnyPattern := regexp.MustCompile(`!(\w+)\.Any\(\)`)
	if notAnyPattern.MatchString(result) {
		result = notAnyPattern.ReplaceAllStringFunc(result, func(match string) string {
			submatches := notAnyPattern.FindStringSubmatch(match)
			if submatches == nil {
				return match
			}
			changed = true
			return submatches[1] + " is []"
		})
	}

	return result, changed
}
