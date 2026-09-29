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

func TestRearrangeOptionsRoundTrip(t *testing.T) {
	options := map[string]string{
		"1": "A",
		"2": "B",
		"3": "C",
		"4": "D",
	}

	optionOrder := map[string]int{
		"1": 3,
		"2": 1,
		"3": 4,
		"4": 2,
	}

	correctKeys := []int{2}
	selectedKeys := []int{4}

	shuffledOptions, translatedCorrect, translatedSelected, err := rearrangeOptions(optionOrder, options, correctKeys, selectedKeys)
	if err != nil {
		t.Fatalf("rearrangeOptions returned an error: %v", err)
	}

	if shuffledOptions["1"] != "C" {
		t.Errorf("displayed key 1 should map to option C, got %s", shuffledOptions["1"])
	}
	if shuffledOptions["2"] != "A" {
		t.Errorf("displayed key 2 should map to option A, got %s", shuffledOptions["2"])
	}
	if shuffledOptions["3"] != "D" {
		t.Errorf("displayed key 3 should map to option D, got %s", shuffledOptions["3"])
	}
	if shuffledOptions["4"] != "B" {
		t.Errorf("displayed key 4 should map to option B, got %s", shuffledOptions["4"])
	}

	if len(translatedCorrect) != 1 || translatedCorrect[0] != 4 {
		t.Errorf("correct key [2] should map to displayed position [4], got %v", translatedCorrect)
	}
	if len(translatedSelected) != 1 || translatedSelected[0] != 3 {
		t.Errorf("selected key [4] should map to displayed position [3], got %v", translatedSelected)
	}
}

func TestRearrangeOptionsIdentityMapping(t *testing.T) {
	options := map[string]string{
		"1": "A",
		"2": "B",
	}

	optionOrder := map[string]int{
		"1": 1,
		"2": 2,
	}

	correctKeys := []int{1}
	selectedKeys := []int{2}

	shuffledOptions, translatedCorrect, translatedSelected, err := rearrangeOptions(optionOrder, options, correctKeys, selectedKeys)
	if err != nil {
		t.Fatalf("rearrangeOptions returned an error: %v", err)
	}

	if shuffledOptions["1"] != "A" {
		t.Errorf("identity mapping should preserve option A at key 1")
	}
	if shuffledOptions["2"] != "B" {
		t.Errorf("identity mapping should preserve option B at key 2")
	}

	if len(translatedCorrect) != 1 || translatedCorrect[0] != 1 {
		t.Errorf("identity mapping should not change correct answer")
	}
	if len(translatedSelected) != 1 || translatedSelected[0] != 2 {
		t.Errorf("identity mapping should not change selected answer")
	}
}

func TestRearrangeOptionsMissingOption(t *testing.T) {
	options := map[string]string{
		"1": "A",
	}

	optionOrder := map[string]int{
		"1": 99,
	}

	_, _, _, err := rearrangeOptions(optionOrder, options, []int{}, []int{})
	if err == nil {
		t.Fatal("expected an error for missing option referenced by option order")
	}
}

func TestRearrangeOptionsInvalidKey(t *testing.T) {
	options := map[string]string{
		"1": "A",
		"2": "B",
	}

	optionOrder := map[string]int{
		"1": 1,
		"2": 2,
	}

	correctKeys := []int{99}
	_, _, _, err := rearrangeOptions(optionOrder, options, correctKeys, []int{})
	if err == nil {
		t.Fatal("expected an error for invalid original key not found in inverse map")
	}
}