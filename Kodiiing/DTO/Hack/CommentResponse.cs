using System.Text.Json.Serialization;

namespace Kodiiing.DTO.Hack;

public sealed class CommentResponse
{
    [JsonPropertyName("hack_id")]
    public required string HackId { get; set; }
    [JsonPropertyName("comment_id")]
    public required string CommentId { get; set; }
}