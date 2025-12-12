package rules

import (
	"regexp"
)

type TupleDeconstruction struct {
	BaseVersionedRule
}

func NewTupleDeconstruction() *TupleDeconstruction {
	return &TupleDeconstruction{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp7, safe: true},
	}
}

func (r *TupleDeconstruction) Name() string {
	return "tuple-deconstruction"
}

func (r *TupleDeconstruction) Description() string {
	return "Use ValueTuple instead of Tuple<T1,T2> (C# 7+)"
}

func (r *TupleDeconstruction) Apply(content string) (string, bool) {
	changed := false
	result := content

	pattern2 := regexp.MustCompile(`new\s+Tuple<\w+,\s*\w+>\s*\(([^)]+)\)`)
	if pattern2.MatchString(result) {
		result = pattern2.ReplaceAllString(result, "(${1})")
		changed = true
	}

	pattern := regexp.MustCompile(`Tuple<(\w+),\s*(\w+)>`)
	if pattern.MatchString(result) {
		result = pattern.ReplaceAllString(result, "(${1}, ${2})")
		changed = true
	}

	return result, changed
}

type DiscardVariable struct {
	BaseVersionedRule
}

func NewDiscardVariable() *DiscardVariable {
	return &DiscardVariable{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp7, safe: false},
	}
}

func (r *DiscardVariable) Name() string {
	return "discard-variable"
}

func (r *DiscardVariable) Description() string {
	return "Use discard (_) for unused out parameters (C# 7+)"
}

func (r *DiscardVariable) Apply(content string) (string, bool) {
	changed := false
	result := content

	// Common pattern: TryParse with unused result variable
	// if (int.TryParse(s, out var unused)) - where 'unused' is never used
	// Convert to: if (int.TryParse(s, out _))

	// Find all 'out var X' patterns and check if X is used elsewhere
	outVarPattern := regexp.MustCompile(`out\s+var\s+(\w+)`)
	matches := outVarPattern.FindAllStringSubmatchIndex(result, -1)

	// Process in reverse to preserve indices
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		varName := result[match[2]:match[3]]

		// Skip if variable name already indicates it's intentionally unused
		if varName == "_" {
			continue
		}

		// Check if the variable is used elsewhere in the content (after the declaration)
		afterDecl := result[match[1]:]
		// Look for usage of the variable (as identifier, not in another out var)
		usagePattern := regexp.MustCompile(`\b` + regexp.QuoteMeta(varName) + `\b`)
		usageMatches := usagePattern.FindAllStringIndex(afterDecl, -1)

		// If no usages found after declaration, it's unused
		if len(usageMatches) == 0 {
			// Replace 'out var X' with 'out _'
			result = result[:match[0]] + "out _" + result[match[1]:]
			changed = true
		}
	}

	return result, changed
}
