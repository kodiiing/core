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

type TaskServiceServer interface {
	// List all task that is available by a certain track ID
	ListTasks(ctx context.Context, req ListTaskRequest) (*ListTaskResponse, *TaskServiceError)
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

	return mux
}