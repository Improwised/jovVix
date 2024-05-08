package quizUtilsHelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatingCode(t *testing.T) {
	t.Run("Test validate code", func(t *testing.T) {
		expected := true
		actual := IsValidCode("123456")
		assert.Equal(t, expected, actual)

		expected = false
		actual = IsValidCode("9999999")
		assert.Equal(t, expected, actual)

	})
}
