package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	commands := map[string]cliCommand {
		"exit": {
			name:			"exit",
			description:	"Exit the Pokedex",
			callback:		commandExit,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		var words []string
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		words = cleanInput(userInput)
		first := strings.ToLower(words[0])
		for _, cmd := range commands {
			if cmd.name == first {
				cmd.callback()
			}
		}
		fmt.Printf("Your command was: %s\n", first)
	} 



}

type cliCommand struct {
	name		string
	description	string
	callback	func() error
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func cleanInput(text string) []string {
	// cleanedInput := []string{}
	cleanedInput := strings.Fields(text)	
	return cleanedInput
}

