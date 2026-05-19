// Handles internet requests
package network

import(
	"io"
	"net/http"
)

// Page results
type MapPage struct {
	Next      *string
	Previous  *string
	Results []struct {
		Name  string 
		URL   string
	}
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

