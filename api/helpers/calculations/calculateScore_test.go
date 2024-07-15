package calculations_test

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/Improwised/quizz-app/api/cli"
	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	"github.com/Improwised/quizz-app/api/database"
	"github.com/Improwised/quizz-app/api/helpers/calculations"
	"github.com/Improwised/quizz-app/api/logger"
	"github.com/Improwised/quizz-app/api/models"
	"github.com/Improwised/quizz-app/api/pkg/structs"
	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestCalculatePointsAndScore(t *testing.T) {
	logger, err := logger.NewRootLogger(true, true)
	assert.Nil(t, err)

	err = os.Chdir("../../")
	assert.Nil(t, err, "error while changing directory to project root")

	err = godotenv.Load(".env.testing")
	assert.Nil(t, err, "error while loading testing env")

	cfg := config.GetConfig()

	migrationCmd := cli.GetMigrationCommandDef(cfg)
	migrationCmd.SetArgs([]string{"up"})
	err = migrationCmd.Execute()
	assert.Nil(t, err, "error while execute migration")

	db, err := database.Connect(cfg.DB)
	assert.Nil(t, err)

	t.Run("check for single correct answer", func(t *testing.T) {
		answers := make([]int, 0)
		answers = append(answers, 1)
		question := models.Question{
			ID:                uuid.New(),
			Question:          "A dummy question",
			Points:            int16(20),
			DurationInSeconds: 20,
			Type:              1,
		}
		options := map[string]string{
			"1": "1",
			"2": "2",
			"3": "3",
			"4": "4",
		}

		jsonAnswers, err := json.Marshal(answers)
		assert.Nil(t, err)

		jsonOptions, err := json.Marshal(options)
		assert.Nil(t, err)

		_, err = db.Insert(constants.QuestionsTable).Rows(
			goqu.Record{
				"id":                  question.ID,
				"question":            question.Question,
				"options":             string(jsonOptions),
				"answers":             string(jsonAnswers),
				"points":              question.Points,
				"duration_in_seconds": question.DurationInSeconds,
				"type":                question.Type,
			},
		).Executor().Exec()
		assert.Nil(t, err)

		answerKey := make([]int, 1)
		answerKey[0] = 1
		userAnswer := structs.ReqAnswerSubmit{
			QuestionId:   question.ID,
			AnswerKeys:   answerKey,
			ResponseTime: 2000,
		}
		points, finalScore, err := calculations.CalculatePointsAndScore(userAnswer, db, logger)

		assert.Nil(t, err)
		assert.NotEmpty(t, points)

		expectedScore := (int(question.Points) * 100) + 500 + ((question.DurationInSeconds - int(userAnswer.ResponseTime/1000)) * int(math.Round(400/float64(question.DurationInSeconds))))

		assert.Equal(t, expectedScore, finalScore)
	})

	t.Run("check for the survey questions", func(t *testing.T) {
		answers := make([]int, 0)
		answers = append(answers, 1, 2, 3, 4)
		question := models.Question{
			ID:                uuid.New(),
			Question:          "A dummy question",
			Points:            int16(20),
			DurationInSeconds: 20,
			Type:              2,
		}
		options := map[string]string{
			"1": "1",
			"2": "2",
			"3": "3",
			"4": "4",
		}

		jsonAnswers, err := json.Marshal(answers)
		assert.Nil(t, err)

		jsonOptions, err := json.Marshal(options)
		assert.Nil(t, err)

		_, err = db.Insert(constants.QuestionsTable).Rows(
			goqu.Record{
				"id":                  question.ID,
				"question":            question.Question,
				"options":             string(jsonOptions),
				"answers":             string(jsonAnswers),
				"points":              question.Points,
				"duration_in_seconds": question.DurationInSeconds,
				"type":                question.Type,
			},
		).Executor().Exec()
		assert.Nil(t, err)

		answerKey := make([]int, 1)
		answerKey[0] = 1
		userAnswer := structs.ReqAnswerSubmit{
			QuestionId:   question.ID,
			AnswerKeys:   answerKey,
			ResponseTime: 2000,
		}
		points, finalScore, err := calculations.CalculatePointsAndScore(userAnswer, db, logger)

		assert.Nil(t, err)
		assert.NotEmpty(t, points)

		expectedScore := (int(question.Points) * 100) + 500 + ((question.DurationInSeconds - int(userAnswer.ResponseTime/1000)) * int(math.Round(400/float64(question.DurationInSeconds))))

		assert.Equal(t, expectedScore, finalScore)
	})

	t.Run("check if no answer is submitted", func(t *testing.T) {
		answers := make([]int, 0)
		question := models.Question{
			ID:                uuid.New(),
			Question:          "A dummy question",
			Points:            int16(20),
			DurationInSeconds: 20,
			Type:              1,
		}
		options := map[string]string{
			"1": "1",
			"2": "2",
			"3": "3",
			"4": "4",
		}

		jsonAnswers, err := json.Marshal(answers)
		assert.Nil(t, err)

		jsonOptions, err := json.Marshal(options)
		assert.Nil(t, err)

		_, err = db.Insert(constants.QuestionsTable).Rows(
			goqu.Record{
				"id":                  question.ID,
				"question":            question.Question,
				"options":             string(jsonOptions),
				"answers":             string(jsonAnswers),
				"points":              question.Points,
				"duration_in_seconds": question.DurationInSeconds,
				"type":                question.Type,
			},
		).Executor().Exec()
		assert.Nil(t, err)

		answerKey := make([]int, 1)
		answerKey[0] = 1
		userAnswer := structs.ReqAnswerSubmit{
			QuestionId:   question.ID,
			AnswerKeys:   answerKey,
			ResponseTime: 2000,
		}
		points, finalScore, err := calculations.CalculatePointsAndScore(userAnswer, db, logger)

		assert.Nil(t, err)
		assert.NotEmpty(t, points)
		fmt.Println(points)

		assert.Equal(t, 0, finalScore)
	})

	t.Cleanup(func() {
		db.Delete(constants.QuestionsTable).Executor().Exec()
	})
}
