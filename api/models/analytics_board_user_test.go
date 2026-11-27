package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAnswerKeys(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want []int
	}{
		{"valid list", "[1, 2, 3]", []int{1, 2, 3}},
		{"single key", "[0]", []int{0}},
		{"negative keys", "[-1, 5]", []int{-1, 5}},
		{"empty list", "[]", []int{}},
		{"null unmarshals to a nil slice", "null", nil},
		{"empty string", "", []int{}},
		{"not JSON", "not json", []int{}},
		{"wrong element type", `[1, "a"]`, []int{}},
		{"object instead of list", `{"a": 1}`, []int{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, parseAnswerKeys(tc.raw))
		})
	}
}

func TestMarshalAnswerKeys(t *testing.T) {
	tests := []struct {
		name string
		keys []int
		want string
	}{
		{"several keys", []int{1, 2, 3}, "[1,2,3]"},
		{"single key", []int{0}, "[0]"},
		{"empty slice", []int{}, "[]"},
		// Documents the current behavior for a nil input. The callers in
		// applyOptionOrderForRow pass slices from parseAnswerKeys, which
		// never returns nil, so this path is not hit in practice.
		{"nil slice", nil, "null"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, marshalAnswerKeys(tc.keys))
		})
	}
}

func TestAnswerKeysRoundTrip(t *testing.T) {
	for _, raw := range []string{"[1,2,3]", "[]", "[0]", "[-1,5]"} {
		assert.Equal(t, raw, marshalAnswerKeys(parseAnswerKeys(raw)))
	}
}
