package pokeapi

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io"
	"github.com/darrenrickard/pokedexcli/internal/pokecache"
)

var BaseURL = "https://pokeapi.co/api/v2/location-area/"
var PokeURL = "https://pokeapi.co/api/v2/pokemon/"

// need to refactor LocationArea and LocationAreaDetailed into an enum Location

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

type PageLinks struct {
	Current     string
	Next 		string
	Previous 	string
}

var Links = PageLinks {
	Current: 		"",
	Next: 			FirstPage,
	Previous: 		"",
}

type LocationAreaDetailed struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

const FirstPage = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"

func FetchPokeLocations(url string, c *pokecache.Cache) ([]string, error) {
	// url is blank if user is on first page
	if url == "" {
		return nil, fmt.Errorf("you're on the first page")
	}

	// keeps track of current link
	Links.Current = url

	var locations []string
	var pokemaps LocationArea
	req, err := http.NewRequest("GET", url , nil)	
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error receiving response: %v", err)
	}
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, res.Body)	
	}
	body, err := io.ReadAll(res.Body) // body is type of []byte, separate this to own function to use globally
	if err != nil {
		return nil, fmt.Errorf("Error reading resource: %v", err)
	}
	if err := json.Unmarshal(body, &pokemaps); err != nil {
		return nil, fmt.Errorf("Error Unmarshal'ing response body:\n%v", err)
	}
	for _, r := range pokemaps.Results {
		locations = append(locations, r.Name)	
	}
	// Add current link and data to cache map
	go c.Add(Links.Current, body)

	Links.Next = pokemaps.Next
	Links.Previous = pokemaps.Previous

	return locations, nil
}


func FetchLocationPokemon(location string, c *pokecache.Cache) ([]string, error) {
	if location == "" {
		return nil, fmt.Errorf("missing location")
	}


	var pokemonList []string
	var locationDetailed LocationAreaDetailed
	req, err := http.NewRequest("GET", location, nil)	
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error receiving response: %v", err)
	}
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, res.Body)	
	}
	body, err := io.ReadAll(res.Body) // body is type of []byte, separate this to own function to use globally
	if err != nil {
		return nil, fmt.Errorf("Error reading resource: %v", err)
	}
	if err := json.Unmarshal(body, &locationDetailed); err != nil {
		return nil, fmt.Errorf("Error Unmarshal'ing response body:\n%v", err)
	}
	for _, val := range locationDetailed.PokemonEncounters{
		pokemonList = append(pokemonList, val.Pokemon.Name)	
	}
	// Add current link and data to cache map
	go c.Add(location, body)

	return pokemonList, nil
}

func FetchPokemon(url string) (Pokemon, error) {
	var pokemon Pokemon 
	if url == "" {
		return pokemon, fmt.Errorf("missing location")
	}
	req, err := http.NewRequest("GET", url, nil)	
	if err != nil {
		return pokemon, fmt.Errorf("Error creating request: %v", err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return pokemon, fmt.Errorf("Error receiving response: %v", err)
	}
	if res.StatusCode > 299 {
		return pokemon, fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, res.Body)	
	}
	body, err := io.ReadAll(res.Body) // body is type of []byte, separate this to own function to use globally
	if err != nil {
		return pokemon, fmt.Errorf("Error reading resource: %v", err)
	}
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return pokemon, fmt.Errorf("Error Unmarshal'ing response body:\n%v", err)
	}
	return pokemon, nil
}

func UnmarshalToList(url string, b []byte) ([]string, error) {
	// move Links update outside of UnmarshalToList()
	// change pokemaps LocationArea to location enum/any
	// move range pokemaps.Results outside of UnmarshalToList()
	Links.Current = url
	var locations []string
	var pokemaps LocationArea
	if err := json.Unmarshal(b, &pokemaps); err != nil {
		return nil, fmt.Errorf("Error Unmarshal'ing cache item data:\n%v", err)
	}
	for _, r := range pokemaps.Results {
		locations = append(locations, r.Name)	
	}
	Links.Next = pokemaps.Next
	Links.Previous = pokemaps.Previous

	return locations, nil
}

func UnmarshalPokemonToList(url string, b []byte) ([]string, error) {
	var pokemon []string
	var location LocationAreaDetailed
	if err := json.Unmarshal(b, &location); err != nil {
		return nil, fmt.Errorf("Error Unmarshal'ing cache item data:\n%v", err)
	}
	for _, r := range location.PokemonEncounters{
		pokemon = append(pokemon, r.Pokemon.Name)	
	}
	return pokemon, nil
}
