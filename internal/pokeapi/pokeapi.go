package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fyzanshaik/pokedex/internal/pokecache"
)

const BASE_URL string = "https://pokeapi.co/api/v2"

type Config struct {
	Next     string
	Previous string
	Cache    *pokecache.Cache
}

type Result struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationArea struct {
	Count    int      `json:"count"`
	Next     string   `json:"next"`
	Previous any      `json:"previous"`
	Results  []Result `json:"results"`
}

// GET https://pokeapi.co/api/v2/location-area/{id or name}/
func GetNextLocations(c *Config) (LocationArea, error) {
	var currentLocationArea LocationArea
	resourceName := "/location-area"
	full_url := BASE_URL + resourceName

	if c.Next != "" {
		full_url = c.Next
	}

	if cachedData, found := c.Cache.Get(full_url); found {
		fmt.Println("Accessing cache for: ", full_url)
		json.Unmarshal(cachedData, &currentLocationArea)
		return currentLocationArea, nil
	}

	// fmt.Println("Current full url: ", full_url)
	res, err := http.Get(full_url)
	if err != nil {
		return LocationArea{}, fmt.Errorf("Error in network request location-area: %w", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationArea{}, fmt.Errorf("Error reading Body: %w", err)
	}

	c.Cache.Add(full_url, body)
	json.Unmarshal(body, &currentLocationArea)

	if val, ok := currentLocationArea.Previous.(string); ok {
		c.Previous = val
	}

	c.Next = currentLocationArea.Next

	return currentLocationArea, nil
}

// GET https://pokeapi.co/api/v2/location-area/{id or name}/
func GetPrevLocations(c *Config) (LocationArea, error) {

	if c.Previous == "" {
		return LocationArea{}, fmt.Errorf("You are at the first location!")
	}
	var currentLocationArea LocationArea

	full_url := c.Previous

	if cachedData, found := c.Cache.Get(full_url); found {
		json.Unmarshal(cachedData, &currentLocationArea)
		return currentLocationArea, nil
	}

	// fmt.Println("Current full url: ", full_url)
	res, err := http.Get(full_url)
	if err != nil {
		return LocationArea{}, fmt.Errorf("Error in network request location-area: %w", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationArea{}, fmt.Errorf("Error reading Body: %w", err)
	}

	c.Cache.Add(full_url, body)

	json.Unmarshal(body, &currentLocationArea)
	if value, ok := currentLocationArea.Previous.(string); ok {
		c.Previous = value
	}
	c.Next = currentLocationArea.Next

	return currentLocationArea, nil
}
