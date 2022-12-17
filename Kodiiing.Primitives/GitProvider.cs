namespace Kodiiing.Primitives;

/// <summary>
/// Specifies the git provider the user signed up from.
/// It is also used for specifying user's profile URL that resides
/// on the upstream provider. For example: https://github.com/elianiva for
/// user elianiva that have GitProvider of Github.
/// </summary>
public enum GitProvider
{
    Github,
    Gitlab
}