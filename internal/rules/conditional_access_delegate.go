package rules

import (
	"regexp"
)

// ConditionalAccessDelegate converts null-check delegate invocation to conditional access
type ConditionalAccessDelegate struct {
	BaseVersionedRule
}

func NewConditionalAccessDelegate() *ConditionalAccessDelegate {
	return &ConditionalAccessDelegate{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp6, safe: true},
	}
}

func (r *ConditionalAccessDelegate) Name() string {
	return "conditional-access-delegate"
}

func (r *ConditionalAccessDelegate) Description() string {
	return "Use conditional access for delegate invocation (C# 6+)"
}

func (r *ConditionalAccessDelegate) Apply(content string) (string, bool) {
	changed := false
	result := content

	// Pattern 1: if (handler != null) handler(args);
	// Go doesn't support backreferences, so we match and verify manually
	pattern1 := regexp.MustCompile(`if\s*\(\s*(\w+)\s*!=\s*null\s*\)\s*(\w+)\s*\(([^)]*)\)\s*;`)
	matches := pattern1.FindAllStringSubmatchIndex(result, -1)

	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		handlerName1 := result[match[2]:match[3]]
		handlerName2 := result[match[4]:match[5]]
		args := result[match[6]:match[7]]

		if handlerName1 == handlerName2 {
			replacement := handlerName1 + "?.Invoke(" + args + ");"
			result = result[:match[0]] + replacement + result[match[1]:]
			changed = true
		}
	}

	// Pattern 2: if (handler != null) { handler(args); }
	pattern2 := regexp.MustCompile(`if\s*\(\s*(\w+)\s*!=\s*null\s*\)\s*\{\s*(\w+)\s*\(([^)]*)\)\s*;\s*\}`)
	matches2 := pattern2.FindAllStringSubmatchIndex(result, -1)

	for i := len(matches2) - 1; i >= 0; i-- {
		match := matches2[i]
		handlerName1 := result[match[2]:match[3]]
		handlerName2 := result[match[4]:match[5]]
		args := result[match[6]:match[7]]

		if handlerName1 == handlerName2 {
			replacement := handlerName1 + "?.Invoke(" + args + ");"
			result = result[:match[0]] + replacement + result[match[1]:]
			changed = true
		}
	}

	// Pattern 3: Multi-line if (handler != null) \n handler(args);
	pattern3 := regexp.MustCompile(`(?m)if\s*\(\s*(\w+)\s*!=\s*null\s*\)\s*\n\s*(\w+)\s*\(([^)]*)\)\s*;`)
	matches3 := pattern3.FindAllStringSubmatchIndex(result, -1)

	for i := len(matches3) - 1; i >= 0; i-- {
		match := matches3[i]
		handlerName1 := result[match[2]:match[3]]
		handlerName2 := result[match[4]:match[5]]
		args := result[match[6]:match[7]]

		if handlerName1 == handlerName2 {
			replacement := handlerName1 + "?.Invoke(" + args + ");"
			result = result[:match[0]] + replacement + result[match[1]:]
			changed = true
		}
	}

	// Pattern 4: if (X != null) X.Invoke(args);
	pattern4 := regexp.MustCompile(`if\s*\(\s*(\w+)\s*!=\s*null\s*\)\s*(\w+)\.Invoke\s*\(([^)]*)\)\s*;`)
	matches4 := pattern4.FindAllStringSubmatchIndex(result, -1)

	for i := len(matches4) - 1; i >= 0; i-- {
		match := matches4[i]
		handlerName1 := result[match[2]:match[3]]
		handlerName2 := result[match[4]:match[5]]
		args := result[match[6]:match[7]]

		if handlerName1 == handlerName2 {
			replacement := handlerName1 + "?.Invoke(" + args + ");"
			result = result[:match[0]] + replacement + result[match[1]:]
			changed = true
		}
	}

	return result, changed
}
