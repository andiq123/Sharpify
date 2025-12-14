package rules

import (
	"regexp"
)


type StopwatchStartNew struct {
	BaseVersionedRule
}

func NewStopwatchStartNew() *StopwatchStartNew {
	return &StopwatchStartNew{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp6, safe: true},
	}
}

func (r *StopwatchStartNew) Name() string {
	return "stopwatch-start-new"
}

func (r *StopwatchStartNew) Description() string {
	return "Use Stopwatch.StartNew() instead of new + Start() (C# 6+)"
}

func (r *StopwatchStartNew) Apply(content string) (string, bool) {
	changed := false

	
	

	
	declPattern := regexp.MustCompile(`(?m)(var|Stopwatch)\s+(\w+)\s*=\s*new\s+Stopwatch\s*\(\s*\)\s*;`)

	matches := declPattern.FindAllStringSubmatchIndex(content, -1)

	
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		varName := content[match[4]:match[5]]
		declEnd := match[1]

		
		remaining := content[declEnd:]
		startPattern := regexp.MustCompile(`^\s*\n?\s*` + regexp.QuoteMeta(varName) + `\.Start\s*\(\s*\)\s*;`)

		startMatch := startPattern.FindStringIndex(remaining)
		if startMatch != nil {
			
			fullEnd := declEnd + startMatch[1]
			newCode := "var " + varName + " = Stopwatch.StartNew();"
			content = content[:match[0]] + newCode + content[fullEnd:]
			changed = true
		}
	}

	return content, changed
}
