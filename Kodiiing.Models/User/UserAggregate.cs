using Kodiiing.Primitives;

namespace Kodiiing.Models.User
{
    public class UserAggregate
    {
        public required int Id { get; set; }
        public required GitProvider GitProvider { get; set; }
        public required string NodeId { get; set; }
        public required string Name { get; set; }
        public required string Username { get; set; }
        public required Uri AvatarUrl { get; set; }
        public required Uri ProfileUrl { get; set; }
        public required string Location { get; set; }
        public required string Email { get; set; }
        public required int PublicRepository { get; set; }
        public required int Followers { get; set; }
        public required int Following { get; set; }
        public required DateTimeOffset GitProviderRegisteredAt { get; set; }
        public required DateTimeOffset KodiiingRegisteredAt { get; set; }
    }
}