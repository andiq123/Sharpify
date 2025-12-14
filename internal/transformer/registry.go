package transformer

import (
	"sort"

	"github.com/andiq123/sharpify/internal/rules"
)


type RuleRegistry struct {
	rules map[string]rules.Rule
}


func NewRegistry() *RuleRegistry {
	r := &RuleRegistry{
		rules: make(map[string]rules.Rule),
	}

	
	r.Register(rules.NewExpressionBody())
	r.Register(rules.NewStringInterpolation())
	r.Register(rules.NewStringConcatToInterpolation())
	r.Register(rules.NewNameofExpression())
	r.Register(rules.NewNullPropagation())
	r.Register(rules.NewVarPattern())
	r.Register(rules.NewStopwatchStartNew())
	r.Register(rules.NewConditionalAccessDelegate())
	r.Register(rules.NewExceptionFilter())
	r.Register(rules.NewLinqCountAny())
	r.Register(rules.NewLinqWhereFirst())

	
	r.Register(rules.NewPatternMatching())
	r.Register(rules.NewDefaultLiteral())
	r.Register(rules.NewTupleDeconstruction())
	r.Register(rules.NewDiscardVariable())
	r.Register(rules.NewSpanSuggestion())
	r.Register(rules.NewThrowExpression())
	r.Register(rules.NewTupleSwap())

	
	r.Register(rules.NewNullCoalescing())
	r.Register(rules.NewIndexRange())
	r.Register(rules.NewSwitchExpression())

	
	r.Register(rules.NewTargetTypedNew())
	r.Register(rules.NewPatternMatchingNull())
	r.Register(rules.NewRecordType())
	r.Register(rules.NewInitOnlyProperty())

	
	r.Register(rules.NewFileScopedNamespace())
	r.Register(rules.NewGlobalUsing())
	r.Register(rules.NewImplicitUsing())
	r.Register(rules.NewThrowHelper())

	
	r.Register(rules.NewRawStringLiteral())
	r.Register(rules.NewRequiredProperty())
	r.Register(rules.NewListPattern())
	r.Register(rules.NewStringIsNullOrEmpty())

	
	r.Register(rules.NewCollectionExpression())
	r.Register(rules.NewPrimaryConstructor())
	r.Register(rules.NewSpreadOperator())

	return r
}


func (r *RuleRegistry) Register(rule rules.Rule) {
	r.rules[rule.Name()] = rule
}


func (r *RuleRegistry) Get(name string) (rules.Rule, bool) {
	rule, ok := r.rules[name]
	return rule, ok
}


func (r *RuleRegistry) All() []rules.Rule {
	result := make([]rules.Rule, 0, len(r.rules))
	for _, rule := range r.rules {
		result = append(result, rule)
	}
	return result
}


func (r *RuleRegistry) AllSafe() []rules.Rule {
	result := make([]rules.Rule, 0)
	for _, rule := range r.rules {
		if vr, ok := rule.(rules.VersionedRule); ok {
			if vr.IsSafe() {
				result = append(result, rule)
			}
		}
	}
	return result
}


func (r *RuleRegistry) GetByVersion(version rules.CSharpVersion, safeOnly bool) []rules.Rule {
	result := make([]rules.Rule, 0)
	for _, rule := range r.rules {
		if vr, ok := rule.(rules.VersionedRule); ok {
			if vr.MinVersion() <= version {
				if !safeOnly || vr.IsSafe() {
					result = append(result, rule)
				}
			}
		}
	}

	
	sort.Slice(result, func(i, j int) bool {
		vi := result[i].(rules.VersionedRule).MinVersion()
		vj := result[j].(rules.VersionedRule).MinVersion()
		if vi != vj {
			return vi < vj
		}
		return result[i].Name() < result[j].Name()
	})

	return result
}


func (r *RuleRegistry) Names() []string {
	result := make([]string, 0, len(r.rules))
	for name := range r.rules {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}


func (r *RuleRegistry) GetByNames(names []string) []rules.Rule {
	result := make([]rules.Rule, 0, len(names))
	for _, name := range names {
		if rule, ok := r.rules[name]; ok {
			result = append(result, rule)
		}
	}
	return result
}


func (r *RuleRegistry) GroupByVersion() map[rules.CSharpVersion][]rules.Rule {
	groups := make(map[rules.CSharpVersion][]rules.Rule)

	for _, rule := range r.rules {
		if vr, ok := rule.(rules.VersionedRule); ok {
			version := vr.MinVersion()
			groups[version] = append(groups[version], rule)
		}
	}

	
	for version := range groups {
		sort.Slice(groups[version], func(i, j int) bool {
			return groups[version][i].Name() < groups[version][j].Name()
		})
	}

	return groups
}
