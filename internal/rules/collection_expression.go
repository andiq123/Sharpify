package rules

import (
	"regexp"
)

// CollectionExpression converts array/list initializers to collection expressions (C# 12+)
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

	// Pattern: new Type[] { ... } -> [...]
	pattern1 := regexp.MustCompile(`new\s+\w+\[\]\s*\{\s*([^}]*)\s*\}`)
	if pattern1.MatchString(result) {
		result = pattern1.ReplaceAllString(result, "[${1}]")
		changed = true
	}

	// Pattern: new List<Type> { ... } -> [...]
	pattern2 := regexp.MustCompile(`new\s+List<\w+>\s*\{\s*([^}]*)\s*\}`)
	if pattern2.MatchString(result) {
		result = pattern2.ReplaceAllString(result, "[${1}]")
		changed = true
	}

	// Pattern: Array.Empty<Type>() -> []
	pattern3 := regexp.MustCompile(`Array\.Empty<\w+>\(\)`)
	if pattern3.MatchString(result) {
		result = pattern3.ReplaceAllString(result, "[]")
		changed = true
	}

	// Pattern: new List<Type>() -> [] (empty list)
	pattern4 := regexp.MustCompile(`new\s+List<\w+>\(\)`)
	if pattern4.MatchString(result) {
		result = pattern4.ReplaceAllString(result, "[]")
		changed = true
	}

	return result, changed
}
