package auth_service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"kodiiing/auth"
	"kodiiing/fgob"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/jackc/pgx/v5"
)

func (d *AuthService) CreateUser(ctx context.Context, user *auth.User) (id int64, err error) {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("starting transaction: %w", err)
	}

	var insertId int64
	err = tx.QueryRow(
		ctx,
		`INSERT INTO users
			(
				provider,
				provider_id,
				name,
				username,
				email,
				profile_url,
				created_at,
				registered_at,
				updated_at,
				updated_by
			)
			VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING
			`,
		user.Provider.ToUint8(),
		user.ID,
		user.Name,
		user.Username,
		user.Email,
		user.ProfileURL,
		user.CreatedAt,
		time.Now(),
		time.Now(),
		"system",
	).Scan(&insertId)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return 0, fmt.Errorf("error rolling back transaction: %w", e)
		}

		return 0, fmt.Errorf("error inserting user: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return 0, fmt.Errorf("error rolling back transaction: %w", e)
		}

		return 0, fmt.Errorf("error committing transaction: %w", err)
	}

	return insertId, nil
}

func (d *AuthService) GetUserById(ctx context.Context, id int64) (auth.User, error) {
	// Is it cached?
	cachedUser, err := d.memory.Get("user:id:" + strconv.FormatInt(id, 10))
	if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
		return auth.User{}, fmt.Errorf("error getting user from cache: %w", err)
	}

	if cachedUser != nil {
		var user auth.User
		err := fgob.Unmarshal(cachedUser, &user)
		if err != nil {
			return auth.User{}, fmt.Errorf("error unmarshalling user from cache: %w", err)
		}

		return user, nil
	}

	tx, err := d.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadOnly})
	if err != nil {
		return auth.User{}, fmt.Errorf("starting transaction: %w", err)
	}

	var user auth.User
	var nullAvatarUrl sql.NullString
	var nullLocation sql.NullString
	err = tx.QueryRow(
		ctx,
		`SELECT
			users.provider,
			users.provider_id,
			users.name,
			users.username,
			users.email,
			users.profile_url,
			users.created_at,
			users.registered_at,
			user_statistics.avatar_url,
			user_statistics.location,
			user_statistics.public_repositories,
			user_statistics.followers,
			user_statistics.following
		FROM
			users
		INNER JOIN
			user_statistics
		ON
			users.id = user_statistics.user_id
		WHERE
			users.id = $1`,
		id,
	).Scan(
		&user.Provider,
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.ProfileURL,
		&user.CreatedAt,
		&user.RegisteredAt,
		&nullAvatarUrl,
		&nullLocation,
		&user.PublicRepository,
		&user.Followers,
		&user.Following,
	)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return auth.User{}, fmt.Errorf("error rolling back transaction: %w", e)
		}

		if errors.Is(err, sql.ErrNoRows) {
			return auth.User{}, auth.ErrUserNotFound
		}

		return auth.User{}, fmt.Errorf("error getting user: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return auth.User{}, fmt.Errorf("error rolling back transaction: %w", e)
		}

		return auth.User{}, fmt.Errorf("error committing transaction: %w", err)
	}

	if nullAvatarUrl.Valid {
		avatarUrl, err := url.Parse(nullAvatarUrl.String)
		if err != nil {
			return auth.User{}, fmt.Errorf("error parsing avatar url: %w", err)
		}

		user.AvatarURL = avatarUrl
	}

	if nullLocation.Valid {
		user.Location = nullLocation.String
	}

	defer func(user auth.User) {
		marshaledUser, err := fgob.Marshal(user)
		if err != nil {
			log.Printf("error marshalling user: %v", err)
			return
		}

		err = d.memory.Set("user:id:"+strconv.FormatInt(id, 10), marshaledUser)
		if err != nil {
			log.Printf("error setting user in cache: %v", err)
		}
	}(user)

	return user, nil
}

