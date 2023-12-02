package user_profile

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

func (u *Repository) Exists(ctx context.Context, userId int64) (bool, error) {
	conn, err := u.db.Acquire(ctx)
	if err != nil {
		return false, fmt.Errorf("acquiring connection from pool: %w", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.ReadCommitted,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return false, fmt.Errorf("creating transaction: %w", err)
	}

	var exists bool
	err = tx.QueryRow(
		ctx,
		`SELECT EXISTS (SELECT * FROM user_profiles WHERE user_id = $1) AS exists`,
		userId,
	).Scan(&exists)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return false, fmt.Errorf("rolling back transaction: %w (%s)", e, err.Error())
		}

		return false, fmt.Errorf("executing insert query: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return false, fmt.Errorf("commiting transaction: %w", err)
	}

	return exists, nil
}
