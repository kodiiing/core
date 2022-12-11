using Kodiiing.Primitives;

namespace Kodiiing.Models.Course.Satisfaction
{
    public class Satisfaction
    {
        public required Task.Task Task { get; set; }
        public required User.User User { get; set; }
        public required Rating SatisfactionLevel { get; set; }
        public string? Comment { get; set; }
    }
}