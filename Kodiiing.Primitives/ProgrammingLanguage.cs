namespace Kodiiing.Primitives;

public class ProgrammingLanguage
{
    public readonly string Language;

    private readonly IEnumerable<string> _validLanguages = new[] { "C#", "JavaScript", "Go" };
    
    public ProgrammingLanguage(string language)
    {
        if (!_validLanguages.Contains(language)) throw new ArgumentException("Invalid language", nameof(language));

        Language = language;
    }
};
