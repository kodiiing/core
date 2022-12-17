using System.Text.Json.Serialization;

namespace Kodiiing.DTO.Hack;

public class UpvoteResponse
{
    [JsonPropertyName("voted")]
    public bool Voted { get; set; }
    [JsonPropertyName("score")]
    public int Score { get; set; }
}