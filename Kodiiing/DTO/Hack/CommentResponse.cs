using System.Text.Json.Serialization;

namespace Kodiiing.DTO.Hack;

public class CommentResponse
{
    [JsonPropertyName("hack_id")]
    public string? HackId { get; set; }
    [JsonPropertyName("comment_id")]
    public string? CommentId { get; set; }
}