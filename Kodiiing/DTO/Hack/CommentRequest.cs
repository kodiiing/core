using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace Kodiiing.DTO.Hack;

public class CommentRequest
{
    [Required, JsonPropertyName("hack_id")]
    public string? HackId { get; set; }
    [Required, JsonPropertyName("parent_id")]
    public string? ParentId { get; set; }
    [Required, JsonPropertyName("text")]
    public string? Text { get; set; }
}