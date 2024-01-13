package service

import (
	"fmt"
	"kodiiing/auth"
	taskRepository "kodiiing/task/repository"
	task_stub "kodiiing/task/stub"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
)

type TaskService struct {
	pool           *pgxpool.Pool
	authentication auth.Authenticate

	taskRepository *taskRepository.Repository
}

type Config struct {
	Pool           *pgxpool.Pool
	Authentication auth.Authenticate
	TaskRepository *taskRepository.Repository
}

var tracer = otel.Tracer("kodiiing/task/service")

func NewTaskService(config *Config) (task_stub.TaskServiceServer, error) {
	if config.Pool == nil {
		return nil, fmt.Errorf("database connection required on task/service module")
	}
	if config.Authentication == nil {
		return nil, fmt.Errorf("authentication service required on task/service module")
	}
	if config.TaskRepository == nil {
		return nil, fmt.Errorf("taskRepository required on task/service module")
	}

	return &TaskService{
		pool:           config.Pool,
		authentication: config.Authentication,
		taskRepository: config.TaskRepository,
	}, nil
}
