using System.Text.Json.Serialization;

namespace Kodiiing.DTO.Auth;

public sealed class LoginResponse
{
    [JsonPropertyName("access_token")]
    public required string AccessToken { get; init; }
    
    [JsonPropertyName("refresh_token")]
    public required string RefreshToken { get; init; }

    [JsonPropertyName("expires_in")] 
    public required int ExpiresIn { get; init; }
    
    [JsonPropertyName("token_type")]
    public required string TokenType { get; set; }
}