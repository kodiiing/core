using Kodiiing.Models.Course.CodeReview;
using Kodiiing.Models.Course.Solution;
using Kodiiing.Models.User;

namespace Kodiiing.ServiceInterfaces;

public interface ICodeReviewService
{
    Task ApplyAsReviewer(UserAggregate user, CancellationToken cancellationToken);
    Task<IEnumerable<SolutionAggregate>> GetAvailableTasksToReview(CancellationToken cancellationToken);

    Task<CodeReview> SubmitTaskReview(UserAggregate user, SolutionAggregate solution, Conversation conversation,
        CancellationToken cancellationToken);

    Task<CodeReview> SubmitReviewComment(UserAggregate user, SolutionAggregate solution, Conversation conversation,
        CancellationToken cancellationToken);

}