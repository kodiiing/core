package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"kodiiing/auth"
	task_stub "kodiiing/task/stub"

	"github.com/jackc/pgx/v5"
)

func (s *TaskService) ListTasks(ctx context.Context, req *task_stub.ListTasksRequest) (*task_stub.ListTasksResponse, *task_stub.TaskServiceError) {
	// Authenticate user
	authenticatedUser, err := s.authentication.Authenticate(ctx, req.Auth.AccessToken)
	if err != nil {
		if errors.Is(err, auth.ErrParameterEmpty) || errors.Is(err, auth.ErrUserNotFound) {
			return &task_stub.ListTasksResponse{}, &task_stub.TaskServiceError{
				StatusCode: http.StatusUnauthorized,
				Error:      fmt.Errorf("unauthenticated: %w", err),
			}
		}

		return nil, &task_stub.TaskServiceError{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Errorf("authenticating user: %w", err),
		}
	}

	// TODO: filter task by track ID
	tasks, err := s.taskRepository.ListTask(ctx, authenticatedUser.ID, req.TrackId)
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

	var responseData task_stub.ListTasksResponse
	for _, task := range tasks {
		taskData := task_stub.Task{
			Id:                fmt.Sprintf("%d", task.Task.Id),
			Title:             task.Task.Title,
			Description:       task.Task.Description,
			Difficulty:        task.Task.Difficulty,
			Completed:         task.Completed,
			Content:           task.Task.Content,
			Author:            task.Task.Author,
			SatisfactionLevel: int32(task.SatisfactionLevel),
		}
		if task.CompletedAt.Valid {
			taskData.CompletedAt = task.CompletedAt.Time.Format(time.RFC3339)
		}

		responseData.Tasks = append(responseData.Tasks, taskData)
	}

	return &responseData, nil
}
