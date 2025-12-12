package rules

import (
	"regexp"
	"strings"
)

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
	pattern := regexp.MustCompile(`(\s+)([A-Z][a-zA-Z0-9_]*(?:<[^>]+>)?)\s+(\w+)\s*=\s*new\s+([A-Z][a-zA-Z0-9_]*(?:<[^>]+>)?)\s*([(\[{])`)

	if !pattern.MatchString(content) {
		return content, false
	}

	changed := false
	result := pattern.ReplaceAllStringFunc(content, func(match string) string {
		submatches := pattern.FindStringSubmatch(match)
		if len(submatches) >= 6 {
			whitespace := submatches[1]
			leftType := submatches[2]
			varName := submatches[3]
			rightType := submatches[4]
			suffix := submatches[5]

			if r.isFieldOrProperty(whitespace, content, match) {
				return match
			}

			if strings.TrimSpace(leftType) == strings.TrimSpace(rightType) {
				changed = true
				return whitespace + "var " + varName + " = new " + rightType + suffix
			}
		}
		return match
	})

	return result, changed
}

func (r *VarPattern) isFieldOrProperty(whitespace string, content string, match string) bool {
	idx := strings.Index(content, match)
	if idx == -1 {
		return false
	}

	lineStart := strings.LastIndex(content[:idx], "\n")
	if lineStart == -1 {
		lineStart = 0
	} else {
		lineStart++
	}

	linePrefix := content[lineStart:idx]

	fieldModifiers := []string{"private", "public", "protected", "internal", "static", "readonly", "const"}
	for _, mod := range fieldModifiers {
		if strings.Contains(linePrefix, mod) {
			return true
		}
	}

	return false
}
