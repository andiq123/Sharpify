# ⚡ Sharpify

**Modernize your C# code instantly.**

Transform legacy C# into modern, clean code with one command.

## Install

```bash
# Download binary from releases
curl -L https://github.com/andiq123/Sharpify/releases/latest/download/sharpify-darwin-arm64 -o sharpify
chmod +x sharpify

# Or build from source
go install github.com/andiq123/sharpify@latest
```

## Usage

```bash
# Improve all C# files in current directory
sharpify -b .

# Preview changes without modifying files
sharpify -b --dry-run .

# Improve a specific file
sharpify -b myfile.cs

# Show what rules will be applied
sharpify -b --verbose .

# List all available rules
sharpify --list-rules
```

## Example

**Before:**
```csharp
namespace MyApp.Services
{
    public class UserService
    {
        private readonly ILogger _logger;
        
        public UserService(ILogger logger)
        {
            if (logger == null)
                throw new ArgumentNullException(nameof(logger));
            _logger = logger;
        }
        
        public User FindUser(int id)
        {
            var users = GetUsers();
            if (users.Count() > 0)
            {
                return users.Where(u => u.Id == id).FirstOrDefault();
            }
            return null;
        }
    }
}
```

**After:**
```csharp
namespace MyApp.Services;

public class UserService(ILogger logger)
{
    private readonly ILogger _logger = logger ?? throw new ArgumentNullException(nameof(logger));
    
    public User FindUser(int id)
    {
        var users = GetUsers();
        if (users.Any())
        {
            return users.FirstOrDefault(u => u.Id == id);
        }
        return null;
    }
}
```

## What it does

| Rule | Before | After |
|------|--------|-------|
| File-scoped namespace | `namespace X { }` | `namespace X;` |
| Throw expression | `if (x == null) throw...` | `x ?? throw` |
| LINQ optimization | `.Count() > 0` | `.Any()` |
| LINQ simplify | `.Where(p).First()` | `.First(p)` |
| Tuple swap | `temp=a; a=b; b=temp` | `(a,b) = (b,a)` |
| Collection expression | `new List<T> {...}` | `[...]` |
| Pattern matching | `x == null` | `x is null` |
| Spread operator | `.Concat().ToList()` | `[..a, ..b]` |
| ThrowIfNull | `if (x is null) throw` | `ArgumentNullException.ThrowIfNull(x)` |

Run `sharpify --list-rules` to see all 36 rules.

## Options

| Flag | Description |
|------|-------------|
| `-b` | Batch mode (non-interactive) |
| `--dry-run` | Preview changes only |
| `--verbose` | Show detailed output |
| `--list-rules` | List all rules |
| `--help` | Show help |

## Interactive Mode

Just run `sharpify` without flags for an interactive experience with menus.

## License

MIT
