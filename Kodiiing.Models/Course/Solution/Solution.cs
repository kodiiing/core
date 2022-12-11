using Kodiiing.Primitives;

namespace Kodiiing.Models.Course.Solution
{
    public class Solution
    {
        public required Guid Id { get; set; }
        public required Task.Task Task { get; set; }
        public required User.User User { get; set; }
        public required SolutionState State { get; set; }
        public required ProgrammingLanguage ProgrammingLanguage { get; set; }
        public int FailedAttempts { get; set; }
        public int SuccessAttempts { get; set; }
        public string? FinalSolution { get; set; }
        public required DateTimeOffset CreatedAt { get; set; }
        public required DateTimeOffset UpdatedAt { get; set; }
    }
}