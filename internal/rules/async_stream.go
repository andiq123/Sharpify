package rules

import (
	"regexp"
)

// AsyncDisposable converts IDisposable to IAsyncDisposable pattern (C# 8+)
type AsyncDisposable struct {
	BaseVersionedRule
}

func NewAsyncDisposable() *AsyncDisposable {
	return &AsyncDisposable{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp8, safe: false},
	}
}

func (r *AsyncDisposable) Name() string {
	return "async-disposable"
}

func (r *AsyncDisposable) Description() string {
	return "Suggest IAsyncDisposable for async resources (C# 8+) [manual review]"
}

func (r *AsyncDisposable) Apply(content string) (string, bool) {
	return content, false
}

// UsingAwait converts using to await using for async disposables (C# 8+)
type UsingAwait struct {
	BaseVersionedRule
}

func NewUsingAwait() *UsingAwait {
	return &UsingAwait{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp8, safe: false},
	}
}

func (r *UsingAwait) Name() string {
	return "await-using"
}

func (r *UsingAwait) Description() string {
	return "Use await using for async disposables (C# 8+) [manual review]"
}

func (r *UsingAwait) Apply(content string) (string, bool) {
	return content, false
}

// IndexRange uses Index and Range types (C# 8+)
type IndexRange struct {
	BaseVersionedRule
}

func NewIndexRange() *IndexRange {
	return &IndexRange{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp8, safe: true},
	}
}

func (r *IndexRange) Name() string {
	return "index-range"
}

func (r *IndexRange) Description() string {
	return "Use Index (^) and Range (..) operators (C# 8+)"
}

func (r *IndexRange) Apply(content string) (string, bool) {
	changed := false
	result := content

	// Pattern: arr[arr.Length - 1] -> arr[^1]
	// Go doesn't support backreferences, so we match and compare manually
	pattern1 := regexp.MustCompile(`(\w+)\s*\[\s*(\w+)\.Length\s*-\s*(\d+)\s*\]`)
	matches := pattern1.FindAllStringSubmatch(result, -1)
	for _, m := range matches {
		if len(m) >= 4 && m[1] == m[2] {
			old := m[0]
			replacement := m[1] + "[^" + m[3] + "]"
			result = regexp.MustCompile(regexp.QuoteMeta(old)).ReplaceAllString(result, replacement)
			changed = true
		}
	}

	// Pattern: str.Substring(0, n) -> str[..n]
	pattern2 := regexp.MustCompile(`(\w+)\.Substring\s*\(\s*0\s*,\s*(\w+)\s*\)`)
	if pattern2.MatchString(result) {
		result = pattern2.ReplaceAllString(result, "${1}[..${2}]")
		changed = true
	}

	// Pattern: str.Substring(n) -> str[n..]
	pattern3 := regexp.MustCompile(`(\w+)\.Substring\s*\(\s*(\w+)\s*\)`)
	if pattern3.MatchString(result) {
		result = pattern3.ReplaceAllString(result, "${1}[${2}..]")
		changed = true
	}

	return result, changed
}
