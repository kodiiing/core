package auth_middleware

import (
	"context"
	"errors"
	"fmt"
	"kodiiing/auth"
	auth_jwt "kodiiing/auth/jwt"
	auth_service "kodiiing/auth/service"
)

type AuthMiddleware struct {
	jwt     *auth_jwt.AuthJwt
	service *auth_service.AuthService
}

func NewAuthMiddleware(service *auth_service.AuthService, jwt *auth_jwt.AuthJwt) *AuthMiddleware {
	return &AuthMiddleware{
		jwt:     jwt,
		service: service,
	}
}

func (a *AuthMiddleware) Authenticate(ctx context.Context, accessToken string) (*auth.User, error) {
	// Make sure accessToken is not empty
	if accessToken == "" {
		return nil, auth.ErrParameterEmpty
	}

	// Parse accessToken as json web token
	userId, err := a.jwt.VerifyAccessToken(accessToken)
	if err != nil {
		return nil, err
	}

	// Get user from service
	user, err := a.service.GetUserById(ctx, userId)
	if err != nil {
		if errors.Is(err, auth.ErrUserNotFound) {
			return nil, auth.ErrUserNotFound
		}

		return nil, fmt.Errorf("getting user: %w", err)
	}

	return &user, nil
}
