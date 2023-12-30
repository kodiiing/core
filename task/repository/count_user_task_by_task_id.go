package repository

import (
	"context"
	"fmt"
)

func (r *Repository) CountUserTaskByTaskId(ctx context.Context, userId, taskId int64) (total int64, err error) {
	var countUserTaskSql = `SELECT COUNT(id) AS total FROM user_tasks WHERE  user_id = $1 AND task_id = $2`

	err = r.db.QueryRow(ctx, countUserTaskSql, userId, taskId).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("error executing count query: %w", err)
	}

	return
}
