// Provides REPL tools to handle user input
package repl

import "strings"

// Clears user input normalizing it into lowercase words without whitespaces
func CleanInput(text string) []string {
	trimmed := strings.TrimSpace(text)
	lowerCase := strings.ToLower(trimmed)	
	return strings.Fields(lowerCase)
}
