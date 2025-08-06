package pokeapi

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io"
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
	Next: 			firstPage,
	Previous: 		"",
}

var firstPage = "https://pokeapi.co/api/v2/location-area/"

func FetchPokeLocations(url string) ([]string, error) {
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
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading resource: %v", err)
	}
	if err := json.Unmarshal(body, &pokemaps); err != nil {
		return nil, fmt.Errorf("Error Unmarshal'ing response body:\n%v", err)
	}
	for _, r := range pokemaps.Results {
		locations = append(locations, r.Name)	
	}

	Links.Next = pokemaps.Next
	Links.Previous = pokemaps.Previous

	return locations, nil
}
