package main

import (
	"testing"
	"time"

	"github.com/fyzanshaik/pokedex/internal/pokeapi"
	"github.com/fyzanshaik/pokedex/internal/pokecache"
)

type testCase struct {
	input    string
	expected []string
}

func TestCleanInput(t *testing.T) {
	cases := []testCase{
		{input: " hello  world  ", expected: []string{"hello", "world"}},
		{input: "one", expected: []string{"one"}},
		{input: "  ", expected: []string{}},
		{input: "HELLO WORLD", expected: []string{"hello", "world"}},
		{input: "Map", expected: []string{"map"}},
		{input: "EXIT", expected: []string{"exit"}},
		{input: "help   me   please", expected: []string{"help", "me", "please"}},
		{input: "", expected: []string{}},
		{input: "\n\t  \r", expected: []string{}},
	}

	for _, tc := range cases {
		actual := cleanInput(tc.input)
		if len(actual) != len(tc.expected) {
			t.Errorf("input %q: expected len %d, got %d", tc.input, len(tc.expected), len(actual))
			continue
		}
		for i := range actual {
			if actual[i] != tc.expected[i] {
				t.Errorf("input %q: at %d, got %q, want %q", tc.input, i, actual[i], tc.expected[i])
				break
			}
		}
	}
}

func TestCommandHelp(t *testing.T) {
	cache := pokecache.NewCache(5 * time.Second)
	config := &pokeapi.Config{
		Next:     "",
		Previous: "",
		Cache:    cache,
	}

	err := commandHelp(config)
	if err != nil {
		t.Errorf("commandHelp should not return an error, got: %v", err)
	}
}

func TestSupportedCommandsMap(t *testing.T) {
	expectedCommands := []string{"exit", "help", "map", "mapb"}

	for _, cmdName := range expectedCommands {
		if _, exists := supportedCommands[cmdName]; !exists {
			t.Errorf("expected command %q to exist in supportedCommands", cmdName)
		}
	}

	// Test that each command has required fields
	for cmdName, cmd := range supportedCommands {
		if cmd.name == "" {
			t.Errorf("command %q should have a name", cmdName)
		}
		if cmd.description == "" {
			t.Errorf("command %q should have a description", cmdName)
		}
		if cmd.callback == nil {
			t.Errorf("command %q should have a callback function", cmdName)
		}
		if cmd.name != cmdName {
			t.Errorf("command map key %q should match command name %q", cmdName, cmd.name)
		}
	}
}

func TestPrintLocations(t *testing.T) {
	testLocations := []pokeapi.Result{
		{Name: "location-1", URL: "https://example.com/1"},
		{Name: "location-2", URL: "https://example.com/2"},
		{Name: "location-3", URL: "https://example.com/3"},
	}

	// This test mainly ensures the function doesn't panic
	// In a real test environment, you might want to capture stdout
	printLocations(testLocations)
	printLocations([]pokeapi.Result{}) // Test empty slice
}

func TestConfigInitialization(t *testing.T) {
	// Test that userConfig is properly initialized
	if userConfig.Cache == nil {
		t.Errorf("userConfig.Cache should not be nil after initialization")
	}

	if userConfig.Next != "" {
		t.Errorf("userConfig.Next should be empty string initially, got %q", userConfig.Next)
	}

	if userConfig.Previous != "" {
		t.Errorf("userConfig.Previous should be empty string initially, got %q", userConfig.Previous)
	}
}

func TestCommandCallbacks(t *testing.T) {
	cache := pokecache.NewCache(5 * time.Second)
	config := &pokeapi.Config{
		Next:     "",
		Previous: "",
		Cache:    cache,
	}

	// Test help command (should not error)
	err := commandHelp(config)
	if err != nil {
		t.Errorf("commandHelp should not return error, got: %v", err)
	}

	// Test mapb command with no previous (should error)
	err = commandMapBack(config)
	if err == nil {
		t.Errorf("commandMapBack should return error when no previous location exists")
	}
}

func TestCleanInputEdgeCases(t *testing.T) {
	edgeCases := []testCase{
		{input: "   ", expected: []string{}},
		{input: "\t\n", expected: []string{}},
		{input: "a", expected: []string{"a"}},
		{input: "A B C D E", expected: []string{"a", "b", "c", "d", "e"}},
		{input: "multiple    spaces    between", expected: []string{"multiple", "spaces", "between"}},
	}

	for _, tc := range edgeCases {
		actual := cleanInput(tc.input)
		if len(actual) != len(tc.expected) {
			t.Errorf("input %q: expected len %d, got %d", tc.input, len(tc.expected), len(actual))
			continue
		}
		for i := range actual {
			if actual[i] != tc.expected[i] {
				t.Errorf("input %q: at %d, got %q, want %q", tc.input, i, actual[i], tc.expected[i])
				break
			}
		}
	}
}
