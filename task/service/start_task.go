package service

import (
	"context"
	"errors"
	"fmt"
	"kodiiing/auth"
	task_stub "kodiiing/task/stub"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

func (s *TaskService) StartTask(ctx context.Context, req *task_stub.StartTaskRequest) (*task_stub.StartTaskResponse, *task_stub.TaskServiceError) {
	// Authenticate user
	authenticatedUser, err := s.authentication.Authenticate(ctx, req.Auth.AccessToken)
	if err != nil {
		if errors.Is(err, auth.ErrParameterEmpty) || errors.Is(err, auth.ErrUserNotFound) {
			return &task_stub.StartTaskResponse{}, &task_stub.TaskServiceError{
				StatusCode: http.StatusUnauthorized,
				Error:      fmt.Errorf("unauthenticated: %w", err),
			}
		}

		return nil, &task_stub.TaskServiceError{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Errorf("authenticating user: %w", err),
		}
	}

	taskId, err := strconv.ParseInt(req.TaskId, 10, 64)
	if err != nil {
		return &task_stub.StartTaskResponse{}, &task_stub.TaskServiceError{
			StatusCode: http.StatusBadRequest,
			Error:      fmt.Errorf("invalid task id"),
		}
	}
	task, err := s.taskRepository.StartTask(ctx, authenticatedUser.ID, taskId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &task_stub.TaskServiceError{
				StatusCode: http.StatusBadRequest,
				Error:      fmt.Errorf("task not found"),
			}
		}

		return nil, &task_stub.TaskServiceError{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}

	responseData := task_stub.StartTaskResponse{
		Task: task_stub.Task{
			Id:                fmt.Sprintf("%d", task.Task.Id),
			Title:             task.Task.Title,
			Description:       task.Task.Description,
			Difficulty:        task.Task.Difficulty,
			Completed:         task.Completed,
			Content:           task.Task.Content,
			Author:            task.Task.Author,
			SatisfactionLevel: int32(task.SatisfactionLevel),
		},
	}
	if task.CompletedAt.Valid {
		responseData.Task.CompletedAt = task.CompletedAt.Time.Format(time.RFC3339)
	}

	return &responseData, nil
}
