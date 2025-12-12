package rules

import (
	"regexp"
	"strings"
)

type SwitchExpression struct {
	BaseVersionedRule
}

func NewSwitchExpression() *SwitchExpression {
	return &SwitchExpression{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp8, safe: false},
	}
}

func (r *SwitchExpression) Name() string {
	return "switch-expression"
}

func (r *SwitchExpression) Description() string {
	return "Convert simple return-only switch statements to switch expressions (C# 8+)"
}

func (r *SwitchExpression) Apply(content string) (string, bool) {
	changed := false
	result := content

	// Match simple switch statements that only contain return statements
	// Pattern: switch (expr) { case value: return result; ... default: return result; }
	switchPattern := regexp.MustCompile(`(?s)(\s*)switch\s*\((\w+)\)\s*\{([^}]+)\}`)

	matches := switchPattern.FindAllStringSubmatchIndex(result, -1)

	// Process matches in reverse to preserve indices
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		indent := result[match[2]:match[3]]
		varName := result[match[4]:match[5]]
		body := result[match[6]:match[7]]

		// Try to convert this switch to an expression
		expr, ok := r.convertToSwitchExpr(varName, body, indent)
		if ok {
			// Replace the switch statement with the expression
			result = result[:match[0]] + expr + result[match[1]:]
			changed = true
		}
	}

	return result, changed
}

func (r *SwitchExpression) convertToSwitchExpr(varName, body, indent string) (string, bool) {
	lines := strings.Split(body, "\n")

	var arms []string
	var hasDefault bool

	currentCase := ""
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Match case label
		caseMatch := regexp.MustCompile(`^case\s+(.+?):$`).FindStringSubmatch(line)
		if caseMatch != nil {
			currentCase = caseMatch[1]
			continue
		}

		// Match default label
		if line == "default:" {
			currentCase = "_"
			hasDefault = true
			continue
		}

		// Match return statement
		returnMatch := regexp.MustCompile(`^return\s+(.+?);$`).FindStringSubmatch(line)
		if returnMatch != nil && currentCase != "" {
			arms = append(arms, currentCase+" => "+returnMatch[1])
			currentCase = ""
			continue
		}

		// If we encounter anything else (break, complex logic), abort
		if strings.HasPrefix(line, "break") {
			continue // Skip breaks
		}

		// Complex statement - can't convert
		if currentCase != "" && line != "{" && line != "}" {
			return "", false
		}
	}

	// Need at least 2 arms and a default
	if len(arms) < 2 || !hasDefault {
		return "", false
	}

	// Build the switch expression
	armsStr := strings.Join(arms, ",\n"+indent+"    ")
	expr := indent + "return " + varName + " switch\n" + indent + "{\n" + indent + "    " + armsStr + "\n" + indent + "};"

	return expr, true
}
