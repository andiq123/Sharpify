package rules

import (
	"regexp"
)

type InitOnlyProperty struct {
	BaseVersionedRule
}

func NewInitOnlyProperty() *InitOnlyProperty {
	return &InitOnlyProperty{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp9, safe: false},
	}
}

func (r *InitOnlyProperty) Name() string {
	return "init-property"
}

func (r *InitOnlyProperty) Description() string {
	return "Suggest init-only properties for immutable types (C# 9+) [manual review]"
}

func (r *InitOnlyProperty) Apply(content string) (string, bool) {
	return content, false
}

type RequiredProperty struct {
	BaseVersionedRule
}

func NewRequiredProperty() *RequiredProperty {
	return &RequiredProperty{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp11, safe: false},
	}
}

func (r *RequiredProperty) Name() string {
	return "required-property"
}

func (r *RequiredProperty) Description() string {
	return "Suggest required modifier for properties (C# 11+) [manual review]"
}

func (r *RequiredProperty) Apply(content string) (string, bool) {
	return content, false
}

type RecordType struct {
	BaseVersionedRule
}

func NewRecordType() *RecordType {
	return &RecordType{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp9, safe: false},
	}
}

func (r *RecordType) Name() string {
	return "record-type"
}

func (r *RecordType) Description() string {
	return "Suggest converting simple classes to records (C# 9+) [manual review]"
}

func (r *RecordType) Apply(content string) (string, bool) {
	return content, false
}

type GlobalUsing struct {
	BaseVersionedRule
}

func NewGlobalUsing() *GlobalUsing {
	return &GlobalUsing{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp10, safe: false},
	}
}

func (r *GlobalUsing) Name() string {
	return "global-using"
}

func (r *GlobalUsing) Description() string {
	return "Suggest global usings for common namespaces (C# 10+) [manual review]"
}

func (r *GlobalUsing) Apply(content string) (string, bool) {
	return content, false
}

type ImplicitUsing struct {
	BaseVersionedRule
}

func NewImplicitUsing() *ImplicitUsing {
	return &ImplicitUsing{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp10, safe: false},
	}
}

func (r *ImplicitUsing) Name() string {
	return "implicit-using"
}

func (r *ImplicitUsing) Description() string {
	return "Remove usings covered by implicit usings (.NET 6+) [manual review]"
}

var implicitUsings = []string{
	"System",
	"System.Collections.Generic",
	"System.IO",
	"System.Linq",
	"System.Net.Http",
	"System.Threading",
	"System.Threading.Tasks",
}

func (r *ImplicitUsing) Apply(content string) (string, bool) {
	changed := false
	result := content

	for _, ns := range implicitUsings {
		pattern := regexp.MustCompile(`(?m)^using\s+` + regexp.QuoteMeta(ns) + `\s*;\s*\n?`)
		if pattern.MatchString(result) {
			result = pattern.ReplaceAllString(result, "")
			changed = true
		}
	}

	return result, changed
}
