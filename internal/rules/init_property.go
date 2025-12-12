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
	return "Add required modifier to non-nullable properties without default values (C# 11+)"
}

func (r *RequiredProperty) Apply(content string) (string, bool) {
	changed := false

	pattern := regexp.MustCompile(`(?m)^(\s*)(public\s+)(string|object|[A-Z][a-zA-Z0-9_]*(?:<[^>]+>)?)\s+(\w+)\s*\{\s*get;\s*(?:set|init);\s*\}(\s*)$`)

	result := pattern.ReplaceAllStringFunc(content, func(match string) string {
		if regexp.MustCompile(`public\s+required\s+`).MatchString(match) {
			return match
		}
		if regexp.MustCompile(`\?\s+\w+\s*\{`).MatchString(match) {
			return match
		}
		if regexp.MustCompile(`=\s*[^;]+;`).MatchString(match) {
			return match
		}
		if regexp.MustCompile(`\b(int|long|double|float|bool|byte|char|decimal|short|uint|ulong|ushort|sbyte|DateTime|Guid)\s+\w+\s*\{`).MatchString(match) {
			return match
		}

		submatches := pattern.FindStringSubmatch(match)
		if submatches == nil {
			return match
		}

		changed = true
		return submatches[1] + submatches[2] + "required " + submatches[3] + " " + submatches[4] + " { get; set; }" + submatches[5]
	})

	return result, changed
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
