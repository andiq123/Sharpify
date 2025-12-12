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
	return "Convert array/list element checks to list patterns (C# 11+)"
}

func (r *ListPattern) Apply(content string) (string, bool) {
	changed := false
	result := content

	// Only convert patterns that improve readability:
	// - Empty check: list.Count == 0 -> list is [] (clearer intent)
	// - !list.Any() -> list is [] (clearer intent)
	//
	// NOT converting (less readable):
	// - list.Count > 0 -> list is [_, ..] (the original is clearer)
	// - list.Any() -> list is [_, ..] (Any() is more idiomatic)
	// - list[0] == x -> list is [x, ..] (unfamiliar syntax for many)

	// Pattern: Check if empty - list.Count == 0 or list.Length == 0
	// Before: if (list.Count == 0) or if (list.Length == 0)
	// After:  if (list is [])
	// This is MORE readable - clearly shows "empty list" intent
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

	// Pattern: !list.Any() -> list is []
	// This is MORE readable - clearly shows "empty list" intent
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
