package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	text = strings.ToLower(text)
	textSlice := strings.Fields(text)
	return textSlice
}

type cliCommands struct {
	name        string
	description string
	callback    func() error
}

const INTRO_STRING string = "Pokedex > "
const USER_INPUT_PREFIX string = "Your command was: "
const WELCOME_STRING string = "Welcome to the Pokedex!"

var supportedCommands map[string]cliCommands

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
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Usage: \n")
	for cmdName, cmd := range supportedCommands {
		fmt.Printf("- %s: %s\n", cmdName, cmd.description)
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

		commandToExpect := userInput[0]
		if command, ok := supportedCommands[commandToExpect]; ok == false {
			fmt.Println("Command not found check listed commands through 'usage'")
		} else {
			if err := supportedCommands[command.name].callback(); err != nil {
				fmt.Println("Error running the command: %w", err)
			}
		}

	}

}
