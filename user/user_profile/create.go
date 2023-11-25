package user_profile

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	user_stub "kodiiing/user/stub"
	"strings"
)

func (u *Repository) Create(ctx context.Context, profile UserProfile) error {
	conn, err := u.db.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("acquiring connection from pool: %w", err)
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadWrite,
	})
	if err != nil {
		return fmt.Errorf("creating transaction: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO 
    	user_profiles 
    (
     user_id,
     join_reason,
     join_reason_other,
     coded_before,
     languages,
     target
    )
    VALUES 
    ($1, $2, $3, $4, $5, $6)`,
		profile.UserID,
		profile.JoinReason,
		sql.NullString{
			String: profile.JoinReasonOther,
			Valid:  profile.JoinReason != user_stub.JoinReasonOther,
		},
		profile.HasCodedBefore,
		sql.NullString{
			String: strings.Join(profile.Languages, ", "),
			Valid:  len(profile.Languages) > 0,
		},
		sql.NullString{
			String: profile.Target,
			Valid:  profile.Target != "",
		},
	)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return fmt.Errorf("rolling back transaction: %w (%s)", e, err.Error())
		}

		return fmt.Errorf("executing insert query: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}

	return nil
}
