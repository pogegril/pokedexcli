// Project root
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/pogegril/pokedexcli/network"
	"github.com/pogegril/pokedexcli/repl"
)

// CLI Commands struct type
type cliCommand struct {
	name                    string
	description             string
	callback func(*config)  error
}
var commands map[string]cliCommand

// User runtime config
type config struct {
	Next                    string
	Previous  		string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Pageholder
	cfg := &config{
		Next: "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
	}

	// Commands' registry
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,	
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    commandMapBack,
		},
	}

	// Program's loop
	for {
		// Prompt & Input handling
		fmt.Printf("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		search := repl.CleanInput(input)	

		// No search
		if (len(search) == 0) {
			continue
		}

		command, exists := commands[search[0]]

		// Command not found
		if !exists {
			fmt.Printf("Unknown command\n")
			continue
		}

		// Execute command
		command.callback(cfg)
	} 
}

// Loads the search results from a byte slice
func unmarshal(data []byte, page *network.MapPage) error {
	return json.Unmarshal(data, page)
}

// Closes the program
func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

// Displays a help message to describe the commands' usage
func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:\n ")

	for name, command := range commands {
		commandInfo := fmt.Sprintf("%s: %s", name, command.description)
		fmt.Println(commandInfo)
	}
	return nil
}

// Displays the next locations
func commandMap(cfg *config) error {
	content, err := network.Search(cfg.Next)
	if err != nil {
		return fmt.Errorf("Failed to get next map locations: %w", err)
	}	

	var page network.MapPage
	err = unmarshal(content, &page)
	if err != nil {
		return fmt.Errorf("Failed to read next map locations: %w", err)
	}	

	for _, result := range page.Results {
		fmt.Println(result.Name)
	}

	if page.Next != nil {
		cfg.Next = *page.Next
	}
	if page.Previous != nil {
		cfg.Previous = *page.Previous
	}
	return err
}

// Displays the previous locations
func commandMapBack(cfg *config) error {
	content, err := network.Search(cfg.Previous)
	if err != nil {
		fmt.Println("You're on the first page")
		return err // Fails silently to not clutter if user is in the first page
	}	

	var page network.MapPage
	err = unmarshal(content, &page)
	if err != nil {
		return err
	}	

	for _, result := range page.Results {
		fmt.Println(result.Name)
	}

	if page.Next != nil {
		cfg.Next = *page.Next
	}
	if page.Previous != nil {
		cfg.Previous = *page.Previous
	} else {
		cfg.Previous = ""
	}
	return err
}
