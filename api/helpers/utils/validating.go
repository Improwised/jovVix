package quizUtilsHelper

import (
	"strconv"

	"github.com/Improwised/quizz-app/api/constants"
)

func IsValidCode(code string) bool {
	if len(code) != 6 {
		return false
	}

	var num int
	var err error
	if num, err = strconv.Atoi(code); err != nil {
		return false
	}

	if num < constants.MinCode || num > constants.MaxCode {
		return false
	}

	return true
}
