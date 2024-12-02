package utils

import (
	"math"
	"testing"

	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/pkg/structs"
	"github.com/stretchr/testify/assert"
)

func TestCalculatePointsAndScore(t *testing.T) {

	t.Run("No answer attempted", func(t *testing.T) {
		userAnswer := structs.ReqAnswerSubmit{AnswerKeys: []int{}}
		answers := []int{1}
		answerPoints := int16(10)
		answerDurationInSeconds := 30
		questionType := constants.SingleAnswer

		points, score := CalculatePointsAndScore(userAnswer, answers, answerPoints, answerDurationInSeconds, questionType)
		assert.False(t, points.Valid)
		assert.Equal(t, 0, score)
	})

	t.Run("Correct answer", func(t *testing.T) {
		userAnswer := structs.ReqAnswerSubmit{AnswerKeys: []int{1}, ResponseTime: 5000}
		answers := []int{1}
		answerPoints := int16(10)
		answerDurationInSeconds := 30
		questionType := constants.SingleAnswer

		// Calculate expected values based on your function logic
		remainingTime := (answerDurationInSeconds * 1000) - userAnswer.ResponseTime
		remainingTimeFloat := math.Round(float64(remainingTime) / 1000)
		timePoints := int(math.Round((remainingTimeFloat * 400) / float64(answerDurationInSeconds)))
		basePoint := 500
		expectedScore := timePoints + basePoint + int(answerPoints*100)

		points, score := CalculatePointsAndScore(userAnswer, answers, answerPoints, answerDurationInSeconds, questionType)
		assert.True(t, points.Valid)
		assert.Equal(t, answerPoints, points.Int16)
		assert.Equal(t, expectedScore, score)
	})

	t.Run("Incorrect answer", func(t *testing.T) {
		userAnswer := structs.ReqAnswerSubmit{AnswerKeys: []int{2}, ResponseTime: 5000}
		answers := []int{1}
		answerPoints := int16(10)
		answerDurationInSeconds := 30
		questionType := constants.SingleAnswer

		points, score := CalculatePointsAndScore(userAnswer, answers, answerPoints, answerDurationInSeconds, questionType)
		assert.True(t, points.Valid)
		assert.Equal(t, 0, score)
	})

	t.Run("Survey question", func(t *testing.T) {
		userAnswer := structs.ReqAnswerSubmit{AnswerKeys: []int{1}, ResponseTime: 5000}
		answers := []int{1}
		answerPoints := int16(10)
		answerDurationInSeconds := 30
		questionType := constants.Survey

		remainingTime := (answerDurationInSeconds * 1000) - userAnswer.ResponseTime
		remainingTimeFloat := math.Round(float64(remainingTime) / 1000)
		timePoints := int(math.Round((remainingTimeFloat * 400) / float64(answerDurationInSeconds)))
		basePoint := 500
		expectedScore := timePoints + basePoint + int(answerPoints*100)

		points, score := CalculatePointsAndScore(userAnswer, answers, answerPoints, answerDurationInSeconds, questionType)
		assert.True(t, points.Valid)
		assert.Equal(t, answerPoints, points.Int16)
		assert.Equal(t, expectedScore, score)
	})
}

func TestCalculateStreakScore(t *testing.T) {

	// Test Case 1: Zero score, streak should reset
	t.Run("Zero score should reset streak", func(t *testing.T) {
		finalScore, newStreak := CalculateStreakScore(5, 0)
		assert.Equal(t, 0, finalScore)
		assert.Equal(t, 0, newStreak)
	})

	// Test Case 2: Zero streak, only increment streak without bonus
	t.Run("Zero streak should not add bonus", func(t *testing.T) {
		finalScore, newStreak := CalculateStreakScore(0, 50)
		assert.Equal(t, 50, finalScore)
		assert.Equal(t, 1, newStreak)
	})

	// Test Case 3: Positive score and positive streak
	t.Run("Positive streak should add bonus", func(t *testing.T) {
		finalScore, newStreak := CalculateStreakScore(3, 200)
		expectedStreakScore := (constants.StreakBaseScore * 3 * 10) / 100
		expectedFinalScore := 200 + constants.StreakBaseScore + expectedStreakScore
		assert.Equal(t, expectedFinalScore, finalScore)
		assert.Equal(t, 4, newStreak)
	})

	// Test Case 4: Negative score, streak should reset
	t.Run("Negative score should reset streak", func(t *testing.T) {
		finalScore, newStreak := CalculateStreakScore(3, -10)
		assert.Equal(t, -10, finalScore)
		assert.Equal(t, 0, newStreak)
	})

	// Test Case 5: Score is positive and streak is zero
	t.Run("Score is positive and streak is zero", func(t *testing.T) {
		finalScore, newStreak := CalculateStreakScore(0, 100)
		assert.Equal(t, 100, finalScore)
		assert.Equal(t, 1, newStreak)
	})

	// Test Case 6: Large streak count
	t.Run("Large streak count", func(t *testing.T) {
		finalScore, newStreak := CalculateStreakScore(10, 500)
		expectedStreakScore := (constants.StreakBaseScore * 10 * 10) / 100
		expectedFinalScore := 500 + constants.StreakBaseScore + expectedStreakScore
		assert.Equal(t, expectedFinalScore, finalScore)
		assert.Equal(t, 11, newStreak)
	})
}
