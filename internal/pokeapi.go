package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const BASE_URL string = "https://pokeapi.co/api/v2"

type Config struct {
	Next     string
	Previous string
}

type LocationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// GET https://pokeapi.co/api/v2/location-area/{id or name}/
func GetNextLocations(c *Config) (LocationArea, error) {
	var currentLocationArea LocationArea
	resourceName := "/location-area"
	full_url := BASE_URL + resourceName

	if c.Next != "" {
		full_url = c.Next
	}
	fmt.Println("Current full url: ", full_url)
	res, err := http.Get(full_url)
	if err != nil {
		return LocationArea{}, fmt.Errorf("Error in network request location-area: %w", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationArea{}, fmt.Errorf("Error reading Body: %w", body)
	}

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

	fmt.Println("Current full url: ", full_url)
	res, err := http.Get(full_url)
	if err != nil {
		return LocationArea{}, fmt.Errorf("Error in network request location-area: %w", err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationArea{}, fmt.Errorf("Error reading Body: %w", body)
	}

	json.Unmarshal(body, &currentLocationArea)
	if value, ok := currentLocationArea.Previous.(string); ok {
		c.Previous = value
	}
	c.Next = currentLocationArea.Next

	return currentLocationArea, nil
}
