package limitron

import (
	"errors"

	"github.com/go-redis/redis/v8"
	"github.com/n0l3r/limitron/algorithms"
	"github.com/n0l3r/limitron/store"
)

type Limiter interface {
	Allow(key string) (bool, error)
}

// NewLimiter creates a new rate limiter based on the provided configuration
func NewLimiter(cfg Config) (Limiter, error) {
	switch cfg.Algorithm {
	case AlgorithmLeakyBucket:
		return NewFixedWindowLimiter(cfg)
	case AlgorithmFixedWindow:
		return NewFixedWindowLimiter(cfg)
	case AlgorithmSlidingWindow:
		return NewSlidingWindowLimiter(cfg)
	case AlgorithmTokenBucket:
		return NewTokenBucketLimiter(cfg)
	default:
		return nil, errors.New("invalid algorithm")
	}
}

// NewTokenBucketLimiter creates a Token Bucket rate limiter
func NewTokenBucketLimiter(cfg Config) (Limiter, error) {
	s, err := buildStore(cfg)
	if err != nil {
		return nil, err
	}
	return &algorithms.TokenBucket{
		Rate:     cfg.Rate,
		Capacity: cfg.Capacity,
		Store:    s,
	}, nil
}

// NewFixedWindowLimiter creates a Fixed Window rate limiter
func NewFixedWindowLimiter(cfg Config) (Limiter, error) {
	s, err := buildStore(cfg)
	if err != nil {
		return nil, err
	}
	return &algorithms.FixedWindow{
		WindowSize: cfg.Rate,
		Limit:      cfg.Capacity,
		Store:      s,
	}, nil
}

// NewSlidingWindowLimiter creates a Sliding Window rate limiter
func NewSlidingWindowLimiter(cfg Config) (Limiter, error) {
	s, err := buildStore(cfg)
	if err != nil {
		return nil, err
	}
	return &algorithms.SlidingWindow{
		WindowSize: cfg.Rate,
		Limit:      cfg.Capacity,
		Store:      s,
	}, nil
}

// NewLeakyBucketLimiter creates a Leaky Bucket rate limiter
func NewLeakyBucketLimiter(cfg Config) (Limiter, error) {
	s, err := buildStore(cfg)
	if err != nil {
		return nil, err
	}
	return &algorithms.LeakyBucket{
		LeakRate: cfg.Rate,
		Capacity: cfg.Capacity,
		Store:    s,
	}, nil
}

// buildStore helper
func buildStore(cfg Config) (store.Store, error) {
	switch cfg.StoreType {
	case StoreTypeMemory:
		return store.NewMemoryStore(), nil
	case StoreTypeRedis:
		var rdb *redis.Client
		if cfg.RedisClient == nil {
			return nil, errors.New("redis client is required")
		}
		return store.NewRedisStore(rdb), nil
	default:
		return nil, errors.New("invalid store type")
	}
}
