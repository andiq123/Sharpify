using System.Diagnostics;
namespace LegacyCompany.BadCodeExamples;
/// <summary>
/// This file demonstrates various legacy C# patterns that Sharpify can modernize.
/// Run: ./sharpify -b testdata/LegacyCode.cs
/// </summary>
public class UserRepository
{
    private readonly IDatabase _database;
    private readonly ILogger _logger;
    private readonly ICache _cache;
    // BAD: Verbose null checks with throw - should use throw expression
    public UserRepository(IDatabase database, ILogger logger, ICache cache)
    {
        _database = database ?? throw new ArgumentNullException(nameof(database));
        _logger = logger ?? throw new ArgumentNullException(nameof(logger));
        _cache = cache ?? throw new ArgumentNullException(nameof(cache));
    }
    // BAD: Count() > 0 instead of Any()
    public bool HasUsers()
    {
        var users = GetAllUsers();
        if (users.Any())
        {
            return true;
        }
        if (users.Any())
        {
            return true;
        }
        if (users.Any())
        {
            return true;
        }
        return false;
    }
    // BAD: Where().First() instead of First(predicate)
    public User FindUserById(int id)
    {
        var users = GetAllUsers();
        var user = users.First(u => u.Id == id);
        var maybeUser = users.FirstOrDefault(u => u.Id == id);
        var singleUser = users.Single(u => u.Id == id);
        var maybeSingle = users.SingleOrDefault(u => u.Id == id);
        return user;
    }
    // BAD: Old temp swap pattern
    public void SwapUsers(ref User a, ref User b)
    {
        (a, b) = (b, a);
    }
    // BAD: Verbose event invocation
    public event EventHandler<UserEventArgs> UserCreated;
    public event EventHandler<UserEventArgs> UserDeleted;
    protected virtual void OnUserCreated(UserEventArgs args)
    {
        UserCreated?.Invoke(this, args);
    }
    protected virtual void OnUserDeleted(UserEventArgs args)
    {
        UserDeleted?.Invoke(this, args);
    }
    // BAD: string.IsNullOrEmpty instead of pattern matching
    public bool ValidateUsername(string username)
    {
        if (username is null or "")
        {
            return false;
        }
        if (!username is null or "")
        {
            return username.Length >= 3;
        }
        return false;
    }
    // BAD: Concat().ToList() instead of spread operator
    public List<User> CombineUserLists(List<User> activeUsers, List<User> inactiveUsers)
    {
        var allUsers = [..activeUsers, ..inactiveUsers];
        var userArray = [..activeUsers, ..inactiveUsers];
        return allUsers;
    }
    // BAD: Block-scoped catch instead of exception filter
    public async Task<User> FetchUserFromApiAsync(int id)
    {
        try
        {
            return await CallExternalApiAsync(id);
        }
        catch (HttpRequestException ex) when (ex.Message.Contains("NotFound"))
        {
            _logger.Log($"User not found: {id}");
            return null;
        }
    }
    // BAD: Old stopwatch pattern
    public void MeasurePerformance()
    {
        var stopwatch = Stopwatch.StartNew();
        DoExpensiveOperation();
        stopwatch.Stop();
        Console.WriteLine($"Elapsed: {stopwatch.ElapsedMilliseconds}");
    }
    // BAD: String concatenation instead of interpolation
    public string GetUserSummary(User user)
    {
        return $"User: {user.Name}" + $" (ID: {user.Id}" + $") - Email: {user.Email}";
    }
    // BAD: (T1, T2) instead of value tuples
    public (string, int) GetUserInfo(User user)
    {
        return (user.Name, user.Id);
    }
    // BAD: old null check patterns
    public void ProcessUser(User user)
    {
        if (user is null) return;
        if (user.Name is not null)
        {
            Console.WriteLine(user.Name);
        }
        var email = user.Email;
        if (email is null)
        {
            email = "no-reply@example.com";
        }
    }
    // BAD: old type check pattern
    public void HandleEntity(object entity)
    {
        if (entity is User)
        {
            var user = entity as User;
            Console.WriteLine(user.Name);
        }
    }
    // BAD: Verbose property getter/setter
    private string _status;
    public string Status
    {
        get
        {
            return _status;
        }
        set
        {
            _status = value;
        }
    }
    // BAD: default(T) instead of default literal
    public User GetDefaultUser()
    {
        User user = default;
        int count = default;
        string name = default;
        return user;
    }
    // BAD: new List<T> { } instead of collection expression
    public List<string> GetDefaultRoles()
    {
        var roles = new List<string> { "user", "guest" };
        var permissions = ["read", "write"];
        return roles;
    }
    // BAD: Traditional namespace style (file-scoped is better)
    // Already using block-scoped namespace in this file
    // BAD: Manual null propagation
    public string GetUserCity(User user)
    {
        if (user is not null)
        {
            if (user.Address is not null)
            {
                return user.Address.City;
            }
        }
        return null;
    }
    // BAD: is null check with verbose throw
    public void ValidateInput(string input)
    {
        ArgumentNullException.ThrowIfNull(input);Console.WriteLine(input);
    }
    private List<User> GetAllUsers() => new List<User>();
    private Task<User> CallExternalApiAsync(int id) => Task.FromResult<User>(null);
    private void DoExpensiveOperation() { }
}
// Supporting classes
public class User
{
    public int Id { get; set; }
    public required string Name { get; set; }
    public required string Email { get; set; }
    public required Address Address { get; set; }
    public bool IsActive { get; set; }
}
public class Address
{
    public required string City { get; set; }
    public required string Street { get; set; }
}
public class UserEventArgs : EventArgs
{
    public required User User { get; set; }
}
// Interfaces
public interface IDatabase { }
public interface ILogger 
{ 
    void Log(string message);
}
public interface ICache { }