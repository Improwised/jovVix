package utils

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/Improwised/jovvix/api/constants"
	validator "gopkg.in/go-playground/validator.v9"
)

const VALIDATE_MESSAGE = "fields are invalid."

// ValidateQuizTitle returns a user-facing validation error for a quiz title.
// Rune counting matches the character limit users see, including non-ASCII titles.
func ValidateQuizTitle(title string) string {
	if strings.TrimSpace(title) == "" {
		return constants.QuizTitleRequired
	}

	if utf8.RuneCountInString(title) > constants.QuizTitleMaxLength {
		return constants.ErrQuizTitleTooLong
	}

	return ""
}

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
