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

const INTRO_STRING string = "Pokedex > "
const USER_INPUT_PREFIX string = "Your command was: "

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(INTRO_STRING)
		scanner.Scan()
		userInput := cleanInput(scanner.Text())[0]

	}

}
