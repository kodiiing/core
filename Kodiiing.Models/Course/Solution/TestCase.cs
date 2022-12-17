namespace Kodiiing.Models.Course.Solution;

public class TestCase
{
    public required string Input { get; set; }
    public required string Expected { get; set; }
    public required string Output { get; set; }
    public required bool Success { get; set; }
    public required bool Hidden { get; set; }
}