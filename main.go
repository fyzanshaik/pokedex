package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fyzanshaik/pokedex/internal/pokeapi"
	"github.com/fyzanshaik/pokedex/internal/pokecache"
)

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	textSlice := strings.Fields(text)
	return textSlice
}

type cliCommands struct {
	name        string
	description string
	callback    func(c *pokeapi.Config) error
}

const INTRO_STRING string = "Pokedex > "
const USER_INPUT_PREFIX string = "Your command was: "
const WELCOME_STRING string = "Welcome to the Pokedex!"

var supportedCommands map[string]cliCommands
var userConfig pokeapi.Config

func init() {

	supportedCommands = map[string]cliCommands{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays 20 locations to explore, each subsequent request displays the next set",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 locations if it exists",
			callback:    commandMapBack,
		},
	}
	interval := time.Duration(time.Second * 10)
	cache := pokecache.NewCache(interval)
	userConfig = pokeapi.Config{
		Next:     "",
		Previous: "",
		Cache:    cache,
	}
}

func printLocations(locations []pokeapi.Result) {
	for i := 0; i < len(locations); i++ {
		fmt.Printf("%d => %s\n", i+1, locations[i].Name)
	}
	fmt.Println()

}

func commandMap(c *pokeapi.Config) error {
	allLocations, err := pokeapi.GetNextLocations(c)
	// fmt.Println(locationArea)
	if err != nil {
		return fmt.Errorf("Error fetching locations: %w", err)
	}
	locations := allLocations.Results
	printLocations(locations)
	return nil
}

func commandMapBack(c *pokeapi.Config) error {
	allLocations, err := pokeapi.GetPrevLocations(c)
	// fmt.Println(locationArea)
	if err != nil {
		return fmt.Errorf("Error fetching locations: %w", err)
	}
	locations := allLocations.Results
	printLocations(locations)
	return nil
}

func commandExit(c *pokeapi.Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *pokeapi.Config) error {
	fmt.Println("Usage:")
	for cmdName, cmd := range supportedCommands {
		fmt.Printf("- %s: %s\n", cmdName, cmd.description)
	}
	fmt.Println()
	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(WELCOME_STRING)
	for {
		fmt.Print(INTRO_STRING)
		scanner.Scan()
		userInput := cleanInput(scanner.Text())
		if len(userInput) == 0 {
			continue
		}
		commandToExpect := userInput[0]
		if command, ok := supportedCommands[commandToExpect]; ok == false {
			fmt.Println("Command not found check listed commands through 'usage'")
		} else {
			if err := supportedCommands[command.name].callback(&userConfig); err != nil {
				fmt.Println(err)
			}
		}

	}
}
