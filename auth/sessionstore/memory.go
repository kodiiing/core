package sessionstore

import (
	"context"
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"strconv"
	"time"
)

type MemoryStore struct {
	cache *bigcache.BigCache
}

func (m *MemoryStore) Revoke(ctx context.Context, accessToken string) error {
	return m.cache.Delete(accessToken)
}

func (m *MemoryStore) Get(ctx context.Context, accessToken string) (id int64, err error) {
	if accessToken == "" {
		return 0, fmt.Errorf("%w: access token", ErrEmptyValue)
	}

	value, err := m.cache.Get(accessToken)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return 0, ErrNotExists
		}

		return 0, fmt.Errorf("acquiring item from cache: %w", err)
	}

	out, err := strconv.ParseInt(string(value), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid value: %w", err)
	}

	return out, nil
}

func (m *MemoryStore) Set(ctx context.Context, accessToken string, id int64) error {
	if accessToken == "" {
		return fmt.Errorf("%w: access token", ErrEmptyValue)
	}

	return m.cache.Set(accessToken, []byte(strconv.FormatInt(id, 10)))
}

func NewMemory() (Store, error) {
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(time.Hour*24))
	if err != nil {
		return nil, fmt.Errorf("creating cache instance: %w", err)
	}

	return &MemoryStore{cache: cache}, nil
}
