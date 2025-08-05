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
	// Declare CLI Commands map
	commands := map[string]cliCommand {}

	// Initialize CLI Commands. Done this way because cant access function arguments at declaration time
	commands["help"] = cliCommand {
		name: 			"help",
		description: 	"Displays a help message",
		callback: 		func() error {return commandHelp(commands, &pokeapi.Links)},
	}
	commands["exit"] = cliCommand{
		name:			"exit",
		description:	"Exit the Pokedex",
		callback:		func() error {return commandExit(&pokeapi.Links)} ,
	}
	commands["map"] = cliCommand{
		name: 			"map",
		description: 	"Displays names of location areas in the Pokemon world",
		callback: 		func() error {return commandMap(&pokeapi.Links)},
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

func commandExit(links* pokeapi.PageLinks) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(commands map[string]cliCommand, links* pokeapi.PageLinks) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:\n ")
	for _, v := range commands {
		msg := fmt.Sprintf("%s: %s", v.name, v.description)
		fmt.Println(msg)
	}
	return nil
}

func commandMap(links* pokeapi.PageLinks) error {
	err := pokeapi.FetchLocationPageLinks(links.Current)
	if err != nil {
		fmt.Println(err)
	}
	locations, err := pokeapi.FetchPokeLocations(links.Current)	
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
