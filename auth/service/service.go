package auth_service

import (
	"database/sql"

	"github.com/allegro/bigcache/v3"
)

type AuthService struct {
	db          *sql.DB
	memory      *bigcache.BigCache
	environment string
}

func NewAuthService(env string, db *sql.DB, memory *bigcache.BigCache) *AuthService {
	return &AuthService{
		environment: env,
		db:          db,
		memory:      memory,
	}
}
