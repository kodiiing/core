using Kodiiing.Models.Hack;
using Kodiiing.Models.User;

namespace Kodiiing.ServiceInterfaces;

public interface IHackService
{
    Task<HackAggregate> Create(UserAggregate user, string title, string content, IEnumerable<string> tags,
        CancellationToken cancellationToken);
    Task<int> Upvote(UserAggregate user, Guid id, CancellationToken cancellationToken);
    Task<HackAggregate> CreateComment(UserAggregate user, Guid hackId, Guid parentId, string content, CancellationToken cancellationToken);
    Task<IEnumerable<HackAggregate>> List(UserAggregate user, CancellationToken cancellationToken);
}