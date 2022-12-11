using Kodiiing.Models.Authentication;
using Kodiiing.Models.User;
using Kodiiing.Primitives;

namespace Kodiiing.RepositoryInterfaces;

public interface IAuthenticationRepository
{
    string GenerateToken(string subject, int expiresIn);
    bool ValidateToken(string token);
}