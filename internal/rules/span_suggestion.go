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
	return "Convert Substring operations to Span/AsSpan for better performance (C# 7.2+)"
}

func (r *SpanSuggestion) Apply(content string) (string, bool) {
	changed := false
	result := content

	// Pattern 1: str.Substring(start, length) -> str.AsSpan(start, length) or str.AsSpan().Slice(start, length)
	// Only safe when used in comparisons or passed to Span-accepting methods
	// For now, convert .Substring() to .AsSpan() when used with .SequenceEqual or similar

	// Pattern 2: str.Substring(start) for comparison -> str.AsSpan(start)
	// Before: if (str.Substring(5) == "test")
	// After:  if (str.AsSpan(5).SequenceEqual("test"))
	// This is risky, so we'll be conservative

	// Pattern 3: string.Concat with Substring -> use Span
	// This requires significant refactoring, skip for safety

	// Pattern 4: Convert StartsWith/EndsWith with Substring to direct check
	// Before: str.Substring(0, prefix.Length) == prefix
	// After:  str.StartsWith(prefix)
	startsWithPattern := regexp.MustCompile(`(\w+)\.Substring\s*\(\s*0\s*,\s*(\w+)\.Length\s*\)\s*==\s*(\w+)`)
	if startsWithPattern.MatchString(result) {
		result = startsWithPattern.ReplaceAllStringFunc(result, func(match string) string {
			submatches := startsWithPattern.FindStringSubmatch(match)
			if submatches == nil {
				return match
			}
			// Verify that the length variable matches the comparison variable
			if submatches[2] != submatches[3] {
				return match
			}
			changed = true
			return submatches[1] + ".StartsWith(" + submatches[2] + ")"
		})
	}

	// Pattern 5: Substring for EndsWith
	// Before: str.Substring(str.Length - suffix.Length) == suffix
	// After:  str.EndsWith(suffix)
	endsWithPattern := regexp.MustCompile(`(\w+)\.Substring\s*\(\s*(\w+)\.Length\s*-\s*(\w+)\.Length\s*\)\s*==\s*(\w+)`)
	if endsWithPattern.MatchString(result) {
		result = endsWithPattern.ReplaceAllStringFunc(result, func(match string) string {
			submatches := endsWithPattern.FindStringSubmatch(match)
			if submatches == nil {
				return match
			}
			// Verify str.Length matches str and suffix.Length matches suffix
			if submatches[1] != submatches[2] || submatches[3] != submatches[4] {
				return match
			}
			changed = true
			return submatches[1] + ".EndsWith(" + submatches[3] + ")"
		})
	}

	// Pattern 6: ToCharArray() in foreach -> iterate directly or use span
	// Before: foreach (char c in str.ToCharArray())
	// After:  foreach (char c in str)
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

	// Pattern 7: new string(charArray) from ToCharArray -> unnecessary
	// Before: new string(str.ToCharArray())
	// After:  str (or just remove the conversion)
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

	// Pattern 8: string.IsNullOrEmpty -> use pattern matching (C# 9+)
	// This is more of a pattern matching thing, skip here

	// Pattern 9: str.ToLower() == str.ToLower() comparisons -> use StringComparison
	// Before: str1.ToLower() == str2.ToLower()
	// After:  string.Equals(str1, str2, StringComparison.OrdinalIgnoreCase)
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

	// Pattern 10: str.ToUpper() == str.ToUpper() comparisons
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
