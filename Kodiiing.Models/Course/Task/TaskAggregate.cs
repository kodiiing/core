using Kodiiing.Primitives;

namespace Kodiiing.Models.Course.Task
{
    public class TaskAggregate
    {
        public required Guid Id { get; set; }
        public required string Title { get; set; }
        public required string Description { get; set; }
        public required Difficulty Difficulty { get; set; }
        public required Author Author { get; set; }
        public required string Prompt { get; set; }
        public required IEnumerable<Implementation> Implementations { get; set; }
    }
}