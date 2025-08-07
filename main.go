package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/darrenrickard/pokedexcli/internal/pokeapi"
	"github.com/darrenrickard/pokedexcli/internal/pokecache"
)

// Data

type cliCommand struct {
	name		string
	description	string
	callback	func() error
}


// End data

func main() {

	// Initialize cache before commands so it can be used in callback functions
	cache := pokecache.NewCache(5*time.Second)

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
		description: 	"Displays next 20 names of location areas in the Pokemon world",
		callback: 		func() error {return commandMap(&pokeapi.Links, &cache)},
	}
	commands["mapb"] = cliCommand{
		name: 			"mapb",
		description: 	"Displays previous 20 names of location areas in the Pokemon world",
		callback: 		func() error {return commandMapb(&pokeapi.Links, &cache)},
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


func commandMap(links* pokeapi.PageLinks, c* pokecache.Cache) error {
	// if links.Next exists in cache.entries { pass cache.entries[links.Next] to a function that unmarshals []byte and returns []string }
	_, exists := c.Entries[links.Next]
	// if link exists in cache, just print cache entry's data 
	var locations []string
	// the entry might exist at evaluation time but when c.Get(links.Next) is called, it might be gone already. refactor this code
	if exists {
		fmt.Println("Getting page from cache!")
		// this functionality can be refactored. doesnt need previous link, exists to determine if the link exists in cache.
		// c1 <- c.Get(links.Next)
		// cachedList, _ := <- c1
		// l, err := pokeapi.UnmarshalToList(cachedList.Data)
		var cachedLocation pokecache.CacheEntryData
		go func() {
			cachedLocation = c.Get(links.Next)		
		}()
		l, err := pokeapi.UnmarshalToList(links.Next, cachedLocation.Data)
		if err != nil {
			fmt.Println(err)
		}
		locations = l
	} else { // else send GET request for data

		l, err := pokeapi.FetchPokeLocations(links.Next, c)	
		if err != nil {
			fmt.Println(err)
		}
		locations = l
	}
	printLocations(locations)
	fmt.Printf("Current: %s\nNext: %s\nPrevious: %s", links.Current, links.Next, links.Previous)
	return nil
}

func printLocations(locations []string) {
	for _, l := range locations {
		fmt.Println(l)
	}
}

func commandMapb(links* pokeapi.PageLinks, c* pokecache.Cache) error {
	if links.Previous == "" {
		fmt.Printf("Current: %s\nNext: %s\nPrevious: %s", links.Current, links.Next, links.Previous)
		return fmt.Errorf("you're on the first page")
	}
	// if links.Next exists in cache.entries { pass cache.entries[links.Next] to a function that unmarshals []byte and returns []string }
	_, exists := c.Entries[links.Previous]
	// if link exists in cache, just print cache entry's data 
	var locations []string
	// c1 := make(chan pokecache.CacheEntryData)
	if exists {
		fmt.Println("Getting page from cache!")
		// this functionality can be refactored. doesnt need previous link, exists to determine if the link exists in cache.
		// c1 <- c.Get(links.Next)
		// cachedList, _ := <- c1
		// l, err := pokeapi.UnmarshalToList(cachedList.Data)
		var cachedLocation pokecache.CacheEntryData
		cachedLocation = c.Get(links.Previous)		
		l, err := pokeapi.UnmarshalToList(links.Previous, cachedLocation.Data)
		if err != nil {
			fmt.Println(err)
		}
		locations = l
	} else { // else send GET request for data

		l, err := pokeapi.FetchPokeLocations(links.Previous, c)	
		if err != nil {
			fmt.Println(err)
		}
		locations = l
	}
	printLocations(locations)
	fmt.Printf("Current: %s\nNext: %s\nPrevious: %s", links.Current, links.Next, links.Previous)
	return nil
	// if links.Previous exists in cache.entries 
	// link, exists := c.Entries[links.Previous]
	// if link exists in cache, just use cache entry's data 
	// var locations []string
	// if exists {
	// 	l, err := pokeapi.UnmarshalToList(link.Val)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	locations = l
	// } else { // else send GET request for data
	//
	// 	l, err := pokeapi.FetchPokeLocations(links.Previous, c)	
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	locations = l
	// }
	// printLocations(locations)
	// return nil
}

func cleanInput(text string) []string {
	cleanedInput := strings.Fields(text)	
	return cleanedInput
}

// End functions
