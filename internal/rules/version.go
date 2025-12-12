package rules

// CSharpVersion represents a C# language version
type CSharpVersion int

const (
	CSharp6 CSharpVersion = iota + 6
	CSharp7
	CSharp8
	CSharp9
	CSharp10
	CSharp11
	CSharp12
	CSharp13
)

// String returns the version as a string
func (v CSharpVersion) String() string {
	switch v {
	case CSharp6:
		return "C# 6.0"
	case CSharp7:
		return "C# 7.x"
	case CSharp8:
		return "C# 8.0"
	case CSharp9:
		return "C# 9.0"
	case CSharp10:
		return "C# 10.0"
	case CSharp11:
		return "C# 11.0"
	case CSharp12:
		return "C# 12.0"
	case CSharp13:
		return "C# 13.0"
	default:
		return "Unknown"
	}
}

// DotNetVersion returns the corresponding .NET version
func (v CSharpVersion) DotNetVersion() string {
	switch v {
	case CSharp6:
		return ".NET Framework 4.6+ / .NET Core 1.0+"
	case CSharp7:
		return ".NET Framework 4.7+ / .NET Core 2.0+"
	case CSharp8:
		return ".NET Core 3.0+ / .NET Standard 2.1"
	case CSharp9:
		return ".NET 5.0+"
	case CSharp10:
		return ".NET 6.0+"
	case CSharp11:
		return ".NET 7.0+"
	case CSharp12:
		return ".NET 8.0+"
	case CSharp13:
		return ".NET 9.0+"
	default:
		return "Unknown"
	}
}

// VersionedRule extends Rule with version information
type VersionedRule interface {
	Rule
	MinVersion() CSharpVersion
	IsSafe() bool // Whether the rule is guaranteed not to break logic
}

// BaseVersionedRule provides common versioned rule functionality
type BaseVersionedRule struct {
	minVersion CSharpVersion
	safe       bool
}

func (r *BaseVersionedRule) MinVersion() CSharpVersion {
	return r.minVersion
}

func (r *BaseVersionedRule) IsSafe() bool {
	return r.safe
}
