using Kodiiing.Primitives;

namespace Kodiiing.Models.Course.CodeReview
{
    public class Conversation
    {
        public required Guid Id { get; set; }
        public required Author Author { get; set; }
        public required string Content { get; set; }
        public required DateTimeOffset CreatedAt { get; set; }
        public required DateTimeOffset UpdatedAt { get; set; }
    }
}