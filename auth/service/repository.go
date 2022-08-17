package auth_service

import (
	"context"
	"database/sql"
	"fmt"
	"kodiiing/auth"
	"log"
	"time"
)

func (d *AuthService) CreateUserRepository(ctx context.Context, userId int64, repositories []auth.Repository) error {
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

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	query, err := tx.PrepareContext(
		ctx,
		`INSERT INTO
			user_repositories
			(
				user_id,
				repository_id,
				provider,
				name,
				url,
				description,
				fork,
				fork_count,
				star_count,
				owner_username,
				created_at,
				last_activity_at,
				updated_at,
				updated_by
			)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
	)
	if err != nil {
		if e := tx.Rollback(); e != nil {
			return fmt.Errorf("failed to rollback transaction: %w", e)
		}

		return fmt.Errorf("failed to prepare query: %w", err)
	}

	for _, repository := range repositories {
		_, err := tx.StmtContext(ctx, query).ExecContext(
			ctx,
			userId,
			repository.ID,
			repository.Provider.ToUint8(),
			repository.Name,
			repository.URL.String(),
			repository.Description,
			repository.Fork,
			repository.ForksCount,
			repository.StarsCount,
			repository.OwnerUsername,
			repository.CreatedAt,
			repository.LastActivityAt,
			time.Now(),
			"system",
		)
		if err != nil {
			if e := tx.Rollback(); e != nil {
				return fmt.Errorf("failed to rollback transaction: %w", e)
			}

			return fmt.Errorf("failed to insert user repository: %w", err)
		}
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
