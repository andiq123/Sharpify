package rules

import (
	"regexp"
	"strings"
)

type CollectionExpression struct {
	BaseVersionedRule
}

func NewCollectionExpression() *CollectionExpression {
	return &CollectionExpression{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp12, safe: true},
	}
}

func (r *CollectionExpression) Name() string {
	return "collection-expression"
}

func (r *CollectionExpression) Description() string {
	return "Use collection expressions (C# 12+)"
}

func (r *CollectionExpression) Apply(content string) (string, bool) {
	changed := false
	result := content

	result, c := r.convertEmptyListWithType(result)
	if c {
		changed = true
	}

	result, c = r.convertEmptyArrayWithType(result)
	if c {
		changed = true
	}

	result, c = r.convertArrayEmptyWithType(result)
	if c {
		changed = true
	}

	result, c = r.convertSimpleArrayInitializer(result)
	if c {
		changed = true
	}

	return result, changed
}

func (r *CollectionExpression) convertEmptyListWithType(content string) (string, bool) {
	pattern := regexp.MustCompile(`(List<[^>]+>\s+\w+\s*=\s*)new\s+List<[^>]+>\(\)`)
	if pattern.MatchString(content) {
		return pattern.ReplaceAllString(content, "${1}[]"), true
	}
	return content, false
}

func (r *CollectionExpression) convertEmptyArrayWithType(content string) (string, bool) {
	pattern := regexp.MustCompile(`((?:string|int|long|double|float|bool|byte|char|decimal|short|uint|ulong|ushort|sbyte|object|[A-Z][a-zA-Z0-9_]*)\[\]\s+\w+\s*=\s*)new\s+(?:string|int|long|double|float|bool|byte|char|decimal|short|uint|ulong|ushort|sbyte|object|[A-Z][a-zA-Z0-9_]*)\[\]\s*\{\s*\}`)
	if pattern.MatchString(content) {
		return pattern.ReplaceAllString(content, "${1}[]"), true
	}
	return content, false
}

func (r *CollectionExpression) convertArrayEmptyWithType(content string) (string, bool) {
	pattern := regexp.MustCompile(`((?:string|int|long|double|float|bool|byte|char|decimal|short|uint|ulong|ushort|sbyte|object|[A-Z][a-zA-Z0-9_]*)\[\]\s+\w+\s*=\s*)Array\.Empty<[^>]+>\(\)`)
	if pattern.MatchString(content) {
		return pattern.ReplaceAllString(content, "${1}[]"), true
	}

	pattern2 := regexp.MustCompile(`(List<[^>]+>\s+\w+\s*=\s*)Array\.Empty<[^>]+>\(\)`)
	if pattern2.MatchString(content) {
		return pattern2.ReplaceAllString(content, "${1}[]"), true
	}
	return content, false
}

func (r *CollectionExpression) convertSimpleArrayInitializer(content string) (string, bool) {
	changed := false
	result := content

	pattern := regexp.MustCompile(`new\s+(string|int|long|double|float|bool|byte|char|decimal|short|uint|ulong|ushort|sbyte)\[\]\s*\{`)

	matches := pattern.FindAllStringIndex(result, -1)
	if len(matches) == 0 {
		return content, false
	}

	for i := len(matches) - 1; i >= 0; i-- {
		startIdx := matches[i][0]
		braceStart := strings.Index(result[startIdx:], "{")
		if braceStart == -1 {
			continue
		}
		braceStart += startIdx

		braceEnd := r.findMatchingBrace(result, braceStart)
		if braceEnd == -1 {
			continue
		}

		innerContent := result[braceStart+1 : braceEnd]

		if strings.Contains(innerContent, "{") {
			continue
		}

		innerContent = strings.TrimSpace(innerContent)
		if innerContent == "" {
			continue
		}

		newExpr := "[" + innerContent + "]"
		result = result[:startIdx] + newExpr + result[braceEnd+1:]
		changed = true
	}

	return result, changed
}

func (r *CollectionExpression) findMatchingBrace(content string, start int) int {
	if start >= len(content) || content[start] != '{' {
		return -1
	}

	count := 1
	for i := start + 1; i < len(content); i++ {
		switch content[i] {
		case '{':
			count++
		case '}':
			count--
			if count == 0 {
				return i
			}
		}
	}
	return -1
}
