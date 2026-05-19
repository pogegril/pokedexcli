// Provides REPL tools to handle user input
package repl

import "testing"

// Runs tests on [ cleanInput(text string) string ]
func TestCleanInput(t *testing.T) {
	cases := []struct {
		input 	 string
		expected []string
	}{
		{
			input: "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input: "  HeLlO  WoRlD  ",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		actual := CleanInput(c.input)
		expected := c.expected

		if len(actual) != len(expected) {
			t.Errorf("Unexpected slice length\n  Actual: %d\n  Expected: %d\n", len(actual), len(expected))
		}
		
		for i := range actual {
			if (actual[i] != expected[i]) {
				t.Errorf("String segments don't match.\n  Actual: %s\n  Expected:  %s\n", actual[i], expected[i])
			}
		}
	}
}
