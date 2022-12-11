using System.Text.Json.Serialization;

namespace Kodiiing.DTO.Auth;

public class LoginResponse
{
    [JsonPropertyName("access_token")]
    public string? AccessToken { get; set; }
}