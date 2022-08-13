package codereview_service

import (
	"context"
	"database/sql"

	codereview_stub "kodiiing/codereview/stub"
)

type CodeReviewService struct {
	db          *sql.DB
	environment string
}

func NewCodeReviewService(env string, db *sql.DB) *CodeReviewService {
	return &CodeReviewService{
		environment: env,
		db:          db,
	}
}

func (d *CodeReviewService) GetAvailableTaskToReview(ctx context.Context, req *codereview_stub.AvailableTaskToReviewRequest) (*codereview_stub.AvailableTaskToReviewResponse, *codereview_stub.CodeReviewServiceError) {
	return &codereview_stub.AvailableTaskToReviewResponse{}, nil
}

func (d *CodeReviewService) SubmitTaskReview(ctx context.Context, req *codereview_stub.SubmitTaskReviewRequest) (*codereview_stub.SubmitTaskReviewResponse, *codereview_stub.CodeReviewServiceError) {
	return &codereview_stub.SubmitTaskReviewResponse{}, nil
}

func (d *CodeReviewService) SubmitReviewComment(ctx context.Context, req *codereview_stub.SubmitReviewCommentRequest) (*codereview_stub.SubmitReviewCommentResponse, *codereview_stub.CodeReviewServiceError) {
	return &codereview_stub.SubmitReviewCommentResponse{}, nil
}

func (d *CodeReviewService) ApplyAsReviewer(ctx context.Context, req *codereview_stub.ApplyAsReviewerRequest) (*codereview_stub.EmptyResponse, *codereview_stub.CodeReviewServiceError) {
	return &codereview_stub.EmptyResponse{}, nil
}
