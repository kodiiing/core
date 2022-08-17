package auth_middleware

import (
	"context"
	"kodiiing/auth"
)

type AuthMiddleware struct{}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (a *AuthMiddleware) Authenticate(ctx context.Context, accessToken string) (*auth.User, error) {
	// TODO
	return &auth.User{}, nil
}
