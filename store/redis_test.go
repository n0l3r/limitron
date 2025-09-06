package store

import (
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestRedisStore_SetGet(t *testing.T) {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	store := NewRedisStore(rdb)

	key := "test-redis-key"
	value := int64(42)

	err := store.Set(key, value, 2)
	if err != nil {
		t.Fatalf("failed to set redis key: %v", err)
	}
	got, err := store.Get(key)
	if err != nil {
		t.Fatalf("failed to get redis key: %v", err)
	}
	if got != value {
		t.Fatalf("expected %d, got %d", value, got)
	}
	time.Sleep(3 * time.Second)
	got, err = store.Get(key)
	if err == nil && got != 0 {
		t.Fatalf("expected key to expire, got %d", got)
	}
}
