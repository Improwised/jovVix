package quizUtilsHelper

import "strconv"

func IsValidCode(code string) bool {
    if len(code) != 6 {
        return false
    }

    if _, err := strconv.Atoi(code); err != nil {
        return false
    }

    // Check if code is between 100000 and 999999
    num, _ := strconv.Atoi(code)
    if num < 100000 || num > 999999 {
        return false
    }

    return true
}