package auth_service

import (
	"context"
	"fmt"
	"kodiiing/auth"
	"time"
)

func (d *AuthService) CreateUserStatistics(ctx context.Context, userId int64, user auth.User) error {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.Exec(
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
		if e := tx.Rollback(ctx); e != nil {
			return fmt.Errorf("failed to rollback transaction: %w", e)
		}

		return fmt.Errorf("failed to insert user statistics: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return fmt.Errorf("failed to rollback transaction: %w", e)
		}

		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
