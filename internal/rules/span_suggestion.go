package rules

import (
	"regexp"
)

type SpanSuggestion struct {
	BaseVersionedRule
}

func NewSpanSuggestion() *SpanSuggestion {
	return &SpanSuggestion{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp7, safe: false},
	}
}

func (r *SpanSuggestion) Name() string {
	return "span-suggestion"
}

func (r *SpanSuggestion) Description() string {
	return "Optimize string operations for better performance (C# 7.2+)"
}

func (r *SpanSuggestion) Apply(content string) (string, bool) {
	changed := false
	result := content

	startsWithPattern := regexp.MustCompile(`(\w+)\.Substring\s*\(\s*0\s*,\s*(\w+)\.Length\s*\)\s*==\s*(\w+)`)
	if startsWithPattern.MatchString(result) {
		result = startsWithPattern.ReplaceAllStringFunc(result, func(match string) string {
			submatches := startsWithPattern.FindStringSubmatch(match)
			if submatches == nil {
				return match
			}
			if submatches[2] != submatches[3] {
				return match
			}
			changed = true
			return submatches[1] + ".StartsWith(" + submatches[2] + ")"
		})
	}

	endsWithPattern := regexp.MustCompile(`(\w+)\.Substring\s*\(\s*(\w+)\.Length\s*-\s*(\w+)\.Length\s*\)\s*==\s*(\w+)`)
	if endsWithPattern.MatchString(result) {
		result = endsWithPattern.ReplaceAllStringFunc(result, func(match string) string {
			submatches := endsWithPattern.FindStringSubmatch(match)
			if submatches == nil {
				return match
			}
			if submatches[1] != submatches[2] || submatches[3] != submatches[4] {
				return match
			}
			changed = true
			return submatches[1] + ".EndsWith(" + submatches[3] + ")"
		})
	}

	toCharArrayPattern := regexp.MustCompile(`foreach\s*\(\s*char\s+(\w+)\s+in\s+(\w+)\.ToCharArray\(\)\s*\)`)
	if toCharArrayPattern.MatchString(result) {
		result = toCharArrayPattern.ReplaceAllStringFunc(result, func(match string) string {
			submatches := toCharArrayPattern.FindStringSubmatch(match)
			if submatches == nil {
				return match
			}
			changed = true
			return "foreach (char " + submatches[1] + " in " + submatches[2] + ")"
		})
	}

	unnecessaryConversion := regexp.MustCompile(`new\s+string\s*\(\s*(\w+)\.ToCharArray\(\)\s*\)`)
	if unnecessaryConversion.MatchString(result) {
		result = unnecessaryConversion.ReplaceAllStringFunc(result, func(match string) string {
			submatches := unnecessaryConversion.FindStringSubmatch(match)
			if submatches == nil {
				return match
			}
			changed = true
			return submatches[1]
		})
	}

	toLowerCompare := regexp.MustCompile(`(\w+)\.ToLower\(\)\s*==\s*(\w+)\.ToLower\(\)`)
	if toLowerCompare.MatchString(result) {
		result = toLowerCompare.ReplaceAllStringFunc(result, func(match string) string {
			submatches := toLowerCompare.FindStringSubmatch(match)
			if submatches == nil {
				return match
			}
			changed = true
			return "string.Equals(" + submatches[1] + ", " + submatches[2] + ", StringComparison.OrdinalIgnoreCase)"
		})
	}

	toUpperCompare := regexp.MustCompile(`(\w+)\.ToUpper\(\)\s*==\s*(\w+)\.ToUpper\(\)`)
	if toUpperCompare.MatchString(result) {
		result = toUpperCompare.ReplaceAllStringFunc(result, func(match string) string {
			submatches := toUpperCompare.FindStringSubmatch(match)
			if submatches == nil {
				return match
			}
			changed = true
			return "string.Equals(" + submatches[1] + ", " + submatches[2] + ", StringComparison.OrdinalIgnoreCase)"
		})
	}

	return result, changed
}
