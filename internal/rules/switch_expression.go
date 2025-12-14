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

	
	
	switchPattern := regexp.MustCompile(`(?s)(\s*)switch\s*\((\w+)\)\s*\{([^}]+)\}`)

	matches := switchPattern.FindAllStringSubmatchIndex(result, -1)

	
	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		indent := result[match[2]:match[3]]
		varName := result[match[4]:match[5]]
		body := result[match[6]:match[7]]

		
		expr, ok := r.convertToSwitchExpr(varName, body, indent)
		if ok {
			
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

		
		caseMatch := regexp.MustCompile(`^case\s+(.+?):$`).FindStringSubmatch(line)
		if caseMatch != nil {
			currentCase = caseMatch[1]
			continue
		}

		
		if line == "default:" {
			currentCase = "_"
			hasDefault = true
			continue
		}

		
		returnMatch := regexp.MustCompile(`^return\s+(.+?);$`).FindStringSubmatch(line)
		if returnMatch != nil && currentCase != "" {
			arms = append(arms, currentCase+" => "+returnMatch[1])
			currentCase = ""
			continue
		}

		
		if strings.HasPrefix(line, "break") {
			continue 
		}

		
		if currentCase != "" && line != "{" && line != "}" {
			return "", false
		}
	}

	
	if len(arms) < 2 || !hasDefault {
		return "", false
	}

	
	armsStr := strings.Join(arms, ",\n"+indent+"    ")
	expr := indent + "return " + varName + " switch\n" + indent + "{\n" + indent + "    " + armsStr + "\n" + indent + "};"

	return expr, true
}
