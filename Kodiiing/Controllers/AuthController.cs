using System.Net.Mime;
using Kodiiing.DTO.Auth;
using Kodiiing.DTO.Common;
using Kodiiing.Models.Authentication;
using Kodiiing.ServiceInterfaces;
using Microsoft.AspNetCore.Mvc;

namespace Kodiiing.Controllers;

[ApiController]
[Route("/auth")]
public class AuthController : ControllerBase
{
    private readonly IAuthenticationService _authenticationService;

    public AuthController(IAuthenticationService authenticationService)
    {
        _authenticationService = authenticationService;
    }

    [Route("/login")]
    [HttpPost]
    [Consumes(MediaTypeNames.Application.Json)]
    [ProducesResponseType(StatusCodes.Status200OK, Type = typeof(JWT))]
    [ProducesResponseType(StatusCodes.Status400BadRequest, Type = typeof(CommonError))]
    public async Task<IActionResult> LoginAsync([FromBody] LoginRequest loginRequest, CancellationToken cancellationToken)
    {
        if (loginRequest.AccessCode == null) return BadRequest(new CommonError("AccessCode is required"));

        JWT token = await _authenticationService.LoginAsync(loginRequest.Provider, loginRequest.AccessCode,
            cancellationToken);

        return Ok(token);
    }
    
    
}