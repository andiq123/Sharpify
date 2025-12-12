package rules

// RawStringLiteral suggests raw string literals (C# 11+)
// This rule is disabled as it requires careful manual review
type RawStringLiteral struct {
	BaseVersionedRule
}

func NewRawStringLiteral() *RawStringLiteral {
	return &RawStringLiteral{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp11, safe: false},
	}
}

func (r *RawStringLiteral) Name() string {
	return "raw-string-literal"
}

func (r *RawStringLiteral) Description() string {
	return "Use raw string literals for complex strings (C# 11+) [manual review]"
}

func (r *RawStringLiteral) Apply(content string) (string, bool) {
	// This transformation is complex and can break strings
	// Disabled for safety - users should manually review
	return content, false
}
