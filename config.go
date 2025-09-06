package limitron

import "github.com/go-redis/redis/v8"

type Config struct {
	Algorithm   string
	Rate        int64
	Capacity    int64
	StoreType   string
	RedisClient *redis.Client
}
