package auth_service

import (
	"database/sql"

	auth_aes "kodiiing/auth/aes"
	auth_jwt "kodiiing/auth/jwt"

	"github.com/allegro/bigcache/v3"
)

type AuthService struct {
	db          *sql.DB
	memory      *bigcache.BigCache
	aes         *auth_aes.Aes
	jwt         *auth_jwt.AuthJwt
	environment string
}

func NewAuthService(env string, db *sql.DB, memory *bigcache.BigCache) *AuthService {
	return &AuthService{
		environment: env,
		db:          db,
		memory:      memory,
	}
}
