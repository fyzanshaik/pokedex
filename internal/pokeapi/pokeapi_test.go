package pokeapi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/fyzanshaik/pokedex/internal/pokecache"
)

// Mock response data for testing
var mockLocationAreaResponse = `{
	"count": 1010,
	"next": "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20",
	"previous": null,
	"results": [
		{
			"name": "canalave-city-area",
			"url": "https://pokeapi.co/api/v2/location-area/1/"
		},
		{
			"name": "eterna-city-area",
			"url": "https://pokeapi.co/api/v2/location-area/2/"
		}
	]
}`

var mockLocationAreaResponsePage2 = `{
	"count": 1010,
	"next": "https://pokeapi.co/api/v2/location-area/?offset=40&limit=20",
	"previous": "https://pokeapi.co/api/v2/location-area/",
	"results": [
		{
			"name": "floaroma-town-area",
			"url": "https://pokeapi.co/api/v2/location-area/21/"
		},
		{
			"name": "solaceon-town-area",
			"url": "https://pokeapi.co/api/v2/location-area/22/"
		}
	]
}`

func TestGetNextLocations(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, mockLocationAreaResponse)
	}))
	defer server.Close()

	// Create cache and config
	cache := pokecache.NewCache(5 * time.Second)
	config := &Config{
		Next:     server.URL,
		Previous: "",
		Cache:    cache,
	}

	// Test getting locations
	locations, err := GetNextLocations(config)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	// Verify response structure
	if len(locations.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(locations.Results))
		return
	}

	if locations.Results[0].Name != "canalave-city-area" {
		t.Errorf("expected first result name to be 'canalave-city-area', got %s", locations.Results[0].Name)
		return
	}

	// Verify config was updated
	if config.Next != "https://pokeapi.co/api/v2/location-area/?offset=20&limit=20" {
		t.Errorf("expected config.Next to be updated")
		return
	}
}

func TestGetNextLocationsWithCache(t *testing.T) {
	requestCount := 0

	// Create a test server that counts requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, mockLocationAreaResponse)
	}))
	defer server.Close()

	// Create cache and config
	cache := pokecache.NewCache(5 * time.Second)
	config := &Config{
		Next:     server.URL,
		Previous: "",
		Cache:    cache,
	}

	// First request - should hit the server
	_, err := GetNextLocations(config)
	if err != nil {
		t.Errorf("expected no error on first request, got %v", err)
		return
	}

	if requestCount != 1 {
		t.Errorf("expected 1 request to server, got %d", requestCount)
		return
	}

	// Reset config to make the same request again
	config.Next = server.URL

	// Second request - should hit the cache
	_, err = GetNextLocations(config)
	if err != nil {
		t.Errorf("expected no error on second request, got %v", err)
		return
	}

	// Request count should still be 1 (cached)
	if requestCount != 1 {
		t.Errorf("expected request count to stay 1 (cached), got %d", requestCount)
		return
	}
}

func TestGetPrevLocations(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, mockLocationAreaResponsePage2)
	}))
	defer server.Close()

	// Create cache and config with Previous set
	cache := pokecache.NewCache(5 * time.Second)
	config := &Config{
		Next:     "",
		Previous: server.URL,
		Cache:    cache,
	}

	// Test getting previous locations
	locations, err := GetPrevLocations(config)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
		return
	}

	// Verify response structure
	if len(locations.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(locations.Results))
		return
	}

	if locations.Results[0].Name != "floaroma-town-area" {
		t.Errorf("expected first result name to be 'floaroma-town-area', got %s", locations.Results[0].Name)
		return
	}
}

func TestGetPrevLocationsNoPrevious(t *testing.T) {
	// Create cache and config without Previous set
	cache := pokecache.NewCache(5 * time.Second)
	config := &Config{
		Next:     "",
		Previous: "",
		Cache:    cache,
	}

	// Should return error when no previous exists
	_, err := GetPrevLocations(config)
	if err == nil {
		t.Errorf("expected error when no previous location exists")
		return
	}

	expectedError := "You are at the first location!"
	if err.Error() != expectedError {
		t.Errorf("expected error message '%s', got '%s'", expectedError, err.Error())
		return
	}
}

