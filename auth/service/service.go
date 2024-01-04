package auth_service

import (
	auth_aes "kodiiing/auth/aes"
	auth_jwt "kodiiing/auth/jwt"
	"kodiiing/auth/provider"
	"kodiiing/auth/sessionstore"

	"github.com/allegro/bigcache/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthService struct {
	pool           *pgxpool.Pool
	memory         *bigcache.BigCache
	aes            *auth_aes.Aes
	jwt            *auth_jwt.AuthJwt
	sessionStore   sessionstore.Store
	githubProvider provider.Authentication
	gitlabProvider provider.Authentication
	environment    string
}

func NewAuthService(env string, pool *pgxpool.Pool, memory *bigcache.BigCache, store sessionstore.Store) *AuthService {
	return &AuthService{
		environment:  env,
		pool:         pool,
		memory:       memory,
		sessionStore: store,
	}
}
