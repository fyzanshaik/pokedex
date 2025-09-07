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
	callback    func(c *pokeapi.Config, args ...string) error
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
		"explore": {
			name:        "explore",
			description: "Explore a location area to find Pokemon. Usage: explore <location-name>",
			callback:    commandExplore,
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

func commandMap(c *pokeapi.Config, args ...string) error {
	allLocations, err := pokeapi.GetNextLocations(c)
	// fmt.Println(locationArea)
	if err != nil {
		return fmt.Errorf("Error fetching locations: %w", err)
	}
	locations := allLocations.Results
	printLocations(locations)
	return nil
}

func commandMapBack(c *pokeapi.Config, args ...string) error {
	allLocations, err := pokeapi.GetPrevLocations(c)
	// fmt.Println(locationArea)
	if err != nil {
		return fmt.Errorf("Error fetching locations: %w", err)
	}
	locations := allLocations.Results
	printLocations(locations)
	return nil
}

func commandExit(c *pokeapi.Config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *pokeapi.Config, args ...string) error {
	fmt.Println("Usage:")
	for cmdName, cmd := range supportedCommands {
		fmt.Printf("- %s: %s\n", cmdName, cmd.description)
	}
	fmt.Println()
	return nil
}

func commandExplore(c *pokeapi.Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a location name. Usage: explore <location-name>")
	}

	locationName := args[0]
	fmt.Printf("Exploring %s...\n", locationName)

	locationInfo, err := pokeapi.GetLocationInformation(c, locationName)
	if err != nil {
		return fmt.Errorf("Error exploring location: %w", err)
	}

	if len(locationInfo.PokemonEncounters) == 0 {
		fmt.Println("Found no Pokemon in this area.")
		return nil
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range locationInfo.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

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
			// Pass remaining arguments to the command
			args := userInput[1:]
			if err := supportedCommands[command.name].callback(&userConfig, args...); err != nil {
				fmt.Println(err)
			}
		}

	}
}
