package utils

import (
	"fmt"
	"regexp"
	"strings"

	validator "gopkg.in/go-playground/validator.v9"
)

const VALIDATE_MESSAGE = "fields are invalid."

func ValidateEmail(email string) (bool, error) {
	return regexp.MatchString("[a-zA-z]+@improwised.com", email)
}

func ValidatorErrorString(err error) string {
	var msg string
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			msg += strings.ToLower(err.Field()) + ","
		}
		msg = strings.TrimSuffix(msg, ",")
		msg = fmt.Sprintf("%s %s", msg, VALIDATE_MESSAGE)
		return msg
	}
	return ""
}

func ValidateGlobalEmail(email string) (bool, error) {
	return regexp.MatchString(`^[\w.-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
}
