package rules

import (
	"regexp"
	"strings"
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

	// Pattern 1: Check first element - list[0] == value or list[0].Equals(value)
	// Before: if (list[0] == "value")
	// After:  if (list is ["value", ..])
	firstElemPattern := regexp.MustCompile(`(\w+)\[0\]\s*==\s*([^&|)]+)`)
	if firstElemPattern.MatchString(result) {
		result = firstElemPattern.ReplaceAllStringFunc(result, func(match string) string {
			submatches := firstElemPattern.FindStringSubmatch(match)
			if submatches == nil {
				return match
			}
			listName := submatches[1]
			value := strings.TrimSpace(submatches[2])
			changed = true
			return listName + " is [" + value + ", ..]"
		})
	}

	// Pattern 2a: Check last element with ^1 index - list[^1] == value
	// Before: if (list[^1] == "value")
	// After:  if (list is [.., "value"])
	lastElemIndexPattern := regexp.MustCompile(`(\w+)\[\s*\^1\s*\]\s*==\s*([^&|)]+)`)
	if lastElemIndexPattern.MatchString(result) {
		result = lastElemIndexPattern.ReplaceAllStringFunc(result, func(match string) string {
			submatches := lastElemIndexPattern.FindStringSubmatch(match)
			if submatches == nil {
				return match
			}
			listName := submatches[1]
			value := strings.TrimSpace(submatches[2])
			changed = true
			return listName + " is [.., " + value + "]"
		})
	}

	// Pattern 2b: Check last element with Count/Length - 1 - list[list.Count - 1] == value
	// Before: if (list[list.Count - 1] == "value")
	// After:  if (list is [.., "value"])
	lastElemCountPattern := regexp.MustCompile(`(\w+)\[\s*(\w+)\.(?:Count|Length)\s*-\s*1\s*\]\s*==\s*([^&|)]+)`)
	if lastElemCountPattern.MatchString(result) {
		result = lastElemCountPattern.ReplaceAllStringFunc(result, func(match string) string {
			submatches := lastElemCountPattern.FindStringSubmatch(match)
			if submatches == nil {
				return match
			}
			listName := submatches[1]
			countVar := submatches[2]
			value := strings.TrimSpace(submatches[3])
			// Only convert if the list name matches the Count/Length variable
			if listName != countVar {
				return match
			}
			changed = true
			return listName + " is [.., " + value + "]"
		})
	}

	// Pattern 3: Check both first and last
	// Before: if (list[0] == first && list[^1] == last)
	// After:  if (list is [var f, .., var l] && f == first && l == last)
	// This is complex and requires careful handling - skip for now as it needs AST

	// Pattern 4: Check if empty - list.Count == 0 or list.Length == 0
	// Before: if (list.Count == 0) or if (list.Length == 0)
	// After:  if (list is [])
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

	// Pattern 5: Check if not empty - list.Count > 0 or list.Length > 0 or list.Any()
	// Before: if (list.Count > 0) or if (list.Any())
	// After:  if (list is [_, ..])
	notEmptyPattern := regexp.MustCompile(`(\w+)\.(?:Count|Length)\s*>\s*0`)
	if notEmptyPattern.MatchString(result) {
		result = notEmptyPattern.ReplaceAllStringFunc(result, func(match string) string {
			submatches := notEmptyPattern.FindStringSubmatch(match)
			if submatches == nil {
				return match
			}
			changed = true
			return submatches[1] + " is [_, ..]"
		})
	}

	// Pattern 6: list.Any() -> list is [_, ..]
	anyPattern := regexp.MustCompile(`(\w+)\.Any\(\)`)
	if anyPattern.MatchString(result) {
		result = anyPattern.ReplaceAllStringFunc(result, func(match string) string {
			submatches := anyPattern.FindStringSubmatch(match)
			if submatches == nil {
				return match
			}
			changed = true
			return submatches[1] + " is [_, ..]"
		})
	}

	// Pattern 7: !list.Any() -> list is []
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

	// Pattern 8: Check single element - list.Count == 1
	// Before: if (list.Count == 1)
	// After:  if (list is [_])
	singlePattern := regexp.MustCompile(`(\w+)\.(?:Count|Length)\s*==\s*1`)
	if singlePattern.MatchString(result) {
		result = singlePattern.ReplaceAllStringFunc(result, func(match string) string {
			submatches := singlePattern.FindStringSubmatch(match)
			if submatches == nil {
				return match
			}
			changed = true
			return submatches[1] + " is [_]"
		})
	}

	return result, changed
}
