package rules


type Rule interface {
	Name() string
	Description() string
	Apply(content string) (string, bool)
}


type RuleResult struct {
	RuleName    string
	Applied     bool
	Description string
}
