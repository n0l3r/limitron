package algorithms

import (
	"time"

	"github.com/n0l3r/limitron/store"
)

// LeakyBucket implements the leaky bucket rate limiting algorithm.
type LeakyBucket struct {
	LeakRate int64 // tokens leaked per second
	Capacity int64 // max tokens in bucket
	Store    store.Store
}

// Allow determines if a request is allowed under the leaky bucket algorithm.
// Returns true if allowed, false otherwise.
func (lb *LeakyBucket) Allow(key string) (bool, error) {
	now := time.Now().Unix()
	bucketKey := key + ":lb"
	lastLeakKey := key + ":lb_last"

	// Get current tokens in bucket
	tokens, err := lb.Store.Get(bucketKey)
	if err != nil {
		return false, err
	}

	// Get last leak timestamp
	lastLeak, err := lb.Store.Get(lastLeakKey)
	if err != nil {
		lastLeak = now
	}

	// Calculate leaked tokens
	elapsed := now - lastLeak
	leaked := elapsed * lb.LeakRate
	if leaked > 0 {
		if tokens < leaked {
			tokens = 0
		} else {
			tokens -= leaked
		}
		lastLeak = now
	}

	allowed := false
	if tokens < lb.Capacity {
		allowed = true
		tokens++
	}

	// Update store with new tokens and leak timestamp
	ttl := int64(60)
	_ = lb.Store.Set(bucketKey, tokens, ttl)
	_ = lb.Store.Set(lastLeakKey, lastLeak, ttl)

	return allowed, nil
}
