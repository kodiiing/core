package sessionstore

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	db *pgxpool.Pool
}

func (p *PostgresStore) Get(ctx context.Context, accessToken string) (id int64, err error) {
	//TODO implement me
	// Read to user_accesstoken table
	panic("implement me")
}

func (p *PostgresStore) Set(ctx context.Context, accessToken string, id int64) error {
	//TODO implement me
	// Read to user_accesstoken table
	panic("implement me")
}

func NewPostgresStore(db *pgxpool.Pool) (Store, error) {
	if db == nil {
		return nil, fmt.Errorf("db is nil")
	}

	return &PostgresStore{db: db}, nil
}
