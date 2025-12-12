# Sharpify

A fast CLI tool to modernize legacy C# codebases. Upgrade your code to use the latest C# features without breaking existing logic.

## Features

- **Simple Interface** - Run, Settings, Rules, Exit
- **Persistent Settings** - Configure once, run instantly
- **Version-Aware** - Applies transformations compatible with your target .NET version
- **Safe by Default** - Conservative transformations that preserve logic
- **Fast** - Written in Go for maximum performance

## Installation

### Go Install

```bash
go install github.com/andiq123/sharpify@latest
```

### Download Binary

Download from [GitHub Releases](https://github.com/andiq123/Sharpify/releases).

| Platform | Download |
|----------|----------|
| macOS (Apple Silicon) | `sharpify-macos-arm64` |
| Windows | `sharpify-windows-amd64.exe` |

## Usage

### Interactive Mode (Default)

```bash
sharpify
```

Menu:
- **Run** - Scan and improve C# files in current directory
- **Settings** - Configure C# version, safe mode, backup
- **Rules** - View available transformation rules
- **Exit**

Settings are saved to `~/.sharpify.json` and persist across sessions.

### Batch Mode

```bash
sharpify -b                           # Improve current directory
sharpify -b ./src                     # Improve specific path
sharpify -b --dry-run ./src           # Preview changes
sharpify -b --rules var-pattern ./src # Apply specific rules
sharpify -b --verbose ./src           # Detailed output
```

## Settings

| Setting | Default | Description |
|---------|---------|-------------|
| C# Version | 12 | Target C# version (6-13) |
| Safe Mode | On | Only apply safe transformations |
| Backup | Off | Create backups before changes |

## Transformation Rules (25 total)

### C# 6+ (.NET Framework 4.6+)
| Rule | Description | Safe |
|------|-------------|------|
| `expression-body` | Expression-bodied members | Yes |
| `string-interpolation` | `string.Format()` to `$"..."` | Yes |
| `nameof-expression` | `nameof()` for exceptions | Yes |
| `null-propagation` | `?.` operator | Yes |
| `var-pattern` | `var` for obvious types | Yes |

### C# 7+ (.NET Core 2.0+)
| Rule | Description | Safe |
|------|-------------|------|
| `pattern-matching` | `is Type` pattern | Yes |
| `default-literal` | `default(T)` to `default` | Yes |
| `tuple-deconstruction` | `Tuple<>` to `()` | Yes |
| `span-suggestion` | Optimize string operations | No |

### C# 8+ (.NET Core 3.0+)
| Rule | Description | Safe |
|------|-------------|------|
| `null-coalescing-assignment` | `??=` operator | Yes |
| `index-range` | `^1` for last element | Yes |

### C# 9+ (.NET 5+)
| Rule | Description | Safe |
|------|-------------|------|
| `target-typed-new` | `new()` syntax | No |
| `pattern-matching-null` | `is null` / `is not null` | Yes |

### C# 10+ (.NET 6+)
| Rule | Description | Safe |
|------|-------------|------|
| `file-scoped-namespace` | `namespace X;` | Yes |
| `implicit-using` | Remove implicit usings | No |

### C# 11+ (.NET 7+)
| Rule | Description | Safe |
|------|-------------|------|
| `list-pattern` | `is []` for empty check | No |
| `required-property` | Add `required` modifier | No |

### C# 12+ (.NET 8+)
| Rule | Description | Safe |
|------|-------------|------|
| `collection-expression` | `[]` syntax | Yes |
| `primary-constructor` | Primary constructors | No |

## Example

### Before
```csharp
using System;
using System.Collections.Generic;

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
    }
}
```

### After (C# 12)
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
}
```

## Backups

When enabled, backups are saved to:
```
.sharpify-backup/YYYYMMDD-HHMMSS/
```

## License

MIT
