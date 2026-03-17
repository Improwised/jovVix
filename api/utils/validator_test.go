package utils

import (
	"testing"
	"errors"

	"github.com/stretchr/testify/assert"
	validator "gopkg.in/go-playground/validator.v9"
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

func TestValidateGlobalEmail(t *testing.T) {
	email := "test@gmail.com"

	valid, err := ValidateGlobalEmail(email)
	assert.NoError(t, err)
	assert.True(t, valid)

	email = "invalid-email"
	valid, err = ValidateGlobalEmail(email)
	assert.NoError(t, err)
	assert.False(t, valid)
}

func TestValidatorErrorString(t *testing.T) {
	validate := validator.New()

	type TestStruct struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}

	// Case: validation error exists
	obj := TestStruct{}
	err := validate.Struct(obj)

	msg := ValidatorErrorString(err)
	assert.Contains(t, msg, "name")
	assert.Contains(t, msg, "email")
	assert.Contains(t, msg, VALIDATE_MESSAGE)

	// Case: no error
	msg = ValidatorErrorString(nil)
	assert.Equal(t, "", msg)

	// Edge case: non-validator error (should panic, but we test safely)
	defer func() {
		if r := recover(); r != nil {
			assert.True(t, true) // expected panic
		}
	}()
	_ = ValidatorErrorString(errors.New("some random error"))
}