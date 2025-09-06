package algorithms

import (
	"testing"

	"github.com/n0l3r/limitron/store"
)

func TestSlidingWindow_Allow(t *testing.T) {
	memStore := store.NewMemoryStore()
	sw := &SlidingWindow{
		WindowSize: 10,
		Limit:      2,
		Store:      memStore,
	}
	key := "test-sw"

	for i := 0; i < int(sw.Limit); i++ {
		allowed, err := sw.Allow(key)
		if err != nil {
			t.Fatalf("sliding window error: %v", err)
		}
		if !allowed {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}
	allowed, _ := sw.Allow(key)
	if allowed {
		t.Fatalf("request exceeding window limit should not be allowed")
	}
}
