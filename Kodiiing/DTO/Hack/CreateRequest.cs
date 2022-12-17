using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace Kodiiing.DTO.Hack;

public class CreateRequest
{
    [Required, JsonPropertyName("title")]
    public string? Title { get; set; }
    
    [Required, JsonPropertyName("text")]
    public string? Text { get; set; }
    
    [Required, JsonPropertyName("tags")]
    public IEnumerable<string>? Tags { get; set; }
}