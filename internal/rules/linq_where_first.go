package rules

import (
	"regexp"
)

// LinqWhereFirst combines Where + First/Single/Last into single predicate call
type LinqWhereFirst struct {
	BaseVersionedRule
}

func NewLinqWhereFirst() *LinqWhereFirst {
	return &LinqWhereFirst{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp6, safe: true},
	}
}

func (r *LinqWhereFirst) Name() string {
	return "linq-where-first"
}

func (r *LinqWhereFirst) Description() string {
	return "Combine Where + First/Single/Last into single call"
}

func (r *LinqWhereFirst) Apply(content string) (string, bool) {
	changed := false
	result := content

	// Pattern: .Where(predicate).First() -> .First(predicate)
	pattern1 := regexp.MustCompile(`\.Where\(([^)]+)\)\.First\(\)`)
	if pattern1.MatchString(result) {
		result = pattern1.ReplaceAllString(result, ".First(${1})")
		changed = true
	}

	// Pattern: .Where(predicate).FirstOrDefault() -> .FirstOrDefault(predicate)
	pattern2 := regexp.MustCompile(`\.Where\(([^)]+)\)\.FirstOrDefault\(\)`)
	if pattern2.MatchString(result) {
		result = pattern2.ReplaceAllString(result, ".FirstOrDefault(${1})")
		changed = true
	}

	// Pattern: .Where(predicate).Single() -> .Single(predicate)
	pattern3 := regexp.MustCompile(`\.Where\(([^)]+)\)\.Single\(\)`)
	if pattern3.MatchString(result) {
		result = pattern3.ReplaceAllString(result, ".Single(${1})")
		changed = true
	}

	// Pattern: .Where(predicate).SingleOrDefault() -> .SingleOrDefault(predicate)
	pattern4 := regexp.MustCompile(`\.Where\(([^)]+)\)\.SingleOrDefault\(\)`)
	if pattern4.MatchString(result) {
		result = pattern4.ReplaceAllString(result, ".SingleOrDefault(${1})")
		changed = true
	}

	// Pattern: .Where(predicate).Last() -> .Last(predicate)
	pattern5 := regexp.MustCompile(`\.Where\(([^)]+)\)\.Last\(\)`)
	if pattern5.MatchString(result) {
		result = pattern5.ReplaceAllString(result, ".Last(${1})")
		changed = true
	}

	// Pattern: .Where(predicate).LastOrDefault() -> .LastOrDefault(predicate)
	pattern6 := regexp.MustCompile(`\.Where\(([^)]+)\)\.LastOrDefault\(\)`)
	if pattern6.MatchString(result) {
		result = pattern6.ReplaceAllString(result, ".LastOrDefault(${1})")
		changed = true
	}

	// Pattern: .Where(predicate).Any() -> .Any(predicate)
	pattern7 := regexp.MustCompile(`\.Where\(([^)]+)\)\.Any\(\)`)
	if pattern7.MatchString(result) {
		result = pattern7.ReplaceAllString(result, ".Any(${1})")
		changed = true
	}

	// Pattern: .Where(predicate).Count() -> .Count(predicate)
	pattern8 := regexp.MustCompile(`\.Where\(([^)]+)\)\.Count\(\)`)
	if pattern8.MatchString(result) {
		result = pattern8.ReplaceAllString(result, ".Count(${1})")
		changed = true
	}

	return result, changed
}
