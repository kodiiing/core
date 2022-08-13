package auth_service

import (
	"context"
	"database/sql"

	auth_stub "kodiiing/auth/stub"

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

func (d *AuthService) Login(ctx context.Context, req *auth_stub.LoginRequest) (*auth_stub.LoginResponse, *auth_stub.AuthenticationServiceError) {
	return &auth_stub.LoginResponse{}, nil
}

func (d *AuthService) Logout(ctx context.Context, req *auth_stub.LogoutRequest) (*auth_stub.EmptyResponse, *auth_stub.AuthenticationServiceError) {
	return &auth_stub.EmptyResponse{}, nil
}
