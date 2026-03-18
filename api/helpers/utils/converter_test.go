package quizUtilsHelper

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	t.Run("returns string representation for supported values", func(t *testing.T) {
		id := uuid.New()

		testCases := []struct {
			name     string
			input    any
			expected string
		}{
			{
				name:     "string input",
				input:    "abcd",
				expected: "abcd",
			},
			{
				name:     "uuid input",
				input:    id,
				expected: id.String(),
			},
			{
				name:     "bool input",
				input:    true,
				expected: "true",
			},
			{
				name:     "integer input",
				input:    42,
				expected: "42",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				actual := GetString(tc.input)
				assert.Equal(t, tc.expected, actual)
				assert.IsType(t, "", actual)
			})
		}
	})
}

func TestGetBool(t *testing.T) {
	t.Run("returns bool value when input is bool", func(t *testing.T) {
		assert.True(t, GetBool(true))
		assert.False(t, GetBool(false))
	})

	t.Run("returns false for non bool input", func(t *testing.T) {
		assert.False(t, GetBool("true"))
		assert.False(t, GetBool(1))
		assert.False(t, GetBool(nil))
	})
}

func TestConvertType(t *testing.T) {
	t.Run("returns converted value when assertion succeeds", func(t *testing.T) {
		value, ok := ConvertType[string]("quiz")

		assert.True(t, ok)
		assert.Equal(t, "quiz", value)
	})

	t.Run("returns zero value and false when assertion fails", func(t *testing.T) {
		value, ok := ConvertType[string](123)

		assert.False(t, ok)
		assert.Equal(t, "", value)
	})

	t.Run("works with other target types", func(t *testing.T) {
		value, ok := ConvertType[bool](true)

		assert.True(t, ok)
		assert.True(t, value)
	})
}