func TestGetPrevLocationsWithCache(t *testing.T) {
	requestCount := 0

	// Create a test server that counts requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, mockLocationAreaResponsePage2)
	}))
	defer server.Close()

	// Create cache and config
	cache := pokecache.NewCache(5 * time.Second)
	config := &Config{
		Next:     "",
		Previous: server.URL,
		Cache:    cache,
	}

	// First request - should hit the server
	_, err := GetPrevLocations(config)
	if err != nil {
		t.Errorf("expected no error on first request, got %v", err)
		return
	}

	if requestCount != 1 {
		t.Errorf("expected 1 request to server, got %d", requestCount)
		return
	}

	// Reset config to make the same request again
	config.Previous = server.URL

	// Second request - should hit the cache
	_, err = GetPrevLocations(config)
	if err != nil {
		t.Errorf("expected no error on second request, got %v", err)
		return
	}

	// Request count should still be 1 (cached)
	if requestCount != 1 {
		t.Errorf("expected request count to stay 1 (cached), got %d", requestCount)
		return
	}
}

func TestCacheExpiration(t *testing.T) {
	requestCount := 0

	// Create a test server that counts requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, mockLocationAreaResponse)
	}))
	defer server.Close()

	// Create cache with very short expiration
	cache := pokecache.NewCache(50 * time.Millisecond)
	config := &Config{
		Next:     server.URL,
		Previous: "",
		Cache:    cache,
	}

	// First request
	_, err := GetNextLocations(config)
	if err != nil {
		t.Errorf("expected no error on first request, got %v", err)
		return
	}

	if requestCount != 1 {
		t.Errorf("expected 1 request to server, got %d", requestCount)
		return
	}

	// Wait for cache to expire
	time.Sleep(100 * time.Millisecond)

	// Reset config to make the same request again
	config.Next = server.URL

	// Second request after expiration - should hit the server again
	_, err = GetNextLocations(config)
	if err != nil {
		t.Errorf("expected no error on second request, got %v", err)
		return
	}

	// Request count should be 2 (cache expired)
	if requestCount != 2 {
		t.Errorf("expected request count to be 2 (cache expired), got %d", requestCount)
		return
	}
}

func TestJSONUnmarshalingFromCache(t *testing.T) {
	// Create cache and manually add JSON data
	cache := pokecache.NewCache(5 * time.Second)
	testURL := "https://test.com/api"

	cache.Add(testURL, []byte(mockLocationAreaResponse))

	config := &Config{
		Next:     testURL,
		Previous: "",
		Cache:    cache,
	}

	// This should use cached data and unmarshal it correctly
	locations, err := GetNextLocations(config)
	if err != nil {
		t.Errorf("expected no error when using cached data, got %v", err)
		return
	}

	// Verify the data was unmarshaled correctly
	if len(locations.Results) != 2 {
		t.Errorf("expected 2 results from cached data, got %d", len(locations.Results))
		return
	}

	if locations.Count != 1010 {
		t.Errorf("expected count 1010 from cached data, got %d", locations.Count)
		return
	}
}

func TestNetworkError(t *testing.T) {
	// Create cache and config with invalid URL
	cache := pokecache.NewCache(5 * time.Second)
	invalidURL := "http://invalid-url-that-does-not-exist-12345.com"
	config := &Config{
		Next:     invalidURL,
		Previous: "",
		Cache:    cache,
	}

	// Should return network error
	_, err := GetNextLocations(config)
	if err == nil {
		t.Errorf("expected network error for invalid URL")
		return
	}

	// Verify it's actually a network error, not a cache error
	if !strings.Contains(err.Error(), "network request") {
		t.Errorf("expected network request error, got: %v", err)
		return
	}
}
