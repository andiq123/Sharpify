package transformer

import (
	"github.com/andiq123/sharpify/internal/rules"
	"github.com/andiq123/sharpify/internal/scanner"
)


type Result struct {
	File         scanner.FileInfo
	NewContent   string
	Changed      bool
	AppliedRules []rules.RuleResult
}


type Transformer struct {
	rules []rules.Rule
}


func New(ruleList []rules.Rule) *Transformer {
	return &Transformer{
		rules: ruleList,
	}
}


func (t *Transformer) Transform(file scanner.FileInfo) Result {
	result := Result{
		File:       file,
		NewContent: file.Content,
		Changed:    false,
	}

	for _, rule := range t.rules {
		newContent, applied := rule.Apply(result.NewContent)
		if applied {
			result.NewContent = newContent
			result.Changed = true
			result.AppliedRules = append(result.AppliedRules, rules.RuleResult{
				RuleName:    rule.Name(),
				Applied:     true,
				Description: rule.Description(),
			})
		}
	}

	return result
}


func (t *Transformer) TransformAll(files []scanner.FileInfo) []Result {
	results := make([]Result, 0, len(files))
	for _, file := range files {
		results = append(results, t.Transform(file))
	}
	return results
}
