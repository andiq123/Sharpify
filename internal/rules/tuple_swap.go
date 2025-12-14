package rules

import (
	"regexp"
	"strings"
)

// TupleSwap converts temp-swap patterns to tuple deconstruction
type TupleSwap struct {
	BaseVersionedRule
}

func NewTupleSwap() *TupleSwap {
	return &TupleSwap{
		BaseVersionedRule: BaseVersionedRule{minVersion: CSharp7, safe: true},
	}
}

func (r *TupleSwap) Name() string {
	return "tuple-swap"
}

func (r *TupleSwap) Description() string {
	return "Use tuple deconstruction for variable swaps (C# 7+)"
}

func (r *TupleSwap) Apply(content string) (string, bool) {
	changed := false
	result := content

	// We need to find the pattern across multiple lines:
	// (var )? temp = a;
	// a = b;
	// b = temp;
	//
	// Go regex doesn't support backreferences, so we scan line by line
	lines := strings.Split(result, "\n")

	for i := 0; i < len(lines)-2; i++ {
		// Try to match: [indent](var )?temp = a;
		line1Pattern := regexp.MustCompile(`^(\s*)(var\s+)?(\w+)\s*=\s*(\w+)\s*;\s*$`)
		match1 := line1Pattern.FindStringSubmatch(lines[i])
		if match1 == nil {
			continue
		}

		indent := match1[1]
		tempVar := match1[3]
		varA := match1[4]

		// Try to match: a = b;
		line2Pattern := regexp.MustCompile(`^\s*(\w+)\s*=\s*(\w+)\s*;\s*$`)
		match2 := line2Pattern.FindStringSubmatch(lines[i+1])
		if match2 == nil {
			continue
		}

		// The first var in line 2 should be varA, the second is varB
		if match2[1] != varA {
			continue
		}
		varB := match2[2]

		// Try to match: b = temp;
		match3 := line2Pattern.FindStringSubmatch(lines[i+2])
		if match3 == nil {
			continue
		}

		// The first var in line 3 should be varB, the second should be tempVar
		if match3[1] != varB || match3[2] != tempVar {
			continue
		}

		// Check if temp variable looks like a temp (temp, tmp, t, swap, etc.)
		lowerTemp := strings.ToLower(tempVar)
		isTempVar := lowerTemp == "temp" || lowerTemp == "tmp" || lowerTemp == "t" ||
			lowerTemp == "swap" || strings.HasPrefix(lowerTemp, "temp") ||
			strings.HasPrefix(lowerTemp, "tmp") || lowerTemp == "aux"

		if !isTempVar {
			continue
		}

		// Replace the three lines with tuple swap
		replacement := indent + "(" + varA + ", " + varB + ") = (" + varB + ", " + varA + ");"
		lines[i] = replacement
		lines[i+1] = "" // Mark for removal
		lines[i+2] = "" // Mark for removal
		changed = true
		i += 2 // Skip the two lines we just processed
	}

	if changed {
		// Filter out empty lines that were marked for removal
		var newLines []string
		for _, line := range lines {
			if line != "" || !changed {
				newLines = append(newLines, line)
			}
		}
		// Keep only non-deleted lines - need to be more careful here
		result = ""
		first := true
		for _, line := range lines {
			if line == "" && changed {
				// Skip empty lines that were replacements
				continue
			}
			if !first {
				result += "\n"
			}
			result += line
			first = false
		}
	}

	return result, changed
}
