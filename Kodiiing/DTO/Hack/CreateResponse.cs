using System.Text.Json.Serialization;

namespace Kodiiing.DTO.Hack;

public class CreateResponse
{
    [JsonPropertyName("id")]
    public string? Id { get; set; }
}