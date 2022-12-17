using Kodiiing.Models.User;
using Kodiiing.Primitives;

namespace Kodiiing.ServiceInterfaces;

public interface IUserService
{
    Task Onboarding(JoinReason joinReason, string? reasonOther, bool codedBefore,
        IEnumerable<string> familiarLanguages, string target, CancellationToken cancellationToken);

    Task<UserAggregate> GetUserByIdAndProvider(string id, GitProvider gitProvider, CancellationToken cancellationToken);
}