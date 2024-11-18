package utils

import (
	"database/sql"
	"math"

	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/pkg/structs"
)

func CalculatePointsAndScore(userAnswer structs.ReqAnswerSubmit, answers []int, answerPoints int16, answerDurationInSeconds, questionType int) (sql.NullInt16, int) {

	var points sql.NullInt16 = sql.NullInt16{}
	var remainingTime int
	var remainingTimeFloat float64
	var timePoints int
	var basePoint int = 500
	var finalScore int = 0

	// check type of the question
	actualAnswerLen := len(answers)
	userAnswerLen := len(userAnswer.AnswerKeys)

	// if not attempted
	if userAnswerLen == 0 {
		return points, finalScore
	}

	points.Valid = true
	// for mcq type question
	if actualAnswerLen == 1 && answerPoints > 0 {
		if answers[0] == userAnswer.AnswerKeys[0] {
			points.Int16 = answerPoints
			remainingTime = (answerDurationInSeconds * 1000) - userAnswer.ResponseTime
			remainingTimeFloat = math.Round(float64(remainingTime) / 1000)
			timePoints = int(math.Round((remainingTimeFloat * 400) / float64(answerDurationInSeconds)))
			finalScore = timePoints + basePoint + int(points.Int16*100)
			return points, finalScore
		}
		return points, finalScore
	} else if questionType == constants.Survey && answerPoints > 0 {
		points.Int16 = answerPoints
		remainingTime = (answerDurationInSeconds * 1000) - userAnswer.ResponseTime
		remainingTimeFloat = math.Round(float64(remainingTime) / 1000)
		timePoints = int(math.Round((remainingTimeFloat * 400) / float64(answerDurationInSeconds)))
		finalScore = timePoints + basePoint + int(points.Int16*100)
		return points, finalScore
	}

	// logic for multiple correct answer to be used in future functionalities
	// fetch options from db for that (make changes in questions.go )
	// else if actualAnswerLen != len(options) {
	// 	// if there are more than 1 correct answers
	// 	for i := 0; i < actualAnswerLen; i++ {
	// 		if answers[i] != userAnswer.AnswerKeys[0] {
	// 			continue
	// 		} else {
	// 			// if answer selected by the user matches with any one of the correct answer (Partial evaluation)
	// 			points.Int16 = answerPoints
	// 			remainingTime = (answerDurationInSeconds * 1000) - userAnswer.ResponseTime
	// 			remainingTimeFloat = math.Round(float64(remainingTime) / 1000)
	// 			timePoints = int(math.Round((remainingTimeFloat * 400) / float64(answerDurationInSeconds)))
	// 			finalScore = timePoints + basePoint + int(points.Int16*100)
	// 			return points, finalScore, nil
	// 		}
	// 	}
	// }

	var noOfMatches int = 0
	for _, actualAnswer := range answers {
		for _, userAnswer := range userAnswer.AnswerKeys {
			if actualAnswer == userAnswer {
				noOfMatches += 1
				if noOfMatches == userAnswerLen {
					break
				}
			}
		}
	}

	points.Int16 = int16(noOfMatches) * answerPoints
	return points, finalScore
}

func CalculateStreakScore(streakCount, score int) (int, int) {

	// reset streak if score is 0
	if score <= 0 {
		return score, 0
	}

	// dont add bonus if streak is 0
	if streakCount <= 0 {
		return score, streakCount + 1
	}

	// Add streak count score in final score
	streakScore := (constants.StreakBaseScore * streakCount * 10) / 100
	finalScore := score + constants.StreakBaseScore + streakScore

	return finalScore, streakCount + 1
}
