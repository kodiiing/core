package auth_service

import (
	"context"
	auth_stub "kodiiing/auth/stub"
)

func (d *AuthService) Logout(ctx context.Context, req *auth_stub.LogoutRequest) (*auth_stub.EmptyResponse, *auth_stub.AuthenticationServiceError) {
	err := d.sessionStore.Revoke(ctx, req.AccessToken)
	if err != nil {
		// Must swallow error
		// TODO: Capture error
	}

	return &auth_stub.EmptyResponse{}, nil
}
