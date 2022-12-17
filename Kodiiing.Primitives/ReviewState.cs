namespace Kodiiing.Primitives;

public enum ReviewState
{
    /// <summary>
    /// Indicate the user does not want their solution
    /// to be reviewed by anyone.
    /// </summary>
    DoNotReview,
    
    /// <summary>
    /// Indicate the user want their solution to be reviewed.
    /// But it's not yet reviewed by anyone.
    /// </summary>
    ReviewRequested,
    
    /// <summary>
    /// Indicate the solution has been reviewed.
    /// </summary>
    Reviewed
}