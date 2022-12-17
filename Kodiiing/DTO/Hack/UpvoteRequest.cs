using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace Kodiiing.DTO.Hack;

public class UpvoteRequest
{
    [Required, JsonPropertyName("id")]
    public string? Id { get; set; }
}