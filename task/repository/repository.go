package repository

import (
	"log"
	"time"

	task_stub "kodiiing/task/stub"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	Id          int64
	Title       string
	Description string
	Difficulty  task_stub.TaskDifficulty
	Content     string
	Author      string
	CreatedAt   time.Time
	CreatedBy   string
	UpdatedAt   time.Time
	UpdatedBy   string
}

type Repository struct {
	db *pgxpool.Pool
}

type Dependency struct {
	DB *pgxpool.Pool
}

func NewTaskRepository(d *Dependency) *Repository {
	if d.DB == nil {
		log.Fatal("[x] database connection required on task/repository module")
	}

	return &Repository{
		db: d.DB,
	}
}
