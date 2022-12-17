
using System.Text.Json.Serialization;

namespace Kodiiing.DTO.Common;

internal class CommonError
{
    [JsonPropertyName("message")]
    public string Message { get; }

    public CommonError(string message)
    {
        Message = message;
    }
}
