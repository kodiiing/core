using Kodiiing.Primitives;

namespace Kodiiing.Models.Hack
{
    public class Hack
    {
        public required Guid Id { get; set; }
        public required string Title { get; set; }
        public required string Content { get; set; }
        public int Upvote { get; set; }
        public required Author Author { get; set; }
        public IEnumerable<Comment>? Comments { get; set; }
        public required DateTimeOffset CreatedAt { get; set; }
        public required DateTimeOffset UpdatedAt { get; set; }
    }
}