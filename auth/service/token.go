package auth_service

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func (d *AuthService) CreateUserAccessToken(ctx context.Context, userId int64, accessToken string, refreshToken string) error {
	// Encrypt both the access and refresh token
	encryptedAccessToken := d.aes.Encrypt(accessToken)
	encryptedRefreshToken := d.aes.Encrypt(refreshToken)

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

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSnapshot, ReadOnly: false})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	_, err = tx.ExecContext(
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
		if e := tx.Rollback(); e != nil {
			return fmt.Errorf("failed to rollback transaction: %w", e)
		}

		return fmt.Errorf("failed to insert user access token: %w", err)
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
