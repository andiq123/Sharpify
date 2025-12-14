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

	
	
	

	
	outVarPattern := regexp.MustCompile(`out\s+var\s+(\w+)`)
	matches := outVarPattern.FindAllStringSubmatchIndex(result, -1)

	
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		varName := result[match[2]:match[3]]

		
		if varName == "_" {
			continue
		}

		
		afterDecl := result[match[1]:]
		
		usagePattern := regexp.MustCompile(`\b` + regexp.QuoteMeta(varName) + `\b`)
		usageMatches := usagePattern.FindAllStringIndex(afterDecl, -1)

		
		if len(usageMatches) == 0 {
			
			result = result[:match[0]] + "out _" + result[match[1]:]
			changed = true
		}
	}

	return result, changed
}
