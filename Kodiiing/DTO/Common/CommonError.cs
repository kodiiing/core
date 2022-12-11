
using System.Text.Json.Serialization;

namespace Kodiiing.DTO.Common
{
    public class CommonError
    {
        [JsonPropertyName("message")]
        public string Message { get; set; }

        public CommonError(string message)
        {
            Message = message;
        }
    }
}