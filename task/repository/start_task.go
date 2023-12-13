package repository

import (
	"context"
	"database/sql"
	"kodiiing/task"
	"time"

	"github.com/jackc/pgx/v5"
)

type StartTaskOut struct {
	Task

	Completed         bool
	CompletedAt       sql.NullTime
	SatisfactionLevel int64
}

func (r *Repository) StartTask(ctx context.Context, userId, taskId int64) (out StartTaskOut, err error) {
	if taskId == 0 || userId == 0 {
		return StartTaskOut{}, pgx.ErrNoRows
	}

	var selectTaskSql = `	
	SELECT
		t.id AS task_id, t.title, t.description, t.difficulty, t.content, u.name AS author,
		t.created_at, t.created_by, t.updated_at, t.updated_by
	FROM
		tasks AS t
		INNER JOIN users AS u ON u.id = t.author`
	err = r.db.QueryRow(ctx, selectTaskSql).Scan(
		&out.Task.Id, &out.Task.Title, &out.Task.Description, &out.Task.Difficulty, &out.Task.Content,
		&out.Task.Author, &out.Task.CreatedAt, &out.Task.CreatedBy, &out.Task.UpdatedAt,
		&out.Task.UpdatedBy,
	)
	if err != nil {
		return
	}

	var insertUserTaskSql = "INSERT INTO user_tasks (task_id, user_id, status, started_at) VALUES ($1, $2, $3, $d) RETURNING finished_at, satisfaction_level"

	err = r.db.QueryRow(ctx, insertUserTaskSql,
		taskId, userId, task.USER_TASK_STATUS_IN_PROGRESS, time.Now().UTC().Format(time.RFC3339),
	).Scan(
		&out.CompletedAt, &out.SatisfactionLevel,
	)
	if err != nil {
		return
	}

	if out.CompletedAt.Valid {
		out.Completed = true
	}

	return
}
