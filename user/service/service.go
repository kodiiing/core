package user_service

import (
	"context"
	"errors"
	"fmt"
	"kodiiing/auth"
	"kodiiing/user/user_profile"
	"net/http"
	"time"

	user_stub "kodiiing/user/stub"
)

type UserService struct {
	environment           string
	userProfileRepository *user_profile.Repository
	authentication        auth.Authenticate
}

func NewUserService(env string, userProfileRepository *user_profile.Repository) user_stub.UserServiceServer {
	return &UserService{environment: env, userProfileRepository: userProfileRepository}
}

func (d *UserService) Onboarding(ctx context.Context, req *user_stub.OnboardingRequest) (*user_stub.EmptyResponse, *user_stub.UserServiceError) {
	// Authenticate user
	authenticatedUser, err := d.authentication.Authenticate(ctx, req.Auth.AccessToken)
	if err != nil {
		if errors.Is(err, auth.ErrParameterEmpty) || errors.Is(err, auth.ErrUserNotFound) {
			return &user_stub.EmptyResponse{}, &user_stub.UserServiceError{
				StatusCode: http.StatusUnauthorized,
				Error:      fmt.Errorf("unauthenticated: %w", err),
			}
		}

		return nil, &user_stub.UserServiceError{
			StatusCode: http.StatusInternalServerError,
			Error:      fmt.Errorf("authenticating user: %w", err),
		}
	}

	// If user has onboarded before, return an empty response
	hasOnboarded, err := d.userProfileRepository.Exists(ctx, authenticatedUser.ID)
	if err != nil {
		return nil, &user_stub.UserServiceError{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}

	if hasOnboarded {
		return &user_stub.EmptyResponse{}, nil
	}

	// User have not onboard yet.`
	err = d.userProfileRepository.Create(ctx, user_profile.UserProfile{
		UserID:          authenticatedUser.ID,
		JoinReason:      req.Reason,
		JoinReasonOther: req.ReasonOther,
		HasCodedBefore:  req.CodedBefore,
		Languages:       req.Languages,
		Target:          req.Target,
		UpdatedAt:       time.Now(),
		UpdatedBy:       "system",
	})
	if err != nil {
		return nil, &user_stub.UserServiceError{
			StatusCode: http.StatusInternalServerError,
			Error:      err,
		}
	}

	return &user_stub.EmptyResponse{}, nil
}
