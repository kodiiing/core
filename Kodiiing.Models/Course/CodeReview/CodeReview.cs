﻿namespace Kodiiing.Models.Course.CodeReview
{
    public class CodeReview
    {
        public required Guid Id { get; set; }
        public required Solution.Solution Solution { get; set; }
        public required IEnumerable<Conversation> Conversations { get; set; }
        public required DateTimeOffset CreatedAt { get; set; }
        public required DateTimeOffset UpdatedAt { get; set; }
    }
}