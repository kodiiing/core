package auth_service

import (
	"context"
	"errors"
	"kodiiing/auth"
	"kodiiing/auth/provider"
	auth_stub "kodiiing/auth/stub"
	"net/http"
)

func (d *AuthService) Login(ctx context.Context, req *auth_stub.LoginRequest) (*auth_stub.LoginResponse, *auth_stub.AuthenticationServiceError) {
	var currentAuthenticationProvider provider.Authentication

	switch req.Provider {
	case auth_stub.ProviderGITHUB:
		currentAuthenticationProvider = d.githubProvider
		break
	case auth_stub.ProviderGITLAB:
		currentAuthenticationProvider = d.gitlabProvider
		break
	default:
		return nil, &auth_stub.AuthenticationServiceError{
			StatusCode: http.StatusBadRequest,
			Error:      errors.New("invalid provider was specified"),
		}
	}

	// Exchange req.AccessCode to the GitHub or GitLab OAuth2 provider
	providerAccessToken, err := currentAuthenticationProvider.AcquireAccessToken(ctx, req.AccessCode)
	if err != nil {
		if errors.Is(err, provider.ErrCodeEmpty) {
			return nil, &auth_stub.AuthenticationServiceError{
				StatusCode: http.StatusBadRequest,
				Error:      err,
			}
		}

		return nil, &auth_stub.AuthenticationServiceError{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}

	// Call profile API
	providerUser, err := currentAuthenticationProvider.GetProfile(ctx, providerAccessToken)
	if err != nil {
		return nil, &auth_stub.AuthenticationServiceError{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}

	// Check if the providerUser is already registered on Kodiiing
	user, err := d.GetUserByProviderId(ctx, providerUser.ID)
	if err != nil && !errors.Is(err, auth.ErrUserNotFound) {
		return nil, &auth_stub.AuthenticationServiceError{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	} else {
		// User not found. We should register the user
		userId, err := d.CreateUser(ctx, &providerUser)
		if err != nil {
			// TODO
		}
		user.ID = userId

		repos, err := currentAuthenticationProvider.GetPublicRepositories(ctx, providerUser.Username)
		if err != nil {
			// TODO
		}

		err = d.CreateUserRepository(ctx, userId, repos)
		if err != nil {
			// TODO
		}

	}

	// Generate access token from ID
	// TODO: generate proper access token
	appAccessToken := ""
	err = d.sessionStore.Set(ctx, appAccessToken, user.ID)
	if err != nil {
		// TODO
	}
	return &auth_stub.LoginResponse{}, nil
}
