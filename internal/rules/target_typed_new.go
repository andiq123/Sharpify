package rules

import (
	"regexp"
	"strings"
)

type TargetTypedNew struct {
	BaseVersionedRule
}

func NewTargetTypedNew() *TargetTypedNew {
	return &TargetTypedNew{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp9, safe: false},
	}
}

func (r *TargetTypedNew) Name() string {
	return "target-typed-new"
}

func (r *TargetTypedNew) Description() string {
	return "Use target-typed new expressions (C# 9+)"
}

func (r *TargetTypedNew) Apply(content string) (string, bool) {
	pattern := regexp.MustCompile(`(\b(?:private|public|protected|internal|static|readonly|\s)+)([A-Z][a-zA-Z0-9_]*(?:<[^>]+>)?)\s+(\w+)\s*=\s*new\s+([A-Z][a-zA-Z0-9_]*(?:<[^>]+>)?)\s*\(([^)]*)\)\s*;`)

	changed := false
	result := pattern.ReplaceAllStringFunc(content, func(match string) string {
		submatches := pattern.FindStringSubmatch(match)
		if len(submatches) >= 6 {
			modifiers := submatches[1]
			leftType := submatches[2]
			varName := submatches[3]
			rightType := submatches[4]
			args := submatches[5]

			if strings.TrimSpace(leftType) == strings.TrimSpace(rightType) {
				changed = true
				return modifiers + leftType + " " + varName + " = new(" + args + ");"
			}
		}
		return match
	})

	return result, changed
}
