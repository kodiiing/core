package task_stub

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type TaskServiceError struct {
	StatusCode int
	Error      error
}

type EmptyResponse struct {
}

type Authentication struct {
	AccessToken string `json:"access_token"`
}

type TaskDifficulty int64

const (
	TASK_DIFFICULTY_UNSPECIFIED TaskDifficulty = 0
	TASK_DIFFICULTY_EASY        TaskDifficulty = 1
	TASK_DIFFICULTY_MEDIUM      TaskDifficulty = 2
	TASK_DIFFICULTY_HARD        TaskDifficulty = 3
)

type Task struct {
	Id                 string         `json:"id"`
	Title              string         `json:"title"`
	Description        string         `json:"description"`
	Difficulty         TaskDifficulty `json:"difficulty"`
	Completed          bool           `json:"completed"`
	Content            string         `json:"content"`
	Author             string         `json:"author"`
	CompletedAt        string         `json:"completed_at"`
	StatisfactionLevel int32          `json:"statisfaction_level"`
}

type ListTaskRequest struct {
	Auth    Authentication `json:"auth"`
	TrackId string         `json:"track_id"`
}

type ListTaskResponse struct {
	Tasks []Task `json:"tasks"`
}

type StartTaskRequest struct {
	Auth   Authentication `json:"auth"`
	TaskId string         `json:"task_id"`
}

type StartTaskResponse struct {
	Task Task `json:"task"`
}

type PostTaskAssessmentRequest struct {
	Auth              Authentication `json:"auth"`
	TaskId            string         `json:"task_id"`
	SatisfactionLevel int32          `json:"satisfaction_level"`
	Comments          string         `json:"comments"`
}

type TaskServiceServer interface {
	// List all task that is available by a certain track ID
	ListTasks(ctx context.Context, req ListTaskRequest) (*ListTaskResponse, *TaskServiceError)

	// Starts a task, will marks the task as "ongoing" when viewed by the current user.
	StartTask(ctx context.Context, req StartTaskRequest) (*StartTaskResponse, *TaskServiceError)

	// Give an assessment to the user about the task, whether they are happy with it or they
	// don't like the given task.
	PostTaskAssessment(ctx context.Context, req PostTaskAssessmentRequest) (*EmptyResponse, *TaskServiceError)
}

func NewTaskServiceServer(implementation TaskServiceServer) *chi.Mux {
	mux := chi.NewMux()
	mux.Post("/List", func(w http.ResponseWriter, r *http.Request) {
		var req ListTaskRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[TaskService - ListError] writing to response stream: %s", e.Error())
			}
			return
		}

		resp, err := implementation.ListTasks(r.Context(), req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[TaskService - ListError] writing to response stream: %s", e.Error())
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[TaskService - ListError] writing to response stream: %s", e.Error())
		}
	})

	mux.Post("/Start", func(w http.ResponseWriter, r *http.Request) {
		var req StartTaskRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[TaskService - ListError] writing to response stream: %s", e.Error())
			}
			return
		}

		resp, err := implementation.StartTask(r.Context(), req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[TaskService - ListError] writing to response stream: %s", e.Error())
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[TaskService - ListError] writing to response stream: %s", e.Error())
		}
	})

	mux.Post("/Assessment", func(w http.ResponseWriter, r *http.Request) {
		var req PostTaskAssessmentRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[TaskService - ListError] writing to response stream: %s", e.Error())
			}
			return
		}

		resp, err := implementation.PostTaskAssessment(r.Context(), req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[TaskService - ListError] writing to response stream: %s", e.Error())
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[TaskService - ListError] writing to response stream: %s", e.Error())
		}
	})

	return mux
}
