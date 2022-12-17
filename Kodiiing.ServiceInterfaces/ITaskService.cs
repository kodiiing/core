using Kodiiing.Models.Course.Satisfaction;
using Kodiiing.Models.Course.Solution;
using Kodiiing.Models.Course.Task;
using Kodiiing.Models.User;

namespace Kodiiing.ServiceInterfaces;

public interface ITaskService
{
    /// <summary>
    /// List all task that is available by a certain track ID
    /// </summary>
    /// <param name="user"></param>
    /// <param name="trackId"></param>
    /// <param name="cancellationToken"></param>
    /// <returns></returns>
    Task<IEnumerable<TaskAggregate>> List(UserAggregate? user, Guid trackId, CancellationToken cancellationToken);
    
    /// <summary>
    /// Starts a task, will marks the task as "ongoing" when viewed by the current user.
    /// </summary>
    /// <param name="user"></param>
    /// <param name="task"></param>
    /// <param name="cancellationToken"></param>
    /// <returns></returns>
    Task Start(UserAggregate user, TaskAggregate task, CancellationToken cancellationToken);
    
    /// <summary>
    /// Executes a code that resides on task if it's a coding task. Will return a test cases result.
    /// </summary>
    /// <param name="user"></param>
    /// <param name="task"></param>
    /// <param name="code"></param>
    /// <param name="cancellationToken"></param>
    /// <returns></returns>
    Task<TestResult> ExecuteCode(UserAggregate user, TaskAggregate task, string code, CancellationToken cancellationToken);

    /// <summary>
    /// Submit a task as a final submission, no more changes after this one.
    /// This should be called after StartTask rpc was called.
    /// </summary>
    /// <param name="user"></param>
    /// <param name="task"></param>
    /// <param name="code"></param>
    /// <param name="cancellationToken"></param>
    /// <returns></returns>
    Task Submit(UserAggregate user, TaskAggregate task, string code, CancellationToken cancellationToken);
    
    /// <summary>
    /// Give an assessment to the user about the task, whether they are happy with it
    /// or they don't like the given task.
    /// </summary>
    /// <param name="user"></param>
    /// <param name="task"></param>
    /// <param name="satisfaction"></param>
    /// <param name="comments"></param>
    /// <param name="cancellationToken"></param>
    /// <returns></returns>
    Task PostTaskAssessment(UserAggregate user, TaskAggregate task, Satisfaction satisfaction, string? comments,
        CancellationToken cancellationToken);

    /// <summary>
    /// Submit task feedback from the user who did the task.
    /// </summary>
    /// <param name="user"></param>
    /// <param name="task"></param>
    /// <param name="feedback"></param>
    /// <param name="cancellationToken"></param>
    /// <returns></returns>
    Task SubmitTaskFeedback(UserAggregate user, TaskAggregate task, string feedback,
        CancellationToken cancellationToken);
}