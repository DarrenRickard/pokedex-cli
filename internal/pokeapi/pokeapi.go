package pokeapi

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io"
	"github.com/darrenrickard/pokedexcli/internal/pokecache"
)

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

const FirstPage = "https://pokeapi.co/api/v2/location-area/"

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

func UnmarshalToList(url string, b []byte) ([]string, error) {
	Links.Current = url
	var locations []string
	var pokemaps LocationArea
	// !!! bug is occuring here after reading from cache once
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
