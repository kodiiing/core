
using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;
using Kodiiing.Primitives;

namespace Kodiiing.DTO.Auth;

public sealed class LoginRequest
{
    [Required, JsonPropertyName("provider")]
    public GitProvider Provider { get; set; }
    
    [Required, JsonPropertyName("access_code")]
    public string? AccessCode { get; set; }
}