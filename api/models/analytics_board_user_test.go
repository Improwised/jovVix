package models

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParseAnswerKeys(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		expected []int
	}{
		{"valid keys", "[1, 3]", []int{1, 3}},
		{"single key", "[2]", []int{2}},
		{"empty array", "[]", []int{}},
		{"keys with spaces", "[ 0, 2, 5 ]", []int{0, 2, 5}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseAnswerKeys(tt.raw); !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("parseAnswerKeys(%q) = %v, want %v", tt.raw, got, tt.expected)
			}
		})
	}
}

func TestParseAnswerKeysReturnsEmptyOnInvalidInput(t *testing.T) {
	for _, raw := range []string{"", "not-json", `"1"`, "[1, x]", "null", "{}"} {
		t.Run(fmt.Sprintf("input %q", raw), func(t *testing.T) {
			if got := parseAnswerKeys(raw); len(got) != 0 {
				t.Fatalf("parseAnswerKeys(%q) = %v, want empty", raw, got)
			}
		})
	}
}

func TestMarshalAnswerKeys(t *testing.T) {
	tests := []struct {
		name     string
		keys     []int
		expected string
	}{
		{"multiple keys", []int{1, 3}, "[1,3]"},
		{"single key", []int{2}, "[2]"},
		{"empty slice", []int{}, "[]"},
		// encoding/json marshals a nil slice as null; pin the current behavior.
		{"nil slice", nil, "null"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := marshalAnswerKeys(tt.keys); got != tt.expected {
				t.Fatalf("marshalAnswerKeys(%v) = %q, want %q", tt.keys, got, tt.expected)
			}
		})
	}
}

func TestAnswerKeysRoundTrip(t *testing.T) {
	keys := []int{0, 2, 5}
	if got := parseAnswerKeys(marshalAnswerKeys(keys)); !reflect.DeepEqual(got, keys) {
		t.Fatalf("round trip of %v = %v", keys, got)
	}
}
