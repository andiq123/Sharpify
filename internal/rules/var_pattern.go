package rules

import (
	"regexp"
	"strings"
)

// VarPattern converts explicit type declarations to var where appropriate
type VarPattern struct {
	BaseVersionedRule
}

func NewVarPattern() *VarPattern {
	return &VarPattern{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp6, safe: true},
	}
}

func (r *VarPattern) Name() string {
	return "var-pattern"
}

func (r *VarPattern) Description() string {
	return "Use var for obvious type declarations (new expressions)"
}

func (r *VarPattern) Apply(content string) (string, bool) {
	// Match: Type variable = new Type(...) -> var variable = new Type(...)
	pattern := regexp.MustCompile(`\b([A-Z][a-zA-Z0-9_]*(?:<[^>]+>)?)\s+(\w+)\s*=\s*new\s+([A-Z][a-zA-Z0-9_]*(?:<[^>]+>)?)\s*([(\[{])`)

	if !pattern.MatchString(content) {
		return content, false
	}

	changed := false
	result := pattern.ReplaceAllStringFunc(content, func(match string) string {
		submatches := pattern.FindStringSubmatch(match)
		if len(submatches) >= 5 {
			leftType := submatches[1]
			varName := submatches[2]
			rightType := submatches[3]
			suffix := submatches[4]

			if strings.TrimSpace(leftType) == strings.TrimSpace(rightType) {
				changed = true
				return "var " + varName + " = new " + rightType + suffix
			}
		}
		return match
	})

	return result, changed
}
