package rules

import (
	"regexp"
)

// StopwatchStartNew converts new Stopwatch() + Start() pattern to Stopwatch.StartNew()
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

	// Pattern: var x = new Stopwatch(); x.Start();
	// Also handles: Stopwatch x = new Stopwatch(); x.Start();

	// First find all variable declarations
	declPattern := regexp.MustCompile(`(?m)(var|Stopwatch)\s+(\w+)\s*=\s*new\s+Stopwatch\s*\(\s*\)\s*;`)

	matches := declPattern.FindAllStringSubmatchIndex(content, -1)

	// Process in reverse order to maintain indices
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		varName := content[match[4]:match[5]]
		declEnd := match[1]

		// Look for varName.Start(); after this declaration
		remaining := content[declEnd:]
		startPattern := regexp.MustCompile(`^\s*\n?\s*` + regexp.QuoteMeta(varName) + `\.Start\s*\(\s*\)\s*;`)

		startMatch := startPattern.FindStringIndex(remaining)
		if startMatch != nil {
			// Found the pattern! Replace both parts
			fullEnd := declEnd + startMatch[1]
			newCode := "var " + varName + " = Stopwatch.StartNew();"
			content = content[:match[0]] + newCode + content[fullEnd:]
			changed = true
		}
	}

	return content, changed
}
