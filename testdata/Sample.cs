namespace MyApp.Services;

public class UserService
{
    private readonly List<string> _users;
    private readonly ILogger _logger;
    private static Dictionary<string, int> _cache = new();

    public UserService(ILogger logger)
    {
        _users = new List<string>();
        _logger = logger;
    }

    public void Initialize()
    {
        List<string> items = [];
        string[] names = ["John", "Jane"];
        int[] numbers = [1, 2, 3];
        var empty = Array.Empty<string>();

        var permissions = new List<PermissionInfo>
        {
            new PermissionInfo { Name = "Permission1", CanRead = true, CanWrite = false, CanDelete = false },
            new PermissionInfo { Name = "Permission2", CanRead = true, CanWrite = true, CanDelete = false }
        };
    }

    public string GetUserName(int id)
    {
        if (id is null)
        {
            throw new ArgumentNullException(nameof(id));
        }
        return $"User {id}";
    }

    public List<string> GetUsers()
    {
        return _users;
    }

    public bool IsValid(object obj)
    {
        if (obj is UserService)
        {
            return true;
        }
        return false;
    }

    public string GetValue(string input)
    {
        input ??= "default";
        return input;
    }

    public int Count => _users.Count;

    public string GetLastUser()
    {
        return _users[^1];
    }

    public string GetPrefix(string text, int length)
    {
        return text[..length];
    }

    public string GetSuffix(string text, int start)
    {
        return text[start..];
    }

    public void ProcessUser(string name)
    {
        string value = default;
        if (name is not null)
        {
            _logger.Log(name);
        }
    }

    public bool TryGetValue(string key, out string unused)
    {
        unused = null;
        return _users.Contains(key);
    }

    public (string, int) GetUserInfo()
    {
        return ("John", 25);
    }
}

public class PermissionInfo
{
    public string Name { get; set; }
    public bool CanRead { get; set; }
    public bool CanWrite { get; set; }
    public bool CanDelete { get; set; }
}

public interface ILogger
{
    void Log(string message);
}
