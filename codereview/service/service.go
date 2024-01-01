package codereview_service

import (
	"context"

	codereview_stub "kodiiing/codereview/stub"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CodeReviewService struct {
	pool        *pgxpool.Pool
	environment string
}

func NewCodeReviewService(env string, pool *pgxpool.Pool) codereview_stub.CodeReviewServiceServer {
	return &CodeReviewService{
		pool:        pool,
		environment: env,
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
