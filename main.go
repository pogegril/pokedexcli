// Project root
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/pogegril/pokedexcli/internal/cache"
	"github.com/pogegril/pokedexcli/internal/network"
	"github.com/pogegril/pokedexcli/internal/repl"
)

// User runtime config
type config struct {
	Commands                *map[string]cliCommand
	Next                    string
	Previous  		string
	Search                  []string
	Cache                   *pokecache.Cache
	Pokedex                 map[string]network.Pokemon
	Rng                     *rand.Rand
}

// CLI Commands struct type
type cliCommand struct {
	name                    string
	description             string
	callback func(*config)  error
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	pokedex := make(map[string]network.Pokemon) 

	// Commands' registry
	commands := map[string]cliCommand{
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
		"explore": {
			name:        "explore",
			description: "Displays pokemon that can be found in this area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch a pokemon to detail",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Displays pokémon's information if caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays caught pokémon",
			callback:    commandPokedex,
		},
	}

	// Config pointers
	cfg := &config{
		Next: "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
		Search: []string{},
		Cache: pokecache.NewCache(10 * time.Second),
		Pokedex: pokedex,
		Rng: rand.New(rand.NewSource(time.Now().UnixNano())),
		Commands: &commands,
	}


	// Program's loop
	for {
		// Prompt & Input handling
		fmt.Printf("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cfg.Search = repl.CleanInput(input)	

		// No search
		if (len(cfg.Search) == 0) {
			continue
		}

		command, exists := commands[cfg.Search[0]]

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
func unmarshal(data []byte, page any) error {
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

	for name, command := range *cfg.Commands {
		commandInfo := fmt.Sprintf("%s: %s", name, command.description)
		fmt.Println(commandInfo)
	}
	return nil
}

// Displays the next locations
func commandMap(cfg *config) error {
	var err error

	// Looks for requested page in cache if possible or requests it
	content, isCached := cfg.Cache.Get(cfg.Next)
	if !isCached {
		content, err = network.Search(cfg.Next)
		if err != nil {
			return fmt.Errorf("Failed to get next map locations: %w", err)
		}	
		cfg.Cache.Add(cfg.Next, content)
	}

	// Load result
	var page network.MapPage
	err = unmarshal(content, &page)
	if err != nil {
		return fmt.Errorf("Failed to read next map locations: %w", err)
	}	

	// Display map page
	for _, result := range page.Results {
		fmt.Println(result.Name)
	}

	// Update page urls
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
	var err error

	// Looks for requested page in cache if possible or requests it
	content, isCached := cfg.Cache.Get(cfg.Previous)
	if !isCached {
		content, err = network.Search(cfg.Previous)
		if err != nil {
			fmt.Println("You're on the first page")
			return err // Fails silently to not clutter if user is in the first page
		}	
		cfg.Cache.Add(cfg.Previous, content)
	}

	// Load result
	var page network.MapPage
	err = unmarshal(content, &page)
	if err != nil {
		return err
	}	

	// Display map page
	for _, result := range page.Results {
		fmt.Println(result.Name)
	}

	// Update page urls
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

// Displays pokemon found in the received location
func commandExplore(cfg *config) error {
	var err error
	url := "https://pokeapi.co/api/v2/location-area/" + cfg.Search[1] 

	// Looks for required page in cache if possible or requests it
	content, isCached := cfg.Cache.Get(url)
	if !isCached {
		content, err = network.Search(url)
		if err != nil {
			return fmt.Errorf("Failed to explore location details: %w", err)
		}	
		cfg.Cache.Add(url, content)
	}

	// Load result
	var page network.ExplorePage
	err = unmarshal(content, &page)
	if err != nil {
		return err
	}	

	// Display pokemon in the area
	for _, encounter := range page.PokemonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}
	return err
}

// Attempts to catch a Pokemon and add it to the Pokedex
func commandCatch(cfg *config) error {
	var err error
	url := "https://pokeapi.co/api/v2/pokemon/" + cfg.Search[1]

	content, isCached := cfg.Cache.Get(url)
	if !isCached {
		content, err = network.Search(url)
		if err != nil {
			return fmt.Errorf("Failed to find pokémon: %w", err)
		}	
		cfg.Cache.Add(url, content)
	}

	// Load result
	var pokemon network.Pokemon
	err = unmarshal(content, &pokemon)
	if err != nil {
		return err
	}	
	fmt.Println("Throwing a Pokeball at " + *pokemon.Name + "...")

	// Pseudo-random catch rate
	r := cfg.Rng.Int31n(101)
	catchRate := int32(75) - (int32(*pokemon.Experience) / int32(10))
	if r < catchRate {
		fmt.Println(*pokemon.Name + " was caught!")
		cfg.Pokedex[*pokemon.Name] = pokemon	
	} else {
		fmt.Println(*pokemon.Name + " escaped!")
	}
	return err
}

// Displays the pokemon's information if caught previously
func commandInspect(cfg *config) error {
	var err error
	url := "https://pokeapi.co/api/v2/pokemon/" + cfg.Search[1]

	content, isCached := cfg.Cache.Get(url)
	if !isCached {
		fmt.Println("You have not caught that Pokémon")
		return nil
	}

	// Load result
	var pokemon network.Pokemon
	err = unmarshal(content, &pokemon)
	if err != nil {
		return err
	}	

	// Display pokemon entry
	fmt.Println("Name: " + *pokemon.Name)

	height := fmt.Sprintf("Height: %d", *pokemon.Height)
	fmt.Println(height)

	weight := fmt.Sprintf("Weight: %d", *pokemon.Weight)
	fmt.Println(weight + "\n")

	fmt.Println("Stats:")
	for _, statEntry := range pokemon.Stats {
		statLine := fmt.Sprintf("  -%s: %d", *statEntry.Stat.Name, *statEntry.Value)
		fmt.Println(statLine)
	}

	fmt.Println("Types")
	for _, typeEntry := range pokemon.Types {
		fmt.Println("  -" + *typeEntry.Type.Name)	
	}
	return err
}

// Displays caught pokemon
func commandPokedex(cfg *config) error {
	fmt.Println("Your Pokedex:")

	for _, pokemon := range cfg.Pokedex {
		fmt.Println("  -" + *pokemon.Name)
	}
	return nil
}
