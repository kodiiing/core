using System.Text.Json.Serialization;

namespace Kodiiing.DTO.Hack;

public sealed class UpvoteResponse
{
    [JsonPropertyName("voted")]
    public required bool Voted { get; set; }
    [JsonPropertyName("score")]
    public required int Score { get; set; }
}