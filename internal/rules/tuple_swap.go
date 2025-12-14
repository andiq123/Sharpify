package rules

import (
	"regexp"
	"strings"
)


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

	
	
	
	
	
	
	lines := strings.Split(result, "\n")

	for i := 0; i < len(lines)-2; i++ {
		
		line1Pattern := regexp.MustCompile(`^(\s*)(var\s+)?(\w+)\s*=\s*(\w+)\s*;\s*$`)
		match1 := line1Pattern.FindStringSubmatch(lines[i])
		if match1 == nil {
			continue
		}

		indent := match1[1]
		tempVar := match1[3]
		varA := match1[4]

		
		line2Pattern := regexp.MustCompile(`^\s*(\w+)\s*=\s*(\w+)\s*;\s*$`)
		match2 := line2Pattern.FindStringSubmatch(lines[i+1])
		if match2 == nil {
			continue
		}

		
		if match2[1] != varA {
			continue
		}
		varB := match2[2]

		
		match3 := line2Pattern.FindStringSubmatch(lines[i+2])
		if match3 == nil {
			continue
		}

		
		if match3[1] != varB || match3[2] != tempVar {
			continue
		}

		
		lowerTemp := strings.ToLower(tempVar)
		isTempVar := lowerTemp == "temp" || lowerTemp == "tmp" || lowerTemp == "t" ||
			lowerTemp == "swap" || strings.HasPrefix(lowerTemp, "temp") ||
			strings.HasPrefix(lowerTemp, "tmp") || lowerTemp == "aux"

		if !isTempVar {
			continue
		}

		
		replacement := indent + "(" + varA + ", " + varB + ") = (" + varB + ", " + varA + ");"
		lines[i] = replacement
		lines[i+1] = "" 
		lines[i+2] = "" 
		changed = true
		i += 2 
	}

	if changed {
		
		result = ""
		first := true
		for _, line := range lines {
			if line == "" {
				
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
