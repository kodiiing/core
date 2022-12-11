namespace Kodiiing.Primitives
{
    public enum SolutionState
    {
        /// <summary>
        /// Started means the user has opened the coding page.
        /// </summary>
        Started,
        
        /// <summary>
        /// Ongoing means the user has submitted an execute test.
        /// If the user left the coding page and haven't finished the
        /// test yet, they will stay at this stage.
        /// </summary>
        Ongoing,
        
        /// <summary>
        /// The coding test is finished with the correct result.
        /// </summary>
        Finished
    }
}