package hack_service

import (
	"context"
	"database/sql"

	hack_stub "kodiiing/hack/stub"

	"github.com/typesense/typesense-go/typesense"
)

type HackService struct {
	environment string
	db          *sql.DB
	search      *typesense.Client
}

func NewHackService(env string, db *sql.DB, search *typesense.Client) *HackService {
	return &HackService{environment: env, db: db, search: search}
}

// Starts a new hack post.
func (d *HackService) Create(ctx context.Context, req *hack_stub.CreateRequest) (*hack_stub.CreateResponse, *hack_stub.HackServiceError) {
	return &hack_stub.CreateResponse{}, nil
}

// Upvote a hack post.
func (d *HackService) Upvote(ctx context.Context, req *hack_stub.UpvoteRequest) (*hack_stub.UpvoteResponse, *hack_stub.HackServiceError) {
	return &hack_stub.UpvoteResponse{}, nil
}

// Comment to a hack post, or reply to an existing comment.
func (d *HackService) Comment(ctx context.Context, req *hack_stub.CommentRequest) (*hack_stub.CommentResponse, *hack_stub.HackServiceError) {
	return &hack_stub.CommentResponse{}, nil
}

// See all hack posts, or maybe with a filter.
func (d *HackService) List(ctx context.Context, req *hack_stub.ListRequest) (*hack_stub.ListResponse, *hack_stub.HackServiceError) {
	return &hack_stub.ListResponse{}, nil
}
