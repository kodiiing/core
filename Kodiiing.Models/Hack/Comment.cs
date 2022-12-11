using Kodiiing.Primitives;

namespace Kodiiing.Models.Hack
{
    public class Comment
    {
        public required Guid Id { get; set; }
        public required string Content { get; set; }
        public required Author Author { get; set; }
        public IEnumerable<Comment>? Replies { get; set; }
        public required DateTimeOffset CreatedAt { get; set; }
    }
}