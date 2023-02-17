using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;
using Kodiiing.Primitives;

namespace Kodiiing.DTO.Hack;

public sealed class UpvoteRequest : Authentication
{
    [Required, JsonPropertyName("id")]
    public string? Id { get; set; }
}