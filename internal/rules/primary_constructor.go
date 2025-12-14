package rules

import (
	"regexp"
	"strings"
)

type PrimaryConstructor struct {
	BaseVersionedRule
}

func NewPrimaryConstructor() *PrimaryConstructor {
	return &PrimaryConstructor{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp12, safe: false},
	}
}

func (r *PrimaryConstructor) Name() string {
	return "primary-constructor"
}

func (r *PrimaryConstructor) Description() string {
	return "Convert simple constructor patterns to primary constructors (C# 12)"
}

func (r *PrimaryConstructor) Apply(content string) (string, bool) {
	changed := false

	
	
	
	
	
	classPattern := regexp.MustCompile(`(?m)((?:public|internal|private|protected)\s+(?:sealed\s+|abstract\s+|partial\s+)*(?:class|struct|record)\s+(\w+)(?:<[^>]+>)?)((\s*:\s*[^{]+)?)\s*\{`)

	matches := classPattern.FindAllStringSubmatchIndex(content, -1)

	for i := len(matches) - 1; i >= 0; i-- {
		match := matches[i]
		classNameAndKeyword := content[match[2]:match[3]] 
		className := content[match[4]:match[5]]           
		inheritance := ""
		if match[6] != -1 && match[7] != -1 {
			inheritance = strings.TrimSpace(content[match[6]:match[7]]) 
		}
		openBracePos := match[1] - 1

		classBodyStart := openBracePos + 1
		classBodyEnd := findMatchingBrace(content, openBracePos)
		if classBodyEnd == -1 {
			continue
		}

		classBody := content[classBodyStart:classBodyEnd]

		constructorPattern := regexp.MustCompile(`(?m)^\s*public\s+` + regexp.QuoteMeta(className) + `\s*\(([^)]*)\)\s*\{([^}]*)\}`)
		ctorMatch := constructorPattern.FindStringSubmatchIndex(classBody)
		if ctorMatch == nil {
			continue
		}

		params := strings.TrimSpace(classBody[ctorMatch[2]:ctorMatch[3]])
		ctorBody := strings.TrimSpace(classBody[ctorMatch[4]:ctorMatch[5]])

		if params == "" {
			continue
		}

		assignments := parseConstructorAssignments(ctorBody)
		if len(assignments) == 0 {
			continue
		}

		paramList := parseParameters(params)
		if len(paramList) == 0 {
			continue
		}

		allParamsAssigned := true
		for _, p := range paramList {
			found := false
			for _, a := range assignments {
				if a.paramName == p.name {
					found = true
					break
				}
			}
			if !found {
				allParamsAssigned = false
				break
			}
		}
		if !allParamsAssigned {
			continue
		}

		if !allFieldsExist(classBody, assignments) {
			continue
		}

		newClassBody := classBody
		ctorFullMatch := classBody[ctorMatch[0]:ctorMatch[1]]
		newClassBody = strings.Replace(newClassBody, ctorFullMatch, "", 1)

		for _, a := range assignments {
			fieldPattern := regexp.MustCompile(`(?m)^\s*private\s+(?:readonly\s+)?` + regexp.QuoteMeta(a.fieldType) + `\s+` + regexp.QuoteMeta(a.fieldName) + `\s*;\s*\n?`)
			if fieldPattern.MatchString(newClassBody) {
				newClassBody = fieldPattern.ReplaceAllString(newClassBody, "")
			}
		}

		var newFieldDecls []string
		for _, a := range assignments {
			newFieldDecls = append(newFieldDecls, "    private readonly "+a.fieldType+" "+a.fieldName+" = "+a.paramName+";")
		}

		newClassBody = strings.TrimLeft(newClassBody, "\n\t ")
		newClassBody = "\n" + strings.Join(newFieldDecls, "\n") + "\n\n    " + newClassBody

		newClassBody = regexp.MustCompile(`\n{3,}`).ReplaceAllString(newClassBody, "\n\n")

		
		newClassDecl := classNameAndKeyword + "(" + params + ")"
		if inheritance != "" {
			newClassDecl += " " + inheritance
		}

		newContent := content[:match[2]] + newClassDecl + "\n{" + newClassBody + content[classBodyEnd:]
		content = newContent
		changed = true
	}

	return content, changed
}

type assignment struct {
	fieldName string
	fieldType string
	paramName string
}

func parseConstructorAssignments(body string) []assignment {
	var result []assignment

	lines := strings.Split(body, ";")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		pattern := regexp.MustCompile(`^(_\w+|\w+)\s*=\s*(\w+)$`)
		match := pattern.FindStringSubmatch(line)
		if match == nil {
			return nil
		}

		result = append(result, assignment{
			fieldName: match[1],
			paramName: match[2],
		})
	}

	return result
}

func parseParameters(params string) []struct{ typ, name string } {
	var result []struct{ typ, name string }

	parts := strings.Split(params, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		tokens := strings.Fields(part)
		if len(tokens) < 2 {
			continue
		}
		result = append(result, struct{ typ, name string }{
			typ:  strings.Join(tokens[:len(tokens)-1], " "),
			name: tokens[len(tokens)-1],
		})
	}

	return result
}

func allFieldsExist(classBody string, assignments []assignment) bool {
	for i := range assignments {
		fieldPattern := regexp.MustCompile(`private\s+(?:readonly\s+)?(\w+(?:<[^>]+>)?)\s+` + regexp.QuoteMeta(assignments[i].fieldName) + `\s*;`)
		match := fieldPattern.FindStringSubmatch(classBody)
		if match == nil {
			return false
		}
		assignments[i].fieldType = match[1]
	}
	return true
}

func findMatchingBrace(content string, openPos int) int {
	if openPos >= len(content) || content[openPos] != '{' {
		return -1
	}

	depth := 1
	for i := openPos + 1; i < len(content); i++ {
		switch content[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}
