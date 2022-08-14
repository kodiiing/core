package provider

import (
	"context"
	"errors"
	"kodiiing/auth"
)

var ErrCodeEmpty = errors.New("code is empty")

type Authentication interface {
	AcquireAccessToken(ctx context.Context, code string) (string, error)
	GetProfile(ctx context.Context, accessToken string) (auth.User, error)
	GetPublicRepositories(ctx context.Context, username string) ([]auth.Repository, error)
}
