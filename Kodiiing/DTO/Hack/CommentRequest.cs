using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;
using Kodiiing.Primitives;

namespace Kodiiing.DTO.Hack;

public sealed class CommentRequest : Authentication
{
    [Required, JsonPropertyName("hack_id")]
    public string? HackId { get; set; }
    [Required, JsonPropertyName("parent_id")]
    public string? ParentId { get; set; }
    [Required, JsonPropertyName("text")]
    public string? Text { get; set; }
}