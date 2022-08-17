package auth_service

import (
	"context"
	"database/sql"
	"fmt"
	"kodiiing/auth"
	"log"
	"time"
)

func (d *AuthService) CreateUserStatistics(ctx context.Context, userId int64, user auth.User) error {
	conn, err := d.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("failed to get connection: %w", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}()

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSnapshot})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO
			user_statistics
			(
				user_id,
				avatar_url,
				location,
				public_repositories,
				followers,
				following,
				created_at,
				updated_at,
				updated_by
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		userId,
		user.AvatarURL.String(),
		user.Location,
		user.PublicRepository,
		user.Followers,
		user.Following,
		user.CreatedAt,
		time.Now(),
		"system",
	)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return fmt.Errorf("failed to rollback transaction: %w", e)
		}

		return fmt.Errorf("failed to insert user statistics: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return fmt.Errorf("failed to rollback transaction: %w", e)
		}

		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
