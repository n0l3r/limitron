package store

import (
	"testing"
	"time"
)

func TestMemoryStore_SetGet(t *testing.T) {
	store := NewMemoryStore()
	key := "foo"
	err := store.Set(key, 123, 1)
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	val, err := store.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != 123 {
		t.Fatalf("expected 123, got %d", val)
	}
}

func TestMemoryStore_Expiry(t *testing.T) {
	store := NewMemoryStore()
	key := "expiry"
	store.Set(key, 55, 1)
	time.Sleep(2 * time.Second)
	val, err := store.Get(key)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != 0 {
		t.Fatalf("expected expired value to be 0, got %d", val)
	}
}

func TestMemoryStore_Incr(t *testing.T) {
	store := NewMemoryStore()
	key := "incr"
	v, err := store.Incr(key, 1, 1)
	if err != nil {
		t.Fatalf("Incr failed: %v", err)
	}
	if v != 1 {
		t.Fatalf("expected 1, got %d", v)
	}
	v, err = store.Incr(key, 5, 1)
	if err != nil {
		t.Fatalf("Incr failed: %v", err)
	}
	if v != 6 {
		t.Fatalf("expected 6, got %d", v)
	}
}
