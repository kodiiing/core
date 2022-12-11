using Kodiiing.Primitives;

namespace Kodiiing.Models.User
{
    public class Repository
    {
        public required int Id { get; set; }
        public required GitProvider Provider { get; set; }
        public required string Name { get; set; }
        public required Uri URL { get; set; }
        public string? Description { get; set; }
        public bool Fork { get; set; }
        public int ForksCount { get; set; }
        public int StarsCount { get; set; }
        public required string OwnerUsername { get; set; }
        public required  DateTimeOffset CreatedAt { get; set; }
        public required DateTimeOffset LastActivityAt { get; set; }
    }
}