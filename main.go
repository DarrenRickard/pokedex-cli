package main

import (
	"bufio"
	"fmt"
	"os"
	"math/rand"
	"strings"
	"time"

	"github.com/darrenrickard/pokedexcli/internal/pokeapi"
	"github.com/darrenrickard/pokedexcli/internal/pokecache"
)

// Data

var Pokedex = map[string]pokeapi.Pokemon{}

type cliCommand struct {
	name		string
	description	string
	callback	func() error
}


// End data

func main() {

	// Initialize cache before commands so it can be used in callback functions
	cache := pokecache.NewCache(5*time.Second)
	var arg string

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
	commands["ccache"] = cliCommand{
		name: "ccache",
		description: "Displays the current cache",
		callback:  func() error {return list(&cache)},
	}
	commands["explore"] = cliCommand{
		name: "explore",
		description: "Displays pokemon located in the given location",
		callback: func() error {return explore(&cache, &arg)},
	}
	commands["catch"] = cliCommand{
		name: "catch",
		description: "Throws a pokeball and attempts to catch the Pokemon, adding it to Poxedex upon success",
		callback: func() error {return catch(&arg)},
	}
	commands["inspect"] = cliCommand{
		name: "inspect",
		description: "Inspect caught pokemon in pokedex",
		callback: func() error {return inspect(&arg)},
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
		// if len(os.Args) > 1 {
		// 	cmdArg := os.Args[1]
		// 	arg = strings.ToLower(cmdArg)
		// }
		if len(words) > 1 {
			arg = strings.ToLower(words[1])
		}

		cmd, exists := commands[first]
		if exists {
			if err := cmd.callback(); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}	
		// reset arg
		arg = ""
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
	var locations []string
	ch := make(chan pokecache.CacheEntryData)
	go func() {
		ch <- c.Get(links.Next)		
	}()
	cachedLocation := <- ch
	if cachedLocation.Exists {
		l, err := pokeapi.UnmarshalToList(links.Next, cachedLocation.Data)
		if err != nil {
			fmt.Println(err)
		}
		locations = l
		fmt.Println("Got page from cache!")
	} else { // else send GET request for data

		l, err := pokeapi.FetchPokeLocations(links.Next, c)	
		if err != nil {
			fmt.Println(err)
		}
		locations = l
	}
	printLocations(locations)
	// fmt.Printf("Current: %s\nNext: %s\nPrevious: %s", links.Current, links.Next, links.Previous)
	return nil
}

func printLocations(locations []string) {
	for _, l := range locations {
		fmt.Println(l)
	}
}

func commandMapb(links* pokeapi.PageLinks, c* pokecache.Cache) error {
	if links.Previous == "" {
		// fmt.Printf("Current: %s\nNext: %s\nPrevious: %s", links.Current, links.Next, links.Previous)
		return fmt.Errorf("you're on the first page")
	}

	var locations []string
	ch := make(chan pokecache.CacheEntryData)
	go func() {
		ch <- c.Get(links.Previous)		
	}()
	cachedLocation := <- ch
	if cachedLocation.Exists {
		l, err := pokeapi.UnmarshalToList(links.Previous, cachedLocation.Data)
		if err != nil {
			fmt.Println(err)
		}
		locations = l
		fmt.Println("Got page from cache!")
	} else { // else send GET request for data

		l, err := pokeapi.FetchPokeLocations(links.Previous, c)	
		if err != nil {
			fmt.Println(err)
		}
		locations = l
	}
	printLocations(locations)
	// fmt.Printf("Current: %s\nNext: %s\nPrevious: %s", links.Current, links.Next, links.Previous)
	return nil
}

func cleanInput(text string) []string {
	cleanedInput := strings.Fields(text)	
	return cleanedInput
}

func list(c* pokecache.Cache) error {
	for key := range c.Entries {
		fmt.Println(key)
	}
	return nil
}

func explore(c* pokecache.Cache, arg* string) error {
	fmt.Printf("Exploring %s...\n", *arg)
	fullUrl := pokeapi.BaseURL + *arg
	var pokemon []string
	ch := make(chan pokecache.CacheEntryData)
	go func() {
		ch <- c.Get(fullUrl)		
	}()
	cachedLocation := <- ch
	if cachedLocation.Exists {
		p, err := pokeapi.UnmarshalPokemonToList(fullUrl, cachedLocation.Data)
		if err != nil {
			fmt.Println(err)
		}
		pokemon = p
		fmt.Println("Got pokemon from cache!")
	} else {
		pokemonList, err := pokeapi.FetchLocationPokemon(fullUrl, c)
		if err != nil {
			fmt.Println(err)
		}
		pokemon = pokemonList
	}
	fmt.Println("Found Pokemon:")
	for _, p := range pokemon {
		fmt.Printf(" - %s\n", p)
	}
	return nil
}

func catch(arg* string) error {
	isCaught := false
	fmt.Printf("Throwing a Pokeball at %s...\n", *arg)
	fullUrl := pokeapi.PokeURL + *arg
	pokemon, err := pokeapi.FetchPokemon(fullUrl)
	if err != nil {
		return fmt.Errorf("Error finding %s: %v", *arg, err)
	}
	rng := rand.Float64()
	maxChance := 0.9
	scale := 1000.0
	catchChance := maxChance - (float64(pokemon.BaseExperience) / scale)

	if catchChance < 0.1 {
		catchChance = 0.1
	} else if catchChance > maxChance {
		catchChance = maxChance
	}

	// Random chance. goodluck.
	isCaught = rng < catchChance

	// if isCaught, add pokemon to pokedex
	if isCaught {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		Pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func inspect(arg* string) error {
	target, exists := Pokedex[*arg]
	if exists {
		fmt.Printf("Name: %s\nHeight: %v\nWeight: %v\n",
			target.Name, target.Height, target.Weight)
		fmt.Println("Stats:")
		for _, stat := range target.Stats {
			fmt.Printf(" -%s: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, t := range target.Types {
			fmt.Printf(" - %s\n",t.Type.Name)
		}
	} else {
		return fmt.Errorf("you have not caught that pokemon")
	}
	return nil
}
// End functions
