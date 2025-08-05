package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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
		// for _, cmd := range commands {
		// 	if cmd.name == first {
		// 		cmd.callback()
		// 		continue
		// 	}
		// }
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

type (
	LocationArea struct {
		Count    int                `json:"count"`
		Next     string             `json:"next"`
		Previous string             `json:"previous"`
		Results  []Results			`json:"results"`
	}

	Results struct {
		Name	string		`json:"name"`	
		URL		string		`json:"url"`
	}
)


func commandMap() error {
	var pokemaps LocationArea
	fmt.Println("Creating Request...")
	req, err := http.NewRequest("GET", "https://pokeapi.co/api/v2/location-area/", nil)	
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("limit", "20")
	req.Header.Set("offset", "20")	
	fmt.Println("Sending Request...")
	client := &http.Client{}
	res, err := client.Do(req)
	fmt.Println("Request sent, Response received")
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}
	if res.StatusCode > 299 {
		return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, res.Body)	
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("Error reading resource: %v", err)
	}
	if err := json.Unmarshal(body, &pokemaps); err != nil {
		return fmt.Errorf("Error Unmarshal'ing response body:\n%v", err)
	}
	for _, r := range pokemaps.Results {
		fmt.Println(r.Name)
	}
	fmt.Println("pokemaps printed")
	return nil
}

func cleanInput(text string) []string {
	cleanedInput := strings.Fields(text)	
	return cleanedInput
}

// End functions
