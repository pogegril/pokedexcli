// Handles internet requests
package network

import(
	"io"
	"net/http"
)

// MapPage results
type MapPage struct {
	Next                  *string `json:"next"`
	Previous      	      *string `json:"previous"`
	Results 	      []struct {
		Name  string `json:"name"` 
	} `json:"results"`
}

// Explore Area page results
type ExplorePage struct {
	PokemonEncounters     []struct{
		Pokemon struct {
			Name  *string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

// Pokemon Entry
type Pokemon struct {
	Name                  *string `json:"name"`
	Experience            *int    `json:"base_experience"`
	Height                *int    `json:"height"`
	Weight                *int    `json:"weight"`

	Stats         []struct {
		Stat struct {
			Name   *string `json:"name"`
		}
		Value  	       *int `json:"base_stat"`
	}                      `json:"stats"`

	Types         []struct {
		Type struct {
			Name   *string `json:"name"`
		}              `json:"type"`
	} 		       `json:"types"`
}


// Searches an url and returns the byte slice received
func Search(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	content, err := io.ReadAll(response.Body)
	response.Body.Close()

	if err != nil {
		return nil, err
		
	}
	
	if response.StatusCode > 299 {
		return content, err
	}
	return content, err
}

