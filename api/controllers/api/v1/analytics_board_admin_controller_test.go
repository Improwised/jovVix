package v1_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Improwised/quizz-app/api/config"
	"github.com/Improwised/quizz-app/api/constants"
	v1 "github.com/Improwised/quizz-app/api/controllers/api/v1"
	"github.com/Improwised/quizz-app/api/database"
	"github.com/Improwised/quizz-app/api/logger"
	"github.com/Improwised/quizz-app/api/pkg/events"
	"github.com/Improwised/quizz-app/api/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewAnalyticsBoardAdminController(t *testing.T) {
	cfg := config.LoadTestEnv()

	db, err := database.Connect(cfg.DB)
	assert.Nil(t, err)

	logger, err := logger.NewRootLogger(true, true)
	assert.Nil(t, err)

	events := events.NewEventBus(logger)

	err = events.SubscribeAll()
	assert.Nil(t, err)
	t.Run("check whether controller is being returned or not", func(t *testing.T) {

		analyticsAdminController, err := v1.NewAnalyticsBoardAdminController(db, logger, events, &cfg)
		assert.Nil(t, err)

		assert.NotNil(t, analyticsAdminController)
	})

}

func TestGetAnalyticsForAdmin(t *testing.T) {

	t.Run("get analytics score for admin with unauthorized user", func(t *testing.T) {
		res, err := userclient.
			R().
			EnableTrace().
			Get("/api/v1/analytics_board/admin")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode(), res)
	})

	t.Run("get analytics score for admin with invalid input (active_quiz_id query params is not given)", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			Get("/api/v1/analytics_board/admin")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, res.StatusCode(), res)
	})

	t.Run("get analytics score for admin with invalid input (active_quiz_id query params is invalid)", func(t *testing.T) {
		res, err := client.
			R().
			EnableTrace().
			SetQueryParam(constants.ActiveQuizId, invaliduserPlayedQuizId).
			Get("/api/v1/analytics_board/admin")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var responseAnalyticsBoardForAdmin utils.ResponseAnalyticsBoardForAdmin
		err = json.Unmarshal(res.Body(), &responseAnalyticsBoardForAdmin.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.Equal(t, "success", responseAnalyticsBoardForAdmin.Body.Status, res)
		assert.Equal(t, 0, len(responseAnalyticsBoardForAdmin.Body.Data), res)
	})

	t.Run("get analytics score for admin with valid input", func(t *testing.T) {
		createQuiz(t, quizTitle)
		generateDemoSession(t)
		playedQuizValidation(t)

		res, err := client.
			R().
			EnableTrace().
			SetQueryParam(constants.ActiveQuizId, sessionId).
			Get("/api/v1/analytics_board/admin")

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode(), res)

		var responseAnalyticsBoardForAdmin utils.ResponseAnalyticsBoardForAdmin
		err = json.Unmarshal(res.Body(), &responseAnalyticsBoardForAdmin.Body)
		assert.Nil(t, err, "Error in parsing response body")

		assert.Equal(t, 5, len(responseAnalyticsBoardForAdmin.Body.Data), res)
		if len(responseAnalyticsBoardForAdmin.Body.Data) >= 1 {
			assert.Equal(t, "success", responseAnalyticsBoardForAdmin.Body.Status, res)
			assert.Equal(t, 5, len(responseAnalyticsBoardForAdmin.Body.Data), res)
			assert.Equal(t, "testcaseuser", responseAnalyticsBoardForAdmin.Body.Data[0].UserName, res)
			assert.Equal(t, "testcaseuser", responseAnalyticsBoardForAdmin.Body.Data[0].FirstName, res)
			assert.Equal(t, "[2]", responseAnalyticsBoardForAdmin.Body.Data[0].CorrectAnswer, res)
			assert.Equal(t, 0, responseAnalyticsBoardForAdmin.Body.Data[0].CalculatedScore, res)
			assert.Equal(t, false, responseAnalyticsBoardForAdmin.Body.Data[0].IsAttend, res)
			assert.Equal(t, -1, responseAnalyticsBoardForAdmin.Body.Data[0].ResponseTime, res)
			assert.Equal(t, 0, responseAnalyticsBoardForAdmin.Body.Data[0].CalculatedPoints, res)
			assert.Equal(t, "Which city is known as the Eternal City?", responseAnalyticsBoardForAdmin.Body.Data[0].Question, res)
			assert.Equal(t, map[string]string{"1": "Paris", "2": "Rome", "3": "Athens", "4": "Cairo"}, responseAnalyticsBoardForAdmin.Body.Data[0].Options, res)
			assert.Equal(t, "text", responseAnalyticsBoardForAdmin.Body.Data[0].QuestionsMedia, res)
			assert.Equal(t, "text", responseAnalyticsBoardForAdmin.Body.Data[0].OptionsMedia, res)
			assert.Equal(t, "", responseAnalyticsBoardForAdmin.Body.Data[0].Resource, res)
			assert.Equal(t, 1, responseAnalyticsBoardForAdmin.Body.Data[0].Points, res)
			assert.Equal(t, 1, responseAnalyticsBoardForAdmin.Body.Data[0].QuestionTypeID, res)
			assert.Equal(t, "single answer", responseAnalyticsBoardForAdmin.Body.Data[0].QuestionType, res)
			assert.Equal(t, 1, responseAnalyticsBoardForAdmin.Body.Data[0].OrderNo, res)

			fmt.Println(responseAnalyticsBoardForAdmin.Body.Data[0].Options)
		}
	})
}
