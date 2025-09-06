package algorithms

import (
	"testing"

	"github.com/n0l3r/limitron/store"
)

func TestTokenBucket_Allow(t *testing.T) {
	memStore := store.NewMemoryStore()
	tb := &TokenBucket{
		Rate:     5,
		Capacity: 5,
		Store:    memStore,
	}
	key := "test-client"

	// Pre-fill bucket
	tb.Store.Set(key, tb.Capacity, 60)

	// Allow requests up to capacity
	for i := 0; i < int(tb.Capacity); i++ {
		allowed, err := tb.Allow(key)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !allowed {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}
	// Next request should be denied
	allowed, err := tb.Allow(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Fatalf("request over capacity should not be allowed")
	}
}
