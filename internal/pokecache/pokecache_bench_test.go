package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkCacheAdd(b *testing.B) {
	cache := NewCache(5 * time.Second)
	testData := []byte("test data for benchmarking")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.Add(key, testData)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cache := NewCache(5 * time.Second)
	testData := []byte("test data for benchmarking")

	// Pre-populate cache with test data
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.Add(key, testData)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i%1000)
		cache.Get(key)
	}
}

func BenchmarkCacheAddGet(b *testing.B) {
	cache := NewCache(5 * time.Second)
	testData := []byte("test data for benchmarking")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key-%d", i)
		cache.Add(key, testData)
		cache.Get(key)
	}
}

func BenchmarkCacheMiss(b *testing.B) {
	cache := NewCache(5 * time.Second)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("nonexistent-key-%d", i)
		cache.Get(key)
	}
}

func BenchmarkCacheWithContention(b *testing.B) {
	cache := NewCache(5 * time.Second)
	testData := []byte("test data for concurrent access")

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			key := fmt.Sprintf("concurrent-key-%d", i%100)
			if i%2 == 0 {
				cache.Add(key, testData)
			} else {
				cache.Get(key)
			}
			i++
		}
	})
}

func BenchmarkLargeDataCache(b *testing.B) {
	cache := NewCache(5 * time.Second)

	// Create large test data (10KB)
	largeData := make([]byte, 10*1024)
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("large-key-%d", i)
		cache.Add(key, largeData)
		cache.Get(key)
	}
}

func BenchmarkCacheEviction(b *testing.B) {
	// Very short interval to trigger frequent eviction
	cache := NewCache(1 * time.Millisecond)
	testData := []byte("test data that will be evicted quickly")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("eviction-key-%d", i)
		cache.Add(key, testData)

		// Small delay to allow some entries to expire
		if i%100 == 0 {
			time.Sleep(2 * time.Millisecond)
		}

		cache.Get(key)
	}
}
