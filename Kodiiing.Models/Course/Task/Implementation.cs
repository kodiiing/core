using Kodiiing.Primitives;

namespace Kodiiing.Models.Course.Task
{
    public abstract class Implementation
    {
        public required ProgrammingLanguage Language { get; set; }
        public required string Placeholder { get; set; }
        public required string TestCode { get; set; }
    }
}