package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type InsertTaskAssessmentIn struct {
	UserId            int64
	TaskId            int64
	SatisfactionLevel int64
	Comments          string
}

func (r *Repository) InsertTaskAssessment(ctx context.Context, data InsertTaskAssessmentIn) (affected int64, err error) {
	if data.UserId == 0 || data.TaskId == 0 {
		return 0, ErrNoRows
	}

	var insertAssessmentSql = `UPDATE user_tasks SET 
		satisfaction_level = $1, 
		comments = $2
	WHERE
		user_id = $3
		AND
		task_id = $4`

	commandTag, err := r.db.Exec(ctx, insertAssessmentSql,
		data.SatisfactionLevel, data.Comments,
		data.UserId, data.TaskId,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNoRows
		}

		return
	}

	return commandTag.RowsAffected(), nil
}
