using Kodiiing.Models.Authentication;
using Kodiiing.Models.User;
using Kodiiing.Primitives;

namespace Kodiiing.ServiceInterfaces;

public interface IAuthenticationService
{
    Task<UserAggregate> Authenticate(string accessToken, CancellationToken cancellationToken);
    Task<JWT> Login(GitProvider gitProvider, string accessCode, CancellationToken cancellationToken);
    Task Logout(string accessToken);
}