package repository

import (
	"context"
	"errors"
	"fmt"
	task_stub "kodiiing/task/stub"

	"github.com/jackc/pgx/v5"
)

func (r *Repository) CountUserTaskByTaskId(ctx context.Context, userId, taskId int64) (total int64, err error) {
	var countUserTaskSql = `SELECT COUNT(id) AS total FROM user_tasks WHERE  user_id = $1 AND task_id = $2`

	err = r.db.QueryRow(ctx, countUserTaskSql, userId, taskId).Scan(&total)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, task_stub.TaskRepositoryErrNoRows
		}

		return 0, fmt.Errorf("error executing count query: %w", err)
	}

	return
}
