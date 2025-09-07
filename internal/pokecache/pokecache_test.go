package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
		{
			key: "https://pokeapi.co/api/v2/location-area",
			val: []byte(`{"count":1010,"next":"https://pokeapi.co/api/v2/location-area/?offset=20&limit=20","previous":null,"results":[{"name":"canalave-city-area","url":"https://pokeapi.co/api/v2/location-area/1/"}]}`),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestGetNonexistentKey(t *testing.T) {
	cache := NewCache(5 * time.Second)

	val, ok := cache.Get("nonexistent-key")
	if ok {
		t.Errorf("expected to not find nonexistent key")
		return
	}
	if val != nil {
		t.Errorf("expected nil value for nonexistent key, got %v", val)
		return
	}
}

func TestCacheUpdate(t *testing.T) {
	cache := NewCache(5 * time.Second)
	key := "https://example.com"

	initialVal := []byte("initial data")
	cache.Add(key, initialVal)

	val, ok := cache.Get(key)
	if !ok {
		t.Errorf("expected to find key after first add")
		return
	}
	if string(val) != string(initialVal) {
		t.Errorf("expected initial value, got %s", string(val))
		return
	}

	updatedVal := []byte("updated data")
	cache.Add(key, updatedVal)

	val, ok = cache.Get(key)
	if !ok {
		t.Errorf("expected to find key after update")
		return
	}
	if string(val) != string(updatedVal) {
		t.Errorf("expected updated value, got %s", string(val))
		return
	}
}

func TestMultipleEntries(t *testing.T) {
	cache := NewCache(5 * time.Second)

	entries := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
		"key3": []byte("value3"),
	}

	for key, val := range entries {
		cache.Add(key, val)
	}

	for key, expectedVal := range entries {
		val, ok := cache.Get(key)
		if !ok {
			t.Errorf("expected to find key %s", key)
			return
		}
		if string(val) != string(expectedVal) {
			t.Errorf("for key %s, expected %s, got %s", key, string(expectedVal), string(val))
			return
		}
	}
}

func TestReapLoopMultipleEntries(t *testing.T) {
	const shortInterval = 10 * time.Millisecond
	cache := NewCache(shortInterval)

	cache.Add("key1", []byte("value1"))

	time.Sleep(shortInterval / 2)

	cache.Add("key2", []byte("value2"))

	time.Sleep(shortInterval)

	_, ok1 := cache.Get("key1")
	if ok1 {
		t.Errorf("expected key1 to be expired")
		return
	}


	time.Sleep(shortInterval)

	_, ok2 := cache.Get("key2")
	if ok2 {
		t.Errorf("expected key2 to be expired after sufficient time")
		return
	}
}

func TestEmptyCache(t *testing.T) {
	cache := NewCache(5 * time.Second)

	val, ok := cache.Get("any-key")
	if ok {
		t.Errorf("expected empty cache to not contain any keys")
		return
	}
	if val != nil {
		t.Errorf("expected nil value from empty cache")
		return
	}
}

func TestCacheCreation(t *testing.T) {
	intervals := []time.Duration{
		1 * time.Second,
		5 * time.Second,
		10 * time.Second,
		100 * time.Millisecond,
	}

	for _, interval := range intervals {
		t.Run(fmt.Sprintf("Interval_%v", interval), func(t *testing.T) {
			cache := NewCache(interval)
			if cache == nil {
				t.Errorf("NewCache should not return nil")
				return
			}

			cache.Add("test", []byte("data"))
			val, ok := cache.Get("test")
			if !ok || string(val) != "data" {
				t.Errorf("cache should be functional after creation")
				return
			}
		})
	}
}
