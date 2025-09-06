package algorithms

import (
	"time"

	"github.com/n0l3r/limitron/store"
)

type TokenBucket struct {
	Rate     int64 // tokens per second
	Capacity int64 // max tokens in bucket
	Store    store.Store
}

func (tb *TokenBucket) Allow(key string) (bool, error) {
	now := time.Now().Unix()
	tokenKey := key + ":tb"
	lastRefillKey := key + ":tb_last"

	// Get last refill time and tokens
	tokens, err := tb.Store.Get(tokenKey)
	if err != nil {
		return false, err
	}
	lastRefill, err := tb.Store.Get(lastRefillKey)
	if err != nil {
		lastRefill = now
	}

	// Refill tokens
	elapsed := now - lastRefill
	newTokens := tokens + elapsed*tb.Rate
	if newTokens > tb.Capacity {
		newTokens = tb.Capacity
	}

	allowed := false
	if newTokens > 0 {
		allowed = true
		newTokens--
	}

	// Update store
	// TTL should be at least several seconds, or use tb.Capacity/tb.Rate
	tb.Store.Set(tokenKey, newTokens, 60)
	tb.Store.Set(lastRefillKey, now, 60)

	return allowed, nil
}
