package algorithms

import (
	"testing"

	"github.com/n0l3r/limitron/store"
)

func TestFixedWindow_Allow(t *testing.T) {
	memStore := store.NewMemoryStore()
	fw := &FixedWindow{
		WindowSize: 10,
		Limit:      2,
		Store:      memStore,
	}
	key := "test-fw"

	for i := 0; i < int(fw.Limit); i++ {
		allowed, err := fw.Allow(key)
		if err != nil {
			t.Fatalf("fixed window error: %v", err)
		}
		if !allowed {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}
	allowed, _ := fw.Allow(key)
	if allowed {
		t.Fatalf("request exceeding window limit should not be allowed")
	}
}
