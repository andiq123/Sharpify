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

	
	

	
	isLambdaContext := func(line string, matchStart int) bool {
		
		prefix := line[:matchStart]
		return regexp.MustCompile(`=>\s*`).MatchString(prefix)
	}

	
	lines := regexp.MustCompile(`(?m)^.*$`).FindAllString(result, -1)
	var processedLines []string

	for _, line := range lines {
		processedLine := line

		
		if regexp.MustCompile(`=>`).MatchString(line) {
			processedLines = append(processedLines, processedLine)
			continue
		}

		
		
		nullEqPattern := regexp.MustCompile(`(\b_?\w+)\s*==\s*null\b`)
		if nullEqPattern.MatchString(processedLine) {
			
			if regexp.MustCompile(`\b(if|while|return)\s*\(`).MatchString(processedLine) ||
				regexp.MustCompile(`[=,\(]\s*_?\w+\s*==\s*null`).MatchString(processedLine) {
				processedLine = nullEqPattern.ReplaceAllString(processedLine, "${1} is null")
				changed = true
			}
		}

		
		nullNeqPattern := regexp.MustCompile(`(\b_?\w+)\s*!=\s*null\b`)
		if nullNeqPattern.MatchString(processedLine) {
			
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

	
	
	pattern5 := regexp.MustCompile(`([=,\(]\s*)(\w+)\s*==\s*null\s*\?`)
	if pattern5.MatchString(result) {
		result = pattern5.ReplaceAllStringFunc(result, func(match string) string {
			
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

	_ = isLambdaContext 

	return result, changed
}
