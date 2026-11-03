package models

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestSuggestPlayerNames(t *testing.T) {
	available := func(name string) (bool, error) { return !strings.EqualFold(name, "taken"), nil }
	suggestions, err := suggestPlayerNames("Smit", "Patel", "Smit", available)
	if err != nil {
		t.Fatal(err)
	}
	if len(suggestions) < 2 || suggestions[0] != "Patel" || suggestions[1] != "SmitPatel" {
		t.Fatalf("unexpected suggestions: %v", suggestions)
	}
	suggestions, err = suggestPlayerNames("Longfirstname", "Taken", "Longnamehere", available)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range suggestions {
		if strings.EqualFold(name, "taken") || utf8.RuneCountInString(name) > 12 {
			t.Fatalf("invalid suggestion %q", name)
		}
	}
	suggestions, err = suggestPlayerNames("Smit", "", "Smit", available)
	if err != nil || len(suggestions) == 0 {
		t.Fatalf("guest suggestions: %v, %v", suggestions, err)
	}
}
