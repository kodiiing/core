// Task provides common functionality for tasks.
package task

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
)

type TaskServiceError struct {
	StatusCode int
	Error      error
}

type ListTasksRequest struct {
	Auth    Authentication `json:"auth"`
	TrackId string         `json:"track_id"`
}

type ListTasksResponse struct {
	Tasks []Task `json:"tasks"`
}

type StartTaskRequest struct {
	Auth   Authentication `json:"auth"`
	TaskId string         `json:"task_id"`
}

type StartTaskResponse struct {
	Task Task `json:"task"`
}

type ExecuteCodeRequest struct {
	Auth   Authentication `json:"auth"`
	TaskId string         `json:"task_id"`
	Code   string         `json:"code"`
}

type ExecuteCodeResponse struct {
	Output          string     `json:"output"`
	TestCases       []TestCase `json:"test_cases"`
	AllowedToSubmit bool       `json:"allowed_to_submit"`
}

type SubmitTaskRequest struct {
	Auth       Authentication `json:"auth"`
	TaskId     string         `json:"task_id"`
	Submission string         `json:"submission"`
}

type SubmitTaskResponse struct {
	NextTaskId string `json:"next_task_id"`
}

type PostTaskAssessmentRequest struct {
	Auth              Authentication `json:"auth"`
	TaskId            string         `json:"task_id"`
	SatisfactionLevel int32          `json:"satisfaction_level"`
	Comments          string         `json:"comments"`
}

type EmptyResponse struct {
}

type SubmitTaskFeedbackRequest struct {
	Auth     Authentication `json:"auth"`
	TaskId   string         `json:"task_id"`
	Feedback string         `json:"feedback"`
}

type SubmitTaskFeedbackResponse struct {
	TaskId   string     `json:"task_id"`
	Feedback []Feedback `json:"feedback"`
}

type Authentication struct {
	AccessToken string `json:"access_token"`
}

type Task struct {
	Id                string         `json:"id"`
	Title             string         `json:"title"`
	Description       string         `json:"description"`
	Difficulty        TaskDifficulty `json:"difficulty"`
	Completed         bool           `json:"completed"`
	Content           string         `json:"content"`
	Author            string         `json:"author"`
	CompletedAt       string         `json:"completed_at"`
	SatisfactionLevel int32          `json:"satisfaction_level"`
}

type TestCase struct {
	Input    string `json:"input"`
	Expected string `json:"expected"`
	Output   string `json:"output"`
	Success  bool   `json:"success"`
	Hidden   bool   `json:"hidden"`
}

type Feedback struct {
	AuthorId   string `json:"author_id"`
	AuthorName string `json:"author_name"`
	Content    string `json:"content"`
	Timestamp  int64  `json:"timestamp"`
}

type TaskDifficulty uint32

const (
	TASK_DIFFICULTY_UNSPECIFIED TaskDifficulty = 0
	TASK_DIFFICULTY_EASY        TaskDifficulty = 1
	TASK_DIFFICULTY_MEDIUM      TaskDifficulty = 2
	TASK_DIFFICULTY_HARD        TaskDifficulty = 3
)

var tracer = otel.Tracer("kodiiing/task/stub")

type TaskServiceServer interface {
	// List all task that is available by a certain track ID
	ListTasks(ctx context.Context, req *ListTasksRequest) (*ListTasksResponse, *TaskServiceError)
	// Starts a task, will marks the task as "ongoing" when viewed by the current user.
	StartTask(ctx context.Context, req *StartTaskRequest) (*StartTaskResponse, *TaskServiceError)
	// Executes a code that resides on task if it's a coding task. Will return a test cases result.
	ExecuteCode(ctx context.Context, req *ExecuteCodeRequest) (*ExecuteCodeResponse, *TaskServiceError)
	// Submit a task as a final submission, no more changes after this one.
	// This should be called after StartTask rpc was called.
	SubmitTask(ctx context.Context, req *SubmitTaskRequest) (*SubmitTaskResponse, *TaskServiceError)
	// Give an assessment to the user about the task, whether they are happy with it or they
	// don't like the given task.
	PostTaskAssessment(ctx context.Context, req *PostTaskAssessmentRequest) (*EmptyResponse, *TaskServiceError)
	// Submit task feedback from the user who did the task. For submitting feedback that comes
	// from the code reviewers, see the codereview proto.
	SubmitTaskFeedback(ctx context.Context, req *SubmitTaskFeedbackRequest) (*SubmitTaskFeedbackResponse, *TaskServiceError)
}

func NewTaskServiceServer(implementation TaskServiceServer) *chi.Mux {
	mux := chi.NewMux()
	mux.Post("/ListTasks", func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracer.Start(r.Context(), "Post.ListTasks")
		defer span.End()

		var req ListTasksRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[TaskService - ListTaskserror] writing to response stream: %s", e.Error())
			}
			return
		}
		resp, err := implementation.ListTasks(ctx, &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[TaskService - ListTaskserror] writing to response stream: %s", e.Error())
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[TaskService - ListTaskserror] writing to response stream: %s", e.Error())
		}
	})

	mux.Post("/StartTask", func(w http.ResponseWriter, r *http.Request) {
		var req StartTaskRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[TaskService - StartTaskerror] writing to response stream: %s", e.Error())
			}
			return
		}
		resp, err := implementation.StartTask(r.Context(), &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[TaskService - StartTaskerror] writing to response stream: %s", e.Error())
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[TaskService - StartTaskerror] writing to response stream: %s", e.Error())
		}
	})

	mux.Post("/ExecuteCode", func(w http.ResponseWriter, r *http.Request) {
		var req ExecuteCodeRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[TaskService - ExecuteCodeerror] writing to response stream: %s", e.Error())
			}
			return
		}
		resp, err := implementation.ExecuteCode(r.Context(), &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[TaskService - ExecuteCodeerror] writing to response stream: %s", e.Error())
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[TaskService - ExecuteCodeerror] writing to response stream: %s", e.Error())
		}
	})

	mux.Post("/SubmitTask", func(w http.ResponseWriter, r *http.Request) {
		var req SubmitTaskRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[TaskService - SubmitTaskerror] writing to response stream: %s", e.Error())
			}
			return
		}
		resp, err := implementation.SubmitTask(r.Context(), &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[TaskService - SubmitTaskerror] writing to response stream: %s", e.Error())
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[TaskService - SubmitTaskerror] writing to response stream: %s", e.Error())
		}
	})

	mux.Post("/PostTaskAssessment", func(w http.ResponseWriter, r *http.Request) {
		var req PostTaskAssessmentRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[TaskService - PostTaskAssessmenterror] writing to response stream: %s", e.Error())
			}
			return
		}
		resp, err := implementation.PostTaskAssessment(r.Context(), &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[TaskService - PostTaskAssessmenterror] writing to response stream: %s", e.Error())
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[TaskService - PostTaskAssessmenterror] writing to response stream: %s", e.Error())
		}
	})

	mux.Post("/SubmitTaskFeedback", func(w http.ResponseWriter, r *http.Request) {
		var req SubmitTaskFeedbackRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[TaskService - SubmitTaskFeedbackerror] writing to response stream: %s", e.Error())
			}
			return
		}
		resp, err := implementation.SubmitTaskFeedback(r.Context(), &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[TaskService - SubmitTaskFeedbackerror] writing to response stream: %s", e.Error())
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[TaskService - SubmitTaskFeedbackerror] writing to response stream: %s", e.Error())
		}
	})

	return mux
}
