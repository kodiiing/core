namespace Kodiiing.Models.Course.Solution;

public class TestResult
{
    public required IEnumerable<TestCase> TestCases { get; set; }
    public required bool AllowedToSubmit { get; set; }
}