package rules

import (
	"regexp"
	"strings"
)

// ThrowHelper converts null checks to modern throw helper methods
type ThrowHelper struct {
	BaseVersionedRule
}

func NewThrowHelper() *ThrowHelper {
	return &ThrowHelper{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp10, safe: true},
	}
}

func (r *ThrowHelper) Name() string {
	return "throw-helper"
}

func (r *ThrowHelper) Description() string {
	return "Use ArgumentNullException.ThrowIfNull (C# 10+)"
}

func (r *ThrowHelper) Apply(content string) (string, bool) {
	changed := false
	result := content

	// Pattern: if (x is null) throw new ArgumentNullException(nameof(x));
	// Go regex doesn't support backreferences, so we match and verify manually
	pattern1 := regexp.MustCompile(`(?m)^\s*if\s*\(\s*(\w+)\s+is\s+null\s*\)\s*\{?\s*\n?\s*throw\s+new\s+ArgumentNullException\s*\(\s*nameof\s*\(\s*(\w+)\s*\)\s*\)\s*;\s*\}?`)
	matches := pattern1.FindAllStringSubmatchIndex(result, -1)

	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		varName1 := result[match[2]:match[3]]
		varName2 := result[match[4]:match[5]]

		// Verify both variable names are the same
		if varName1 == varName2 {
			replacement := "        ArgumentNullException.ThrowIfNull(" + varName1 + ");"
			result = result[:match[0]] + replacement + result[match[1]:]
			changed = true
		}
	}

	// Pattern 2: if (x == null) throw new ArgumentNullException(nameof(x)); (standalone)
	pattern2 := regexp.MustCompile(`(?m)^\s*if\s*\(\s*(\w+)\s*==\s*null\s*\)\s*\{?\s*\n?\s*throw\s+new\s+ArgumentNullException\s*\(\s*nameof\s*\(\s*(\w+)\s*\)\s*\)\s*;\s*\}?`)
	matches2 := pattern2.FindAllStringSubmatchIndex(result, -1)

	for i := len(matches2) - 1; i >= 0; i-- {
		match := matches2[i]
		varName1 := result[match[2]:match[3]]
		varName2 := result[match[4]:match[5]]
		matchEnd := match[1]

		// Verify both variable names are the same
		if varName1 != varName2 {
			continue
		}

		// Check if there's an assignment on the next line (to avoid conflict with throw-expression)
		restOfContent := result[matchEnd:]
		assignmentPattern := regexp.MustCompile(`^\s*\n\s*\w+\s*=\s*` + regexp.QuoteMeta(varName1))
		if !assignmentPattern.MatchString(restOfContent) {
			replacement := "        ArgumentNullException.ThrowIfNull(" + varName1 + ");"
			result = result[:match[0]] + replacement + result[match[1]:]
			changed = true
		}
	}

	// Pattern 3: if (string.IsNullOrEmpty(x)) throw new ArgumentException(..., nameof(x));
	pattern3 := regexp.MustCompile(`(?m)^\s*if\s*\(\s*string\.IsNullOrEmpty\s*\(\s*(\w+)\s*\)\s*\)\s*\{?\s*\n?\s*throw\s+new\s+Argument(?:Null)?Exception\s*\([^)]*nameof\s*\(\s*(\w+)\s*\)[^)]*\)\s*;\s*\}?`)
	matches3 := pattern3.FindAllStringSubmatchIndex(result, -1)

	for i := len(matches3) - 1; i >= 0; i-- {
		match := matches3[i]
		varName1 := result[match[2]:match[3]]
		varName2 := result[match[4]:match[5]]

		if varName1 == varName2 {
			replacement := "        ArgumentException.ThrowIfNullOrEmpty(" + varName1 + ");"
			result = result[:match[0]] + replacement + result[match[1]:]
			changed = true
		}
	}

	_ = strings.TrimSpace // Silence unused import warning

	return result, changed
}
