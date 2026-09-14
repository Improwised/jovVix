package models

import (
	"strconv"
	"testing"
)

func TestRandomOptionOrderContainsEachOriginalOptionOnce(t *testing.T) {
	options := map[string]string{
		"1": "first",
		"2": "second",
		"3": "third",
		"4": "fourth",
	}

	optionOrder, err := randomOptionOrder(options)
	if err != nil {
		t.Fatalf("randomOptionOrder returned an error: %v", err)
	}
	if len(optionOrder) != len(options) {
		t.Fatalf("got %d displayed options, want %d", len(optionOrder), len(options))
	}

	seenOriginalKeys := make(map[int]bool, len(options))
	for displayedKey, originalKey := range optionOrder {
		if _, err := strconv.Atoi(displayedKey); err != nil {
			t.Fatalf("displayed key %q is not numeric", displayedKey)
		}
		if seenOriginalKeys[originalKey] {
			t.Fatalf("original option key %d appears more than once", originalKey)
		}
		if _, ok := options[strconv.Itoa(originalKey)]; !ok {
			t.Fatalf("original option key %d does not exist", originalKey)
		}
		seenOriginalKeys[originalKey] = true
	}
}

func TestRandomOptionOrderRejectsNonNumericOptionKeys(t *testing.T) {
	_, err := randomOptionOrder(map[string]string{"a": "invalid"})
	if err == nil {
		t.Fatal("expected an error for a non-numeric option key")
	}
}
