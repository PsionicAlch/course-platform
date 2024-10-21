package validators_test

import (
	"testing"

	"github.com/PsionicAlch/psionicalch-home/internal/validators"
)

func TestEmptyString(t *testing.T) {
	words := []struct {
		word  string
		empty bool
	}{
		{"", true},
		{"hello, world!", false},
	}

	for _, testWord := range words {
		if testWord.empty != validators.EmptyString(testWord.word) {
			var state string
			if !testWord.empty {
				state = "filled"
			} else {
				state = "empty"
			}

			t.Fatalf("%s was %s when it should not have been", testWord.word, state)
		}
	}
}
