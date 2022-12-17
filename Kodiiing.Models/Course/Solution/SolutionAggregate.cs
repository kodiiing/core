using Kodiiing.Models.Course.Task;
using Kodiiing.Models.User;
using Kodiiing.Primitives;

namespace Kodiiing.Models.Course.Solution
{
    public class SolutionAggregate
    {
        public required Guid Id { get; set; }
        public required TaskAggregate TaskAggregate { get; set; }
        public required UserAggregate UserAggregate { get; set; }
        public required SolutionState State { get; set; }
        public required ProgrammingLanguage ProgrammingLanguage { get; set; }
        public int FailedAttempts { get; set; }
        public int SuccessAttempts { get; set; }
        public string? FinalSolution { get; set; }
        public ReviewState ReviewState { get; set; }
        public required DateTimeOffset CreatedAt { get; set; }
        public required DateTimeOffset UpdatedAt { get; set; }
    }
}