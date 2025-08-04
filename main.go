package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		var words []string
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		words = cleanInput(userInput)
		fmt.Printf("Your command was: %s\n", words[0])
	} 



}

func cleanInput(text string) []string {
	// cleanedInput := []string{}
	cleanedInput := strings.Fields(text)	
	return cleanedInput
}

