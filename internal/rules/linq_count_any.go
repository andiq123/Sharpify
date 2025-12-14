package rules

import (
	"regexp"
)


type LinqCountAny struct {
	BaseVersionedRule
}

func NewLinqCountAny() *LinqCountAny {
	return &LinqCountAny{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp6, safe: true},
	}
}

func (r *LinqCountAny) Name() string {
	return "linq-count-any"
}

func (r *LinqCountAny) Description() string {
	return "Use Any() instead of Count() > 0 for performance"
}

func (r *LinqCountAny) Apply(content string) (string, bool) {
	changed := false
	result := content

	
	pattern1 := regexp.MustCompile(`\.Count\(\)\s*>\s*0`)
	if pattern1.MatchString(result) {
		result = pattern1.ReplaceAllString(result, ".Any()")
		changed = true
	}

	
	pattern2 := regexp.MustCompile(`\.Count\(\)\s*>=\s*1`)
	if pattern2.MatchString(result) {
		result = pattern2.ReplaceAllString(result, ".Any()")
		changed = true
	}

	
	pattern3 := regexp.MustCompile(`\.Count\(\)\s*!=\s*0`)
	if pattern3.MatchString(result) {
		result = pattern3.ReplaceAllString(result, ".Any()")
		changed = true
	}

	
	pattern4 := regexp.MustCompile(`(\w+)\.Count\(\)\s*==\s*0`)
	if pattern4.MatchString(result) {
		result = pattern4.ReplaceAllString(result, "!${1}.Any()")
		changed = true
	}

	
	pattern5 := regexp.MustCompile(`(\w+)\.Count\(\)\s*<\s*1`)
	if pattern5.MatchString(result) {
		result = pattern5.ReplaceAllString(result, "!${1}.Any()")
		changed = true
	}

	return result, changed
}
