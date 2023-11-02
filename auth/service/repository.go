package auth_service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kodiiing/auth"
	"kodiiing/fgob"
	"log"
	"strconv"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/jackc/pgx/v5"
)

func (d *AuthService) CreateUserRepository(ctx context.Context, userId int64, repositories []auth.Repository) error {
	db, err := d.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection from pool: %w", err)
	}
	defer db.Release()

	tx, err := db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	var sql = `
	INSERT INTO
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
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	for _, repository := range repositories {
		_, err := tx.Exec(
			ctx,
			sql,
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
			if e := tx.Rollback(ctx); e != nil {
				return fmt.Errorf("failed to rollback transaction: %w", e)
			}

			return fmt.Errorf("failed to insert user repository: %w", err)
		}
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

func (d *AuthService) GetUserRepositoryByUserId(ctx context.Context, userId int64) ([]auth.Repository, error) {
	cachedRepository, err := d.memory.Get("user:repository:id:" + strconv.FormatInt(userId, 10))
	if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
		return []auth.Repository{}, fmt.Errorf("failed to get user repository from cache: %w", err)
	}

	if cachedRepository != nil {
		var repositories []auth.Repository
		err := fgob.Unmarshal(cachedRepository, &repositories)
		if err != nil {
			return []auth.Repository{}, fmt.Errorf("failed to unmarshal user repository from cache: %w", err)
		}

		return repositories, nil
	}

	conn, err := d.pool.Acquire(ctx)
	if err != nil {
		return []auth.Repository{}, fmt.Errorf("failed to acquire connection from pool: %w", err)
	}
	defer conn.Release()

	rows, err := conn.Query(
		ctx,
		`SELECT
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
			last_activity_at
		FROM
			user_repositories
		WHERE
			user_id = $1`,
		userId,
	)
	if err != nil {
		return []auth.Repository{}, fmt.Errorf("failed to query user repositories: %w", err)
	}
	defer rows.Close()

	var repositories []auth.Repository
	for rows.Next() {
		var repository auth.Repository
		var nullDescription sql.NullString
		err := rows.Scan(
			&repository.ID,
			&repository.Provider,
			&repository.Name,
			&repository.URL,
			&nullDescription,
			&repository.Fork,
			&repository.ForksCount,
			&repository.StarsCount,
			&repository.OwnerUsername,
			&repository.CreatedAt,
			&repository.LastActivityAt,
		)
		if err != nil {
			return []auth.Repository{}, fmt.Errorf("failed to scan user repository: %w", err)
		}

		if nullDescription.Valid {
			repository.Description = nullDescription.String
		}

		repositories = append(repositories, repository)
	}

	go func(repositories []auth.Repository) {
		marshaledRepositories, err := fgob.Marshal(repositories)
		if err != nil {
			log.Printf("failed to marshal user repositories: %v", err)
			return
		}

		err = d.memory.Set("user:repository:id:"+strconv.FormatInt(userId, 10), marshaledRepositories)
		if err != nil {
			log.Printf("failed to set user repository in cache: %v", err)
		}
	}(repositories)

	return repositories, nil
}
