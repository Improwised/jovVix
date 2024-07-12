package calculations

import (
	"database/sql"
	"math"

	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/structs"
	"github.com/doug-martin/goqu/v9"
	"go.uber.org/zap"
)

func CalculatePointsAndScore(userAnswer structs.ReqAnswerSubmit, db *goqu.Database, logger *zap.Logger) (sql.NullInt16, int, error) {

	var points sql.NullInt16 = sql.NullInt16{}
	var remainingTime int
	var remainingTimeFloat float64
	var timePoints int
	var basePoint int = 500
	var finalScore int = 0

	questionModel := models.InitQuestionModel(db, logger)
 
	answers, answerPoints, answerDurationInSeconds, options, err :=questionModel.GetAnswersPointsDurationOptions(userAnswer.QuestionId.String())
	if err != nil {
		return points, finalScore, err
	}

	// check type of the question
	actualAnswerLen := len(answers)
	userAnswerLen := len(userAnswer.AnswerKeys)

	// if not attempted
	if userAnswerLen == 0 {
		return points, finalScore, nil
	}

	points.Valid = true
	// for mcq type question
	if actualAnswerLen == 1 {
		if answers[0] == userAnswer.AnswerKeys[0] {
			points.Int16 = answerPoints
			remainingTime = (answerDurationInSeconds * 1000) - userAnswer.ResponseTime
			remainingTimeFloat = math.Round(float64(remainingTime) / 1000)
			timePoints = int(math.Round((remainingTimeFloat * 400) / float64(answerDurationInSeconds)))
			finalScore = timePoints + basePoint + int(points.Int16*100)
			return points, finalScore, nil
		}
		return points, finalScore, nil
	} else if actualAnswerLen != len(options) {
		// if there are more than 1 correct answers
		for i := 0; i < actualAnswerLen; i++ {
			if answers[i] != userAnswer.AnswerKeys[0] {
				continue
			} else {
				// if answer selected by the user matches with any one of the correct answer (Partial evaluation)
				points.Int16 = answerPoints
				remainingTime = (answerDurationInSeconds * 1000) - userAnswer.ResponseTime
				remainingTimeFloat = math.Round(float64(remainingTime) / 1000)
				timePoints = int(math.Round((remainingTimeFloat * 400) / float64(answerDurationInSeconds)))
				finalScore = timePoints + basePoint + int(points.Int16*100)
				return points, finalScore, nil
			}
		}
	}

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
	return points, finalScore, nil
}