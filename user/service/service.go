package user_service

import (
	"context"

	user_stub "kodiiing/user/stub"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	environment string
	pool        *pgxpool.Pool
}

func NewUserService(env string, pool *pgxpool.Pool) *UserService {
	return &UserService{environment: env, pool: pool}
}

func (d *UserService) Onboarding(ctx context.Context, req *user_stub.OnboardingRequest) (*user_stub.EmptyResponse, *user_stub.UserServiceError) {
	return &user_stub.EmptyResponse{}, nil
}
