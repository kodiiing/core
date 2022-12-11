using Kodiiing.Models.Authentication;
using Kodiiing.Models.User;
using Kodiiing.Primitives;

namespace Kodiiing.ServiceInterfaces
{
    public interface IAuthenticationService
    {
        Task<User> AuthenticateAsync(string accessToken, CancellationToken cancellationToken);
        Task<JWT> LoginAsync(GitProvider gitProvider, string accessCode, CancellationToken cancellationToken);
        Task LogoutAsync(string accessToken);
    }
}