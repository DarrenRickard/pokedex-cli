package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"github.com/darrenrickard/pokedexcli/internal/pokeapi"
	
)

// Data

type cliCommand struct {
	name		string
	description	string
	callback	func() error
}

// End data

func main() {
	// Initialize CLI Commands
	commands := map[string]cliCommand {
		"exit": {
			name:			"exit",
			description:	"Exit the Pokedex",
			callback:		commandExit,
		},
		"help": {
			name:			"help",
			description: 	"Displays a help message",
			callback:		nil,
		},
		"map": {
			name: 			"map",
			description: 	"Displays names of location areas in the Pokemon world",
			callback: 		commandMap,
		},
	}
	// Assigns "help" callback function after map has been initialized to access the map in callback
	commands["help"] = cliCommand {
		name: 			"help",
		description: 	"Displays a help message",
		callback: 		func() error {return commandHelp(commands)},
	}

	// Main program loop
	scanner := bufio.NewScanner(os.Stdin)
	for {
		var words []string
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		words = cleanInput(userInput)
		first := strings.ToLower(words[0])
		cmd, exists := commands[first]
		if exists {
			if err := cmd.callback(); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}	
	} 
	// End main program loop
}

// Functions 

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(commands map[string]cliCommand) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n ")
	for _, v := range commands {
		msg := fmt.Sprintf("%s: %s", v.name, v.description)
		fmt.Println(msg)
	}
	return nil
}

func commandMap() error {
	url := "https://pokeapi.co/api/v2/location-area/"
	locations, err := pokeapi.FetchPokeLocations(url)
	if err != nil {
		fmt.Println(err)
	}
	for _, l := range locations {
		fmt.Println(l)
	}
	return nil
}

func cleanInput(text string) []string {
	cleanedInput := strings.Fields(text)	
	return cleanedInput
}

// End functions
