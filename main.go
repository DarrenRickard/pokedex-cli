package main

import "strings"
import "fmt"

func main() {
	fmt.Println(cleanInput("   hello world  !"))

}

func cleanInput(text string) []string {
	// cleanedInput := []string{}
	cleanedInput := strings.Fields(text)	
	return cleanedInput
}

