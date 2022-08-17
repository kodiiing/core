package auth_service

import (
	"context"
	auth_stub "kodiiing/auth/stub"
)

func (d *AuthService) Login(ctx context.Context, req *auth_stub.LoginRequest) (*auth_stub.LoginResponse, *auth_stub.AuthenticationServiceError) {
	return &auth_stub.LoginResponse{}, nil
}