func (d *AuthService) GetUserByUsername(ctx context.Context, username string) (auth.User, error) {
	if username == "" {
		return auth.User{}, auth.ErrParameterEmpty
	}

	// Is it cached?
	cachedUser, err := d.memory.Get("user:username:" + username)
	if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
		return auth.User{}, fmt.Errorf("error getting user from cache: %w", err)
	}

	if cachedUser != nil {
		var user auth.User
		err := fgob.Unmarshal(cachedUser, &user)
		if err != nil {
			return auth.User{}, fmt.Errorf("error unmarshalling user from cache: %w", err)
		}

		return user, nil
	}

	tx, err := d.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadOnly})
	if err != nil {
		return auth.User{}, fmt.Errorf("starting transaction: %w", err)
	}

	var user auth.User
	var nullAvatarUrl sql.NullString
	var nullLocation sql.NullString
	err = tx.QueryRow(
		ctx,
		`SELECT
			users.provider,
			users.provider_id,
			users.name,
			users.username,
			users.email,
			users.profile_url,
			users.created_at,
			users.registered_at,
			user_statistics.avatar_url,
			user_statistics.location,
			user_statistics.public_repositories,
			user_statistics.followers,
			user_statistics.following
		FROM
			users
		INNER JOIN
			user_statistics
		ON
			users.id = user_statistics.user_id
		WHERE
			users.username = $1`,
		username,
	).Scan(
		&user.Provider,
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.ProfileURL,
		&user.CreatedAt,
		&user.RegisteredAt,
		&nullAvatarUrl,
		&nullLocation,
		&user.PublicRepository,
		&user.Followers,
		&user.Following,
	)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return auth.User{}, fmt.Errorf("error rolling back transaction: %w", e)
		}

		if errors.Is(err, sql.ErrNoRows) {
			return auth.User{}, auth.ErrUserNotFound
		}

		return auth.User{}, fmt.Errorf("error getting user: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return auth.User{}, fmt.Errorf("error rolling back transaction: %w", e)
		}

		return auth.User{}, fmt.Errorf("error committing transaction: %w", err)
	}

	if nullAvatarUrl.Valid {
		avatarUrl, err := url.Parse(nullAvatarUrl.String)
		if err != nil {
			return auth.User{}, fmt.Errorf("error parsing avatar url: %w", err)
		}

		user.AvatarURL = avatarUrl
	}

	if nullLocation.Valid {
		user.Location = nullLocation.String
	}

	defer func(user auth.User) {
		marshaledUser, err := fgob.Marshal(user)
		if err != nil {
			log.Printf("error marshalling user: %v", err)
			return
		}

		err = d.memory.Set("user:username:"+username, marshaledUser)
		if err != nil {
			log.Printf("error setting user in cache: %v", err)
		}
	}(user)

	return user, nil
}

func (d *AuthService) GetUserByEmail(ctx context.Context, email string) (auth.User, error) {
	if email == "" {
		return auth.User{}, auth.ErrParameterEmpty
	}

	// Is it cached?
	cachedUser, err := d.memory.Get("user:email:" + email)
	if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
		return auth.User{}, fmt.Errorf("error getting user from cache: %w", err)
	}

	if cachedUser != nil {
		var user auth.User
		err := fgob.Unmarshal(cachedUser, &user)
		if err != nil {
			return auth.User{}, fmt.Errorf("error unmarshalling user from cache: %w", err)
		}

		return user, nil
	}

	tx, err := d.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted, AccessMode: pgx.ReadOnly})
	if err != nil {
		return auth.User{}, fmt.Errorf("starting transaction: %w", err)
	}

	var user auth.User
	var nullAvatarUrl sql.NullString
	var nullLocation sql.NullString
	err = tx.QueryRow(
		ctx,
		`SELECT
			users.provider,
			users.provider_id,
			users.name,
			users.username,
			users.email,
			users.profile_url,
			users.created_at,
			users.registered_at,
			user_statistics.avatar_url,
			user_statistics.location,
			user_statistics.public_repositories,
			user_statistics.followers,
			user_statistics.following
		FROM
			users
		INNER JOIN
			user_statistics
		ON
			users.id = user_statistics.user_id
		WHERE
			users.email = $1`,
		email,
	).Scan(
		&user.Provider,
		&user.ID,
		&user.Name,
		&user.Username,
		&user.Email,
		&user.ProfileURL,
		&user.CreatedAt,
		&user.RegisteredAt,
		&nullAvatarUrl,
		&nullLocation,
		&user.PublicRepository,
		&user.Followers,
		&user.Following,
	)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return auth.User{}, fmt.Errorf("error rolling back transaction: %w", e)
		}

		if errors.Is(err, sql.ErrNoRows) {
			return auth.User{}, auth.ErrUserNotFound
		}

		return auth.User{}, fmt.Errorf("error getting user: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		if e := tx.Rollback(ctx); e != nil {
			return auth.User{}, fmt.Errorf("error rolling back transaction: %w", e)
		}

		return auth.User{}, fmt.Errorf("error committing transaction: %w", err)
	}

	if nullAvatarUrl.Valid {
		avatarUrl, err := url.Parse(nullAvatarUrl.String)
		if err != nil {
			return auth.User{}, fmt.Errorf("error parsing avatar url: %w", err)
		}

		user.AvatarURL = avatarUrl
	}

	if nullLocation.Valid {
		user.Location = nullLocation.String
	}

	defer func(user auth.User) {
		marshaledUser, err := fgob.Marshal(user)
		if err != nil {
			log.Printf("error marshalling user: %v", err)
			return
		}

		err = d.memory.Set("user:email:"+email, marshaledUser)
		if err != nil {
			log.Printf("error setting user in cache: %v", err)
		}
	}(user)

	return user, nil
}
