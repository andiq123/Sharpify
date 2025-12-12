package rules

import (
	"regexp"
	"strings"
)

// FileScopedNamespace converts block-scoped namespaces to file-scoped (C# 10+)
type FileScopedNamespace struct {
	BaseVersionedRule
}

func NewFileScopedNamespace() *FileScopedNamespace {
	return &FileScopedNamespace{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp10, safe: true},
	}
}

func (r *FileScopedNamespace) Name() string {
	return "file-scoped-namespace"
}

func (r *FileScopedNamespace) Description() string {
	return "Convert block-scoped namespace to file-scoped namespace (C# 10+)"
}

func (r *FileScopedNamespace) Apply(content string) (string, bool) {
	// Check if already file-scoped
	fileScopedPattern := regexp.MustCompile(`(?m)^namespace\s+[\w.]+\s*;`)
	if fileScopedPattern.MatchString(content) {
		return content, false
	}

	// Match traditional namespace declaration with opening brace
	pattern := regexp.MustCompile(`(?m)^(\s*namespace\s+[\w.]+)\s*\n?\s*\{`)
	if !pattern.MatchString(content) {
		return content, false
	}

	// Only convert if there's exactly one namespace in the file
	allNamespaces := regexp.MustCompile(`(?m)^(\s*)namespace\s+[\w.]+`)
	matches := allNamespaces.FindAllString(content, -1)
	if len(matches) != 1 {
		return content, false
	}

	result := r.convertToFileScoped(content)
	if result == content {
		return content, false
	}

	return result, true
}

func (r *FileScopedNamespace) convertToFileScoped(content string) string {
	lines := strings.Split(content, "\n")
	var result []string
	inNamespace := false
	braceCount := 0
	namespaceStartIdx := -1

	nsPattern := regexp.MustCompile(`^(\s*)namespace\s+([\w.]+)\s*\{?\s*$`)

	for i, line := range lines {
		if !inNamespace {
			matches := nsPattern.FindStringSubmatch(line)
			if matches != nil {
				namespaceName := matches[2]

				if strings.Contains(line, "{") {
					braceCount = 1
					inNamespace = true
					namespaceStartIdx = i
					result = append(result, "namespace "+namespaceName+";")
					result = append(result, "")
					continue
				}

				if i+1 < len(lines) && strings.TrimSpace(lines[i+1]) == "{" {
					braceCount = 1
					inNamespace = true
					namespaceStartIdx = i + 1
					result = append(result, "namespace "+namespaceName+";")
					result = append(result, "")
					continue
				}
			}
		}

		if inNamespace && i == namespaceStartIdx {
			continue
		}

		if inNamespace {
			for _, ch := range line {
				if ch == '{' {
					braceCount++
				} else if ch == '}' {
					braceCount--
				}
			}

			if braceCount == 0 {
				inNamespace = false
				continue
			}

			result = append(result, r.removeIndent(line))
		} else if namespaceStartIdx == -1 || i < namespaceStartIdx {
			result = append(result, line)
		} else {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

func (r *FileScopedNamespace) removeIndent(line string) string {
	if strings.TrimSpace(line) == "" {
		return ""
	}

	if strings.HasPrefix(line, "\t") {
		return strings.TrimPrefix(line, "\t")
	}
	if strings.HasPrefix(line, "    ") {
		return strings.TrimPrefix(line, "    ")
	}

	return line
}
