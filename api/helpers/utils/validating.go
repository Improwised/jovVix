package quizUtilsHelper

import "strconv"

func IsValidCode(code string) bool {
	if len(code) != 6 {
		return false
	}

	var num int
	var err error
	if num, err = strconv.Atoi(code); err != nil {
		return false
	}

	if num < 100000 || num > 999999 {
		return false
	}

	return true
}
