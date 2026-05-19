// Handles internet requests
package network

import(
	"fmt"
	"io"
	"net/http"
)

// Searches an url and returns the byte slice received
func Search(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch %s: %w", url, err)	
	}

	content, err := io.ReadAll(response.Body)
	response.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("Failed to read content: %w", err)
		
	}
	
	if response.StatusCode > 299 {
		return content, fmt.Errorf("Invalid response: %d\n%w", response.StatusCode, err)
	}

	return content, err
}

