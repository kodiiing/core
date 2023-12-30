package sessionstore

import (
	"context"
	"errors"
)

var ErrEmptyValue = errors.New("empty value")
var ErrNotExists = errors.New("not exist")

type Store interface {
	Get(ctx context.Context, accessToken string) (id int64, err error)
	Set(ctx context.Context, accessToken string, id int64) error
}
