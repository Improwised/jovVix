package quizUtilsHelper

import (
	"fmt"
)

func GetString(text any) string {
	return fmt.Sprint(text)
}

func GetBool(flag any) bool {
	var value bool
	var ok bool
	if value, ok = flag.(bool); ok {
		return value
	}
	return false
}

func ConvertType[T any](x any) (T, bool) {

	v, ok := x.(T) // convert, then assert
	if !ok {
		return v, ok
	}
	return v, ok
}
