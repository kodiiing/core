package sessionstore

import "context"

type RedisStore struct{}

func (r *RedisStore) Revoke(ctx context.Context, accessToken string) error {
	//TODO implement me
	panic("implement me")
}

func (r *RedisStore) Get(ctx context.Context, accessToken string) (id int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisStore) Set(ctx context.Context, accessToken string, id int64) error {
	//TODO implement me
	panic("implement me")
}

func NewRedisStore() (Store, error) {
	return &RedisStore{}, nil
}
