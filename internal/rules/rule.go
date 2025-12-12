package rules

// Rule defines the interface for transformation rules (Open/Closed Principle)
type Rule interface {
	Name() string
	Description() string
	Apply(content string) (string, bool)
}

// RuleResult holds the result of applying a rule
type RuleResult struct {
	RuleName    string
	Applied     bool
	Description string
}
