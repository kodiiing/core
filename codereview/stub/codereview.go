package codereview_stub

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CodeReviewServiceError struct {
	StatusCode int
	Error      error
}

type AvailableTaskToReviewRequest struct {
	Auth Authentication `json:"auth"`
}

type AvailableTaskToReviewResponse struct {
	TaskAnswers []TaskAnswer `json:"task_answers"`
}

type SubmitTaskReviewRequest struct {
	Auth         Authentication `json:"auth"`
	TaskAnswerId string         `json:"task_answer_id"`
	Content      string         `json:"content"`
}

type SubmitTaskReviewResponse struct {
	TaskAnswerId string     `json:"task_answer_id"`
	Feedback     []Feedback `json:"feedback"`
}

type SubmitReviewCommentRequest struct {
	Auth           Authentication `json:"auth"`
	TaskAnswerId   string         `json:"task_answer_id"`
	FeedbackId     string         `json:"feedback_id"`
	ConversationId string         `json:"conversation_id"`
	Content        string         `json:"content"`
}

type SubmitReviewCommentResponse struct {
	TaskAnswerId string     `json:"task_answer_id"`
	Feedback     []Feedback `json:"feedback"`
}

type ApplyAsReviewerRequest struct {
	Auth Authentication `json:"auth"`
}

type EmptyResponse struct {
}

type Authentication struct {
	AccessToken string `json:"access_token"`
}

type Task struct {
	Id                string `json:"id"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	Difficulty        string `json:"difficulty"`
	Completed         bool   `json:"completed"`
	Content           string `json:"content"`
	Author            string `json:"author"`
	CompletedAt       string `json:"completed_at"`
	SatisfactionLevel int32  `json:"satisfaction_level"`
}

type Author struct {
	Name       string `json:"name"`
	ProfileUrl string `json:"profile_url"`
	PictureUrl string `json:"picture_url"`
}

type TaskAnswer struct {
	Id          string `json:"id"`
	UserId      string `json:"user_id"`
	UserName    string `json:"user_name"`
	SubmittedAt string `json:"submitted_at"`
	Content     string `json:"content"`
	Task        Task   `json:"task"`
}

type Conversation struct {
	Id        string `json:"id"`
	Author    Author `json:"author"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

type Feedback struct {
	Id            string         `json:"id"`
	Author        Author         `json:"author"`
	Content       string         `json:"content"`
	Conversations []Conversation `json:"conversations"`
	CreatedAt     string         `json:"created_at"`
}

type CodeReviewServiceServer interface {
	GetAvailableTaskToReview(ctx context.Context, req *AvailableTaskToReviewRequest) (*AvailableTaskToReviewResponse, *CodeReviewServiceError)
	SubmitTaskReview(ctx context.Context, req *SubmitTaskReviewRequest) (*SubmitTaskReviewResponse, *CodeReviewServiceError)
	SubmitReviewComment(ctx context.Context, req *SubmitReviewCommentRequest) (*SubmitReviewCommentResponse, *CodeReviewServiceError)
	ApplyAsReviewer(ctx context.Context, req *ApplyAsReviewerRequest) (*EmptyResponse, *CodeReviewServiceError)
}

func NewCodeReviewServiceServer(implementation CodeReviewServiceServer) *chi.Mux {
	mux := chi.NewMux()
	mux.Post("/GetAvailableTaskToReview", func(w http.ResponseWriter, r *http.Request) {
		var req AvailableTaskToReviewRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[CodeReviewService - GetAvailableTaskToReviewerror] writing to response stream: %s", e.Error())
			}
			return
		}
		resp, err := implementation.GetAvailableTaskToReview(r.Context(), &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[CodeReviewService - GetAvailableTaskToReviewerror] writing to response stream: %s", e.Error())
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[CodeReviewService - GetAvailableTaskToReviewerror] writing to response stream: %s", e.Error())
		}
	})

	mux.Post("/SubmitTaskReview", func(w http.ResponseWriter, r *http.Request) {
		var req SubmitTaskReviewRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[CodeReviewService - SubmitTaskReviewerror] writing to response stream: %s", e.Error())
			}
			return
		}
		resp, err := implementation.SubmitTaskReview(r.Context(), &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[CodeReviewService - SubmitTaskReviewerror] writing to response stream: %s", e.Error())
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[CodeReviewService - SubmitTaskReviewerror] writing to response stream: %s", e.Error())
		}
	})

	mux.Post("/SubmitReviewComment", func(w http.ResponseWriter, r *http.Request) {
		var req SubmitReviewCommentRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[CodeReviewService - SubmitReviewCommenterror] writing to response stream: %s", e.Error())
			}
			return
		}
		resp, err := implementation.SubmitReviewComment(r.Context(), &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[CodeReviewService - SubmitReviewCommenterror] writing to response stream: %s", e.Error())
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[CodeReviewService - SubmitReviewCommenterror] writing to response stream: %s", e.Error())
		}
	})

	mux.Post("/ApplyAsReviewer", func(w http.ResponseWriter, r *http.Request) {
		var req ApplyAsReviewerRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[CodeReviewService - ApplyAsReviewererror] writing to response stream: %s", e.Error())
			}
			return
		}
		resp, err := implementation.ApplyAsReviewer(r.Context(), &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[CodeReviewService - ApplyAsReviewererror] writing to response stream: %s", e.Error())
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[CodeReviewService - ApplyAsReviewererror] writing to response stream: %s", e.Error())
		}
	})

	return mux
}
