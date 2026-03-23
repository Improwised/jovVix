package quizUtilsHelper

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/Improwised/jovvix/api/constants"
)

func TestIsValidCode(t *testing.T) {

	t.Run("valid code", func(t *testing.T) {
		validCode := strconv.Itoa(constants.MinInvitationCode)
		assert.True(t, IsValidCode(validCode))
	})

	t.Run("invalid length", func(t *testing.T) {
		assert.False(t, IsValidCode("12345"))
		assert.False(t, IsValidCode("1234567"))
	})

	t.Run("non numeric code", func(t *testing.T) {
		assert.False(t, IsValidCode("12a456"))
		assert.False(t, IsValidCode("abcdef"))
	})

	t.Run("out of range", func(t *testing.T) {
		// below min
		belowMin := strconv.Itoa(constants.MinInvitationCode - 1)
		assert.False(t, IsValidCode(belowMin))

		// above max
		aboveMax := strconv.Itoa(constants.MaxInvitationCode + 1)
		assert.False(t, IsValidCode(aboveMax))
	})
}
