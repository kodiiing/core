// Hack is a mini-forum discussion about a topic. It is encouraged
// to the users (and reviewers) to give their tips-and-trick,
// or other interesting thoughts that are best expressed
// in a short text manner. Imagine this like Hacker News
// (the news.ycombinator.com) or Reddit, but combine it
// with Twitter.
package hack_stub

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HackServiceError struct {
	StatusCode int
	Error      error
}

type CreateRequest struct {
	Title string         `json:"title"`
	Text  string         `json:"text"`
	Tags  []string       `json:"tags"`
	Auth  Authentication `json:"auth"`
}

type CreateResponse struct {
	Id string `json:"id"`
}

type UpvoteRequest struct {
	Id   string         `json:"id"`
	Auth Authentication `json:"auth"`
}

type UpvoteResponse struct {
	Upvoted bool  `json:"upvoted"`
	Score   int64 `json:"score"`
}

type CommentRequest struct {
	HackId   string         `json:"hack_id"`
	ParentId string         `json:"parent_id"`
	Text     string         `json:"text"`
	Auth     Authentication `json:"auth"`
}

type CommentResponse struct {
	HackId    string `json:"hack_id"`
	CommentId string `json:"comment_id"`
}

type ListRequest struct {
	Page   uint32       `json:"page"`
	SortBy SortCriteria `json:"sort_by"`
}

type ListResponse struct {
	Hacks        []Hack `json:"hacks"`
	TotalResults uint64 `json:"total_results"`
	CurrentPage  uint32 `json:"current_page"`
	TotalPage    uint32 `json:"total_page"`
}

type Authentication struct {
	AccessToken string `json:"access_token"`
}

type Author struct {
	Name       string `json:"name"`
	ProfileUrl string `json:"profile_url"`
	PictureUrl string `json:"picture_url"`
}

type Hack struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Upvotes   int64     `json:"upvotes"`
	Author    Author    `json:"author"`
	Comments  []Comment `json:"comments"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type Comment struct {
	Id        string    `json:"id"`
	Content   string    `json:"content"`
	Author    Author    `json:"author"`
	Replies   []Comment `json:"replies"`
	CreatedAt string    `json:"created_at"`
}

type SortCriteria uint32

const (
	SortCriteriaScore        SortCriteria = 0
	SortCriteriaComments     SortCriteria = 1
	SortCriteriaCreated_date SortCriteria = 2
	SortCriteriaUpdated_date SortCriteria = 3
)

type HackServiceServer interface {
	// Starts a new hack post.
	Create(ctx context.Context, req *CreateRequest) (*CreateResponse, *HackServiceError)
	// Upvote a hack post.
	Upvote(ctx context.Context, req *UpvoteRequest) (*UpvoteResponse, *HackServiceError)
	// Comment to a hack post, or reply to an existing comment.
	Comment(ctx context.Context, req *CommentRequest) (*CommentResponse, *HackServiceError)
	// See all hack posts, or maybe with a filter.
	List(ctx context.Context, req *ListRequest) (*ListResponse, *HackServiceError)
}

func NewHackServiceServer(implementation HackServiceServer) *chi.Mux {
	mux := chi.NewMux()
	mux.Post("/Create", func(w http.ResponseWriter, r *http.Request) {
		var req CreateRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[HackService - Createerror] writing to response stream: %s", e.Error())
			}
			return
		}
		resp, err := implementation.Create(r.Context(), &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[HackService - Createerror] writing to response stream: %s", e.Error())
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[HackService - Createerror] writing to response stream: %s", e.Error())
		}
	})

	mux.Post("/Upvote", func(w http.ResponseWriter, r *http.Request) {
		var req UpvoteRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[HackService - Upvoteerror] writing to response stream: %s", e.Error())
			}
			return
		}
		resp, err := implementation.Upvote(r.Context(), &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[HackService - Upvoteerror] writing to response stream: %s", e.Error())
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[HackService - Upvoteerror] writing to response stream: %s", e.Error())
		}
	})

	mux.Post("/Comment", func(w http.ResponseWriter, r *http.Request) {
		var req CommentRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[HackService - Commenterror] writing to response stream: %s", e.Error())
			}
			return
		}
		resp, err := implementation.Comment(r.Context(), &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[HackService - Commenterror] writing to response stream: %s", e.Error())
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[HackService - Commenterror] writing to response stream: %s", e.Error())
		}
	})

	mux.Post("/List", func(w http.ResponseWriter, r *http.Request) {
		var req ListRequest
		e := json.NewDecoder(r.Body).Decode(&req)
		if e != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": e.Error(),
			})
			if e != nil {
				log.Printf("[HackService - Listerror] writing to response stream: %s", e.Error())
			}
			return
		}
		resp, err := implementation.List(r.Context(), &req)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(err.StatusCode)
			e := json.NewEncoder(w).Encode(map[string]string{
				"message": err.Error.Error(),
			})
			if e != nil {
				log.Printf("[HackService - Listerror] writing to response stream: %s", e.Error())
			}
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		e = json.NewEncoder(w).Encode(resp)
		if e != nil {
			log.Printf("[HackService - Listerror] writing to response stream: %s", e.Error())
		}
	})

	return mux
}
