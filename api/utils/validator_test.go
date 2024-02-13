package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmail(t *testing.T) {
	email := "xyz@improwised.com"

	valid, err := ValidateEmail(email)
	assert.NoError(t, err)
	assert.True(t, valid)

	email = "xyz@improwisd.com"
	valid, err = ValidateEmail(email)
	assert.NoError(t, err)
	assert.False(t, valid)
}
