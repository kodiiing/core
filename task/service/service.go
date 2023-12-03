package service

import (
	"kodiiing/auth"
	taskRepository "kodiiing/task/repository"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TaskService struct {
	pool           *pgxpool.Pool
	authentication auth.Authenticate

	taskRepository *taskRepository.Repository
}

type Dependency struct {
	Pool           *pgxpool.Pool
	Authentication auth.Authenticate
	taskRepository *taskRepository.Repository
}

func NewTaskService(d *Dependency) *TaskService {
	if d.Pool == nil {
		log.Fatal("[x] database connection required on task/service module")
	}
	if d.Authentication == nil {
		log.Fatal("[x] authentication service required on task/service module")
	}
	if d.taskRepository == nil {
		log.Fatal("[x] taskRepository required on task/service module")
	}

	return &TaskService{
		pool:           d.Pool,
		authentication: d.Authentication,
		taskRepository: d.taskRepository,
	}
}
