// Test file to verify all Sharpify rules work correctly
// Each section tests a specific rule

namespace OldNamespace;

// ========================================
// Test: Primary Constructor (C# 12)
// Should transform class with simple ctor assignments
// ========================================
public class UserService(ILogger logger, IRepository repository)
{
private readonly ILogger _logger = logger;
private readonly IRepository _repository = repository;

public void DoWork() { }
}

// ========================================
// Test: Target-Typed New (C# 9+)
// ========================================
public class TargetTypedNewTest
{
    private List<string> _items = [];
    private Dictionary<string, int> _map = new();

    public void Method()
    {
        List<int> numbers = [];
        var sb = new StringBuilder();
    }
}

// ========================================
// Test: Null Coalescing (C# 8+)
// ========================================
public class NullCoalescingTest
{
    public void Method(string input)
    {
        if (input is null)
        {
            input = "default";
        }

        string value;
        if (input is not null)
        {
            value = input;
        }
        else
        {
            value = "fallback";
        }
    }
}

// ========================================
// Test: Pattern Matching (C# 7+)
// ========================================
public class PatternMatchingTest
{
    public bool IsString(object obj)
    {
        if (obj is string)
        {
            return true;
        }
        return false;
    }

    public void ProcessType(object item)
    {
        if (item is UserService)
        {
            var service = (UserService)item;
            service.DoWork();
        }
    }

    public bool CheckNull(object value)
    {
        if (value is not null)
        {
            return true;
        }
        return false;
    }
}

// ========================================
// Test: Expression Body (C# 6+)
// ========================================
public class ExpressionBodyTest
{
    private string _name;

    public string GetName()
    {
        return _name;
    }

    public int GetLength()
    {
        return _name.Length;
    }
}

// ========================================
// Test: String Interpolation (C# 6+)
// ========================================
public class StringInterpolationTest
{
    public string Format(string name, int age)
    {
        return $"Name: {name}, Age: {age}";
    }

    public string Concat(string a, string b)
    {
        return a + " - " + b;
    }
}

// ========================================
// Test: Var Pattern
// ========================================
public class VarPatternTest
{
    public void Method()
    {
        string text = "hello";
        int number = 42;
        List<string> items = [];
    }
}

// ========================================
// Test: Collection Expression (C# 12+)
// ========================================
public class CollectionExpressionTest
{
    public void Method()
    {
        var list = new List<int> { 1, 2, 3 };
        var array = [4, 5, 6];
        var empty = new List<string>();
    }
}

// ========================================
// Test: Null Propagation (C# 6+)
// ========================================
public class NullPropagationTest
{
    public int GetLength(string text)
    {
        if (text is not null)
        {
            return text.Length;
        }
        return 0;
    }
}

// ========================================
// Test: Switch Expression (C# 8+)
// ========================================
public class SwitchExpressionTest
{
    public string GetDayName(int day)
    {
        return day switch

        {

            1 => "Monday",

            2 => "Tuesday",

            3 => "Wednesday",

            _ => "Unknown"

        };
    }
}

// ========================================
// Test: Index/Range (C# 8+)
// ========================================
public class IndexRangeTest
{
    public char GetLastChar(string text)
    {
        return text[^1];
    }

    public string GetSubstring(string text)
    {
        return text.Substring(0, 5);
    }
}

// ========================================
// Test: NameOf (C# 6+)
// ========================================
public class NameofTest
{
    public void Method(string param)
    {
        if (param is null)
            throw new ArgumentNullException(nameof(param));
    }
}

// ========================================
// Test: Default Literal (C# 7.1+)
// ========================================
public class DefaultLiteralTest
{
    public void Method()
    {
        int number = default;
        string text = default;
        List<int> items = default;
    }
}

// Interfaces for testing
public interface ILogger { void Log(string msg); }
public interface IRepository { }
