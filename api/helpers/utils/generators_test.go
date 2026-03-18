package quizUtilsHelper

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRandomString(t *testing.T) {
	t.Run("correct length", func(t *testing.T) {
		str := GenerateRandomString(10)
		assert.Len(t, str, 10)
	})

	t.Run("different outputs", func(t *testing.T) {
		str1 := GenerateRandomString(10)
		str2 := GenerateRandomString(10)
		assert.NotEqual(t, str1, str2)
	})
}

func TestGenerateRandomInt(t *testing.T) {
	t.Run("within range", func(t *testing.T) {
		min := 10
		max := 20

		for i := 0; i < 100; i++ {
			val := GenerateRandomInt(min, max)
			assert.GreaterOrEqual(t, val, min)
			assert.LessOrEqual(t, val, max)
		}
	})
}

func TestGenerateNewStringHavingSuffixName(t *testing.T) {
	t.Run("adds suffix correctly", func(t *testing.T) {
		main := "hello"
		result := GenerateNewStringHavingSuffixName(main, 5, 20)

		assert.True(t, strings.HasPrefix(result, "hello"))
		assert.Contains(t, result, "_")
		assert.Len(t, result, len(main)+5)
	})

	t.Run("respects max length", func(t *testing.T) {
		main := "verylongstring"
		result := GenerateNewStringHavingSuffixName(main, 5, 10)

		assert.LessOrEqual(t, len(result), 10)
	})
}

func TestGenerateID(t *testing.T) {
	t.Run("returns increasing values", func(t *testing.T) {
		id1 := GenerateID()
		time.Sleep(1 * time.Nanosecond)
		id2 := GenerateID()

		assert.Greater(t, id2, id1)
	})
}