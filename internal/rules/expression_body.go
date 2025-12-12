package rules

import (
	"regexp"
)

type ExpressionBody struct {
	BaseVersionedRule
}

func NewExpressionBody() *ExpressionBody {
	return &ExpressionBody{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp6, safe: true},
	}
}

func (r *ExpressionBody) Name() string {
	return "expression-body"
}

func (r *ExpressionBody) Description() string {
	return "Use expression-bodied members for simple methods/properties (C# 6+)"
}

func (r *ExpressionBody) Apply(content string) (string, bool) {
	changed := false
	result := content

	propPattern := regexp.MustCompile(`((?:public|private|protected|internal|static|\s)+\w+(?:<[^>]+>)?\s+\w+)\s*\{\s*get\s*\{\s*return\s+([^;]+);\s*\}\s*\}`)
	if propPattern.MatchString(result) {
		result = propPattern.ReplaceAllString(result, "${1} => ${2};")
		changed = true
	}

	return result, changed
}
