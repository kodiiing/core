package provider

import (
	"context"
	"errors"
)

var ErrCodeEmpty = errors.New("code is empty")

type Authentication interface {
	AcquireAccessToken(ctx context.Context, code string) (string, error)
}
