# Sharpify

A fast, interactive CLI tool to modernize legacy C# codebases. Upgrade your code to use the latest C# features without breaking existing logic.

```
  ____  _                       _  __
 / ___|| |__   __ _ _ __ _ __ (_)/ _|_   _
 \___ \| '_ \ / _' | '__| '_ \| | |_| | | |
  ___) | | | | (_| | |  | |_) | |  _| |_| |
 |____/|_| |_|\__,_|_|  | .__/|_|_|  \__, |
                        |_|          |___/
```

## Features

- **Interactive Mode** - Visual menus to select target C# version and rules
- **Version-Aware** - Only applies transformations compatible with your target .NET version
- **Safe by Default** - Conservative transformations that preserve logic
- **Backup System** - Automatic backups before any changes
- **File Preview** - Review changes file-by-file with diff view
- **Fast** - Written in Go for maximum performance

## Installation

### Go Install

```bash
go install github.com/andiq123/sharpify@latest
```

### Download Binary

Download the latest release from [GitHub Releases](https://github.com/andiq123/Sharpify/releases).

| Platform | Download |
|----------|----------|
| macOS (Apple Silicon) | `sharpify-macos-arm64` |
| Windows | `sharpify-windows-amd64.exe` |

## Usage

### Interactive Mode (Default)

Simply run `sharpify` to start the interactive mode:

```bash
sharpify
```

You'll be guided through:
1. Selecting your target C# version (.NET 5-9)
2. Choosing which transformation rules to apply
3. Reviewing and applying changes file-by-file

### Batch Mode

For CI/CD or scripting:

```bash
# Preview changes
sharpify -b --dry-run ./src

# Apply all changes
sharpify -b ./src

# Apply specific rules only
sharpify -b --rules file-scoped-namespace,pattern-matching ./src

# Verbose output
sharpify -b --verbose ./MyProject
```

## Transformation Rules

### C# 6+ (.NET Framework 4.6+)

| Rule | Description |
|------|-------------|
| `expression-body` | Convert simple methods/properties to expression-bodied members |
| `string-interpolation` | Convert `string.Format()` to `$"..."` |
| `nameof-expression` | Use `nameof()` for `ArgumentNullException` parameters |
| `null-propagation` | Use `?.` null-conditional operator |

### C# 7+ (.NET Core 2.0+)

| Rule | Description |
|------|-------------|
| `pattern-matching` | Use `is Type` pattern matching |
| `default-literal` | Convert `default(T)` to `default` |
| `tuple-deconstruction` | Convert `Tuple<T1,T2>` to `(T1, T2)` |

### C# 8+ (.NET Core 3.0+)

| Rule | Description |
|------|-------------|
| `null-coalescing-assignment` | Use `??=` operator |
| `index-range` | Use `^1` and `[..]` syntax |

### C# 9+ (.NET 5+)

| Rule | Description |
|------|-------------|
| `target-typed-new` | Convert `Type x = new Type()` to `Type x = new()` |
| `pattern-matching-null` | Use `is null` and `is not null` |

### C# 10+ (.NET 6+)

| Rule | Description |
|------|-------------|
| `file-scoped-namespace` | Convert `namespace X { }` to `namespace X;` |
| `implicit-using` | Remove usings covered by implicit usings |

### C# 12+ (.NET 8+)

| Rule | Description |
|------|-------------|
| `collection-expression` | Convert `new List<T>()` to `[]` |

## Examples

### Before

```csharp
using System;
using System.Collections.Generic;
using System.Linq;

namespace MyApp.Services
{
    public class UserService
    {
        private readonly List<string> _users;

        public UserService()
        {
            _users = new List<string>();
        }

        public int Count { get { return _users.Count; } }

        public string GetLast()
        {
            return _users[_users.Length - 1];
        }

        public void Process(string name)
        {
            if (name == null)
            {
                throw new ArgumentNullException("name");
            }
            Console.WriteLine(string.Format("Hello {0}", name));
        }
    }
}
```

### After (targeting C# 12 / .NET 8)

```csharp
namespace MyApp.Services;

public class UserService
{
    private readonly List<string> _users;

    public UserService()
    {
        _users = [];
    }

    public int Count => _users.Count;

    public string GetLast()
    {
        return _users[^1];
    }

    public void Process(string name)
    {
        if (name is null)
        {
            throw new ArgumentNullException(nameof(name));
        }
        Console.WriteLine($"Hello {name}");
    }
}
```

## Safety

Sharpify prioritizes **safe transformations** that won't change your code's behavior:

- **Safe rules** (enabled by default): Guaranteed to preserve logic
- **Experimental rules** (opt-in): May require manual review

Use Settings in interactive mode to toggle between safe-only and all rules.

## Backups

Before modifying any file, Sharpify creates a backup in:

```
.sharpify-backup/YYYYMMDD-HHMMSS/
```

You can restore original files from this directory if needed.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
