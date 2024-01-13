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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func (s *TaskService) ListTasks(ctx context.Context, req *task_stub.ListTasksRequest) (*task_stub.ListTasksResponse, *task_stub.TaskServiceError) {
	ctx, span := tracer.Start(ctx, "TaskService.ListTasks")
	defer span.End()

	// Authenticate user
	span.AddEvent("authenticating user")
	authenticatedUser, err := s.authentication.Authenticate(ctx, req.Auth.AccessToken)
	if err != nil {
		span.SetStatus(codes.Error, "error when authenticating user")
		if errors.Is(err, auth.ErrParameterEmpty) || errors.Is(err, auth.ErrUserNotFound) {
			return &task_stub.ListTasksResponse{}, &task_stub.TaskServiceError{
				StatusCode: http.StatusUnauthorized,
				Error:      fmt.Errorf("unauthenticated: %w", err),
			}
		}

		span.RecordError(err, trace.WithAttributes(attribute.String("body", "tes")))
		return nil, &task_stub.TaskServiceError{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Errorf("authenticating user: %w", err),
		}
	}

	// TODO: filter task by track ID
	span.AddEvent("find task list")
	tasks, err := s.taskRepository.ListTask(ctx, authenticatedUser.ID, req.TrackId)
	if err != nil {
		// span.SetStatus(codes.Error, "error when getting task list") // ini keknya gaperlu record error, karena udah di level repo. nanti jadi dobel
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &task_stub.TaskServiceError{
				StatusCode: http.StatusBadRequest,
				Error:      fmt.Errorf("task not found"),
			}
		}

		// span.RecordError(err, trace.WithStackTrace(true)) // ini keknya gaperlu record error, karena udah di level repo. nanti jadi dobel
		return nil, &task_stub.TaskServiceError{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}

	var responseData task_stub.ListTasksResponse
	for _, task := range tasks {
		taskData := task_stub.Task{
			Id:          fmt.Sprintf("%d", task.Task.Id),
			Title:       task.Task.Title,
			Description: task.Task.Description,
			Difficulty:  task.Task.Difficulty,
			Completed:   task.Completed,
			Content:     task.Task.Content,
			Author:      task.Task.Author,
		}

		if task.SatisfactionLevel.Valid {
			taskData.SatisfactionLevel = int32(task.SatisfactionLevel.Int64)
		}
		if task.CompletedAt.Valid {
			taskData.CompletedAt = task.CompletedAt.Time.Format(time.RFC3339)
		}

		responseData.Tasks = append(responseData.Tasks, taskData)
	}

	span.SetStatus(codes.Ok, "success getting tasks")
	return &responseData, nil
}
