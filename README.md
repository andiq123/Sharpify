# Sharpify

Modernize legacy C# code with a fast, interactive CLI tool.

## Quick Start

```bash
# Install
go install github.com/andiq123/sharpify@latest

# Run (interactive)
sharpify
```

Or download from [GitHub Releases](https://github.com/andiq123/Sharpify/releases).

## Features

- **Interactive UI** - Simple menu: Run → Settings → Rules → Exit
- **Smart Path Selection** - Prompts for project path if no C# files found
- **Rule Toggle** - Enable/disable rules before each run
- **Persistent Config** - Settings saved to `~/.sharpify.json`
- **Version-Aware** - Only applies rules compatible with your C# version
- **Safe by Default** - Conservative transformations that won't break code

## Usage

### Interactive Mode

```bash
sharpify
```

1. **Run** - Scan files, toggle rules, apply changes
2. **Settings** - C# version (6-13), safe mode, backups
3. **Rules** - View all available transformations
4. **Exit**

### Batch Mode

```bash
sharpify -b                    # Current directory
sharpify -b ./src              # Specific path  
sharpify -b --dry-run          # Preview only
sharpify -b --verbose          # Detailed output
```

## Transformation Rules

| Version | Examples |
|---------|----------|
| **C# 6+** | Expression bodies, string interpolation, `?.` operator, `nameof()` |
| **C# 7+** | Pattern matching, `default` literal, tuple deconstruction |
| **C# 8+** | `??=` operator, `^1` index, ranges |
| **C# 9+** | `new()` syntax, `is null` / `is not null` |
| **C# 10+** | File-scoped namespaces |
| **C# 11+** | List patterns, `required` properties |
| **C# 12+** | Collection expressions `[]`, primary constructors |

Run `sharpify --list-rules` for the full list.

## Example

**Before:**
```csharp
namespace MyApp.Services
{
    public class UserService
    {
        private readonly List<string> _users = new List<string>();
        public int Count { get { return _users.Count; } }
        public string Last() { return _users[_users.Length - 1]; }
    }
}
```

**After (C# 12):**
```csharp
namespace MyApp.Services;

public class UserService
{
    private readonly List<string> _users = [];
    public int Count => _users.Count;
    public string Last() => _users[^1];
}
```

## License

MIT
