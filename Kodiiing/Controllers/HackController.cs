using Kodiiing.ServiceInterfaces;
using Microsoft.AspNetCore.Mvc;

namespace Kodiiing.Controllers;

[ApiController]
[Route("/hack")]
public class HackController
{
    private readonly IHackService _hackService;

    public HackController(IHackService hackService)
    {
        _hackService = hackService;
    }
}