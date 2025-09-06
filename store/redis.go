package store

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisStore) Get(key string) (int64, error) {
	return r.client.Get(r.ctx, key).Int64()
}

func (r *RedisStore) Set(key string, value int64, expiry int64) error {
	return r.client.Set(r.ctx, key, value, time.Duration(expiry)*time.Second).Err()
}

func (r *RedisStore) Incr(key string, n int64, expiry int64) (int64, error) {
	res := r.client.IncrBy(r.ctx, key, n)
	r.client.Expire(r.ctx, key, time.Duration(expiry)*time.Second)
	return res.Val(), res.Err()
}
