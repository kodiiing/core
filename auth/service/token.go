package auth_service

import (
	"context"
	"fmt"
	"time"
)

func (d *AuthService) CreateUserAccessToken(ctx context.Context, userId int64, accessToken string, refreshToken string) error {
	// Encrypt both the access and refresh token
	encryptedAccessToken := d.aes.Encrypt(accessToken)
	encryptedRefreshToken := d.aes.Encrypt(refreshToken)

	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO
			user_accesstoken
			(
				user_id,
				access_token,
				refresh_token,
				created_at,
				updated_at,
				updated_by
			)
		VALUES
			($1, $2, $3, $4, $4, $5)
		ON CONFLICT (user_id) DO UPDATE
			SET
				access_token = $2,
				refresh_token = $3,
				updated_at = $4,
				updated_by = $5`,
		userId,
		encryptedAccessToken,
		encryptedRefreshToken,
		time.Now(),
		"system",
	)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return fmt.Errorf("failed to rollback transaction: %w", e)
		}

		return fmt.Errorf("failed to insert user access token: %w", err)
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
