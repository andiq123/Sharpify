package rules

import (
	"regexp"
	"strings"
)


type StringConcatToInterpolation struct {
	BaseVersionedRule
}

func NewStringConcatToInterpolation() *StringConcatToInterpolation {
	return &StringConcatToInterpolation{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp6, safe: false}, 
	}
}

func (r *StringConcatToInterpolation) Name() string {
	return "string-concat-interpolation"
}

func (r *StringConcatToInterpolation) Description() string {
	return "Convert string concatenation to interpolation (C# 6+)"
}

func (r *StringConcatToInterpolation) Apply(content string) (string, bool) {
	changed := false

	
	

	
	pattern1 := regexp.MustCompile(`"([^"]*?)"\s*\+\s*([a-zA-Z_][a-zA-Z0-9_]*(?:\.[a-zA-Z_][a-zA-Z0-9_]*)*)`)

	
	pattern2 := regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*(?:\.[a-zA-Z_][a-zA-Z0-9_]*)*)\s*\+\s*"([^"]*?)"`)

	
	result := pattern1.ReplaceAllStringFunc(content, func(match string) string {
		submatches := pattern1.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}
		literal := submatches[1]
		variable := submatches[2]

		
		
		if strings.HasPrefix(literal, "{") || strings.HasSuffix(literal, "}") {
			return match
		}

		changed = true
		return `$"` + literal + `{` + variable + `}"`
	})

	
	
	result = pattern2.ReplaceAllStringFunc(result, func(match string) string {
		submatches := pattern2.FindStringSubmatch(match)
		if len(submatches) < 3 {
			return match
		}
		variable := submatches[1]
		literal := submatches[2]

		
		if variable == "string" || variable == "String" {
			return match
		}

		
		if strings.Contains(match, `$"`) {
			return match
		}

		changed = true
		return `$"{` + variable + `}` + literal + `"`
	})

	return result, changed
}
