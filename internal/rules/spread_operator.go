package rules

import (
	"regexp"
)


type SpreadOperator struct {
	BaseVersionedRule
}

func NewSpreadOperator() *SpreadOperator {
	return &SpreadOperator{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp12, safe: true},
	}
}

func (r *SpreadOperator) Name() string {
	return "spread-operator"
}

func (r *SpreadOperator) Description() string {
	return "Use spread operator in collection expressions (C# 12+)"
}

func (r *SpreadOperator) Apply(content string) (string, bool) {
	changed := false
	result := content

	
	pattern1 := regexp.MustCompile(`(\w+)\.Concat\s*\(\s*(\w+)\s*\)\.ToList\s*\(\s*\)`)
	if pattern1.MatchString(result) {
		result = pattern1.ReplaceAllString(result, "[..${1}, ..${2}]")
		changed = true
	}

	
	pattern2 := regexp.MustCompile(`(\w+)\.Concat\s*\(\s*(\w+)\s*\)\.ToArray\s*\(\s*\)`)
	if pattern2.MatchString(result) {
		result = pattern2.ReplaceAllString(result, "[..${1}, ..${2}]")
		changed = true
	}

	
	pattern3 := regexp.MustCompile(`Enumerable\.Concat\s*\(\s*(\w+)\s*,\s*(\w+)\s*\)\.ToList\s*\(\s*\)`)
	if pattern3.MatchString(result) {
		result = pattern3.ReplaceAllString(result, "[..${1}, ..${2}]")
		changed = true
	}

	
	
	
	pattern4 := regexp.MustCompile(`new\s+List<\w+>\s*\{\s*(\w+)\s*\}\.Concat\s*\(\s*(\w+)\s*\)\.ToList\s*\(\s*\)`)
	if pattern4.MatchString(result) {
		result = pattern4.ReplaceAllString(result, "[${1}, ..${2}]")
		changed = true
	}

	
	pattern5 := regexp.MustCompile(`(\w+)\.Append\s*\(\s*(\w+)\s*\)\.ToArray\s*\(\s*\)`)
	if pattern5.MatchString(result) {
		result = pattern5.ReplaceAllString(result, "[..${1}, ${2}]")
		changed = true
	}

	
	pattern6 := regexp.MustCompile(`(\w+)\.Prepend\s*\(\s*(\w+)\s*\)\.ToArray\s*\(\s*\)`)
	if pattern6.MatchString(result) {
		result = pattern6.ReplaceAllString(result, "[${2}, ..${1}]")
		changed = true
	}

	return result, changed
}
