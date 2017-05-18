package main

import (
	"testing"
	"time"
)

func TestCachem(t *testing.T) {
	item := &Item{
		url:       "test",
		html:      "",
		updatedAt: time.Now().Add(-20 * time.Minute),
	}
	cache := newCache()
	cache.add(item)
	if len(cache.queue) != len(cache.entries) {
		t.Fatalf("Cache inconsistent: %d(queue) vs %d(entries)", len(cache.queue), len(cache.entries))
	}
	cache.cleanUp(100)
	if len(cache.queue) != len(cache.entries) {
		t.Fatalf("Cache inconsistent: %d(queue) vs %d(entries)", len(cache.queue), len(cache.entries))
	}
	if len(cache.queue) != 1 {
		t.Fatalf("An entry shoud not be deleted.")
	}
	cache.cleanUp(0)
	if len(cache.queue) != 0 {
		t.Fatalf("An entry should be deleted.")
	}
}
