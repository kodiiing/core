package user_service

import (
	"context"
	"database/sql"

	user_stub "kodiiing/user/stub"
)

type UserService struct {
	environment string
	db          *sql.DB
}

func NewUserService(env string, db *sql.DB) *UserService {
	return &UserService{environment: env, db: db}
}

func (d *UserService) Onboarding(ctx context.Context, req *user_stub.OnboardingRequest) (*user_stub.EmptyResponse, *user_stub.UserServiceError) {
	return &user_stub.EmptyResponse{}, nil
}
