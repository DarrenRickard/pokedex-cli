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
	Current: 		"https://pokeapi.co/api/v2/location-area/",
	Next: 			"",
	Previous: 		"",
}

func FetchLocationPageLinks(url string) error {
	var pokemaps LocationArea
	fmt.Println("Creating Request...")
	req, err := http.NewRequest("GET", url , nil)	
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
	// for _, r := range pokemaps.Results {
	// 	locations = append(locations, r.Name)
	// }
	Links.Next = pokemaps.Next
	Links.Previous = pokemaps.Previous
	return nil
}

func FetchPokeLocations(url string) ([]string, error) {
	var locations []string
	var pokemaps LocationArea
	fmt.Println("Creating Request...")
	req, err := http.NewRequest("GET", url , nil)	
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("limit", "20")
	req.Header.Set("offset", "20")	
	fmt.Println("Sending Request...")
	client := &http.Client{}
	res, err := client.Do(req)
	fmt.Println("Request sent, Response received")
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
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
	return locations, nil
}
