package rules

import (
	"regexp"
)

type PatternMatching struct {
	BaseVersionedRule
}

func NewPatternMatching() *PatternMatching {
	return &PatternMatching{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp7, safe: true},
	}
}

func (r *PatternMatching) Name() string {
	return "pattern-matching"
}

func (r *PatternMatching) Description() string {
	return "Use pattern matching for type checks (C# 7+)"
}

func (r *PatternMatching) Apply(content string) (string, bool) {
	changed := false
	result := content

	pattern := regexp.MustCompile(`\(\s*(\w+)\s+as\s+(\w+)\s*\)\s*!=\s*null`)
	if pattern.MatchString(result) {
		result = pattern.ReplaceAllString(result, "${1} is ${2}")
		changed = true
	}

	return result, changed
}

type PatternMatchingNull struct {
	BaseVersionedRule
}

func NewPatternMatchingNull() *PatternMatchingNull {
	return &PatternMatchingNull{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp9, safe: true},
	}
}

func (r *PatternMatchingNull) Name() string {
	return "pattern-matching-null"
}

func (r *PatternMatchingNull) Description() string {
	return "Use 'is null' and 'is not null' patterns (C# 9+)"
}

func (r *PatternMatchingNull) Apply(content string) (string, bool) {
	changed := false
	result := content

	// Convert null checks in statement contexts (if, while, etc.)
	// Avoid converting inside lambdas (=>) as it may break Expression trees in LINQ

	// Helper function to check if we're in a lambda context
	isLambdaContext := func(line string, matchStart int) bool {
		// Check if there's a => before this match on the same line
		prefix := line[:matchStart]
		return regexp.MustCompile(`=>\s*`).MatchString(prefix)
	}

	// Process line by line to avoid cross-line issues and check lambda context
	lines := regexp.MustCompile(`(?m)^.*$`).FindAllString(result, -1)
	var processedLines []string

	for _, line := range lines {
		processedLine := line

		// Skip lines that look like they're in lambda/LINQ context
		if regexp.MustCompile(`=>`).MatchString(line) {
			processedLines = append(processedLines, processedLine)
			continue
		}

		// Pattern for identifier == null (with word boundary to catch field names like _redisDatabase)
		// Match: _identifier == null or identifier == null
		nullEqPattern := regexp.MustCompile(`(\b_?\w+)\s*==\s*null\b`)
		if nullEqPattern.MatchString(processedLine) {
			// Only apply in if/while/return/assignment contexts
			if regexp.MustCompile(`\b(if|while|return)\s*\(`).MatchString(processedLine) ||
				regexp.MustCompile(`[=,\(]\s*_?\w+\s*==\s*null`).MatchString(processedLine) {
				processedLine = nullEqPattern.ReplaceAllString(processedLine, "${1} is null")
				changed = true
			}
		}

		// Pattern for identifier != null
		nullNeqPattern := regexp.MustCompile(`(\b_?\w+)\s*!=\s*null\b`)
		if nullNeqPattern.MatchString(processedLine) {
			// Only apply in if/while/return/assignment contexts
			if regexp.MustCompile(`\b(if|while|return)\s*\(`).MatchString(processedLine) ||
				regexp.MustCompile(`[=,\(]\s*_?\w+\s*!=\s*null`).MatchString(processedLine) {
				processedLine = nullNeqPattern.ReplaceAllString(processedLine, "${1} is not null")
				changed = true
			}
		}

		processedLines = append(processedLines, processedLine)
	}

	if changed {
		result = ""
		for i, line := range processedLines {
			if i > 0 {
				result += "\n"
			}
			result += line
		}
	}

	// Handle ternary conditionals at assignment level (safe context)
	// x == null ? a : b - only when not preceded by =>
	pattern5 := regexp.MustCompile(`([=,\(]\s*)(\w+)\s*==\s*null\s*\?`)
	if pattern5.MatchString(result) {
		result = pattern5.ReplaceAllStringFunc(result, func(match string) string {
			// Skip if this looks like it's in a lambda context
			if regexp.MustCompile(`=>\s*`).MatchString(match) {
				return match
			}
			submatches := pattern5.FindStringSubmatch(match)
			if submatches == nil {
				return match
			}
			changed = true
			return submatches[1] + submatches[2] + " is null ?"
		})
	}

	_ = isLambdaContext // Silence unused warning for now

	return result, changed
}
