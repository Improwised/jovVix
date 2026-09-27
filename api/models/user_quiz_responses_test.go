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

func TestTranslateAnswerKeysToPlayer(t *testing.T) {
	optionOrder := map[string]int{"1": 3, "2": 1, "3": 2}

	got, err := translateAnswerKeysToPlayer(optionOrder, []int{2})
	if err != nil {
		t.Fatalf("translateAnswerKeysToPlayer returned an error: %v", err)
	}
	if len(got) != 1 || got[0] != 3 {
		t.Fatalf("got displayed answer keys %v, want [3]", got)
	}
}

func TestTranslateAnswerKeysToPlayerRejectsMissingOriginalKey(t *testing.T) {
	_, err := translateAnswerKeysToPlayer(map[string]int{"1": 3, "2": 1}, []int{2})
	if err == nil {
		t.Fatal("expected an error when the answer key is missing from the option order")
	}
}
