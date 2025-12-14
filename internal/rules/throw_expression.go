package rules

import (
	"regexp"
)

// ThrowExpression converts null-check-throw-assign patterns to throw expressions
type ThrowExpression struct {
	BaseVersionedRule
}

func NewThrowExpression() *ThrowExpression {
	return &ThrowExpression{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp7, safe: true},
	}
}

func (r *ThrowExpression) Name() string {
	return "throw-expression"
}

func (r *ThrowExpression) Description() string {
	return "Use throw expressions for null checks (C# 7+)"
}

func (r *ThrowExpression) Apply(content string) (string, bool) {
	changed := false
	result := content

	// Pattern: if (x == null) throw new ArgumentNullException(nameof(x)); followed by _field = x;
	// Go doesn't support backreferences, so we match and verify manually
	pattern1 := regexp.MustCompile(`(?m)if\s*\(\s*(\w+)\s*==\s*null\s*\)\s*\{?\s*throw\s+new\s+ArgumentNullException\s*\(\s*nameof\s*\(\s*(\w+)\s*\)\s*\)\s*;\s*\}?\s*\n\s*(_?\w+)\s*=\s*(\w+)\s*;`)
	matches := pattern1.FindAllStringSubmatchIndex(result, -1)

	// Process in reverse to preserve indices
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		varName := result[match[2]:match[3]]     // first capture (x in null check)
		nameofVar := result[match[4]:match[5]]   // nameof argument
		fieldName := result[match[6]:match[7]]   // _field
		assignedVar := result[match[8]:match[9]] // x in assignment

		// Verify all variable references are consistent
		if varName == nameofVar && varName == assignedVar {
			replacement := fieldName + " = " + varName + " ?? throw new ArgumentNullException(nameof(" + varName + "));"
			result = result[:match[0]] + replacement + result[match[1]:]
			changed = true
		}
	}

	// Pattern 2: if (x == null) throw new ArgumentNullException("x"); followed by _field = x;
	pattern2 := regexp.MustCompile(`(?m)if\s*\(\s*(\w+)\s*==\s*null\s*\)\s*\{?\s*throw\s+new\s+ArgumentNullException\s*\(\s*"(\w+)"\s*\)\s*;\s*\}?\s*\n\s*(_?\w+)\s*=\s*(\w+)\s*;`)
	matches2 := pattern2.FindAllStringSubmatchIndex(result, -1)

	for i := len(matches2) - 1; i >= 0; i-- {
		match := matches2[i]
		varName := result[match[2]:match[3]]
		strArg := result[match[4]:match[5]]
		fieldName := result[match[6]:match[7]]
		assignedVar := result[match[8]:match[9]]

		if varName == assignedVar {
			replacement := fieldName + " = " + varName + " ?? throw new ArgumentNullException(\"" + strArg + "\");"
			result = result[:match[0]] + replacement + result[match[1]:]
			changed = true
		}
	}

	// Pattern 3: Handle this.field = value pattern
	pattern3 := regexp.MustCompile(`(?m)if\s*\(\s*(\w+)\s*==\s*null\s*\)\s*\{?\s*throw\s+new\s+ArgumentNullException\s*\(\s*nameof\s*\(\s*(\w+)\s*\)\s*\)\s*;\s*\}?\s*\n\s*(this\.)?(\w+)\s*=\s*(\w+)\s*;`)
	matches3 := pattern3.FindAllStringSubmatchIndex(result, -1)

	for i := len(matches3) - 1; i >= 0; i-- {
		match := matches3[i]
		varName := result[match[2]:match[3]]
		nameofVar := result[match[4]:match[5]]
		thisPrefix := ""
		if match[6] != -1 && match[7] != -1 {
			thisPrefix = result[match[6]:match[7]]
		}
		fieldName := result[match[8]:match[9]]
		assignedVar := result[match[10]:match[11]]

		if varName == nameofVar && varName == assignedVar {
			replacement := thisPrefix + fieldName + " = " + varName + " ?? throw new ArgumentNullException(nameof(" + varName + "));"
			result = result[:match[0]] + replacement + result[match[1]:]
			changed = true
		}
	}

	return result, changed
}
