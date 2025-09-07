package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/chzyer/readline"
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
		"catch": {
			name:        "catch",
			description: "Attempt to catch a Pokemon. Usage: catch <pokemon-name>",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a caught Pokemon. Usage: inspect <pokemon-name>",
			callback:    commandInspect,
		},
		"pokedx": {
			name:        "pokedx",
			description: "List all caught Pokemon in your Pokedex",
			callback:    commandPokedx,
		},
	}

	// rand.Seed(time.Now().UnixNano())

	interval := time.Duration(time.Second * 10)
	cache := pokecache.NewCache(interval)
	userConfig = pokeapi.Config{
		Next:          "",
		Previous:      "",
		Cache:         cache,
		CaughtPokemon: make(map[string]pokeapi.Pokemon),
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

func commandCatch(c *pokeapi.Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a Pokemon name. Usage: catch <pokemon-name>")
	}

	pokemonName := args[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	if _, exists := c.CaughtPokemon[pokemonName]; exists {
		fmt.Printf("You have already caught %s!\n", pokemonName)
		return nil
	}

	pokemon, err := pokeapi.GetPokemon(c, pokemonName)
	if err != nil {
		return fmt.Errorf("Error getting Pokemon data: %w", err)
	}

	catchThreshold := 50
	if pokemon.BaseExperience > 0 {
		catchThreshold = 60 - (pokemon.BaseExperience / 4)
		if catchThreshold < 10 {
			catchThreshold = 10
		}
	}

	roll := rand.Intn(100)

	if roll < catchThreshold {
		fmt.Printf("%s was caught!\n", pokemonName)
		c.CaughtPokemon[pokemonName] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

func commandInspect(c *pokeapi.Config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("you must provide a Pokemon name. Usage: inspect <pokemon-name>")
	}

	pokemonName := args[0]

	pokemon, exists := c.CaughtPokemon[pokemonName]
	if !exists {
		return fmt.Errorf("you have not caught that pokemon")
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Printf("Stats:\n")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, typeInfo := range pokemon.Types {
		fmt.Printf("  - %s\n", typeInfo.Type.Name)
	}

	return nil
}

func commandPokedx(c *pokeapi.Config, args ...string) error {
	fmt.Println("Your Pokedx:")

	if len(c.CaughtPokemon) == 0 {
		fmt.Println("  - You haven't caught any Pokemon yet!")
		return nil
	}

	for name := range c.CaughtPokemon {
		fmt.Printf("  - %s\n", name)
	}

	return nil
}

func main() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:            INTRO_STRING,
		HistoryFile:       "/tmp/.pokedex_history",
		AutoComplete:      completer,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
		HistorySearchFold: true,
		FuncFilterInputRune: func(r rune) (rune, bool) {
			switch r {
			case readline.CharCtrlZ:
				return r, false
			}
			return r, true
		},
	})
	if err != nil {
		fmt.Printf("Error initializing readline: %v\n", err)
		fmt.Println("Falling back to basic input mode...")
		fallbackInputLoop()
		return
	}
	defer rl.Close()

	fmt.Println(WELCOME_STRING)
	fmt.Println("Use UP/DOWN arrows to navigate command history, TAB for autocomplete")

	for {
		line, err := rl.Readline()
		if err != nil {
			if err == readline.ErrInterrupt {
				fmt.Println("\nUse 'exit' to quit")
				continue
			}
			break
		}

		userInput := cleanInput(line)
		if len(userInput) == 0 {
			continue
		}

		commandToExpect := userInput[0]
		if command, ok := supportedCommands[commandToExpect]; !ok {
			fmt.Println("Command not found. Type 'help' to see available commands")
		} else {
			args := userInput[1:]
			if err := supportedCommands[command.name].callback(&userConfig, args...); err != nil {
				fmt.Println(err)
			}
		}
	}
}

func fallbackInputLoop() {
	fmt.Println("Note: Command history and autocomplete not available in fallback mode")

	for {
		fmt.Print(INTRO_STRING)
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			continue
		}

		userInput := cleanInput(input)
		if len(userInput) == 0 {
			continue
		}

		commandToExpect := userInput[0]
		if command, ok := supportedCommands[commandToExpect]; !ok {
			fmt.Println("Command not found. Type 'help' to see available commands")
		} else {
			args := userInput[1:]
			if err := supportedCommands[command.name].callback(&userConfig, args...); err != nil {
				fmt.Println(err)
			}
		}
	}
}

var completer = readline.NewPrefixCompleter(
	readline.PcItem("help"),
	readline.PcItem("exit"),
	readline.PcItem("map"),
	readline.PcItem("mapb"),
	readline.PcItem("explore",
		readline.PcItem("pastoria-city-area"),
		readline.PcItem("canalave-city-area"),
		readline.PcItem("eterna-city-area"),
		readline.PcItem("floaroma-town-area"),
		readline.PcItem("solaceon-town-area"),
	),
	readline.PcItem("catch",
		readline.PcItem("pikachu"),
		readline.PcItem("charizard"),
		readline.PcItem("blastoise"),
		readline.PcItem("venusaur"),
		readline.PcItem("mewtwo"),
		readline.PcItem("mew"),
		readline.PcItem("magikarp"),
		readline.PcItem("gyarados"),
		readline.PcItem("tentacool"),
		readline.PcItem("tentacruel"),
		readline.PcItem("remoraid"),
		readline.PcItem("octillery"),
		readline.PcItem("wingull"),
		readline.PcItem("pelipper"),
		readline.PcItem("shellos"),
		readline.PcItem("gastrodon"),
	),
	readline.PcItem("inspect",
		readline.PcItem("pikachu"),
		readline.PcItem("charizard"),
		readline.PcItem("blastoise"),
		readline.PcItem("venusaur"),
		readline.PcItem("mewtwo"),
		readline.PcItem("mew"),
		readline.PcItem("magikarp"),
		readline.PcItem("gyarados"),
	),
	readline.PcItem("pokedx"),
)
