package quizUtilsHelper

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetString(t *testing.T) {
	t.Run("Test get string", func(t *testing.T) {
		expected := "abcd"
		actual := GetString("abcd")
		assert.Equal(t, expected, actual)

		id := uuid.New()
		actual = GetString(id)
		assert.IsType(t, expected, actual)

		expected = "true"
		actual = GetString(true)
		assert.Equal(t, expected, actual)
	})
}
