
using System.Text.Json.Serialization;

namespace Kodiiing.DTO.Common;

internal sealed class CommonError
{
    [JsonPropertyName("message")]
    public required string Message { get; set; }
}
