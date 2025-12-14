package rules

import (
	"regexp"
	"strings"
)


type ExceptionFilter struct {
	BaseVersionedRule
}

func NewExceptionFilter() *ExceptionFilter {
	return &ExceptionFilter{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp6, safe: true},
	}
}

func (r *ExceptionFilter) Name() string {
	return "exception-filter"
}

func (r *ExceptionFilter) Description() string {
	return "Use exception filters (catch when) (C# 6+)"
}

func (r *ExceptionFilter) Apply(content string) (string, bool) {
	changed := false
	result := content

	
	

	
	

	
	
	pattern1 := regexp.MustCompile(`catch\s*\(\s*(\w+)\s+(\w+)\s*\)\s*\{\s*if\s*\(\s*!\s*([^)]+\([^)]*\)|[^)]+)\s*\)\s*throw\s*;\s*`)

	matches := pattern1.FindAllStringSubmatchIndex(result, -1)
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		excType := result[match[2]:match[3]]
		excVar := result[match[4]:match[5]]
		condition := result[match[6]:match[7]]

		
		condition = strings.TrimSpace(condition)

		replacement := "catch (" + excType + " " + excVar + ") when (" + condition + ")\n        {\n            "
		result = result[:match[0]] + replacement + result[match[1]:]
		changed = true
	}

	
	pattern2 := regexp.MustCompile(`catch\s*\(\s*(\w+)\s+(\w+)\s*\)\s*\{\s*if\s*\(\s*(\w+\.\w+)\s*!=\s*(\w+(?:\.\w+)*)\s*\)\s*throw\s*;\s*`)
	matches2 := pattern2.FindAllStringSubmatchIndex(result, -1)

	for i := len(matches2) - 1; i >= 0; i-- {
		match := matches2[i]
		excType := result[match[2]:match[3]]
		excVar := result[match[4]:match[5]]
		leftSide := result[match[6]:match[7]]
		rightSide := result[match[8]:match[9]]

		replacement := "catch (" + excType + " " + excVar + ") when (" + leftSide + " == " + rightSide + ")\n        {\n            "
		result = result[:match[0]] + replacement + result[match[1]:]
		changed = true
	}

	
	pattern3 := regexp.MustCompile(`catch\s*\(\s*(\w+)\s+(\w+)\s*\)\s*\{\s*if\s*\(\s*(\w+(?:\.\w+)?)\s*==\s*false\s*\)\s*throw\s*;\s*`)
	matches3 := pattern3.FindAllStringSubmatchIndex(result, -1)

	for i := len(matches3) - 1; i >= 0; i-- {
		match := matches3[i]
		excType := result[match[2]:match[3]]
		excVar := result[match[4]:match[5]]
		condition := result[match[6]:match[7]]

		replacement := "catch (" + excType + " " + excVar + ") when (" + condition + ")\n        {\n            "
		result = result[:match[0]] + replacement + result[match[1]:]
		changed = true
	}

	return result, changed
}
