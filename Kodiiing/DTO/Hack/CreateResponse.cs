using System.Text.Json.Serialization;

namespace Kodiiing.DTO.Hack;

public sealed class CreateResponse
{
    [JsonPropertyName("id")]
    public required string Id { get; set; }
}