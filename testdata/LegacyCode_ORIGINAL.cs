using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Net.Http;
using System.Threading.Tasks;

namespace LegacyCompany.BadCodeExamples
{
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
            if (database == null)
                throw new ArgumentNullException(nameof(database));
            _database = database;

            if (logger == null)
                throw new ArgumentNullException(nameof(logger));
            _logger = logger;

            if (cache == null)
                throw new ArgumentNullException(nameof(cache));
            _cache = cache;
        }

        // BAD: Count() > 0 instead of Any()
        public bool HasUsers()
        {
            var users = GetAllUsers();
            if (users.Count() > 0)
            {
                return true;
            }
            if (users.Count() != 0)
            {
                return true;
            }
            if (users.Count() >= 1)
            {
                return true;
            }
            return false;
        }

        // BAD: Where().First() instead of First(predicate)
        public User FindUserById(int id)
        {
            var users = GetAllUsers();
            var user = users.Where(u => u.Id == id).First();
            var maybeUser = users.Where(u => u.Id == id).FirstOrDefault();
            var singleUser = users.Where(u => u.Id == id).Single();
            var maybeSingle = users.Where(u => u.Id == id).SingleOrDefault();
            return user;
        }

        // BAD: Old temp swap pattern
        public void SwapUsers(ref User a, ref User b)
        {
            var temp = a;
            a = b;
            b = temp;
        }

        // BAD: Verbose event invocation
        public event EventHandler<UserEventArgs> UserCreated;
        public event EventHandler<UserEventArgs> UserDeleted;

        protected virtual void OnUserCreated(UserEventArgs args)
        {
            if (UserCreated != null)
                UserCreated(this, args);
        }

        protected virtual void OnUserDeleted(UserEventArgs args)
        {
            if (UserDeleted != null)
            {
                UserDeleted(this, args);
            }
        }

        // BAD: string.IsNullOrEmpty instead of pattern matching
        public bool ValidateUsername(string username)
        {
            if (string.IsNullOrEmpty(username))
            {
                return false;
            }
            if (!string.IsNullOrEmpty(username))
            {
                return username.Length >= 3;
            }
            return false;
        }

        // BAD: Concat().ToList() instead of spread operator
        public List<User> CombineUserLists(List<User> activeUsers, List<User> inactiveUsers)
        {
            var allUsers = activeUsers.Concat(inactiveUsers).ToList();
            var userArray = activeUsers.Concat(inactiveUsers).ToArray();
            return allUsers;
        }

        // BAD: Block-scoped catch instead of exception filter
        public async Task<User> FetchUserFromApiAsync(int id)
        {
            try
            {
                return await CallExternalApiAsync(id);
            }
            catch (HttpRequestException ex) { if (!ex.Message.Contains("NotFound")) throw;
                _logger.Log("User not found: " + id);
                return null;
            }
        }

        // BAD: Old stopwatch pattern
        public void MeasurePerformance()
        {
            var stopwatch = new Stopwatch();
            stopwatch.Start();

            DoExpensiveOperation();

            stopwatch.Stop();
            Console.WriteLine("Elapsed: " + stopwatch.ElapsedMilliseconds);
        }

        // BAD: String concatenation instead of interpolation
        public string GetUserSummary(User user)
        {
            return "User: " + user.Name + " (ID: " + user.Id + ") - Email: " + user.Email;
        }

        // BAD: Tuple<T1,T2> instead of value tuples
        public Tuple<string, int> GetUserInfo(User user)
        {
            return new Tuple<string, int>(user.Name, user.Id);
        }

        // BAD: old null check patterns
        public void ProcessUser(User user)
        {
            if (user == null) return;
            
            if (user.Name != null)
            {
                Console.WriteLine(user.Name);
            }

            var email = user.Email;
            if (email == null)
            {
                email = "no-reply@example.com";
            }
        }

        // BAD: old type check pattern
        public void HandleEntity(object entity)
        {
            if ((entity as User) != null)
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
            User user = default(User);
            int count = default(int);
            string name = default(string);
            return user;
        }

        // BAD: new List<T> { } instead of collection expression
        public List<string> GetDefaultRoles()
        {
            var roles = new List<string> { "user", "guest" };
            var permissions = new string[] { "read", "write" };
            return roles;
        }

        // BAD: Traditional namespace style (file-scoped is better)
        // Already using block-scoped namespace in this file

        // BAD: Manual null propagation
        public string GetUserCity(User user)
        {
            if (user != null)
            {
                if (user.Address != null)
                {
                    return user.Address.City;
                }
            }
            return null;
        }

        // BAD: is null check with verbose throw
        public void ValidateInput(string input)
        {
            if (input is null)
                throw new ArgumentNullException(nameof(input));
            
            Console.WriteLine(input);
        }

        private List<User> GetAllUsers() => new List<User>();
        private Task<User> CallExternalApiAsync(int id) => Task.FromResult<User>(null);
        private void DoExpensiveOperation() { }
    }

    // Supporting classes
    public class User
    {
        public int Id { get; set; }
        public string Name { get; set; }
        public string Email { get; set; }
        public Address Address { get; set; }
        public bool IsActive { get; set; }
    }

    public class Address
    {
        public string City { get; set; }
        public string Street { get; set; }
    }

    public class UserEventArgs : EventArgs
    {
        public User User { get; set; }
    }

    // Interfaces
    public interface IDatabase { }
    public interface ILogger 
    { 
        void Log(string message);
    }
    public interface ICache { }
}
