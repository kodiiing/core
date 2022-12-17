using Kodiiing.Models.Course.Task;
using Kodiiing.Models.User;
using Kodiiing.Primitives;

namespace Kodiiing.Models.Course.Satisfaction
{
    public class Satisfaction
    {
        public required TaskAggregate TaskAggregate { get; set; }
        public required UserAggregate UserAggregate { get; set; }
        public required Rating SatisfactionLevel { get; set; }
        public string? Comment { get; set; }
    }
}