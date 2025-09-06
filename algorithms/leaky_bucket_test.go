package algorithms

import (
	"testing"
	"time"

	"github.com/n0l3r/limitron/store"
)

func TestLeakyBucket(t *testing.T) {
	memStore := store.NewMemoryStore()
	lb := &LeakyBucket{
		LeakRate: 1, // 1 token per second
		Capacity: 3, // max 3 tokens
		Store:    memStore,
	}

	key := "test-leakybucket"

	// Bucket is empty at start, allow up to capacity
	for i := int64(1); i <= lb.Capacity; i++ {
		allowed, err := lb.Allow(key)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !allowed {
			t.Errorf("Attempt %d should be allowed", i)
		}
	}

	// Next attempt should be denied (bucket full)
	allowed, err := lb.Allow(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if allowed {
		t.Errorf("Attempt after capacity should NOT be allowed")
	}

	// Wait for leak (1 second per token)
	time.Sleep(1100 * time.Millisecond)
	allowed, err = lb.Allow(key)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !allowed {
		t.Errorf("Attempt after leak should be allowed")
	}

	// Fill bucket again
	for i := int64(1); i < lb.Capacity; i++ {
		allowed, _ = lb.Allow(key)
	}
	// Denied again
	allowed, _ = lb.Allow(key)
	if allowed {
		t.Errorf("Attempt after refilling should NOT be allowed")
	}
}
