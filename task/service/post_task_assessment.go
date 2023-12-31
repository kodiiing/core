package service

import (
	"context"
	"errors"
	"fmt"
	"kodiiing/auth"
	"kodiiing/task/repository"
	task_stub "kodiiing/task/stub"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func ValidatePostTaskAssessmentReq(req task_stub.PostTaskAssessmentRequest) *task_stub.TaskServiceError {
	_, err := strconv.ParseInt(req.TaskId, 10, 64)
	if err != nil {
		return &task_stub.TaskServiceError{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Errorf("invalid task id"),
		}
	}

	if req.SatisfactionLevel < 0 || req.SatisfactionLevel > 5 {
		return &task_stub.TaskServiceError{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Errorf("satisfaction level must be between 1 and 5"),
		}
	}

	if len(req.Comments) > 511 {
		return &task_stub.TaskServiceError{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Errorf("comments too long"),
		}
	}

	return nil
}

func (s *TaskService) PostTaskAssessment(ctx context.Context, req task_stub.PostTaskAssessmentRequest) (*task_stub.EmptyResponse, *task_stub.TaskServiceError) {
	// authenticate user
	authenticatedUser, err := s.authentication.Authenticate(ctx, req.Auth.AccessToken)
	if err != nil {
		if errors.Is(err, auth.ErrParameterEmpty) || errors.Is(err, auth.ErrUserNotFound) {
			return &task_stub.EmptyResponse{}, &task_stub.TaskServiceError{
				StatusCode: http.StatusUnauthorized,
				Error:      fmt.Errorf("unauthenticated: %w", err),
			}
		}

		return nil, &task_stub.TaskServiceError{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Errorf("authenticating user: %w", err),
		}
	}

	// validate request
	if err := ValidatePostTaskAssessmentReq(req); err != nil {
		return &task_stub.EmptyResponse{}, err
	}

	// insert assessment
	taskId, _ := strconv.ParseInt(req.TaskId, 10, 64)
	affected, err := s.taskRepository.InsertTaskAssessment(ctx, repository.InsertTaskAssessmentIn{
		UserId:            authenticatedUser.ID,
		TaskId:            taskId,
		SatisfactionLevel: int64(req.SatisfactionLevel),
		Comments:          req.Comments,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &task_stub.TaskServiceError{
				StatusCode: http.StatusNotFound,
				Error:      fmt.Errorf("task not found"),
			}
		}

		return nil, &task_stub.TaskServiceError{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}

	if affected == 0 {
		return nil, &task_stub.TaskServiceError{
			StatusCode: http.StatusNotFound,
			Error:      fmt.Errorf("task not found"),
		}
	}

	return &task_stub.EmptyResponse{}, nil
}