namespace MyApp.Services;

public class UserService
{
    private readonly List<string> _users;
    private readonly ILogger _logger;

    public UserService(ILogger logger)
    {
        _users = [];
        _logger = logger;
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
        if (input is null) input = "default";
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

public interface ILogger
{
    void Log(string message);
}
