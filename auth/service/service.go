package auth_service

import (
	auth_aes "kodiiing/auth/aes"
	auth_jwt "kodiiing/auth/jwt"
	auth_stub "kodiiing/auth/stub"

	"github.com/allegro/bigcache/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthService struct {
	pool        *pgxpool.Pool
	memory      *bigcache.BigCache
	aes         *auth_aes.Aes
	jwt         *auth_jwt.AuthJwt
	environment string
}

func NewAuthService(env string, pool *pgxpool.Pool, memory *bigcache.BigCache) auth_stub.AuthenticationServiceServer {
	return &AuthService{
		environment: env,
		pool:        pool,
		memory:      memory,
	}
}
