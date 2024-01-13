package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type ListTaskOut struct {
	Task

	Completed         bool
	CompletedAt       sql.NullTime
	SatisfactionLevel sql.NullInt64
}

func (r *Repository) ListTask(ctx context.Context, userId int64, trackId string) (out []ListTaskOut, err error) {
	// TODO: filter by track ID
	if userId == 0 {
		return []ListTaskOut{}, pgx.ErrNoRows
	}

	ctx, span := tracer.Start(ctx, "Repository.ListTask")
	defer span.End()

	span.AddEvent("creating database transaction")
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return
	}

	var findTaskSql = `
	SELECT
		t.id AS task_id, t.title, t.description, t.difficulty, t.content, t.author AS author,
		t.created_at, t.created_by, t.updated_at, t.updated_by,
		ut.finished_at, ut.satisfaction_level,
		CASE
			WHEN
				ut.finished_at IS NULL
			THEN
				false
			ELSE
				true
		END AS completed
	FROM
		tasks AS t
		LEFT JOIN user_tasks AS ut ON ut.task_id = t.id AND ut.user_id = $1`

	span.AddEvent("finding task lists")
	rows, err := tx.Query(ctx, findTaskSql, userId)
	if err != nil {
		span.SetStatus(codes.Error, "error when finding task lists")
		span.RecordError(err, trace.WithStackTrace(true))
		if e := tx.Rollback(ctx); e != nil {
			return out, fmt.Errorf("rolling back transaction: %w (%s)", e, err.Error())
		}

		return out, fmt.Errorf("executing insert query: %w", err)
	}

	for rows.Next() {
		var row ListTaskOut
		err = rows.Scan(
			&row.Task.Id, &row.Task.Title, &row.Task.Description, &row.Task.Difficulty, &row.Task.Content,
			&row.Task.Author, &row.Task.CreatedAt, &row.Task.CreatedBy, &row.Task.UpdatedAt, &row.Task.UpdatedBy,
			&row.CompletedAt, &row.SatisfactionLevel, &row.Completed,
		)
		if err != nil {
			if e := tx.Rollback(ctx); e != nil {
				return out, fmt.Errorf("rolling back transaction: %w (%s)", e, err.Error())
			}

			return []ListTaskOut{}, fmt.Errorf("executing select query: %w", err)
		}

		out = append(out, row)
	}

	span.AddEvent("committing transaction")
	err = tx.Commit(ctx)
	if err != nil {
		span.SetStatus(codes.Error, "error during committing transaction")
		span.RecordError(err, trace.WithStackTrace(true))
		if e := tx.Rollback(ctx); e != nil {
			return out, fmt.Errorf("rolling back transaction: %w (%s)", e, err.Error())
		}

		return out, fmt.Errorf("executing select query: %w", err)
	}

	span.SetStatus(codes.Ok, "success finding task lists")
	return
}
