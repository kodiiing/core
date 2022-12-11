using System.Text.Json.Serialization;

namespace Kodiiing.Primitives;

public class Authentication
{
    [JsonPropertyName("access_token")]
    public string? AccessToken { get; set; }
}