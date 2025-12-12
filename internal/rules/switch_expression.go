package rules

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
	return "Convert switch statements to switch expressions (C# 8+) [manual review]"
}

func (r *SwitchExpression) Apply(content string) (string, bool) {
	return content, false
}
